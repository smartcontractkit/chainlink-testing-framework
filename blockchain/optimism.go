package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Handles specific issues with the Metis EVM chain: https://docs.metis.io/

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Metis network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Metis network
type OptimismClient struct {
	*EthereumClient
}

// NewOptimismClient returns an instantiated instance of the Metis client that has connected to the server
func NewOptimismClient(networkSettings *config.ETHNetwork) (EVMClient, error) {
	client, err := NewEthereumClient(networkSettings)
	log.Info().Str("Network Name", client.GetNetworkName()).Msg("Using custom Optimism client")
	return &OptimismClient{client.(*EthereumClient)}, err
}

// NewOptimismMultiNodeClient returns an instantiated instance of all Metis clients connected to all nodes
func NewOptimismMultiNodeClient(
	_ string,
	networkConfig map[string]interface{},
	urls []*url.URL,
) (EVMClient, error) {
	networkSettings := &config.ETHNetwork{}
	err := UnmarshalNetworkConfig(networkConfig, networkSettings)
	if err != nil {
		return nil, err
	}
	log.Info().
		Interface("URLs", networkSettings.URLs).
		Msg("Connecting multi-node client")

	multiNodeClient := &EthereumMultinodeClient{}
	for _, envURL := range urls {
		networkSettings.URLs = append(networkSettings.URLs, envURL.String())
	}
	for idx, networkURL := range networkSettings.URLs {
		networkSettings.URL = networkURL
		ec, err := NewOptimismClient(networkSettings)
		if err != nil {
			return nil, err
		}
		ec.SetID(idx)
		multiNodeClient.Clients = append(multiNodeClient.Clients, ec)
	}
	multiNodeClient.DefaultClient = multiNodeClient.Clients[0]
	return &OptimismMultinodeClient{multiNodeClient}, nil
}

// Fund sends some ETH to an address using the default wallet
func (m *OptimismClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(m.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	suggestedGasPrice, err := m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(m.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	nonce, err := m.GetNonce(context.Background(), common.HexToAddress(m.DefaultWallet.Address()))
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(m.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    utils.EtherToWei(amount),
		GasPrice: suggestedGasPrice,
		Gas:      22000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ETH").
		Str("From", m.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := m.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return m.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an EVM chain
func (m *OptimismClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := m.TransactionOpts(m.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice, err = m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	if m.NetworkConfig.GasEstimationBuffer > 0 {
		log.Debug().
			Str("Contract Name", contractName).
			Msg("Bumping Suggested Gas Price")
	}

	contractAddress, transaction, contractInstance, err := deployer(opts, m.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := m.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", m.DefaultWallet.Address()).
		Str("Total Gas Cost (ETH)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", m.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}
