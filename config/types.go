package config

import "github.com/smartcontractkit/seth"

type SethConfig interface {
	GetSethConfig() *seth.Config
}

type NamedConfigurations interface {
	GetConfigurationNames() []string
}

type GlobalTestConfig interface {
	GetChainlinkImageConfig() *ChainlinkImageConfig
	GetLoggingConfig() *LoggingConfig
	GetNetworkConfig() *NetworkConfig
	GetPrivateEthereumNetworkConfig() *EthereumNetworkConfig
	GetPyroscopeConfig() *PyroscopeConfig
	GetNodeConfig() *NodeConfig
	SethConfig
}
