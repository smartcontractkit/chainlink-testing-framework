/*
This should be removed when we migrate all Ethereum client code to Seth
*/
package blockchain

// Contains implementations for multi and single node ethereum clients
import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
)

const (
	MaxTimeoutForFinality = 15 * time.Minute
	DefaultDialTimeout    = 1 * time.Minute
)

// EthereumClient wraps the client and the BlockChain network to interact with an EVM based Blockchain
type EthereumClient struct {
	ID                   int
	Client               *ethclient.Client
	rawRPC               *rpc.Client
	NetworkConfig        EVMNetwork
	Wallets              []*EthereumWallet
	DefaultWallet        *EthereumWallet
	NonceSettings        *NonceSettings
	FinalizedHeader      atomic.Pointer[FinalizedHeader]
	headerSubscriptions  map[string]HeaderEventSubscription
	subscriptionMutex    *sync.Mutex
	queueTransactions    bool
	gasStats             *GasStats
	connectionIssueCh    chan time.Time
	connectionRestoredCh chan time.Time
	doneChan             chan struct{}
	l                    zerolog.Logger
	subscriptionWg       sync.WaitGroup
}

// newEVMClient creates an EVM client for a single node/URL
func newEVMClient(networkSettings EVMNetwork, logger zerolog.Logger) (EVMClient, error) {
	logger.Info().
		Str("Name", networkSettings.Name).
		Str("URL", networkSettings.URL).
		Int64("Chain ID", networkSettings.ChainID).
		Bool("Simulated", networkSettings.Simulated).
		Bool("Supports EIP-1559", networkSettings.SupportsEIP1559).
		Bool("Finality Tag", networkSettings.FinalityTag).
		Msg("Connecting client")
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDialTimeout)
	defer cancel()
	raw, err := rpc.DialOptions(ctx, networkSettings.URL)
	if err != nil {
		return nil, err
	}
	cl := ethclient.NewClient(raw)

	ec := &EthereumClient{
		NetworkConfig:       networkSettings,
		Client:              cl,
		rawRPC:              raw,
		Wallets:             make([]*EthereumWallet, 0),
		headerSubscriptions: map[string]HeaderEventSubscription{},
		subscriptionMutex:   &sync.Mutex{},
		queueTransactions:   false,

		connectionIssueCh:    make(chan time.Time, 100), // buffered to prevent blocking, size is probably overkill, but some tests might not care
		connectionRestoredCh: make(chan time.Time, 100),
		doneChan:             make(chan struct{}),
		l:                    logger,
	}

	if ec.NetworkConfig.Simulated {
		ec.NonceSettings = newNonceSettings()
	} else { // un-simulated chain means potentially running tests in parallel, need to share nonces
		ec.NonceSettings = useGlobalNonceManager(ec.GetChainID())
	}

	if err := ec.LoadWallets(networkSettings); err != nil {
		return nil, err
	}
	ec.gasStats = NewGasStats(ec.ID)

	// Initialize header subscription or polling
	if err := ec.InitializeHeaderSubscription(); err != nil {
		return nil, err
	}
	// Check if the chain supports EIP-1559
	// https://eips.ethereum.org/EIPS/eip-1559
	if networkSettings.SupportsEIP1559 {
		ec.l.Debug().Msg("Network supports EIP-1559, using Dynamic transactions")
	} else {
		ec.l.Debug().Msg("Network does NOT support EIP-1559, using Legacy transactions")
	}

	return wrapSingleClient(networkSettings, ec), nil
}

// SyncNonce sets the NonceMu and Nonces value based on existing EVMClient
// it ensures the instance of EthereumClient is synced with passed EVMClient's nonce updates.
func (e *EthereumClient) SyncNonce(c EVMClient) {
	n := c.GetNonceSetting()
	n.NonceMu.Lock()
	defer n.NonceMu.Unlock()
	e.NonceSettings.NonceMu = n.NonceMu
	e.NonceSettings.Nonces = n.Nonces
}

// Get returns the underlying client type to be used generically across the framework for switching
// network types
func (e *EthereumClient) Get() interface{} {
	return e
}

// GetNetworkName retrieves the ID of the network that the client interacts with
func (e *EthereumClient) GetNetworkName() string {
	return e.NetworkConfig.Name
}

// NetworkSimulated returns true if the network is a simulated geth instance, false otherwise
func (e *EthereumClient) NetworkSimulated() bool {
	return e.NetworkConfig.Simulated
}

func (e *EthereumClient) GetNonceSetting() NonceSettings {
	e.NonceSettings.NonceMu.Lock()
	defer e.NonceSettings.NonceMu.Unlock()
	return NonceSettings{
		NonceMu: e.NonceSettings.NonceMu,
		Nonces:  e.NonceSettings.Nonces,
	}
}

// GetChainID retrieves the ChainID of the network that the client interacts with
func (e *EthereumClient) GetChainID() *big.Int {
	return big.NewInt(e.NetworkConfig.ChainID)
}

// GetClients not used, only applicable to EthereumMultinodeClient
func (e *EthereumClient) GetClients() []EVMClient {
	return []EVMClient{e}
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetDefaultWallet() *EthereumWallet {
	return e.DefaultWallet
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetWallets() []*EthereumWallet {
	return e.Wallets
}

// GetWalletByAddress returns the Ethereum wallet by address if it exists, else returns nil
func (e *EthereumClient) GetWalletByAddress(address common.Address) *EthereumWallet {
	for _, w := range e.Wallets {
		if w.address == address {
			return w
		}
	}
	return nil
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetNetworkConfig() *EVMNetwork {
	return &e.NetworkConfig
}

// SetID sets client id, only used for multi-node networks
func (e *EthereumClient) SetID(id int) {
	e.ID = id
}

// SetDefaultWallet sets default wallet
func (e *EthereumClient) SetDefaultWallet(num int) error {
	if num >= len(e.GetWallets()) {
		return fmt.Errorf("no wallet #%d found for default client", num)
	}
	e.DefaultWallet = e.Wallets[num]
	e.l.Debug().Str("Address", e.DefaultWallet.Address()).Int("Index", num).Msg("Set default wallet")
	return nil
}

// SetDefaultWalletByAddress sets default wallet by address if it exists, else returns error
func (e *EthereumClient) SetDefaultWalletByAddress(address common.Address) error {
	w := e.GetWalletByAddress(address)
	if w == nil {
		return fmt.Errorf("no wallet found for address %s", address.Hex())
	}
	e.DefaultWallet = w
	e.l.Debug().Int("Client", e.ID).Str("Address", e.DefaultWallet.Address()).Msg("Set default wallet")
	return nil
}

// SetWallets sets all wallets to be used by the client
func (e *EthereumClient) SetWallets(wallets []*EthereumWallet) {
	e.Wallets = wallets
}

// LoadWallets loads wallets from config
func (e *EthereumClient) LoadWallets(cfg EVMNetwork) error {
	pkStrings := cfg.PrivateKeys
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

// BalanceAt returns the ETH balance of the specified address
func (e *EthereumClient) BalanceAt(ctx context.Context, address common.Address) (*big.Int, error) {
	return e.Client.BalanceAt(ctx, address, nil)
}

// SwitchNode not used, only applicable to EthereumMultinodeClient
func (e *EthereumClient) SwitchNode(_ int) error {
	return nil
}

// HeaderHashByNumber gets header hash by block number
func (e *EthereumClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	h, err := e.HeaderByNumber(ctx, bn)
	if err != nil {
		return "", err
	}
	return h.Hash.String(), nil
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *EthereumClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	h, err := e.HeaderByNumber(ctx, bn)
	if err != nil {
		return 0, err
	}
	timestamp := h.Timestamp.UTC().Unix()
	if timestamp < 0 {
		return 0, fmt.Errorf("negative timestamp value: %d", timestamp)
	}
	return uint64(timestamp), nil
}

// BlockNumber gets latest block number
func (e *EthereumClient) LatestBlockNumber(ctx context.Context) (uint64, error) {
	bn, err := e.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

func (e *EthereumClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// for instant txs, wait for permission to go
	if e.NetworkConfig.MinimumConfirmations <= 0 {
		fromAddr, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return err
		}
		<-e.NonceSettings.registerInstantTransaction(fromAddr.Hex(), tx.Nonce())
	}
	return e.Client.SendTransaction(ctx, tx)
}

// Fund sends some ETH to an address using the default wallet
func (e *EthereumClient) Fund(
	toAddress string,
	amount *big.Float,
	gasEstimations GasEstimations,
) error {
	privateKey, err := crypto.HexToECDSA(e.DefaultWallet.PrivateKey())
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	to := common.HexToAddress(toAddress)

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(e.DefaultWallet.Address()))
	if err != nil {
		return err
	}

	tx, err := e.NewTx(privateKey, nonce, to, conversions.EtherToWei(amount), gasEstimations)
	if err != nil {
		return err
	}

	e.l.Info().
		Str("Token", "ETH").
		Str("From", e.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Hash", tx.Hash().Hex()).
		Uint64("Nonce", tx.Nonce()).
		Str("Network Name", e.GetNetworkName()).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", tx.Cost().Uint64()).
		Msg("Funding Address")
	if err := e.SendTransaction(context.Background(), tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = fmt.Errorf("using nonce %d err: %w", nonce, err)
		}
		return err
	}

	return e.ProcessTransaction(tx)
}

// ReturnFunds achieves a lazy method of fund return as too many guarantees get too complex
func (e *EthereumClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	fromAddress, err := conversions.PrivateKeyToAddress(fromKey)
	if err != nil {
		return err
	}
	nonce, err := e.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	balance, err := e.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err

	}
	gasLimit, err := e.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &e.DefaultWallet.address,
	})
	if err != nil {
		e.l.Warn().Int("Default", 21_000).Msg("Could not estimate gas for return funds transaction, using default")
		gasLimit = 21_000
	}
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	totalGasCost := new(big.Int).Mul(big.NewInt(0).SetUint64(gasLimit), gasPrice)
	toSend := new(big.Int).Sub(balance, totalGasCost)

	// Negative values can happen on large gas costs, and we can't send negative ETH
	if toSend.Cmp(big.NewInt(0)) <= 0 {
		e.l.Error().Str("From", fromAddress.Hex()).
			Uint64("Estimated Gas Cost", totalGasCost.Uint64()).
			Uint64("Balance", balance.Uint64()).
			Str("To Send", toSend.String()).
			Msg("Insufficient funds to return at current gas prices")
		return fmt.Errorf("insufficient funds to return at current gas prices, balance: %s", balance.String())
	}

	tx := types.NewTransaction(nonce, e.DefaultWallet.address, toSend, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(e.GetChainID()), fromKey)
	if err != nil {
		return err
	}

	// There are several errors that can happen when trying to drain an ETH wallet.
	// Ultimately, we want our money back, so we'll catch errors and react to them until we get something through,
	// or an error we don't recognize.

	// Try to send the money back, see if we hit any issues, and if we recognize them
	fundReturnErr := e.Client.SendTransaction(context.Background(), signedTx)

	// Handle overshot error
	overshotRe := regexp.MustCompile(`overshot (\d+)`)
	for fundReturnErr != nil && strings.Contains(fundReturnErr.Error(), "overshot") {
		submatches := overshotRe.FindStringSubmatch(fundReturnErr.Error())
		if len(submatches) < 1 {
			return fmt.Errorf("error parsing overshot amount in error: %w", err)
		}
		numberString := submatches[1]
		overshotAmount, err := strconv.Atoi(numberString)
		if err != nil {
			return err
		}
		toSend.Sub(toSend, big.NewInt(int64(overshotAmount)))
		tx := types.NewTransaction(nonce, e.DefaultWallet.address, toSend, gasLimit, gasPrice, nil)
		signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(e.GetChainID()), fromKey)
		if err != nil {
			return err
		}
		fundReturnErr = e.Client.SendTransaction(context.Background(), signedTx)
	}

	// Handle insufficient funds error
	// We don't get an overshot calculation, we just know it was too much, so subtract by 1 GWei and try again
	for fundReturnErr != nil && (strings.Contains(fundReturnErr.Error(), "insufficient funds") || strings.Contains(fundReturnErr.Error(), "gas too low")) {
		toSend.Sub(toSend, big.NewInt(GWei))
		if toSend.Cmp(big.NewInt(0)) <= 0 {
			e.l.Error().Str("From", fromAddress.Hex()).
				Uint64("Estimated Gas Cost", totalGasCost.Uint64()).
				Uint64("Balance", balance.Uint64()).
				Str("To Send", toSend.String()).
				Msg("Insufficient funds to return at current gas prices")
			return fmt.Errorf("insufficient funds to return at current gas prices, balance: %s : %w", balance.String(), fundReturnErr)
		}
		gasLimit += 21_000 // Add 21k gas for each attempt in case gas limit is too low, this has happened in weird L2 scenarios
		tx := types.NewTransaction(nonce, e.DefaultWallet.address, toSend, gasLimit, gasPrice, nil)
		signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(e.GetChainID()), fromKey)
		if err != nil {
			return err
		}
		fundReturnErr = e.Client.SendTransaction(context.Background(), signedTx)
	}

	e.l.Info().
		Uint64("Funds", toSend.Uint64()).
		Str("From", fromAddress.Hex()).
		Str("To", e.DefaultWallet.Address()).
		Msg("Returning funds to Default Wallet")
	return fundReturnErr
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
	gasCostPerOperationETH := conversions.WeiToEther(gasCostPerOperationWei)
	// total Wei needed for all TXs = total value for TX * number of TXs
	totalWeiForAllOperations := big.NewInt(1).Mul(gasCostPerOperationWei, bigAmountOfOperations)
	totalEthForAllOperations := conversions.WeiToEther(totalWeiForAllOperations)

	e.l.Debug().
		Int("Number of Operations", amountOfOperations).
		Uint64("Gas Limit per Operation", gasLimit).
		Str("Value per Operation (ETH)", gasCostPerOperationETH.String()).
		Str("Total (ETH)", totalEthForAllOperations.String()).
		Msg("Calculated ETH for Chainlink Operations")

	return totalEthForAllOperations, nil
}

func (e *EthereumClient) RawJsonRPCCall(ctx context.Context, result interface{}, method string, params ...interface{}) error {
	err := e.rawRPC.CallContext(ctx, &result, method, params...)

	return err
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
	if !e.NetworkConfig.SupportsEIP1559 {
		opts.GasPrice, err = e.EstimateGasPrice()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	contractAddress, transaction, contractInstance, err := deployer(opts, e.Client)
	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = fmt.Errorf("using nonce %d err: %w", opts.Nonce.Uint64(), err)
		}
		return nil, nil, nil, err
	}

	if err = e.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	e.l.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", e.DefaultWallet.Address()).
		Str("Total Gas Cost", conversions.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", e.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// LoadContract load already deployed contract instance
func (e *EthereumClient) LoadContract(contractName string, contractAddress common.Address, loader ContractLoader) (interface{}, error) {
	contractInstance, err := loader(contractAddress, e.Client)
	if err != nil {
		return nil, err
	}
	e.l.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("Network Name", e.NetworkConfig.Name).
		Msg("Loaded contract instance")
	return contractInstance, err
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (e *EthereumClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
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
	opts.Nonce = big.NewInt(0).SetUint64(nonce)

	if e.NetworkConfig.MinimumConfirmations <= 0 { // Wait for your turn to send on an L2 chain
		<-e.NonceSettings.registerInstantTransaction(from.Address(), nonce)
	}
	// if the gas limit is less than the default gas limit, use the default
	if e.NetworkConfig.DefaultGasLimit > opts.GasLimit {
		opts.GasLimit = e.NetworkConfig.DefaultGasLimit
	}
	if !e.NetworkConfig.SupportsEIP1559 {
		opts.GasPrice, err = e.EstimateGasPrice()
		if err != nil {
			return nil, err
		}
	}
	return opts, nil
}

func (e *EthereumClient) NewTx(
	fromPrivateKey *ecdsa.PrivateKey,
	nonce uint64,
	to common.Address,
	value *big.Int,
	gasEstimations GasEstimations,
) (*types.Transaction, error) {
	var (
		tx  *types.Transaction
		err error
	)
	if e.NetworkConfig.SupportsEIP1559 {
		tx, err = types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(e.GetChainID()), &types.DynamicFeeTx{
			ChainID:   e.GetChainID(),
			Nonce:     nonce,
			To:        &to,
			Value:     value,
			GasTipCap: gasEstimations.GasTipCap,
			GasFeeCap: gasEstimations.GasFeeCap,
			Gas:       gasEstimations.GasUnits,
		})
		if err != nil {
			return nil, err
		}
	} else {
		tx, err = types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(e.GetChainID()), &types.LegacyTx{
			Nonce:    nonce,
			To:       &to,
			Value:    value,
			GasPrice: gasEstimations.GasPrice,
			Gas:      gasEstimations.GasUnits,
		})
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}

// MarkTxAsSent On an L2 chain, indicate the tx has been sent
func (e *EthereumClient) MarkTxAsSentOnL2(tx *types.Transaction) error {
	if e.NetworkConfig.MinimumConfirmations > 0 {
		return nil
	}
	fromAddr, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return err
	}
	e.NonceSettings.sentInstantTransaction(fromAddr.Hex())
	return nil
}

// ProcessTransaction will queue or wait on a transaction depending on whether parallel transactions are enabled
func (e *EthereumClient) ProcessTransaction(tx *types.Transaction) error {
	e.l.Trace().Str("Hash", tx.Hash().Hex()).Msg("Processing Tx")
	var txConfirmer HeaderEventSubscription
	if e.GetNetworkConfig().MinimumConfirmations <= 0 {
		err := e.MarkTxAsSentOnL2(tx)
		if err != nil {
			return err
		}
		txConfirmer = NewInstantConfirmer(e, tx.Hash(), nil, nil, e.l)
	} else {
		txConfirmer = NewTransactionConfirmer(e, tx, e.GetNetworkConfig().MinimumConfirmations, e.l)
	}

	e.AddHeaderEventSubscription(tx.Hash().String(), txConfirmer)

	if !e.queueTransactions { // For sequential transactions
		e.l.Debug().Str("Hash", tx.Hash().String()).Msg("Waiting for TX to confirm before moving on")
		defer e.DeleteHeaderEventSubscription(tx.Hash().String())
		return txConfirmer.Wait()
	}
	return nil
}

func (e *EthereumClient) GetEthClient() *ethclient.Client {
	return e.Client
}

// ProcessEvent will queue or wait on an event depending on whether parallel transactions are enabled
func (e *EthereumClient) ProcessEvent(name string, event *types.Log, confirmedChan chan bool, errorChan chan error) error {
	var eventConfirmer HeaderEventSubscription
	if e.GetNetworkConfig().MinimumConfirmations <= 0 {
		eventConfirmer = NewInstantConfirmer(e, event.TxHash, confirmedChan, errorChan, e.l)
	} else {
		eventConfirmer = NewEventConfirmer(name, e, event, e.GetNetworkConfig().MinimumConfirmations, confirmedChan, errorChan)
	}

	subscriptionHash := fmt.Sprintf("%s-%s", event.TxHash.Hex(), name) // Many events can occupy the same tx hash
	e.AddHeaderEventSubscription(subscriptionHash, eventConfirmer)

	if !e.queueTransactions { // For sequential transactions
		e.l.Debug().Str("Hash", event.Address.Hex()).Msg("Waiting for Event to confirm before moving on")
		defer e.DeleteHeaderEventSubscription(subscriptionHash)
		return eventConfirmer.Wait()
	}
	return nil
}

// PollFinalizedHeader continuously polls the latest finalized header and stores it in the client
func (e *EthereumClient) PollFinality() error {
	if e.NetworkConfig.FinalityDepth > 0 {
		return fmt.Errorf("finality depth is greater than zero. no need to poll for finality")
	}
	f := newGlobalFinalizedHeaderManager(e)
	if f == nil {
		return fmt.Errorf("could not create finalized header manager")
	}
	e.FinalizedHeader.Store(f)
	e.AddHeaderEventSubscription(FinalizedHeaderKey, f)
	return nil
}

// CancelFinalityPolling stops polling for the latest finalized header
func (e *EthereumClient) CancelFinalityPolling() {
	if _, ok := e.headerSubscriptions[FinalizedHeaderKey]; ok {
		e.DeleteHeaderEventSubscription(FinalizedHeaderKey)
	}
}

// WaitForFinalizedTx waits for a transaction to be finalized
// If the network is simulated, it will return immediately
// otherwise it waits for the transaction to be finalized and returns the block number and time of the finalization
func (e *EthereumClient) WaitForFinalizedTx(txHash common.Hash) (*big.Int, time.Time, error) {
	receipt, err := e.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("error getting receipt: %w in network %s tx %s", err, e.GetNetworkName(), txHash.Hex())
	}
	txHdr, err := e.HeaderByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("error getting header: %w in network %s tx %s", err, e.GetNetworkName(), txHash.Hex())
	}
	finalizer := NewTransactionFinalizer(e, txHdr, receipt.TxHash)
	key := "txFinalizer-" + txHash.String()
	e.AddHeaderEventSubscription(key, finalizer)
	defer e.DeleteHeaderEventSubscription(key)

	if !e.Client.Client().SupportsSubscriptions() {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-finalizer.context.Done():
					return
				case <-ticker.C:
					latestHeader, err := e.GetLatestFinalizedBlockHeader(context.Background())
					if err != nil {
						e.l.Err(err).Msg("Error fetching latest finalized header via HTTP polling")
					}
					if latestHeader == nil {
						e.l.Error().Msg("Latest finalized header is nil")
						continue
					}
					if latestHeader.Time > math.MaxInt64 {
						e.l.Error().Msg("Latest finalized header time is too large")
						continue
					}

					nodeHeader := NodeHeader{
						// NodeID: 0, // Assign appropriate NodeID if needed
						SafeEVMHeader: SafeEVMHeader{
							Hash:      latestHeader.Hash(),
							Number:    latestHeader.Number,
							Timestamp: time.Unix(int64(latestHeader.Time), 0),
							BaseFee:   latestHeader.BaseFee,
						},
					}

					err = finalizer.ReceiveHeader(nodeHeader)
					if err != nil {
						e.l.Err(err).Msg("Finalizer received error during HTTP polling")
					}
				}
			}
		}()
	}
	err = finalizer.Wait()
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("error waiting for finalization: %w in network %s tx %s", err, e.GetNetworkName(), txHash.Hex())
	}
	return finalizer.FinalizedBy, finalizer.FinalizedAt, nil
}

// IsTxHeadFinalized checks if the transaction is finalized on chain
// in case of network with finality tag if the tx is not finalized it returns false,
// the latest finalized header number and the time at which it was finalized
// if the tx is finalized it returns true, the finalized header number by which the tx was considered finalized and the time at which it was finalized
func (e *EthereumClient) IsTxHeadFinalized(txHdr, header *SafeEVMHeader) (bool, *big.Int, time.Time, error) {
	if e.NetworkConfig.FinalityDepth > 0 {
		if header.Number.Cmp(new(big.Int).Add(txHdr.Number,
			big.NewInt(0).SetUint64(e.NetworkConfig.FinalityDepth))) > 0 {
			return true, header.Number, header.Timestamp, nil
		}
		return false, nil, time.Time{}, nil
	}
	fHead := e.FinalizedHeader.Load()
	if fHead != nil {
		latestFinalized := fHead.LatestFinalized.Load().(*big.Int)
		latestFinalizedAt := fHead.FinalizedAt.Load().(time.Time)
		if latestFinalized.Cmp(txHdr.Number) >= 0 {
			return true, latestFinalized, latestFinalizedAt, nil
		}
		return false, latestFinalized, latestFinalizedAt, nil
	}
	return false, nil, time.Time{}, fmt.Errorf("no finalized head found. start polling for finalized header, network %s", e.GetNetworkName())
}

// IsTxConfirmed checks if the transaction is confirmed on chain or not
func (e *EthereumClient) IsTxConfirmed(txHash common.Hash) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	tx, isPending, err := e.Client.TransactionByHash(ctx, txHash)
	cancel()
	if err != nil {
		if errors.Is(err, ethereum.NotFound) { // not found is fine, it's not on chain yet
			return false, nil
		}
		return !isPending, err
	}
	if !isPending && e.NetworkConfig.MinimumConfirmations > 0 { // Instant chains don't bother with this receipt nonsense
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		receipt, err := e.Client.TransactionReceipt(ctx, txHash)
		cancel()
		if err != nil {
			if errors.Is(err, ethereum.NotFound) { // not found is fine, it's not on chain yet
				return false, nil
			}
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
			to := "(none)"
			if tx.To() != nil {
				to = tx.To().Hex()
			}
			from := "(unknown)"
			fromAddr, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
			if err == nil {
				from = fromAddr.Hex()
			}
			reason, err := e.ErrorReason(e.Client, tx, receipt)
			if err != nil {
				e.l.Warn().Str("TX Hash", txHash.Hex()).
					Str("To", to).
					Str("From", from).
					Uint64("Nonce", tx.Nonce()).
					Str("Error extracting reason", err.Error()).
					Msg("Transaction failed and was reverted! Unable to retrieve reason!")
			} else {
				e.l.Warn().Str("TX Hash", txHash.Hex()).
					Str("Revert reason", reason).
					Str("To", to).
					Str("From", from).
					Msg("Transaction failed and was reverted!")
			}
			return false, fmt.Errorf("transaction failed and was reverted")
		}
	}
	return !isPending, err
}

// IsEventConfirmed returns if eth client can confirm that the event happened
func (e *EthereumClient) IsEventConfirmed(event *types.Log) (confirmed, removed bool, err error) {
	if event.Removed {
		return false, event.Removed, nil
	}
	eventTx, isPending, err := e.Client.TransactionByHash(context.Background(), event.TxHash)
	if err != nil {
		return false, event.Removed, err
	}
	if isPending {
		return false, event.Removed, nil
	}
	eventReceipt, err := e.Client.TransactionReceipt(context.Background(), eventTx.Hash())
	if err != nil {
		return false, event.Removed, err
	}
	if eventReceipt.Status == 0 { // Failed event tx
		reason, err := e.ErrorReason(e.Client, eventTx, eventReceipt)
		if err != nil {
			e.l.Warn().Str("TX Hash", eventTx.Hash().Hex()).
				Str("Error extracting reason", err.Error()).
				Msg("Transaction failed and was reverted! Unable to retrieve reason!")
		} else {
			e.l.Warn().Str("TX Hash", eventTx.Hash().Hex()).
				Str("Revert reason", reason).
				Msg("Transaction failed and was reverted!")
		}
		return false, event.Removed, err
	}
	headerByNumber, err := e.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(event.BlockNumber))
	if err != nil || headerByNumber == nil {
		return false, event.Removed, err
	}
	if headerByNumber.Hash != event.BlockHash {
		return false, event.Removed, nil
	}

	return true, event.Removed, nil
}

// GetTxReceipt returns the receipt of the transaction if available, error otherwise
func (e *EthereumClient) GetTxReceipt(txHash common.Hash) (*types.Receipt, error) {
	receipt, err := e.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// RevertReasonFromTx returns the revert reason for the transaction error by parsing through abi defined error list
func (e *EthereumClient) RevertReasonFromTx(txHash common.Hash, abiString string) (string, interface{}, error) {
	tx, _, err := e.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return "", nil, err
	}
	re, err := e.GetTxReceipt(txHash)
	if err != nil {
		return "", nil, err
	}
	errData, err := e.ErrorReason(e.Client, tx, re)
	if err != nil {
		return "", nil, err
	}
	if errData == "" {
		return "", nil, fmt.Errorf("no revert reason found")
	}
	data, err := hex.DecodeString(errData[2:])
	if err != nil {
		return "", nil, err
	}
	jsonABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return "", nil, err
	}
	for errName, abiError := range jsonABI.Errors {
		if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
			// Found a matching error
			v, err := abiError.Unpack(data)
			if err != nil {
				return "", nil, err
			}
			return errName, v, nil
		}
	}
	return "", nil, fmt.Errorf("revert Reason could not be found for given abistring")
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
func (e *EthereumClient) ParallelTransactions(enabled bool) {
	e.queueTransactions = enabled
}

// Close tears down the current open Ethereum client and unsubscribes from the manager, if any
func (e *EthereumClient) Close() error {
	// close(e.NonceSettings.doneChan)
	close(e.doneChan)
	e.subscriptionWg.Wait()

	chainManagerRegistry.Lock()
	defer chainManagerRegistry.Unlock()

	mgr, exists := chainManagerRegistry.managers[e.GetChainID().Int64()]
	if exists {
		mgr.unsubscribe(e)
		// If no more subscribers remain, we can shut down the manager
		if len(mgr.subscribers) == 0 {
			mgr.shutdown()
			removeChainManager(e.GetChainID().Int64())
		}
	}
	return nil
}

// EstimateTransactionGasCost estimates the current total gas cost for a simple transaction
func (e *EthereumClient) EstimateTransactionGasCost() (*big.Int, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return gasPrice.Mul(gasPrice, big.NewInt(21000)), err
}

// GasStats retrieves all information on gas spent by this client
func (e *EthereumClient) GasStats() *GasStats {
	return e.gasStats
}

// EstimateGas estimates all gas values based on call data
func (e *EthereumClient) EstimateGas(callMsg ethereum.CallMsg) (GasEstimations, error) {
	var (
		gasUnits  uint64
		gasTipCap *big.Int
		gasFeeCap *big.Int
		err       error
	)
	if callMsg.To == nil && callMsg.Data == nil {
		return GasEstimations{}, fmt.Errorf("you've called EstimateGas with a nil To address, this can cause weird errors and inaccuracies, provide a To address in your callMsg, or use EstimateGasPrice for just a price")
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.NetworkConfig.Timeout.Duration)
	// Gas Units
	gasUnits, err = e.Client.EstimateGas(ctx, callMsg)
	cancel()
	if err != nil {
		return GasEstimations{}, err
	}

	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)
	// Legacy Gas Price
	ctx, cancel = context.WithTimeout(context.Background(), e.NetworkConfig.Timeout.Duration)
	gasPrice, err := e.Client.SuggestGasPrice(ctx)
	cancel()
	if err != nil {
		return GasEstimations{}, err
	}
	gasPrice.Add(gasPrice, gasPriceBuffer)

	if e.NetworkConfig.SupportsEIP1559 {
		// GasTipCap
		ctx, cancel := context.WithTimeout(context.Background(), e.NetworkConfig.Timeout.Duration)
		gasTipCap, err = e.Client.SuggestGasTipCap(ctx)
		cancel()
		if err != nil {
			return GasEstimations{}, err
		}
		gasTipCap.Add(gasTipCap, gasPriceBuffer)

		// GasFeeCap
		ctx, cancel = context.WithTimeout(context.Background(), e.NetworkConfig.Timeout.Duration)
		latestHeader, err := e.HeaderByNumber(ctx, nil)
		cancel()
		if err != nil {
			return GasEstimations{}, err
		}
		baseFeeMult := big.NewInt(1).Mul(latestHeader.BaseFee, big.NewInt(2))
		gasFeeCap = baseFeeMult.Add(baseFeeMult, gasTipCap)
	} else {
		gasFeeCap = gasPrice
		gasTipCap = gasPrice
	}

	// Total Gas Cost
	totalGasCost := big.NewInt(0).Mul(gasFeeCap, new(big.Int).SetUint64(gasUnits))

	return GasEstimations{
		GasUnits:     gasUnits,
		GasPrice:     gasPrice,
		GasFeeCap:    gasFeeCap,
		GasTipCap:    gasTipCap,
		TotalGasCost: totalGasCost,
	}, nil
}

func (e *EthereumClient) EstimateGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.NetworkConfig.Timeout.Duration)
	defer cancel()
	return e.Client.SuggestGasPrice(ctx)
}

// ConnectionIssue returns a channel that will receive a timestamp when the connection is lost
func (e *EthereumClient) ConnectionIssue() chan time.Time {
	return e.connectionIssueCh
}

// ConnectionRestored returns a channel that will receive a timestamp when the connection is restored
func (e *EthereumClient) ConnectionRestored() chan time.Time {
	return e.connectionRestoredCh
}

func (e *EthereumClient) Backend() bind.ContractBackend {
	return e.Client
}

func (e *EthereumClient) DeployBackend() bind.DeployBackend {
	return e.Client
}

func (e *EthereumClient) SubscribeNewHeaders(
	ctx context.Context,
	headerChan chan *SafeEVMHeader,
) (ethereum.Subscription, error) {
	clientSub, err := e.rawRPC.EthSubscribe(ctx, headerChan, "newHeads")
	if err != nil {
		return nil, err
	}

	return clientSub, err
}

// HeaderByNumber retrieves a Safe EVM header by number, nil for latest
func (e *EthereumClient) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*SafeEVMHeader, error) {
	var head *SafeEVMHeader
	err := e.rawRPC.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

// HeaderByHash retrieves a Safe EVM header by hash
func (e *EthereumClient) HeaderByHash(ctx context.Context, hash common.Hash) (*SafeEVMHeader, error) {
	var head *SafeEVMHeader
	err := e.rawRPC.CallContext(ctx, &head, "eth_getBlockByHash", hash, false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

// toBlockNumArg translates a block number to the correct argument for the RPC call
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *EthereumClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	e.subscriptionMutex.Lock()
	defer e.subscriptionMutex.Unlock()
	e.headerSubscriptions[key] = subscriber
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *EthereumClient) DeleteHeaderEventSubscription(key string) {
	e.subscriptionMutex.Lock()
	defer e.subscriptionMutex.Unlock()
	delete(e.headerSubscriptions, key)
}

// WaitForEvents is a blocking function that waits for all event subscriptions that have been queued within the client.
func (e *EthereumClient) WaitForEvents() error {
	e.l.Debug().Msg("Waiting for blockchain events to finish before continuing")
	queuedEvents := e.GetHeaderSubscriptions()

	g := errgroup.Group{}

	for subName, sub := range queuedEvents {
		g.Go(func() error {
			defer func() {
				// if the subscription is complete, delete it from the queue
				if sub.Complete() {
					e.DeleteHeaderEventSubscription(subName)
				}
			}()
			return sub.Wait()
		})
	}

	return g.Wait()
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (e *EthereumClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return e.Client.SubscribeFilterLogs(ctx, q, ch)
}

// FilterLogs executes a filter query
func (e *EthereumClient) FilterLogs(ctx context.Context, filterQuery ethereum.FilterQuery) ([]types.Log, error) {
	return e.Client.FilterLogs(ctx, filterQuery)
}

// GetLatestFinalizedBlockHeader returns the latest finalized block header
// if finality tag is enabled, it returns the latest finalized block header
// otherwise it returns the block header for the block obtained by latest block number - finality depth
func (e *EthereumClient) GetLatestFinalizedBlockHeader(ctx context.Context) (*types.Header, error) {
	if e.NetworkConfig.FinalityTag {
		return e.Client.HeaderByNumber(ctx, big.NewInt(rpc.FinalizedBlockNumber.Int64()))
	}
	if e.NetworkConfig.FinalityDepth == 0 {
		return nil, fmt.Errorf("finality depth is 0 and finality tag is not enabled")
	}
	header, err := e.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	latestBlockNumber := header.Number.Uint64()
	finalizedBlockNumber := latestBlockNumber - e.NetworkConfig.FinalityDepth
	if finalizedBlockNumber > math.MaxInt64 {
		return nil, fmt.Errorf("finalized block number %d exceeds int64 range", finalizedBlockNumber)
	}
	return e.Client.HeaderByNumber(ctx, big.NewInt(int64(finalizedBlockNumber)))
}

// EstimatedFinalizationTime returns the estimated time it takes for a block to be finalized
// for networks with finality tag enabled, it returns the time between the current and next finalized block
// for networks with finality depth enabled, it returns the time to mine blocks equal to finality depth
func (e *EthereumClient) EstimatedFinalizationTime(ctx context.Context) (time.Duration, error) {
	if e.NetworkConfig.TimeToReachFinality.Duration != 0 {
		return e.NetworkConfig.TimeToReachFinality.Duration, nil
	}
	e.l.Info().Msg("TimeToReachFinality is not provided. Calculating estimated finalization time")
	if e.NetworkConfig.FinalityTag {
		return e.TimeBetweenFinalizedBlocks(ctx, MaxTimeoutForFinality)
	}
	blckTime, err := e.AvgBlockTime(ctx)
	if err != nil {
		return 0, err
	}
	if e.NetworkConfig.FinalityDepth == 0 {
		return 0, fmt.Errorf("finality depth is 0 and finality tag is not enabled")
	}
	timeBetween := time.Duration(e.NetworkConfig.FinalityDepth) * blckTime //nolint
	e.l.Info().
		Str("Time", timeBetween.String()).
		Str("Network", e.GetNetworkName()).
		Msg("Estimated finalization time")
	return timeBetween, nil

}

// TimeBetweenFinalizedBlocks is used to calculate the time between finalized blocks for chains with finality tag enabled
func (e *EthereumClient) TimeBetweenFinalizedBlocks(ctx context.Context, maxTimeToWait time.Duration) (time.Duration, error) {
	if !e.NetworkConfig.FinalityTag {
		return 0, fmt.Errorf("finality tag is not enabled; cannot calculate time between finalized blocks")
	}
	currentFinalizedHeader, err := e.GetLatestFinalizedBlockHeader(ctx)
	if err != nil {
		return 0, err
	}
	hdrChannel := make(chan *types.Header)
	var sub ethereum.Subscription
	sub, err = e.Client.SubscribeNewHead(ctx, hdrChannel)
	if err != nil {
		return 0, err
	}
	defer sub.Unsubscribe()
	c, cancel := context.WithTimeout(ctx, maxTimeToWait)
	defer cancel()
	for {
		select {
		case <-c.Done():
			return 0, fmt.Errorf("timed out waiting for next finalized block. If the finality time is more than %s, provide it as TimeToReachFinality in Network config: %w", maxTimeToWait, c.Err())
		case <-hdrChannel:
			// a new header is received now query the finalized block
			nextFinalizedHeader, err := e.GetLatestFinalizedBlockHeader(ctx)
			if err != nil {
				return 0, err
			}
			if nextFinalizedHeader.Number.Cmp(currentFinalizedHeader.Number) > 0 {
				timeBetween := time.Unix(int64(nextFinalizedHeader.Time), 0).Sub(time.Unix(int64(currentFinalizedHeader.Time), 0)) //nolint
				e.l.Info().
					Str("Time", timeBetween.String()).
					Str("Network", e.GetNetworkName()).
					Msg("Time between finalized blocks")
				return timeBetween, nil
			}
		case subErr := <-sub.Err():
			return 0, subErr
		}
	}
}

// AvgBlockTime calculates the average block time over the last 100 blocks for non-simulated networks
// and the last 10 blocks for simulated networks.
func (e *EthereumClient) AvgBlockTime(ctx context.Context) (time.Duration, error) {
	header, err := e.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	latestBlockNumber := header.Number.Uint64()
	numBlocks := uint64(100) // Number of blocks to consider for calculating block time
	if e.NetworkSimulated() {
		numBlocks = uint64(10)
	}
	startBlockNumber := latestBlockNumber - numBlocks + 1
	if startBlockNumber <= 0 {
		return 0, fmt.Errorf("not enough blocks mined to calculate block time")
	}
	totalTime := time.Duration(0)
	var previousHeader *types.Header
	previousHeader, err = e.Client.HeaderByNumber(ctx, big.NewInt(int64(startBlockNumber-1))) //nolint
	if err != nil {
		return totalTime, err
	}
	for i := startBlockNumber; i <= latestBlockNumber; i++ {
		hdr, err := e.Client.HeaderByNumber(ctx, big.NewInt(int64(i))) //nolint
		if err != nil {
			return totalTime, err
		}

		blockTime := time.Unix(int64(hdr.Time), 0)                    //nolint
		previousBlockTime := time.Unix(int64(previousHeader.Time), 0) //nolint
		blockDuration := blockTime.Sub(previousBlockTime)
		totalTime += blockDuration
		previousHeader = hdr
	}

	averageBlockTime := totalTime / time.Duration(numBlocks) //nolint

	return averageBlockTime, nil
}

// InitializeHeaderSubscription initializes either subscription-based or polling-based header processing
func (e *EthereumClient) InitializeHeaderSubscription() error {
	if e.Client.Client().SupportsSubscriptions() {
		return e.subscribeToNewHeaders()
	}
	// Fallback to polling if subscriptions are not supported
	e.l.Info().Str("Network", e.NetworkConfig.Name).Msg("Subscriptions not supported. Using polling for new headers.")

	// Acquire (or create) a manager for this chain
	mgr := getOrCreateChainManager(
		e.GetChainID().Int64(),
		15*time.Second, // or e.NetworkConfig.PollInterval, etc.
		e.NetworkConfig,
		e.l,
		e.Client,
		e.rawRPC,
	)

	// Subscribe
	mgr.subscribe(e)
	// Start polling if not started
	mgr.startPolling()

	return nil
}

// EthereumMultinodeClient wraps the client and the BlockChain network to interact with an EVM based Blockchain with multiple nodes
type EthereumMultinodeClient struct {
	DefaultClient EVMClient
	Clients       []EVMClient
}

func (e *EthereumMultinodeClient) GetEthClient() *ethclient.Client {
	return e.DefaultClient.GetEthClient()
}

func (e *EthereumMultinodeClient) Backend() bind.ContractBackend {
	return e.DefaultClient.Backend()
}

func (e *EthereumMultinodeClient) DeployBackend() bind.DeployBackend {
	return e.DefaultClient.DeployBackend()
}

func (e *EthereumMultinodeClient) SubscribeNewHeaders(
	ctx context.Context,
	headerChan chan *SafeEVMHeader,
) (ethereum.Subscription, error) {
	return e.DefaultClient.SubscribeNewHeaders(ctx, headerChan)
}

// HeaderByNumber retrieves a Safe EVM header by number, nil for latest
func (e *EthereumMultinodeClient) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*SafeEVMHeader, error) {
	return e.DefaultClient.HeaderByNumber(ctx, number)
}

// HeaderByHash retrieves a Safe EVM header by hash
func (e *EthereumMultinodeClient) HeaderByHash(ctx context.Context, hash common.Hash) (*SafeEVMHeader, error) {
	return e.DefaultClient.HeaderByHash(ctx, hash)
}

func (e *EthereumMultinodeClient) RawJsonRPCCall(ctx context.Context, result interface{}, method string, params ...interface{}) error {
	err := e.DefaultClient.RawJsonRPCCall(ctx, result, method, params)

	return err
}

// LoadContract load already deployed contract instance
func (e *EthereumMultinodeClient) LoadContract(contractName string, address common.Address, loader ContractLoader) (interface{}, error) {
	return e.DefaultClient.LoadContract(contractName, address, loader)
}

// EstimateCostForChainlinkOperations calculates TXs cost as a dirty estimation based on transactionLimit for that network
func (e *EthereumMultinodeClient) EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	return e.DefaultClient.EstimateCostForChainlinkOperations(amountOfOperations)
}

func (e *EthereumMultinodeClient) GetHeaderSubscriptions() map[string]HeaderEventSubscription {
	return e.DefaultClient.GetHeaderSubscriptions()
}

// NewEVMClientFromNetwork returns a multi-node EVM client connected to the specified network
func NewEVMClientFromNetwork(networkSettings EVMNetwork, logger zerolog.Logger) (EVMClient, error) {
	ecl := &EthereumMultinodeClient{}
	for idx, networkURL := range networkSettings.URLs {
		networkSettings.URL = networkURL
		ec, err := newEVMClient(networkSettings, logger)
		if err != nil {
			return nil, err
		}
		// a call to BalanceAt (can be any on chain call) to ensure the client is connected
		_, err = ec.BalanceAt(context.Background(), ec.GetDefaultWallet().address)
		if err == nil {
			ec.SetID(idx)
			ecl.Clients = append(ecl.Clients, ec)
			break
		}
	}
	if len(ecl.Clients) == 0 {
		return nil, fmt.Errorf("no networks or URLS found for network %s", networkSettings.Name)
	}
	ecl.DefaultClient = ecl.Clients[0]
	wrappedClient := wrapMultiClient(networkSettings, ecl)
	// required in Geth when you need to call "simulate" transactions from nodes
	if ecl.NetworkSimulated() {
		zero := common.HexToAddress("0x0")
		gasEstimations, err := wrappedClient.EstimateGas(ethereum.CallMsg{
			To: &zero,
		})
		if err != nil {
			return nil, err
		}
		if err := ecl.Fund("0x0", big.NewFloat(1000), gasEstimations); err != nil {
			return nil, err
		}
	}
	return wrappedClient, nil
}

// NewEVMClient returns a multi-node EVM client connected to the specified network
// Note: This should mostly be deprecated in favor of ConnectEVMClient. This is really only used when needing to connect
// to simulated networks
func NewEVMClient(networkSettings EVMNetwork, env *environment.Environment, logger zerolog.Logger) (EVMClient, error) {
	if env == nil {
		return nil, fmt.Errorf("environment nil, use ConnectEVMClient or provide a non-nil environment")
	}

	if networkSettings.Simulated {
		if _, ok := env.URLs[networkSettings.Name]; !ok {
			return nil, fmt.Errorf("network %s not found in environment", networkSettings.Name)
		}
		networkSettings.URLs = env.URLs[networkSettings.Name]
		networkSettings.HTTPURLs = env.URLs[networkSettings.Name+"_http"]
	}

	return ConnectEVMClient(networkSettings, logger)
}

// ConnectEVMClient returns a multi-node EVM client connected to a specified network, using only URLs.
// Should mostly be used for inside K8s, non-simulated tests.
func ConnectEVMClient(networkSettings EVMNetwork, logger zerolog.Logger) (EVMClient, error) {
	var urls []string
	if len(networkSettings.URLs) == 0 {
		if len(networkSettings.HTTPURLs) == 0 {
			return nil, fmt.Errorf("no URLs provided to connect to network")
		}
		logger.Warn().Msg("You are running using only HTTP RPC URLs.")
		urls = networkSettings.HTTPURLs
	} else {
		urls = networkSettings.URLs
	}

	ecl := &EthereumMultinodeClient{}

	for idx, networkURL := range urls {
		networkSettings.URL = networkURL
		ec, err := newEVMClient(networkSettings, logger)

		if err != nil {
			logger.Info().
				Err(err).
				Str("URL Suffix", networkURL[len(networkURL)-6:]).
				Msg("failed to create new EVM client")
			continue
		}
		// a call to BalanceAt to ensure the client is connected
		b, err := ec.BalanceAt(context.Background(), ec.GetDefaultWallet().address)
		if err == nil {
			ec.SetID(idx)
			ecl.Clients = append(ecl.Clients, ec)

			logger.Info().
				Uint64("Balance", b.Uint64()).
				Str("Address", ec.GetDefaultWallet().address.Hex()).
				Msg("Default address balance")

			if networkSettings.Simulated && b.Cmp(big.NewInt(0)) == 0 {
				noBalanceErr := fmt.Errorf("default wallet %s has no balance", ec.GetDefaultWallet().address.Hex())
				logger.Err(noBalanceErr).
					Msg("Ending test before it fails anyway")

				return nil, noBalanceErr
			}

			break
		}
	}
	if len(ecl.Clients) == 0 {
		return nil, fmt.Errorf("failed to create new EVM client")
	}
	ecl.DefaultClient = ecl.Clients[0]
	wrappedClient := wrapMultiClient(networkSettings, ecl)
	// required in Geth when you need to call "simulate" transactions from nodes
	if ecl.NetworkSimulated() {
		zero := common.HexToAddress("0x0")
		gasEstimations, err := wrappedClient.EstimateGas(ethereum.CallMsg{
			To: &zero,
		})
		if err != nil {
			return nil, err
		}
		if err := ecl.Fund("0x0", big.NewFloat(1000), gasEstimations); err != nil {
			return nil, err
		}
	}

	return wrappedClient, nil
}

// ConcurrentEVMClient returns a multi-node EVM client connected to a specified network
// It is used for concurrent interactions from different threads with the same network and from same owner
// account. This ensures that correct nonce value is fetched when an instance of EVMClient is initiated using this method.
// This is mainly useful for simulated networks as we don't use global nonce manager for them.
func ConcurrentEVMClient(networkSettings EVMNetwork, env *environment.Environment, existing EVMClient, logger zerolog.Logger) (EVMClient, error) {
	// if not simulated use the NewEVMClient
	if !networkSettings.Simulated {
		return ConnectEVMClient(networkSettings, logger)
	}
	ecl := &EthereumMultinodeClient{}
	if env != nil {
		if _, ok := env.URLs[existing.GetNetworkConfig().Name]; !ok {
			return nil, fmt.Errorf("network %s not found in environment", existing.GetNetworkConfig().Name)
		}
		networkSettings.URLs = env.URLs[existing.GetNetworkConfig().Name]
	}
	for idx, networkURL := range networkSettings.URLs {
		networkSettings.URL = networkURL
		ec, err := newEVMClient(networkSettings, logger)
		if err != nil {
			logger.Info().
				Err(err).
				Str("URL Suffix", networkURL[len(networkURL)-6:]).
				Msg("failed to create new EVM client")
			continue
		}
		// a call to BalanceAt (can be any on chain call) to ensure the client is connected
		_, err = ec.BalanceAt(context.Background(), ec.GetDefaultWallet().address)
		if err == nil {
			ec.SyncNonce(existing)
			ec.SetID(idx)
			ecl.Clients = append(ecl.Clients, ec)
			break
		}
	}
	if len(ecl.Clients) == 0 {
		return nil, fmt.Errorf("failed to create new EVM client")
	}
	ecl.DefaultClient = ecl.Clients[0]
	wrappedClient := wrapMultiClient(networkSettings, ecl)
	ecl.SetWallets(existing.GetWallets())
	if err := ecl.SetDefaultWalletByAddress(existing.GetDefaultWallet().address); err != nil {
		return nil, err
	}
	// no need to fund the account as it is already funded in the existing client
	return wrappedClient, nil
}

// SetDefaultWalletByAddress sets default wallet by address if it exists, else returns error
func (e *EthereumMultinodeClient) SetDefaultWalletByAddress(address common.Address) error {
	return e.DefaultClient.SetDefaultWalletByAddress(address)
}

// GetWalletByAddress returns the Ethereum wallet by address if it exists, else returns nil
func (e *EthereumMultinodeClient) GetWalletByAddress(address common.Address) *EthereumWallet {
	return e.DefaultClient.GetWalletByAddress(address)
}

// Get gets default client as an interface{}
func (e *EthereumMultinodeClient) Get() interface{} {
	return e.DefaultClient
}
func (e *EthereumMultinodeClient) GetNonceSetting() NonceSettings {
	return e.DefaultClient.GetNonceSetting()
}

func (e *EthereumMultinodeClient) SyncNonce(c EVMClient) {
	e.DefaultClient.SyncNonce(c)
}

// GetNetworkName gets the ID of the chain that the clients are connected to
func (e *EthereumMultinodeClient) GetNetworkName() string {
	return e.DefaultClient.GetNetworkName()
}

// GetNetworkType retrieves the type of network this is running on
func (e *EthereumMultinodeClient) NetworkSimulated() bool {
	return e.DefaultClient.NetworkSimulated()
}

// GetChainID retrieves the ChainID of the network that the client interacts with
func (e *EthereumMultinodeClient) GetChainID() *big.Int {
	return e.DefaultClient.GetChainID()
}

// GetClients gets clients for all nodes connected
func (e *EthereumMultinodeClient) GetClients() []EVMClient {
	cl := make([]EVMClient, 0)
	cl = append(cl, e.Clients...)
	return cl
}

// GetDefaultWallet returns the default wallet for the network
func (e *EthereumMultinodeClient) GetDefaultWallet() *EthereumWallet {
	return e.DefaultClient.GetDefaultWallet()
}

// GetWallets returns the default wallet for the network
func (e *EthereumMultinodeClient) GetWallets() []*EthereumWallet {
	return e.DefaultClient.GetWallets()
}

// GetNetworkConfig return the network config
func (e *EthereumMultinodeClient) GetNetworkConfig() *EVMNetwork {
	return e.DefaultClient.GetNetworkConfig()
}

// SetID sets client ID in a multi-node environment
func (e *EthereumMultinodeClient) SetID(id int) {
	e.DefaultClient.SetID(id)
}

// SetDefaultWallet sets default wallet
func (e *EthereumMultinodeClient) SetDefaultWallet(num int) error {
	return e.DefaultClient.SetDefaultWallet(num)
}

// SetWallets sets the default client's wallets
func (e *EthereumMultinodeClient) SetWallets(wallets []*EthereumWallet) {
	e.DefaultClient.SetWallets(wallets)
}

// LoadWallets loads wallets using private keys provided in the config
func (e *EthereumMultinodeClient) LoadWallets(cfg EVMNetwork) error {
	pkStrings := cfg.PrivateKeys
	wallets := make([]*EthereumWallet, 0)
	for _, pks := range pkStrings {
		w, err := NewEthereumWallet(pks)
		if err != nil {
			return err
		}
		wallets = append(wallets, w)
	}
	for _, c := range e.Clients {
		c.SetWallets(wallets)
	}
	return nil
}

// BalanceAt returns the ETH balance of the specified address
func (e *EthereumMultinodeClient) BalanceAt(ctx context.Context, address common.Address) (*big.Int, error) {
	return e.DefaultClient.BalanceAt(ctx, address)
}

// SwitchNode sets default client to perform calls to the network
func (e *EthereumMultinodeClient) SwitchNode(clientID int) error {
	if clientID > len(e.Clients) {
		return fmt.Errorf("client for node %d not found", clientID)
	}
	e.DefaultClient = e.Clients[clientID]
	return nil
}

// HeaderHashByNumber gets header hash by block number
func (e *EthereumMultinodeClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	return e.DefaultClient.HeaderHashByNumber(ctx, bn)
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *EthereumMultinodeClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	return e.DefaultClient.HeaderTimestampByNumber(ctx, bn)
}

// LatestBlockNumber gets the latest block number from the default client
func (e *EthereumMultinodeClient) LatestBlockNumber(ctx context.Context) (uint64, error) {
	return e.DefaultClient.LatestBlockNumber(ctx)
}

// SendTransaction wraps ethereum's SendTransaction to make it safe with instant transaction types
func (e *EthereumMultinodeClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return e.DefaultClient.SendTransaction(ctx, tx)
}

// Fund funds a specified address with ETH from the given wallet
func (e *EthereumMultinodeClient) Fund(toAddress string, nativeAmount *big.Float, gasEstimations GasEstimations) error {
	return e.DefaultClient.Fund(toAddress, nativeAmount, gasEstimations)
}

func (e *EthereumMultinodeClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	return e.DefaultClient.ReturnFunds(fromKey)
}

// DeployContract deploys a specified contract
func (e *EthereumMultinodeClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	return e.DefaultClient.DeployContract(contractName, deployer)
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (e *EthereumMultinodeClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	return e.DefaultClient.TransactionOpts(from)
}

func (e *EthereumMultinodeClient) NewTx(
	fromPrivateKey *ecdsa.PrivateKey,
	nonce uint64,
	to common.Address,
	value *big.Int,
	gasEstimations GasEstimations,
) (*types.Transaction, error) {
	return e.DefaultClient.NewTx(fromPrivateKey, nonce, to, value, gasEstimations)
}

// ProcessTransaction returns the result of the default client's processed transaction
func (e *EthereumMultinodeClient) ProcessTransaction(tx *types.Transaction) error {
	return e.DefaultClient.ProcessTransaction(tx)
}

// ProcessEvent returns the result of the default client's processed event
func (e *EthereumMultinodeClient) ProcessEvent(name string, event *types.Log, confirmedChan chan bool, errorChan chan error) error {
	return e.DefaultClient.ProcessEvent(name, event, confirmedChan, errorChan)
}

// IsTxConfirmed returns the default client's transaction confirmations
func (e *EthereumMultinodeClient) IsTxConfirmed(txHash common.Hash) (bool, error) {
	return e.DefaultClient.IsTxConfirmed(txHash)
}

// IsEventConfirmed returns if the default client can confirm the event has happened
func (e *EthereumMultinodeClient) IsEventConfirmed(event *types.Log) (confirmed, removed bool, err error) {
	return e.DefaultClient.IsEventConfirmed(event)
}

// IsTxFinalized returns if the default client can confirm the transaction has been finalized
func (e *EthereumMultinodeClient) IsTxHeadFinalized(txHdr, header *SafeEVMHeader) (bool, *big.Int, time.Time, error) {
	return e.DefaultClient.IsTxHeadFinalized(txHdr, header)
}

// WaitForTxTobeFinalized waits for the transaction to be finalized
func (e *EthereumMultinodeClient) WaitForFinalizedTx(txHash common.Hash) (*big.Int, time.Time, error) {
	return e.DefaultClient.WaitForFinalizedTx(txHash)
}

// MarkTxAsSentOnL2 marks the transaction as sent on L2
func (e *EthereumMultinodeClient) MarkTxAsSentOnL2(tx *types.Transaction) error {
	return e.DefaultClient.MarkTxAsSentOnL2(tx)
}

// PollFinality polls for finality
func (e *EthereumMultinodeClient) PollFinality() error {
	return e.DefaultClient.PollFinality()
}

// StopPollingForFinality stops polling for finality
func (e *EthereumMultinodeClient) CancelFinalityPolling() {
	e.DefaultClient.CancelFinalityPolling()
}

// GetTxReceipt returns the receipt of the transaction if available, error otherwise
func (e *EthereumMultinodeClient) GetTxReceipt(txHash common.Hash) (*types.Receipt, error) {
	return e.DefaultClient.GetTxReceipt(txHash)
}

func (e *EthereumMultinodeClient) RevertReasonFromTx(txHash common.Hash, abiString string) (string, interface{}, error) {
	return e.DefaultClient.RevertReasonFromTx(txHash, abiString)
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

func (e *EthereumMultinodeClient) EstimateTransactionGasCost() (*big.Int, error) {
	return e.DefaultClient.EstimateTransactionGasCost()
}

// GasStats gets gas stats instance
func (e *EthereumMultinodeClient) GasStats() *GasStats {
	return e.DefaultClient.GasStats()
}

func (e *EthereumMultinodeClient) EstimateGas(callMsg ethereum.CallMsg) (GasEstimations, error) {
	return e.DefaultClient.EstimateGas(callMsg)
}

func (e *EthereumMultinodeClient) EstimateGasPrice() (*big.Int, error) {
	return e.DefaultClient.EstimateGasPrice()
}

func (e *EthereumMultinodeClient) ConnectionIssue() chan time.Time {
	return e.DefaultClient.ConnectionIssue()
}
func (e *EthereumMultinodeClient) ConnectionRestored() chan time.Time {
	return e.DefaultClient.ConnectionRestored()
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *EthereumMultinodeClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	e.DefaultClient.AddHeaderEventSubscription(key, subscriber)
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (e *EthereumMultinodeClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, logs chan<- types.Log) (ethereum.Subscription, error) {
	return e.DefaultClient.SubscribeFilterLogs(ctx, q, logs)
}

// FilterLogs executes a filter query
func (e *EthereumMultinodeClient) FilterLogs(ctx context.Context, filterQuery ethereum.FilterQuery) ([]types.Log, error) {
	return e.DefaultClient.FilterLogs(ctx, filterQuery)
}

func (e *EthereumMultinodeClient) GetLatestFinalizedBlockHeader(ctx context.Context) (*types.Header, error) {
	return e.DefaultClient.GetLatestFinalizedBlockHeader(ctx)
}

func (e *EthereumMultinodeClient) EstimatedFinalizationTime(ctx context.Context) (time.Duration, error) {
	return e.DefaultClient.EstimatedFinalizationTime(ctx)
}

func (e *EthereumMultinodeClient) AvgBlockTime(ctx context.Context) (time.Duration, error) {
	return e.DefaultClient.AvgBlockTime(ctx)
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *EthereumMultinodeClient) DeleteHeaderEventSubscription(key string) {
	e.DefaultClient.DeleteHeaderEventSubscription(key)
}

// WaitForEvents is a blocking function that waits for all event subscriptions for all clients
func (e *EthereumMultinodeClient) WaitForEvents() error {
	g := errgroup.Group{}
	for _, c := range e.Clients {
		g.Go(func() error {
			return c.WaitForEvents()
		})
	}
	return g.Wait()
}

func (e *EthereumMultinodeClient) ErrorReason(b ethereum.ContractCaller, tx *types.Transaction, receipt *types.Receipt) (string, error) {
	return e.DefaultClient.ErrorReason(b, tx, receipt)
}

// InitializeHeaderSubscription initializes either subscription-based or polling-based header processing
func (e *EthereumMultinodeClient) InitializeHeaderSubscription() error {
	return e.DefaultClient.InitializeHeaderSubscription()
}
