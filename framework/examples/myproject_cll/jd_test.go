package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgJD struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	JD          *jd.Input         `toml:"jd" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestJDAndNodeSet(t *testing.T) {
	in, err := framework.Load[CfgJD](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err)
	_, err = jd.NewJD(in.JD)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
