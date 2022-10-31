package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// ArbitrumMultinodeClient represents a multi-node, EVM compatible client for the Arbitrum network
type ArbitrumMultinodeClient struct {
	*EthereumMultinodeClient
}

// ArbitrumClient represents a single node, EVM compatible client for the Arbitrum network
type ArbitrumClient struct {
	*EthereumClient
}

// Fund sends some ARB to an address using the default wallet
func (a *ArbitrumClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(a.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// Arbitrum uses legacy transactions and gas estimations
	suggestedGasPrice, err := a.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	nonce, err := a.GetNonce(context.Background(), common.HexToAddress(a.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	gas, err := a.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(a.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    utils.EtherToWei(amount),
		GasPrice: suggestedGasPrice,
		Gas:      gas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ARB").
		Str("From", a.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := a.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return a.ProcessTransaction(tx)
}

// Fund sends some ARB to an address using the default wallet
func (a *ArbitrumClient) ReturnFunds(fromPrivateKey *ecdsa.PrivateKey) error {
	to := common.HexToAddress(a.DefaultWallet.Address())

	// Arbitrum uses legacy transactions and gas estimations
	suggestedGasPrice, err := a.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fromAddress, err := utils.PrivateKeyToAddress(fromPrivateKey)
	if err != nil {
		return err
	}

	balance, err := a.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	gas, err := a.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}
	balance.Sub(balance, big.NewInt(1).Mul(suggestedGasPrice, big.NewInt(0).SetUint64(gas)))

	nonce, err := a.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(a.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    balance,
		GasPrice: suggestedGasPrice,
		Gas:      gas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ARB").
		Str("From", fromAddress.Hex()).
		Str("Amount", balance.String()).
		Msg("Returning Funds to Default Wallet")
	if err := a.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return a.ProcessTransaction(tx)
}
