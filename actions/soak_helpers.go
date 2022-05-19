package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

// Soak Test helpers

// BuildGoTests builds the go tests to run, and returns a path to it, along with remote config options
func BuildGoTests(executablePath, testsPath string) (string, error) {
	exePath := filepath.Join(executablePath, "remote.test")
	compileCmd := exec.Command("go", "test", "-ldflags=-s -w", "-c", testsPath, "-o", exePath) // #nosec G204
	compileCmd.Env = os.Environ()
	compileCmd.Env = append(compileCmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

	log.Info().Str("Test Directory", testsPath).Msg("Compiling tests")
	compileOut, err := compileCmd.CombinedOutput()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", compileCmd.String()).
		Msg("Ran command")
	if err != nil {
		return "", fmt.Errorf("Env: %s\nCommand: %s\nCommand Output: %s, %w", compileCmd.Env, compileCmd.String(), string(compileOut), err)
	}

	_, err = os.Stat(exePath)
	if err != nil {
		return "", fmt.Errorf("Expected '%s' to exist, %w", exePath, err)
	}
	return exePath, nil
}

// runs a soak test based on the tag, launching as many chainlink nodes as necessary
func RunSoakTest(testTag, namespacePrefix string, chainlinkReplicas int) error {
	LoadConfigs()
	soakTestsPath := filepath.Join(utils.SoakRoot, "tests")
	exePath, err := BuildGoTests(utils.ProjectRoot, soakTestsPath)
	if err != nil {
		return err
	}

	runnerHelmValues := environment.CommonRemoteRunnerValues(
		testTag, // Name of the test to run
		config.ProjectConfig.RemoteRunnerConfig.SlackAPIKey,  // API key to use to upload artifacts to slack
		config.ProjectConfig.RemoteRunnerConfig.SlackChannel, // Slack Channel to upload test artifacts to
		config.ProjectConfig.RemoteRunnerConfig.SlackUserID,  // Slack user to notify on completion
	)
	env, err := environment.DeployRemoteRunnerEnvironment(
		environment.NewChainlinkConfig(
			environment.ChainlinkReplicas(chainlinkReplicas, config.ChainlinkVals()),
			namespacePrefix,
			config.GethNetworks()...,
		),
		filepath.Join(utils.SuiteRoot, "framework.yaml"), // Path of the framework config
		filepath.Join(utils.SuiteRoot, "networks.yaml"),  // Path to the networks config
		exePath, // Path to the executable test file
		runnerHelmValues,
	)
	if err != nil {
		return fmt.Errorf("Error launching soak test environment %w", err)
	}
	log.Info().Str("Namespace", env.Namespace).
		Str("Environment File", fmt.Sprintf("%s.%s", env.Namespace, "yaml")).
		Msg("Soak Test Successfully Launched. Save the environment file to collect logs when test is done.")
	return nil
}
