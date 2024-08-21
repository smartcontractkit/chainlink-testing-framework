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

func init() {
	initDefaultLogging()
}

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

// GetLogger instantiates a logger that takes into account the test context and the log level
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
