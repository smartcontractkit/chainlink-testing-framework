package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgLogs struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestLogsSmoke(t *testing.T) {
	in, err := framework.Load[CfgLogs](t)
	require.NoError(t, err)
	// most simple checks, save all the logs and check (CRIT|PANIC|FATAL) log levels
	//t.Cleanup(func() {
	//	err := framework.SaveAndCheckLogs(t)
	//	require.NoError(t, err)
	//})
	t.Cleanup(func() {
		// save all the logs to default directory "logs/docker-$test_name"
		logs, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
		// check that CL nodes has no errors (CRIT|PANIC|FATAL) levels
		err = framework.CheckCLNodeContainerErrors()
		require.NoError(t, err)
		// do custom assertions
		for _, l := range logs {
			matches, err := framework.SearchLogFile(l, " name=HeadReporter version=\\d")
			require.NoError(t, err)
			_ = matches
		}
	})

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
