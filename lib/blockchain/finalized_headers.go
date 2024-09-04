package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	minBlocksToFinalize = 10
	finalityTimeout     = 90 * time.Minute
	FinalizedHeaderKey  = "finalizedHeads"
	logNotifyFrequency  = 2 * time.Minute
)

var (
	globalFinalizedHeaderManager = sync.Map{}
)

// FinalizedHeader is an implementation of the HeaderEventSubscription interface.
// It keeps track of the latest finalized header for a network.
type FinalizedHeader struct {
	lggr              zerolog.Logger
	LatestFinalized   atomic.Value // *big.Int
	FinalizedAt       atomic.Value // time.Time
	client            EVMClient
	headerUpdateMutex *sync.Mutex
}

// Wait is not a blocking call.
func (f *FinalizedHeader) Wait() error {
	return nil
}

// Complete returns false as the HeaderEventSubscription should run for the entire course of the test.
func (f *FinalizedHeader) Complete() bool {
	return false
}

// ReceiveHeader is called whenever a new header is received.
// During the course of test whenever a new header is received, ReceiveHeader checks if there is a new finalized header tagged.
func (f *FinalizedHeader) ReceiveHeader(header NodeHeader) error {
	fLatest := f.LatestFinalized.Load().(*big.Int)
	// assumption : it will take at least minBlocksToFinalize num of blocks to finalize
	// this is to reduce the number of calls to the client
	if new(big.Int).Sub(header.Number, fLatest).Cmp(big.NewInt(minBlocksToFinalize)) <= 0 {
		return nil
	}
	f.headerUpdateMutex.Lock()
	defer f.headerUpdateMutex.Unlock()
	if f.FinalizedAt.Load() != nil {
		fTime := f.FinalizedAt.Load().(time.Time)
		// if the time difference between the new header and the last finalized header is less than 10s, ignore
		if header.Timestamp.Sub(fTime) <= 10*time.Second {
			return nil
		}
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), f.client.GetNetworkConfig().Timeout.Duration)
	lastFinalized, err := f.client.GetLatestFinalizedBlockHeader(ctx)
	ctxCancel()
	if err != nil {
		return fmt.Errorf("error getting latest finalized block header: %w", err)
	}

	if lastFinalized.Number.Cmp(fLatest) > 0 {
		f.LatestFinalized.Store(lastFinalized.Number)
		f.FinalizedAt.Store(header.Timestamp)
		f.lggr.Info().
			Str("Finalized Header", lastFinalized.Number.String()).
			Str("Finalized At", header.Timestamp.String()).
			Msg("new finalized header received")
	}
	return nil
}

// newGlobalFinalizedHeaderManager is a global manager for finalized headers per network.
// It is used to keep track of the latest finalized header for each network.
func newGlobalFinalizedHeaderManager(evmClient EVMClient) *FinalizedHeader {
	// if finality depth is greater than 0, there is no need to track finalized headers return nil
	if evmClient.GetNetworkConfig().FinalityDepth > 0 {
		return nil
	}
	f, ok := globalFinalizedHeaderManager.Load(evmClient.GetChainID().String())
	isFinalizedHeaderObsolete := false
	var fHeader *FinalizedHeader
	if f != nil {
		now := time.Now().UTC()
		// if the last finalized header is older than an hour
		lastFinalizedAt := f.(*FinalizedHeader).FinalizedAt.Load().(time.Time)
		isFinalizedHeaderObsolete = now.Sub(lastFinalizedAt) > 1*time.Hour
		fHeader = f.(*FinalizedHeader)
	}

	// if there is no finalized header for this network or the last finalized header is older than 1 hour
	if !ok || isFinalizedHeaderObsolete {
		mu := &sync.Mutex{}
		mu.Lock()
		defer mu.Unlock()
		ctx, ctxCancel := context.WithTimeout(context.Background(), evmClient.GetNetworkConfig().Timeout.Duration)
		lastFinalized, err := evmClient.GetLatestFinalizedBlockHeader(ctx)
		ctxCancel()
		if err != nil {
			log.Err(fmt.Errorf("error getting latest finalized block header %w", err)).Msg("NewFinalizedHeader")
			return nil
		}
		fHeader = &FinalizedHeader{
			lggr:              log.With().Str("Network", evmClient.GetNetworkName()).Logger(),
			client:            evmClient,
			headerUpdateMutex: mu,
		}
		fHeader.LatestFinalized.Store(lastFinalized.Number)
		fHeader.FinalizedAt.Store(time.Now().UTC())
		globalFinalizedHeaderManager.Store(evmClient.GetChainID().String(), fHeader)
		fHeader.lggr.Info().
			Str("Finalized Header", lastFinalized.Number.String()).
			Str("Finalized At", time.Now().UTC().String()).
			Msg("new finalized header received")
	}

	return fHeader
}

// TransactionFinalizer is an implementation of HeaderEventSubscription that waits for a transaction to be finalized.
type TransactionFinalizer struct {
	lggr          zerolog.Logger
	client        EVMClient
	doneChan      chan struct{}
	context       context.Context
	cancel        context.CancelFunc
	networkConfig *EVMNetwork
	complete      bool
	completeMu    sync.Mutex
	txHdr         *SafeEVMHeader
	txHash        common.Hash
	FinalizedBy   *big.Int
	FinalizedAt   time.Time
	lastLogUpdate time.Time
}

func NewTransactionFinalizer(client EVMClient, txHdr *SafeEVMHeader, txHash common.Hash) *TransactionFinalizer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), finalityTimeout)
	tf := &TransactionFinalizer{
		lggr: log.With().
			Str("txHash", txHash.String()).
			Str("Tx Block", txHdr.Number.String()).
			Str("Network", client.GetNetworkName()).
			Logger(),
		client:        client,
		doneChan:      make(chan struct{}, 1),
		context:       ctx,
		cancel:        ctxCancel,
		networkConfig: client.GetNetworkConfig(),
		complete:      false,
		txHdr:         txHdr,
		txHash:        txHash,
	}

	return tf
}

func (tf *TransactionFinalizer) ReceiveHeader(header NodeHeader) error {
	isFinalized, by, at, err := tf.client.IsTxHeadFinalized(tf.txHdr, &header.SafeEVMHeader)
	if err != nil {
		return err
	}
	if isFinalized {
		tf.lggr.Info().
			Str("Finalized Block", header.Number.String()).
			Str("Tx block", tf.txHdr.Number.String()).
			Msg("Found finalized log")
		tf.complete = true
		tf.FinalizedBy = by
		tf.FinalizedAt = at
		tf.doneChan <- struct{}{}
	} else {
		lgEvent := tf.lggr.Info()
		// if the transaction is not finalized, notify every logNotifyFrequency duration
		if time.Now().UTC().Sub(tf.lastLogUpdate) < logNotifyFrequency {
			return nil
		}
		if tf.networkConfig.FinalityDepth > 0 {
			lgEvent.
				Str("Current Block", header.Number.String()).
				Uint64("Finality Depth", tf.networkConfig.FinalityDepth)
		} else {
			lgEvent.
				Str("Last Finalized Block", by.String())
		}
		lgEvent.Msg("Still Waiting for transaction log to be finalized")
		tf.lastLogUpdate = time.Now().UTC()
	}
	return nil
}

// Wait is a blocking function that waits until the transaction is finalized or the context is cancelled
func (tf *TransactionFinalizer) Wait() error {
	defer func() {
		tf.completeMu.Lock()
		tf.complete = true
		tf.completeMu.Unlock()
	}()

	if tf.Complete() {
		tf.cancel()
		return nil
	}

	for {
		select {
		case <-tf.doneChan:
			tf.cancel()
			return nil
		case <-tf.context.Done():
			return fmt.Errorf("timeout waiting for transaction to be finalized: %s network %s", tf.txHash, tf.client.GetNetworkName())
		}
	}
}

// Complete returns if the finalizer has completed or not
func (tf *TransactionFinalizer) Complete() bool {
	tf.completeMu.Lock()
	defer tf.completeMu.Unlock()
	return tf.complete
}
