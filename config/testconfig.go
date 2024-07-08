package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/seth"
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

// Read config values from environment variables
func (c *TestConfig) ReadConfigValuesFromEnvVars() error {
	logger := logging.GetTestLogger(nil)

	walletKeys := mergeMaps(loadEnvVarSingleMap(`(.+)_WALLET_KEY$`), loadEnvVarGroupedMap(`(.+)_WALLET_KEY_(\d+)$`))
	if len(walletKeys) > 0 {
		if c.Network == nil {
			c.Network = &NetworkConfig{}
		}
		c.Network.WalletKeys = walletKeys
	}
	rpcHttpUrls := mergeMaps(loadEnvVarSingleMap(`(.+)_RPC_HTTP_URL$`), loadEnvVarGroupedMap(`(.+)_RPC_HTTP_URL_(\d+)$`))
	if len(rpcHttpUrls) > 0 {
		if c.Network == nil {
			c.Network = &NetworkConfig{}
		}
		c.Network.RpcHttpUrls = rpcHttpUrls
	}
	rpcWsUrls := mergeMaps(loadEnvVarSingleMap(`(.+)_RPC_WS_URL$`), loadEnvVarGroupedMap(`(.+)_RPC_WS_URL_(\d+)$`))
	if len(rpcWsUrls) > 0 {
		if c.Network == nil {
			c.Network = &NetworkConfig{}
		}
		c.Network.RpcWsUrls = rpcWsUrls
	}

	chainlinkImage, err := readEnvVarValue("CHAINLINK_IMAGE", String)
	if err != nil {
		return err
	}
	if chainlinkImage != nil && chainlinkImage.(string) != "" {
		if c.ChainlinkImage == nil {
			c.ChainlinkImage = &ChainlinkImageConfig{}
		}
		image := chainlinkImage.(string)
		logger.Debug().Msgf("Using CHAINLINK_IMAGE env var to override ChainlinkImage.Image")
		c.ChainlinkImage.Image = &image
	}

	chainlinkUpgradeImage, err := readEnvVarValue("CHAINLINK_UPGRADE_IMAGE", String)
	if err != nil {
		return err
	}
	if chainlinkUpgradeImage != nil && chainlinkUpgradeImage.(string) != "" {
		if c.ChainlinkUpgradeImage == nil {
			c.ChainlinkUpgradeImage = &ChainlinkImageConfig{}
		}
		image := chainlinkUpgradeImage.(string)
		logger.Debug().Msgf("Using CHAINLINK_UPGRADE_IMAGE env var to override ChainlinkUpgradeImage.Image")
		c.ChainlinkUpgradeImage.Image = &image
	}

	lokiTenantID, err := readEnvVarValue("LOKI_TENANT_ID", String)
	if err != nil {
		return err
	}
	if lokiTenantID != nil && lokiTenantID.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		id := lokiTenantID.(string)
		logger.Debug().Msgf("Using LOKI_TENANT_ID env var to override Logging.Loki.TenantId")
		c.Logging.Loki.TenantId = &id
	}

	lokiEndpoint, err := readEnvVarValue("LOKI_ENDPOINT", String)
	if err != nil {
		return err
	}
	if lokiEndpoint != nil && lokiEndpoint.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		endpoint := lokiEndpoint.(string)
		logger.Debug().Msgf("Using LOKI_ENDPOINT env var to override Logging.Loki.Endpoint")
		c.Logging.Loki.Endpoint = &endpoint
	}

	lokiBasicAuth, err := readEnvVarValue("LOKI_BASIC_AUTH", String)
	if err != nil {
		return err
	}
	if lokiBasicAuth != nil && lokiBasicAuth.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		basicAuth := lokiBasicAuth.(string)
		logger.Debug().Msgf("Using LOKI_BASIC_AUTH env var to override Logging.Loki.BasicAuth")
		c.Logging.Loki.BasicAuth = &basicAuth
	}

	lokiBearerToken, err := readEnvVarValue("LOKI_BEARER_TOKEN", String)
	if err != nil {
		return err
	}
	if lokiBearerToken != nil && lokiBearerToken.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		bearerToken := lokiBearerToken.(string)
		logger.Debug().Msgf("Using LOKI_BEARER_TOKEN env var to override Logging.Loki.BearerToken")
		c.Logging.Loki.BearerToken = &bearerToken
	}

	grafanaBaseUrl, err := readEnvVarValue("GRAFANA_BASE_URL", String)
	if err != nil {
		return err
	}
	if grafanaBaseUrl != nil && grafanaBaseUrl.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Grafana == nil {
			c.Logging.Grafana = &GrafanaConfig{}
		}
		baseUrl := grafanaBaseUrl.(string)
		logger.Debug().Msgf("Using GRAFANA_BASE_URL env var to override Logging.Grafana.BaseUrl")
		c.Logging.Grafana.BaseUrl = &baseUrl
	}

	grafanaDashboardUrl, err := readEnvVarValue("GRAFANA_DASHBOARD_URL", String)
	if err != nil {
		return err
	}
	if grafanaDashboardUrl != nil && grafanaDashboardUrl.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Grafana == nil {
			c.Logging.Grafana = &GrafanaConfig{}
		}
		dashboardUrl := grafanaDashboardUrl.(string)
		logger.Debug().Msgf("Using GRAFANA_DASHBOARD_URL env var to override Logging.Grafana.DashboardUrl")
		c.Logging.Grafana.DashboardUrl = &dashboardUrl
	}

	grafanaBearerToken, err := readEnvVarValue("GRAFANA_BEARER_TOKEN", String)
	if err != nil {
		return err
	}
	if grafanaBearerToken != nil && grafanaBearerToken.(string) != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Grafana == nil {
			c.Logging.Grafana = &GrafanaConfig{}
		}
		bearerToken := grafanaBearerToken.(string)
		logger.Debug().Msgf("Using GRAFANA_BEARER_TOKEN env var to override Logging.Grafana.BearerToken")
		c.Logging.Grafana.BearerToken = &bearerToken
	}

	pyroscopeServerUrl, err := readEnvVarValue("PYROSCOPE_SERVER_URL", String)
	if err != nil {
		return err
	}
	if pyroscopeServerUrl != nil && pyroscopeServerUrl.(string) != "" {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		serverUrl := pyroscopeServerUrl.(string)
		logger.Debug().Msgf("Using PYROSCOPE_SERVER_URL env var to override Pyroscope.ServerUrl")
		c.Pyroscope.ServerUrl = &serverUrl
	}

	pyroscopeKey, err := readEnvVarValue("PYROSCOPE_KEY", String)
	if err != nil {
		return err
	}
	if pyroscopeKey != nil && pyroscopeKey.(string) != "" {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		key := pyroscopeKey.(string)
		logger.Debug().Msgf("Using PYROSCOPE_KEY env var to override Pyroscope.Key")
		c.Pyroscope.Key = &key
	}

	pyroscopeEnvironment, err := readEnvVarValue("PYROSCOPE_ENVIRONMENT", String)
	if err != nil {
		return err
	}
	if pyroscopeEnvironment != nil && pyroscopeEnvironment.(string) != "" {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		environment := pyroscopeEnvironment.(string)
		logger.Debug().Msgf("Using PYROSCOPE_ENVIRONMENT env var to override Pyroscope.Environment")
		c.Pyroscope.Environment = &environment
	}

	pyroscopeEnabled, err := readEnvVarValue("PYROSCOPE_ENABLED", Boolean)
	if err != nil {
		return err
	}
	if pyroscopeEnabled != nil {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		enabled := pyroscopeEnabled.(bool)
		logger.Debug().Msgf("Using PYROSCOPE_ENABLED env var to override Pyroscope.Enabled")
		c.Pyroscope.Enabled = &enabled
	}

	return nil
}

// loadEnvVarGroupedMap scans all environment variables, matches them against
// a specified pattern, and returns a map of grouped values based on the pattern.
// The grouping is defined by the first capture group of the regex.
func loadEnvVarGroupedMap(pattern string) map[string][]string {
	logger := logging.GetTestLogger(nil)
	re := regexp.MustCompile(pattern)
	groupedVars := make(map[string][]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 && value != "" {
			group := matches[1] // Use the first capture group for grouping
			groupedVars[group] = append(groupedVars[group], value)
			logger.Debug().Msgf("Will override test config from env var '%s'", key)
		}
	}
	return groupedVars
}

func loadEnvVarSingleMap(pattern string) map[string]string {
	logger := logging.GetTestLogger(nil)
	re := regexp.MustCompile(pattern)
	singleVars := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 && value != "" {
			group := matches[1] // Use the first capture group for grouping
			singleVars[group] = value
			logger.Debug().Msgf("Will override test config from env var '%s'", key)
		}
	}
	return singleVars
}

// Merges a map[string]string with a map[string][]string and returns a new map[string][]string.
// Elements from the single map are inserted at index 0 in the slice of the new map.
func mergeMaps(single map[string]string, multi map[string][]string) map[string][]string {
	newMap := make(map[string][]string)

	// First, copy all elements from the multi map to the new map
	for key, values := range multi {
		newMap[key] = make([]string, len(values))
		copy(newMap[key], values)
	}

	// Next, insert or prepend the elements from the single map
	for key, value := range single {
		if existingValues, exists := newMap[key]; exists {
			// Prepend the value from the single map
			newMap[key] = append([]string{value}, existingValues...)
		} else {
			// Initialize a new slice if the key does not exist
			newMap[key] = []string{value}
		}
	}

	return newMap
}

type EnvValueType int

const (
	String EnvValueType = iota
	Integer
	Boolean
	Float
)

// readEnvVarValue reads an environment variable and returns the value parsed according to the specified type.
func readEnvVarValue(envVarName string, valueType EnvValueType) (interface{}, error) {
	// Get the environment variable value
	value, isSet := os.LookupEnv(envVarName)
	if !isSet {
		return nil, nil
	}
	if isSet && value == "" {
		return "", nil // Return "" if the environment variable is not set
	}

	// Parse the value according to the specified type
	switch valueType {
	case Integer:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("error converting value to integer: %v", err)
		}
		return intVal, nil
	case Boolean:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("error converting value to boolean: %v", err)
		}
		return boolVal, nil
	case Float:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting value to float: %v", err)
		}
		return floatVal, nil
	default: // String or unrecognized type
		return value, nil
	}
}
