package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v70/github"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// GetGHArtifactLinkCmd fetches the artifact link from GitHub API.
var GetGHArtifactLinkCmd = &cobra.Command{
	Use:   "get-gh-artifact",
	Short: "Get artifact link from GitHub API",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag values
		githubRepo, _ := cmd.Flags().GetString("github-repository")
		githubRunID, _ := cmd.Flags().GetInt64("github-run-id")
		artifactName, _ := cmd.Flags().GetString("failed-tests-artifact-name")

		// Get the GitHub token from environment variable
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			log.Error().Msg("GITHUB_TOKEN environment variable is not set")
			os.Exit(ErrorExitCode)
		}

		// Fetch artifact link from GitHub API with retry logic
		artifactLink, err := fetchArtifactLinkWithRetry(githubToken, githubRepo, githubRunID, artifactName, 5, 5*time.Second)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching artifact link")
			os.Exit(ErrorExitCode)
		}

		fmt.Println(artifactLink)
	},
}

func init() {
	GetGHArtifactLinkCmd.Flags().String("github-repository", "", "The GitHub repository in the format owner/repo (required)")
	GetGHArtifactLinkCmd.Flags().Int64("github-run-id", 0, "The GitHub Actions run ID (required)")
	GetGHArtifactLinkCmd.Flags().String("failed-tests-artifact-name", "failed-test-results-with-logs.json", "The name of the failed tests artifact (default 'failed-test-results-with-logs.json')")

	if err := GetGHArtifactLinkCmd.MarkFlagRequired("github-repository"); err != nil {
		log.Error().Err(err).Msg("Error marking github-repository flag as required")
		os.Exit(ErrorExitCode)
	}
	if err := GetGHArtifactLinkCmd.MarkFlagRequired("github-run-id"); err != nil {
		log.Error().Err(err).Msg("Error marking github-run-id flag as required")
		os.Exit(ErrorExitCode)
	}
}

// fetchArtifactLink uses the GitHub API to retrieve the artifact link.
func fetchArtifactLink(githubToken, githubRepo string, githubRunID int64, artifactName string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Split owner and repo from the provided repository string.
	repoParts := strings.SplitN(githubRepo, "/", 2)
	if len(repoParts) != 2 {
		return "", fmt.Errorf("invalid format for --github-repository, expected owner/repo")
	}
	owner, repo := repoParts[0], repoParts[1]

	opts := &github.ListOptions{PerPage: 100} // maximum per page allowed by GitHub
	var allArtifacts []*github.Artifact

	// Paginate through all artifacts.
	for {
		artifacts, resp, err := client.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, githubRunID, opts)
		if err != nil {
			return "", fmt.Errorf("error listing artifacts: %w", err)
		}

		allArtifacts = append(allArtifacts, artifacts.Artifacts...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Search for the artifact by name.
	for _, artifact := range allArtifacts {
		if artifact.GetName() == artifactName {
			artifactID := artifact.GetID()
			artifactURL := fmt.Sprintf("https://github.com/%s/%s/actions/runs/%d/artifacts/%d",
				owner, repo, githubRunID, artifactID)
			return artifactURL, nil
		}
	}

	return "", fmt.Errorf("artifact '%s' not found in the workflow run", artifactName)
}

// fetchArtifactLinkWithRetry attempts to fetch the artifact link with retry logic.
func fetchArtifactLinkWithRetry(githubToken, githubRepo string, githubRunID int64, artifactName string, maxRetries int, delay time.Duration) (string, error) {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		link, err := fetchArtifactLink(githubToken, githubRepo, githubRunID, artifactName)
		if err == nil {
			return link, nil
		}

		lastErr = err
		if attempt == maxRetries {
			break
		}

		log.Printf("[Attempt %d/%d] Artifact not yet available. Retrying in %s...", attempt, maxRetries, delay)
		time.Sleep(delay)
	}

	return "", fmt.Errorf("failed to fetch artifact link after %d retries: %w", maxRetries, lastErr)
}
