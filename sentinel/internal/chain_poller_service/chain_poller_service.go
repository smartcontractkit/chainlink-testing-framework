// File: internal/chain_poller_service/chain_poller_service.go
package chain_poller_service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal/chain_poller"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal/subscription_manager"
)

// ChainPollerServiceConfig holds the configuration for the ChainPollerService.
type ChainPollerServiceConfig struct {
	PollInterval     time.Duration
	Logger           *zerolog.Logger
	BlockchainClient api.BlockchainClient
	ChainID          int64
}

// ChainPollerService orchestrates the polling process and log broadcasting.
type ChainPollerService struct {
	config          ChainPollerServiceConfig
	SubscriptionMgr *subscription_manager.SubscriptionManager
	ChainPoller     chain_poller.ChainPollerInterface
	ChainID         int64
	LastBlock       *big.Int
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	started         bool
	mu              sync.Mutex
}

func (eps *ChainPollerService) SubscriptionManager() *subscription_manager.SubscriptionManager {
	return eps.SubscriptionMgr
}

// NewChainPollerService initializes a new ChainPollerService.
func NewChainPollerService(cfg ChainPollerServiceConfig) (*ChainPollerService, error) {
	if cfg.PollInterval <= 0 {
		return nil, fmt.Errorf("poll interval must be positive")
	}
	if cfg.BlockchainClient == nil {
		return nil, fmt.Errorf("blockchain client cannot be nil")
	}
	if cfg.ChainID < 1 {
		return nil, fmt.Errorf("chainid missing")
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Create a subscrition manager
	subscription_manager := subscription_manager.NewSubscriptionManager(subscription_manager.SubscriptionManagerConfig{Logger: cfg.Logger, ChainID: cfg.ChainID})
	chain_poller, err := chain_poller.NewChainPoller(chain_poller.ChainPollerConfig{
		BlockchainClient: cfg.BlockchainClient,
		Logger:           cfg.Logger,
		ChainID:          cfg.ChainID})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize ChainPoller: %w", err)
	}

	l := cfg.Logger.With().Str("Component", "ChainPollerService").Logger().With().Int64("ChainID", cfg.ChainID).Logger()

	cfg.Logger = &l

	// Initialize lastBlock as the latest block at startup
	latestBlock, err := cfg.BlockchainClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	if latestBlock == 0 {
		return nil, errors.New("blockchain has no blocks")
	}
	lastBlock := new(big.Int).Sub(new(big.Int).SetUint64(latestBlock), big.NewInt(1))

	return &ChainPollerService{
		config:          cfg,
		SubscriptionMgr: subscription_manager,
		ChainPoller:     chain_poller,
		ChainID:         cfg.ChainID,
		LastBlock:       lastBlock,
	}, nil
}

// Start begins the polling loop.
func (eps *ChainPollerService) Start() {
	eps.mu.Lock()
	defer eps.mu.Unlock()

	if eps.started {
		eps.config.Logger.Warn().Msg("ChainPollerService already started")
		return
	}

	eps.ctx, eps.cancel = context.WithCancel(context.Background())
	eps.started = true
	eps.wg.Add(1)
	go eps.pollingLoop()
	eps.config.Logger.Info().Dur("Poll interval", eps.config.PollInterval).Msg("ChainPollerService started")
}

// Stop gracefully stops the polling loop.
func (eps *ChainPollerService) Stop() {
	eps.mu.Lock()
	defer eps.mu.Unlock()

	if !eps.started {
		return
	}
	eps.SubscriptionMgr.Close()
	eps.cancel()
	eps.wg.Wait()
	eps.started = false

	eps.config.Logger.Info().Msg("ChainPollerService stopped")
}

// pollingLoop runs the periodic polling process.
func (eps *ChainPollerService) pollingLoop() {
	defer eps.wg.Done()

	ticker := time.NewTicker(eps.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-eps.ctx.Done():
			eps.config.Logger.Info().Msg("Polling loop terminating")
			return
		case <-ticker.C:
			eps.pollCycle()
		}
	}
}

// pollCycle performs a single polling cycle: fetching logs and broadcasting them.
func (eps *ChainPollerService) pollCycle() {
	startTime := time.Now()
	eps.config.Logger.Debug().Msg("Starting polling cycle")

	// Fetch the latest block number
	latestBlock, err := eps.config.BlockchainClient.BlockNumber(eps.ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			eps.config.Logger.Debug().Msg("Fetching latest block number canceled due to shutdown")
		} else {
			eps.config.Logger.Error().Err(err).Msg("Failed to get latest block")
		}

		return
	}

	toBlock := latestBlock
	fromBlock := new(big.Int).Add(eps.LastBlock, big.NewInt(1))

	// Ensure fromBlock is not greater than toBlock
	if fromBlock.Cmp(new(big.Int).SetUint64(toBlock)) > 0 {
		eps.config.Logger.Warn().Msg(fmt.Sprintf("fromBlock (%s) is greater than toBlock (%s), skipping poll", fromBlock.String(), new(big.Int).SetUint64(toBlock).String()))
		return
	}

	// Get current subscriptions
	subscriptions := eps.SubscriptionMgr.GetAddressesAndTopics()

	if len(subscriptions) == 0 {
		// Update the last processed block to toBlock
		eps.LastBlock = new(big.Int).SetUint64(toBlock)
		eps.config.Logger.Debug().Msg("No active subscriptions, skipping polling cycle")
		return
	}

	// Construct filter queries with the same fromBlock and toBlock
	var filterQueries []api.FilterQuery
	for _, eventKey := range subscriptions {
		filterQueries = append(filterQueries, api.FilterQuery{
			FromBlock: fromBlock.Uint64(),
			ToBlock:   toBlock,
			Addresses: []common.Address{eventKey.Address},
			Topics:    [][]common.Hash{{eventKey.Topic}},
		})
	}

	// Fetch logs using the Chain Poller
	ctx, cancel := context.WithTimeout(eps.ctx, 10*time.Second)
	defer cancel()

	logs, err := eps.ChainPoller.FilterLogs(ctx, filterQueries)
	if err != nil {
		eps.config.Logger.Error().Err(err).Msg("Error during polling")
		return
	}

	eps.config.Logger.Debug().
		Int("Number of fetched logs", len(logs)).
		Uint64("FromBlock", fromBlock.Uint64()).
		Uint64("ToBlock", toBlock).
		Uint64("Number of blocks", toBlock-fromBlock.Uint64()).
		Msg(("Fetched logs from blockchain"))

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
			eps.SubscriptionMgr.BroadcastLog(eventKey, log)
		}
	}

	// Update the last processed block to toBlock
	eps.LastBlock = new(big.Int).SetUint64(toBlock)

	duration := time.Since(startTime)
	eps.config.Logger.Debug().Dur("Duration", duration).Msg("Completed polling cycle")
}
