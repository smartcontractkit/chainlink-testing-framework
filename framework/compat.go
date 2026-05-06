package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/sync/errgroup"

	"github.com/Masterminds/semver/v3"
)

/*
 * This file contains functions to verify backward/forward/inter-product compatibility of products
 * which are using devenv.
 */

const (
	CISummaryFile = "ci_summary.txt"
)

// UpgradeContext contains all the data needed for an upgrade test run
type UpgradeContext struct {
	ProductName      string
	Refs             []string
	DonNodes         int
	UpgradeNodes     int
	NodeNameTemplate string
	Registry         string
	Buildcmd         string
	Envcmd           string
	Testcmd          string
	SkipPull         bool
}

// UpgradeNRollingSummaryTemplate holds the data for rendering a rolling N upgrade summary
type UpgradeNRollingSummaryTemplate struct {
	Total    int
	Earliest string
	Latest   string
	Sequence []string
}

// WriteRollingNUpgradeSummary renders an upgrade summary and writes it to ci_summary.txt
func WriteRollingNUpgradeSummary(tmpl UpgradeNRollingSummaryTemplate) error {
	r, err := RenderTemplate(`
Testing upgrade sequence to {{.Latest}} from previous versions:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{{- range .Sequence}}
  • {{.}}
{{- end}}
	`, tmpl)
	if err != nil {
		return fmt.Errorf("failed to render upgrade summary: %w", err)
	}
	return os.WriteFile(CISummaryFile, []byte(r), 0o600)
}

// UpgradeSOTDONSummary holds the data for rendering a SOT DON upgrade summary
type UpgradeSOTDONSummary struct {
	ProductName    string
	TotalRefs      int
	DONSize        int
	Earliest       string
	Latest         string
	Sequence       []string
	SequenceChunks [][]string
}

// WriteSOTDONUpgradeSummary renders an upgrade summary and writes it to ci_summary.txt
func WriteSOTDONUpgradeSummary(tmpl UpgradeSOTDONSummary) error {
	r, err := RenderTemplate(`
Testing upgrade sequence to {{.Latest}} for DON versions from RANE SOT, for product {{.ProductName}}:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{{- range $index, $version := .Sequence}}
  • {{$version}}
{{- end}}

Mapping {{.TotalRefs}} available unique versions to the current DON size of {{.DONSize}} nodes:
{{- range $chunkIndex, $chunk := .SequenceChunks}}
Sequence #{{$chunkIndex}}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  {{- range $versionIndex, $version := $chunk}}
    • {{$version}} -> {{$.Latest}}
  {{- end}}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{{- end}}`, tmpl)
	if err != nil {
		return fmt.Errorf("failed to render upgrade summary: %w", err)
	}
	return os.WriteFile(CISummaryFile, []byte(r), 0o600)
}

func RunInitialUpgradeSetup(ctx context.Context, u UpgradeContext) error {
	L.Info().Strs("Sequence", u.Refs).Msg("Running upgrade sequence")
	// first env var is used by devenv, second by simple nodeset to override the image
	os.Setenv("CHAINLINK_IMAGE", fmt.Sprintf("%s:%s", u.Registry, u.Refs[0]))
	os.Setenv("CTF_CHAINLINK_IMAGE", fmt.Sprintf("%s:%s", u.Registry, u.Refs[0]))
	if _, err := ExecCmdWithContext(ctx, u.Buildcmd); err != nil {
		return err
	}
	if _, err := ExecCmdWithContext(ctx, u.Envcmd); err != nil {
		return err
	}
	if _, err := ExecCmdWithContext(ctx, u.Testcmd); err != nil {
		return err
	}
	return nil
}

// prepareRefsForDON chunks all available refs to DONs
func prepareRefsForDON(slice []string, donSize int) [][]string {
	if donSize <= 0 || len(slice) == 0 {
		return [][]string{}
	}
	var chunks [][]string
	for i := 0; i < len(slice); i += donSize {
		end := i + donSize
		if end > len(slice) {
			chunks = append(chunks, slice[i:])
		} else {
			chunks = append(chunks, slice[i:end])
		}
	}
	return chunks
}

// UpgradeNProductUniqueVersionsRolling models multiple product DONs as one having all unique versions
// first upgrades the current version of each node to modelled DON version in parallel and tests it once
// then upgrades each node to the target version and tests it after each upgrade
// [v1.0, v1.0, v1.1, v1.1] -> test it
// [v1.0, v1.1, v1.3, v1.4] -> get N unique refs and boot up DON to simulate a real one
// [vX.X, v1.1, v1.3, v1.4] -> upgrade a single node, test
// [vX.X, vX.X, v1.3, v1.4] -> upgrade a single node, test
// ... etc
// [v1.0, v1.1, v1.3, v1.4] -> get N unique refs and boot up DON to simulate a real one
// [vX.X, v1.1, v1.3, v1.4] -> upgrade a single node, test
// [vX.X, vX.X, v1.3, v1.4] -> upgrade a single node, test
// ... etc
func UpgradeNProductUniqueVersionsRolling(ctx context.Context, u UpgradeContext) error {
	// first model a DON from unique refs, grab N refs to fit the DON size
	donRefs := prepareRefsForDON(u.Refs, u.DonNodes)

	// this is used to form CI summary
	if err := WriteSOTDONUpgradeSummary(
		UpgradeSOTDONSummary{
			ProductName:    u.ProductName,
			TotalRefs:      len(u.Refs),
			DONSize:        u.DonNodes,
			Earliest:       u.Refs[0],
			Latest:         u.Refs[len(u.Refs)-1],
			Sequence:       u.Refs[:len(u.Refs)-1],
			SequenceChunks: donRefs,
		},
	); err != nil {
		return err
	}

	// upgrade to DON versions, then upgrade to target version one by one
	eg := errgroup.Group{}
	for _, refChunk := range donRefs {
		// reboot before each chunk (DON) iteration to clean up the data
		// this is needed because we can only upgrade up
		if err := RunInitialUpgradeSetup(ctx, u); err != nil {
			return err
		}
		// for each chunk, model the DON versions and test once after upgrade
		for i, ref := range refChunk {
			img := fmt.Sprintf("%s:%s", u.Registry, ref)
			eg.Go(func() error {
				if !u.SkipPull {
					if _, err := ExecCmdWithContext(ctx, fmt.Sprintf("docker pull %s", img)); err != nil {
						return fmt.Errorf("failed to pull image %s: %w", img, err)
					}
				}
				return UpgradeContainer(ctx, fmt.Sprintf(u.NodeNameTemplate, i), img)
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		// run the tests
		if _, err := ExecCmdWithContext(ctx, u.Testcmd); err != nil {
			return err
		}

		// if all good upgrade each node to the latest version
		img := fmt.Sprintf("%s:%s", u.Registry, u.Refs[len(u.Refs)-1])
		if !u.SkipPull {
			if _, err := ExecCmdWithContext(ctx, fmt.Sprintf("docker pull %s", img)); err != nil {
				return fmt.Errorf("failed to pull image %s: %w", img, err)
			}
		}
		for i := range u.DonNodes {
			if err := UpgradeContainer(ctx, fmt.Sprintf(u.NodeNameTemplate, i), img); err != nil {
				return err
			}
		}
		// run the tests again
		if _, err := ExecCmdWithContext(ctx, u.Testcmd); err != nil {
			return err
		}
	}
	return nil
}

// UpgradeNRolling upgrades nodes to the given refs in rolling fashion, testing with the given testcmd after each upgrade
// [v1.0, v1.0, v1.1, v1.1] -> test
// [v1.0, v1.0, v1.2, v1.2] -> test
// [v1.0, v1.0, v1.3, v1.3] -> test
// etc
func UpgradeNRolling(ctx context.Context, u UpgradeContext) error {
	// this is used to form CI summary
	if err := WriteRollingNUpgradeSummary(
		UpgradeNRollingSummaryTemplate{
			Total:    len(u.Refs),
			Earliest: u.Refs[0],
			Latest:   u.Refs[len(u.Refs)-1],
			Sequence: u.Refs[:len(u.Refs)-1],
		},
	); err != nil {
		return err
	}

	if err := RunInitialUpgradeSetup(ctx, u); err != nil {
		return err
	}

	for _, tag := range u.Refs {
		L.Info().
			Int("Nodes", u.UpgradeNodes).
			Str("Version", tag).
			Msg("Upgrading nodes")
		img := fmt.Sprintf("%s:%s", u.Registry, tag)
		if !u.SkipPull {
			if _, err := ExecCmdWithContext(ctx, fmt.Sprintf("docker pull %s", img)); err != nil {
				return fmt.Errorf("failed to pull image %s: %w", img, err)
			}
		}
		for i := range u.UpgradeNodes {
			if err := UpgradeContainer(ctx, fmt.Sprintf(u.NodeNameTemplate, i), img); err != nil {
				return err
			}
		}
		if _, err := ExecCmdWithContext(ctx, u.Testcmd); err != nil {
			return err
		}
	}
	return nil
}

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
	l.Info().Msg("Container successfully rebooted with new image")
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
		NOP  string `json:"nop"`
		Jobs []struct {
			Product string `json:"product"`
		} `json:"jobs"`
		Version string `json:"version"`
	} `json:"nodes"`
}

// fetchSOTSource fetches the SOT source from the given URL and returns the parsed response body
func fetchSOTSource(url string) (RaneSOTResponseBody, error) {
	L.Info().
		Str("URL", url).
		Msg("Fetching versions from snapshot")
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return RaneSOTResponseBody{}, fmt.Errorf("failed to fetch SOT URL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return RaneSOTResponseBody{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RaneSOTResponseBody{}, fmt.Errorf("failed to read SOT URL response body: %w", err)
	}
	var response RaneSOTResponseBody

	if err := json.Unmarshal(body, &response); err != nil {
		return RaneSOTResponseBody{}, fmt.Errorf("failed to parse SOT response body: %w", err)
	}
	return response, nil
}

// FindNOPsVersionsByProduct finds all the versions of a given product from a RANE SOT source across all the NOPs
func FindNOPsVersionsByProduct(url string, product string, exclude []string) ([]string, error) {
	src, err := fetchSOTSource(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SOT data: %w", err)
	}
	refs := make([]string, 0)
	products := make(map[string]bool)
	seenRefs := make(map[string]bool, 0)
	for _, n := range src.Nodes {
		for _, j := range n.Jobs {
			products[j.Product] = true
			if j.Product == product {
				if _, ok := seenRefs[n.Version]; !ok {
					refs = append(refs, n.Version)
					seenRefs[n.Version] = true
				}
			}
		}
	}
	semverTags := FilterSemverTags(refs, []string{}, exclude)
	slices.Reverse(semverTags)
	L.Info().Any("Products", slices.Collect(maps.Keys(products))).Msg("Found products")
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
