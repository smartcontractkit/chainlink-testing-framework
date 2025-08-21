package reports

import (
	"math"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGenerateSummaryData tests the GenerateSummaryData function.
func TestGenerateSummaryData(t *testing.T) {
	tests := []struct {
		name       string
		testReport *TestReport
		expected   *SummaryData
	}{
		{
			name: "All tests passed",
			testReport: &TestReport{
				Results: []TestResult{
					{PassRatio: 1.0, Runs: 10, Successes: 10},
					{PassRatio: 1.0, Runs: 5, Successes: 5},
				},
				MaxPassRatio: 1.0,
			},
			expected: &SummaryData{
				UniqueTestsRun:         2,
				TestRunCount:           10,
				PanickedTests:          0,
				RacedTests:             0,
				FlakyTests:             0,
				FlakyTestPercent:       "0%", // no flaky tests
				TotalRuns:              15,
				PassedRuns:             15,
				FailedRuns:             0,
				SkippedRuns:            0,
				PassPercent:            "100%",
				UniqueSkippedTestCount: 0,
			},
		},
		{
			name: "Some flaky tests",
			testReport: &TestReport{
				Results: []TestResult{
					{PassRatio: 0.8, Runs: 10, Successes: 8, Failures: 2},
					{PassRatio: 1.0, Runs: 5, Successes: 5},
					{PassRatio: 0.5, Runs: 4, Successes: 2, Failures: 2},
				},
				MaxPassRatio: 0.9,
			},
			expected: &SummaryData{
				UniqueTestsRun: 3,
				TestRunCount:   10,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     2,
				// 2/3 => 66.666...%
				FlakyTestPercent: "66.6667%",
				TotalRuns:        19,
				PassedRuns:       15,
				FailedRuns:       4,
				SkippedRuns:      0,
				// 15/19 => ~78.947...
				PassPercent:            "78.9474%",
				UniqueSkippedTestCount: 0,
			},
		},
		{
			name: "Tests with panics and races",
			testReport: &TestReport{
				Results: []TestResult{
					{PassRatio: 1.0, Runs: 5, Successes: 5, Panic: true},
					{PassRatio: 0.9, Runs: 10, Successes: 9, Failures: 1, Race: true},
					{PassRatio: 1.0, Runs: 3, Successes: 3},
				},
				MaxPassRatio: 1.0,
			},
			expected: &SummaryData{
				UniqueTestsRun: 3,
				TestRunCount:   10,
				PanickedTests:  1,
				RacedTests:     1,
				FlakyTests:     2,
				// 2/3 => ~66.666...
				FlakyTestPercent: "66.6667%",
				TotalRuns:        18,
				PassedRuns:       17,
				FailedRuns:       1,
				SkippedRuns:      0,
				// 17/18 => ~94.444...
				PassPercent:            "94.4444%",
				UniqueSkippedTestCount: 0,
			},
		},
		{
			name: "No tests ran",
			testReport: &TestReport{
				MaxPassRatio: 1.0,
				Results:      []TestResult{},
			},
			expected: &SummaryData{
				UniqueTestsRun:   0,
				TestRunCount:     0,
				PanickedTests:    0,
				RacedTests:       0,
				FlakyTests:       0,
				FlakyTestPercent: "0%",
				TotalRuns:        0,
				PassedRuns:       0,
				FailedRuns:       0,
				SkippedRuns:      0,
				// With zero runs, we default passRatio to "100%"
				PassPercent:            "100%",
				UniqueSkippedTestCount: 0,
			},
		},
		{
			name: "Skipped tests included in total but not executed",
			testReport: &TestReport{
				Results: []TestResult{
					// Skipped test with no runs should be counted in UniqueSkippedTestCount.
					{PassRatio: -1.0, Runs: 0, Successes: 0, Skips: 1, Skipped: true},
					{PassRatio: 0.7, Runs: 10, Successes: 7, Failures: 3},
				},
				MaxPassRatio: 0.8,
			},
			expected: &SummaryData{
				UniqueTestsRun:         2,
				TestRunCount:           10,
				PanickedTests:          0,
				RacedTests:             0,
				FlakyTests:             1,
				FlakyTestPercent:       "50%",
				TotalRuns:              10,
				PassedRuns:             7,
				FailedRuns:             3,
				SkippedRuns:            1,
				PassPercent:            "70%",
				UniqueSkippedTestCount: 1,
			},
		},
		{
			name: "Mixed skipped and executed tests",
			testReport: &TestReport{
				Results: []TestResult{
					// Skipped test should be counted for UniqueSkippedTestCount.
					{PassRatio: -1.0, Runs: 0, Successes: 0, Skips: 1, Skipped: true},
					{PassRatio: 0.9, Runs: 10, Successes: 9, Failures: 1},
					{PassRatio: 0.5, Runs: 4, Successes: 2, Failures: 2},
				},
				MaxPassRatio: 0.85,
			},
			expected: &SummaryData{
				UniqueTestsRun:         3,
				TestRunCount:           10,
				PanickedTests:          0,
				RacedTests:             0,
				FlakyTests:             1,
				FlakyTestPercent:       "33.3333%",
				TotalRuns:              14,
				PassedRuns:             11,
				FailedRuns:             3,
				SkippedRuns:            1,
				PassPercent:            "78.5714%",
				UniqueSkippedTestCount: 1,
			},
		},
		{
			name: "Tiny flake ratio that is exactly 0.01%",
			testReport: &TestReport{
				Results: func() []TestResult {
					// 9,999 total:
					//  - 9,998 stable tests => pass ratio = 1.0
					//  - 1 flaky test => pass ratio = 0.5
					const total = 9999
					tests := make([]TestResult, total)
					for i := 0; i < total-1; i++ {
						tests[i] = TestResult{
							PassRatio: 1.0,
							Runs:      10,
							Successes: 10,
						}
					}
					tests[total-1] = TestResult{
						PassRatio: 0.5, // 1 success, 1 failure
						Runs:      2,
						Successes: 1,
						Failures:  1,
					}
					return tests
				}(),
				MaxPassRatio: 1.0,
			},
			expected: &SummaryData{
				UniqueTestsRun:         9999,
				TestRunCount:           10,
				PanickedTests:          0,
				RacedTests:             0,
				FlakyTests:             1,
				FlakyTestPercent:       "0.01%",
				TotalRuns:              (9998 * 10) + 2,
				PassedRuns:             (9998 * 10) + 1,
				FailedRuns:             1,
				SkippedRuns:            0,
				PassPercent:            "99.999%",
				UniqueSkippedTestCount: 0,
			},
		},
		{
			name: "Duplicate skipped tests",
			testReport: &TestReport{
				Results: []TestResult{
					// Two entries for "TestA" should count as one unique skipped test.
					{TestName: "TestA", PassRatio: -1.0, Runs: 0, Skips: 1, Skipped: true},
					{TestName: "TestA", PassRatio: -1.0, Runs: 0, Skips: 1, Skipped: true},
					// A different test "TestB"
					{TestName: "TestB", PassRatio: -1.0, Runs: 0, Skips: 1, Skipped: true},
					// This test was executed so it should not count as skipped.
					{TestName: "TestC", PassRatio: 1.0, Runs: 5, Successes: 5, Skipped: false},
				},
				MaxPassRatio: 1.0,
			},
			expected: &SummaryData{
				UniqueTestsRun:         4,
				TestRunCount:           5,
				PanickedTests:          0,
				RacedTests:             0,
				FlakyTests:             0,
				FlakyTestPercent:       "0%",
				TotalRuns:              5,
				PassedRuns:             5,
				FailedRuns:             0,
				SkippedRuns:            3, // Sum of Skips for all skipped tests.
				PassPercent:            "100%",
				UniqueSkippedTestCount: 2, // Only "TestA" and "TestB" count.
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testReport.GenerateSummaryData()
			assert.Equal(t, tc.expected, tc.testReport.SummaryData, "Summary data does not match expected")
		})
	}
}

// TestFilterTests tests the FilterTests function.
func TestFilterTests(t *testing.T) {
	testResults := []TestResult{
		{TestName: "TestA", PassRatio: 1.0, Skipped: false},
		{TestName: "TestB", PassRatio: 0.8, Skipped: false},
		{TestName: "TestC", PassRatio: 0.7, Skipped: true},
		{TestName: "TestD", PassRatio: 0.6, Skipped: false},
	}

	// Filter tests with PassRatio < 0.9 and not skipped
	filtered := FilterTests(testResults, func(tr TestResult) bool {
		return !tr.Skipped && tr.PassRatio < 0.9
	})

	expected := []TestResult{
		{TestName: "TestB", PassRatio: 0.8, Skipped: false},
		{TestName: "TestD", PassRatio: 0.6, Skipped: false},
	}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Expected %+v, got %+v", expected, filtered)
	}
}

func TestFilterFailedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.5, Skipped: false},
		{TestName: "Test2", PassRatio: 0.9, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true}, // Skipped test
	}

	failedTests := FilterTests(results, func(tr TestResult) bool {
		return !tr.Skipped && tr.PassRatio < 0.6
	})
	expected := []TestResult{
		{TestName: "Test1", PassRatio: 0.5, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
	}

	if !reflect.DeepEqual(failedTests, expected) {
		t.Errorf("Expected failed tests %+v, got %+v", expected, failedTests)
	}
}

func TestFilterPassedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true}, // Skipped test
	}

	passedTests := FilterTests(results, func(tr TestResult) bool {
		return !tr.Skipped && tr.PassRatio >= 0.6
	})
	expected := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: false},
	}

	if !reflect.DeepEqual(passedTests, expected) {
		t.Errorf("Expected passed tests %+v, got %+v", expected, passedTests)
	}
}

func TestFilterSkippedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: true},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true},
	}

	skippedTests := FilterTests(results, func(tr TestResult) bool {
		return tr.Skipped
	})
	expected := []TestResult{
		{TestName: "Test2", PassRatio: 1.0, Skipped: true},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true},
	}

	if !reflect.DeepEqual(skippedTests, expected) {
		t.Errorf("Expected skipped tests %+v, got %+v", expected, skippedTests)
	}
}

// TestAggregate tests the Aggregate function.
func TestAggregate(t *testing.T) {
	report1 := &TestReport{
		GoProject: "ProjectX",
		Results: []TestResult{
			{
				TestName:    "TestA",
				TestPackage: "pkg1",
				Runs:        2,
				Successes:   2,
				PassRatio:   1.0,
			},
			{
				TestName:    "TestB",
				TestPackage: "pkg1",
				Runs:        2,
				Successes:   1,
				Failures:    1,
				PassRatio:   0.5,
			},
		},
	}

	report2 := &TestReport{
		GoProject: "ProjectX",
		Results: []TestResult{
			{
				TestName:    "TestA",
				TestPackage: "pkg1",
				Runs:        3,
				Successes:   3,
				PassRatio:   1.0,
			},
			{
				TestName:    "TestC",
				TestPackage: "pkg2",
				Runs:        3,
				Successes:   2,
				Failures:    1,
				PassRatio:   0.6667,
			},
		},
	}

	// Create channels for test results and errors.
	resultsChan := make(chan []TestResult)
	errChan := make(chan error)

	// Launch a goroutine to send the test results into the results channel.
	go func() {
		resultsChan <- report1.Results
		resultsChan <- report2.Results
		close(resultsChan)
	}()

	// No errors to send; close the error channel.
	go func() {
		close(errChan)
	}()

	// Call the updated aggregate function.
	aggregatedResults, err := aggregate(resultsChan, errChan)
	if err != nil {
		t.Fatalf("Error aggregating reports: %v", err)
	}

	expectedResults := []TestResult{
		{
			TestName:    "TestA",
			TestPackage: "pkg1",
			Runs:        5,
			Successes:   5,
			Failures:    0,
			PassRatio:   1.0,
		},
		{
			TestName:    "TestB",
			TestPackage: "pkg1",
			Runs:        2,
			Successes:   1,
			Failures:    1,
			PassRatio:   0.5,
		},
		{
			TestName:    "TestC",
			TestPackage: "pkg2",
			Runs:        3,
			Successes:   2,
			Failures:    1,
			PassRatio:   0.6667,
		},
	}

	// Sort both slices by TestName to ensure the order matches for comparison.
	sort.Slice(expectedResults, func(i, j int) bool {
		return expectedResults[i].TestName < expectedResults[j].TestName
	})
	sort.Slice(aggregatedResults, func(i, j int) bool {
		return aggregatedResults[i].TestName < aggregatedResults[j].TestName
	})

	// Compare the aggregated results with expected results.
	for i, result := range aggregatedResults {
		expected := expectedResults[i]
		if result.TestName != expected.TestName ||
			result.TestPackage != expected.TestPackage ||
			result.Runs != expected.Runs ||
			result.Successes != expected.Successes ||
			result.Failures != expected.Failures ||
			math.Abs(result.PassRatio-expected.PassRatio) > 0.0001 {
			t.Errorf("Mismatch in aggregated result for test %s. Expected %+v, got %+v", expected.TestName, expected, result)
		}
	}
}

func TestAggregateOutputs(t *testing.T) {
	report1 := &TestReport{
		GoProject:   "ProjectX",
		SummaryData: &SummaryData{UniqueTestsRun: 1},
		Results: []TestResult{
			{
				TestName:    "TestOutput",
				TestPackage: "pkg1",
				Runs:        1,
				Successes:   1,
				PassedOutputs: map[string][]string{
					"run1": {"Output from report1 test run"},
				},
				PackageOutputs: []string{"Package output from report1"},
			},
		},
	}

	report2 := &TestReport{
		GoProject:   "ProjectX",
		SummaryData: &SummaryData{UniqueTestsRun: 1},
		Results: []TestResult{
			{
				TestName:    "TestOutput",
				TestPackage: "pkg1",
				Runs:        1,
				Successes:   1,
				PassedOutputs: map[string][]string{
					"run2": {"Output from report2 test run"},
				},
				PackageOutputs: []string{"Package output from report2"},
			},
		},
	}

	// Create channels for results and errors.
	resultsChan := make(chan []TestResult)
	errChan := make(chan error)

	// Launch a goroutine to send the results into the resultsChan.
	go func() {
		resultsChan <- report1.Results
		resultsChan <- report2.Results
		close(resultsChan)
	}()

	// Launch a goroutine to close the error channel (no errors to report).
	go func() {
		close(errChan)
	}()

	// Call the new aggregate function.
	aggregatedResults, err := aggregate(resultsChan, errChan)
	if err != nil {
		t.Fatalf("Error aggregating reports: %v", err)
	}

	if len(aggregatedResults) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(aggregatedResults))
	}

	result := aggregatedResults[0]

	expectedOutputs := map[string][]string{
		"run1": {"Output from report1 test run"},
		"run2": {"Output from report2 test run"},
	}

	expectedPackageOutputs := []string{
		"Package output from report1",
		"Package output from report2",
	}

	if !reflect.DeepEqual(result.PassedOutputs, expectedOutputs) {
		t.Errorf("Expected Outputs %v, got %v", expectedOutputs, result.PassedOutputs)
	}

	if !reflect.DeepEqual(result.PackageOutputs, expectedPackageOutputs) {
		t.Errorf("Expected PackageOutputs %v, got %v", expectedPackageOutputs, result.PackageOutputs)
	}
}

func TestAggregateIdenticalOutputs(t *testing.T) {
	report1 := &TestReport{
		GoProject:   "ProjectX",
		SummaryData: &SummaryData{UniqueTestsRun: 1},
		Results: []TestResult{
			{
				TestName:    "TestIdenticalOutput",
				TestPackage: "pkg1",
				Runs:        1,
				Successes:   1,
				PassedOutputs: map[string][]string{
					"run1": {"Identical output"},
				},
				PackageOutputs: []string{"Identical package output"},
			},
		},
	}

	report2 := &TestReport{
		GoProject:   "ProjectX",
		SummaryData: &SummaryData{UniqueTestsRun: 1},
		Results: []TestResult{
			{
				TestName:    "TestIdenticalOutput",
				TestPackage: "pkg1",
				Runs:        1,
				Successes:   1,
				PassedOutputs: map[string][]string{
					"run1": {"Identical output"},
				},
				PackageOutputs: []string{"Identical package output"},
			},
		},
	}

	// Create channels for results and errors.
	resultsChan := make(chan []TestResult)
	errChan := make(chan error)

	// Send the results from both reports.
	go func() {
		resultsChan <- report1.Results
		resultsChan <- report2.Results
		close(resultsChan)
	}()

	// Close error channel since there are no errors.
	go func() {
		close(errChan)
	}()

	aggregatedResults, err := aggregate(resultsChan, errChan)
	if err != nil {
		t.Fatalf("Error aggregating reports: %v", err)
	}

	if len(aggregatedResults) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(aggregatedResults))
	}

	result := aggregatedResults[0]

	expectedOutputs := map[string][]string{
		"run1": {"Identical output", "Identical output"},
	}

	expectedPackageOutputs := []string{
		"Identical package output",
		"Identical package output",
	}

	if !reflect.DeepEqual(result.PassedOutputs, expectedOutputs) {
		t.Errorf("Expected Outputs %v, got %v", expectedOutputs, result.PassedOutputs)
	}

	if !reflect.DeepEqual(result.PackageOutputs, expectedPackageOutputs) {
		t.Errorf("Expected PackageOutputs %v, got %v", expectedPackageOutputs, result.PackageOutputs)
	}
}

// TestAvgDuration tests the avgDuration function.
func TestAvgDuration(t *testing.T) {
	durations := []time.Duration{
		time.Second,
		2 * time.Second,
		3 * time.Second,
	}
	expected := 2 * time.Second

	avg := avgDuration(durations)
	if avg != expected {
		t.Errorf("Expected average duration %v, got %v", expected, avg)
	}

	// Test with empty slice
	avg = avgDuration([]time.Duration{})
	if avg != 0 {
		t.Errorf("Expected average duration 0, got %v", avg)
	}
}

func TestAggregate_AllSkippedTests(t *testing.T) {
	report1 := &TestReport{
		GoProject:   "ProjectX",
		SummaryData: &SummaryData{UniqueTestsRun: 3},
		Results: []TestResult{
			{
				TestName:    "TestSkipped",
				TestPackage: "pkg1",
				Skipped:     true,
				Runs:        0,
				Skips:       3,
				PassRatio:   -1, // -1 indicates undefined
			},
		},
	}

	report2 := &TestReport{
		GoProject:   "ProjectX",
		SummaryData: &SummaryData{UniqueTestsRun: 2},
		Results: []TestResult{
			{
				TestName:    "TestSkipped",
				TestPackage: "pkg1",
				Skipped:     true,
				Runs:        0,
				Skips:       2,
				PassRatio:   -1,
			},
		},
	}

	// Create channels for results and errors.
	resultsChan := make(chan []TestResult)
	errChan := make(chan error)

	// Send results from both reports.
	go func() {
		resultsChan <- report1.Results
		resultsChan <- report2.Results
		close(resultsChan)
	}()

	// Close error channel since there are no errors.
	go func() {
		close(errChan)
	}()

	aggregatedResults, err := aggregate(resultsChan, errChan)
	if err != nil {
		t.Fatalf("Error aggregating reports: %v", err)
	}

	expectedResult := TestResult{
		TestName:    "TestSkipped",
		TestPackage: "pkg1",
		Skipped:     true,
		Runs:        0,
		Skips:       5,
		PassRatio:   -1,
	}

	if len(aggregatedResults) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(aggregatedResults))
	}

	result := aggregatedResults[0]

	if result.TestName != expectedResult.TestName {
		t.Errorf("Expected TestName %v, got %v", expectedResult.TestName, result.TestName)
	}
	if result.TestPackage != expectedResult.TestPackage {
		t.Errorf("Expected TestPackage %v, got %v", expectedResult.TestPackage, result.TestPackage)
	}
	if result.Skipped != expectedResult.Skipped {
		t.Errorf("Expected Skipped %v, got %v", expectedResult.Skipped, result.Skipped)
	}
	if result.Runs != expectedResult.Runs {
		t.Errorf("Expected Runs %v, got %v", expectedResult.Runs, result.Runs)
	}
	if result.Skips != expectedResult.Skips {
		t.Errorf("Expected Skips %v, got %v", expectedResult.Skips, result.Skips)
	}
	if result.PassRatio != expectedResult.PassRatio {
		t.Errorf("Expected PassRatio %v, got %v", expectedResult.PassRatio, result.PassRatio)
	}
}
