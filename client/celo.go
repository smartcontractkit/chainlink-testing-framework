package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/pkg/errors"

	celoContracts "github.com/smartcontractkit/integrations-framework/contracts/celo"

	"github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/crypto"
	"github.com/celo-org/celo-blockchain/ethclient"
	"github.com/rs/zerolog/log"
)

type CeloBlock struct {
	*types.Block
}

func (e *CeloBlock) GetHash() HashInterface {
	return e.Hash()
}

// CeloClients wraps the client and the BlockChain network to interact with an Celo EVM based Blockchain with multiple nodes
type CeloClients struct {
	DefaultClient *CeloClient
	Clients       []*CeloClient
}

// GetNetworkName gets the ID of the chain that the clients are connected to
func (e *CeloClients) GetNetworkName() string {
	return e.DefaultClient.GetNetworkName()
}

// GetID gets client ID, node number it's connected to
func (e *CeloClients) GetID() int {
	return e.DefaultClient.ID
}

// GasStats gets gas stats instance
func (e *CeloClients) GasStats() *GasStats {
	return e.DefaultClient.gasStats
}

// SetDefaultClient sets default client to perform calls to the network
func (e *CeloClients) SetDefaultClient(clientID int) error {
	if clientID > len(e.Clients) {
		return fmt.Errorf("client for node %d not found", clientID)
	}
	e.DefaultClient = e.Clients[clientID]
	return nil
}

// GetClients gets clients for all nodes connected
func (e *CeloClients) GetClients() []BlockchainClient {
	cl := make([]BlockchainClient, 0)
	for _, c := range e.Clients {
		cl = append(cl, c)
	}
	return cl
}

// SetID sets client ID (node)
func (e *CeloClients) SetID(id int) {
	e.DefaultClient.SetID(id)
}

// BlockNumber gets block number
func (e *CeloClients) BlockNumber(ctx context.Context) (uint64, error) {
	return e.DefaultClient.BlockNumber(ctx)
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *CeloClients) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	return e.DefaultClient.HeaderTimestampByNumber(ctx, bn)
}

// HeaderHashByNumber gets header hash by block number
func (e *CeloClients) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	return e.DefaultClient.HeaderHashByNumber(ctx, bn)
}

// Get gets default client as an interface{}
func (e *CeloClients) Get() interface{} {
	return e.DefaultClient
}

// CalculateTxGas calculates tx gas cost accordingly gas used plus buffer, converts it to big.Float for funding
func (e *CeloClients) CalculateTxGas(gasUsedValue *big.Int) (*big.Float, error) {
	return e.DefaultClient.CalculateTxGas(gasUsedValue)
}

// Fund funds a specified address with LINK token and or ETH from the given wallet
func (e *CeloClients) Fund(fromWallet BlockchainWallet, toAddress string, nativeAmount, linkAmount *big.Float) error {
	return e.DefaultClient.Fund(fromWallet, toAddress, nativeAmount, linkAmount)
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (e *CeloClients) ParallelTransactions(enabled bool) {
	for _, c := range e.Clients {
		c.ParallelTransactions(enabled)
	}
}

// Close tears down the all the clients
func (e *CeloClients) Close() error {
	for _, c := range e.Clients {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *CeloClients) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	for _, c := range e.Clients {
		c.AddHeaderEventSubscription(key, subscriber)
	}
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *CeloClients) DeleteHeaderEventSubscription(key string) {
	for _, c := range e.Clients {
		c.DeleteHeaderEventSubscription(key)
	}
}

// WaitForEvents is a blocking function that waits for all event subscriptions for all clients
func (e *CeloClients) WaitForEvents() error {
	g := errgroup.Group{}
	for _, c := range e.Clients {
		c := c
		g.Go(func() error {
			return c.WaitForEvents()
		})
	}
	return g.Wait()
}

// CeloClient wraps the client and the BlockChain network to interact with an EVM based Blockchain
type CeloClient struct {
	ID                  int
	Client              *ethclient.Client
	Network             BlockchainNetwork
	BorrowNonces        bool
	NonceMu             *sync.Mutex
	Nonces              map[string]uint64
	txQueue             chan common.Hash
	headerSubscriptions map[string]HeaderEventSubscription
	mutex               *sync.Mutex
	queueTransactions   bool
	gasStats            *GasStats
	doneChan            chan struct{}
}

// GetID gets client ID, node number it's connected to
func (e *CeloClient) GetID() int {
	return e.ID
}

// SetDefaultClient not used, only applicable to CeloClients
func (e *CeloClient) SetDefaultClient(_ int) error {
	return nil
}

// GetClients not used, only applicable to CeloClients
func (e *CeloClient) GetClients() []BlockchainClient {
	return []BlockchainClient{e}
}

// SuggestGasPrice gets suggested gas price
func (e *CeloClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := e.Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

// SetID sets client id, useful for multi-node networks
func (e *CeloClient) SetID(id int) {
	e.ID = id
}

// BlockNumber gets latest block number
func (e *CeloClient) BlockNumber(ctx context.Context) (uint64, error) {
	bn, err := e.Client.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return bn.NumberU64(), nil
}

// HeaderHashByNumber gets header hash by block number
func (e *CeloClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	h, err := e.Client.HeaderByNumber(ctx, bn)
	if err != nil {
		return "", err
	}
	return h.Hash().String(), nil
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *CeloClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	h, err := e.Client.HeaderByNumber(ctx, bn)
	if err != nil {
		return 0, err
	}
	return h.Time, nil
}

// CeloContractDeployer acts as a go-between function for general contract deployment
type CeloContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)

// NewCeloClient returns an instantiated instance of the Celo client that has connected to the server
func NewCeloClient(network BlockchainNetwork) (*CeloClient, error) {
	cl, err := ethclient.Dial(network.LocalURL())
	if err != nil {
		return nil, err
	}

	ec := &CeloClient{
		Network:             network,
		Client:              cl,
		BorrowNonces:        true,
		NonceMu:             &sync.Mutex{},
		Nonces:              make(map[string]uint64),
		txQueue:             make(chan common.Hash, 64), // Max buffer of 64 tx
		headerSubscriptions: map[string]HeaderEventSubscription{},
		mutex:               &sync.Mutex{},
		queueTransactions:   false,
		doneChan:            make(chan struct{}),
	}
	ec.gasStats = NewGasStats(ec.ID)
	go ec.newHeadersLoop()
	return ec, nil
}

// NewCeloClients returns an instantiated instance of all Celo client connected to all nodes
func NewCeloClients(network BlockchainNetwork) (*CeloClients, error) {
	ecl := &CeloClients{Clients: make([]*CeloClient, 0)}
	for idx, url := range network.URLs() {
		network.SetLocalURL(url)
		ec, err := NewCeloClient(network)
		if err != nil {
			return nil, err
		}
		ec.SetID(idx)
		ecl.Clients = append(ecl.Clients, ec)
	}
	ecl.DefaultClient = ecl.Clients[0]
	return ecl, nil
}

// GetNetworkName retrieves the ID of the network that the client interacts with
func (e *CeloClient) GetNetworkName() string {
	return e.Network.ID()
}

// Close tears down the current open Celo client
func (e *CeloClient) Close() error {
	e.doneChan <- struct{}{}
	e.Client.Close()
	return nil
}

// SuggestGasPrice gets suggested gas price
func (e *CeloClients) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := e.DefaultClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

// BorrowedNonces allows to handle nonces concurrently without requesting them every time
func (e *CeloClient) BorrowedNonces(n bool) {
	e.BorrowNonces = n
}

// GetNonce keep tracking of nonces per address, add last nonce for addr if the map is empty
func (e *CeloClient) GetNonce(ctx context.Context, addr common.Address) (uint64, error) {
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
		e.Nonces[addr.Hex()] += 1
		return e.Nonces[addr.Hex()], nil
	}
	lastNonce, err := e.Client.PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}

// Get returns the underlying client type to be used generically across the framework for switching
// network types
func (e *CeloClient) Get() interface{} {
	return e
}

// CalculateTxGas calculates tx gas cost accordingly gas used plus buffer, converts it to big.Float for funding
func (e *CeloClient) CalculateTxGas(gasUsed *big.Int) (*big.Float, error) {
	gp, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gpWei := gp.Mul(gp, OneGWei)
	log.Debug().Int64("Gas price", gp.Int64()).Msg("Suggested gas price")
	buf := big.NewInt(int64(e.Network.Config().GasEstimationBuffer))
	gasUsedWithBuf := gasUsed.Add(gasUsed, buf)
	cost := big.NewInt(1).Mul(gpWei, gasUsedWithBuf)
	log.Debug().Int64("TX Gas cost", cost.Int64()).Msg("Estimated tx gas cost with buffer")
	bf := new(big.Float).SetInt(cost)
	return big.NewFloat(1).Quo(bf, OneEth), nil
}

// GasStats gets gas stats instance
func (e *CeloClient) GasStats() *GasStats {
	return e.gasStats
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (e *CeloClient) ParallelTransactions(enabled bool) {
	e.queueTransactions = enabled
}

// Fund funds a specified address with LINK token and or ETH from the given wallet
func (e *CeloClient) Fund(
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
		linkInstance, err := celoContracts.NewLinkToken(linkAddress, e.Client)
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
func (e *CeloClient) SendTransaction(
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
	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return common.Hash{}, err
	}

	tx := types.NewTransaction(
		nonce,
		*callMsg.To,
		callMsg.Value,
		callMsg.Gas,
		callMsg.GasPrice,
		callMsg.FeeCurrency,
		callMsg.GatewayFeeRecipient,
		callMsg.GatewayFee,
		callMsg.Data,
	)

	txSigned, err := types.SignTx(tx, types.NewEIP155Signer(e.Network.ChainID()), privateKey)
	if err != nil {
		return common.Hash{}, err
	}
	if err := e.Client.SendTransaction(context.Background(), txSigned); err != nil {
		return common.Hash{}, err
	}
	return txSigned.Hash(), e.ProcessTransaction(txSigned.Hash())
}

// ProcessTransaction will queue or wait on a transaction depending on whether queue transactions is enabled
func (e *CeloClient) ProcessTransaction(txHash common.Hash) error {
	var txConfirmer HeaderEventSubscription
	if e.Network.Config().MinimumConfirmations == 0 {
		txConfirmer = &CeloInstantConfirmations{}
	} else {
		txConfirmer = NewCeloTransactionConfirmer(e, txHash, e.Network.Config().MinimumConfirmations)
	}

	e.AddHeaderEventSubscription(txHash.String(), txConfirmer)

	if !e.queueTransactions {
		defer e.DeleteHeaderEventSubscription(txHash.String())
		if err := txConfirmer.Wait(); err != nil {
			return err
		}
	}
	return nil
}

// DeployContract acts as a general contract deployment tool to an Celo chain
func (e *CeloClient) DeployContract(
	fromWallet BlockchainWallet,
	contractName string,
	deployer CeloContractDeployer,
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
		Str("Gas Cost", transaction.Cost().String()).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// TransactionCallMessage returns a filled Celo CallMsg object with suggest gas price and limit
func (e *CeloClient) TransactionCallMessage(
	from BlockchainWallet,
	to common.Address,
	value *big.Int,
	data []byte,
) (*celo.CallMsg, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	msg := celo.CallMsg{
		From:     common.HexToAddress(from.Address()),
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	msg.Gas = e.Network.Config().TransactionLimit + e.Network.Config().GasEstimationBuffer
	return &msg, nil
}

// TransactionOpts return the base binding transaction options to create a new valid tx for contract deployment
func (e *CeloClient) TransactionOpts(
	from BlockchainWallet,
	to common.Address,
	value *big.Int,
	data []byte,
) (*bind.TransactOpts, error) {
	callMsg, err := e.TransactionCallMessage(from, to, value, data)
	if err != nil {
		return nil, err
	}
	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	// TODO koteld: should we use e.Client.ChainID(context.Background()) or e.Network.ChainID() here?
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
func (e *CeloClient) WaitForEvents() error {
	queuedEvents := e.GetHeaderSubscriptions()
	g := errgroup.Group{}

	for subName, sub := range queuedEvents {
		subName := subName
		sub := sub
		g.Go(func() error {
			defer e.DeleteHeaderEventSubscription(subName)
			return sub.Wait()
		})
	}
	return g.Wait()
}

// GetHeaderSubscriptions returns a duplicate map of the queued transactions
func (e *CeloClient) GetHeaderSubscriptions() map[string]HeaderEventSubscription {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	newMap := map[string]HeaderEventSubscription{}
	for k, v := range e.headerSubscriptions {
		newMap[k] = v
	}
	return newMap
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *CeloClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.headerSubscriptions[key] = subscriber
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *CeloClient) DeleteHeaderEventSubscription(key string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.headerSubscriptions, key)
}

func (e *CeloClient) newHeadersLoop() {
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

func (e *CeloClient) subscribeToNewHeaders() error {
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
			e.receiveHeader(header)
		case <-e.doneChan:
			return nil
		}
	}
}

func (e *CeloClient) receiveHeader(header *types.Header) {
	log.Debug().
		Str("Network", e.Network.ID()).
		Int("Node", e.ID).
		Str("Hash", header.Hash().String()).
		Str("Number", header.Number.String()).
		Msg("Received block header")

	subs := e.GetHeaderSubscriptions()
	block, err := e.Client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Err(fmt.Errorf("error fetching block by number: %v", err))
	}

	celoBlock := &CeloBlock{block}

	g := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		g.Go(func() error {
			return sub.ReceiveBlock(NodeBlock{NodeID: e.ID, BlockInterface: celoBlock})
		})
	}
	if err := g.Wait(); err != nil {
		log.Err(fmt.Errorf("error on sending block to recivers: %v", err))
	}
}

func (e *CeloClient) isTxConfirmed(txHash common.Hash) (bool, error) {
	tx, isPending, err := e.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return !isPending, err
	}
	if !isPending {
		receipt, err := e.Client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return !isPending, err
		}
		e.gasStats.AddClientTXData(TXGasData{
			TXHash:            txHash.String(),
			Value:             tx.Value().Uint64(),
			GasLimit:          tx.Gas(),
			GasUsed:           receipt.GasUsed,
			GasPrice:          tx.GasPrice().Uint64(),
			CumulativeGasUsed: receipt.CumulativeGasUsed,
		})
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
func (e *CeloClient) errorReason(
	b celo.ContractCaller,
	tx *types.Transaction,
	receipt *types.Receipt,
) (string, error) {
	chID, err := e.Client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(chID))
	if err != nil {
		return "", err
	}
	callMsg := celo.CallMsg{
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

// CeloTransactionConfirmer is an implementation of HeaderEventSubscription that checks whether tx are confirmed
type CeloTransactionConfirmer struct {
	minConfirmations int
	confirmations    int
	eth              *CeloClient
	txHash           common.Hash
	doneChan         chan struct{}
	context          context.Context
	cancel           context.CancelFunc
}

// NewCeloTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewCeloTransactionConfirmer(eth *CeloClient, txHash common.Hash, minConfirmations int) *CeloTransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), eth.Network.Config().Timeout)
	tc := &CeloTransactionConfirmer{
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

// ReceiveBlock the implementation of the HeaderEventSubscription that receives each block and checks
// tx confirmation
func (t *CeloTransactionConfirmer) ReceiveBlock(block NodeBlock) error {
	if block.BlockInterface == nil {
		log.Info().Msg("Received nil block")
		return nil
	}
	confirmationLog := log.Debug().Str("Network", t.eth.Network.ID()).
		Str("Block Hash", block.GetHash().Hex()).
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
func (t *CeloTransactionConfirmer) Wait() error {
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

// CeloInstantConfirmations is a no-op confirmer as all transactions are instantly mined so no confs are needed
type CeloInstantConfirmations struct{}

// ReceiveBlock is a no-op
func (i *CeloInstantConfirmations) ReceiveBlock(block NodeBlock) error {
	return nil
}

// Wait is a no-op
func (i *CeloInstantConfirmations) Wait() error {
	return nil
}
