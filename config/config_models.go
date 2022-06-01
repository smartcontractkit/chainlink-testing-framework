package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"gopkg.in/yaml.v2"
)

// Config general framework config
type Config struct {
	FrameworkConfig    *FrameworkConfig    `envconfig:"FRAMEWORK_CONFIG_FILE" default:"../framework.yaml"`
	NetworksConfig     *NetworksConfig     `envconfig:"NETWORKS_CONFIG_FILE" default:"../networks.yaml"`
	RemoteRunnerConfig *RemoteRunnerConfig `envconfig:"REMOTE_RUNNER_CONFIG_FILE" required:"false" default:"../remote_runner_config.yaml"`
	EnvironmentConfig  *environment.Config `envconfig:"ENVIRONMENT_CONFIG_FILE"`
}

// FrameworkConfig common framework config
type FrameworkConfig struct {
	KeepEnvironments string         `envconfig:"KEEP_ENVIRONMENTS" yaml:"keep_environments"`
	Logging          *LoggingConfig `envconfig:"LOGGING" yaml:"logging"`
	ChainlinkImage   string         `yaml:"chainlink_image" envconfig:"CHAINLINK_IMAGE"`
	ChainlinkVersion string         `yaml:"chainlink_version" envconfig:"CHAINLINK_VERSION"`
	// ChainlinkEnvValues uses interface{} as the value because it's needed for proper helmchart merges
	ChainlinkEnvValues map[string]interface{} `envconfig:"CHAINLINK_ENV_VALUES" yaml:"chainlink_env_values"`
}

func (m *FrameworkConfig) Decode(path string) error {
	// Marshal YAML first, then "envconfig" tags of that struct got marshalled
	if err := unmarshalYAML(path, &m); err != nil {
		return err
	}
	return envconfig.Process("", m)
}

// NetworksConfig holds all network configurations
type NetworksConfig struct {
	SelectedNetworks   []string            `envconfig:"SELECTED_NETWORKS" yaml:"selected_networks"`
	EVMNetworkSettings *EVMNetworkSettings `envconfig:"EVM_NETWORKS" yaml:"evm_networks"`
	DefaultKeyStore    string
	ConfigFileLocation string
}

func (m *NetworksConfig) Decode(path string) error {
	// Marshal YAML first, then "envconfig" tags of that struct got marshalled
	if err := unmarshalYAML(path, &m); err != nil {
		return err
	}
	return envconfig.Process("", m)
}

// EVMNetwork data to configure fully ETH compatible network
type EVMNetwork struct {
	Name                      string        `envconfig:"name" yaml:"name"`
	ChainID                   int64         `envconfig:"chain_id" yaml:"chain_id"`
	URLs                      []*EVMUrls    `envconfig:"urls" yaml:"urls"`
	Type                      string        `envconfig:"type" yaml:"type"`
	PrivateKeys               []string      `envconfig:"private_keys" yaml:"private_keys"`
	ChainlinkTransactionLimit uint64        `envconfig:"chainlink_transaction_limit" yaml:"chainlink_transaction_limit"`
	Timeout                   time.Duration `envconfig:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations      int           `envconfig:"minimum_confirmations" yaml:"minimum_confirmations"`
	GasEstimationBuffer       uint64        `envconfig:"gas_estimation_buffer" yaml:"gas_estimation_buffer"`
	ContractsDeployed         bool          `envconfig:"contracts_deployed" yaml:"contracts_deployed"`
}

type EVMUrls struct {
	WebSocket string `envconfig:"ws_url" yaml:"ws_url"`
	HTTP      string `envconfig:"http_url" yaml:"http_url"`
}

// EVMNetworkSettings holds settings for all EVM networks
type EVMNetworkSettings map[string]EVMNetwork

// GetNetworkSettings pulls network settings for an EVM network name. Returns false if network isn't found
func (e *EVMNetworkSettings) GetNetworkSettings(networkName string) (*EVMNetwork, bool) {
	network, exists := (*e)[networkName]
	return &network, exists
}

// Decode is used by envconfig to initialize the custom Charts type with populated values
// This function will take a JSON object representing charts, and unmarshal it into the existing object to "merge" the
// two
func (e *EVMNetworkSettings) Decode(value string) error {
	// Support the use of files for unmarshaling charts JSON
	if _, err := os.Stat(value); err == nil {
		b, err := os.ReadFile(value)
		if err != nil {
			return err
		}
		value = string(b)
	}
	networkSettings := EVMNetworkSettings{}
	if err := yaml.Unmarshal([]byte(value), &networkSettings); err != nil {
		return fmt.Errorf("failed to unmarshal YAML, either a file path specific doesn't exist, or the YAML is invalid: %v", err)
	}
	return mergo.Merge(&e, networkSettings, mergo.WithOverride)
}

// TerraNetwork data to configure Terra network
type TerraNetwork struct {
	Name                      string        `envconfig:"name" yaml:"name"`
	ChainName                 string        `envconfig:"chain_name" yaml:"chain_name"`
	Mnemonics                 []string      `envconfig:"mnemonic" yaml:"mnemonic"`
	Currency                  string        `envconfig:"currency" yaml:"currency"`
	Type                      string        `envconfig:"type" yaml:"type"`
	ChainlinkTransactionLimit uint64        `envconfig:"chainlink_transaction_limit" yaml:"chainlink_transaction_limit"`
	Timeout                   time.Duration `envconfig:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations      int           `envconfig:"minimum_confirmations" yaml:"minimum_confirmations"`
}

// RemoteRunnerConfig reads the config file for remote test runs
type RemoteRunnerConfig struct {
	TestRegex     string   `envconfig:"TEST_REGEX" yaml:"test_regex"`
	TestDirectory string   `envconfig:"TEST_DIRECTORY" yaml:"test_directory"`
	SlackAPIKey   string   `envconfig:"SLACK_API_KEY" yaml:"slack_api_key"`
	SlackChannel  string   `envconfig:"SLACK_CHANNEL" yaml:"slack_channel"`
	SlackUserID   string   `envconfig:"SLACK_USER_ID" yaml:"slack_user_id"`
	CustomEnvVars []string `envconfig:"CUSTOM_ENV_VARS" yaml:"custom_env_vars"`
}

func (m *RemoteRunnerConfig) Decode(path string) error {
	// Marshal YAML first, then "envconfig" tags of that struct got marshalled
	if err := unmarshalYAML(path, &m); err != nil {
		return err
	}
	return envconfig.Process("", m)
}

// LoggingConfig for logging
type LoggingConfig struct {
	WritePodLogs string `envconfig:"WRITE_POD_LOGS" yaml:"write_pod_logs"`
	Level        int8   `envconfig:"LEVEL" yaml:"level"`
}

// ChartOverrides enables building json styled chart overrides for the deployed chart values and environment variables
type ChartOverrides struct {
	GethChartOverride       *GethChart      `json:"geth,omitempty"`
	ChainlinkChartOverrride *ChainlinkChart `json:"chainlink,omitempty"`
}

// GethChart holds the overall geth chart values
type GethChart struct {
	Values *GethValuesWrapper `json:"values,omitempty"`
}

// GethValuesWrapper geth values wrapper
type GethValuesWrapper struct {
	GethVals *GethValues   `json:"geth,omitempty"`
	Args     []interface{} `json:"args,omitempty"`
}

// GethValues wraps all values
type GethValues struct {
	Image *GethImage `json:"image,omitempty"`
}

// GethImage defines geth image and version
type GethImage struct {
	Image   string `json:"image,omitempty" yaml:"geth_image"`
	Version string `json:"version,omitempty" yaml:"geth_version"`
}

// ChainlinkChart holds the overall geth chart values
type ChainlinkChart struct {
	Values *ChainlinkValuesWrapper `json:"values,omitempty"`
}

// ChainlinkValuesWrapper Chainlink values wrapper
type ChainlinkValuesWrapper struct {
	ChainlinkVals        *ChainlinkValues  `json:"chainlink,omitempty"`
	EnvironmentVariables map[string]string `json:"env,omitempty" yaml:"chainlink_env_values"`
}

// ChainlinkValues wraps all values
type ChainlinkValues struct {
	Image *ChainlinkImage `json:"image,omitempty"`
}

// ChainlinkImage defines chainlink image and version
type ChainlinkImage struct {
	Image   string `json:"image,omitempty" yaml:"chainlink_image"`
	Version string `json:"version,omitempty" yaml:"chainlink_version"`
}

func unmarshalYAML(path string, to interface{}) error {
	ap, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	log.Info().Str("Path", ap).Msg("Decoding config")
	f, err := ioutil.ReadFile(ap)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, to)
}
