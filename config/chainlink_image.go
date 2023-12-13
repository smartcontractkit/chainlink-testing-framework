package config

import (
	"github.com/pkg/errors"
)

type ChainlinkImageConfig struct {
	Image   *string `toml:"image"`
	Version *string `toml:"version"`
}

func (c *ChainlinkImageConfig) Validate() error {
	if c.Image == nil || *c.Image == "" {
		return errors.New("chainlink image name must be set")
	}

	if c.Version == nil || *c.Version == "" {
		return errors.New("chainlink version must be set")
	}

	return nil
}

func (c *ChainlinkImageConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case ChainlinkImageConfig:
		if asCfg.Image != nil && *asCfg.Image != "" {
			c.Image = asCfg.Image
		}
		if asCfg.Version != nil && *asCfg.Version != "" {
			c.Version = asCfg.Version
		}
		return nil
	case *ChainlinkImageConfig:
		if asCfg == nil {
			return nil
		}
		if asCfg.Image != nil && *asCfg.Image != "" {
			c.Image = asCfg.Image
		}
		if asCfg.Version != nil && *asCfg.Version != "" {
			c.Version = asCfg.Version
		}
		return nil
	default:
		return errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (c *ChainlinkImageConfig) Default() error {
	return nil
}
