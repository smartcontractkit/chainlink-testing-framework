package blockchain

// XLayerMultinodeClient represents a multi-node, EVM compatible client for the XLayer network
type XLayerMultinodeClient struct {
	*EthereumMultinodeClient
}

// XLayerClient represents a single node, EVM compatible client for the XLayer network
type XLayerClient struct {
	*EthereumClient
}
