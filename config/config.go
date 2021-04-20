package config

import (
	"github.com/spf13/viper"
)

// ConfigurationType refers to the different ways that configurations can be set
type ConfigurationType string

// Config is the overall config for the framework, holding configurations for supported networks
type Config struct {
	Networks map[string]*NetworkConfig `mapstructure:"networks"`
	// "EthHardhat": conf, etc...
	DefaultKeyStore string
}

// NetworkConfig holds the basic values that identify a blockchain network and contains private keys on the network
type NetworkConfig struct {
	Name        string   `mapstructure:"name"`
	URL         string   `mapstructure:"url"`
	ChainID     int64    `mapstructure:"chain_id"`
	RawKeys     []string `mapstructure:"private_keys"`
	PrivateKeys PrivateKeyStore
}

const (
	LocalConfig  ConfigurationType = "local"
	SecretConfig ConfigurationType = "secret"
)

// NewConfig creates a new configuration instance via viper from env vars, config file, or a secret store
func NewConfig(configType ConfigurationType) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	// File acts as defaults
	v.SetConfigName("networks")
	v.AddConfigPath("./config/")
	v.AddConfigPath("../config/")
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := &Config{}
	err := v.Unmarshal(conf)
	for _, networkConf := range conf.Networks {
		networkConf.PrivateKeys = NewPrivateKeyStore(configType, networkConf.RawKeys, networkConf.Name)
	}
	return conf, err
}

// PrivateKeyStore enables access, through a variety of methods, to private keys for use in blockchain networks
type PrivateKeyStore interface {
	Fetch() ([]string, error)
}

// NewPrivateKeyStore returns a keystore of a specific type, depending on where it should source its keys from
func NewPrivateKeyStore(configType ConfigurationType, keys []string, networkName string) PrivateKeyStore {
	switch configType {
	case LocalConfig:
		return &LocalStore{keys}
	case SecretConfig:
		return &SecretStore{networkName}
	}
	return nil
}

// FileStore retrieves keys defined in a networks.yml config file
type LocalStore struct {
	rawKeys []string
}

// Fetch private keys from local environment variables or a config file
func (l *LocalStore) Fetch() ([]string, error) {
	return l.rawKeys, nil
}

// SecretStore retrieves keys from an encrypted secret storage service TBD
type SecretStore struct {
	networkName string
}

// Fetch private keys from env variables or a secret management system
func (s *SecretStore) Fetch() ([]string, error) {
	// TODO: Set up connection with whatever secret store we choose
	// Connect to secrets service / local encryption setup
	// Fetch keys based on the networkName
	// Return them
	return []string{""}, nil
}
