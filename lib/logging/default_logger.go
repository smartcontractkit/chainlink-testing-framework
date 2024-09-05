package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogLevelEnvVar = "TEST_LOG_LEVEL"
)

var (
	L zerolog.Logger
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
	L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}
