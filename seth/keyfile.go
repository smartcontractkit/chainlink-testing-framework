package seth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/sync/errgroup"
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
		return "", "", fmt.Errorf("failed to cast generated public key to ECDSA type.\n"+
			"This is an internal error in the crypto.GenerateKey() function.\n"+
			"Expected type: *ecdsa.PublicKey, got: %T\n"+
			"Please report this issue: https://github.com/smartcontractkit/chainlink-testing-framework/issues",
			publicKey)
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
		if err := c.validateAddressesKeyNum(0); err != nil {
			return err
		}
		toAddr = c.Addresses[0].Hex()
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	gasPrice, err := c.GetSuggestedLegacyFees(ctx, Priority_Standard)
	if err != nil {
		gasPrice = big.NewInt(c.Cfg.Network.GasPrice)
	}

	if len(c.Addresses) == 1 {
		return fmt.Errorf("no ephemeral addresses found to return funds from.\n"+
			"Current addresses count: %d (only root key present)\n"+
			"This indicates either:\n"+
			"  1. Key file doesn't contain ephemeral addresses\n"+
			"  2. Wrong key file was loaded\n"+
			"  3. Ephemeral keys were never created (set 'ephemeral_addresses_number' > 0 in config)",
			len(c.Addresses))
	}

	eg, egCtx := errgroup.WithContext(ctx)
	for i := 1; i < len(c.Addresses); i++ {
		idx := i //nolint
		eg.Go(func() error {
			ctx, balanceCancel := context.WithTimeout(egCtx, c.Cfg.Network.TxnTimeout.Duration())
			balance, err := c.Client.BalanceAt(ctx, c.Addresses[idx], nil)
			balanceCancel()
			if err != nil {
				L.Error().Err(err).Msg("Error getting balance")
				return err
			}

			var gasLimit int64
			//nolint
			gasLimitRaw, err := c.EstimateGasLimitForFundTransfer(c.Addresses[idx], common.HexToAddress(toAddr), balance)
			if err != nil {
				gasLimit = c.Cfg.Network.TransferGasFee
			} else {
				gasLimit = mustSafeInt64(gasLimitRaw)
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

	return eg.Wait()
}
