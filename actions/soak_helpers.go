// Package actions contains some general functions that help setup tests
package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

// BuildGoTests builds the go tests using native go cross-compilation to run, and returns a path to the test executable
// along with its size in bytes.
//
//	Note: currentProjectRootPath and currentSoakTestRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
//	when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
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

/** This gets complicated with recent refactoring. Seeing how it's a niche use case, we'll remove it for now.
BuildGoTestsWithDocker builds the go tests to run using docker, and returns a path to the test executable, along with
remote config options. This version usually takes longer to run, but eliminates issues with cross-compilation.
 Note: executablePath and projectRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
 when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
func BuildGoTestsWithDocker(executablePath, testsPath, projectRootPath string) (fs.FileInfo, error) {
	dockerfilePath, err := filepath.Abs("./Dockerfile.compiler")
	if err != nil {
		return nil, err
	}
	testTargetDir := filepath.Join(executablePath, "generated_test_dir")
	finalTestDestination := filepath.Join(executablePath, "remote.test")
	// Clean up old test files if they're around
	if _, err := os.Stat(finalTestDestination); err == nil {
		if err = os.Remove(finalTestDestination); err != nil {
			return nil, err
		}
	}

	// Get the relative paths to directories needed by docker
	relativeTestDirectoryToRootPath, err := filepath.Rel(executablePath, testsPath)
	if err != nil {
		return nil, err
	}
	log.Info().Str("path", relativeTestDirectoryToRootPath).Msg("docker build arg testDirectory")
	relativeProjectRootPathToRunningTest, err := filepath.Rel(projectRootPath, executablePath)

	if err != nil {
		return nil, err
	}
	log.Info().Str("path", relativeProjectRootPathToRunningTest).Msg("docker build arg projectRootPath")

	// TODO: Docker has a Go API, but it was oddly complicated and not at all documented, and kept failing.
	// So for now, we're doing the tried and true method of plain commands.
	dockerBuildCmd := exec.Command("docker",
		"build",
		"-t",
		"test-compiler",
		"--build-arg",
		fmt.Sprintf("testDirectory=./%s", relativeTestDirectoryToRootPath),
		"--build-arg",
		fmt.Sprintf("projectRootPath=./%s", relativeProjectRootPathToRunningTest),
		"-f",
		dockerfilePath,
		"--output",
		testTargetDir,
		executablePath) // #nosec G204
	dockerBuildCmd.Env = os.Environ()
	log.Info().Str("Docker File", dockerfilePath).Msg("Compiling tests using Docker")
	compileOut, err := dockerBuildCmd.CombinedOutput()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", dockerBuildCmd.String()).
		Msg("Ran command")
	if err != nil {
		return nil, err
	}

	err = os.Rename(filepath.Join(testTargetDir, "remote.test"), finalTestDestination)
	if err != nil {
		return nil, err
	}
	err = os.Remove(testTargetDir)
	if err != nil {
		return nil, err
	}

	exeFileInfo, err := os.Stat(finalTestDestination)
	if err != nil {
		return nil, fmt.Errorf("Expected '%s' to exist, %w", finalTestDestination, err)
	}
	return exeFileInfo, nil
}
**/
