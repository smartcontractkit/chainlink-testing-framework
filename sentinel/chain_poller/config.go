// File: chain_poller/config.go
package chain_poller

import (
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
)

// ChainEventPollerConfig holds configuration for ChainEventPoller
type ChainEventPollerConfig struct {
	BlockchainClient BlockchainClient
	PollInterval     time.Duration
	Logger           internal.Logger
	ChainID          int64
}
