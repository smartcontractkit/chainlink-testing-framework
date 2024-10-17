package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Runner struct {
	Verbose  bool   // If true, provides detailed logging.
	Dir      string // Directory to run commands in.
	Count    int    // Number of times to run the tests.
	UseRace  bool   // Enable race detector.
	FailFast bool   // Stop on first test failure.
}

// RunTests executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTests(packages []string) ([]TestResult, error) {
	var allResults []TestResult
	var errors []string

	for _, p := range packages {
		testResults, err := r.runTestPackage(p)
		allResults = append(allResults, testResults...)

		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to run tests in package %s: %v", p, err))
		}
	}

	if len(errors) > 0 {
		return allResults, fmt.Errorf("some tests failed: %s", strings.Join(errors, "; "))
	}

	return allResults, nil
}

type TestResult struct {
	TestName  string
	PassRatio float64
	Runs      int
	Passed    bool
}

// runTestPackage executes the test command for a single test package.
func (r *Runner) runTestPackage(testPackage string) ([]TestResult, error) {
	args := []string{"test", "-json"} // Enable JSON output
	if r.Count > 0 {
		args = append(args, "-count", fmt.Sprint(r.Count))
	}
	if r.UseRace {
		args = append(args, "-race")
	}
	if r.FailFast {
		args = append(args, "-failfast")
	}
	args = append(args, testPackage)

	// Construct the command
	cmd := exec.Command("go", args...)
	// cmd.Env = append(cmd.Env, "GOFLAGS=-extldflags=-Wl,-ld_classic") // Ensure modules are enabled
	cmd.Dir = r.Dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()

	// Parse results
	results, parseErr := parseTestResults(out.Bytes())
	if parseErr != nil {
		return results, fmt.Errorf("failed to parse test results for %s: %v", testPackage, parseErr)
	}

	if err != nil {
		return results, fmt.Errorf("test command failed at %s: %w", testPackage, err)
	}

	return results, nil
}

// parseTestResults analyzes the JSON output from 'go test -json' to determine test results
func parseTestResults(data []byte) ([]TestResult, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	testDetails := make(map[string]map[string]int) // Holds run and pass counts for each test

	for scanner.Scan() {
		var entry struct {
			Action string `json:"Action"`
			Test   string `json:"Test"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, fmt.Errorf("failed to parse json test output: %w", err)
		}

		// Skip processing if the test name is empty
		if entry.Test == "" {
			continue
		}

		if _, exists := testDetails[entry.Test]; !exists {
			testDetails[entry.Test] = map[string]int{"run": 0, "pass": 0}
		}

		if entry.Action == "run" {
			testDetails[entry.Test]["run"]++
		}
		if entry.Action == "pass" {
			testDetails[entry.Test]["pass"]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading standard input: %w", err)
	}

	var results []TestResult
	for testName, counts := range testDetails {
		runs := counts["run"]
		passes := counts["pass"]
		passRatio := 0.0
		if runs > 0 {
			passRatio = float64(passes) * 100 / float64(runs)
		}
		results = append(results, TestResult{
			TestName:  testName,
			PassRatio: passRatio,
			Passed:    passes == runs,
			Runs:      runs,
		})
	}

	return results, nil
}
