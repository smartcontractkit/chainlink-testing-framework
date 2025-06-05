// TODO: This is a PoC, need to be cleaned up much more before game time
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var MakeIssueCmd = &cobra.Command{
	Use:   "make-issue",
	Short: "Make an issue to skip identified flaky tests",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if githubToken == "" {
			githubToken = os.Getenv("GITHUB_TOKEN")
		}
		if githubToken == "" {
			return fmt.Errorf("GitHub token not set, set GITHUB_TOKEN env var or use --githubToken flag")
		}
		return nil
	},
	RunE: makeIssue,
}

func makeIssue(cmd *cobra.Command, args []string) error {
	// Set up GitHub client
	tok := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	token := oauth2.NewClient(context.Background(), tok)
	graphqlClient := githubv4.NewClient(token)

	// https://docs.github.com/en/copilot/using-github-copilot/coding-agent/using-copilot-to-work-on-an-issue#assigning-an-issue-to-copilot-via-the-github-api

	// Get suggestedActors from the API to check if Copilot is enabled
	var suggestedActorsQuery struct {
		Repository struct {
			SuggestedActors struct {
				Nodes []struct {
					Login    string
					Typename string `graphql:"__typename"`
					Bot      struct {
						ID githubv4.ID
					} `graphql:"... on Bot"`
					User struct {
						ID githubv4.ID
					} `graphql:"... on User"`
				}
			} `graphql:"suggestedActors(capabilities: [CAN_BE_ASSIGNED], first: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String("smartcontractkit"), // Replace with your repo owner
		"name":  githubv4.String("chainlink"),        // Replace with your repo name
	}

	err := graphqlClient.Query(context.Background(), &suggestedActorsQuery, variables)
	if err != nil {
		return fmt.Errorf("failed to get suggested actors: %w", err)
	}

	// Check if Copilot coding agent is enabled
	// If enabled, the first node will have login "copilot-swe-agent"
	copilotFound := false
	for _, actor := range suggestedActorsQuery.Repository.SuggestedActors.Nodes {
		if actor.Login == "copilot-swe-agent" {
			copilotFound = true
			break // Found it, no need to continue
		}
	}

	if copilotFound {
		fmt.Println("✅ Copilot coding agent is available for this user and repository")
	} else {
		fmt.Println("❌ Copilot coding agent is not enabled for this user/repository")
		os.Exit(ErrorExitCode)
	}

	return nil
}

func init() {
	MakeIssueCmd.Flags().StringVarP(&githubToken, "githubToken", "t", "", "GitHub token to use for creating the issue (can be set with GITHUB_TOKEN env var)")
}
