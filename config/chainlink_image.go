package config

import (
	"errors"
)

type ChainlinkImageConfig struct {
	Image   *string `toml:"image"`
	Version *string `toml:"version"`
}

// Validate checks that the chainlink image config is valid, which means that
// both image and version are set and non-empty
func (c *ChainlinkImageConfig) Validate() error {
	if c.Image == nil || *c.Image == "" {
		return errors.New("chainlink image name must be set")
	}

	if c.Version == nil || *c.Version == "" {
		return errors.New("chainlink version must be set")
	}

	return nil
}
