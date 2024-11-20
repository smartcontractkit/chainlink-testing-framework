package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
)

var (
	panicRe = regexp.MustCompile(`^panic:`)
)

type Runner struct {
	ProjectPath          string   // Path to the Go project directory.
	Verbose              bool     // If true, provides detailed logging.
	RunCount             int      // Number of times to run the tests.
	UseRace              bool     // Enable race detector.
	FailFast             bool     // Stop on first test failure.
	SkipTests            []string // Test names to exclude.
	SelectedTestPackages []string // Explicitly selected packages to run.
	CollectRawOutput     bool     // Collect test output for later inspection.
	rawOutput            bytes.Buffer
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

// RawOutput retrieves the raw output from the test runs, if CollectRawOutput enabled.
func (r *Runner) RawOutput() bytes.Buffer {
	return r.rawOutput
}

type exitCoder interface {
	ExitCode() int
}

// runTests runs the tests for a given package and returns the path to the output file.
func (r *Runner) runTests(packageName string) (string, bool, error) {
	args := []string{"test", packageName, "-json", "-count=1"}
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
	if r.CollectRawOutput {
		cmd.Stdout = io.MultiWriter(tmpFile, &r.rawOutput)
		cmd.Stderr = io.MultiWriter(tmpFile, &r.rawOutput)
	} else {
		cmd.Stdout = tmpFile
		cmd.Stderr = tmpFile
	}
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

		scanner := bufio.NewScanner(file)
		var precedingLines []string // Store preceding lines for context
		var followingLines []string // To collect lines after an error

		for scanner.Scan() {
			line := scanner.Text()
			precedingLines = append(precedingLines, line)

			// Limit precedingLines to the last 15 lines
			if len(precedingLines) > 15 {
				precedingLines = precedingLines[1:]
			}

			var entry struct {
				Action  string  `json:"Action"`
				Test    string  `json:"Test"`
				Package string  `json:"Package"`
				Output  string  `json:"Output"`
				Elapsed float64 `json:"Elapsed"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
				// Collect 15 lines after the error for more context
				for scanner.Scan() && len(followingLines) < 15 {
					followingLines = append(followingLines, scanner.Text())
				}

				// Combine precedingLines and followingLines to provide 15 lines before and after
				context := append(precedingLines, followingLines...)
				return nil, fmt.Errorf("failed to parse json test output near lines:\n%s\nerror: %w", strings.Join(context, "\n"), err)
			}

			// Only create TestResult for test-level entries
			var result *reports.TestResult
			if entry.Test != "" {
				// Determine the key
				key := entry.Package + "/" + entry.Test

				if _, exists := testDetails[key]; !exists {
					testDetails[key] = &reports.TestResult{
						TestName:       entry.Test,
						TestPackage:    entry.Package,
						PassRatio:      0,
						Outputs:        []string{},
						PackageOutputs: []string{},
					}
				}
				result = testDetails[key]
			}

			// Collect outputs
			if entry.Output != "" {
				if entry.Test != "" {
					// Test-level output
					result.Outputs = append(result.Outputs, strings.TrimSpace(entry.Output))
				} else {
					// Package-level output
					// Append to PackageOutputs of all TestResults in the same package
					for _, res := range testDetails {
						if res.TestPackage == entry.Package {
							res.PackageOutputs = append(res.PackageOutputs, entry.Output)
						}
					}
				}
			}

			switch entry.Action {
			case "run":
				if entry.Test != "" {
					result.Runs++
				}
			case "pass":
				if entry.Test != "" {
					result.Durations = append(result.Durations, entry.Elapsed)
					result.Successes++
				}
			case "fail":
				if entry.Test != "" {
					result.Durations = append(result.Durations, entry.Elapsed)
					result.Failures++
				}
			case "output":
				// Output already handled above
				if panicRe.MatchString(entry.Output) {
					if entry.Test != "" {
						// Test-level panic
						result.Panicked = true
						result.Durations = append(result.Durations, entry.Elapsed)
						result.Panics++
					} else {
						// Package-level panic
						// Mark PackagePanicked for all TestResults in the package
						for _, res := range testDetails {
							if res.TestPackage == entry.Package {
								res.PackagePanicked = true
							}
						}
					}
				}
			case "skip":
				if entry.Test != "" {
					result.Skipped = true
					result.Skips++
					result.Durations = append(result.Durations, entry.Elapsed)
				}
			}
			if entry.Test != "" {
				result.PassRatio = float64(result.Successes) / float64(result.Runs)
				result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading test output file: %w", err)
		}
		// Clean up file after parsing
		if err = file.Close(); err != nil {
			log.Printf("WARN: failed to close file: %v", err)
		}
		if err = os.Remove(filePath); err != nil {
			log.Printf("WARN: failed to delete file: %v", err)
		}
	}

	var results []reports.TestResult
	for _, result := range testDetails {
		results = append(results, *result)
	}

	return results, nil
}
