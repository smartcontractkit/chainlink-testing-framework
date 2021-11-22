package client

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/config"
	"math/big"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

var (
	// OneGWei represents 1 GWei
	OneGWei = big.NewInt(1e9)
	// OneEth represents 1 Ethereum
	OneEth = big.NewFloat(1e18)
)

// EthereumMultinodeClient wraps the client and the BlockChain network to interact with an EVM based Blockchain with multiple nodes
type EthereumMultinodeClient struct {
	DefaultClient *EthereumClient
	Clients       []*EthereumClient
}

// CalculateTXSCost calculates TXs cost as a dirty estimation based on transactionLimit for that network
func (e *EthereumMultinodeClient) CalculateTXSCost(txs int64) (*big.Float, error) {
	return e.DefaultClient.CalculateTXSCost(txs)
}

// LoadWallets loads wallets from config
func (e *EthereumMultinodeClient) LoadWallets(cfg interface{}) error {
	pkStrings := cfg.(config.ETHNetwork).PrivateKeys
	wallets := make([]*EthereumWallet, 0)
	for _, pks := range pkStrings {
		w, err := NewEthereumWallet(pks)
		if err != nil {
			return err
		}
		wallets = append(wallets, w)
	}
	for _, c := range e.Clients {
		c.Wallets = wallets
	}
	return nil
}

// SetWallet sets default wallet
func (e *EthereumMultinodeClient) SetWallet(num int) error {
	if num > len(e.DefaultClient.Wallets) {
		return fmt.Errorf("no wallet #%d found for default client", num)
	}
	e.DefaultClient.DefaultWallet = e.DefaultClient.Wallets[num]
	return nil
}

// GetNetworkName gets the ID of the chain that the clients are connected to
func (e *EthereumMultinodeClient) GetNetworkName() string {
	return e.DefaultClient.GetNetworkName()
}

// GasStats gets gas stats instance
func (e *EthereumMultinodeClient) GasStats() *GasStats {
	return e.DefaultClient.gasStats
}

// SwitchNode sets default client to perform calls to the network
func (e *EthereumMultinodeClient) SwitchNode(clientID int) error {
	if clientID > len(e.Clients) {
		return fmt.Errorf("client for node %d not found", clientID)
	}
	e.DefaultClient = e.Clients[clientID]
	return nil
}

// GetClients gets clients for all nodes connected
func (e *EthereumMultinodeClient) GetClients() []BlockchainClient {
	cl := make([]BlockchainClient, 0)
	for _, c := range e.Clients {
		cl = append(cl, c)
	}
	return cl
}

// SetID sets client ID (node)
func (e *EthereumMultinodeClient) SetID(id int) {
	e.DefaultClient.SetID(id)
}

// BlockNumber gets block number
func (e *EthereumMultinodeClient) BlockNumber(ctx context.Context) (uint64, error) {
	return e.DefaultClient.BlockNumber(ctx)
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *EthereumMultinodeClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	return e.DefaultClient.HeaderTimestampByNumber(ctx, bn)
}

// HeaderHashByNumber gets header hash by block number
func (e *EthereumMultinodeClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	return e.DefaultClient.HeaderHashByNumber(ctx, bn)
}

// Get gets default client as an interface{}
func (e *EthereumMultinodeClient) Get() interface{} {
	return e.DefaultClient
}

// CalculateTxGas calculates tx gas cost accordingly gas used plus buffer, converts it to big.Float for funding
func (e *EthereumMultinodeClient) CalculateTxGas(gasUsedValue *big.Int) (*big.Float, error) {
	return e.DefaultClient.CalculateTxGas(gasUsedValue)
}

// Fund funds a specified address with LINK token and or ETH from the given wallet
func (e *EthereumMultinodeClient) Fund(toAddress string, nativeAmount *big.Float) error {
	return e.DefaultClient.Fund(toAddress, nativeAmount)
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (e *EthereumMultinodeClient) ParallelTransactions(enabled bool) {
	for _, c := range e.Clients {
		c.ParallelTransactions(enabled)
	}
}

// Close tears down the all the clients
func (e *EthereumMultinodeClient) Close() error {
	for _, c := range e.Clients {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *EthereumMultinodeClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	for _, c := range e.Clients {
		c.AddHeaderEventSubscription(key, subscriber)
	}
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *EthereumMultinodeClient) DeleteHeaderEventSubscription(key string) {
	for _, c := range e.Clients {
		c.DeleteHeaderEventSubscription(key)
	}
}

// WaitForEvents is a blocking function that waits for all event subscriptions for all clients
func (e *EthereumMultinodeClient) WaitForEvents() error {
	g := errgroup.Group{}
	for _, c := range e.Clients {
		c := c
		g.Go(func() error {
			return c.WaitForEvents()
		})
	}
	return g.Wait()
}

// EthereumClient wraps the client and the BlockChain network to interact with an EVM based Blockchain
type EthereumClient struct {
	ID                  int
	Client              *ethclient.Client
	NetworkConfig       *config.ETHNetwork
	Wallets             []*EthereumWallet
	DefaultWallet       *EthereumWallet
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

// CalculateTXSCost calculates TXs cost as a dirty estimation based on transactionLimit for that network
func (e *EthereumClient) CalculateTXSCost(txs int64) (*big.Float, error) {
	txsLimit := e.NetworkConfig.TransactionLimit
	gasPrice, err := e.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gpFloat := big.NewFloat(1).SetInt(gasPrice)
	oneGWei := big.NewFloat(1).SetInt(OneGWei)
	gpGWei := big.NewFloat(1).Quo(gpFloat, oneGWei)
	log.Debug().Str("Gas price (GWei)", gpGWei.String()).Msg("Suggested gas price")
	txl := big.NewFloat(1).SetUint64(txsLimit)
	oneTx := big.NewFloat(1).Mul(txl, gpFloat)
	transactions := big.NewFloat(1).SetInt64(txs)
	totalWei := big.NewFloat(1).Mul(oneTx, transactions)
	totalETH := big.NewFloat(1).Quo(totalWei, OneEth)
	log.Debug().Str("ETH", totalETH.String()).Int64("TXs", txs).Msg("Calculated required ETH")
	return totalETH, nil
}

// LoadWallets loads wallets from config
func (e *EthereumClient) LoadWallets(cfg interface{}) error {
	pkStrings := cfg.(*config.ETHNetwork).PrivateKeys
	for _, pks := range pkStrings {
		w, err := NewEthereumWallet(pks)
		if err != nil {
			return err
		}
		e.Wallets = append(e.Wallets, w)
	}
	if len(e.Wallets) == 0 {
		return fmt.Errorf("no private keys found to load wallets")
	}
	e.DefaultWallet = e.Wallets[0]
	return nil
}

// SetWallet sets default wallet
func (e *EthereumClient) SetWallet(num int) error {
	if num > len(e.Wallets) {
		return fmt.Errorf("no wallet #%d found for default client", num)
	}
	e.DefaultWallet = e.Wallets[num]
	return nil
}

// SwitchNode not used, only applicable to EthereumMultinodeClient
func (e *EthereumClient) SwitchNode(_ int) error {
	return nil
}

// GetClients not used, only applicable to EthereumMultinodeClient
func (e *EthereumClient) GetClients() []BlockchainClient {
	return []BlockchainClient{e}
}

// SuggestGasPrice gets suggested gas price
func (e *EthereumClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := e.Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

// SetID sets client id, useful for multi-node networks
func (e *EthereumClient) SetID(id int) {
	e.ID = id
}

// BlockNumber gets latest block number
func (e *EthereumClient) BlockNumber(ctx context.Context) (uint64, error) {
	bn, err := e.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

// HeaderHashByNumber gets header hash by block number
func (e *EthereumClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	h, err := e.Client.HeaderByNumber(ctx, bn)
	if err != nil {
		return "", err
	}
	return h.Hash().String(), nil
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *EthereumClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	h, err := e.Client.HeaderByNumber(ctx, bn)
	if err != nil {
		return 0, err
	}
	return h.Time, nil
}

// ContractDeployer acts as a go-between function for general contract deployment
type ContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)

// NewEthereumClient returns an instantiated instance of the Ethereum client that has connected to the server
func NewEthereumClient(networkSettings *config.ETHNetwork) (*EthereumClient, error) {
	log.Info().
		Str("ID", networkSettings.ID).
		Str("URL", networkSettings.URL).
		Interface("Settings", networkSettings).
		Msg("Connecting client")
	cl, err := ethclient.Dial(networkSettings.URL)
	if err != nil {
		return nil, err
	}

	ec := &EthereumClient{
		NetworkConfig:       networkSettings,
		Client:              cl,
		BorrowNonces:        true,
		NonceMu:             &sync.Mutex{},
		Wallets:             make([]*EthereumWallet, 0),
		Nonces:              make(map[string]uint64),
		txQueue:             make(chan common.Hash, 64), // Max buffer of 64 tx
		headerSubscriptions: map[string]HeaderEventSubscription{},
		mutex:               &sync.Mutex{},
		queueTransactions:   false,
		doneChan:            make(chan struct{}),
	}
	if err := ec.LoadWallets(networkSettings); err != nil {
		return nil, err
	}
	ec.gasStats = NewGasStats(ec.ID)
	go ec.newHeadersLoop()
	return ec, nil
}

// NewEthereumMultiNodeClient returns an instantiated instance of all Ethereum client connected to all nodes
func NewEthereumMultiNodeClient(_ string,
	networkConfig map[string]interface{},
	urls []*url.URL,
) (BlockchainClient, error) {
	networkSettings := &config.ETHNetwork{}
	err := UnmarshalNetworkConfig(networkConfig, networkSettings)
	if err != nil {
		return nil, err
	}
	log.Info().
		Interface("URLs", networkSettings.URLs).
		Msg("Connecting multi-node client")

	ecl := &EthereumMultinodeClient{}
	for _, envURL := range urls {
		networkSettings.URLs = append(networkSettings.URLs, envURL.String())
	}
	for idx, networkURL := range networkSettings.URLs {
		networkSettings.URL = networkURL
		ec, err := NewEthereumClient(networkSettings)
		if err != nil {
			return nil, err
		}
		ec.SetID(idx)
		ecl.Clients = append(ecl.Clients, ec)
	}
	ecl.DefaultClient = ecl.Clients[0]
	return ecl, nil
}

// EthereumMultiNodeURLs returns the websocket URLs for a deployed Ethereum multi-node setup
func EthereumMultiNodeURLs(e *environment.Environment) ([]*url.URL, error) {
	return e.Charts.Connections("geth").LocalURLsByPort("ws-rpc", environment.WS)
}

// GetNetworkName retrieves the ID of the network that the client interacts with
func (e *EthereumClient) GetNetworkName() string {
	return e.NetworkConfig.ID
}

// Close tears down the current open Ethereum client
func (e *EthereumClient) Close() error {
	e.doneChan <- struct{}{}
	e.Client.Close()
	return nil
}

// SuggestGasPrice gets suggested gas price
func (e *EthereumMultinodeClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := e.DefaultClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

// BorrowedNonces allows to handle nonces concurrently without requesting them every time
func (e *EthereumClient) BorrowedNonces(n bool) {
	e.BorrowNonces = n
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
func (e *EthereumClient) Get() interface{} {
	return e
}

// CalculateTxGas calculates tx gas cost accordingly gas used plus buffer, converts it to big.Float for funding
func (e *EthereumClient) CalculateTxGas(gasUsed *big.Int) (*big.Float, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background()) // Wei
	if err != nil {
		return nil, err
	}
	buffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)
	gasUsedWithBuffer := gasUsed.Add(gasUsed, buffer)
	cost := big.NewFloat(0).SetInt(big.NewInt(1).Mul(gasPrice, gasUsedWithBuffer))
	costInEth := big.NewFloat(0).Quo(cost, OneEth)
	costInEthFloat, _ := costInEth.Float64()

	log.Debug().Float64("ETH", costInEthFloat).Msg("Estimated tx gas cost with buffer")
	return costInEth, nil
}

// GasStats gets gas stats instance
func (e *EthereumClient) GasStats() *GasStats {
	return e.gasStats
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (e *EthereumClient) ParallelTransactions(enabled bool) {
	e.queueTransactions = enabled
}

// Fund sends some ETH to an address
func (e *EthereumClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	ethAddress := common.HexToAddress(toAddress)
	if amount != nil && big.NewFloat(0).Cmp(amount) != 0 {
		eth := big.NewFloat(1).Mul(OneEth, amount)
		log.Info().
			Str("Token", "ETH").
			Str("From", e.DefaultWallet.Address()).
			Str("To", toAddress).
			Str("Amount", amount.String()).
			Msg("Funding Address")
		_, err := e.SendTransaction(e.DefaultWallet, ethAddress, eth, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendTransaction sends a specified amount of WEI from a selected wallet to an address, and blocks until the
// transaction completes
func (e *EthereumClient) SendTransaction(
	from *EthereumWallet,
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

	tx, err := types.SignNewTx(privateKey, types.NewEIP2930Signer(big.NewInt(e.NetworkConfig.ChainID)), &types.LegacyTx{
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
	if e.NetworkConfig.MinimumConfirmations == 0 {
		txConfirmer = &InstantConfirmations{}
	} else {
		txConfirmer = NewTransactionConfirmer(e, txHash, e.NetworkConfig.MinimumConfirmations)
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

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (e *EthereumClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := e.TransactionOpts(e.DefaultWallet, common.Address{}, big.NewInt(0), nil)
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
		Str("From", e.DefaultWallet.Address()).
		Str("Gas Cost", transaction.Cost().String()).
		Str("NetworkConfig", e.NetworkConfig.ID).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// TransactionCallMessage returns a filled Ethereum CallMsg object with suggest gas price and limit
func (e *EthereumClient) TransactionCallMessage(
	from *EthereumWallet,
	to common.Address,
	value *big.Int,
	data []byte,
) (*ethereum.CallMsg, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(from.Address()),
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	msg.Gas = e.NetworkConfig.TransactionLimit + e.NetworkConfig.GasEstimationBuffer
	return &msg, nil
}

// TransactionOpts return the base binding transaction options to create a new valid tx for contract deployment
func (e *EthereumClient) TransactionOpts(
	from *EthereumWallet,
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
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(e.NetworkConfig.ChainID))
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
				Str("NetworkConfig", e.NetworkConfig.ID).
				Msgf("Error while subscribing to headers: %v", err.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}
	log.Debug().Str("NetworkConfig", e.NetworkConfig.ID).Msg("Stopped subscribing to new headers")
}

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
			return err
		case header := <-headerChannel:
			e.receiveHeader(header)
		case <-e.doneChan:
			log.Debug().Str("Network", e.NetworkConfig.Name).Msg("Subscription cancelled")
			return nil
		}
	}
}

func (e *EthereumClient) receiveHeader(header *types.Header) {
	log.Debug().
		Str("NetworkConfig", e.NetworkConfig.ID).
		Int("Node", e.ID).
		Str("Hash", header.Hash().String()).
		Str("Number", header.Number.String()).
		Msg("Received block header")

	subs := e.GetHeaderSubscriptions()
	block, err := e.Client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Err(fmt.Errorf("error fetching block by number: %v", err))
	}

	g := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		g.Go(func() error {
			return sub.ReceiveBlock(NodeBlock{NodeID: e.ID, Block: block})
		})
	}
	if err := g.Wait(); err != nil {
		log.Err(fmt.Errorf("error on sending block to recivers: %v", err))
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), eth.NetworkConfig.Timeout)
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

// ReceiveBlock the implementation of the HeaderEventSubscription that receives each block and checks
// tx confirmation
func (t *TransactionConfirmer) ReceiveBlock(block NodeBlock) error {
	if block.Block == nil {
		// strange, but happening on Kovan
		log.Info().Msg("Received nil block")
		return nil
	}
	confirmationLog := log.Debug().Str("NetworkConfig", t.eth.NetworkConfig.ID).
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

// ReceiveBlock is a no-op
func (i *InstantConfirmations) ReceiveBlock(block NodeBlock) error {
	return nil
}

// Wait is a no-op
func (i *InstantConfirmations) Wait() error {
	return nil
}
