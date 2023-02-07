package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func (o *OptimismClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	var tx *types.Transaction
	var err error
	for attempt := 1; attempt < 10; attempt++ {
		tx, err = o.attemptReturn(fromKey, attempt)
		if err == nil {
			return o.ProcessTransaction(tx)
		}
		log.Debug().Err(err).Int("Attempt", attempt+1).Msg("Error returning funds from Chainlink node, trying again")
	}
	return err
}

// a single fund return attempt, further attempts exponentially raise the error margin for fund returns
func (o *OptimismClient) attemptReturn(fromKey *ecdsa.PrivateKey, attemptCount int) (*types.Transaction, error) {
	to := common.HexToAddress(o.DefaultWallet.Address())

	suggestedGasTipCap, err := o.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}
	latestHeader, err := o.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	baseFeeMult := big.NewInt(1).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := baseFeeMult.Add(baseFeeMult, suggestedGasTipCap)
	// optimism is being cagey for now about these values, makes it tricky to estimate properly
	// https://community.optimism.io/docs/developers/bedrock/how-is-bedrock-different/#eip-1559
	gasFeeCap.Add(gasFeeCap, big.NewInt(int64(math.Pow(float64(attemptCount), 3)*1000))) // exponentially increase error margin

	fromAddress, err := utils.PrivateKeyToAddress(fromKey)
	if err != nil {
		return nil, err
	}
	balance, err := o.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, err
	}
	estGas, err := o.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}
	balance.Sub(balance, big.NewInt(1).Mul(gasFeeCap, big.NewInt(0).SetUint64(estGas)))

	nonce, err := o.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	estimatedGas, err := o.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}

	tx, err := types.SignNewTx(fromKey, types.LatestSignerForChainID(o.GetChainID()), &types.DynamicFeeTx{
		ChainID:   o.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     balance,
		GasTipCap: suggestedGasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
	})
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Token", "ETH").
		Str("Amount", balance.String()).
		Str("From", fromAddress.Hex()).
		Msg("Returning Funds to Default Wallet")
	return tx, o.SendTransaction(context.Background(), tx)
}
