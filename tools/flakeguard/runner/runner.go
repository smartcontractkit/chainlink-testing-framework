package runner

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
)

type Runner struct {
	ProjectPath          string   // Path to the Go project directory.
	Verbose              bool     // If true, provides detailed logging.
	RunCount             int      // Number of times to run the tests.
	UseRace              bool     // Enable race detector.
	FailFast             bool     // Stop on first test failure.
	SkipTests            []string // Test names to exclude.
	SelectedTestPackages []string // Explicitly selected packages to run.
}

// RunTests executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTests() ([]reports.TestResult, error) {
	var jsonFilePaths []string
	for _, p := range r.SelectedTestPackages {
		for i := 0; i < r.RunCount; i++ {
			jsonFilePath, passed, err := r.runTests(p)
			if err != nil {
				return nil, fmt.Errorf("failed to run tests in package %s: %w", p, err)
			}
			jsonFilePaths = append(jsonFilePaths, jsonFilePath)
			if !passed && r.FailFast {
				break
			}
		}
	}

	return parseTestResults(jsonFilePaths)
}

type exitCoder interface {
	ExitCode() int
}

// runTests runs the tests for a given package and returns the path to the output file.
func (r *Runner) runTests(packageName string) (string, bool, error) {
	args := []string{"test", packageName, "-json", "-count=1"} // Enable JSON output
	if r.UseRace {
		args = append(args, "-race")
	}
	if len(r.SkipTests) > 0 {
		skipPattern := strings.Join(r.SkipTests, "|")
		args = append(args, fmt.Sprintf("-skip=%s", skipPattern))
	}

	if r.Verbose {
		log.Printf("Running command: go %s\n", strings.Join(args, " "))
	}

	// Create a temporary file to store the output
	tmpFile, err := os.CreateTemp("", "test-output-*.json")
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Run the command with output directed to the file
	cmd := exec.Command("go", args...)
	cmd.Dir = r.ProjectPath
	cmd.Stdout = tmpFile
	cmd.Stderr = tmpFile

	err = cmd.Run()
	if err != nil {
		var exErr exitCoder
		// Check if the error is due to a non-zero exit code
		if errors.As(err, &exErr) && exErr.ExitCode() == 0 {
			return "", false, fmt.Errorf("test command failed at %s: %w", packageName, err)
		}
		return tmpFile.Name(), false, nil // Test failed
	}

	return tmpFile.Name(), true, nil // Test succeeded
}

// parseTestResults reads the test output files and returns the parsed test results.
func parseTestResults(filePaths []string) ([]reports.TestResult, error) {
	testDetails := make(map[string]*reports.TestResult) // Holds run, pass counts, and other details for each test

	// Process each file
	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open test output file: %w", err)
		}
		defer os.Remove(filePath) // Clean up file after parsing
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var entry struct {
				Action  string  `json:"Action"`
				Test    string  `json:"Test"`
				Package string  `json:"Package"`
				Output  string  `json:"Output"`
				Elapsed float64 `json:"Elapsed"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
				return nil, fmt.Errorf("failed to parse json test output: %s, err: %w", scanner.Text(), err)
			}

			// Skip processing if the test name is empty
			if entry.Test == "" {
				continue
			}

			key := entry.Package + "/" + entry.Test // Create a unique key using package and test name
			if _, exists := testDetails[key]; !exists {
				testDetails[key] = &reports.TestResult{
					TestName:    entry.Test,
					TestPackage: entry.Package,
					Runs:        0,
					PassRatio:   0,
					Outputs:     []string{},
				}
			}

			result := testDetails[key]
			switch entry.Action {
			case "run":
				result.Runs++
			case "pass":
				result.PassRatio = (result.PassRatio*float64(result.Runs-1) + 1) / float64(result.Runs)
				result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
				result.Durations = append(result.Durations, entry.Elapsed)
			case "output":
				result.Outputs = append(result.Outputs, entry.Output)
			case "fail":
				result.PassRatio = (result.PassRatio * float64(result.Runs-1)) / float64(result.Runs)
				result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
				result.Durations = append(result.Durations, entry.Elapsed)
			case "skip":
				result.Skipped = true
				result.Runs++
				result.Durations = append(result.Durations, entry.Elapsed)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading test output file: %w", err)
		}
	}

	var results []reports.TestResult
	for _, result := range testDetails {
		results = append(results, *result)
	}

	return results, nil
}
