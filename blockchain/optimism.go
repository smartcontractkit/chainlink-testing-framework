package blockchain

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Optimism network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Optimism network
type OptimismClient struct {
	*EthereumClient
}
