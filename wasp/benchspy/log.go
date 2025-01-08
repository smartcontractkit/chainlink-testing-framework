package benchspy

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogLevelEnvVar = "BENCHSPY_LOG_LEVEL"
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

	output := zerolog.ConsoleWriter{Out: os.Stderr}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\033[38;5;136m%v \033[0m", i) // Dark gold color for message
	}

	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("\033[38;5;136m%v \033[0m", i) // Dark gold color for field value
	}

	L = log.Output(output).Level(lvl)
}
