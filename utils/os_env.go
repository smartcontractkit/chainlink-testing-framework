package utils

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-env/config"
)

// GetEnv returns the value of the environment variable named by the key
// and sets the environment variable up to be used in the remote runner
func GetEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val != "" {
		prefixedKey := fmt.Sprintf("%s%s", config.EnvVarPrefix, key)
		if os.Getenv(prefixedKey) != "" {
			return val, fmt.Errorf("environment variable collision with prefixed key, Original: %s, Prefixed: %s", key, prefixedKey)
		}
		err := os.Setenv(prefixedKey, val)
		if err != nil {
			return val, fmt.Errorf("failed to set environment variable %s: %v", prefixedKey, err)
		}
	}
	return val, nil
}

// SetupEnvVarsForRemoteRunner sets up the environment variables in the list to propagate to the remote runner
func SetupEnvVarsForRemoteRunner(envVars []string) error {
	for _, envVar := range envVars {
		_, err := GetEnv(envVar)
		if err != nil {
			return err
		}
	}
	return nil
}
