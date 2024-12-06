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

		// Get flag values
		reportResultsPath, _ := cmd.Flags().GetString("results-path")
		reportOutputPath, _ := cmd.Flags().GetString("output-path")
		reportMaxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		reportCodeOwnersPath, _ := cmd.Flags().GetString("codeowners-path")
		reportRepoPath, _ := cmd.Flags().GetString("repo-path")
		generatePRComment, _ := cmd.Flags().GetBool("generate-pr-comment")

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

		// Generate GitHub summary markdown
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Generating GitHub summary markdown..."
		s.Start()

		err = generateGitHubSummaryMarkdown(aggregatedReport, filepath.Join(outputDir, "all-tests"))
		if err != nil {
			s.Stop()
			return fmt.Errorf("error generating GitHub summary markdown: %w", err)
		}
		s.Stop()
		fmt.Println("GitHub summary markdown generated successfully.")

		if generatePRComment {
			// Retrieve required flags
			currentBranch, _ := cmd.Flags().GetString("current-branch")
			currentCommitSHA, _ := cmd.Flags().GetString("current-commit-sha")
			baseBranch, _ := cmd.Flags().GetString("base-branch")
			repoURL, _ := cmd.Flags().GetString("repo-url")
			actionRunID, _ := cmd.Flags().GetString("action-run-id")

			// Validate that required flags are provided
			missingFlags := []string{}
			if currentBranch == "" {
				missingFlags = append(missingFlags, "--current-branch")
			}
			if currentCommitSHA == "" {
				missingFlags = append(missingFlags, "--current-commit-sha")
			}
			if repoURL == "" {
				missingFlags = append(missingFlags, "--repo-url")
			}
			if actionRunID == "" {
				missingFlags = append(missingFlags, "--action-run-id")
			}
			if len(missingFlags) > 0 {
				return fmt.Errorf("the following flags are required when --generate-pr-comment is set: %s", strings.Join(missingFlags, ", "))
			}

			// Generate PR comment markdown
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " Generating PR comment markdown..."
			s.Start()

			err = generatePRCommentMarkdown(aggregatedReport, filepath.Join(outputDir, "all-tests"), baseBranch, currentBranch, currentCommitSHA, repoURL, actionRunID)
			if err != nil {
				s.Stop()
				return fmt.Errorf("error generating PR comment markdown: %w", err)
			}
			s.Stop()
			fmt.Println("PR comment markdown generated successfully.")
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

		fmt.Printf("Reports generated at: %s\n", reportOutputPath)

		return nil
	},
}

func init() {
	ReportCmd.Flags().StringP("results-path", "p", "", "Path to the folder containing JSON test result files (required)")
	ReportCmd.Flags().StringP("output-path", "o", "./report", "Path to output the generated report files")
	ReportCmd.Flags().Float64P("max-pass-ratio", "", 1.0, "The maximum pass ratio threshold for a test to be considered flaky")
	ReportCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")
	ReportCmd.Flags().StringP("repo-path", "", ".", "The path to the root of the repository/project")
	ReportCmd.Flags().Bool("generate-pr-comment", false, "Set to true to generate PR comment markdown")
	ReportCmd.Flags().String("base-branch", "develop", "The base branch to compare against (used in PR comment)")
	ReportCmd.Flags().String("current-branch", "", "The current branch name (required if generate-pr-comment is set)")
	ReportCmd.Flags().String("current-commit-sha", "", "The current commit SHA (required if generate-pr-comment is set)")
	ReportCmd.Flags().String("repo-url", "", "The repository URL (required if generate-pr-comment is set)")
	ReportCmd.Flags().String("action-run-id", "", "The GitHub Actions run ID (required if generate-pr-comment is set)")

	ReportCmd.MarkFlagRequired("results-path")
}

func generateGitHubSummaryMarkdown(report *reports.TestReport, outputPath string) error {
	fs := reports.OSFileSystem{}
	mdFileName := outputPath + "-summary.md"
	mdFile, err := fs.Create(mdFileName)
	if err != nil {
		return fmt.Errorf("error creating GitHub summary markdown file: %w", err)
	}
	defer mdFile.Close()
	reports.GenerateGitHubSummaryMarkdown(mdFile, report, 1.0)
	return nil
}

func generatePRCommentMarkdown(report *reports.TestReport, outputPath, baseBranch, currentBranch, currentCommitSHA, repoURL, actionRunID string) error {
	fs := reports.OSFileSystem{}
	mdFileName := outputPath + "-pr-comment.md"
	mdFile, err := fs.Create(mdFileName)
	if err != nil {
		return fmt.Errorf("error creating PR comment markdown file: %w", err)
	}
	defer mdFile.Close()
	reports.GeneratePRCommentMarkdown(mdFile, report, 1.0, baseBranch, currentBranch, currentCommitSHA, repoURL, actionRunID)
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