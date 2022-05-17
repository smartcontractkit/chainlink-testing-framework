package soak_runner

//revive:disable:dot-imports
import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

var _ = Describe("OCR Soak Setup @setup-soak", func() {
	It("Deploys soak test runner", func() {
		actions.LoadConfigs()
		exePath := buildGoTests()

		env, err := environment.DeployRemoteRunnerEnvironment(
			environment.NewChainlinkConfig(
				environment.ChainlinkReplicas(6, config.ChainlinkVals()),
				"chainlink-soak",
				config.GethNetworks()...,
			),
			"@soak-ocr", // Name of the test to run
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
	})
})
