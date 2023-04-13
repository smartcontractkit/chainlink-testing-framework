package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Handles specific issues with the Quorum EVM chain: https://docs.Quorum.com/

// QuorumMultinodeClient represents a multi-node, EVM compatible client for the Quorum network
type QuorumMultinodeClient struct {
	*EthereumMultinodeClient
}

// QuorumClient represents a single node, EVM compatible client for the Quorum network
type QuorumClient struct {
	*EthereumClient
}

// Fund overrides ethereum's fund to account for Quorum's gas specifications
func (k *QuorumClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	privateKey, err := crypto.HexToECDSA(k.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return err
	}
	// Don't bump gas for Quorum
	gasPrice, err := k.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	nonce, err := k.GetNonce(context.Background(), k.DefaultWallet.address)
	if err != nil {
		return err
	}
	log.Warn().
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Setting GasTipCap = SuggestedGasPrice for Quorum network")
	estimatedGas, err := k.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}
	// https://docs.Quorum.com/Quorum/design/transaction-fees#gas
	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(k.GetChainID()), &types.DynamicFeeTx{
		ChainID:   k.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     utils.EtherToWei(amount),
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "Q").
		Str("From", k.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(estimatedGas)).Uint64()).
		Msg("Funding Address")
	if err := k.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	return k.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (k *QuorumClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := k.TransactionOpts(k.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}

	// Don't bump gas for Quorum
	// https://docs.Quorum.com/Quorum/design/transaction-fees#unit-price
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
		Str("Total Gas Cost (Q)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

func (k *QuorumClient) ReturnFunds(fromPrivateKey *ecdsa.PrivateKey) error {
	to := common.HexToAddress(k.DefaultWallet.Address())

	// Don't bump gas for Quorum
	gasPrice, err := k.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	fromAddress, err := utils.PrivateKeyToAddress(fromPrivateKey)
	if err != nil {
		return err
	}

	nonce, err := k.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	balance, err := k.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	estimatedGas, err := k.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}
	balance.Sub(balance, big.NewInt(1).Mul(gasPrice, big.NewInt(0).SetUint64(estimatedGas)))
	// https://docs.Quorum.com/Quorum/design/transaction-fees#gas
	tx, err := types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(k.GetChainID()), &types.DynamicFeeTx{
		ChainID:   k.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     balance,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "Q").
		Str("From", fromAddress.Hex()).
		Str("Amount", balance.String()).
		Msg("Returning Funds to Default Wallet")
	if err := k.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	return k.ProcessTransaction(tx)
}
