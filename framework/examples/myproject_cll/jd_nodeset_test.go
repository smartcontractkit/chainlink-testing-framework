package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgJDNodeSet struct {
	Blockchains        []*blockchain.Input `toml:"blockchains" validate:"required"`
	NodeSets           []*ns.Input         `toml:"nodesets" validate:"required"`
	MockedDataProvider *fake.Input         `toml:"data_provider" validate:"required"`
	JD                 *jd.Input           `toml:"jd" validate:"required"`
}

func TestJDNodeSet(t *testing.T) {
	in, err := framework.Load[CfgJDNodeSet](t)
	require.NoError(t, err)

	blockchainInfo, err := blockchain.NewBlockchainNetwork(in.Blockchains[0])
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockedDataProvider)
	require.NoError(t, err)
	nodeSetInfo, err := ns.NewSharedDBNodeSet(in.NodeSets[0], blockchainInfo)
	require.NoError(t, err)

	jobDistributorInfo, err := jd.NewJD(in.JD)
	require.NoError(t, err)

	t.Run("test changesets with forked network/JD state", func(t *testing.T) {
		clClients, err := clclient.New(nodeSetInfo.CLNodes)
		require.NoError(t, err)
		_ = clClients
		// create some jobs
		//_, _, err = c[0].CreateJobRaw()
		//require.NoError(t, err)
		_ = nodeSetInfo
		_ = blockchainInfo
		_ = jobDistributorInfo
	})
}
