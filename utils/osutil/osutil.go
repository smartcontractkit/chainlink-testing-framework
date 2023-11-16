package osutil

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
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
