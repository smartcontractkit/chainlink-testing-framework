// File: subscription_manager/subscription_manager_test.go
package subscription_manager

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupSubscriptionManager initializes a SubscriptionManager with a MockLogger for testing.
func setupSubscriptionManager() *SubscriptionManager {
	mockLogger := internal.NewMockLogger()
	return NewSubscriptionManager(mockLogger, 1)
}

func TestSubscriptionManager_Subscribe(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// Reset logs before Subscribe operation
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

	// Valid subscription
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Assert the expected log message after Subscribe
	expectedCacheLog := "ChainID=1 Cache invalidated due to subscription change"
	expectedSubscriptionLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedCacheLog), "Expected cache invalidation log")
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog), "Expected log message for valid subscription")

	// Reset logs before invalid subscription attempts
	mockLogger.Reset()

	// Invalid subscription with empty address
	_, err = manager.Subscribe(common.Address{}, topic)
	assert.Error(t, err)
	expectedWarn1 := "Attempted to subscribe with an empty address"
	assert.True(t, mockLogger.ContainsLog(expectedWarn1), "Expected warning for empty address")

	// Reset logs before next invalid subscription
	mockLogger.Reset()

	// Invalid subscription with empty topic
	_, err = manager.Subscribe(address, common.Hash{})
	assert.Error(t, err)
	expectedWarn2 := "Attempted to subscribe with an empty topic"
	assert.True(t, mockLogger.ContainsLog(expectedWarn2), "Expected warning for empty topic")

	// Check registry state
	manager.registryMutex.RLock()
	defer manager.registryMutex.RUnlock()
	assert.Len(t, manager.registry, 1, "Registry should contain one event key")
	assert.Len(t, manager.registry[internal.EventKey{Address: address, Topic: topic}], 1, "EventKey should have one subscriber")
}

func TestSubscriptionManager_MultipleSubscribers(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	// Reset logs before subscriptions
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

	// Subscribe first consumer
	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for first subscription
	expectedSubscriptionLog1 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog1), "Expected log message for first subscription")

	// Reset logs before second subscription
	mockLogger.Reset()

	// Subscribe second consumer
	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for second subscription
	expectedSubscriptionLog2 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=2 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog2), "Expected log message for second subscription")

	// Reset logs before broadcasting
	mockLogger.Reset()

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

	// Assert broadcast log message
	expectedBroadcastLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s Log broadcasted to all subscribers",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedBroadcastLog), "Expected broadcast log message")
}

func TestSubscriptionManager_Unsubscribe(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Reset logs before Unsubscribe operation
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

	// Unsubscribe existing channel
	err = manager.Unsubscribe(address, topic, ch)
	assert.NoError(t, err)

	// Assert the expected log message after Unsubscribe
	expectedRemoveLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s RemainingSubscribers=0 Subscription removed",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedRemoveLog), "Expected log message for unsubscribing")

	// Reset logs before attempting to unsubscribe again
	mockLogger.Reset()

	// Try unsubscribing again (should fail)
	err = manager.Unsubscribe(address, topic, ch)
	assert.Error(t, err)

	// Assert the expected warning log message
	expectedWarn := "Attempted to unsubscribe from a non-existent EventKey"
	assert.True(t, mockLogger.ContainsLog(expectedWarn), "Expected warning for unsubscribing a non-existent subscriber")

	// Reset logs before unsubscribing a non-existent EventKey
	mockLogger.Reset()

	// Unsubscribe non-existent event key
	otherCh := make(chan internal.Log)
	err = manager.Unsubscribe(address, topic, otherCh)
	assert.Error(t, err)

	// Assert the expected warning log message
	assert.True(t, mockLogger.ContainsLog(expectedWarn), "Expected warning for unsubscribing a non-existent subscriber")

	// Check registry state
	manager.registryMutex.RLock()
	defer manager.registryMutex.RUnlock()
	assert.Len(t, manager.registry, 0, "Registry should be empty after unsubscribing")
}

func TestSubscriptionManager_UnsubscribeSelective(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	// Subscribe multiple consumers to the same EventKey
	mockLogger := manager.logger.(*internal.MockLogger)

	// Reset logs before subscriptions
	mockLogger.Reset()

	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for first subscription
	expectedSubscriptionLog1 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog1), "Expected log message for first subscription")

	// Reset logs before second subscription
	mockLogger.Reset()

	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for second subscription
	expectedSubscriptionLog2 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=2 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog2), "Expected log message for second subscription")

	// Reset logs before Unsubscribe operation
	mockLogger.Reset()

	// Unsubscribe one consumer and ensure the other remains
	err = manager.Unsubscribe(address, topic, ch1)
	require.NoError(t, err)

	// Assert log for selective unsubscription
	expectedSelectiveRemoveLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s RemainingSubscribers=1 Subscription removed",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSelectiveRemoveLog), "Expected log message for selective unsubscription")

	// Check registry state
	manager.registryMutex.RLock()
	subscribers := manager.registry[eventKey]
	manager.registryMutex.RUnlock()

	assert.Len(t, subscribers, 1, "There should be one remaining channel after unsubscription")
	assert.Equal(t, ch2, subscribers[0], "The remaining channel should be the second subscriber")

	// Reset logs before final Unsubscribe operation
	mockLogger.Reset()

	// Unsubscribe the last consumer and ensure the registry is cleaned up
	err = manager.Unsubscribe(address, topic, ch2)
	require.NoError(t, err)

	// Assert log for final unsubscription
	expectedFinalRemoveLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s RemainingSubscribers=0 Subscription removed",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedFinalRemoveLog), "Expected log message for final unsubscription")

	// Check registry state
	manager.registryMutex.RLock()
	_, exists := manager.registry[eventKey]
	manager.registryMutex.RUnlock()

	assert.False(t, exists, "The EventKey should no longer exist in the registry after the last subscriber unsubscribes")
}

func TestSubscriptionManager_BroadcastLog(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Reset logs before BroadcastLog operation
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

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

	// Assert broadcast log message
	expectedBroadcastLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s Log broadcasted to all subscribers",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedBroadcastLog), "Expected log message for broadcasting")
}

func TestSubscriptionManager_BroadcastToAllSubscribers(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	eventKey := internal.EventKey{Address: address, Topic: topic}

	// Subscribe multiple consumers to the same EventKey
	mockLogger := manager.logger.(*internal.MockLogger)

	// Reset logs before subscriptions
	mockLogger.Reset()

	ch1, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for first subscription
	expectedSubscriptionLog1 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog1), "Expected log message for first subscription")

	// Reset logs before second subscription
	mockLogger.Reset()

	ch2, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for second subscription
	expectedSubscriptionLog2 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=2 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog2), "Expected log message for second subscription")

	// Reset logs before third subscription
	mockLogger.Reset()

	ch3, err := manager.Subscribe(address, topic)
	require.NoError(t, err)

	// Assert log for third subscription
	expectedSubscriptionLog3 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=3 New subscription added",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog3), "Expected log message for third subscription")

	// Reset logs before broadcasting
	mockLogger.Reset()

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

	// Assert broadcast log message
	expectedBroadcastLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s Log broadcasted to all subscribers",
		address.Hex(),
		topic.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedBroadcastLog), "Expected log message for broadcasting to all subscribers")
}

func TestSubscriptionManager_GetAddressesAndTopics(t *testing.T) {
	manager := setupSubscriptionManager()

	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	// Reset logs before subscriptions
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

	_, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	// Assert log for first subscription
	expectedSubscriptionLog1 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address1.Hex(),
		topic1.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog1), "Expected log message for first subscription")

	// Reset logs before second subscription
	mockLogger.Reset()

	_, err = manager.Subscribe(address2, topic2)
	require.NoError(t, err)

	// Assert log for second subscription
	expectedSubscriptionLog2 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address2.Hex(),
		topic2.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog2), "Expected log message for second subscription")

	// Fetch addresses and topics
	result := manager.GetAddressesAndTopics()

	// Verify addresses and topics
	assert.Contains(t, result, address1, "Address1 should be in the cache")
	assert.Contains(t, result, address2, "Address2 should be in the cache")
	assert.ElementsMatch(t, result[address1], []common.Hash{topic1}, "Cache should contain topic1 for address1")
	assert.ElementsMatch(t, result[address2], []common.Hash{topic2}, "Cache should contain topic2 for address2")

	// Assert cache log message
	expectedCacheLog := fmt.Sprintf(
		"ChainID=1 UniqueAddresses=%d Cached address-topic pairs",
		len(result),
	)
	assert.True(t, mockLogger.ContainsLog(expectedCacheLog), "Expected cache initialization log message")
}

func TestSubscriptionManager_CacheInvalidation(t *testing.T) {
	manager := setupSubscriptionManager()

	address1 := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic1 := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	address2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")
	topic2 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	// Reset logs before first subscription
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

	// Subscribe to an initial event
	_, err := manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	// Assert log for first subscription
	expectedSubscriptionLog1 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address1.Hex(),
		topic1.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog1), "Expected log message for first subscription")

	// Assert cache initialization log
	expectedCacheLog := "ChainID=1 Cache invalidated due to subscription change"
	assert.True(t, mockLogger.ContainsLog(expectedCacheLog), "Expected cache invalidation log after first subscription")

	// Reset logs before second subscription
	mockLogger.Reset()

	// Add another subscription and ensure cache invalidation
	ch, err := manager.Subscribe(address2, topic2)
	require.NoError(t, err)

	// Assert log for second subscription
	expectedSubscriptionLog2 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=1 New subscription added",
		address2.Hex(),
		topic2.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog2), "Expected log message for second subscription")

	// Assert cache invalidation log
	assert.True(t, mockLogger.ContainsLog(expectedCacheLog), "Expected cache invalidation log after second subscription")

	// Check updated cache
	updatedCache := manager.GetAddressesAndTopics()
	require.Contains(t, updatedCache, address1, "Address1 should still be in the cache")
	require.Contains(t, updatedCache, address2, "Address2 should now be in the cache")
	assert.ElementsMatch(t, updatedCache[address1], []common.Hash{topic1}, "Cache should still contain topic1 for address1")
	assert.ElementsMatch(t, updatedCache[address2], []common.Hash{topic2}, "Cache should contain topic2 for address2")

	// Reset logs before adding an extra subscription
	mockLogger.Reset()

	// Add an extra subscription for address1/topic1
	_, err = manager.Subscribe(address1, topic1)
	require.NoError(t, err)

	// Assert log for extra subscription
	expectedSubscriptionLog3 := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s SubscriberCount=2 New subscription added",
		address1.Hex(),
		topic1.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedSubscriptionLog3), "Expected log message for extra subscription")

	// Reset logs before Unsubscribe operation
	mockLogger.Reset()

	// Unsubscribe from address2/topic2
	err = manager.Unsubscribe(address2, topic2, ch)
	require.NoError(t, err)

	// Assert log for unsubscription
	expectedRemoveLog := fmt.Sprintf(
		"ChainID=1 Address=%s Topic=%s RemainingSubscribers=0 Subscription removed",
		address2.Hex(),
		topic2.Hex(),
	)
	assert.True(t, mockLogger.ContainsLog(expectedRemoveLog), "Expected log message for unsubscribing address2/topic2")

	// Assert cache invalidation log
	assert.True(t, mockLogger.ContainsLog(expectedCacheLog), "Expected cache invalidation log after unsubscription")

	// Check final cache
	finalCache := manager.GetAddressesAndTopics()
	require.Contains(t, finalCache, address1, "Address1 should still be in the cache")
	assert.ElementsMatch(t, finalCache[address1], []common.Hash{topic1}, "Cache should still contain topic1 for address1")
	require.NotContains(t, finalCache[address2], topic2, "Topic2 should be removed for address2 after unsubscription")
}

func TestSubscriptionManager_Close(t *testing.T) {
	manager := setupSubscriptionManager()

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Reset logs before Close operation
	mockLogger := manager.logger.(*internal.MockLogger)
	mockLogger.Reset()

	// Close the SubscriptionManager
	manager.Close()

	// Verify channel is closed
	_, open := <-ch
	assert.False(t, open, "Channel should be closed after Close()")

	// Assert close log message
	expectedCloseLog := "ChainID=1 SubscriptionManager closed, all subscriber channels have been closed"
	assert.True(t, mockLogger.ContainsLog(expectedCloseLog), "Expected close log message")

	// Verify registry is empty
	manager.registryMutex.RLock()
	defer manager.registryMutex.RUnlock()
	assert.Len(t, manager.registry, 0, "Registry should be empty after Close()")
}
