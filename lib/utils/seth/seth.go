package seth

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	pkg_seth "github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
)

var ErrInsufficientEphemeralKeys = `
Error: Insufficient Ephemeral Addresses for Simulated Network

To operate on a simulated network, you must configure at least one ephemeral address. Currently, %d ephemeral address(es) are set. Please update your TOML configuration file as follows to meet this requirement:
[Seth] ephemeral_addresses_number = 1

This adjustment ensures that your setup is minimally viable. Although it is highly recommended to use at least 20 ephemeral addresses.
`

var ErrInsufficientStaticKeys = `
Error: Insufficient Private Keys for Live Network

To run this test on a live network, you must either:
1. Set at least two private keys in the '[Network.WalletKeys]' section of your TOML configuration file. Example format:
   [Network.WalletKeys]
   NETWORK_NAME=["PRIVATE_KEY_1", "PRIVATE_KEY_2"]
2. Set at least two private keys in the '[Network.EVMNetworks.NETWORK_NAME] section of your TOML configuration file. Example format:
   evm_keys=["PRIVATE_KEY_1", "PRIVATE_KEY_2"]

Currently, only %d private key/s is/are set.

Recommended Action:
Distribute your funds across multiple private keys and update your configuration accordingly. Even though 1 private key is sufficient for testing, it is highly recommended to use at least 10 private keys.
`

var noOpSethConfigFn = func(cfg *pkg_seth.Config) error { return nil }

type ConfigFunction = func(*pkg_seth.Config) error

// OneEphemeralKeysLiveTestnetCheckFn checks whether there's at least one ephemeral key on a simulated network or at least one static key on a live network,
// and that there are no ephemeral keys on a live network. Root key is excluded from the check.
var OneEphemeralKeysLiveTestnetCheckFn = func(sethCfg *pkg_seth.Config) error {
	concurrency := sethCfg.GetMaxConcurrency()

	if sethCfg.IsSimulatedNetwork() {
		if concurrency < 1 {
			return fmt.Errorf(ErrInsufficientEphemeralKeys, 0)
		}

		return nil
	}

	if sethCfg.EphemeralAddrs != nil && int(*sethCfg.EphemeralAddrs) > 0 {
		ephMsg := `
			Error: Ephemeral Addresses Detected on Live Network

			Ephemeral addresses are currently set for use on a live network, which is not permitted. The number of ephemeral addresses set is %d. Please make the following update to your TOML configuration file to correct this:
			'[Seth] ephemeral_addresses_number = 0'

			Additionally, ensure the following requirements are met to run this test on a live network:
			1. Use more than one private key in your network configuration.
			`

		return errors.New(ephMsg)
	}

	if concurrency < 1 {
		return fmt.Errorf(ErrInsufficientStaticKeys, len(sethCfg.Network.PrivateKeys))
	}

	return nil
}

// OneEphemeralKeysLiveTestnetAutoFixFn checks whether there's at least one ephemeral key on a simulated network or at least one static key on a live network,
// and that there are no ephemeral keys on a live network (if ephemeral keys count is different from zero, it will disable them). Root key is excluded from the check.
var OneEphemeralKeysLiveTestnetAutoFixFn = func(sethCfg *pkg_seth.Config) error {
	concurrency := sethCfg.GetMaxConcurrency()

	if sethCfg.IsSimulatedNetwork() {
		if concurrency < 1 {
			return fmt.Errorf(ErrInsufficientEphemeralKeys, 0)
		}

		return nil
	}

	if sethCfg.EphemeralAddrs != nil && int(*sethCfg.EphemeralAddrs) > 0 {
		var zero int64 = 0
		sethCfg.EphemeralAddrs = &zero
	}

	if concurrency < 1 {
		return fmt.Errorf(ErrInsufficientStaticKeys, len(sethCfg.Network.PrivateKeys))
	}

	return nil
}

// GetChainClient returns a seth client for the given network after validating the config
func GetChainClient(c config.SethConfig, network blockchain.EVMNetwork) (*pkg_seth.Client, error) {
	return GetChainClientWithConfigFunction(c, network, noOpSethConfigFn)
}

// GetChainClientWithConfigFunction returns a seth client for the given network after validating the config and applying the config function
func GetChainClientWithConfigFunction(c config.SethConfig, network blockchain.EVMNetwork, configFn ConfigFunction) (*pkg_seth.Client, error) {
	readSethCfg := c.GetSethConfig()
	if readSethCfg == nil {
		return nil, errors.New("Seth config not found")
	}

	sethCfg, err := MergeSethAndEvmNetworkConfigs(network, *readSethCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error merging seth and evm network configs")
	}

	err = configFn(&sethCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error applying seth config function")
	}

	err = ValidateSethNetworkConfig(sethCfg.Network)
	if err != nil {
		return nil, errors.Wrapf(err, "Error validating seth network config")
	}

	chainClient, err := pkg_seth.NewClientWithConfig(&sethCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating seth client")
	}

	return chainClient, nil
}

// MergeSethAndEvmNetworkConfigs merges EVMNetwork to Seth config. If Seth config already has Network settings,
// it will return unchanged Seth config that was passed to it. If the network is simulated, it will
// use Geth-specific settings. Otherwise, it will use the chain ID to find the correct network settings.
// If no match is found it will return error.
func MergeSethAndEvmNetworkConfigs(evmNetwork blockchain.EVMNetwork, sethConfig pkg_seth.Config) (pkg_seth.Config, error) {
	if sethConfig.Network != nil {
		return sethConfig, nil
	}

	var sethNetwork *pkg_seth.Network

	mergeSimulatedNetworks := func(evmNetwork blockchain.EVMNetwork, sethNetwork pkg_seth.Network) *pkg_seth.Network {
		sethNetwork.PrivateKeys = evmNetwork.PrivateKeys
		if len(sethNetwork.URLs) == 0 {
			if len(evmNetwork.URLs) > 0 {
				sethNetwork.URLs = evmNetwork.URLs
			} else {
				sethNetwork.URLs = evmNetwork.HTTPURLs
			}
		}
		// important since Besu doesn't support EIP-1559, but other EVM clients do
		sethNetwork.EIP1559DynamicFees = evmNetwork.SupportsEIP1559
		// might be needed for cases, when node is incapable of estimating gas limit (e.g. Geth < v1.10.0)
		if evmNetwork.DefaultGasLimit != 0 {
			sethNetwork.GasLimit = evmNetwork.DefaultGasLimit
		}
		return &sethNetwork
	}

	for _, conf := range sethConfig.Networks {
		if evmNetwork.Simulated && evmNetwork.Name == pkg_seth.ANVIL && conf.Name == pkg_seth.ANVIL {
			// Merge Anvil network
			sethNetwork = mergeSimulatedNetworks(evmNetwork, *conf)
			break
		} else if evmNetwork.Simulated && conf.Name == pkg_seth.GETH {
			// Merge all other simulated Geth networks
			sethNetwork = mergeSimulatedNetworks(evmNetwork, *conf)
			break
		} else if isSameNetwork(conf, evmNetwork) {
			conf.PrivateKeys = evmNetwork.PrivateKeys
			if len(conf.URLs) == 0 {
				if len(evmNetwork.URLs) > 0 {
					conf.URLs = evmNetwork.URLs
				} else {
					conf.URLs = evmNetwork.HTTPURLs
				}
			}

			sethNetwork = conf
			break
		}
	}

	// If the network is not found, try to find the default network and replace it with the EVM network
	if sethNetwork == nil {
		for _, conf := range sethConfig.Networks {
			if conf.Name == fmt.Sprint(pkg_seth.DefaultNetworkName) {
				conf.Name = evmNetwork.Name
				conf.PrivateKeys = evmNetwork.PrivateKeys
				if len(evmNetwork.URLs) > 0 {
					conf.URLs = evmNetwork.URLs
				} else {
					conf.URLs = evmNetwork.HTTPURLs
				}
				sethNetwork = conf
				break
			}
		}
	}

	// If the network is still not found, return an error
	if sethNetwork == nil {
		msg := `Failed to build network config for chain ID %d. This could be the result of various reasons:
1. You are running tests for a network that hasn't been defined in known_networks.go and you have not defined it under [Network.EVMNetworks.NETWORK_NAME] in TOML
3. You have not defined Seth network settings for the chain ID %d in TOML under [Seth.Networks]
2. You have not defined a Seth Default network in your TOML config file under [Seth.Networks] using name %s`

		return pkg_seth.Config{}, fmt.Errorf(msg, evmNetwork.ChainID, evmNetwork.ChainID, pkg_seth.DefaultNetworkName)
	}
	sethConfig.Network = sethNetwork

	return sethConfig, nil
}

// MustReplaceSimulatedNetworkUrlWithK8 replaces the simulated network URL with the K8 URL and returns the network.
// If the network is not simulated, it will return the network unchanged.
func MustReplaceSimulatedNetworkUrlWithK8(l zerolog.Logger, network blockchain.EVMNetwork, testEnvironment environment.Environment) blockchain.EVMNetwork {
	if !network.Simulated {
		return network
	}

	networkKeys := []string{"Simulated Geth", "Simulated-Geth"}
	var keyToUse string

	for _, key := range networkKeys {
		_, ok := testEnvironment.URLs[key]
		if ok {
			keyToUse = key
			break
		}
	}

	if keyToUse == "" {
		for k := range testEnvironment.URLs {
			l.Info().Str("Network", k).Msg("Available networks")
		}
		panic("no network settings for Simulated Geth")
	}

	network.URLs = testEnvironment.URLs[keyToUse]

	return network
}

// ValidateSethNetworkConfig validates the Seth network config
func ValidateSethNetworkConfig(cfg *pkg_seth.Network) error {
	if cfg == nil {
		return errors.New("network cannot be nil")
	}
	if len(cfg.URLs) == 0 {
		return errors.New("URLs are required")
	}
	if len(cfg.PrivateKeys) == 0 {
		return errors.New("PrivateKeys are required")
	}
	if cfg.TransferGasFee == 0 {
		return errors.New("TransferGasFee needs to be above 0. It's the gas fee for a simple transfer transaction")
	}
	if cfg.TxnTimeout.Duration() == 0 {
		return errors.New("TxnTimeout needs to be above 0. It's the timeout for a transaction")
	}
	if cfg.EIP1559DynamicFees {
		if cfg.GasFeeCap == 0 {
			return errors.New("GasFeeCap needs to be above 0. It's the maximum fee per gas for a transaction (including tip)")
		}
		if cfg.GasTipCap == 0 {
			return errors.New("GasTipCap needs to be above 0. It's the maximum tip per gas for a transaction")
		}
		if cfg.GasFeeCap <= cfg.GasTipCap {
			return errors.New("GasFeeCap needs to be above GasTipCap (as it is base fee + tip cap)")
		}
	} else {
		if cfg.GasPrice == 0 {
			return errors.New("GasPrice needs to be above 0. It's the price of gas for a transaction")
		}
	}

	return nil
}

const RootKeyNum = 0

// AvailableSethKeyNum returns the available Seth address index
// If there are multiple addresses, it will return any synced key
// Otherwise it will return the root key
func AvailableSethKeyNum(client *pkg_seth.Client) int {
	if len(client.Addresses) > 1 {
		return client.AnySyncedKey()
	}
	return RootKeyNum
}

func isSameNetwork(conf *pkg_seth.Network, network blockchain.EVMNetwork) bool {
	if strings.EqualFold(conf.Name, fmt.Sprint(network.Name)) {
		return true
	}

	re := regexp.MustCompile(`[\s-]+`)
	cleanSethName := re.ReplaceAllString(conf.Name, "_")
	cleanNetworkName := re.ReplaceAllString(fmt.Sprint(network.Name), "_")

	return strings.EqualFold(cleanSethName, cleanNetworkName)
}
