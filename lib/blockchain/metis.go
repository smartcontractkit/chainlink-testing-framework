package blockchain

// Handles specific issues with the Metis EVM chain: https://docs.metis.io/

// MetisMultinodeClient represents a multi-node, EVM compatible client for the Metis network
type MetisMultinodeClient struct {
	*EthereumMultinodeClient
}

// MetisClient represents a single node, EVM compatible client for the Metis network
type MetisClient struct {
	*EthereumClient
}
