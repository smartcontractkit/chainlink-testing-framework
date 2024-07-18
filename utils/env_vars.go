package utils

import (
	"fmt"
	"os"
	"regexp"
)

// mustResolveEnvPlaceholder checks if the input string is an environment variable placeholder and resolves it.
func MustResolveEnvPlaceholder(input string) string {
	envVarName, hasEnvVar := LookupEnvVarName(input)
	if hasEnvVar {
		value, set := os.LookupEnv(envVarName)
		if !set {
			fmt.Fprintf(os.Stderr, "Error resolving '%s'. Environment variable '%s' not set or is empty\n", input, envVarName)
			os.Exit(1)
		}
		return value
	}
	return input
}

func LookupEnvVarName(input string) (string, bool) {
	re := regexp.MustCompile(`^{{ env\.([a-zA-Z_]+) }}$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}
