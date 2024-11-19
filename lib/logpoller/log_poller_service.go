package logpoller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// LogPollerService manages multiple ChainLogPollers across different chains
type LogPollerService struct {
	chainPollers map[int64]*ChainLogPoller    // Keyed by Chain ID
	cancelFuncs  map[int64]context.CancelFunc // Keyed by Chain ID
	logger       logging.Logger
	mutex        sync.RWMutex
}

// LogPollerServiceConfig holds configuration for LogPollerService
type LogPollerServiceConfig struct {
	SethClients         []*seth.Client          // List of SETH clients for different chains
	PollIntervals       map[int64]time.Duration // Custom poll intervals keyed by Chain ID
	Logger              *logging.Logger         // Pointer to Logger instance
	DefaultPollInterval time.Duration           // Fallback poll interval if not specified in PollIntervals
}

// NewLogPollerService initializes a new LogPollerService
func NewLogPollerService(cfg LogPollerServiceConfig) (*LogPollerService, error) {
	// Initialize the logger
	logging.Init()
	var logger logging.Logger
	if cfg.Logger == nil {
		defaultLogger := logging.GetLogger(nil, "CLIENT_LOG_POLLER_LOG_LEVEL")
		logger = logging.Logger(defaultLogger)
	} else {
		logger = logging.Logger(*cfg.Logger)
	}

	// Set DefaultPollInterval to 30 seconds if not provided
	if cfg.DefaultPollInterval <= 0 {
		cfg.DefaultPollInterval = 30 * time.Second
		logger.Info().
			Dur("DefaultPollInterval", cfg.DefaultPollInterval).
			Msg("DefaultPollInterval not set. Using default value of 30 seconds.")
	}

	service := &LogPollerService{
		chainPollers: make(map[int64]*ChainLogPoller),
		cancelFuncs:  make(map[int64]context.CancelFunc),
		logger:       logger,
	}

	// Validate that SethClients are provided
	if len(cfg.SethClients) == 0 {
		return nil, fmt.Errorf("no SETH clients provided in configuration")
	}

	// Initialize chain pollers
	for _, client := range cfg.SethClients {
		chainID := int64(client.ChainID)
		chainLogger := logger.With().Int64("ChainID", chainID).Logger()

		// Duplicate Chain ID Check
		if _, exists := service.chainPollers[chainID]; exists {
			service.logger.Warn().
				Int64("ChainID", chainID).
				Msg("Duplicate ChainID detected during initialization; skipping poller creation")
			continue
		}

		// Custom Poll Interval Handling
		pollInterval, exists := cfg.PollIntervals[chainID]
		if !exists {
			pollInterval = cfg.DefaultPollInterval // Fallback to default if not specified
		}

		// Create a new context for each poller
		ctx, cancel := context.WithCancel(context.Background())

		chainPoller, err := NewChainLogPoller(ChainLogPollerConfig{
			Seth:         client,
			PollInterval: pollInterval,
			Logger:       chainLogger,
		})
		if err != nil {
			cancel() // Cancel the context if poller creation fails
			return nil, fmt.Errorf("failed to create ChainLogPoller for chain %d: %w", chainID, err)
		}

		service.chainPollers[chainID] = chainPoller
		service.cancelFuncs[chainID] = cancel

		// Start polling
		chainPoller.Start(ctx)
	}

	return service, nil
}

// StopPoller stops the polling process for a specific chain by cancelling its context and removes it from the LogPoller Service.
func (s *LogPollerService) StopPoller(chainID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cancel, exists := s.cancelFuncs[chainID]
	if !exists {
		s.logger.Error().
			Int64("ChainID", chainID).
			Msg("No poller found to stop")
		return fmt.Errorf("no poller found for chain ID %d", chainID)
	}

	// Cancel the context to stop the poller
	cancel()

	// Wait for the poller to finish
	poller := s.chainPollers[chainID]
	poller.Wait()

	// Remove the poller and cancel function from the maps
	delete(s.chainPollers, chainID)
	delete(s.cancelFuncs, chainID)

	s.logger.Info().
		Int64("ChainID", chainID).
		Msg("Polling stopped for chain")

	return nil
}

// StopAllPollers stops and removes all active pollers
func (s *LogPollerService) StopAllPollers() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for chainID, cancel := range s.cancelFuncs {
		cancel()
		poller := s.chainPollers[chainID]
		poller.Wait()
		delete(s.chainPollers, chainID)
		delete(s.cancelFuncs, chainID)
		s.logger.Info().
			Int64("ChainID", chainID).
			Msg("Stopped and removed poller for chain")
	}
}

// AddPoller adds a new ChainLogPoller for the specified chain ID and starts polling.
func (s *LogPollerService) AddPoller(chainID int64, sethClient *seth.Client, pollInterval time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.chainPollers[chainID]; exists {
		s.logger.Info().
			Int64("ChainID", chainID).
			Msg("Chain Poller already exists.")
		return nil
	}

	chainLogger := s.logger.With().Int64("ChainID", chainID).Logger()

	// Create a new context for the poller
	ctx, cancel := context.WithCancel(context.Background())

	chainPoller, err := NewChainLogPoller(ChainLogPollerConfig{
		Seth:         sethClient,
		PollInterval: pollInterval,
		Logger:       chainLogger,
	})
	if err != nil {
		cancel() // Cancel the context if poller creation fails
		return fmt.Errorf("failed to create ChainLogPoller for chain %d: %w", chainID, err)
	}

	s.chainPollers[chainID] = chainPoller
	s.cancelFuncs[chainID] = cancel

	// Start polling
	chainPoller.Start(ctx)

	s.logger.Info().
		Int64("ChainID", chainID).
		Msg("Added and started polling for new chain")

	return nil
}

// Subscribe allows consumers to subscribe to events on a specific chain
func (s *LogPollerService) Subscribe(chainID int64, address common.Address, topic common.Hash) (chan LogEvent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	poller, exists := s.chainPollers[chainID]
	if !exists {
		return nil, fmt.Errorf("no poller found for chain ID %d", chainID)
	}

	ch, err := poller.SubscriptionManager.Subscribe(address, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to event: %w", err)
	}

	s.logger.Info().
		Int64("ChainID", chainID).
		Str("Address", address.Hex()).
		Str("Topic", topic.Hex()).
		Msg("Consumer subscribed to event")

	return ch, nil
}

// Unsubscribe allows consumers to unsubscribe from events on a specific chain
func (s *LogPollerService) Unsubscribe(chainID int64, address common.Address, topic common.Hash, ch chan LogEvent) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	poller, exists := s.chainPollers[chainID]
	if !exists {
		return fmt.Errorf("no poller found for chain ID %d", chainID)
	}

	err := poller.SubscriptionManager.Unsubscribe(address, topic, ch)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from event: %w", err)
	}

	s.logger.Info().
		Int64("ChainID", chainID).
		Str("Address", address.Hex()).
		Str("Topic", topic.Hex()).
		Msg("Consumer unsubscribed from event")

	return nil
}
