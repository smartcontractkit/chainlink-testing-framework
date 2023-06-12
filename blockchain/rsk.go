package blockchain

// Handles specific issues with the RSK EVM chain: https://developers.rsk.co/rsk/node/architecture/json-rpc/

// RSKMultinodeClient represents a multi-node, EVM compatible client for the RSK network
type RSKMultinodeClient struct {
	*EthereumMultinodeClient
}

// RSKClient represents a single node, EVM compatible client for the RSK network
type RSKClient struct {
	*EthereumClient
}
