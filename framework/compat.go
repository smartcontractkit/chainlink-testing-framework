package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"

	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

/*
 * This file contains functions to verify backward/forward/inter-product compatibility of products
 * which are using devenv.
 */

// UpgradeContainer stops a container, removes it, and creates a new one with the specified image
func UpgradeContainer(ctx context.Context, containerName, newImage string) error {
	L = L.With().
		Str("Container", containerName).
		Str("Image", newImage).
		Logger()
	L.Info().Msg("Starting container reboot with new image")
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		L.Error().Err(err).Msg("Failed to create Docker client")
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()
	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		L.Error().Err(err).Msg("Failed to inspect container")
		return fmt.Errorf("failed to inspect container %s: %w", containerName, err)
	}
	L.Info().Msg("Stopping container")
	stopOpts := container.StopOptions{}
	if err := cli.ContainerStop(ctx, containerName, stopOpts); err != nil {
		L.Error().Err(err).Msg("Failed to stop container")
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}
	L.Info().Msg("Container stopped successfully")
	L.Info().Msg("Removing container")
	// keep the volumes
	removeOpts := container.RemoveOptions{RemoveVolumes: false}
	if err := cli.ContainerRemove(ctx, containerName, removeOpts); err != nil {
		L.Error().Err(err).Msg("Failed to remove container")
		return fmt.Errorf("failed to remove container %s: %w", containerName, err)
	}
	L.Info().Msg("Container removed successfully")
	L.Info().Msg("Pulling new image")
	pullReader, err := cli.ImagePull(ctx, newImage, image.PullOptions{})
	if err != nil {
		L.Error().Err(err).Msg("Failed to pull image")
		return fmt.Errorf("failed to pull image %s: %w", newImage, err)
	}
	defer pullReader.Close()

	// Optionally log pull progress (can be verbose)
	if L.GetLevel() <= zerolog.DebugLevel {
		io.Copy(os.Stdout, pullReader)
	} else {
		// Just consume the reader to ensure pull completes
		io.Copy(io.Discard, pullReader)
	}
	L.Info().Msg("Image pulled successfully")

	// Create new container with same configuration but new image
	L.Info().Msg("Creating new container with updated image")
	inspect.Config.Image = newImage

	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"ctf": {
				NetworkID: "ctf",
			},
		},
	}
	createResp, err := cli.ContainerCreate(
		ctx,
		inspect.Config,     // Use original config with updated image
		inspect.HostConfig, // Keep same host config
		networkingConfig,
		nil,
		containerName,
	)
	if err != nil {
		L.Error().Err(err).Msg("Failed to create container")
		return fmt.Errorf("failed to create container with image %s: %w", newImage, err)
	}
	L.Debug().
		Str("ContainerID", createResp.ID).
		Msg("Container created")
	L.Info().Msg("Starting new container")
	startOpts := container.StartOptions{}
	if err := cli.ContainerStart(ctx, createResp.ID, startOpts); err != nil {
		L.Error().Err(err).Msg("Failed to start container")
		return fmt.Errorf("failed to start container %s: %w", containerName, err)
	}
	L.Info().
		Str("ContainerID", createResp.ID[:12]).
		Msg("Container successfully rebooted with new image")
	return nil
}

// RestoreToDevelop restores git back to the develop branch
func RestoreToDevelop() error {
	output, err := ExecCmd(L, "git checkout develop")
	if err != nil {
		return fmt.Errorf("failed to checkout develop branch: %w", err)
	}

	L.Info().
		Str("Branch", "develop").
		Str("output", string(output)).
		Msg("Successfully restored to develop branch")
	return nil
}

// RollbackToEarliestSemverTag gets all semver tags, sorts them, and rolls back to the earliest tag
// returns all the tags starting from the oldest one
func RollbackToEarliestSemverTag(tagsBack int, include, exclude []string) ([]string, error) {
	output, err := ExecCmd(L, "git tag --list")
	if err != nil {
		return nil, fmt.Errorf("failed to list git tags: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		return nil, fmt.Errorf("no tags found in repository")
	}

	sortedDesc := SortSemverTags(tags, include, exclude)
	if len(sortedDesc) == 0 {
		return nil, fmt.Errorf("no valid semver tags found")
	}

	remainingTags := sortedDesc
	if len(sortedDesc) > tagsBack {
		remainingTags = sortedDesc[:tagsBack]
	}
	earliestTag := remainingTags[len(remainingTags)-1]

	L.Info().
		Int("TotalValidTags", len(sortedDesc)).
		Strs("SelectedTags", remainingTags).
		Str("EarliestTag", earliestTag).
		Msg("Selected previous tag")

	output, err = ExecCmd(L, "git checkout "+earliestTag)
	if err != nil {
		L.Error().
			Str("Tag", earliestTag).
			Err(err).
			Msg("Failed to checkout tag")
		return nil, fmt.Errorf("failed to checkout tag %s: %w", earliestTag, err)
	}

	L.Info().
		Str("Tag", earliestTag).
		Msg("Successfully rolled back to tag")
	return remainingTags, nil
}

type RaneSOTResponseBody struct {
	Nodes []struct {
		NOP     string `json:"nop"`
		Version string `json:"version"`
	} `json:"nodes"`
}

// GetTagsFromURL fetches tags from a JSON endpoint and applies filtering
func GetTagsFromURL(url, imageTagSuffix string, nopSuffixes, ignores []string) ([]string, error) {
	L.Info().
		Str("URL", url).
		Str("ImageTagSuffix", imageTagSuffix).
		Strs("NOPs", nopSuffixes).
		Strs("IgnoreSuffix", ignores).
		Msg("Fetching tags from snapshot")
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var response RaneSOTResponseBody

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	tags := make([]string, 0)
	seenTags := make(map[string]bool, 0)
	// check all the nodes for uniq tags
	for _, node := range response.Nodes {
		version := node.Version
		nop := node.NOP
		// skip if we are not interested in this NOP
		for _, ns := range nopSuffixes {
			if !strings.Contains(nop, ns) {
				continue
			}
		}

		// skip if version is empty
		if version == "" {
			continue
		}

		// skip if we ignore some images
		ignored := false
		for _, ignore := range ignores {
			if strings.Contains(version, ignore) {
				ignored = true
				break
			}
		}

		if strings.Contains(version, imageTagSuffix) && !ignored {
			if _, ok := seenTags[version]; ok {
				continue
			}
			tags = append(tags, version)
			seenTags[version] = true
		}
	}
	return tags, nil
}

// GetECRRepositoryTags returns a list of image tags from an ECR repository
func GetECRRepositoryTags(suffix, repoName, registryID, region string, ignores []string) ([]string, error) {
	fmt.Printf("Fetching tags for repository: %s in region %s\n", repoName, region)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	client := ecr.NewFromConfig(cfg)

	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
		RegistryId:     aws.String(registryID),
		MaxResults:     aws.Int32(1000), // Max allowed is 1000
	}

	var allTags []string
	pageCount := 0
	imageCount := 0
	taggedImageCount := 0

	paginator := ecr.NewDescribeImagesPaginator(client, input)
	for paginator.HasMorePages() {
		pageCount++
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to describe images page %d: %w", pageCount, err)
		}

		fmt.Printf("Processing page %d with %d images\n", pageCount, len(page.ImageDetails))

		for _, image := range page.ImageDetails {
			imageCount++
			if len(image.ImageTags) > 0 {
				taggedImageCount++
				for _, t := range image.ImageTags {
					ignored := false
					for _, ignore := range ignores {
						if strings.Contains(t, ignore) {
							ignored = true
						}
					}
					if strings.Contains(t, suffix) && !ignored {
						allTags = append(allTags, t)
					}
				}
			}
		}
	}
	return allTags, nil
}

// SortSemverTags parses valid versions and returns them sorted from latest to lowest
func SortSemverTags(versions []string, include []string, exclude []string) []string {
	parsedVersions := make([]*semver.Version, 0)
	for _, v := range versions {
		parsed, err := semver.NewVersion(v)
		if err != nil {
			L.Debug().
				Str("Tag", v).
				Msg("Skipping invalid semver tag")
			continue
		}
		parsedVersions = append(parsedVersions, parsed)
	}

	sort.Slice(parsedVersions, func(i, j int) bool {
		return parsedVersions[i].GreaterThan(parsedVersions[j])
	})

	result := make([]string, len(parsedVersions))
	for i, v := range parsedVersions {
		result[i] = v.Original()
	}

	// ignore non GA tags
	tags := make([]string, 0)
	for _, r := range result {
		excluded := false
		for _, f := range exclude {
			if strings.Contains(r, f) {
				excluded = true
			}
		}
		included := false
		for _, f := range include {
			if strings.Contains(r, f) {
				included = true
			}
		}
		if included && !excluded {
			tags = append(tags, r)
		}
	}
	L.Info().
		Strs("Include", include).
		Strs("Exclude", exclude).
		Msg("Applied filters")
	return tags
}
