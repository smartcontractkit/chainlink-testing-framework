package config

import (
	_ "embed"
	"encoding/base64"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

//go:embed tomls/default.toml
var DefaultNetworkConfig []byte

type NetworkConfig struct {
	Network *struct {
		SelectedNetworks []string            `toml:"selected_networks"`
		RpcHttpUrls      map[string][]string `toml:"RpcHttpUrls"`
		WsRpcsUrls       map[string][]string `toml:"RpcWsUrls"`
		WalletKeys       map[string][]string `toml:"WalletKeys"`
	} `toml:"Network"`
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
	if len(n.Network.SelectedNetworks) == 0 {
		return errors.New("selected_networks must be set")
	}

	for _, network := range n.Network.SelectedNetworks {
		network = strings.ToUpper(network)
		if strings.Contains(network, "SIMULATED") {
			// we don't need to validate RPC endpoints or private keys for simulated networks
			continue
		}

		if _, ok := n.Network.RpcHttpUrls[network]; !ok {
			return errors.Errorf("At least one HTTP RPC endpoint for %s network must be set", network)
		}

		if _, ok := n.Network.WsRpcsUrls[network]; !ok {
			return errors.Errorf("At least one WS RPC endpoint for %s network must be set", network)
		}

		if _, ok := n.Network.WalletKeys[network]; !ok {
			return errors.Errorf("At least one private key of funding wallet for %s network must be set", network)
		}
	}

	for i, network := range n.Network.SelectedNetworks {
		n.Network.SelectedNetworks[i] = strings.ToUpper(network)
	}

	return nil
}

func (n *NetworkConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case NetworkConfig:
		if asCfg.Network.SelectedNetworks != nil {
			n.Network.SelectedNetworks = asCfg.Network.SelectedNetworks
		}
		if asCfg.Network.RpcHttpUrls != nil {
			if n.Network.RpcHttpUrls == nil || len(n.Network.RpcHttpUrls) == 0 {
				n.Network.RpcHttpUrls = asCfg.Network.RpcHttpUrls
			} else {
				for network, urls := range asCfg.Network.RpcHttpUrls {
					n.Network.RpcHttpUrls[network] = urls
				}
			}
		}
		if asCfg.Network.WsRpcsUrls != nil {
			if n.Network.WsRpcsUrls == nil || len(n.Network.WsRpcsUrls) == 0 {
				n.Network.WsRpcsUrls = asCfg.Network.WsRpcsUrls
			} else {
				for network, urls := range asCfg.Network.WsRpcsUrls {
					n.Network.WsRpcsUrls[network] = urls
				}
			}
		}
		if asCfg.Network.WalletKeys != nil {
			if n.Network.WalletKeys == nil || len(n.Network.WalletKeys) == 0 {
				n.Network.WalletKeys = asCfg.Network.WalletKeys
			} else {
				for network, urls := range asCfg.Network.WalletKeys {
					n.Network.WalletKeys[network] = urls
				}
			}
		}

		return nil
	case *NetworkConfig:
		if asCfg.Network.SelectedNetworks != nil {
			n.Network.SelectedNetworks = asCfg.Network.SelectedNetworks
		}
		if asCfg.Network.RpcHttpUrls != nil {
			n.Network.RpcHttpUrls = asCfg.Network.RpcHttpUrls
		}
		if asCfg.Network.WsRpcsUrls != nil {
			n.Network.WsRpcsUrls = asCfg.Network.WsRpcsUrls
		}
		if asCfg.Network.WalletKeys != nil {
			n.Network.WalletKeys = asCfg.Network.WalletKeys
		}

		return nil
	default:
		return errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (n *NetworkConfig) Default() error {
	if err := toml.Unmarshal(DefaultNetworkConfig, n); err != nil {
		return errors.Wrapf(err, "error unmarshaling config")
	}
	return nil
}
