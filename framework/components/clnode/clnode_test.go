package clnode_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

type testCase struct {
	name      string
	input     *clnode.Input
	assertion func(t *testing.T, output *clnode.Output)
}

func checkBasicOutputs(t *testing.T, output *clnode.Output) {
	require.NotNil(t, output)
	require.NotNil(t, output.Node)
	require.Contains(t, output.Node.HostURL, "127.0.0.1")
	require.Contains(t, output.Node.DockerURL, "cl-node")
	require.Contains(t, output.Node.DockerP2PUrl, "cl-node")
	require.NotNil(t, output.PostgreSQL)
	require.Contains(t, output.PostgreSQL.Url, "postgresql://chainlink:thispasswordislongenough@127.0.0.1")
	require.Contains(t, output.PostgreSQL.DockerInternalURL, "postgresql://chainlink:thispasswordislongenough@ns-postgresql")
}

func TestComponentDockerNodeWithSharedDB(t *testing.T) {
	testCases := []testCase{
		{
			name: "basic use case",
			input: &clnode.Input{
				DbInput: &postgres.Input{
					Image:      "postgres:15.6",
					Port:       16000,
					VolumeName: "a",
				},
				Node: &clnode.NodeInput{
					Image: "public.ecr.aws/chainlink/chainlink:v2.17.0",
					Name:  "cl-node-1",
				},
			},
			assertion: func(t *testing.T, output *clnode.Output) {
				checkBasicOutputs(t, output)
			},
		},
	}

	for _, tc := range testCases {
		err := framework.DefaultNetwork(&sync.Once{})
		require.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {
			pgOut, err := postgres.NewPostgreSQL(tc.input.DbInput)
			require.NoError(t, err)
			output, err := clnode.NewNode(tc.input, pgOut)
			require.NoError(t, err)
			tc.assertion(t, output)
		})
	}
}

func TestComponentDockerNodeWithDB(t *testing.T) {
	testCases := []testCase{
		{
			name: "basic use case",
			input: &clnode.Input{
				DbInput: &postgres.Input{
					Image:      "postgres:15.6",
					Port:       15000,
					VolumeName: "b",
				},
				Node: &clnode.NodeInput{
					Image: "public.ecr.aws/chainlink/chainlink:v2.17.0",
					Name:  "cl-node-2",
				},
			},
			assertion: func(t *testing.T, output *clnode.Output) {
				checkBasicOutputs(t, output)
			},
		},
	}

	for _, tc := range testCases {
		err := framework.DefaultNetwork(&sync.Once{})
		require.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {
			output, err := clnode.NewNodeWithDB(tc.input)
			require.NoError(t, err)
			tc.assertion(t, output)
		})
	}
}
