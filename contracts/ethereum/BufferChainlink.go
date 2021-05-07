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

// BufferChainlinkABI is the input ABI used to generate the binding from.
const BufferChainlinkABI = "[]"

// BufferChainlinkBin is the compiled bytecode used for deploying new contracts.
var BufferChainlinkBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220965cfee118e33e79c4e981093a1adc42825116732b18dfadb01e8a8a1ae9eef064736f6c63430006060033"

// DeployBufferChainlink deploys a new Ethereum contract, binding an instance of BufferChainlink to it.
func DeployBufferChainlink(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BufferChainlink, error) {
	parsed, err := abi.JSON(strings.NewReader(BufferChainlinkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BufferChainlinkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BufferChainlink{BufferChainlinkCaller: BufferChainlinkCaller{contract: contract}, BufferChainlinkTransactor: BufferChainlinkTransactor{contract: contract}, BufferChainlinkFilterer: BufferChainlinkFilterer{contract: contract}}, nil
}

// BufferChainlink is an auto generated Go binding around an Ethereum contract.
type BufferChainlink struct {
	BufferChainlinkCaller     // Read-only binding to the contract
	BufferChainlinkTransactor // Write-only binding to the contract
	BufferChainlinkFilterer   // Log filterer for contract events
}

// BufferChainlinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type BufferChainlinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BufferChainlinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BufferChainlinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BufferChainlinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BufferChainlinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BufferChainlinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BufferChainlinkSession struct {
	Contract     *BufferChainlink  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BufferChainlinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BufferChainlinkCallerSession struct {
	Contract *BufferChainlinkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// BufferChainlinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BufferChainlinkTransactorSession struct {
	Contract     *BufferChainlinkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// BufferChainlinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type BufferChainlinkRaw struct {
	Contract *BufferChainlink // Generic contract binding to access the raw methods on
}

// BufferChainlinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BufferChainlinkCallerRaw struct {
	Contract *BufferChainlinkCaller // Generic read-only contract binding to access the raw methods on
}

// BufferChainlinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BufferChainlinkTransactorRaw struct {
	Contract *BufferChainlinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBufferChainlink creates a new instance of BufferChainlink, bound to a specific deployed contract.
func NewBufferChainlink(address common.Address, backend bind.ContractBackend) (*BufferChainlink, error) {
	contract, err := bindBufferChainlink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BufferChainlink{BufferChainlinkCaller: BufferChainlinkCaller{contract: contract}, BufferChainlinkTransactor: BufferChainlinkTransactor{contract: contract}, BufferChainlinkFilterer: BufferChainlinkFilterer{contract: contract}}, nil
}

// NewBufferChainlinkCaller creates a new read-only instance of BufferChainlink, bound to a specific deployed contract.
func NewBufferChainlinkCaller(address common.Address, caller bind.ContractCaller) (*BufferChainlinkCaller, error) {
	contract, err := bindBufferChainlink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BufferChainlinkCaller{contract: contract}, nil
}

// NewBufferChainlinkTransactor creates a new write-only instance of BufferChainlink, bound to a specific deployed contract.
func NewBufferChainlinkTransactor(address common.Address, transactor bind.ContractTransactor) (*BufferChainlinkTransactor, error) {
	contract, err := bindBufferChainlink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BufferChainlinkTransactor{contract: contract}, nil
}

// NewBufferChainlinkFilterer creates a new log filterer instance of BufferChainlink, bound to a specific deployed contract.
func NewBufferChainlinkFilterer(address common.Address, filterer bind.ContractFilterer) (*BufferChainlinkFilterer, error) {
	contract, err := bindBufferChainlink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BufferChainlinkFilterer{contract: contract}, nil
}

// bindBufferChainlink binds a generic wrapper to an already deployed contract.
func bindBufferChainlink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BufferChainlinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BufferChainlink *BufferChainlinkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BufferChainlink.Contract.BufferChainlinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BufferChainlink *BufferChainlinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BufferChainlink.Contract.BufferChainlinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BufferChainlink *BufferChainlinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BufferChainlink.Contract.BufferChainlinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BufferChainlink *BufferChainlinkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BufferChainlink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BufferChainlink *BufferChainlinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BufferChainlink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BufferChainlink *BufferChainlinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BufferChainlink.Contract.contract.Transact(opts, method, params...)
}
