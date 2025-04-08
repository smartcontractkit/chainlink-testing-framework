package chaos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgChaos struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestChaos(t *testing.T) {
	in, err := framework.Load[CfgChaos](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	c, err := clclient.New(out.CLNodes)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		command  string
		wait     time.Duration
		validate func(c []*clclient.ChainlinkClient) error
	}{
		{
			name:    "Reboot the pods",
			wait:    1 * time.Minute,
			command: "stop --duration=20s --restart re2:don-node0",
			validate: func(c []*clclient.ChainlinkClient) error {
				_, _, err := c[0].ReadBridges()
				return err
			},
		},
		{
			name:    "Introduce network delay",
			wait:    1 * time.Minute,
			command: "netem --tc-image=gaiadocker/iproute2 --duration=1m delay --time=1000 re2:don-node.*",
			validate: func(c []*clclient.ChainlinkClient) error {
				_, _, err := c[0].ReadBridges()
				return err
			},
		},
	}

	// Start WASP load test here, apply average load profile that you expect in production!
	// Configure timeouts and validate all the test cases until the test ends

	// Run chaos test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			_, err = chaos.ExecPumba(tc.command, tc.wait)
			require.NoError(t, err)
			err = tc.validate(c)
			require.NoError(t, err)
		})
	}
}
