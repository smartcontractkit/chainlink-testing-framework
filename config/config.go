// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/helmenv/environment"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"
)

// ConfigurationType refers to the different ways that configurations can be set
type ConfigurationType string

// Configs
const (
	LocalConfig ConfigurationType = "local"

	DefaultGeth     string = "geth"
	PerformanceGeth string = "geth_performance"
	RealisticGeth   string = "geth_realistic"
)

// NetworkSettings is a map that holds configuration for each individual network
type NetworkSettings map[string]map[string]interface{}

var ProjectConfig Config
var ProjectConfigDirectory string

// ChainlinkVals formats Chainlink values set in the framework config to be passed to Chainlink deployments
func ChainlinkVals() map[string]interface{} {
	if ProjectConfig.FrameworkConfig == nil {
		log.Error().Msg("ProjectFrameworkSettings not set!")
		return nil
	}
	values := map[string]interface{}{}
	if len(ProjectConfig.FrameworkConfig.ChainlinkEnvValues) > 0 {
		values["env"] = ProjectConfig.FrameworkConfig.ChainlinkEnvValues
	}
	if ProjectConfig.FrameworkConfig.ChainlinkImage != "" {
		values["chainlink"] = map[string]interface{}{
			"image": map[string]interface{}{
				"image":   ProjectConfig.FrameworkConfig.ChainlinkImage,
				"version": ProjectConfig.FrameworkConfig.ChainlinkVersion,
			},
		}
	}
	return values
}

// GethNetworks builds the proper geth network settings to use based on the selected_networks config
func GethNetworks() []environment.SimulatedNetwork {
	if ProjectConfig.NetworksConfig == nil {
		log.Error().Msg("ProjectNetworkSettings not set!")
		return nil
	}
	var gethNetworks []environment.SimulatedNetwork
	for _, network := range ProjectConfig.NetworksConfig.SelectedNetworks {
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

func LoadFromEnv() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := envconfig.Process("", &ProjectConfig); err != nil {
		return err
	}
	log.Logger = log.Logger.Level(zerolog.Level(ProjectConfig.FrameworkConfig.Logging.Level))
	return nil
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
