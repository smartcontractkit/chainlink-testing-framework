package config

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

type PyroscopeConfig struct {
	Pyroscope *struct {
		Enabled     *bool   `toml:"enabled"`
		ServerUrl   *string `toml:"server_url"`
		Key         *string `toml:"key"`
		Environment *string `toml:"environment"`
	} `toml:"Pyroscope"`
}

func (n *PyroscopeConfig) ReadDecoded(configDecoded string) error {
	if configDecoded != "" {
		return nil
	}

	var cfg PyroscopeConfig
	err := toml.Unmarshal([]byte(configDecoded), &cfg)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling pyroscope config")
	}

	err = n.ApplyOverrides(cfg)
	if err != nil {
		return errors.Wrapf(err, "error applying overrides from decoded pyroscope config file to config")
	}

	return nil
}

func (c *PyroscopeConfig) ReadSecrets() error {
	return nil
}

func (c *PyroscopeConfig) Validate() error {
	if c.Pyroscope != nil && c.Pyroscope.Enabled != nil && *c.Pyroscope.Enabled {
		if c.Pyroscope.ServerUrl == nil || *c.Pyroscope.ServerUrl == "" {
			return errors.New("pyroscope server url must be set")
		}
		// if c.Pyroscope.Key == "" {
		// 	return errors.New("pyroscope key must be set")
		// }
		// if c.Pyroscope.Environment == "" {
		// 	return errors.New("pyroscope environment must be set")
		// }
	}

	return nil
}

func (c *PyroscopeConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case PyroscopeConfig:
		if asCfg.Pyroscope == nil {
			return nil
		}
		if asCfg.Pyroscope.Enabled != nil {
			c.Pyroscope.Enabled = asCfg.Pyroscope.Enabled
		}
		if asCfg.Pyroscope.ServerUrl != nil {
			c.Pyroscope.ServerUrl = asCfg.Pyroscope.ServerUrl
		}
		if asCfg.Pyroscope.Key != nil {
			c.Pyroscope.Key = asCfg.Pyroscope.Key
		}
		if asCfg.Pyroscope.Environment != nil {
			c.Pyroscope.Environment = asCfg.Pyroscope.Environment
		}
		return nil
	case *PyroscopeConfig:
		if asCfg.Pyroscope == nil {
			return nil
		}
		if asCfg.Pyroscope == nil {
			return nil
		}
		if asCfg.Pyroscope.Enabled != nil {
			c.Pyroscope.Enabled = asCfg.Pyroscope.Enabled
		}
		if asCfg.Pyroscope.ServerUrl != nil {
			c.Pyroscope.ServerUrl = asCfg.Pyroscope.ServerUrl
		}
		if asCfg.Pyroscope.Key != nil {
			c.Pyroscope.Key = asCfg.Pyroscope.Key
		}
		if asCfg.Pyroscope.Environment != nil {
			c.Pyroscope.Environment = asCfg.Pyroscope.Environment
		}
		return nil
	default:
		return errors.Errorf("cannot apply overrides to pyroscope config from type %T", from)
	}
}

func (c *PyroscopeConfig) Default() error {
	return nil
}
