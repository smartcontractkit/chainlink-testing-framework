package blockchain

// LineaMultinodeClient represents a multi-node, EVM compatible client for the Linea network
type LineaMultinodeClient struct {
	*EthereumMultinodeClient
}

// LineaClient represents a single node, EVM compatible client for the Linea network
type LineaClient struct {
	*EthereumClient
}
