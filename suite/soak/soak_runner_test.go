package soak_runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/utils"
	"github.com/stretchr/testify/require"
)

func TestSoak(t *testing.T) {
	actions.LoadConfigs(utils.ProjectRoot)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	exePath, remoteConfig := buildGoTests(t)

	env, err := environment.DeployLongTestEnvironment(
		environment.NewChainlinkConfig(
			environment.ChainlinkReplicas(6, config.ChainlinkVals()),
			"chainlink-soak",
			config.GethNetworks()...,
		),
		tools.ChartsRoot,
		remoteConfig.TestRegex,                             // Name of the test to run
		remoteConfig.SlackAPIKey,                           // API key to use to upload artifacts to slack
		remoteConfig.SlackChannel,                          // Slack Channel to upload test artifacts to
		remoteConfig.SlackUserID,                           // Slack user to notify on completion
		filepath.Join(utils.ProjectRoot, "framework.yaml"), // Path of the framework config
		filepath.Join(utils.ProjectRoot, "networks.yaml"),  // Path to the networks config
		exePath, // Path to the executable test file
	)
	require.NoError(t, err)
	require.NotNil(t, env)
	log.Info().Str("Namespace", env.Namespace).
		Str("Environment File", fmt.Sprintf("%s.%s", env.Namespace, "yaml")).
		Msg("Soak Test Successfully Launched. Save the environment file to collect logs when test is done.")
}

// Builds the go tests to run, and returns a path to it, along with remote config options
func buildGoTests(t *testing.T) (string, *config.RemoteRunnerConfig) {
	exePath := filepath.Join(utils.ProjectRoot, "remote.test")
	remoteConfig, err := config.ReadWriteRemoteRunnerConfig()
	require.NoError(t, err)
	compileCmd := exec.Command("go", "test", "-c", remoteConfig.TestDirectory, "-o", exePath) // #nosec G204
	compileCmd.Env = os.Environ()
	compileCmd.Env = append(compileCmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

	log.Info().Str("Test Directory", remoteConfig.TestDirectory).Msg("Compiling tests")
	compileOut, err := compileCmd.Output()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", compileCmd.String()).
		Msg("Ran command")
	require.NoError(t, err, "Env: %s\nCommand: %s\nCommand Output: %s", compileCmd.Env, compileCmd.String(), compileOut)

	_, err = os.Stat(exePath)
	require.NoError(t, err, "Expected '%s' to exist", exePath)
	return exePath, remoteConfig
}
