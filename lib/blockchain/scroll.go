package blockchain

// ScrollMultinodeClient represents a multi-node, EVM compatible client for the Scroll network
type ScrollMultinodeClient struct {
	*EthereumMultinodeClient
}

// ScrollClient represents a single node, EVM compatible client for the Scroll network
type ScrollClient struct {
	*EthereumClient
}
