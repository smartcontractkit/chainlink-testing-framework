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

// Config holds configuration relevant to command execution.
// This might need adjustment based on what parts of the original Runner's config
// are purely related to execution flags vs. parsing/reporting.
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
	RawOutputDir      string // Directory for raw command output
}

// Executor defines the interface for running test commands.
type Executor interface {
	// RunTestPackage executes 'go test' for a specific package and configuration.
	// It writes the raw JSON output to a temporary file and returns the path.
	// It also returns whether the test command exited with code 0.
	RunTestPackage(cfg Config, packageName string, runIndex int) (outputFilePath string, passed bool, err error)

	// RunCmd executes an arbitrary command expected to produce Go test JSON output.
	// It writes the raw JSON output to a temporary file and returns the path.
	// It also returns whether the command exited with code 0.
	RunCmd(cfg Config, testCmd []string, runIndex int) (outputFilePath string, passed bool, err error)
}

type exitCoder interface {
	ExitCode() int
}

// commandExecutor uses os/exec to run commands.
type commandExecutor struct{}

// NewCommandExecutor creates a default executor using os/exec.
func NewCommandExecutor() Executor {
	return &commandExecutor{}
}

// RunTestPackage runs the tests for a given package and returns the path to the output file.
// This version captures stdout explicitly before writing to the file.
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

	// Create a temporary file to get a unique path, then close it.
	saniPackageName := filepath.Base(packageName)
	tmpFile, err := os.CreateTemp(cfg.RawOutputDir, fmt.Sprintf("test-output-%s-%d-*.json", saniPackageName, runIndex))
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFilePath := tmpFile.Name()          // Get the unique path
	if err := tmpFile.Close(); err != nil { // Close the handle immediately
		log.Warn().Err(err).Str("file", tempFilePath).Msg("Failed to close temporary file handle after creation")
		// Continue execution, as the path is still valid for writing later
	}

	// Get absolute path for logging
	absPath, absErr := filepath.Abs(tempFilePath)
	if absErr != nil {
		log.Warn().Err(absErr).Str("relative_path", tempFilePath).Msg("Failed to get absolute path for log message, using relative path")
		absPath = tempFilePath // Fallback to relative path
	}

	if cfg.Verbose {
		log.Info().Str("raw_output_file", absPath).Str("command", fmt.Sprintf("go %s", strings.Join(args, " "))).Msg("Running command")
	}

	// Prepare the command
	cmd := exec.Command("go", args...)
	cmd.Dir = cfg.ProjectPath
	// Capture stderr separately (useful for build errors, etc.)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// Run the command and capture stdout
	stdoutBytes, err := cmd.Output() // This runs the command and waits for completion

	// Always write captured stdout to the file, even if the command failed (it might contain partial JSON or error info)
	writeErr := os.WriteFile(tempFilePath, stdoutBytes, 0644)
	if writeErr != nil {
		log.Error().Err(writeErr).Str("file", tempFilePath).Msg("Failed to write captured stdout to temp file")
		// If writing failed, the file is likely unusable. Attempt to remove it and return error.
		_ = os.Remove(tempFilePath) // Best effort removal
		return "", false, fmt.Errorf("failed to write command output to %s: %w", tempFilePath, writeErr)
	}

	// Now check the error from cmd.Output()
	if err != nil {
		// Log the captured stderr for debugging, especially on errors
		if stderrStr := stderrBuf.String(); stderrStr != "" {
			log.Error().Str("package", packageName).Str("stderr", stderrStr).Msg("Command failed with error and stderr output")
		}

		var exitErr *exec.ExitError // Check specifically for non-zero exit errors
		if errors.As(err, &exitErr) {
			// Non-zero exit code => test failure (the output file should contain details)
			return tempFilePath, false, nil // Return path, indicate failure, nil error (command ran)
		}
		// Otherwise, it's likely an error *running* the command itself (e.g., not found)
		_ = os.Remove(tempFilePath) // Remove the empty/potentially misleading file
		return "", false, fmt.Errorf("test command execution failed for package %s: %w", packageName, err)
	}

	// Command ran successfully and exited with 0
	return tempFilePath, true, nil
}

// runCmd runs the user-supplied command once, captures its JSON output,
// and returns the temp file path, whether the test passed, and an error if any.
// Apply similar logic: capture output then write.
func (e *commandExecutor) RunCmd(cfg Config, testCmd []string, runIndex int) (tempFilePath string, passed bool, err error) {
	if len(testCmd) == 0 {
		return "", false, errors.New("test command cannot be empty")
	}

	err = os.MkdirAll(cfg.RawOutputDir, 0o755)
	if err != nil {
		return "", false, fmt.Errorf("failed to create raw output directory: %w", err)
	}
	// Create temp file path
	tmpFile, err := os.CreateTemp(cfg.RawOutputDir, fmt.Sprintf("test-output-cmd-run%d-*.json", runIndex+1))
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file for command run: %w", err)
	}
	tempFilePath = tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		log.Warn().Err(err).Str("file", tempFilePath).Msg("Failed to close temporary file handle after creation (cmd run)")
	}

	// Get absolute path for logging
	absPath, absErr := filepath.Abs(tempFilePath)
	if absErr != nil {
		log.Warn().Err(absErr).Str("relative_path", tempFilePath).Msg("Failed to get absolute path for log message (cmd run), using relative path")
		absPath = tempFilePath // Fallback to relative path
	}

	if cfg.Verbose {
		log.Info().Str("raw_output_file", absPath).Str("command", strings.Join(testCmd, " ")).Msg("Running custom command")
	}

	cmd := exec.Command(testCmd[0], testCmd[1:]...) //nolint:gosec
	cmd.Dir = cfg.ProjectPath
	// Capture stderr
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// Capture stdout
	stdoutBytes, err := cmd.Output() // Run and capture

	// Write captured output
	writeErr := os.WriteFile(tempFilePath, stdoutBytes, 0644)
	if writeErr != nil {
		log.Error().Err(writeErr).Str("file", tempFilePath).Msg("Failed to write captured stdout to temp file (cmd run)")
		_ = os.Remove(tempFilePath)
		return "", false, fmt.Errorf("failed to write command output to %s: %w", tempFilePath, writeErr)
	}

	// Check execution error
	if err != nil {
		if stderrStr := stderrBuf.String(); stderrStr != "" {
			log.Error().Str("command", strings.Join(testCmd, " ")).Str("stderr", stderrStr).Msg("Custom command failed with error and stderr output")
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// Non-zero exit code => command failed (likely test failure)
			passed = false
			return tempFilePath, passed, nil // Return path, pass status false, nil error (command ran)
		} else {
			// Some other error running the command (e.g., command not found)
			_ = os.Remove(tempFilePath)
			return "", false, fmt.Errorf("error running test command %v: %w", testCmd, err)
		}
	}

	// Otherwise, command exited 0 => passed
	passed = true
	return tempFilePath, passed, nil
}
