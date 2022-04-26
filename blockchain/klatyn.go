package blockchain

import (
	"context"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/utils"
)

type KlaytnMultinodeClient struct {
	Client
}

type KlaytnClient struct {
	*EthereumClient
}

// NewKlaytnClient returns an instantiated instance of the Klaytn client that has connected to the server
func NewKlaytnClient(networkSettings *config.ETHNetwork) (*KlaytnClient, error) {
	client, err := NewEthereumClient(networkSettings)
	return &KlaytnClient{client}, err
}

// NewKlaytnMultiNodeClient returns an instantiated instance of all Klaytn clients connected to all nodes
func NewKlaytnMultiNodeClient(
	_ string,
	networkConfig map[string]interface{},
	urls []*url.URL,
) (Client, error) {
	client, err := NewEthereumMultiNodeClient("", networkConfig, urls)
	return &KlaytnMultinodeClient{client}, err
}

// SendTransaction override for Klaytn's gas specifications
// https://docs.klaytn.com/klaytn/design/transaction-fees#unit-price
func (k *KlaytnClient) SendTransaction(
	from *EthereumWallet,
	to common.Address,
	value *big.Float,
) (common.Hash, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return common.Hash{}, err
	}
	// Don't bump gas for Klaytn
	gasPrice, err := k.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	nonce, err := k.GetNonce(context.Background(), from.address)
	if err != nil {
		return common.Hash{}, err
	}
	// https://docs.klaytn.com/klaytn/design/transaction-fees#gas
	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(k.GetChainID()), &types.DynamicFeeTx{
		ChainID:   k.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     utils.EtherToWei(value),
		GasTipCap: gasPrice,
		Gas:       22000,
	})
	if err != nil {
		return common.Hash{}, err
	}

	log.Warn().
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Setting GasTipCap = SuggestedGasPrice for Klaytn network")
	if err := k.Client.SendTransaction(context.Background(), tx); err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), k.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (k *KlaytnClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := k.TransactionOpts(k.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}

	// Don't bump gas for Klaytn
	// https://docs.klaytn.com/klaytn/design/transaction-fees#unit-price
	log.Warn().
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Setting GasTipCap = SuggestedGasPrice for Klaytn network")
	opts.GasTipCap, err = k.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	contractAddress, transaction, contractInstance, err := deployer(opts, k.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := k.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", k.DefaultWallet.Address()).
		Str("Total Gas Cost (KLAY)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}
