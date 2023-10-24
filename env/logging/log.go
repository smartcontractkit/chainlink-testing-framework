package logging

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/config"
)

// People often call this multiple times
var loggingMu sync.Mutex

func getTestLevel() zerolog.Level {
	lvlStr := os.Getenv(config.EnvVarLogLevel)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	return lvl
}

func Init() {
	loggingMu.Lock()
	defer loggingMu.Unlock()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.00"}).With().Timestamp().Logger().Level(getTestLevel())
}

func GetTestLogger(t *testing.T) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	return zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.00"}).With().Timestamp().Logger().Level(getTestLevel())
}
