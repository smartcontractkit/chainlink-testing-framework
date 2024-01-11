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

func (c *ChainlinkImageConfig) ApplyOverrides(from *ChainlinkImageConfig) error {
	if from == nil {
		return nil
	}
	if from.Image != nil && *from.Image != "" {
		c.Image = from.Image
	}
	if from.Version != nil && *from.Version != "" {
		c.Version = from.Version
	}
	return nil
}

func (c *ChainlinkImageConfig) Default() error {
	return nil
}
