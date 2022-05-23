package config

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// Config general framework config
type Config struct {
	RemoteRunnerConfig *RemoteRunnerConfig `envconfig:"REMOTE_RUNNER_CONFIG_FILE" required:"false" default:"../remote_runner_config.yaml"`
}

// ETHNetwork data to configure fully ETH compatible network
type ETHNetwork struct {
	ContractsDeployed         bool          `envconfig:"contracts_deployed" yaml:"contracts_deployed"`
	Name                      string        `envconfig:"name" yaml:"name"`
	ChainID                   int64         `envconfig:"chain_id" yaml:"chain_id"`
	URL                       string        `envconfig:"url" yaml:"url"`
	URLs                      []string      `envconfig:"urls" yaml:"urls"`
	Type                      string        `envconfig:"type" yaml:"type"`
	PrivateKeys               []string      `envconfig:"private_keys" yaml:"private_keys"`
	ChainlinkTransactionLimit uint64        `envconfig:"chainlink_transaction_limit" yaml:"chainlink_transaction_limit"`
	Timeout                   time.Duration `envconfig:"transaction_timeout" yaml:"transaction_timeout"`
	MinimumConfirmations      int           `envconfig:"minimum_confirmations" yaml:"minimum_confirmations"`
	GasEstimationBuffer       uint64        `envconfig:"gas_estimation_buffer" yaml:"gas_estimation_buffer"`
	BlockGasLimit             uint64        `envconfig:"block_gas_limit" yaml:"block_gas_limit"`
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
