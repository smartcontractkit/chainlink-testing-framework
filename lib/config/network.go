package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

type AnvilConfig struct {
	URL                *string                  `toml:"url,omitempty"`                   // Needed if you want to fork a network. URL is the URL of the node to fork from. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BlockNumber        *int64                   `toml:"block_number,omitempty"`          // Needed if fork URL is provided for forking. BlockNumber is the block number to fork from. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BlockTime          *int64                   `toml:"block_time,omitempty"`            // how frequent blocks are mined. By default, it automatically generates a new block as soon as a transaction is submitted. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BlockGaslimit      *int64                   `toml:"block_gaslimit,omitempty"`        //  BlockGaslimit is the gas limit for each block. Refer to https://book.getfoundry.sh/reference/anvil/#options
	CodeSize           *int64                   `toml:"code_size,omitempty"`             //  CodeSize is the size of the code in bytes. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BaseFee            *int64                   `toml:"base_fee,omitempty"`              //  BaseFee is the base fee for block. Refer to https://book.getfoundry.sh/reference/anvil/#options
	Retries            *int                     `toml:"retries,omitempty"`               //  Needed if fork URL is provided for forking. Number of retry requests for spurious networks (timed out requests). Refer to https://book.getfoundry.sh/reference/anvil/#options
	Timeout            *int64                   `toml:"timeout,omitempty"`               //  Needed if fork URL is provided for forking. Timeout in ms for requests sent to remote JSON-RPC server in forking mode. Refer to https://book.getfoundry.sh/reference/anvil/#options
	ComputePerSecond   *int64                   `toml:"compute_per_second,omitempty"`    // Needed if fork URL is provided for forking. Sets the number of assumed available compute units per second for this provider. Refer to https://book.getfoundry.sh/reference/anvil/#options
	RateLimitDisabled  *bool                    `toml:"rate_limit_disabled,omitempty"`   // Needed if fork URL is provided for forking. Rate limiting for this node’s provider. If set to true the node will start with --no-rate-limit Refer to https://book.getfoundry.sh/reference/anvil/#options
	NoOfAccounts       *int                     `toml:"no_of_accounts,omitempty"`        // Number of accounts to generate. Refer to https://book.getfoundry.sh/reference/anvil/#options
	EnableTracing      *bool                    `toml:"enable_tracing,omitempty"`        // Enable tracing for the node. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BlocksToKeepInMem  *int64                   `toml:"blocks_to_keep_in_mem,omitempty"` // Refer to --transaction-block-keeper option in https://book.getfoundry.sh/reference/anvil/#options
	GasSpikeSimulation GasSpikeSimulationConfig `toml:"GasSpikeSimulation,omitempty"`
	GasLimitSimulation GasLimitSimulationConfig `toml:"GasLimitSimulation,omitempty"`
}

// NetworkConfig is the configuration for the networks to be used
type NetworkConfig struct {
	SelectedNetworks []string `toml:"selected_networks,omitempty"`
	// EVMNetworks is the configuration for the EVM networks, key is the network name as declared in selected_networks slice.
	// if not set, it will try to find the network from defined networks in MappedNetworks under known_networks.go
	EVMNetworks map[string]*blockchain.EVMNetwork `toml:"EVMNetworks,omitempty"`
	// AnvilConfigs is the configuration for forking from a node,
	// key is the network name as declared in selected_networks slice
	AnvilConfigs map[string]*AnvilConfig `toml:"AnvilConfigs,omitempty"`
	// GethReorgConfig is the configuration for handling reorgs on Simulated Geth
	GethReorgConfig ReorgConfig `toml:"GethReorgConfig,omitempty"`
	// RpcHttpUrls is the RPC HTTP endpoints for each network,
	// key is the network name as declared in selected_networks slice
	RpcHttpUrls map[string][]string `toml:"RpcHttpUrls,omitempty"`
	// RpcWsUrls is the RPC WS endpoints for each network,
	// key is the network name as declared in selected_networks slice
	RpcWsUrls map[string][]string `toml:"RpcWsUrls,omitempty"`
	// WalletKeys is the private keys for the funding wallets for each network,
	// key is the network name as declared in selected_networks slice
	WalletKeys map[string][]string `toml:"WalletKeys,omitempty"`
}

type ReorgConfig struct {
	Enabled     bool                   `toml:"enabled,omitempty"`
	Depth       int                    `toml:"depth,omitempty"`
	DelayCreate blockchain.StrDuration `toml:"delay_create,omitempty"` // Delay before creating, expressed in Go duration format (e.g., "1m", "30s")
}

// GasSpikeSimulation is the configuration for simulating gas spikes on the network
type GasSpikeSimulationConfig struct {
	Enabled           bool                   `toml:"enabled,omitempty"`
	StartGasPrice     int64                  `toml:"start_gas_price,omitempty"`
	GasRisePercentage float64                `toml:"gas_rise_percentage,omitempty"`
	GasSpike          bool                   `toml:"gas_spike,omitempty"`
	DelayCreate       blockchain.StrDuration `toml:"delay_create,omitempty"` // Delay before creating, expressed in Go duration format (e.g., "1m", "30s")
	Duration          blockchain.StrDuration `toml:"duration,omitempty"`     // Duration of the gas simulation, expressed in Go duration format (e.g., "1m", "30s")
}

// GasLimitSimulationConfig is the configuration for simulating gas limit changes on the network
type GasLimitSimulationConfig struct {
	Enabled                bool                   `toml:"enabled,omitempty"`
	NextGasLimitPercentage float64                `toml:"next_gas_limit_percentage,omitempty"` // Percentage of last gasUsed in previous block creating congestion
	DelayCreate            blockchain.StrDuration `toml:"delay_create,omitempty"`              // Delay before creating, expressed in Go duration format (e.g., "1m", "30s")
	Duration               blockchain.StrDuration `toml:"duration,omitempty"`                  // Duration of the gas simulation, expressed in Go duration format (e.g., "1m", "30s")
}

func (n NetworkConfig) IsSimulatedGethSelected() bool {
	for _, network := range n.SelectedNetworks {
		if strings.ToLower(network) == "simulated" {
			return true
		}
	}
	return false
}

// OverrideURLsAndKeysFromEVMNetwork applies the URLs and keys from the EVMNetworks to the NetworkConfig
// it overrides the URLs and Keys present in RpcHttpUrls, RpcWsUrls and WalletKeys in the NetworkConfig
// with the URLs and Keys provided in the EVMNetworks
func (n *NetworkConfig) OverrideURLsAndKeysFromEVMNetwork() {
	if n.EVMNetworks == nil {
		return
	}
	for name, evmNetwork := range n.EVMNetworks {
		if len(evmNetwork.URLs) > 0 {
			logging.L.Warn().Str("network", name).Msg("found URLs in EVMNetwork. overriding RPC URLs in RpcWsUrls with EVMNetwork URLs")
			if n.RpcWsUrls == nil {
				n.RpcWsUrls = make(map[string][]string)
			}
			n.RpcWsUrls[name] = evmNetwork.URLs
		}
		if len(evmNetwork.HTTPURLs) > 0 {
			logging.L.Warn().Str("network", name).Msg("found HTTPURLs in EVMNetwork. overriding RPC URLs in RpcHttpUrls with EVMNetwork HTTP URLs")
			if n.RpcHttpUrls == nil {
				n.RpcHttpUrls = make(map[string][]string)
			}
			n.RpcHttpUrls[name] = evmNetwork.HTTPURLs
		}
		if len(evmNetwork.PrivateKeys) > 0 {
			logging.L.Warn().Str("network", name).Msg("found PrivateKeys in EVMNetwork. overriding wallet keys in WalletKeys with EVMNetwork private keys")
			if n.WalletKeys == nil {
				n.WalletKeys = make(map[string][]string)
			}
			n.WalletKeys[name] = evmNetwork.PrivateKeys
		}
	}
}

// Validate checks if all required fields are set, meaning that there must be at least
// 1 selected network and unless it's a simulated network, there must be at least 1
// rpc endpoint for HTTP and WS and 1 private key for funding wallet
func (n *NetworkConfig) Validate() error {
	if len(n.SelectedNetworks) == 0 {
		return errors.New("selected_networks must be set")
	}

	for _, network := range n.SelectedNetworks {
		if strings.Contains(network, "SIMULATED") {
			// we don't need to validate RPC endpoints or private keys for simulated networks
			continue
		}
		for name, evmNetwork := range n.EVMNetworks {
			if evmNetwork.ClientImplementation == "" {
				return fmt.Errorf("client implementation for %s network must be set", name)
			}
			if evmNetwork.ChainID == 0 {
				return fmt.Errorf("chain ID for %s network must be set", name)
			}
		}
		if n.AnvilConfigs != nil {
			if _, ok := n.AnvilConfigs[network]; ok {
				// we don't need to validate RPC endpoints or private keys for forked networks
				continue
			}
		}

		// If the network is not forked, we need to validate RPC endpoints and private keys
		_, httpOk := n.RpcHttpUrls[network]
		_, wsOk := n.RpcWsUrls[network]
		// WS can be present but only if HTTP is also available
		if wsOk && !httpOk {
			return fmt.Errorf("WS RPC endpoint for %s network is set without an HTTP endpoint; only HTTP or both HTTP and WS are allowed", network)
		}

		// Validate that there is at least one HTTP endpoint
		if !httpOk {
			return fmt.Errorf("at least one HTTP RPC endpoint for %s network must be set", network)
		}

		if _, ok := n.WalletKeys[network]; !ok {
			return fmt.Errorf("at least one private key of funding wallet for %s network must be set", network)
		}
	}

	return nil
}

// UpperCaseNetworkNames converts all network name keys for wallet keys, rpc endpoints maps and
// selected network slice to upper case
func (n *NetworkConfig) UpperCaseNetworkNames() {
	var upperCaseMapKeys = func(m map[string][]string) {
		newMap := make(map[string][]string)
		for key := range m {
			newMap[strings.ToUpper(key)] = m[key]
			delete(m, key)
		}
		for key := range newMap {
			m[key] = newMap[key]
		}
	}

	upperCaseMapKeys(n.RpcHttpUrls)
	upperCaseMapKeys(n.RpcWsUrls)
	upperCaseMapKeys(n.WalletKeys)

	for network := range n.EVMNetworks {
		if network != strings.ToUpper(network) {
			n.EVMNetworks[strings.ToUpper(network)] = n.EVMNetworks[network]
			delete(n.EVMNetworks, network)
		}
	}

	for network := range n.AnvilConfigs {
		if network != strings.ToUpper(network) {
			n.AnvilConfigs[strings.ToUpper(network)] = n.AnvilConfigs[network]
			delete(n.AnvilConfigs, network)
		}
	}

	for i, network := range n.SelectedNetworks {
		n.SelectedNetworks[i] = strings.ToUpper(network)
	}
}
