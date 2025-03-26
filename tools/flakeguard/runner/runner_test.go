package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
)

var (
	defaultTestRunCount  = 5
	flakyTestPackagePath = "./example_test_package"
	debugDir             = "debug_outputs"
)

type expectedTestResult struct {
	allSuccesses  bool
	someSuccesses bool
	allFailures   bool
	someFailures  bool
	allSkips      bool
	testPanic     bool
	packagePanic  bool
	race          bool
	maximumRuns   int

	exactRuns       *int
	minimumRuns     *int
	exactPassRate   *float64
	minimumPassRate *float64
	maximumPassRate *float64

	seen bool
}

func TestPrettyProjectPath(t *testing.T) {
	t.Parallel()

	prettyPath, err := utils.GetGoProjectName("./")
	require.NoError(t, err)
	assert.Equal(t, "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard", prettyPath)
}

func TestRun(t *testing.T) {
	var (
		zeroRuns        = 0
		oneCount        = 1
		successPassRate = 1.0
		failPassRate    = 0.0
	)
	testCases := []struct {
		name          string
		runner        Runner
		expectedTests map[string]*expectedTestResult
	}{
		{
			name: "default",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         defaultTestRunCount,
				GoTestRaceFlag:   false,
				SkipTests:        []string{"TestPanic", "TestFlakyPanic", "TestSubTestsSomePanic", "TestTimeout"},
				FailFast:         false,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				"TestFlaky": {
					exactRuns:       &defaultTestRunCount,
					minimumPassRate: &failPassRate,
					maximumPassRate: &successPassRate,
					someSuccesses:   true,
					someFailures:    true,
				},
				"TestFail": {
					exactRuns:     &defaultTestRunCount,
					exactPassRate: &failPassRate,
					allFailures:   true,
				},
				"TestFailLargeOutput": {
					exactRuns:     &defaultTestRunCount,
					exactPassRate: &failPassRate,
					allFailures:   true,
				},
				"TestPass": {
					exactRuns:     &defaultTestRunCount,
					exactPassRate: &successPassRate,
					allSuccesses:  true,
				},
				"TestSkipped": {
					exactRuns:     &zeroRuns,
					exactPassRate: &successPassRate,
					allSkips:      true,
				},
				"TestRace": {
					exactRuns:     &defaultTestRunCount,
					exactPassRate: &successPassRate,
					allSuccesses:  true,
				},
				"TestSubTestsAllPass": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
				"TestFailInParentAfterSubTests": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestFailInParentAfterSubTests/Pass1": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
				"TestFailInParentAfterSubTests/Pass2": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
				"TestFailInParentBeforeSubTests": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestSubTestsAllPass/Pass1": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
				"TestSubTestsAllPass/Pass2": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
				"TestSubTestsAllFail": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestSubTestsAllFail/Fail1": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestSubTestsAllFail/Fail2": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestSubTestsSomeFail": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestSubTestsSomeFail/Pass": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
				"TestSubTestsSomeFail/Fail": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
			},
		},
		{
			name: "always panic",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         defaultTestRunCount,
				GoTestRaceFlag:   false,
				SkipTests:        []string{},
				SelectTests:      []string{"TestPanic"},
				FailFast:         false,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				"TestPanic": {
					packagePanic: true,
					testPanic:    true,
					maximumRuns:  defaultTestRunCount,
				},
			},
		},
		{
			name: "flaky panic",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         defaultTestRunCount,
				GoTestRaceFlag:   false,
				GoTestCountFlag:  &oneCount,
				SkipTests:        []string{},
				SelectTests:      []string{"TestFlakyPanic"},
				FailFast:         false,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				"TestFlakyPanic": {
					packagePanic: true,
					testPanic:    true,
					maximumRuns:  defaultTestRunCount,
				},
			},
		},
		{
			name: "subtest panic",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         defaultTestRunCount,
				GoTestRaceFlag:   false,
				SkipTests:        []string{},
				SelectTests:      []string{"TestSubTestsSomePanic"},
				FailFast:         false,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				"TestSubTestsSomePanic": {
					packagePanic: true,
					testPanic:    true,
					maximumRuns:  defaultTestRunCount,
				},
				"TestSubTestsSomePanic/Pass": {
					packagePanic: true,
					allSuccesses: true,
					maximumRuns:  defaultTestRunCount,
				},
				"TestSubTestsSomePanic/Panic": {
					packagePanic: true,
					testPanic:    true,
					maximumRuns:  defaultTestRunCount,
				},
			},
		},
		{
			name: "failfast",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         defaultTestRunCount,
				GoTestRaceFlag:   false,
				SkipTests:        []string{},
				SelectTests:      []string{"TestFail", "TestPass"},
				FailFast:         true,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				"TestFail": {
					exactRuns:   &oneCount,
					allFailures: true,
				},
				"TestPass": {
					exactRuns:    &oneCount,
					allSuccesses: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testResults, err := tc.runner.RunTestPackages([]string{flakyTestPackagePath})
			require.NoError(t, err)

			t.Cleanup(func() {
				if !t.Failed() {
					return
				}
				if err := os.MkdirAll(debugDir, 0755); err != nil {
					t.Logf("error creating directory: %v", err)
					return
				}
				saniTName := strings.ReplaceAll(t.Name(), "/", "_")
				resultsFileName := filepath.Join(debugDir, fmt.Sprintf("test_results_%s.json", saniTName))
				jsonResults, err := json.Marshal(testResults)
				if err != nil {
					t.Logf("error marshalling test report: %v", err)
					return
				}
				err = os.WriteFile(resultsFileName, jsonResults, 0644) //nolint:gosec
				if err != nil {
					t.Logf("error writing test results: %v", err)
					return
				}
				for packageName, rawOutput := range tc.runner.RawOutputs() {
					saniPackageName := filepath.Base(packageName)
					rawJSONOutputFileName := filepath.Join(debugDir, fmt.Sprintf("raw_output_%s_%s.json", saniTName, saniPackageName))
					err = os.WriteFile(rawJSONOutputFileName, rawOutput.Bytes(), 0644) //nolint:gosec
					if err != nil {
						t.Logf("error writing raw JSON output: %v", err)
					}
				}
			})

			assert.Equal(t, len(tc.expectedTests), len(testResults), "unexpected number of test results")
			for _, result := range testResults {
				t.Run(fmt.Sprintf("checking results of %s", result.TestName), func(t *testing.T) {
					require.NotNil(t, result, "test result was nil")
					expected, ok := tc.expectedTests[result.TestName]
					require.True(t, ok, "unexpected test name: %s", result.TestName)
					require.False(t, expected.seen, "test '%s' was seen multiple times", result.TestName)
					expected.seen = true

					if !expected.testPanic { // Panics end up wrecking durations
						assert.Len(t, result.Durations, result.Runs, "test '%s' has a mismatch of runs %d and duration counts %d",
							result.TestName, result.Runs, len(result.Durations),
						)
						assert.False(t, result.Panic, "test '%s' should not have panicked", result.TestName)
					}
					resultCounts := result.Successes + result.Failures
					assert.Equal(t, result.Runs, resultCounts,
						"test '%s' doesn't match Runs count with results counts\n%s", result.TestName, result.Runs, resultsString(result),
					)

					if expected.minimumRuns != nil {
						assert.GreaterOrEqual(t, result.Runs, *expected.minimumRuns, "test '%s' had fewer runs than expected", result.TestName)
					}
					if expected.exactRuns != nil {
						assert.Equal(t, *expected.exactRuns, result.Runs, "test '%s' had an unexpected number of runs", result.TestName)
					} else {
						assert.LessOrEqual(t, result.Runs, expected.maximumRuns, "test '%s' had more runs than expected", result.TestName)
					}
					if expected.exactPassRate != nil {
						assert.Equal(t, *expected.exactPassRate, result.PassRatio, "test '%s' had an unexpected pass ratio", result.TestName)
					}
					if expected.minimumPassRate != nil {
						assert.Greater(t, result.PassRatio, *expected.minimumPassRate, "test '%s' had a pass ratio below the minimum", result.TestName)
					}
					if expected.maximumPassRate != nil {
						assert.Less(t, result.PassRatio, *expected.maximumPassRate, "test '%s' had a pass ratio above the maximum", result.TestName)
					}
					if expected.allSuccesses {
						assert.Equal(t, result.Successes, result.Runs, "test '%s' has %d total runs and should have passed all runs, only passed %d\n%s", result.TestName, result.Runs, result.Successes, resultsString(result))
						assert.Zero(t, result.Failures, "test '%s' has %d total runs and should have passed all runs, but failed some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Skips, "test '%s' has %d total runs and should have passed all runs, but skipped some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.False(t, result.Panic, "test '%s' should not have panicked\n%s", result.TestName, resultsString(result))
						assert.False(t, result.Race, "test '%s' should not have raced\n%s", result.TestName, resultsString(result))
					}
					if expected.someSuccesses {
						assert.Greater(t, result.Successes, 0, "test '%s' has %d total runs and should have passed some runs, passed none\n%s", result.TestName, result.Runs, resultsString(result))
					}
					if expected.allFailures {
						assert.Equal(t, result.Failures, result.Runs, "test '%s' has %d total runs and should have failed all runs, only failed %d\n%s", result.TestName, result.Runs, result.Failures, resultsString(result))
						assert.Zero(t, result.Successes, "test '%s' has %d total runs and should have failed all runs, but succeeded some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Skips, "test '%s' has %d total runs and should have failed all runs, but skipped some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.False(t, result.Race, "test '%s' should not have raced\n%s", result.TestName, resultsString(result))
					}
					if expected.packagePanic {
						assert.True(t, result.PackagePanic, "test '%s' should have package panicked", result.TestName)
					}
					if expected.testPanic {
						assert.True(t, result.Panic, "test '%s' should have panicked", result.TestName)
						assert.True(t, result.PackagePanic, "test '%s' should have package panicked", result.TestName)
						expected.someFailures = true
					}
					if expected.someFailures {
						assert.Greater(t, result.Failures, 0, "test '%s' has %d total runs and should have failed some runs, failed none\n%s", result.TestName, result.Runs, resultsString(result))
					}
					if expected.allSkips {
						assert.Equal(t, 0, result.Runs, "test '%s' has %d total runs and should have skipped all of them, no runs expected\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Successes, "test '%s' has %d total runs and should have skipped all runs, but succeeded some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Failures, "test '%s' has %d total runs and should have skipped all runs, but panicked some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.False(t, result.Panic, "test '%s' should not have panicked\n%s", result.TestName, resultsString(result))
						assert.False(t, result.Race, "test '%s' should not have raced\n%s", result.TestName, resultsString(result))
					}
					if expected.race {
						assert.True(t, result.Race, "test '%s' should have a data race\n%s", result.TestName, resultsString(result))
						assert.False(t, result.Panic, "test '%s' should not have panicked\n%s", result.TestName, resultsString(result))
						assert.Zero(t, result.Successes, "test '%s' has %d total runs and should have raced all runs, but succeeded some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Failures, "test '%s' has %d total runs and should have raced all runs, but panicked some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Skips, "test '%s' has %d total runs and should have raced all runs, but skipped some\n%s", result.TestName, result.Runs, resultsString(result))
						assert.Zero(t, result.Skips, "test '%s' has %d total runs and should have raced all runs, but panicked some\n%s", result.TestName, result.Runs, resultsString(result))
					}
				})
			}

			allTestsRun := []string{}
			for testName, expected := range tc.expectedTests {
				if expected.seen {
					allTestsRun = append(allTestsRun, testName)
				}
			}
			for testName, expected := range tc.expectedTests {
				require.True(t, expected.seen, "expected test '%s' not found in test runs\nAll tests run: %s", testName, strings.Join(allTestsRun, ", "))
			}
		})
	}
}

func resultsString(result reports.TestResult) string {
	resultCounts := result.Successes + result.Failures + result.Skips
	return fmt.Sprintf("Runs: %d\nPanicked: %t\nRace: %t\nSuccesses: %d\nFailures: %d\nSkips: %d\nTotal Results: %d",
		result.Runs, result.Panic, result.Race, result.Successes, result.Failures, result.Skips, resultCounts)
}

// TODO: Running the failing test here fools tools like gotestfmt into thinking we actually ran a failing test
// as the output gets piped out to the console. This a confusing annoyance that I'd like to fix, but it's not crucial.
func TestFailedOutputs(t *testing.T) {
	t.Parallel()

	runner := Runner{
		ProjectPath:      "./",
		Verbose:          true,
		RunCount:         1,
		SelectTests:      []string{"TestFail"}, // This test is known to fail consistently
		CollectRawOutput: true,
	}

	testResults, err := runner.RunTestPackages([]string{flakyTestPackagePath})
	require.NoError(t, err, "running tests should not produce an unexpected error")

	require.Equal(t, 1, len(testResults), "unexpected number of test runs")

	var testFailResult *reports.TestResult
	for i := range testResults {
		if testResults[i].TestName == "TestFail" {
			testFailResult = &testResults[i]
			break
		}
	}
	require.NotNil(t, testFailResult, "expected TestFail result not found in report")

	require.NotEmpty(t, testFailResult.FailedOutputs, "expected failed outputs for TestFail")

	// Verify that each run (in this case, only one) has some non-empty output
	for runID, outputs := range testFailResult.FailedOutputs {
		t.Logf("Failed outputs for run %s: %v", runID, outputs)
		require.NotEmpty(t, outputs, "Failed outputs should not be empty for TestFail")
	}
}

func TestSkippedTests(t *testing.T) {
	t.Parallel()

	runner := Runner{
		ProjectPath:      "./",
		Verbose:          true,
		RunCount:         1,
		SelectTests:      []string{"TestSkipped"}, // Known skipping test
		CollectRawOutput: true,
	}

	testResults, err := runner.RunTestPackages([]string{flakyTestPackagePath})
	require.NoError(t, err, "running tests should not produce an unexpected error")

	var testSkipResult *reports.TestResult
	for i := range testResults {
		if testResults[i].TestName == "TestSkipped" {
			testSkipResult = &testResults[i]
			break
		}
	}
	require.NotNil(t, testSkipResult, "expected 'TestSkipped' result not found in report")

	// Check that the test was properly marked as skipped
	require.True(t, testSkipResult.Skipped, "test 'TestSkipped' should be marked as skipped")
	require.Equal(t, 0, testSkipResult.Failures, "test 'TestSkipped' should have no failures")
	require.Equal(t, 0, testSkipResult.Successes, "test 'TestSkipped' should have no successes")
	require.Equal(t, 1, testSkipResult.Skips, "test 'TestSkipped' should have exactly one skip recorded")
}

func TestOmitOutputsOnSuccess(t *testing.T) {
	t.Parallel()

	runner := Runner{
		ProjectPath:          "./",
		Verbose:              true,
		RunCount:             1,
		SelectTests:          []string{"TestPass"}, // Known passing test
		CollectRawOutput:     true,
		OmitOutputsOnSuccess: true,
	}

	testResults, err := runner.RunTestPackages([]string{flakyTestPackagePath})
	require.NoError(t, err, "running tests should not produce an unexpected error")

	var testPassResult *reports.TestResult
	for i := range testResults {
		if testResults[i].TestName == "TestPass" {
			testPassResult = &testResults[i]
			break
		}
	}
	require.NotNil(t, testPassResult, "expected 'TestPass' result not found in report")
	require.Empty(t, testPassResult.PassedOutputs, "expected no passed outputs due to OmitOutputsOnSuccess")
	require.Empty(t, testPassResult.Outputs, "expected no captured outputs due to OmitOutputsOnSuccess and a successful test")
}
