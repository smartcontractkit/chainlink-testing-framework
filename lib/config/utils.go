package config

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// BytesToAnyTomlStruct unmarshals the given bytes into the given target struct, apart from unmarshalling the root
// it also looks for given configuration name and unmarshals it into the target struct overwriting the root.
// Target needs to be a struct with exported fields with `toml:"field_name"` tags.
func BytesToAnyTomlStruct(logger zerolog.Logger, filename, configurationName string, target any, content []byte) error {
	var someToml map[string]interface{}
	err := toml.Unmarshal(content, &someToml)
	if err != nil {
		return err
	}

	if configurationName == "" {
		logger.Debug().Msgf("No configuration name provided, will read only default configuration.")

		err := toml.Unmarshal(content, target)

		if err != nil {
			return errors.Wrapf(err, "error unmarshalling config")
		}

		logger.Debug().Msgf("Successfully unmarshalled %s config file", filename)

		return nil
	}

	if _, ok := someToml[configurationName]; !ok {
		logger.Debug().Msgf("Config file %s does not contain configuration named '%s'. Won't read anything", filename, configurationName)
		return nil
	}

	marshalled, err := toml.Marshal(someToml[configurationName])
	if err != nil {
		return err
	}

	err = toml.Unmarshal(marshalled, target)
	if err != nil {
		return err
	}

	logger.Debug().Msgf("Configuration named '%s' read successfully.", configurationName)

	return nil
}
