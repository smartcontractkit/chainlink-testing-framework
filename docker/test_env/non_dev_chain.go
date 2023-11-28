package test_env

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

type NonDevNode interface {
	GetInternalHttpUrl() string
	GetInternalWsUrl() string
	GetEVMClient() blockchain.EVMClient
	WithTestInstance(t *testing.T) NonDevNode
	Start() error
	ConnectToClient() error
}

type PrivateChain interface {
	GetPrimaryNode() NonDevNode
	GetNodes() []NonDevNode
	GetNetworkConfig() *blockchain.EVMNetwork
	GetDockerNetworks() []string
}
