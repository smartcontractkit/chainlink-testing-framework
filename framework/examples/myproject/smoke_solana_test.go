package examples

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgSolana struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestSolanaSmoke(t *testing.T) {
	in, err := framework.Load[CfgSolana](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect chainlink nodes
		_ = bc.Nodes[0].DockerInternalHTTPUrl
		// use host URL to deploy contracts
		c := client.NewClient(bc.Nodes[0].HostHTTPUrl)
		latestSlot, err := c.GetSlotWithConfig(context.Background(), client.GetSlotConfig{Commitment: "processed"})
		require.NoError(t, err)
		fmt.Printf("Latest slot: %v\n", latestSlot)
	})
}
