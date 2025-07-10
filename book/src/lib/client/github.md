# Github

This small client makes it easy to get `N` latest releases or tags from any Github.com repository. To use it, all you need to have
is a properly scoped access token.

```go
publicRepoClient := NewGithubClient(WITHOUT_TOKEN)

// "smartcontractkit", "chainlink"
latestCLReleases, err := publicRepoClient.ListLatestCLCoreReleases(10)
if err != nil {
    panic(err)
}

// "smartcontractkit", "chainlink"
latestCLTags, err := publicRepoClient.ListLatestCLCoreTags(10)
if err != nil {
    panic(err)
}

privateRepoClient := NewGithubClient("my-secret-PAT")
myLatestReleases, err := privateRepoClient.ListLatestReleases("my-org", "my-private-repo", 5)
if err != nil {
    panic(err)
}
```

There's really not much more to it...