package blockchain

// BttcMultinodeClient represents a multi-node, EVM compatible client for the Kava network
type MantleSepoliaMultinodeClient struct {
	*EthereumMultinodeClient
}

// BttcClient represents a single node, EVM compatible client for the Kava network
type MantleSepoliaClient struct {
	*EthereumClient
}
