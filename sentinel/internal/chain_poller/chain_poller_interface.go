// File: internal/chain_poller/chain_poller_interface.go
package chain_poller

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
)

type ChainPollerInterface interface {
	FilterLogs(ctx context.Context, filterQueries []api.FilterQuery) ([]api.Log, error)
}
