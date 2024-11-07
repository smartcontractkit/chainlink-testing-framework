package examples

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type CfgUpgradeAll struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestUpgradeAll(t *testing.T) {
	in, err := framework.Load[CfgReload](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)

	// deploy first time
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	c, err := clclient.NewCLDefaultClients(out.CLNodes, framework.L)
	require.NoError(t, err)
	_, _, err = c[0].CreateJobRaw(testJob)
	require.NoError(t, err)

	in.NodeSet.NodeSpecs[0].Node.Image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
	in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides = `
										[Log]
										level = 'info'
`

	out, err = ns.UpgradeNodeSet(in.NodeSet, bc, dp.BaseURLDocker, 10*time.Second)
	require.NoError(t, err)

	jobs, _, err := c[0].ReadJobs()
	require.NoError(t, err)
	fmt.Println(jobs)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.HostURL)
			require.NotEmpty(t, n.Node.HostP2PURL)
		}
	})
}
