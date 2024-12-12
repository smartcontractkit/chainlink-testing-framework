// File: internal/subscription_manager/subscription_manager_test.go
package subscription_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to initialize a SubscriptionManager for testing
func setupSubscriptionManager(t *testing.T) *SubscriptionManager {
	t.Helper()
	testLogger := logging.GetTestLogger(t)
	return NewSubscriptionManager(SubscriptionManagerConfig{Logger: &testLogger, ChainID: 1})
}

// Helper function to create an EventKey
func createEventKey(address common.Address, topic common.Hash) internal.EventKey {
	return internal.EventKey{Address: address, Topic: topic}
}

// Helper function to create a log event
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

// Helper function to assert registry state
func assertRegistryState(t *testing.T, manager *SubscriptionManager, expectedLength int, expectedEventKeys map[internal.EventKey]int) {
	t.Helper()
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	assert.Len(t, manager.registry, expectedLength, "Registry should have the expected number of EventKeys")
	for ek, expectedCount := range expectedEventKeys {
		actualCount, exists := manager.registry[ek]
		assert.True(t, exists, "Expected EventKey %v to exist", ek)
		assert.Len(t, actualCount, expectedCount, "Expected %d subscribers for EventKey %v, got %d", expectedCount, ek, len(actualCount))
	}
}

func TestSubscriptionManager_Subscribe(t *testing.T) {
	manager := setupSubscriptionManager(t)

	tests := []struct {
		name      string
		address   common.Address
		topic     common.Hash
		expectErr bool
	}{
		{
			name:      "Valid subscription",
			address:   common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			topic:     common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			expectErr: false,
		},
		{
			name:      "Invalid subscription with empty address",
			address:   common.Address{},
			topic:     common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			expectErr: true,
		},
		{
			name:      "Invalid subscription with empty topic",
			address:   common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			topic:     common.Hash{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, err := manager.Subscribe(tt.address, tt.topic)
			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, ch)
			} else {
				require.NoError(t, err)
				require.NotNil(t, ch)
			}
		})
	}

	// After subscriptions, check registry state
	expectedEventKeys := map[internal.EventKey]int{
		createEventKey(common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"), common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")): 1,
	}
	assertRegistryState(t, manager, 1, expectedEventKeys)
}

func TestSubscriptionManager_MultipleSubscribers(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := createEventKey(address, topic)

	// Subscribe first consumer
	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Subscribe second consumer
	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Verify that the list of channels grows
	expectedEventKeys := map[internal.EventKey]int{
		eventKey: 2,
	}
	assertRegistryState(t, manager, 1, expectedEventKeys)

	// Broadcast a log and ensure both channels receive it
	logEvent := createLog(
		1,
		common.HexToHash("0x1234"),
		address,
		[]common.Hash{topic},
		[]byte("log data"),
		0,
	)

	manager.BroadcastLog(eventKey, logEvent)

	receivedLog1 := <-ch1
	receivedLog2 := <-ch2

	assert.Equal(t, logEvent, receivedLog1, "Subscriber 1 should receive the log")
	assert.Equal(t, logEvent, receivedLog2, "Subscriber 2 should receive the log")
}

func TestSubscriptionManager_Unsubscribe(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Unsubscribe existing channel
	err = manager.Unsubscribe(address, topic, ch)
	assert.NoError(t, err)

	// Try unsubscribing again (should fail)
	err = manager.Unsubscribe(address, topic, ch)
	assert.Error(t, err)

	// Unsubscribe non-existent event key
	otherCh := make(chan api.Log)
	err = manager.Unsubscribe(address, topic, otherCh)
	assert.Error(t, err)

	// Check registry state
	expectedEventKeys := map[internal.EventKey]int{}
	assertRegistryState(t, manager, 0, expectedEventKeys)
}

func TestSubscriptionManager_UnsubscribeSelective(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := createEventKey(address, topic)

	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Unsubscribe one consumer and ensure the other remains
	err = manager.Unsubscribe(address, topic, ch1)
	require.NoError(t, err)

	// Check registry state
	expectedEventKeys := map[internal.EventKey]int{
		eventKey: 1,
	}
	assertRegistryState(t, manager, 1, expectedEventKeys)

	// Unsubscribe the last consumer and ensure the registry is cleaned up
	err = manager.Unsubscribe(address, topic, ch2)
	require.NoError(t, err)

	expectedEventKeys = map[internal.EventKey]int{}
	assertRegistryState(t, manager, 0, expectedEventKeys)
}

func TestSubscriptionManager_BroadcastLog(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := createEventKey(address, topic)

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	logEvent := createLog(
		1,
		common.HexToHash("0x1234"),
		address,
		[]common.Hash{topic},
		[]byte("log data"),
		0,
	)

	// Broadcast log event
	manager.BroadcastLog(eventKey, logEvent)

	// Verify the channel received the event
	receivedLog := <-ch
	assert.Equal(t, logEvent, receivedLog, "Subscriber should receive the broadcasted log")
}

func TestSubscriptionManager_BroadcastToAllSubscribers(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address1 := common.HexToAddress("0x9999567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xaaadefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey1 := createEventKey(address1, topic1)

	address2 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic2 := common.HexToHash("0xaaadefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey2 := createEventKey(address2, topic2)

	ch1, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	ch2, err := manager.Subscribe(address2, topic2)
	require.NoError(t, err)

	ch3, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	// Broadcast a log and ensure all channels receive it
	logEvent1 := createLog(
		2,
		common.HexToHash("0x5678"),
		address1,
		[]common.Hash{topic1},
		[]byte("another log data"),
		0,
	)

	logEvent2 := createLog(
		3,
		common.HexToHash("0x2345"),
		address2,
		[]common.Hash{topic2},
		[]byte("another log data 2"),
		0,
	)

	manager.BroadcastLog(eventKey1, logEvent1)
	manager.BroadcastLog(eventKey2, logEvent2)

	receivedLog1 := <-ch1
	receivedLog2 := <-ch2
	receivedLog3 := <-ch3

	assert.Equal(t, logEvent1, receivedLog1, "Subscriber 1 should receive the log")
	assert.Equal(t, logEvent2, receivedLog2, "Subscriber 2 should receive the log")
	assert.Equal(t, logEvent1, receivedLog3, "Subscriber 3 should receive the log")
}

func TestSubscriptionManager_GetAddressesAndTopics(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	ek1 := createEventKey(address1, topic1)
	ek2 := createEventKey(address2, topic2)

	_, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	_, err = manager.Subscribe(address2, topic2)
	require.NoError(t, err)

	// Fetch addresses and topics as EventKeys
	result := manager.GetAddressesAndTopics()

	// Verify the slice contains the expected EventKeys
	assert.Contains(t, result, ek1, "EventKey1 should be in the result")
	assert.Contains(t, result, ek2, "EventKey2 should be in the result")
	assert.Len(t, result, 2, "There should be two unique EventKeys")
}

func TestSubscriptionManager_Cache(t *testing.T) {
	manager := setupSubscriptionManager(t)
	assert.False(t, manager.cacheInitialized, "Cache should not be initialized when Subscription Manager is initialized.")

	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	ek1 := createEventKey(address1, topic1)
	ek2 := createEventKey(address2, topic2)

	// Initialize expected slice of EventKeys
	expectedCache := []internal.EventKey{}

	// Step 1: Subscribe to an event
	_, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)
	assert.False(t, manager.cacheInitialized, "Cache should not be initialized after Subscribe.")

	// Update expected slice
	expectedCache = append(expectedCache, ek1)

	// Verify cache matches expected slice
	cache := manager.GetAddressesAndTopics()
	assert.True(t, manager.cacheInitialized, "Cache should be initialized after GetAddressesAndTopics() is called.")
	assert.ElementsMatch(t, expectedCache, cache, "Cache should match the expected slice of EventKeys.")

	// Step 2: Add another subscription
	_, err = manager.Subscribe(address2, topic2)
	require.NoError(t, err)
	assert.False(t, manager.cacheInitialized, "Cache should be invalidated after Subscribe.")

	// Update expected slice
	expectedCache = append(expectedCache, ek2)

	// Verify cache matches updated slice
	cache = manager.GetAddressesAndTopics()
	assert.True(t, manager.cacheInitialized, "Cache should be reinitialized after GetAddressesAndTopics() is called.")
	assert.ElementsMatch(t, expectedCache, cache, "Cache should match the updated slice of EventKeys.")

	// Step 3: Add a duplicate subscription for address1/topic1
	_, err = manager.Subscribe(address1, topic1)
	require.NoError(t, err)
	assert.False(t, manager.cacheInitialized, "Cache should be invalidated after Subscribe.")

	// No change to expected slice since it's a duplicate subscription

	cache = manager.GetAddressesAndTopics()
	assert.True(t, manager.cacheInitialized, "Cache should be reinitialized after GetAddressesAndTopics() is called.")
	assert.ElementsMatch(t, expectedCache, cache, "Cache should remain unchanged for duplicate subscriptions.")

	// Step 4: Unsubscribe from address2/topic2
	// Retrieve the subscriber channel for ek2
	manager.mu.RLock()
	ch := manager.registry[ek2][0]
	manager.mu.RUnlock()

	err = manager.Unsubscribe(address2, topic2, ch)
	require.NoError(t, err)
	assert.False(t, manager.cacheInitialized, "Cache should be invalidated after Unsubscribe.")

	// Update expected slice
	// Remove ek2 from expectedCache
	for i, ek := range expectedCache {
		if ek.Address == ek2.Address && ek.Topic == ek2.Topic {
			expectedCache = append(expectedCache[:i], expectedCache[i+1:]...)
			break
		}
	}

	// Verify cache matches updated slice
	cache = manager.GetAddressesAndTopics()
	assert.True(t, manager.cacheInitialized, "Cache should be reinitialized after GetAddressesAndTopics() is called.")
	assert.ElementsMatch(t, expectedCache, cache, "Cache should match the updated slice of EventKeys after unsubscription.")

	// Step 5: Unsubscribe from non-existent subscription
	err = manager.Unsubscribe(address2, topic2, ch)
	assert.Error(t, err, "Unsubscribing a non-existent subscription should return an error.")

	// Ensure expected slice remains unchanged
	cache = manager.GetAddressesAndTopics()
	assert.True(t, manager.cacheInitialized, "Cache should remain initialized after an invalid unsubscribe attempt.")
	assert.ElementsMatch(t, expectedCache, cache, "Cache should remain unchanged for invalid unsubscribe attempts.")
	manager.mu.RLock()
	assert.Len(t, manager.registry[ek1], 2, "EventKey should have two subscribers")
	manager.mu.RUnlock()

	// Step 6: Unsubscribe from address1, topic1, ch2
	// Retrieve the second subscriber's channel for ek1
	manager.mu.RLock()
	ch2 := manager.registry[ek1][1]
	manager.mu.RUnlock()

	err = manager.Unsubscribe(address1, topic1, ch2)
	require.NoError(t, err)
	assert.False(t, manager.cacheInitialized, "Cache should be invalidated after Unsubscribe.")

	// Verify cache remains unchanged for the remaining subscriber
	cache = manager.GetAddressesAndTopics()
	assert.True(t, manager.cacheInitialized, "Cache should be reinitialized after GetAddressesAndTopics() is called.")
	assert.ElementsMatch(t, expectedCache, cache, "Cache should remain unchanged for duplicate subscriptions.")
	manager.mu.RLock()
	assert.Len(t, manager.registry[ek1], 1, "EventKey should have one remaining subscriber")
	manager.mu.RUnlock()
}

func TestSubscriptionManager_Close(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Close the SubscriptionManager
	manager.Close()

	// Verify channel is closed
	_, open := <-ch
	assert.False(t, open, "Channel should be closed after Close()")

	// Verify registry is empty
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	assert.Len(t, manager.registry, 0, "Registry should be empty after Close()")
}
