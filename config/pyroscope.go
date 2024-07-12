package config

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/net"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

type PyroscopeConfig struct {
	Enabled     *bool   `toml:"enabled"`
	ServerUrl   *string `toml:"server_url"`
	Key         *string `toml:"key_secret"`
	Environment *string `toml:"environment"`
}

func (c *PyroscopeConfig) LoadFromEnv() error {
	logger := logging.GetTestLogger(nil)

	if c.ServerUrl == nil {
		serverUrl, err := readEnvVarValue("E2E_TEST_PYROSCOPE_SERVER_URL", String)
		if err != nil {
			return err
		}
		if serverUrl != nil && serverUrl.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_PYROSCOPE_SERVER_URL env var to override PyroscopeConfig.ServerUrl")
			c.ServerUrl = ptr.Ptr(serverUrl.(string))
		}
	}
	if c.Key == nil {
		key, err := readEnvVarValue("E2E_TEST_PYROSCOPE_KEY", String)
		if err != nil {
			return err
		}
		if key != nil && key.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_PYROSCOPE_KEY env var to override PyroscopeConfig.Key")
			c.Key = ptr.Ptr(key.(string))
		}
	}
	if c.Environment == nil {
		environment, err := readEnvVarValue("E2E_TEST_PYROSCOPE_ENVIRONMENT", String)
		if err != nil {
			return err
		}
		if environment != nil && environment.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_PYROSCOPE_ENVIRONMENT env var to override PyroscopeConfig.Environment")
			c.Environment = ptr.Ptr(environment.(string))
		}
	}
	if c.Enabled == nil {
		enabled, err := readEnvVarValue("E2E_TEST_PYROSCOPE_ENABLED", Boolean)
		if err != nil {
			return err
		}
		if enabled != nil {
			logger.Debug().Msg("Using E2E_TEST_PYROSCOPE_ENABLED env var to override PyroscopeConfig.Enabled")
			c.Enabled = ptr.Ptr(enabled.(bool))
		}
	}
	return nil
}

// Validate checks that the pyroscope config is valid, which means that
// server url, environment and key are set and non-empty, but only if
// pyroscope is enabled
func (c *PyroscopeConfig) Validate() error {
	if c.Enabled != nil && *c.Enabled {
		if c.ServerUrl == nil {
			return errors.New("pyroscope server url must be set")
		}
		if !net.IsValidURL(*c.ServerUrl) {
			return fmt.Errorf("invalid pyroscope server url %s", *c.ServerUrl)
		}
		if c.Environment == nil || *c.Environment == "" {
			return errors.New("pyroscope environment must be set")
		}
		if c.Key == nil || *c.Key == "" {
			return errors.New("pyroscope key must be set")
		}
	}

	return nil
}
