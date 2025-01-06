// File: internal/subscription_manager/subscription_manager.go
package subscription_manager

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"

	"github.com/ethereum/go-ethereum/common"
)

type SubscriptionManagerConfig struct {
	Logger  *zerolog.Logger
	ChainID int64
}

// SubscriptionManager manages subscriptions for a specific chain.
type SubscriptionManager struct {
	mu                sync.RWMutex // Protects all shared fields below
	registry          map[internal.EventKey][]chan api.Log
	logger            zerolog.Logger
	chainID           int64
	cachedEventKeys   []internal.EventKey
	cacheInitialized  bool
	channelBufferSize int
	closing           bool // Indicates if the manager is shutting down
	wg                sync.WaitGroup
}

// NewSubscriptionManager initializes a new SubscriptionManager.
func NewSubscriptionManager(cfg SubscriptionManagerConfig) *SubscriptionManager {
	subscriptionManagerLogger := cfg.Logger.With().Str("Component", "SubscriptionManager").Logger()

	return &SubscriptionManager{
		registry:          make(map[internal.EventKey][]chan api.Log),
		logger:            subscriptionManagerLogger,
		chainID:           cfg.ChainID,
		channelBufferSize: 3,
	}
}

// Subscribe registers a new subscription and returns a channel for events.
func (sm *SubscriptionManager) Subscribe(address common.Address, topic common.Hash) (chan api.Log, error) {
	if address == (common.Address{}) {
		sm.logger.Warn().Msg("Attempted to subscribe with an empty address")
		return nil, errors.New("address cannot be empty")
	}
	if topic == (common.Hash{}) {
		sm.logger.Warn().Msg("Attempted to subscribe with an empty topic")
		return nil, errors.New("topic cannot be empty")
	}

	sm.invalidateCache()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	eventKey := internal.EventKey{Address: address, Topic: topic}
	newChan := make(chan api.Log, sm.channelBufferSize)
	sm.registry[eventKey] = append(sm.registry[eventKey], newChan)

	sm.logger.Info().
		Int64("ChainID", sm.chainID).
		Str("Address", address.Hex()).
		Str("Topic", topic.Hex()).
		Int64("SubscriberCount", int64(len(sm.registry[eventKey]))).
		Msg("New subscription added")

	return newChan, nil
}

// Unsubscribe removes a subscription and closes the channel.
func (sm *SubscriptionManager) Unsubscribe(address common.Address, topic common.Hash, ch chan api.Log) error {
	eventKey := internal.EventKey{Address: address, Topic: topic}
	sm.mu.RLock()
	subscribers, exists := sm.registry[eventKey]
	sm.mu.RUnlock()
	if !exists {
		sm.logger.Warn().
			Int64("ChainID", sm.chainID).
			Str("Address", address.Hex()).
			Str("Topic", topic.Hex()).
			Msg("Attempted to unsubscribe from a non-existent EventKey")
		return errors.New("event key does not exist")
	}

	found := false // Flag to track if the subscriber was found

	for i, subscriber := range subscribers {
		if subscriber == ch {
			sm.invalidateCache()
			// Remove the subscriber from the list
			sm.mu.Lock()
			sm.registry[eventKey] = append(subscribers[:i], subscribers[i+1:]...)
			sm.mu.Unlock()
			sm.logger.Info().
				Int64("ChainID", sm.chainID).
				Str("Address", address.Hex()).
				Str("Topic", topic.Hex()).
				Int64("RemainingSubscribers", int64(len(sm.registry[eventKey]))).
				Msg("Subscription removed")
			found = true
			break
		}
	}

	if !found {
		// Subscriber channel was not found in the registry
		sm.logger.Warn().
			Int64("ChainID", sm.chainID).
			Hex("Address", []byte(address.Hex())).
			Hex("Topic", []byte(topic.Hex())).
			Msg("Attempted to unsubscribe a non-existent subscriber")
		return errors.New("subscriber channel not found")
	}

	sm.mu.Lock()
	if len(sm.registry[eventKey]) == 0 {
		// Clean up the map if there are no more subscribers
		delete(sm.registry, eventKey)
		sm.logger.Debug().
			Int64("ChainID", sm.chainID).
			Str("Address", address.Hex()).
			Str("Topic", topic.Hex()).
			Msg("No remaining subscribers, removing EventKey from registry")
	}
	sm.mu.Unlock()
	sm.wg.Wait() // Wait for all sends to complete before closing the channel

	close(ch) // Safely close the channel
	return nil
}

// BroadcastLog sends the log event to all relevant subscribers.
func (sm *SubscriptionManager) BroadcastLog(eventKey internal.EventKey, log api.Log) {
	// Check if the manager is closing
	sm.mu.RLock()
	if sm.closing {
		defer sm.mu.RUnlock()
		sm.logger.Debug().
			Interface("EventKey", eventKey).
			Msg("SubscriptionManager is closing, skipping broadcast")
		return
	}
	// Retrieve subscribers
	subscribers, exists := sm.registry[eventKey]
	sm.mu.RUnlock()

	if !exists {
		sm.logger.Debug().
			Interface("EventKey", eventKey).
			Msg("EventKey not found in registry")
		return
	}

	for _, ch := range subscribers {
		sm.wg.Add(1)
		go func(ch chan api.Log) {
			defer sm.wg.Done()
			select {
			case ch <- log:
			case <-time.After(100 * time.Millisecond): // Prevent blocking forever
				sm.logger.Warn().
					Int64("ChainID", sm.chainID).
					Msg("Log broadcast to channel timed out")
			}
		}(ch)
	}
	sm.logger.Debug().
		Int64("ChainID", sm.chainID).
		Int("Subscribers", len(subscribers)).
		Str("Address", eventKey.Address.Hex()).
		Str("Topic", eventKey.Topic.Hex()).
		Msg("Log broadcasted to all subscribers")
}

// GetAddressesAndTopics retrieves all unique EventKeys.
// Implements caching: caches the result after the first call and invalidates it upon subscription changes.
// Returns a slice of EventKeys, each containing a unique address-topic pair.
func (sm *SubscriptionManager) GetAddressesAndTopics() []internal.EventKey {
	sm.mu.RLock()
	if sm.cacheInitialized {
		defer sm.mu.RUnlock()
		return sm.cachedEventKeys
	}
	sm.mu.RUnlock()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	eventKeys := make([]internal.EventKey, 0, len(sm.registry))
	for eventKey := range sm.registry {
		eventKeys = append(eventKeys, eventKey)
	}

	// Update the cache
	sm.cachedEventKeys = eventKeys
	sm.cacheInitialized = true

	sm.logger.Debug().
		Int64("ChainID", sm.chainID).
		Int("UniqueEventKeys", len(sm.cachedEventKeys)).
		Msg("Cached EventKeys")

	return sm.cachedEventKeys
}

// invalidateCache invalidates the cached addresses and topics.
func (sm *SubscriptionManager) invalidateCache() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.cacheInitialized = false
	sm.cachedEventKeys = nil

	sm.logger.Debug().
		Int64("ChainID", sm.chainID).
		Msg("Cache invalidated due to subscription change")
}

// Close gracefully shuts down the SubscriptionManager by closing all subscriber channels.
func (sm *SubscriptionManager) Close() {
	// Set the closing flag under sendMutex
	sm.mu.Lock()
	sm.closing = true // Signal that the manager is closing
	sm.mu.Unlock()

	sm.wg.Wait() // Wait for all sends to complete before closing the channels

	sm.mu.Lock()
	for eventKey, subscribers := range sm.registry {
		for _, ch := range subscribers {
			close(ch)
		}
		delete(sm.registry, eventKey)
	}
	sm.mu.Unlock()
	sm.invalidateCache()

	sm.logger.Info().
		Int64("ChainID", sm.chainID).
		Msg("SubscriptionManager closed, all subscriber channels have been closed")
}
