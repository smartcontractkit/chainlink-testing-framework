package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	flake_git "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/git"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/golang"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
)

const (
	openAIKeyEnvVar = "OPENAI_API_KEY"
)

var (
	repoPath    string
	localDBPath string
	openAIKey   string
)

var MakePRCmd = &cobra.Command{
	Use:   "make-pr",
	Short: "Make a PR to skip identified flaky tests",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("openAIKey").Changed {
			openAIKey = cmd.Flag("openAIKey").Value.String()
		} else if os.Getenv(openAIKeyEnvVar) != "" {
			openAIKey = os.Getenv(openAIKeyEnvVar)
		} else {
			fmt.Printf("%s is not set, cannot use LLM to skip or fix flaky tests\n", openAIKeyEnvVar)
		}
	},
	RunE: makePR,
}

func makePR(cmd *cobra.Command, args []string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}

	db, err := localdb.LoadDBWithPath(localDBPath)
	if err != nil {
		return fmt.Errorf("failed to load local db: %w", err)
	}

	currentlyFlakyEntries := db.GetAllCurrentlyFlakyEntries()

	owner, repoName, defaultBranch, err := flake_git.GetOwnerRepoDefaultBranchFromLocalRepo(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get repo info: %w", err)
	}

	targetRepoWorktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to open repo's worktree: %w", err)
	}

	// First checkout default branch and pull latest
	err = targetRepoWorktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(defaultBranch),
	})
	if err != nil {
		if errors.Is(err, git.ErrUnstagedChanges) {
			fmt.Println("Local repo has unstaged changes, please commit or stash them before running this command")
		}
		return fmt.Errorf("failed to checkout default branch %s: %w", defaultBranch, err)
	}

	fmt.Printf("Fetching latest changes from default branch '%s', tap your yubikey if it's blinking...", defaultBranch)
	err = repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to fetch latest: %w", err)
	}
	fmt.Println(" ✅")

	fmt.Printf("Pulling latest changes from default branch '%s', tap your yubikey if it's blinking...", defaultBranch)
	err = targetRepoWorktree.Pull(&git.PullOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}
	fmt.Println(" ✅")

	// Create and checkout new branch
	branchName := fmt.Sprintf("flakeguard-skip-%s", time.Now().Format("20060102150405"))
	err = targetRepoWorktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout new branch: %w", err)
	}

	cleanUpBranch := true
	defer func() {
		if cleanUpBranch {
			fmt.Printf("Cleaning up branch %s...", branchName)
			// First checkout default branch
			err = targetRepoWorktree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(defaultBranch),
				Force:  true, // Force checkout to discard any changes for a clean default branch
			})
			if err != nil {
				fmt.Printf("Failed to checkout default branch: %v\n", err)
				return
			}
			// Then delete the local branch
			err = repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(branchName))
			if err != nil {
				fmt.Printf("Failed to remove local branch: %v\n", err)
				return
			}
			fmt.Println(" ✅")
		}
	}()

	if len(currentlyFlakyEntries) == 0 {
		fmt.Println("No flaky tests found!")
		return nil
	}

	jiraTickets := []string{}
	testsToSkip := []*golang.SkipTest{}
	for _, entry := range currentlyFlakyEntries {
		testsToSkip = append(testsToSkip, &golang.SkipTest{
			Package:    entry.TestPackage,
			Name:       entry.TestName,
			JiraTicket: entry.JiraTicket,
		})
		jiraTickets = append(jiraTickets, entry.JiraTicket)
	}

	err = golang.SkipTests(repoPath, openAIKey, testsToSkip)
	if err != nil {
		return fmt.Errorf("failed to modify code to skip tests: %w", err)
	}

	_, err = targetRepoWorktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	fmt.Print("Committing changes, tap your yubikey if it's blinking...")
	commitHash, err := targetRepoWorktree.Commit(fmt.Sprintf("Skips flaky %d tests", len(testsToSkip)), &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	fmt.Println(" ✅")

	fmt.Print("Pushing changes to remote, tap your yubikey if it's blinking...")
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}
	fmt.Println(" ✅")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var (
		skippedTestsPRBody        strings.Builder
		alreadySkippedTestsPRBody strings.Builder
		errorSkippingTestsPRBody  strings.Builder
		llmSkippedTestsPRBody     strings.Builder
	)

	for _, test := range testsToSkip {
		if test.ErrorSkipping != nil {
			errorSkippingTestsPRBody.WriteString(fmt.Sprintf("- Package: `%s`\n", test.Package))
			errorSkippingTestsPRBody.WriteString(fmt.Sprintf("  Test: `%s`\n", test.Name))
			errorSkippingTestsPRBody.WriteString(fmt.Sprintf("  Ticket: [%s](https://%s/browse/%s)\n", test.JiraTicket, os.Getenv("JIRA_DOMAIN"), test.JiraTicket))
			errorSkippingTestsPRBody.WriteString(fmt.Sprintf("  Error: %s\n\n", test.ErrorSkipping))
		} else if test.SimplySkipped {
			skippedTestsPRBody.WriteString(fmt.Sprintf("- Package: `%s`\n", test.Package))
			skippedTestsPRBody.WriteString(fmt.Sprintf("  Test: `%s`\n", test.Name))
			skippedTestsPRBody.WriteString(fmt.Sprintf("  Ticket: [%s](https://%s/browse/%s)\n", test.JiraTicket, os.Getenv("JIRA_DOMAIN"), test.JiraTicket))
			skippedTestsPRBody.WriteString(fmt.Sprintf("  [View skip in PR](https://github.com/%s/%s/pull/%s/files#diff-%sL%d)\n\n", owner, repoName, branchName, commitHash, test.Line))
		} else if test.AlreadySkipped {
			alreadySkippedTestsPRBody.WriteString(fmt.Sprintf("- Package: `%s`\n", test.Package))
			alreadySkippedTestsPRBody.WriteString(fmt.Sprintf("  Test: `%s`\n", test.Name))
			alreadySkippedTestsPRBody.WriteString(fmt.Sprintf("  Ticket: [%s](https://%s/browse/%s)\n", test.JiraTicket, os.Getenv("JIRA_DOMAIN"), test.JiraTicket))
		} else if test.LLMSkipped {
			llmSkippedTestsPRBody.WriteString(fmt.Sprintf("- Package: `%s`\n", test.Package))
			llmSkippedTestsPRBody.WriteString(fmt.Sprintf("  Test: `%s`\n", test.Name))
			llmSkippedTestsPRBody.WriteString(fmt.Sprintf("  Ticket: [%s](https://%s/browse/%s)\n", test.JiraTicket, os.Getenv("JIRA_DOMAIN"), test.JiraTicket))
			llmSkippedTestsPRBody.WriteString(fmt.Sprintf("  [View skip in PR](https://github.com/%s/%s/pull/%s/files#diff-%sL%d)\n\n", owner, repoName, branchName, commitHash, test.Line))
		}
	}
	body := fmt.Sprintf(`## Tests That I Failed to Skip, Need Manual Intervention

%s

## Tests Skipped Using Simple AST Parsing

%s

## Tests Skipped Using LLM Assistance

%s

## Tests That Were Already Skipped

%s`, errorSkippingTestsPRBody.String(), skippedTestsPRBody.String(), llmSkippedTestsPRBody.String(), alreadySkippedTestsPRBody.String())

	pr := &github.NewPullRequest{
		Title:               github.Ptr(fmt.Sprintf("[%s] Flakeguard: Skip flaky tests", strings.Join(jiraTickets, "] ["))),
		Head:                github.Ptr(branchName),
		Base:                github.Ptr(defaultBranch),
		Body:                github.Ptr(body),
		MaintainerCanModify: github.Ptr(true),
	}

	fmt.Println("PR Preview:")
	fmt.Println("================================================")
	fmt.Println(pr.Title)
	fmt.Println("--------------------------------")
	fmt.Printf("Merging '%s' into '%s'\n", branchName, defaultBranch)
	fmt.Println(pr.Body)
	fmt.Println("================================================")

	fmt.Printf("To preview the code changes in the GitHub UI, visit: https://github.com/%s/%s/compare/%s...%s\n", owner, repoName, defaultBranch, branchName)
	fmt.Print("Would you like to create the PR automatically from the CLI? (y/N): ")

	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		return err
	}

	if strings.ToLower(confirm) != "y" {
		fmt.Println("Exiting. Please use the GitHub UI to create the PR.")
		return nil
	}

	createdPR, resp, err := client.PullRequests.Create(ctx, owner, repoName, pr)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read github response body while trying to create PR: %s\n%w", resp.Status, err)
		}
		return fmt.Errorf("failed to create PR, got bad status: %s\n%s", resp.Status, string(body))
	}

	cleanUpBranch = false
	fmt.Printf("PR created! https://github.com/%s/%s/pull/%d\n", owner, repoName, createdPR.GetNumber())
	return nil
}

func init() {
	MakePRCmd.Flags().StringVarP(&repoPath, "repoPath", "r", ".", "Local path to the repository to make the PR for")
	MakePRCmd.Flags().StringVarP(&openAIKey, "openAIKey", "k", "", fmt.Sprintf("OpenAI API key for using an LLM to help create the PR (can be set via %s environment variable)", openAIKeyEnvVar))
}
