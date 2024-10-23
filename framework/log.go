package framework

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvVarLogLevel = "CTF_LOG_LEVEL"
)

var (
	L zerolog.Logger
)

func init() {
	initDefaultLogging()
}

func initDefaultLogging() {
	lvlStr := os.Getenv(EnvVarLogLevel)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}
