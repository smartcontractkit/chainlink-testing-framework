package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	// import for side effect of sql packages
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client *github.Client
}

const WITHOUT_TOKEN = ""

// NewGithubClient creates a new instance of GithubClient, optionally authenticating with a personal access token.
// This is useful for making authenticated requests to the GitHub API, helping to avoid rate limits.
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

// ListLatestReleases lists the latest releases for a given repository
func (g *GithubClient) ListLatestReleases(org, repository string, count int) ([]*github.RepositoryRelease, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancelFn()
	releases, _, err := g.client.Repositories.ListReleases(ctx, org, repository, &github.ListOptions{PerPage: count})
	return releases, err
}

// ListLatestCLCoreReleases lists the latest releases for the Chainlink core repository
func (g *GithubClient) ListLatestCLCoreReleases(count int) ([]*github.RepositoryRelease, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancelFn()
	releases, _, err := g.client.Repositories.ListReleases(ctx, "smartcontractkit", "chainlink", &github.ListOptions{PerPage: count})
	return releases, err
}

// ListLatestCLCoreTags lists the latest tags for the Chainlink core repository
func (g *GithubClient) ListLatestCLCoreTags(count int) ([]*github.RepositoryTag, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancelFn()
	tags, _, err := g.client.Repositories.ListTags(ctx, "smartcontractkit", "chainlink", &github.ListOptions{PerPage: count})
	return tags, err
}

func (g *GithubClient) DownloadAssetFromRelease(owner, repository, releaseTag, assetName string) ([]byte, error) {
	var content []byte

	// assuming 180s is enough to fetch releases, find the asset we need and download it
	// some assets might be 30+ MB, so we need to give it some time (for really slow connections)
	ctx, cancelFn := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancelFn()
	ghReleases, _, err := g.client.Repositories.ListReleases(ctx, owner, repository, &github.ListOptions{PerPage: 20})
	if err != nil {
		return content, errors.Wrapf(err, "failed to list releases for %s", repository)
	}

	var ghRelease *github.RepositoryRelease
	for _, release := range ghReleases {
		if release.TagName == nil {
			continue
		}

		if *release.TagName == releaseTag {
			ghRelease = release
			break
		}
	}

	if ghRelease == nil {
		return content, errors.New("failed to find release with tag: " + releaseTag)
	}

	var assetID int64
	for _, asset := range ghRelease.Assets {
		if strings.Contains(asset.GetName(), assetName) {
			assetID = asset.GetID()
			break
		}
	}

	if assetID == 0 {
		return content, fmt.Errorf("failed to find asset %s for %s", assetName, *ghRelease.TagName)
	}

	asset, _, err := g.client.Repositories.DownloadReleaseAsset(ctx, owner, repository, assetID, g.client.Client())
	if err != nil {
		return content, errors.Wrapf(err, "failed to download asset %s for %s", assetName, *ghRelease.TagName)
	}

	content, err = io.ReadAll(asset)
	if err != nil {
		return content, err
	}

	return content, nil
}
