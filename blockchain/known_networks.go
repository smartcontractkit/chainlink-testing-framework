package blockchain

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
)

// KnownEVMNetworks holds all previously known EVM networks for easy connections
type KnownEVMNetworks map[int64]*RegisteredEVMNetwork

// GetNetwork retrieves an EVM network based on the chainId
func (k *KnownEVMNetworks) GetNetwork(chainId int64) *RegisteredEVMNetwork {
	network, known := (*k)[chainId]
	if !known {
		log.Warn().Int64("Chain ID", chainId).Msg("Unknown EVM chain ID, using defaults")
		return &RegisteredEVMNetwork{
			Name:              "Generic EVM Testnet",
			ConnectClientFunc: NewEthereumMultiNodeClient,
		}
	}
	return network
}

// RegisteredEVMNetwork represents very basic information about a known evm network
type RegisteredEVMNetwork struct {
	Name              string
	ChainAttributes   *client.EVMChainAttributes
	ConnectClientFunc NewEVMClientFn
}

// RegisteredEVMNetworks holds all known and registered EVM networks.
// ChainID : Network details
var RegisteredEVMNetworks KnownEVMNetworks = map[int64]*RegisteredEVMNetwork{
	// Simulated
	1337: {
		Name:              "Simulated Geth Network",
		ConnectClientFunc: NewEthereumMultiNodeClient,
	},

	// Testnets
	42: {
		Name:              "Kovan Testnet",
		ConnectClientFunc: NewEthereumMultiNodeClient,
	},
	69: {
		Name:              "Optimism Testnet",
		ConnectClientFunc: NewEthereumMultiNodeClient,
	},
	588: {
		Name:              "Metis Testnet",
		ConnectClientFunc: NewMetisMultiNodeClient,
	},
	1001: {
		Name:              "Klaytn Testnet",
		ConnectClientFunc: NewKlaytnMultiNodeClient,
	},

	// Mainnets
	1: {
		Name:              "Ethereum Mainnet",
		ConnectClientFunc: NewEthereumMultiNodeClient,
	},
	10: {
		Name:              "Optimism Mainnet",
		ConnectClientFunc: NewEthereumMultiNodeClient,
	},
	1088: {
		Name:              "Metis Mainnet",
		ConnectClientFunc: NewMetisMultiNodeClient,
	},
	8217: {
		Name:              "Klaytn Mainnet",
		ConnectClientFunc: NewKlaytnMultiNodeClient,
	},
}
