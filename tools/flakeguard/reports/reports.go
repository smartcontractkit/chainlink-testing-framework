package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	SelectedTests []string
	Results       []TestResult
}

// TestResult represents the result of a single test being run through flakeguard.
type TestResult struct {
	TestName       string
	TestPackage    string
	PackagePanic   bool            // Indicates a package-level panic
	Panic          bool            // Indicates a test-level panic
	Timeout        bool            // Indicates if the test timed out
	Race           bool            // Indicates if the test caused a data race
	Skipped        bool            // Indicates if the test was skipped
	PassRatio      float64         // Pass ratio in decimal format like 0.5
	Runs           int             // Count of how many times the test was run
	Failures       int             // Count of how many times the test failed
	Successes      int             // Count of how many times the test passed
	Skips          int             // Count of how many times the test was skipped
	Outputs        []string        `json:"outputs,omitempty"` // Stores outputs for a test
	Durations      []time.Duration // Stores elapsed time for each run of the test
	PackageOutputs []string        `json:"package_outputs,omitempty"` // Stores package-level outputs
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
		selectedTests = map[string]struct{}{}
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
			for _, test := range report.SelectedTests {
				selectedTests[test] = struct{}{}
			}
			// Process each test results
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
					existingResult.Panic = existingResult.Panic || result.Panic
					existingResult.Race = existingResult.Race || result.Race
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
	for test := range selectedTests {
		fullReport.SelectedTests = append(fullReport.SelectedTests, test)
	}

	var (
		aggregatedResults = make([]TestResult, 0, len(testMap))
		allSuccesses      int
	)
	for _, result := range testMap {
		aggregatedResults = append(aggregatedResults, result)
		allSuccesses += result.Successes
	}

	sortTestResults(aggregatedResults)
	fullReport.Results = aggregatedResults

	return fullReport, nil
}

// PrintTests prints tests in a pretty format
func PrintTests(
	w io.Writer,
	tests []TestResult,
	maxPassRatio float64,
) (runs, passes, fails, skips, panickedTests, racedTests, flakyTests int) {
	sortTestResults(tests)
	headers := []string{
		"**Test**",
		"**Pass Ratio**",
		"**Runs**",
		"**Panicked?**",
		"**Timed Out?**",
		"**Race?**",
		"**Successes**",
		"**Failures**",
		"**Skips**",
		"**Package**",
		"**Package Panicked?**",
		"**Avg Duration**",
	}

	// Build test rows and summary data
	rows := [][]string{}
	for _, test := range tests {
		if test.PassRatio < maxPassRatio {
			rows = append(rows, []string{
				test.TestName,
				fmt.Sprintf("%.2f%%", test.PassRatio*100),
				fmt.Sprintf("%d", test.Runs),
				fmt.Sprintf("%t", test.Panic),
				fmt.Sprintf("%t", test.Timeout),
				fmt.Sprintf("%t", test.Race),
				fmt.Sprintf("%d", test.Successes),
				fmt.Sprintf("%d", test.Failures),
				fmt.Sprintf("%d", test.Skips),
				test.TestPackage,
				fmt.Sprintf("%t", test.PackagePanic),
				avgDuration(test.Durations).String(),
			})
		}

		runs += test.Runs
		passes += test.Successes
		fails += test.Failures
		skips += test.Skips
		if test.Panic {
			panickedTests++
			flakyTests++
		} else if test.Race {
			racedTests++
			flakyTests++
		} else if test.PassRatio < maxPassRatio {
			flakyTests++
		}
	}

	var passRatioStr string
	if runs == 0 || passes == runs {
		passRatioStr = "100%"
	} else {
		percentage := float64(passes) / float64(runs) * 100
		truncatedPercentage := math.Floor(percentage*100) / 100 // Truncate to 2 decimal places
		passRatioStr = fmt.Sprintf("%.2f%%", truncatedPercentage)
	}

	// Print out summary data
	summaryData := [][]string{
		{"**Category**", "**Total**"},
		{"**Tests**", fmt.Sprint(len(tests))},
		{"**Panicked Tests**", fmt.Sprint(panickedTests)},
		{"**Raced Tests**", fmt.Sprint(racedTests)},
		{"**Flaky Tests**", fmt.Sprint(flakyTests)},
		{"**Pass Ratio**", passRatioStr},
		{"**Runs**", fmt.Sprint(runs)},
		{"**Passes**", fmt.Sprint(passes)},
		{"**Failures**", fmt.Sprint(fails)},
		{"**Skips**", fmt.Sprint(skips)},
	}
	colWidths := make([]int, len(summaryData[0]))

	for _, row := range summaryData {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	if len(rows) == 0 {
		fmt.Fprintf(w, "No tests found under pass ratio of %.2f%%\n", maxPassRatio*100)
		return
	}

	printRow := func(cells []string) {
		fmt.Fprintf(w, "| %-*s | %-*s |\n", colWidths[0], cells[0], colWidths[1], cells[1])
	}
	printSeparator := func() {
		fmt.Fprintf(w, "|-%s-|-%s-|\n", strings.Repeat("-", colWidths[0]), strings.Repeat("-", colWidths[1]))
	}
	printRow(summaryData[0])
	printSeparator()
	for _, row := range summaryData[1:] {
		printRow(row)
	}
	fmt.Fprintln(w)

	// Print out test data
	colWidths = make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	printRow = func(cells []string) {
		var buffer bytes.Buffer
		for i, cell := range cells {
			buffer.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], cell))
		}
		fmt.Fprintln(w, "|"+buffer.String())
	}

	printSeparator = func() {
		var buffer bytes.Buffer
		for _, width := range colWidths {
			buffer.WriteString(" " + strings.Repeat("-", width) + " |")
		}
		fmt.Fprintln(w, "|"+buffer.String())
	}

	printRow(headers)
	printSeparator()
	for _, row := range rows {
		printRow(row)
	}
	return
}

// MarkdownSummary builds a summary of test results in markdown format, handy for reporting in CI and Slack
func MarkdownSummary(w io.Writer, testReport *TestReport, maxPassRatio float64) {
	var (
		avgPassRatio = 1.0
		testsData    = bytes.NewBuffer(nil)
		tests        = testReport.Results
	)

	rows := [][]string{
		{"**Setting**", "**Value**"},
		{"Project", testReport.GoProject},
		{"Max Pass Ratio", fmt.Sprintf("%.2f%%", maxPassRatio*100)},
		{"Test Run Count", fmt.Sprintf("%d", testReport.TestRunCount)},
		{"Race Detection", fmt.Sprintf("%t", testReport.RaceDetection)},
	}
	if len(testReport.ExcludedTests) > 0 {
		rows = append(rows, []string{"Excluded Tests", strings.Join(testReport.ExcludedTests, ", ")})
	}
	if len(testReport.SelectedTests) > 0 {
		rows = append(rows, []string{"Selected Tests", strings.Join(testReport.SelectedTests, ", ")})
	}
	colWidths := make([]int, len(rows[0]))

	// Calculate column widths
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	printRow := func(cells []string) {
		fmt.Fprintf(w, "| %-*s | %-*s |\n", colWidths[0], cells[0], colWidths[1], cells[1])
	}
	printSeparator := func() {
		fmt.Fprintf(w, "|-%s-|-%s-|\n", strings.Repeat("-", colWidths[0]), strings.Repeat("-", colWidths[1]))
	}
	fmt.Fprint(w, "# Flakeguard Summary\n\n")
	// Print settings data
	printRow(rows[0])
	printSeparator()
	for _, row := range rows[1:] {
		printRow(row)
	}
	fmt.Fprintln(w)

	if len(tests) == 0 {
		fmt.Fprintln(w, "## No tests ran :warning:")
		return
	}

	allRuns, passes, _, _, _, _, _ := PrintTests(testsData, tests, maxPassRatio)
	if allRuns > 0 {
		avgPassRatio = float64(passes) / float64(allRuns)
	}
	if avgPassRatio < maxPassRatio {
		fmt.Fprint(w, "## Found Flaky Tests :x:\n\n")
	} else {
		fmt.Fprint(w, "## No Flakes Found :white_check_mark:\n\n")
	}
	fmt.Fprint(w, testsData.String())
}

// Helper function to save filtered results and logs to specified paths
func SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath string, report *TestReport) error {
	if outputResultsPath != "" {
		if err := os.MkdirAll(filepath.Dir(outputResultsPath), 0755); err != nil { //nolint:gosec
			return fmt.Errorf("error creating output directory: %w", err)
		}
		jsonFileName := strings.TrimSuffix(outputResultsPath, filepath.Ext(outputResultsPath)) + ".json"
		mdFileName := strings.TrimSuffix(outputResultsPath, filepath.Ext(outputResultsPath)) + ".md"
		// no pointer to avoid destroying the original report
		if err := saveReportNoLogs(jsonFileName, *report); err != nil {
			return fmt.Errorf("error writing filtered results to file: %w", err)
		}
		summaryFile, err := os.Create(mdFileName)
		if err != nil {
			return fmt.Errorf("error creating markdown file: %w", err)
		}
		defer summaryFile.Close()
		MarkdownSummary(summaryFile, report, 1.0)
		fmt.Printf("Test results saved to %s and summary to %s\n", jsonFileName, mdFileName)
	} else {
		fmt.Println("No failed tests found based on the specified threshold and min pass ratio.")
	}

	if outputLogsPath != "" {
		if err := os.MkdirAll(filepath.Dir(outputLogsPath), 0755); err != nil { //nolint:gosec
			return fmt.Errorf("error creating output directory: %w", err)
		}
		// no pointer to avoid destroying the original report
		if err := saveReport(outputLogsPath, *report); err != nil {
			return fmt.Errorf("error writing filtered logs to file: %w", err)
		}
		fmt.Printf("Test logs saved to %s\n", outputLogsPath)
	}
	return nil
}

// saveReportNoLogs saves the test results to JSON without logs
// as outputs can take up a lot of space and are not always needed.
// Outputs can be saved separately using saveTestOutputs
func saveReportNoLogs(filePath string, report TestReport) error {
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
func saveReport(filePath string, report TestReport) error {
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

// sortTestResults sorts results by TestPackage, TestName, and PassRatio for consistent comparison and pretty printing
func sortTestResults(results []TestResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].TestPackage != results[j].TestPackage {
			return results[i].TestPackage < results[j].TestPackage
		}
		if results[i].TestName != results[j].TestName {
			return results[i].TestName < results[j].TestName
		}
		return results[i].PassRatio < results[j].PassRatio
	})
}
