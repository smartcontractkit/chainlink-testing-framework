package executor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type Config struct {
	ProjectPath       string
	Verbose           bool
	GoTestCountFlag   *int
	GoTestRaceFlag    bool
	GoTestTimeoutFlag string
	Tags              []string
	UseShuffle        bool
	ShuffleSeed       string
	SkipTests         []string
	SelectTests       []string
	RawOutputDir      string
}

type Executor interface {
	RunTestPackage(cfg Config, packageName string, runIndex int) (outputFilePath string, passed bool, err error)
	RunCmd(cfg Config, testCmd []string, runIndex int) (outputFilePath string, passed bool, err error)
}

type exitCoder interface {
	ExitCode() int
}

type commandExecutor struct{}

func NewCommandExecutor() Executor {
	return &commandExecutor{}
}

// RunTestPackage runs the tests for a given package and returns the path to the output file.
func (e *commandExecutor) RunTestPackage(cfg Config, packageName string, runIndex int) (string, bool, error) {
	args := []string{"test", packageName, "-json"}
	if cfg.GoTestCountFlag != nil {
		args = append(args, fmt.Sprintf("-count=%d", *cfg.GoTestCountFlag))
	}
	if cfg.GoTestRaceFlag {
		args = append(args, "-race")
	}
	if cfg.GoTestTimeoutFlag != "" {
		args = append(args, fmt.Sprintf("-timeout=%s", cfg.GoTestTimeoutFlag))
	}
	if len(cfg.Tags) > 0 {
		args = append(args, fmt.Sprintf("-tags=%s", strings.Join(cfg.Tags, ",")))
	}
	if cfg.UseShuffle {
		if cfg.ShuffleSeed != "" {
			args = append(args, fmt.Sprintf("-shuffle=%s", cfg.ShuffleSeed))
		} else {
			args = append(args, "-shuffle=on")
		}
	}
	if len(cfg.SkipTests) > 0 {
		skipPattern := strings.Join(cfg.SkipTests, "|")
		args = append(args, fmt.Sprintf("-skip=%s", skipPattern))
	}
	if len(cfg.SelectTests) > 0 {
		selectPattern := strings.Join(cfg.SelectTests, "$|^")
		args = append(args, fmt.Sprintf("-run=^%s$", selectPattern))
	}

	err := os.MkdirAll(cfg.RawOutputDir, 0o755)
	if err != nil {
		return "", false, fmt.Errorf("failed to create raw output directory: %w", err)
	}

	saniPackageName := filepath.Base(packageName)
	tmpFile, err := os.CreateTemp(cfg.RawOutputDir, fmt.Sprintf("test-output-%s-%d-*.json", saniPackageName, runIndex))
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFilePath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		log.Warn().Err(err).Str("file", tempFilePath).Msg("Failed to close temporary file handle after creation")
	}

	absPath, absErr := filepath.Abs(tempFilePath)
	if absErr != nil {
		log.Warn().Err(absErr).Str("relative_path", tempFilePath).Msg("Failed to get absolute path for log message, using relative path")
		absPath = tempFilePath
	}

	if cfg.Verbose {
		log.Info().Str("raw_output_file", absPath).Str("command", fmt.Sprintf("go %s", strings.Join(args, " "))).Msg("Running command")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = cfg.ProjectPath
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	stdoutBytes, err := cmd.Output()

	writeErr := os.WriteFile(tempFilePath, stdoutBytes, 0644)
	if writeErr != nil {
		log.Error().Err(writeErr).Str("file", tempFilePath).Msg("Failed to write captured stdout to temp file")
		_ = os.Remove(tempFilePath)
		return "", false, fmt.Errorf("failed to write command output to %s: %w", tempFilePath, writeErr)
	}

	if err != nil {
		if stderrStr := stderrBuf.String(); stderrStr != "" {
			log.Error().Str("package", packageName).Str("stderr", stderrStr).Msg("Command failed with error and stderr output")
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return tempFilePath, false, nil
		}
		_ = os.Remove(tempFilePath)
		return "", false, fmt.Errorf("test command execution failed for package %s: %w", packageName, err)
	}

	return tempFilePath, true, nil
}

// RunCmd runs the user-supplied command once, captures its JSON output
func (e *commandExecutor) RunCmd(cfg Config, testCmd []string, runIndex int) (tempFilePath string, passed bool, err error) {
	if len(testCmd) == 0 {
		return "", false, errors.New("test command cannot be empty")
	}

	err = os.MkdirAll(cfg.RawOutputDir, 0o755)
	if err != nil {
		return "", false, fmt.Errorf("failed to create raw output directory: %w", err)
	}
	tmpFile, err := os.CreateTemp(cfg.RawOutputDir, fmt.Sprintf("test-output-cmd-run%d-*.json", runIndex+1))
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file for command run: %w", err)
	}
	tempFilePath = tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		log.Warn().Err(err).Str("file", tempFilePath).Msg("Failed to close temporary file handle after creation (cmd run)")
	}

	absPath, absErr := filepath.Abs(tempFilePath)
	if absErr != nil {
		log.Warn().Err(absErr).Str("relative_path", tempFilePath).Msg("Failed to get absolute path for log message (cmd run), using relative path")
		absPath = tempFilePath
	}

	if cfg.Verbose {
		log.Info().Str("raw_output_file", absPath).Str("command", strings.Join(testCmd, " ")).Msg("Running custom command")
	}

	cmd := exec.Command(testCmd[0], testCmd[1:]...) //nolint:gosec
	cmd.Dir = cfg.ProjectPath
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	stdoutBytes, err := cmd.Output()

	writeErr := os.WriteFile(tempFilePath, stdoutBytes, 0644)
	if writeErr != nil {
		log.Error().Err(writeErr).Str("file", tempFilePath).Msg("Failed to write captured stdout to temp file (cmd run)")
		_ = os.Remove(tempFilePath)
		return "", false, fmt.Errorf("failed to write command output to %s: %w", tempFilePath, writeErr)
	}

	if err != nil {
		if stderrStr := stderrBuf.String(); stderrStr != "" {
			log.Error().Str("command", strings.Join(testCmd, " ")).Str("stderr", stderrStr).Msg("Custom command failed with error and stderr output")
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			passed = false
			return tempFilePath, passed, nil
		} else {
			_ = os.Remove(tempFilePath)
			return "", false, fmt.Errorf("error running test command %v: %w", testCmd, err)
		}
	}

	passed = true
	return tempFilePath, passed, nil
}
