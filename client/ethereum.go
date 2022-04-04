package client

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/config"

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
	OneGWei = big.NewFloat(1e9)
	// OneEth represents 1 Ethereum
	OneEth = big.NewFloat(1e18)
)

// EthereumMultinodeClient wraps the client and the BlockChain network to interact with an EVM based Blockchain with multiple nodes
type EthereumMultinodeClient struct {
	DefaultClient *EthereumClient
	Clients       []*EthereumClient
}

func (e *EthereumMultinodeClient) ContractsDeployed() bool {
	return e.DefaultClient.ContractsDeployed()
}

func (e *EthereumMultinodeClient) EstimateTransactionGasCost() (*big.Int, error) {
	return e.DefaultClient.EstimateTransactionGasCost()
}

// EstimateCostForChainlinkOperations calculates TXs cost as a dirty estimation based on transactionLimit for that network
func (e *EthereumMultinodeClient) EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	return e.DefaultClient.EstimateCostForChainlinkOperations(amountOfOperations)
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

// DefaultWallet returns the default wallet for the network
func (e *EthereumMultinodeClient) GetDefaultWallet() *EthereumWallet {
	return e.DefaultClient.DefaultWallet
}

// GetNetworkName gets the ID of the chain that the clients are connected to
func (e *EthereumMultinodeClient) GetNetworkName() string {
	return e.DefaultClient.GetNetworkName()
}

// GetNetworkType retrieves the type of network this is running on
func (e *EthereumMultinodeClient) GetNetworkType() string {
	return e.DefaultClient.GetNetworkType()
}

// GetChainID retrieves the ChainID of the network that the client interacts with
func (e *EthereumMultinodeClient) GetChainID() int64 {
	return e.DefaultClient.GetChainID()
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

func (e *EthereumClient) ContractsDeployed() bool {
	return e.NetworkConfig.ContractsDeployed
}

// EstimateTransactionGasCost estimates the current total gas cost for a simple transaction
func (e *EthereumClient) EstimateTransactionGasCost() (*big.Int, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return gasPrice.Mul(gasPrice, big.NewInt(21000)), err
}

// EstimateCostForChainlinkOperations calculates required amount of ETH for amountOfOperations Chainlink operations
// based on the network's suggested gas price and the chainlink gas limit. This is fairly imperfect and should be used
// as only a rough, upper-end estimate instead of an exact calculation.
// See https://ethereum.org/en/developers/docs/gas/#post-london for info on how gas calculation works
func (e *EthereumClient) EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	bigAmountOfOperations := big.NewInt(int64(amountOfOperations))
	gasPriceInWei, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// https://ethereum.stackexchange.com/questions/19665/how-to-calculate-transaction-fee
	// total gas limit = chainlink gas limit + gas limit buffer
	gasLimit := e.NetworkConfig.GasEstimationBuffer + e.NetworkConfig.ChainlinkTransactionLimit
	// gas cost for TX = total gas limit * estimated gas price
	gasCostPerOperationWei := big.NewInt(1).Mul(big.NewInt(1).SetUint64(gasLimit), gasPriceInWei)
	gasCostPerOperationWeiFloat := big.NewFloat(1).SetInt(gasCostPerOperationWei)
	gasCostPerOperationETH := big.NewFloat(1).Quo(gasCostPerOperationWeiFloat, OneEth)
	// total Wei needed for all TXs = total value for TX * number of TXs
	totalWeiForAllOperations := big.NewInt(1).Mul(gasCostPerOperationWei, bigAmountOfOperations)
	totalWeiForAllOperationsFloat := big.NewFloat(1).SetInt(totalWeiForAllOperations)
	totalEthForAllOperations := big.NewFloat(1).Quo(totalWeiForAllOperationsFloat, OneEth)

	log.Debug().
		Int("Number of Operations", amountOfOperations).
		Uint64("Gas Limit per Operation", gasLimit).
		Str("Value per Operation (ETH)", gasCostPerOperationETH.String()).
		Str("Total (ETH)", totalEthForAllOperations.String()).
		Msg("Calculated ETH for Chainlink Operations")

	return totalEthForAllOperations, nil
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
func NewEthereumMultiNodeClient(
	_ string,
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

// SimulatedEthereumURLs returns the websocket URLs for a simulated geth network
func SimulatedEthereumURLs(e *environment.Environment) ([]*url.URL, error) {
	return e.Charts.Connections("geth").LocalURLsByPort("ws-rpc", environment.WS)
}

// SimulatedEthereumURLs returns the websocket URLs for a simulated geth network
func SimulatedSoakEthereumURLs(e *environment.Environment) ([]*url.URL, error) {
	return e.Charts.Connections("geth").RemoteURLsByPort("ws-rpc", environment.WS)
}

// LiveEthTestnetURLs indicates that there are no urls to fetch, except from the network config
func LiveEthTestnetURLs(e *environment.Environment) ([]*url.URL, error) {
	return []*url.URL{}, nil
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetDefaultWallet() *EthereumWallet {
	return e.DefaultWallet
}

// GetNetworkName retrieves the ID of the network that the client interacts with
func (e *EthereumClient) GetNetworkName() string {
	return e.NetworkConfig.ID
}

// GetNetworkType retrieves the type of network this is running on
func (e *EthereumClient) GetNetworkType() string {
	return e.NetworkConfig.Type
}

// GetChainID retrieves the ChainID of the network that the client interacts with
func (e *EthereumClient) GetChainID() int64 {
	return e.NetworkConfig.ChainID
}

// Close tears down the current open Ethereum client
func (e *EthereumClient) Close() error {
	e.doneChan <- struct{}{}
	e.Client.Close()
	return nil
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
		e.Nonces[addr.Hex()]++
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
		wei := big.NewFloat(1).Mul(OneEth, amount)
		log.Info().
			Str("Token", "ETH").
			Str("From", e.DefaultWallet.Address()).
			Str("To", toAddress).
			Str("Amount", amount.String()).
			Msg("Funding Address")
		_, err := e.SendTransaction(e.DefaultWallet, ethAddress, wei)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendTransaction sends a specified amount of ETH from a selected wallet to an address
func (e *EthereumClient) SendTransaction(
	from *EthereumWallet,
	to common.Address,
	value *big.Float,
) (common.Hash, error) {
	weiValue, _ := value.Int(nil)
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid private key: %v", err)
	}
	suggestedGasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return common.Hash{}, err
	}

	// TODO: Update from LegacyTx to DynamicFeeTx
	tx, err := types.SignNewTx(privateKey, types.NewEIP2930Signer(big.NewInt(e.NetworkConfig.ChainID)),
		&types.LegacyTx{
			To:       &to,
			Value:    weiValue,
			Data:     nil,
			Gas:      21000,
			GasPrice: suggestedGasPrice,
			Nonce:    nonce,
		})
	if err != nil {
		return common.Hash{}, err
	}
	if e.NetworkConfig.GasEstimationBuffer > 0 {
		log.Debug().
			Uint64("Suggested Gas Price Wei", big.NewInt(0).Sub(suggestedGasPrice, gasPriceBuffer).Uint64()).
			Uint64("Bumped Gas Price Wei", suggestedGasPrice.Uint64()).
			Str("TX Hash", tx.Hash().Hex()).
			Msg("Bumping Suggested Gas Price")
	}
	if err := e.Client.SendTransaction(context.Background(), tx); err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), e.ProcessTransaction(tx)
}

// ProcessTransaction will queue or wait on a transaction depending on whether parallel transactions are enabled
func (e *EthereumClient) ProcessTransaction(tx *types.Transaction) error {
	var txConfirmer HeaderEventSubscription
	if e.NetworkConfig.MinimumConfirmations == 0 {
		txConfirmer = &InstantConfirmations{}
	} else {
		txConfirmer = NewTransactionConfirmer(e, tx, e.NetworkConfig.MinimumConfirmations)
	}

	e.AddHeaderEventSubscription(tx.Hash().String(), txConfirmer)

	if !e.queueTransactions { // For sequential transactions
		log.Debug().Str("Hash", tx.Hash().String()).Msg("Waiting for TX to confirm before moving on")
		defer e.DeleteHeaderEventSubscription(tx.Hash().String())
		return txConfirmer.Wait()
	}
	return nil
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (e *EthereumClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := e.TransactionOpts(e.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	suggestedTipCap, err := e.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)
	opts.GasTipCap = suggestedTipCap.Add(gasPriceBuffer, suggestedTipCap)
	contractAddress, transaction, contractInstance, err := deployer(opts, e.Client)
	if err != nil {
		return nil, nil, nil, err
	}
	if e.NetworkConfig.GasEstimationBuffer > 0 {
		log.Debug().
			Uint64("Suggested Gas Tip Cap", big.NewInt(0).Sub(suggestedTipCap, gasPriceBuffer).Uint64()).
			Uint64("Bumped Gas Price", suggestedTipCap.Uint64()).
			Str("Contract Name", contractName).
			Msg("Bumping Suggested Gas Price")
	}
	if err := e.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}
	totalGasCostWeiFloat := big.NewFloat(1).SetInt(transaction.Cost())
	totalGasCostGwei := big.NewFloat(1).Quo(totalGasCostWeiFloat, OneGWei)

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", e.DefaultWallet.Address()).
		Str("Total Gas Cost (GWei)", totalGasCostGwei.String()).
		Str("Network", e.NetworkConfig.ID).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
// abigen will handle gas estimation for us on the backend.
func (e *EthereumClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(e.NetworkConfig.ChainID))
	if err != nil {
		return nil, err
	}
	opts.From = common.HexToAddress(from.Address())
	opts.Context = context.Background()

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return nil, err
	}
	opts.Nonce = big.NewInt(int64(nonce))

	return opts, nil
}

// WaitForEvents is a blocking function that waits for all event subscriptions that have been queued within the client.
func (e *EthereumClient) WaitForEvents() error {
	log.Debug().Msg("Waiting for blockchain events to finish before continuing")
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
				Err(err).
				Str("NetworkName", e.NetworkConfig.Name).
				Msg("Error while subscribing to headers")
			time.Sleep(time.Second)
			continue
		}
		break
	}
	log.Debug().Str("NetworkName", e.NetworkConfig.Name).Msg("Stopped subscribing to new headers")
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
		Str("NetworkName", e.NetworkConfig.Name).
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
		log.Err(fmt.Errorf("error on sending block to receivers: %v", err))
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
		if receipt.Status == 0 { // 0 indicates failure, 1 indicates success
			reason, err := e.errorReason(e.Client, tx, receipt)
			if err != nil {
				log.Warn().Str("TX Hash", txHash.Hex()).
					Str("To", tx.To().Hex()).
					Uint64("Nonce", tx.Nonce()).
					Msg("Transaction failed and was reverted! Unable to retrieve reason!")
				return false, err
			}
			log.Warn().Str("TX Hash", txHash.Hex()).
				Str("To", tx.To().Hex()).
				Str("Revert reason", reason).
				Msg("Transaction failed and was reverted!")
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
	tx               *types.Transaction
	doneChan         chan struct{}
	context          context.Context
	cancel           context.CancelFunc
}

// NewTransactionConfirmer returns a new instance of the transaction confirmer that waits for on-chain minimum
// confirmations
func NewTransactionConfirmer(eth *EthereumClient, tx *types.Transaction, minConfirmations int) *TransactionConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), eth.NetworkConfig.Timeout)
	tc := &TransactionConfirmer{
		minConfirmations: minConfirmations,
		confirmations:    0,
		eth:              eth,
		tx:               tx,
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
	confirmationLog := log.Debug().Str("Network", t.eth.NetworkConfig.ID).
		Str("Block Hash", block.Hash().Hex()).
		Str("Block Number", block.Number().String()).
		Str("Tx Hash", t.tx.Hash().String()).
		Uint64("Nonce", t.tx.Nonce()).
		Int("Minimum Confirmations", t.minConfirmations)
	isConfirmed, err := t.eth.isTxConfirmed(t.tx.Hash())
	if err != nil {
		return err
	} else if isConfirmed {
		t.confirmations++
	}
	if t.confirmations == t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).
			Msg("Transaction confirmations met")
		t.doneChan <- struct{}{}
	} else if t.confirmations <= t.minConfirmations {
		confirmationLog.Int("Current Confirmations", t.confirmations).
			Msg("Waiting on minimum confirmations")
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
			return fmt.Errorf("timeout waiting for transaction to confirm: %s", t.tx.Hash())
		}
	}
}

// InstantConfirmations is a no-op confirmer as all transactions are instantly mined so no confirmations are needed
type InstantConfirmations struct{}

// ReceiveBlock is a no-op
func (i *InstantConfirmations) ReceiveBlock(block NodeBlock) error {
	return nil
}

// Wait is a no-op
func (i *InstantConfirmations) Wait() error {
	return nil
}
