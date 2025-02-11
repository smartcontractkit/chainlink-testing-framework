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

type SummaryData struct {
	TotalTests     int     `json:"total_tests"`
	PanickedTests  int     `json:"panicked_tests"`
	RacedTests     int     `json:"raced_tests"`
	FlakyTests     int     `json:"flaky_tests"`
	FlakyTestRatio string  `json:"flaky_test_ratio"`
	TotalRuns      int     `json:"total_runs"`
	PassedRuns     int     `json:"passed_runs"`
	FailedRuns     int     `json:"failed_runs"`
	SkippedRuns    int     `json:"skipped_runs"`
	PassRatio      string  `json:"pass_ratio"`
	MaxPassRatio   float64 `json:"max_pass_ratio"`
}

var GenerateReportCmd = &cobra.Command{
	Use:   "generate-report",
	Short: "Generate reports from an aggregated test results",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
		}

		// Load the aggregated report
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Loading aggregated test report..."
		s.Start()

		aggregatedReport := &reports.TestReport{}
		reportFile, err := os.Open(aggregatedResultsPath)
		if err != nil {
			s.Stop()
			return fmt.Errorf("error opening aggregated test report: %w", err)
		}
		defer reportFile.Close()

		if err := json.NewDecoder(reportFile).Decode(aggregatedReport); err != nil {
			s.Stop()
			return fmt.Errorf("error decoding aggregated test report: %w", err)
		}
		s.Stop()
		log.Info().Msg("Successfully loaded aggregated test report")

		// Load the summary data to check for failed tests
		var summaryData SummaryData

		if summaryPath == "" {
			return fmt.Errorf("--summary-path is required")
		}

		summaryFile, err := os.Open(summaryPath)
		if err != nil {
			return fmt.Errorf("error opening summary JSON file: %w", err)
		}
		defer summaryFile.Close()

		if err := json.NewDecoder(summaryFile).Decode(&summaryData); err != nil {
			return fmt.Errorf("error decoding summary JSON file: %w", err)
		}

		// Check if there are failed tests
		hasFailedTests := summaryData.FailedRuns > 0

		var artifactLink string
		if hasFailedTests && githubRepo != "" && githubRunID != 0 && artifactName != "" {
			// Fetch artifact link from GitHub API
			artifactLink, err = fetchArtifactLinkWithRetry(githubToken, githubRepo, githubRunID, artifactName, 5, 5*time.Second)
			if err != nil {
				return fmt.Errorf("error fetching artifact link: %w", err)
			}
		}

		// Create output directory if it doesn't exist
		if err := fs.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}

		// Generate GitHub summary markdown
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Generating GitHub summary markdown..."
		s.Start()

		err = generateGitHubSummaryMarkdown(aggregatedReport, filepath.Join(outputDir, "all-test"), artifactLink, artifactName)
		if err != nil {
			s.Stop()
			return fmt.Errorf("error generating GitHub summary markdown: %w", err)
		}
		s.Stop()
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
				return fmt.Errorf("the following flags are required when --generate-pr-comment is set: %s", strings.Join(missingFlags, ", "))
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
				return fmt.Errorf("error generating PR comment markdown: %w", err)
			}
			s.Stop()
			log.Info().Msg("PR comment markdown generated successfully")
		}

		log.Info().Str("output", outputDir).Msg("Reports generated successfully")

		return nil
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
		log.Fatal().Err(err).Msg("Error marking flag as required")
	}
	if err := GenerateReportCmd.MarkFlagRequired("summary-path"); err != nil {
		log.Fatal().Err(err).Msg("Error marking flag as required")
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
	opts := &github.ListOptions{PerPage: 500}
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

func fetchArtifactLinkWithRetry(
	githubToken, githubRepo string,
	githubRunID int64, artifactName string,
	maxRetries int, delay time.Duration,
) (string, error) {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		link, err := fetchArtifactLink(githubToken, githubRepo, githubRunID, artifactName)
		if err == nil {
			// Found the artifact link successfully
			return link, nil
		}

		// If this was our last attempt, return the error
		lastErr = err
		if attempt == maxRetries {
			break
		}

		// Otherwise wait and retry
		log.Printf("[Attempt %d/%d] Artifact not yet available. Retrying in %s...", attempt, maxRetries, delay)
		time.Sleep(delay)
	}

	return "", fmt.Errorf("failed to fetch artifact link after %d retries: %w", maxRetries, lastErr)
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
