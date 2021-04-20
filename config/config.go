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
	Name            string          `mapstructure:"name"`
	URL             string          `mapstructure:"url"`
	ChainID         int             `mapstructure:"chain_id"`
	PrivateKeyStore PrivateKeyStore `mapstructure:"private_keys"`
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
		v.AddConfigPath("./config/")
		v.AddConfigPath("../config/") // Not a huge fan of this, alternatives?
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
	for _, networkConf := range conf.Networks {
		networkConf.PrivateKeyStore = &FileStore{} // TODO: Adjust as needed for config type sent in
	}
	return conf, err
}

type PrivateKeyStore interface {
	Fetch() (string, error)
}

type EnvStore struct{}

func (e *EnvStore) Fetch() (string, error) {
	return "", nil
}

type FileStore struct{}

func (f *FileStore) Fetch() (string, error) {
	return "", nil
}
