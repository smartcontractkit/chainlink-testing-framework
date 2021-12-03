// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"

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
	EnvironmentFile  string         `mapstructure:"environment_file" yaml:"environment_file"`
	ChainlinkImage   string         `mapstructure:"chainlink_image" yaml:"chainlink_image"`
	ChainlinkVersion string         `mapstructure:"chainlink_version" yaml:"chainlink_version"`
	GethImage        string         `mapstructure:"geth_image" yaml:"geth_image"`
	GethVersion      string         `mapstructure:"geth_version" yaml:"geth_version"`
}

// NetworkSettings is a map that holds configuration for each individual network
type NetworkSettings map[string]map[string]interface{}

// Decode is used by envconfig to initialise the custom Charts type with populated values
// This function will take a JSON object representing charts, and unmarshal it into the existing object to "merge" the
// two
func (n NetworkSettings) Decode(value string) error {
	// Support the use of files for unmarshaling charts JSON
	if _, err := os.Stat(value); err == nil {
		b, err := os.ReadFile(value)
		if err != nil {
			return err
		}
		value = string(b)
	}
	networkSettings := NetworkSettings{}
	if err := yaml.Unmarshal([]byte(value), &networkSettings); err != nil {
		return fmt.Errorf("failed to unmarshal YAML, either a file path specific doesn't exist, or the YAML is invalid: %v", err)
	}
	return mergo.Merge(&n, networkSettings, mergo.WithOverride)
}

// NetworksConfig is network configurations
type NetworksConfig struct {
	SelectedNetworks   []string        `mapstructure:"selected_networks" yaml:"selected_networks" envconfig:"selected_networks"`
	NetworkSettings    NetworkSettings `mapstructure:"networks" yaml:"networks" envconfig:"network_settings"`
	DefaultKeyStore    string
	ConfigFileLocation string
}

// LoggingConfig for logging
type LoggingConfig struct {
	Level int8 `mapstructure:"level" yaml:"logging"`
}

// ETHNetwork data to configure fully ETH compatible network
type ETHNetwork struct {
	External                  bool          `mapstructure:"external" yaml:"external"`
	Name                      string        `mapstructure:"name" yaml:"name"`
	ID                        string        `mapstructure:"id" yaml:"id"`
	ChainID                   int64         `mapstructure:"chain_id" yaml:"chain_id"`
	URL                       string        `mapstructure:"url" yaml:"url"`
	URLs                      []string      `mapstructure:"urls" yaml:"urls"`
	Type                      string        `mapstructure:"type" yaml:"type"`
	PrivateKeys               []string      `mapstructure:"private_keys" yaml:"private_keys"`
	ChainlinkTransactionLimit uint64        `mapstructure:"chainlink_transaction_limit" yaml:"chainlink_transaction_limit"`
	Timeout                   time.Duration `mapstructure:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations      int           `mapstructure:"minimum_confirmations" yaml:"minimum_confirmations"`
	GasEstimationBuffer       uint64        `mapstructure:"gas_estimation_buffer" yaml:"gas_estimation_buffer"`
	BlockGasLimit             uint64        `mapstructure:"block_gas_limit" yaml:"block_gas_limit"`
}

// TerraNetwork data to configure Terra network
type TerraNetwork struct {
	Name                      string        `mapstructure:"name" yaml:"name"`
	ChainName                 string        `mapstructure:"chain_name" yaml:"chain_name"`
	Mnemonics                 []string      `mapstructure:"mnemonic" yaml:"mnemonic"`
	Currency                  string        `mapstructure:"currency" yaml:"currency"`
	Type                      string        `mapstructure:"type" yaml:"type"`
	ChainlinkTransactionLimit uint64        `mapstructure:"chainlink_transaction_limit" yaml:"chainlink_transaction_limit"`
	Timeout                   time.Duration `mapstructure:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations      int           `mapstructure:"minimum_confirmations" yaml:"minimum_confirmations"`
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

	// Allow the networks config to be overridden when this codebase is imported as a library
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
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
