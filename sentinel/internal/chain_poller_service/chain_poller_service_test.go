// File: internal/chain_poller_service/chain_poller_service_test.go
package chain_poller_service_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal/chain_poller_service"
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

func setup(t *testing.T) (chain_poller_service.ChainPollerServiceConfig, *internal.MockBlockchainClient) {
	mockBlockchainClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		Logger:           &testLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

	return config, mockBlockchainClient
}

func TestChainPollerService_Initialization(t *testing.T) {
	config, mockBlockchainClient := setup(t)

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Verify initial LastBlock is set correctly
	assert.Equal(t, big.NewInt(99), chainPollerService.LastBlock)

	// Assert that BlockNumber was called once
	mockBlockchainClient.AssertCalled(t, "BlockNumber", mock.Anything)
}

func TestChainPollerService_Initialization_InvalidBlockchainClient(t *testing.T) {
	testLogger := logging.GetTestLogger(t)

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		Logger:           &testLogger,
		ChainID:          1,
		BlockchainClient: nil,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.Error(t, err)
	assert.Nil(t, chainPollerService)
	assert.Equal(t, "blockchain client cannot be nil", err.Error())
}

func TestChainPollerService_StartAndStop(t *testing.T) {
	config, _ := setup(t)

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Start the service
	chainPollerService.Start()

	// Allow some time for polling loop to start
	time.Sleep(10 * time.Millisecond)

	// Stop the service
	chainPollerService.Stop()
}

func TestChainPollerService_PollCycle_FetchAndBroadcast(t *testing.T) {
	config, mockBlockchainClient := setup(t)

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Setup a subscriber
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	logCh, err := chainPollerService.SubscriptionMgr.Subscribe(address, topic)
	require.NoError(t, err)
	defer chainPollerService.SubscriptionMgr.Unsubscribe(address, topic, logCh)

	// Define the expected toBlock
	toBlock := uint64(110)

	// Define the expected filter query
	filterQuery := api.FilterQuery{
		FromBlock: chainPollerService.LastBlock.Uint64() + 1,
		ToBlock:   toBlock,
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}

	// Define fetched logs
	fetchedLogs := []api.Log{
		{
			BlockNumber: 105,
			TxHash:      common.HexToHash("0xdeadbeef"),
			Address:     address,
			Topics:      []common.Hash{topic},
			Data:        []byte("test log data"),
			Index:       0,
		},
	}

	mockBlockchainClient.On("FilterLogs", mock.Anything, filterQuery).Return(fetchedLogs, nil).Once()

	// Mock BlockchainClient.BlockNumber for the next poll
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(toBlock, nil).Once()

	// Start the polling service
	chainPollerService.Start()

	// Allow some time for polling cycle to execute
	time.Sleep(150 * time.Millisecond)

	// Stop the polling service
	chainPollerService.Stop()

	//Assert that the fetched log was broadcasted to the subscriber
	select {
	case receivedLog := <-logCh:
		assert.Equal(t, fetchedLogs[0], receivedLog, "Received log should match the fetched log")
	default:
		t.Fatal("Did not receive the expected log")
	}
}

func TestChainPollerService_PollCycle_NoSubscriptions(t *testing.T) {
	config, mockBlockchainClient := setup(t)

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Mock BlockchainClient.BlockNumber for the next poll
	toBlock := uint64(110)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(toBlock, nil).Once()

	// Start the polling service
	chainPollerService.Start()

	// Allow some time for polling cycle to execute
	time.Sleep(150 * time.Millisecond)

	// Stop the polling service
	chainPollerService.Stop()
}

func TestChainPollerService_StopWithoutStart(t *testing.T) {
	config, _ := setup(t)

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Attempt to stop without starting
	chainPollerService.Stop()
}

func TestChainPollerService_MultipleStartCalls(t *testing.T) {
	config, _ := setup(t)

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Start the service first time
	chainPollerService.Start()

	// Start the service second time
	chainPollerService.Start()

	// Stop the service
	chainPollerService.Stop()
}
