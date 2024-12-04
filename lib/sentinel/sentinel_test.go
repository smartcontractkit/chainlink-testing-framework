// File: sentinel_test.go
package sentinel_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel"
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

func setupSentinel() (*sentinel.Sentinel, *internal.MockLogger) {
	logger := internal.NewMockLogger()
	s := sentinel.NewSentinel(sentinel.SentinelConfig{
		Logger: logger,
	})

	return s, logger
}

func setupChainPollerServiceConfig(l *internal.MockLogger, chainID int64) (*chain_poller_service.ChainPollerServiceConfig, *internal.MockBlockchainClient) {
	mockBlockchainClient := new(internal.MockBlockchainClient)

	mockChainPoller := new(MockChainPoller)

	return &chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     100 * time.Millisecond,
		ChainPoller:      mockChainPoller,
		Logger:           l,
		ChainID:          chainID,
		BlockchainClient: mockBlockchainClient,
	}, mockBlockchainClient
}

func TestNewSentinel_NoErrors(t *testing.T) {
	s, _ := setupSentinel()
	defer s.Close()
	require.NotNil(t, s, "Sentinel should not be nil")
}

func TestAddRemoveChain(t *testing.T) {
	s, logger := setupSentinel()
	defer s.Close()

	config1, mockClient1 := setupChainPollerServiceConfig(logger, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChainPollerServiceConfig(logger, 1)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	// Add and remove a chain
	require.NoError(t, s.AddChain(*config1), "Should add chain without error")
	require.NoError(t, s.RemoveChain(1), "Should remove chain without error")

	// Add another chain
	require.NoError(t, s.AddChain(*config2), "Should add another chain without error")
}

func TestAddChain_SubscribeUnsubscribeEvent(t *testing.T) {
	s, logger := setupSentinel()
	defer s.Close()

	config, mockClient1 := setupChainPollerServiceConfig(logger, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	require.NoError(t, s.AddChain(*config), "Should add chain without error")

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

	ch, err := s.Subscribe(1, address, topic)
	require.NoError(t, err, "Should subscribe without error")
	defer s.Unsubscribe(1, address, topic, ch)

	// Wait to ensure subscription is ready
	time.Sleep(50 * time.Millisecond)

	// Simulate log broadcast
	eventKey := internal.EventKey{Address: address, Topic: topic}
	log := internal.Log{Address: address, Topics: []common.Hash{topic}, Data: []byte("event data")}

	chainService, exists := s.GetService(1) // Add a helper to retrieve the service.
	require.True(t, exists, "Chain service should exist")

	chainService.SubscriptionMgr.BroadcastLog(eventKey, log)

	select {
	case receivedLog := <-ch:
		assert.Equal(t, log, receivedLog, "Received log should match the broadcasted log")
	default:
		t.Fatal("Log not received")
	}

	require.NoError(t, s.Unsubscribe(1, address, topic, ch), "Should unsubscribe without error")
}

func TestAddChains_MultipleConsumers(t *testing.T) {
	s, logger := setupSentinel()
	defer s.Close()

	config1, mockClient1 := setupChainPollerServiceConfig(logger, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChainPollerServiceConfig(logger, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	require.NoError(t, s.AddChain(*config1), "Should add chain 1 without error")
	require.NoError(t, s.AddChain(*config2), "Should add chain 2 without error")

	// Chain 1 subscribers
	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

	ch1, err := s.Subscribe(1, address1, topic1)
	require.NoError(t, err, "Subscriber 1 should subscribe without error")
	defer s.Unsubscribe(1, address1, topic1, ch1)

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	ch2, err := s.Subscribe(1, address2, topic2)
	require.NoError(t, err, "Subscriber 2 should subscribe without error")
	defer s.Unsubscribe(1, address2, topic2, ch2)

	// Chain 2 subscriber
	address3 := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	topic3 := common.HexToHash("0xcafebabecafebabecafebabecafebabecafebabe")

	ch3, err := s.Subscribe(2, address3, topic3)
	require.NoError(t, err, "Subscriber 3 should subscribe without error")
	defer s.Unsubscribe(2, address3, topic3, ch3)

	// Broadcast events
	eventKey1 := internal.EventKey{Address: address1, Topic: topic1}
	log1 := internal.Log{Address: address1, Topics: []common.Hash{topic1}, Data: []byte("log1")}
	chainService, exists := s.GetService(1)
	require.True(t, exists, "Chain service should exist")
	chainService.SubscriptionMgr.BroadcastLog(eventKey1, log1)

	eventKey3 := internal.EventKey{Address: address3, Topic: topic3}
	log3 := internal.Log{Address: address3, Topics: []common.Hash{topic3}, Data: []byte("log3")}
	chainService2, exists := s.GetService(2)
	require.True(t, exists, "Chain service should exist")
	chainService2.SubscriptionMgr.BroadcastLog(eventKey3, log3)

	select {
	case receivedLog := <-ch1:
		assert.Equal(t, log1, receivedLog, "Subscriber 1 should receive the correct log")
	default:
		t.Fatal("Subscriber 1 did not receive the log")
	}

	select {
	case receivedLog := <-ch3:
		assert.Equal(t, log3, receivedLog, "Subscriber 3 should receive the correct log")
	default:
		t.Fatal("Subscriber 3 did not receive the log")
	}
}

func TestAddChains_RemoveAndValidate(t *testing.T) {
	s, logger := setupSentinel()
	defer s.Close()

	config1, mockClient1 := setupChainPollerServiceConfig(logger, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChainPollerServiceConfig(logger, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	require.NoError(t, s.AddChain(*config1), "Should add chain 1 without error")
	require.NoError(t, s.AddChain(*config2), "Should add chain 2 without error")

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

	ch, err := s.Subscribe(1, address, topic)
	require.NoError(t, err, "Should subscribe without error")
	defer s.Unsubscribe(1, address, topic, ch)

	require.NoError(t, s.RemoveChain(1), "Should remove chain 1 without error")

	select {
	case <-ch:
		t.Fatal("Channel should be closed after chain removal")
	default:
	}
}

func TestAddMultipleChains_CloseSentinel(t *testing.T) {
	s, logger := setupSentinel()

	config1, mockClient1 := setupChainPollerServiceConfig(logger, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChainPollerServiceConfig(logger, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	require.NoError(t, s.AddChain(*config1), "Should add chain 1 without error")
	require.NoError(t, s.AddChain(*config2), "Should add chain 2 without error")

	s.Close()

	assert.False(t, s.HasServices(), "All chains should be cleaned up after Close")
}
