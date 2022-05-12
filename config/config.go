// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ConfigurationType refers to the different ways that configurations can be set
type ConfigurationType string

// Configs
const (
	LocalConfig  ConfigurationType = "local"
	SecretConfig ConfigurationType = "secret"

	DefaultGeth     string = "geth"
	PerformanceGeth string = "geth_performance"
	RealisticGeth   string = "geth_realistic"
)

// NetworkSettings is a map that holds configuration for each individual network
type NetworkSettings map[string]map[string]interface{}

var ProjectFrameworkSettings *FrameworkConfig
var ProjectNetworkSettings *NetworksConfig
var ProjectConfigDirectory string

// ChainlinkVals formats Chainlink values set in the framework config to be passed to Chainlink deployments
func ChainlinkVals() map[string]interface{} {
	if ProjectFrameworkSettings == nil {
		log.Error().Msg("ProjectFrameworkSettings not set!")
		return nil
	}
	values := map[string]interface{}{}
	if len(ProjectFrameworkSettings.ChainlinkEnvValues) > 0 {
		values["env"] = ProjectFrameworkSettings.ChainlinkEnvValues
	}
	if ProjectFrameworkSettings.ChainlinkImage != "" {
		values["chainlink"] = map[string]interface{}{
			"image": map[string]interface{}{
				"image":   ProjectFrameworkSettings.ChainlinkImage,
				"version": ProjectFrameworkSettings.ChainlinkVersion,
			},
		}
	}
	return values
}

// GethNetworks builds the proper geth network settings to use based on the selected_networks config
func GethNetworks() []environment.SimulatedNetwork {
	if ProjectNetworkSettings == nil {
		log.Error().Msg("ProjectNetworkSettings not set!")
		return nil
	}
	var gethNetworks []environment.SimulatedNetwork
	for _, network := range ProjectNetworkSettings.SelectedNetworks {
		switch network {
		case DefaultGeth:
			gethNetworks = append(gethNetworks, environment.DefaultGeth)
		case PerformanceGeth:
			gethNetworks = append(gethNetworks, environment.PerformanceGeth)
		case RealisticGeth:
			gethNetworks = append(gethNetworks, environment.RealisticGeth)
		}
	}
	return gethNetworks
}

// Decode is used by envconfig to initialize the custom Charts type with populated values
// This function will take a JSON object representing charts, and unmarshal it into the existing object to "merge" the
// two
func (n NetworkSettings) Decode(value string) error {
	// Support the use of files for unmarshaling charts JSON
	if _, err := os.Stat(value); err == nil {
		b, err := os.ReadFile(value)
		if err != nil {
			return err
		}
		value = string(b)
	}
	networkSettings := NetworkSettings{}
	if err := yaml.Unmarshal([]byte(value), &networkSettings); err != nil {
		return fmt.Errorf("failed to unmarshal YAML, either a file path specific doesn't exist, or the YAML is invalid: %v", err)
	}
	return mergo.Merge(&n, networkSettings, mergo.WithOverride)
}

func defaultViper(dir string, file string) *viper.Viper {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName(file)
	v.SetConfigType("yaml")

	v.AddConfigPath(dir)
	v.AddConfigPath(".")
	v.AddConfigPath(utils.ProjectRoot) // Default
	return v
}

// LoadFrameworkConfig loads framework config
func LoadFrameworkConfig(cfgPath string) (*FrameworkConfig, error) {
	dir, file := path.Split(cfgPath)
	v := defaultViper(dir, file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	usedDirPath, _ := path.Split(v.ConfigFileUsed())
	log.Info().
		Str("File", v.ConfigFileUsed()).
		Str("Directory Used", usedDirPath).
		Str("Hint", "If this is an unexpected file or path, it's likely that the provided one was unable to resolve and so a default was used").
		Msg("Loaded framework config file")
	var err error
	ProjectConfigDirectory, err = filepath.Abs(usedDirPath)
	if err != nil {
		return nil, err
	}

	var cfg *FrameworkConfig
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	ProjectFrameworkSettings = cfg
	return ProjectFrameworkSettings, err
}

// LoadNetworksConfig loads networks config
func LoadNetworksConfig(cfgPath string) (*NetworksConfig, error) {
	dir, file := path.Split(cfgPath)
	v := defaultViper(dir, file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	log.Info().
		Str("File", v.ConfigFileUsed()).
		Msg("Loaded networks config file")
	var cfg *NetworksConfig
	err := v.Unmarshal(&cfg)

	// Allow the networks config to be overridden when this codebase is imported as a library
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	ProjectNetworkSettings = cfg
	return ProjectNetworkSettings, err
}

// PrivateKeyStore enables access, through a variety of methods, to private keys for use in blockchain networks
type PrivateKeyStore interface {
	Fetch() ([]string, error)
}

// LocalStore retrieves keys defined in a networks.yaml file, or from environment variables
type LocalStore struct {
	RawKeys []string
}

// Fetch private keys from local environment variables or a config file
func (l *LocalStore) Fetch() ([]string, error) {
	if l.RawKeys == nil {
		return nil, errors.New("no keys found, ensure your configuration is properly set")
	}
	return l.RawKeys, nil
}

// ReadWriteRemoteRunnerConfig looks for an already existing remote config to read from, or asks the user to build one
func ReadWriteRemoteRunnerConfig(configFileLocation string) (*RemoteRunnerConfig, error) {
	config, err := readRemoteRunnerConfig(configFileLocation)
	return config, err
}

// Reads in the runner config
func readRemoteRunnerConfig(configLocation string) (*RemoteRunnerConfig, error) {
	var config *RemoteRunnerConfig
	absoluteConfigFileLocation, err := filepath.Abs(configLocation)
	if err != nil {
		log.Fatal().
			Str("Path", configLocation).
			Msg("Unable to resolve path to an absolute path")
		return nil, err
	}

	remoteRunnerConfig := filepath.Join(absoluteConfigFileLocation, "remote_runner_config.yaml")
	if os.Getenv("REMOTE_RUNNER_CONFIG_FILE") != "" {
		remoteRunnerConfig = os.Getenv("REMOTE_RUNNER_CONFIG_FILE")
	}

	remoteViper := viper.New()
	remoteViper.SetConfigFile(remoteRunnerConfig)
	if err := remoteViper.ReadInConfig(); err != nil {
		return nil, err
	}
	err = remoteViper.Unmarshal(&config)
	log.Info().Str("File", remoteRunnerConfig).Msg("Read Remote Runner Config")
	return config, err
}
