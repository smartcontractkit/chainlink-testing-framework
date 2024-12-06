// File: blockchain_client_wrapper/seth_wrapper.go
package sentinel

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// SethClientWrapper wraps a Seth client to implement the BlockchainClient interface.
type SethClientWrapper struct {
	client *seth.Client
}

// NewSethClientWrapper wraps an existing Seth client.
func NewSethClientWrapper(client *seth.Client) internal.BlockchainClient {
	return &SethClientWrapper{client: client}
}

// Ensure SethClientWrapper implements BlockchainClient.
var _ internal.BlockchainClient = (*SethClientWrapper)(nil)

// BlockNumber retrieves the latest block number.
func (s *SethClientWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return s.client.Client.BlockNumber(ctx)
}

// BlockByNumber retrieves a block by number.
func (s *SethClientWrapper) FilterLogs(ctx context.Context, query internal.FilterQuery) ([]internal.Log, error) {
	fromBlock := new(big.Int).SetUint64(query.FromBlock)
	toBlock := new(big.Int).SetUint64(query.ToBlock)

	// Convert internal.FilterQuery to ethereum.FilterQuery
	ethQuery := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}

	// Fetch logs using Seth's client
	ethLogs, err := s.client.Client.FilterLogs(ctx, ethQuery)
	if err != nil {
		return nil, err
	}

	// Convert []types.Log to []internal.Log
	internalLogs := make([]internal.Log, len(ethLogs))
	for i, ethLog := range ethLogs {
		internalLogs[i] = internal.Log{
			Address:     ethLog.Address,
			Topics:      ethLog.Topics,
			Data:        ethLog.Data,
			BlockNumber: ethLog.BlockNumber,
			TxHash:      ethLog.TxHash,
			Index:       ethLog.Index,
		}
	}

	return internalLogs, nil
}
