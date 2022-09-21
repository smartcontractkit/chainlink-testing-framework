package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Optimism network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Optimism network
type OptimismClient struct {
	*EthereumClient
}

// Fund sends some OP to an address using the default wallet
func (o *OptimismClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(o.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// Optimism uses legacy transactions and gas estimations
	suggestedGasPrice, err := o.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	nonce, err := o.GetNonce(context.Background(), common.HexToAddress(o.DefaultWallet.Address()))
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(o.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    utils.EtherToWei(amount),
		GasPrice: suggestedGasPrice,
		Gas:      21000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "OP").
		Str("From", o.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := o.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return o.ProcessTransaction(tx)
}

// DeployContract deploys smart contracts specifically on Optimism
func (o *OptimismClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := o.TransactionOpts(o.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}

	// Optimism uses legacy transactions and gas estimations, is behind London fork as of 04/27/2022
	suggestedGasPrice, err := o.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(o.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	opts.GasPrice = suggestedGasPrice

	contractAddress, transaction, contractInstance, err := deployer(opts, o.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = o.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", o.DefaultWallet.Address()).
		Str("Total Gas Cost (OP)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", o.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}
