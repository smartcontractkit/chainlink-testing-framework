package blockchain

// GnosisMultinodeClient represents a multi-node, EVM compatible client for the Gnosis network
type GnosisMultinodeClient struct {
	*EthereumMultinodeClient
}

// GnosisClient represents a single node, EVM compatible client for the Gnosis network
type GnosisClient struct {
	*EthereumClient
}
