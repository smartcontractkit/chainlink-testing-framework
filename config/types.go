package config

import "github.com/smartcontractkit/seth"

type SethConfig interface {
	GetSethConfig() *seth.Config
}

type NamedConfiguration interface {
	GetConfigurationName() string
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
