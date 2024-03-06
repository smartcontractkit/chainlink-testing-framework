package test_env

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/client"
)

const AUTOMATIC_LATEST_TAG = "latest_available"
const AUTOMATIC_STABLE_LATEST_TAG = "latest_stable"

func fetchLatestIfNeed(dockerImageWithVersion string) (string, error) {
	imageParts := strings.Split(dockerImageWithVersion, ":")
	if len(imageParts) != 2 {
		return "", fmt.Errorf("expected correctly formatted docker image, but got '%s'", dockerImageWithVersion)
	}

	ghRepo, err := GetGithubRepositoryFromDockerImage(dockerImageWithVersion)
	if err != nil {
		return "", err
	}

	if (imageParts[1] == AUTOMATIC_LATEST_TAG) || (imageParts[1] == AUTOMATIC_STABLE_LATEST_TAG) {
		repoParts := strings.Split(ghRepo, "/")
		if len(repoParts) != 2 {
			return "", fmt.Errorf("full github repository must have org and repository names, but '%s' does not", ghRepo)
		}

		ghClient := client.NewGithubClient(client.WITHOUT_TOKEN)
		latestTags, err := ghClient.ListLatestReleases(repoParts[0], repoParts[1], 20)
		if err != nil {
			return "", err
		}

		var latestTag string
		for _, tag := range latestTags {
			if imageParts[1] == AUTOMATIC_STABLE_LATEST_TAG {
				if tag.Prerelease != nil && *tag.Prerelease {
					continue
				}
				if tag.Draft != nil && *tag.Draft {
					continue
				}
			}
			if tag.TagName != nil {
				latestTag = *tag.TagName
				break
			}
		}

		if latestTag == "" {
			return "", fmt.Errorf("no latest tag found for %s", ghRepo)
		}

		return fmt.Sprintf("%s:%s", imageParts[0], latestTag), nil
	}

	return dockerImageWithVersion, nil
}
