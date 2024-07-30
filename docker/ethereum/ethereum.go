package ethereum

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-testing-framework/config/types"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/docker"
)

type Fork string

const (
	EthereumFork_Shanghai Fork = "Shanghai"
	EthereumFork_Deneb    Fork = "Deneb"
)

func (e Fork) GetValidFutureForks() ([]Fork, error) {
	switch e {
	case EthereumFork_Shanghai:
		return []Fork{EthereumFork_Deneb}, nil
	case EthereumFork_Deneb:
		return []Fork{}, nil
	default:
		return []Fork{}, fmt.Errorf("unknown fork: %s", e)
	}
}

const (
	MsgInvalidDockerImageFormat = "invalid docker image format: %s"
	MsgUnsupportedDockerImage   = "unsupported docker image: %s"
)

func GetLastSupportedForkForEthereumClient(imageWithVersion string) (Fork, error) {
	executionLayer, err := GetExecutionLayerFromDockerImage(imageWithVersion)
	if err != nil {
		return "", err
	}

	version, err := docker_utils.GetSemverFromImage(imageWithVersion)
	if err != nil {
		return "", err
	}

	var constraint *semver.Constraints

	switch executionLayer {
	case types.ExecutionLayer_Besu:
		constraint, err = semver.NewConstraint("<24.1")
	case types.ExecutionLayer_Geth:
		constraint, err = semver.NewConstraint("<1.13.12")
	case types.ExecutionLayer_Erigon:
		constraint, err = semver.NewConstraint("<v2.59.0")
	case types.ExecutionLayer_Nethermind:
		constraint, err = semver.NewConstraint("<v1.26.0")
	case types.ExecutionLayer_Reth:
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

// GetExecutionLayerFromDockerImage returns the execution layer based on the Docker image
func GetExecutionLayerFromDockerImage(imageWithVersion string) (types.ExecutionLayer, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(MsgInvalidDockerImageFormat, imageWithVersion)
	}

	switch {
	case strings.Contains(parts[0], GethBaseImageName):
		return types.ExecutionLayer_Geth, nil
	case strings.Contains(parts[0], BesuBaseImageName):
		return types.ExecutionLayer_Besu, nil
	case strings.Contains(parts[0], NethermindBaseImageName):
		return types.ExecutionLayer_Nethermind, nil
	case strings.Contains(parts[0], ErigonBaseImageName):
		return types.ExecutionLayer_Erigon, nil
	case strings.Contains(parts[0], RethBaseImageName):
		return types.ExecutionLayer_Reth, nil
	default:
		return "", fmt.Errorf(MsgUnsupportedDockerImage, parts[0])
	}
}

// GetGithubRepositoryFromEthereumClientDockerImage returns the GitHub repository name based on the Docker image
func GetGithubRepositoryFromEthereumClientDockerImage(imageWithVersion string) (string, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf(MsgInvalidDockerImageFormat, imageWithVersion)
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
		return "", fmt.Errorf(MsgUnsupportedDockerImage, parts[0])
	}
}
