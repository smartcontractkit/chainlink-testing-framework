package test_env

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/config/types"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/docker"
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
func GetEthereumVersionFromImage(executionLayer types.ExecutionLayer, imageWithVersion string) (config.EthereumVersion, error) {
	version, err := docker_utils.GetSemverFromImage(imageWithVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse docker image and extract version: %s", imageWithVersion)
	}

	var constraint *semver.Constraints

	switch executionLayer {
	case types.ExecutionLayer_Geth:
		constraint, err = semver.NewConstraint("<1.13.0")
	case types.ExecutionLayer_Besu:
		constraint, err = semver.NewConstraint("<23.1")
	case types.ExecutionLayer_Erigon:
		constraint, err = semver.NewConstraint("<v2.41.0")
	case types.ExecutionLayer_Nethermind:
		constraint, err = semver.NewConstraint("<1.17.0")
	case types.ExecutionLayer_Reth:
		return config.EthereumVersion_Eth2, nil
	default:
		return "", fmt.Errorf(MsgUnsupportedExecutionLayer, executionLayer)
	}

	if err != nil {
		return "", errors.New("failed to parse semver constraint for comparison")
	}

	if constraint.Check(version) {
		return config.EthereumVersion_Eth1, nil
	}
	return config.EthereumVersion_Eth2, nil
}

// UniqueStringSlice returns a deduplicated slice of strings
func UniqueStringSlice(slice []string) []string {
	stringSet := make(map[string]struct{})
	deduplicated := make([]string, 0)

	for _, el := range slice {
		if _, exists := stringSet[el]; exists {
			continue
		}

		stringSet[el] = struct{}{}
		deduplicated = append(deduplicated, el)
	}

	return deduplicated
}
