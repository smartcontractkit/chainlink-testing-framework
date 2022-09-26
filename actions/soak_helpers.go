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
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// BuildGoTests Builds the tests inside of a docker container to get a consistently working binary regardless of your machine.
// Builds the go tests and returns a path to the test executable along with its size in bytes.
// All paths should be relative to the repoRootPath.
// The repoRootPath needs to be relative to where you are running go test from.
// This is a temporary solution to m1 macs having issues with compiling cosmwasm compared to x84_64 macs
//
//	There are also potential issues with how go is installed on a machine (via brew, executable, asdf, nix) behaving
//		non deterministically for different types of installations
//
// Note: executablePath and projectRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
// when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
func BuildGoTests(repoRootPath, executablePath, testsPath, goProjectPath string) (string, int64, error) {
	logging.Init()
	// build all the paths
	absRepoRootPath, err := filepath.Abs(repoRootPath)
	if err != nil {
		return "", 0, err
	}
	absTestsPath := filepath.Join(absRepoRootPath, testsPath)
	absGoProjectPath := filepath.Join(absRepoRootPath, goProjectPath)
	absExecutablePath := filepath.Join(absRepoRootPath, executablePath)
	dockerFileName := "./Dockerfile.compiler"
	dockerfilePath, err := filepath.Abs(filepath.Join(utils.ProjectRoot, dockerFileName))
	if err != nil {
		return "", 0, err
	}
	absDockerFilePath := filepath.Join(absExecutablePath, dockerFileName)
	builtExecutablePath := filepath.Join(absExecutablePath, "remote.test")
	// Clean up old test files if they're around
	if _, err := os.Stat(builtExecutablePath); err == nil {
		if err = os.Remove(builtExecutablePath); err != nil {
			return "", 0, err
		}
	}

	// Get the relative paths to directories needed by docker
	relativeTestDirectoryToGoProjectPath, err := filepath.Rel(absGoProjectPath, absTestsPath)
	if err != nil {
		return "", 0, err
	}
	relativeGoProjectPath, err := filepath.Rel(absRepoRootPath, absGoProjectPath)
	if err != nil {
		return "", 0, err
	}

	// log all the paths for debugging
	log.Info().
		Str("dockerFilePath", absDockerFilePath).
		Str("testsPath", absTestsPath).
		Str("goProjectPath", absGoProjectPath).
		Str("executablePath", absExecutablePath).
		Str("repoRootPath", absRepoRootPath).
		Str("relativeTestDirectoryToGoProjectPath", relativeTestDirectoryToGoProjectPath).
		Str("relativeGoProjectPath", relativeGoProjectPath).
		Msg("Paths")

	// copy the dockerfile over to the repo so we can get the repos context during build
	err = copyDockerfile(dockerfilePath, absDockerFilePath)
	if err != nil {
		return "", 0, err
	}

	// change to the repos root path so docker will get the context it needs
	originalPath, err := os.Getwd()
	if err != nil {
		return "", 0, err
	}
	log.Info().Str("Path", originalPath).Msg("Originating Path For Test Run")
	err = os.Chdir(absRepoRootPath)
	if err != nil {
		return "", 0, err
	}

	// build the docker command to run
	dockerBuildCmd := exec.Command("docker",
		"build",
		"-t",
		"test-compiler",
		"--build-arg",
		fmt.Sprintf("testDirectory=./%s", relativeTestDirectoryToGoProjectPath),
		"--build-arg",
		fmt.Sprintf("goProjectPath=./%s", relativeGoProjectPath),
		"-f",
		absDockerFilePath,
		"--output",
		absExecutablePath,
		absRepoRootPath) // #nosec G204
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
		return "", 0, err
	}
	finished := time.Now()
	log.Info().Str("total", fmt.Sprintf("%v", finished.Sub(started))).Msg("Docker Command Run Time")

	exeFileInfo, err := os.Stat(builtExecutablePath)
	if err != nil {
		return "", 0, fmt.Errorf("expected '%s' to exist, %w", builtExecutablePath, err)
	}
	log.Info().Str("Path", builtExecutablePath).Int64("File Size (bytes)", exeFileInfo.Size()).Msg("Compiled tests")

	// change back to original directory so any tests running after will behave like normal
	err = os.Chdir(originalPath)
	if err != nil {
		return "", 0, err
	}

	return builtExecutablePath, exeFileInfo.Size(), nil
}

// copyDockerfile copies the dockerfile from the chainlink-testing-framework to the repo it needs to run in
func copyDockerfile(source, destination string) error {
	// clean out old dockerfile first
	if _, err := os.Stat(destination); err == nil {
		if err = os.Remove(destination); err != nil {
			return err
		}
	}

	bytesRead, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	err = os.WriteFile(destination, bytesRead, 0600)
	return err
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

// TriggerRemoteTest copies the executable to the remote-test-runner and starts the run
func TriggerRemoteTest(exePath string, testEnvironment *environment.Environment) error {
	logging.Init()

	// Leaving commented out for now since it may be useful for the final solution to compiling cosmwasm into the tests
	// Add gcompat package, required by libwasmvm
	// _, _, err := testEnvironment.Client.ExecuteInPod(testEnvironment.Cfg.Namespace, "remote-test-runner", "remote-test-runner", []string{"apk", "add", "gcompat"})
	// if err != nil {
	// 	return errors.Wrap(err, "Error adding gcompat")
	// }
	// Copy libwasmvm dependency of chainlink core
	// _, _, errOut, err := testEnvironment.Client.CopyToPod(
	// 	testEnvironment.Cfg.Namespace,
	// 	os.Getenv("GOPATH")+"/pkg/mod/github.com/!cosm!wasm/wasmvm@v1.0.0/api/libwasmvm.x86_64.so",
	// 	fmt.Sprintf("%s/%s:/usr/lib/libwasmvm.x86_64.so", testEnvironment.Cfg.Namespace, "remote-test-runner"),
	// 	"remote-test-runner")
	// if err != nil {
	// 	return errors.Wrap(err, errOut.String())
	// }
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
