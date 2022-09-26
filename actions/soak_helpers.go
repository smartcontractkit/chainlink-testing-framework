// Package actions contains some general functions that help setup tests
package actions

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// BuildGoTests builds the go tests using native go cross-compilation to run, and returns a path to the test executable
// along with its size in bytes.
//
//	Note: currentProjectRootPath and currentSoakTestRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
//	when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
func BuildGoTests(executablePath, testsPath, projectRootPath, repoRootPath string) (string, int64, error) {
	logging.Init()
	// absExecutablePath, err := filepath.Abs(executablePath)
	// if err != nil {
	// 	return "", 0, err
	// }
	// absTestsPath, err := filepath.Abs(testsPath)
	// if err != nil {
	// 	return "", 0, err
	// }
	// absProjectRootPath, err := filepath.Abs(projectRootPath)
	// if err != nil {
	// 	return "", 0, err
	// }
	// log.Info().
	// 	Str("Test Directory", absTestsPath).
	// 	Str("Executable Path", absExecutablePath).
	// 	Str("Project Root Path", absProjectRootPath).
	// 	Msg("Compiling tests")

	// dest, size, err :=  buildLocal(absExecutablePath, absTestsPath, absProjectRootPath)
	// dest, size, err := buildDocker(absExecutablePath, absTestsPath, absProjectRootPath)
	dest, size, err := buildDocker(executablePath, testsPath, projectRootPath, repoRootPath)
	return dest, size, err
}

// func buildLocal(executablePath, testsPath, projectRootPath string) (string, int64, error) {
// 	exeFile := filepath.Join(executablePath, "remote.test")
// 	compileCmd := exec.Command("go", "test", "-ldflags=-s -w", "-c", testsPath, "-o", exeFile) // #nosec G204
// 	compileCmd.Env = os.Environ()
// 	compileCmd.Env = append(compileCmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

// 	compileOut, err := compileCmd.CombinedOutput()
// 	log.Debug().
// 		Str("Output", string(compileOut)).
// 		Str("Command", compileCmd.String()).
// 		Msg("Ran command")
// 	if err != nil {
// 		return "", 0, fmt.Errorf("env: %s\nCommand: %s\nCommand Output: %s, %w",
// 			compileCmd.Env, compileCmd.String(), string(compileOut), err)
// 	}

// 	exeFileInfo, err := os.Stat(exeFile)
// 	if err != nil {
// 		return "", 0, fmt.Errorf("expected '%s' to exist, %w", exeFile, err)
// 	}
// 	log.Info().Str("Path", exeFile).Int64("File Size (bytes)", exeFileInfo.Size()).Msg("Compiled tests")
// 	return exeFile, exeFileInfo.Size(), nil
// }

// BuildGoTestsWithDocker Builds the tests inside of a docker container to get a consistently working binary regardless of your machine.
// This is a temporary solution to m1 macs having issues with compiling cosmwasm compared to x84_64 macs
//
//	There are also potential issues with how go is installed on a machine (via brew, executable, asdf, nix) behaving
//		non deterministically for different types of installations
//
// Note: executablePath and projectRootPath are not interchangeable with utils.ProjectRoot and utils.SoakRoot
// when running in outside repositories. Keep an eye on when you need paths leading to this go package vs the current running project.
func buildDocker(executablePath, testsPath, goProjectPath, repoRootPath string) (string, int64, error) {
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

	// clean out old dockerfile
	if _, err := os.Stat(absDockerFilePath); err == nil {
		if err = os.Remove(absDockerFilePath); err != nil {
			return "", 0, err
		}
	}
	// copy the dockerfile over to the repo so we can get the repos context during build
	err = copyDockerfile(dockerfilePath, absDockerFilePath)
	if err != nil {
		return "", 0, err
	}

	// change to the repos root path so docker will get the context it needs consistently
	originalPath, err := os.Getwd()
	if err != nil {
		return "", 0, err
	}
	log.Info().Str("Path", originalPath).Msg("Originating Path")
	err = os.Chdir(absRepoRootPath)
	if err != nil {
		return "", 0, err
	}

	// TODO: Docker has a Go API, but it was oddly complicated and not at well documented, and kept failing.
	// So for now, we're doing the tried and true method of plain commands.
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
		absRepoRootPath)
	// executablePath) // #nosec G204
	dockerBuildCmd.Env = os.Environ()
	log.Info().Str("Docker File", absDockerFilePath).Msg("Compiling tests using Docker")
	log.Info().Str("cmd", dockerBuildCmd.String()).Msg("Docker commands")
	// compileOut, err := dockerBuildCmd.CombinedOutput()
	stderr, err := dockerBuildCmd.StderrPipe()
	if err != nil {
		return "", 0, err
	}
	stdout, err := dockerBuildCmd.StdoutPipe()
	if err != nil {
		return "", 0, err
	}

	// start the command and wrap stderr and stdout
	err = dockerBuildCmd.Start()
	if err != nil {
		return "", 0, err
	}
	go readStdPipe(stderr, "Error")
	go readStdPipe(stdout, "Output")

	// wait for the command to finish
	err = dockerBuildCmd.Wait()
	if err != nil {
		return "", 0, err
	}

	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	panic(err)
	// }
	// td := fmt.Sprintf("./%s", relativeTestDirectoryToRootPath)
	// prf := fmt.Sprintf("./%s", relativeProjectRootPathToRunningTest)

	// opt := &types.ImageBuildOptions{
	// 	Tags: []string{
	// 		"test-compiler",
	// 	},
	// 	Dockerfile: absDockerFilePath,
	// 	BuildArgs: map[string]*string{
	// 		"testDirectory":   &td,
	// 		"projectRootPath": &prf,
	// 	},
	// 	Outputs: []types.ImageBuildOutput{
	// 		{
	// 			Type: "local",
	// 			Attrs: map[string]string{
	// 				"dest": testTargetDir,
	// 			},
	// 		},
	// 	},
	// }
	// buildResponse, err := cli.ImageBuild(context.Background(), nil, *opt)
	// if err == nil {
	// 	fmt.Printf("Error, %v", err)
	// }
	// if err != nil {
	// 	fmt.Printf("%s", err.Error())
	// }
	// defer buildResponse.Body.Close()
	// fmt.Printf("********* %s **********\n", buildResponse.OSType)

	// err = os.Rename(filepath.Join(testTargetDir, "remote.test"), finalTestDestination)
	// if err != nil {
	// 	return "", 0, err
	// }
	// err = os.Remove(testTargetDir)
	// if err != nil {
	// 	return "", 0, err
	// }

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

func copyDockerfile(source, destination string) error {
	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destination, bytesRead, 0644)
	return err
}

func readStdPipe(writer io.ReadCloser, prependOutput string) {
	reader := bufio.NewReader(writer)
	line, err := reader.ReadString('\n')
	for err == nil {
		log.Info().Str(prependOutput, line).Msg("Docker")
		line, err = reader.ReadString('\n')
	}
}

// func DockerBuild(testTargetDir, dockerfilePath, relativeProjectRootPathToRunningTest, relativeTestDirectoryToRootPath string) {
// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	if err != nil {
// 		panic(err)
// 	}
// 	td := fmt.Sprintf("./%s", relativeTestDirectoryToRootPath)
// 	prf := fmt.Sprintf("./%s", relativeProjectRootPathToRunningTest)

// 	opt := &types.ImageBuildOptions{
// 		Tags: []string{
// 			"test-compiler",
// 		},
// 		Dockerfile: dockerfilePath,
// 		BuildArgs: map[string]*string{
// 			"testDirectory":   &td,
// 			"projectRootPath": &prf,
// 		},
// 		Outputs: []types.ImageBuildOutput{
// 			{
// 				Type: "local",
// 				Attrs: map[string]string{
// 					"dest": testTargetDir,
// 				},
// 			},
// 		},
// 	}
// 	buildResponse, err := cli.ImageBuild(context.Background(), nil, *opt)
// 	if err == nil {
// 		fmt.Printf("Error, %v", err)
// 	}
// 	if err != nil {
// 		fmt.Printf("%s", err.Error())
// 	}
// 	defer buildResponse.Body.Close()
// 	fmt.Printf("********* %s **********\n", buildResponse.OSType)
// }

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
