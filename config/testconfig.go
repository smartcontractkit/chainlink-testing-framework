package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
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
func (c *TestConfig) ReadFromEnvVar() error {
	logger := logging.GetTestLogger(nil)

	lokiTenantID := MustReadEnvVar_String(E2E_TEST_LOKI_TENANT_ID_ENV)
	if lokiTenantID != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.TenantId", E2E_TEST_LOKI_TENANT_ID_ENV)
		c.Logging.Loki.TenantId = &lokiTenantID
	}

	lokiEndpoint := MustReadEnvVar_String(E2E_TEST_LOKI_ENDPOINT_ENV)
	if lokiEndpoint != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.Endpoint", E2E_TEST_LOKI_ENDPOINT_ENV)
		c.Logging.Loki.Endpoint = &lokiEndpoint
	}

	lokiBasicAuth := MustReadEnvVar_String(E2E_TEST_LOKI_BASIC_AUTH_ENV)
	if lokiBasicAuth != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.BasicAuth", E2E_TEST_LOKI_BASIC_AUTH_ENV)
		c.Logging.Loki.BasicAuth = &lokiBasicAuth
	}

	lokiBearerToken := MustReadEnvVar_String(E2E_TEST_LOKI_BEARER_TOKEN_ENV)
	if lokiBearerToken != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Loki == nil {
			c.Logging.Loki = &LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.BearerToken", E2E_TEST_LOKI_BEARER_TOKEN_ENV)
		c.Logging.Loki.BearerToken = &lokiBearerToken
	}

	grafanaBaseUrl := MustReadEnvVar_String(E2E_TEST_GRAFANA_BASE_URL_ENV)
	if grafanaBaseUrl != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Grafana == nil {
			c.Logging.Grafana = &GrafanaConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Grafana.BaseUrl", E2E_TEST_GRAFANA_BASE_URL_ENV)
		c.Logging.Grafana.BaseUrl = &grafanaBaseUrl
	}

	grafanaDashboardUrl := MustReadEnvVar_String(E2E_TEST_GRAFANA_DASHBOARD_URL_ENV)
	if grafanaDashboardUrl != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Grafana == nil {
			c.Logging.Grafana = &GrafanaConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Grafana.DashboardUrl", E2E_TEST_GRAFANA_DASHBOARD_URL_ENV)
		c.Logging.Grafana.DashboardUrl = &grafanaDashboardUrl
	}

	grafanaBearerToken := MustReadEnvVar_String(E2E_TEST_GRAFANA_BEARER_TOKEN_ENV)
	if grafanaBearerToken != "" {
		if c.Logging == nil {
			c.Logging = &LoggingConfig{}
		}
		if c.Logging.Grafana == nil {
			c.Logging.Grafana = &GrafanaConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Grafana.BearerToken", E2E_TEST_GRAFANA_BEARER_TOKEN_ENV)
		c.Logging.Grafana.BearerToken = &grafanaBearerToken
	}

	pyroscopeEnabled := MustReadEnvVar_Boolean(E2E_TEST_PYROSCOPE_ENABLED_ENV)
	if pyroscopeEnabled != nil {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Pyroscope.Enabled", E2E_TEST_PYROSCOPE_ENABLED_ENV)
		c.Pyroscope.Enabled = pyroscopeEnabled
	}

	pyroscopeServerUrl := MustReadEnvVar_String(E2E_TEST_PYROSCOPE_SERVER_URL_ENV)
	if pyroscopeServerUrl != "" {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Pyroscope.ServerUrl", E2E_TEST_PYROSCOPE_SERVER_URL_ENV)
		c.Pyroscope.ServerUrl = &pyroscopeServerUrl
	}

	pyroscopeKey := MustReadEnvVar_String(E2E_TEST_PYROSCOPE_KEY_ENV)
	if pyroscopeKey != "" {
		if c.Pyroscope == nil {
			c.Pyroscope = &PyroscopeConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Pyroscope.Key", E2E_TEST_PYROSCOPE_KEY_ENV)
		c.Pyroscope.Key = &pyroscopeKey
	}

	walletKeys := ReadEnvVarGroupedMap(E2E_TEST_WALLET_KEY_ENV, E2E_TEST_WALLET_KEYS_ENV)
	if len(walletKeys) > 0 {
		if c.Network == nil {
			c.Network = &NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s and/or %s env vars to override Network.WalletKeys", E2E_TEST_WALLET_KEY_ENV, E2E_TEST_WALLET_KEYS_ENV)
		c.Network.WalletKeys = walletKeys
	}

	rpcHttpUrls := ReadEnvVarGroupedMap(E2E_TEST_RPC_HTTP_URL_ENV, E2E_TEST_RPC_HTTP_URLS_ENV)
	if len(rpcHttpUrls) > 0 {
		if c.Network == nil {
			c.Network = &NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s and/or %s env vars to override Network.RpcHttpUrls", E2E_TEST_RPC_HTTP_URL_ENV, E2E_TEST_RPC_HTTP_URLS_ENV)
		c.Network.RpcHttpUrls = rpcHttpUrls
	}

	rpcWsUrls := ReadEnvVarGroupedMap(E2E_TEST_RPC_WS_URL_ENV, E2E_TEST_RPC_WS_URLS_ENV)
	if len(rpcWsUrls) > 0 {
		if c.Network == nil {
			c.Network = &NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s and/or %s env vars to override Network.RpcWsUrls", E2E_TEST_RPC_WS_URL_ENV, E2E_TEST_RPC_WS_URLS_ENV)
		c.Network.RpcWsUrls = rpcWsUrls
	}

	chainlinkImage := MustReadEnvVar_String(E2E_TEST_CHAINLINK_IMAGE_ENV)
	if chainlinkImage != "" {
		if c.ChainlinkImage == nil {
			c.ChainlinkImage = &ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override ChainlinkImage.Image", E2E_TEST_CHAINLINK_IMAGE_ENV)
		c.ChainlinkImage.Image = &chainlinkImage
	}

	chainlinkUpgradeImage := MustReadEnvVar_String(E2E_TEST_CHAINLINK_UPGRADE_IMAGE_ENV)
	if chainlinkUpgradeImage != "" {
		if c.ChainlinkUpgradeImage == nil {
			c.ChainlinkUpgradeImage = &ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override ChainlinkUpgradeImage.Image", E2E_TEST_CHAINLINK_UPGRADE_IMAGE_ENV)
		c.ChainlinkUpgradeImage.Image = &chainlinkUpgradeImage
	}

	return nil
}

func LoadSecretEnvsFromFiles() error {
	logger := logging.GetTestLogger(nil)

	// Load existing environment variables into a map
	existingEnv := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		existingEnv[pair[0]] = pair[1]
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrapf(err, "error getting user home directory")
	}
	homePath := fmt.Sprintf("%s/.testsecrets", homeDir)
	etcPath := "/etc/e2etests/.testsecrets"
	testsecretsPath := []string{etcPath, homePath}

	for _, path := range testsecretsPath {
		logger.Debug().Msgf("Checking for test secrets file at %s", path)

		// Load variables from the env file
		envMap, err := godotenv.Read(path)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Debug().Msgf("No test secrets file found at %s", path)
				continue
			}
			return errors.Wrapf(err, "error reading test secrets file at %s", path)
		}

		// Set env vars from file only if they are not already set
		for key, value := range envMap {
			if _, exists := existingEnv[key]; !exists {
				logger.Debug().Msgf("Setting env var %s from %s file", key, path)
				os.Setenv(key, value)
			} else {
				logger.Debug().Msgf("Env var %s already set, not overriding it from %s file", key, path)
			}
		}
	}

	return nil
}
