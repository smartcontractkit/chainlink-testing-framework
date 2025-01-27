package blockchain

import (
	"context"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

type ChainHeaderManager struct {
	chainID      int64
	pollInterval time.Duration
	networkCfg   EVMNetwork
	logger       zerolog.Logger

	ethClient *ethclient.Client
	rpcClient *rpc.Client

	done chan struct{}
	wg   sync.WaitGroup

	headersChan chan *SafeEVMHeader

	mu          sync.RWMutex
	subscribers map[*EthereumClient]struct{}

	lastProcessed uint64

	started bool
}

var (
	chainManagerRegistry = struct {
		sync.Mutex
		managers map[int64]*ChainHeaderManager
	}{
		managers: make(map[int64]*ChainHeaderManager),
	}
)

// getOrCreateChainManager returns an existing manager if found, otherwise creates one.
func getOrCreateChainManager(
	chainID int64,
	pollInterval time.Duration,
	networkCfg EVMNetwork,
	logger zerolog.Logger,
	ethClient *ethclient.Client,
	rpcClient *rpc.Client,
) *ChainHeaderManager {
	chainManagerRegistry.Lock()
	defer chainManagerRegistry.Unlock()

	if mgr, exists := chainManagerRegistry.managers[chainID]; exists {
		return mgr
	}

	mgr := newChainHeaderManager(chainID, pollInterval, networkCfg, logger, ethClient, rpcClient)
	chainManagerRegistry.managers[chainID] = mgr
	return mgr
}

func removeChainManager(chainID int64) {
	chainManagerRegistry.Lock()
	defer chainManagerRegistry.Unlock()
	delete(chainManagerRegistry.managers, chainID)
}

// newChainHeaderManager creates the manager but does not start polling automatically
func newChainHeaderManager(
	chainID int64,
	pollInterval time.Duration,
	networkCfg EVMNetwork,
	logger zerolog.Logger,
	ethClient *ethclient.Client,
	rpcClient *rpc.Client,
) *ChainHeaderManager {
	return &ChainHeaderManager{
		chainID:      chainID,
		pollInterval: pollInterval,
		networkCfg:   networkCfg,
		logger:       logger,
		ethClient:    ethClient,
		rpcClient:    rpcClient,
		subscribers:  make(map[*EthereumClient]struct{}),
		headersChan:  make(chan *SafeEVMHeader, 1000), // Buffer to handle rapid blocks
		done:         make(chan struct{}),
	}
}

// startPolling initiates the two background goroutines (poll + fan-out).
func (m *ChainHeaderManager) startPolling() {
	if m.started {
		return
	}
	m.started = true

	// Attempt an initial fetch of the latest block, so we know where to begin
	initCtx, cancel := context.WithTimeout(context.Background(), m.networkCfg.Timeout.Duration)
	defer cancel()
	latestHeader, err := m.ethClient.HeaderByNumber(initCtx, nil)
	if err != nil {
		m.logger.Error().
			Int64("ChainID", m.chainID).
			Err(err).
			Msg("Failed initial fetch of the latest header, manager won't start polling")
		return
	}
	safeLatest := convertToSafeEVMHeader(latestHeader)
	m.lastProcessed = safeLatest.Number.Uint64() - 1

	m.logger.Info().
		Int64("ChainID", m.chainID).
		Uint64("InitialBlock", m.lastProcessed).
		Msg("ChainHeaderManager starting polling")

	m.wg.Add(2)
	go m.pollRoutine()
	go m.fanOutRoutine()
}

// pollRoutine fetches new headers at a fixed interval and sends them down m.headersChan
func (m *ChainHeaderManager) pollRoutine() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.done:
			m.logger.Debug().
				Int64("ChainID", m.chainID).
				Msg("pollRoutine: shutting down")
			return
		case <-ticker.C:
			if err := m.fetchAndQueueNewHeaders(); err != nil {
				m.logger.Error().
					Int64("ChainID", m.chainID).
					Err(err).
					Msg("pollRoutine: error fetching new headers")
			}
		}
	}
}

// fanOutRoutine receives newly fetched headers from m.headersChan and distributes them
func (m *ChainHeaderManager) fanOutRoutine() {
	defer m.wg.Done()

	for {
		select {
		case <-m.done:
			m.logger.Debug().
				Int64("ChainID", m.chainID).
				Msg("fanOutRoutine: shutting down")
			return
		case hdr := <-m.headersChan:
			m.mu.RLock()
			for sub := range m.subscribers {
				err := sub.receiveHeader(hdr)
				if err != nil {
					m.logger.Err(err).Msg("Finalizer received error during HTTP polling")
				}
			}
			m.mu.RUnlock()
		}
	}
}

// fetchAndQueueNewHeaders fetches the latest header and then loops over any missing blocks
func (m *ChainHeaderManager) fetchAndQueueNewHeaders() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.networkCfg.Timeout.Duration)
	defer cancel()

	latest, err := m.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	latestNum := latest.Number.Uint64()

	// We already processed up to X, we process X+1..latest
	for blockNum := m.lastProcessed + 1; blockNum <= latestNum; blockNum++ {
		if blockNum > math.MaxInt64 {
			m.logger.Error().Int64("ChainID", m.chainID).
				Uint64("BlockNumber", blockNum).
				Msg("blockNum exceeds int64 max, skipping")
			continue
		}
		blockCtx, blockCancel := context.WithTimeout(context.Background(), m.networkCfg.Timeout.Duration)
		blockHdr, err := m.ethClient.HeaderByNumber(blockCtx, big.NewInt(int64(blockNum)))
		blockCancel()
		if err != nil {
			m.logger.Error().
				Int64("ChainID", m.chainID).
				Err(err).
				Uint64("BlockNumber", blockNum).
				Msg("Could not fetch block header in range")
			continue
		}
		safeHdr := convertToSafeEVMHeader(blockHdr)
		m.headersChan <- safeHdr
		m.lastProcessed = blockNum
	}
	return nil
}

// subscribe attaches an EthereumClient to our manager
func (m *ChainHeaderManager) subscribe(client *EthereumClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers[client] = struct{}{}
}

// unsubscribe removes a subscriber from the manager
func (m *ChainHeaderManager) unsubscribe(client *EthereumClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.subscribers, client)
}

// shutdown stops the goroutines and closes channels.
func (m *ChainHeaderManager) shutdown() {
	close(m.done)
	m.wg.Wait()
	close(m.headersChan)
}

func convertToSafeEVMHeader(hdr *types.Header) *SafeEVMHeader {
	if hdr == nil {
		return nil
	}
	var safeTime int64
	if hdr.Time > math.MaxInt64 {
		safeTime = math.MaxInt64
	} else {
		safeTime = int64(hdr.Time)
	}
	return &SafeEVMHeader{
		Hash:      hdr.Hash(),
		Number:    hdr.Number,
		BaseFee:   hdr.BaseFee,
		Timestamp: time.Unix(safeTime, 0),
	}
}
