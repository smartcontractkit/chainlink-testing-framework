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

	err = n.ApplyOverrides(&cfg)
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

func (n *NetworkConfig) ApplyOverrides(from *NetworkConfig) error {
	if from != nil {
		return nil
	}
	if from.SelectedNetworks != nil {
		n.SelectedNetworks = from.SelectedNetworks
	}
	if from.RpcHttpUrls != nil {
		if n.RpcHttpUrls == nil || len(n.RpcHttpUrls) == 0 {
			n.RpcHttpUrls = from.RpcHttpUrls
		} else {
			for network, urls := range from.RpcHttpUrls {
				n.RpcHttpUrls[network] = urls
			}
		}
	}
	if from.WsRpcsUrls != nil {
		if n.WsRpcsUrls == nil || len(n.WsRpcsUrls) == 0 {
			n.WsRpcsUrls = from.WsRpcsUrls
		} else {
			for network, urls := range from.WsRpcsUrls {
				n.WsRpcsUrls[network] = urls
			}
		}
	}
	if from.WalletKeys != nil {
		if n.WalletKeys == nil || len(n.WalletKeys) == 0 {
			n.WalletKeys = from.WalletKeys
		} else {
			for network, urls := range from.WalletKeys {
				n.WalletKeys[network] = urls
			}
		}
	}

	return nil
}

func (n *NetworkConfig) Default() error {
	return nil
}
