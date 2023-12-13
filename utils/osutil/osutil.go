package osutil

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
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

const DEFAULT_STOP_FILE_NAME = ".root_dir"

func FindFile(filename, stopFile string) ([]byte, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Find "stopFile" to determine the starting point
	rootDirPath, err := findFileRecursivelyWithLimit(currentDir, stopFile, "", 2)
	if err != nil {
		return nil, err // Return an error if "stopFile" is not found within the limit
	}

	// Use the location of "stopFile" as the starting point to find "filename"
	configFileContent, err := findFileRecursively(rootDirPath, filename, "")
	if err != nil {
		return nil, err
	}

	return configFileContent, nil
}

func findFileRecursively(startDir, targetFileName, stopFileName string) ([]byte, error) {
	var fileContent []byte

	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has the specified name
		if info.IsDir() {
			return nil // Skip directories
		}

		if info.Name() == targetFileName {
			// Read the content of the file
			fileContent, err = os.ReadFile(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if fileContent == nil {
		return nil, os.ErrNotExist // File not found
	}

	return fileContent, nil
}

func findFileRecursivelyWithLimit(startDir, targetFileName, stopFileName string, limit int) (string, error) {
	var filePath string
	parentLevel := 0

	for parentLevel <= limit {
		err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if the file has the specified name
			if !info.IsDir() && info.Name() == targetFileName {
				// Set the filePath when the target file is found
				filePath = path
			}

			return nil
		})

		if err != nil {
			return "", err
		}

		if filePath != "" {
			break // Exit the loop if ".root_dir" is found
		}

		// Move to the parent directory
		startDir = filepath.Dir(startDir)
		parentLevel++
	}

	if filePath == "" {
		return "", os.ErrNotExist // File not found
	}

	return filepath.Dir(filePath), nil
}
