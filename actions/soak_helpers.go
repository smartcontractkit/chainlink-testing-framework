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
	LoadConfigs()
	dockerfilePath := filepath.Join(utils.SoakRoot, "Dockerfile.compiler")
	testTargetDir := filepath.Join(utils.ProjectRoot, "generated_test_dir")
	finalTestDestination := filepath.Join(utils.ProjectRoot, "remote.test")
	// Clean up old test files if they're around
	if _, err := os.Stat(finalTestDestination); err == nil {
		if err = os.Remove(finalTestDestination); err != nil {
			return "", nil
		}
	}

	// TODO: Docker has a Go API, but it was oddly complicated and not at all documented, and kept failing.
	// So for now, we're doing the tried and true method of plain commands.
	dockerBuildCmd := exec.Command("docker", "build", "-t", "test-compiler", "-f",
		dockerfilePath, "--output", testTargetDir, utils.ProjectRoot) // #nosec G204
	dockerBuildCmd.Env = os.Environ()
	log.Info().Str("Docker File", dockerfilePath).Msg("Compiling tests")
	compileOut, err := dockerBuildCmd.CombinedOutput()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", dockerBuildCmd.String()).
		Msg("Ran command")
	if err != nil {
		return "", err
	}

	err = os.Rename(filepath.Join(testTargetDir, "remote.test"), finalTestDestination)
	if err != nil {
		return "", err
	}
	err = os.Remove(testTargetDir)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(finalTestDestination)
	if err != nil {
		return "", fmt.Errorf("Expected '%s' to exist, %w", finalTestDestination, err)
	}
	return finalTestDestination, nil
}

// runs a soak test based on the tag, launching as many chainlink nodes as necessary
func RunSoakTest(testTag, namespacePrefix string, chainlinkReplicas int) error {
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
