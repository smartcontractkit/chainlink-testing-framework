package soak

//revive:disable:dot-imports
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/gomega"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

// Builds the go tests to run, and returns a path to it, along with remote config options
func buildGoTests() string {
	exePath := filepath.Join(utils.ProjectRoot, "remote.test")
	compileCmd := exec.Command("go", "test", "-c", utils.SoakRoot, "-o", exePath) // #nosec G204
	compileCmd.Env = os.Environ()
	compileCmd.Env = append(compileCmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

	log.Info().Str("Test Directory", utils.SuiteRoot).Msg("Compiling tests")
	compileOut, err := compileCmd.Output()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", compileCmd.String()).
		Msg("Ran command")
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Env: %s\nCommand: %s\nCommand Output: %s", compileCmd.Env, compileCmd.String(), compileOut))

	_, err = os.Stat(exePath)
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Expected '%s' to exist", exePath))
	return exePath
}

func runSoakTest(testTag, namespacePrefix string, chainlinkReplicas int) {
	actions.LoadConfigs()
	exePath := buildGoTests()

	env, err := environment.DeployRemoteRunnerEnvironment(
		environment.NewChainlinkConfig(
			environment.ChainlinkReplicas(chainlinkReplicas, config.ChainlinkVals()),
			namespacePrefix,
			config.GethNetworks()...,
		),
		testTag, // Name of the test to run
		config.ProjectConfig.RemoteRunnerConfig.SlackAPIKey,  // API key to use to upload artifacts to slack
		config.ProjectConfig.RemoteRunnerConfig.SlackChannel, // Slack Channel to upload test artifacts to
		config.ProjectConfig.RemoteRunnerConfig.SlackUserID,  // Slack user to notify on completion
		filepath.Join(utils.SuiteRoot, "framework.yaml"),     // Path of the framework config
		filepath.Join(utils.SuiteRoot, "networks.yaml"),      // Path to the networks config
		exePath, // Path to the executable test file
	)
	Expect(err).ShouldNot(HaveOccurred())
	log.Info().Str("Namespace", env.Namespace).
		Str("Environment File", fmt.Sprintf("%s.%s", env.Namespace, "yaml")).
		Msg("Soak Test Successfully Launched. Save the environment file to collect logs when test is done.")
}
