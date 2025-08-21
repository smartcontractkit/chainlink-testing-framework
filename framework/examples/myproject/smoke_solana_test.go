package examples

import (
	"context"
	"fmt"
	"testing"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type CfgSolana struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestSolanaSmoke(t *testing.T) {
	in, err := framework.Load[CfgSolana](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect Chainlink nodes
		_ = bc.Nodes[0].InternalHTTPUrl
		// use host URL to deploy contracts
		c := client.NewClient(bc.Nodes[0].ExternalHTTPUrl)
		latestSlot, err := c.GetSlotWithConfig(context.Background(), client.GetSlotConfig{Commitment: "processed"})
		require.NoError(t, err)
		fmt.Printf("Latest slot: %v\n", latestSlot)
	})
}
