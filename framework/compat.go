package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"github.com/Masterminds/semver/v3"
)

/*
 * This file contains functions to verify backward/forward/inter-product compatibility of products
 * which are using devenv.
 */

// UpgradeContainer stops a container, removes it, and creates a new one with the specified image
func UpgradeContainer(ctx context.Context, containerName, newImage string) error {
	l := L.With().
		Str("Container", containerName).
		Str("Image", newImage).
		Logger()
	l.Debug().Msg("Upgrading container")
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
	l.Debug().Msg("Stopping container")
	stopOpts := container.StopOptions{}
	if err := cli.ContainerStop(ctx, containerName, stopOpts); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}
	l.Debug().Msg("Removing container")
	// keep the volumes
	removeOpts := container.RemoveOptions{RemoveVolumes: false}
	if err := cli.ContainerRemove(ctx, containerName, removeOpts); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerName, err)
	}

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
	l.Debug().
		Str("ContainerID", createResp.ID).
		Msg("Container created")
	l.Debug().Msg("Starting new container")
	startOpts := container.StartOptions{}
	if err := cli.ContainerStart(ctx, createResp.ID, startOpts); err != nil {
		return fmt.Errorf("failed to start container %s: %w", containerName, err)
	}
	l.Debug().
		Str("ContainerID", createResp.ID[:12]).
		Msg("Container successfully rebooted with new image")
	return nil
}

// FindSemVerRefSequence gets all semver tags, sorts them, and rolls back to the earliest tag
// returns all the tags starting from the oldest one
func FindSemVerRefSequence(tagsBack int, include, exclude []string) ([]string, error) {
	output, err := ExecCmd("git tag --list")
	if err != nil {
		return nil, fmt.Errorf("failed to list git tags: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		return nil, fmt.Errorf("no tags found in repository")
	}

	sortedDesc := FilterSemverTags(tags, include, exclude)
	if len(sortedDesc) == 0 {
		return nil, fmt.Errorf("no valid semver tags found")
	}

	remainingTags := sortedDesc
	if len(sortedDesc) > tagsBack {
		remainingTags = sortedDesc[:tagsBack]
	}

	slices.Reverse(remainingTags)
	return remainingTags, nil
}

func CheckOut(ref string) error {
	_, err := ExecCmd("git checkout " + ref)
	if err != nil {
		return fmt.Errorf("failed to checkout ref %s: %w", ref, err)
	}
	L.Info().
		Str("Ref", ref).
		Msg("Successfully rolled back to ref")
	return nil
}

type RaneSOTResponseBody struct {
	Nodes []struct {
		NOP     string `json:"nop"`
		Version string `json:"version"`
	} `json:"nodes"`
}

// FindNOPRefs fetches NOP tags from a RANE SOT source
func FindNOPRefs(url string, nopName string, exclude []string) ([]string, error) {
	L.Info().
		Str("URL", url).
		Str("NOP", nopName).
		Strs("Exclude", exclude).
		Msg("Fetching refs from snapshot")
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SOT URL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read SOT URL response body: %w", err)
	}
	var response RaneSOTResponseBody

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse SOT response body: %w", err)
	}

	refs := make([]string, 0)
	nops := make([]string, 0)
	seenRefs := make(map[string]bool, 0)
	seenNOPs := make(map[string]bool, 0)
	// check all the nodes for uniq tags
	for _, node := range response.Nodes {
		version := node.Version
		nop := node.NOP

		// add uniq NOP
		if _, ok := seenNOPs[nop]; !ok {
			nops = append(nops, nop)
			seenNOPs[nop] = true
		}
		// continue if it's not the NOP we need
		if nop != nopName {
			continue
		}
		// skip if version is empty
		if version == "" {
			continue
		}
		// add uniq version
		if _, ok := seenRefs[version]; !ok {
			refs = append(refs, version)
			seenRefs[version] = true
		}
	}
	semverTags := FilterSemverTags(refs, []string{}, exclude)
	slices.Reverse(semverTags)
	L.Info().Strs("NOPs", nops).Msg("Scanned NOPs")
	return semverTags, nil
}

// FilterSemverTags parses valid versions and returns them sorted from latest to lowest
func FilterSemverTags(versions []string, include []string, exclude []string) []string {
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
		Strs("Include", include).
		Strs("Exclude", exclude).
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
