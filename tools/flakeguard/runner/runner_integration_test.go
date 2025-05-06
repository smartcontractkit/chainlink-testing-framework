// Integration tests for the runner package, executing real tests.
package runner_test

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
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
)

var (
	flakyTestPackagePath = "./example_test_package"
	debugDir             = "_debug_outputs_integration"
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

func TestRunIntegration(t *testing.T) {
	var (
		zeroRuns        = 0
		oneCount        = 1
		defaultRunCount = 3
		successPassRate = 1.0
		failPassRate    = 0.0
	)
	testCases := []struct {
		name           string
		cfg            runnerConfig
		expectedTests  map[string]*expectedTestResult
		expectBuildErr bool
	}{
		{
			name: "default (integration)",
			cfg: runnerConfig{
				ProjectPath:      "../",
				RunCount:         defaultRunCount,
				SkipTests:        []string{"TestPanic", "TestFlakyPanic", "TestSubTestsSomePanic", "TestTimeout"},
				GoTestCountFlag:  &oneCount,
				OmitOutputs:      true,
				IgnoreSubtestErr: false,
				Tags:             []string{"example_package_tests"},
			},
			expectedTests: map[string]*expectedTestResult{
				"TestFlaky":                           {exactRuns: &defaultRunCount, someSuccesses: true, someFailures: true},
				"TestFail":                            {exactRuns: &defaultRunCount, allFailures: true, exactPassRate: &failPassRate},
				"TestFailLargeOutput":                 {exactRuns: &defaultRunCount, allFailures: true, exactPassRate: &failPassRate},
				"TestPass":                            {exactRuns: &defaultRunCount, allSuccesses: true, exactPassRate: &successPassRate},
				"TestSkipped":                         {exactRuns: &zeroRuns, allSkips: true, exactPassRate: &successPassRate},
				"TestRace":                            {exactRuns: &defaultRunCount, allSuccesses: true, exactPassRate: &successPassRate},
				"TestSubTestsAllPass":                 {exactRuns: &defaultRunCount, allSuccesses: true},
				"TestSubTestsAllPass/Pass1":           {exactRuns: &defaultRunCount, allSuccesses: true},
				"TestSubTestsAllPass/Pass2":           {exactRuns: &defaultRunCount, allSuccesses: true},
				"TestFailInParentAfterSubTests":       {exactRuns: &defaultRunCount, allFailures: true},
				"TestFailInParentAfterSubTests/Pass1": {exactRuns: &defaultRunCount, allSuccesses: true},
				"TestFailInParentAfterSubTests/Pass2": {exactRuns: &defaultRunCount, allSuccesses: true},
				"TestFailInParentBeforeSubTests":      {exactRuns: &defaultRunCount, allFailures: true},
				"TestSubTestsAllFail":                 {exactRuns: &defaultRunCount, allFailures: true},
				"TestSubTestsAllFail/Fail1":           {exactRuns: &defaultRunCount, allFailures: true},
				"TestSubTestsAllFail/Fail2":           {exactRuns: &defaultRunCount, allFailures: true},
				"TestSubTestsSomeFail":                {exactRuns: &defaultRunCount, allFailures: true},
				"TestSubTestsSomeFail/Pass":           {exactRuns: &defaultRunCount, allSuccesses: true},
				"TestSubTestsSomeFail/Fail":           {exactRuns: &defaultRunCount, allFailures: true},
			},
		},
		{
			name: "race (integration)",
			cfg: runnerConfig{
				ProjectPath:      "../",
				RunCount:         defaultRunCount,
				SelectTests:      []string{"TestRace"},
				GoTestRaceFlag:   true,
				OmitOutputs:      true,
				IgnoreSubtestErr: false,
				Tags:             []string{"example_package_tests"},
			},
			expectedTests: map[string]*expectedTestResult{
				"TestRace": {race: true, maximumRuns: defaultRunCount, allFailures: true},
			},
		},
		{
			name: "always panic (integration)",
			cfg: runnerConfig{
				ProjectPath:     "../",
				RunCount:        defaultRunCount,
				SelectTests:     []string{"TestPanic"},
				GoTestCountFlag: &oneCount,
				OmitOutputs:     true,
				Tags:            []string{"example_package_tests"},
			},
			expectedTests: map[string]*expectedTestResult{
				"TestPanic": {packagePanic: true, testPanic: true, maximumRuns: defaultRunCount, allFailures: true},
			},
		},
		{
			name: "flaky panic (integration)",
			cfg: runnerConfig{
				ProjectPath:     "../",
				RunCount:        defaultRunCount,
				SelectTests:     []string{"TestFlakyPanic"},
				GoTestCountFlag: &oneCount,
				OmitOutputs:     true,
				Tags:            []string{"example_package_tests"},
			},
			expectedTests: map[string]*expectedTestResult{
				// This test panics on first run, passes on second. We run 3 times.
				// Expect PackagePanic=true, TestPanic=true (as it panicked at least once)
				// Expect some failures (at least 1), some successes (at least 1).
				// Exact runs should be defaultRunCount.
				"TestFlakyPanic": {exactRuns: &defaultRunCount, packagePanic: true, testPanic: true, someSuccesses: true, someFailures: true},
			},
		},
		{
			name: "subtest panic (integration)",
			cfg: runnerConfig{
				ProjectPath:     "../",
				RunCount:        defaultRunCount,
				SelectTests:     []string{"TestSubTestsSomePanic"},
				GoTestCountFlag: &oneCount,
				OmitOutputs:     true,
				Tags:            []string{"example_package_tests"},
			},
			expectedTests: map[string]*expectedTestResult{
				"TestSubTestsSomePanic":       {exactRuns: &defaultRunCount, packagePanic: true, testPanic: true, allFailures: true}, // Parent fails due to subtest panic
				"TestSubTestsSomePanic/Pass":  {exactRuns: &defaultRunCount, packagePanic: true, testPanic: true, allFailures: true}, // Inherits panic, successes become failures
				"TestSubTestsSomePanic/Panic": {exactRuns: &defaultRunCount, packagePanic: true, testPanic: true, allFailures: true}, // Panics directly
			},
		},
		{
			name: "failfast (integration)",
			cfg: runnerConfig{
				ProjectPath:     "../",
				RunCount:        defaultRunCount, // Will try 3 times, but fail-fast stops early
				SelectTests:     []string{"TestFail", "TestPass"},
				GoTestCountFlag: &oneCount,
				FailFast:        true,
				OmitOutputs:     true,
				Tags:            []string{"example_package_tests"},
			},
			expectedTests: map[string]*expectedTestResult{
				// Only one execution attempt happens because FailFast=true and TestFail fails.
				"TestFail": {exactRuns: &oneCount, allFailures: true},
				"TestPass": {exactRuns: &oneCount, allSuccesses: true},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			absProjectPath, err := filepath.Abs(tc.cfg.ProjectPath)
			require.NoError(t, err)

			tempDir, err := os.MkdirTemp("", "flakeguard-test")
			require.NoError(t, err)

			testRunner := runner.NewRunner(
				absProjectPath,
				false,
				tc.cfg.RunCount,
				tc.cfg.GoTestCountFlag,
				tc.cfg.GoTestRaceFlag,
				tc.cfg.GoTestTimeoutFlag,
				tc.cfg.Tags,
				tc.cfg.UseShuffle,
				tc.cfg.ShuffleSeed,
				tc.cfg.FailFast,
				tc.cfg.SkipTests,
				tc.cfg.SelectTests,
				tc.cfg.IgnoreSubtestErr,
				tc.cfg.OmitOutputs,
				tempDir,
				nil, // Use default executor
				nil, // Use default parser
			)

			testResults, err := testRunner.RunTestPackages([]string{"./runner/example_test_package"})

			if tc.expectBuildErr {
				require.Error(t, err)
				return
			}
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
				jsonResults, err := json.MarshalIndent(testResults, "", "  ")
				if err != nil {
					t.Logf("error marshalling test report: %v", err)
					return
				}
				err = os.WriteFile(resultsFileName, jsonResults, 0644) //nolint:gosec
				if err != nil {
					t.Logf("error writing test results: %v", err)
					return
				}
				t.Logf("Saved failing test results to %s", resultsFileName)
			})

			checkTestResults(t, tc.expectedTests, testResults)
		})
	}
}

// Helper function to check results against expectations
func checkTestResults(t *testing.T, expectedTests map[string]*expectedTestResult, actualResults []reports.TestResult) {
	t.Helper()
	assert.Equal(t, len(expectedTests), len(actualResults), "unexpected number of test results recorded")

	for _, result := range actualResults {
		t.Run(fmt.Sprintf("checking results of %s", result.TestName), func(t *testing.T) {
			require.NotNil(t, result, "test result was nil")
			expected, ok := expectedTests[result.TestName]
			require.True(t, ok, "unexpected test name found in results: %s", result.TestName)
			require.False(t, expected.seen, "test '%s' was seen multiple times", result.TestName)
			expected.seen = true

			if !expected.testPanic {
				assert.False(t, result.Panic, "test '%s' should not have panicked", result.TestName)
			}

			if expected.minimumRuns != nil {
				assert.GreaterOrEqual(t, result.Runs, *expected.minimumRuns, "test '%s' had fewer runs (%d) than expected minimum (%d)", result.TestName, result.Runs, *expected.minimumRuns)
			}
			if expected.exactRuns != nil {
				assert.Equal(t, *expected.exactRuns, result.Runs, "test '%s' had an unexpected number of runs", result.TestName)
			} else {
				assert.LessOrEqual(t, result.Runs, expected.maximumRuns, "test '%s' had more runs (%d) than expected maximum (%d)", result.TestName, result.Runs, expected.maximumRuns)
			}
			if expected.exactPassRate != nil {
				assert.InDelta(t, *expected.exactPassRate, result.PassRatio, 0.001, "test '%s' had an unexpected pass ratio", result.TestName)
			}
			if expected.minimumPassRate != nil {
				assert.Greater(t, result.PassRatio, *expected.minimumPassRate, "test '%s' had a pass ratio below the minimum", result.TestName)
			}
			if expected.maximumPassRate != nil {
				assert.Less(t, result.PassRatio, *expected.maximumPassRate, "test '%s' had a pass ratio above the maximum", result.TestName)
			}
			if expected.allSuccesses {
				assert.Equal(t, result.Runs, result.Successes, "test '%s' has %d runs and should have passed all, only passed %d", result.TestName, result.Runs, result.Successes)
				assert.Zero(t, result.Failures, "test '%s' has %d runs and should have passed all, but failed %d", result.TestName, result.Runs, result.Failures)
				assert.False(t, result.Panic, "test '%s' should not have panicked", result.TestName)
				assert.False(t, result.Race, "test '%s' should not have raced", result.TestName)
			}
			if expected.someSuccesses {
				assert.Greater(t, result.Successes, 0, "test '%s' has %d runs and should have passed some runs, passed none", result.TestName, result.Runs)
			}
			if expected.allFailures {
				assert.Equal(t, result.Runs, result.Failures, "test '%s' has %d runs and should have failed all, only failed %d", result.TestName, result.Runs, result.Failures)
				assert.Zero(t, result.Successes, "test '%s' has %d runs and should have failed all, but succeeded %d", result.TestName, result.Runs, result.Successes)
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
				assert.Greater(t, result.Failures, 0, "test '%s' has %d runs and should have failed some runs, failed none", result.TestName, result.Runs)
			}
			if expected.allSkips {
				assert.Equal(t, 0, result.Runs, "test '%s' has %d runs and should have skipped all of them, no runs expected", result.TestName, result.Runs)
				assert.True(t, result.Skipped, "test '%s' should be marked skipped", result.TestName)
				assert.Zero(t, result.Successes, "test '%s' should have skipped all runs, but succeeded some", result.TestName)
				assert.Zero(t, result.Failures, "test '%s' should have skipped all runs, but failed some", result.TestName)
				assert.False(t, result.Panic, "test '%s' should not have panicked", result.TestName)
				assert.False(t, result.Race, "test '%s' should not have raced", result.TestName)
			}
			if expected.race {
				assert.True(t, result.Race, "test '%s' should have a data race", result.TestName)
				assert.GreaterOrEqual(t, result.Failures, 1, "test '%s' should have failed due to race", result.TestName)
			}
		})
	}

	allTestsRun := []string{}
	for testName, expected := range expectedTests {
		if expected.seen {
			allTestsRun = append(allTestsRun, testName)
		}
	}
	for testName, expected := range expectedTests {
		require.True(t, expected.seen, "expected test '%s' not found in test runs\nAll tests run: %s", testName, strings.Join(allTestsRun, ", "))
	}
}

type runnerConfig struct {
	ProjectPath       string
	RunCount          int
	GoTestCountFlag   *int
	GoTestRaceFlag    bool
	GoTestTimeoutFlag string
	Tags              []string
	UseShuffle        bool
	ShuffleSeed       string
	FailFast          bool
	SkipTests         []string
	SelectTests       []string
	OmitOutputs       bool
	IgnoreSubtestErr  bool
}
