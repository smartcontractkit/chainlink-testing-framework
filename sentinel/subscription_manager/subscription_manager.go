// File: subscription_manager/subscription_manager.go
package subscription_manager

import (
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"

	"github.com/ethereum/go-ethereum/common"
)

// SubscriptionManager manages subscriptions for a specific chain.
type SubscriptionManager struct {
	registry          map[internal.EventKey][]chan internal.Log
	registryMutex     sync.RWMutex
	logger            internal.Logger
	chainID           int64
	addressTopicCache map[common.Address][]common.Hash
	cacheInitialized  bool
	cacheMutex        sync.RWMutex
	channelBufferSize int

	closing     bool       // Indicates if the manager is shutting down
	activeSends int        // Tracks active sends in BroadcastLog
	cond        *sync.Cond // Used to coordinate between BroadcastLog and Close
}

// NewSubscriptionManager initializes a new SubscriptionManager.
func NewSubscriptionManager(logger internal.Logger, chainID int64) *SubscriptionManager {
	mu := &sync.Mutex{}

	return &SubscriptionManager{
		registry:          make(map[internal.EventKey][]chan internal.Log),
		logger:            logger,
		chainID:           chainID,
		channelBufferSize: 3,
		cond:              sync.NewCond(mu),
	}
}

// Subscribe registers a new subscription and returns a channel for events.
func (sm *SubscriptionManager) Subscribe(address common.Address, topic common.Hash) (chan internal.Log, error) {
	if address == (common.Address{}) {
		sm.logger.Warn("Attempted to subscribe with an empty address")
		return nil, errors.New("address cannot be empty")
	}
	if topic == (common.Hash{}) {
		sm.logger.Warn("Attempted to subscribe with an empty topic")
		return nil, errors.New("topic cannot be empty")
	}

	sm.registryMutex.Lock()
	defer sm.registryMutex.Unlock()

	eventKey := internal.EventKey{Address: address, Topic: topic}
	newChan := make(chan internal.Log, sm.channelBufferSize)
	sm.registry[eventKey] = append(sm.registry[eventKey], newChan)

	sm.invalidateCache()

	sm.logger.Infof("ChainID=%d Address=%s Topic=%s SubscriberCount=%d New subscription added",
		sm.chainID, address.Hex(), topic.Hex(), len(sm.registry[eventKey]))

	return newChan, nil
}

// Unsubscribe removes a subscription and closes the channel.
func (sm *SubscriptionManager) Unsubscribe(address common.Address, topic common.Hash, ch chan internal.Log) error {
	sm.registryMutex.Lock()
	defer sm.registryMutex.Unlock()

	eventKey := internal.EventKey{Address: address, Topic: topic}
	subscribers, exists := sm.registry[eventKey]
	if !exists {
		sm.logger.Warnf("ChainID=%d Address=%s Topic=%s Attempted to unsubscribe from a non-existent EventKey",
			sm.chainID, address.Hex(), topic.Hex())
		return errors.New("event key does not exist")
	}

	found := false // Flag to track if the subscriber was found

	for i, subscriber := range subscribers {
		if subscriber == ch {
			// Remove the subscriber from the list
			sm.registry[eventKey] = append(subscribers[:i], subscribers[i+1:]...)
			close(ch)
			sm.logger.Infof("ChainID=%d Address=%s Topic=%s RemainingSubscribers=%d Subscription removed",
				sm.chainID, address.Hex(), topic.Hex(), len(sm.registry[eventKey]))
			found = true
			break
		}
	}

	if !found {
		// Subscriber channel was not found in the registry
		sm.logger.Warnf("ChainID=%d Address=%s Topic=%s Attempted to unsubscribe a non-existent subscriber",
			sm.chainID, address.Hex(), topic.Hex())
		return errors.New("subscriber channel not found")
	}

	if len(sm.registry[eventKey]) == 0 {
		// Clean up the map if there are no more subscribers
		delete(sm.registry, eventKey)
		sm.logger.Debugf("ChainID=%d Address=%s Topic=%s No remaining subscribers, removing EventKey from registry",
			sm.chainID, address.Hex(), topic.Hex())
	}

	sm.invalidateCache()

	return nil
}

// BroadcastLog sends the log event to all relevant subscribers.
func (sm *SubscriptionManager) BroadcastLog(eventKey internal.EventKey, log internal.Log) {
	sm.registryMutex.RLock()
	subscribers, exists := sm.registry[eventKey]
	sm.registryMutex.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subscribers {
		sm.cond.L.Lock()
		if sm.closing {
			// If the manager is closing, skip sending logs
			sm.cond.L.Unlock()
			return
		}
		sm.activeSends++
		sm.cond.L.Unlock()
		go func(ch chan internal.Log) {
			defer func() {
				sm.cond.L.Lock()
				sm.activeSends--
				sm.cond.Broadcast() // Notify Close() when all sends are done
				sm.cond.L.Unlock()
			}()
			ch <- log
		}(ch)
	}

	sm.logger.Debugf("ChainID=%d Address=%s Topic=%s Log broadcasted to all subscribers",
		sm.chainID, eventKey.Address.Hex(), eventKey.Topic.Hex())
}

// GetAddressesAndTopics retrieves all unique addresses and their associated topics.
// Implements caching: caches the result after the first call and invalidates it upon subscription changes.
// Returns a map where each key is an address and the value is a slice of topics.
func (sm *SubscriptionManager) GetAddressesAndTopics() map[common.Address][]common.Hash {
	sm.cacheMutex.RLock()
	if sm.cacheInitialized {
		defer sm.cacheMutex.RUnlock()
		return sm.addressTopicCache
	}
	sm.cacheMutex.RUnlock()

	sm.registryMutex.RLock()
	defer sm.registryMutex.RUnlock()

	addressTopicMap := make(map[common.Address]map[common.Hash]struct{})

	for eventKey := range sm.registry {
		topicSet, exists := addressTopicMap[eventKey.Address]
		if !exists {
			topicSet = make(map[common.Hash]struct{})
			addressTopicMap[eventKey.Address] = topicSet
		}
		topicSet[eventKey.Topic] = struct{}{}
	}

	result := make(map[common.Address][]common.Hash)
	for addr, topics := range addressTopicMap {
		topicList := make([]common.Hash, 0, len(topics))
		for topic := range topics {
			topicList = append(topicList, topic)
		}
		result[addr] = topicList
	}

	// Update cache
	sm.cacheMutex.Lock()
	sm.addressTopicCache = result
	sm.cacheInitialized = true
	sm.cacheMutex.Unlock()

	sm.logger.Debugf("ChainID=%d UniqueAddresses=%d Cached address-topic pairs",
		sm.chainID, len(sm.addressTopicCache))

	return sm.addressTopicCache
}

// invalidateCache invalidates the cached addresses and topics.
func (sm *SubscriptionManager) invalidateCache() {
	sm.cacheMutex.Lock()
	sm.cacheInitialized = false
	sm.addressTopicCache = nil
	sm.cacheMutex.Unlock()

	sm.logger.Debugf("ChainID=%d Cache invalidated due to subscription change", sm.chainID)
}

// Close gracefully shuts down the SubscriptionManager by closing all subscriber channels.
func (sm *SubscriptionManager) Close() {
	sm.registryMutex.Lock()
	sm.closing = true // Signal that the manager is closing
	sm.registryMutex.Unlock()

	// Wait for all active sends to complete
	sm.cond.L.Lock()
	for sm.activeSends > 0 {
		sm.cond.Wait()
	}
	sm.cond.L.Unlock()

	sm.registryMutex.Lock()
	defer sm.registryMutex.Unlock()

	for eventKey, subscribers := range sm.registry {
		for _, ch := range subscribers {
			close(ch)
		}
		delete(sm.registry, eventKey)
	}

	sm.invalidateCache()

	sm.logger.Infof("ChainID=%d SubscriptionManager closed, all subscriber channels have been closed", sm.chainID)
}
