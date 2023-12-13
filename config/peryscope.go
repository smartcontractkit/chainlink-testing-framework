package config

import (
	"github.com/pkg/errors"
)

type PyroscopeConfig struct {
	Enabled     *bool   `toml:"enabled"`
	ServerUrl   *string `toml:"server_url"`
	Key         *string `toml:"key"`
	Environment *string `toml:"environment"`
}

func (c *PyroscopeConfig) Validate() error {
	if c.Enabled != nil && *c.Enabled {
		if c.ServerUrl == nil {
			return errors.New("pyroscope server url must be set")
		}
		if !isValidURL(*c.ServerUrl) {
			return errors.Errorf("invalid pyroscope server url %s", *c.ServerUrl)
		}
		if c.Environment == nil {
			return errors.New("pyroscope environment must be set")
		}
	}

	return nil
}

func (c *PyroscopeConfig) ApplyOverrides(from *PyroscopeConfig) error {
	if from == nil {
		return nil
	}
	if from.Enabled != nil {
		c.Enabled = from.Enabled
	}
	if from.ServerUrl != nil {
		c.ServerUrl = from.ServerUrl
	}
	if from.Key != nil {
		c.Key = from.Key
	}
	if from.Environment != nil {
		c.Environment = from.Environment
	}
	return nil
}

func (c *PyroscopeConfig) Default() error {
	return nil
}
