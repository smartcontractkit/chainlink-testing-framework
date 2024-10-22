package networks

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
)

const (
	pyroscopeTOML = `[Pyroscope]
ServerAddress = '%s'
Environment = '%s'`
	//nolint:gosec //ignoring G101
	secretTOML = `
[Mercury.Credentials.cred1]
URL = '%s'
Username = '%s'
Password = '%s'
`
)

// AddNetworksConfig adds EVM network configurations to a base config TOML. Useful for adding networks with default
// settings. See AddNetworkDetailedConfig for adding more detailed network configuration.
func AddNetworksConfig(baseTOML string, pyroscopeConfig *config.PyroscopeConfig, networks ...blockchain.EVMNetwork) string {
	networksToml := ""
	for _, network := range networks {
		networksToml = fmt.Sprintf("%s\n\n%s", networksToml, network.MustChainlinkTOML(""))
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s", baseTOML, pyroscopeSettings(pyroscopeConfig), networksToml)
}

func AddSecretTomlConfig(url, username, password string) string {
	return fmt.Sprintf(secretTOML, url, username, password)
}

// AddNetworkDetailedConfig adds EVM config to a base TOML. Also takes a detailed network config TOML where values like
// using transaction forwarders can be included.
// See https://github.com/smartcontractkit/chainlink/blob/develop/docs/CONFIG.md#EVM
func AddNetworkDetailedConfig(baseTOML string, pyroscopeConfig *config.PyroscopeConfig, detailedNetworkConfig string, network blockchain.EVMNetwork) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", baseTOML, pyroscopeSettings(pyroscopeConfig), network.MustChainlinkTOML(detailedNetworkConfig))
}

func pyroscopeSettings(config *config.PyroscopeConfig) string {
	if config == nil || config.Enabled == nil || !*config.Enabled {
		return ""
	}
	return fmt.Sprintf(pyroscopeTOML, *config.ServerUrl, *config.Environment)
}
