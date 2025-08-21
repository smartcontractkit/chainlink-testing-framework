//go:build !integration_tests
// +build !integration_tests

package runner_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/executor"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/parser"
)

type mockExecutor struct {
	RunTestPackageFn func(cfg executor.Config, packageName string, runIndex int) (outputFilePath string, passed bool, err error)
	RunCmdFn         func(cfg executor.Config, testCmd []string, runIndex int) (outputFilePath string, passed bool, err error)

	RunTestPackageCalls []executor.Config
	RunCmdCalls         [][]string
}

func (m *mockExecutor) RunTestPackage(cfg executor.Config, packageName string, runIndex int) (string, bool, error) {
	m.RunTestPackageCalls = append(m.RunTestPackageCalls, cfg)
	if m.RunTestPackageFn != nil {
		return m.RunTestPackageFn(cfg, packageName, runIndex)
	}
	return fmt.Sprintf("mock_output_%s_%d.json", packageName, runIndex), true, nil
}

func (m *mockExecutor) RunCmd(cfg executor.Config, testCmd []string, runIndex int) (string, bool, error) {
	m.RunCmdCalls = append(m.RunCmdCalls, testCmd)
	if m.RunCmdFn != nil {
		return m.RunCmdFn(cfg, testCmd, runIndex)
	}
	return fmt.Sprintf("mock_cmd_output_%d.json", runIndex), true, nil
}

type mockParser struct {
	ParseFilesFn func(rawFilePaths []string, runPrefix string, expectedRuns int, cfg parser.Config) ([]reports.TestResult, []string, error)

	ParseFilesCalls [][]string
	LastParseCfg    parser.Config
}

func (m *mockParser) ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, cfg parser.Config) ([]reports.TestResult, []string, error) {
	m.ParseFilesCalls = append(m.ParseFilesCalls, rawFilePaths)
	m.LastParseCfg = cfg
	if m.ParseFilesFn != nil {
		return m.ParseFilesFn(rawFilePaths, runPrefix, expectedRuns, cfg)
	}
	return []reports.TestResult{{TestName: "DefaultMockTest"}}, rawFilePaths, nil
}

func TestRunner_RunTestPackages(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		runCount          int
		failFast          bool
		packages          []string
		executorResponses map[string]struct {
			passed bool
			err    error
		}
		expectedExecCalls int
		expectedParseArgs struct {
			fileCount int
			cfg       parser.Config
		}
		expectedResultCount int
		expectedError       bool
	}{
		{
			name:     "Happy path - 2 runs, 2 packages",
			runCount: 2,
			failFast: false,
			packages: []string{"pkgA", "pkgB"},
			executorResponses: map[string]struct {
				passed bool
				err    error
			}{
				"pkgA-0": {passed: true, err: nil},
				"pkgA-1": {passed: true, err: nil},
				"pkgB-0": {passed: true, err: nil},
				"pkgB-1": {passed: true, err: nil},
			},
			expectedExecCalls: 4,
			expectedParseArgs: struct {
				fileCount int
				cfg       parser.Config
			}{fileCount: 4, cfg: parser.Config{IgnoreParentFailuresOnSubtests: false, OmitOutputsOnSuccess: false}},
			expectedResultCount: 1,
			expectedError:       false,
		},
		{
			name:     "FailFast stops execution",
			runCount: 5,
			failFast: true,
			packages: []string{"pkgA", "pkgB"},
			executorResponses: map[string]struct {
				passed bool
				err    error
			}{
				"pkgA-0": {passed: false, err: nil},
			},
			expectedExecCalls: 1,
			expectedParseArgs: struct {
				fileCount int
				cfg       parser.Config
			}{fileCount: 1, cfg: parser.Config{IgnoreParentFailuresOnSubtests: false, OmitOutputsOnSuccess: true}},
			expectedResultCount: 1,
			expectedError:       false,
		},
		{
			name:     "Executor error stops execution",
			runCount: 3,
			failFast: false,
			packages: []string{"pkgA", "pkgB"},
			executorResponses: map[string]struct {
				passed bool
				err    error
			}{
				"pkgA-0": {passed: true, err: nil},
				"pkgA-1": {passed: true, err: fmt.Errorf("executor boom")},
			},
			expectedExecCalls:   2,
			expectedResultCount: 0,
			expectedError:       true,
		},
		{
			name:     "Parser error propagated",
			runCount: 1,
			failFast: false,
			packages: []string{"pkgA"},
			executorResponses: map[string]struct {
				passed bool
				err    error
			}{
				"pkgA-0": {passed: true, err: nil},
			},
			expectedExecCalls:   1,
			expectedResultCount: 0,
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockExec := &mockExecutor{
				RunTestPackageFn: func(cfg executor.Config, pkg string, idx int) (string, bool, error) {
					key := fmt.Sprintf("%s-%d", pkg, idx)
					resp, ok := tc.executorResponses[key]
					if !ok {
						return fmt.Sprintf("mock_%s_%d.json", pkg, idx), true, nil
					}
					return fmt.Sprintf("mock_%s_%d.json", pkg, idx), resp.passed, resp.err
				},
			}
			mockParse := &mockParser{}
			if tc.name == "Parser error propagated" {
				mockParse.ParseFilesFn = func(_ []string, _ string, _ int, _ parser.Config) ([]reports.TestResult, []string, error) {
					return nil, nil, fmt.Errorf("parser failed")
				}
			}

			r := runner.NewRunner(
				".",
				false, // Verbose
				tc.runCount,
				nil,   // goTestCountFlag
				false, // goTestRaceFlag
				"",    // goTestTimeoutFlag
				nil,   // tags
				false, // useShuffle
				"",    // shuffleSeed
				tc.failFast,
				nil, // skipTests
				nil, // selectTests
				tc.expectedParseArgs.cfg.IgnoreParentFailuresOnSubtests,
				tc.expectedParseArgs.cfg.OmitOutputsOnSuccess,
				"",
				mockExec,
				mockParse,
			)

			actualResults, err := r.RunTestPackages(tc.packages)

			assert.Len(t, mockExec.RunTestPackageCalls, tc.expectedExecCalls, "Unexpected number of executor calls")

			if tc.expectedError {
				assert.Error(t, err)
				if tc.name == "Executor error stops execution" {
					assert.Len(t, mockParse.ParseFilesCalls, 0, "Parser should not be called on executor error")
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, mockParse.ParseFilesCalls, 1, "Parser should be called once on success/parser error")
				if len(mockParse.ParseFilesCalls) > 0 {
					assert.Len(t, mockParse.ParseFilesCalls[0], tc.expectedParseArgs.fileCount, "Parser called with wrong number of files")
					assert.Equal(t, tc.expectedParseArgs.cfg, mockParse.LastParseCfg, "Parser called with wrong config")
				}
				assert.Len(t, actualResults, tc.expectedResultCount, "Unexpected number of results returned")
			}
		})
	}
}

func TestRunner_RunTestCmd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		runCount          int
		failFast          bool
		cmd               []string
		executorResponses []struct {
			passed bool
			err    error
		}
		expectedExecCalls int
		expectedParseArgs struct {
			fileCount int
			cfg       parser.Config
		}
		expectedResultCount int
		expectedError       bool
	}{
		{
			name:     "Happy path - 2 runs",
			runCount: 2,
			failFast: false,
			cmd:      []string{"go", "test", "./..."},
			executorResponses: []struct {
				passed bool
				err    error
			}{
				{passed: true, err: nil},
				{passed: true, err: nil},
			},
			expectedExecCalls: 2,
			expectedParseArgs: struct {
				fileCount int
				cfg       parser.Config
			}{fileCount: 2, cfg: parser.Config{IgnoreParentFailuresOnSubtests: false, OmitOutputsOnSuccess: false}},
			expectedResultCount: 1,
			expectedError:       false,
		},
		{
			name:     "FailFast stops execution",
			runCount: 3,
			failFast: true,
			cmd:      []string{"go", "test", "./..."},
			executorResponses: []struct {
				passed bool
				err    error
			}{
				{passed: false, err: nil}, // Fails on first run
			},
			expectedExecCalls: 1,
			expectedParseArgs: struct {
				fileCount int
				cfg       parser.Config
			}{fileCount: 1, cfg: parser.Config{IgnoreParentFailuresOnSubtests: false, OmitOutputsOnSuccess: true}},
			expectedResultCount: 1,
			expectedError:       false,
		},
		{
			name:     "Executor error stops execution",
			runCount: 3,
			failFast: false,
			cmd:      []string{"go", "test", "./..."},
			executorResponses: []struct {
				passed bool
				err    error
			}{
				{passed: true, err: nil},
				{passed: false, err: fmt.Errorf("exec boom")}, // Error on second run
			},
			expectedExecCalls:   2,
			expectedResultCount: 0,
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			execCallCount := 0
			mockExec := &mockExecutor{
				RunCmdFn: func(cfg executor.Config, cmd []string, idx int) (string, bool, error) {
					execCallCount++
					if idx < len(tc.executorResponses) {
						resp := tc.executorResponses[idx]
						return fmt.Sprintf("mock_cmd_%d.json", idx), resp.passed, resp.err
					}
					return fmt.Sprintf("mock_cmd_%d.json", idx), true, nil
				},
			}
			mockParse := &mockParser{}

			r := runner.NewRunner(
				".", tc.failFast, // other fields default/nil
				tc.runCount, nil, false, "", nil, false, "", tc.failFast, nil, nil,
				tc.expectedParseArgs.cfg.IgnoreParentFailuresOnSubtests,
				tc.expectedParseArgs.cfg.OmitOutputsOnSuccess,
				"",
				mockExec, mockParse,
			)

			actualResults, err := r.RunTestCmd(tc.cmd)

			assert.Equal(t, tc.expectedExecCalls, execCallCount, "Unexpected number of executor RunCmd calls")

			if tc.expectedError {
				assert.Error(t, err)
				assert.Len(t, mockParse.ParseFilesCalls, 0, "Parser should not be called on executor error")
			} else {
				assert.NoError(t, err)
				assert.Len(t, mockParse.ParseFilesCalls, 1, "Parser should be called once on success")
				if len(mockParse.ParseFilesCalls) > 0 {
					assert.Len(t, mockParse.ParseFilesCalls[0], tc.expectedParseArgs.fileCount, "Parser called with wrong number of files")
					assert.Equal(t, tc.expectedParseArgs.cfg, mockParse.LastParseCfg, "Parser called with wrong config")
				}
				assert.Len(t, actualResults, tc.expectedResultCount, "Unexpected number of results returned")
			}
		})
	}
}

func TestRunner_RerunFailedTests(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		initialFailedTests []reports.TestResult
		rerunCount         int
		executorResponses  map[string]struct {
			passed bool
			err    error
		}
		expectedExecCalls int
		expectedParseArgs struct {
			fileCount int
			cfg       parser.Config
		}
		expectedFinalResultCount int
		expectedError            bool
	}{
		{
			name: "Rerun successful",
			initialFailedTests: []reports.TestResult{
				{TestName: "TestFailA", TestPackage: "pkgA", Failures: 1, Runs: 1},
				{TestName: "TestFailB", TestPackage: "pkgB", Failures: 1, Runs: 1},
			},
			rerunCount: 2,
			executorResponses: map[string]struct {
				passed bool
				err    error
			}{
				"pkgA-0": {passed: true, err: nil},
				"pkgB-0": {passed: false, err: nil},
				"pkgB-1": {passed: true, err: nil},
			},
			expectedExecCalls: 3,
			expectedParseArgs: struct {
				fileCount int
				cfg       parser.Config
			}{fileCount: 3, cfg: parser.Config{OmitOutputsOnSuccess: false}},
			expectedFinalResultCount: 1,
			expectedError:            false,
		},
		{
			name:               "No failed tests to rerun",
			initialFailedTests: []reports.TestResult{},
			rerunCount:         3,
			expectedExecCalls:  0,
			expectedParseArgs: struct {
				fileCount int
				cfg       parser.Config
			}{fileCount: 0},
			expectedError: false,
		},
		{
			name: "Executor error during rerun",
			initialFailedTests: []reports.TestResult{
				{TestName: "TestFailA", TestPackage: "pkgA", Failures: 1, Runs: 1},
			},
			rerunCount: 1,
			executorResponses: map[string]struct {
				passed bool
				err    error
			}{
				"pkgA-0": {passed: false, err: fmt.Errorf("exec rerun boom")},
			},
			expectedExecCalls: 1,
			expectedError:     true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			execCallCount := 0
			mockExec := &mockExecutor{
				RunTestPackageFn: func(cfg executor.Config, pkg string, idx int) (string, bool, error) {
					execCallCount++
					key := fmt.Sprintf("%s-%d", pkg, idx)
					resp, ok := tc.executorResponses[key]
					if !ok {
						return fmt.Sprintf("mock_rerun_%s_%d.json", pkg, idx), true, nil
					}
					// Check if config forces count=1 and has the right -run pattern
					assert.NotNil(t, cfg.GoTestCountFlag, "Rerun should force GoTestCountFlag")
					if cfg.GoTestCountFlag != nil {
						assert.Equal(t, 1, *cfg.GoTestCountFlag, "Rerun should force count=1")
					}
					assert.Nil(t, cfg.SkipTests, "Rerun should clear SkipTests")
					require.Len(t, cfg.SelectTests, 1, "Rerun should set exactly one SelectTests pattern")
					// Basic check if pattern looks right (contains test names from input)
					for _, failedTest := range tc.initialFailedTests {
						if failedTest.TestPackage == pkg {
							assert.Contains(t, cfg.SelectTests[0], regexp.QuoteMeta(failedTest.TestName))
						}
					}
					return fmt.Sprintf("mock_rerun_%s_%d.json", pkg, idx), resp.passed, resp.err
				},
			}
			mockParse := &mockParser{}

			r := runner.NewRunner(".", false, 0, nil, false, "", nil, false, "", false, nil, nil, false, false,
				"",
				mockExec, mockParse)

			actualResults, _, err := r.RerunFailedTests(tc.initialFailedTests, tc.rerunCount)

			assert.Equal(t, tc.expectedExecCalls, execCallCount, "Unexpected number of executor calls")

			if tc.expectedError {
				assert.Error(t, err)
				assert.Len(t, mockParse.ParseFilesCalls, 0, "Parser should not be called on rerun executor error")
			} else {
				assert.NoError(t, err)
				if tc.expectedExecCalls > 0 {
					assert.Len(t, mockParse.ParseFilesCalls, 1, "Parser should be called once after reruns")
					if len(mockParse.ParseFilesCalls) > 0 {
						assert.Len(t, mockParse.ParseFilesCalls[0], tc.expectedParseArgs.fileCount, "Parser called with wrong number of files")
						assert.Equal(t, r.IgnoreParentFailuresOnSubtests, mockParse.LastParseCfg.IgnoreParentFailuresOnSubtests, "Parser IgnoreParentFailures mismatch")
						assert.Equal(t, r.OmitOutputsOnSuccess, mockParse.LastParseCfg.OmitOutputsOnSuccess, "Parser OmitOutputsOnSuccess mismatch")
					}
					assert.Len(t, actualResults, tc.expectedFinalResultCount, "Unexpected number of results returned from rerun parse")
				} else {
					assert.Len(t, mockParse.ParseFilesCalls, 0, "Parser should not be called if no reruns executed")
					assert.Empty(t, actualResults, "No results expected if no reruns executed")
				}
			}
		})
	}
}
