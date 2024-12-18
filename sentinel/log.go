// File: log.go
package sentinel

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

const (
	LogLevelEnvVar = "SENTINEL_LOG_LEVEL"
)

// GetLogger instantiates a logger that takes into account the test context and the log level
func GetLogger(t *testing.T, componentName string) zerolog.Logger {
	return logging.GetLogger(t, LogLevelEnvVar).With().Str("Component", componentName).Logger()
}
