package runner

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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
	Count                int      // -test.count flag value.
	UseShuffle           bool     // -test.shuffle flag value.
	JsonOutput           bool     // -test.json flag value.
	FailFast             bool     // Stop on first test failure.
	SkipTests            []string // Test names to exclude.
	SelectedTestPackages []string // Explicitly selected packages to run.
}

// RunTests executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTests() ([]reports.TestResult, error) {
	var jsonFilePaths []string
	for _, p := range r.SelectedTestPackages {
		binaryPath, err := r.buildTestBinary(p)
		if err != nil {
			return nil, fmt.Errorf("failed to build test binary for package %s: %w", p, err)
		}
		defer os.Remove(binaryPath) // Clean up test binary after running tests
		for i := 0; i < r.RunCount; i++ {
			jsonFilePath, passed, err := r.runTestBinary(binaryPath)
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

// buildTestBinary builds the test binary for a given package and returns the path to the binary.
func (r *Runner) buildTestBinary(packageName string) (string, error) {
	binaryPath := fmt.Sprintf("%s/test-binary-%s", os.TempDir(), strings.ReplaceAll(packageName, "/", "-"))
	// Compile-time flags for building the binary
	args := []string{"test", "-c", packageName, "-o", binaryPath}
	if r.UseRace {
		args = append(args, "-race")
	}

	if r.Verbose {
		log.Printf("Building test binary with command: go %s\n", strings.Join(args, " "))
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = r.ProjectPath
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build test binary for %s: %w", packageName, err)
	}

	return binaryPath, nil
}

// runTestBinary runs the tests for a given package and returns the path to the output file.
func (r *Runner) runTestBinary(binaryPath string) (string, bool, error) {
	// Runtime flags for executing the binary
	args := []string{binaryPath}
	if r.UseShuffle {
		args = append(args, "-test.shuffle=on")
	}
	if r.JsonOutput {
		args = append(args, "-test.json")
	}
	if r.Count > 0 {
		args = append(args, fmt.Sprintf("-test.count=%d", r.Count))
	}
	if len(r.SkipTests) > 0 {
		skipPattern := strings.Join(r.SkipTests, "|")
		args = append(args, fmt.Sprintf("-test.skip=%s", skipPattern))
	}

	if r.Verbose {
		log.Printf("Running command: %s\n", strings.Join(args, " "))
	}

	// Create a temporary file to store the output
	tmpFile, err := os.CreateTemp("", "test-output-*.json")
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = r.ProjectPath
	cmd.Stdout = tmpFile
	cmd.Stderr = tmpFile

	err = cmd.Run()
	if err != nil {
		var exErr exitCoder
		// Check if the error is due to a non-zero exit code
		if errors.As(err, &exErr) && exErr.ExitCode() == 0 {
			return "", false, fmt.Errorf("test command failed for binary at %s: %w", binaryPath, err)
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

			// Only create TestResult for test-level entries
			var result *reports.TestResult
			if entry.Test != "" {
				// Determine the key
				key := entry.Package + "/" + entry.Test

				if _, exists := testDetails[key]; !exists {
					testDetails[key] = &reports.TestResult{
						TestName:       entry.Test,
						TestPackage:    entry.Package,
						Runs:           0,
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
					result.Outputs = append(result.Outputs, entry.Output)
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
					result.PassRatio = (result.PassRatio*float64(result.Runs-1) + 1) / float64(result.Runs)
					result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
					result.Durations = append(result.Durations, entry.Elapsed)
				}
			case "fail":
				if entry.Test != "" {
					result.PassRatio = (result.PassRatio * float64(result.Runs-1)) / float64(result.Runs)
					result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
					result.Durations = append(result.Durations, entry.Elapsed)
				}
			case "output":
				// Output already handled above
				if panicRe.MatchString(entry.Output) {
					if entry.Test != "" {
						// Test-level panic
						result.Panicked = true
						result.PassRatio = (result.PassRatio * float64(result.Runs-1)) / float64(result.Runs)
						result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
						result.Durations = append(result.Durations, entry.Elapsed)
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
					result.Runs++
					result.Durations = append(result.Durations, entry.Elapsed)
				}
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
