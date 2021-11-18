// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ConfigurationType refers to the different ways that configurations can be set
type ConfigurationType string

// Configs
const (
	LocalConfig  ConfigurationType = "local"
	SecretConfig ConfigurationType = "secret"
)

// FrameworkConfig common framework config
type FrameworkConfig struct {
	KeepEnvironments string         `mapstructure:"keep_environments" yaml:"keep_environments"`
	Logging          *LoggingConfig `mapstructure:"logging" yaml:"logging"`
}

// NetworksConfig is network configurations
type NetworksConfig struct {
	SelectedNetworks   []string                          `mapstructure:"selected_networks" yaml:"selected_networks"`
	NetworkSettings    map[string]map[string]interface{} `mapstructure:"networks" yaml:"networks"`
	DefaultKeyStore    string
	ConfigFileLocation string
}

// LoggingConfig for logging
type LoggingConfig struct {
	Level int8 `mapstructure:"level" yaml:"logging"`
}

// ETHNetwork data to configure fully ETH compatible network
type ETHNetwork struct {
	External             bool          `mapstructure:"external" yaml:"external"`
	Name                 string        `mapstructure:"name" yaml:"name"`
	ID                   string        `mapstructure:"id" yaml:"id"`
	ChainID              int64         `mapstructure:"chain_id" yaml:"chain_id"`
	URL                  string        `mapstructure:"url" yaml:"url"`
	URLs                 []string      `mapstructure:"urls" yaml:"urls"`
	Type                 string        `mapstructure:"type" yaml:"type"`
	PrivateKeys          []string      `mapstructure:"private_keys" yaml:"private_keys"`
	TransactionLimit     uint64        `mapstructure:"transaction_limit" yaml:"transaction_limit"`
	Timeout              time.Duration `mapstructure:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations int           `mapstructure:"minimum_confirmations" yaml:"minimum_confirmations"`
	GasEstimationBuffer  uint64        `mapstructure:"gas_estimation_buffer" yaml:"gas_estimation_buffer"`
	BlockGasLimit        uint64        `mapstructure:"block_gas_limit" yaml:"block_gas_limit"`
}

// TerraNetwork data to configure Terra network
type TerraNetwork struct {
	Name                 string        `mapstructure:"name" yaml:"name"`
	ChainName            string        `mapstructure:"chain_name" yaml:"chain_name"`
	Mnemonics            []string      `mapstructure:"mnemonic" yaml:"mnemonic"`
	Currency             string        `mapstructure:"currency" yaml:"currency"`
	Type                 string        `mapstructure:"type" yaml:"type"`
	TransactionLimit     uint64        `mapstructure:"transaction_limit" yaml:"transaction_limit"`
	Timeout              time.Duration `mapstructure:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations int           `mapstructure:"minimum_confirmations" yaml:"minimum_confirmations"`
}

func defaultViper(dir string, file string) *viper.Viper {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName(file)
	if dir == "" {
		v.AddConfigPath(".")
	} else {
		v.AddConfigPath(dir)
	}
	v.SetConfigType("yaml")
	return v
}

// LoadFrameworkConfig loads framework config
func LoadFrameworkConfig(cfgPath string) (*FrameworkConfig, error) {
	dir, file := path.Split(cfgPath)
	log.Info().
		Str("Dir", dir).
		Str("File", file).
		Msg("Loading config file")
	v := defaultViper(dir, file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg *FrameworkConfig
	err := v.Unmarshal(&cfg)
	return cfg, err
}

// LoadNetworksConfig loads networks config
func LoadNetworksConfig(cfgPath string) (*NetworksConfig, error) {
	dir, file := path.Split(cfgPath)
	log.Info().
		Str("Dir", dir).
		Str("File", file).
		Msg("Loading config file")
	v := defaultViper(dir, file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg *NetworksConfig
	err := v.Unmarshal(&cfg)
	return cfg, err
}

// PrivateKeyStore enables access, through a variety of methods, to private keys for use in blockchain networks
type PrivateKeyStore interface {
	Fetch() ([]string, error)
}

// LocalStore retrieves keys defined in a networks.yaml file, or from environment variables
type LocalStore struct {
	RawKeys []string
}

// Fetch private keys from local environment variables or a config file
func (l *LocalStore) Fetch() ([]string, error) {
	if l.RawKeys == nil {
		return nil, errors.New("no keys found, ensure your configuration is properly set")
	}
	return l.RawKeys, nil
}
