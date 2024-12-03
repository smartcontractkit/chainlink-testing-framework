// File: chain_poller/interface.go
package chain_poller

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
)

// BlockchainClient defines the required methods for interacting with a blockchain.
type BlockchainClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, query internal.FilterQuery) ([]internal.Log, error)
}

// ChainPollerInterface defines the methods that ChainPoller must implement.
type ChainPollerInterface interface {
	Poll(ctx context.Context, filterQueries []internal.FilterQuery) ([]internal.Log, error)
}
