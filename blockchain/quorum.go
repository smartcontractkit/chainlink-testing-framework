package blockchain

// Handles specific issues with the Quorum EVM chain: https://docs.Quorum.com/

// QuorumMultinodeClient represents a multi-node, EVM compatible client for the Quorum network
type QuorumMultinodeClient struct {
	*EthereumMultinodeClient
}

// QuorumClient represents a single node, EVM compatible client for the Quorum network
type QuorumClient struct {
	*EthereumClient
}
