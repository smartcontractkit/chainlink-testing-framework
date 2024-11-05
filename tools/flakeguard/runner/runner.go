package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	RunAllTestPackages   bool     // Run all test packages if true.
	SelectedTestPackages []string // Explicitly selected packages to run.
}

// RunTests executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTests() ([]reports.TestResult, error) {
	var jsonOutputs [][]byte
	packages := r.SelectedTestPackages
	if r.RunAllTestPackages {
		packages = []string{"./..."}
	}

	for _, p := range packages {
		for i := 0; i < r.RunCount; i++ {
			jsonOutput, passed, err := r.runTests(p)
			if err != nil {
				return nil, fmt.Errorf("failed to run tests in package %s: %w", p, err)
			}
			jsonOutputs = append(jsonOutputs, jsonOutput)
			if !passed && r.FailFast {
				break
			}
		}
	}

	return parseTestResults(jsonOutputs)
}

type exitCoder interface {
	ExitCode() int
}

// runTestPackage executes the test command for a single test package.
// It returns the command output, a boolean indicating success, and any error encountered.
func (r *Runner) runTests(packageName string) ([]byte, bool, error) {
	args := []string{"test", packageName, "-json", "-count=1"} // Enable JSON output for parsing
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
	cmd := exec.Command("go", args...)
	cmd.Dir = r.ProjectPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		var exErr exitCoder
		// Check if the error is due to a non-zero exit code
		if errors.As(err, &exErr) && exErr.ExitCode() == 0 {
			return nil, false, fmt.Errorf("test command failed at %s: %w", packageName, err)
		}
		return out.Bytes(), false, nil // Test failed
	}

	return out.Bytes(), true, nil // Test succeeded
}

// parseTestResults analyzes multiple JSON outputs from 'go test -json' commands to determine test results.
// It accepts a slice of []byte where each []byte represents a separate JSON output from a test run.
// This function aggregates results across multiple test runs, summing runs and passes for each test.
func parseTestResults(datas [][]byte) ([]reports.TestResult, error) {
	testDetails := make(map[string]*reports.TestResult) // Holds run, pass counts, and other details for each test

	// Process each data set
	for _, data := range datas {
		scanner := bufio.NewScanner(bytes.NewReader(data))
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
				result.Durations = append(result.Durations, entry.Elapsed)
			case "output":
				result.Outputs = append(result.Outputs, entry.Output)
			case "fail":
				result.PassRatio = (result.PassRatio * float64(result.Runs-1)) / float64(result.Runs)
				result.Durations = append(result.Durations, entry.Elapsed)
			case "skip":
				result.Skipped = true
				result.Runs++
				result.Durations = append(result.Durations, entry.Elapsed)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading standard input: %w", err)
		}
	}

	var results []reports.TestResult
	for _, result := range testDetails {
		results = append(results, *result)
	}

	return results, nil
}
