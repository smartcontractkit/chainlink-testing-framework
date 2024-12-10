// File: blockchain_client_wrapper/geth_wrapper.go
package blockchain_client_wrapper

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
)

// GethWrapper wraps a Geth client to implement the BlockchainClient interface.
type GethWrapper struct {
	client *ethclient.Client
}

// NewGethClientWrapper wraps an existing Geth client.
func NewGethClientWrapper(client *ethclient.Client) api.BlockchainClient {
	return &GethWrapper{client: client}
}

// BlockNumber retrieves the latest block number.
func (g *GethWrapper) BlockNumber(ctx context.Context) (uint64, error) {
	return g.client.BlockNumber(ctx)
}

// BlockByNumber retrieves a block by number.
func (g *GethWrapper) FilterLogs(ctx context.Context, query api.FilterQuery) ([]api.Log, error) {
	fromBlock := new(big.Int).SetUint64(query.FromBlock)
	toBlock := new(big.Int).SetUint64(query.ToBlock)

	// Convert FilterQuery to ethereum.FilterQuery
	ethQuery := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}

	// Fetch logs using geth client
	ethLogs, err := g.client.FilterLogs(ctx, ethQuery)
	if err != nil {
		return nil, err
	}

	// Convert []types.Log to []Log
	internalLogs := make([]api.Log, len(ethLogs))
	for i, ethLog := range ethLogs {
		internalLogs[i] = api.Log{
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
