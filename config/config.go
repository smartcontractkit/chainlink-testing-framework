package config

import (
	"github.com/spf13/viper"
)

type ConfigurationType string

type Config struct {
	Networks map[string]*NetworkConfig `mapstructure:"networks"`
	// "EthHardhat": conf, etc...
	DefaultKeyStore string
}

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

// NewConfig creates a new configuration instance via viper from env vars, conig file, or a secret store
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
		err := v.ReadInConfig()
		if err != nil {
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

type PrivateKeyStore interface {
	Fetch() ([]string, error)
}

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

type EnvStore struct {
	rawKeys []string
}

func (e *EnvStore) Fetch() ([]string, error) {
	return e.rawKeys, nil
}

type FileStore struct {
	rawKeys []string
}

func (f *FileStore) Fetch() ([]string, error) {
	return f.rawKeys, nil
}

type SecretStore struct{}

func (s *SecretStore) Fetch() ([]string, error) {
	// TODO: Set up connection with whatever secret store we choose
	return []string{""}, nil
}
