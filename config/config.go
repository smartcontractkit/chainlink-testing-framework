// Package config enables loading and utilizing configuration options for different blockchain networks
package config

import (
	"errors"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// ConfigurationType refers to the different ways that configurations can be set
type ConfigurationType string

const (
	LocalConfig ConfigurationType = "local"
)

var ProjectConfig Config
var ProjectConfigDirectory string

// LoadFromEnv loads all config files and environment variables
func LoadFromEnv() error {
	return envconfig.Process("", &ProjectConfig)
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
