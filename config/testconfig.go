package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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
func (c *TestConfig) ReadConfigValuesFromEnvVars() error {
	logger := logging.GetTestLogger(nil)

	if c.Network == nil {
		c.Network = &NetworkConfig{}
	}
	err := c.Network.LoadFromEnv()
	if err != nil {
		return errors.Wrap(err, "error loading network config from env")
	}

	chainlinkImage, err := readEnvVarValue("E2E_TEST_CHAINLINK_IMAGE", String)
	if err != nil {
		return err
	}
	if chainlinkImage != nil && chainlinkImage.(string) != "" {
		if c.ChainlinkImage == nil {
			c.ChainlinkImage = &ChainlinkImageConfig{}
		}
		image := chainlinkImage.(string)
		logger.Debug().Msgf("Using E2E_TEST_CHAINLINK_IMAGE env var to override ChainlinkImage.Image")
		c.ChainlinkImage.Image = &image
	}

	chainlinkUpgradeImage, err := readEnvVarValue("E2E_TEST_CHAINLINK_UPGRADE_IMAGE", String)
	if err != nil {
		return err
	}
	if chainlinkUpgradeImage != nil && chainlinkUpgradeImage.(string) != "" {
		if c.ChainlinkUpgradeImage == nil {
			c.ChainlinkUpgradeImage = &ChainlinkImageConfig{}
		}
		image := chainlinkUpgradeImage.(string)
		logger.Debug().Msgf("Using E2E_TEST_CHAINLINK_UPGRADE_IMAGE env var to override ChainlinkUpgradeImage.Image")
		c.ChainlinkUpgradeImage.Image = &image
	}

	if c.Logging == nil {
		c.Logging = &LoggingConfig{}
	}
	err = c.Logging.LoadFromEnv()
	if err != nil {
		return errors.Wrap(err, "error loading logging config from env")
	}
	if c.Pyroscope == nil {
		c.Pyroscope = &PyroscopeConfig{}
	}
	err = c.Pyroscope.LoadFromEnv()
	if err != nil {
		return errors.Wrap(err, "error loading pyroscope config from env")
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

// LoadSecretDotEnvFiles loads environment variables from .testsecrets files in /etc/e2etests and the user's home directory
func LoadSecretDotEnvFiles() error {
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
