// File: internal/interface.go
package internal

import (
	"context"
)

// BlockchainClient defines the required methods for interacting with a blockchain.
type BlockchainClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, query FilterQuery) ([]Log, error)
}
