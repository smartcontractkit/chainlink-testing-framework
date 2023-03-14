package utils

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-env/config"
)

// GetEnv returns the value of the environment variable named by the key
// and sets the environment variable up to be used in the remote runner
func GetEnv(key string) string {
	val := os.Getenv(key)
	if val != "" {
		os.Setenv(fmt.Sprintf("%s%s", config.EnvVarPrefix, key), val)
	}
	return val
}
