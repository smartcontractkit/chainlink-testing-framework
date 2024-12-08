// File: internal/default_logger.go
package internal

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )

// // DefaultLogger is the default implementation of the Logger interface using zerolog.
// type DefaultLogger struct {
// 	logger zerolog.Logger
// }

// // NewDefaultLogger initializes and returns a DefaultLogger.
// func NewDefaultLogger() *DefaultLogger {
// 	lvlStr := os.Getenv("SENTINEL_LOG_LEVEL")
// 	if lvlStr == "" {
// 		lvlStr = "info"
// 	}
// 	lvl, err := zerolog.ParseLevel(strings.ToLower(lvlStr))
// 	if err != nil {
// 		panic(fmt.Sprintf("failed to parse log level: %v", err))
// 	}

// 	// Configure zerolog to output to the console with a human-friendly format.
// 	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
// 		Level(lvl).
// 		With().
// 		Timestamp().
// 		Logger()

// 	return &DefaultLogger{
// 		logger: logger,
// 	}
// }

// // Debug logs a message at debug level.
// func (dl *DefaultLogger) Debug(args ...interface{}) {
// 	dl.logger.Debug().Msg(fmt.Sprint(args...))
// }

// // Info logs a message at info level.
// func (dl *DefaultLogger) Info(args ...interface{}) {
// 	dl.logger.Info().Msg(fmt.Sprint(args...))
// }

// // Warn logs a message at warn level.
// func (dl *DefaultLogger) Warn(args ...interface{}) {
// 	dl.logger.Warn().Msg(fmt.Sprint(args...))
// }

// // Error logs a message at error level.
// func (dl *DefaultLogger) Error(args ...interface{}) {
// 	dl.logger.Error().Msg(fmt.Sprint(args...))
// }

// // Debugf logs a formatted message at debug level.
// func (dl *DefaultLogger) Debugf(format string, args ...interface{}) {
// 	dl.logger.Debug().Msgf(format, args...)
// }

// // Infof logs a formatted message at info level.
// func (dl *DefaultLogger) Infof(format string, args ...interface{}) {
// 	dl.logger.Info().Msgf(format, args...)
// }

// // Warnf logs a formatted message at warn level.
// func (dl *DefaultLogger) Warnf(format string, args ...interface{}) {
// 	dl.logger.Warn().Msgf(format, args...)
// }

// // Errorf logs a formatted message at error level.
// func (dl *DefaultLogger) Errorf(format string, args ...interface{}) {
// 	dl.logger.Error().Msgf(format, args...)
// }
