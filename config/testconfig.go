package config

import (
	"github.com/smartcontractkit/seth"
)

type GlobalTestConfig interface {
	GetChainlinkImageConfig() *ChainlinkImageConfig
	GetLoggingConfig() *LoggingConfig
	GetNetworkConfig() *NetworkConfig
	GetPrivateEthereumNetworkConfig() *EthereumNetworkConfig
	GetPyroscopeConfig() *PyroscopeConfig
	SethConfig
}

func (c *TestConfig) GetLoggingConfig() *LoggingConfig {
	return c.Logging
}

func (c TestConfig) GetNetworkConfig() *NetworkConfig {
	return c.Network
}

func (c TestConfig) GetChainlinkImageConfig() *ChainlinkImageConfig {
	return c.ChainlinkImage
}

func (c TestConfig) GetPrivateEthereumNetworkConfig() *EthereumNetworkConfig {
	return c.PrivateEthereumNetwork
}

func (c TestConfig) GetPyroscopeConfig() *PyroscopeConfig {
	return c.Pyroscope
}

type NamedConfiguration interface {
	GetConfigurationName() string
}

type SethConfig interface {
	GetSethConfig() *seth.Config
}

type TestConfig struct {
	ChainlinkImage         *ChainlinkImageConfig  `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ChainlinkImageConfig  `toml:"ChainlinkUpgradeImage"`
	Logging                *LoggingConfig         `toml:"Logging"`
	Network                *NetworkConfig         `toml:"Network"`
	Pyroscope              *PyroscopeConfig       `toml:"Pyroscope"`
	PrivateEthereumNetwork *EthereumNetworkConfig `toml:"PrivateEthereumNetwork"`
	WaspConfig             *WaspAutoBuildConfig   `toml:"WaspAutoBuild"`
	Seth                   *seth.Config           `toml:"Seth"`
}
