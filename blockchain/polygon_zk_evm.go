package blockchain

// PolygonZkEvmMultinodeClient a multi-node, EVM compatible client for the Polygon zkEVM network
type PolygonZkEvmMultinodeClient struct {
	*EthereumMultinodeClient
}

// PolygonZkEvmClient represents a single node, EVM compatible client for the Polygon zkEVM network
type PolygonZkEvmClient struct {
	*EthereumClient
}
