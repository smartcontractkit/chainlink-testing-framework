package ethereum

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/docker"
)

type Fork string

const (
	EthereumFork_Shanghai Fork = "Shanghai"
	EthereumFork_Deneb    Fork = "Deneb"
)

// ValidFutureForks returns the list of valid future forks for the given Ethereum fork
func (e Fork) ValidFutureForks() ([]Fork, error) {
	switch e {
	case EthereumFork_Shanghai:
		return []Fork{EthereumFork_Deneb}, nil
	case EthereumFork_Deneb:
		return []Fork{}, nil
	default:
		return []Fork{}, fmt.Errorf("unknown fork: %s", e)
	}
}

// LastSupportedForkForEthereumClient returns the last supported fork for the given Ethereum client. It supports only eth2 clients.
func LastSupportedForkForEthereumClient(imageWithVersion string) (Fork, error) {
	ethereumVersion, err := VersionFromImage(imageWithVersion)
	if err != nil {
		return "", err
	}

	//nolint:staticcheck //ignore SA1019
	if ethereumVersion == config_types.EthereumVersion_Eth1 || ethereumVersion == config_types.EthereumVersion_Eth1_Legacy {
		return "", fmt.Errorf("ethereum version '%s' is not supported", ethereumVersion)
	}

	executionLayer, err := ExecutionLayerFromDockerImage(imageWithVersion)
	if err != nil {
		return "", err
	}

	version, err := docker_utils.GetSemverFromImage(imageWithVersion)
	if err != nil {
		return "", err
	}

	var constraint *semver.Constraints

	switch executionLayer {
	case config_types.ExecutionLayer_Besu:
		constraint, err = semver.NewConstraint("<24.1")
	case config_types.ExecutionLayer_Geth:
		constraint, err = semver.NewConstraint("<1.13.12")
	case config_types.ExecutionLayer_Erigon:
		constraint, err = semver.NewConstraint("<v2.59.0")
	case config_types.ExecutionLayer_Nethermind:
		constraint, err = semver.NewConstraint("<v1.26.0")
	case config_types.ExecutionLayer_Reth:
		constraint, err = semver.NewConstraint("<v1.0.0")
	default:
		return "", fmt.Errorf("unsupported execution layer: %s", executionLayer)
	}

	if err != nil {
		return "", fmt.Errorf("failed to parse semver constraint for comparison: %w", err)
	}

	if constraint.Check(version) {
		return EthereumFork_Shanghai, nil
	}
	return EthereumFork_Deneb, nil
}

// VersionFromImage returns the consensus type based on the Docker image version
func VersionFromImage(imageWithVersion string) (config_types.EthereumVersion, error) {
	version, err := docker_utils.GetSemverFromImage(imageWithVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse docker image and extract version: %s", imageWithVersion)
	}

	var constraint *semver.Constraints

	executionLayer, err := ExecutionLayerFromDockerImage(imageWithVersion)
	if err != nil {
		return "", err
	}

	switch executionLayer {
	case config_types.ExecutionLayer_Geth:
		constraint, err = semver.NewConstraint("<1.13.0")
	case config_types.ExecutionLayer_Besu:
		constraint, err = semver.NewConstraint("<23.1")
	case config_types.ExecutionLayer_Erigon:
		constraint, err = semver.NewConstraint("<v2.41.0")
	case config_types.ExecutionLayer_Nethermind:
		constraint, err = semver.NewConstraint("<1.17.0")
	case config_types.ExecutionLayer_Reth:
		return config_types.EthereumVersion_Eth2, nil
	default:
		return "", fmt.Errorf(config_types.MsgUnsupportedExecutionLayer, executionLayer)
	}

	if err != nil {
		return "", errors.New("failed to parse semver constraint for comparison")
	}

	if constraint.Check(version) {
		return config_types.EthereumVersion_Eth1, nil
	}
	return config_types.EthereumVersion_Eth2, nil
}

// ExecutionLayerFromDockerImage returns the execution layer based on the Docker image
func ExecutionLayerFromDockerImage(imageWithVersion string) (config_types.ExecutionLayer, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(config_types.MsgInvalidDockerImageFormat, imageWithVersion)
	}

	switch {
	case strings.Contains(parts[0], GethBaseImageName):
		return config_types.ExecutionLayer_Geth, nil
	case strings.Contains(parts[0], BesuBaseImageName):
		return config_types.ExecutionLayer_Besu, nil
	case strings.Contains(parts[0], NethermindBaseImageName):
		return config_types.ExecutionLayer_Nethermind, nil
	case strings.Contains(parts[0], ErigonBaseImageName):
		return config_types.ExecutionLayer_Erigon, nil
	case strings.Contains(parts[0], RethBaseImageName):
		return config_types.ExecutionLayer_Reth, nil
	default:
		return "", fmt.Errorf(config_types.MsgUnsupportedDockerImage, parts[0])
	}
}

// GithubRepositoryFromEthereumClientDockerImage returns the GitHub repository name based on the Docker image
func GithubRepositoryFromEthereumClientDockerImage(imageWithVersion string) (string, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(config_types.MsgInvalidDockerImageFormat, imageWithVersion)
	}

	switch {
	case strings.Contains(parts[0], GethBaseImageName):
		return gethGitRepo, nil
	case strings.Contains(parts[0], BesuBaseImageName):
		return besuGitRepo, nil
	case strings.Contains(parts[0], NethermindBaseImageName):
		return nethermindGitRepo, nil
	case strings.Contains(parts[0], ErigonBaseImageName):
		return erigonGitRepo, nil
	case strings.Contains(parts[0], RethBaseImageName):
		return rethGitRepo, nil
	default:
		return "", fmt.Errorf(config_types.MsgUnsupportedDockerImage, parts[0])
	}
}
