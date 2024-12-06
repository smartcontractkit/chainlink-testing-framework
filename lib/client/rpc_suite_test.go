// nolint
package client

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

func TestRPCSuite(t *testing.T) {
	t.Run("(anvil) test we can modulate next block base fee per gas over a duration", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "1"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		printGasPrices(t, client)
		// set a base fee
		anvilClient := NewRPCClient(ac.URL, nil)
		// set fee for the next block
		err = anvilClient.AnvilSetNextBlockBaseFeePerGas([]interface{}{"2000000000"})
		require.NoError(t, err)
		// mine a block
		err = anvilClient.AnvilMine(nil)
		require.NoError(t, err)
		blockNumber, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		require.NoError(t, err)
		// check the base fee of the block
		require.Equal(t, "2000000000", block.BaseFee().String(), "expected base fee to be 20 gwei")
		logger := logging.GetTestLogger(t)
		err = anvilClient.ModulateBaseFeeOverDuration(logger, 2000000000, 0.5, 20*time.Second, true)
		require.NoError(t, err)
		// mine a block
		err = anvilClient.AnvilMine(nil)
		require.NoError(t, err)
		blockNumber, err = client.BlockNumber(context.Background())
		require.NoError(t, err)
		block, err = client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		require.NoError(t, err)
		// check the base fee of the block
		require.Equal(t, "3000000000", block.BaseFee().String(), "expected base fee to be 30 gwei")
		err = anvilClient.ModulateBaseFeeOverDuration(logger, 3000000000, 0.25, 15*time.Second, false)
		require.NoError(t, err)
		// mine a block
		err = anvilClient.AnvilMine(nil)
		require.NoError(t, err)
		blockNumber, err = client.BlockNumber(context.Background())
		require.NoError(t, err)
		block, err = client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		require.NoError(t, err)
		// check the base fee of the block
		require.Equal(t, "2250000000", block.BaseFee().String(), "expected base fee to be 30 gwei")
	})
}
