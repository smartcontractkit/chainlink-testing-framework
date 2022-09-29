package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

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
func (m *ArbitrumClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(m.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// Arbitrum uses legacy transactions and gas estimations
	suggestedGasPrice, err := m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	nonce, err := m.GetNonce(context.Background(), common.HexToAddress(m.DefaultWallet.Address()))
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(m.GetChainID()), &types.LegacyTx{
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
		Str("Token", "ARB").
		Str("From", m.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := m.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return m.ProcessTransaction(tx)
}

// Fund sends some ARB to an address using the default wallet
func (m *ArbitrumClient) ReturnFunds(fromPrivateKey *ecdsa.PrivateKey) error {
	to := common.HexToAddress(m.DefaultWallet.Address())

	// Arbitrum uses legacy transactions and gas estimations
	suggestedGasPrice, err := m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fromAddress, err := utils.PrivateKeyToAddress(fromPrivateKey)
	if err != nil {
		return err
	}

	balance, err := m.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	balance.Sub(balance, big.NewInt(1).Mul(suggestedGasPrice, big.NewInt(21000)))

	nonce, err := m.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(m.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    balance,
		GasPrice: suggestedGasPrice,
		Gas:      21000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ARB").
		Str("From", fromAddress.Hex()).
		Str("Amount", balance.String()).
		Msg("Returning Funds to Default Wallet")
	if err := m.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return m.ProcessTransaction(tx)
}
