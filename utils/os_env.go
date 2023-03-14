package utils

import (
	"fmt"
	"os"
)

// GetEnv returns the value of the environment variable named by the key
// and sets the environment variable up to be used in the remote runner
func GetEnv(key string) string {
	val := os.Getenv(key)
	if val != "" {
		os.Setenv(fmt.Sprintf("TEST_%s", key), val)
	}
	return val
}
