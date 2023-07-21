package utils

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

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
