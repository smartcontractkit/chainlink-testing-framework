// File: chain_poller_service/chain_poller_service_test.go
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

	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/chain_poller"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/chain_poller_service"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/internal"
)

// MockChainPoller implements the ChainPollerInterface for testing.
type MockChainPoller struct {
	mock.Mock
}

func (m *MockChainPoller) Poll(ctx context.Context, filterQueries []internal.FilterQuery) ([]internal.Log, error) {
	args := m.Called(ctx, filterQueries)
	if logs, ok := args.Get(0).([]internal.Log); ok {
		return logs, args.Error(1)
	}
	return nil, args.Error(1)
}

// Ensure MockChainPoller implements ChainPollerInterface
var _ chain_poller.ChainPollerInterface = (*MockChainPoller)(nil)

func TestChainPollerService_Initialization(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Verify initial LastBlock is set correctly
	assert.Equal(t, big.NewInt(99), chainPollerService.LastBlock)

	// Assert that BlockNumber was called once
	mockBlockchainClient.AssertCalled(t, "BlockNumber", mock.Anything)
}

func TestChainPollerService_Initialization_InvalidConfig(t *testing.T) {
	mockLogger := internal.NewMockLogger()

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval: 100 * time.Millisecond,
		ChainPoller:  nil, // Invalid
		Logger:       mockLogger,
		ChainID:      1,
		// BlockchainClient is missing
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.Error(t, err)
	assert.Nil(t, chainPollerService)
	assert.Equal(t, "chain poller cannot be nil", err.Error())
}

func TestChainPollerService_Initialization_InvalidBlockchainClient(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockLogger := internal.NewMockLogger()

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: nil,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.Error(t, err)
	assert.Nil(t, chainPollerService)
	assert.Equal(t, "blockchain client cannot be nil", err.Error())
}

func TestChainPollerService_StartAndStop(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Start the service
	chainPollerService.Start()

	// Allow some time for polling loop to start
	time.Sleep(10 * time.Millisecond)

	// Stop the service
	chainPollerService.Stop()

	assert.True(t, mockLogger.ContainsLog("ChainPollerService started with poll interval: 100ms"))
	assert.True(t, mockLogger.ContainsLog("Polling loop terminating"))
	assert.True(t, mockLogger.ContainsLog("ChainPollerService stopped"))
}

func TestChainPollerService_PollCycle_FetchAndBroadcast(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	// Initialize ChainPollerService
	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

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
	filterQuery := internal.FilterQuery{
		FromBlock: chainPollerService.LastBlock.Uint64() + 1,
		ToBlock:   toBlock,
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}

	// Define fetched logs
	fetchedLogs := []internal.Log{
		{
			BlockNumber: 105,
			TxHash:      common.HexToHash("0xdeadbeef"),
			Address:     address,
			Topics:      []common.Hash{topic},
			Data:        []byte("test log data"),
			Index:       0,
		},
	}

	mockChainPoller.On("Poll", mock.Anything, []internal.FilterQuery{filterQuery}).Return(fetchedLogs, nil).Once()

	// Mock BlockchainClient.BlockNumber for the next poll
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(toBlock, nil).Once()

	// Start the polling service
	chainPollerService.Start()

	// Allow some time for polling cycle to execute
	time.Sleep(150 * time.Millisecond)

	// Stop the polling service
	chainPollerService.Stop()

	// Assert that the fetched log was broadcasted to the subscriber
	select {
	case receivedLog := <-logCh:
		assert.Equal(t, fetchedLogs[0], receivedLog, "Received log should match the fetched log")
	default:
		t.Fatal("Did not receive the expected log")
	}

	assert.True(t, mockLogger.ContainsLog("ChainPollerService started with poll interval: 100ms"))
	assert.True(t, mockLogger.ContainsLog("Starting polling cycle"))
	assert.True(t, mockLogger.ContainsLog("Fetched 1 logs from blockchain"))
	assert.True(t, mockLogger.ContainsLog("Completed polling cycle in"))
}

func TestChainPollerService_PollCycle_NoSubscriptions(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	// Initialize ChainPollerService
	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

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

	// Assert that Poll was not called
	mockChainPoller.AssertNotCalled(t, "Poll", mock.Anything, mock.Anything)

	// Assert that the logger logged the no active subscriptions message
	assert.True(t, mockLogger.ContainsLog("No active subscriptions, skipping polling cycle"))
}

func TestChainPollerService_PollCycle_MultipleLogs(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	// Initialize ChainPollerService
	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Verify initial LastBlock is set correctly
	assert.Equal(t, big.NewInt(99), chainPollerService.LastBlock)

	// Setup subscribers
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	logCh1, err := chainPollerService.SubscriptionMgr.Subscribe(address, topic1)
	require.NoError(t, err)
	defer chainPollerService.SubscriptionMgr.Unsubscribe(address, topic1, logCh1)

	logCh2, err := chainPollerService.SubscriptionMgr.Subscribe(address, topic2)
	require.NoError(t, err)
	defer chainPollerService.SubscriptionMgr.Unsubscribe(address, topic2, logCh2)

	// Define the expected toBlock
	toBlock := uint64(110)

	// Define the expected filter queries (same fromBlock and toBlock)
	filterQueries := []internal.FilterQuery{
		{
			FromBlock: 100,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
			Topics:    [][]common.Hash{{topic1}},
		},
		{
			FromBlock: 100,
			ToBlock:   toBlock,
			Addresses: []common.Address{address},
			Topics:    [][]common.Hash{{topic2}},
		},
	}

	// Define fetched logs
	fetchedLogs := []internal.Log{
		{
			BlockNumber: 105,
			TxHash:      common.HexToHash("0xdeadbeef"),
			Address:     address,
			Topics:      []common.Hash{topic1},
			Data:        []byte("test log data 1"),
			Index:       0,
		},
		{
			BlockNumber: 106,
			TxHash:      common.HexToHash("0xfeedface"),
			Address:     address,
			Topics:      []common.Hash{topic2},
			Data:        []byte("test log data 2"),
			Index:       1,
		},
	}

	// Mock ChainPoller.Poll to return fetchedLogs
	mockChainPoller.On("Poll", mock.Anything, filterQueries).Return(fetchedLogs, nil).Once()

	// Mock BlockchainClient.BlockNumber for the next poll
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(toBlock, nil).Once()

	// Start the polling service
	chainPollerService.Start()

	// Allow some time for polling cycle to execute
	time.Sleep(150 * time.Millisecond)
	chainPollerService.Stop()

	// Assert that the fetched logs were broadcasted to the subscribers
	select {
	case receivedLog := <-logCh1:
		assert.Equal(t, fetchedLogs[0], receivedLog, "Received log1 should match the fetched log1")
	default:
		t.Fatal("Did not receive the expected log1")
	}

	select {
	case receivedLog := <-logCh2:
		assert.Equal(t, fetchedLogs[1], receivedLog, "Received log2 should match the fetched log2")
	default:
		t.Fatal("Did not receive the expected log2")
	}

	// Assert that the mocks were called as expected
	mockChainPoller.AssertCalled(t, "Poll", mock.Anything, filterQueries)
	mockBlockchainClient.AssertCalled(t, "BlockNumber", mock.Anything)
	assert.True(t, mockLogger.ContainsLog("Starting polling cycle"))
	assert.True(t, mockLogger.ContainsLog("Fetched 2 logs from blockchain"))
	assert.True(t, mockLogger.ContainsLog("Completed polling cycle in"))
}

func TestChainPollerService_StopWithoutStart(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	// Initialize ChainPollerService
	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Attempt to stop without starting
	chainPollerService.Stop()
	// Verify that no logger calls were made
	assert.Equal(t, 0, mockLogger.NumLogs())
}

func TestChainPollerService_MultipleStartCalls(t *testing.T) {
	mockChainPoller := new(MockChainPoller)
	mockBlockchainClient := new(internal.MockBlockchainClient)
	mockLogger := internal.NewMockLogger()

	// Mock BlockchainClient.BlockNumber during initialization
	initialLastBlockNum := uint64(100)
	mockBlockchainClient.On("BlockNumber", mock.Anything).Return(initialLastBlockNum, nil).Once()

	config := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           mockLogger,
		ChainID:          1,
		BlockchainClient: mockBlockchainClient,
	}

	chainPollerService, err := chain_poller_service.NewChainPollerService(config)
	require.NoError(t, err)
	require.NotNil(t, chainPollerService)

	// Start the service first time
	chainPollerService.Start()

	// Start the service second time
	chainPollerService.Start()

	// Stop the service
	chainPollerService.Stop()

	assert.True(t, mockLogger.ContainsLog("ChainPollerService started with poll interval: 100ms"))
	assert.True(t, mockLogger.ContainsLog("ChainPollerService already started"))
	assert.True(t, mockLogger.ContainsLog("Polling loop terminating"))
	assert.True(t, mockLogger.ContainsLog("ChainPollerService stopped"))
}
