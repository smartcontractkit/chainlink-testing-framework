package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Aggregate test results and generate reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := reports.OSFileSystem{}

		// Get flag values directly using cmd.Flags().Get* methods
		reportResultsPath, _ := cmd.Flags().GetString("results-path")
		reportOutputPath, _ := cmd.Flags().GetString("output-path")
		reportFormats, _ := cmd.Flags().GetString("format")
		reportMaxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		reportCodeOwnersPath, _ := cmd.Flags().GetString("codeowners-path")
		reportRepoPath, _ := cmd.Flags().GetString("repo-path")

		// Split the formats into a slice
		formats := strings.Split(reportFormats, ",")

		// Start spinner for loading test reports
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Loading test reports..."
		s.Start()

		// Load test reports from JSON files
		testReports, err := reports.LoadReports(reportResultsPath)
		if err != nil {
			s.Stop()
			return fmt.Errorf("error loading test reports: %w", err)
		}
		s.Stop()
		fmt.Println("Test reports loaded successfully.")

		// Start spinner for aggregating reports
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Aggregating test reports..."
		s.Start()

		// Aggregate the reports
		aggregatedReport, err := reports.Aggregate(testReports...)
		if err != nil {
			s.Stop()
			return fmt.Errorf("error aggregating test reports: %w", err)
		}
		s.Stop()
		fmt.Println("Test reports aggregated successfully.")

		// Start spinner for mapping test results to paths
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Mapping test results to paths..."
		s.Start()

		// Map test results to test paths
		err = reports.MapTestResultsToPaths(aggregatedReport, reportRepoPath)
		if err != nil {
			s.Stop()
			return fmt.Errorf("error mapping test results to paths: %w", err)
		}
		s.Stop()
		fmt.Println("Test results mapped to paths successfully.")

		// Map test results to code owners if codeOwnersPath is provided
		if reportCodeOwnersPath != "" {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " Mapping test results to code owners..."
			s.Start()

			err = reports.MapTestResultsToOwners(aggregatedReport, reportCodeOwnersPath)
			if err != nil {
				s.Stop()
				return fmt.Errorf("error mapping test results to code owners: %w", err)
			}
			s.Stop()
			fmt.Println("Test results mapped to code owners successfully.")
		}

		// Exclude outputs and package outputs from the aggregated report of all tests
		for i := range aggregatedReport.Results {
			aggregatedReport.Results[i].Outputs = nil
			aggregatedReport.Results[i].PackageOutputs = nil
		}

		// Create output directory if it doesn't exist
		outputDir := reportOutputPath
		if err := fs.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}

		// Save the aggregated report (all tests)
		allTestsReportPath := filepath.Join(outputDir, "all-tests-report.json")
		if err := reports.SaveReport(fs, allTestsReportPath, *aggregatedReport); err != nil {
			return fmt.Errorf("error saving all tests report: %w", err)
		}
		fmt.Printf("All tests report saved to %s\n", allTestsReportPath)

		// Generate and save the reports (all tests) in specified formats
		for _, format := range formats {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = fmt.Sprintf(" Generating all tests report in format %s...", format)
			s.Start()

			if err := generateReport(aggregatedReport, format, filepath.Join(outputDir, "all-tests")); err != nil {
				s.Stop()
				return fmt.Errorf("error generating all tests report in format %s: %w", format, err)
			}
			s.Stop()
			fmt.Printf("All tests report in format %s generated successfully.\n", format)
		}

		// Filter failed tests (PassRatio < maxPassRatio and not skipped)
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Filtering failed tests..."
		s.Start()

		failedTests := reports.FilterTests(aggregatedReport.Results, func(tr reports.TestResult) bool {
			return !tr.Skipped && tr.PassRatio < reportMaxPassRatio
		})
		s.Stop()
		fmt.Println("Failed tests filtered successfully.")

		// For failed tests, include outputs and package outputs
		for i := range failedTests {
			// Retrieve outputs and package outputs from original reports
			failedTests[i].Outputs = getOriginalOutputs(testReports, failedTests[i].TestName, failedTests[i].TestPackage)
			failedTests[i].PackageOutputs = getOriginalPackageOutputs(testReports, failedTests[i].TestName, failedTests[i].TestPackage)
		}

		// Create a new report for failed tests
		failedReport := &reports.TestReport{
			GoProject:     aggregatedReport.GoProject,
			TestRunCount:  aggregatedReport.TestRunCount,
			RaceDetection: aggregatedReport.RaceDetection,
			ExcludedTests: aggregatedReport.ExcludedTests,
			SelectedTests: aggregatedReport.SelectedTests,
			Results:       failedTests,
		}

		// Save the failed tests report
		failedTestsReportPath := filepath.Join(outputDir, "failed-tests-report.json")
		if err := reports.SaveReport(fs, failedTestsReportPath, *failedReport); err != nil {
			return fmt.Errorf("error saving failed tests report: %w", err)
		}
		fmt.Printf("Failed tests report saved to %s\n", failedTestsReportPath)

		// Generate and save the reports for failed tests in specified formats
		for _, format := range formats {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = fmt.Sprintf(" Generating failed tests report in format %s...", format)
			s.Start()

			if err := generateReport(failedReport, format, filepath.Join(outputDir, "failed-tests")); err != nil {
				s.Stop()
				return fmt.Errorf("error generating failed tests report in format %s: %w", format, err)
			}
			s.Stop()
			fmt.Printf("Failed tests report in format %s generated successfully.\n", format)
		}

		fmt.Printf("Reports generated at: %s\n", reportOutputPath)

		return nil
	},
}

func init() {
	ReportCmd.Flags().StringP("results-path", "p", "", "Path to the folder containing JSON test result files (required)")
	ReportCmd.Flags().StringP("output-path", "o", "./report", "Path to output the generated report files")
	ReportCmd.Flags().StringP("format", "f", "markdown,json", "Comma-separated list of report formats (markdown,json)")
	ReportCmd.Flags().Float64P("max-pass-ratio", "", 1.0, "The maximum pass ratio threshold for a test to be considered flaky")
	ReportCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")
	ReportCmd.Flags().StringP("repo-path", "", ".", "The path to the root of the repository/project")
	ReportCmd.MarkFlagRequired("results-path")
}

func generateReport(report *reports.TestReport, format, outputPath string) error {
	fs := reports.OSFileSystem{}
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "markdown":
		mdFileName := outputPath + ".md"
		mdFile, err := fs.Create(mdFileName)
		if err != nil {
			return fmt.Errorf("error creating markdown file: %w", err)
		}
		defer mdFile.Close()
		reports.GenerateMarkdownSummary(mdFile, report, 1.0)
	case "json":
		jsonFileName := outputPath + ".json"
		if err := reports.SaveReportNoLogs(fs, jsonFileName, *report); err != nil {
			return fmt.Errorf("error saving JSON report: %w", err)
		}
	default:
		return fmt.Errorf("unsupported report format: %s", format)
	}

	// Generate summary JSON
	summaryData := reports.GenerateSummaryData(report.Results, 1.0)
	summaryFileName := outputPath + "-summary.json"
	if err := reports.SaveSummaryAsJSON(fs, summaryFileName, summaryData); err != nil {
		return fmt.Errorf("error saving summary JSON: %w", err)
	}

	return nil
}

// Helper functions to retrieve original outputs and package outputs
func getOriginalOutputs(reports []*reports.TestReport, testName, testPackage string) []string {
	for _, report := range reports {
		for _, result := range report.Results {
			if result.TestName == testName && result.TestPackage == testPackage {
				return result.Outputs
			}
		}
	}
	return nil
}

func getOriginalPackageOutputs(reports []*reports.TestReport, testName, testPackage string) []string {
	for _, report := range reports {
		for _, result := range report.Results {
			if result.TestName == testName && result.TestPackage == testPackage {
				return result.PackageOutputs
			}
		}
	}
	return nil
}
