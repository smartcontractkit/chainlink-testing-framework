package test_utils

import (
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// CopyConfig creates a deep copy of Seth config
func CopyConfig(config *seth.Config) (*seth.Config, error) {
	marshalled, err := toml.Marshal(config)
	if err != nil {
		return nil, err
	}

	var configCopy seth.Config
	err = toml.Unmarshal(marshalled, &configCopy)
	if err != nil {
		return nil, err
	}

	return &configCopy, nil
}
