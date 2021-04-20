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
	EnvironmentConfig ConfigurationType = "env"
	FileConfig        ConfigurationType = "file"
	SecretConfig      ConfigurationType = "secret"
)

// NewConfig creates a new configuration instance via viper from env vars, config file, or a secret store
func NewConfig(configType ConfigurationType) (*Config, error) {
	v := viper.New()

	switch configType {
	case EnvironmentConfig:
		v.AutomaticEnv()
	case FileConfig:
		v.SetConfigName("networks")
		v.AddConfigPath("./config/")
		v.AddConfigPath("../config/") // Not a huge fan of this, alternatives?
		v.SetConfigType("yml")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	case SecretConfig:
		// Deal with secret store
	}

	conf := &Config{}
	err := v.Unmarshal(conf)
	for _, networkConf := range conf.Networks {
		networkConf.PrivateKeys = NewPrivateKeyStore(configType, networkConf.RawKeys)
	}
	return conf, err
}

// PrivateKeyStore enables access, through a variety of methods, to private keys for use in blockchain networks
type PrivateKeyStore interface {
	Fetch() ([]string, error)
}

// NewPrivateKeyStore returns a keystore of a specific type, depending on where it should source its keys from
func NewPrivateKeyStore(configType ConfigurationType, keys []string) PrivateKeyStore {
	switch configType {
	case EnvironmentConfig:
		return &EnvStore{keys}
	case FileConfig:
		return &FileStore{keys}
	case SecretConfig:
		return &SecretStore{}
	}
	return nil
}

// EnvStore retrieves keys dictated in environment variables
type EnvStore struct {
	rawKeys []string
}

func (e *EnvStore) Fetch() ([]string, error) {
	return e.rawKeys, nil
}

// FileStore retrieves keys defined in a networks.yml config file
type FileStore struct {
	rawKeys []string
}

func (f *FileStore) Fetch() ([]string, error) {
	return f.rawKeys, nil
}

// SecretStore retrieves keys from an encrypted secret storage service TBD
type SecretStore struct{}

func (s *SecretStore) Fetch() ([]string, error) {
	// TODO: Set up connection with whatever secret store we choose
	return []string{""}, nil
}
