package utils

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"fmt"

	envConf "github.com/smartcontractkit/chainlink-env/config"
)

// GetTestLogger instantiates a logger that takes into account the test context and the log level
func GetTestLogger(t *testing.T) zerolog.Logger {
	lvlStr := os.Getenv(envConf.EnvVarLogLevel)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	require.NoError(t, err, "error parsing log level")
	l := zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Logger()
	return l
}

// GetLogger instantiates a logger that takes into account the test context and the log level
func GetLogger(t *testing.T, envVarName string) zerolog.Logger {
	lvlStr := os.Getenv(envVarName)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse log level: %s", err))
	}
	if t == nil {
		return zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Logger()
	}
	return zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Logger()
}
