package client

import (
	"context"
	"net/http"

	// import for side effect of sql packages
	_ "github.com/lib/pq"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client *github.Client
}

const WITHOUT_TOKEN = ""

func NewGithubClient(token string) *GithubClient {
	// Optional: Authenticate with a personal access token if necessary
	// This is recommended to avoid rate limits for unauthenticated requests
	var tc *http.Client
	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(tc)

	return &GithubClient{
		client: client,
	}
}

func (g *GithubClient) ListLatestReleases(org, repository string, count int) ([]*github.RepositoryRelease, error) {
	ctx := context.Background()
	releases, _, err := g.client.Repositories.ListReleases(ctx, org, repository, &github.ListOptions{PerPage: count})
	return releases, err
}

func (g *GithubClient) ListLatestCLCoreReleases(count int) ([]*github.RepositoryRelease, error) {
	ctx := context.Background()
	releases, _, err := g.client.Repositories.ListReleases(ctx, "smartcontractkit", "chainlink", &github.ListOptions{PerPage: count})
	return releases, err
}

func (g *GithubClient) ListLatestCLCoreTags(count int) ([]*github.RepositoryTag, error) {
	ctx := context.Background()
	tags, _, err := g.client.Repositories.ListTags(ctx, "smartcontractkit", "chainlink", &github.ListOptions{PerPage: count})
	return tags, err
}
