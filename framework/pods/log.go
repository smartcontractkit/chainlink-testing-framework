package pods

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvVarLogLevel = "PODS_LOG_LEVEL"
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
	if err != nil { // coverage-ignore
		panic(fmt.Sprintf("invalid log level: %s", lvlStr))
	}
	L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}
