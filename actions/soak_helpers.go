package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/stretchr/testify/require"
)

// Soak Test helpers

// BuildGoTests builds the go tests to run, and returns a path to it, along with remote config options
func BuildGoTests(executablePath, testsPath, projectRootPath string) (string, int64, error) {
	logging.Init()
	absExecutablePath, err := filepath.Abs(executablePath)
	if err != nil {
		return "", 0, err
	}
	absTestsPath, err := filepath.Abs(testsPath)
	if err != nil {
		return "", 0, err
	}
	absProjectRootPath, err := filepath.Abs(projectRootPath)
	if err != nil {
		return "", 0, err
	}
	log.Info().
		Str("Test Directory", absTestsPath).
		Str("Executable Path", absExecutablePath).
		Str("Project Root Path", absProjectRootPath).
		Msg("Compiling tests")

	exeFile := filepath.Join(absExecutablePath, "remote.test")
	compileCmd := exec.Command("go", "test", "-ldflags=-s -w", "-c", absTestsPath, "-o", exeFile) // #nosec G204
	compileCmd.Env = os.Environ()
	compileCmd.Env = append(compileCmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

	compileOut, err := compileCmd.CombinedOutput()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", compileCmd.String()).
		Msg("Ran command")
	if err != nil {
		return "", 0, fmt.Errorf("Env: %s\nCommand: %s\nCommand Output: %s, %w",
			compileCmd.Env, compileCmd.String(), string(compileOut), err)
	}

	exeFileInfo, err := os.Stat(exeFile)
	if err != nil {
		return "", 0, fmt.Errorf("Expected '%s' to exist, %w", exeFile, err)
	}
	log.Info().Str("Path", exeFile).Int64("File Size (bytes)", exeFileInfo.Size()).Msg("Compiled tests")
	return exeFile, exeFileInfo.Size(), nil
}

// RunSoakTest runs a soak test based on the tag, launching as many chainlink nodes as necessary
func RunSoakTest(t *testing.T, testTag, namespacePrefix string, chainlinkValues map[string]interface{}) {
	logging.Init()

	exeFile, exeFileSize, err := BuildGoTests("./", "./tests", "../")
	require.NoError(t, err, "Error building go tests")
	env := environment.New(&environment.Config{
		TTL:             999 * time.Hour,
		Labels:          []string{fmt.Sprintf("envType=%s", pkg.EnvTypeEVM5RemoteRunner), "testType=soak"},
		NamespacePrefix: namespacePrefix,
	})
	err = env.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(remotetestrunner.New(map[string]interface{}{
			"remote_test_runner": map[string]interface{}{
				"test_name":      testTag,
				"env_namespace":  env.Cfg.Namespace,
				"test_file_size": fmt.Sprint(exeFileSize),
			},
		})).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, chainlinkValues)).
		Run()
	require.NoError(t, err, "Error launching test environment")
	err = TriggerRemoteTest(exeFile, env)
	require.NoError(t, err, "Error activating remote test")
}

// TriggerRemoteTest copies the executable to the remote-test-runner and starts the run
func TriggerRemoteTest(exePath string, testEnvironment *environment.Environment) error {
	logging.Init()

	_, _, errOut, err := testEnvironment.Client.CopyToPod(
		testEnvironment.Cfg.Namespace,
		exePath,
		fmt.Sprintf("%s/%s:/root/remote.test", testEnvironment.Cfg.Namespace, "remote-test-runner"),
		"remote-test-runner")
	if err != nil {
		return errors.Wrap(err, errOut.String())
	}
	log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Remote Test Triggered on 'remote-test-runner'")
	return nil
}
