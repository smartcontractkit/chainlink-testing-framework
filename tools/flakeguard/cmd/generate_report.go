package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var GenerateReportCmd = &cobra.Command{
	Use:   "generate-report",
	Short: "Generate test reports from aggregated results that can be posted to GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		fs := reports.OSFileSystem{}

		// Get flag values
		aggregatedResultsPath, _ := cmd.Flags().GetString("aggregated-results-path")
		summaryPath, _ := cmd.Flags().GetString("summary-path")
		outputDir, _ := cmd.Flags().GetString("output-path")
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		generatePRComment, _ := cmd.Flags().GetBool("generate-pr-comment")
		githubRepo, _ := cmd.Flags().GetString("github-repository")
		githubRunID, _ := cmd.Flags().GetInt64("github-run-id")
		artifactName, _ := cmd.Flags().GetString("failed-tests-artifact-name")

		// Get the GitHub token from environment variable
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			log.Error().Msg("GITHUB_TOKEN environment variable is not set")
			os.Exit(ErrorExitCode)
		}

		// Load the aggregated report
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Loading aggregated test report..."
		s.Start()

		aggregatedReport := &reports.TestReport{}
		reportFile, err := os.Open(aggregatedResultsPath)
		if err != nil {
			s.Stop()
			fmt.Println()
			log.Error().Err(err).Msg("Error opening aggregated test report")
			os.Exit(ErrorExitCode)
		}
		defer reportFile.Close()

		if err := json.NewDecoder(reportFile).Decode(aggregatedReport); err != nil {
			s.Stop()
			fmt.Println()
			log.Error().Err(err).Msg("Error decoding aggregated test report")
			os.Exit(ErrorExitCode)
		}
		s.Stop()
		fmt.Println()
		log.Info().Msg("Successfully loaded aggregated test report")

		// Load the summary data to check for failed tests
		var summaryData reports.SummaryData

		if summaryPath == "" {
			log.Error().Msg("Summary path is required")
			os.Exit(ErrorExitCode)
		}

		summaryFile, err := os.Open(summaryPath)
		if err != nil {
			log.Error().Err(err).Msg("Error opening summary JSON file")
			os.Exit(ErrorExitCode)
		}
		defer summaryFile.Close()

		if err := json.NewDecoder(summaryFile).Decode(&summaryData); err != nil {
			log.Error().Err(err).Msg("Error decoding summary JSON file")
			os.Exit(ErrorExitCode)
		}

		// Check if there are failed tests
		hasFailedTests := summaryData.FailedRuns > 0

		var artifactLink string
		if hasFailedTests {
			// Fetch artifact link from GitHub API
			artifactLink, err = fetchArtifactLink(githubToken, githubRepo, githubRunID, artifactName)
			if err != nil {
				log.Error().Err(err).Msg("Error fetching artifact link")
				os.Exit(ErrorExitCode)
			}
		} else {
			// No failed tests, set artifactLink to empty string
			artifactLink = ""
			log.Debug().Msg("No failed tests found. Skipping artifact link generation")
		}

		// Create output directory if it doesn't exist
		if err := fs.MkdirAll(outputDir, 0755); err != nil {
			log.Error().Err(err).Msg("Error creating output directory")
			os.Exit(ErrorExitCode)
		}

		// Generate GitHub summary markdown
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Generating GitHub summary markdown..."
		s.Start()

		err = generateGitHubSummaryMarkdown(aggregatedReport, filepath.Join(outputDir, "all-test"), artifactLink, artifactName)
		if err != nil {
			s.Stop()
			fmt.Println()
			log.Error().Err(err).Msg("Error generating GitHub summary markdown")
			os.Exit(ErrorExitCode)
		}
		s.Stop()
		fmt.Println()
		log.Info().Msg("GitHub summary markdown generated successfully")

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
				log.Error().Strs("missing flags", missingFlags).Msg("Not all required flags are provided for --generate-pr-comment")
				os.Exit(ErrorExitCode)
			}

			// Generate PR comment markdown
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " Generating PR comment markdown..."
			s.Start()

			err = generatePRCommentMarkdown(
				aggregatedReport,
				filepath.Join(outputDir, "all-test"),
				baseBranch,
				currentBranch,
				currentCommitSHA,
				repoURL,
				actionRunID,
				artifactName,
				artifactLink,
				maxPassRatio,
			)
			if err != nil {
				s.Stop()
				fmt.Println()
				log.Error().Err(err).Msg("Error generating PR comment markdown")
				os.Exit(ErrorExitCode)
			}
			s.Stop()
			fmt.Println()
			log.Info().Msg("PR comment markdown generated successfully")
		}

		log.Info().Str("output", outputDir).Msg("Reports generated successfully")
	},
}

func init() {
	GenerateReportCmd.Flags().StringP("aggregated-results-path", "i", "", "Path to the aggregated JSON report file (required)")
	GenerateReportCmd.Flags().StringP("summary-path", "s", "", "Path to the summary JSON file (required)")
	GenerateReportCmd.Flags().StringP("output-path", "o", "./report", "Path to output the generated report files")
	GenerateReportCmd.Flags().Float64P("max-pass-ratio", "", 1.0, "The maximum pass ratio threshold for a test to be considered flaky")
	GenerateReportCmd.Flags().Bool("generate-pr-comment", false, "Set to true to generate PR comment markdown")
	GenerateReportCmd.Flags().String("base-branch", "develop", "The base branch to compare against (used in PR comment)")
	GenerateReportCmd.Flags().String("current-branch", "", "The current branch name (required if generate-pr-comment is set)")
	GenerateReportCmd.Flags().String("current-commit-sha", "", "The current commit SHA (required if generate-pr-comment is set)")
	GenerateReportCmd.Flags().String("repo-url", "", "The repository URL (required if generate-pr-comment is set)")
	GenerateReportCmd.Flags().String("action-run-id", "", "The GitHub Actions run ID (required if generate-pr-comment is set)")
	GenerateReportCmd.Flags().String("github-repository", "", "The GitHub repository in the format owner/repo (required)")
	GenerateReportCmd.Flags().Int64("github-run-id", 0, "The GitHub Actions run ID (required)")
	GenerateReportCmd.Flags().String("failed-tests-artifact-name", "failed-test-results-with-logs.json", "The name of the failed tests artifact (default 'failed-test-results-with-logs.json')")

	if err := GenerateReportCmd.MarkFlagRequired("aggregated-results-path"); err != nil {
		log.Error().Err(err).Msg("Error marking flag as required")
		os.Exit(ErrorExitCode)
	}
	if err := GenerateReportCmd.MarkFlagRequired("summary-path"); err != nil {
		log.Error().Err(err).Msg("Error marking flag as required")
		os.Exit(ErrorExitCode)
	}
	if err := GenerateReportCmd.MarkFlagRequired("github-repository"); err != nil {
		log.Error().Err(err).Msg("Error marking flag as required")
		os.Exit(ErrorExitCode)
	}
	if err := GenerateReportCmd.MarkFlagRequired("github-run-id"); err != nil {
		log.Error().Err(err).Msg("Error marking flag as required")
		os.Exit(ErrorExitCode)
	}
}

func fetchArtifactLink(githubToken, githubRepo string, githubRunID int64, artifactName string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Split the repository into owner and repo
	repoParts := strings.SplitN(githubRepo, "/", 2)
	if len(repoParts) != 2 {
		return "", fmt.Errorf("invalid format for --github-repository, expected owner/repo")
	}
	owner, repo := repoParts[0], repoParts[1]

	// List artifacts for the workflow run
	opts := &github.ListOptions{PerPage: 100}
	artifacts, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, githubRunID, opts)
	if err != nil {
		return "", fmt.Errorf("error listing artifacts: %w", err)
	}

	// Find the artifact
	for _, artifact := range artifacts.Artifacts {
		if artifact.GetName() == artifactName {
			// Construct the artifact URL using the artifact ID
			artifactID := artifact.GetID()
			artifactURL := fmt.Sprintf("https://github.com/%s/%s/actions/runs/%d/artifacts/%d", owner, repo, githubRunID, artifactID)
			return artifactURL, nil
		}
	}

	return "", fmt.Errorf("artifact '%s' not found in the workflow run", artifactName)
}

func generateGitHubSummaryMarkdown(report *reports.TestReport, outputPath, artifactLink, artifactName string) error {
	fs := reports.OSFileSystem{}
	mdFileName := outputPath + "-summary.md"
	mdFile, err := fs.Create(mdFileName)
	if err != nil {
		return fmt.Errorf("error creating GitHub summary markdown file: %w", err)
	}
	defer mdFile.Close()

	// Generate the summary markdown
	reports.GenerateGitHubSummaryMarkdown(mdFile, report, 1.0, artifactName, artifactLink)

	return nil
}

func generatePRCommentMarkdown(
	report *reports.TestReport,
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
	mdFileName := outputPath + "-pr-comment.md"
	mdFile, err := fs.Create(mdFileName)
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
