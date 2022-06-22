// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
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

// LoadFromEnv loads all config files and environment variables
func LoadFromEnv() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := envconfig.Process("", &ProjectConfig); err != nil {
		return err
	}
	log.Logger = log.Logger.Level(zerolog.Level(ProjectConfig.FrameworkConfig.Logging.Level))
	return nil
}

// LoadRemoteEnv loads environment variables when running on remote test runner
func LoadRemoteEnv() error {
	err := LoadFromEnv()
	if strings.Contains(err.Error(), "envconfig.Process: assigning REMOTE_RUNNER_CONFIG_FILE to RemoteRunnerConfig") {
		// a remote runner no longer needs the remote config file
		return nil
	}
	return err
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
