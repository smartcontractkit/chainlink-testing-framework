// File: sentinel.go
package sentinel

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal/chain_poller_service"
)

// SentinelConfig holds configuration for the Sentinel.
type SentinelConfig struct {
	Logger zerolog.Logger
}

type AddChainConfig struct {
	ChainID          int64
	PollInterval     time.Duration
	BlockchainClient api.BlockchainClient
}

type Sentinel struct {
	config   SentinelConfig
	mu       sync.RWMutex
	services map[int64]*chain_poller_service.ChainPollerService // Map of chainID to ChianPollerService
}

// NewSentinel initializes and returns a new Sentinel instance.
func NewSentinel(cfg SentinelConfig) *Sentinel {
	cfg.Logger = cfg.Logger.With().Str("component", "Sentinel").Logger()
	cfg.Logger.Info().Msg("Initializing Sentinel")
	return &Sentinel{
		config:   cfg,
		services: make(map[int64]*chain_poller_service.ChainPollerService),
	}
}

// AddChain adds a new chain to Sentinel.
func (s *Sentinel) AddChain(acc AddChainConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[acc.ChainID]; exists {
		return fmt.Errorf("chain with ID %d already exists", acc.ChainID)
	}

	cfg := chain_poller_service.ChainPollerServiceConfig{
		PollInterval:     acc.PollInterval,
		ChainID:          acc.ChainID,
		Logger:           &s.config.Logger,
		BlockchainClient: acc.BlockchainClient,
	}

	eps, err := chain_poller_service.NewChainPollerService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize ChainPollerService: %w", err)
	}
	s.services[cfg.ChainID] = eps
	s.config.Logger.Info().Int64("ChainID", cfg.ChainID).Msg("Added new chain")
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
	s.config.Logger.Info().Msg("Removed chain")
	return nil
}

// Subscribe subscribes to events for a specific chain.
func (s *Sentinel) Subscribe(chainID int64, address common.Address, topic common.Hash) (chan api.Log, error) {
	s.mu.RLock()
	eps, exists := s.services[chainID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("chain with ID %d does not exist", chainID)
	}

	return eps.SubscriptionManager().Subscribe(address, topic)
}

// Unsubscribe unsubscribes from events for a specific chain.
func (s *Sentinel) Unsubscribe(chainID int64, address common.Address, topic common.Hash, ch chan api.Log) error {
	s.mu.RLock()
	eps, exists := s.services[chainID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("chain with ID %d does not exist", chainID)
	}

	return eps.SubscriptionManager().Unsubscribe(address, topic, ch)
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
