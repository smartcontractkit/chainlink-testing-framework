package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	prettyPath, err := prettyProjectPath("./")
	require.NoError(t, err)
	assert.Equal(t, "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard", prettyPath)
}

func TestRun(t *testing.T) {
	var (
		zeroRuns        = 0
		oneRun          = 1
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
				UseRace:          false,
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
				"TestParentWithFailingParentAndSubtest": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestParentWithFailingParentAndSubtest/FailingSubtest": {
					exactRuns:   &defaultTestRunCount,
					allFailures: true,
				},
				"TestParentWithFailingParentAndSubtest/PassingSubtest": {
					exactRuns:    &defaultTestRunCount,
					allSuccesses: true,
				},
			},
		},
		{
			name: "always panic",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         defaultTestRunCount,
				UseRace:          false,
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
				UseRace:          false,
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
				UseRace:          false,
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
				UseRace:          false,
				SkipTests:        []string{},
				SelectTests:      []string{"TestFail", "TestPass"},
				FailFast:         true,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				"TestFail": {
					exactRuns:   &oneRun,
					allFailures: true,
				},
				"TestPass": {
					exactRuns:    &oneRun,
					allSuccesses: true,
				},
			},
		},
		{
			name: "subtest fails but parent is not failing",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         2, // run it a couple times
				UseRace:          false,
				SkipTests:        []string{},
				SelectTests:      []string{"TestParentWithFailingSubtest"},
				FailFast:         false,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				// The parent test
				"TestParentWithFailingSubtest": {
					// We expect the parent test to pass every run, because there's no parent-level fail
					allSuccesses: true,
					// or exactPassRate: &successPassRate,
					exactRuns: &[]int{2}[0], // 2 runs
				},
				// The failing subtest
				"TestParentWithFailingSubtest/FailingSubtest": {
					allFailures: true,
					exactRuns:   &[]int{2}[0],
				},
				// The passing subtest
				"TestParentWithFailingSubtest/PassingSubtest": {
					allSuccesses: true,
					exactRuns:    &[]int{2}[0],
				},
			},
		},
		{
			name: "parent fails and subtest fails",
			runner: Runner{
				ProjectPath:      "./",
				Verbose:          true,
				RunCount:         2, // run it a couple times
				UseRace:          false,
				SkipTests:        []string{},
				SelectTests:      []string{"TestParentWithFailingParentAndSubtest"},
				FailFast:         false,
				CollectRawOutput: true,
			},
			expectedTests: map[string]*expectedTestResult{
				// The parent test: we expect it to fail every run because it prints an error
				"TestParentWithFailingParentAndSubtest": {
					exactRuns:   &[]int{2}[0],
					allFailures: true,
				},
				// The failing subtest: it fails on all runs
				"TestParentWithFailingParentAndSubtest/FailingSubtest": {
					exactRuns:   &[]int{2}[0],
					allFailures: true,
				},
				// The passing subtest: it passes on all runs
				"TestParentWithFailingParentAndSubtest/PassingSubtest": {
					exactRuns:    &[]int{2}[0],
					allSuccesses: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testReport, err := tc.runner.RunTestPackages([]string{flakyTestPackagePath})
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
				jsonResults, err := json.Marshal(testReport)
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

			if tc.runner.FailFast {
				require.Equal(t, 1, testReport.SummaryData.TestRunCount, "unexpected number of unique tests run")
			} else {
				require.Equal(t, tc.runner.RunCount, testReport.SummaryData.TestRunCount, "unexpected number of test runs")
			}

			require.Equal(t, tc.runner.UseRace, testReport.RaceDetection, "unexpected race usage")

			assert.Equal(t, len(tc.expectedTests), len(testReport.Results), "unexpected number of test results")
			for _, result := range testReport.Results {
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
						"test '%s' doesn't match Runs count with results counts\n%s", result.TestName, resultsString(result),
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

func TestAttributePanicToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		packageName      string
		expectedTestName string
		expectedTimeout  bool
		panicEntries     []entry
	}{
		{
			name:             "properly attributed panic",
			expectedTestName: "TestPanic",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries:     properlyAttributedPanicEntries,
		},
		{
			name:             "improperly attributed panic",
			expectedTestName: "TestPanic",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries:     improperlyAttributedPanicEntries,
		},
		{
			name:             "timeout panic",
			expectedTestName: "TestTimedOut",
			expectedTimeout:  true,
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries:     timedOutPanicEntries,
		},
		{
			name:             "subtest panic",
			expectedTestName: "TestSubTestsSomePanic",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries:     subTestPanicEntries,
		},
		{
			name:         "empty",
			packageName:  "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries: []entry{},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testName, timeout, err := attributePanicToTest(tc.packageName, tc.panicEntries)
			assert.Equal(t, tc.expectedTimeout, timeout, "test timeout not correctly discovered")
			if tc.expectedTestName == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTestName, testName, "test panic not attributed correctly")
			}
		})
	}
}

func TestAttributeRaceToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		packageName      string
		expectedTestName string
		raceEntries      []entry
	}{
		{
			name:             "properly attributed race",
			expectedTestName: "TestRace",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries:      properlyAttributedRaceEntries,
		},
		{
			name:             "improperly attributed race",
			expectedTestName: "TestRace",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries:      improperlyAttributedRaceEntries,
		},
		{
			name:        "empty",
			packageName: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries: []entry{
				{},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testName, err := attributeRaceToTest(tc.packageName, tc.raceEntries)
			if tc.expectedTestName == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTestName, testName, "test race not attributed correctly")
			}
		})
	}
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

	testReport, err := runner.RunTestPackages([]string{flakyTestPackagePath})
	require.NoError(t, err, "running tests should not produce an unexpected error")

	require.Equal(t, 1, testReport.SummaryData.TotalRuns, "unexpected number of test runs")

	var testFailResult *reports.TestResult
	for i := range testReport.Results {
		if testReport.Results[i].TestName == "TestFail" {
			testFailResult = &testReport.Results[i]
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

	testReport, err := runner.RunTestPackages([]string{flakyTestPackagePath})
	require.NoError(t, err, "running tests should not produce an unexpected error")

	require.Equal(t, 0, testReport.SummaryData.TotalRuns, "unexpected number of test runs")
	require.Equal(t, 1, len(testReport.Results), "unexpected number of test results")
	require.Equal(t, 0, testReport.SummaryData.TestRunCount, "unexpected test run count")

	var testSkipResult *reports.TestResult
	for i := range testReport.Results {
		if testReport.Results[i].TestName == "TestSkipped" {
			testSkipResult = &testReport.Results[i]
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

	testReport, err := runner.RunTestPackages([]string{flakyTestPackagePath})
	require.NoError(t, err, "running tests should not produce an unexpected error")

	require.Equal(t, 1, testReport.SummaryData.TotalRuns, "unexpected number of test runs")

	var testPassResult *reports.TestResult
	for i := range testReport.Results {
		if testReport.Results[i].TestName == "TestPass" {
			testPassResult = &testReport.Results[i]
			break
		}
	}
	require.NotNil(t, testPassResult, "expected 'TestPass' result not found in report")
	require.Empty(t, testPassResult.PassedOutputs, "expected no passed outputs due to OmitOutputsOnSuccess")
	require.Empty(t, testPassResult.Outputs, "expected no captured outputs due to OmitOutputsOnSuccess and a successful test")
}

func TestGetOrCreateTestResult(t *testing.T) {
	tests := []struct {
		key                 string
		expectedTestPackage string
		expectedTestName    string
	}{
		{
			key:                 "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/TestParentWithFailingSubtest",
			expectedTestPackage: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			expectedTestName:    "TestParentWithFailingSubtest",
		},
		{
			key:                 "somepackage/TestFunction",
			expectedTestPackage: "somepackage",
			expectedTestName:    "TestFunction",
		},
		{
			key:                 "TestFunctionWithoutPackage",
			expectedTestPackage: "",
			expectedTestName:    "TestFunctionWithoutPackage",
		},
		{
			key:                 "smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/TestParentWithFailingSubtest/SubA/SubB",
			expectedTestPackage: "smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			expectedTestName:    "TestParentWithFailingSubtest/SubA/SubB",
		},
	}

	for _, tc := range tests {
		testResultsMap := make(map[string]*reports.TestResult)
		result := getOrCreateTestResult(testResultsMap, tc.key)
		if result.TestPackage != tc.expectedTestPackage {
			t.Errorf("For key %q, expected TestPackage %q, got %q", tc.key, tc.expectedTestPackage, result.TestPackage)
		}
		if result.TestName != tc.expectedTestName {
			t.Errorf("For key %q, expected TestName %q, got %q", tc.key, tc.expectedTestName, result.TestName)
		}
		// Verify that subsequent lookups return the same instance.
		duplicate := getOrCreateTestResult(testResultsMap, tc.key)
		if duplicate != result {
			t.Errorf("Subsequent call did not return the same instance for key %q", tc.key)
		}
	}
}

func TestFailedLinesIndicateRealParentFailure(t *testing.T) {
	tests := []struct {
		name           string
		parentName     string
		failOuts       map[string][]string
		expectRealFail bool
	}{
		{
			name:       "Only subtest failure – no parent failure",
			parentName: "TestParentWithFailingSubtest",
			failOuts: map[string][]string{
				"stdout": {
					"=== RUN   TestParentWithFailingSubtest\n",
					"=== RUN   TestParentWithFailingSubtest/FailingSubtest\n",
					"    example_tests_test.go:323: This subtest always fails.\n",
					"=== RUN   TestParentWithFailingSubtest/PassingSubtest\n",
					"--- FAIL: TestParentWithFailingSubtest (0.00s)\n",
					"    --- FAIL: TestParentWithFailingSubtest/FailingSubtest (0.00s)\n",
					"    --- PASS: TestParentWithFailingSubtest/PassingSubtest (0.00s)\n",
					"FAIL\n",
					"FAIL\n",
				},
			},
			expectRealFail: false,
		},
		{
			name:       "Parent and subtest failure – parent error summary present (from JSON output)",
			parentName: "TestParentWithFailingSubtest",
			failOuts: map[string][]string{
				"stdout": {
					"=== RUN   TestParentWithFailingSubtest\n",
					"    example_tests_test.go:329: parent fails\n",
					"--- FAIL: TestParentWithFailingSubtest (0.00s)\n",
					// (Subtest outputs would be collected separately.)
				},
			},
			expectRealFail: true,
		},
		{
			name:       "No .go error lines",
			parentName: "TestParent",
			failOuts: map[string][]string{
				"stderr": {"Some error occurred", "Another error message"},
			},
			expectRealFail: false,
		},
		{
			name:       "Subtest error only (with parent's name embedded)",
			parentName: "TestParent",
			failOuts: map[string][]string{
				"stdout": {"    example_tests_test.go:45: TestParent/SubTest error\n"},
			},
			expectRealFail: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := failedLinesIndicateRealParentFailure(tc.parentName, tc.failOuts)
			if result != tc.expectRealFail {
				t.Errorf("For %q, expected %v, got %v. failOuts: %v", tc.name, tc.expectRealFail, result, tc.failOuts)
			}
		})
	}
}

func TestNormalizeParentFailures(t *testing.T) {
	// Keys for our test results.
	parentKey := "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/TestParentWithFailingSubtest"
	subFailKey := parentKey + "/FailingSubtest"
	subPassKey := parentKey + "/PassingSubtest"

	// Test case 1: Parent test with only subtest-level failures.
	// In this scenario (from JSON output) the parent's outputs do not include a parent-level error.
	testsMap1 := map[string]*reports.TestResult{
		parentKey: {
			TestName:    "TestParentWithFailingSubtest",
			TestPackage: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			Runs:        2,
			Failures:    2,
			Successes:   0,
			PassRatio:   0,
			FailedOutputs: map[string][]string{
				"stdout": {
					"=== RUN   TestParentWithFailingSubtest\n",
					"    --- FAIL: TestParentWithFailingSubtest/FailingSubtest (0.00s)\n",
					"    example_tests_test.go:323: This subtest always fails.\n",
					"FAIL\n",
				},
			},
		},
		subFailKey: {
			TestName:    "TestParentWithFailingSubtest/FailingSubtest",
			TestPackage: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			Runs:        2,
			Failures:    2,
			Successes:   0,
			PassRatio:   0,
			FailedOutputs: map[string][]string{
				"stdout": {
					"=== RUN   TestParentWithFailingSubtest/FailingSubtest\n",
					"    example_tests_test.go:323: This subtest always fails.\n",
					"    --- FAIL: TestParentWithFailingSubtest/FailingSubtest (0.00s)\n",
					"FAIL\n",
				},
			},
		},
		subPassKey: {
			TestName:      "TestParentWithFailingSubtest/PassingSubtest",
			TestPackage:   "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			Runs:          2,
			Failures:      0,
			Successes:     2,
			PassRatio:     1,
			FailedOutputs: map[string][]string{},
		},
	}

	normalizeParentFailures(testsMap1)

	// For testsMap1, no parent-level error is detected so the parent's failures are cleared.
	parentRes1 := testsMap1[parentKey]
	if parentRes1.Failures != 0 {
		t.Errorf("expected parent Failures to be 0, got %d", parentRes1.Failures)
	}
	if parentRes1.Successes != parentRes1.Runs {
		t.Errorf("expected parent Successes (%d) to equal Runs (%d)", parentRes1.Successes, parentRes1.Runs)
	}
	if parentRes1.PassRatio != 1.0 {
		t.Errorf("expected parent PassRatio to be 1.0, got %f", parentRes1.PassRatio)
	}
	if len(parentRes1.FailedOutputs) != 0 {
		t.Errorf("expected parent FailedOutputs to be cleared, got %v", parentRes1.FailedOutputs)
	}

	// Test case 2: Parent test with both subtest and parent-level failures.
	// Here the JSON output for the parent includes a line indicating a parent-level error:
	// "    example_tests_test.go:329: parent fails\n"
	testsMap2 := map[string]*reports.TestResult{
		parentKey: {
			TestName:    "TestParentWithFailingSubtest",
			TestPackage: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			Runs:        2,
			Failures:    2,
			Successes:   0,
			PassRatio:   0,
			FailedOutputs: map[string][]string{
				"stdout": {
					"=== RUN   TestParentWithFailingSubtest\n",
					"    example_tests_test.go:329: parent fails\n",
					"--- FAIL: TestParentWithFailingSubtest (0.00s)\n",
				},
			},
		},
		subFailKey: {
			TestName:    "TestParentWithFailingSubtest/FailingSubtest",
			TestPackage: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			Runs:        2,
			Failures:    2,
			Successes:   0,
			PassRatio:   0,
			FailedOutputs: map[string][]string{
				"stdout": {
					"=== RUN   TestParentWithFailingSubtest/FailingSubtest\n",
					"    example_tests_test.go:323: This subtest always fails.\n",
					"    --- FAIL: TestParentWithFailingSubtest/FailingSubtest (0.00s)\n",
				},
			},
		},
		subPassKey: {
			TestName:      "TestParentWithFailingSubtest/PassingSubtest",
			TestPackage:   "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package",
			Runs:          2,
			Failures:      0,
			Successes:     2,
			PassRatio:     1,
			FailedOutputs: map[string][]string{},
		},
	}

	normalizeParentFailures(testsMap2)

	// In testsMap2, because a parent-level error is detected, the parent's failure counts remain unchanged.
	parentRes2 := testsMap2[parentKey]
	if parentRes2.Failures != 2 {
		t.Errorf("expected parent Failures to remain 2, got %d", parentRes2.Failures)
	}
	if parentRes2.Runs != 2 {
		t.Errorf("expected parent Runs to remain 2, got %d", parentRes2.Runs)
	}
	if parentRes2.PassRatio != 0 {
		t.Errorf("expected parent PassRatio to remain 0, got %f", parentRes2.PassRatio)
	}
	expectedOutputs2 := map[string][]string{
		"stdout": {
			"=== RUN   TestParentWithFailingSubtest\n",
			"    example_tests_test.go:329: parent fails\n",
			"--- FAIL: TestParentWithFailingSubtest (0.00s)\n",
		},
	}
	if !reflect.DeepEqual(parentRes2.FailedOutputs, expectedOutputs2) {
		t.Errorf("expected parent FailedOutputs to remain unchanged, got %v", parentRes2.FailedOutputs)
	}

	// Subtest entries remain unchanged.
	subFailRes := testsMap2[subFailKey]
	if subFailRes.Failures != 2 || subFailRes.Successes != 0 {
		t.Errorf("expected subtest (FailingSubtest) counts to remain unchanged, got Failures=%d, Successes=%d",
			subFailRes.Failures, subFailRes.Successes)
	}
	subPassRes := testsMap2[subPassKey]
	if subPassRes.Failures != 0 || subPassRes.Successes != 2 {
		t.Errorf("expected subtest (PassingSubtest) counts to remain unchanged, got Failures=%d, Successes=%d",
			subPassRes.Failures, subPassRes.Successes)
	}
}

var (
	improperlyAttributedPanicEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "panic: This test intentionally panics [recovered]\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\tpanic: This test intentionally panics\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "goroutine 25 [running]:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "testing.tRunner.func1.2({0x1008cde80, 0x1008f7d90})\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "testing.tRunner.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "panic({0x1008cde80?, 0x1008f7d90?})\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestPanic(0x140000b6ea0?)\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\t/Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:51 +0x30\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "testing.tRunner(0x140000b6ea0, 0x1008f73d0)\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0xe4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "created by testing.(*T).Run in goroutine 1\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x314\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Output: "FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package\t0.170s\n"},
		{Action: "fail", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Elapsed: 0.171},
	}
	properlyAttributedPanicEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "panic: This test intentionally panics [recovered]\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\tpanic: This test intentionally panics\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "goroutine 25 [running]:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "testing.tRunner.func1.2({0x1008cde80, 0x1008f7d90})\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "testing.tRunner.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "panic({0x1008cde80?, 0x1008f7d90?})\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestPanic(0x140000b6ea0?)\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\t/Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:51 +0x30\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "testing.tRunner(0x140000b6ea0, 0x1008f73d0)\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0xe4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "created by testing.(*T).Run in goroutine 1\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPanic", Output: "\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x314\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Output: "FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package\t0.170s\n"},
		{Action: "fail", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Elapsed: 0.171},
	}
	subTestPanicEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "panic: This subtest always panics [recovered]"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "panic: This subtest always panics"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "goroutine 23 [running]:"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "testing.tRunner.func1.2({0x100489e80, 0x1004b3e30})"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "	/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "testing.tRunner.func1()"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "	/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "panic({0x100489e80?, 0x1004b3e30?})"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "	/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestSubTestsSomePanic.func2(0x140000c81a0?)"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "	/Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:43 +0x30"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "testing.tRunner(0x140000c81a0, 0x1004b34d0)"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "	/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0xe4"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "created by testing.(*T).Run in goroutine 6"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSubTestsAllPass/Pass2", Output: "	/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x314"},
	}
	timedOutPanicEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "panic: test timed out after 10m0s"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "running tests"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "TestTimedOut (10m0s)"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "goroutine 397631 [running]:"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "testing.(*M).startAlarm.func1()"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "	/opt/hostedtoolcache/go/1.23.3/x64/src/testing/testing.go:2373 +0x385"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "created by time.goFunc"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestTimeout", Output: "/opt/hostedtoolcache/go/1.23.3/x64/src/time/sleep.go:215 +0x2d"},
	}

	improperlyAttributedRaceEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Read at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0x94\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Previous write at 0x00c000292028 by goroutine 12:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 12 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Previous write at 0x00c000292028 by goroutine 14:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 14 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Read at 0x00c000292028 by goroutine 19:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:68 +0xb8\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Previous write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 19 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "    testing.go:1399: race detected during execution of test\n"},
	}
	properlyAttributedRaceEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Read at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0x94\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Previous write at 0x00c000292028 by goroutine 12:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 12 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Previous write at 0x00c000292028 by goroutine 14:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 14 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Read at 0x00c000292028 by goroutine 19:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:68 +0xb8\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Previous write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 19 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "    testing.go:1399: race detected during execution of test\n"},
	}
)
