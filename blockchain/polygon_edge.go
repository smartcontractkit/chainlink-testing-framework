package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// PolygonEdgeMultinodeClient represents a multi-node, EVM compatible client for the Klaytn network
type PolygonEdgeMultinodeClient struct {
	*EthereumMultinodeClient
}

// PolygonEdgeClient represents a single node, EVM compatible client for the Polygon edge network
type PolygonEdgeClient struct {
	*EthereumClient
}

func (k *PolygonEdgeClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	privateKey, err := crypto.HexToECDSA(k.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return err
	}
	nonce, err := k.GetNonce(context.Background(), k.DefaultWallet.address)
	if err != nil {
		return err
	}
	log.Warn().
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Setting GasTipCap = SuggestedGasPrice for Polygon edge network")
	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(k.GetChainID()), &types.DynamicFeeTx{
		ChainID: k.GetChainID(),
		Nonce:   nonce,
		To:      &to,
		Value:   utils.EtherToWei(amount),
		Gas:     100000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("From", k.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := k.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	return k.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (k *PolygonEdgeClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := k.TransactionOpts(k.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasTipCap = nil
	opts.GasPrice = nil

	contractAddress, transaction, contractInstance, err := deployer(opts, k.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = k.ProcessTransaction(transaction); err != nil {
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
