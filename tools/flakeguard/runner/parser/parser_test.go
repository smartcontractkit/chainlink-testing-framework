package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
)

func jsonLine(action, pkg, test, output string, elapsed float64) string {
	fixedTime, _ := time.Parse(time.RFC3339Nano, "2024-01-01T10:00:00.000Z")
	entry := struct {
		Time    time.Time
		Action  string
		Package string
		Test    string  `json:",omitempty"`
		Output  string  `json:",omitempty"`
		Elapsed float64 `json:",omitempty"`
	}{
		Time:    fixedTime,
		Action:  action,
		Package: pkg,
		Test:    test,
		Output:  output,
	}
	if action == "pass" || action == "fail" {
		entry.Elapsed = elapsed
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		panic(fmt.Sprintf("test setup error: failed to marshal jsonLine: %v", err))
	}
	return string(jsonBytes)
}

func buildOutput(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}

func TestParseTestResults_Basic(t *testing.T) {
	t.Parallel()

	pkg1 := "github.com/test/package1"
	pkg2 := "github.com/test/package2"

	testCases := []struct {
		name             string
		inputFiles       map[string]string
		cfg              Config
		expectedResults  map[string]reports.TestResult
		expectedErrorIs  error
		expectedErrorMsg string
	}{
		{
			name: "Single Test Pass",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("run", pkg1, "TestPass", "", 0),
					jsonLine("output", pkg1, "TestPass", "output line 1\n", 0),
					jsonLine("pass", pkg1, "TestPass", "", 1.23),
				),
			},
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/%s", pkg1, "TestPass"): {
					TestName:    "TestPass",
					TestPackage: pkg1,
					Runs:        1,
					Successes:   1,
					PassRatio:   1.0,
				},
			},
		},
		{
			name: "Single Test Fail",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("run", pkg1, "TestFail", "", 0),
					jsonLine("output", pkg1, "TestFail", "fail output\n", 0),
					jsonLine("fail", pkg1, "TestFail", "", 2.34),
				),
			},
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/%s", pkg1, "TestFail"): {
					TestName:    "TestFail",
					TestPackage: pkg1,
					Runs:        1,
					Failures:    1,
					PassRatio:   0.0,
				},
			},
		},
		{
			name: "Single Test Skip",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("run", pkg1, "TestSkip", "", 0),
					jsonLine("output", pkg1, "TestSkip", "skip reason\n", 0),
					jsonLine("skip", pkg1, "TestSkip", "", 0),
				),
			},
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/%s", pkg1, "TestSkip"): {
					TestName:    "TestSkip",
					TestPackage: pkg1,
					Runs:        0,
					Skips:       1,
					Skipped:     true,
					PassRatio:   1.0,
				},
			},
		},
		{
			name: "Mixed Pass Fail Skip Multiple Runs",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("run", pkg1, "TestA", "", 0),
					jsonLine("pass", pkg1, "TestA", "", 1.0),
					jsonLine("run", pkg1, "TestB", "", 0),
					jsonLine("fail", pkg1, "TestB", "", 1.0),
					jsonLine("run", pkg2, "TestC", "", 0),
					jsonLine("skip", pkg2, "TestC", "", 0),
				),
				"run2.json": buildOutput(
					jsonLine("run", pkg1, "TestA", "", 0),
					jsonLine("fail", pkg1, "TestA", "", 1.1), // TestA fails on run 2
					jsonLine("run", pkg1, "TestB", "", 0),
					jsonLine("pass", pkg1, "TestB", "", 1.1), // TestB passes on run 2
					jsonLine("run", pkg2, "TestC", "", 0),
					jsonLine("skip", pkg2, "TestC", "", 0),
				),
			},
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/%s", pkg1, "TestA"): {TestName: "TestA", TestPackage: pkg1, Runs: 2, Successes: 1, Failures: 1, PassRatio: 0.5},
				fmt.Sprintf("%s/%s", pkg1, "TestB"): {TestName: "TestB", TestPackage: pkg1, Runs: 2, Successes: 1, Failures: 1, PassRatio: 0.5},
				fmt.Sprintf("%s/%s", pkg2, "TestC"): {TestName: "TestC", TestPackage: pkg2, Runs: 0, Skips: 2, Skipped: true, PassRatio: 1.0},
			},
		},
		{
			name: "Build Failure",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("build-fail", "", "", "compile error message", 0),
				),
			},
			cfg:             Config{},
			expectedResults: nil,
			expectedErrorIs: ErrBuild,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parser := NewParser().(*defaultParser)

			tempDir := t.TempDir()
			filePaths := make([]string, 0, len(tc.inputFiles))
			for name, content := range tc.inputFiles {
				fpath := filepath.Join(tempDir, name)
				err := os.WriteFile(fpath, []byte(content), 0644)
				require.NoError(t, err, "Failed to write temp file %s", name)
				filePaths = append(filePaths, fpath)
			}

			actualResults, err := parser.parseTestResults(filePaths, "run", len(filePaths), tc.cfg)

			if tc.expectedErrorIs != nil {
				require.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tc.expectedErrorIs, "Error mismatch")
				if tc.expectedErrorMsg != "" {
					assert.ErrorContains(t, err, tc.expectedErrorMsg, "Error message mismatch")
				}
				assert.Nil(t, actualResults, "Results should be nil on error")
			} else {
				require.NoError(t, err, "Expected no error but got: %v", err)
				require.NotNil(t, actualResults, "Results should not be nil on success")
				require.Equal(t, len(tc.expectedResults), len(actualResults), "Unexpected number of results")

				actualResultsMap := make(map[string]reports.TestResult)
				for _, res := range actualResults {
					key := fmt.Sprintf("%s/%s", res.TestPackage, res.TestName)
					actualResultsMap[key] = res
				}

				for key, expected := range tc.expectedResults {
					actual, ok := actualResultsMap[key]
					require.True(t, ok, "Expected result for key '%s' not found", key)

					// Compare relevant fields
					assertResultBasic(t, key, expected, actual)
					if strings.HasSuffix(key, "TestPass") {
						assert.NotEmpty(t, actual.Durations, "TestPass should have duration")
						if !tc.cfg.OmitOutputsOnSuccess {
							assert.Contains(t, actual.PassedOutputs["run1"], "output line 1\n", "TestPass missing expected output")
						}
					} else if strings.HasSuffix(key, "TestFail") {
						assert.NotEmpty(t, actual.Durations, "TestFail should have duration")
						assert.Contains(t, actual.FailedOutputs["run1"], "fail output\n", "TestFail missing expected output")
					} else if strings.HasSuffix(key, "TestSkip") {
						assert.Empty(t, actual.Durations, "TestSkip should have no duration")
						assert.Empty(t, actual.PassedOutputs, "TestSkip should have no passed output")
						assert.Empty(t, actual.FailedOutputs, "TestSkip should have no failed output")
					}
					if expected.TestName == "TestA" || expected.TestName == "TestB" {
						assert.Len(t, actual.Durations, 2, "%s should have 2 durations", expected.TestName)
					}
					if expected.TestName == "TestC" {
						assert.Empty(t, actual.Durations, "TestC should have 0 durations")
					}
				}
			}
		})
	}
}

func TestParseTestResults_OutputHandling(t *testing.T) {
	t.Parallel()

	pkg1 := "github.com/test/outputpkg"

	testCases := []struct {
		name            string
		inputFile       string
		cfg             Config
		expectedPassOut map[string][]string
		expectedFailOut map[string][]string
		expectedPkgOut  []string
	}{
		{
			name: "OmitOutputsOnSuccess=true",
			inputFile: buildOutput(
				jsonLine("run", pkg1, "TestPass", "", 0),
				jsonLine("output", pkg1, "TestPass", "pass output 1", 0),
				jsonLine("pass", pkg1, "TestPass", "", 1.0),
				jsonLine("run", pkg1, "TestFail", "", 0),
				jsonLine("output", pkg1, "TestFail", "fail output 1", 0),
				jsonLine("fail", pkg1, "TestFail", "", 1.0),
				jsonLine("output", pkg1, "", "package output 1", 0),
			),
			cfg:             Config{OmitOutputsOnSuccess: true},
			expectedPassOut: map[string][]string{},
			expectedFailOut: map[string][]string{
				"run1": {"fail output 1"},
			},
			expectedPkgOut: []string{"package output 1"},
		},
		{
			name: "OmitOutputsOnSuccess=false",
			inputFile: buildOutput(
				jsonLine("run", pkg1, "TestPass", "", 0),
				jsonLine("output", pkg1, "TestPass", "pass output 1", 0),
				jsonLine("pass", pkg1, "TestPass", "", 1.0),
				jsonLine("run", pkg1, "TestFail", "", 0),
				jsonLine("output", pkg1, "TestFail", "fail output 1", 0),
				jsonLine("fail", pkg1, "TestFail", "", 1.0),
				jsonLine("output", pkg1, "", "package output 1", 0),
				jsonLine("output", pkg1, "", "package output 2", 0),
			),
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedPassOut: map[string][]string{
				"run1": {"pass output 1"},
			},
			expectedFailOut: map[string][]string{
				"run1": {"fail output 1"},
			},
			expectedPkgOut: []string{"package output 1", "package output 2"},
		},
		{
			name: "No test-specific output",
			inputFile: buildOutput(
				jsonLine("run", pkg1, "TestPass", "", 0),
				jsonLine("pass", pkg1, "TestPass", "", 1.0),
				jsonLine("run", pkg1, "TestFail", "", 0),
				jsonLine("fail", pkg1, "TestFail", "", 1.0),
				jsonLine("output", pkg1, "", "package output only", 0),
			),
			cfg:             Config{OmitOutputsOnSuccess: false},
			expectedPassOut: map[string][]string{},
			expectedFailOut: map[string][]string{
				"run1": {"--- TEST FAILED (no specific output captured) ---"},
			},
			expectedPkgOut: []string{"package output only"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser().(*defaultParser)
			tempDir := t.TempDir()
			fpath := filepath.Join(tempDir, "run1.json")
			err := os.WriteFile(fpath, []byte(tc.inputFile), 0644)
			require.NoError(t, err)

			actualResults, err := parser.parseTestResults([]string{fpath}, "run", 1, tc.cfg)
			require.NoError(t, err)
			require.NotEmpty(t, actualResults)

			passResult := findResult(t, actualResults, "TestPass")
			failResult := findResult(t, actualResults, "TestFail")

			if passResult != nil {
				assert.Equal(t, len(tc.expectedPassOut), len(passResult.PassedOutputs), "PassedOutputs length mismatch for TestPass")
				if len(tc.expectedPassOut) > 0 {
					assert.Equal(t, tc.expectedPassOut["run1"], passResult.PassedOutputs["run1"], "PassedOutputs content mismatch for TestPass")
				}
				assert.Empty(t, passResult.Outputs, "General Outputs map should be empty after processing TestPass")
				assert.Equal(t, tc.expectedPkgOut, passResult.PackageOutputs, "PackageOutputs mismatch for TestPass")
			}

			if failResult != nil {
				assert.Equal(t, len(tc.expectedFailOut), len(failResult.FailedOutputs), "FailedOutputs length mismatch for TestFail")
				if len(tc.expectedFailOut) > 0 {
					assert.Equal(t, tc.expectedFailOut["run1"], failResult.FailedOutputs["run1"], "FailedOutputs content mismatch for TestFail")
				}
				assert.Empty(t, failResult.Outputs, "General Outputs map should be empty after processing TestFail")
				assert.Equal(t, tc.expectedPkgOut, failResult.PackageOutputs, "PackageOutputs mismatch for TestFail")
			}
		})
	}
}

func TestParseTestResults_Subtests(t *testing.T) {
	t.Parallel()

	pkg := "github.com/test/subtestpkg"

	testCases := []struct {
		name            string
		inputFile       string
		cfg             Config
		expectedResults map[string]reports.TestResult
	}{
		{
			name: "Parent and Subtest Pass",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestParent", "", 0),
				jsonLine("run", pkg, "TestParent/SubPass", "", 0),
				jsonLine("output", pkg, "TestParent/SubPass", "sub output", 0),
				jsonLine("pass", pkg, "TestParent/SubPass", "", 0.5),
				jsonLine("output", pkg, "TestParent", "parent output after sub", 0),
				jsonLine("pass", pkg, "TestParent", "", 1.0),
			),
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParent", pkg):         {TestName: "TestParent", TestPackage: pkg, Runs: 1, Successes: 1, PassRatio: 1.0},
				fmt.Sprintf("%s/TestParent/SubPass", pkg): {TestName: "TestParent/SubPass", TestPackage: pkg, Runs: 1, Successes: 1, PassRatio: 1.0},
			},
		},
		{
			name: "Parent Pass, Subtest Fail",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestParent", "", 0),
				jsonLine("run", pkg, "TestParent/SubFail", "", 0),
				jsonLine("output", pkg, "TestParent/SubFail", "sub fail output", 0),
				jsonLine("fail", pkg, "TestParent/SubFail", "", 0.6),
				jsonLine("output", pkg, "TestParent", "parent output after sub fail", 0),
				jsonLine("pass", pkg, "TestParent", "", 1.2),
			),
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParent", pkg):         {TestName: "TestParent", TestPackage: pkg, Runs: 1, Successes: 1, PassRatio: 1.0},
				fmt.Sprintf("%s/TestParent/SubFail", pkg): {TestName: "TestParent/SubFail", TestPackage: pkg, Runs: 1, Failures: 1, PassRatio: 0.0},
			},
		},
		{
			name: "Parent Fail Before Subtest",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestParentFailEarly", "", 0),
				jsonLine("output", pkg, "TestParentFailEarly", "parent fail output", 0),
				jsonLine("fail", pkg, "TestParentFailEarly", "", 0.1),
			),
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParentFailEarly", pkg): {TestName: "TestParentFailEarly", TestPackage: pkg, Runs: 1, Failures: 1, PassRatio: 0.0},
			},
		},
		{
			name: "Parent Fail After Subtest",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestParentFailLate", "", 0),
				jsonLine("run", pkg, "TestParentFailLate/SubPass", "", 0),
				jsonLine("pass", pkg, "TestParentFailLate/SubPass", "", 0.5),
				jsonLine("output", pkg, "TestParentFailLate", "parent fail output later", 0),
				jsonLine("fail", pkg, "TestParentFailLate", "", 1.5),
			),
			cfg: Config{OmitOutputsOnSuccess: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParentFailLate", pkg):         {TestName: "TestParentFailLate", TestPackage: pkg, Runs: 1, Failures: 1, PassRatio: 0.0},
				fmt.Sprintf("%s/TestParentFailLate/SubPass", pkg): {TestName: "TestParentFailLate/SubPass", TestPackage: pkg, Runs: 1, Successes: 1, PassRatio: 1.0},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser().(*defaultParser)
			tempDir := t.TempDir()
			fpath := filepath.Join(tempDir, "run1.json")
			err := os.WriteFile(fpath, []byte(tc.inputFile), 0644)
			require.NoError(t, err)

			actualResults, err := parser.parseTestResults([]string{fpath}, "run", 1, tc.cfg)
			require.NoError(t, err)
			require.Equal(t, len(tc.expectedResults), len(actualResults), "Unexpected number of results")

			actualResultsMap := resultsToMap(actualResults)
			for key, expected := range tc.expectedResults {
				actual, ok := actualResultsMap[key]
				require.True(t, ok, "Expected result for key '%s' not found", key)
				assertResultBasic(t, key, expected, actual)
				if strings.Contains(key, "SubPass") {
					assert.Len(t, actual.Durations, 1, "SubPass should have 1 duration")
					if !tc.cfg.OmitOutputsOnSuccess {
						if tc.name == "Parent and Subtest Pass" {
							require.Contains(t, actual.PassedOutputs, "run1", "PassedOutputs map missing run1 key for %s in %s", key, tc.name)
							assert.Contains(t, actual.PassedOutputs["run1"], "sub output", "SubPass missing expected output in %s", tc.name)
						} else {
							assert.Empty(t, actual.PassedOutputs["run1"], "PassedOutputs[run1] should be empty for %s in %s", key, tc.name)
						}
					}
				} else if strings.Contains(key, "SubFail") {
					assert.Len(t, actual.Durations, 1, "SubFail should have 1 duration")
					require.Contains(t, actual.FailedOutputs, "run1", "FailedOutputs map missing run1 key for %s", key)
					assert.Contains(t, actual.FailedOutputs["run1"], "sub fail output", "SubFail missing expected output")
				}
			}
		})
	}
}

func TestParseTestResults_Durations(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/durationpkg"

	inputFile := buildOutput(
		jsonLine("run", pkg, "TestA", "", 0),
		jsonLine("pass", pkg, "TestA", "", 1.5),
		jsonLine("run", pkg, "TestB", "", 0),
		jsonLine("fail", pkg, "TestB", "", 2.5),
		jsonLine("run", pkg, "TestC", "", 0),
		jsonLine("pass", pkg, "TestC", "", 0),
		jsonLine("run", pkg, "TestD", "", 0),
		jsonLine("skip", pkg, "TestD", "", 0),
	)

	parser := NewParser().(*defaultParser)
	tempDir := t.TempDir()
	fpath := filepath.Join(tempDir, "run1.json")
	err := os.WriteFile(fpath, []byte(inputFile), 0644)
	require.NoError(t, err)

	actualResults, err := parser.parseTestResults([]string{fpath}, "run", 1, Config{})
	require.NoError(t, err)

	resultsMap := resultsToMap(actualResults)

	resA, ok := resultsMap[fmt.Sprintf("%s/TestA", pkg)]
	require.True(t, ok, "TestA not found")
	require.Len(t, resA.Durations, 1, "TestA should have 1 duration")
	assert.Equal(t, int64(1500), resA.Durations[0].Milliseconds(), "TestA duration mismatch")

	resB, ok := resultsMap[fmt.Sprintf("%s/TestB", pkg)]
	require.True(t, ok, "TestB not found")
	require.Len(t, resB.Durations, 1, "TestB should have 1 duration")
	assert.Equal(t, int64(2500), resB.Durations[0].Milliseconds(), "TestB duration mismatch")

	resC, ok := resultsMap[fmt.Sprintf("%s/TestC", pkg)]
	require.True(t, ok, "TestC not found")
	require.Len(t, resC.Durations, 1, "TestC should have 1 duration")
	assert.Equal(t, int64(0), resC.Durations[0].Milliseconds(), "TestC duration mismatch")

	resD, ok := resultsMap[fmt.Sprintf("%s/TestD", pkg)]
	require.True(t, ok, "TestD not found")
	assert.Empty(t, resD.Durations, "TestD should have 0 durations")
}

func TestParseTestResults_PanicRace(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/panicracepkg"

	panicOutput := []string{
		"panic: This test intentionally panics",
		"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestRegularPanic(...)",
	}
	raceOutput := []string{
		"WARNING: DATA RACE",
		"  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestDataRace(...)",
	}
	timeoutOutput := []string{
		"panic: test timed out after 1m0s",
		"running tests:",
		"\tTestTimeoutCulprit (1m0s)",
	}

	testCases := []struct {
		name            string
		inputFile       string
		cfg             Config
		expectedResults map[string]reports.TestResult
	}{
		{
			name: "Regular Panic",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestRegularPanic", "", 0),
				jsonLine("output", pkg, "TestRegularPanic", panicOutput[0], 0),
				jsonLine("output", pkg, "TestRegularPanic", panicOutput[1], 0),
				jsonLine("fail", pkg, "TestRegularPanic", "", 0.5),
			),
			cfg: Config{},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestRegularPanic", pkg): {
					TestName: "TestRegularPanic", TestPackage: pkg, Runs: 1, Failures: 1, Panic: true,
				},
			},
		},
		{
			name: "Data Race",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestDataRace", "", 0),
				jsonLine("output", pkg, "TestDataRace", raceOutput[0], 0),
				jsonLine("output", pkg, "TestDataRace", raceOutput[1], 0),
				jsonLine("fail", pkg, "TestDataRace", "", 0.6),
			),
			cfg: Config{},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestDataRace", pkg): {
					TestName: "TestDataRace", TestPackage: pkg, Runs: 1, Failures: 1, Race: true,
				},
			},
		},
		{
			name: "Timeout Panic",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestTimeoutCulprit", "", 0),
				jsonLine("output", pkg, "TestTimeoutCulprit", timeoutOutput[0], 0),
				jsonLine("output", pkg, "TestTimeoutCulprit", timeoutOutput[1], 0),
				jsonLine("output", pkg, "TestTimeoutCulprit", timeoutOutput[2], 0),
				jsonLine("fail", pkg, "TestTimeoutCulprit", "", 60.1),
			),
			cfg: Config{},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestTimeoutCulprit", pkg): {
					TestName: "TestTimeoutCulprit", TestPackage: pkg, Runs: 1, Failures: 1, Panic: true, Timeout: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser().(*defaultParser)
			tempDir := t.TempDir()
			fpath := filepath.Join(tempDir, "run1.json")
			err := os.WriteFile(fpath, []byte(tc.inputFile), 0644)
			require.NoError(t, err)

			actualResults, err := parser.parseTestResults([]string{fpath}, "run", 1, tc.cfg)
			require.NoError(t, err)
			require.Equal(t, len(tc.expectedResults), len(actualResults), "Unexpected number of results")

			actualResultsMap := resultsToMap(actualResults)
			for key, expected := range tc.expectedResults {
				actual, ok := actualResultsMap[key]
				require.True(t, ok, "Expected result for key '%s' not found", key)
				assertResultBasic(t, key, expected, actual)
				assert.Equal(t, expected.Panic, actual.Panic, "Panic flag mismatch for %s", key)
				assert.Equal(t, expected.Race, actual.Race, "Race flag mismatch for %s", key)
				assert.Equal(t, expected.Timeout, actual.Timeout, "Timeout flag mismatch for %s", key)

				// Check if panic/race output was added to FailedOutputs
				if expected.Panic || expected.Race {
					outputs, ok := actual.FailedOutputs["run1"]
					require.True(t, ok, "FailedOutputs map missing run1 key for %s", key)
					require.NotEmpty(t, outputs, "FailedOutputs should contain panic/race info for %s", key)

					if expected.Panic {
						assert.Contains(t, outputs[0], "PANIC DETECTED", "Missing PANIC marker for %s", key)
						if tc.name == "Regular Panic" {
							assert.Contains(t, outputs, panicOutput[0])
							assert.Contains(t, outputs, panicOutput[1])
						} else if tc.name == "Timeout Panic" {
							assert.Contains(t, outputs, timeoutOutput[0])
							assert.Contains(t, outputs, timeoutOutput[1])
							assert.Contains(t, outputs, timeoutOutput[2])
						}
					} else if expected.Race {
						assert.Contains(t, outputs[0], "RACE DETECTED", "Missing RACE marker for %s", key)
						assert.Contains(t, outputs, raceOutput[0])
						assert.Contains(t, outputs, raceOutput[1])
					}
				}
			}
		})
	}
}

func TestParseTestResults_RunCountCorrection(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/runcountpkg"

	// Simulate a panic happening, which often causes a 'fail' event without a preceding 'pass'
	// for the same test in the same run, leading to potential overcounting if not handled.
	// This example simulates 2 expected runs, but the panic causes 3 fail events for TestA.
	inputFileRun1 := buildOutput(
		jsonLine("run", pkg, "TestA", "", 0),
		jsonLine("output", pkg, "TestA", "panic: Error in TestA", 0),
		jsonLine("output", pkg, "TestA", "github.com/test/runcountpkg.TestA(...)", 0),
		jsonLine("fail", pkg, "TestA", "", 0.1), // Fail from panic
		jsonLine("run", pkg, "TestB", "", 0),
		jsonLine("pass", pkg, "TestB", "", 0.2),
	)
	inputFileRun2 := buildOutput(
		jsonLine("run", pkg, "TestA", "", 0),
		jsonLine("pass", pkg, "TestA", "", 1.1), // Passes on run 2
		jsonLine("run", pkg, "TestB", "", 0),
		jsonLine("pass", pkg, "TestB", "", 1.2),
	)

	parser := NewParser().(*defaultParser)
	tempDir := t.TempDir()
	filePaths := []string{
		filepath.Join(tempDir, "run1.json"),
		filepath.Join(tempDir, "run2.json"),
	}
	err := os.WriteFile(filePaths[0], []byte(inputFileRun1), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filePaths[1], []byte(inputFileRun2), 0644)
	require.NoError(t, err)

	// Pass 2 as totalExpectedRunsPerTest
	actualResults, err := parser.parseTestResults(filePaths, "run", 2, Config{})
	require.NoError(t, err)

	resultsMap := resultsToMap(actualResults)

	resA, ok := resultsMap[fmt.Sprintf("%s/TestA", pkg)]
	require.True(t, ok, "TestA not found")
	// Final runs = 2 (processed run1 fail, processed run2 pass)
	assert.Equal(t, 2, resA.Runs, "TestA Runs should be 2")
	assert.Equal(t, 1, resA.Successes, "TestA Successes should be 1")
	assert.Equal(t, 1, resA.Failures, "TestA Failures should be 1")
	assert.True(t, resA.Panic, "TestA should be marked panicked") // Panic flag from attribution
	assert.InDelta(t, 0.5, resA.PassRatio, 0.001, "TestA PassRatio mismatch")

	resB, ok := resultsMap[fmt.Sprintf("%s/TestB", pkg)]
	require.True(t, ok, "TestB not found")
	// Final runs = 2 (processed run1 pass, processed run2 pass)
	assert.Equal(t, 2, resB.Runs, "TestB Runs should be 2")
	assert.Equal(t, 2, resB.Successes, "TestB Successes should be 2")
	assert.Equal(t, 0, resB.Failures, "TestB Failures should be 0")
	assert.False(t, resB.Panic, "TestB should not be panicked")
	assert.Equal(t, 1.0, resB.PassRatio, "TestB PassRatio mismatch")
}

func TestParseTestResults_RunCountCorrectionRefined(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/runcountpkg2"

	testCases := []struct {
		name                string
		inputFiles          map[string]string
		expectedTotalRuns   int
		expectedResultTestA reports.TestResult
	}{
		{
			name: "Panic within expected runs",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("output", pkg, "TestA", "panic: Error", 0), // Panic
					jsonLine("output", pkg, "TestA", "github.com/test/pkg.TestA(...)", 0),
					jsonLine("fail", pkg, "TestA", "", 0.1),
				),
				"run2.json": buildOutput(
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("pass", pkg, "TestA", "", 0.2),
				),
				"run3.json": buildOutput(
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("pass", pkg, "TestA", "", 0.3),
				),
			},
			expectedTotalRuns: 3,
			expectedResultTestA: reports.TestResult{ // Runs=3, Success=2, Fail=1, Panic=true
				TestName: "TestA", TestPackage: pkg, Runs: 3, Successes: 2, Failures: 1, Panic: true, PassRatio: 2.0 / 3.0,
			},
		},
		{
			name: "Panic exceeding expected runs (capped)",
			inputFiles: map[string]string{
				"run1.json": buildOutput( // This run fails due to panic
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("output", pkg, "TestA", "panic: Error", 0),
					jsonLine("output", pkg, "TestA", "github.com/test/pkg.TestA(...)", 0),
					jsonLine("fail", pkg, "TestA", "", 0.1),
					jsonLine("fail", pkg, "TestA", "", 0.11),
				),
				"run2.json": buildOutput(
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("pass", pkg, "TestA", "", 0.2),
				),
			},
			expectedTotalRuns: 2, // Only expected 2 runs total
			expectedResultTestA: reports.TestResult{ // Expect correction: Runs=2, Success=1, Fail=1, Panic=true
				TestName: "TestA", TestPackage: pkg, Runs: 2, Successes: 1, Failures: 1, Panic: true, PassRatio: 0.5,
			},
		},
		{
			name: "Normal overcount (no panic/race, capped)",
			inputFiles: map[string]string{
				"run1.json": buildOutput(
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("pass", pkg, "TestA", "", 0.1),
					jsonLine("pass", pkg, "TestA", "", 0.11),
				),
				"run2.json": buildOutput(
					jsonLine("run", pkg, "TestA", "", 0),
					jsonLine("fail", pkg, "TestA", "", 0.2),
				),
			},
			expectedTotalRuns: 2,
			expectedResultTestA: reports.TestResult{
				TestName: "TestA", TestPackage: pkg, Runs: 2, Successes: 1, Failures: 1, Panic: false, PassRatio: 0.5,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser().(*defaultParser)
			tempDir := t.TempDir()
			filePaths := make([]string, 0, len(tc.inputFiles))
			for name, content := range tc.inputFiles {
				fpath := filepath.Join(tempDir, name)
				err := os.WriteFile(fpath, []byte(content), 0644)
				require.NoError(t, err)
				filePaths = append(filePaths, fpath)
			}

			actualResults, err := parser.parseTestResults(filePaths, "run", tc.expectedTotalRuns, Config{})
			require.NoError(t, err)

			resultsMap := resultsToMap(actualResults)
			actualA, ok := resultsMap[fmt.Sprintf("%s/TestA", pkg)]
			require.True(t, ok, "TestA not found")
			assertResultBasic(t, "TestA", tc.expectedResultTestA, actualA)
			assert.Equal(t, tc.expectedResultTestA.Panic, actualA.Panic, "TestA Panic mismatch")

		})
	}
}

func TestParseTestResults_PanicInheritance(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/panicinheritpkg"

	parentPanicInput := buildOutput(
		jsonLine("run", pkg, "TestParentPanics", "", 0),
		jsonLine("run", pkg, "TestParentPanics/SubPass", "", 0),
		jsonLine("pass", pkg, "TestParentPanics/SubPass", "", 0.1),
		jsonLine("run", pkg, "TestParentPanics/SubFail", "", 0),
		jsonLine("fail", pkg, "TestParentPanics/SubFail", "", 0.2),
		jsonLine("output", pkg, "TestParentPanics", "panic: Parent panics here!", 0),
		jsonLine("output", pkg, "TestParentPanics", "github.com/test/panicinheritpkg.TestParentPanics(...)", 0),
		jsonLine("fail", pkg, "TestParentPanics", "", 0.3),
	)

	testCases := []struct {
		name            string
		cfg             Config
		expectedResults map[string]reports.TestResult
	}{
		{
			name: "Inheritance Enabled (Default)",
			cfg:  Config{IgnoreParentFailuresOnSubtests: false},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParentPanics", pkg):         {TestName: "TestParentPanics", TestPackage: pkg, Runs: 1, Failures: 1, Panic: true},
				fmt.Sprintf("%s/TestParentPanics/SubPass", pkg): {TestName: "TestParentPanics/SubPass", TestPackage: pkg, Runs: 1, Successes: 0, Failures: 1, Panic: true},
				fmt.Sprintf("%s/TestParentPanics/SubFail", pkg): {TestName: "TestParentPanics/SubFail", TestPackage: pkg, Runs: 1, Failures: 1, Panic: true},
			},
		},
		{
			name: "IgnoreParentFailures (No Transform Effect Here)",
			cfg:  Config{IgnoreParentFailuresOnSubtests: true},
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParentPanics", pkg):         {TestName: "TestParentPanics", TestPackage: pkg, Runs: 1, Failures: 1, Panic: true},
				fmt.Sprintf("%s/TestParentPanics/SubPass", pkg): {TestName: "TestParentPanics/SubPass", TestPackage: pkg, Runs: 1, Successes: 0, Failures: 1, Panic: true},
				fmt.Sprintf("%s/TestParentPanics/SubFail", pkg): {TestName: "TestParentPanics/SubFail", TestPackage: pkg, Runs: 1, Failures: 1, Panic: true},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser().(*defaultParser)
			tempDir := t.TempDir()
			fpath := filepath.Join(tempDir, "run1.json")
			err := os.WriteFile(fpath, []byte(parentPanicInput), 0644)
			require.NoError(t, err)

			actualResults, err := parser.parseTestResults([]string{fpath}, "run", 1, tc.cfg)
			require.NoError(t, err)
			require.Equal(t, len(tc.expectedResults), len(actualResults), "Unexpected number of results")

			actualResultsMap := resultsToMap(actualResults)
			for key, expected := range tc.expectedResults {
				actual, ok := actualResultsMap[key]
				require.True(t, ok, "Expected result for key '%s' not found", key)
				assertResultBasic(t, key, expected, actual)
				assert.Equal(t, expected.Panic, actual.Panic, "Panic flag mismatch for %s", key)
			}
		})
	}
}

func TestParseTestResults_JSONErrors(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/jsonerrpkg"

	inputFile := strings.Join([]string{
		jsonLine("run", pkg, "TestBeforeError", "", 0),                                   // Valid line
		jsonLine("pass", pkg, "TestBeforeError", "", 1.0),                                // Valid line
		`{"Action":"run","Package":"github.com/test/jsonerrpkg","Test":"TestWithError"}`, // Missing fields
		`this is not json`, // Invalid line
		jsonLine("run", pkg, "TestAfterError", "", 0),    // Valid line
		jsonLine("pass", pkg, "TestAfterError", "", 1.0), // Valid line
	}, "\n") + "\n"

	parser := NewParser().(*defaultParser)
	tempDir := t.TempDir()
	fpath := filepath.Join(tempDir, "run1.json")
	err := os.WriteFile(fpath, []byte(inputFile), 0644)
	require.NoError(t, err)

	actualResults, err := parser.parseTestResults([]string{fpath}, "run", 1, Config{})
	require.NoError(t, err, "Parsing should continue despite invalid JSON lines")

	require.Len(t, actualResults, 2, "Expected results only from tests with terminal actions")
	resultsMap := resultsToMap(actualResults)

	resBefore, okBefore := resultsMap[fmt.Sprintf("%s/TestBeforeError", pkg)]
	assert.True(t, okBefore, "TestBeforeError should be parsed")
	assert.Equal(t, 1, resBefore.Runs, "TestBeforeError Runs mismatch")
	assert.Equal(t, 1, resBefore.Successes, "TestBeforeError Successes mismatch")

	resAfter, okAfter := resultsMap[fmt.Sprintf("%s/TestAfterError", pkg)]
	assert.True(t, okAfter, "TestAfterError should be parsed")
	assert.Equal(t, 1, resAfter.Runs, "TestAfterError Runs mismatch")
	assert.Equal(t, 1, resAfter.Successes, "TestAfterError Successes mismatch")

	_, okMid := resultsMap[fmt.Sprintf("%s/TestWithError", pkg)]
	assert.False(t, okMid, "TestWithError should not be in final results")
}

func TestParseFiles_Transformation(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/transformpkg"

	inputFile := buildOutput(
		jsonLine("run", pkg, "TestParentTransform", "", 0),
		jsonLine("output", pkg, "TestParentTransform", "parent output", 0),
		jsonLine("run", pkg, "TestParentTransform/SubFail", "", 0),
		jsonLine("output", pkg, "TestParentTransform/SubFail", "sub fail output", 0),
		jsonLine("fail", pkg, "TestParentTransform/SubFail", "", 0.1),
		jsonLine("fail", pkg, "TestParentTransform", "", 0.2),
	)

	parser := NewParser()
	tempDir := t.TempDir()
	fpath := filepath.Join(tempDir, "run1.json")
	err := os.WriteFile(fpath, []byte(inputFile), 0644)
	require.NoError(t, err)

	cfg := Config{IgnoreParentFailuresOnSubtests: true, OmitOutputsOnSuccess: false}
	actualResults, _, err := parser.ParseFiles([]string{fpath}, "run", 1, cfg)
	require.NoError(t, err)

	require.Len(t, actualResults, 2, "Expected 2 results after transformation")
	resultsMap := resultsToMap(actualResults)

	parentRes, okP := resultsMap[fmt.Sprintf("%s/TestParentTransform", pkg)]
	require.True(t, okP, "Parent test not found")
	assert.Equal(t, 1, parentRes.Runs, "Parent Runs mismatch")
	assert.Equal(t, 1, parentRes.Successes, "Parent Successes mismatch (should pass)")
	assert.Equal(t, 0, parentRes.Failures, "Parent Failures mismatch (should pass)")
	assert.Equal(t, 1.0, parentRes.PassRatio, "Parent PassRatio mismatch")
	assert.False(t, parentRes.Panic, "Parent Panic mismatch")
	require.Contains(t, parentRes.PassedOutputs, "run1", "Parent PassedOutputs missing run1")
	assert.Contains(t, parentRes.PassedOutputs["run1"], "parent output", "Parent output missing from PassedOutputs")
	assert.NotContains(t, parentRes.PassedOutputs["run1"][0], "=== PASS", "Parent output should not be transformed unless original contained FAIL markers")

	subRes, okS := resultsMap[fmt.Sprintf("%s/TestParentTransform/SubFail", pkg)]
	require.True(t, okS, "Subtest not found")
	assert.Equal(t, 1, subRes.Runs, "Subtest Runs mismatch")
	assert.Equal(t, 0, subRes.Successes, "Subtest Successes mismatch")
	assert.Equal(t, 1, subRes.Failures, "Subtest Failures mismatch")
	assert.Equal(t, 0.0, subRes.PassRatio, "Subtest PassRatio mismatch")
}

func TestParseTestResults_EmptyOrIncomplete(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/empty"

	testCases := []struct {
		name          string
		inputFiles    map[string]string
		numExpResults int
		expError      bool
	}{
		{
			name:          "Empty File",
			inputFiles:    map[string]string{"run1.json": ""},
			numExpResults: 0,
			expError:      false,
		},
		{
			name:          "Only Run Action",
			inputFiles:    map[string]string{"run1.json": buildOutput(jsonLine("run", pkg, "TestOnlyRun", "", 0))},
			numExpResults: 0,
			expError:      false,
		},
		{
			name:          "Run and Output Only",
			inputFiles:    map[string]string{"run1.json": buildOutput(jsonLine("run", pkg, "TestRunOutput", "", 0), jsonLine("output", pkg, "TestRunOutput", "out", 0))},
			numExpResults: 0,
			expError:      false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser().(*defaultParser)
			tempDir := t.TempDir()
			filePaths := make([]string, 0, len(tc.inputFiles))
			for name, content := range tc.inputFiles {
				fpath := filepath.Join(tempDir, name)
				err := os.WriteFile(fpath, []byte(content), 0644)
				require.NoError(t, err)
				filePaths = append(filePaths, fpath)
			}

			actualResults, err := parser.parseTestResults(filePaths, "run", 1, Config{})

			if tc.expError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, actualResults, tc.numExpResults)
			}
		})
	}
}

// findResult finds a specific test result by name from a slice.
func findResult(t *testing.T, results []reports.TestResult, testName string) *reports.TestResult {
	t.Helper()
	for i := range results {
		if results[i].TestName == testName {
			return &results[i]
		}
	}
	return nil
}

// resultsToMap converts a slice of results to a map keyed by "package/testName".
func resultsToMap(results []reports.TestResult) map[string]reports.TestResult {
	m := make(map[string]reports.TestResult, len(results))
	for _, res := range results {
		key := fmt.Sprintf("%s/%s", res.TestPackage, res.TestName)
		m[key] = res
	}
	return m
}

// assertResultBasic performs basic assertions on core result fields.
func assertResultBasic(t *testing.T, key string, expected, actual reports.TestResult) {
	t.Helper()
	assert.Equal(t, expected.TestName, actual.TestName, "TestName mismatch for %s", key)
	assert.Equal(t, expected.TestPackage, actual.TestPackage, "TestPackage mismatch for %s", key)
	assert.Equal(t, expected.Runs, actual.Runs, "Runs mismatch for %s", key)
	assert.Equal(t, expected.Successes, actual.Successes, "Successes mismatch for %s", key)
	assert.Equal(t, expected.Failures, actual.Failures, "Failures mismatch for %s", key)
	assert.Equal(t, expected.Skips, actual.Skips, "Skips mismatch for %s", key)
	assert.Equal(t, expected.Skipped, actual.Skipped, "Skipped flag mismatch for %s", key)
	assert.InDelta(t, expected.PassRatio, actual.PassRatio, 0.001, "PassRatio mismatch for %s", key)
}

// TestParseFiles_WithTransformationScenarios verifies the interaction between the parser and the transformer.
func TestParseFiles_IgnoreParentFailures(t *testing.T) {
	t.Parallel()
	pkg := "github.com/test/transformpkg"

	testCases := []struct {
		name            string
		inputFile       string
		expectedResults map[string]reports.TestResult
	}{
		{
			name: "Parent only fails due to subtest",
			inputFile: buildOutput(
				jsonLine("run", pkg, "TestParentTransform", "", 0),
				jsonLine("output", pkg, "TestParentTransform", "parent setup output", 0), // Regular output
				jsonLine("run", pkg, "TestParentTransform/SubFail", "", 0),
				jsonLine("output", pkg, "TestParentTransform/SubFail", "sub fail output", 0),
				jsonLine("output", pkg, "TestParentTransform/SubFail", "--- FAIL: TestParentTransform/SubFail (0.1s)", 0), // Subtest fail marker
				jsonLine("fail", pkg, "TestParentTransform/SubFail", "", 0.1),
				jsonLine("output", pkg, "TestParentTransform", "--- FAIL: TestParentTransform (0.2s)", 0), // Parent fail marker (due to subtest)
				jsonLine("fail", pkg, "TestParentTransform", "", 0.2),
			),
			expectedResults: map[string]reports.TestResult{
				fmt.Sprintf("%s/TestParentTransform", pkg):         {TestName: "TestParentTransform", TestPackage: pkg, Runs: 1, Successes: 1, Failures: 0, PassRatio: 1.0},
				fmt.Sprintf("%s/TestParentTransform/SubFail", pkg): {TestName: "TestParentTransform/SubFail", TestPackage: pkg, Runs: 1, Successes: 0, Failures: 1, PassRatio: 0.0},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser()
			tempDir := t.TempDir()
			fpath := filepath.Join(tempDir, "run1.json")
			err := os.WriteFile(fpath, []byte(tc.inputFile), 0644)
			require.NoError(t, err)

			cfg := Config{IgnoreParentFailuresOnSubtests: true, OmitOutputsOnSuccess: false}
			actualResults, _, err := parser.ParseFiles([]string{fpath}, "run", 1, cfg)
			require.NoError(t, err)

			require.Equal(t, len(tc.expectedResults), len(actualResults), "Unexpected number of results")
			actualResultsMap := resultsToMap(actualResults)
			for key, expected := range tc.expectedResults {
				actual, ok := actualResultsMap[key]
				require.True(t, ok, "Expected result for key '%s' not found", key)
				assertResultBasic(t, key, expected, actual)
				if expected.TestName == "TestParentTransform" {
					require.Contains(t, actual.PassedOutputs, "run1", "PassedOutputs missing run1 for transformed parent")
					assert.Contains(t, actual.PassedOutputs["run1"], "parent setup output", "Parent original output missing")
					assert.Contains(t, actual.PassedOutputs["run1"], "--- PASS: TestParentTransform (0.2s)", "Parent output marker not transformed to PASS")
				}
				if expected.TestName == "TestParentTransform/SubFail" {
					require.Contains(t, actual.FailedOutputs, "run1", "FailedOutputs missing run1 for failing subtest %s", key)
					assert.Contains(t, actual.FailedOutputs["run1"], "sub fail output", "Subtest fail output missing for %s", key)
					assert.Contains(t, actual.FailedOutputs["run1"], "--- FAIL: TestParentTransform/SubFail (0.1s)", "Subtest fail marker missing for %s", key)
				}
			}
		})
	}
}
