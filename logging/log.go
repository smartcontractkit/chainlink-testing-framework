package logging

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/stretchr/testify/require"
)

// CustomT wraps testing.T for two puposes:
// 1. it implements Write to override the default logger
// 2. it implements Printf to implement the testcontainers-go/Logging interface
type CustomT struct {
	*testing.T
	L zerolog.Logger
}

func (ct *CustomT) Write(p []byte) (n int, err error) {
	str := string(p)
	if strings.TrimSpace(str) == "" {
		return len(p), nil
	}
	ct.T.Log(strings.TrimSuffix(str, "\n"))
	return len(p), nil
}

// Printf implements the testcontainers-go/Logging interface.
func (ct CustomT) Printf(format string, v ...interface{}) {
	ct.L.Info().Msgf(format, v...)
}

func Init() {
	l := GetLogger(nil, config.EnvVarLogLevel)
	log.Logger = l
}

func GetLogger(t *testing.T, envVarName string) zerolog.Logger {
	lvlStr := os.Getenv(envVarName)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		if t != nil {
			require.NoError(t, err, "error parsing log level")
		} else {
			panic(fmt.Sprintf("failed to parse log level: %s", err))
		}
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	if t != nil {
		ct := &CustomT{T: t}
		return zerolog.New(ct).Output(zerolog.ConsoleWriter{Out: ct, TimeFormat: "15:04:05.00"}).Level(lvl).With().Timestamp().Logger()
	}
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.00"}).Level(lvl).With().Timestamp().Logger()
}

func GetTestLogger(t *testing.T) zerolog.Logger {
	return GetLogger(t, config.EnvVarLogLevel)
}
