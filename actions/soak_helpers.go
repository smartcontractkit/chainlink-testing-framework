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
func BuildGoTests(testTargetDir, finalTestDestination string) (string, string, error) {
	dockerfilePath := filepath.Join(utils.SoakRoot, "Dockerfile.compiler")
	// Clean up old test files if they're around
	if _, err := os.Stat(finalTestDestination); err == nil {
		if err = os.Remove(finalTestDestination); err != nil {
			return "", "", nil
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
		return "", "", err
	}

	err = os.Rename(filepath.Join(testTargetDir, "remote.test"), finalTestDestination)
	if err != nil {
		return "", "", err
	}
	err = os.Remove(testTargetDir)
	if err != nil {
		return "", "", err
	}

	fileInfo, err := os.Stat(finalTestDestination)
	if err != nil {
		return "", "", fmt.Errorf("expected '%s' to exist, %w", finalTestDestination, err)
	}
	return finalTestDestination, strconv.Itoa(int(fileInfo.Size())), nil
}

// RunSoakTest runs a soak test based on the tag, launching as many chainlink nodes as necessary
func RunSoakTest(testTargetDir, finalTestDestination, testTag, namespacePrefix string, chainlinkReplicas int) error {
	LoadConfigs()
	_, fileSize, err := BuildGoTests(testTargetDir, finalTestDestination)
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
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": chainlinkReplicas,
		})).
		Run()
	if err != nil {
		return err
	}
	_, _, errOut, err := env.Client.CopyToPod(
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
