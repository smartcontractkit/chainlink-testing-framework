package reports

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GenerateFlakyTestsTable(
	results []TestResult,
	expectedPassRatio float64,
	markdown bool,
) [][]string {
	p := message.NewPrinter(language.English)
	sortTestResults(results)

	// Headers in the requested order
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
		"Code Owners",
	}

	// Format headers for Markdown if needed
	if markdown {
		for i, header := range headers {
			headers[i] = fmt.Sprintf("**%s**", header)
		}
	}

	// Initialize the table with headers
	table := [][]string{headers}

	for _, result := range results {
		// Exclude skipped tests and only include tests below the expected pass ratio
		if !result.Skipped && result.PassRatio < expectedPassRatio {
			row := []string{
				result.TestName,
				formatPassRatio(result.PassRatio),
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

			// Add code owners
			owners := "Unknown"
			if len(result.CodeOwners) > 0 {
				owners = strings.Join(result.CodeOwners, ", ")
			}
			row = append(row, owners)

			table = append(table, row)
		}
	}
	return table
}

func formatPassRatio(passRatio float64) string {
	if passRatio < 0 {
		return "N/A" // Handle undefined pass ratios (e.g., skipped tests)
	}
	return fmt.Sprintf("%.2f%%", passRatio*100)
}

func GenerateGitHubSummaryMarkdown(w io.Writer, testReport *TestReport, maxPassRatio float64) {
	settingsTable := buildSettingsTable(testReport, maxPassRatio)
	fmt.Fprint(w, "# Flakeguard Summary\n\n")
	printTable(w, settingsTable)
	fmt.Fprintln(w)

	if len(testReport.Results) == 0 {
		fmt.Fprintln(w, "## No tests ran :warning:")
		return
	}

	summary := GenerateSummaryData(testReport.Results, maxPassRatio)
	if summary.FlakyTests > 0 {
		fmt.Fprintln(w, "## Found Flaky Tests :x:")
	} else {
		fmt.Fprintln(w, "## No Flakes Found :white_check_mark:")
	}

	RenderResults(w, testReport.Results, maxPassRatio, true)
}

func GeneratePRCommentMarkdown(w io.Writer, testReport *TestReport, maxPassRatio float64, baseBranch, currentBranch, currentCommitSHA, repoURL, actionRunID string) {
	fmt.Fprint(w, "# Flakeguard Summary\n\n")

	// Construct additional info
	additionalInfo := fmt.Sprintf(
		"Ran new or updated tests between `%s` and %s (`%s`).",
		baseBranch,
		currentCommitSHA,
		currentBranch,
	)

	// Construct the links
	viewDetailsLink := fmt.Sprintf("[View Flaky Detector Details](%s/actions/runs/%s)", repoURL, actionRunID)
	compareChangesLink := fmt.Sprintf("[Compare Changes](%s/compare/%s...%s#files_bucket)", repoURL, baseBranch, currentCommitSHA)
	linksLine := fmt.Sprintf("%s | %s", viewDetailsLink, compareChangesLink)

	// Include additional information
	fmt.Fprintln(w, additionalInfo)
	fmt.Fprintln(w) // Add an extra newline for formatting

	// Include the links
	fmt.Fprintln(w, linksLine)
	fmt.Fprintln(w) // Add an extra newline for formatting

	if len(testReport.Results) == 0 {
		fmt.Fprintln(w, "## No tests ran :warning:")
		return
	}

	// Add the flaky tests section
	if GenerateSummaryData(testReport.Results, maxPassRatio).FlakyTests > 0 {
		fmt.Fprintln(w, "## Found Flaky Tests :x:")
	} else {
		fmt.Fprintln(w, "## No Flakes Found :white_check_mark:")
	}

	resultsTable := GenerateFlakyTestsTable(testReport.Results, maxPassRatio, true)
	renderTestResultsTable(w, resultsTable, true)
}

func buildSettingsTable(testReport *TestReport, maxPassRatio float64) [][]string {
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
	return rows
}

func RenderResults(
	w io.Writer,
	tests []TestResult,
	maxPassRatio float64,
	markdown bool,
) {
	resultsTable := GenerateFlakyTestsTable(tests, maxPassRatio, markdown)
	summary := GenerateSummaryData(tests, maxPassRatio)
	renderSummaryTable(w, summary, markdown)
	renderTestResultsTable(w, resultsTable, markdown)
}

func renderSummaryTable(w io.Writer, summary SummaryData, markdown bool) {
	summaryData := [][]string{
		{"Category", "Total"},
		{"Tests", fmt.Sprintf("%d", summary.TotalTests)},
		{"Panicked Tests", fmt.Sprintf("%d", summary.PanickedTests)},
		{"Raced Tests", fmt.Sprintf("%d", summary.RacedTests)},
		{"Flaky Tests", fmt.Sprintf("%d", summary.FlakyTests)},
		{"Flaky Test Ratio", summary.FlakyTestRatio},
		{"Runs", fmt.Sprintf("%d", summary.TotalRuns)},
		{"Passes", fmt.Sprintf("%d", summary.PassedRuns)},
		{"Failures", fmt.Sprintf("%d", summary.FailedRuns)},
		{"Skips", fmt.Sprintf("%d", summary.SkippedRuns)},
		{"Pass Ratio", summary.PassRatio},
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
	printTable(w, summaryData)
	fmt.Fprintln(w)
}

func renderTestResultsTable(w io.Writer, table [][]string, markdown bool) {
	if len(table) <= 1 {
		fmt.Fprintln(w, "No tests found under the specified pass ratio threshold.")
		return
	}
	printTable(w, table)
}

func printTable(w io.Writer, table [][]string) {
	colWidths := calculateColumnWidths(table)
	separator := buildSeparator(colWidths)

	for i, row := range table {
		printRow(w, row, colWidths)
		if i == 0 {
			fmt.Fprintln(w, separator)
		}
	}
}

func calculateColumnWidths(table [][]string) []int {
	colWidths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	return colWidths
}

func buildSeparator(colWidths []int) string {
	var buffer bytes.Buffer
	for _, width := range colWidths {
		buffer.WriteString("|-")
		buffer.WriteString(strings.Repeat("-", width))
		buffer.WriteString("-")
	}
	buffer.WriteString("|")
	return buffer.String()
}

func printRow(w io.Writer, row []string, colWidths []int) {
	var buffer bytes.Buffer
	for i, cell := range row {
		buffer.WriteString(fmt.Sprintf("| %-*s ", colWidths[i], cell))
	}
	buffer.WriteString("|")
	fmt.Fprintln(w, buffer.String())
}
