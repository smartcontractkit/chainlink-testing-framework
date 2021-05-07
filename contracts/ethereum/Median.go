// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MedianABI is the input ABI used to generate the binding from.
const MedianABI = "[]"

// MedianBin is the compiled bytecode used for deploying new contracts.
var MedianBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212207a93346690778f30f0cea1edf7aa63ef2bfa9735eb54c08082beeeb225c1ee7f64736f6c63430006060033"

// DeployMedian deploys a new Ethereum contract, binding an instance of Median to it.
func DeployMedian(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Median, error) {
	parsed, err := abi.JSON(strings.NewReader(MedianABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MedianBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Median{MedianCaller: MedianCaller{contract: contract}, MedianTransactor: MedianTransactor{contract: contract}, MedianFilterer: MedianFilterer{contract: contract}}, nil
}

// Median is an auto generated Go binding around an Ethereum contract.
type Median struct {
	MedianCaller     // Read-only binding to the contract
	MedianTransactor // Write-only binding to the contract
	MedianFilterer   // Log filterer for contract events
}

// MedianCaller is an auto generated read-only Go binding around an Ethereum contract.
type MedianCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedianTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MedianTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedianFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MedianFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedianSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MedianSession struct {
	Contract     *Median           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MedianCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MedianCallerSession struct {
	Contract *MedianCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MedianTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MedianTransactorSession struct {
	Contract     *MedianTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MedianRaw is an auto generated low-level Go binding around an Ethereum contract.
type MedianRaw struct {
	Contract *Median // Generic contract binding to access the raw methods on
}

// MedianCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MedianCallerRaw struct {
	Contract *MedianCaller // Generic read-only contract binding to access the raw methods on
}

// MedianTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MedianTransactorRaw struct {
	Contract *MedianTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMedian creates a new instance of Median, bound to a specific deployed contract.
func NewMedian(address common.Address, backend bind.ContractBackend) (*Median, error) {
	contract, err := bindMedian(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Median{MedianCaller: MedianCaller{contract: contract}, MedianTransactor: MedianTransactor{contract: contract}, MedianFilterer: MedianFilterer{contract: contract}}, nil
}

// NewMedianCaller creates a new read-only instance of Median, bound to a specific deployed contract.
func NewMedianCaller(address common.Address, caller bind.ContractCaller) (*MedianCaller, error) {
	contract, err := bindMedian(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MedianCaller{contract: contract}, nil
}

// NewMedianTransactor creates a new write-only instance of Median, bound to a specific deployed contract.
func NewMedianTransactor(address common.Address, transactor bind.ContractTransactor) (*MedianTransactor, error) {
	contract, err := bindMedian(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MedianTransactor{contract: contract}, nil
}

// NewMedianFilterer creates a new log filterer instance of Median, bound to a specific deployed contract.
func NewMedianFilterer(address common.Address, filterer bind.ContractFilterer) (*MedianFilterer, error) {
	contract, err := bindMedian(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MedianFilterer{contract: contract}, nil
}

// bindMedian binds a generic wrapper to an already deployed contract.
func bindMedian(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MedianABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Median *MedianRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Median.Contract.MedianCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Median *MedianRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Median.Contract.MedianTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Median *MedianRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Median.Contract.MedianTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Median *MedianCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Median.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Median *MedianTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Median.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Median *MedianTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Median.Contract.contract.Transact(opts, method, params...)
}
