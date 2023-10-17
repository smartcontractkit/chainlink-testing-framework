package blockchain

// FantomMultinodeClient represents a multi-node, EVM compatible client for the Fantom network
type FantomMultinodeClient struct {
	*EthereumMultinodeClient
}

// FantomClient represents a single node, EVM compatible client for the Fantom network
type FantomClient struct {
	*EthereumClient
}
