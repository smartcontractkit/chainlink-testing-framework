package blockchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// How many successful connection polls to consider the connection restored after a disconnection
const targetSuccessCount = 3

// TransactionConfirmer is an implementation of HeaderEventSubscription that checks whether tx are confirmed
type TransactionConfirmer struct {
	minConfirmations      int
	confirmations         int
	client                EVMClient
	tx                    *types.Transaction
	doneChan              chan struct{}
	revertChan            chan struct{}
	context               context.Context
	cancel                context.CancelFunc
	networkConfig         *EVMNetwork
	lastReceivedHeaderNum uint64
	complete              bool
	completeMu            sync.Mutex
	l                     zerolog.Logger
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(client EVMClient, tx *types.Transaction, minConfirmations int, logger zerolog.Logger) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	tc := &TransactionConfirmer{
		minConfirmations: minConfirmations,
		confirmations:    0,
		client:           client,
		tx:               tx,
		doneChan:         make(chan struct{}, 1),
		revertChan:       make(chan struct{}, 1),
		context:          ctx,
		cancel:           ctxCancel,
		networkConfig:    client.GetNetworkConfig(),
		complete:         false,
		l:                logger,
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
	confirmationLog := t.l.Debug().
		Str("Network Name", t.networkConfig.Name).
		Str("Header Hash", header.Hash.Hex()).
		Str("Header Number", header.Number.String()).
		Str("Tx Hash", t.tx.Hash().String()).
		Uint64("Nonce", t.tx.Nonce()).
		Int("Minimum Confirmations", t.minConfirmations)
	isConfirmed, err := t.client.IsTxConfirmed(t.tx.Hash())
	if err != nil {
		if err.Error() == "not found" {
			confirmationLog.Msg("Transaction not found on chain yet. Waiting to confirm.")
			return nil
		}
		if strings.Contains(err.Error(), "transaction failed and was reverted") {
			t.revertChan <- struct{}{}
		}
		return err
	}
	if isConfirmed {
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
		case <-t.revertChan:
			return fmt.Errorf("transaction reverted: %s network %s", t.tx.Hash(), t.client.GetNetworkName())
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
	client        EVMClient
	txHash        common.Hash
	complete      bool // tracks if the subscription is completed or not
	context       context.Context
	cancel        context.CancelFunc
	newHeaderChan chan struct{}
	// For events
	confirmed     bool // tracks the confirmation status of the subscription
	confirmedChan chan bool
	log           zerolog.Logger
}

func NewInstantConfirmer(
	client EVMClient,
	txHash common.Hash,
	confirmedChan chan bool,
	_ chan error,
	logger zerolog.Logger,
) *InstantConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout.Duration)
	return &InstantConfirmer{
		client:  client,
		txHash:  txHash,
		context: ctx,
		cancel:  ctxCancel,
		// For events
		confirmedChan: confirmedChan,
		log:           logger,
	}
}

// ReceiveHeader does a quick check on if the tx is confirmed already
func (l *InstantConfirmer) ReceiveHeader(_ NodeHeader) error {
	return nil
}

// Wait checks every header if the tx has been included on chain or not
func (l *InstantConfirmer) Wait() error {
	defer func() {
		l.complete = true
		l.cancel()
	}()
	confirmed, err := l.checkConfirmed()
	if err != nil {
		return err
	}
	if confirmed {
		return nil
	}

	poll := time.NewTicker(time.Millisecond * 250)
	for {
		select {
		case <-poll.C:
			confirmed, err := l.checkConfirmed()
			if err != nil {
				return err
			}
			if confirmed {
				return nil
			}
		case <-l.newHeaderChan:
			confirmed, err := l.checkConfirmed()
			if err != nil {
				return err
			}
			if confirmed {
				return nil
			}
		case <-l.context.Done():
			return fmt.Errorf("timeout waiting for instant transaction to confirm after %s: %s",
				l.client.GetNetworkConfig().Timeout.String(), l.txHash.Hex())
		}
	}
}

func (l *InstantConfirmer) checkConfirmed() (bool, error) {
	confirmed, err := l.client.IsTxConfirmed(l.txHash)
	if err != nil {
		return false, err
	}
	l.confirmed = confirmed
	if confirmed {
		go func() {
			if l.confirmedChan != nil {
				l.confirmedChan <- true
			}
		}()
	}
	l.log.Trace().Bool("Confirmed", confirmed).Str("Hash", l.txHash.Hex()).Msg("Checked if transaction confirmed")
	return confirmed, nil
}

// Complete returns if the transaction is complete or not
func (l *InstantConfirmer) Complete() bool {
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
	e.l.Info().Str("Network", e.NetworkConfig.Name).Msg("Subscribed to new block headers")

	e.subscriptionWg.Add(1)
	go e.headerSubscriptionLoop(subscription, headerChannel)
	return nil
}

// headerSubscriptionLoop receives new headers, and handles subscription errors when they pop up
func (e *EthereumClient) headerSubscriptionLoop(subscription ethereum.Subscription, headerChannel chan *SafeEVMHeader) {
	defer e.subscriptionWg.Done()
	lastHeaderNumber := uint64(0)
	for {
		select {
		case err := <-subscription.Err(): // Most subscription errors are temporary RPC downtime, so let's poll to resubscribe
			e.l.Error().
				Str("URL Suffix", e.NetworkConfig.URL[len(e.NetworkConfig.URL)-6:]).
				Str("Network", e.NetworkConfig.Name).
				Err(err).
				Msg("Error while subscribed to new headers, likely RPC downtime. Attempting to resubscribe")
			e.connectionIssueCh <- time.Now()
			subscription = e.resubscribeLoop(headerChannel, lastHeaderNumber)
			e.connectionRestoredCh <- time.Now()
		case header := <-headerChannel:
			err := e.receiveHeader(header)
			if err != nil {
				e.l.Error().Str("Network", e.NetworkConfig.Name).Err(err).Msg("Error receiving header, possible RPC issues")
			} else {
				lastHeaderNumber = header.Number.Uint64()
			}
		case <-e.doneChan:
			e.l.Debug().Str("Network", e.NetworkConfig.Name).Msg("Subscription cancelled")
			e.Client.Close()
			return
		}
	}
}

// resubscribeLoop polls the RPC connection until it comes back up, and then resubscribes to new headers
func (e *EthereumClient) resubscribeLoop(headerChannel chan *SafeEVMHeader, lastHeaderNumber uint64) ethereum.Subscription {
	rpcDegradedTime := time.Now()
	pollInterval := time.Millisecond * 500
	checkTicker, informTicker := time.NewTicker(pollInterval), time.NewTicker(time.Second*10)
	reconnectAttempts, consecutiveSuccessCount := 0, 0
	e.l.Debug().Str("Network", e.NetworkConfig.Name).Str("Poll Interval", pollInterval.String()).Msg("Attempting to resubscribe to new headers")

	for {
		select {
		case <-checkTicker.C:
			reconnectAttempts++
			_, err := e.LatestBlockNumber(context.Background())
			if err != nil {
				e.l.Trace().Err(err).Msg("Error trying to resubscribe to new headers, likely RPC down")
				consecutiveSuccessCount = 0
				continue
			}
			consecutiveSuccessCount++
			e.l.Debug().
				Str("Network", e.NetworkConfig.Name).
				Str("URL Suffix", e.NetworkConfig.URL[len(e.NetworkConfig.URL)-6:]).
				Int("Target Success Count", targetSuccessCount).
				Int("Consecutive Success Count", consecutiveSuccessCount).
				Msg("RPC connection seems to be healthy, still checking")

			if consecutiveSuccessCount >= targetSuccessCount { // Make sure node is actually back up and not just a blip
				subscription, err := e.SubscribeNewHeaders(context.Background(), headerChannel)
				if err != nil {
					e.l.Error().Err(err).Msg("Error resubscribing to new headers, RPC connection is still coming up")
					consecutiveSuccessCount = 0
					continue
				}
				// No error on resubscription, RPC connection restored, back to regularly scheduled programming
				checkTicker.Stop()
				informTicker.Stop()
				e.l.Info().
					Int("Reconnect Attempts", reconnectAttempts).
					Str("URL Suffix", e.NetworkConfig.URL[len(e.NetworkConfig.URL)-6:]).
					Str("Time waiting", time.Since(rpcDegradedTime).String()).
					Msg("RPC connection and subscription restored")
				go e.backfillMissedBlocks(lastHeaderNumber, headerChannel)
				return subscription
			}
		case <-informTicker.C:
			e.l.Warn().
				Str("Network", e.NetworkConfig.Name).
				Str("URL Suffix", e.NetworkConfig.URL[len(e.NetworkConfig.URL)-6:]).
				Int("Reconnect Attempts", reconnectAttempts).
				Str("Time waiting", time.Since(rpcDegradedTime).String()).
				Str("RPC URL", e.NetworkConfig.URL).
				Msg("RPC connection still down, waiting for it to come back up")
		}
	}
}

// backfillMissedBlocks checks if there are any missed blocks since a bad connection was detected, and if so, backfills them
// to our header channel
func (e *EthereumClient) backfillMissedBlocks(lastBlockSeen uint64, headerChannel chan *SafeEVMHeader) {
	start := time.Now()
	latestBlockNumber, err := e.LatestBlockNumber(context.Background())
	if err != nil {
		e.l.Error().Str("Network", e.NetworkConfig.Name).Err(err).Msg("Error getting latest block number. Unable to backfill missed blocks. Subscription likely degraded")
		return
	}
	if latestBlockNumber <= lastBlockSeen {
		e.l.Info().Msg("No missed blocks to backfill")
		return
	}
	e.l.Info().
		Str("Network", e.NetworkConfig.Name).
		Uint64("Last Block Seen", lastBlockSeen).
		Uint64("Blocks Behind", latestBlockNumber-lastBlockSeen).
		Uint64("Latest Block", latestBlockNumber).
		Msg("Backfilling missed blocks since RPC connection issues")
	for i := lastBlockSeen + 1; i <= latestBlockNumber; i++ {
		header, err := e.HeaderByNumber(context.Background(), big.NewInt(int64(i))) //nolint
		if err != nil {
			e.l.Err(err).Uint64("Number", i).Msg("Error getting header, unable to backfill and process it")
			return
		}
		e.l.Trace().Str("Network", e.NetworkConfig.Name).Str("Hash", header.Hash.Hex()).Uint64("Number", i).Msg("Backfilling header")
		headerChannel <- header
	}
	e.l.Info().
		Str("Network", e.NetworkConfig.Name).
		Uint64("Backfilled blocks", latestBlockNumber-lastBlockSeen).
		Str("Time", time.Since(start).String()).
		Msg("Finished backfilling missed blocks")
}

// receiveHeader takes in a new header from the chain, and sends the header to all active header subscriptions
func (e *EthereumClient) receiveHeader(header *SafeEVMHeader) error {
	if header == nil {
		e.l.Trace().Msg("Received Nil Header")
		return nil
	}
	headerValue := *header

	e.l.Trace().
		Str("NetworkName", e.NetworkConfig.Name).
		Int("Node", e.ID).
		Str("Hash", headerValue.Hash.String()).
		Str("Number", headerValue.Number.String()).
		Msg("Received Header")

	safeHeader := NodeHeader{NodeID: e.ID, SafeEVMHeader: headerValue}
	subs := e.GetHeaderSubscriptions()
	e.l.Trace().
		Str("NetworkName", e.NetworkConfig.Name).
		Int("Node", e.ID).
		Interface("Map", subs).Msg("Active Header Subscriptions")

	g := errgroup.Group{}
	for _, sub := range subs {
		g.Go(func() error {
			return sub.ReceiveHeader(safeHeader)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error on sending block header to receivers: '%w'", err)
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
			e.l.Trace().
				Str("NetworkName", e.NetworkConfig.Name).
				Int("Node", e.ID).
				Uint("Recently Removed", subsRemoved).
				Int("Active", len(e.GetHeaderSubscriptions())).
				Msg("Updated Header Subscriptions")
		}
	}
	return nil
}

// ErrorReason decodes tx revert reason
func (e *EthereumClient) ErrorReason(
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
		return "", errors.New("couldn't find revert reason. CallContract did not fail")
	}
	return RPCErrorFromError(txError)
}

func RPCErrorFromError(txError error) (string, error) {
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
