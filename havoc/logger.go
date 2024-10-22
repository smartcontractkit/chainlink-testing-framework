package havoc

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Default logger
var Logger zerolog.Logger

func init() {
	// Default logger
	Logger = CreateLogger(LoggerConfig{
		LogOutput: os.Getenv("CHAOS_LOG_OUTPUT"),
		LogLevel:  os.Getenv("CHAOS_LOG_LEVEL"),
		LogType:   "chaos",
	})
}

type LoggerConfig struct {
	LogOutput string // "json-console" for JSON output, empty or "console" for human-friendly console output
	LogLevel  string // Log level (e.g., "info", "debug", "error")
	LogType   string // Custom log type identifier
}

// Create initializes a zerolog.Logger based on the specified configuration.
func CreateLogger(config LoggerConfig) zerolog.Logger {
	// Parse the log level
	lvl, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err) // Consider more graceful error handling based on your application's requirements
	}

	switch config.LogOutput {
	case "json-console":
		// Configure for JSON console output
		return zerolog.New(os.Stderr).Level(lvl).With().Timestamp().Str("type", config.LogType).Logger()
	default:
		// Configure for console (human-friendly) output
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(lvl).With().Timestamp().Str("type", config.LogType).Logger()
	}
}
