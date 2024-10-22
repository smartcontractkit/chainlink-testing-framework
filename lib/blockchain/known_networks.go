package blockchain

import (
	"github.com/rs/zerolog/log"
)

// ClientImplementation represents the type of EVM client implementation for the framework to use
type ClientImplementation string

const (
	// Ethereum uses the standard EVM implementation, and is considered default
	EthereumClientImplementation     ClientImplementation = "Ethereum"
	MetisClientImplementation        ClientImplementation = "Metis"
	KlaytnClientImplementation       ClientImplementation = "Klaytn"
	OptimismClientImplementation     ClientImplementation = "Optimism"
	ArbitrumClientImplementation     ClientImplementation = "Arbitrum"
	PolygonClientImplementation      ClientImplementation = "Polygon"
	RSKClientImplementation          ClientImplementation = "RSK"
	CeloClientImplementation         ClientImplementation = "Celo"
	QuorumClientImplementation       ClientImplementation = "Quorum"
	ScrollClientImplementation       ClientImplementation = "Scroll"
	BSCClientImplementation          ClientImplementation = "BSC"
	LineaClientImplementation        ClientImplementation = "Linea"
	PolygonZkEvmClientImplementation ClientImplementation = "PolygonZkEvm"
	FantomClientImplementation       ClientImplementation = "Fantom"
	WeMixClientImplementation        ClientImplementation = "WeMix"
	KromaClientImplementation        ClientImplementation = "Kroma"
	GnosisClientImplementation       ClientImplementation = "Gnosis"
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
	case ScrollClientImplementation:
		wrappedEc = &ScrollClient{client}
	case BSCClientImplementation:
		wrappedEc = &BSCClient{client}
	case LineaClientImplementation:
		wrappedEc = &LineaClient{client}
	case PolygonZkEvmClientImplementation:
		wrappedEc = &PolygonZkEvmClient{client}
	case FantomClientImplementation:
		wrappedEc = &FantomClient{client}
	case WeMixClientImplementation:
		wrappedEc = &WeMixClient{client}
	case KromaClientImplementation:
		wrappedEc = &KromaClient{client}
	case GnosisClientImplementation:
		wrappedEc = &GnosisClient{client}
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
	case ScrollClientImplementation:
		logMsg.Msg("Using Scroll Client")
		wrappedEc = &ScrollMultinodeClient{client}
	case BSCClientImplementation:
		logMsg.Msg("Using BSC Client")
		wrappedEc = &BSCMultinodeClient{client}
	case LineaClientImplementation:
		logMsg.Msg("Using Linea Client")
		wrappedEc = &LineaMultinodeClient{client}
	case PolygonZkEvmClientImplementation:
		logMsg.Msg("Using Polygon zkEVM Client")
		wrappedEc = &PolygonZkEvmMultinodeClient{client}
	case FantomClientImplementation:
		logMsg.Msg("Using Fantom Client")
		wrappedEc = &FantomMultinodeClient{client}
	case WeMixClientImplementation:
		logMsg.Msg("Using WeMix Client")
		wrappedEc = &WeMixMultinodeClient{client}
	case KromaClientImplementation:
		logMsg.Msg("Using Kroma Client")
		wrappedEc = &KromaMultinodeClient{client}
	case GnosisClientImplementation:
		logMsg.Msg("Using Gnosis Client")
		wrappedEc = &GnosisMultinodeClient{client}
	default:
		log.Warn().Str("Network", networkSettings.Name).Msg("Unknown client implementation, defaulting to standard Ethereum client")
		wrappedEc = client
	}
	return wrappedEc
}
