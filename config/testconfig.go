package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/seth"
	"github.com/spf13/viper"
)

func (c *TestConfig) GetLoggingConfig() *LoggingConfig {
	return c.Logging
}

func (c *TestConfig) GetNodeConfig() *NodeConfig {
	return c.NodeConfig
}

func (c TestConfig) GetNetworkConfig() *NetworkConfig {
	return c.Network
}

func (c TestConfig) GetChainlinkImageConfig() *ChainlinkImageConfig {
	return c.ChainlinkImage
}

func (c TestConfig) GetPrivateEthereumNetworkConfig() *EthereumNetworkConfig {
	return c.PrivateEthereumNetwork
}

func (c TestConfig) GetPyroscopeConfig() *PyroscopeConfig {
	return c.Pyroscope
}

type TestConfig struct {
	ChainlinkImage         *ChainlinkImageConfig  `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ChainlinkImageConfig  `toml:"ChainlinkUpgradeImage"`
	Logging                *LoggingConfig         `toml:"Logging"`
	Network                *NetworkConfig         `toml:"Network"`
	Pyroscope              *PyroscopeConfig       `toml:"Pyroscope"`
	PrivateEthereumNetwork *EthereumNetworkConfig `toml:"PrivateEthereumNetwork"`
	WaspConfig             *WaspAutoBuildConfig   `toml:"WaspAutoBuild"`
	Seth                   *seth.Config           `toml:"Seth"`
	NodeConfig             *NodeConfig            `toml:"NodeConfig"`
}

var viperLock sync.Mutex

// Read config values from environment variables
func (c *TestConfig) ReadConfigValuesFromEnvVars() error {
	err := c.readEnvVarGroups(
		"Network.WalletKeys",
		`TEST_CONFIG_(.+)_WALLET_KEY_(\d+)$`,
	)
	if err != nil {
		return err
	}
	err = c.readEnvVarGroups(
		"Network.RpcHttpUrls",
		`TEST_CONFIG_(.+)_RPC_HTTP_URL_(\d+)$`,
	)
	if err != nil {
		return err
	}
	err = c.readEnvVarGroups(
		"Network.RpcWsUrls",
		`TEST_CONFIG_(.+)_RPC_WS_URL_(\d+)$`,
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"ChainlinkImage.Image",
		"TEST_CONFIG_CHAINLINK_IMAGE",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"ChainlinkUpgradeImage.Image",
		"TEST_CONFIG_CHAINLINK_UPGRADE_IMAGE",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Loki.TenantId",
		"TEST_CONFIG_LOKI_TENANT_ID",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Loki.Endpoint",
		"TEST_CONFIG_LOKI_ENDPOINT",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Loki.BasicAuth",
		"TEST_CONFIG_LOKI_BASIC_AUTH",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Loki.BearerToken",
		"TEST_CONFIG_LOKI_BEARER_TOKEN",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Grafana.BaseUrl",
		"TEST_CONFIG_GRAFANA_BASE_URL",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Grafana.DashboardUrl",
		"TEST_CONFIG_GRAFANA_DASHBOARD_URL",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Logging.Grafana.BearerToken",
		"TEST_CONFIG_GRAFANA_BEARER_TOKEN",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Pyroscope.ServerUrl",
		"TEST_CONFIG_PYROSCOPE_SERVER_URL",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Pyroscope.Key",
		"TEST_CONFIG_PYROSCOPE_KEY",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Pyroscope.Environment",
		"TEST_CONFIG_PYROSCOPE_ENVIRONMENT",
		EnvValueType(String),
	)
	if err != nil {
		return err
	}
	err = c.readSingleEnvVar(
		"Pyroscope.Enabled",
		"TEST_CONFIG_PYROSCOPE_ENABLED",
		EnvValueType(Boolean),
	)
	if err != nil {
		return err
	}
	return nil
}

// Read env vars for map[string][]string
func (c *TestConfig) readEnvVarGroups(key string, regexpStr string) error {
	logger := logging.GetTestLogger(nil)
	re := regexp.MustCompile(regexpStr)

	envVars := make(map[string][]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]

		matches := re.FindStringSubmatch(key)
		if matches != nil {
			key := matches[1]
			envVars[key] = append(envVars[key], value)
		}
	}

	if len(envVars) == 0 {
		logger.Debug().Msgf("Not setting test config key '%s' because environment variables that match '%s' regex not found", key, regexpStr)
		return nil
	}

	for network, keys := range envVars {
		keyPath := fmt.Sprintf("%s.%s", key, network)
		viperLock.Lock()
		viper.Set(keyPath, keys)
		viperLock.Unlock()
		logger.Debug().Msgf("Setting test config key '%s' from env var", keyPath)
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("error reading test config values from environment variables. Unable to unmarshal config: %v", err)
	}
	return nil
}

type EnvValueType int

const (
	String EnvValueType = iota
	Integer
	Boolean
	Float
)

// Read env var for single value
func (c *TestConfig) readSingleEnvVar(key, envVarName string, valueType EnvValueType) error {
	logger := logging.GetTestLogger(nil)

	// Get the environment variable value
	value := os.Getenv(envVarName)
	if value == "" {
		logger.Debug().Msgf("Not setting test config key '%s' because environment variable '%s' not found", key, envVarName)
		return nil
	}

	viperLock.Lock()
	defer viperLock.Unlock()

	// Parse the value according to the specified type and set it in Viper
	switch valueType {
	case Integer:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("error converting value to integer: %v", err)
		}
		viper.Set(key, intVal)
	case Boolean:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("error converting value to boolean: %v", err)
		}
		viper.Set(key, boolVal)
	case Float:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("error converting value to float: %v", err)
		}
		viper.Set(key, floatVal)
	default: // String or unrecognized type
		viper.Set(key, value)
	}

	logger.Debug().Msgf("Setting test config key '%s' from environment variable '%s'", key, envVarName)

	// Unmarshal the configuration into the TestConfig struct
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("error reading test config values from environment variables. Unable to unmarshal config: %v", err)
	}

	return nil
}
