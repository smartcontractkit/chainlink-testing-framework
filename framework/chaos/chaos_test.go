package chaos_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
)

func TestChaos(t *testing.T) {
	err := framework.DefaultNetwork(nil)
	require.NoError(t, err)
	_, err = rpc.StartAnvil([]string{"--balance", "1", "--block-time", "5"})
	require.NoError(t, err)

	dtc, err := chaos.NewDockerChaos(t.Context())
	require.NoError(t, err)

	containerName := "anvil"

	tests := []struct {
		name          string
		containerName string
		cmd           string
		value         string
		wantErr       bool
	}{
		{
			name:          "pause container",
			containerName: containerName,
			cmd:           chaos.CmdPause,
		},
		{
			name:          "delay container",
			containerName: containerName,
			cmd:           chaos.CmdDelay,
			value:         "8000ms",
		},
		{
			name:          "loss container",
			containerName: containerName,
			cmd:           chaos.CmdLoss,
			value:         "50%",
		},
		{
			name:          "corrupt traffic in container",
			containerName: containerName,
			cmd:           chaos.CmdCorrupt,
			value:         "10%",
		},
		{
			name:          "duplicate traffic in container",
			containerName: containerName,
			cmd:           chaos.CmdDuplicate,
			value:         "50%",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err = dtc.Chaos(tc.containerName, tc.cmd, tc.value)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			time.Sleep(5 * time.Second)

			err = dtc.RemoveAll()
			require.NoError(t, err)
		})
	}
}
