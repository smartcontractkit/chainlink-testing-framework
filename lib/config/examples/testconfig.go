package testconfig

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/barkimedes/go-deepcopy"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	k8s_config "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
)

type TestConfig struct {
	ChainlinkImage         *ctf_config.ChainlinkImageConfig `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ctf_config.ChainlinkImageConfig `toml:"ChainlinkUpgradeImage"`
	Logging                *ctf_config.LoggingConfig        `toml:"Logging"`
	Network                *ctf_config.NetworkConfig        `toml:"Network"`
	Pyroscope              *ctf_config.PyroscopeConfig      `toml:"Pyroscope"`
	PrivateEthereumNetwork *ctf_test_env.EthereumNetwork    `toml:"PrivateEthereumNetwork"`
	ConfigurationName      string                           `toml:"-"`
}

// Saves Test Config to a local file
func (c *TestConfig) Save() (string, error) {
	filePath := fmt.Sprintf("test_config-%s.toml", uuid.New())

	content, err := toml.Marshal(*c)
	if err != nil {
		return "", errors.Wrapf(err, "error marshaling test config")
	}

	err = os.WriteFile(filePath, content, 0600)
	if err != nil {
		return "", errors.Wrapf(err, "error writing test config")
	}

	return filePath, nil
}

// Returns a deep copy of the Test Config or panics on error
func (c TestConfig) MustCopy() TestConfig {
	return deepcopy.MustAnything(c).(TestConfig)
}

type Common struct {
	ChainlinkNodeFunding *float64 `toml:"chainlink_node_funding"`
}

func (c *Common) Validate() error {
	if c.ChainlinkNodeFunding != nil && *c.ChainlinkNodeFunding < 0 {
		return fmt.Errorf("chainlink node funding must be positive")
	}

	return nil
}

type Product string

const (
	Automation    Product = "automation"
	Cron          Product = "cron"
	DirectRequest Product = "direct_request"
	Flux          Product = "flux"
	ForwarderOcr  Product = "forwarder_ocr"
	ForwarderOcr2 Product = "forwarder_ocr2"
	Functions     Product = "functions"
	Keeper        Product = "keeper"
	LogPoller     Product = "log_poller"
	Node          Product = "node"
	OCR           Product = "ocr"
	OCR2          Product = "ocr2"
	OCR2VRF       Product = "ocr2vrf"
	RunLog        Product = "runlog"
	VRF           Product = "vrf"
	VRFv2         Product = "vrfv2"
	VRFv2Plus     Product = "vrfv2plus"
)

const TestTypeEnvVarName = "TEST_TYPE"

func GetConfigurationNameFromEnv() (string, error) {
	testType := os.Getenv(TestTypeEnvVarName)
	if testType == "" {
		return "", fmt.Errorf("%s env var not set", TestTypeEnvVarName)
	}

	return cases.Title(language.English, cases.NoLower).String(testType), nil
}

const (
	Base64OverrideEnvVarName = k8s_config.EnvBase64ConfigOverride
	NoKey                    = "NO_KEY"
)

func GetConfig(configurationName string, product Product) (TestConfig, error) {
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
	testConfig.ConfigurationName = configurationName
	logger.Debug().Msgf("Will apply configuration named '%s' if it is found in any of the configs", configurationName)

	logger.Info().Msg("Reading and applying configs from file system")
	for _, fileName := range fileNames {
		logger.Debug().Msgf("Looking for config file %s", fileName)
		filePath, err := osutil.FindFile(fileName, osutil.DEFAULT_STOP_FILE_NAME, 2)

		if err != nil && errors.Is(err, os.ErrNotExist) {
			logger.Debug().Msgf("Config file %s not found", fileName)
			continue
		} else if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error looking for file %s", filePath)
		}
		logger.Debug().Str("location", filePath).Msgf("Found config file %s", fileName)

		content, err := readFile(filePath)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error reading file %s", filePath)
		}

		err = ctf_config.BytesToAnyTomlStruct(logger, fileName, configurationName, &testConfig, content)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error reading file %s", filePath)
		}
	}

	configEncoded, isSet := os.LookupEnv(Base64OverrideEnvVarName)
	if isSet && configEncoded != "" {
		logger.Debug().Msgf("Base64 config override from environment variable '%s' found", Base64OverrideEnvVarName)
		decoded, err := base64.StdEncoding.DecodeString(configEncoded)
		if err != nil {
			return TestConfig{}, err
		}

		err = toml.Unmarshal(decoded, &testConfig)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error unmarshalling base64 config")
		}
	} else {
		logger.Debug().Msg("Base64 config override from environment variable not found")
	}

	// it needs some custom logic, so we do it separately
	err := testConfig.readNetworkConfiguration()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error reading network config")
	}

	logger.Debug().Msg("Validating test config")
	err = testConfig.Validate()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error validating test config")
	}

	logger.Debug().Msg("Correct test config constructed successfully")
	return testConfig, nil
}

func (c *TestConfig) readNetworkConfiguration() error {
	// currently we need to read that kind of secrets only for network configuration
	if c.Network == nil {
		c.Network = &ctf_config.NetworkConfig{}
	}

	c.Network.UpperCaseNetworkNames()
	c.Network.OverrideURLsAndKeysFromEVMNetwork()

	// this is the only value we need to generate dynamically before starting a new simulated chain
	if c.PrivateEthereumNetwork != nil && c.PrivateEthereumNetwork.EthereumChainConfig != nil {
		c.PrivateEthereumNetwork.EthereumChainConfig.GenerateGenesisTimestamp()
	}

	return nil
}

func (c *TestConfig) Validate() error {
	if c.ChainlinkImage == nil {
		return fmt.Errorf("chainlink image config must be set")
	}
	if err := c.ChainlinkImage.Validate(); err != nil {
		return errors.Wrapf(err, "chainlink image config validation failed")
	}
	if c.ChainlinkUpgradeImage != nil {
		if err := c.ChainlinkUpgradeImage.Validate(); err != nil {
			return errors.Wrapf(err, "chainlink upgrade image config validation failed")
		}
	}
	if err := c.Network.Validate(); err != nil {
		return errors.Wrapf(err, "network config validation failed")
	}

	if c.Logging == nil {
		return fmt.Errorf("logging config must be set")
	}

	if err := c.Logging.Validate(); err != nil {
		return errors.Wrapf(err, "logging config validation failed")
	}

	if c.Pyroscope != nil {
		if err := c.Pyroscope.Validate(); err != nil {
			return errors.Wrapf(err, "pyroscope config validation failed")
		}
	}

	if c.PrivateEthereumNetwork != nil {
		if err := c.PrivateEthereumNetwork.Validate(); err != nil {
			return errors.Wrapf(err, "private ethereum network config validation failed")
		}
	}

	return nil
}

func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file %s", filePath)
	}

	return content, nil
}
