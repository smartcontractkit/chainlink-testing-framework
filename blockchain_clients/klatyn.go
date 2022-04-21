package blockchainclient

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type KlatynMultinodeClient struct {
	DefaultClient *EthereumClient
	Clients       []*EthereumClient
}

type KlatynClient struct {
	ethClient *EthereumClient
}

// SendTransaction override for Klatyn's gas specifications
// https://docs.klaytn.com/klaytn/design/transaction-fees#unit-price
func (k *KlatynClient) SendTransaction(
	from *EthereumWallet,
	to common.Address,
	value *big.Float,
) (common.Hash, error) {
	weiValue, _ := value.Int(nil)
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid private key: %v", err)
	}
	suggestedGasPrice, err := k.ethClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	log.Warn().
		Str("Network ID", k.ethClient.NetworkConfig.ID).
		Msg("Not bumping gas price while running on a Klaytn network.")

	nonce, err := k.ethClient.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return common.Hash{}, err
	}

	// TODO: Update from LegacyTx to DynamicFeeTx
	tx, err := types.SignNewTx(privateKey, types.NewEIP2930Signer(big.NewInt(k.ethClient.NetworkConfig.ChainID)),
		&types.LegacyTx{
			To:       &to,
			Value:    weiValue,
			Data:     nil,
			Gas:      21000,
			GasPrice: suggestedGasPrice,
			Nonce:    nonce,
		})
	if err != nil {
		return common.Hash{}, err
	}
	if k.ethClient.NetworkConfig.GasEstimationBuffer > 0 {
		log.Debug().
			Uint64("Suggested Gas Price Wei", suggestedGasPrice.Uint64()).
			Str("TX Hash", tx.Hash().Hex()).
			Msg("Bumping Suggested Gas Price")
	}
	if err := k.ethClient.Client.SendTransaction(context.Background(), tx); err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), k.ethClient.ProcessTransaction(tx)
}
