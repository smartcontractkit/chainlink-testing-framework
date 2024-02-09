package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

const (
	Base64NetworkConfigEnvVarName = "BASE64_NETWORK_CONFIG"
)

type ForkConfig struct {
	URL              string `toml:"url"`          // URL is the URL of the node to fork from
	BlockNumber      int64  `toml:"block_number"` // BlockNumber is the block number to fork from
	BlockTime        int64  `toml:"block_time"`
	Retries          int    `toml:"retries"`
	Timeout          int64  `toml:"timeout"`
	ComputePerSecond int64  `toml:"compute_per_second"`
	RateLimitEnabled bool   `toml:"rate_limit_enabled"`
}

// NetworkConfig is the configuration for the networks to be used
type NetworkConfig struct {
	SelectedNetworks []string `toml:"selected_networks"`
	// EVMNetworks is the configuration for the EVM networks, key is the network name as declared in selected_networks slice.
	// if not set, it will try to find the network from defined networks in MappedNetworks under known_networks.go
	EVMNetworks map[string]*blockchain.EVMNetwork `toml:"evm_networks,omitempty"`
	// ForkConfigs is the configuration for forking from a node,
	// key is the network name as declared in selected_networks slice
	ForkConfigs map[string]ForkConfig `toml:"fork_config,omitempty"`
	// RpcHttpUrls is the RPC HTTP endpoints for each network,
	// key is the network name as declared in selected_networks slice
	RpcHttpUrls map[string][]string `toml:"RpcHttpUrls"`
	// RpcWsUrls is the RPC WS endpoints for each network,
	// key is the network name as declared in selected_networks slice
	RpcWsUrls map[string][]string `toml:"RpcWsUrls"`
	// WalletKeys is the private keys for the funding wallets for each network,
	// key is the network name as declared in selected_networks slice
	WalletKeys map[string][]string `toml:"WalletKeys"`
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

	return nil
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
		if n.ForkConfigs != nil {
			if _, ok := n.ForkConfigs[network]; ok {
				if evmConfig, exists := n.EVMNetworks[network]; !exists || evmConfig == nil {
					return fmt.Errorf("fork config for %s network is set, but no corresponding EVM network is defined", network)
				}
				if n.ForkConfigs[network].URL == "" {
					return fmt.Errorf("fork config for %s network must have a URL", network)
				}
				if n.ForkConfigs[network].BlockNumber == 0 {
					return fmt.Errorf("fork config for %s network must have a block number", network)
				}
				// we don't need to validate RPC endpoints or private keys for forked networks
				continue
			}
		}
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

	for i, network := range n.SelectedNetworks {
		n.SelectedNetworks[i] = strings.ToUpper(network)
		if _, ok := n.EVMNetworks[network]; ok {
			n.EVMNetworks[strings.ToUpper(network)] = n.EVMNetworks[network]
			delete(n.EVMNetworks, network)
		}
		if _, ok := n.ForkConfigs[network]; ok {
			n.ForkConfigs[strings.ToUpper(network)] = n.ForkConfigs[network]
			delete(n.ForkConfigs, network)
		}
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
	if defaults.ForkConfigs != nil {
		if n.ForkConfigs == nil || len(n.ForkConfigs) == 0 {
			n.ForkConfigs = defaults.ForkConfigs
		} else {
			for network, cfg := range defaults.ForkConfigs {
				if _, ok := n.ForkConfigs[network]; !ok {
					n.ForkConfigs[network] = cfg
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
