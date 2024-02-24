package blockchain

// MantleMultinodeClient represents a multi-node, EVM compatible client for the Mantle network
type MantleMultinodeClient struct {
	*EthereumMultinodeClient
}

// MantleClient represents a single node, EVM compatible client for the Mantle network
type MantleClient struct {
	*EthereumClient
}
