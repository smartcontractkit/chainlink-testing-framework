package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

const (
	Base64NetworkConfigEnvVarName = "BASE64_NETWORK_CONFIG"
)

type AnvilConfig struct {
	URL               *string `toml:"url,omitempty"`                 // Needed if you want to fork a network. URL is the URL of the node to fork from. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BlockNumber       *int64  `toml:"block_number,omitempty"`        // Needed if fork URL is provided for forking. BlockNumber is the block number to fork from. Refer to https://book.getfoundry.sh/reference/anvil/#options
	BlockTime         *int64  `toml:"block_time,omitempty"`          // how frequent blocks are mined. By default, it automatically generates a new block as soon as a transaction is submitted. Refer to https://book.getfoundry.sh/reference/anvil/#options
	Retries           *int    `toml:"retries,omitempty"`             //  Needed if fork URL is provided for forking. Number of retry requests for spurious networks (timed out requests). Refer to https://book.getfoundry.sh/reference/anvil/#options
	Timeout           *int64  `toml:"timeout,omitempty"`             //  Needed if fork URL is provided for forking. Timeout in ms for requests sent to remote JSON-RPC server in forking mode. Refer to https://book.getfoundry.sh/reference/anvil/#options
	ComputePerSecond  *int64  `toml:"compute_per_second,omitempty"`  // Needed if fork URL is provided for forking. Sets the number of assumed available compute units per second for this provider. Refer to https://book.getfoundry.sh/reference/anvil/#options
	RateLimitDisabled *bool   `toml:"rate_limit_disabled,omitempty"` // Needed if fork URL is provided for forking. Rate limiting for this node’s provider. If set to true the node will start with --no-rate-limit Refer to https://book.getfoundry.sh/reference/anvil/#options
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

func (n *NetworkConfig) applySecrets() error {
	encodedEndpoints, isSet := os.LookupEnv(Base64NetworkConfigEnvVarName)
	if !isSet {
		return nil
	}

	err := n.applyBase64Enconded(encodedEndpoints)
	if err != nil {
		return fmt.Errorf("error reading network encoded endpoints: %w", err)
	}

	return nil
}

func (n *NetworkConfig) applyDecoded(configDecoded string) error {
	if configDecoded == "" {
		return nil
	}

	var cfg NetworkConfig
	err := toml.Unmarshal([]byte(configDecoded), &cfg)
	if err != nil {
		return fmt.Errorf("error unmarshaling network config: %w", err)
	}

	cfg.UpperCaseNetworkNames()

	err = n.applyDefaults(&cfg)
	if err != nil {
		return fmt.Errorf("error applying overrides from decoded network config file to config: %w", err)
	}
	n.OverrideURLsAndKeysFromEVMNetwork()

	return nil
}

// OverrideURLsAndKeysFromEVMNetwork applies the URLs and keys from the EVMNetworks to the NetworkConfig
// it overrides the URLs and Keys present in RpcHttpUrls, RpcWsUrls and WalletKeys in the NetworkConfig
// with the URLs and Keys provided in the EVMNetworks
func (n *NetworkConfig) OverrideURLsAndKeysFromEVMNetwork() {
	if n.EVMNetworks == nil {
		return
	}
	for name, evmNetwork := range n.EVMNetworks {
		if evmNetwork.URLs != nil && len(evmNetwork.URLs) > 0 {
			logging.L.Warn().Str("network", name).Msg("found URLs in EVMNetwork. overriding RPC URLs in RpcWsUrls with EVMNetwork URLs")
			if n.RpcWsUrls == nil {
				n.RpcWsUrls = make(map[string][]string)
			}
			n.RpcWsUrls[name] = evmNetwork.URLs
		}
		if evmNetwork.HTTPURLs != nil && len(evmNetwork.HTTPURLs) > 0 {
			logging.L.Warn().Str("network", name).Msg("found HTTPURLs in EVMNetwork. overriding RPC URLs in RpcHttpUrls with EVMNetwork HTTP URLs")
			if n.RpcHttpUrls == nil {
				n.RpcHttpUrls = make(map[string][]string)
			}
			n.RpcHttpUrls[name] = evmNetwork.HTTPURLs
		}
		if evmNetwork.PrivateKeys != nil && len(evmNetwork.PrivateKeys) > 0 {
			logging.L.Warn().Str("network", name).Msg("found PrivateKeys in EVMNetwork. overriding wallet keys in WalletKeys with EVMNetwork private keys")
			if n.WalletKeys == nil {
				n.WalletKeys = make(map[string][]string)
			}
			n.WalletKeys[name] = evmNetwork.PrivateKeys
		}
	}
}

func (n *NetworkConfig) applyBase64Enconded(configEncoded string) error {
	if configEncoded == "" {
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(configEncoded)
	if err != nil {
		return err
	}

	return n.applyDecoded(string(decoded))
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
		// if the network is not forked, we need to validate RPC endpoints and private keys
		if _, ok := n.RpcHttpUrls[network]; !ok {
			return fmt.Errorf("at least one HTTP RPC endpoint for %s network must be set", network)
		}

		if _, ok := n.RpcWsUrls[network]; !ok {
			return fmt.Errorf("at least one WS RPC endpoint for %s network must be set", network)
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

func (n *NetworkConfig) applyDefaults(defaults *NetworkConfig) error {
	if defaults == nil {
		return nil
	}

	if defaults.SelectedNetworks != nil {
		n.SelectedNetworks = defaults.SelectedNetworks
	}
	if defaults.EVMNetworks != nil {
		if n.EVMNetworks == nil || len(n.EVMNetworks) == 0 {
			n.EVMNetworks = defaults.EVMNetworks
		} else {
			for network, cfg := range defaults.EVMNetworks {
				if _, ok := n.EVMNetworks[network]; !ok {
					n.EVMNetworks[network] = cfg
				}
			}
		}
	}
	if defaults.AnvilConfigs != nil {
		if n.AnvilConfigs == nil || len(n.AnvilConfigs) == 0 {
			n.AnvilConfigs = defaults.AnvilConfigs
		} else {
			for network, cfg := range defaults.AnvilConfigs {
				if _, ok := n.AnvilConfigs[network]; !ok {
					n.AnvilConfigs[network] = cfg
				}
			}
		}
	}
	if defaults.RpcHttpUrls != nil {
		if n.RpcHttpUrls == nil || len(n.RpcHttpUrls) == 0 {
			n.RpcHttpUrls = defaults.RpcHttpUrls
		} else {
			for network, urls := range defaults.RpcHttpUrls {
				if _, ok := n.RpcHttpUrls[network]; !ok {
					n.RpcHttpUrls[network] = urls
				}
			}
		}
	}
	if defaults.RpcWsUrls != nil {
		if n.RpcWsUrls == nil || len(n.RpcWsUrls) == 0 {
			n.RpcWsUrls = defaults.RpcWsUrls
		} else {
			for network, urls := range defaults.RpcWsUrls {
				if _, ok := n.RpcWsUrls[network]; !ok {
					n.RpcWsUrls[network] = urls
				}
			}
		}
	}
	if defaults.WalletKeys != nil {
		if n.WalletKeys == nil || len(n.WalletKeys) == 0 {
			n.WalletKeys = defaults.WalletKeys
		} else {

			for network, keys := range defaults.WalletKeys {
				if _, ok := n.WalletKeys[network]; !ok {
					n.WalletKeys[network] = keys
				}
			}
		}
	}

	return nil
}

// Default applies default values to the network config after reading it
// from BASE64_NETWORK_CONFIG env var. It will only fill in the gaps, not override
// meaning that if you provided WS RPC endpoint in your network config, but not the
// HTTP one, then only HTTP will be taken from default config (provided it's there)
func (n *NetworkConfig) Default() error {
	return n.applySecrets()
}
