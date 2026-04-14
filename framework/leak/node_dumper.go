package leak

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var containerNameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

const (
	DefaultAdminProfilesDir       = "admin-profiles"
	DefaultNodeProfileDumpTimeout = 5 * time.Minute
)

// DumpNodeProfiles runs chainlink profile collection in each running container
// with a name containing namePattern and saves ./profiles as dst/profile-<container-name>.tar.
func DumpNodeProfiles(ctx context.Context, namePattern, dst string) error {
	f.L.Info().
		Str("NamePattern", namePattern).
		Str("DestinationDir", dst).
		Msg("Dumping node profiles by container name pattern")

	if strings.TrimSpace(namePattern) == "" {
		return fmt.Errorf("container name pattern must not be empty")
	}
	if strings.TrimSpace(dst) == "" {
		return fmt.Errorf("destination path must not be empty")
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory %q: %w", dst, err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()
	dc, err := f.NewDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create framework docker client: %w", err)
	}

	containers, err := runningContainers(ctx, cli)
	if err != nil {
		return err
	}

	var errs []error
	for _, c := range containers {
		if !strings.Contains(c.name, namePattern) {
			continue
		}

		// Keep destination names safe and filesystem-friendly.
		safeName := containerNameSanitizer.ReplaceAllString(c.name, "_")
		targetArchivePath := filepath.Join(dst, fmt.Sprintf("profile-%s.tar", safeName))

		f.L.Info().Str("ContainerName", c.name).Msg("Collecting node profile")

		out, execErr := dc.ExecContainerWithContext(
			ctx,
			c.name,
			[]string{"chainlink", "admin", "profile", "-seconds", "1", "-output_dir", "./profiles"},
		)
		if execErr != nil {
			errs = append(errs, fmt.Errorf("failed to execute profile command in container %s: %w, output: %s", c.name, execErr, strings.TrimSpace(out)))
			continue
		}

		profilesPath := path.Clean(path.Join(c.workingDir, "profiles"))
		if copyErr := dc.CopyFromContainerToTarWithContext(ctx, c.name, profilesPath, targetArchivePath); copyErr != nil {
			errs = append(errs, fmt.Errorf("failed to copy profiles archive from container %s to %s: %w", c.name, targetArchivePath, copyErr))
			continue
		}

		f.L.Info().Str("ContainerName", c.name).Str("Destination", targetArchivePath).Msg("Profiles copied as archive")
	}

	return errors.Join(errs...)
}

type runningContainer struct {
	name       string
	workingDir string
}

func runningContainers(ctx context.Context, cli *client.Client) ([]runningContainer, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list running Docker containers: %w", err)
	}

	res := make([]runningContainer, 0, len(containers))
	for _, c := range containers {
		name := firstContainerName(c.Names)
		if name == "" {
			continue
		}

		inspect, inspectErr := cli.ContainerInspect(ctx, c.ID)
		if inspectErr != nil {
			return nil, fmt.Errorf("failed to inspect container %s: %w", name, inspectErr)
		}
		workingDir := "/"
		if inspect.Config != nil && inspect.Config.WorkingDir != "" {
			workingDir = inspect.Config.WorkingDir
		}
		res = append(res, runningContainer{
			name:       name,
			workingDir: workingDir,
		})
	}
	return res, nil
}

func firstContainerName(names []string) string {
	for _, n := range names {
		if n == "" {
			continue
		}
		return strings.TrimPrefix(n, "/")
	}
	return ""
}
