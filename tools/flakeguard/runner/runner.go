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
	Verbose  bool // If true, provides detailed logging.
	RunCount int  // Number of times to run the tests.
	UseRace  bool // Enable race detector.
	FailFast bool // Stop on first test failure.
}

// RunTests executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTests(packages []string) ([]reports.TestResult, error) {
	var jsonOutputs [][]byte

	for _, p := range packages {
		for i := 0; i < r.RunCount; i++ {
			jsonOutput, err := r.runTestPackage(p)
			if err != nil {
				return nil, fmt.Errorf("failed to run tests in package %s: %w", p, err)
			}
			jsonOutputs = append(jsonOutputs, jsonOutput)
		}
	}

	// TODO: make sure the test name includes package name
	return parseTestResults(jsonOutputs)
}

type exitCoder interface {
	ExitCode() int
}

// runTestPackage executes the test command for a single test package.
func (r *Runner) runTestPackage(testPackage string) ([]byte, error) {
	args := []string{"test", "-json", "-count=1"} // Enable JSON output
	// if r.Count > 0 {
	// 	args = append(args, fmt.Sprintf("-count=%d", r.Count))
	// }
	if r.UseRace {
		args = append(args, "-race")
	}
	if r.FailFast {
		args = append(args, "-failfast")
	}
	args = append(args, testPackage)

	if r.Verbose {
		log.Printf("Running command: go %s\n", strings.Join(args, " "))
	}
	cmd := exec.Command("go", args...)

	// cmd.Env = append(cmd.Env, "GOFLAGS=-extldflags=-Wl,-ld_classic") // Ensure modules are enabled
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		var exErr exitCoder
		if errors.As(err, &exErr) && exErr.ExitCode() == 0 {
			return nil, fmt.Errorf("test command failed at %s: %w", testPackage, err)
		}
	}

	return out.Bytes(), nil
}

// parseTestResults analyzes multiple JSON outputs from 'go test -json' commands to determine test results.
// It accepts a slice of []byte where each []byte represents a separate JSON output from a test run.
// This function aggregates results across multiple test runs, summing runs and passes for each test.
func parseTestResults(datas [][]byte) ([]reports.TestResult, error) {
	testDetails := make(map[string]map[string]interface{}) // Holds run, pass counts, and other details for each test

	// Process each data set
	for _, data := range datas {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			var entry struct {
				Action  string `json:"Action"`
				Test    string `json:"Test"`
				Package string `json:"Package"`
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
				testDetails[key] = map[string]interface{}{
					"run":     0,
					"pass":    0,
					"name":    entry.Test,
					"package": entry.Package,
				}
			}

			details := testDetails[key]
			if entry.Action == "run" {
				details["run"] = details["run"].(int) + 1
			}
			if entry.Action == "pass" {
				details["pass"] = details["pass"].(int) + 1
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading standard input: %w", err)
		}
	}

	var results []reports.TestResult
	for _, details := range testDetails {
		runs := details["run"].(int)
		passes := details["pass"].(int)
		passRatio := 0.0
		if runs > 0 {
			passRatio = float64(passes) / float64(runs)
		}
		results = append(results, reports.TestResult{
			TestName:    details["name"].(string),
			TestPackage: details["package"].(string),
			PassRatio:   passRatio,
			Runs:        runs,
		})
	}

	return results, nil
}
