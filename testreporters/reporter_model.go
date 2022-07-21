package testreporters

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-env/environment"
	"k8s.io/apimachinery/pkg/api/errors"
)

// TestReporter is a general interface for all test reporters
type TestReporter interface {
	WriteReport(folderLocation string) error
	SendSlackNotification(slackClient *slack.Client) error
	SetNamespace(namespace string)
}

const (
	// DefaultArtifactsDir default artifacts dir
	DefaultArtifactsDir string = "logs"
)

// attempts to download the logs of all ephemeral test deployments onto the test runner, also writing a test report
// if one is provided
func WriteTeardownLogs(env *environment.Environment, optionalTestReporter TestReporter) error {
	if ginkgo.CurrentSpecReport().Failed() || optionalTestReporter != nil {
		testFilename := strings.Split(ginkgo.CurrentSpecReport().FileName(), ".")[0]
		_, testName := filepath.Split(testFilename)
		logsPath := filepath.Join(DefaultArtifactsDir, fmt.Sprintf("%s-%d", testName, time.Now().Unix()))
		if err := env.Artifacts.DumpTestResult(logsPath, "chainlink"); err != nil {
			log.Warn().Err(err).Msg("Error trying to collect pod logs")
			if errors.IsForbidden(err) {
				log.Warn().Msg("Unable to gather logs from a remote_test_runner instance. Working on improving this.")
			} else {
				return err
			}
		}
		if err := SendReport(env, logsPath, optionalTestReporter); err != nil {
			log.Warn().Err(err).Msg("Error writing test report")
		}
	}
	return nil
}

// if provided, writes a test report and sends a Slack notification
func SendReport(env *environment.Environment, logsPath string, optionalTestReporter TestReporter) error {
	if optionalTestReporter != nil {
		log.Info().Msg("Writing Test Report")
		optionalTestReporter.SetNamespace(env.Cfg.Namespace)
		err := optionalTestReporter.WriteReport(logsPath)
		if err != nil {
			return err
		}
		err = optionalTestReporter.SendSlackNotification(nil)
		if err != nil {
			return err
		}
	}
	return nil
}
