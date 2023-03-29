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
			return "", fmt.Errorf("environment variable collision with prefixed key, Original: %s, Prefixed: %s", key, prefixedKey)
		}
		os.Setenv(prefixedKey, val)
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
