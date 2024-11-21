package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TestResult represents the result of a single test being run through flakeguard.
type TestResult struct {
	TestName            string
	TestPackage         string
	Panicked            bool            // Indicates a test-level panic
	PackagePanicked     bool            // Indicates a package-level panic
	PassRatio           float64         // Pass ratio in decimal format like 0.5
	PassRatioPercentage string          // Pass ratio in percentage format like "50%"
	Skipped             bool            // Indicates if the test was skipped
	Runs                int             // Count of how many times the test was run
	Failures            int             // Count of how many times the test failed
	Successes           int             // Count of how many times the test passed
	Panics              int             // Count of how many times the test panicked
	Skips               int             // Count of how many times the test was skipped
	Outputs             []string        // Stores outputs for a test
	Durations           []time.Duration // Stores elapsed time for each run of the test
	PackageOutputs      []string        // Stores package-level outputs
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

// FilterFlakyTests returns a slice of TestResult where the pass ratio is between the min pass ratio and the threshold.
func FilterFlakyTests(testResults []TestResult, minPassRatio, threshold float64) []TestResult {
	var flakyTests []TestResult
	for _, test := range testResults {
		if test.PassRatio >= minPassRatio && test.PassRatio < threshold && !test.Skipped {
			flakyTests = append(flakyTests, test)
		}
	}
	return flakyTests
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

// AggregateTestResults aggregates all JSON test results.
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

					// Calculate total successful runs for correct pass ratio calculation
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

	// Sort by PassRatio in ascending order
	sort.Slice(aggregatedResults, func(i, j int) bool {
		return aggregatedResults[i].PassRatio < aggregatedResults[j].PassRatio
	})

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
			durationsStr[i] = duration.String()
		}
		fmt.Fprintf(w, "Durations: %s\n", strings.Join(durationsStr, ", "))
		fmt.Fprintf(w, "Outputs:\n%s\n", strings.Join(test.Outputs, ""))
	}
}

// Helper function to save filtered results and logs to specified paths
func SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath string, failedResults []TestResult) {
	if outputResultsPath != "" {
		if err := saveResults(outputResultsPath, failedResults); err != nil {
			log.Fatalf("Error writing failed results to file: %v", err)
		}
		fmt.Printf("Test results saved to %s\n", outputResultsPath)
	} else {
		fmt.Println("No failed tests found based on the specified threshold and min pass ratio.")
	}

	if outputLogsPath != "" {
		if err := saveTestOutputs(outputLogsPath, failedResults); err != nil {
			log.Fatalf("Error writing failed logs to file: %v", err)
		}
		fmt.Printf("Test logs saved to %s\n", outputLogsPath)
	}
}

// Helper function to save results to JSON file
func saveResults(filePath string, results []TestResult) error {
	// Define a struct type without Outputs and PackageOutputs
	type filteredTestResult struct {
		TestName            string
		TestPackage         string
		Panicked            bool
		PackagePanicked     bool
		PassRatio           float64
		PassRatioPercentage string
		Skipped             bool
		Runs                int
		Durations           []time.Duration
	}

	var filteredResults []filteredTestResult
	for _, r := range results {
		filteredResults = append(filteredResults, filteredTestResult{
			TestName:            r.TestName,
			TestPackage:         r.TestPackage,
			Panicked:            r.Panicked,
			PackagePanicked:     r.PackagePanicked,
			PassRatio:           r.PassRatio,
			PassRatioPercentage: r.PassRatioPercentage,
			Skipped:             r.Skipped,
			Runs:                r.Runs,
			Durations:           r.Durations,
		})
	}

	data, err := json.MarshalIndent(filteredResults, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling results: %v", err)
	}
	return os.WriteFile(filePath, data, 0644) //nolint:gosec
}

// Helper function to save test names, packages, and outputs to JSON file
func saveTestOutputs(filePath string, results []TestResult) error {
	// Define a struct type with only the required fields
	type outputOnlyResult struct {
		TestName       string
		TestPackage    string
		Outputs        []string
		PackageOutputs []string
	}

	// Convert results to the filtered struct
	var outputResults []outputOnlyResult
	for _, r := range results {
		outputResults = append(outputResults, outputOnlyResult{
			TestName:       r.TestName,
			TestPackage:    r.TestPackage,
			Outputs:        r.Outputs,
			PackageOutputs: r.PackageOutputs,
		})
	}

	data, err := json.MarshalIndent(outputResults, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling outputs: %v", err)
	}
	return os.WriteFile(filePath, data, 0644) //nolint:gosec
}
