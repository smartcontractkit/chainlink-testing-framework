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

// SafeMathChainlinkABI is the input ABI used to generate the binding from.
const SafeMathChainlinkABI = "[]"

// SafeMathChainlinkBin is the compiled bytecode used for deploying new contracts.
var SafeMathChainlinkBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212209c64aa16a93bb5443357041705f36bddc4c00d371cb430bf3fc73a2c103ddb3364736f6c63430006060033"

// DeploySafeMathChainlink deploys a new Ethereum contract, binding an instance of SafeMathChainlink to it.
func DeploySafeMathChainlink(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMathChainlink, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathChainlinkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathChainlinkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMathChainlink{SafeMathChainlinkCaller: SafeMathChainlinkCaller{contract: contract}, SafeMathChainlinkTransactor: SafeMathChainlinkTransactor{contract: contract}, SafeMathChainlinkFilterer: SafeMathChainlinkFilterer{contract: contract}}, nil
}

// SafeMathChainlink is an auto generated Go binding around an Ethereum contract.
type SafeMathChainlink struct {
	SafeMathChainlinkCaller     // Read-only binding to the contract
	SafeMathChainlinkTransactor // Write-only binding to the contract
	SafeMathChainlinkFilterer   // Log filterer for contract events
}

// SafeMathChainlinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathChainlinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathChainlinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathChainlinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathChainlinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathChainlinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathChainlinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathChainlinkSession struct {
	Contract     *SafeMathChainlink // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SafeMathChainlinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathChainlinkCallerSession struct {
	Contract *SafeMathChainlinkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// SafeMathChainlinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathChainlinkTransactorSession struct {
	Contract     *SafeMathChainlinkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// SafeMathChainlinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathChainlinkRaw struct {
	Contract *SafeMathChainlink // Generic contract binding to access the raw methods on
}

// SafeMathChainlinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathChainlinkCallerRaw struct {
	Contract *SafeMathChainlinkCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathChainlinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathChainlinkTransactorRaw struct {
	Contract *SafeMathChainlinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMathChainlink creates a new instance of SafeMathChainlink, bound to a specific deployed contract.
func NewSafeMathChainlink(address common.Address, backend bind.ContractBackend) (*SafeMathChainlink, error) {
	contract, err := bindSafeMathChainlink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMathChainlink{SafeMathChainlinkCaller: SafeMathChainlinkCaller{contract: contract}, SafeMathChainlinkTransactor: SafeMathChainlinkTransactor{contract: contract}, SafeMathChainlinkFilterer: SafeMathChainlinkFilterer{contract: contract}}, nil
}

// NewSafeMathChainlinkCaller creates a new read-only instance of SafeMathChainlink, bound to a specific deployed contract.
func NewSafeMathChainlinkCaller(address common.Address, caller bind.ContractCaller) (*SafeMathChainlinkCaller, error) {
	contract, err := bindSafeMathChainlink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathChainlinkCaller{contract: contract}, nil
}

// NewSafeMathChainlinkTransactor creates a new write-only instance of SafeMathChainlink, bound to a specific deployed contract.
func NewSafeMathChainlinkTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathChainlinkTransactor, error) {
	contract, err := bindSafeMathChainlink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathChainlinkTransactor{contract: contract}, nil
}

// NewSafeMathChainlinkFilterer creates a new log filterer instance of SafeMathChainlink, bound to a specific deployed contract.
func NewSafeMathChainlinkFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathChainlinkFilterer, error) {
	contract, err := bindSafeMathChainlink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathChainlinkFilterer{contract: contract}, nil
}

// bindSafeMathChainlink binds a generic wrapper to an already deployed contract.
func bindSafeMathChainlink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathChainlinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMathChainlink *SafeMathChainlinkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMathChainlink.Contract.SafeMathChainlinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMathChainlink *SafeMathChainlinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMathChainlink.Contract.SafeMathChainlinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMathChainlink *SafeMathChainlinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMathChainlink.Contract.SafeMathChainlinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMathChainlink *SafeMathChainlinkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMathChainlink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMathChainlink *SafeMathChainlinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMathChainlink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMathChainlink *SafeMathChainlinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMathChainlink.Contract.contract.Transact(opts, method, params...)
}
