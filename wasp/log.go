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

// init initializes the logging configuration for the application. 
// It sets the log level based on the environment variable specified by LogLevelEnvVar. 
// If the environment variable is not set, it defaults to "info". 
// In case of an invalid log level, the function will panic, ensuring that the application does not start with an incorrect logging configuration.
func init() {
	initDefaultLogging()
}

// initDefaultLogging initializes the logging configuration for the application. 
// It retrieves the log level from the environment variable specified by LogLevelEnvVar. 
// If the environment variable is not set, it defaults to the "info" level. 
// The function sets up the logger to output to standard error with the specified log level. 
// If the log level cannot be parsed, the function will panic, terminating the application.
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

// GetLogger returns a zerolog.Logger instance configured for logging. 
// It sets the log level based on the environment variable specified by LogLevelEnvVar, defaulting to "info" if not set. 
// If the provided testing.T pointer is not nil, the logger will be configured to write to the test output; 
// otherwise, it will log to standard error. 
// The logger will include a timestamp and a string field indicating the component name.
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
