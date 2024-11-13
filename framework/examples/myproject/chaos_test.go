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
	"time"
)

type CfgChaos struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func verifyServices(t *testing.T, c []*clclient.ChainlinkClient) {
	_, _, err := c[0].ReadBridges()
	require.NoError(t, err)
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

	t.Run("run the cluster and test various chaos scenarios", func(t *testing.T) {
		// Here are examples of using Pumba (https://github.com/alexei-led/pumba)
		// for simplicity we allow users to run commands "as is", read their docs to learn more
		// second parameter is experiment wait time

		// Restart the container
		_, err = chaos.ExecPumba("stop --duration=20s --restart re2:node0", 30*time.Second)
		require.NoError(t, err)
		verifyServices(t, c)

		// Simulate poor network with 1s delay
		_, err = chaos.ExecPumba("netem --tc-image=gaiadocker/iproute2 --duration=1m delay --time=1000 re2:node.*", 30*time.Second)
		require.NoError(t, err)
		verifyServices(t, c)

		// Stress container CPU (TODO: it is not portable, works only in CI or Linux VM, cgroups are required)
		//_, err = chaos.ExecPumba(`stress --stress-image=alexeiled/stress-ng:latest-ubuntu --duration=30s --stressors="--cpu 1 --vm 2 --vm-bytes 1G" node0`, 30*time.Second)
		//require.NoError(t, err)
		//verifyServices(t, c)
	})
}
