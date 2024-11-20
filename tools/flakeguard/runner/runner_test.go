package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockGoRunner struct {
	RunFunc func(dir string, args []string) (string, bool, error)
}

func (m MockGoRunner) RunCommand(dir string, args []string) (string, bool, error) {
	return m.RunFunc(dir, args)
}

type expectedTestResult struct {
	TestResult *reports.TestResult
	seen       bool
}

var (
	defaultRuns       = 5
	defaultTestRunner = Runner{
		ProjectPath:          "./",
		Verbose:              true,
		RunCount:             defaultRuns,
		UseRace:              false,
		SkipTests:            []string{"TestPanic"},
		FailFast:             false,
		SelectedTestPackages: []string{"./flaky_test_package"},
		CollectRawOutput:     true,
	}
)

func TestRunDefault(t *testing.T) {
	expectedResults := map[string]*expectedTestResult{
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
	}

	testResults, err := defaultTestRunner.RunTests()
	require.NoError(t, err)
	t.Cleanup(func() {
		if t.Failed() {
			resultsFileName := fmt.Sprintf("flaky_test_results_%s.json", t.Name())
			t.Logf("Writing test results to %s", resultsFileName)
			jsonResults, err := json.Marshal(testResults)
			require.NoError(t, err)
			err = os.WriteFile(resultsFileName, jsonResults, 0644) //nolint:gosec
			require.NoError(t, err)
		}
	})

	for _, result := range testResults {
		t.Run(fmt.Sprintf("checking results of %s", result.TestName), func(t *testing.T) {
			expected, ok := expectedResults[result.TestName]
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

	for _, expected := range expectedResults {
		assert.True(t, expected.seen, "expected test '%s' not found in test runs", expected.TestResult.TestName)
	}
}

func TestRunWithPanics(t *testing.T) {
	panicRunner := defaultTestRunner
	panicRunner.SkipTests = []string{}

	expectedResults := map[string]*expectedTestResult{
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
	}

	testResults, err := panicRunner.RunTests()
	require.NoError(t, err)
	t.Cleanup(func() {
		if t.Failed() {
			resultsFileName := fmt.Sprintf("flaky_test_results_%s.json", t.Name())
			t.Logf("Writing test results to %s", resultsFileName)
			jsonResults, err := json.Marshal(testResults)
			require.NoError(t, err)
			err = os.WriteFile(resultsFileName, jsonResults, 0644) //nolint:gosec
			require.NoError(t, err)
		}
	})

	for _, result := range testResults {
		t.Run(fmt.Sprintf("checking results of %s", result.TestName), func(t *testing.T) {
			expected, ok := expectedResults[result.TestName]
			// Sanity checks
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

	for _, expected := range expectedResults {
		assert.True(t, expected.seen, "expected test '%s' not found in test runs", expected.TestResult.TestName)
	}
}
