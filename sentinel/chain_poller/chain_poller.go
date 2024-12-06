// File: chain_poller/chain_poller.go
package chain_poller

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
)

// ChainPollerConfig holds the configuration for the ChainPoller.
type ChainPollerConfig struct {
	BlockchainClient BlockchainClient
	Logger           internal.Logger
	ChainID          int64
}

// ChainPoller is responsible for polling logs from the blockchain.
type ChainPoller struct {
	Config ChainPollerConfig
}

// NewChainPoller initializes a new ChainPoller.
func NewChainPoller(cfg ChainPollerConfig) (*ChainPoller, error) {
	if cfg.BlockchainClient == nil {
		return nil, errors.New("blockchain client cannot be nil")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &ChainPoller{
		Config: cfg,
	}, nil
}

// Poll fetches logs from the blockchain based on the provided filter queries.
func (cp *ChainPoller) Poll(ctx context.Context, filterQueries []internal.FilterQuery) ([]internal.Log, error) {
	var allLogs []internal.Log

	for _, query := range filterQueries {
		logs, err := cp.Config.BlockchainClient.FilterLogs(ctx, query)
		if err != nil {
			cp.Config.Logger.Error("Failed to filter logs", "error", err, "query", query)
			continue
		}
		allLogs = append(allLogs, logs...)
	}

	return allLogs, nil
}

var _ ChainPollerInterface = (*ChainPoller)(nil)
