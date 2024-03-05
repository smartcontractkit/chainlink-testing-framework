package test_env

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	tc "github.com/testcontainers/testcontainers-go"
)

const (
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`

	DEFAULT_EVM_NODE_HTTP_PORT = "8544"
	DEFAULT_EVM_NODE_WS_PORT   = "8545"
)

type EthereumVersion = string

const (
	EthereumVersion_Eth1 = "eth1"
	EthereumVersion_Eth2 = "eth2"
)

type ExecutionClient interface {
	GetContainerName() string
	StartContainer() (blockchain.EVMNetwork, error)
	GetContainer() *tc.Container
	GetContainerType() ContainerType
	GetInternalExecutionURL() string
	GetExternalExecutionURL() string
	GetInternalHttpUrl() string
	GetInternalWsUrl() string
	GetExternalHttpUrl() string
	GetExternalWsUrl() string
	GetEthereumVersion() EthereumVersion
	WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error
	WithTestInstance(t *testing.T) ExecutionClient
}
