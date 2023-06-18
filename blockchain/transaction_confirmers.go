package blockchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// TransactionConfirmer is an implementation of HeaderEventSubscription that checks whether tx are confirmed
type TransactionConfirmer struct {
	minConfirmations      int
	confirmations         int
	client                EVMClient
	tx                    *types.Transaction
	doneChan              chan struct{}
	context               context.Context
	cancel                context.CancelFunc
	networkConfig         *EVMNetwork
	lastReceivedHeaderNum uint64
	complete              bool
	completeMu            sync.Mutex
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(client EVMClient, tx *types.Transaction, minConfirmations int) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	tc := &TransactionConfirmer{
		minConfirmations: minConfirmations,
		confirmations:    0,
		client:           client,
		tx:               tx,
		doneChan:         make(chan struct{}, 1),
		context:          ctx,
		cancel:           ctxCancel,
		networkConfig:    client.GetNetworkConfig(),
		complete:         false,
	}
	return tc
}

// ReceiveHeader the implementation of the HeaderEventSubscription that receives each header and checks
// tx confirmation
func (t *TransactionConfirmer) ReceiveHeader(header NodeHeader) error {
	if header.Number.Uint64() <= t.lastReceivedHeaderNum {
		return nil // Header with same number mined, disregard for confirming
	}
	t.lastReceivedHeaderNum = header.Number.Uint64()
	confirmationLog := log.Debug().
		Str("Network Name", t.networkConfig.Name).
		Str("Header Hash", header.Hash.Hex()).
		Str("Header Number", header.Number.String()).
		Str("Tx Hash", t.tx.Hash().String()).
		Uint64("Nonce", t.tx.Nonce()).
		Int("Minimum Confirmations", t.minConfirmations)
	isConfirmed, err := t.client.IsTxConfirmed(t.tx.Hash())
	if err != nil {
		return err
	} else if isConfirmed {
		t.confirmations++
	}
	if t.confirmations >= t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).Msg("Transaction confirmations met")
		t.complete = true
		t.doneChan <- struct{}{}
	} else {
		confirmationLog.Int("Current Confirmations", t.confirmations).Msg("Waiting on minimum confirmations")
	}
	return nil
}

// Wait is a blocking function that waits until the transaction is complete
func (t *TransactionConfirmer) Wait() error {
	defer func() {
		t.completeMu.Lock()
		t.complete = true
		t.completeMu.Unlock()
	}()

	if t.Complete() {
		t.cancel()
		return nil
	}

	for {
		select {
		case <-t.doneChan:
			t.cancel()
			return nil
		case <-t.context.Done():
			return fmt.Errorf("timeout waiting for transaction to confirm: %s network %s", t.tx.Hash(), t.client.GetNetworkName())
		}
	}
}

// Complete returns if the confirmer has completed or not
func (t *TransactionConfirmer) Complete() bool {
	t.completeMu.Lock()
	defer t.completeMu.Unlock()
	return t.complete
}

// InstantConfirmer is a near-instant confirmation method, primarily for optimistic L2s that have near-instant finalization
type InstantConfirmer struct {
	client       EVMClient
	txHash       common.Hash
	complete     bool // tracks if the subscription is completed or not
	completeChan chan struct{}
	completeMu   sync.Mutex
	context      context.Context
	cancel       context.CancelFunc
	// For events
	confirmed     bool // tracks the confirmation status of the subscription
	confirmedChan chan bool
	errorChan     chan error
}

func NewInstantConfirmer(
	client EVMClient,
	txHash common.Hash,
	confirmedChan chan bool,
	errorChan chan error,
) *InstantConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	return &InstantConfirmer{
		client:       client,
		txHash:       txHash,
		completeChan: make(chan struct{}, 1),
		context:      ctx,
		cancel:       ctxCancel,
		// For events
		confirmedChan: confirmedChan,
		errorChan:     errorChan,
	}
}

// ReceiveHeader does a quick check on if the tx is confirmed already
func (l *InstantConfirmer) ReceiveHeader(_ NodeHeader) error {
	var err error
	l.confirmed, err = l.client.IsTxConfirmed(l.txHash)
	if err != nil {
		if err.Error() == "not found" {
			log.Debug().Str("Tx", l.txHash.Hex()).Msg("Transaction not found on chain yet. Waiting to confirm.")
			return err
		}
		log.Error().Str("Tx", l.txHash.Hex()).Err(err).Msg("Error checking tx confirmed")
		if l.errorChan != nil {
			l.errorChan <- err
		}
		return err
	}
	log.Debug().Bool("Confirmed", l.confirmed).Str("Tx", l.txHash.Hex()).Msg("Instant Confirmation")
	if l.confirmed {
		l.completeChan <- struct{}{}
		if l.confirmedChan != nil {
			l.confirmedChan <- l.confirmed
		}
	}
	return nil
}

// Wait checks every header if the tx has been included on chain or not
func (l *InstantConfirmer) Wait() error {
	defer func() {
		l.completeMu.Lock()
		l.complete = true
		l.completeMu.Unlock()
	}()

	for {
		select {
		case <-l.completeChan:
			l.cancel()
			return nil
		case <-l.context.Done():
			return fmt.Errorf("timeout waiting for instant transaction to confirm after %s: %s",
				l.client.GetNetworkConfig().Timeout.String(), l.txHash.Hex())
		}
	}
}

// Complete returns if the transaction is complete or not
func (l *InstantConfirmer) Complete() bool {
	l.completeMu.Lock()
	defer l.completeMu.Unlock()
	return l.complete
}

// EventConfirmer confirms that an event is confirmed by a certain amount of headers
type EventConfirmer struct {
	eventName             string
	minConfirmations      int
	confirmations         int
	client                EVMClient
	event                 *types.Log
	waitChan              chan struct{}
	errorChan             chan error
	confirmedChan         chan bool
	context               context.Context
	cancel                context.CancelFunc
	lastReceivedHeaderNum uint64
	complete              bool
}

// NewEventConfirmer returns a new instance of the event confirmer that waits for on-chain minimum
// confirmations
func NewEventConfirmer(
	eventName string,
	client EVMClient,
	event *types.Log,
	minConfirmations int,
	confirmedChan chan bool,
	errorChan chan error,
) *EventConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	tc := &EventConfirmer{
		eventName:        eventName,
		minConfirmations: minConfirmations,
		confirmations:    0,
		client:           client,
		event:            event,
		waitChan:         make(chan struct{}, 1),
		errorChan:        errorChan,
		confirmedChan:    confirmedChan,
		context:          ctx,
		cancel:           ctxCancel,
		complete:         false,
	}
	return tc
}

// ReceiveHeader will attempt to confirm an event for the chain's configured minimum confirmed headers. Errors encountered
// are sent along the eventErrorChan, and the result of confirming the event is sent to eventConfirmedChan.
func (e *EventConfirmer) ReceiveHeader(header NodeHeader) error {
	if header.Number.Uint64() <= e.lastReceivedHeaderNum {
		return nil
	}
	e.lastReceivedHeaderNum = header.Number.Uint64()
	confirmed, removed, err := e.client.IsEventConfirmed(e.event)
	if err != nil {
		e.errorChan <- err
		return err
	}
	if removed {
		e.confirmedChan <- false
		e.complete = true
		return nil
	}
	if confirmed {
		e.confirmations++
	}
	if e.confirmations >= e.minConfirmations {
		e.confirmedChan <- true
		e.complete = true
	}
	return nil
}

// Wait until the event fully presents as complete
func (e *EventConfirmer) Wait() error {
	defer func() { e.complete = true }()
	for {
		select {
		case <-e.waitChan:
			e.cancel()
			return nil
		case <-e.context.Done():
			return fmt.Errorf("timeout waiting for event to confirm after %s: %s",
				e.client.GetNetworkConfig().Timeout.String(), e.event.TxHash.Hex())
		}
	}
}

// Complete returns if the confirmer is done, whether confirmation was successful or not
func (e *EventConfirmer) Complete() bool {
	return e.complete
}

// GetHeaderSubscriptions returns a duplicate map of the queued transactions
func (e *EthereumClient) GetHeaderSubscriptions() map[string]HeaderEventSubscription {
	e.subscriptionMutex.Lock()
	defer e.subscriptionMutex.Unlock()

	newMap := map[string]HeaderEventSubscription{}
	for k, v := range e.headerSubscriptions {
		newMap[k] = v
	}
	return newMap
}

// subscribeToNewHeaders kicks off the primary header subscription for the test loop
func (e *EthereumClient) subscribeToNewHeaders() error {
	headerChannel := make(chan *SafeEVMHeader)
	subscription, err := e.SubscribeNewHeaders(context.Background(), headerChannel)
	if err != nil {
		return err
	}
	log.Info().Str("Network", e.NetworkConfig.Name).Msg("Subscribed to new block headers")

	go e.headerSubscriptionLoop(subscription, headerChannel)
	return nil
}

// headerSubscriptionLoop receives new headers, and handles subscription errors when they pop up
func (e *EthereumClient) headerSubscriptionLoop(subscription ethereum.Subscription, headerChannel chan *SafeEVMHeader) {
	lastHeaderNumber := uint64(0)
	for {
		select {
		case err := <-subscription.Err(): // Most subscription errors are temporary RPC downtime, so let's poll to resubscribe
			log.Error().Err(err).Msg("Error while subscribed to new headers, likely RPC errors. Attempting Reconnect.")
			subscription.Unsubscribe()

			subscription = e.resubscribeLoop(headerChannel, lastHeaderNumber)
		case header := <-headerChannel:
			lastHeaderNumber = header.Number.Uint64()
			e.receiveHeader(header)
		case <-e.doneChan:
			log.Debug().Str("Network", e.NetworkConfig.Name).Msg("Subscription cancelled")
			e.Client.Close()
			return
		}
	}
}

// resubscribeLoop polls the RPC connection until it comes back up, and then resubscribes to new headers
func (e *EthereumClient) resubscribeLoop(headerChannel chan *SafeEVMHeader, lastHeaderNumber uint64) ethereum.Subscription {
	rpcDegradedTime, rpcDegradedNotifyTime := time.Now(), time.Now()
	timeout, ticker := time.NewTimer(e.NetworkConfig.Timeout.Duration), time.NewTicker(time.Second)
	for { // Loop to resubscribe to new headers,
		select {
		case <-timeout.C: // Timeout waiting for RPC to come back up
			log.Error().Str("Time waiting", time.Since(rpcDegradedTime).String()).Msg("RPC connection still down, timed out waiting for it to come back up")
			e.Client.Close()
			return nil
		case <-ticker.C: // Poll the RPC connection to see if it's back up
			if time.Since(rpcDegradedNotifyTime) >= time.Second*10 { // Periodically inform that we're still waiting for RPC to come back up
				log.Warn().Str("Time waiting", time.Since(rpcDegradedTime).String()).Msg("RPC connection still down, waiting for it to come back up")
				rpcDegradedNotifyTime = time.Now()
			}
			subscription, err := e.SubscribeNewHeaders(context.Background(), headerChannel)
			if err == nil { // No error on resubscription, RPC connection restored, back to regularly scheduled programming
				ticker.Stop()
				log.Info().Str("Time waiting", time.Since(rpcDegradedTime).String()).Msg("RPC connection and subscription restored")
				err = e.backfillMissedBlocks(lastHeaderNumber, headerChannel)
				if err != nil {
					log.Error().Err(err).Msg("Error backfilling missed blocks, subscriptions may be out of sync")
				}
				return subscription
			}
			log.Trace().Err(err).Msg("Error trying to resubscribe to new headers, likely RPC down")
		}
	}
}

// backfillMissedBlocks checks if there are any missed blocks since a bad connection was detected, and if so, backfills them
// to our header channel
func (e *EthereumClient) backfillMissedBlocks(lastBlockSeen uint64, headerChannel chan *SafeEVMHeader) error {
	start := time.Now()
	latestBlockNumber, err := e.LatestBlockNumber(context.Background())
	if err != nil {
		return err
	}
	if latestBlockNumber <= lastBlockSeen {
		log.Info().Msg("No missed blocks to backfill")
		return nil
	}
	log.Info().
		Uint64("Last Block Seen", lastBlockSeen).
		Uint64("Latest Block", latestBlockNumber).
		Msg("Backfilling missed blocks since RPC connection issues")
	for i := lastBlockSeen + 1; i <= latestBlockNumber; i++ {
		header, err := e.HeaderByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		log.Trace().Uint64("Number", i).Str("Hash", header.Hash.Hex()).Msg("Backfilling block")
		headerChannel <- header
	}
	log.Info().
		Uint64("Backfilled blocks", latestBlockNumber-lastBlockSeen).
		Str("Time", time.Since(start).String()).
		Msg("Finished backfilling missed blocks")
	return nil
}

// receiveHeader takes in a new header from the chain, and sends the header to all active header subscriptions
func (e *EthereumClient) receiveHeader(header *SafeEVMHeader) {
	if header == nil {
		log.Debug().Msg("Received Nil Header")
		return
	}
	headerValue := *header

	suggestedPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		suggestedPrice = big.NewInt(0)
		log.Err(err).
			Str("Header Hash", headerValue.Hash.String()).
			Msg("Error retrieving Suggested Gas Price for new block header")
	}
	log.Debug().
		Str("NetworkName", e.NetworkConfig.Name).
		Int("Node", e.ID).
		Str("Hash", headerValue.Hash.String()).
		Str("Number", headerValue.Number.String()).
		Str("Gas Price", suggestedPrice.String()).
		Msg("Received block header")

	subs := e.GetHeaderSubscriptions()

	g := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		g.Go(func() error {
			return sub.ReceiveHeader(NodeHeader{NodeID: e.ID, SafeEVMHeader: headerValue})
		})
	}
	if err := g.Wait(); err != nil {
		log.Err(fmt.Errorf("error on sending block header to receivers: %v", err))
	}
	if len(subs) > 0 {
		var subsRemoved uint
		for key, sub := range subs { // Cleanup subscriptions that might not have Wait called on them
			if sub.Complete() {
				subsRemoved++
				e.DeleteHeaderEventSubscription(key)
			}
		}
		if subsRemoved > 0 {
			log.Trace().
				Uint("Recently Removed", subsRemoved).
				Int("Active", len(e.GetHeaderSubscriptions())).
				Msg("Updated Header Subscriptions")
		}
	}
}

// errorReason decodes tx revert reason
func (e *EthereumClient) errorReason(
	b ethereum.ContractCaller,
	tx *types.Transaction,
	receipt *types.Receipt,
) (string, error) {
	chID, err := e.Client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	var msg *core.Message
	if e.NetworkConfig.SupportsEIP1559 {
		msg, err = core.TransactionToMessage(tx, types.LatestSignerForChainID(chID), nil)
	} else {
		msg, err = core.TransactionToMessage(tx, types.NewEIP155Signer(chID), nil)
	}
	if err != nil {
		return "", err
	}

	callMsg := ethereum.CallMsg{
		From:     msg.From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	_, txError := b.CallContract(context.Background(), callMsg, receipt.BlockNumber)
	if txError == nil {
		return "", errors.Wrap(err, "no error in CallContract")
	}
	errBytes, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(errBytes, &callErr)
	if err != nil {
		return "", err
	}
	// If the error data is blank
	if len(callErr.Data) == 0 {
		return callErr.Data, nil
	}
	// Some nodes prepend "Reverted " and we also remove the 0x
	trimmed := strings.TrimPrefix(callErr.Data, "Reverted ")[2:]
	data, err := hex.DecodeString(trimmed)
	if err != nil {
		return "", err
	}
	revert, err := abi.UnpackRevert(data)
	// If we can't decode the revert reason, return the raw data
	if err != nil {
		return callErr.Data, nil
	}
	return revert, nil
}
