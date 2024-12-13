// File: internal/chain_poller/chain_poller.go
package chain_poller

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
)

// ChainPollerConfig holds the configuration for the ChainPoller.
type ChainPollerConfig struct {
	BlockchainClient api.BlockchainClient
	Logger           *zerolog.Logger
	ChainID          int64
}

// ChainPoller is responsible for polling logs from the blockchain.
type ChainPoller struct {
	blockchainClient api.BlockchainClient
	logger           zerolog.Logger
	chainID          int64
}

// NewChainPoller initializes a new ChainPoller.
func NewChainPoller(cfg ChainPollerConfig) (*ChainPoller, error) {
	if cfg.BlockchainClient == nil {
		return nil, errors.New("blockchain client cannot be nil")
	}
	if cfg.Logger == nil {
		return nil, errors.New("no logger passed")
	}
	if cfg.ChainID < 1 {
		return nil, errors.New("chain ID not set")
	}

	logger := cfg.Logger.With().Str("Component", "ChainPoller").Logger().With().Int64("ChainID", cfg.ChainID).Logger()

	return &ChainPoller{
		blockchainClient: cfg.BlockchainClient,
		logger:           logger,
		chainID:          cfg.ChainID,
	}, nil
}

// Poll fetches logs from the blockchain based on the provided filter queries.
func (cp *ChainPoller) FilterLogs(ctx context.Context, filterQueries []api.FilterQuery) ([]api.Log, error) {
	var allLogs []api.Log

	for _, query := range filterQueries {
		logs, err := cp.blockchainClient.FilterLogs(ctx, query)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				cp.logger.Debug().
					Int64("ChainID", cp.chainID).
					Msg("Log filtering canceled due to shutdown")
				return allLogs, nil
			}
			cp.logger.Error().Err(err).Interface("query", query).Msg("Failed to filter logs")
			continue
		}
		allLogs = append(allLogs, logs...)
	}

	return allLogs, nil
}
