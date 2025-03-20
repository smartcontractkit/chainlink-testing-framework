package reports

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFlakyTestsTable(t *testing.T) {
	report := TestReport{
		Results: []TestResult{
			{
				TestName:    "TestFlaky",
				PassRatio:   0.5,
				Skipped:     false,
				Runs:        2,
				Successes:   1,
				Failures:    1,
				TestPackage: "pkg1",
				CodeOwners:  []string{"owner1"},
			},
			{
				TestName:    "TestSkipped",
				PassRatio:   -1.0,
				Skipped:     true,
				Runs:        0,
				Skips:       1,
				TestPackage: "pkg2",
				CodeOwners:  []string{"owner2"},
			},
		},
		MaxPassRatio: 0.9,
	}

	markdown := false

	table := GenerateFlakyTestsTable(report, markdown)

	// Verify headers
	expectedHeaders := []string{
		"Name", "Pass Ratio", "Panicked?", "Timed Out?", "Race?", "Runs",
		"Successes", "Failures", "Skips", "Package", "Package Panicked?",
		"Avg Duration", "Code Owners",
	}
	assert.Equal(t, expectedHeaders, table[0], "Expected headers to match")

	// Verify rows (only TestFlaky should appear)
	assert.Len(t, table, 2, "Expected 2 rows in table (headers + 1 data row)")

	expectedRow := []string{
		"TestFlaky",
		"50%",
		"false",
		"false",
		"false",
		"2",
		"1",
		"1",
		"0",
		"pkg1",
		"false",
		"0s",
		"owner1",
	}
	assert.Equal(t, expectedRow, table[1], "Expected row to match")
}

func TestGenerateGitHubSummaryMarkdown(t *testing.T) {
	maxPassRatio := 0.9
	testReport := TestReport{
		GoProject:     "ProjectX",
		SummaryData:   &SummaryData{UniqueTestsRun: 2, FlakyTests: 1},
		RaceDetection: true,
		Results: []TestResult{
			{
				TestName:    "TestA",
				PassRatio:   0.8,
				Runs:        5,
				Successes:   4,
				Failures:    1,
				TestPackage: "pkg1",
				CodeOwners:  []string{"owner1"},
				Durations:   []time.Duration{time.Second, time.Second, time.Second, time.Second, time.Second},
			},
			{
				TestName:    "TestB",
				PassRatio:   1.0,
				Runs:        3,
				Successes:   3,
				Failures:    0,
				TestPackage: "pkg2",
				CodeOwners:  []string{"owner2"},
				Durations:   []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second},
			},
		},
		MaxPassRatio: maxPassRatio,
	}

	var buffer bytes.Buffer

	GenerateGitHubSummaryMarkdown(&buffer, testReport, maxPassRatio, "", "")

	output := buffer.String()

	// Check that the summary includes the expected headings
	assert.Contains(t, output, "# Flakeguard Summary", "Expected markdown summary to contain '# Flakeguard Summary'")
	assert.Contains(t, output, "## Found Flaky Tests :x:", "Expected markdown summary to contain '## Found Flaky Tests :x:'")
	assert.Contains(t, output, "| **Name**", "Expected markdown table headers for test results")
	assert.Contains(t, output, "| TestA ", "Expected markdown table to include TestA")
	assert.NotContains(t, output, "| TestB ", "Markdown table should not include TestB")
}

// TestGeneratePRCommentMarkdown tests the GeneratePRCommentMarkdown function.
func TestGeneratePRCommentMarkdown(t *testing.T) {
	maxPassRatio := 0.9
	testReport := TestReport{
		GoProject:     "ProjectX",
		SummaryData:   &SummaryData{UniqueTestsRun: 3, FlakyTests: 1},
		RaceDetection: true,
		MaxPassRatio:  maxPassRatio,
		Results: []TestResult{
			{
				TestName:    "TestA",
				PassRatio:   0.8,
				Runs:        5,
				Successes:   4,
				Failures:    1,
				TestPackage: "pkg1",
				CodeOwners:  []string{"owner1"},
				Durations:   []time.Duration{time.Second, time.Second, time.Second, time.Second, time.Second},
			},
			{
				TestName:    "TestB",
				PassRatio:   1.0,
				Runs:        3,
				Successes:   3,
				Failures:    0,
				TestPackage: "pkg2",
				CodeOwners:  []string{"owner2"},
				Durations:   []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second},
			},
		},
	}

	var buffer bytes.Buffer
	baseBranch := "develop"
	currentBranch := "feature-branch"
	currentCommitSHA := "abcdef1234567890"
	repoURL := "https://github.com/example/repo"
	actionRunID := "123456789"

	GeneratePRCommentMarkdown(&buffer, testReport, maxPassRatio, baseBranch, currentBranch, currentCommitSHA, repoURL, actionRunID, "", "")

	output := buffer.String()

	// Check that the output includes the expected headings and links
	assert.Contains(t, output, "# Flakeguard Summary", "Expected markdown summary to contain '# Flakeguard Summary'")
	assert.Contains(t, output, fmt.Sprintf("Ran new or updated tests between `%s` and %s (`%s`).", baseBranch, currentCommitSHA, currentBranch), "Expected markdown to contain the additional info line with branches and commit SHA")
	assert.Contains(t, output, fmt.Sprintf("[View Flaky Detector Details](%s/actions/runs/%s)", repoURL, actionRunID), "Expected markdown to contain the 'View Flaky Detector Details' link")
	assert.Contains(t, output, fmt.Sprintf("[Compare Changes](%s/compare/%s...%s#files_bucket)", repoURL, baseBranch, currentCommitSHA), "Expected markdown to contain the 'Compare Changes' link")
	assert.Contains(t, output, "## Found Flaky Tests :x:", "Expected markdown summary to contain '## Found Flaky Tests :x:'")
	assert.Contains(t, output, "| **Name**", "Expected markdown table headers for test results")
	assert.Contains(t, output, "| TestA ", "Expected markdown table to include TestA")
	assert.NotContains(t, output, "| TestB ", "Markdown table should not include TestB")
}

// TestPrintTable tests the printTable function.
func TestPrintTable(t *testing.T) {
	table := [][]string{
		{"Header1", "Header2", "Header3"},
		{"Row1Col1", "Row1Col2", "Row1Col3"},
		{"Row2Col1", "Row2Col2", "Row2Col3"},
	}

	var buffer bytes.Buffer
	printTable(&buffer, table, false)

	output := buffer.String()

	expected := `| Header1  | Header2  | Header3  |
|----------|----------|----------|
| Row1Col1 | Row1Col2 | Row1Col3 |
| Row2Col1 | Row2Col2 | Row2Col3 |
`

	if output != expected {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expected, output)
	}
}

func TestRenderResults(t *testing.T) {
	testcases := []struct {
		name                   string
		testReport             TestReport
		expectedSummary        *SummaryData
		expectedStringsContain []string
	}{
		{
			name: "single flaky test",
			testReport: TestReport{
				Results: []TestResult{
					{
						TestName:    "Test1",
						TestPackage: "package1",
						PassRatio:   0.75,
						Successes:   3,
						Failures:    1,
						Skipped:     false,
						Runs:        4,
						Durations: []time.Duration{
							time.Millisecond * 1200,
							time.Millisecond * 900,
							time.Millisecond * 1100,
							time.Second,
						},
					},
				},
				MaxPassRatio: 0.9,
			},
			expectedSummary: &SummaryData{
				UniqueTestsRun:   1,
				TestRunCount:     4,
				PanickedTests:    0,
				RacedTests:       0,
				FlakyTests:       1,
				FlakyTestPercent: "100%",
				TotalRuns:        4,
				PassedRuns:       3,
				FailedRuns:       1,
				SkippedRuns:      0,
				PassPercent:      "75%",
			},
			expectedStringsContain: []string{"Test1", "package1", "75%", "false", "1.05s", "4", "0"},
		},
		// Add more test cases as needed
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate the summary data
			tc.testReport.GenerateSummaryData()

			var buf bytes.Buffer
			RenderTestReport(&buf, tc.testReport, false, false)
			output := buf.String()

			// Verify summary data
			assert.Equal(t, tc.expectedSummary, tc.testReport.SummaryData, "Summary data does not match expected")

			// Verify output content
			for _, expected := range tc.expectedStringsContain {
				assert.Contains(t, output, expected, "Expected output to contain %q", expected)
			}
		})
	}
}
