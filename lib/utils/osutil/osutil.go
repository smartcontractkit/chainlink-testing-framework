package osutil

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
)

// GetEnv returns the value of the environment variable named by the key
// and sets the environment variable up to be used in the remote runner
func GetEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val != "" {
		prefixedKey := fmt.Sprintf("%s%s", config.EnvVarPrefix, key)
		if os.Getenv(prefixedKey) != "" && os.Getenv(prefixedKey) != val {
			return val, fmt.Errorf("environment variable collision with prefixed key, Original: %s=%s, Prefixed: %s=%s", key, val, prefixedKey, os.Getenv(prefixedKey))
		}
		err := os.Setenv(prefixedKey, val)
		if err != nil {
			return val, fmt.Errorf("failed to set environment variable %s err: %w", prefixedKey, err)
		}
	}
	return val, nil
}

// SetupEnvVarsForRemoteRunner sets up the environment variables in the list to propagate to the remote runner
func SetupEnvVarsForRemoteRunner(envVars []string) error {
	for _, envVar := range envVars {
		_, err := GetEnv(envVar)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecCmd executes a command and logs the output
func ExecCmd(l zerolog.Logger, command string) error {
	return ExecCmdWithContext(context.Background(), l, command)
}

// ExecCmdWithContext executes a command with ctx and logs the output
func ExecCmdWithContext(ctx context.Context, l zerolog.Logger, command string) error {
	return ExecCmdWithOptions(ctx, l, command, func(m string) {
		l.Debug().Str("Text", m).Msg("Std Pipe")
	})
}

// readStdPipe continuously read a pipe from the command
func readStdPipe(pipe io.ReadCloser, outputFunction func(string)) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if outputFunction != nil {
			outputFunction(m)
		}
	}
}

// ExecCmdWithOptions executes a command with ctx and logs the output with a custom logging func
func ExecCmdWithOptions(ctx context.Context, l zerolog.Logger, command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
	l.Info().Interface("Command", c).Msg("Executing command")
	cmd := exec.CommandContext(ctx, c[0], c[1:]...) // #nosec: G204
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go readStdPipe(stderr, outputFunction)
	go readStdPipe(stdout, outputFunction)
	return cmd.Wait()
}

func GetAbsoluteFolderPath(folder string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, folder), nil
}

const (
	DEFAULT_STOP_FILE_NAME         = ".root_dir"
	ErrStopFileNotFoundWithinLimit = "stop file not found in any parent directory within search limit"
)

// FindFile looks for given file in the current directory and its parent directories first by locating
// top level parent folder where the search should begin (defined by stopFile, which cannot be located
// further "up" than limit parent folders) and then by searching from there for the file all subdirectories
func FindFile(filename, stopFile string, limit int) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	currentFilePath := filepath.Join(workingDir, filename)
	if _, err := os.Stat(currentFilePath); err == nil {
		return currentFilePath, nil
	}

	rootDirPath, err := findTopParentFolderWithLimit(workingDir, stopFile, limit)
	if err != nil {
		return "", err
	}

	configFilePath, err := findFileInSubfolders(rootDirPath, filename)
	if err != nil {
		return "", err
	}
	return configFilePath, nil
}

func findTopParentFolderWithLimit(startDir, stopFileName string, limit int) (string, error) {
	currentDir := startDir

	for level := 0; level < limit; level++ {
		stopFilePath := filepath.Join(currentDir, stopFileName)
		if _, err := os.Stat(stopFilePath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("%s: %s", ErrStopFileNotFoundWithinLimit, stopFileName)
}

func findFileInSubfolders(startDir, targetFileName string) (string, error) {
	var filePath string
	ErrFileFound := "file found"

	err := filepath.WalkDir(startDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == targetFileName {
			filePath = path
			return errors.New(ErrFileFound)
		}

		return nil
	})

	if err != nil && err.Error() != ErrFileFound {
		return "", err
	}

	if filePath == "" {
		return "", os.ErrNotExist
	}

	return filePath, nil
}

// FindDirectoriesContainingFile finds all directories containing a file matching the given regular expression
func FindDirectoriesContainingFile(dir string, r *regexp.Regexp) ([]string, error) {
	foundDirs := []string{}
	// Walk through all files in the directory and its sub-directories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if r.MatchString(info.Name()) {
			foundDirs = append(foundDirs, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the directory: %v\n", err)
		return nil, err
	}
	return foundDirs, nil
}
