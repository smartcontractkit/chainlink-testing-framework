package reports

import (
	"math"
	"reflect"
	"sort"
	"testing"
	"time"
)

// TestGenerateSummaryData tests the GenerateSummaryData function.
func TestGenerateSummaryData(t *testing.T) {
	tests := []struct {
		name         string
		testResults  []TestResult
		maxPassRatio float64
		expected     SummaryData
	}{
		{
			name: "All tests passed",
			testResults: []TestResult{
				{PassRatio: 1.0, Runs: 10, Successes: 10},
				{PassRatio: 1.0, Runs: 5, Successes: 5},
			},
			maxPassRatio: 1.0,
			expected: SummaryData{
				TotalTests:     2,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     0,
				FlakyTestRatio: "0%", // no flaky tests
				TotalRuns:      15,
				PassedRuns:     15,
				FailedRuns:     0,
				SkippedRuns:    0,
				PassRatio:      "100%",
				MaxPassRatio:   1.0,
			},
		},
		{
			name: "Some flaky tests",
			testResults: []TestResult{
				{PassRatio: 0.8, Runs: 10, Successes: 8, Failures: 2},
				{PassRatio: 1.0, Runs: 5, Successes: 5},
				{PassRatio: 0.5, Runs: 4, Successes: 2, Failures: 2},
			},
			maxPassRatio: 0.9,
			expected: SummaryData{
				TotalTests:    3,
				PanickedTests: 0,
				RacedTests:    0,
				FlakyTests:    2,
				// 2/3 => 66.666...%
				FlakyTestRatio: "66.6667%",
				TotalRuns:      19,
				PassedRuns:     15,
				FailedRuns:     4, // total failures
				SkippedRuns:    0,
				// 15/19 => ~78.947...
				PassRatio:    "78.9474%",
				MaxPassRatio: 0.9,
			},
		},
		{
			name: "Tests with panics and races",
			testResults: []TestResult{
				{PassRatio: 1.0, Runs: 5, Successes: 5, Panic: true},
				{PassRatio: 0.9, Runs: 10, Successes: 9, Failures: 1, Race: true},
				{PassRatio: 1.0, Runs: 3, Successes: 3},
			},
			maxPassRatio: 1.0,
			expected: SummaryData{
				TotalTests:    3,
				PanickedTests: 1,
				RacedTests:    1,
				FlakyTests:    2,
				// 2/3 => ~66.666...
				FlakyTestRatio: "66.6667%",
				TotalRuns:      18,
				PassedRuns:     17,
				FailedRuns:     1,
				SkippedRuns:    0,
				// 17/18 => ~94.444...
				PassRatio:    "94.4444%",
				MaxPassRatio: 1.0,
			},
		},
		{
			name:         "No tests ran",
			testResults:  []TestResult{},
			maxPassRatio: 1.0,
			expected: SummaryData{
				TotalTests:     0,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     0,
				FlakyTestRatio: "0%", // no tests => 0%
				TotalRuns:      0,
				PassedRuns:     0,
				FailedRuns:     0,
				SkippedRuns:    0,
				// With zero runs, we default passRatio to "100%"
				PassRatio:    "100%",
				MaxPassRatio: 1.0,
			},
		},
		{
			name: "Skipped tests included in total but not executed",
			testResults: []TestResult{
				{PassRatio: -1.0, Runs: 0, Successes: 0, Skips: 1, Skipped: true},
				{PassRatio: 0.7, Runs: 10, Successes: 7, Failures: 3},
			},
			maxPassRatio: 0.8,
			expected: SummaryData{
				TotalTests:     2,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     1,     // second test has ratio=0.7 => "flaky"
				FlakyTestRatio: "50%", // 1 out of 2 => 50%
				TotalRuns:      10,
				PassedRuns:     7,
				FailedRuns:     3,
				SkippedRuns:    1, // from first test
				PassRatio:      "70%",
				MaxPassRatio:   0.8,
			},
		},
		{
			name: "Mixed skipped and executed tests",
			testResults: []TestResult{
				{PassRatio: -1.0, Runs: 0, Successes: 0, Skips: 1, Skipped: true},
				{PassRatio: 0.9, Runs: 10, Successes: 9, Failures: 1},
				{PassRatio: 0.5, Runs: 4, Successes: 2, Failures: 2},
			},
			maxPassRatio: 0.85,
			expected: SummaryData{
				TotalTests:     3,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     1,          // last test has ratio=0.5 => "flaky"
				FlakyTestRatio: "33.3333%", // 1 out of 3 => 33.333...
				TotalRuns:      14,         // 10 + 4
				PassedRuns:     11,         // 9 + 2
				FailedRuns:     3,          // 1 + 2
				SkippedRuns:    1,          // from first test
				// 11/14 => 78.5714...
				PassRatio:    "78.5714%",
				MaxPassRatio: 0.85,
			},
		},
		{
			name: "Tiny flake ratio that is exactly 0.01%",
			testResults: func() []TestResult {
				// 9,999 total:
				//  - 9,998 stable => pass=1.0
				//  - 1 flaky => pass=0.5
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
					PassRatio: 0.5, // 1 success, 1 fail
					Runs:      2,
					Successes: 1,
					Failures:  1,
				}
				return tests
			}(),
			maxPassRatio: 1.0,
			expected: SummaryData{
				TotalTests:     9999,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     1,
				FlakyTestRatio: "0.01%",
				TotalRuns:      (9998 * 10) + 2,
				PassedRuns:     (9998 * 10) + 1,
				FailedRuns:     1,
				SkippedRuns:    0,
				PassRatio:      "99.999%",
				MaxPassRatio:   1.0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			summary := GenerateSummaryData(tc.testResults, tc.maxPassRatio)
			if !reflect.DeepEqual(summary, tc.expected) {
				t.Errorf("Test %q failed.\nExpected: %+v\nGot:      %+v",
					tc.name, tc.expected, summary)
			}
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
		GoProject:    "ProjectX",
		TestRunCount: 2,
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
		GoProject:    "ProjectX",
		TestRunCount: 3,
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

	aggregateOptions := &aggregateOptions{reportID: "123"}
	aggregatedReport, err := aggregateFromReports(aggregateOptions, report1, report2)
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

	// Sort results for comparison
	sort.Slice(expectedResults, func(i, j int) bool {
		return expectedResults[i].TestName < expectedResults[j].TestName
	})
	sort.Slice(aggregatedReport.Results, func(i, j int) bool {
		return aggregatedReport.Results[i].TestName < aggregatedReport.Results[j].TestName
	})

	for i, result := range aggregatedReport.Results {
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
		GoProject:    "ProjectX",
		TestRunCount: 1,
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
		GoProject:    "ProjectX",
		TestRunCount: 1,
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

	aggregateOptions := &aggregateOptions{reportID: "123"}
	aggregatedReport, err := aggregateFromReports(aggregateOptions, report1, report2)
	if err != nil {
		t.Fatalf("Error aggregating reports: %v", err)
	}

	if len(aggregatedReport.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(aggregatedReport.Results))
	}

	result := aggregatedReport.Results[0]

	expectedOutputs := map[string][]string{
		"run1": {
			"Output from report1 test run",
		},
		"run2": {
			"Output from report2 test run",
		},
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
		GoProject:    "ProjectX",
		TestRunCount: 1,
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
		GoProject:    "ProjectX",
		TestRunCount: 1,
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

	aggregateOptions := &aggregateOptions{reportID: "123"}
	aggregatedReport, err := aggregateFromReports(aggregateOptions, report1, report2)
	if err != nil {
		t.Fatalf("Error aggregating reports: %v", err)
	}

	if len(aggregatedReport.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(aggregatedReport.Results))
	}

	result := aggregatedReport.Results[0]

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
		GoProject:    "ProjectX",
		TestRunCount: 3,
		Results: []TestResult{
			{
				TestName:    "TestSkipped",
				TestPackage: "pkg1",
				Skipped:     true,
				Runs:        0,
				Skips:       3,
				PassRatio:   -1, // 1 indicate undefined
			},
		},
	}

	report2 := &TestReport{
		GoProject:    "ProjectX",
		TestRunCount: 2,
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

	aggregateOptions := &aggregateOptions{reportID: "123"}
	aggregatedReport, err := aggregateFromReports(aggregateOptions, report1, report2)
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

	if len(aggregatedReport.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(aggregatedReport.Results))
	}

	result := aggregatedReport.Results[0]

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
