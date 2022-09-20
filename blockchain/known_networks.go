package blockchain

import (
	"github.com/rs/zerolog/log"
)

// knownNetworks lets us toggle which blockchain client to utilize depending on the supplied Chain ID for the network.
// Use resources like https://chainlist.org/ or official docs on each chain to distinguish each one.
var knownNetworks = map[int64]string{
	1:        "Ethereum", // Mainnet
	5:        "Ethereum", // Goerli
	1337:     "Ethereum", // Geth Dev
	2337:     "Ethereum", // Geth Dev #2 used in reorg tests
	11155111: "Ethereum", // Sepolia

	588:  "Metis", // Stardust (testnet)
	1088: "Metis", // Andoromeda (mainnet)

	1001: "Klaytn", // Testnet
	8217: "Klaytn", // Mainnet

	80001: "Mumbai",
	100:   "edge",

	421611: "Arbitrum", // Rinkeby
	421613: "Arbitrum", // Goerli
}

// wrapSingleClient Wraps a single EVM client in its appropriate implementation, based on the chain ID
func wrapSingleClient(networkSettings *EVMNetwork, client *EthereumClient) EVMClient {
	networkType, known := knownNetworks[networkSettings.ChainID]
	if !known {
		log.Warn().
			Str("Network", networkSettings.Name).
			Int64("Based on Chain ID", networkSettings.ChainID).
			Msg("Unrecognized Chain ID. Defaulting to a Standard Ethereum Client.")
		return client
	}

	var wrappedEc EVMClient
	switch networkType {
	case "Ethereum":
		wrappedEc = client
	case "Metis":
		wrappedEc = &MetisClient{client}
	case "edge":
		wrappedEc = &PolygonEdgeClient{client}
	case "Klaytn":
		wrappedEc = &KlaytnClient{client}
	case "Arbitrum":
		wrappedEc = &ArbitrumClient{client}
	default:
		wrappedEc = client
	}
	return wrappedEc
}

// wrapMultiClient Wraps a multi-node EVM client in its appropriate implementation, based on the chain ID
func wrapMultiClient(networkSettings *EVMNetwork, client *EthereumMultinodeClient) EVMClient {
	networkType, known := knownNetworks[networkSettings.ChainID]
	if !known {
		log.Warn().
			Str("Network", networkSettings.Name).
			Interface("URLs", networkSettings.URLs).
			Int64("Based on Chain ID", networkSettings.ChainID).
			Msg("Unrecognized Chain ID. Defaulting to a Standard Ethereum Client.")
		return client
	}

	var wrappedEc EVMClient
	logMsg := log.Info().Str("Network", networkSettings.Name).Int64("Based on Chain ID", networkSettings.ChainID)
	switch networkType {
	case "Ethereum":
		logMsg.Msg("Using Standard Ethereum Client")
		wrappedEc = client
	case "edge":
		logMsg.Msg("Using Polygon edge client")
		wrappedEc = &PolygonEdgeMultinodeClient{client}
	case "Metis":
		logMsg.Msg("Using Metis Client")
		wrappedEc = &MetisMultinodeClient{client}
	case "Klaytn":
		logMsg.Msg("Using Klaytn Client")
		wrappedEc = &KlaytnMultinodeClient{client}
	case "Arbitrum":
		logMsg.Msg("Using Arbitrum Client")
		wrappedEc = &ArbitrumMultinodeClient{client}
	default:
		logMsg.Msg("Using Standard Ethereum Client")
		wrappedEc = client
	}
	return wrappedEc
}
