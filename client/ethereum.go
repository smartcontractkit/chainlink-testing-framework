package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"math/big"
	"sync"
	"time"

	ethContracts "github.com/smartcontractkit/integrations-framework/contracts/ethereum"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

var OneEth = big.NewFloat(1e18)

// EthereumClient wraps the client and the BlockChain network to interact with an EVM based Blockchain
type EthereumClient struct {
	Client  *ethclient.Client
	Network BlockchainNetwork

	txQueue             chan common.Hash
	headerSubscriptions map[string]HeaderEventSubscription
	mutex               *sync.Mutex
	queueTransactions   bool
	doneChan            chan struct{}
}

// ContractDeployer acts as a go-between function for general contract deployment
type ContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)

// NewEthereumClient returns an instantiated instance of the Ethereum client that has connected to the server
func NewEthereumClient(network BlockchainNetwork) (*EthereumClient, error) {
	cl, err := ethclient.Dial(network.URL())
	if err != nil {
		return nil, err
	}

	ec := &EthereumClient{
		Network:             network,
		Client:              cl,
		txQueue:             make(chan common.Hash, 64), // Max buffer of 64 tx
		headerSubscriptions: map[string]HeaderEventSubscription{},
		mutex:               &sync.Mutex{},
		queueTransactions:   false,
		doneChan:            make(chan struct{}),
	}
	go ec.newHeadersLoop()
	return ec, err
}

// Close tears down the current open Ethereum client
func (e *EthereumClient) Close() error {
	e.doneChan <- struct{}{}
	e.Client.Close()
	return nil
}

// Get returns the underlying client type to be used generically across the framework for switching
// network types
func (e *EthereumClient) Get() interface{} {
	return e
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (e *EthereumClient) ParallelTransactions(enabled bool) {
	e.queueTransactions = enabled
}

// Fund funds a specified address with LINK token and or ETH from the given wallet
func (e *EthereumClient) Fund(
	fromWallet BlockchainWallet,
	toAddress string,
	ethAmount, linkAmount *big.Float,
) error {
	ethAddress := common.HexToAddress(toAddress)
	// Send ETH if not 0
	if ethAmount != nil && big.NewFloat(0).Cmp(ethAmount) != 0 {
		eth := big.NewFloat(1).Mul(OneEth, ethAmount)
		log.Info().
			Str("Token", "ETH").
			Str("From", fromWallet.Address()).
			Str("To", toAddress).
			Str("Amount", eth.String()).
			Msg("Funding Address")
		_, err := e.SendTransaction(fromWallet, ethAddress, eth, nil)
		if err != nil {
			return err
		}
	}

	// Send LINK if not 0
	if linkAmount != nil && big.NewFloat(0).Cmp(linkAmount) != 0 {
		link := big.NewFloat(1).Mul(OneLINK, linkAmount)
		log.Info().
			Str("Token", "LINK").
			Str("From", fromWallet.Address()).
			Str("To", toAddress).
			Str("Amount", link.String()).
			Msg("Funding Address")
		linkAddress := common.HexToAddress(e.Network.Config().LinkTokenAddress)
		linkInstance, err := ethContracts.NewLinkToken(linkAddress, e.Client)
		if err != nil {
			return err
		}
		opts, err := e.TransactionOpts(fromWallet, ethAddress, nil, nil)
		if err != nil {
			return err
		}
		linkInt, _ := link.Int(nil)
		_, err = linkInstance.Transfer(opts, ethAddress, linkInt)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendTransaction sends a specified amount of WEI from a selected wallet to an address, and blocks until the
// transaction completes
func (e *EthereumClient) SendTransaction(
	from BlockchainWallet,
	to common.Address,
	value *big.Float,
	data []byte,
) (common.Hash, error) {
	intVal, _ := value.Int(nil)
	callMsg, err := e.TransactionCallMessage(from, to, intVal, data)
	if err != nil {
		return common.Hash{}, err
	}
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid private key: %v", err)
	}
	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := types.SignNewTx(privateKey, types.NewEIP2930Signer(e.Network.ChainID()), &types.LegacyTx{
		To:       callMsg.To,
		Value:    callMsg.Value,
		Data:     callMsg.Data,
		GasPrice: callMsg.GasPrice,
		Gas:      callMsg.Gas,
		Nonce:    nonce,
	})
	if err != nil {
		return common.Hash{}, err
	}
	if err := e.Client.SendTransaction(context.Background(), tx); err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), e.ProcessTransaction(tx.Hash())
}

// ProcessTransaction will queue or wait on a transaction depending on whether queue transactions is enabled
func (e *EthereumClient) ProcessTransaction(txHash common.Hash) error {
	var txConfirmer HeaderEventSubscription
	if e.Network.Config().MinimumConfirmations == 0 {
		txConfirmer = &InstantConfirmations{}
	} else {
		txConfirmer = NewTransactionConfirmer(e, txHash, e.Network.Config().MinimumConfirmations)
	}

	e.AddHeaderEventSubscription(txHash.String(), txConfirmer)

	if !e.queueTransactions {
		if err := txConfirmer.Wait(); err != nil {
			return err
		}
		e.DeleteHeaderEventSubscription(txHash.String())
	}
	return nil
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (e *EthereumClient) DeployContract(
	fromWallet BlockchainWallet,
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := e.TransactionOpts(fromWallet, common.Address{}, big.NewInt(0), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	contractAddress, transaction, contractInstance, err := deployer(opts, e.Client)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := e.ProcessTransaction(transaction.Hash()); err != nil {
		return nil, nil, nil, err
	}
	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", fromWallet.Address()).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// TransactionCallMessage returns a filled Ethereum CallMsg object with suggest gas price and limit
func (e *EthereumClient) TransactionCallMessage(
	from BlockchainWallet,
	to common.Address,
	value *big.Int,
	data []byte,
) (*ethereum.CallMsg, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	log.Debug().Str("Suggested Gas Price", gasPrice.String())
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(from.Address()),
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	msg.Gas = e.Network.Config().TransactionLimit + e.Network.Config().GasEstimationBuffer
	log.Debug().Uint64("Gas Limit", e.Network.Config().TransactionLimit).Uint64("Limit + Buffer", msg.Gas)
	return &msg, nil
}

// TransactionOpts return the base binding transaction options to create a new valid tx for contract deployment
func (e *EthereumClient) TransactionOpts(
	from BlockchainWallet,
	to common.Address,
	value *big.Int,
	data []byte,
) (*bind.TransactOpts, error) {
	callMsg, err := e.TransactionCallMessage(from, to, value, data)
	if err != nil {
		return nil, err
	}
	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, e.Network.ChainID())
	if err != nil {
		return nil, err
	}
	opts.From = callMsg.From
	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = value
	opts.GasPrice = callMsg.GasPrice
	opts.GasLimit = callMsg.Gas
	opts.Context = context.Background()

	return opts, nil
}

// WaitForEvents is a blocking function that waits for all event subscriptions that have been queued within the client.
func (e *EthereumClient) WaitForEvents() error {
	var errGroup error

	queuedTx := e.GetHeaderSubscriptions()
	wg := sync.WaitGroup{}
	wg.Add(len(queuedTx))

	for txHash, sub := range queuedTx {
		sub := sub
		txHash := txHash
		go func() {
			defer wg.Done()
			if err := sub.Wait(); err != nil {
				errGroup = multierror.Append(errGroup, err)
			}
			e.DeleteHeaderEventSubscription(txHash)
		}()
	}
	wg.Wait()

	return errGroup
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

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *EthereumClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.headerSubscriptions[key] = subscriber
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *EthereumClient) DeleteHeaderEventSubscription(key string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.headerSubscriptions, key)
}

func (e *EthereumClient) newHeadersLoop() {
	for {
		if err := e.subscribeToNewHeaders(); err != nil {
			log.Error().
				Str("Network", e.Network.ID()).
				Msgf("Error while subscribing to headers: %v", err.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}
	log.Debug().Str("Network", e.Network.ID()).Msg("Stopped subscribing to new headers")
}

func (e *EthereumClient) subscribeToNewHeaders() error {
	headerChannel := make(chan *types.Header)
	subscription, err := e.Client.SubscribeNewHead(context.Background(), headerChannel)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	log.Info().Str("Network", e.Network.ID()).Msg("Subscribed to new block headers")

	for {
		select {
		case err := <-subscription.Err():
			return err
		case header := <-headerChannel:
			log.Debug().
				Str("Network", e.Network.ID()).
				Str("Block Number", header.Number.String()).
				Msg("Received block header")

			queuedTx := e.GetHeaderSubscriptions()
			for _, sub := range queuedTx {
				if err := sub.ReceiveHeader(header); err != nil {
					log.Error().Msgf("Error while sending header to receiver: %v", err)
				}
			}
		case <-e.doneChan:
			return nil
		}
	}
}

func (e *EthereumClient) isTxConfirmed(txHash common.Hash) (bool, error) {
	tx, isPending, err := e.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return !isPending, err
	}
	if !isPending {
		receipt, err := e.Client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return !isPending, err
		}
		if receipt.Status == 0 {
			log.Warn().Str("TX Hash", txHash.Hex()).Msg("Transaction failed and was reverted!")
			reason, err := e.errorReason(e.Client, tx, receipt)
			if err != nil {
				return false, err
			}
			log.Debug().Str("Revert reason", reason).Send()
		}
	}
	return !isPending, err
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

// TransactionConfirmer is an implementation of HeaderEventSubscription that checks whether tx are confirmed
type TransactionConfirmer struct {
	minConfirmations int
	confirmations    int
	eth              *EthereumClient
	txHash           common.Hash
	doneChan         chan struct{}
	context          context.Context
	cancel           context.CancelFunc
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(eth *EthereumClient, txHash common.Hash, minConfirmations int) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), eth.Network.Config().Timeout)
	tc := &TransactionConfirmer{
		minConfirmations: minConfirmations,
		confirmations:    0,
		eth:              eth,
		txHash:           txHash,
		doneChan:         make(chan struct{}, 1),
		context:          ctx,
		cancel:           ctxCancel,
	}
	return tc
}

// ReceiveHeader the implementation of the HeaderEventSubscription that receives each block header and checks
// tx confirmation
func (t *TransactionConfirmer) ReceiveHeader(header *types.Header) error {
	block, err := t.eth.Client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		return err
	}
	confirmationLog := log.Debug().Str("Network", t.eth.Network.ID()).
		Str("Block Hash", block.Hash().Hex()).
		Str("Block Number", block.Number().String()).Str("Tx Hash", t.txHash.Hex()).
		Int("Minimum Confirmations", t.minConfirmations)
	isConfirmed, err := t.eth.isTxConfirmed(t.txHash)
	if err != nil {
		return err
	} else if isConfirmed {
		t.confirmations++
	}
	if t.confirmations == t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).Msg("Transaction confirmations met")
		t.doneChan <- struct{}{}
	} else if t.confirmations <= t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).Msg("Waiting on minimum confirmations")
	}
	return nil
}

// Wait is a blocking function that waits until the transaction is complete
func (t *TransactionConfirmer) Wait() error {
	for {
		select {
		case <-t.doneChan:
			t.cancel()
			return nil
		case <-t.context.Done():
			return fmt.Errorf("timeout waiting for transaction to confirm: %s", t.txHash.String())
		}
	}
}

// InstantConfirmations is a no-op confirmer as all transactions are instantly mined so no confs are needed
type InstantConfirmations struct{}

// ReceiveHeader is a no-op
func (i *InstantConfirmations) ReceiveHeader(*types.Header) error {
	return nil
}

// Wait is a no-op
func (i *InstantConfirmations) Wait() error {
	return nil
}

