package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

// BuildGoTestsDocker builds the go tests to run using docker, and returns a path to the test executable, along with
// remote config options. This version usually takes longer to run, but eliminates issues with cross-compilation.
//  Note: currentProjectRootPath and currentSoakTestRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
//  when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
func BuildGoTestsDocker(currentProjectRootPath, currentSoakTestRootPath, testsPath string) (string, error) {
	LoadConfigs()
	dockerfilePath := filepath.Join(utils.SoakRoot, "Dockerfile.compiler")
	testTargetDir := filepath.Join(currentProjectRootPath, "generated_test_dir")
	finalTestDestination := filepath.Join(currentProjectRootPath, "remote.test")
	// Clean up old test files if they're around
	if _, err := os.Stat(finalTestDestination); err == nil {
		if err = os.Remove(finalTestDestination); err != nil {
			return "", err
		}
	}

	// Get the relative paths to directories needed by docker
	relativeTestDirectoryToRootPath, err := filepath.Rel(currentProjectRootPath, testsPath)
	if err != nil {
		return "", err
	}
	log.Info().Str("path", relativeTestDirectoryToRootPath).Msg("docker build arg testDirectory")
	relativeProjectRootPathToRunningTest, err := filepath.Rel(currentSoakTestRootPath, currentProjectRootPath)

	if err != nil {
		return "", err
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
		currentProjectRootPath) // #nosec G204
	dockerBuildCmd.Env = os.Environ()
	log.Info().Str("Docker File", dockerfilePath).Msg("Compiling tests")
	compileOut, err := dockerBuildCmd.CombinedOutput()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", dockerBuildCmd.String()).
		Msg("Ran command")
	if err != nil {
		return "", err
	}

	err = os.Rename(filepath.Join(testTargetDir, "remote.test"), finalTestDestination)
	if err != nil {
		return "", err
	}
	err = os.Remove(testTargetDir)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(finalTestDestination)
	if err != nil {
		return "", fmt.Errorf("Expected '%s' to exist, %w", finalTestDestination, err)
	}
	return finalTestDestination, nil
}

// BuildGoTests builds the go tests using native go cross-compilation to run, and returns a path to the test executable
// along with remote config options. If you run into compilation issues using this function, utilize the BuildGoTestsDocker
// function.
//  Note: currentProjectRootPath and currentSoakTestRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
//  when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
func BuildGoTests(executablePath, testsPath string) (string, error) {
	exePath := filepath.Join(executablePath, "remote.test")
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
		return "", fmt.Errorf("Env: %s\nCommand: %s\nCommand Output: %s, %w", compileCmd.Env, compileCmd.String(), string(compileOut), err)
	}

	_, err = os.Stat(exePath)
	if err != nil {
		return "", fmt.Errorf("Expected '%s' to exist, %w", exePath, err)
	}
	return exePath, nil
}

// RunSoakTest runs a soak test based on the tag, launching as many chainlink nodes as necessary
//  Note: This function will only work for tests running from this repository since paths in utils
//  only point to this package/repository structure. Tests in outside repositories will need their own run function
func RunSoakTest(testTag, namespacePrefix string, chainlinkReplicas int, dockerCompile bool) error {
	soakTestsPath := filepath.Join(utils.SoakRoot, "tests")
	var exePath string
	var err error
	if dockerCompile {
		exePath, err = BuildGoTestsDocker(utils.ProjectRoot, utils.SoakRoot, soakTestsPath)
	} else {
		exePath, err = BuildGoTests(utils.ProjectRoot, soakTestsPath)
	}

	if err != nil {
		return err
	}

	runnerHelmValues := environment.CommonRemoteRunnerValues(
		testTag, // Name of the test to run
		config.ProjectConfig.RemoteRunnerConfig.SlackAPIKey,  // API key to use to upload artifacts to slack
		config.ProjectConfig.RemoteRunnerConfig.SlackChannel, // Slack Channel to upload test artifacts to
		config.ProjectConfig.RemoteRunnerConfig.SlackUserID,  // Slack user to notify on completion
	)
	env, err := environment.DeployRemoteRunnerEnvironment(
		environment.NewChainlinkConfig(
			environment.ChainlinkReplicas(chainlinkReplicas, config.ChainlinkVals()),
			namespacePrefix,
			config.GethNetworks()...,
		),
		filepath.Join(utils.SuiteRoot, "framework.yaml"), // Path of the framework config
		filepath.Join(utils.SuiteRoot, "networks.yaml"),  // Path to the networks config
		exePath, // Path to the executable test file
		runnerHelmValues,
	)
	if err != nil {
		return fmt.Errorf("Error launching soak test environment %w", err)
	}
	log.Info().Str("Namespace", env.Namespace).
		Str("Environment File", fmt.Sprintf("%s.%s", env.Namespace, "yaml")).
		Msg("Soak Test Successfully Launched. Save the environment file to collect logs when test is done.")
	return nil
}
