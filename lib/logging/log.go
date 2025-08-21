package logging

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	tclog "github.com/testcontainers/testcontainers-go/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
)

const afterTestEndedMsg = "LOG AFTER TEST ENDED"

// Logger is an alias for zerolog.Logger, exposed through the logging package
type Logger = zerolog.Logger

// CustomT wraps testing.T for two purposes:
// 1. it implements Write to override the default logger
// 2. it implements Printf to implement the testcontainers-go/Logging interface
// The reason for both of these is that go parallel testing causes the logs to get mixed up,
// so we need to override the default logger to *testing.T.Log to ensure that the logs are
// properly associated with the tests running. The testcontainers-go/Logging interface complicates
// this more since it needs a struct with L to hold the logger and needs to override Printf.
type CustomT struct {
	*testing.T
	L     zerolog.Logger
	ended bool
}

// Write writes the contents of p to the logger associated with CustomT.
// It handles empty input gracefully and logs a warning if called after the test has ended.
// Returns the number of bytes written and any write error encountered.
func (ct *CustomT) Write(p []byte) (n int, err error) {
	str := string(p)
	if strings.TrimSpace(str) == "" {
		return len(p), nil
	}
	if ct.ended {
		l := GetTestLogger(nil)
		l.Warn().Msgf("%s %s: %s", afterTestEndedMsg, ct.Name(), string(p))
		return len(p), nil
	}
	ct.T.Log(strings.TrimSuffix(str, "\n"))
	return len(p), nil
}

// Printf implements the testcontainers-go/Logging interface.
func (ct CustomT) Printf(format string, v ...interface{}) {
	if ct.ended {
		s := "%s: "
		formatted := fmt.Sprintf("%s %s%s", afterTestEndedMsg, s, format)
		l := GetTestLogger(nil)
		l.Warn().Msgf(formatted, ct.Name(), v)
	} else {
		ct.L.Info().Msgf(format, v...)
	}
}

// Init initializes the logging system by setting up a logger with the specified log level.
// It ensures that logging is configured before any application components are initialized.
func Init() {
	l := GetLogger(nil, config.EnvVarLogLevel)
	log.Logger = l
}

// GetLogger returns a logger that will write to the testing.T.Log function using the env var provided for the log level.
// nil can be passed for t to get a logger that is not associated with a go test.
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
		// Use Cleanup function to set ended to true once the test completes
		t.Cleanup(func() {
			ct.ended = true
		})
		return zerolog.New(ct).Output(zerolog.ConsoleWriter{Out: ct, TimeFormat: "15:04:05.00"}).Level(lvl).With().Timestamp().Logger()
	}
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.00"}).Level(lvl).With().Timestamp().Logger()
}

// GetTestLogger returns a logger that will write to the testing.T.Log function using the env var for log level.
// nil can be passed for t to get a logger that is not associated with a go test.
func GetTestLogger(t *testing.T) zerolog.Logger {
	return GetLogger(t, config.EnvVarLogLevel)
}

// GetTestContainersGoTestLogger returns a logger that will write to the testing.T.Log function using the env var for log level
// for logs that testcontainers-go will log out. nil can be passed to this and it will be treated as the default tc.Logger
func GetTestContainersGoTestLogger(t *testing.T) tclog.Logger {
	if t != nil {
		return CustomT{
			T: t,
			L: GetTestLogger(t),
		}
	}
	return tclog.Default()
}

// SplitStringIntoChunks takes a string and splits it into chunks of a specified size.
func SplitStringIntoChunks(s string, chunkSize int) []string {
	// Length of the string.
	strLen := len(s)

	// Number of chunks needed.
	numChunks := (strLen + chunkSize - 1) / chunkSize

	// Slice to hold the chunks.
	chunks := make([]string, numChunks)

	// Loop to create chunks.
	for i := 0; i < numChunks; i++ {
		// Calculate the start and end indices of the chunk.
		start := i * chunkSize
		end := start + chunkSize

		// If the end index goes beyond the string length, adjust it to the string length.
		if end > strLen {
			end = strLen
		}

		// Slice the string and add the chunk to the slice.
		chunks[i] = s[start:end]
	}

	return chunks
}
