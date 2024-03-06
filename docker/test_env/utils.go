package test_env

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
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
func GetEthereumVersionFromImage(executionLayer ExecutionLayer, imageWithVersion string) (EthereumVersion, error) {
	version, err := GetComparableVersionFromDockerImage(imageWithVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse docker image and extract version: %s", imageWithVersion)
	}
	switch executionLayer {
	case ExecutionLayer_Geth:
		if version < 113 {
			return EthereumVersion_Eth1, nil
		} else {
			return EthereumVersion_Eth2, nil
		}
	case ExecutionLayer_Besu:
		if version < 231 {
			return EthereumVersion_Eth1, nil
		} else {
			return EthereumVersion_Eth2, nil
		}
	case ExecutionLayer_Erigon:
		if version < 241 {
			return EthereumVersion_Eth1, nil
		} else {
			return EthereumVersion_Eth2, nil
		}
	case ExecutionLayer_Nethermind:
		if version < 117 {
			return EthereumVersion_Eth1, nil
		} else {
			return EthereumVersion_Eth2, nil
		}
	}

	return "", fmt.Errorf("unsupported execution layer: %s", executionLayer)
}

// GetComparableVersionFromDockerImage returns version in xy format removing all non-numeric characters
// and patch version if present. So x.y.z becomes xy.
func GetComparableVersionFromDockerImage(imageWithVersion string) (int, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return -1, fmt.Errorf("invalid docker image format: %s", imageWithVersion)
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
		return -1, fmt.Errorf("failed to pase docker version to an integer: %s", cleanedVersion)
	}

	return version, nil
}

// GetGithubRepositoryFromDockerImage returns the GitHub repository name based on the Docker image
func GetGithubRepositoryFromDockerImage(imageWithVersion string) (string, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid docker image format: %s", imageWithVersion)
	}

	switch parts[0] {
	case "ethereum/client-go":
		return "ethereum/go-ethereum", nil
	case "hyperledger/besu":
		return "hyperledger/besu", nil
	case "nethermind/nethermind":
		return "NethermindEth/nethermind", nil
	case "thorax/erigon":
		return "ledgerwatch/erigon", nil
	default:
		return "", fmt.Errorf("unsupported docker image: %s", parts[0])
	}
}

func GetExecutionLayerFromDockerImage(imageWithVersion string) (ExecutionLayer, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid docker image format: %s", imageWithVersion)
	}

	switch parts[0] {
	case "ethereum/client-go":
		return ExecutionLayer_Geth, nil
	case "hyperledger/besu":
		return ExecutionLayer_Besu, nil
	case "nethermind/nethermind":
		return ExecutionLayer_Nethermind, nil
	case "thorax/erigon":
		return ExecutionLayer_Erigon, nil
	default:
		return "", fmt.Errorf("unsupported docker image: %s", parts[0])
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
