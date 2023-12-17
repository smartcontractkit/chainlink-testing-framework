package blockchain

// BttcMultinodeClient represents a multi-node, EVM compatible client for the Kava network
type BttcMultinodeClient struct {
	*EthereumMultinodeClient
}

// BttcClient represents a single node, EVM compatible client for the Kava network
type BttcClient struct {
	*EthereumClient
}
