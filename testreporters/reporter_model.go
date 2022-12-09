package testreporters

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-env/environment"
)

// TestReporter is a general interface for all test reporters
type TestReporter interface {
	WriteReport(folderLocation string) error
	SendSlackNotification(t *testing.T, slackClient *slack.Client) error
	SetNamespace(namespace string)
}

const (
	// DefaultArtifactsDir default artifacts dir
	DefaultArtifactsDir string = "logs"
)

// attempts to download the logs of all ephemeral test deployments onto the test runner, also writing a test report
// if one is provided
func WriteTeardownLogs(t *testing.T, env *environment.Environment, optionalTestReporter TestReporter) error {
	if t.Failed() || optionalTestReporter != nil {
		logsPath := filepath.Join(DefaultArtifactsDir, fmt.Sprintf("%s-%s-%d", t.Name(), env.Cfg.Namespace, time.Now().Unix()))
		if err := env.Artifacts.DumpTestResult(logsPath, "chainlink"); err != nil {
			log.Warn().Err(err).Msg("Error trying to collect pod logs")
			return err
		}
		if err := SendReport(t, env, logsPath, optionalTestReporter); err != nil {
			log.Warn().Err(err).Msg("Error writing test report")
		}
	}
	return nil
}

// if provided, writes a test report and sends a Slack notification
func SendReport(t *testing.T, env *environment.Environment, logsPath string, optionalTestReporter TestReporter) error {
	if optionalTestReporter != nil {
		log.Info().Msg("Writing Test Report")
		optionalTestReporter.SetNamespace(env.Cfg.Namespace)
		err := optionalTestReporter.WriteReport(logsPath)
		if err != nil {
			return err
		}
		err = optionalTestReporter.SendSlackNotification(t, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
