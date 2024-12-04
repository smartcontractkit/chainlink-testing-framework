package reports

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterFailedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.5, Skipped: false},
		{TestName: "Test2", PassRatio: 0.9, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true}, // Skipped test
	}

	failedTests := FilterFailedTests(results, 0.6)
	expected := []string{"Test1", "Test3"}

	require.Equal(t, len(expected), len(failedTests), "not as many failed tests as expected")

	for i, test := range failedTests {
		assert.Equal(t, expected[i], test.TestName, "wrong test name")
	}
}

func TestFilterPassedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true}, // Skipped test
	}

	passedTests := FilterPassedTests(results, 0.6)
	expected := []string{"Test1", "Test2"}

	require.Equal(t, len(expected), len(passedTests), "not as many passed tests as expected")

	for i, test := range passedTests {
		assert.Equal(t, expected[i], test.TestName, "wrong test name")
	}
}

func TestFilterSkippedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: true},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true},
	}

	skippedTests := FilterSkippedTests(results)
	expected := []string{"Test2", "Test4"}

	require.Equal(t, len(expected), len(skippedTests), "not as many skipped tests as expected")

	for i, test := range skippedTests {
		assert.Equal(t, expected[i], test.TestName, "wrong test name")
	}
}

func TestPrintTests(t *testing.T) {
	testcases := []struct {
		name                   string
		testResults            []TestResult
		maxPassRatio           float64
		expectedRuns           int
		expectedPasses         int
		expectedFails          int
		expectedSkippedTests   int
		expectedPanickedTests  int
		expectedRacedTests     int
		expectedFlakyTests     int
		expectedStringsContain []string
	}{
		{
			name: "single flaky test",
			testResults: []TestResult{
				{
					TestName:    "Test1",
					TestPackage: "package1",
					PassRatio:   0.75,
					Successes:   3,
					Failures:    1,
					Skipped:     false,
					Runs:        4,
					Durations:   []time.Duration{time.Millisecond * 1200, time.Millisecond * 900, time.Millisecond * 1100, time.Second},
				},
			},
			maxPassRatio:           1.0,
			expectedRuns:           4,
			expectedPasses:         3,
			expectedFails:          1,
			expectedSkippedTests:   0,
			expectedPanickedTests:  0,
			expectedRacedTests:     0,
			expectedFlakyTests:     1,
			expectedStringsContain: []string{"Test1", "package1", "75.00%", "false", "1.05s", "4", "0"},
		},
		{
			name: "multiple passing tests",
			testResults: []TestResult{
				{
					TestName:    "Test1",
					TestPackage: "package1",
					PassRatio:   1.0,
					Skipped:     false,
					Successes:   4,
					Runs:        4,
					Durations:   []time.Duration{time.Millisecond * 1200, time.Millisecond * 900, time.Millisecond * 1100, time.Second},
				},
				{
					TestName:    "Test2",
					TestPackage: "package1",
					PassRatio:   1.0,
					Skipped:     false,
					Successes:   4,
					Runs:        4,
					Durations:   []time.Duration{time.Millisecond * 1200, time.Millisecond * 900, time.Millisecond * 1100, time.Second},
				},
			},
			maxPassRatio:           1.0,
			expectedRuns:           8,
			expectedPasses:         8,
			expectedFails:          0,
			expectedSkippedTests:   0,
			expectedPanickedTests:  0,
			expectedRacedTests:     0,
			expectedFlakyTests:     0,
			expectedStringsContain: []string{},
		},
	}

	for _, testCase := range testcases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer

			runs, passes, fails, skips, panickedTests, racedTests, flakyTests := PrintTests(&buf, tc.testResults, tc.maxPassRatio, false)
			assert.Equal(t, tc.expectedRuns, runs, "wrong number of runs")
			assert.Equal(t, tc.expectedPasses, passes, "wrong number of passes")
			assert.Equal(t, tc.expectedFails, fails, "wrong number of failures")
			assert.Equal(t, tc.expectedSkippedTests, skips, "wrong number of skips")
			assert.Equal(t, tc.expectedPanickedTests, panickedTests, "wrong number of panicked tests")
			assert.Equal(t, tc.expectedRacedTests, racedTests, "wrong number of raced tests")
			assert.Equal(t, tc.expectedFlakyTests, flakyTests, "wrong number of flaky tests")

			// Get the output as a string
			output := buf.String()
			for _, expected := range tc.expectedStringsContain {
				assert.Contains(t, output, expected, "output does not contain expected string")
			}
		})
	}

}

func TestAggregateTestResults(t *testing.T) {
	// Create a temporary directory for test JSON files
	tempDir, err := os.MkdirTemp("", "aggregatetestresults")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Test cases
	testCases := []struct {
		description    string
		inputReports   []*TestReport
		expectedReport *TestReport
	}{
		{
			description: "Unique test results",
			inputReports: []*TestReport{
				{
					TestRunCount:  2, // 2 runs of A and 4 runs of B will add up to 6 total runs. Not quite ideal.
					GoProject:     "project1",
					RaceDetection: false,
					Results: []TestResult{
						{
							TestName:    "TestA",
							TestPackage: "pkgA",
							PassRatio:   1,
							Skipped:     false,
							Runs:        2,
							Successes:   2,
							Durations:   []time.Duration{time.Millisecond * 10, time.Millisecond * 20},
							Outputs:     []string{"Output1", "Output2"},
						},
					},
				},
				{
					TestRunCount:  4,
					GoProject:     "project1",
					RaceDetection: false,
					Results: []TestResult{
						{
							TestName:    "TestB",
							TestPackage: "pkgB",
							PassRatio:   0.5,
							Skipped:     false,
							Runs:        4,
							Successes:   2,
							Failures:    2,
							Durations:   []time.Duration{time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50},
							Outputs:     []string{"Output3", "Output4", "Output5", "Output6"},
						},
					},
				},
			},
			expectedReport: &TestReport{
				TestRunCount:  6,
				GoProject:     "project1",
				RaceDetection: false,
				Results: []TestResult{
					{
						TestName:    "TestA",
						TestPackage: "pkgA",
						PassRatio:   1,
						Skipped:     false,
						Runs:        2,
						Successes:   2,
						Durations:   []time.Duration{time.Millisecond * 10, time.Millisecond * 20},
						Outputs:     []string{"Output1", "Output2"},
					},
					{
						TestName:    "TestB",
						TestPackage: "pkgB",
						PassRatio:   0.5,
						Skipped:     false,
						Runs:        4,
						Successes:   2,
						Failures:    2,
						Durations:   []time.Duration{time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50},
						Outputs:     []string{"Output3", "Output4", "Output5", "Output6"},
					},
				},
			},
		},
		{
			description: "Duplicate test results with aggregation",
			inputReports: []*TestReport{
				{
					TestRunCount:  2,
					GoProject:     "project2",
					RaceDetection: false,
					Results: []TestResult{
						{
							TestName:    "TestC",
							TestPackage: "pkgC",
							PassRatio:   1,
							Skipped:     false,
							Runs:        2,
							Successes:   2,
							Durations:   []time.Duration{time.Millisecond * 100, time.Millisecond * 100},
							Outputs:     []string{"Output7", "Output8"},
						},
					},
				},
				{
					TestRunCount:  2,
					GoProject:     "project2",
					RaceDetection: false,
					Results: []TestResult{
						{
							TestName:    "TestC",
							TestPackage: "pkgC",
							PassRatio:   1,
							Skipped:     false,
							Runs:        0,
							Skips:       2,
							Durations:   []time.Duration{time.Millisecond * 200, time.Millisecond * 200},
							Outputs:     []string{"Output9", "Output10"},
						},
					},
				},
			},
			expectedReport: &TestReport{
				TestRunCount:  4,
				GoProject:     "project2",
				RaceDetection: false,
				Results: []TestResult{
					{
						TestName:    "TestC",
						TestPackage: "pkgC",
						PassRatio:   1.0,
						Skipped:     false,
						Runs:        2,
						Successes:   2,
						Durations:   []time.Duration{time.Millisecond * 100, time.Millisecond * 100, time.Millisecond * 200, time.Millisecond * 200},
						Outputs:     []string{"Output7", "Output8", "Output9", "Output10"},
					},
				},
			},
		},
		{
			description: "All Skipped test results",
			inputReports: []*TestReport{
				{
					TestRunCount:  3,
					GoProject:     "project3",
					RaceDetection: false,
					Results: []TestResult{
						{
							TestName:    "TestD",
							TestPackage: "pkgD",
							PassRatio:   1,
							Skipped:     true,
							Runs:        0,
							Durations:   []time.Duration{time.Millisecond * 100, time.Millisecond * 200, time.Millisecond * 100},
							Outputs:     []string{"Output11", "Output12", "Output13"},
						},
					},
				},
				{
					TestRunCount:  2,
					GoProject:     "project3",
					RaceDetection: false,
					Results: []TestResult{
						{
							TestName:    "TestD",
							TestPackage: "pkgD",
							PassRatio:   1,
							Skipped:     true,
							Runs:        0,
							Durations:   []time.Duration{time.Millisecond * 150, time.Millisecond * 150},
							Outputs:     []string{"Output14", "Output15"},
						},
					},
				},
			},
			expectedReport: &TestReport{
				TestRunCount:  5,
				GoProject:     "project3",
				RaceDetection: false,
				Results: []TestResult{
					{
						TestName:    "TestD",
						TestPackage: "pkgD",
						PassRatio:   1,
						Skipped:     true, // Should remain true as all runs are skipped
						Runs:        0,
						Durations:   []time.Duration{time.Millisecond * 100, time.Millisecond * 200, time.Millisecond * 100, time.Millisecond * 150, time.Millisecond * 150},
						Outputs:     []string{"Output11", "Output12", "Output13", "Output14", "Output15"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			finalReport, err := Aggregate(tc.inputReports...)
			if err != nil {
				t.Fatalf("AggregateTestResults failed: %v", err)
			}

			sortTestResults(finalReport.Results)
			sortTestResults(tc.expectedReport.Results)

			assert.Equal(t, tc.expectedReport.TestRunCount, finalReport.TestRunCount, "TestRunCount mismatch")
			assert.Equal(t, tc.expectedReport.GoProject, finalReport.GoProject, "GoProject mismatch")
			assert.Equal(t, tc.expectedReport.RaceDetection, finalReport.RaceDetection, "RaceDetection mismatch")

			require.Equal(t, len(tc.expectedReport.Results), len(finalReport.Results), "number of results mismatch")
			for i, expected := range tc.expectedReport.Results {
				got := finalReport.Results[i]
				assert.Equal(t, expected.TestName, got.TestName, "TestName mismatch")
				assert.Equal(t, expected.TestPackage, got.TestPackage, "TestPackage mismatch")
				assert.Equal(t, expected.Runs, got.Runs, "Runs mismatch")
				assert.Equal(t, expected.Skipped, got.Skipped, "Skipped mismatch")
				assert.Equal(t, expected.PassRatio, got.PassRatio, "PassRatio mismatch")
				assert.Equal(t, len(expected.Durations), len(got.Durations), "Durations mismatch")
				assert.Equal(t, len(expected.Outputs), len(got.Outputs), "Outputs mismatch")
			}
		})
	}
}
