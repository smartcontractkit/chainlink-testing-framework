package test_env

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
)

const (
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`

	DEFAULT_EVM_NODE_HTTP_PORT = "8544"
	DEFAULT_EVM_NODE_WS_PORT   = "8545"
)

type ConsensusMechanism string

const (
	ConsensusMechanism_PoW ConsensusMechanism = "pow"
	ConsensusMechanism_PoS ConsensusMechanism = "pos"
	ConsensusMechanism_PoA ConsensusMechanism = "poa"
)

type ExecutionClient interface {
	GetContainerName() string
	StartContainer() (blockchain.EVMNetwork, error)
	GetContainer() *tc.Container
	GetInternalExecutionURL() string
	GetExternalExecutionURL() string
	GetInternalHttpUrl() string
	GetInternalWsUrl() string
	GetExternalHttpUrl() string
	GetExternalWsUrl() string
	GetEthereumVersion() config_types.EthereumVersion
	GetConsensusMechanism() ConsensusMechanism
	WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error
	WithTestInstance(t *testing.T) ExecutionClient
}

type UnsupportedVersion struct {
	DockerImage string
	Reason      string
}

var UNSUPPORTED_VERSIONS = []UnsupportedVersion{
	{DockerImage: fmt.Sprintf("%s:1.20.0", ethereum.NethermindBaseImageName),
		Reason: "1.20.0 was replaced with 1.20.1, for more info check https://github.com/NethermindEth/nethermind/releases/tag/1.20.0",
	},
	{DockerImage: fmt.Sprintf("%s:v1.9.0", ethereum.GethBaseImageName),
		Reason: "v1.9.0 randomly drops websocket connections, for more info check https://github.com/ethereum/go-ethereum/issues/19001",
	},
}

// IsDockerImageVersionSupported checks if the given docker image version is supported and if not returns the reason why
func IsDockerImageVersionSupported(imageWithVersion string) (bool, string, error) {
	parts := strings.Split(imageWithVersion, ":")
	if len(parts) != 2 {
		return false, "", fmt.Errorf("invalid docker image format: %s", imageWithVersion)
	}

	for _, unsp := range UNSUPPORTED_VERSIONS {
		if strings.Contains(imageWithVersion, unsp.DockerImage) {
			return false, unsp.Reason, nil
		}
	}
	return true, "", nil
}
