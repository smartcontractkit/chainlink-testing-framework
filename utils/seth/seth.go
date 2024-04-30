package seth

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	pkg_seth "github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
)

var INSUFFICIENT_EPHEMERAL_KEYS = `
Error: Insufficient Ephemeral Addresses for Simulated Network

To operate on a simulated network, you must configure at least one ephemeral address. Currently, %d ephemeral address(es) are set. Please update your TOML configuration file as follows to meet this requirement:
[Seth] ephemeral_addresses_number = 1

This adjustment ensures that your setup is minimaly viable. Although it is highly recommended to use at least 20 ephemeral addresses.
`

var INSUFFICIENT_STATIC_KEYS = `
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
// and that there are no epehemeral keys on a live network. Root key is excluded from the check.
var OneEphemeralKeysLiveTestnetCheckFn = func(sethCfg *pkg_seth.Config) error {
	concurrency := sethCfg.GetMaxConcurrency()

	if sethCfg.IsSimulatedNetwork() {
		if concurrency < 1 {
			return fmt.Errorf(INSUFFICIENT_EPHEMERAL_KEYS, 0)
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
		return fmt.Errorf(INSUFFICIENT_STATIC_KEYS, len(sethCfg.Network.PrivateKeys))
	}

	return nil
}

// OneEphemeralKeysLiveTestnetAutoFixFn checks whether there's at least one ephemeral key on a simulated network or at least one static key on a live network,
// and that there are no epehemeral keys on a live network (if ephemeral keys count is different from zero, it will disable them). Root key is excluded from the check.
var OneEphemeralKeysLiveTestnetAutoFixFn = func(sethCfg *pkg_seth.Config) error {
	concurrency := sethCfg.GetMaxConcurrency()

	if sethCfg.IsSimulatedNetwork() {
		if concurrency < 1 {
			return fmt.Errorf(INSUFFICIENT_EPHEMERAL_KEYS, 0)
		}

		return nil
	}

	if sethCfg.EphemeralAddrs != nil && int(*sethCfg.EphemeralAddrs) > 0 {
		var zero int64 = 0
		sethCfg.EphemeralAddrs = &zero
	}

	if concurrency < 1 {
		return fmt.Errorf(INSUFFICIENT_STATIC_KEYS, len(sethCfg.Network.PrivateKeys))
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
		return nil, fmt.Errorf("Seth config not found")
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
// use Geth-specific settings. Otherwise it will use the chain ID to find the correct network settings.
// If no match is found it will return error.
func MergeSethAndEvmNetworkConfigs(evmNetwork blockchain.EVMNetwork, sethConfig pkg_seth.Config) (pkg_seth.Config, error) {
	if sethConfig.Network != nil {
		return sethConfig, nil
	}

	var sethNetwork *pkg_seth.Network

	for _, conf := range sethConfig.Networks {
		if evmNetwork.Simulated {
			if conf.Name == pkg_seth.GETH {
				conf.PrivateKeys = evmNetwork.PrivateKeys
				if len(conf.URLs) == 0 {
					conf.URLs = evmNetwork.URLs
				}
				// important since Besu doesn't support EIP-1559, but other EVM clients do
				conf.EIP1559DynamicFees = evmNetwork.SupportsEIP1559

				// might be needed for cases, when node is incapable of estimating gas limit (e.g. Geth < v1.10.0)
				if evmNetwork.DefaultGasLimit != 0 {
					conf.GasLimit = evmNetwork.DefaultGasLimit
				}

				sethNetwork = conf
				break
			}
		} else if conf.ChainID == fmt.Sprint(evmNetwork.ChainID) {
			conf.PrivateKeys = evmNetwork.PrivateKeys
			if len(conf.URLs) == 0 {
				conf.URLs = evmNetwork.URLs
			}

			sethNetwork = conf
			break
		}
	}

	if sethNetwork == nil {
		return pkg_seth.Config{}, fmt.Errorf("No matching EVM network found for chain ID %d. If it's a new network please define it as [Network.EVMNetworks.NETWORK_NAME] in TOML", evmNetwork.ChainID)
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
		return fmt.Errorf("Network cannot be nil")
	}
	if cfg.ChainID == "" {
		return fmt.Errorf("ChainID is required")
	}
	_, err := strconv.Atoi(cfg.ChainID)
	if err != nil {
		return fmt.Errorf("ChainID needs to be a number")
	}
	if cfg.URLs == nil || len(cfg.URLs) == 0 {
		return fmt.Errorf("URLs are required")
	}
	if cfg.PrivateKeys == nil || len(cfg.PrivateKeys) == 0 {
		return fmt.Errorf("PrivateKeys are required")
	}
	if cfg.TransferGasFee == 0 {
		return fmt.Errorf("TransferGasFee needs to be above 0. It's the gas fee for a simple transfer transaction")
	}
	if cfg.TxnTimeout.Duration() == 0 {
		return fmt.Errorf("TxnTimeout needs to be above 0. It's the timeout for a transaction")
	}
	if cfg.EIP1559DynamicFees {
		if cfg.GasFeeCap == 0 {
			return fmt.Errorf("GasFeeCap needs to be above 0. It's the maximum fee per gas for a transaction (including tip)")
		}
		if cfg.GasTipCap == 0 {
			return fmt.Errorf("GasTipCap needs to be above 0. It's the maximum tip per gas for a transaction")
		}
		if cfg.GasFeeCap <= cfg.GasTipCap {
			return fmt.Errorf("GasFeeCap needs to be above GasTipCap (as it is base fee + tip cap)")
		}
	} else {
		if cfg.GasPrice == 0 {
			return fmt.Errorf("GasPrice needs to be above 0. It's the price of gas for a transaction")
		}
	}

	return nil
}
