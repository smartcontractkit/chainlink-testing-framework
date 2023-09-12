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

// CustomT wraps testing.T and provides a bytes.Buffer to capture logs.
type customT struct {
	*testing.T
}

func (ct *customT) Write(p []byte) (n int, err error) {
	str := string(p)
	if strings.TrimSpace(str) == "" {
		return len(p), nil
	}
	ct.T.Log(strings.TrimSuffix(str, "\n"))
	return len(p), nil
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
		ct := &customT{T: t}
		return zerolog.New(ct).Output(zerolog.ConsoleWriter{Out: ct, TimeFormat: "15:04:05.00"}).Level(lvl).With().Timestamp().Logger()
	}
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.00"}).Level(lvl).With().Timestamp().Logger()
}

func GetTestLogger(t *testing.T) zerolog.Logger {
	return GetLogger(t, config.EnvVarLogLevel)
}
