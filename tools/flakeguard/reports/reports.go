package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TestResult struct {
	TestName            string
	TestPackage         string
	PassRatio           float64 // Pass ratio in decimal format like 0.5
	PassRatioPercentage string  // Pass ratio in percentage format like "50%"
	Skipped             bool    // Indicates if the test was skipped
	Runs                int
	Outputs             []string  // Stores outputs for a test
	Durations           []float64 // Stores elapsed time in seconds for each run of the test
}

// FilterFailedTests returns a slice of TestResult where the pass ratio is below the specified threshold.
func FilterFailedTests(results []TestResult, threshold float64) []TestResult {
	var failedTests []TestResult
	for _, result := range results {
		if !result.Skipped && result.PassRatio < threshold {
			failedTests = append(failedTests, result)
		}
	}
	return failedTests
}

// FilterPassedTests returns a slice of TestResult where the tests passed and were not skipped.
func FilterPassedTests(results []TestResult, threshold float64) []TestResult {
	var passedTests []TestResult
	for _, result := range results {
		if !result.Skipped && result.PassRatio >= threshold {
			passedTests = append(passedTests, result)
		}
	}
	return passedTests
}

// FilterSkippedTests returns a slice of TestResult where the tests were skipped.
func FilterSkippedTests(results []TestResult) []TestResult {
	var skippedTests []TestResult
	for _, result := range results {
		if result.Skipped {
			skippedTests = append(skippedTests, result)
		}
	}
	return skippedTests
}

// Helper function to aggregate all JSON test results from a folder
func AggregateTestResults(folderPath string) ([]TestResult, error) {
	// Map to hold unique tests based on their TestName and TestPackage
	testMap := make(map[string]TestResult)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Read file content
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			// Parse JSON data into TestResult slice
			var results []TestResult
			if jsonErr := json.Unmarshal(data, &results); jsonErr != nil {
				return jsonErr
			}
			// Process each result
			for _, result := range results {
				// Unique key for each test based on TestName and TestPackage
				key := result.TestName + "|" + result.TestPackage
				if existingResult, found := testMap[key]; found {
					// Aggregate runs, durations, and outputs
					totalRuns := existingResult.Runs + result.Runs
					existingResult.Durations = append(existingResult.Durations, result.Durations...)
					existingResult.Outputs = append(existingResult.Outputs, result.Outputs...)

					// Calculate total successful runs and aggregate pass ratio
					successfulRuns := existingResult.PassRatio*float64(existingResult.Runs) + result.PassRatio*float64(result.Runs)
					existingResult.Runs = totalRuns
					existingResult.PassRatio = successfulRuns / float64(totalRuns)
					existingResult.Skipped = existingResult.Skipped && result.Skipped // Mark as skipped only if all occurrences are skipped

					// Update the map with the aggregated result
					testMap[key] = existingResult

				} else {
					// Add new entry to the map
					testMap[key] = result
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading files: %v", err)
	}

	// Convert map to slice of TestResult and set PassRatioPercentage
	aggregatedResults := make([]TestResult, 0, len(testMap))
	for _, result := range testMap {
		result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
		aggregatedResults = append(aggregatedResults, result)
	}
	return aggregatedResults, nil
}

// PrintTests prints tests in a pretty format
func PrintTests(tests []TestResult, w io.Writer) {
	for i, test := range tests {
		fmt.Fprintf(w, "\n--- Test %d ---\n", i+1)
		fmt.Fprintf(w, "TestName: %s\n", test.TestName)
		fmt.Fprintf(w, "TestPackage: %s\n", test.TestPackage)
		fmt.Fprintf(w, "PassRatio: %.2f\n", test.PassRatio)
		fmt.Fprintf(w, "Skipped: %v\n", test.Skipped)
		fmt.Fprintf(w, "Runs: %d\n", test.Runs)
		durationsStr := make([]string, len(test.Durations))
		for i, duration := range test.Durations {
			durationsStr[i] = fmt.Sprintf("%.2fs", duration)
		}
		fmt.Fprintf(w, "Durations: %s\n", strings.Join(durationsStr, ", "))
		fmt.Fprintf(w, "Outputs:\n%s\n", strings.Join(test.Outputs, ""))
	}
}
