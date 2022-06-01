package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

// Handles specific issues with the Metis EVM chain: https://docs.metis.io/

// MetisMultinodeClient represents a multi-node, EVM compatible client for the Metis network
type MetisMultinodeClient struct {
	*EthereumMultinodeClient
}

// MetisClient represents a single node, EVM compatible client for the Metis network
type MetisClient struct {
	*EthereumClient
}

// NewMetisClient returns an instantiated instance of the Metis client that has connected to the server
func NewMetisClient(networkSettings *config.EVMNetwork, urls *config.EVMUrls) (EVMClient, error) {
	client, err := NewEthereumClient(networkSettings, urls)
	log.Info().Str("Network Name", client.GetNetworkName()).Msg("Using custom Metis client")
	return &MetisClient{client.(*EthereumClient)}, err
}

// NewMetisMultinodeClient returns an instantiated instance of all Metis clients connected to all nodes
func NewMetisMultiNodeClient(
	networkConfig *config.EVMNetwork,
	testEnvironment *environment.Environment,
) (EVMClient, error) {
	log.Info().
		Interface("URLs", networkConfig.URLs).
		Msg("Connecting multi-node client")

	multiNodeClient := &EthereumMultinodeClient{}
	if len(networkConfig.URLs) == 0 {
		simURLs, err := simulatedEthereumURLs(testEnvironment)
		if err != nil {
			return nil, err
		}
		networkConfig.URLs = simURLs
	}
	for idx, networkURLs := range networkConfig.URLs {
		ec, err := NewMetisClient(networkConfig, networkURLs)
		if err != nil {
			return nil, err
		}
		ec.SetID(idx)
		multiNodeClient.Clients = append(multiNodeClient.Clients, ec)
	}
	multiNodeClient.DefaultClient = multiNodeClient.Clients[0]
	return &MetisMultinodeClient{multiNodeClient}, nil
}

// Fund sends some ETH to an address using the default wallet
func (m *MetisClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(m.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// Metis uses legacy transactions and gas estimations, is behind London fork as of 04/27/2022
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
		Str("Token", "METIS").
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
func (m *MetisClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := m.TransactionOpts(m.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}

	// Metis uses legacy transactions and gas estimations, is behind London fork as of 04/27/2022
	suggestedGasPrice, err := m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(m.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	opts.GasPrice = suggestedGasPrice

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
		Str("Total Gas Cost (METIS)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", m.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}
