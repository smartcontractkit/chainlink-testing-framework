package logpoller

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// ChainLogPollerConfig holds configuration for ChainLogPoller
type ChainLogPollerConfig struct {
	Seth         *seth.Client
	PollInterval time.Duration
	Logger       logging.Logger
}

// ChainLogPoller polls logs from a specific blockchain and manages subscriptions
type ChainLogPoller struct {
	seth                *seth.Client
	logger              logging.Logger
	pollInterval        time.Duration
	lastBlock           *big.Int
	SubscriptionManager *SubscriptionManager
	started             bool
	startMutex          sync.Mutex
	wg                  sync.WaitGroup
}

// NewChainLogPoller initializes a new ChainLogPoller
func NewChainLogPoller(cfg ChainLogPollerConfig) (*ChainLogPoller, error) {
	// Initialize the lastBlock to the latest block
	latestBlock, err := cfg.Seth.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	latestBlockPlus := new(big.Int).SetUint64(latestBlock)
	// Subtract 1 to start polling from the latest block at initialization
	lastBlock := new(big.Int).Sub(latestBlockPlus, big.NewInt(1))

	// Initialize SubscriptionManager
	subMgr := NewSubscriptionManager(cfg.Logger, int64(cfg.Seth.ChainID))

	return &ChainLogPoller{
		seth:                cfg.Seth,
		logger:              cfg.Logger,
		pollInterval:        cfg.PollInterval,
		lastBlock:           lastBlock,
		SubscriptionManager: subMgr,
	}, nil
}

// Start begins the polling process only if not already started
func (clp *ChainLogPoller) Start(ctx context.Context) {
	clp.startMutex.Lock()
	defer clp.startMutex.Unlock()

	if clp.started {
		clp.logger.Warn().Msg("Poller already started")
	}

	clp.started = true
	clp.wg.Add(1)
	go clp.pollWithSubscriptions(ctx)

	clp.logger.Info().
		Str("PollInterval", clp.pollInterval.String()).
		Msg("Poller started.")
}

// pollWithSubscriptions handles the periodic polling
func (clp *ChainLogPoller) pollWithSubscriptions(ctx context.Context) {
	defer clp.wg.Done()
	ticker := time.NewTicker(clp.pollInterval)
	defer ticker.Stop()

	clp.logger.Info().Msg("Polling loop started.")

	for {
		select {
		case <-ctx.Done():
			clp.logger.Info().Msg("Shutting down ChainLogPoller")
			return
		case <-ticker.C:
			clp.logger.Debug().Msg("Starting polling cycle.")
			clp.pollLogs(ctx)
		}
	}
}

// Wait waits for the polling goroutine to finish
func (clp *ChainLogPoller) Wait() {
	clp.wg.Wait()
}

// pollLogs fetches logs from the blockchain and broadcasts them to subscribers
func (clp *ChainLogPoller) pollLogs(ctx context.Context) {
	startTime := time.Now()
	latestBlock, err := clp.seth.Client.BlockNumber(ctx)
	if err != nil {
		clp.logger.Error().Err(err).Msg("Failed to get latest block")
		return
	}
	fromBlock := new(big.Int).Add(clp.lastBlock, big.NewInt(1))
	toBlock := new(big.Int).SetUint64(latestBlock)

	addresses, topics := clp.SubscriptionManager.GetAddressesAndTopics()
	if len(addresses) == 0 || len(topics) == 0 {
		// No active subscriptions, skip polling
		clp.logger.Debug().
			Int("Addresses", len(addresses)).
			Int("Topics", len(topics)).
			Msg("No active subscriptions, skipping poll")
		return
	}

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
		Topics:    topics,
	}

	logs, err := clp.seth.Client.FilterLogs(ctx, query)
	if err != nil {
		clp.logger.Error().Err(err).Msg("Failed to filter logs")
		return
	}
	clp.logger.Debug().
		Int64("ChainID", clp.SubscriptionManager.chainID).
		Str("FromBlock", fromBlock.String()).
		Str("ToBlock", toBlock.String()).
		Int("LogsFetched", len(logs)).
		Msg("Fetched logs from blockchain")

	if len(logs) == 0 {
		clp.logger.Info().
			Int64("ChainID", clp.SubscriptionManager.chainID).
			Str("FromBlock", fromBlock.String()).
			Str("ToBlock", toBlock.String()).
			Msg("No new logs found in the current polling cycle")
		return
	}

	for _, vLog := range logs {
		if len(vLog.Topics) == 0 {
			continue // Skip logs without topics
		}
		// Iterate over all topics in the log
		for _, topic := range vLog.Topics {
			eventKey := EventKey{
				Address: vLog.Address,
				Topic:   topic,
			}

			logEvent := LogEvent{
				BlockNumber: vLog.BlockNumber,
				TxHash:      vLog.TxHash,
				Data:        vLog.Data,
			}

			clp.SubscriptionManager.BroadcastLog(eventKey, logEvent)
		}
	}
	clp.lastBlock = toBlock

	duration := time.Since(startTime)

	clp.logger.Debug().
		Int64("ChainID", clp.SubscriptionManager.chainID).
		Str("FromBlock", fromBlock.String()).
		Str("ToBlock", toBlock.String()).
		Str("LastBlock", clp.lastBlock.String()).
		Dur("PollingDuration", duration).
		Msg("Completed polling cycle")

}
