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
	fakeURL      string
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
	require.Contains(t, output.CLNodes[0].PostgreSQL.DockerInternalURL, "postgresql://chainlink:thispasswordislongenough@postgresql-")
	require.Contains(t, output.CLNodes[0].Node.HostURL, "127.0.0.1")
	require.Contains(t, output.CLNodes[0].Node.DockerURL, "node")
	require.Contains(t, output.CLNodes[0].Node.DockerP2PUrl, "node")

	require.Contains(t, output.CLNodes[1].PostgreSQL.Url, "postgresql://chainlink:thispasswordislongenough@127.0.0.1")
	require.Contains(t, output.CLNodes[1].PostgreSQL.DockerInternalURL, "postgresql://chainlink:thispasswordislongenough@postgresql-")
	require.Contains(t, output.CLNodes[1].Node.HostURL, "127.0.0.1")
	require.Contains(t, output.CLNodes[1].Node.DockerURL, "node")
	require.Contains(t, output.CLNodes[1].Node.DockerP2PUrl, "node")
}

func TestDockerNodeSetSharedDB(t *testing.T) {
	testCases := []testCase{
		{
			name:    "2 nodes cluster, override mode 'all'",
			fakeURL: "http://example.com",
			bcInput: &blockchain.Input{
				Type:      "anvil",
				Image:     "f4hrenh9it/foundry",
				PullImage: true,
				Port:      "8545",
				ChainID:   "31337",
			},
			nodeSetInput: &ns.Input{
				Nodes:        2,
				OverrideMode: "all",
				NodeSpecs: []*clnode.Input{
					{
						DataProviderURL: "http://example.com",
						DbInput: &postgres.Input{
							Image:     "postgres:15.6",
							PullImage: true,
						},
						Node: &clnode.NodeInput{
							Image:     "public.ecr.aws/chainlink/chainlink:v2.17.0",
							Name:      "cl-node",
							PullImage: true,
						},
					},
				},
			},
			assertion: func(t *testing.T, output *ns.Output) {
				checkBasicOutputs(t, output)
			},
		},
		{
			name:    "2 nodes cluster, override mode 'each'",
			fakeURL: "http://example.com",
			bcInput: &blockchain.Input{
				Type:      "anvil",
				Image:     "f4hrenh9it/foundry",
				PullImage: true,
				Port:      "8546",
				ChainID:   "31337",
			},
			nodeSetInput: &ns.Input{
				Nodes:        2,
				OverrideMode: "each",
				NodeSpecs: []*clnode.Input{
					{
						DataProviderURL: "http://example.com",
						DbInput: &postgres.Input{
							Image:     "postgres:15.6",
							PullImage: true,
						},
						Node: &clnode.NodeInput{
							Image:     "public.ecr.aws/chainlink/chainlink:v2.17.0",
							Name:      "cl-node-1",
							PullImage: true,
							UserConfigOverrides: `
[Log]
level = 'info'
`,
						},
					},
					{
						DataProviderURL: "http://example.com",
						DbInput: &postgres.Input{
							Image:     "postgres:15.6",
							PullImage: true,
						},
						Node: &clnode.NodeInput{
							Image:     "public.ecr.aws/chainlink/chainlink:v2.17.0",
							Name:      "cl-node-2",
							PullImage: true,
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
		err := framework.DefaultNetwork(t, &sync.Once{})
		require.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {
			bc, err := blockchain.NewBlockchainNetwork(tc.bcInput)
			require.NoError(t, err)
			output, err := ns.NewSharedDBNodeSet(tc.nodeSetInput, bc, tc.fakeURL)
			require.NoError(t, err)
			tc.assertion(t, output)
		})
	}
}
