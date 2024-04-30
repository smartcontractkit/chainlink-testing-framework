package test_env

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
)

func NatPortFormat(port string) string {
	return fmt.Sprintf("%s/tcp", port)
}

func NatPort(port string) nat.Port {
	return nat.Port(NatPortFormat(port))
}

// GetHost returns the host of a container, if localhost then force ipv4 localhost
// to avoid ipv6 docker bugs https://github.com/moby/moby/issues/42442 https://github.com/moby/moby/issues/42375
func GetHost(ctx context.Context, container tc.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}
	// if localhost then force it to ipv4 localhost
	if host == "localhost" {
		host = "127.0.0.1"
	}
	return host, nil
}

// GetEndpoint returns the endpoint of a container, if localhost then force ipv4 localhost
// to avoid ipv6 docker bugs https://github.com/moby/moby/issues/42442 https://github.com/moby/moby/issues/42375
func GetEndpoint(ctx context.Context, container tc.Container, endpointType string) (string, error) {
	endpoint, err := container.Endpoint(ctx, endpointType)
	if err != nil {
		return "", err
	}
	return strings.Replace(endpoint, "localhost", "127.0.0.1", 1), nil
}

func FormatHttpUrl(host string, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}

func FormatWsUrl(host string, port string) string {
	return fmt.Sprintf("ws://%s:%s", host, port)
}

// GetEthereumVersionFromImage returns the consensus type based on the Docker image version
func GetEthereumVersionFromImage(executionLayer config.ExecutionLayer, imageWithVersion string) (config.EthereumVersion, error) {
	version, err := GetComparableVersionFromDockerImage(imageWithVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse docker image and extract version: %s", imageWithVersion)
	}
	switch executionLayer {
	case config.ExecutionLayer_Geth:
		if version < 113 {
			return config.EthereumVersion_Eth1, nil
		} else {
			return config.EthereumVersion_Eth2, nil
		}
	case config.ExecutionLayer_Besu:
		if version < 231 {
			return config.EthereumVersion_Eth1, nil
		} else {
			return config.EthereumVersion_Eth2, nil
		}
	case config.ExecutionLayer_Erigon:
		if version < 241 {
			return config.EthereumVersion_Eth1, nil
		} else {
			return config.EthereumVersion_Eth2, nil
		}
	case config.ExecutionLayer_Nethermind:
		if version < 117 {
			return config.EthereumVersion_Eth1, nil
		} else {
			return config.EthereumVersion_Eth2, nil
		}
	}

	return "", fmt.Errorf(MsgUnsupportedExecutionLayer, executionLayer)
}

// GetComparableVersionFromDockerImage returns version in xy format removing all non-numeric characters
// and patch version if present. So x.y.z becomes xy.
func GetComparableVersionFromDockerImage(imageWithVersion string) (int, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return -1, fmt.Errorf(MsgInvalidDockerImageFormat, imageWithVersion)
	}

	re := regexp.MustCompile("[a-zA-Z]")
	cleanedVersion := re.ReplaceAllString(parts[1], "")
	if idx := strings.Index(cleanedVersion, "-"); idx != -1 {
		cleanedVersion = string(cleanedVersion[:idx])
	}
	// remove patch version if present
	if count := strings.Count(cleanedVersion, "."); count > 1 {
		cleanedVersion = string(cleanedVersion[:strings.LastIndex(cleanedVersion, ".")])
	}
	version, err := strconv.Atoi(strings.Replace(cleanedVersion, ".", "", -1))
	if err != nil {
		return -1, fmt.Errorf("failed to parse docker version to an integer: %s", cleanedVersion)
	}

	return version, nil
}

// GetGithubRepositoryFromEthereumClientDockerImage returns the GitHub repository name based on the Docker image
func GetGithubRepositoryFromEthereumClientDockerImage(imageWithVersion string) (string, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(MsgInvalidDockerImageFormat, imageWithVersion)
	}

	switch {
	case strings.Contains(parts[0], gethBaseImageName):
		return gethGitRepo, nil
	case strings.Contains(parts[0], besuBaseImageName):
		return besuGitRepo, nil
	case strings.Contains(parts[0], nethermindBaseImageName):
		return nethermindGitRepo, nil
	case strings.Contains(parts[0], erigonBaseImageName):
		return erigonGitRepo, nil
	default:
		return "", fmt.Errorf(MsgUnsupportedDockerImage, parts[0])
	}
}

// GetExecutionLayerFromDockerImage returns the execution layer based on the Docker image
func GetExecutionLayerFromDockerImage(imageWithVersion string) (config.ExecutionLayer, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(MsgInvalidDockerImageFormat, imageWithVersion)
	}

	switch {
	case strings.Contains(parts[0], gethBaseImageName):
		return config.ExecutionLayer_Geth, nil
	case strings.Contains(parts[0], besuBaseImageName):
		return config.ExecutionLayer_Besu, nil
	case strings.Contains(parts[0], nethermindBaseImageName):
		return config.ExecutionLayer_Nethermind, nil
	case strings.Contains(parts[0], erigonBaseImageName):
		return config.ExecutionLayer_Erigon, nil
	default:
		return "", fmt.Errorf(MsgUnsupportedDockerImage, parts[0])
	}
}

// UniqueStringSlice returns a deduplicated slice of strings
func UniqueStringSlice(slice []string) []string {
	addressSet := make(map[string]struct{})
	deduplicated := make([]string, 0)

	for _, el := range slice {
		if _, exists := addressSet[el]; exists {
			continue
		}

		addressSet[el] = struct{}{}
		deduplicated = append(deduplicated, el)
	}

	return deduplicated
}
