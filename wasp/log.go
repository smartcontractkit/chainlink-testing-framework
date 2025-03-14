package wasp

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogLevelEnvVar = "WASP_LOG_LEVEL"
)

// init initializes the default logging configuration for the package by setting the logging level and output destination.
func init() {
	initDefaultLogging()
}

// initDefaultLogging configures the default logger using the LogLevelEnvVar environment variable.
// It sets the logging output to standard error and defaults to the "info" level if the variable is unset.
func initDefaultLogging() {
	lvlStr := os.Getenv(LogLevelEnvVar)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}

// GetLogger returns a zerolog.Logger configured with the specified component name and log level.
// If a *testing.T is provided, the logger integrates with test output.
// Use it to enable consistent logging across components with environment-based log level control.
func GetLogger(t *testing.T, componentName string) zerolog.Logger {
	lvlStr := os.Getenv(LogLevelEnvVar)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse log level: %s", err))
	}
	if t != nil {
		return zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Str("Component", componentName).Logger()
	}
	return zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Str("Component", componentName).Logger()
}
