// File: chain_poller/interface.go
package chain_poller

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/internal"
)

// ChainPollerInterface defines the methods that ChainPoller must implement.
type ChainPollerInterface interface {
	Poll(ctx context.Context, filterQueries []internal.FilterQuery) ([]internal.Log, error)
}
