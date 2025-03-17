package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
	"github.com/spf13/cobra"
)

var GenerateTestReportCmd = &cobra.Command{
	Use:   "generate-test-report",
	Short: "Generate test report based on test results",
	Run: func(cmd *cobra.Command, args []string) {
		fs := reports.OSFileSystem{}

		// Get flag values
		testResultsDir, _ := cmd.Flags().GetString("test-results-dir")
		outputDir, _ := cmd.Flags().GetString("output-path")
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		projectPath, _ := cmd.Flags().GetString("project-path")
		repoPath, _ := cmd.Flags().GetString("repo-path")
		codeOwnersPath, _ := cmd.Flags().GetString("codeowners-path")
		useRace, _ := cmd.Flags().GetBool("race")
		repoURL, _ := cmd.Flags().GetString("repo-url")
		branchName, _ := cmd.Flags().GetString("branch-name")
		headSHA, _ := cmd.Flags().GetString("head-sha")
		baseSHA, _ := cmd.Flags().GetString("base-sha")
		githubWorkflowName, _ := cmd.Flags().GetString("github-workflow-name")
		githubWorkflowRunURL, _ := cmd.Flags().GetString("github-workflow-run-url")
		reportID, _ := cmd.Flags().GetString("report-id")
		rerunOfReportID, _ := cmd.Flags().GetString("rerun-of-report-id")
		genReportID, _ := cmd.Flags().GetBool("gen-report-id")

		goProject, err := utils.GetGoProjectName(projectPath)
		if err != nil {
			log.Warn().Err(err).Str("projectPath", goProject).Msg("Failed to get pretty project path")
		}

		initialDirSize, err := getDirSize(testResultsDir)
		if err != nil {
			log.Error().Err(err).Str("path", testResultsDir).Msg("Error getting initial directory size")
			// intentionally don't exit here, as we can still proceed with the aggregation
		}

		// Ensure the output directory exists
		if err := fs.MkdirAll(outputDir, 0755); err != nil {
			log.Error().Err(err).Str("path", outputDir).Msg("Error creating output directory")
			os.Exit(ErrorExitCode)
		}

		// Start spinner for loading test results
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Aggregating test results..."
		s.Start()
		fmt.Println()

		// Load test results from JSON files and aggregate them
		aggregatedResults, err := reports.LoadAndAggregate(testResultsDir)
		if err != nil {
			s.Stop()
			log.Error().Err(err).Stack().Msg("Error aggregating test results")
			os.Exit(ErrorExitCode)
		}
		s.Stop()
		log.Debug().Msg("Successfully loaded and aggregated test results")

		// Start spinner for mapping test results to paths
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Filter failed tests..."
		s.Start()

		failedTests := reports.FilterTests(aggregatedResults, func(tr reports.TestResult) bool {
			return !tr.Skipped && tr.PassRatio < maxPassRatio
		})
		s.Stop()

		// Check if there are any failed tests
		if len(failedTests) > 0 {
			log.Info().Int("count", len(failedTests)).Msg("Found failed tests")

			// Create a new report for failed tests with logs
			failedReportWithLogs, err := reports.NewTestReport(failedTests,
				reports.WithGoProject(goProject),
				reports.WithProjectPath(projectPath),
				reports.WithRepoPath(repoPath),
				reports.WithCodeOwnersPath(codeOwnersPath),
				reports.WithRerunOfReportID(rerunOfReportID),
				reports.WithReportID(reportID),
				reports.WithGeneratedReportID(genReportID),
				reports.WithGoRaceDetection(useRace),
				reports.WithBranchName(branchName),
				reports.WithBaseSha(baseSHA),
				reports.WithHeadSha(headSHA),
				reports.WithRepoURL(repoURL),
				reports.WithGitHubWorkflowName(githubWorkflowName),
				reports.WithGitHubWorkflowRunURL(githubWorkflowRunURL),
			)
			if err != nil {
				log.Error().Stack().Err(err).Msg("Error creating failed tests report with logs")
				os.Exit(ErrorExitCode)
			}

			// Save the failed tests report with logs
			failedTestsReportWithLogsPath := filepath.Join(outputDir, "failed-test-report-with-logs.json")
			if err := reports.SaveReport(fs, failedTestsReportWithLogsPath, failedReportWithLogs); err != nil {
				log.Error().Stack().Err(err).Msg("Error saving failed tests report with logs")
				os.Exit(ErrorExitCode)
			}
			log.Info().Str("path", failedTestsReportWithLogsPath).Msg("Failed tests report with logs saved")

			// Remove logs from test results for the report without logs
			for i := range failedReportWithLogs.Results {
				failedReportWithLogs.Results[i].PassedOutputs = nil
				failedReportWithLogs.Results[i].FailedOutputs = nil
				failedReportWithLogs.Results[i].PackageOutputs = nil
			}

			// Save the failed tests report without logs
			failedTestsReportNoLogsPath := filepath.Join(outputDir, "failed-test-report.json")
			if err := reports.SaveReport(fs, failedTestsReportNoLogsPath, failedReportWithLogs); err != nil {
				log.Error().Stack().Err(err).Msg("Error saving failed tests report without logs")
				os.Exit(ErrorExitCode)
			}
			log.Info().Str("path", failedTestsReportNoLogsPath).Msg("Failed tests report without logs saved")
		} else {
			log.Info().Msg("No failed tests found. Skipping generation of failed tests reports")
		}

		// Remove logs from test results for the aggregated report
		for i := range aggregatedResults {
			aggregatedResults[i].PassedOutputs = nil
			aggregatedResults[i].FailedOutputs = nil
			aggregatedResults[i].PackageOutputs = nil
		}

		aggregatedReport, err := reports.NewTestReport(aggregatedResults,
			reports.WithGoProject(goProject),
			reports.WithProjectPath(projectPath),
			reports.WithRepoPath(repoPath),
			reports.WithCodeOwnersPath(codeOwnersPath),
			reports.WithRerunOfReportID(rerunOfReportID),
			reports.WithReportID(reportID),
			reports.WithGeneratedReportID(genReportID),
			reports.WithGoRaceDetection(useRace),
			reports.WithBranchName(branchName),
			reports.WithBaseSha(baseSHA),
			reports.WithHeadSha(headSHA),
			reports.WithRepoURL(repoURL),
			reports.WithGitHubWorkflowName(githubWorkflowName),
			reports.WithGitHubWorkflowRunURL(githubWorkflowRunURL),
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Error creating aggregated test report")
			os.Exit(ErrorExitCode)
		}

		// Save the aggregated report to the output directory
		aggregatedReportPath := filepath.Join(outputDir, "all-test-report.json")
		if err := reports.SaveReport(fs, aggregatedReportPath, aggregatedReport); err != nil {
			log.Error().Stack().Err(err).Msg("Error saving aggregated test report")
			os.Exit(ErrorExitCode)
		}
		log.Info().Str("path", aggregatedReportPath).Msg("All tests report without logs saved")

		finalDirSize, err := getDirSize(testResultsDir)
		if err != nil {
			log.Error().Err(err).Str("path", testResultsDir).Msg("Error getting final directory size")
			// intentionally don't exit here, as we can still proceed with the aggregation
		}
		diskSpaceUsed := byteCountSI(finalDirSize - initialDirSize)
		log.Info().Str("disk space used", diskSpaceUsed).Msg("Aggregation complete")
	},
}

func init() {
	GenerateTestReportCmd.Flags().StringP("test-results-dir", "p", "", "Path to the folder containing JSON test result files (required)")
	GenerateTestReportCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	GenerateTestReportCmd.Flags().StringP("output-path", "o", "./report", "Path to output the aggregated results (directory)")
	GenerateTestReportCmd.Flags().Float64P("max-pass-ratio", "", 1.0, "The maximum pass ratio threshold for a test to be considered flaky")
	GenerateTestReportCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")
	GenerateTestReportCmd.Flags().StringP("repo-path", "", ".", "The path to the root of the repository/project")
	GenerateTestReportCmd.Flags().String("repo-url", "", "The repository URL")
	GenerateTestReportCmd.Flags().String("branch-name", "", "Branch name for the test report")
	GenerateTestReportCmd.Flags().String("head-sha", "", "Head commit SHA for the test report")
	GenerateTestReportCmd.Flags().String("base-sha", "", "Base commit SHA for the test report")
	GenerateTestReportCmd.Flags().String("github-workflow-name", "", "GitHub workflow name for the test report")
	GenerateTestReportCmd.Flags().String("github-workflow-run-url", "", "GitHub workflow run URL for the test report")
	GenerateTestReportCmd.Flags().String("report-id", "", "Optional identifier for the test report. Will be generated if not provided")
	GenerateTestReportCmd.Flags().Bool("gen-report-id", false, "Generate a random report ID")
	GenerateTestReportCmd.Flags().String("rerun-of-report-id", "", "Optional identifier for the report this is a rerun of")
	GenerateTestReportCmd.Flags().Bool("race", false, "Enable the race detector")

	if err := GenerateTestReportCmd.MarkFlagRequired("test-results-dir"); err != nil {
		log.Fatal().Err(err).Msg("Error marking flag as required")
	}
}

// getDirSize returns the size of a directory in bytes
// helpful for tracking how much data is being produced on disk
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// byteCountSI returns a human-readable byte count (decimal SI units)
func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
