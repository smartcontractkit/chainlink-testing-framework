package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type CfgFork struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestFork(t *testing.T) {
	in, err := framework.Load[CfgFork](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	_, err = ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err)

	t.Run("test some contracts with fork", func(t *testing.T) {
		// deploy your contracts with 0s blocks
		_ = rpc.New(bc.Nodes[0].HostHTTPUrl, nil)

		// start a miner so CL nodes can update blocks
		miner := rpc.NewRemoteAnvilMiner(bc.Nodes[0].HostHTTPUrl, nil)
		miner.MinePeriodically(1 * time.Second)
	})
}
