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

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	TestPath       string          // Path to the test file
	CodeOwners     []string        // Owners of the test
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

// Aggregate aggregates multiple test reports into a single report.
func Aggregate(reportsToAggregate ...*TestReport) (*TestReport, error) {
	var (
		// Map to hold unique tests based on their TestName and TestPackage
		// Key: TestName|TestPackage, Value: TestResult
		testMap       = make(map[string]TestResult)
		fullReport    = &TestReport{}
		excludedTests = map[string]struct{}{}
		selectedTests = map[string]struct{}{}
	)

	// Read all JSON files in the folder
	for _, report := range reportsToAggregate {
		if fullReport.GoProject == "" {
			fullReport.GoProject = report.GoProject
		} else if fullReport.GoProject != report.GoProject {
			return nil, fmt.Errorf("reports with different Go projects found, expected %s, got %s", fullReport.GoProject, report.GoProject)
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

func TestResultsTable(
	results []TestResult,
	expectedPassRatio float64,
	includeCodeOwners bool,
	markdown bool,
) (resultsTable [][]string, runs, passes, fails, skips, panickedTests, racedTests, flakyTests int) {
	p := message.NewPrinter(language.English) // For formatting numbers
	sortTestResults(results)

	headers := []string{
		"Name",
		"Pass Ratio",
		"Panicked?",
		"Timed Out?",
		"Race?",
		"Runs",
		"Successes",
		"Failures",
		"Skips",
		"Package",
		"Package Panicked?",
		"Avg Duration",
	}

	if includeCodeOwners {
		headers = append(headers, "Code Owners")
	}
	if markdown {
		for i, header := range headers {
			headers[i] = fmt.Sprintf("**%s**", header)
		}
	}

	resultsTable = [][]string{}
	resultsTable = append(resultsTable, headers)
	for _, result := range results {
		if result.PassRatio < expectedPassRatio {
			row := []string{
				result.TestName,
				fmt.Sprintf("%.2f%%", result.PassRatio*100),
				fmt.Sprintf("%t", result.Panic),
				fmt.Sprintf("%t", result.Timeout),
				fmt.Sprintf("%t", result.Race),
				p.Sprintf("%d", result.Runs),
				p.Sprintf("%d", result.Successes),
				p.Sprintf("%d", result.Failures),
				p.Sprintf("%d", result.Skips),
				result.TestPackage,
				fmt.Sprintf("%t", result.PackagePanic),
				avgDuration(result.Durations).String(),
			}

			if includeCodeOwners {
				owners := "Unknown"
				if len(result.CodeOwners) > 0 {
					owners = strings.Join(result.CodeOwners, ", ")
				}
				row = append(row, owners)
			}

			resultsTable = append(resultsTable, row)
		}

		runs += result.Runs
		passes += result.Successes
		fails += result.Failures
		skips += result.Skips
		if result.Panic {
			panickedTests++
			flakyTests++
		} else if result.Race {
			racedTests++
			flakyTests++
		} else if result.PassRatio < expectedPassRatio {
			flakyTests++
		}
	}
	return
}

// PrintTests prints tests in a pretty format
func PrintResults(
	w io.Writer,
	tests []TestResult,
	maxPassRatio float64,
	markdown bool,
	includeCodeOwners bool, // Include code owners in the output. Set to true if test results have code owners
) (runs, passes, fails, skips, panickedTests, racedTests, flakyTests int) {
	var (
		resultsTable  [][]string
		passRatioStr  string
		flakeRatioStr string
		p             = message.NewPrinter(language.English) // For formatting numbers
	)
	resultsTable, runs, passes, fails, skips, panickedTests, racedTests, flakyTests = TestResultsTable(tests, maxPassRatio, markdown, includeCodeOwners)
	// Print out summary data
	if runs == 0 || passes == runs {
		passRatioStr = "100%"
		flakeRatioStr = "0%"
	} else {
		passPercentage := float64(passes) / float64(runs) * 100
		truncatedPassPercentage := math.Floor(passPercentage*100) / 100 // Truncate to 2 decimal places
		flakePercentage := float64(flakyTests) / float64(len(tests)) * 100
		truncatedFlakePercentage := math.Floor(flakePercentage*100) / 100 // Truncate to 2 decimal places
		passRatioStr = fmt.Sprintf("%.2f%%", truncatedPassPercentage)
		flakeRatioStr = fmt.Sprintf("%.2f%%", truncatedFlakePercentage)
	}
	summaryData := [][]string{
		{"Category", "Total"},
		{"Tests", p.Sprint(len(tests))},
		{"Panicked Tests", p.Sprint(panickedTests)},
		{"Raced Tests", p.Sprint(racedTests)},
		{"Flaky Tests", p.Sprint(flakyTests)},
		{"Flaky Test Ratio", flakeRatioStr},
		{"Runs", p.Sprint(runs)},
		{"Passes", p.Sprint(passes)},
		{"Failures", p.Sprint(fails)},
		{"Skips", p.Sprint(skips)},
		{"Pass Ratio", passRatioStr},
	}
	if markdown {
		for i, row := range summaryData {
			if i == 0 {
				summaryData[i] = []string{"**Category**", "**Total**"}
			} else {
				summaryData[i] = []string{fmt.Sprintf("**%s**", row[0]), row[1]}
			}
		}
	}

	colWidths := make([]int, len(summaryData[0]))

	for _, row := range summaryData {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	if len(resultsTable) <= 1 {
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
	resultsHeaders := resultsTable[0]
	colWidths = make([]int, len(resultsHeaders))
	for i, header := range resultsHeaders {
		colWidths[i] = len(header)
	}
	for rowNum := 1; rowNum < len(resultsTable); rowNum++ {
		for i, cell := range resultsTable[rowNum] {
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

	printRow(resultsHeaders)
	printSeparator()
	for rowNum := 1; rowNum < len(resultsTable); rowNum++ {
		printRow(resultsTable[rowNum])
	}
	return
}

// MarkdownSummary builds a summary of test results in markdown format, handy for reporting in CI and Slack
func MarkdownSummary(w io.Writer, testReport *TestReport, maxPassRatio float64, includeCodeOwners bool) {
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

	allRuns, passes, _, _, _, _, _ := PrintResults(testsData, tests, maxPassRatio, true, includeCodeOwners)
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
func SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath string, report *TestReport, includeCodeOwners bool) error {
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
		MarkdownSummary(summaryFile, report, 1.0, includeCodeOwners)
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
		// Compare TestPackage first
		if results[i].TestPackage != results[j].TestPackage {
			return results[i].TestPackage < results[j].TestPackage
		}

		// Split TestName into components for hierarchical comparison
		iParts := strings.Split(results[i].TestName, "/")
		jParts := strings.Split(results[j].TestName, "/")

		// Compare each part of the TestName hierarchically
		for k := 0; k < len(iParts) && k < len(jParts); k++ {
			if iParts[k] != jParts[k] {
				return iParts[k] < jParts[k]
			}
		}

		// If all compared parts are equal, the shorter name (parent) comes first
		if len(iParts) != len(jParts) {
			return len(iParts) < len(jParts)
		}

		// Finally, compare PassRatio if everything else is equal
		return results[i].PassRatio < results[j].PassRatio
	})
}
