package client

// Generalized blockchain client for interaction with multiple different blockchains
type BlockchainClient interface {
	DeployStorageContract() string
}
