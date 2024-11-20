package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedTestResult struct {
	TestResult *reports.TestResult
	seen       bool
}

const (
	defaultRuns          = 5
	flakyTestPackagePath = "./example_test_package"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		runner   Runner
		expected map[string]*expectedTestResult
	}{
		{
			name: "default",
			runner: Runner{
				ProjectPath:          "./",
				Verbose:              true,
				RunCount:             defaultRuns,
				UseRace:              false,
				SkipTests:            []string{"TestPanic"},
				FailFast:             false,
				SelectedTestPackages: []string{flakyTestPackagePath},
				CollectRawOutput:     true,
			},
			expected: map[string]*expectedTestResult{
				"TestFlaky": {
					TestResult: &reports.TestResult{
						TestName: "TestFlaky",
						Panicked: false,
						Skipped:  false,
					},
				},
				"TestFail": {
					TestResult: &reports.TestResult{
						TestName:  "TestFail",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 0,
						Failures:  defaultRuns,
					},
				},
				"TestPass": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
						Successes: defaultRuns,
					},
				},
				"TestSkipped": {
					TestResult: &reports.TestResult{
						TestName:  "TestSkipped",
						Panicked:  false,
						Skipped:   true,
						PassRatio: 0,
					},
				},
				"TestRace": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
						Successes: defaultRuns,
					},
				},
			},
		},
		{
			name: "panic",
			runner: Runner{
				ProjectPath:          "./",
				Verbose:              true,
				RunCount:             defaultRuns,
				UseRace:              false,
				SkipTests:            []string{},
				FailFast:             false,
				SelectedTestPackages: []string{flakyTestPackagePath},
				CollectRawOutput:     true,
			},
			expected: map[string]*expectedTestResult{
				"TestFlaky": {
					TestResult: &reports.TestResult{
						TestName: "TestFlaky",
						Panicked: true,
						Skipped:  false,
					},
				},
				"TestFail": {
					TestResult: &reports.TestResult{
						TestName:  "TestFail",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 0,
						Failures:  defaultRuns,
					},
				},
				"TestPass": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  true,
						Skipped:   false,
						PassRatio: 1,
						Successes: defaultRuns,
					},
				},
				"TestSkipped": {
					TestResult: &reports.TestResult{
						TestName:  "TestSkipped",
						Panicked:  true,
						Skipped:   true,
						PassRatio: 0,
					},
				},
				"TestPanic": {
					TestResult: &reports.TestResult{
						TestName: "TestPanic",
						Panicked: true,
						Skipped:  false,
					},
				},
				"TestRace": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
						Successes: defaultRuns,
					},
				},
			},
		},
		{
			name: "race",
			runner: Runner{
				ProjectPath:          "./",
				Verbose:              true,
				RunCount:             defaultRuns,
				UseRace:              true,
				SkipTests:            []string{"TestPanic"},
				FailFast:             false,
				SelectedTestPackages: []string{flakyTestPackagePath},
				CollectRawOutput:     true,
			},
			expected: map[string]*expectedTestResult{
				"TestFlaky": {
					TestResult: &reports.TestResult{
						TestName: "TestFlaky",
						Panicked: false,
						Skipped:  false,
					},
				},
				"TestFail": {
					TestResult: &reports.TestResult{
						TestName:  "TestFail",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 0,
						Failures:  defaultRuns,
					},
				},
				"TestPass": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
						Successes: defaultRuns,
					},
				},
				"TestSkipped": {
					TestResult: &reports.TestResult{
						TestName:  "TestSkipped",
						Panicked:  false,
						Skipped:   true,
						PassRatio: 0,
					},
				},
				"TestRace": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 0,
						Failures:  defaultRuns,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testResults, err := tc.runner.RunTests()
			require.NoError(t, err)
			cleanup(t, testResults, tc.runner)

			for _, result := range testResults {
				t.Run(fmt.Sprintf("checking results of %s", result.TestName), func(t *testing.T) {
					expected, ok := tc.expected[result.TestName]
					require.True(t, ok, "unexpected test result: %s", result.TestName)
					require.False(t, expected.seen, "test '%s' was seen multiple times", result.TestName)
					expected.seen = true

					assert.Equal(t, defaultRuns, result.Runs, "test '%s' had an unexpected number of runs", result.TestName)
					assert.Len(t, result.Durations, result.Runs, "test '%s' has a mismatch of runs and duration counts", result.TestName, defaultRuns)
					resultCounts := result.Successes + result.Failures + result.Panics + result.Skips
					assert.Equal(t, result.Runs, resultCounts,
						"test '%s' doesn't match Runs count with results counts\nRuns: %d\nSuccesses: %d\nFailures: %d\nPanics: %d\nSkips: %d\nTotal: %d",
						result.TestName, result.Runs, result.Successes, result.Failures, result.Panics, result.Skips, resultCounts,
					)
					assert.Equal(t, expected.TestResult.Panicked, result.Panicked, "test '%s' had an unexpected panic result", result.TestName)
					assert.Equal(t, expected.TestResult.Skipped, result.Skipped, "test '%s' had an unexpected skipped result", result.TestName)

					if result.TestName == "TestFlaky" {
						assert.Greater(t, result.Successes, 0, "flaky test '%s' should have passed some", result.TestName)
						assert.Greater(t, result.Failures, 0, "flaky test '%s' should have failed some", result.TestName)
						assert.Greater(t, result.PassRatio, float64(0), "flaky test '%s' should have a flaky pass ratio", result.TestName)
						assert.Less(t, result.PassRatio, float64(1), "flaky test '%s' should have a flaky pass ratio", result.TestName)
					} else {
						assert.Equal(t, expected.TestResult.PassRatio, result.PassRatio, "test '%s' had an unexpected pass ratio", result.TestName)
						assert.Equal(t, expected.TestResult.Successes, result.Successes, "test '%s' had an unexpected number of successes", result.TestName)
						assert.Equal(t, expected.TestResult.Failures, result.Failures, "test '%s' had an unexpected number of failures", result.TestName)
					}
				})
			}

			for _, expected := range tc.expected {
				assert.True(t, expected.seen, "expected test '%s' not found in test runs", expected.TestResult.TestName)
			}
		})
	}
}

func cleanup(t *testing.T, testResults []reports.TestResult, runner Runner) {
	t.Helper()

	t.Cleanup(func() {
		if t.Failed() {
			saniTName := strings.ReplaceAll(t.Name(), "/", "_")
			resultsFileName := fmt.Sprintf("debug_test_results_%s.json", saniTName)
			jsonResults, err := json.Marshal(testResults)
			if err != nil {
				t.Logf("error marshalling test results: %v", err)
			}
			err = os.WriteFile(resultsFileName, jsonResults, 0644) //nolint:gosec
			if err != nil {
				t.Logf("error writing test results: %v", err)
			}
			for packageName, rawOutput := range runner.RawOutputs() {
				saniPackageName := filepath.Base(packageName)
				rawOutputFileName := fmt.Sprintf("debug_raw_output_%s_%s.json", saniTName, saniPackageName)
				err = os.WriteFile(rawOutputFileName, rawOutput.Bytes(), 0644) //nolint:gosec
				if err != nil {
					t.Logf("error writing raw output: %v", err)
				}
			}
		}
	})
}
