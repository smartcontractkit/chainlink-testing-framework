package blockchain

// ArbitrumMultinodeClient represents a multi-node, EVM compatible client for the Arbitrum network
type ArbitrumMultinodeClient struct {
	*EthereumMultinodeClient
}

// ArbitrumClient represents a single node, EVM compatible client for the Arbitrum network
type ArbitrumClient struct {
	*EthereumClient
}
