package blockchain

// WeMixMultinodeClient represents a multi-node, EVM compatible client for the WeMix network
type WeMixMultinodeClient struct {
	*EthereumMultinodeClient
}

// WeMixClient represents a single node, EVM compatible client for the WeMix network
type WeMixClient struct {
	*EthereumClient
}
