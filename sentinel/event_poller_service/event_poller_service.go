// File: event_poller_service/event_poller_service.go
package event_poller_service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/chain_poller"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/subscription_manager"
)

// EventPollerServiceConfig holds the configuration for the EventPollerService.
type EventPollerServiceConfig struct {
	PollInterval     time.Duration
	ChainPoller      chain_poller.ChainPollerInterface
	SubscriptionMgr  *subscription_manager.SubscriptionManager
	Logger           internal.Logger
	ChainID          int64
	BlockchainClient internal.BlockchainClient
}

// EventPollerService orchestrates the polling process and log broadcasting.
type EventPollerService struct {
	config    EventPollerServiceConfig
	LastBlock *big.Int
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	started   bool
	mu        sync.Mutex
}

// NewEventPollerService initializes a new EventPollerService.
func NewEventPollerService(cfg EventPollerServiceConfig) (*EventPollerService, error) {
	if cfg.ChainPoller == nil {
		return nil, fmt.Errorf("chain poller cannot be nil")
	}
	if cfg.SubscriptionMgr == nil {
		return nil, fmt.Errorf("subscription manager cannot be nil")
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	if cfg.PollInterval <= 0 {
		return nil, fmt.Errorf("poll interval must be positive")
	}
	if cfg.BlockchainClient == nil {
		return nil, fmt.Errorf("blockchain client cannot be nil")
	}

	// Initialize lastBlock as the latest block at startup
	latestBlock, err := cfg.BlockchainClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	if latestBlock == 0 {
		return nil, errors.New("blockchain has no blocks")
	}
	lastBlock := new(big.Int).Sub(new(big.Int).SetUint64(latestBlock), big.NewInt(1))

	return &EventPollerService{
		config:    cfg,
		LastBlock: lastBlock,
	}, nil
}

// Start begins the polling loop.
func (eps *EventPollerService) Start() {
	eps.mu.Lock()
	defer eps.mu.Unlock()

	if eps.started {
		eps.config.Logger.Warn("EventPollerService already started")
		return
	}

	eps.ctx, eps.cancel = context.WithCancel(context.Background())
	eps.started = true
	eps.wg.Add(1)
	go eps.pollingLoop()

	eps.config.Logger.Info(fmt.Sprintf("EventPollerService started with poll interval: %s", eps.config.PollInterval.String()))
}

// Stop gracefully stops the polling loop.
func (eps *EventPollerService) Stop() {
	eps.mu.Lock()
	defer eps.mu.Unlock()

	if !eps.started {
		return
	}

	eps.cancel()
	eps.wg.Wait()
	eps.started = false

	eps.config.Logger.Info("EventPollerService stopped")
}

// pollingLoop runs the periodic polling process.
func (eps *EventPollerService) pollingLoop() {
	defer eps.wg.Done()

	ticker := time.NewTicker(eps.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-eps.ctx.Done():
			eps.config.Logger.Info("Polling loop terminating")
			return
		case <-ticker.C:
			eps.pollCycle()
		}
	}
}

// pollCycle performs a single polling cycle: fetching logs and broadcasting them.
func (eps *EventPollerService) pollCycle() {
	startTime := time.Now()
	eps.config.Logger.Debug("Starting polling cycle")

	// Fetch the latest block number
	latestBlock, err := eps.config.BlockchainClient.BlockNumber(eps.ctx)
	if err != nil {
		eps.config.Logger.Error("Failed to get latest block", "error", err)
		return
	}

	toBlock := latestBlock
	fromBlock := new(big.Int).Add(eps.LastBlock, big.NewInt(1))

	// Ensure fromBlock is not greater than toBlock
	if fromBlock.Cmp(new(big.Int).SetUint64(toBlock)) > 0 {
		eps.config.Logger.Warn(fmt.Sprintf("fromBlock (%s) is greater than toBlock (%s), skipping poll", fromBlock.String(), new(big.Int).SetUint64(toBlock).String()))
		return
	}

	// Get current subscriptions
	subscriptions := eps.config.SubscriptionMgr.GetAddressesAndTopics()

	if len(subscriptions) == 0 {
		// Update the last processed block to toBlock
		eps.LastBlock = new(big.Int).SetUint64(toBlock)
		eps.config.Logger.Debug("No active subscriptions, skipping polling cycle")
		return
	}

	// Construct filter queries with the same fromBlock and toBlock
	var filterQueries []internal.FilterQuery
	for address, topics := range subscriptions { // 'topics' is []common.Hash
		for _, topic := range topics { // Iterate over each topic
			filterQueries = append(filterQueries, internal.FilterQuery{
				FromBlock: fromBlock.Uint64(),
				ToBlock:   toBlock,
				Addresses: []common.Address{address},
				Topics:    [][]common.Hash{{topic}}, // Separate query per topic
			})
		}
	}

	// Fetch logs using the Chain Poller
	ctx, cancel := context.WithTimeout(eps.ctx, 10*time.Second)
	defer cancel()

	logs, err := eps.config.ChainPoller.Poll(ctx, filterQueries)
	if err != nil {
		eps.config.Logger.Error("Error during polling", "error", err)
		return
	}

	eps.config.Logger.Debug(fmt.Sprintf("Fetched %d logs from blockchain", len(logs)))

	// Broadcast each log to subscribers
	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue // Skip logs without topics
		}

		for _, topic := range log.Topics {
			eventKey := internal.EventKey{
				Address: log.Address,
				Topic:   topic,
			}
			eps.config.SubscriptionMgr.BroadcastLog(eventKey, log)
		}
	}

	// Update the last processed block to toBlock
	eps.LastBlock = new(big.Int).SetUint64(toBlock)

	duration := time.Since(startTime)
	eps.config.Logger.Debug(fmt.Sprintf("Completed polling cycle in %s", duration.String()))
}
