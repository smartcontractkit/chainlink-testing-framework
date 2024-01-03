package config

import (
	"dario.cat/mergo"
)

func MustConfigOverrideChainlinkVersion(config *ChainlinkImageConfig, target interface{}) {
	if config == nil {
		panic("[ChainlinkImageConfig] must be present")
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

func MightConfigOverridePyroscopeKey(config *PyroscopeConfig, target interface{}) {
	if config == nil || *config.Key == "" {
		return
	}

	env := make(map[string]string)
	env["CL_PYROSCOPE_AUTH_TOKEN"] = *config.Key

	if err := mergo.Merge(target, map[string]interface{}{
		"env": env,
	}, mergo.WithOverride); err != nil {
		panic(err)
	}
}
