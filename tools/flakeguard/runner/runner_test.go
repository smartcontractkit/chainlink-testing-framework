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

func TestRun(t *testing.T) {
	t.Parallel()

	runs := 5

	runner := Runner{
		ProjectPath:          "./",
		Verbose:              true,
		RunCount:             runs,
		UseRace:              false,
		SkipTests:            []string{},
		FailFast:             false,
		SelectedTestPackages: []string{"./flaky_test_package"},
	}

	expectedResults := map[string]*struct {
		*reports.TestResult
		seen bool
	}{
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
				Failures:  runs,
			},
		},
		"TestPass": {
			TestResult: &reports.TestResult{
				TestName:  "TestPass",
				Panicked:  false,
				Skipped:   false,
				PassRatio: 1,
				Successes: runs,
			},
		},
		// "TestPanic": {
		// 	TestResult: &reports.TestResult{
		// 		TestName:  "TestPanic",
		// 		Panicked:  true,
		// 		Skipped:   false,
		// 		PassRatio: 0,
		// 	},
		// },
		"TestSkipped": {
			TestResult: &reports.TestResult{
				TestName:  "TestSkipped",
				Panicked:  false,
				Skipped:   true,
				PassRatio: 0,
			},
		},
	}

	testResults, err := runner.RunTests()
	require.NoError(t, err)
	t.Cleanup(func() {
		if t.Failed() {
			t.Log("Writing test results to flaky_test_results.json")
			jsonResults, err := json.Marshal(testResults)
			require.NoError(t, err)
			err = os.WriteFile("flaky_test_results.json", jsonResults, 0644) //nolint:gosec
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

			assert.Equal(t, runs, result.Runs, "test '%s' had an unexpected number of runs", result.TestName)
			assert.Len(t, result.Durations, runs, "test '%s' had an unexpected number of durations as it was run %d times", result.TestName, runs)
			if result.TestName == "TestSlow" {
				for _, duration := range result.Durations {
					assert.GreaterOrEqual(t, duration, float64(1), "slow test '%s' should have a duration of at least 2s", result.TestName)
				}
			}
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
