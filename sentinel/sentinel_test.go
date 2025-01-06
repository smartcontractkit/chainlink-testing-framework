// File: sentinel_test.go
package sentinel

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
)

// Helper function to initialize a Sentinel instance for testing.
func setupSentinel(t *testing.T) *Sentinel {
	t.Helper()
	s := NewSentinel(SentinelConfig{t: t})

	// Ensure Sentinel is closed after the test.
	t.Cleanup(func() {
		s.Close()
	})

	return s
}

// Helper function to set up a chain with a mock blockchain client.
func setupChain(t *testing.T, chainID int64) (*AddChainConfig, *internal.MockBlockchainClient) {
	t.Helper()
	mockClient := new(internal.MockBlockchainClient)

	config := &AddChainConfig{
		BlockchainClient: mockClient,
		PollInterval:     100 * time.Millisecond,
		ChainID:          chainID,
	}

	return config, mockClient
}

// Helper function to create an EventKey.
func createEventKey(address common.Address, topic common.Hash) internal.EventKey {
	return internal.EventKey{Address: address, Topic: topic}
}

// Helper function to create a log event.
func createLog(blockNumber uint64, txHash common.Hash, address common.Address, topics []common.Hash, data []byte, index uint) api.Log {
	return api.Log{
		BlockNumber: blockNumber,
		TxHash:      txHash,
		Address:     address,
		Topics:      topics,
		Data:        data,
		Index:       index,
	}
}

func TestNewSentinel_NoErrors(t *testing.T) {
	s := setupSentinel(t)
	require.NotNil(t, s, "Sentinel should not be nil")
}

func TestAddRemoveChain(t *testing.T) {
	s := setupSentinel(t)

	// Setup two chains with the same ChainID to test removal.
	config1, mockClient1 := setupChain(t, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChain(t, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	// Add first chain.
	require.NoError(t, s.AddChain(*config1), "Should add chain without error")

	// Add second chain.
	require.NoError(t, s.AddChain(*config2), "Should add another chain without error")

	// Remove first chain.
	require.NoError(t, s.RemoveChain(1), "Should remove chain without error")

	// Attempt to add a chain with the same ChainID again.
	require.Error(t, s.AddChain(*config2), "chain with ID 2 already exists")
}

func TestAddChain_SubscribeUnsubscribeEvent(t *testing.T) {
	s := setupSentinel(t)

	config, mockClient := setupChain(t, 1)
	mockClient.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	// Add chain.
	require.NoError(t, s.AddChain(*config), "Should add chain without error")

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

	// Subscribe to an event.
	ch, err := s.Subscribe(1, address, topic)
	require.NoError(t, err, "Should subscribe without error")
	defer s.Unsubscribe(1, address, topic, ch)

	// Wait briefly to ensure subscription is registered.
	time.Sleep(50 * time.Millisecond)

	// Simulate log broadcast.
	eventKey := createEventKey(address, topic)
	log := createLog(1, common.HexToHash("0x1234"), address, []common.Hash{topic}, []byte("event data"), 0)

	// Retrieve the chain service to access the Subscription Manager.
	chainService, exists := s.services[1]
	require.True(t, exists, "Chain service should exist")

	// Broadcast the log.
	chainService.SubscriptionMgr.BroadcastLog(eventKey, log)

	// Verify the subscriber receives the log.
	select {
	case receivedLog := <-ch:
		assert.Equal(t, log, receivedLog, "Received log should match the broadcasted log")
	case <-time.After(1 * time.Second):
		t.Fatal("Subscriber did not receive the log")
	}

	// Unsubscribe and ensure no errors.
	require.NoError(t, s.Unsubscribe(1, address, topic, ch), "Should unsubscribe without error")
}

func TestAddChains_MultipleConsumers(t *testing.T) {
	s := setupSentinel(t)

	// Setup two different chains.
	config1, mockClient1 := setupChain(t, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChain(t, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	// Add both chains.
	require.NoError(t, s.AddChain(*config1), "Should add chain 1 without error")
	require.NoError(t, s.AddChain(*config2), "Should add chain 2 without error")

	// Subscribe to events on chain 1.
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

	// Subscribe to an event on chain 2.
	address3 := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	topic3 := common.HexToHash("0xcafebabecafebabecafebabecafebabecafebabe")

	ch3, err := s.Subscribe(2, address3, topic3)
	require.NoError(t, err, "Subscriber 3 should subscribe without error")
	defer s.Unsubscribe(2, address3, topic3, ch3)

	// Broadcast logs to both chains.
	logEvent1 := createLog(2, common.HexToHash("0x5678"), address1, []common.Hash{topic1}, []byte("another log data"), 0)
	logEvent2 := createLog(3, common.HexToHash("0x2345"), address2, []common.Hash{topic2}, []byte("another log data 2"), 0)
	logEvent3 := createLog(4, common.HexToHash("0x3456"), address3, []common.Hash{topic3}, []byte("another log data 3"), 0)

	chainService1, exists1 := s.services[1]
	require.True(t, exists1, "Chain service 1 should exist")
	chainService2, exists2 := s.services[2]
	require.True(t, exists2, "Chain service 2 should exist")

	chainService1.SubscriptionMgr.BroadcastLog(createEventKey(address1, topic1), logEvent1)
	chainService1.SubscriptionMgr.BroadcastLog(createEventKey(address2, topic2), logEvent2)
	chainService2.SubscriptionMgr.BroadcastLog(createEventKey(address3, topic3), logEvent3)

	// Verify subscribers receive their respective logs.
	select {
	case receivedLog := <-ch1:
		assert.Equal(t, logEvent1, receivedLog, "Subscriber 1 should receive the correct log")
	case <-time.After(1 * time.Second):
		t.Fatal("Subscriber 1 did not receive the log")
	}

	select {
	case receivedLog := <-ch2:
		assert.Equal(t, logEvent2, receivedLog, "Subscriber 2 should receive the correct log")
	case <-time.After(1 * time.Second):
		t.Fatal("Subscriber 2 did not receive the log")
	}

	select {
	case receivedLog := <-ch3:
		assert.Equal(t, logEvent3, receivedLog, "Subscriber 3 should receive the correct log")
	case <-time.After(1 * time.Second):
		t.Fatal("Subscriber 3 did not receive the log")
	}
}

func TestAddChains_RemoveAndValidate(t *testing.T) {
	s := setupSentinel(t)

	// Setup two chains.
	config1, mockClient1 := setupChain(t, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChain(t, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	// Add both chains.
	require.NoError(t, s.AddChain(*config1), "Should add chain 1 without error")
	require.NoError(t, s.AddChain(*config2), "Should add chain 2 without error")

	// Subscribe to an event on chain 1.
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

	ch, err := s.Subscribe(1, address, topic)
	require.NoError(t, err, "Should subscribe without error")
	defer s.Unsubscribe(1, address, topic, ch)

	// Remove chain 1.
	require.NoError(t, s.RemoveChain(1), "Should remove chain 1 without error")

	// Verify that the subscriber's channel is closed.
	select {
	case _, open := <-ch:
		assert.False(t, open, "Channel should be closed after chain removal")
	default:
		t.Fatal("Channel was not closed after chain removal")
	}
}

func TestAddMultipleChains_CloseSentinel(t *testing.T) {
	s := setupSentinel(t)

	// Setup two chains.
	config1, mockClient1 := setupChain(t, 1)
	mockClient1.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	config2, mockClient2 := setupChain(t, 2)
	mockClient2.On("BlockNumber", mock.Anything).Return(uint64(100), nil).Once()

	// Add both chains.
	require.NoError(t, s.AddChain(*config1), "Should add chain 1 without error")
	require.NoError(t, s.AddChain(*config2), "Should add chain 2 without error")

	// Close Sentinel.
	s.Close()

	// Verify that all chains are cleaned up.
	assert.False(t, len(s.services) > 0, "All chains should be cleaned up after Close")
}
