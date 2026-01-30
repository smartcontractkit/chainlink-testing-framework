package seth

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const (
	LogLevelEnvVar = "SETH_LOG_LEVEL"
)

func newLogger() zerolog.Logger {
	lvlStr := strings.TrimSpace(os.Getenv(LogLevelEnvVar))
	if lvlStr == "" {
		lvlStr = zerolog.InfoLevel.String()
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Timestamp().Logger()
}

// NewLogger returns a zerolog.Logger configured using the environment defaults.
func NewLogger() zerolog.Logger {
	return newLogger()
}
