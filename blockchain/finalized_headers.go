package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	minBlocksToFinalize = 10
	finalityTimeout     = 45 * time.Minute
	finalizedHeaderKey  = "finalizedHeads"
)

var (
	globalFinalizedHeaderManager = sync.Map{}
)

// FinalizedHeader is an implementation of the HeaderEventSubscription interface.
// It keeps track of the latest finalized header for a network.
type FinalizedHeader struct {
	mutex           *sync.Mutex
	LatestFinalized *big.Int
	FinalizedAt     time.Time
	client          EVMClient
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
	if header.Number.Cmp(f.LatestFinalized) > 0 &&
		// assumption : it will take at least minBlocksToFinalize num of blocks to finalize
		// this is to reduce the number of calls to the client
		new(big.Int).Sub(header.Number, f.LatestFinalized).Cmp(big.NewInt(minBlocksToFinalize)) <= 0 {
		return nil
	}
	ctx, ctxCancel := context.WithTimeout(context.Background(), f.client.GetNetworkConfig().Timeout.Duration)
	lastFinalized, err := f.client.GetLatestFinalizedBlockHeader(ctx)
	ctxCancel()
	if err != nil {
		return fmt.Errorf("error getting latest finalized block header - network %s", f.client.GetNetworkName())
	}
	if lastFinalized.Number.Cmp(f.LatestFinalized) > 0 {
		f.mutex.Lock()
		f.LatestFinalized = lastFinalized.Number
		f.FinalizedAt = header.Timestamp
		f.mutex.Unlock()
		log.Info().
			Str("Finalized Header", lastFinalized.Number.String()).
			Str("Network", f.client.GetNetworkName()).
			Str("Finalized At", f.FinalizedAt.String()).
			Msg("new finalized header received")
	}
	return nil
}

// newGlobalFinalizedHeaderManager is a global manager for finalized headers per network.
// It is used to keep track of the latest finalized header for each network.
func newGlobalFinalizedHeaderManager(client EVMClient) *FinalizedHeader {
	if client.NetworkSimulated() || client.GetNetworkConfig().FinalityDepth > 0 {
		return nil
	}
	f, ok := globalFinalizedHeaderManager.Load(client.GetChainID().String())
	now := time.Now().UTC()
	// if there is no finalized header for this network or the last finalized header is older than 1 hour
	if !ok || f != nil && now.Sub(f.(*FinalizedHeader).FinalizedAt) > 1*time.Hour {
		ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
		lastFinalized, err := client.GetLatestFinalizedBlockHeader(ctx)
		ctxCancel()
		if err != nil {
			log.Err(fmt.Errorf("error getting latest finalized block header")).Msg("NewFinalizedHeader")
			return nil
		}
		f := &FinalizedHeader{
			mutex:           &sync.Mutex{},
			LatestFinalized: lastFinalized.Number,
			FinalizedAt:     time.Now().UTC(),
			client:          client,
		}
		globalFinalizedHeaderManager.Store(client.GetChainID().String(), f)
		client.AddHeaderEventSubscription(finalizedHeaderKey, f)
	}
	return f.(*FinalizedHeader)
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
}

func NewTransactionFinalizer(client EVMClient, txHdr *SafeEVMHeader, txHash common.Hash) *TransactionFinalizer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), finalityTimeout)
	tf := &TransactionFinalizer{
		lggr:          log.With().Str("txHash", txHash.String()).Str("Tx Block", txHdr.Number.String()).Logger(),
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
	// if simulated network, return
	if tf.client.NetworkSimulated() {
		return nil
	}
	isFinalized, by, at, err := tf.client.IsTxFinalized(tf.txHdr, &header.SafeEVMHeader)
	if err != nil {
		return err
	}
	lgEvent := tf.lggr.Info()
	if isFinalized {
		tf.lggr.Info().
			Str("Finalized Block", header.Number.String()).
			Str("Tx block", tf.txHdr.Number.String()).
			Msg("Found finalized log")
		tf.complete = true
		tf.doneChan <- struct{}{}
		tf.FinalizedBy = by
		tf.FinalizedAt = at
	} else {
		if tf.networkConfig.FinalityDepth > 0 {
			lgEvent.
				Str("Current Block", header.Number.String()).
				Uint64("Finality Depth", tf.networkConfig.FinalityDepth)
		} else {
			lgEvent.
				Str("Last Finalized Block", by.String())
		}
		lgEvent.Msg("Still Waiting for transaction log to be finalized")
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
