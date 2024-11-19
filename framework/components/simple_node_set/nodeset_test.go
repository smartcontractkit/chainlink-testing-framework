package simple_node_set_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

type testCase struct {
	name         string
	funding      float64
	bcInput      *blockchain.Input
	nodeSetInput *ns.Input
	assertion    func(t *testing.T, output *ns.Output)
}

func checkBasicOutputs(t *testing.T, output *ns.Output) {
	require.NotNil(t, output)
	require.NotNil(t, output.CLNodes)
	require.Len(t, output.CLNodes, 2)
	require.Contains(t, output.CLNodes[0].PostgreSQL.Url, "postgresql://chainlink:thispasswordislongenough@127.0.0.1")
	require.Contains(t, output.CLNodes[0].PostgreSQL.DockerInternalURL, "postgresql://chainlink:thispasswordislongenough@ns-postgresql-")
	require.Contains(t, output.CLNodes[0].Node.HostURL, "127.0.0.1")
	require.Contains(t, output.CLNodes[0].Node.DockerURL, "node")
	require.Contains(t, output.CLNodes[0].Node.DockerP2PUrl, "node")

	require.Contains(t, output.CLNodes[1].PostgreSQL.Url, "postgresql://chainlink:thispasswordislongenough@127.0.0.1")
	require.Contains(t, output.CLNodes[1].PostgreSQL.DockerInternalURL, "postgresql://chainlink:thispasswordislongenough@ns-postgresql-")
	require.Contains(t, output.CLNodes[1].Node.HostURL, "127.0.0.1")
	require.Contains(t, output.CLNodes[1].Node.DockerURL, "node")
	require.Contains(t, output.CLNodes[1].Node.DockerP2PUrl, "node")
}

func TestComponentDockerNodeSetSharedDB(t *testing.T) {
	testCases := []testCase{
		{
			name: "2 nodes cluster, override mode 'all'",
			bcInput: &blockchain.Input{
				Type:    "anvil",
				Image:   "f4hrenh9it/foundry",
				Port:    "8545",
				ChainID: "31337",
			},
			nodeSetInput: &ns.Input{
				Nodes:        2,
				OverrideMode: "all",
				DbInput: &postgres.Input{
					Image: "postgres:15.6",
				},
				NodeSpecs: []*clnode.Input{
					{
						Node: &clnode.NodeInput{
							Image: "public.ecr.aws/chainlink/chainlink:v2.17.0",
							Name:  "cl-node",
						},
					},
				},
			},
			assertion: func(t *testing.T, output *ns.Output) {
				checkBasicOutputs(t, output)
			},
		},
		{
			name: "2 nodes cluster, override mode 'each'",
			bcInput: &blockchain.Input{
				Type:    "anvil",
				Image:   "f4hrenh9it/foundry",
				Port:    "8546",
				ChainID: "31337",
			},
			nodeSetInput: &ns.Input{
				Nodes:              2,
				OverrideMode:       "each",
				HTTPPortRangeStart: 20000,
				P2PPortRangeStart:  22000,
				DbInput: &postgres.Input{
					Image: "postgres:15.6",
					Port:  14000,
				},
				NodeSpecs: []*clnode.Input{
					{
						Node: &clnode.NodeInput{
							Image: "public.ecr.aws/chainlink/chainlink:v2.17.0",
							Name:  "cl-node-1",
							UserConfigOverrides: `
[Log]
level = 'info'
`,
						},
					},
					{
						Node: &clnode.NodeInput{
							Image: "public.ecr.aws/chainlink/chainlink:v2.17.0",
							Name:  "cl-node-2",
							UserConfigOverrides: `
[Log]
level = 'info'
`,
						},
					},
				},
			},
			assertion: func(t *testing.T, output *ns.Output) {
				checkBasicOutputs(t, output)
			},
		},
	}

	for _, tc := range testCases {
		err := framework.DefaultNetwork(&sync.Once{})
		require.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {
			bc, err := blockchain.NewBlockchainNetwork(tc.bcInput)
			require.NoError(t, err)
			output, err := ns.NewSharedDBNodeSet(tc.nodeSetInput, bc)
			require.NoError(t, err)
			tc.assertion(t, output)
		})
	}
}
