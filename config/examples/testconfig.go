package examples

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	k8s_config "github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

type TestConfig struct {
	ChainlinkImage         *ctf_config.ChainlinkImageConfig `toml:"ChainlinkImage"`
	Logging                *ctf_config.LoggingConfig        `toml:"Logging"`
	Network                *ctf_config.NetworkConfig        `toml:"Network"`
	Pyroscope              *ctf_config.PyroscopeConfig      `toml:"Pyroscope"`
	PrivateEthereumNetwork *ctf_test_env.EthereumNetwork    `toml:"PrivateEthereumNetwork"`
}

func GetConfig(configurationName string, product string) (TestConfig, error) {
	logger := logging.GetTestLogger(nil)

	configurationName = strings.ReplaceAll(configurationName, "/", "_")
	configurationName = strings.ReplaceAll(configurationName, " ", "_")
	configurationName = cases.Title(language.English, cases.NoLower).String(configurationName)
	fileNames := []string{
		"default.toml",
		fmt.Sprintf("%s.toml", product),
		"overrides.toml",
	}

	testConfig := TestConfig{}
	maybeTestConfigs := []TestConfig{}

	logger.Debug().Msgf("Will apply configuration named '%s' if it is found in any of the configs", configurationName)

	for _, fileName := range fileNames {
		logger.Debug().Msgf("Looking for config file %s", fileName)
		filePath, err := osutil.FindFile(fileName, osutil.DEFAULT_STOP_FILE_NAME, 2)

		if err != nil && errors.Is(err, os.ErrNotExist) {
			logger.Debug().Msgf("Config file %s not found", fileName)
			continue
		}
		logger.Debug().Str("location", filePath).Msgf("Found config file %s", fileName)

		content, err := readFile(filePath)
		if err != nil {
			return TestConfig{}, fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		var readConfig TestConfig
		err = toml.Unmarshal(content, &readConfig)
		if err != nil {
			return TestConfig{}, fmt.Errorf("error unmarshaling config: %w", err)
		}

		logger.Debug().Msgf("Successfully unmarshalled config file %s", fileName)
		maybeTestConfigs = append(maybeTestConfigs, readConfig)

		var someToml map[string]interface{}
		err = toml.Unmarshal(content, &someToml)
		if err != nil {
			return TestConfig{}, err
		}

		if _, ok := someToml[configurationName]; !ok {
			logger.Debug().Msgf("Config file %s does not contain configuration named '%s', skipping.", fileName, configurationName)
			continue
		}

		marshalled, err := toml.Marshal(someToml[configurationName])
		if err != nil {
			return TestConfig{}, err
		}

		err = toml.Unmarshal(marshalled, &readConfig)
		if err != nil {
			return TestConfig{}, err
		}

		logger.Debug().Msgf("Configuration named '%s' read successfully.", configurationName)
		maybeTestConfigs = append(maybeTestConfigs, readConfig)
	}

	configEncoded, isSet := os.LookupEnv(k8s_config.EnvBase64ConfigOverride)
	if isSet && configEncoded != "" {
		decoded, err := base64.StdEncoding.DecodeString(configEncoded)
		if err != nil {
			return TestConfig{}, err
		}

		var base64override TestConfig
		err = toml.Unmarshal(decoded, &base64override)
		if err != nil {
			return TestConfig{}, fmt.Errorf("error unmarshaling base64 config: %w", err)
		}

		logger.Debug().Msgf("Applying base64 config override from environment variable %s", k8s_config.EnvBase64ConfigOverride)
		maybeTestConfigs = append(maybeTestConfigs, base64override)
	} else {
		logger.Debug().Msg("Base64 config override from environment variable not found")
	}

	// currently we need to read that kind of secrets only for network configuration
	testConfig.Network = &ctf_config.NetworkConfig{}
	err := testConfig.Network.ApplySecrets()
	if err != nil {
		return TestConfig{}, fmt.Errorf("error applying secrets to network config: %w", err)
	}

	for i := range maybeTestConfigs {
		err := testConfig.ApplyOverrides(&maybeTestConfigs[i])
		if err != nil {
			return TestConfig{}, fmt.Errorf("error applying overrides to test config: %w", err)
		}
	}

	err = testConfig.Validate()
	if err != nil {
		return TestConfig{}, fmt.Errorf("error validating test config: %w", err)
	}

	return testConfig, nil
}

func (c *TestConfig) ApplyOverrides(from *TestConfig) error {
	if from == nil {
		return nil
	}

	if from.ChainlinkImage != nil {
		if c.ChainlinkImage == nil {
			c.ChainlinkImage = from.ChainlinkImage
		} else {
			err := c.ChainlinkImage.ApplyOverrides(from.ChainlinkImage)
			if err != nil {
				return fmt.Errorf("error applying overrides to chainlink image config: %w", err)
			}
		}
	}

	if from.Logging != nil {
		if c.Logging == nil {
			c.Logging = from.Logging
		} else {
			err := c.Logging.ApplyOverrides(from.Logging)
			if err != nil {
				return fmt.Errorf("error applying overrides to logging config: %w", err)
			}
		}
	}

	if from.Network != nil {
		if c.Network == nil {
			c.Network = from.Network
		} else {
			err := c.Network.ApplyOverrides(from.Network)
			if err != nil {
				return fmt.Errorf("error applying overrides to network config: %w", err)
			}
		}
	}

	if from.Pyroscope != nil {
		if c.Pyroscope == nil {
			c.Pyroscope = from.Pyroscope
		} else {
			err := c.Pyroscope.ApplyOverrides(from.Pyroscope)
			if err != nil {
				return fmt.Errorf("error applying overrides to pyroscope config: %w", err)
			}
		}
	}

	if from.PrivateEthereumNetwork != nil {
		if c.PrivateEthereumNetwork == nil {
			c.PrivateEthereumNetwork = from.PrivateEthereumNetwork
		} else {
			err := c.PrivateEthereumNetwork.ApplyOverrides(from.PrivateEthereumNetwork)
			if err != nil {
				return fmt.Errorf("error applying overrides to private ethereum network config: %w", err)
			}
		}
		c.PrivateEthereumNetwork.EthereumChainConfig.GenerateGenesisTimestamp()
	}

	return nil
}

func (c *TestConfig) Validate() error {
	if c.ChainlinkImage == nil {
		return fmt.Errorf("chainlink image config must be set")
	}
	if err := c.ChainlinkImage.Validate(); err != nil {
		return fmt.Errorf("chainlink image config validation failed: %w", err)
	}
	if err := c.Network.Validate(); err != nil {
		return fmt.Errorf("network config validation failed: %w", err)
	}
	if c.Logging == nil {
		return errors.New("logging config must be set")
	}
	if err := c.Logging.Validate(); err != nil {
		return fmt.Errorf("logging config validation failed: %w", err)
	}
	if c.Pyroscope != nil {
		if err := c.Pyroscope.Validate(); err != nil {
			return fmt.Errorf("pyroscope config validation failed: %w", err)
		}
	}
	if c.PrivateEthereumNetwork != nil {
		if err := c.PrivateEthereumNetwork.Validate(); err != nil {
			return fmt.Errorf("private ethereum network config validation failed: %w", err)
		}
	}

	return nil
}

func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return content, nil
}
