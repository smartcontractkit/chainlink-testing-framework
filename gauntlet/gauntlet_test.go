package gauntlet_test

import (
	"os"
	"testing"

	gauntlet "github.com/smartcontractkit/chainlink-testing-framework/gauntlet"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestNewGauntletStruct(t *testing.T) {
	t.Parallel()
	g, err := gauntlet.NewGauntlet()
	require.NoError(t, err, "Error getting new Gauntlet struct")
	require.Contains(t, g.Network, "test", "Network did not contain a test")
}

func TestProperlyFormattedFlag(t *testing.T) {
	t.Parallel()
	g, err := gauntlet.NewGauntlet()
	require.NoError(t, err, "Error getting new Gauntlet struct")
	require.Equal(t, "--flag=value", g.Flag("flag", "value"))
}

func TestExecuteCommand(t *testing.T) {
	t.Parallel()
	g, err := gauntlet.NewGauntlet()
	require.NoError(t, err, "Error getting new Gauntlet struct")

	options := gauntlet.ExecCommandOptions{
		ErrHandling:       []string{},
		CheckErrorsInRead: true,
	}
	out, err := g.ExecCommand([]string{}, options)

	require.Error(t, err, "The command should technically always fail because we don't have access to a gauntlet"+
		"executable, if it passed without error then we have an issue")
	require.Contains(t, out, "yarn", "Did not contain expected output")
}

func TestExpectedError(t *testing.T) {
	t.Parallel()
	g, err := gauntlet.NewGauntlet()
	require.NoError(t, err, "Error getting new Gauntlet struct")

	options := gauntlet.ExecCommandOptions{
		ErrHandling: []string{"yarn"},
		RetryCount:  1,
	}
	_, err = g.ExecCommandWithRetries([]string{}, options)
	require.Error(t, err, "Failed to find the expected error")
}
