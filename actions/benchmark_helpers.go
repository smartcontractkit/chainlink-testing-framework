package actions

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

// Benchmark Test helpers

// runs a benchmark test based on the tag, launching as many chainlink nodes as necessary
func RunBenchmarkTest(testTag, namespacePrefix string, chainlinkReplicas int) error {
	LoadConfigs()
	benchmarkTestsPath := filepath.Join(utils.BenchmarkRoot, "tests")
	exePath, err := BuildGoTests(utils.ProjectRoot, benchmarkTestsPath)
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
		return fmt.Errorf("Error launching benchmark test environment %w", err)
	}
	log.Info().Str("Namespace", env.Namespace).
		Str("Environment File", fmt.Sprintf("%s.%s", env.Namespace, "yaml")).
		Msg("Benchmark Test Successfully Launched. Save the environment file to collect logs when test is done.")
	return nil
}
