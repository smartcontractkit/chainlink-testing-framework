package config

import (
	_ "embed"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

//go:embed tomls/default.toml
var DefaultChainlinkImageConfig []byte

type ChainlinkImageConfig struct {
	ChainlinkImage *struct {
		Image   *string `toml:"image"`
		Version *string `toml:"version"`
	} `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage *struct {
		Image   *string `toml:"image"`
		Version *string `toml:"version"`
	} `toml:"ChainlinkUpgradeImage"`
}

func (c *ChainlinkImageConfig) ReadSecrets() error {
	return nil
}

func (c *ChainlinkImageConfig) Validate() error {
	if c.ChainlinkImage == nil {
		return errors.New("chainlink image information must be set")
	}

	if c.ChainlinkImage.Image == nil || *c.ChainlinkImage.Image == "" {
		return errors.New("chainlink image name must be set")
	}

	if c.ChainlinkImage.Version == nil || *c.ChainlinkImage.Version == "" {
		return errors.New("chainlink version must be set")
	}

	if c.ChainlinkUpgradeImage != nil {
		if c.ChainlinkUpgradeImage.Image == nil || *c.ChainlinkUpgradeImage.Image == "" {
			return errors.New("chainlink upgrade image name must be set")
		}

		if c.ChainlinkUpgradeImage.Version == nil || *c.ChainlinkUpgradeImage.Version == "" {
			return errors.New("chainlink upgrade version must be set")
		}
	}

	return nil
}

func (c *ChainlinkImageConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case ChainlinkImageConfig:
		if asCfg.ChainlinkImage.Image != nil && *asCfg.ChainlinkImage.Image != "" {
			c.ChainlinkImage.Image = asCfg.ChainlinkImage.Image
		}
		if asCfg.ChainlinkImage.Version != nil && *asCfg.ChainlinkImage.Version != "" {
			c.ChainlinkImage.Version = asCfg.ChainlinkImage.Version
		}
		if asCfg.ChainlinkUpgradeImage.Image != nil && *asCfg.ChainlinkUpgradeImage.Image != "" {
			c.ChainlinkUpgradeImage.Image = asCfg.ChainlinkUpgradeImage.Image
		}
		if asCfg.ChainlinkImage.Version != nil && *asCfg.ChainlinkUpgradeImage.Version != "" {
			c.ChainlinkUpgradeImage.Version = asCfg.ChainlinkUpgradeImage.Version
		}
		return nil
	case *ChainlinkImageConfig:
		if asCfg == nil {
			return nil
		}
		if asCfg.ChainlinkImage.Image != nil && *asCfg.ChainlinkImage.Image != "" {
			c.ChainlinkImage.Image = asCfg.ChainlinkImage.Image
		}
		if asCfg.ChainlinkImage.Version != nil && *asCfg.ChainlinkImage.Version != "" {
			c.ChainlinkImage.Version = asCfg.ChainlinkImage.Version
		}
		if asCfg.ChainlinkUpgradeImage.Image != nil && *asCfg.ChainlinkUpgradeImage.Image != "" {
			c.ChainlinkUpgradeImage.Image = asCfg.ChainlinkUpgradeImage.Image
		}
		if asCfg.ChainlinkImage.Version != nil && *asCfg.ChainlinkUpgradeImage.Version != "" {
			c.ChainlinkUpgradeImage.Version = asCfg.ChainlinkUpgradeImage.Version
		}
		return nil
	default:
		return errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (c *ChainlinkImageConfig) Default() error {
	if err := toml.Unmarshal(DefaultChainlinkImageConfig, c); err != nil {
		return errors.Wrapf(err, "error unmarshaling config")
	}

	return nil
}
