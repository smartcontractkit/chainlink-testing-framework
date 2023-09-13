package test_env

import "github.com/smartcontractkit/chainlink-testing-framework/blockchain"

type NonDevNode interface {
	GetInternalHttpUrl() string
	GetInternalWsUrl() string
	GetEVMClient() blockchain.EVMClient
	Start() error
	ConnectToClient() error
}

type PrivateChain interface {
	GetPrimaryNode() NonDevNode
	GetNodes() []NonDevNode
	GetNetworkConfig() *blockchain.EVMNetwork
	GetDockerNetworks() []string
}
