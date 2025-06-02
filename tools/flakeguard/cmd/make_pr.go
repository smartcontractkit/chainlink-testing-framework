package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v72/github"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
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
		return err
	}

	db, err := localdb.LoadDBWithPath(localDBPath)
	if err != nil {
		return err
	}

	entries := db.GetAllEntries()

	branchName := fmt.Sprintf("flakeguard-skip-%s", time.Now().Format("20060102150405"))
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: true,
	})
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Find the test function in the codebase and skip it

		entry.SkippedAt = time.Now()
		db.UpsertEntry(entry.TestPackage, entry.TestName, entry.JiraTicket, entry.SkippedAt, entry.AssigneeID)

	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}
	_, err = w.Commit("Skips flaky tests", &git.CommitOptions{})
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner := "your-org"
	repoName := "your-repo"
	pr := &github.NewPullRequest{
		Title:               github.Ptr("Skip flaky tests"),
		Head:                github.Ptr(branchName),
		Base:                github.Ptr("main"),
		Body:                github.Ptr("This PR skips flaky tests."),
		MaintainerCanModify: github.Ptr(true),
	}
	_, _, err = client.PullRequests.Create(ctx, owner, repoName, pr)
	if err != nil {
		return err
	}

	fmt.Println("PR created!")
	return nil
}

func init() {
	MakePRCmd.Flags().StringVarP(&repoPath, "repo", "r", ".", "Path to the repository to make the PR in")
}
