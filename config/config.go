// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"fmt"
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

// Config is the overall config for the framework, holding configurations for supported networks
type Config struct {
	Networks           []string                 `mapstructure:"networks" yaml:"networks"`
	Logging            *LoggingConfig           `mapstructure:"logging" yaml:"logging"`
	NetworkConfigs     map[string]NetworkConfig `mapstructure:"network_configs" yaml:"network_configs"`
	Retry              *RetryConfig             `mapstructure:"retry" yaml:"retry"`
	Apps               AppConfig                `mapstructure:"apps" yaml:"apps"`
	Kubernetes         KubernetesConfig         `mapstructure:"kubernetes" yaml:"kubernetes"`
	KeepEnvironments   string                   `mapstructure:"keep_environments" yaml:"keep_environments"`
	Prometheus         *PrometheusConfig        `mapstructure:"prometheus" yaml:"prometheus"`
	Contracts          *ContractsConfig         `mapstructure:"contracts" yaml:"contracts"`
	DefaultKeyStore    string
	ConfigFileLocation string
}

// PrometheusConfig for prometheus
type PrometheusConfig struct {
	URL string `mapstructure:"url" yaml:"url"`
}

// LoggingConfig for logging
type LoggingConfig struct {
	Level int8 `mapstructure:"level" yaml:"logging"`
}

// GetNetworkConfig finds a specified network config based on its name
func (c *Config) GetNetworkConfig(name string) (NetworkConfig, error) {
	if network, ok := c.NetworkConfigs[name]; ok {
		return network, nil
	}
	return NetworkConfig{}, fmt.Errorf("no supported network of name '%s' was found. Ensure that the config for it exists.", name)
}

// ContractsConfig contracts sources config
type ContractsConfig struct {
	Ethereum EthereumSources `mapstructure:"ethereum" yaml:"ethereum"`
}

// EthereumSources sources to generate bindings to ethereum contracts
type EthereumSources struct {
	ExecutablePath string     `mapstructure:"executable_path" yaml:"executable_path"`
	OutPath        string     `mapstructure:"out_path" yaml:"out_path"`
	Sources        SourcesMap `mapstructure:"sources" yaml:"sources"`
}

// SourcesMap describes different sources types, local or remote s3 (external)
type SourcesMap struct {
	Local    LocalSource     `mapstructure:"local" yaml:"local"`
	External ExternalSources `mapstructure:"external" yaml:"external"`
}

// ExternalSources are sources downloaded from remote
type ExternalSources struct {
	RootPath     string                    `mapstructure:"path" yaml:"path"`
	Region       string                    `mapstructure:"region" yaml:"region"`
	S3URL        string                    `mapstructure:"s3_path" yaml:"s3_path"`
	Repositories map[string]ExternalSource `mapstructure:"repositories" yaml:"repositories"`
}

// LocalSource local contracts artifacts directory
type LocalSource struct {
	Path string `mapstructure:"path" yaml:"path"`
}

// ExternalSource remote contracts artifacts source directory
type ExternalSource struct {
	Path   string `mapstructure:"path" yaml:"path"`
	Commit string `mapstructure:"commit" yaml:"commit"`
}

// NetworkConfig holds the basic values that identify a blockchain network and contains private keys on the network
type NetworkConfig struct {
	Name                 string   `mapstructure:"name" yaml:"name"`
	ChainName            string   `mapstructure:"chain_name" yaml:"chain_name"`
	Mnemonics            []string `mapstructure:"mnemonic" yaml:"mnemonic"`
	Currency             string   `mapstructure:"currency" yaml:"currency"`
	ClusterURL           string
	LocalURL             string
	URLS                 []string      `mapstructure:"urls" yaml:"urls"`
	ChainID              int64         `mapstructure:"chain_id" yaml:"chain_id"`
	Type                 string        `mapstructure:"type" yaml:"type"`
	SecretPrivateKeys    bool          `mapstructure:"secret_private_keys" yaml:"secret_private_keys"`
	SecretPrivateURL     bool          `mapstructure:"secret_private_url" yaml:"secret_private_url"`
	NamespaceForSecret   string        `mapstructure:"namespace_for_secret" yaml:"namespace_for_secret"`
	PrivateKeys          []string      `mapstructure:"private_keys" yaml:"private_keys"`
	PrivateURL           string        `mapstructure:"private_url" yaml:"private_url"`
	TransactionLimit     uint64        `mapstructure:"transaction_limit" yaml:"transaction_limit"`
	Timeout              time.Duration `mapstructure:"transaction_timeout" yaml:"transaction_timeout"`
	LinkTokenAddress     string        `mapstructure:"link_token_address" yaml:"link_token_address"`
	MinimumConfirmations int           `mapstructure:"minimum_confirmations" yaml:"minimum_confirmations"`
	GasEstimationBuffer  uint64        `mapstructure:"gas_estimation_buffer" yaml:"gas_estimation_buffer"`
	BlockGasLimit        uint64        `mapstructure:"block_gas_limit" yaml:"block_gas_limit"`
	RPCPort              uint16        `mapstructure:"rpc_port" yaml:"rpc_port"`
	PrivateKeyStore      PrivateKeyStore
}

// KubernetesConfig holds the configuration for how the framework interacts with the k8s cluster
type KubernetesConfig struct {
	QPS               float32       `mapstructure:"qps" yaml:"qps"`
	Burst             int           `mapstructure:"burst" yaml:"burst"`
	DeploymentTimeout time.Duration `mapstructure:"deployment_timeout" yaml:"deployment_timeout"`
}

// AppConfig holds all the configuration for the core apps that are deployed for testing
type AppConfig struct {
	Chainlink        StandardConfig `mapstructure:"chainlink" yaml:"chainlink"`
	Geth             StandardConfig `mapstructure:"geth" yaml:"geth"`
	Adapter          StandardConfig `mapstructure:"adapter" yaml:"adapter"`
	Postgres         StandardConfig `mapstructure:"postgres" yaml:"postgres"`
	Otpe             StandardConfig `mapstructure:"otpe" yaml:"otpe"`
	Explorer         StandardConfig `mapstructure:"explorer" yaml:"explorer"`
	AtlasEvm         StandardConfig `mapstructure:"atlas-evm" yaml:"atlas-evm"`
	CpSchemaRegistry StandardConfig `mapstructure:"cp-schema-registry" yaml:"cp-schema-registry"`
	Prometheus       StandardConfig `mapstructure:"prometheus" yaml:"prometheus"`
	KafkaRest        StandardConfig `mapstructure:"kafka-rest" yaml:"kafka-rest"`
}

// ResourcesConfig hols the resource usage configuration for a pod
type ResourcesConfig struct {
	Memory string `mapstructure:"memory" yaml:"memory"`
	Cpu    string `mapstructure:"cpu" yaml:"cpu"`
}

// StandardConfig holds the configuration for an app to be deployed
type StandardConfig struct {
	Image    string          `mapstructure:"image" yaml:"image"`
	Version  string          `mapstructure:"version" yaml:"version"`
	Requests ResourcesConfig `mapstructure:"requests" yaml:"requests"`
	Limits   ResourcesConfig `mapstructure:"limits" yaml:"limits"`
}

// NewConfig creates a new configuration instance via viper from env vars, config file, or a secret store
func NewConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := &Config{
		ConfigFileLocation: strings.TrimRight(v.ConfigFileUsed(), "config.yml"),
	}
	log.Info().Str("File Location", v.ConfigFileUsed()).Msg("Loading config file")
	err := v.Unmarshal(conf)
	return conf, err
}

// PrivateKeyStore enables access, through a variety of methods, to private keys for use in blockchain networks
type PrivateKeyStore interface {
	Fetch() ([]string, error)
}

// LocalStore retrieves keys defined in a config.yml file, or from environment variables
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

// RetryConfig holds config for retry attempts and delays
type RetryConfig struct {
	Attempts    uint          `mapstructure:"attempts" yaml:"attempts"`
	LinearDelay time.Duration `mapstructure:"linear_delay" yaml:"linear_delay"`
}
