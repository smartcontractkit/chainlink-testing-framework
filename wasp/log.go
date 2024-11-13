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

// init initializes the default logging configuration for the application.
// It sets the logging level based on the environment variable LogLevelEnvVar,
// defaulting to "info" if the variable is not set.
func init() {
	initDefaultLogging()
}

// initDefaultLogging initializes the default logging configuration for the application.
// It sets the log level based on the environment variable specified by LogLevelEnvVar.
// If the environment variable is not set, it defaults to the "info" level.
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

// GetLogger returns a zerolog.Logger configured with the specified component name.
// It sets the log level based on the LogLevelEnvVar environment variable, defaulting to "info".
// If a testing.T object is provided, the logger outputs to a test writer; otherwise, it outputs to stderr.
// The logger includes a timestamp and the component name in its context.
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
	} else {
		return zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Str("Component", componentName).Logger()
	}
}
