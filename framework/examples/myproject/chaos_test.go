package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgChaos struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestChaos(t *testing.T) {
	in, err := framework.Load[CfgChaos](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	c, err := clclient.NewCLDefaultClients(out.CLNodes, framework.L)
	require.NoError(t, err)

	t.Run("run the cluster and simulate slow network", func(t *testing.T) {
		// example commands for Pumba:
		// stop --duration=1s --restart re2:node0                                            # stop one container for 1s and restart
		// netem --tc-image=gaiadocker/iproute2 --duration=1m delay --time=300 re2:node.*   # slow network
		_, err = chaos.ExecPumba("stop --duration=1s --restart re2:node0")
		require.NoError(t, err)
		_, _, err = c[0].ReadBridges()
		require.NoError(t, err)
	})
}
