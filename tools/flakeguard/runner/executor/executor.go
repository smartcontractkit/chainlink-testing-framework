package executor

import (
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
// This is moved from the original Runner.
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
	// Create a temporary file to store the output
	saniPackageName := filepath.Base(packageName)
	tmpFile, err := os.CreateTemp(cfg.RawOutputDir, fmt.Sprintf("test-output-%s-%d-*.json", saniPackageName, runIndex))
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close() // Close the file descriptor, but the file remains

	if cfg.Verbose {
		log.Info().Str("raw_output_file", tmpFile.Name()).Str("command", fmt.Sprintf("go %s", strings.Join(args, " "))).Msg("Running command")
	}

	// Run the command with output directed to the file
	cmd := exec.Command("go", args...)
	cmd.Dir = cfg.ProjectPath
	cmd.Stdout = tmpFile
	cmd.Stderr = os.Stderr // Capture stderr for build errors, etc.

	err = cmd.Run()
	if err != nil {
		var exErr exitCoder
		// Check if the error is due to a non-zero exit code (test failure)
		if errors.As(err, &exErr) {
			// Non-zero exit code => test failure (or other non-zero exit)
			return tmpFile.Name(), false, nil // Return path, indicate failure, nil error (command ran)
		}
		// Otherwise, it's an error running the command itself
		return "", false, fmt.Errorf("test command execution failed for package %s: %w", packageName, err)
	}

	// Command ran and exited with 0
	return tmpFile.Name(), true, nil
}

// runCmd runs the user-supplied command once, captures its JSON output,
// and returns the temp file path, whether the test passed, and an error if any.
// This is moved from the original Runner.
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
	defer tmpFile.Close() // Close the file descriptor

	tempFilePath = tmpFile.Name()

	if cfg.Verbose {
		log.Info().Str("raw_output_file", tempFilePath).Str("command", strings.Join(testCmd, " ")).Msg("Running custom command")
	}

	cmd := exec.Command(testCmd[0], testCmd[1:]...) //nolint:gosec
	cmd.Dir = cfg.ProjectPath
	cmd.Stdout = tmpFile
	cmd.Stderr = os.Stderr // Capture stderr

	err = cmd.Run()

	// Determine pass/fail from exit code
	var ec exitCoder
	if errors.As(err, &ec) {
		// Non-zero exit code => command failure (likely test failure)
		passed = ec.ExitCode() == 0
		return tempFilePath, passed, nil // Return path, pass status, nil error (command ran)
	} else if err != nil {
		// Some other error running the command (e.g., command not found)
		return "", false, fmt.Errorf("error running test command %v: %w", testCmd, err)
	}

	// Otherwise, command exited 0 => passed
	passed = true
	return tempFilePath, passed, nil
}
