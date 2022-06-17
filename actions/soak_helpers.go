package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/remotetestrunner"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Soak Test helpers

// BuildGoTests builds the go tests to run, and returns a path to it, along with remote config options
func BuildGoTests(testsPath, exePath string) (string, string, error) {
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
		return "", "", fmt.Errorf("Env: %s\nCommand: %s\nCommand Output: %s, %w", compileCmd.Env, compileCmd.String(), string(compileOut), err)
	}

	fileInfo, err := os.Stat(exePath)
	if err != nil {
		return "", "", fmt.Errorf("expected '%s' to exist, %w", exePath, err)
	}
	return exePath, strconv.Itoa(int(fileInfo.Size())), nil
}

// RunSoakTest runs a soak test based on the tag, launching as many chainlink nodes as necessary
func RunSoakTest(testDirPath, exePath, testTag, namespacePrefix string, chainlinkReplicas int) error {
	LoadConfigs()
	_, fileSize, err := BuildGoTests(testDirPath, exePath)
	if err != nil {
		return err
	}
	env := environment.New(&environment.Config{
		TTL:       168 * time.Hour,
		Labels:    []string{fmt.Sprintf("envType=%s", pkg.EnvTypeEVM5RemoteRunner)},
		Namespace: namespacePrefix,
	})
	err = env.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(remotetestrunner.New(map[string]interface{}{
			"remote_test_runner": map[string]interface{}{
				"test_name":      testTag,
				"env_namespace":  env.Cfg.Namespace,
				"slack_api":      config.ProjectConfig.RemoteRunnerConfig.SlackAPIKey,
				"slack_channel":  config.ProjectConfig.RemoteRunnerConfig.SlackChannel,
				"slack_user_id":  config.ProjectConfig.RemoteRunnerConfig.SlackUserID,
				"test_file_size": fileSize,
				"access_port":    8080,
			},
		})).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(map[string]interface{}{
			"replicas": chainlinkReplicas,
		})).
		Run()
	if err != nil {
		return err
	}
	_, _, errOut, err := env.Client.CopyToPod(
		env.Cfg.Namespace,
		filepath.Join(utils.SuiteRoot, "framework.yaml"),
		fmt.Sprintf("%s/%s:/root/framework.yaml", env.Cfg.Namespace, "remote-test-runner"),
		"remote-test-runner")
	if err != nil {
		return errors.Wrap(err, errOut.String())
	}
	_, _, errOut, err = env.Client.CopyToPod(
		env.Cfg.Namespace,
		filepath.Join(utils.SuiteRoot, "networks.yaml"),
		fmt.Sprintf("%s/%s:/root/networks.yaml", env.Cfg.Namespace, "remote-test-runner"),
		"remote-test-runner")
	if err != nil {
		return errors.Wrap(err, errOut.String())
	}
	_, _, errOut, err = env.Client.CopyToPod(
		env.Cfg.Namespace,
		filepath.Join(utils.ProjectRoot, "remote.test"),
		fmt.Sprintf("%s/%s:/root/remote.test", env.Cfg.Namespace, "remote-test-runner"),
		"remote-test-runner")
	if err != nil {
		return errors.Wrap(err, errOut.String())
	}
	log.Info().Str("Namespace", env.Cfg.Namespace).
		Str("Environment File", fmt.Sprintf("%s.%s", env.Cfg.Namespace, "yaml")).
		Msg("Soak Test Successfully Launched. Save the environment file to collect logs when test is done.")
	return nil
}
