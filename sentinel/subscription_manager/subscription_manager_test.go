// File: subscription_manager/subscription_manager_test.go
package subscription_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupSubscriptionManager initializes a SubscriptionManager with a MockLogger for testing.
func setupSubscriptionManager(t *testing.T) *SubscriptionManager {
	testLogger := logging.GetTestLogger(t)
	return NewSubscriptionManager(testLogger, 1)
}

func TestSubscriptionManager_Subscribe(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// Valid subscription
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Invalid subscription with empty address
	_, err = manager.Subscribe(common.Address{}, topic)
	assert.Error(t, err)

	// Invalid subscription with empty topic
	_, err = manager.Subscribe(address, common.Hash{})
	assert.Error(t, err)

	// Check registry state
	manager.registryMutex.RLock()
	defer manager.registryMutex.RUnlock()
	assert.Len(t, manager.registry, 1, "Registry should contain one event key")
	assert.Len(t, manager.registry[internal.EventKey{Address: address, Topic: topic}], 1, "EventKey should have one subscriber")
}

func TestSubscriptionManager_MultipleSubscribers(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	// Subscribe first consumer
	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Subscribe second consumer
	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Verify that the list of channels grows
	manager.registryMutex.RLock()
	subscribers := manager.registry[eventKey]
	manager.registryMutex.RUnlock()
	assert.Len(t, subscribers, 2, "There should be two channels subscribed to the EventKey")

	// Broadcast a log and ensure both channels receive it
	logEvent := internal.Log{
		BlockNumber: 1,
		TxHash:      common.HexToHash("0x1234"),
		Data:        []byte("log data"),
		Address:     address,
		Topics:      []common.Hash{topic},
		Index:       0,
	}

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
	otherCh := make(chan internal.Log)
	err = manager.Unsubscribe(address, topic, otherCh)
	assert.Error(t, err)

	// Check registry state
	manager.registryMutex.RLock()
	defer manager.registryMutex.RUnlock()
	assert.Len(t, manager.registry, 0, "Registry should be empty after unsubscribing")
}

func TestSubscriptionManager_UnsubscribeSelective(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Unsubscribe one consumer and ensure the other remains
	err = manager.Unsubscribe(address, topic, ch1)
	require.NoError(t, err)

	// Check registry state
	manager.registryMutex.RLock()
	subscribers := manager.registry[eventKey]
	manager.registryMutex.RUnlock()

	assert.Len(t, subscribers, 1, "There should be one remaining channel after unsubscription")
	assert.Equal(t, ch2, subscribers[0], "The remaining channel should be the second subscriber")

	// Unsubscribe the last consumer and ensure the registry is cleaned up
	err = manager.Unsubscribe(address, topic, ch2)
	require.NoError(t, err)

	// Check registry state
	manager.registryMutex.RLock()
	_, exists := manager.registry[eventKey]
	manager.registryMutex.RUnlock()

	assert.False(t, exists, "The EventKey should no longer exist in the registry after the last subscriber unsubscribes")
}

func TestSubscriptionManager_BroadcastLog(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	logEvent := internal.Log{
		BlockNumber: 1,
		TxHash:      common.HexToHash("0x1234"),
		Data:        []byte("log data"),
		Address:     address,
		Topics:      []common.Hash{topic},
		Index:       0,
	}

	// Broadcast log event
	manager.BroadcastLog(eventKey, logEvent)

	// Verify the channel received the event
	receivedLog := <-ch
	assert.Equal(t, logEvent, receivedLog, "Subscriber should receive the broadcasted log")
}

func TestSubscriptionManager_BroadcastToAllSubscribers(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	ch3, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Broadcast a log and ensure all channels receive it
	logEvent := internal.Log{
		BlockNumber: 2,
		TxHash:      common.HexToHash("0x5678"),
		Data:        []byte("another log data"),
		Address:     address,
		Topics:      []common.Hash{topic},
		Index:       0,
	}

	manager.BroadcastLog(eventKey, logEvent)

	receivedLog1 := <-ch1
	receivedLog2 := <-ch2
	receivedLog3 := <-ch3

	assert.Equal(t, logEvent, receivedLog1, "Subscriber 1 should receive the log")
	assert.Equal(t, logEvent, receivedLog2, "Subscriber 2 should receive the log")
	assert.Equal(t, logEvent, receivedLog3, "Subscriber 3 should receive the log")
}

func TestSubscriptionManager_GetAddressesAndTopics(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	_, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	_, err = manager.Subscribe(address2, topic2)
	require.NoError(t, err)

	// Fetch addresses and topics
	result := manager.GetAddressesAndTopics()

	// Verify addresses and topics
	assert.Contains(t, result, address1, "Address1 should be in the cache")
	assert.Contains(t, result, address2, "Address2 should be in the cache")
	assert.ElementsMatch(t, result[address1], []common.Hash{topic1}, "Cache should contain topic1 for address1")
	assert.ElementsMatch(t, result[address2], []common.Hash{topic2}, "Cache should contain topic2 for address2")
}

func TestSubscriptionManager_CacheInvalidation(t *testing.T) {
	manager := setupSubscriptionManager(t)

	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	// Subscribe to an initial event
	_, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	// Add another subscription and ensure cache invalidation
	ch, err := manager.Subscribe(address2, topic2)
	require.NoError(t, err)

	// Check updated cache
	updatedCache := manager.GetAddressesAndTopics()
	require.Contains(t, updatedCache, address1, "Address1 should still be in the cache")
	require.Contains(t, updatedCache, address2, "Address2 should now be in the cache")
	assert.ElementsMatch(t, updatedCache[address1], []common.Hash{topic1}, "Cache should still contain topic1 for address1")
	assert.ElementsMatch(t, updatedCache[address2], []common.Hash{topic2}, "Cache should contain topic2 for address2")

	// Add an extra subscription for address1/topic1
	_, err = manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	// Unsubscribe from address2/topic2
	err = manager.Unsubscribe(address2, topic2, ch)
	require.NoError(t, err)

	// Check final cache
	finalCache := manager.GetAddressesAndTopics()
	require.Contains(t, finalCache, address1, "Address1 should still be in the cache")
	assert.ElementsMatch(t, finalCache[address1], []common.Hash{topic1}, "Cache should still contain topic1 for address1")
	require.NotContains(t, finalCache[address2], topic2, "Topic2 should be removed for address2 after unsubscription")
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
	manager.registryMutex.RLock()
	defer manager.registryMutex.RUnlock()
	assert.Len(t, manager.registry, 0, "Registry should be empty after Close()")
}
