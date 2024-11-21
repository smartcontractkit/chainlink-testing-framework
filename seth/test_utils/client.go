package test_utils

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

var zero int64

// NewClientWithAddresses creates a new Seth client with the given number of addresses. Each address is funded with the given amount of native tokens.
// If funding is `nil` then  each address is funded with the calculated with the amount of ETH calculated by dividing the total balance of root key by the number of addresses (minus root key buffer amount).
func NewClientWithAddresses(t *testing.T, addressCount int, funding *big.Int) *seth.Client {
	cfg, err := seth.ReadConfig()
	require.NoError(t, err, "failed to read config")

	cfg.EphemeralAddrs = &zero

	c, err := seth.NewClientWithConfig(cfg)
	require.NoError(t, err, "failed to initialize seth")

	var privateKeys []string
	var addresses []string
	for i := 0; i < addressCount; i++ {
		addr, pk, err := seth.NewAddress()
		require.NoError(t, err, "failed to generate new address")

		privateKeys = append(privateKeys, pk)
		addresses = append(addresses, addr)
	}

	gasPrice, err := c.GetSuggestedLegacyFees(context.Background(), seth.Priority_Standard)
	if err != nil {
		gasPrice = big.NewInt(c.Cfg.Network.GasPrice)
	}

	if funding == nil {
		bd, err := c.CalculateSubKeyFunding(int64(addressCount), gasPrice.Int64(), *cfg.RootKeyFundsBuffer)
		require.NoError(t, err, "failed to calculate subkey funding")

		funding = bd.AddrFunding
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, egCtx := errgroup.WithContext(ctx)
	// root key is element 0 in ephemeral
	for _, addr := range addresses {
		eg.Go(func() error {
			return c.TransferETHFromKey(egCtx, 0, addr, funding, gasPrice)
		})
	}
	err = eg.Wait()
	require.NoError(t, err, "failed to transfer funds to subkeys")

	// Add root private key to the list of private keys
	pksToUse := []string{cfg.Network.PrivateKeys[0]}
	pksToUse = append(pksToUse, privateKeys...)
	// Set funded private keys in config and create a new Seth client to simulate a situation, in which PKs were passed in config to a new client
	cfg.Network.PrivateKeys = pksToUse

	newClient, err := seth.NewClientWithConfig(cfg)
	require.NoError(t, err, "failed to initialize new Seth with private keys")

	return newClient
}

// NewPrivateKeyWithFunds generates a new private key and funds it with the given amount of native tokens.
func NewPrivateKeyWithFunds(t *testing.T, c *seth.Client, funds *big.Int) string {
	addr, pk, err := seth.NewAddress()
	require.NoError(t, err, "failed to generate new address")

	gasPrice, err := c.GetSuggestedLegacyFees(context.Background(), seth.Priority_Standard)
	if err != nil {
		gasPrice = big.NewInt(c.Cfg.Network.GasPrice)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Cfg.Network.TxnTimeout.Duration())
	err = c.TransferETHFromKey(ctx, 0, addr, funds, gasPrice)
	defer cancel()
	require.NoError(t, err, "failed to transfer funds to subkeys")

	return pk
}

// TransferAllFundsBetweenKeyAndAddress transfers all funds key at specified index has to the given address.
func TransferAllFundsBetweenKeyAndAddress(client *seth.Client, keyNum int, toAddress common.Address) error {
	err := client.NonceManager.UpdateNonces()
	if err != nil {
		return err
	}

	gasPrice, err := client.GetSuggestedLegacyFees(context.Background(), seth.Priority_Standard)
	if err != nil {
		gasPrice = big.NewInt(client.Cfg.Network.GasPrice)
	}

	balance, err := client.Client.BalanceAt(context.Background(), client.Addresses[0], nil)
	if err != nil {
		return err
	}

	toTransfer := new(big.Int).Sub(balance, big.NewInt(0).Mul(gasPrice, big.NewInt(client.Cfg.Network.TransferGasFee)))

	ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	return client.TransferETHFromKey(ctx, keyNum, toAddress.Hex(), toTransfer, gasPrice)
}
