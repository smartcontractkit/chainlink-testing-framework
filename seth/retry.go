package seth

import (
	"context"
	"fmt"
	"github.com/holiman/uint256"
	"math/big"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

/* these are the common errors of RPCs */

const (
	ErrRPCConnectionRefused = "connection refused"
)

const (
	ErrRetryTimeout = "retry timeout"
)

// RetryTxAndDecode executes transaction several times, retries if connection is lost and decodes all the data
func (m *Client) RetryTxAndDecode(f func() (*types.Transaction, error)) (*DecodedTransaction, error) {
	var tx *types.Transaction
	err := retry.Do(
		func() error {
			var err error
			tx, err = f()
			return err
		}, retry.OnRetry(func(i uint, _ error) {
			L.Debug().Uint("Attempt", i).Msg("Retrying transaction...")
		}),
		retry.DelayType(retry.FixedDelay),
		retry.Attempts(10), retry.Delay(time.Duration(1)*time.Second), retry.RetryIf(func(err error) bool {
			return strings.Contains(err.Error(), ErrRPCConnectionRefused)
		}),
	)

	if err != nil {
		return &DecodedTransaction{}, errors.New(ErrRetryTimeout)
	}

	dt, err := m.Decode(tx, nil)
	if err != nil {
		return &DecodedTransaction{}, errors.Wrap(err, "error decoding transaction")
	}

	return dt, nil
}

// GasBumpStrategyFn is a function that returns a new gas price based on the previous one
type GasBumpStrategyFn = func(previousGasPrice *big.Int) *big.Int

// NoOpGasBumpStrategyFn is a default gas bump strategy that does nothing
var NoOpGasBumpStrategyFn = func(previousGasPrice *big.Int) *big.Int {
	return previousGasPrice
}

// PriorityBasedGasBumpingStrategyFn is a function that returns a gas bump strategy based on the priority.
// For Fast priority it bumps gas price by 30%, for Standard by 15%, for Slow by 5% and for the rest it does nothing.
var PriorityBasedGasBumpingStrategyFn = func(priority string) GasBumpStrategyFn {
	switch priority {
	case Priority_Degen:
		// +100%
		return func(gasPrice *big.Int) *big.Int {
			return gasPrice.Mul(gasPrice, big.NewInt(2))
		}
	case Priority_Fast:
		// +30%
		return func(gasPrice *big.Int) *big.Int {
			gasPriceFloat, _ := gasPrice.Float64()
			newGasPriceFloat := big.NewFloat(0.0).Mul(big.NewFloat(gasPriceFloat), big.NewFloat(1.3))
			newGasPrice, _ := newGasPriceFloat.Int64()
			return big.NewInt(newGasPrice)
		}
	case Priority_Standard:
		// 15%
		return func(gasPrice *big.Int) *big.Int {
			gasPriceFloat, _ := gasPrice.Float64()
			newGasPriceFloat := big.NewFloat(0.0).Mul(big.NewFloat(gasPriceFloat), big.NewFloat(1.15))
			newGasPrice, _ := newGasPriceFloat.Int64()
			return big.NewInt(newGasPrice)
		}
	case Priority_Slow:
		// 5%
		return func(gasPrice *big.Int) *big.Int {
			gasPriceFloat, _ := gasPrice.Float64()
			newGasPriceFloat := big.NewFloat(0.0).Mul(big.NewFloat(gasPriceFloat), big.NewFloat(1.05))
			newGasPrice, _ := newGasPriceFloat.Int64()
			return big.NewInt(newGasPrice)
		}
	default:
		return func(gasPrice *big.Int) *big.Int {
			return gasPrice
		}
	}
}

// prepareReplacementTransaction bumps gas price of the transaction if it wasn't confirmed in time. It returns a signed replacement transaction.
// Errors might be returned, because transaction was no longer pending, max gas price was reached or there was an error sending the transaction (e.g. nonce too low, meaning that original transaction was mined).
var prepareReplacementTransaction = func(client *Client, tx *types.Transaction) (*types.Transaction, error) {
	L.Warn().Msgf("Transaction wasn't confirmed in %s. Bumping gas", client.Cfg.Network.TxnTimeout.String())

	ctxPending, cancelPending := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	_, isPending, err := client.Client.TransactionByHash(ctxPending, tx.Hash())
	defer cancelPending()
	if err != nil {
		return nil, err
	}

	if !isPending {
		L.Debug().Str("Tx hash", tx.Hash().Hex()).Msg("Transaction was confirmed before bumping gas")
		return nil, errors.New("transaction was confirmed before bumping gas")
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}

	senderPkIdx := -1
	for j, maybeSender := range client.Addresses {
		if maybeSender == sender {
			senderPkIdx = j
			break
		}
	}

	if senderPkIdx == -1 {
		return nil, fmt.Errorf("sender address '%s' not found in loaded private keys", sender)
	}

	maxGasPrice := big.NewInt(client.Cfg.GasBump.MaxGasPrice)
	privateKey := client.PrivateKeys[senderPkIdx]
	var replacementTx *types.Transaction

	var checkMaxPrice = func(gasPrice, maxGasPrice *big.Int) error {
		if !client.Cfg.HasMaxBumpGasPrice() {
			L.Debug().Msg("Max gas price for gas bump is not set, skipping check")
			return nil
		}

		if gasPrice.Cmp(maxGasPrice) > 0 {
			return fmt.Errorf("bumped gas price %s is higher than max gas price %s", gasPrice.String(), big.NewInt(client.Cfg.GasBump.MaxGasPrice).String())
		}

		return nil
	}

	switch tx.Type() {
	case types.LegacyTxType:
		gasPrice := client.Cfg.GasBump.StrategyFn(tx.GasPrice())
		if err := checkMaxPrice(gasPrice, maxGasPrice); err != nil {
			return nil, err
		}
		L.Warn().Interface("Old gas price", tx.GasPrice()).Interface("New gas price", gasPrice).Msg("Bumping gas price for legacy transaction")
		txData := &types.LegacyTx{
			Nonce:    tx.Nonce(),
			To:       tx.To(),
			Value:    tx.Value(),
			Gas:      tx.Gas(),
			GasPrice: gasPrice,
			Data:     tx.Data(),
		}
		replacementTx, err = types.SignNewTx(privateKey, signer, txData)
	case types.DynamicFeeTxType:
		gasFeeCap := client.Cfg.GasBump.StrategyFn(tx.GasFeeCap())
		gasTipCap := client.Cfg.GasBump.StrategyFn(tx.GasTipCap())
		if err := checkMaxPrice(big.NewInt(0).Add(gasFeeCap, gasTipCap), maxGasPrice); err != nil {
			return nil, err
		}
		L.Warn().Interface("Old gas fee cap", tx.GasFeeCap()).Interface("New gas fee cap", gasFeeCap).Interface("Old gas tip cap", tx.GasTipCap()).Interface("New gas tip cap", gasTipCap).Msg("Bumping gas fee cap and tip cap for EIP-1559 transaction")
		txData := &types.DynamicFeeTx{
			Nonce:     tx.Nonce(),
			To:        tx.To(),
			Value:     tx.Value(),
			Gas:       tx.Gas(),
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Data:      tx.Data(),
		}

		replacementTx, err = types.SignNewTx(privateKey, signer, txData)
	case types.BlobTxType:
		if tx.To() == nil {
			return nil, fmt.Errorf("blob tx with nil recipient is not supported")
		}
		gasFeeCap := client.Cfg.GasBump.StrategyFn(tx.GasFeeCap())
		gasTipCap := client.Cfg.GasBump.StrategyFn(tx.GasTipCap())
		blobFeeCap := client.Cfg.GasBump.StrategyFn(tx.BlobGasFeeCap())
		if err := checkMaxPrice(big.NewInt(0).Add(gasFeeCap, big.NewInt(0).Add(gasTipCap, blobFeeCap)), maxGasPrice); err != nil {
			return nil, err
		}

		L.Warn().Interface("Old gas fee cap", tx.GasFeeCap()).Interface("Old max fee per blob", tx.BlobGasFeeCap()).Interface("New max fee per blob", blobFeeCap).Interface("New gas fee cap", gasFeeCap).Interface("Old gas tip cap", tx.GasTipCap()).Interface("New gas tip cap", gasTipCap).Msg("Bumping gas fee cap and tip cap for Blob transaction")
		txData := &types.BlobTx{
			Nonce:      tx.Nonce(),
			To:         *tx.To(),
			Value:      uint256.NewInt(tx.Value().Uint64()),
			Gas:        tx.Gas(),
			GasFeeCap:  uint256.NewInt(gasFeeCap.Uint64()),
			GasTipCap:  uint256.NewInt(gasTipCap.Uint64()),
			BlobFeeCap: uint256.NewInt(blobFeeCap.Uint64()),
			BlobHashes: tx.BlobHashes(),
			Data:       tx.Data(),
		}

		replacementTx, err = types.SignNewTx(privateKey, signer, txData)
	case types.AccessListTxType:
		gasPrice := client.Cfg.GasBump.StrategyFn(tx.GasPrice())
		if err := checkMaxPrice(gasPrice, maxGasPrice); err != nil {
			return nil, err
		}
		L.Warn().Interface("Old gas price", tx.GasPrice()).Interface("New gas price", gasPrice).Msg("Bumping gas price for access list transaction")

		txData := &types.AccessListTx{
			Nonce:      tx.Nonce(),
			To:         tx.To(),
			Value:      tx.Value(),
			Gas:        tx.Gas(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		}

		replacementTx, err = types.SignNewTx(privateKey, signer, txData)

	default:
		return nil, fmt.Errorf("unsupported tx type %d", tx.Type())
	}

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	err = client.Client.SendTransaction(ctx, replacementTx)
	if err != nil {
		return nil, err
	}

	return replacementTx, nil
}
