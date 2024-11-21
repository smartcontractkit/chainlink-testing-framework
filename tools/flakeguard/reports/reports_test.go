package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	tests := []TestResult{
		{
			TestName:    "Test1",
			TestPackage: "package1",
			PassRatio:   0.75,
			Skipped:     false,
			Runs:        4,
			Outputs:     []string{"Output1", "Output2"},
			Durations:   []time.Duration{time.Millisecond * 1200, time.Millisecond * 900, time.Millisecond * 1100, time.Second},
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
		"Durations: 1.2s, 900ms, 1.1s, 1s",
		"Outputs:\nOutput1Output2",
	}

	for _, expected := range expectedContains {
		assert.Contains(t, output, expected, "printed test output doesn't contain expected string")
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
	t.Helper()

	filePath := filepath.Join(dir, filename)
	fileData, err := json.Marshal(data)
	require.NoError(t, err)
	err = os.WriteFile(filePath, fileData, 0644) //nolint:gosec
	require.NoError(t, err)
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
						Durations:           []time.Duration{time.Millisecond * 10, time.Millisecond * 20},
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
						Durations:           []time.Duration{time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50},
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
					Durations:           []time.Duration{time.Millisecond * 10, time.Millisecond * 20},
					Outputs:             []string{"Output1", "Output2"},
				},
				{
					TestName:            "TestB",
					TestPackage:         "pkgB",
					PassRatio:           0.5,
					PassRatioPercentage: "50%",
					Skipped:             false,
					Runs:                4,
					Durations:           []time.Duration{time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50, time.Millisecond * 50},
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
						Durations:           []time.Duration{time.Millisecond * 100, time.Millisecond * 100},
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
						Durations:           []time.Duration{time.Millisecond * 200, time.Millisecond * 200},
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
					Durations:           []time.Duration{time.Millisecond * 100, time.Millisecond * 100, time.Millisecond * 200, time.Millisecond * 200},
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
						Durations:           []time.Duration{time.Millisecond * 100, time.Millisecond * 200, time.Millisecond * 100},
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
						Durations:           []time.Duration{time.Millisecond * 150, time.Millisecond * 150},
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
					Durations:           []time.Duration{time.Millisecond * 100, time.Millisecond * 200, time.Millisecond * 100, time.Millisecond * 150, time.Millisecond * 150},
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
			require.Equal(t, len(tc.expectedOutput), len(result), "number of results mismatch")

			for i, expected := range tc.expectedOutput {
				got := result[i]
				assert.Equal(t, expected.TestName, got.TestName, "TestName mismatch")
				assert.Equal(t, expected.TestPackage, got.TestPackage, "TestPackage mismatch")
				assert.Equal(t, expected.Runs, got.Runs, "Runs mismatch")
				assert.Equal(t, expected.Skipped, got.Skipped, "Skipped mismatch")
				assert.Equal(t, expected.PassRatio, got.PassRatio, "PassRatio mismatch")
				assert.Equal(t, expected.PassRatioPercentage, got.PassRatioPercentage, "PassRatioPercentage mismatch")
				assert.Equal(t, len(expected.Durations), len(got.Durations), "Durations mismatch")
				assert.Equal(t, len(expected.Outputs), len(got.Outputs), "Outputs mismatch")
			}
		})
	}
}
