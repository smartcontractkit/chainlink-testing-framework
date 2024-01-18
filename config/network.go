package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/pelletier/go-toml/v2"
)

const (
	Base64NetworkConfigEnvVarName = "BASE64_NETWORK_CONFIG"
)

type NetworkConfig struct {
	SelectedNetworks []string            `toml:"selected_networks"`
	RpcHttpUrls      map[string][]string `toml:"RpcHttpUrls"`
	RpcWsUrls        map[string][]string `toml:"RpcWsUrls"`
	WalletKeys       map[string][]string `toml:"WalletKeys"`
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
	}
}

func (n *NetworkConfig) applyDefaults(defaults *NetworkConfig) error {
	if defaults == nil {
		return nil
	}

	if defaults.SelectedNetworks != nil {
		n.SelectedNetworks = defaults.SelectedNetworks
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
