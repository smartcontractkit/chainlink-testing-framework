package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var GenerateReportCmd = &cobra.Command{
	Use:   "generate-github-report",
	Short: "Generate Github reports from Flakeguard test report",
	Run: func(cmd *cobra.Command, args []string) {
		fs := reports.OSFileSystem{}

		// Get flag values
		flakeguardReportPath, _ := cmd.Flags().GetString("flakeguard-report")
		summaryReportPath, _ := cmd.Flags().GetString("summary-report-md-path")
		prCommentPath, _ := cmd.Flags().GetString("pr-comment-md-path")
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		failedLogsURL, _ := cmd.Flags().GetString("failed-logs-url")

		failedLogsArtifactName := "failed-test-results-with-logs.json"

		// Load the test report from file
		testReport := &reports.TestReport{}
		reportFile, err := os.Open(flakeguardReportPath)
		if err != nil {
			fmt.Println()
			log.Error().Err(err).Msg("Error opening aggregated test report")
			os.Exit(ErrorExitCode)
		}
		defer reportFile.Close()

		if err := json.NewDecoder(reportFile).Decode(testReport); err != nil {
			fmt.Println()
			log.Error().Err(err).Msg("Error decoding aggregated test report")
			os.Exit(ErrorExitCode)
		}
		fmt.Println()
		log.Info().Msg("Successfully loaded aggregated test report")

		// Create directory for summary markdown if needed
		summaryDir := filepath.Dir(summaryReportPath)
		if err := fs.MkdirAll(summaryDir, 0755); err != nil {
			log.Error().Err(err).Msg("Error creating directory for summary report markdown")
			os.Exit(ErrorExitCode)
		}

		// Generate GitHub summary markdown
		err = generateGitHubSummaryMarkdown(*testReport, summaryReportPath, failedLogsURL, failedLogsArtifactName)
		if err != nil {
			fmt.Println()
			log.Error().Err(err).Msg("Error generating GitHub summary markdown")
			os.Exit(ErrorExitCode)
		}
		fmt.Println()
		log.Info().
			Str("path", summaryReportPath).
			Msg("GitHub summary markdown generated successfully")

		// Generate PR comment markdown if prCommentPath is provided
		if prCommentPath != "" {
			// Retrieve required flags for PR comment
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
				log.Error().Strs("missing flags", missingFlags).Msg("Not all required flags are provided for PR comment generation")
				os.Exit(ErrorExitCode)
			}

			// Create directory for PR comment markdown if needed
			prCommentDir := filepath.Dir(prCommentPath)
			if err := fs.MkdirAll(prCommentDir, 0755); err != nil {
				log.Error().Err(err).Msg("Error creating directory for PR comment markdown")
				os.Exit(ErrorExitCode)
			}

			err = generatePRCommentMarkdown(
				*testReport,
				prCommentPath,
				baseBranch,
				currentBranch,
				currentCommitSHA,
				repoURL,
				actionRunID,
				failedLogsArtifactName,
				failedLogsURL,
				maxPassRatio,
			)
			if err != nil {
				fmt.Println()
				log.Error().Err(err).Msg("Error generating PR comment markdown")
				os.Exit(ErrorExitCode)
			}
			fmt.Println()
			log.Info().
				Str("path", prCommentPath).
				Msg("PR comment markdown generated successfully")
		}
	},
}

func init() {
	GenerateReportCmd.Flags().StringP("flakeguard-report", "i", "", "Path to the flakeguard test report JSON file (required)")
	GenerateReportCmd.Flags().String("summary-report-md-path", "", "Path to output the generated summary markdown file (required)")
	GenerateReportCmd.Flags().String("pr-comment-md-path", "", "Path to output the generated PR comment markdown file (optional)")
	GenerateReportCmd.Flags().Float64P("max-pass-ratio", "", 1.0, "The maximum pass ratio threshold for a test to be considered flaky")
	GenerateReportCmd.Flags().String("base-branch", "develop", "The base branch to compare against (used in PR comment)")
	GenerateReportCmd.Flags().String("current-branch", "", "The current branch name (required if pr-comment-md-path is provided)")
	GenerateReportCmd.Flags().String("current-commit-sha", "", "The current commit SHA (required if pr-comment-md-path is provided)")
	GenerateReportCmd.Flags().String("repo-url", "", "The repository URL (required if pr-comment-md-path is provided)")
	GenerateReportCmd.Flags().String("action-run-id", "", "The GitHub Actions run ID (required if pr-comment-md-path is provided)")
	GenerateReportCmd.Flags().String("github-repository", "", "The GitHub repository in the format owner/repo (required)")
	GenerateReportCmd.Flags().Int64("github-run-id", 0, "The GitHub Actions run ID (required)")
	GenerateReportCmd.Flags().String("failed-logs-url", "", "Optional URL linking to additional logs for failed tests")

	if err := GenerateReportCmd.MarkFlagRequired("flakeguard-report"); err != nil {
		log.Error().Err(err).Msg("Error marking flag as required")
		os.Exit(ErrorExitCode)
	}
	if err := GenerateReportCmd.MarkFlagRequired("summary-report-md-path"); err != nil {
		log.Error().Err(err).Msg("Error marking flag as required")
		os.Exit(ErrorExitCode)
	}
}

func generateGitHubSummaryMarkdown(report reports.TestReport, outputPath, artifactLink, artifactName string) error {
	fs := reports.OSFileSystem{}
	mdFile, err := fs.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating GitHub summary markdown file: %w", err)
	}
	defer mdFile.Close()

	// Generate the summary markdown
	reports.GenerateGitHubSummaryMarkdown(mdFile, report, 1.0, artifactName, artifactLink)

	return nil
}

func generatePRCommentMarkdown(
	report reports.TestReport,
	outputPath,
	baseBranch,
	currentBranch,
	currentCommitSHA,
	repoURL,
	actionRunID,
	artifactName,
	artifactLink string,
	maxPassRatio float64,
) error {
	fs := reports.OSFileSystem{}
	mdFile, err := fs.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating PR comment markdown file: %w", err)
	}
	defer mdFile.Close()

	reports.GeneratePRCommentMarkdown(
		mdFile,
		report,
		maxPassRatio,
		baseBranch,
		currentBranch,
		currentCommitSHA,
		repoURL,
		actionRunID,
		artifactName,
		artifactLink,
	)

	return nil
}
