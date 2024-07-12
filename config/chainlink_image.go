package config

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

type ChainlinkImageConfig struct {
	Image           *string `toml:"image"`
	Version         *string `toml:"version"`
	PostgresVersion *string `toml:"postgres_version,omitempty"`
}

func (c *ChainlinkImageConfig) LoadFromEnv(envName string) error {
	logger := logging.GetTestLogger(nil)

	image, err := readEnvVarValue(envName, String)
	if err != nil {
		return err
	}
	if image != nil && image.(string) != "" {
		logger.Debug().Msgf("Using %s env var to override ChainlinkImage.Image", envName)
		c.Image = ptr.Ptr(image.(string))
	}

	return nil
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
