package reports

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGenerateFlakyTestsTable(t *testing.T) {
	results := []TestResult{
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
	}

	expectedPassRatio := 0.9
	markdown := false

	table := GenerateFlakyTestsTable(results, expectedPassRatio, markdown)

	// Verify headers
	expectedHeaders := []string{
		"Name", "Pass Ratio", "Panicked?", "Timed Out?", "Race?", "Runs",
		"Successes", "Failures", "Skips", "Package", "Package Panicked?",
		"Avg Duration", "Code Owners",
	}
	if !reflect.DeepEqual(table[0], expectedHeaders) {
		t.Errorf("Expected headers %+v, got %+v", expectedHeaders, table[0])
	}

	// Verify rows (only TestFlaky should appear)
	if len(table) != 2 { // 1 header row + 1 data row
		t.Fatalf("Expected table length 2 (headers + 1 row), got %d", len(table))
	}

	expectedRow := []string{
		"TestFlaky",
		"50.00%",
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
	if !reflect.DeepEqual(table[1], expectedRow) {
		t.Errorf("Expected row %+v, got %+v", expectedRow, table[1])
	}
}

// TestGenerateGitHubSummaryMarkdown tests the GenerateGitHubSummaryMarkdown function.
func TestGenerateGitHubSummaryMarkdown(t *testing.T) {
	testReport := &TestReport{
		GoProject:     "ProjectX",
		TestRunCount:  3,
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
	}

	var buffer bytes.Buffer
	maxPassRatio := 0.9

	GenerateGitHubSummaryMarkdown(&buffer, testReport, maxPassRatio, "", "")

	output := buffer.String()

	// Check that the summary includes the expected headings
	if !strings.Contains(output, "# Flakeguard Summary") {
		t.Error("Expected markdown summary to contain '# Flakeguard Summary'")
	}
	if !strings.Contains(output, "## Found Flaky Tests :x:") {
		t.Error("Expected markdown summary to contain '## Found Flaky Tests :x:'")
	}
	if !strings.Contains(output, "| **Name**") {
		t.Error("Expected markdown table headers for test results")
	}
	if !strings.Contains(output, "| TestA ") {
		t.Error("Expected markdown table to include TestA")
	}
	if strings.Contains(output, "| TestB ") {
		t.Error("Did not expect markdown table to include TestB since its pass ratio is above the threshold")
	}
}

// TestGeneratePRCommentMarkdown tests the GeneratePRCommentMarkdown function.
func TestGeneratePRCommentMarkdown(t *testing.T) {
	testReport := &TestReport{
		GoProject:     "ProjectX",
		TestRunCount:  3,
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
	}

	var buffer bytes.Buffer
	maxPassRatio := 0.9
	baseBranch := "develop"
	currentBranch := "feature-branch"
	currentCommitSHA := "abcdef1234567890"
	repoURL := "https://github.com/example/repo"
	actionRunID := "123456789"

	GeneratePRCommentMarkdown(&buffer, testReport, maxPassRatio, baseBranch, currentBranch, currentCommitSHA, repoURL, actionRunID, "", "")

	output := buffer.String()

	// Check that the output includes the expected headings and links
	if !strings.Contains(output, "# Flakeguard Summary") {
		t.Error("Expected markdown summary to contain '# Flakeguard Summary'")
	}
	if !strings.Contains(output, fmt.Sprintf("Ran new or updated tests between `%s` and %s (`%s`).", baseBranch, currentCommitSHA, currentBranch)) {
		t.Error("Expected markdown to contain the additional info line with branches and commit SHA")
	}
	if !strings.Contains(output, fmt.Sprintf("[View Flaky Detector Details](%s/actions/runs/%s)", repoURL, actionRunID)) {
		t.Error("Expected markdown to contain the 'View Flaky Detector Details' link")
	}
	if !strings.Contains(output, fmt.Sprintf("[Compare Changes](%s/compare/%s...%s#files_bucket)", repoURL, baseBranch, currentCommitSHA)) {
		t.Error("Expected markdown to contain the 'Compare Changes' link")
	}
	if !strings.Contains(output, "## Found Flaky Tests :x:") {
		t.Error("Expected markdown summary to contain '## Found Flaky Tests :x:'")
	}
	if !strings.Contains(output, "| **Name**") {
		t.Error("Expected markdown table headers for test results")
	}
	if !strings.Contains(output, "| TestA ") {
		t.Error("Expected markdown table to include TestA")
	}
	if strings.Contains(output, "| TestB ") {
		t.Error("Did not expect markdown table to include TestB since its pass ratio is above the threshold")
	}
}

// TestPrintTable tests the printTable function.
func TestPrintTable(t *testing.T) {
	table := [][]string{
		{"Header1", "Header2", "Header3"},
		{"Row1Col1", "Row1Col2", "Row1Col3"},
		{"Row2Col1", "Row2Col2", "Row2Col3"},
	}

	var buffer bytes.Buffer
	printTable(&buffer, table)

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
		testResults            []TestResult
		maxPassRatio           float64
		expectedSummary        SummaryData
		expectedStringsContain []string
	}{
		{
			name: "single flaky test",
			testResults: []TestResult{
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
			maxPassRatio: 0.9,
			expectedSummary: SummaryData{
				TotalTests:     1,
				PanickedTests:  0,
				RacedTests:     0,
				FlakyTests:     1,
				FlakyTestRatio: "100%",
				TotalRuns:      4,
				PassedRuns:     3,
				FailedRuns:     1,
				SkippedRuns:    0,
				PassRatio:      "75%",
				MaxPassRatio:   0.9,
			},
			expectedStringsContain: []string{"Test1", "package1", "75%", "false", "1.05s", "4", "0"},
		},
		// Add more test cases as needed
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			RenderResults(&buf, tc.testResults, tc.maxPassRatio, false)
			output := buf.String()

			// Generate the summary data
			summary := GenerateSummaryData(tc.testResults, tc.maxPassRatio)

			// Verify summary data
			if summary.TotalTests != tc.expectedSummary.TotalTests {
				t.Errorf("Expected TotalTests %v, got %v", tc.expectedSummary.TotalTests, summary.TotalTests)
			}
			if summary.TotalRuns != tc.expectedSummary.TotalRuns {
				t.Errorf("Expected TotalRuns %v, got %v", tc.expectedSummary.TotalRuns, summary.TotalRuns)
			}
			if summary.PassedRuns != tc.expectedSummary.PassedRuns {
				t.Errorf("Expected PassedRuns %v, got %v", tc.expectedSummary.PassedRuns, summary.PassedRuns)
			}
			if summary.FailedRuns != tc.expectedSummary.FailedRuns {
				t.Errorf("Expected FailedRuns %v, got %v", tc.expectedSummary.FailedRuns, summary.FailedRuns)
			}
			if summary.FlakyTests != tc.expectedSummary.FlakyTests {
				t.Errorf("Expected FlakyTests %v, got %v", tc.expectedSummary.FlakyTests, summary.FlakyTests)
			}
			if summary.PassRatio != tc.expectedSummary.PassRatio {
				t.Errorf("Expected PassRatio %v, got %v", tc.expectedSummary.PassRatio, summary.PassRatio)
			}

			// Verify output content
			for _, expected := range tc.expectedStringsContain {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but it did not", expected)
				}
			}
		})
	}
}
