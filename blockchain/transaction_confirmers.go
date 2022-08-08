package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// TransactionConfirmer is an implementation of HeaderEventSubscription that checks whether tx are confirmed
type TransactionConfirmer struct {
	minConfirmations     int
	confirmations        int
	client               EVMClient
	tx                   *types.Transaction
	doneChan             chan struct{}
	context              context.Context
	cancel               context.CancelFunc
	networkConfig        *EVMNetwork
	lastReceivedBlockNum uint64
	complete             bool
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(client EVMClient, tx *types.Transaction, minConfirmations int) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout)
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

// ReceiveBlock the implementation of the HeaderEventSubscription that receives each block and checks
// tx confirmation
func (t *TransactionConfirmer) ReceiveBlock(block NodeBlock) error {
	if block.NumberU64() <= t.lastReceivedBlockNum {
		return nil // Block with same number mined, disregard for confirming
	}
	t.lastReceivedBlockNum = block.NumberU64()
	confirmationLog := log.Debug().
		Str("Network Name", t.networkConfig.Name).
		Str("Block Hash", block.Hash().Hex()).
		Str("Block Number", block.Number().String()).
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
	defer func() { t.complete = true }()
	for {
		select {
		case <-t.doneChan:
			t.cancel()
			return nil
		case <-t.context.Done():
			return fmt.Errorf("timeout waiting for transaction to confirm: %s", t.tx.Hash())
		}
	}
}

// Complete returns if the confirmer has completed or not
func (t *TransactionConfirmer) Complete() bool {
	return t.complete
}

// L2TxConfirmer is a near-instant confirmation method, primarily for optimistic L2s that have near-instant finalization
type L2TxConfirmer struct {
	client   EVMClient
	txHash   common.Hash
	complete bool
	context  context.Context
	cancel   context.CancelFunc
}

func NewL2TxConfirmer(client EVMClient, txHash common.Hash) *L2TxConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout)
	return &L2TxConfirmer{
		client:  client,
		txHash:  txHash,
		context: ctx,
		cancel:  ctxCancel,
	}
}

// ReceiveBlock does a quick check on if the tx is confirmed already
func (l *L2TxConfirmer) ReceiveBlock(block NodeBlock) error {
	confirmed, err := l.client.IsTxConfirmed(l.txHash)
	if err != nil {
		return err
	}
	l.complete = confirmed
	return nil
}

// Wait checks every 50 milliseconds if the tx has been confirmed or not
func (l *L2TxConfirmer) Wait() error {
	countdown := time.NewTimer(0)
	defer func() {
		l.cancel()
		countdown.Stop()
		l.complete = true
	}()
	if l.complete {
		return nil
	}

	for {
		select {
		case <-countdown.C:
			confirmed, err := l.client.IsTxConfirmed(l.txHash)
			if err != nil {
				return err
			}
			if confirmed {
				log.Debug().Str("Hash", l.txHash.Hex()).Msg("L2 Tx Confirmed")
				return nil
			}
			countdown = time.NewTimer(time.Millisecond * 50)
		case <-l.context.Done():
			return fmt.Errorf("timeout waiting for transaction to confirm: %s", l.txHash.Hex())
		}
	}
}

// Complete is a no-op
func (l *L2TxConfirmer) Complete() bool {
	return l.complete
}

// EventConfirmer confirms that an event is confirmed by a certain amount of blocks
type EventConfirmer struct {
	eventName            string
	minConfirmations     int
	confirmations        int
	client               EVMClient
	event                *types.Log
	waitChan             chan struct{}
	errorChan            chan error
	confirmedChan        chan bool
	context              context.Context
	cancel               context.CancelFunc
	lastReceivedBlockNum uint64
	complete             bool
}

// NewEventConfirmer returns a new instance of the event confirmer that waits for on-chain minimum
// confirmations
func NewEventConfirmer(
	eventName string,
	client EVMClient,
	event *types.Log,
	minConfirmations int,
	errorChan chan error,
	confirmedChan chan bool,
) *EventConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), client.GetNetworkConfig().Timeout)
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

// ProcessEvent will attempt to confirm an event for the chain's configured minimum confirmed blocks. Errors encountered
// are sent along the eventErrorChan, and the result of confirming the event is sent to eventConfirmedChan.
func (e *EventConfirmer) ReceiveBlock(block NodeBlock) error {
	if block.NumberU64() <= e.lastReceivedBlockNum {
		return nil
	}
	e.lastReceivedBlockNum = block.NumberU64()
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
			return fmt.Errorf("timeout waiting for transaction to confirm: %s", e.event.TxHash.Hex())
		}
	}
}

// Complete returns if the event has officially been confirmed (true or false)
func (e *EventConfirmer) Complete() bool {
	return e.complete
}

// GetNonce keep tracking of nonces per address, add last nonce for addr if the map is empty
func (e *EthereumClient) GetNonce(ctx context.Context, addr common.Address) (uint64, error) {
	if e.BorrowNonces {
		e.NonceMu.Lock()
		defer e.NonceMu.Unlock()
		if _, ok := e.Nonces[addr.Hex()]; !ok {
			lastNonce, err := e.Client.PendingNonceAt(ctx, addr)
			if err != nil {
				return 0, err
			}
			e.Nonces[addr.Hex()] = lastNonce
			return lastNonce, nil
		}
		e.Nonces[addr.Hex()]++
		return e.Nonces[addr.Hex()], nil
	}
	lastNonce, err := e.Client.PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}

// GetHeaderSubscriptions returns a duplicate map of the queued transactions
func (e *EthereumClient) GetHeaderSubscriptions() map[string]HeaderEventSubscription {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	newMap := map[string]HeaderEventSubscription{}
	for k, v := range e.headerSubscriptions {
		newMap[k] = v
	}
	return newMap
}

// subscribeToNewHeaders
func (e *EthereumClient) subscribeToNewHeaders() error {
	headerChannel := make(chan *types.Header)
	subscription, err := e.Client.SubscribeNewHead(context.Background(), headerChannel)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	log.Info().Str("Network", e.NetworkConfig.Name).Msg("Subscribed to new block headers")

	for {
		select {
		case err := <-subscription.Err():
			log.Error().Err(err).Msg("Error while subscribed to new headers, restarting subscription")
			subscription.Unsubscribe()

			subscription, err = e.Client.SubscribeNewHead(context.Background(), headerChannel)
			if err != nil {
				return err
			}
		case header := <-headerChannel:
			e.receiveHeader(header)
		case <-e.doneChan:
			log.Debug().Str("Network", e.NetworkConfig.Name).Msg("Subscription cancelled")
			return nil
		}
	}
}

// receiveHeader
func (e *EthereumClient) receiveHeader(header *types.Header) {
	if header == nil {
		log.Debug().Msg("Received Nil block")
		return
	}
	suggestedPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		suggestedPrice = big.NewInt(0)
		log.Err(err).
			Str("Block Hash", header.Hash().String()).
			Msg("Error retrieving Suggested Gas Price for new block header")
	}
	log.Debug().
		Str("NetworkName", e.NetworkConfig.Name).
		Int("Node", e.ID).
		Str("Hash", header.Hash().String()).
		Str("Number", header.Number.String()).
		Str("Gas Price", suggestedPrice.String()).
		Msg("Received block header")

	subs := e.GetHeaderSubscriptions()
	block, err := e.Client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		log.Err(fmt.Errorf("error fetching block by hash: %v", err))
	}
	if block == nil || header == nil {
		log.Debug().Msg("Received Nil block")
		return
	}

	g := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		g.Go(func() error {
			return sub.ReceiveBlock(NodeBlock{NodeID: e.ID, Block: *block})
		})
	}
	if err := g.Wait(); err != nil {
		log.Err(fmt.Errorf("error on sending block to receivers: %v", err))
	}
	for key, sub := range subs { // Cleanup subscriptions that might not have Wait called on them
		if sub.Complete() {
			e.DeleteHeaderEventSubscription(key)
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
	msg, err := tx.AsMessage(types.NewEIP155Signer(chID), nil)
	if err != nil {
		return "", err
	}
	callMsg := ethereum.CallMsg{
		From:     msg.From(),
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	res, err := b.CallContract(context.Background(), callMsg, receipt.BlockNumber)
	if err != nil {
		return "", errors.Wrap(err, "CallContract")
	}
	return abi.UnpackRevert(res)
}
