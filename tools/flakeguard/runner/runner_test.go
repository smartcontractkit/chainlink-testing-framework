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

const (
	defaultRuns          = 5
	flakyTestPackagePath = "./example_test_package"
	debugDir             = "debug_outputs"
)

type expectedTestResult struct {
	TestResult *reports.TestResult
	seen       bool
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name            string
		runner          Runner
		expectedRuns    int
		expectedResults map[string]*expectedTestResult
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
			expectedRuns: defaultRuns,
			expectedResults: map[string]*expectedTestResult{
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
						Skips:     defaultRuns,
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
			expectedRuns: defaultRuns,
			expectedResults: map[string]*expectedTestResult{
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
						Panics:   defaultRuns,
					},
				},
				"TestRace": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
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
			expectedRuns: defaultRuns,
			expectedResults: map[string]*expectedTestResult{
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
		{
			name: "failfast",
			runner: Runner{
				ProjectPath:          "./",
				Verbose:              true,
				RunCount:             defaultRuns,
				UseRace:              false,
				SkipTests:            []string{"TestPanic", "TestFlaky"}, // Flaky test introduces too much variability for failfast
				FailFast:             true,
				SelectedTestPackages: []string{flakyTestPackagePath},
				CollectRawOutput:     true,
			},
			expectedRuns: 1,
			expectedResults: map[string]*expectedTestResult{
				"TestFail": {
					TestResult: &reports.TestResult{
						TestName:  "TestFail",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 0,
						Failures:  1,
					},
				},
				"TestPass": {
					TestResult: &reports.TestResult{
						TestName:  "TestPass",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
						Successes: 1,
					},
				},
				"TestSkipped": {
					TestResult: &reports.TestResult{
						TestName:  "TestSkipped",
						Panicked:  false,
						Skipped:   true,
						PassRatio: 0,
						Skips:     1,
					},
				},
				"TestRace": {
					TestResult: &reports.TestResult{
						TestName:  "TestRace",
						Panicked:  false,
						Skipped:   false,
						PassRatio: 1,
						Successes: 1,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testResults, err := tc.runner.RunTests()
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
					t.Logf("error marshalling test results: %v", err)
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

			for _, result := range testResults {
				t.Run(fmt.Sprintf("checking results of %s", result.TestName), func(t *testing.T) {
					expected, ok := tc.expectedResults[result.TestName]
					require.True(t, ok, "unexpected test result: %s", result.TestName)
					require.False(t, expected.seen, "test '%s' was seen multiple times", result.TestName)
					expected.seen = true

					assert.Len(t, result.Durations, result.Runs, "test '%s' has a mismatch of runs and duration counts", result.TestName, defaultRuns)
					resultCounts := result.Successes + result.Failures + result.Panics + result.Skips
					assert.Equal(t, result.Runs, resultCounts,
						"test '%s' doesn't match Runs count with results counts\nRuns: %d\nSuccesses: %d\nFailures: %d\nPanics: %d\nSkips: %d\nTotal: %d",
						result.TestName, result.Runs, result.Successes, result.Failures, result.Panics, result.Skips, resultCounts,
					)
					assert.Equal(t, expected.TestResult.Panicked, result.Panicked, "test '%s' had an unexpected panic result", result.TestName)
					assert.Equal(t, expected.TestResult.Skipped, result.Skipped, "test '%s' had an unexpected skipped result", result.TestName)
					assert.Equal(t, expected.TestResult.Panics, result.Panics, "test '%s' had an unexpected number of panics", result.TestName)
					assert.Equal(t, expected.TestResult.Skips, result.Skips, "test '%s' had an unexpected number of skips", result.TestName)

					if tc.name == "panic" {
						assert.True(t, result.PackagePanicked, "test '%s' should have package panicked", result.TestName)
					} else if result.TestName == "TestFlaky" {
						assert.Greater(t, result.Successes, 0, "flaky test '%s' should have passed some", result.TestName)
						assert.Greater(t, result.Failures, 0, "flaky test '%s' should have failed some", result.TestName)
						assert.Greater(t, result.PassRatio, float64(0), "flaky test '%s' should have a flaky pass ratio", result.TestName)
						assert.Less(t, result.PassRatio, float64(1), "flaky test '%s' should have a flaky pass ratio", result.TestName)
					} else {
						assert.Equal(t, tc.expectedRuns, result.Runs, "test '%s' had an unexpected number of runs", result.TestName)
						assert.Equal(t, expected.TestResult.PassRatio, result.PassRatio, "test '%s' had an unexpected pass ratio", result.TestName)
						assert.Equal(t, expected.TestResult.Successes, result.Successes, "test '%s' had an unexpected number of successes", result.TestName)
						assert.Equal(t, expected.TestResult.Failures, result.Failures, "test '%s' had an unexpected number of failures", result.TestName)
					}
				})
			}

			for _, expected := range tc.expectedResults {
				assert.True(t, expected.seen, "expected test '%s' not found in test runs", expected.TestResult.TestName)
			}
		})
	}
}

func TestAttributePanicToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		packageName      string
		expectedTestName string
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
			name:         "no panic",
			packageName:  "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries: noPanicRaceEntries,
		},
		{
			name:        "empty",
			packageName: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			panicEntries: []entry{
				{},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testName, err := attributePanicToTest(tc.packageName, tc.panicEntries)
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
			name:        "no race",
			packageName: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries: noPanicRaceEntries,
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
	noPanicRaceEntries = []entry{
		{Action: "start", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package"},
		{Action: "run", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass", Output: "=== RUN   TestPass\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass", Output: "=== PAUSE TestPass\n"},
		{Action: "pause", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass"},
		{Action: "run", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail", Output: "=== RUN   TestFail\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail", Output: "=== PAUSE TestFail\n"},
		{Action: "pause", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail"},
		{Action: "run", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "=== RUN   TestFlaky\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "=== PAUSE TestFlaky\n"},
		{Action: "pause", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky"},
		{Action: "run", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped", Output: "=== RUN   TestSkipped\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped", Output: "=== PAUSE TestSkipped\n"},
		{Action: "pause", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped"},
		{Action: "run", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "=== RUN   TestRace\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "=== PAUSE TestRace\n"},
		{Action: "pause", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace"},
		{Action: "cont", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass", Output: "=== CONT  TestPass\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass", Output: "    example_tests_test.go:11: This test always passes\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass", Output: "--- PASS: TestPass (0.00s)\n"},
		{Action: "pass", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestPass", Elapsed: 0},
		{Action: "cont", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail", Output: "=== CONT  TestFail\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail", Output: "    example_tests_test.go:16: This test always fails\n"},
		{Action: "cont", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped", Output: "=== CONT  TestSkipped\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail", Output: "--- FAIL: TestFail (0.00s)\n"},
		{Action: "fail", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFail", Elapsed: 0},
		{Action: "cont", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "=== CONT  TestRace\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "    example_tests_test.go:56: This test should trigger a failure if run with the -race flag, but otherwise pass\n"},
		{Action: "cont", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "=== CONT  TestFlaky\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped", Output: "    example_tests_test.go:46: This test is intentionally skipped\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped", Output: "--- SKIP: TestSkipped (0.00s)\n"},
		{Action: "skip", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestSkipped", Elapsed: 0},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "    example_tests_test.go:80: Final value of sharedCounter: 46430\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "--- PASS: TestRace (0.00s)\n"},
		{Action: "pass", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Elapsed: 0},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "    example_tests_test.go:31: This is a designed flaky test working as intended\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "--- FAIL: TestFlaky (0.00s)\n"},
		{Action: "fail", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Elapsed: 0},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Output: "FAIL\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Output: "FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package\t0.138s\n"},
		{Action: "fail", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Elapsed: 0.138},
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
