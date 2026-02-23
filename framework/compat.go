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
	L.Info().Msg("Upgrading container")
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()
	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to inspect container %s: %w", containerName, err)
	}
	L.Debug().Msg("Stopping container")
	stopOpts := container.StopOptions{}
	if err := cli.ContainerStop(ctx, containerName, stopOpts); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}
	L.Debug().Msg("Removing container")
	// keep the volumes
	removeOpts := container.RemoveOptions{RemoveVolumes: false}
	if err := cli.ContainerRemove(ctx, containerName, removeOpts); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerName, err)
	}
	L.Debug().Msg("Pulling new image")
	pullReader, err := cli.ImagePull(ctx, newImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", newImage, err)
	}
	defer pullReader.Close()

	// log pull process for debug
	if L.GetLevel() <= zerolog.DebugLevel {
		io.Copy(os.Stdout, pullReader)
	} else {
		io.Copy(io.Discard, pullReader)
	}
	L.Debug().Msg("Image pulled successfully")

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
		inspect.Config,
		inspect.HostConfig,
		networkingConfig,
		nil,
		containerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create container with image %s: %w", newImage, err)
	}
	L.Debug().
		Str("ContainerID", createResp.ID).
		Msg("Container created")
	L.Debug().Msg("Starting new container")
	startOpts := container.StartOptions{}
	if err := cli.ContainerStart(ctx, createResp.ID, startOpts); err != nil {
		return fmt.Errorf("failed to start container %s: %w", containerName, err)
	}
	L.Info().
		Str("ContainerID", createResp.ID[:12]).
		Msg("Container successfully rebooted with new image")
	return nil
}

// RestoreToBranch restores git back to the develop branch
func RestoreToBranch(baseBranch string) error {
	_, err := ExecCmd(fmt.Sprintf("git checkout %s", baseBranch))
	if err != nil {
		return fmt.Errorf("failed to checkout develop branch: %w", err)
	}
	L.Info().
		Str("Branch", "develop").
		Msg("Successfully restored to develop branch")
	return nil
}

// RollbackToEarliestSemverTag gets all semver tags, sorts them, and rolls back to the earliest tag
// returns all the tags starting from the oldest one
func RollbackToEarliestSemverTag(tagsBack int, include, exclude []string) ([]string, error) {
	output, err := ExecCmd("git tag --list")
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

	_, err = ExecCmd("git checkout " + earliestTag)
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

// SortSemverTags parses valid versions and returns them sorted from latest to lowest
func SortSemverTags(versions []string, include []string, exclude []string) []string {
	if len(versions) == 0 {
		return []string{}
	}
	// parse semver, exclude invalid tags
	parsedVersions := make([]*semver.Version, 0, len(versions))
	for _, v := range versions {
		parsed, err := semver.NewVersion(v)
		if err != nil {
			L.Trace().Str("tag", v).Msg("Skipping invalid semver tag")
			continue
		}
		parsedVersions = append(parsedVersions, parsed)
	}

	// descending order
	sort.Slice(parsedVersions, func(i, j int) bool {
		return parsedVersions[i].GreaterThan(parsedVersions[j])
	})

	// fliter include/exclude
	filtered := filterVersions(parsedVersions, include, exclude)
	L.Info().
		Strs("include", include).
		Strs("exclude", exclude).
		Msg("Applied filters")
	return filtered
}

// filterVersions applies include/exclude filters to parsed versions
func filterVersions(versions []*semver.Version, include, exclude []string) []string {
	if len(include) == 0 {
		include = []string{""} // Include all by default
	}
	result := make([]string, 0, len(versions))

	for _, v := range versions {
		original := v.Original()
		if matchesFilter(original, include, exclude) {
			result = append(result, original)
		}
	}
	return result
}

// matchesFilter checks if a version string matches the include/exclude criteria
func matchesFilter(version string, include, exclude []string) bool {
	for _, pattern := range exclude {
		if strings.Contains(version, pattern) {
			return false
		}
	}
	for _, pattern := range include {
		if strings.Contains(version, pattern) {
			return true
		}
	}
	return false
}

func matchesAny(s string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
