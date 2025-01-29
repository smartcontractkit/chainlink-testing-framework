// Package blockchain handles connections to various blockchains
package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EVMClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
// EVMClient can be connected to a single or multiple nodes,
type EVMClient interface {
	// Getters
	Get() interface{}
	GetNetworkName() string
	NetworkSimulated() bool
	GetChainID() *big.Int
	GetClients() []EVMClient
	GetDefaultWallet() *EthereumWallet
	GetWallets() []*EthereumWallet
	GetWalletByAddress(address common.Address) *EthereumWallet
	GetNetworkConfig() *EVMNetwork
	GetNonceSetting() NonceSettings

	GetHeaderSubscriptions() map[string]HeaderEventSubscription

	// Setters
	SetID(id int)
	SetDefaultWallet(num int) error
	SetDefaultWalletByAddress(address common.Address) error
	SetWallets([]*EthereumWallet)
	LoadWallets(ns EVMNetwork) error
	SwitchNode(node int) error
	SyncNonce(c EVMClient)

	// On-chain Operations
	BalanceAt(ctx context.Context, address common.Address) (*big.Int, error)
	HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error)
	HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error)
	LatestBlockNumber(ctx context.Context) (uint64, error)
	GetLatestFinalizedBlockHeader(ctx context.Context) (*types.Header, error)
	AvgBlockTime(ctx context.Context) (time.Duration, error)
	EstimatedFinalizationTime(ctx context.Context) (time.Duration, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	Fund(toAddress string, amount *big.Float, gasEstimations GasEstimations) error
	ReturnFunds(fromKey *ecdsa.PrivateKey) error
	DeployContract(
		contractName string,
		deployer ContractDeployer,
	) (*common.Address, *types.Transaction, interface{}, error)
	// TODO: Load and Deploy need to both agree to use an address pointer, there's unnecessary casting going on
	LoadContract(contractName string, address common.Address, loader ContractLoader) (interface{}, error)
	TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error)
	NewTx(
		fromPrivateKey *ecdsa.PrivateKey,
		nonce uint64,
		to common.Address,
		value *big.Int,
		gasEstimations GasEstimations,
	) (*types.Transaction, error)
	ProcessTransaction(tx *types.Transaction) error

	MarkTxAsSentOnL2(tx *types.Transaction) error
	ProcessEvent(name string, event *types.Log, confirmedChan chan bool, errorChan chan error) error
	IsEventConfirmed(event *types.Log) (confirmed, removed bool, err error)
	IsTxConfirmed(txHash common.Hash) (bool, error)
	IsTxHeadFinalized(txHdr, header *SafeEVMHeader) (bool, *big.Int, time.Time, error)
	WaitForFinalizedTx(txHash common.Hash) (*big.Int, time.Time, error)
	PollFinality() error
	CancelFinalityPolling()
	GetTxReceipt(txHash common.Hash) (*types.Receipt, error)
	RevertReasonFromTx(txHash common.Hash, abiString string) (string, interface{}, error)
	ErrorReason(b ethereum.ContractCaller, tx *types.Transaction, receipt *types.Receipt) (string, error)

	ParallelTransactions(enabled bool)
	Close() error
	Backend() bind.ContractBackend
	DeployBackend() bind.DeployBackend
	// Deal with wrapped headers
	SubscribeNewHeaders(ctx context.Context, headerChan chan *SafeEVMHeader) (ethereum.Subscription, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*SafeEVMHeader, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*SafeEVMHeader, error)

	// Gas Operations
	EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error)
	EstimateTransactionGasCost() (*big.Int, error)
	GasStats() *GasStats
	// EstimateGas provides all gas stats needed, best for estimating gas and prices for a specific transaction
	EstimateGas(callMsg ethereum.CallMsg) (GasEstimations, error)
	// EstimateGasPrice provides a plain gas price estimate, best for quick checks and contract deployments
	EstimateGasPrice() (*big.Int, error)

	// Connection Status
	// ConnectionIssue returns a channel that will receive a timestamp when the connection is lost
	ConnectionIssue() chan time.Time
	// ConnectionRestored returns a channel that will receive a timestamp when the connection is restored
	ConnectionRestored() chan time.Time

	// Event Subscriptions
	AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription)
	DeleteHeaderEventSubscription(key string)
	WaitForEvents() error
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)

	// Polling Events
	FilterLogs(ctx context.Context, filterQuery ethereum.FilterQuery) ([]types.Log, error)

	RawJsonRPCCall(ctx context.Context, result interface{}, method string, params ...interface{}) error

	GetEthClient() *ethclient.Client

	InitializeHeaderSubscription() error
}

// NodeHeader header with the ID of the node that received it
type NodeHeader struct {
	NodeID int
	SafeEVMHeader
}

// SafeEVMHeader is a wrapper for the EVM header, to allow for the addition/removal of fields without breaking the interface
type SafeEVMHeader struct {
	Hash      common.Hash
	Number    *big.Int
	Timestamp time.Time
	BaseFee   *big.Int
}

// GasEstimations is a wrapper for the gas estimations
type GasEstimations struct {
	GasUnits     uint64   // How many units of gas the transaction will use
	GasPrice     *big.Int // Gas price of the transaction (for Legacy transactions)
	GasTipCap    *big.Int // Gas tip cap of the transaction (for DynamicFee transactions)
	GasFeeCap    *big.Int // Gas fee cap of the transaction (for DynamicFee transactions)
	TotalGasCost *big.Int // Total gas cost of the transaction (gas units * total gas price)
}

// UnmarshalJSON enables Geth to unmarshal block headers into our custom type
func (h *SafeEVMHeader) UnmarshalJSON(bs []byte) error {
	type head struct {
		Hash      common.Hash    `json:"hash"`
		Number    *hexutil.Big   `json:"number"`
		Timestamp hexutil.Uint64 `json:"timestamp"`
		BaseFee   *hexutil.Big   `json:"baseFeePerGas"`
	}

	var jsonHead head
	err := json.Unmarshal(bs, &jsonHead)
	if err != nil {
		return err
	}

	if jsonHead.Number == nil {
		*h = SafeEVMHeader{}
		return nil
	}

	h.Hash = jsonHead.Hash
	h.Number = (*big.Int)(jsonHead.Number)
	h.Timestamp = time.Unix(int64(jsonHead.Timestamp), 0).UTC() //nolint
	h.BaseFee = (*big.Int)(jsonHead.BaseFee)
	return nil
}

// HeaderEventSubscription is an interface for allowing callbacks when the client receives a new header
type HeaderEventSubscription interface {
	ReceiveHeader(header NodeHeader) error
	Wait() error
	Complete() bool
}

// ContractDeployer acts as a go-between function for general contract deployment
type ContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)

// ContractLoader helps loading already deployed contracts
type ContractLoader func(address common.Address, backend bind.ContractBackend) (
	interface{},
	error,
)
