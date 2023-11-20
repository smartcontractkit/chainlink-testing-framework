package blockchain

// KavaMultinodeClient represents a multi-node, EVM compatible client for the Kava network
type KavaMultinodeClient struct {
	*EthereumMultinodeClient
}

// KavaClient represents a single node, EVM compatible client for the Kava network
type KavaClient struct {
	*EthereumClient
}
