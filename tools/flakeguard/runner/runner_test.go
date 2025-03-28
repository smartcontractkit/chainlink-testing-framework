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

func TestNoCompileTests(t *testing.T) {
	// Test that we are not swallowing test compilation errors
	t.Parallel()

	runner := Runner{
		ProjectPath:    "./",
		Verbose:        true,
		RunCount:       1,
		GoTestRaceFlag: false,
		FailFast:       false,
	}

	_, err := runner.RunTestPackages([]string{"./example_bad_test_package"})
	require.Error(t, err)
	require.ErrorIs(t, err, buildErr, "expected a compile error")
	require.NotErrorIs(t, err, failedToShowBuildErr, "should be able to print out build errors")
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
				ProjectPath:    "./",
				Verbose:        true,
				RunCount:       defaultTestRunCount,
				GoTestRaceFlag: false,
				SkipTests:      []string{"TestPanic", "TestFlakyPanic", "TestSubTestsSomePanic", "TestTimeout"},
				FailFast:       false,
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
				ProjectPath:    "./",
				Verbose:        true,
				RunCount:       defaultTestRunCount,
				GoTestRaceFlag: false,
				SkipTests:      []string{},
				SelectTests:    []string{"TestPanic"},
				FailFast:       false,
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
				ProjectPath:     "./",
				Verbose:         true,
				RunCount:        defaultTestRunCount,
				GoTestRaceFlag:  false,
				GoTestCountFlag: &oneCount,
				SkipTests:       []string{},
				SelectTests:     []string{"TestFlakyPanic"},
				FailFast:        false,
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
				ProjectPath:    "./",
				Verbose:        true,
				RunCount:       defaultTestRunCount,
				GoTestRaceFlag: false,
				SkipTests:      []string{},
				SelectTests:    []string{"TestSubTestsSomePanic"},
				FailFast:       false,
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
				ProjectPath:    "./",
				Verbose:        true,
				RunCount:       defaultTestRunCount,
				GoTestRaceFlag: false,
				SkipTests:      []string{},
				SelectTests:    []string{"TestFail", "TestPass"},
				FailFast:       true,
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

func TestAttributePanicToTest(t *testing.T) {
	t.Parallel()

	// Test cases: each test case contains a slice of output strings.
	testCases := []struct {
		name             string
		expectedTestName string
		expectedTimeout  bool
		outputs          []string
	}{
		{
			name:             "properly attributed panic",
			expectedTestName: "TestPanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This test intentionally panics [recovered]",
				"\tpanic: This test intentionally panics",
				"goroutine 25 [running]:",
				"testing.tRunner.func1.2({0x1008cde80, 0x1008f7d90})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc",
				"testing.tRunner.func1()",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334",
				"panic({0x1008cde80?, 0x1008f7d90?})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestPanic(0x140000b6ea0?)",
			},
		},
		{
			name:             "improperly attributed panic",
			expectedTestName: "TestPanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This test intentionally panics [recovered]",
				"TestPanic(0x140000b6ea0?)",
				"goroutine 25 [running]:",
				"testing.tRunner.func1.2({0x1008cde80, 0x1008f7d90})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc",
				"testing.tRunner.func1()",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334",
				"panic({0x1008cde80?, 0x1008f7d90?})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestPanic(0x140000b6ea0?)",
			},
		},
		{
			name:             "log after test complete panic",
			expectedTestName: "Test_workflowRegisteredHandler/skips_fetch_if_secrets_url_is_missing",
			expectedTimeout:  false,
			outputs: []string{
				"panic: Log in goroutine after Test_workflowRegisteredHandler/skips_fetch_if_secrets_url_is_missing has completed: 2025-03-28T17:18:16.703Z\tDEBUG\tCapabilitiesRegistry\tcapabilities/registry.go:69\tget capability\t{\"version\": \"unset@unset\", \"id\": \"basic-test-trigger@1.0.0\"}",
				"goroutine 646 [running]:",
				"testing.(*common).logDepth(0xc000728000, {0xc0001b9400, 0x9b}, 0x3)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1064 +0x69f",
				"testing.(*common).log(...)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1046",
				"testing.(*common).Logf(0xc000728000, {0x6000752, 0x2}, {0xc001070430, 0x1, 0x1})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1097 +0x9f",
				"go.uber.org/zap/zaptest.TestingWriter.Write({{0x7fb811aa2818?, 0xc000728000?}, 0x20?}, {0xc001074000, 0x9c, 0x400})",
				"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.27.0/zaptest/logger.go:146 +0x11d",
				"go.uber.org/zap/zapcore.(*ioCore).Write(0xc000bff1d0, {0xff, {0xc1f1d45629e99436, 0x252667087c, 0x87b3fa0}, {0x602a730, 0x14}, {0x601d42f, 0xe}, {0x1, ...}, ...}, ...)",
				"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.27.0/zapcore/core.go:99 +0x18e",
				"go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc00106dba0, {0xc00101d400, 0x1, 0x2})",
				"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.27.0/zapcore/entry.go:253 +0x1ed",
				"go.uber.org/zap.(*SugaredLogger).log(0xc0001e48b8, 0xff, {0x601d42f, 0xe}, {0x0, 0x0, 0x0}, {0xc00101bf40, 0x2, 0x2})",
				"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.27.0/sugar.go:355 +0x12d",
				"go.uber.org/zap.(*SugaredLogger).Debugw(...)",
				"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.27.0/sugar.go:251",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities.(*Registry).Get(0xc000ab88c0, {0x20?, 0x87bb320?}, {0xc0013282a0, 0x18})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/registry.go:69 +0x1cf",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities.(*Registry).GetTrigger(0xc000ab88c0, {0x67d38a8, 0xc0011f22d0}, {0xc0013282a0, 0x18})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/registry.go:80 +0x6f",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).resolveWorkflowCapabilities(0xc000e75188, {0x67d38a8, 0xc0011f22d0})",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:198 +0x173",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).init.func1()",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:348 +0x2aa",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.retryable({0x67d38a8, 0xc0011f22d0}, {0x680c850, 0xc000e08210}, 0x1388, 0x0, 0xc000f0bf08)",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/retry.go:45 +0x402",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).init(0xc000e75188, {0x67d38a8, 0xc0011f22d0})",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:339 +0x225",
				"created by github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).Start.func1 in goroutine 390",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:179 +0xf37",
				"FAIL\tgithub.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer\t159.643s",
			},
		},
		{
			name:             "timeout panic with obvious culprit",
			expectedTestName: "TestTimedOut",
			expectedTimeout:  true,
			outputs: []string{
				"panic: test timed out after 10m0s",
				"running tests",
				"\tTestNoTimeout (9m59s)",
				"\tTestTimedOut (10m0s)",
				"goroutine 397631 [running]:",
				"testing.(*M).startAlarm.func1()",
				"\t/opt/hostedtoolcache/go/1.23.3/x64/src/testing/testing.go:2373 +0x385",
				"created by time.goFunc",
				"/opt/hostedtoolcache/go/1.23.3/x64/src/time/sleep.go:215 +0x2d",
			},
		},
		{
			name:             "subtest panic",
			expectedTestName: "TestSubTestsSomePanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This subtest always panics [recovered]",
				"panic: This subtest always panics",
				"goroutine 23 [running]:",
				"testing.tRunner.func1.2({0x100489e80, 0x1004b3e30})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc",
				"testing.tRunner.func1()",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334",
				"panic({0x100489e80?, 0x1004b3e30?})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestSubTestsSomePanic.func2(0x140000c81a0?)",
			},
		},
		{
			name:             "memory_test panic extraction",
			expectedTestName: "TestJobClientJobAPI",
			expectedTimeout:  false,
			outputs: []string{
				"panic: freeport: cannot allocate port block [recovered]",
				"\tpanic: freeport: cannot allocate port block",
				"goroutine 321 [running]:",
				"testing.tRunner.func1.2({0x5e0dd80, 0x72ebb40})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1734 +0x21c",
				"testing.tRunner.func1()",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1737 +0x35e",
				"panic({0x5e0dd80?, 0x72ebb40?})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/runtime/panic.go:787 +0x132",
				"github.com/hashicorp/consul/sdk/freeport.alloc()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:274 +0xad",
				"github.com/hashicorp/consul/sdk/freeport.initialize()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:124 +0x2d7",
				"sync.(*Once).doSlow(0xc0018eb600?, 0xc000da4a98?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:78 +0xab",
				"sync.(*Once).Do(...)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:69",
				"github.com/hashicorp/consul/sdk/freeport.Take(0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:303 +0xe5",
				"github.com/hashicorp/consul/sdk/freeport.GetN({0x7337708, 0xc000683dc0}, 0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:427 +0x48",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory_test.TestJobClientJobAPI(0xc000683dc0)",
				"\t/home/runner/work/chainlink/chainlink/deployment/environment/memory/job_service_client_test.go:116 +0xc6",
				"testing.tRunner(0xc000683dc0, 0x6d6c838)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1792 +0xf4",
				"created by testing.(*T).Run in goroutine 1",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1851 +0x413",
			},
		},
		{
			name:             "changeset_test panic extraction",
			expectedTestName: "TestDeployBalanceReader",
			expectedTimeout:  false,
			outputs: []string{
				"panic: freeport: cannot allocate port block [recovered]",
				"\tpanic: freeport: cannot allocate port block",
				"goroutine 378 [running]:",
				"testing.tRunner.func1.2({0x6063f40, 0x76367f0})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1734 +0x21c",
				"testing.tRunner.func1()",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1737 +0x35e",
				"panic({0x6063f40?, 0x76367f0?})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/runtime/panic.go:787 +0x132",
				"github.com/hashicorp/consul/sdk/freeport.alloc()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:274 +0xad",
				"github.com/hashicorp/consul/sdk/freeport.initialize()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:124 +0x2d7",
				"sync.(*Once).doSlow(0xa94f820?, 0xa8000a?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:78 +0xab",
				"sync.(*Once).Do(...)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:69",
				"github.com/hashicorp/consul/sdk/freeport.Take(0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:303 +0xe5",
				"github.com/hashicorp/consul/sdk/freeport.GetN({0x7684150, 0xc000583c00}, 0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:427 +0x48",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory.NewNodes(0xc000583c00, 0xff, 0xc001583d10, 0xc005aa0030, 0x1, 0x0, {0x0, {0x0, 0x0, 0x0, ...}, ...}, ...)",
				"\t/home/runner/work/chainlink/chainlink/deployment/environment/memory/environment.go:177 +0xa5",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory.NewMemoryEnvironment(_, {_, _}, _, {0x2, 0x0, 0x0, 0x1, 0x0, {0x0, ...}})",
				"\t/home/runner/work/chainlink/chainlink/deployment/environment/memory/environment.go:223 +0x10c",
				"github.com/smartcontractkit/chainlink/deployment/keystone/changeset_test.TestDeployBalanceReader(0xc000583c00)",
				"\t/home/runner/work/chainlink/chainlink/deployment/keystone/changeset/deploy_balance_reader_test.go:23 +0xf5",
				"testing.tRunner(0xc000583c00, 0x70843d0)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1792 +0xf4",
				"created by testing.(*T).Run in goroutine 1",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1851 +0x413",
				"    logger.go:146: 03:14:04.485880684\tINFO\tDeployed KeystoneForwarder 1.0.0 chain selector 909606746561742123 addr 0x72B66019aCEdc35F7F6e58DF94De95f3cBCC5971\t{\"version\": \"(devel)@unset\"}",
				"    logger.go:146: 03:14:04.486035865\tINFO\tdeploying forwarder\t{\"version\": \"(devel)@unset\", \"chainSelector\": 5548718428018410741}",
				"    logger.go:146: 2025-03-08T03:14:04.490Z\tINFO\tchangeset/jd_register_nodes.go:91\tregistered node\t{\"version\": \"unset@unset\", \"name\": \"node1\", \"id\": \"node:{id:\\\"895776f5ba0cc11c570a47b5cc3dbb8771da9262cfb545cd5d48251796af7f\\\"  public_key:\\\"895776f5ba0cc11c570a47b5cc3dbb8771da9262cfb545cd5d48251796af7f\\\"  is_enabled:true  is_connected:true  labels:{key:\\\"product\\\"  value:\\\"test-product\\\"}  labels:{key:\\\"environment\\\"  value:\\\"test-env\\\"}  labels:{key:\\\"nodeType\\\"  value:\\\"bootstrap\\\"}  labels:{key:\\\"don-0-don1\\\"}\"}",
			},
		},
		{
			name:             "empty",
			expectedTestName: "",
			expectedTimeout:  false,
			outputs:          []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testName, timeout, err := attributePanicToTest(tc.outputs)
			assert.Equal(t, tc.expectedTimeout, timeout, "timeout flag mismatch")
			if tc.expectedTestName == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTestName, testName, "test name mismatch")
			}
		})
	}
}

func TestFailToAttributePanicToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		expectedTimeout bool
		expectedError   error
		outputs         []string
	}{
		{
			name:            "no test name in panic",
			expectedTimeout: false,
			expectedError:   ErrFailedToAttributePanicToTest,
			outputs: []string{
				"panic: reflect: Elem of invalid type bool",
				"goroutine 104182 [running]:",
				"reflect.elem(0xc0569d9998?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/reflect/type.go:733 +0x9a",
				"reflect.(*rtype).Elem(0xa4dd940?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/reflect/type.go:737 +0x15",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.setPollingFilterOverrides(0x0, {0xc052040510, 0x1, 0xc?})",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:942 +0x492",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.(*ContractReaderService).addEventRead(_, _, {_, _}, {_, _}, {{0xc0544c4270, 0x9}, {0xc0544c4280, 0xc}, ...}, ...)",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:605 +0x13d",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.(*ContractReaderService).initNamespace(0xc054472540, 0xc01c37d440?)",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:443 +0x28b",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.NewContractReaderService({0x7fcf8b532040?, 0xc015b223e0?}, {0xc6ac960, 0xc05464e470}, {0xc0544384e0?, {0xc01c37d440?, 0xc054163b84?, 0xc054163b80?}}, {0x7fcf8071c7a0, 0xc0157928c0})",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:97 +0x287",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana.(*Relayer).NewContractReader(0xc015b2e150, {0x4d0102030cb384f5?, 0xb938300b5ca1aa13?}, {0xc05469c000, 0x1eedf, 0x20000})",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/relay.go:160 +0x205",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator.(*pluginOracleCreator).createReadersAndWriters(_, {_, _}, {_, _}, _, {0x3, {0x0, 0xa, 0x93, ...}, ...}, ...)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/oraclecreator/plugin.go:446 +0x338",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator.(*pluginOracleCreator).Create(0xc033a69ad0, {0xc6f5a10, 0xc02e4f9a40}, 0x3, {0x3, {0x0, 0xa, 0x93, 0x8f, 0x67, ...}, ...})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/oraclecreator/plugin.go:215 +0xc0c",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.createDON({0xc6f5a10, 0xc02e4f9a40}, {0x7fcf8b533ad0, 0xc015b97340}, {0xb6, 0x5e, 0x31, 0xd0, 0x35, 0xef, ...}, ...)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:367 +0x451",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).processAdded(0xc015723080, {0xc6f5a10, 0xc02e4f9a40}, 0xc053de2ff0)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:254 +0x239",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).processDiff(0xc015723080, {0xc6f5a10, 0xc02e4f9a40}, {0xc053de2ff0?, 0xc053de3020?, 0xc053de3050?})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:192 +0x68",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).tick(0xc015723080, {0xc6f5a10, 0xc02e4f9a40})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:178 +0x20b",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).monitor(0xc015723080)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:152 +0x112",
				"created by github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).Start.func1 in goroutine 1335",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:134 +0xa5",
				"FAIL\tgithub.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana\t184.801s",
			},
		},
		{
			name:            "fail to parse timeout duration",
			expectedTimeout: true,
			expectedError:   ErrFailedToParseTimeoutDuration,
			outputs: []string{
				"panic: test timed out after malformedDurationStr\n",
				"\trunning tests:\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain (22s)\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain/Remote_chains_owned_by_MCMS (22s)\n",
				"\t\tTestAlmostPanicTime (9m59s)\n",
				"\t\tTestConnectNewChain (1m1s)\n",
				"\t\tTestConnectNewChain/Use_production_router_(with_MCMS) (1m1s)\n",
				"\t\tTestJobSpecChangeset (0s)\n",
				"\t\tTest_ActiveCandidate (1m1s)\n",
				"goroutine 971967 [running]:\n",
				"testing.(*M).startAlarm.func1()\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2484 +0x605\n",
				"created by time.goFunc\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/time/sleep.go:215 +0x45\n",
				"goroutine 1 [chan receive]:\n",
				"testing.tRunner.func1()\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1753 +0x965\n",
				"testing.tRunner(0xc0013dac40, 0xc0025b7ae0)\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1798 +0x25f\n",
				"testing.runTests(0xc0010a0b70, {0x14366840, 0x25, 0x25}, {0x3?, 0x0?, 0x146214a0?})\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2277 +0x96d\n",
				"testing.(*M).Run(0xc0014732c0)\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2142 +0xeeb\n",
				"main.main()\n",
				"\t_testmain.go:119 +0x165\n",
			},
		},
		{
			name:            "fail to parse test duration",
			expectedTimeout: true,
			expectedError:   ErrDetectedTimeoutFailedParse,
			outputs: []string{
				"panic: test timed out after 10m0s\n",
				"\trunning tests:\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain (malformedDurationStr)\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain/Remote_chains_owned_by_MCMS (22s)\n",
				"\t\tTestAlmostPanicTime (9m59s)\n",
				"\t\tTestConnectNewChain (1m1s)\n",
				"\t\tTestConnectNewChain/Use_production_router_(with_MCMS) (1m1s)\n",
				"\t\tTestJobSpecChangeset (0s)\n",
				"\t\tTest_ActiveCandidate (1m1s)\n",
				"goroutine 971967 [running]:\n",
				"testing.(*M).startAlarm.func1()\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2484 +0x605\n",
				"created by time.goFunc\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/time/sleep.go:215 +0x45\n",
				"goroutine 1 [chan receive]:\n",
				"testing.tRunner.func1()\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1753 +0x965\n",
				"testing.tRunner(0xc0013dac40, 0xc0025b7ae0)\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1798 +0x25f\n",
				"testing.runTests(0xc0010a0b70, {0x14366840, 0x25, 0x25}, {0x3?, 0x0?, 0x146214a0?})\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2277 +0x96d\n",
				"testing.(*M).Run(0xc0014732c0)\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2142 +0xeeb\n",
				"main.main()\n",
				"\t_testmain.go:119 +0x165\n",
			},
		},
		{
			name:            "timeout panic without obvious culprit",
			expectedTimeout: true,
			expectedError:   ErrDetectedTimeoutFailedAttribution,
			outputs: []string{
				"panic: test timed out after 10m0s\n",
				"\trunning tests:\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain (22s)\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain/Remote_chains_owned_by_MCMS (22s)\n",
				"\t\tTestAlmostPanicTime (9m59s)\n",
				"\t\tTestConnectNewChain (1m1s)\n",
				"\t\tTestConnectNewChain/Use_production_router_(with_MCMS) (1m1s)\n",
				"\t\tTestJobSpecChangeset (0s)\n",
				"\t\tTest_ActiveCandidate (1m1s)\n",
				"goroutine 971967 [running]:\n",
				"testing.(*M).startAlarm.func1()\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2484 +0x605\n",
				"created by time.goFunc\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/time/sleep.go:215 +0x45\n",
				"goroutine 1 [chan receive]:\n",
				"testing.tRunner.func1()\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1753 +0x965\n",
				"testing.tRunner(0xc0013dac40, 0xc0025b7ae0)\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1798 +0x25f\n",
				"testing.runTests(0xc0010a0b70, {0x14366840, 0x25, 0x25}, {0x3?, 0x0?, 0x146214a0?})\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2277 +0x96d\n",
				"testing.(*M).Run(0xc0014732c0)\n",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:2142 +0xeeb\n",
				"main.main()\n",
				"\t_testmain.go:119 +0x165\n",
			},
		},
		{
			name:            "possible regex trip-up",
			expectedTimeout: false,
			expectedError:   ErrFailedToAttributePanicToTest,
			outputs: []string{
				"panic: runtime error: invalid memory address or nil pointer dereference\n",
				"[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x21589cc]\n",
				"\n",
				"goroutine 3048 [running]:\n",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*MeteringReport).Message(0x0)\n",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/metering.go:147 +0x6c\n",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.newTestEngine.func4(0x0)\n", // Possible regex trip-up
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine_test.go:230 +0x1e5\n",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).handleStepUpdate(0xc002a8ce08, {0x533a008, 0xc002aee730}, {{0xc001428f00, 0xe}, {0xc000ebef60, 0x22}, {0x4cdb827, 0x9}, 0xc001426338, ...}, ...)\n",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:662 +0x72f\n",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).stepUpdateLoop(0xc002a8ce08, {0x533a008, 0xc002aee730}, {0xc002746540, 0xe}, 0xc00275e000, 0xc00274e4b0)\n",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:545 +0x418\n",
				"created by github.com/smartcontractkit/chainlink/v2/core/services/workflows.(*Engine).resumeInProgressExecutions in goroutine 3314\n",
				"\t/home/runner/work/chainlink/chainlink/core/services/workflows/engine.go:437 +0x511\n",
				"    logger.go:146: 2025-03-21T17:15:55.491Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.491Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.Confirmer\ttxmgr/confirmer.go:265\tFinished CheckForConfirmation\t{\"version\": \"unset@unset\", \"headNum\": 230, \"time\": \"14.767649ms\", \"id\": \"confirmer\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.Confirmer\ttxmgr/confirmer.go:271\tFinished ProcessStuckTransactions\t{\"version\": \"unset@unset\", \"headNum\": 230, \"time\": \"1.543µs\", \"id\": \"confirmer\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.493Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.496Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.496Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.496Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.496Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.496Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.496Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.498Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.499Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.504Z\tDEBUG\tEVM.1337.Txm.Confirmer\ttxmgr/confirmer.go:277\tFinished RebroadcastWhereNecessary\t{\"version\": \"unset@unset\", \"headNum\": 230, \"time\": \"11.026144ms\", \"id\": \"confirmer\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.504Z\tDEBUG\tEVM.1337.Txm.Confirmer\ttxmgr/confirmer.go:278\tprocessHead finish\t{\"version\": \"unset@unset\", \"headNum\": 230, \"id\": \"confirmer\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.506Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 443}\n",
				"    logger.go:146: 2025-03-21T17:15:55.506Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.508Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1099\tUnfinalized log query\t{\"version\": \"unset@unset\", \"logs\": 7, \"currentBlockNumber\": 443, \"blockHash\": \"0xbbb8232c79d104d6da1cd97f9725a44e3fc3dd660a519ba139590f5367ec2b8f\", \"timestamp\": \"2025-03-21T17:21:47.000Z\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.515Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 560}\n",
				"    logger.go:146: 2025-03-21T17:15:55.516Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 560, \"finalized\": 558}\n",
				"    logger.go:146: 2025-03-21T17:15:55.517Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 443}\n",
				"    logger.go:146: 2025-03-21T17:15:55.517Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1099\tUnfinalized log query\t{\"version\": \"unset@unset\", \"logs\": 0, \"currentBlockNumber\": 560, \"blockHash\": \"0xe04f73ebcb6fca77f11860b9f03f0c4dec77ffffa8202366ee1e512155b0e108\", \"timestamp\": \"2025-03-21T17:23:44.000Z\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.518Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.519Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\"]}, \"oracleID\": 2, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.519Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\"]}, \"oracleID\": 0, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.519Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 1, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.519Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 3, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.519Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:167\tChannel is not reportable\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"channelID\": 1, \"err\": \"ChannelID: 1; Reason: ChannelID: 1; Reason: IsReportable=false; not valid yet (observationsTimestampSeconds=1742577354, validAfterSeconds=1742577354)\", \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.519Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:298\tGenerated outcome\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"outcome\": {\"LifeCycleStage\":\"production\",\"ObservationTimestampNanoseconds\":1742577355265352958,\"ChannelDefinitions\":{\"1\":{\"reportFormat\":\"json\",\"streams\":[{\"streamId\":52,\"aggregator\":\"median\"}],\"opts\":null}},\"ValidAfterNanoseconds\":{\"1\":1742577354000000000},\"StreamAggregates\":{\"52\":{\"median\":\"2976.39\"}}}, \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.520Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\"]}, \"oracleID\": 2, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.520Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\"]}, \"oracleID\": 0, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.520Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 1, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.520Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 3, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.520Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:167\tChannel is not reportable\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"channelID\": 1, \"err\": \"ChannelID: 1; Reason: ChannelID: 1; Reason: IsReportable=false; not valid yet (observationsTimestampSeconds=1742577354, validAfterSeconds=1742577354)\", \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.520Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:298\tGenerated outcome\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"outcome\": {\"LifeCycleStage\":\"production\",\"ObservationTimestampNanoseconds\":1742577355265352958,\"ChannelDefinitions\":{\"1\":{\"reportFormat\":\"json\",\"streams\":[{\"streamId\":52,\"aggregator\":\"median\"}],\"opts\":null}},\"ValidAfterNanoseconds\":{\"1\":1742577354000000000},\"StreamAggregates\":{\"52\":{\"median\":\"2976.39\"}}}, \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.521Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\"]}, \"oracleID\": 2, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.521Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\"]}, \"oracleID\": 0, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.521Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 1, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.521Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 3, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.521Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:167\tChannel is not reportable\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"channelID\": 1, \"err\": \"ChannelID: 1; Reason: ChannelID: 1; Reason: IsReportable=false; not valid yet (observationsTimestampSeconds=1742577354, validAfterSeconds=1742577354)\", \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.521Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:298\tGenerated outcome\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"outcome\": {\"LifeCycleStage\":\"production\",\"ObservationTimestampNanoseconds\":1742577355265352958,\"ChannelDefinitions\":{\"1\":{\"reportFormat\":\"json\",\"streams\":[{\"streamId\":52,\"aggregator\":\"median\"}],\"opts\":null}},\"ValidAfterNanoseconds\":{\"1\":1742577354000000000},\"StreamAggregates\":{\"52\":{\"median\":\"2976.39\"}}}, \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1099\tUnfinalized log query\t{\"version\": \"unset@unset\", \"logs\": 7, \"currentBlockNumber\": 443, \"blockHash\": \"0xbbb8232c79d104d6da1cd97f9725a44e3fc3dd660a519ba139590f5367ec2b8f\", \"timestamp\": \"2025-03-21T17:21:47.000Z\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\"]}, \"oracleID\": 2, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\"]}, \"oracleID\": 0, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 1, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:352\tGot observations from peer\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"stage\": \"Outcome\", \"sv\": {\"52\":[\"2976.39\",\"2976.39\",\"2976.39\",\"2976.39\"]}, \"oracleID\": 3, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:167\tChannel is not reportable\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"channelID\": 1, \"err\": \"ChannelID: 1; Reason: ChannelID: 1; Reason: IsReportable=false; not valid yet (observationsTimestampSeconds=1742577354, validAfterSeconds=1742577354)\", \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.522Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_outcome.go:298\tGenerated outcome\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"outcome\": {\"LifeCycleStage\":\"production\",\"ObservationTimestampNanoseconds\":1742577355265352958,\"ChannelDefinitions\":{\"1\":{\"reportFormat\":\"json\",\"streams\":[{\"streamId\":52,\"aggregator\":\"median\"}],\"opts\":null}},\"ValidAfterNanoseconds\":{\"1\":1742577354000000000},\"StreamAggregates\":{\"52\":{\"median\":\"2976.39\"}}}, \"stage\": \"Outcome\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.528Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_reports.go:51\tReportable channels\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"reportableChannels\": [1], \"unreportableChannels\": null, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.528Z\tDEBUG\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1.ReportingPlugin\tllo/plugin_reports.go:72\tEmitting report\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"channelID\": 1, \"report\": {\"ConfigDigest\":\"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\",\"SeqNr\":78,\"ChannelID\":1,\"ValidAfterNanoseconds\":1742577354000000000,\"ObservationTimestampNanoseconds\":1742577355265352958,\"Values\":[\"2976.39\"],\"Specimen\":false}, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.529Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_reports.go:51\tReportable channels\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"reportableChannels\": [1], \"unreportableChannels\": null, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.529Z\tDEBUG\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1.ReportingPlugin\tllo/plugin_reports.go:72\tEmitting report\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"channelID\": 1, \"report\": {\"ConfigDigest\":\"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\",\"SeqNr\":78,\"ChannelID\":1,\"ValidAfterNanoseconds\":1742577354000000000,\"ObservationTimestampNanoseconds\":1742577355265352958,\"Values\":[\"2976.39\"],\"Specimen\":false}, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.528Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_reports.go:51\tReportable channels\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"reportableChannels\": [1], \"unreportableChannels\": null, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.532Z\tDEBUG\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1.ReportingPlugin\tllo/plugin_reports.go:72\tEmitting report\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"channelID\": 1, \"report\": {\"ConfigDigest\":\"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\",\"SeqNr\":78,\"ChannelID\":1,\"ValidAfterNanoseconds\":1742577354000000000,\"ObservationTimestampNanoseconds\":1742577355265352958,\"Values\":[\"2976.39\"],\"Specimen\":false}, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.529Z\tDEBUG\toracle_streams_3.EVM.1337.Relayer.job-3.LLO-888333\tllo/transmitter.go:138\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"report\": {\"Report\":\"eyJDb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwiU2VxTnIiOjc4LCJDaGFubmVsSUQiOjEsIlZhbGlkQWZ0ZXJOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NDAwMDAwMDAwMCwiT2JzZXJ2YXRpb25UaW1lc3RhbXBOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NTI2NTM1Mjk1OCwiVmFsdWVzIjpbeyJ0IjowLCJ2IjoiMjk3Ni4zOSJ9XSwiU3BlY2ltZW4iOmZhbHNlfQ==\",\"Info\":{\"LifeCycleStage\":\"production\",\"ReportFormat\":\"json\"}}, \"sigs\": [{\"Signature\":\"CT9+T7PVUZ8Al7MZir9fQOdUPmZndjrYlnhaZSaoBcga+lfRMCYbkeiMpUiv8Jt1d9DUZrgsdm8T",
				"nGKP1EmWAQE=\",\"Signer\":2},{\"Signature\":\"USh7s2xt+5M5OqCHm86lgeGj8g+4dq597bvWeXj4hiJri7Nvohgf4jBTqxzQhFlrdqkST1ysYbhvpkDXkIWhEwE=\",\"Signer\":3}]}\n",
				"    logger.go:146: 2025-03-21T17:15:55.529Z\tDEBUG\toracle_streams_2.EVM.1337.Relayer.job-3.LLO-888333\tllo/transmitter.go:138\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"report\": {\"Report\":\"eyJDb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwiU2VxTnIiOjc4LCJDaGFubmVsSUQiOjEsIlZhbGlkQWZ0ZXJOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NDAwMDAwMDAwMCwiT2JzZXJ2YXRpb25UaW1lc3RhbXBOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NTI2NTM1Mjk1OCwiVmFsdWVzIjpbeyJ0IjowLCJ2IjoiMjk3Ni4zOSJ9XSwiU3BlY2ltZW4iOmZhbHNlfQ==\",\"Info\":{\"LifeCycleStage\":\"production\",\"ReportFormat\":\"json\"}}, \"sigs\": [{\"Signature\":\"CT9+T7PVUZ8Al7MZir9fQOdUPmZndjrYlnhaZSaoBcga+lfRMCYbkeiMpUiv8Jt1d9DUZrgsdm8T",
				"nGKP1EmWAQE=\",\"Signer\":2},{\"Signature\":\"USh7s2xt+5M5OqCHm86lgeGj8g+4dq597bvWeXj4hiJri7Nvohgf4jBTqxzQhFlrdqkST1ysYbhvpkDXkIWhEwE=\",\"Signer\":3}]}\n",
				"    logger.go:146: 2025-03-21T17:15:55.532Z\tDEBUG\toracle_streams_1.EVM.1337.Relayer.job-3.LLO-888333\tllo/transmitter.go:138\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"report\": {\"Report\":\"eyJDb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwiU2VxTnIiOjc4LCJDaGFubmVsSUQiOjEsIlZhbGlkQWZ0ZXJOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NDAwMDAwMDAwMCwiT2JzZXJ2YXRpb25UaW1lc3RhbXBOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NTI2NTM1Mjk1OCwiVmFsdWVzIjpbeyJ0IjowLCJ2IjoiMjk3Ni4zOSJ9XSwiU3BlY2ltZW4iOmZhbHNlfQ==\",\"Info\":{\"LifeCycleStage\":\"production\",\"ReportFormat\":\"json\"}}, \"sigs\": [{\"Signature\":\"2xZ13fSjHHLkkZFPL1qsOR6uzGdLBM3QmnaWy97LP71wjtECYly7BCxcFXLDY4BjsTO/LojDFmmq",
				"Ts0IIeR6LgA=\",\"Signer\":1},{\"Signature\":\"CT9+T7PVUZ8Al7MZir9fQOdUPmZndjrYlnhaZSaoBcga+lfRMCYbkeiMpUiv8Jt1d9DUZrgsdm8TnGKP1EmWAQE=\",\"Signer\":2}]}\n",
				"    logger.go:146: 2025-03-21T17:15:55.533Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561}\n",
				"    logger.go:146: 2025-03-21T17:15:55.533Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_reports.go:51\tReportable channels\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"reportableChannels\": [1], \"unreportableChannels\": null, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.533Z\tDEBUG\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1.ReportingPlugin\tllo/plugin_reports.go:72\tEmitting report\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"lloProtocolVersion\": 0, \"lifeCycleStage\": \"production\", \"channelID\": 1, \"report\": {\"ConfigDigest\":\"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\",\"SeqNr\":78,\"ChannelID\":1,\"ValidAfterNanoseconds\":1742577354000000000,\"ObservationTimestampNanoseconds\":1742577355265352958,\"Values\":[\"2976.39\"],\"Specimen\":false}, \"stage\": \"Report\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.533Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 560, \"finalized\": 558}\n",
				"    logger.go:146: 2025-03-21T17:15:55.533Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561, \"latestBlockNumber\": 560}\n",
				"    logger.go:146: 2025-03-21T17:15:55.534Z\tDEBUG\toracle_streams_3.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.LLOMercuryTransmitter\tmercurytransmitter/transmitter.go:315\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"reportFormat\": \"json\", \"reportLifeCycleStage\": \"production\", \"transmissionHash\": \"773e4fcec57212c299c183b740e86ca6026439c6b2c7f159c33e901e9b5ca37c\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.534Z\tINFO\toracle_streams_3.OCR2.offchainreporting2.9a9f0afb-b437-4026-8133-f52c4bc053e9.LLO.1\tllo/suppressed_logger.go:51\t🚀 successfully invoked ContractTransmitter.Transmit\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"oid\": 3, \"configDigest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"proto\": \"transmission\", \"seqNr\": 78, \"index\": 0}\n",
				"    logger.go:146: 2025-03-21T17:15:55.534Z\tDEBUG\toracle_streams_1.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.LLOMercuryTransmitter\tmercurytransmitter/transmitter.go:315\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"reportFormat\": \"json\", \"reportLifeCycleStage\": \"production\", \"transmissionHash\": \"33d4af4eef92950ac4483c2ceae526d2e5a4a1cc0b5b9fe7074e819a3b2b474f\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.534Z\tINFO\toracle_streams_1.OCR2.offchainreporting2.5c86b148-b15a-4541-82b9-ae0079f35304.LLO.1\tllo/suppressed_logger.go:51\t🚀 successfully invoked ContractTransmitter.Transmit\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"oid\": 1, \"proto\": \"transmission\", \"configDigest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"index\": 0, \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.534Z\tDEBUG\toracle_streams_0.EVM.1337.Relayer.job-3.LLO-888333\tllo/transmitter.go:138\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"report\": {\"Report\":\"eyJDb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwiU2VxTnIiOjc4LCJDaGFubmVsSUQiOjEsIlZhbGlkQWZ0ZXJOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NDAwMDAwMDAwMCwiT2JzZXJ2YXRpb25UaW1lc3RhbXBOYW5vc2Vjb25kcyI6MTc0MjU3NzM1NTI2NTM1Mjk1OCwiVmFsdWVzIjpbeyJ0IjowLCJ2IjoiMjk3Ni4zOSJ9XSwiU3BlY2ltZW4iOmZhbHNlfQ==\",\"Info\":{\"LifeCycleStage\":\"production\",\"ReportFormat\":\"json\"}}, \"sigs\": [{\"Signature\":\"CT9+T7PVUZ8Al7MZir9fQOdUPmZndjrYlnhaZSaoBcga+lfRMCYbkeiMpUiv8Jt1d9DUZrgsdm8T",
				"nGKP1EmWAQE=\",\"Signer\":2},{\"Signature\":\"USh7s2xt+5M5OqCHm86lgeGj8g+4dq597bvWeXj4hiJri7Nvohgf4jBTqxzQhFlrdqkST1ysYbhvpkDXkIWhEwE=\",\"Signer\":3}]}\n",
				"    logger.go:146: 2025-03-21T17:15:55.534Z\tDEBUG\toracle_streams_3.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.\"127.0.0.1:46557\"\tmercurytransmitter/server.go:231\tTransmit report success\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"b3567e847b4b38827a78d4c289aa559674bee38859064ded8dce99aa1a2f99ce\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"serverURL\": \"127.0.0.1:46557\", \"req.Payload\": \"eyJjb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwic2VxTnIiOjc4LCJyZXBvcnQiOnsiQ29uZmlnRGlnZXN0IjoiMDAwOTQ2Y2E3OTJiM2U3OWFjZTA4MTdkOWI5YWQ5ZTM5N2RiZGMxOGFmOGVlOWM1YmI3NTFiNTg2M2RkYjg2ZCIsIlNlcU5yIjo3OCwiQ2hhbm5lbElEIjoxLCJWYWxpZEFmdGVyTmFub3NlY29uZHMiOjE3NDI1NzczNTQwMDAwMDAwMDAsIk9ic2VydmF0aW9uVGltZXN0YW1wTmFub3NlY29uZHMiOjE3NDI1NzczNTUyNjUzNTI5NTgsIlZhbHVlcyI6W3sidCI6MCwidiI6IjI5NzYuMzkifV0sIlNwZWNpbWVuIjpmYWxzZX0sInNpZ3MiOlt7IlNpZ25hdHVyZSI6IkNUOS",
				"tUN1BWVVo4QWw3TVppcjlmUU9kVVBtWm5kanJZbG5oYVpTYW9CY2dhK2xmUk1DWWJrZWlNcFVpdjhKdDFkOURVWnJnc2RtOFRuR0tQMUVtV0FRRT0iLCJTaWduZXIiOjJ9LHsiU2lnbmF0dXJlIjoiVVNoN3MyeHQrNU01T3FDSG04NmxnZUdqOGcrNGRxNTk3YnZXZVhqNGhpSnJpN052b2hnZjRqQlRxeHpRaEZscmRxa1NUMXlzWWJodnBrRFhrSVdoRXdFPSIsIlNpZ25lciI6M31dfQ==\", \"req.ReportFormat\": 2}\n",
				"    logger.go:146: 2025-03-21T17:15:55.535Z\tDEBUG\toracle_streams_1.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.\"127.0.0.1:46557\"\tmercurytransmitter/server.go:231\tTransmit report success\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"43e4d609de4ea5422f3796de3874abe56e755c2b6ca575c05707f6a341402bf9\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"serverURL\": \"127.0.0.1:46557\", \"req.Payload\": \"eyJjb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwic2VxTnIiOjc4LCJyZXBvcnQiOnsiQ29uZmlnRGlnZXN0IjoiMDAwOTQ2Y2E3OTJiM2U3OWFjZTA4MTdkOWI5YWQ5ZTM5N2RiZGMxOGFmOGVlOWM1YmI3NTFiNTg2M2RkYjg2ZCIsIlNlcU5yIjo3OCwiQ2hhbm5lbElEIjoxLCJWYWxpZEFmdGVyTmFub3NlY29uZHMiOjE3NDI1NzczNTQwMDAwMDAwMDAsIk9ic2VydmF0aW9uVGltZXN0YW1wTmFub3NlY29uZHMiOjE3NDI1NzczNTUyNjUzNTI5NTgsIlZhbHVlcyI6W3sidCI6MCwidiI6IjI5NzYuMzkifV0sIlNwZWNpbWVuIjpmYWxzZX0sInNpZ3MiOlt7IlNpZ25hdHVyZSI6IjJ4Wj",
				"EzZlNqSEhMa2taRlBMMXFzT1I2dXpHZExCTTNRbW5hV3k5N0xQNzF3anRFQ1lseTdCQ3hjRlhMRFk0QmpzVE8vTG9qREZtbXFUczBJSWVSNkxnQT0iLCJTaWduZXIiOjF9LHsiU2lnbmF0dXJlIjoiQ1Q5K1Q3UFZVWjhBbDdNWmlyOWZRT2RVUG1abmRqcllsbmhhWlNhb0JjZ2ErbGZSTUNZYmtlaU1wVWl2OEp0MWQ5RFVacmdzZG04VG5HS1AxRW1XQVFFPSIsIlNpZ25lciI6Mn1dfQ==\", \"req.ReportFormat\": 2}\n",
				"    logger.go:146: 2025-03-21T17:15:55.535Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561}\n",
				"    logger.go:146: 2025-03-21T17:15:55.536Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 560, \"finalized\": 558}\n",
				"    logger.go:146: 2025-03-21T17:15:55.536Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561, \"latestBlockNumber\": 560}\n",
				"    logger.go:146: 2025-03-21T17:15:55.536Z\tDEBUG\toracle_streams_0.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.LLOMercuryTransmitter\tmercurytransmitter/transmitter.go:315\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"reportFormat\": \"json\", \"reportLifeCycleStage\": \"production\", \"transmissionHash\": \"773e4fcec57212c299c183b740e86ca6026439c6b2c7f159c33e901e9b5ca37c\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.537Z\tINFO\toracle_streams_0.OCR2.offchainreporting2.a01a923c-0104-4a46-a904-14bc4abc096a.LLO.1\tllo/suppressed_logger.go:51\t🚀 successfully invoked ContractTransmitter.Transmit\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"index\": 0, \"configDigest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"oid\": 0, \"proto\": \"transmission\", \"seqNr\": 78}\n",
				"    logger.go:146: 2025-03-21T17:15:55.537Z\tDEBUG\toracle_streams_0.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.\"127.0.0.1:46557\"\tmercurytransmitter/server.go:231\tTransmit report success\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"68ab5d04d12c8d2127639d2e32c294e8e849fa9b608f5ad0a650bc0e7386d448\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"serverURL\": \"127.0.0.1:46557\", \"req.Payload\": \"eyJjb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwic2VxTnIiOjc4LCJyZXBvcnQiOnsiQ29uZmlnRGlnZXN0IjoiMDAwOTQ2Y2E3OTJiM2U3OWFjZTA4MTdkOWI5YWQ5ZTM5N2RiZGMxOGFmOGVlOWM1YmI3NTFiNTg2M2RkYjg2ZCIsIlNlcU5yIjo3OCwiQ2hhbm5lbElEIjoxLCJWYWxpZEFmdGVyTmFub3NlY29uZHMiOjE3NDI1NzczNTQwMDAwMDAwMDAsIk9ic2VydmF0aW9uVGltZXN0YW1wTmFub3NlY29uZHMiOjE3NDI1NzczNTUyNjUzNTI5NTgsIlZhbHVlcyI6W3sidCI6MCwidiI6IjI5NzYuMzkifV0sIlNwZWNpbWVuIjpmYWxzZX0sInNpZ3MiOlt7IlNpZ25hdHVyZSI6IkNUOS",
				"tUN1BWVVo4QWw3TVppcjlmUU9kVVBtWm5kanJZbG5oYVpTYW9CY2dhK2xmUk1DWWJrZWlNcFVpdjhKdDFkOURVWnJnc2RtOFRuR0tQMUVtV0FRRT0iLCJTaWduZXIiOjJ9LHsiU2lnbmF0dXJlIjoiVVNoN3MyeHQrNU01T3FDSG04NmxnZUdqOGcrNGRxNTk3YnZXZVhqNGhpSnJpN052b2hnZjRqQlRxeHpRaEZscmRxa1NUMXlzWWJodnBrRFhrSVdoRXdFPSIsIlNpZ25lciI6M31dfQ==\", \"req.ReportFormat\": 2}\n",
				"    logger.go:146: 2025-03-21T17:15:55.539Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444, \"latestBlockNumber\": 443}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tDEBUG\toracle_streams_2.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.LLOMercuryTransmitter\tmercurytransmitter/transmitter.go:315\tTransmit report\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"digest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"seqNr\": 78, \"reportFormat\": \"json\", \"reportLifeCycleStage\": \"production\", \"transmissionHash\": \"773e4fcec57212c299c183b740e86ca6026439c6b2c7f159c33e901e9b5ca37c\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tINFO\toracle_streams_2.OCR2.offchainreporting2.a2742fc2-23d8-408d-9bce-d78b539b9f44.LLO.1\tllo/suppressed_logger.go:51\t🚀 successfully invoked ContractTransmitter.Transmit\t{\"version\": \"unset@unset\", \"jobID\": 3, \"jobName\": \"feed-1\", \"contractID\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"evmChainID\": \"1337\", \"donID\": 888333, \"channelDefinitionsContractAddress\": \"0xE278738AaB5aA4Cb17F16Ada3D197A2FdE7D935c\", \"instanceType\": \"Green\", \"proto\": \"transmission\", \"configDigest\": \"000946ca792b3e79ace0817d9b9ad9e397dbdc18af8ee9c5bb751b5863ddb86d\", \"oid\": 2, \"seqNr\": 78, \"index\": 0}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.540Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444, \"latestBlockNumber\": 443}\n",
				"    logger.go:146: 2025-03-21T17:15:55.541Z\tDEBUG\toracle_streams_2.EVM.1337.Relayer.job-3.LLO-888333.LLOMercuryTransmitter.\"127.0.0.1:46557\"\tmercurytransmitter/server.go:231\tTransmit report success\t{\"version\": \"unset@unset\", \"evmChainID\": \"1337\", \"donID\": 888333, \"transmitterID\": \"2e0b415eb4d97f389ff6d6c33eaadf0cc4613171ebdc59d1e87f91539985a7ce\", \"configMode\": \"bluegreen\", \"configuratorAddress\": \"0xc78dbd2D4bfCE2fDA728461C5f1b67222a4031B6\", \"donID\": 888333, \"serverURL\": \"127.0.0.1:46557\", \"req.Payload\": \"eyJjb25maWdEaWdlc3QiOiIwMDA5NDZjYTc5MmIzZTc5YWNlMDgxN2Q5YjlhZDllMzk3ZGJkYzE4YWY4ZWU5YzViYjc1MWI1ODYzZGRiODZkIiwic2VxTnIiOjc4LCJyZXBvcnQiOnsiQ29uZmlnRGlnZXN0IjoiMDAwOTQ2Y2E3OTJiM2U3OWFjZTA4MTdkOWI5YWQ5ZTM5N2RiZGMxOGFmOGVlOWM1YmI3NTFiNTg2M2RkYjg2ZCIsIlNlcU5yIjo3OCwiQ2hhbm5lbElEIjoxLCJWYWxpZEFmdGVyTmFub3NlY29uZHMiOjE3NDI1NzczNTQwMDAwMDAwMDAsIk9ic2VydmF0aW9uVGltZXN0YW1wTmFub3NlY29uZHMiOjE3NDI1NzczNTUyNjUzNTI5NTgsIlZhbHVlcyI6W3sidCI6MCwidiI6IjI5NzYuMzkifV0sIlNwZWNpbWVuIjpmYWxzZX0sInNpZ3MiOlt7IlNpZ25hdHVyZSI6IkNUOS",
				"tUN1BWVVo4QWw3TVppcjlmUU9kVVBtWm5kanJZbG5oYVpTYW9CY2dhK2xmUk1DWWJrZWlNcFVpdjhKdDFkOURVWnJnc2RtOFRuR0tQMUVtV0FRRT0iLCJTaWduZXIiOjJ9LHsiU2lnbmF0dXJlIjoiVVNoN3MyeHQrNU01T3FDSG04NmxnZUdqOGcrNGRxNTk3YnZXZVhqNGhpSnJpN052b2hnZjRqQlRxeHpRaEZscmRxa1NUMXlzWWJodnBrRFhrSVdoRXdFPSIsIlNpZ25lciI6M31dfQ==\", \"req.ReportFormat\": 2}\n",
				"    logger.go:146: 2025-03-21T17:15:55.548Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561}\n",
				"    logger.go:146: 2025-03-21T17:15:55.549Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 560, \"finalized\": 558}\n",
				"    logger.go:146: 2025-03-21T17:15:55.549Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561, \"latestBlockNumber\": 560}\n",
				"    logger.go:146: 2025-03-21T17:15:55.549Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561}\n",
				"    logger.go:146: 2025-03-21T17:15:55.550Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 560, \"finalized\": 558}\n",
				"    logger.go:146: 2025-03-21T17:15:55.550Z\tDEBUG\tEVM.1000.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 561, \"latestBlockNumber\": 560}\n",
				"    logger.go:146: 2025-03-21T17:15:55.558Z\tWARN\tsyncer/handler.go:769\tworkflow spec not found\t{\"version\": \"unset@unset\", \"workflowID\": \"004b077cb5debdd46c7fcedb10b182e1f89880e9891f7e308563dd3bdb08b85b\"}\n",
				"--- PASS: Test_workflowDeletedHandler/success_deleting_non-existing_workflow_spec (2.97s)\n",
				"--- PASS: Test_workflowDeletedHandler (6.48s)\n",
				"=== RUN   Test_workflowPausedActivatedUpdatedHandler\n",
				"=== RUN   Test_workflowPausedActivatedUpdatedHandler/success_pausing_activating_and_updating_existing_engine_and_spec\n",
				"    logger.go:146: 2025-03-21T17:15:55.579Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444}\n",
				"    logger.go:146: 2025-03-21T17:15:55.580Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.580Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444, \"latestBlockNumber\": 443}\n",
				"    logger.go:146: 2025-03-21T17:15:55.589Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.589Z\tDEBUG\tEVM.1337.Txm.TxmStore.TxmStore\tlogger/logger.go:199\tNew logger: TxmStore\t{\"version\": \"unset@unset\"}\n",
				"    logger.go:146: 2025-03-21T17:15:55.599Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444}\n",
				"    logger.go:146: 2025-03-21T17:15:55.600Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.600Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444, \"latestBlockNumber\": 443}\n",
				"    logger.go:146: 2025-03-21T17:15:55.611Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1035\tPolling for logs\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444}\n",
				"    logger.go:146: 2025-03-21T17:15:55.611Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1141\tLatest blocks read from chain\t{\"version\": \"unset@unset\", \"latest\": 443, \"finalized\": 441}\n",
				"    logger.go:146: 2025-03-21T17:15:55.611Z\tDEBUG\tEVM.1337.LogPoller\tlogpoller/log_poller.go:1047\tNo new blocks since last poll\t{\"version\": \"unset@unset\", \"currentBlockNumber\": 444, \"latestBlockNumber\": 443}\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testName, timeout, err := attributePanicToTest(tc.outputs)
			assert.Equal(t, tc.expectedTimeout, timeout, "timeout flag mismatch")
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError, "error mismatch")
			assert.Empty(t, testName, "test name should be empty")
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
		ProjectPath: "./",
		Verbose:     true,
		RunCount:    1,
		SelectTests: []string{"TestFail"}, // This test is known to fail consistently
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
		ProjectPath: "./",
		Verbose:     true,
		RunCount:    1,
		SelectTests: []string{"TestSkipped"}, // Known skipping test
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

var (
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
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "    testing.go:1399: race detected during execution of test\n"},
	}
)
