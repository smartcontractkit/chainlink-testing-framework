package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttributePanicToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		expectedTestName string
		expectedTimeout  bool
		expectedError    error
		outputs          []string
	}{
		{
			name:             "properly attributed panic",
			expectedTestName: "TestPanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This test intentionally panics [recovered]",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestPanic(0x140000b6ea0?)",
			},
		},
		{
			name:             "improperly attributed panic (but still findable)",
			expectedTestName: "TestPanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This test intentionally panics [recovered]",
				"TestPanic(0x140000b6ea0?)",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestPanic(0x140000b6ea0?)",
			},
		},
		{
			name:             "log after test complete panic",
			expectedTestName: "Test_workflowRegisteredHandler/skips_fetch_if_secrets_url_is_missing",
			expectedTimeout:  false,
			outputs: []string{
				"panic: Log in goroutine after Test_workflowRegisteredHandler/skips_fetch_if_secrets_url_is_missing has completed: ...",
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
			},
		},
		{
			name:             "subtest panic",
			expectedTestName: "TestSubTestsSomePanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This subtest always panics",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestSubTestsSomePanic.func2(0x140000c81a0?)",
			},
		},
		{
			name:             "memory_test panic extraction",
			expectedTestName: "TestJobClientJobAPI",
			expectedTimeout:  false,
			outputs: []string{
				"panic: freeport: cannot allocate port block",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory_test.TestJobClientJobAPI(0xc000683dc0)",
			},
		},
		{
			name:             "changeset_test panic extraction",
			expectedTestName: "TestDeployBalanceReader",
			expectedTimeout:  false,
			outputs: []string{
				"panic: freeport: cannot allocate port block",
				"github.com/smartcontractkit/chainlink/deployment/keystone/changeset_test.TestDeployBalanceReader(0xc000583c00)",
			},
		},
		{
			name:          "empty output",
			expectedError: ErrFailedToAttributePanicToTest,
			outputs:       []string{},
		},
		{
			name:          "no test name in panic",
			expectedError: ErrFailedToAttributePanicToTest,
			outputs: []string{
				"panic: reflect: Elem of invalid type bool",
			},
		},
		{
			name:            "fail to parse timeout duration",
			expectedTimeout: true,
			expectedError:   ErrFailedToParseTimeoutDuration,
			outputs: []string{
				"panic: test timed out after malformedDurationStr\n",
			},
		},
		{
			name:            "timeout panic without obvious culprit",
			expectedTimeout: true,
			expectedError:   ErrDetectedTimeoutFailedAttribution,
			outputs: []string{
				"panic: test timed out after 10m0s\n",
				"\trunning tests:\n",
				"\t\tTestAlmostPanicTime (9m59s)\n",
			},
		},
		{
			name:          "possible regex trip-up (no TestXxx)",
			expectedError: ErrFailedToAttributePanicToTest,
			outputs: []string{
				"panic: runtime error: invalid memory address or nil pointer dereference\n",
				"github.com/smartcontractkit/chainlink/v2/core/services/workflows.newTestEngine.func4(0x0)",
			},
		},
		{
			name:             "Panic with multiple Test names in stack",
			expectedTestName: "TestInner",
			expectedTimeout:  false,
			outputs: []string{
				"panic: Something went wrong in helper",
				"main.helperFunction()",
				"main.TestInner(0xc00...)",
				"main.TestOuter(0xc00...)",
			},
		},
		{
			name:             "Timeout with multiple matching durations",
			expectedTestName: "TestA",
			expectedTimeout:  true,
			outputs: []string{
				"panic: test timed out after 5m0s",
				"running tests:",
				"\tTestA (5m0s)",
				"\tTestB (4m59s)",
				"\tTestC (5m1s)",
				"\tTestD (5m0s)",
			},
		},
		{
			name:            "fail to parse test duration in timeout list",
			expectedTimeout: true,
			expectedError:   ErrDetectedTimeoutFailedParse,
			outputs: []string{
				"panic: test timed out after 10m0s\n",
				"\trunning tests:\n",
				"\t\tTestAddAndPromoteCandidatesForNewChain (malformedDurationStr)\n",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testName, timeout, err := AttributePanicToTest(tc.outputs)

			assert.Equal(t, tc.expectedTimeout, timeout, "Timeout flag mismatch")
			if tc.expectedError != nil {
				require.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tc.expectedError, "Error mismatch")
				assert.Empty(t, testName, "Test name should be empty on error")
			} else {
				require.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tc.expectedTestName, testName, "Test name mismatch")
			}
		})
	}
}

func TestAttributeRaceToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		expectedTestName string
		expectedError    error
		outputs          []string
	}{
		{
			name:             "properly attributed race",
			expectedTestName: "TestRace",
			outputs: []string{
				"WARNING: DATA RACE",
				"  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()",
			},
		},
		{
			name:          "empty output",
			expectedError: ErrFailedToAttributeRaceToTest,
			outputs:       []string{},
		},
		{
			name:          "no test name in race output",
			expectedError: ErrFailedToAttributeRaceToTest,
			outputs: []string{
				"WARNING: DATA RACE",
				"  main.main.func1()",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testName, err := AttributeRaceToTest(tc.outputs)

			if tc.expectedError != nil {
				require.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tc.expectedError, "Error mismatch")
				assert.Empty(t, testName, "Test name should be empty on error")
			} else {
				require.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tc.expectedTestName, testName, "Test name mismatch")
			}
		})
	}
}
