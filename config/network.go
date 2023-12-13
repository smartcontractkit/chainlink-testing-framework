package config

import (
	"encoding/base64"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

type NetworkConfig struct {
	SelectedNetworks []string            `toml:"selected_networks"`
	RpcHttpUrls      map[string][]string `toml:"RpcHttpUrls"`
	WsRpcsUrls       map[string][]string `toml:"RpcWsUrls"`
	WalletKeys       map[string][]string `toml:"WalletKeys"`
}

func (n *NetworkConfig) ApplySecrets() error {
	encodedEndpoints, err := osutil.GetEnv("BASE64_NETWORK_CONFIG")
	if err != nil {
		return err
	}

	err = n.ApplyBase64Enconded(encodedEndpoints)
	if err != nil {
		return errors.Wrapf(err, "error reading network encoded endpoints")
	}

	return nil
}

func (n *NetworkConfig) ApplyDecoded(configDecoded string) error {
	if configDecoded == "" {
		return nil
	}

	var cfg NetworkConfig
	err := toml.Unmarshal([]byte(configDecoded), &cfg)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling network config")
	}

	err = n.ApplyOverrides(cfg)
	if err != nil {
		return errors.Wrapf(err, "error applying overrides from decoded network config file to config")
	}

	return nil
}

func (n *NetworkConfig) ApplyBase64Enconded(configEncoded string) error {
	if configEncoded == "" {
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(configEncoded)
	if err != nil {
		return err
	}

	return n.ApplyDecoded(string(decoded))
}

func (n *NetworkConfig) Validate() error {
	if len(n.SelectedNetworks) == 0 {
		return errors.New("selected_networks must be set")
	}

	for _, network := range n.SelectedNetworks {
		network = strings.ToUpper(network)
		if strings.Contains(network, "SIMULATED") {
			// we don't need to validate RPC endpoints or private keys for simulated networks
			continue
		}

		if _, ok := n.RpcHttpUrls[network]; !ok {
			return errors.Errorf("At least one HTTP RPC endpoint for %s network must be set", network)
		}

		if _, ok := n.WsRpcsUrls[network]; !ok {
			return errors.Errorf("At least one WS RPC endpoint for %s network must be set", network)
		}

		if _, ok := n.WalletKeys[network]; !ok {
			return errors.Errorf("At least one private key of funding wallet for %s network must be set", network)
		}
	}

	for i, network := range n.SelectedNetworks {
		n.SelectedNetworks[i] = strings.ToUpper(network)
	}

	return nil
}

func (n *NetworkConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case NetworkConfig:
		if asCfg.SelectedNetworks != nil {
			n.SelectedNetworks = asCfg.SelectedNetworks
		}
		if asCfg.RpcHttpUrls != nil {
			if n.RpcHttpUrls == nil || len(n.RpcHttpUrls) == 0 {
				n.RpcHttpUrls = asCfg.RpcHttpUrls
			} else {
				for network, urls := range asCfg.RpcHttpUrls {
					n.RpcHttpUrls[network] = urls
				}
			}
		}
		if asCfg.WsRpcsUrls != nil {
			if n.WsRpcsUrls == nil || len(n.WsRpcsUrls) == 0 {
				n.WsRpcsUrls = asCfg.WsRpcsUrls
			} else {
				for network, urls := range asCfg.WsRpcsUrls {
					n.WsRpcsUrls[network] = urls
				}
			}
		}
		if asCfg.WalletKeys != nil {
			if n.WalletKeys == nil || len(n.WalletKeys) == 0 {
				n.WalletKeys = asCfg.WalletKeys
			} else {
				for network, urls := range asCfg.WalletKeys {
					n.WalletKeys[network] = urls
				}
			}
		}

		return nil
	case *NetworkConfig:
		if asCfg.SelectedNetworks != nil {
			n.SelectedNetworks = asCfg.SelectedNetworks
		}
		if asCfg.RpcHttpUrls != nil {
			n.RpcHttpUrls = asCfg.RpcHttpUrls
		}
		if asCfg.WsRpcsUrls != nil {
			n.WsRpcsUrls = asCfg.WsRpcsUrls
		}
		if asCfg.WalletKeys != nil {
			n.WalletKeys = asCfg.WalletKeys
		}

		return nil
	default:
		return errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (n *NetworkConfig) Default() error {
	return nil
}
