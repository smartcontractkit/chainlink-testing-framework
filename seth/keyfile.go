package seth

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"math/big"
)

// NewAddress creates a new address
func NewAddress() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	L.Info().
		Str("Addr", address).
		Msg("New address created")
	return address, hexutil.Encode(privateKeyBytes)[2:], nil
}

// ReturnFunds returns funds to the root key from all other keys
func ReturnFunds(c *Client, toAddr string) error {
	if toAddr == "" {
		toAddr = c.Addresses[0].Hex()
	}

	gasPrice, err := c.GetSuggestedLegacyFees(context.Background(), Priority_Standard)
	if err != nil {
		gasPrice = big.NewInt(c.Cfg.Network.GasPrice)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, egCtx := errgroup.WithContext(ctx)

	if len(c.Addresses) == 1 {
		return errors.New("No addresses to return funds from. Have you passed correct key file?")
	}

	for i := 1; i < len(c.Addresses); i++ {
		idx := i
		eg.Go(func() error {
			balance, err := c.Client.BalanceAt(context.Background(), c.Addresses[idx], nil)
			if err != nil {
				L.Error().Err(err).Msg("Error getting balance")
				return err
			}

			var gasLimit int64
			gasLimitRaw, err := c.EstimateGasLimitForFundTransfer(c.Addresses[idx], common.HexToAddress(toAddr), balance)
			if err != nil {
				gasLimit = c.Cfg.Network.TransferGasFee
			} else {
				gasLimit = int64(gasLimitRaw)
			}

			networkTransferFee := gasPrice.Int64() * gasLimit
			fundsToReturn := new(big.Int).Sub(balance, big.NewInt(networkTransferFee))

			if fundsToReturn.Cmp(big.NewInt(0)) == -1 {
				L.Warn().
					Str("Key", c.Addresses[idx].Hex()).
					Interface("Balance", balance).
					Interface("NetworkFee", networkTransferFee).
					Interface("FundsToReturn", fundsToReturn).
					Msg("Insufficient funds to return. Skipping.")
				return nil
			}

			L.Info().
				Str("Key", c.Addresses[idx].Hex()).
				Interface("Balance", balance).
				Interface("NetworkFee", c.Cfg.Network.GasPrice*gasLimit).
				Interface("GasLimit", gasLimit).
				Interface("GasPrice", gasPrice).
				Interface("FundsToReturn", fundsToReturn).
				Msg("Returning funds from address")

			return c.TransferETHFromKey(
				egCtx,
				idx,
				toAddr,
				fundsToReturn,
				gasPrice,
			)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
