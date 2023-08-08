package blockchain

// BSCMultinodeClient represents a multi-node, EVM compatible client for the BSC network
type BSCMultinodeClient struct {
	*EthereumMultinodeClient
}

// BSCClient represents a single node, EVM compatible client for the BSC network
type BSCClient struct {
	*EthereumClient
}
