// File: sentinel.go
package sentinel

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/chain_poller_service"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/sentinel/internal"
)

// SentinelConfig holds configuration for the Sentinel.
type SentinelConfig struct {
	Logger internal.Logger
}

type Sentinel struct {
	config   SentinelConfig
	mu       sync.RWMutex
	services map[int64]*chain_poller_service.ChainPollerService // Map of chainID to ChianPollerService
}

// NewSentinel initializes and returns a new Sentinel instance.
func NewSentinel(cfg SentinelConfig) *Sentinel {
	return &Sentinel{
		config:   cfg,
		services: make(map[int64]*chain_poller_service.ChainPollerService),
	}
}

// AddChain adds a new chain to Sentinel.
func (s *Sentinel) AddChain(cfg chain_poller_service.ChainPollerServiceConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[cfg.ChainID]; exists {
		return fmt.Errorf("chain with ID %d already exists", cfg.ChainID)
	}

	eps, err := chain_poller_service.NewChainPollerService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize ChainPollerService: %w", err)
	}

	s.services[cfg.ChainID] = eps
	eps.Start()
	return nil
}

// RemoveChain removes a chain from Sentinel.
func (s *Sentinel) RemoveChain(chainID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	eps, exists := s.services[chainID]
	if !exists {
		return fmt.Errorf("chain with ID %d does not exist", chainID)
	}

	eps.Stop()
	delete(s.services, chainID)
	s.config.Logger.Info(fmt.Sprintf("Removed chain with ID %d", chainID))
	return nil
}

// Subscribe subscribes to events for a specific chain.
func (s *Sentinel) Subscribe(chainID int64, address common.Address, topic common.Hash) (chan internal.Log, error) {
	s.mu.RLock()
	eps, exists := s.services[chainID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("chain with ID %d does not exist", chainID)
	}

	return eps.SubscriptionManager().Subscribe(address, topic)
}

// Unsubscribe unsubscribes from events for a specific chain.
func (s *Sentinel) Unsubscribe(chainID int64, address common.Address, topic common.Hash, ch chan internal.Log) error {
	s.mu.RLock()
	eps, exists := s.services[chainID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("chain with ID %d does not exist", chainID)
	}

	return eps.SubscriptionManager().Unsubscribe(address, topic, ch)
}

// GetService gets the chain poller service for a chain id.
func (s *Sentinel) GetService(chainID int64) (*chain_poller_service.ChainPollerService, bool) {
	s.mu.RLock()
	eps, exists := s.services[chainID]
	s.mu.RUnlock()

	return eps, exists
}

// HasServices returns true if there is at least 1 service running.
func (s *Sentinel) HasServices() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.services) > 0
}

// Close shuts down all chains and the global registry.
func (s *Sentinel) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, eps := range s.services {
		eps.Stop()
		delete(s.services, eps.ChainID)
	}
}
