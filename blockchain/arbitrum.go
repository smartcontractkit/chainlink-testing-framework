package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
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
	gasPriceBuffer := big.NewInt(0).SetUint64(a.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

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
		Uint64("Estimated Gas Cost", new(big.Int).Mul(suggestedGasPrice, new(big.Int).SetUint64(gas)).Uint64()).
		Msg("Funding Address")
	if err := a.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return a.ProcessTransaction(tx)
}

func (a *ArbitrumClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	var tx *types.Transaction
	var err error
	for attempt := 1; attempt < 10; attempt++ {
		tx, err = attemptArbReturn(a, fromKey, attempt)
		if err == nil {
			return a.ProcessTransaction(tx)
		}
		log.Debug().Err(err).Int("Attempt", attempt+1).Msg("Error returning funds from Chainlink node, trying again")
	}
	return err
}

// a single fund return attempt, further attempts exponentially raise the error margin for fund returns
func attemptArbReturn(a *ArbitrumClient, fromKey *ecdsa.PrivateKey, attemptCount int) (*types.Transaction, error) {
	to := common.HexToAddress(a.DefaultWallet.Address())

	// Arbitrum uses legacy transactions and gas estimations
	suggestedGasPrice, err := a.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	suggestedGasPrice.Add(suggestedGasPrice, big.NewInt(int64(math.Pow(float64(attemptCount), 2)*1000))) // exponentially increase error margin
	fromAddress, err := utils.PrivateKeyToAddress(fromKey)
	if err != nil {
		return nil, err
	}

	balance, err := a.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, err
	}
	gas, err := a.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}
	balance.Sub(balance, big.NewInt(1).Mul(suggestedGasPrice, big.NewInt(0).SetUint64(gas)))

	nonce, err := a.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	tx, err := types.SignNewTx(fromKey, types.LatestSignerForChainID(a.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    balance,
		GasPrice: suggestedGasPrice,
		Gas:      gas,
	})
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("Token", "ARB").
		Str("From", fromAddress.Hex()).
		Str("Amount", balance.String()).
		Msg("Returning Funds to Default Wallet")
	return tx, a.SendTransaction(context.Background(), tx)
}
