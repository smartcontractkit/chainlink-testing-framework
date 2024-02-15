package blockchain

// BttcMultinodeClient represents a multi-node, EVM compatible client for the Kava network
type MantleGoerliMultinodeClient struct {
	*EthereumMultinodeClient
}

// BttcClient represents a single node, EVM compatible client for the Kava network
type MantleGoerliClient struct {
	*EthereumClient
}
