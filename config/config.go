package config

import (
	"math/big"

	"github.com/spf13/viper"
)

type ConfigurationType string

type Config struct {
	Networks map[string]*NetworkConfig `mapstructure:"networks"`
	// "EthHardhat": conf, etc...
	DefaultKeyStore string
}

type NetworkConfig struct {
	Name            string   `mapstructure:"name"`
	URL             string   `mapstructure:"url"`
	ChainID         *big.Int `mapstructure:"chain_id"`
	PrivateKeyStore PrivateKeyStore
}

const (
	EnvironmentVariables ConfigurationType = "env"
	ConfigurationFile    ConfigurationType = "file"
	SecretStore          ConfigurationType = "secret"
)

// NewConfig creates a new configuration instance via viper from env vars, conig file, or a secret store
func NewConfig(configType ConfigurationType) (*Config, error) {
	v := viper.New()

	switch configType {
	case EnvironmentVariables:
		v.AutomaticEnv()
	case ConfigurationFile:
		v.SetConfigName("networks")
		v.SetConfigType("yml")
	case SecretStore:
		// Deal with secret store
	}
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = v.Unmarshal(conf)
	return conf, err
}

type PrivateKeyStore interface {
	Fetch() (string, error)
}

func NewPrivateKeyStore() {

}

type EnvStore struct{}

func (e *EnvStore) Fetch() (string, error) {
	return "", nil
}

type FileStore struct{}

func (f *FileStore) Fetch() (string, error) {
	return "", nil
}
