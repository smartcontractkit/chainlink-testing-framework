package config

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/net"
)

type PyroscopeConfig struct {
	Enabled     *bool   `toml:"enabled"`
	ServerUrl   *string `toml:"-"`
	Key         *string `toml:"-"`
	Environment *string `toml:"-"`
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
