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

// ArbitrumMultinodeClient represents a multi-node, EVM compatible client for the Arbitrum network
type ArbitrumMultinodeClient struct {
	*EthereumMultinodeClient
}

// ArbitrumClient represents a single node, EVM compatible client for the Arbitrum network
type ArbitrumClient struct {
	*EthereumClient
}

// Fund sends some ETH to an address using the default wallet
func (a *ArbitrumClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(a.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// Metis uses legacy transactions and gas estimations, is behind London fork as of 04/27/2022
	suggestedGasPrice, err := a.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(a.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	nonce, err := a.GetNonce(context.Background(), common.HexToAddress(a.DefaultWallet.Address()))
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(a.GetChainID()), &types.LegacyTx{
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

// ProcessTransaction adds tx hash to a separate list to check for on waiting. Arbitrum is a near-instant
// Optimistic Rollup that requires we keep MinConfirmations = 0 in most cases, but can take up to a couple seconds to
// actually confirm. IsPending also seems to be improperly implemented in Arbitrum, so need a separate strategy to
// confirm transactions.
func (a *ArbitrumClient) IsTxConfirmed(txHash common.Hash) (bool, error) {
	tx, _, err := a.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return false, err
	}
	receipt, err := a.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return false, err
	}
	a.gasStats.AddClientTXData(TXGasData{
		TXHash:            txHash.String(),
		Value:             tx.Value().Uint64(),
		GasLimit:          tx.Gas(),
		GasUsed:           receipt.GasUsed,
		GasPrice:          tx.GasPrice().Uint64(),
		CumulativeGasUsed: receipt.CumulativeGasUsed,
	})
	latestBlockNum, err := a.Client.BlockNumber(context.Background())
	if err != nil {
		return false, err
	}

	return receipt.BlockNumber.Uint64() >= latestBlockNum, err
}
