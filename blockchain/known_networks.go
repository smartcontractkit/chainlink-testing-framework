package blockchain

import (
	"github.com/rs/zerolog/log"
)

// ClientImplementation represents the type of EVM client implementation for the framework to use
type ClientImplementation string

const (
	// Ethereum uses the standard EVM implementation, and is considered default
	EthereumClientImplementation ClientImplementation = "Ethereum"
	MetisClientImplementation    ClientImplementation = "Metis"
	KlaytnClientImplementation   ClientImplementation = "Klaytn"
	OptimismClientImplementation ClientImplementation = "Optimism"
	ArbitrumClientImplementation ClientImplementation = "Arbitrum"
	PolygonClientImplementation  ClientImplementation = "Polygon"
	RSKClientImplementation      ClientImplementation = "RSK"
	CeloClientImplementation     ClientImplementation = "Celo"
	QuorumClientImplementation   ClientImplementation = "Quorum"
)

// wrapSingleClient Wraps a single EVM client in its appropriate implementation, based on the chain ID
func wrapSingleClient(networkSettings EVMNetwork, client *EthereumClient) EVMClient {
	var wrappedEc EVMClient
	switch networkSettings.ClientImplementation {
	case EthereumClientImplementation:
		wrappedEc = client
	case MetisClientImplementation:
		wrappedEc = &MetisClient{client}
	case PolygonClientImplementation:
		wrappedEc = &PolygonClient{client}
	case KlaytnClientImplementation:
		wrappedEc = &KlaytnClient{client}
	case ArbitrumClientImplementation:
		wrappedEc = &ArbitrumClient{client}
	case OptimismClientImplementation:
		wrappedEc = &OptimismClient{client}
	case RSKClientImplementation:
		wrappedEc = &RSKClient{client}
	case CeloClientImplementation:
		wrappedEc = &CeloClient{client}
	case QuorumClientImplementation:
		wrappedEc = &QuorumClient{client}
	default:
		wrappedEc = client
	}
	return wrappedEc
}

// wrapMultiClient Wraps a multi-node EVM client in its appropriate implementation, based on the chain ID
func wrapMultiClient(networkSettings EVMNetwork, client *EthereumMultinodeClient) EVMClient {
	var wrappedEc EVMClient
	logMsg := log.Info().Str("Network", networkSettings.Name)
	switch networkSettings.ClientImplementation {
	case EthereumClientImplementation:
		logMsg.Msg("Using Standard Ethereum Client")
		wrappedEc = client
	case PolygonClientImplementation:
		logMsg.Msg("Using Polygon Client")
		wrappedEc = &PolygonMultinodeClient{client}
	case MetisClientImplementation:
		logMsg.Msg("Using Metis Client")
		wrappedEc = &MetisMultinodeClient{client}
	case KlaytnClientImplementation:
		logMsg.Msg("Using Klaytn Client")
		wrappedEc = &KlaytnMultinodeClient{client}
	case ArbitrumClientImplementation:
		logMsg.Msg("Using Arbitrum Client")
		wrappedEc = &ArbitrumMultinodeClient{client}
	case OptimismClientImplementation:
		logMsg.Msg("Using Optimism Client")
		wrappedEc = &OptimismMultinodeClient{client}
	case RSKClientImplementation:
		logMsg.Msg("Using RSK Client")
		wrappedEc = &RSKMultinodeClient{client}
	case CeloClientImplementation:
		logMsg.Msg("Using Celo Client")
		wrappedEc = &CeloMultinodeClient{client}
	case QuorumClientImplementation:
		logMsg.Msg("Using Quorum Client")
		wrappedEc = &QuorumMultinodeClient{client}
	default:
		log.Warn().Str("Network", networkSettings.Name).Msg("Unknown client implementation, defaulting to standard Ethereum client")
		wrappedEc = client
	}
	return wrappedEc
}
