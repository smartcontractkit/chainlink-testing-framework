package logpoller

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

// EventKey uniquely identifies an event subscription based on address and topic
type EventKey struct {
	Address common.Address
	Topic   common.Hash
}

// LogEvent represents a single log event fetched from the blockchain
type LogEvent struct {
	BlockNumber uint64      // The block number where the event was included
	TxHash      common.Hash // The transaction hash that emitted the event
	Data        []byte      // The data payload of the event
}

// SubscriptionManager manages subscriptions for a specific chain
type SubscriptionManager struct {
	registry          map[EventKey][]chan LogEvent
	registryMutex     sync.RWMutex
	logger            logging.Logger
	chainID           int64
	addressesCache    []common.Address
	topicsCache       [][]common.Hash
	cacheInitialized  bool
	cacheMutex        sync.RWMutex
	channelBufferSize int
}

// NewSubscriptionManager initializes a new SubscriptionManager
func NewSubscriptionManager(logger logging.Logger, chainID int64) *SubscriptionManager {
	return &SubscriptionManager{
		registry:          make(map[EventKey][]chan LogEvent),
		logger:            logger,
		chainID:           chainID,
		channelBufferSize: 3,
	}
}

// Subscribe registers a new subscription and returns a channel for events
func (sm *SubscriptionManager) Subscribe(address common.Address, topic common.Hash) (chan LogEvent, error) {
	if address == (common.Address{}) {
		sm.logger.Warn().
			Int64("ChainID", sm.chainID).
			Msg("Attempted to subscribe with an empty address")
		return nil, errors.New("address cannot be empty")
	}
	if topic == (common.Hash{}) {
		sm.logger.Warn().
			Int64("ChainID", sm.chainID).
			Msg("Attempted to subscribe with an empty topic")
		return nil, errors.New("topic cannot be empty")
	}

	sm.registryMutex.Lock()
	defer sm.registryMutex.Unlock()

	eventKey := EventKey{Address: address, Topic: topic}
	newChan := make(chan LogEvent, sm.channelBufferSize)
	sm.registry[eventKey] = append(sm.registry[eventKey], newChan)

	sm.invalidateCache()

	sm.logger.Info().
		Int64("ChainID", sm.chainID).
		Str("Address", address.Hex()).
		Str("Topic", topic.Hex()).
		Int("Subscriber Count", len(sm.registry[eventKey])).
		Msg("New subscription added")

	return newChan, nil
}

// Unsubscribe removes a subscription and closes the channel
func (sm *SubscriptionManager) Unsubscribe(address common.Address, topic common.Hash, ch chan LogEvent) error {
	sm.registryMutex.Lock()
	defer sm.registryMutex.Unlock()

	eventKey := EventKey{Address: address, Topic: topic}
	subscribers, exists := sm.registry[eventKey]
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
			// Remove the subscriber from the list
			sm.registry[eventKey] = append(subscribers[:i], subscribers[i+1:]...)
			close(ch)
			sm.logger.Info().
				Int64("ChainID", sm.chainID).
				Str("Address", address.Hex()).
				Str("Topic", topic.Hex()).
				Int("Remaining Subscribers", len(sm.registry[eventKey])).
				Msg("Subscription removed")
			found = true
			break
		}
	}

	if !found {
		// Subscriber channel was not found in the registry
		sm.logger.Warn().
			Int64("ChainID", sm.chainID).
			Str("Address", address.Hex()).
			Str("Topic", topic.Hex()).
			Msg("Attempted to unsubscribe a non-existent subscriber")
		return errors.New("subscriber channel not found")
	}

	if len(sm.registry[eventKey]) == 0 {
		// Clean up the map if there are no more subscribers
		delete(sm.registry, eventKey)
		sm.logger.Debug().
			Int64("ChainID", sm.chainID).
			Str("Address", address.Hex()).
			Str("Topic", topic.Hex()).
			Msg("No remaining subscribers, removing EventKey from registry")
	}

	sm.invalidateCache()

	return nil
}

// BroadcastLog sends the log event to all relevant subscribers
func (sm *SubscriptionManager) BroadcastLog(eventKey EventKey, log LogEvent) {
	sm.registryMutex.RLock()
	subscribers, exists := sm.registry[eventKey]
	sm.registryMutex.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subscribers {
		// Blocking send to ensure no logs are dropped
		select {
		case ch <- log:
			// Successfully sent
		}
	}

	sm.logger.Debug().
		Int64("ChainID", sm.chainID).
		Str("Address", eventKey.Address.Hex()).
		Str("Topic", eventKey.Topic.Hex()).
		Msg("Log broadcasted to all subscribers")
}

// GetAddressesAndTopics retrieves all unique addresses and their associated topics
// Implements caching: caches the result after the first call and invalidates it upon subscription changes
// Returns a slice of addresses and a slice of topic slices, where each topics[i] corresponds to addresses[i]
func (sm *SubscriptionManager) GetAddressesAndTopics() ([]common.Address, [][]common.Hash) {
	sm.cacheMutex.RLock()
	defer sm.cacheMutex.RUnlock()
	if sm.cacheInitialized {
		return sm.addressesCache, sm.topicsCache
	}

	sm.registryMutex.RLock()
	defer sm.registryMutex.RUnlock()

	addressSet := make(map[common.Address]struct{})
	topicSet := make(map[common.Hash]struct{})

	for eventKey := range sm.registry {
		addressSet[eventKey.Address] = struct{}{}
		topicSet[eventKey.Topic] = struct{}{}
	}

	addresses := make([]common.Address, 0, len(addressSet))
	for addr := range addressSet {
		addresses = append(addresses, addr)
	}

	topics := make([][]common.Hash, 1)
	topics[0] = make([]common.Hash, 0, len(topicSet))
	for topic := range topicSet {
		topics[0] = append(topics[0], topic)
	}

	// Update cache
	sm.addressesCache = addresses
	sm.topicsCache = topics
	sm.cacheInitialized = true

	sm.logger.Debug().
		Int64("ChainID", sm.chainID).
		Int("Unique Addresses", len(addresses)).
		Int("Unique Topics", len(topicSet)).
		Msg("Cached addresses and topics")

	return addresses, topics
}

// invalidateCache invalidates the cached addresses and topics
func (sm *SubscriptionManager) invalidateCache() {
	sm.cacheMutex.Lock()
	defer sm.cacheMutex.Unlock()
	sm.cacheInitialized = false
	sm.addressesCache = nil
	sm.topicsCache = nil
	sm.logger.Debug().
		Int64("ChainID", sm.chainID).
		Msg("Cache invalidated due to subscription change")
}

// Close gracefully shuts down the SubscriptionManager by closing all subscriber channels
func (sm *SubscriptionManager) Close() {
	sm.registryMutex.Lock()
	defer sm.registryMutex.Unlock()

	for eventKey, subscribers := range sm.registry {
		for _, ch := range subscribers {
			close(ch)
		}
		delete(sm.registry, eventKey)
	}

	sm.invalidateCache()

	sm.logger.Info().
		Int64("ChainID", sm.chainID).
		Msg("SubscriptionManager closed, all subscriber channels have been closed")
}
