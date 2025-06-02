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

var (
	repoPath    string
	localDBPath string
)

var MakePRCmd = &cobra.Command{
	Use:   "make-pr",
	Short: "Make a PR to skip identified flaky tests",
	RunE:  makePR,
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

	fmt.Print("Fetching latest changes from default branch, tap your yubikey if it's blinking...")
	err = repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to fetch latest: %w", err)
	}
	fmt.Println(" ✅")

	fmt.Print("Pulling latest changes from default branch, tap your yubikey if it's blinking...")
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
			err = targetRepoWorktree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(defaultBranch),
			})
			if err != nil {
				fmt.Printf("Failed to clean up branch: %v\n", err)
			}
			err = repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(branchName))
			if err != nil {
				fmt.Printf("Failed to remove branch: %v\n", err)
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

	err = golang.SkipTests(repoPath, testsToSkip)
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
	)

	for _, test := range testsToSkip {
		if test.Skipped {
			skippedTestsPRBody.WriteString(fmt.Sprintf("- Package: `%s`\n", test.Package))
			skippedTestsPRBody.WriteString(fmt.Sprintf("  Test: `%s`\n", test.Name))
			skippedTestsPRBody.WriteString(fmt.Sprintf("  Ticket: [%s](https://%s/browse/%s)\n", test.JiraTicket, os.Getenv("JIRA_DOMAIN"), test.JiraTicket))
			skippedTestsPRBody.WriteString(fmt.Sprintf("  [View skip in PR](https://github.com/%s/%s/pull/%s/files#diff-%sL%d)\n\n", owner, repoName, branchName, commitHash, test.Line))
		} else {
			alreadySkippedTestsPRBody.WriteString(fmt.Sprintf("- Package: `%s`\n", test.Package))
			alreadySkippedTestsPRBody.WriteString(fmt.Sprintf("  Test: `%s`\n", test.Name))
			alreadySkippedTestsPRBody.WriteString(fmt.Sprintf("  Ticket: [%s](https://%s/browse/%s)\n", test.JiraTicket, os.Getenv("JIRA_DOMAIN"), test.JiraTicket))
		}
	}

	pr := &github.NewPullRequest{
		Title:               github.Ptr(fmt.Sprintf("[%s] Flakeguard: Skip flaky tests", strings.Join(jiraTickets, "] ["))),
		Head:                github.Ptr(branchName),
		Base:                github.Ptr(defaultBranch),
		Body:                github.Ptr(fmt.Sprintf("## Tests Skipped\n\n%s\n\n## Tests Already Skipped\n\n%s", skippedTestsPRBody.String(), alreadySkippedTestsPRBody.String())),
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
}
