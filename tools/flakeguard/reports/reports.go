package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TestReport represents the report of all tests run through flakeguard.
type TestReport struct {
	GoProject     string
	TestRunCount  int
	RaceDetection bool
	ExcludedTests []string
	Results       []TestResult
}

// TestResult represents the result of a single test being run through flakeguard.
type TestResult struct {
	TestName        string
	TestPackage     string
	Panicked        bool            // Indicates a test-level panic
	Skipped         bool            // Indicates if the test was skipped
	PackagePanicked bool            // Indicates a package-level panic
	PassRatio       float64         // Pass ratio in decimal format like 0.5
	Runs            int             // Count of how many times the test was run
	Failures        int             // Count of how many times the test failed
	Successes       int             // Count of how many times the test passed
	Panics          int             // Count of how many times the test panicked
	Races           int             // Count of how many times the test encountered a data race
	Skips           int             // Count of how many times the test was skipped
	Outputs         []string        `json:"outputs,omitempty"` // Stores outputs for a test
	Durations       []time.Duration // Stores elapsed time for each run of the test
	PackageOutputs  []string        `json:"package_outputs,omitempty"` // Stores package-level outputs
}

// FilterFailedTests returns a slice of TestResult where the pass ratio is below the specified threshold.
func FilterFailedTests(results []TestResult, maxPassRatio float64) []TestResult {
	var failedTests []TestResult
	for _, result := range results {
		if !result.Skipped && result.PassRatio < maxPassRatio {
			failedTests = append(failedTests, result)
		}
	}
	return failedTests
}

// FilterFlakyTests returns a slice of TestResult where the pass ratio is between the min pass ratio and the threshold.
func FilterFlakyTests(testResults []TestResult, maxPassRatio float64) []TestResult {
	var flakyTests []TestResult
	for _, test := range testResults {
		if test.PassRatio < maxPassRatio && !test.Skipped {
			flakyTests = append(flakyTests, test)
		}
	}
	return flakyTests
}

// FilterPassedTests returns a slice of TestResult where the tests passed and were not skipped.
func FilterPassedTests(results []TestResult, maxPassRatio float64) []TestResult {
	var passedTests []TestResult
	for _, result := range results {
		if !result.Skipped && result.PassRatio >= maxPassRatio {
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
func AggregateTestResults(folderPath string) (*TestReport, error) {
	var (
		// Map to hold unique tests based on their TestName and TestPackage
		// Key: TestName|TestPackage, Value: TestResult
		testMap       = make(map[string]TestResult)
		fullReport    = &TestReport{}
		excludedTests = map[string]struct{}{}
	)

	// Read all JSON files in the folder
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
			var report TestReport
			if jsonErr := json.Unmarshal(data, &report); jsonErr != nil {
				return jsonErr
			}
			if fullReport.GoProject == "" {
				fullReport.GoProject = report.GoProject
			} else if fullReport.GoProject != report.GoProject {
				return fmt.Errorf("multiple projects found in the results folder, expected %s, got %s", fullReport.GoProject, report.GoProject)
			}
			fullReport.TestRunCount += report.TestRunCount
			fullReport.RaceDetection = report.RaceDetection && fullReport.RaceDetection
			for _, test := range report.ExcludedTests {
				excludedTests[test] = struct{}{}
			}
			// Process each test result
			for _, result := range report.Results {
				// Unique key for each test based on TestName and TestPackage
				key := result.TestName + "|" + result.TestPackage
				if existingResult, found := testMap[key]; found {
					// Aggregate runs, durations, and outputs
					existingResult.Runs = existingResult.Runs + result.Runs
					existingResult.Durations = append(existingResult.Durations, result.Durations...)
					existingResult.Outputs = append(existingResult.Outputs, result.Outputs...)
					existingResult.PackageOutputs = append(existingResult.PackageOutputs, result.PackageOutputs...)
					existingResult.Successes += result.Successes
					existingResult.Failures += result.Failures
					existingResult.Panics += result.Panics
					existingResult.Races += result.Races
					existingResult.Skips += result.Skips
					existingResult.PassRatio = 1.0
					if existingResult.Runs > 0 {
						existingResult.PassRatio = float64(existingResult.Successes) / float64(existingResult.Runs)
					}

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
	// Aggregate
	for test := range excludedTests {
		fullReport.ExcludedTests = append(fullReport.ExcludedTests, test)
	}

	var (
		aggregatedResults = make([]TestResult, 0, len(testMap))
		allSuccesses      int
	)
	for _, result := range testMap {
		aggregatedResults = append(aggregatedResults, result)
		allSuccesses += result.Successes
	}

	// Sort by PassRatio in ascending order
	sort.Slice(aggregatedResults, func(i, j int) bool {
		return aggregatedResults[i].PassRatio < aggregatedResults[j].PassRatio
	})
	fullReport.Results = aggregatedResults

	return fullReport, nil
}

// PrintTests prints tests in a pretty format
func PrintTests(w io.Writer, tests []TestResult, maxPassRatio float64) (allRuns, passes, fails, skips, races, panics, flakes int) {
	fmt.Fprintln(w, "| Test Name | Test Package | Pass Ratio | Skipped | Runs | Successes | Failures | Panics | Races | Skips | Avg Duration |")
	fmt.Fprintln(w, "|-----------|--------------|------------|---------|------|-----------|----------|--------|-------|-------|--------------|")
	for _, test := range tests {
		if test.PassRatio >= maxPassRatio {
			continue
		}
		fmt.Fprintf(w, "| %s | %s | %.2f%% | %v | %d | %d | %d | %d | %d | %d | %s |\n",
			test.TestName, test.TestPackage, test.PassRatio*100, test.Skipped, test.Runs, test.Successes, test.Failures, test.Panics, test.Races, test.Skips, avgDuration(test.Durations).String())
		allRuns += test.Runs
		passes += test.Successes
		fails += test.Failures
		skips += test.Skips
		races += test.Races
		panics += test.Panics
		flakes += fails + races + panics
	}
	return
}

// TestsSummary builds a summary of test results in markdown format, handy for reporting in CI and Slack
func TestsSummary(w io.Writer, testReport *TestReport, maxPassRatio float64) {
	tests := testReport.Results
	fmt.Fprintln(w, "# Flakeguard Summary")
	fmt.Fprintln(w, "| **Setting** | **Value** |")
	fmt.Fprintln(w, "|-------------|-----------|")
	fmt.Fprintf(w, "| Go Project | %s |\n", testReport.GoProject)
	fmt.Fprintf(w, "| Max Pass Ratio | %.2f%% |\n", maxPassRatio*100)
	fmt.Fprintf(w, "| Test Run Count | %d |\n", testReport.TestRunCount)
	fmt.Fprintf(w, "| Race Detection | %t |\n", testReport.RaceDetection)
	fmt.Fprintf(w, "| Excluded Tests | %s |\n", strings.Join(testReport.ExcludedTests, ", "))
	fmt.Fprintln(w, "|-------------|-----------|")
	if len(tests) == 0 {
		fmt.Fprintln(w, "## No tests ran :warning:")
		return
	}
	var (
		avgPassRatio = 1.0
		testsData    = bytes.NewBuffer(nil)
	)
	for _, test := range tests {
		fmt.Fprintf(testsData, "| %s | %s | %.2f%% | %v | %d | %d | %d | %d | %d | %d | %s |\n",
			test.TestName, test.TestPackage, test.PassRatio*100, test.Skipped, test.Runs, test.Successes,
			test.Failures, test.Panics, test.Races, test.Skips, avgDuration(test.Durations).String(),
		)
	}
	allRuns, passes, fails, skips, races, panics, flakes := PrintTests(testsData, tests, maxPassRatio)
	if allRuns > 0 {
		avgPassRatio = float64(passes) / float64(allRuns)
	}
	if avgPassRatio < maxPassRatio {
		fmt.Fprintln(w, "## Found Flaky Tests :x:")
	} else {
		fmt.Fprintln(w, "## No Flakes Found :white_check_mark:")
	}
	fmt.Fprintf(w, "Ran `%d` tests `%d` times with a `%.2f%%` pass ratio and found `%d` flaky tests\n", len(tests), allRuns, avgPassRatio*100, flakes)
	fmt.Fprintf(w, "### Results")
	separator := "|----------------|------------------|--------------------|------------------|-----------------|-----------------|"
	fmt.Fprintln(w, "| **Total Runs** | **Total Passes** | **Total Failures** | **Total Panics** | **Total Races** | **Total Skips** |")
	fmt.Fprintln(w, separator)
	fmt.Fprintf(w, "| %d | %d | %d | %d | %d | %d |\n", allRuns, passes, fails, panics, races, skips)
	fmt.Fprintln(w, separator)
	if avgPassRatio < maxPassRatio {
		fmt.Fprintln(w, "### Flakes")
		separator = "|---------------|------------------|----------------|-------------|----------|---------------|--------------|------------|-----------|-----------|------------------|"
		fmt.Fprintln(w, "| **Test Name** | **Test Package** | **Pass Ratio** | **Skipped** | **Runs** | **Successes** | **Failures** | **Panics** | **Races** | **Skips** | **Avg Duration** |")
		fmt.Fprintln(w, separator)
		fmt.Fprint(w, testsData.String())
		fmt.Fprintln(w, separator)
	}
}

// Helper function to save filtered results and logs to specified paths
func SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath string, report *TestReport) error {
	if outputResultsPath != "" {
		if err := os.MkdirAll(filepath.Dir(outputResultsPath), 0755); err != nil { //nolint:gosec
			return fmt.Errorf("error creating output directory: %w", err)
		}
		jsonFileName := strings.TrimSuffix(outputResultsPath, filepath.Ext(outputResultsPath)) + ".json"
		mdFileName := strings.TrimSuffix(outputResultsPath, filepath.Ext(outputResultsPath)) + ".md"
		if err := saveReportNoLogs(jsonFileName, report); err != nil {
			return fmt.Errorf("error writing filtered results to file: %w", err)
		}
		summaryFile, err := os.Open(mdFileName)
		if err != nil {
			return fmt.Errorf("error opening markdown file: %w", err)
		}
		defer summaryFile.Close()
		TestsSummary(summaryFile, report, 1.0)
		fmt.Printf("Test results saved to %s and summary to %s\n", jsonFileName, mdFileName)
	} else {
		fmt.Println("No failed tests found based on the specified threshold and min pass ratio.")
	}

	if outputLogsPath != "" {
		if err := os.MkdirAll(filepath.Dir(outputLogsPath), 0755); err != nil { //nolint:gosec
			return fmt.Errorf("error creating output directory: %w", err)
		}
		if err := saveReport(outputLogsPath, report); err != nil {
			return fmt.Errorf("error writing filtered logs to file: %w", err)
		}
		fmt.Printf("Test logs saved to %s\n", outputLogsPath)
	}
	return nil
}

// saveReportNoLogs saves the test results to JSON without logs
// as outputs can take up a lot of space and are not always needed.
// Outputs can be saved separately using saveTestOutputs
func saveReportNoLogs(filePath string, report *TestReport) error {
	var filteredResults []TestResult
	for _, r := range report.Results {
		filteredResult := r
		filteredResult.Outputs = nil
		filteredResult.PackageOutputs = nil
		filteredResults = append(filteredResults, filteredResult)
	}
	report.Results = filteredResults

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling results: %v", err)
	}
	return os.WriteFile(filePath, data, 0644) //nolint:gosec
}

// saveReport saves the test results to JSON
func saveReport(filePath string, report *TestReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling outputs: %v", err)
	}
	return os.WriteFile(filePath, data, 0644) //nolint:gosec
}

// avgDuration calculates the average duration from a slice of time.Duration
func avgDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
