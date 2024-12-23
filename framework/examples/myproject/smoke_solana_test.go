package examples

import (
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

func TestSmokeSolana(t *testing.T) {
	in, err := framework.Load[CfgSolana](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		_ = bc
		// ...
	})
}
