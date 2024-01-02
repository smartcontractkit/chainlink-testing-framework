package config

import (
	"dario.cat/mergo"
)

func MustConfigOverrideChainlinkVersion(config *ChainlinkImageConfig, target interface{}) {
	if config == nil {
		panic("ChainlinkImageConfig must not be nil")
	}
	if config.Image != nil && *config.Image != "" && config.Version != nil && *config.Version != "" {
		if err := mergo.Merge(target, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   *config.Image,
					"version": *config.Version,
				},
			},
		}, mergo.WithOverride); err != nil {
			panic(err)
		}
	}
}
