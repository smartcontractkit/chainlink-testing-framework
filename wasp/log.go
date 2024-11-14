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

// init initializes the application's default logging configuration.
// It sets the logging level based on the LOG_LEVEL environment variable,
// defaulting to "info" if the variable is not set.
// If an invalid log level is provided, the application panics.
func init() {
	initDefaultLogging()
}

// initDefaultLogging initializes the application's logging system.
// It sets the log level based on the LOG_LEVEL environment variable, defaulting to "info" if not set,
// and configures logs to be output to standard error.
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

// GetLogger returns a zerolog.Logger configured for the specified component.
// If the testing object t is non-nil, the logger integrates with the test output.
// The log level is determined by the LogLevelEnvVar environment variable or defaults to "info".
// The logger includes a timestamp and the component name as a field.
// It outputs logs to standard error in a console-friendly format.
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
