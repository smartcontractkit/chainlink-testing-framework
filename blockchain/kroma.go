package blockchain

// KromaMultinodeClient represents a multi-node, EVM compatible client for the Kroma network
type KromaMultinodeClient struct {
	*EthereumMultinodeClient
}

// KromaClient represents a single node, EVM compatible client for the Kroma network
type KromaClient struct {
	*EthereumClient
}
