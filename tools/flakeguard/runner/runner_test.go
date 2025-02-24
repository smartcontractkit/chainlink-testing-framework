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

func TestZeroOutParentFailsIfSubtestOnlyFails(t *testing.T) {

	// 1) All subtests fail, but parent has no parent-level fail line
	t.Run("AllSubtestsFailButNoParentFailLine => zero out parent's fail", func(t *testing.T) {
		// We'll simulate "TestParentAllFailSubtests" from the logs.
		// The parent is failing, but all lines mention subtests (FailA, FailB).
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParentAllFailSubtests": {
				TestName:    "TestParentAllFailSubtests",
				TestPackage: "pkg",
				Runs:        2, // total runs
				Failures:    1, // parent is marked fail in run2
				Successes:   1,
				PassRatio:   0.5,
				FailedOutputs: map[string][]string{
					"run2": {
						// Lines from your logs referencing subtests only:
						"=== CONT  TestParentAllFailSubtests/FailA",
						"    example_tests_test.go:246: This subtest always fails",
						"    --- FAIL: TestParentAllFailSubtests/FailA (0.00s)",
						"=== CONT  TestParentAllFailSubtests/FailB",
						"    example_tests_test.go:250: This subtest always fails",
						"    --- FAIL: TestParentAllFailSubtests/FailB (0.00s)",
						// Notice we do NOT include a line like
						// "FAIL: TestParentAllFailSubtests (0.00s)"
						// That would be a genuine parent-level fail.
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParentAllFailSubtests": {"FailA", "FailB"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParentAllFailSubtests"]
		// Because all fail lines reference subtests only,
		// we expect parent's failure to be zeroed out:
		assert.Equal(t, 0, parent.Failures)
		assert.Equal(t, 1, parent.Runs, "2 => 1 after removing that 1 failure")
		assert.InDelta(t, 1.0, parent.PassRatio, 0.0001)
		assert.Empty(t, parent.FailedOutputs)
	})

	// 2) Some subtests fail, parent has no real parent-level line => zero out
	t.Run("SomeSubtestsFailButNoParentFailLine => zero out parent's fail", func(t *testing.T) {
		// Example: "TestParentSomeFailSubtests" => partial fail in subtest "Fail".
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParentSomeFailSubtests": {
				TestName:    "TestParentSomeFailSubtests",
				TestPackage: "pkg",
				Runs:        3,
				Failures:    1,
				Successes:   2,
				PassRatio:   2.0 / 3.0,
				FailedOutputs: map[string][]string{
					"run2": {
						// Real lines from logs referencing subtest "Fail" only:
						"=== CONT  TestParentSomeFailSubtests/Fail",
						"    example_tests_test.go:265: This subtest fails",
						"    --- FAIL: TestParentSomeFailSubtests/Fail (0.00s)",
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParentSomeFailSubtests": {"Pass", "Fail"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParentSomeFailSubtests"]
		assert.Equal(t, 0, parent.Failures, "parent test should have 0 failures now")
		assert.Equal(t, 2, parent.Runs, "3 => 2 after removing 1 failure")
		assert.InDelta(t, 1.0, parent.PassRatio, 0.0001)
		assert.Empty(t, parent.FailedOutputs)
	})

	// 3) Parent fails before subtests => genuine fail remains
	t.Run("ParentFailsBeforeSubtests => remains failing", func(t *testing.T) {
		// Example: “TestParentOwnFailBeforeSubtests” => logs show a parent-level line
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParentOwnFailBeforeSubtests": {
				TestName:    "TestParentOwnFailBeforeSubtests",
				TestPackage: "pkg",
				Runs:        1,
				Failures:    1,
				Successes:   0,
				PassRatio:   0.0,
				FailedOutputs: map[string][]string{
					"run1": {
						// A genuine parent-level fail line:
						"Error in TestParentOwnFailBeforeSubtests: parent test fails immediately",
						// This also references subtest, but that doesn't matter;
						// parent line is enough to keep it failing.
						"Error in TestParentOwnFailBeforeSubtests/Pass1: subtest never ran",
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParentOwnFailBeforeSubtests": {"Pass1", "Pass2"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParentOwnFailBeforeSubtests"]
		assert.Equal(t, 1, parent.Failures, "still fails, genuine parent-level line")
		assert.Equal(t, 1, parent.Runs)
		assert.InDelta(t, 0.0, parent.PassRatio, 0.0001)
		assert.NotEmpty(t, parent.FailedOutputs)
	})

	// 4) Parent fails after subtests pass => remains failing
	t.Run("ParentFailsAfterSubtests => genuine fail remains", func(t *testing.T) {
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParentOwnFailAfterSubtests": {
				TestName:    "TestParentOwnFailAfterSubtests",
				TestPackage: "pkg",
				Runs:        2,
				Failures:    1,
				Successes:   1,
				PassRatio:   0.5,
				FailedOutputs: map[string][]string{
					"run2": {
						// Contains parent-level line referencing "TestParentOwnFailAfterSubtests" alone:
						"Error in TestParentOwnFailAfterSubtests: parent test fails after subtests pass",
						// Subtests lines also appear, but the parent-level line is enough
						"Error in TestParentOwnFailAfterSubtests/Pass1: subtest passes",
						"Error in TestParentOwnFailAfterSubtests/Pass2: subtest passes",
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParentOwnFailAfterSubtests": {"Pass1", "Pass2"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParentOwnFailAfterSubtests"]
		assert.Equal(t, 1, parent.Failures, "still fails, found genuine parent-level line")
		assert.Equal(t, 2, parent.Runs)
		assert.InDelta(t, 0.5, parent.PassRatio, 0.0001)
		assert.NotEmpty(t, parent.FailedOutputs)
	})

	// 5) Nested subtests: parent lines mention deeper "TestParent/Nest1/Nest2" => zero out
	t.Run("NestedSubtests => parent's lines mention deeper sub-subtest only => zero out", func(t *testing.T) {
		testDetails := map[string]*reports.TestResult{
			"pkg/TestNestedSubtests": {
				TestName:    "TestNestedSubtests",
				TestPackage: "pkg",
				Runs:        2,
				Failures:    1,
				Successes:   1,
				PassRatio:   0.5,
				FailedOutputs: map[string][]string{
					"run1": {
						"Error in TestNestedSubtests/Level1/Level2Fail: sub-subtest fails",
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestNestedSubtests": {"Level1"}, // e.g., "Level1" might also have "Level2Fail"
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestNestedSubtests"]
		assert.Equal(t, 0, parent.Failures)
		assert.Equal(t, 1, parent.Runs, "2 => 1 after zeroing out the parent's fail")
		assert.InDelta(t, 1.0, parent.PassRatio, 0.0001)
		assert.Empty(t, parent.FailedOutputs)
	})

	// 6) “Simpler” tests from your prior suite

	t.Run("Parent fails only referencing subtests => zero out parent's fails", func(t *testing.T) {
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParent": {
				TestName:    "TestParent",
				TestPackage: "pkg",
				Runs:        3,
				Failures:    1,
				Successes:   2,
				PassRatio:   2.0 / 3.0,
				FailedOutputs: map[string][]string{
					"run1": {
						"Error in TestParent/SubtestA: something bad",
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParent": {"SubtestA"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParent"]
		assert.Equal(t, 0, parent.Failures)
		assert.Equal(t, 2, parent.Runs)
		assert.InDelta(t, 1.0, parent.PassRatio, 0.0001)
		assert.Empty(t, parent.FailedOutputs)
	})

	t.Run("Parent with genuine fail => remains failing", func(t *testing.T) {
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParent": {
				TestName:    "TestParent",
				TestPackage: "pkg",
				Runs:        2,
				Failures:    1,
				Successes:   1,
				PassRatio:   0.5,
				FailedOutputs: map[string][]string{
					"run1": {
						// Genuine parent-level line
						"Error in TestParent: parent-level assertion failed",
					},
				},
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParent": {"SubtestA"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParent"]
		assert.Equal(t, 1, parent.Failures)
		assert.Equal(t, 2, parent.Runs)
		assert.InDelta(t, 0.5, parent.PassRatio, 0.0001)
		assert.NotEmpty(t, parent.FailedOutputs)
	})

	t.Run("Parent with zero failures => no change", func(t *testing.T) {
		testDetails := map[string]*reports.TestResult{
			"pkg/TestParent": {
				TestName:    "TestParent",
				TestPackage: "pkg",
				Runs:        3,
				Failures:    0,
				Successes:   3,
				PassRatio:   1.0,
			},
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParent": {"SubtestA", "SubtestB"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		parent := testDetails["pkg/TestParent"]
		assert.Equal(t, 3, parent.Runs)
		assert.Equal(t, 0, parent.Failures)
		assert.InDelta(t, 1.0, parent.PassRatio, 0.0001)
	})

	t.Run("Parent not found in testDetails => no crash", func(t *testing.T) {
		testDetails := map[string]*reports.TestResult{
			// some other test, but not pkg/TestParent
		}
		testsWithSubTests := map[string][]string{
			"pkg/TestParent": {"SubtestA"},
		}

		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)
		// just confirm no panic, nothing to assert
	})

	t.Run("NestedSubtests => only sub-subtest lines => parent's fail is zeroed out", func(t *testing.T) {
		// This mimics a scenario where Go initially flags the parent "TestNestedSubtests" as failed,
		// but all failing lines actually reference sub-subtests only:
		testDetails := map[string]*reports.TestResult{
			"pkg/TestNestedSubtests": {
				TestName:    "TestNestedSubtests",
				TestPackage: "pkg",
				// Suppose we had exactly 1 run that ended in 'fail' (incorrectly marking the parent):
				Runs:      1,
				Failures:  1,
				Successes: 0,
				PassRatio: 0.0,
				FailedOutputs: map[string][]string{
					"run1": {
						// All lines mention sub-subtest path "TestNestedSubtests/Level1/Level2Fail"
						// There's NO line referencing "TestNestedSubtests" alone.
						"=== CONT  TestNestedSubtests/Level1/Level2Fail",
						"    example_tests_test.go:315: This sub-subtest fails",
						"    --- FAIL: TestNestedSubtests/Level1/Level2Fail (0.00s)",
					},
				},
			},
			// (Optionally, if you track each sub-subtest individually, you'd have
			// "pkg/TestNestedSubtests/Level1/Level2Fail": {...} for completeness,
			// but it's not necessary to show the parent's zero-out logic.)
		}

		// Indicate the parent has a subtest "Level1," which might have its own children:
		testsWithSubTests := map[string][]string{
			"pkg/TestNestedSubtests": {"Level1"},
			// If your code needs deeper nesting, you could do:
			// "pkg/TestNestedSubtests/Level1": {"Level2Fail", "Level2Pass"},
		}

		// Run your function that zeros out parent fails if lines are subtest-only
		zeroOutParentFailsIfSubtestOnlyFails(testDetails, testsWithSubTests)

		// Now check that the parent is zeroed out, since no line references "TestNestedSubtests" alone
		parent := testDetails["pkg/TestNestedSubtests"]
		// The parent's failure should be removed:
		assert.Equal(t, 0, parent.Failures,
			"parent test should have 0 failures now, since only sub-subtest lines appear")
		// If the parent had 1 run all failing, but that was purely from sub-subtest => we remove it
		// So the parent's total runs become 0
		assert.Equal(t, 0, parent.Runs,
			"parent test's runs should go from 1 => 0 after removing that purely sub-subtest failure")
		// With 0 runs, we typically set passRatio = 1.0 in your logic (like a 'no-run is pass' fallback)
		assert.InDelta(t, 1.0, parent.PassRatio, 0.0001,
			"if there are 0 runs left, pass ratio is 1.0 by default")
		assert.Empty(t, parent.FailedOutputs,
			"failed outputs should be cleared after zeroing out the parent's fail")
	})
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
