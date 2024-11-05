package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
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

	if len(failedTests) != len(expected) {
		t.Fatalf("expected %d failed tests, got %d", len(expected), len(failedTests))
	}

	for i, test := range failedTests {
		if test.TestName != expected[i] {
			t.Errorf("expected test %s, got %s", expected[i], test.TestName)
		}
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

	if len(passedTests) != len(expected) {
		t.Fatalf("expected %d passed tests, got %d", len(expected), len(passedTests))
	}

	for i, test := range passedTests {
		if test.TestName != expected[i] {
			t.Errorf("expected test %s, got %s", expected[i], test.TestName)
		}
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

	if len(skippedTests) != len(expected) {
		t.Fatalf("expected %d skipped tests, got %d", len(expected), len(skippedTests))
	}

	for i, test := range skippedTests {
		if test.TestName != expected[i] {
			t.Errorf("expected test %s, got %s", expected[i], test.TestName)
		}
	}
}

func TestPrintTests(t *testing.T) {
	tests := []TestResult{
		{
			TestName:    "Test1",
			TestPackage: "package1",
			PassRatio:   0.75,
			Skipped:     false,
			Runs:        4,
			Outputs:     []string{"Output1", "Output2"},
			Durations:   []float64{1.2, 0.9, 1.1, 1.0},
		},
	}

	// Use a buffer to capture the output
	var buf bytes.Buffer

	// Call PrintTests with the buffer
	PrintTests(tests, &buf)

	// Get the output as a string
	output := buf.String()
	expectedContains := []string{
		"TestName: Test1",
		"TestPackage: package1",
		"PassRatio: 0.75",
		"Skipped: false",
		"Runs: 4",
		"Durations: 1.20s, 0.90s, 1.10s, 1.00s",
		"Outputs:\nOutput1Output2",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, but it did not", expected)
		}
	}
}

// Sorts TestResult slice by TestName and TestPackage for consistent comparison
func sortTestResults(results []TestResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].TestName == results[j].TestName {
			return results[i].TestPackage < results[j].TestPackage
		}
		return results[i].TestName < results[j].TestName
	})
}

// Helper function to write a JSON file for testing
func writeTempJSONFile(t *testing.T, dir string, filename string, data interface{}) string {
	filePath := filepath.Join(dir, filename)
	fileData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}
	return filePath
}

func TestAggregateTestResults(t *testing.T) {
	// Create a temporary directory for test JSON files
	tempDir, err := os.MkdirTemp("", "aggregatetestresults")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	testCases := []struct {
		description    string
		inputFiles     []interface{}
		expectedOutput []TestResult
	}{
		{
			description: "Unique test results without aggregation",
			inputFiles: []interface{}{
				[]TestResult{
					{
						TestName:            "TestA",
						TestPackage:         "pkgA",
						PassRatio:           1,
						PassRatioPercentage: "100%",
						Skipped:             false,
						Runs:                2,
						Durations:           []float64{0.01, 0.02},
						Outputs:             []string{"Output1", "Output2"},
					},
				},
				[]TestResult{
					{
						TestName:            "TestB",
						TestPackage:         "pkgB",
						PassRatio:           0.5,
						PassRatioPercentage: "50%",
						Skipped:             false,
						Runs:                4,
						Durations:           []float64{0.05, 0.05, 0.05, 0.05},
						Outputs:             []string{"Output3", "Output4", "Output5", "Output6"},
					},
				},
			},
			expectedOutput: []TestResult{
				{
					TestName:            "TestA",
					TestPackage:         "pkgA",
					PassRatio:           1,
					PassRatioPercentage: "100%",
					Skipped:             false,
					Runs:                2,
					Durations:           []float64{0.01, 0.02},
					Outputs:             []string{"Output1", "Output2"},
				},
				{
					TestName:            "TestB",
					TestPackage:         "pkgB",
					PassRatio:           0.5,
					PassRatioPercentage: "50%",
					Skipped:             false,
					Runs:                4,
					Durations:           []float64{0.05, 0.05, 0.05, 0.05},
					Outputs:             []string{"Output3", "Output4", "Output5", "Output6"},
				},
			},
		},
		{
			description: "Duplicate test results with aggregation",
			inputFiles: []interface{}{
				[]TestResult{
					{
						TestName:            "TestC",
						TestPackage:         "pkgC",
						PassRatio:           1,
						PassRatioPercentage: "100%",
						Skipped:             false,
						Runs:                2,
						Durations:           []float64{0.1, 0.1},
						Outputs:             []string{"Output7", "Output8"},
					},
				},
				[]TestResult{
					{
						TestName:            "TestC",
						TestPackage:         "pkgC",
						PassRatio:           0.5,
						PassRatioPercentage: "50%",
						Skipped:             false,
						Runs:                2,
						Durations:           []float64{0.2, 0.2},
						Outputs:             []string{"Output9", "Output10"},
					},
				},
			},
			expectedOutput: []TestResult{
				{
					TestName:            "TestC",
					TestPackage:         "pkgC",
					PassRatio:           0.75, // Calculated as (2*1 + 2*0.5) / 4
					PassRatioPercentage: "75%",
					Skipped:             false,
					Runs:                4,
					Durations:           []float64{0.1, 0.1, 0.2, 0.2},
					Outputs:             []string{"Output7", "Output8", "Output9", "Output10"},
				},
			},
		},
		{
			description: "All Skipped test results",
			inputFiles: []interface{}{
				[]TestResult{
					{
						TestName:            "TestD",
						TestPackage:         "pkgD",
						PassRatio:           1,
						PassRatioPercentage: "100%",
						Skipped:             true,
						Runs:                3,
						Durations:           []float64{0.1, 0.2, 0.1},
						Outputs:             []string{"Output11", "Output12", "Output13"},
					},
				},
				[]TestResult{
					{
						TestName:            "TestD",
						TestPackage:         "pkgD",
						PassRatio:           1,
						PassRatioPercentage: "100%",
						Skipped:             true,
						Runs:                2,
						Durations:           []float64{0.15, 0.15},
						Outputs:             []string{"Output14", "Output15"},
					},
				},
			},
			expectedOutput: []TestResult{
				{
					TestName:            "TestD",
					TestPackage:         "pkgD",
					PassRatio:           1,
					PassRatioPercentage: "100%",
					Skipped:             true, // Should remain true as all runs are skipped
					Runs:                5,
					Durations:           []float64{0.1, 0.2, 0.1, 0.15, 0.15},
					Outputs:             []string{"Output11", "Output12", "Output13", "Output14", "Output15"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Write input files to the temporary directory
			for i, inputData := range tc.inputFiles {
				writeTempJSONFile(t, tempDir, fmt.Sprintf("input%d.json", i), inputData)
			}

			// Run AggregateTestResults
			result, err := AggregateTestResults(tempDir)
			if err != nil {
				t.Fatalf("AggregateTestResults failed: %v", err)
			}

			// Sort both result and expectedOutput for consistent comparison
			sortTestResults(result)
			sortTestResults(tc.expectedOutput)

			// Compare the result with the expected output
			if len(result) != len(tc.expectedOutput) {
				t.Fatalf("Expected %d results, got %d", len(tc.expectedOutput), len(result))
			}

			for i, expected := range tc.expectedOutput {
				got := result[i]
				if got.TestName != expected.TestName || got.TestPackage != expected.TestPackage || got.Runs != expected.Runs || got.Skipped != expected.Skipped {
					t.Errorf("Result %d - expected %+v, got %+v", i, expected, got)
				}
				if got.PassRatio != expected.PassRatio {
					t.Errorf("Result %d - expected PassRatio %f, got %f", i, expected.PassRatio, got.PassRatio)
				}
				if got.PassRatioPercentage != expected.PassRatioPercentage {
					t.Errorf("Result %d - expected PassRatioPercentage %s, got %s", i, expected.PassRatioPercentage, got.PassRatioPercentage)
				}
				if len(got.Durations) != len(expected.Durations) {
					t.Errorf("Result %d - expected %d durations, got %d", i, len(expected.Durations), len(got.Durations))
				}
				if len(got.Outputs) != len(expected.Outputs) {
					t.Errorf("Result %d - expected %d outputs, got %d", i, len(expected.Outputs), len(got.Outputs))
				}
			}
		})
	}
}
