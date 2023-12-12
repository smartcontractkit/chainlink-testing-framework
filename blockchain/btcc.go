package blockchain

// BTCCMultinodeClient represents a multi-node, EVM compatible client for the BTCC network
type BTCCMultinodeClient struct {
	*EthereumMultinodeClient
}

// BTCCClient represents a single node, EVM compatible client for the BTCC network
type BTCCClient struct {
	*EthereumClient
}
