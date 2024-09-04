package docker

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

const MsgInvalidDockerImageFormat = "invalid docker image format: %s"

// GetSemverFromImage returns a semver version from a docker image string
func GetSemverFromImage(imageWithVersion string) (*semver.Version, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf(MsgInvalidDockerImageFormat, imageWithVersion)
	}

	parsedVersion, err := semver.NewVersion(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse docker version to a semver: %s", parts[1])
	}

	return parsedVersion, nil
}
