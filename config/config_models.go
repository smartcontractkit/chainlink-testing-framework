package config

import "time"

// FrameworkConfig common framework config
type FrameworkConfig struct {
	KeepEnvironments string         `mapstructure:"keep_environments" yaml:"keep_environments"`
	Logging          *LoggingConfig `mapstructure:"logging" yaml:"logging"`
	EnvironmentFile  string         `mapstructure:"environment_file" yaml:"environment_file"`
	ChainlinkImage   string         `mapstructure:"chainlink_image" yaml:"chainlink_image"`
	ChainlinkVersion string         `mapstructure:"chainlink_version" yaml:"chainlink_version"`
	// ChainlinkEnvValues uses interface{} as the value because it's needed for proper helmchart merges
	ChainlinkEnvValues map[string]interface{} `mapstructure:"chainlink_env_values" yaml:"chainlink_env_values"`
}

// ETHNetwork data to configure fully ETH compatible network
type ETHNetwork struct {
	ContractsDeployed         bool          `mapstructure:"contracts_deployed" yaml:"contracts_deployed"`
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

// NetworksConfig is network configurations
type NetworksConfig struct {
	SelectedNetworks   []string        `mapstructure:"selected_networks" yaml:"selected_networks" envconfig:"selected_networks"`
	NetworkSettings    NetworkSettings `mapstructure:"networks" yaml:"networks" envconfig:"network_settings"`
	DefaultKeyStore    string
	ConfigFileLocation string
}

// LoggingConfig for logging
type LoggingConfig struct {
	WritePodLogs string `mapstructure:"write_pod_logs" yaml:"write_pod_logs"`
	Level        int8   `mapstructure:"level" yaml:"level"`
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

// RemoteRunnerConfig reads the config file for remote test runs
type RemoteRunnerConfig struct {
	TestRegex     string `mapstructure:"test_regex" yaml:"test_regex"`
	TestDirectory string `mapstructure:"test_directory" yaml:"test_directory"`
	SlackAPIKey   string `mapstructure:"slack_api_key" yaml:"slack_api_key"`
	SlackChannel  string `mapstructure:"slack_channel" yaml:"slack_channel"`
	SlackUserID   string `mapstructure:"slack_user_id" yaml:"slack_user_id"`
}
