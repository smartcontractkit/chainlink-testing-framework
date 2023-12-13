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

func (c *PyroscopeConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case PyroscopeConfig:
		if asCfg.Enabled != nil {
			c.Enabled = asCfg.Enabled
		}
		if asCfg.ServerUrl != nil {
			c.ServerUrl = asCfg.ServerUrl
		}
		if asCfg.Key != nil {
			c.Key = asCfg.Key
		}
		if asCfg.Environment != nil {
			c.Environment = asCfg.Environment
		}
		return nil
	case *PyroscopeConfig:
		if asCfg == nil {
			return nil
		}
		if asCfg.Enabled != nil {
			c.Enabled = asCfg.Enabled
		}
		if asCfg.ServerUrl != nil {
			c.ServerUrl = asCfg.ServerUrl
		}
		if asCfg.Key != nil {
			c.Key = asCfg.Key
		}
		if asCfg.Environment != nil {
			c.Environment = asCfg.Environment
		}
		return nil
	default:
		return errors.Errorf("cannot apply overrides to pyroscope config from type %T", from)
	}
}

func (c *PyroscopeConfig) Default() error {
	return nil
}
