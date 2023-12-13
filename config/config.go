package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

// GenericConfig is an interface for all product based config types to implement
type GenericConfig interface {
	ApplySecrets() error
	Validate() error
	ApplyOverrides(from interface{}) error
	ApplyDecoded(configDecoded string) error
	ApplyBase64Enconded(configEncoded string) error
	Default() error
}

func Read[T GenericConfig](configEnvPrefix string, cfg T) error {
	// Read default config
	err := cfg.Default()
	if err != nil {
		return err
	}

	// Use config from file
	configPath, isSet := os.LookupEnv(fmt.Sprintf("%s_TOML_PATH", strings.ToUpper(configEnvPrefix)))
	if isSet {
		d, err := os.ReadFile(configPath)
		if err != nil {
			return errors.Wrapf(err, "error reading config file: ")
		}

		var fileConfig T
		err = toml.Unmarshal(d, &fileConfig)
		if err != nil {
			return errors.Wrapf(err, "error unmarshaling config")
		}

		err = cfg.ApplyOverrides(fileConfig)
		if err != nil {
			return errors.Wrapf(err, "error applying overrides from config file to default config")
		}
	}

	// Use base64 encoded config if set (e.g. on Github CI)
	configEncoded, isSet := os.LookupEnv(fmt.Sprintf("%s_BASE64_TOML_CONTENT", configEnvPrefix))
	if isSet && configEncoded != "" {
		decoded, err := base64.StdEncoding.DecodeString(configEncoded)
		if err != nil {
			return err
		}

		var base64override GenericConfig
		err = toml.Unmarshal(decoded, &base64override)
		if err != nil {
			return errors.Wrapf(err, "error unmarshaling base64 config")
		}

		err = cfg.ApplyOverrides(base64override)
		if err != nil {
			return errors.Wrapf(err, "error applying overrides from base64 config file to config")
		}
	}

	err = cfg.ApplySecrets()
	if err != nil {
		return errors.Wrapf(err, "error reading secrets")
	}

	err = cfg.Validate()
	if err != nil {
		return errors.Wrapf(err, "error validating config")
	}

	return nil
}
