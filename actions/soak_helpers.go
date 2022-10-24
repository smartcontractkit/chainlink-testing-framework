// Package actions contains some general functions that help setup tests
package actions

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

const containerName = "remote-test-runner"
const tarFileName = "ztarrepo.tar.gz"

// BasicRunnerValuesSetup Basic values needed to run a soak test in the remote runner
func BasicRunnerValuesSetup(focus, namespace, testDirectory string) map[string]interface{} {
	return map[string]interface{}{
		"focus":         focus,
		"env_namespace": namespace,
		"test_dir":      testDirectory,
	}
}

// TriggerRemoteTest copies the executable to the remote-test-runner and starts the run
func TriggerRemoteTest(repoSourcePath string, testEnvironment *environment.Environment) error {
	logging.Init()

	// Leaving commented out for now since it may be useful for the final solution to compiling cosmwasm into the tests
	// Add gcompat package, required by libwasmvm
	// _, _, err := testEnvironment.Client.ExecuteInPod(testEnvironment.Cfg.Namespace, containerName, containerName, []string{"apk", "add", "gcompat"})
	// if err != nil {
	// 	return errors.Wrap(err, "Error adding gcompat")
	// }
	// Copy libwasmvm dependency of chainlink core
	// _, _, errOut, err := testEnvironment.Client.CopyToPod(
	// 	testEnvironment.Cfg.Namespace,
	// 	os.Getenv("GOPATH")+"/pkg/mod/github.com/!cosm!wasm/wasmvm@v1.0.0/api/libwasmvm.x86_64.so",
	// 	fmt.Sprintf("%s/%s:/usr/lib/libwasmvm.x86_64.so", testEnvironment.Cfg.Namespace, containerName),
	// 	"remote-test-runner")
	// if err != nil {
	// 	return errors.Wrap(err, errOut.String())
	// }

	// tar the repo
	tarPath, _, err := tarRepo(repoSourcePath)
	if err != nil {
		return err
	}
	log.Debug().Str("Path", tarPath).Msg("Tar file to copy to pod")

	// copy the repo containing the test to the pod
	_, _, errOut, err := testEnvironment.Client.CopyToPod(
		testEnvironment.Cfg.Namespace,
		tarPath,
		fmt.Sprintf("%s/%s:/root/", testEnvironment.Cfg.Namespace, containerName),
		containerName)
	if err != nil {
		return errors.Wrap(err, errOut.String())
	}

	// create start file in pod to start the test
	_, _, err = testEnvironment.Client.ExecuteInPod(
		testEnvironment.Cfg.Namespace,
		containerName,
		containerName,
		[]string{
			"touch",
			"/root/start.txt",
		},
	)
	if err != nil {
		return err
	}

	log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg(fmt.Sprintf("Remote Test Triggered on '%s'", containerName))
	return nil
}

// tarRep Uses the gnu tar command from a docker image to consistently compress the code
// using .gitingore and some other excludes to reduce the file size
func tarRepo(repoRootPath string) (string, int64, error) {
	absRepoRootPath, err := filepath.Abs(repoRootPath)
	if err != nil {
		return "", 0, err
	}
	err = os.Chdir(absRepoRootPath)
	if err != nil {
		return "", 0, err
	}

	absTarFilePath := filepath.Join(absRepoRootPath, tarFileName)
	// Clean up old test files if they're around
	if _, err := os.Stat(absTarFilePath); err == nil {
		if err = os.Remove(absTarFilePath); err != nil {
			return "", 0, err
		}
	}

	dockerBuildCmd := exec.Command("docker",
		"run",
		"--rm",
		"-v",
		fmt.Sprintf("%s:/usr/src", absRepoRootPath),
		"tateexon/tar@sha256:093ca3ba5dbdd906d03425cbec3296705a6aa316db2071f4c8b487504de9e129",
		fmt.Sprintf("--exclude=%s", tarFileName),
		"--exclude=.git",
		"--exclude-vcs-ignores",
		"--ignore-failed-read",
		"-czvf",
		fmt.Sprintf("/usr/src/%s", tarFileName),
		"/usr/src/",
	) // #nosec G204
	dockerBuildCmd.Env = os.Environ()
	log.Info().Str("cmd", dockerBuildCmd.String()).Msg("Docker command")
	stderr, err := dockerBuildCmd.StderrPipe()
	if err != nil {
		return "", 0, err
	}
	stdout, err := dockerBuildCmd.StdoutPipe()
	if err != nil {
		return "", 0, err
	}

	// start the command and wrap stderr and stdout
	started := time.Now()
	err = dockerBuildCmd.Start()
	if err != nil {
		return "", 0, err
	}
	go readStdPipeDocker(stderr, "Error")
	go readStdPipeDocker(stdout, "Output")

	// wait for the command to finish
	err = dockerBuildCmd.Wait()
	if err != nil {
		log.Info().Err(err).Msg("Ignoring error since it always fails for directory changing.")
	}

	finished := time.Now()
	log.Info().Str("total", fmt.Sprintf("%v", finished.Sub(started))).Msg("Docker Command Run Time")

	tarFileInfo, err := os.Stat(absTarFilePath)
	if err != nil {
		return "", 0, fmt.Errorf("expected '%s' to exist, %w", absTarFilePath, err)
	}
	log.Info().Str("Path", absTarFilePath).Int64("File Size (bytes)", tarFileInfo.Size()).Msg("Compiled tests")

	return absTarFilePath, tarFileInfo.Size(), nil
}

// readStdPipeDocker continuously read a pipe from the docker command
func readStdPipeDocker(writer io.ReadCloser, prependOutput string) {
	reader := bufio.NewReader(writer)
	line, err := reader.ReadString('\n')
	for err == nil {
		log.Info().Str(prependOutput, line).Msg("Docker")
		line, err = reader.ReadString('\n')
	}
}
