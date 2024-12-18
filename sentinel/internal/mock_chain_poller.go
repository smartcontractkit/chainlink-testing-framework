// File: internal/mock_chain_poller.go
package internal

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/stretchr/testify/mock"
)

// MockChainPoller implements the ChainPollerInterface for testing.
type MockChainPoller struct {
	mock.Mock
}

func (m *MockChainPoller) FilterLogs(ctx context.Context, filterQueries []api.FilterQuery) ([]api.Log, error) {
	args := m.Called(ctx, filterQueries)
	if logs, ok := args.Get(0).([]api.Log); ok {
		return logs, args.Error(1)
	}
	return nil, args.Error(1)
}
