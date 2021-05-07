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

// SafeMath64ABI is the input ABI used to generate the binding from.
const SafeMath64ABI = "[]"

// SafeMath64Bin is the compiled bytecode used for deploying new contracts.
var SafeMath64Bin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220dff2d28f4dc3bcf140c028cc56945af3b9a013071aff50bca693c7a6a6097d8164736f6c63430006060033"

// DeploySafeMath64 deploys a new Ethereum contract, binding an instance of SafeMath64 to it.
func DeploySafeMath64(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath64, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath64ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMath64Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath64{SafeMath64Caller: SafeMath64Caller{contract: contract}, SafeMath64Transactor: SafeMath64Transactor{contract: contract}, SafeMath64Filterer: SafeMath64Filterer{contract: contract}}, nil
}

// SafeMath64 is an auto generated Go binding around an Ethereum contract.
type SafeMath64 struct {
	SafeMath64Caller     // Read-only binding to the contract
	SafeMath64Transactor // Write-only binding to the contract
	SafeMath64Filterer   // Log filterer for contract events
}

// SafeMath64Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMath64Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath64Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMath64Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath64Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMath64Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath64Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMath64Session struct {
	Contract     *SafeMath64       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMath64CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMath64CallerSession struct {
	Contract *SafeMath64Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SafeMath64TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMath64TransactorSession struct {
	Contract     *SafeMath64Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SafeMath64Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMath64Raw struct {
	Contract *SafeMath64 // Generic contract binding to access the raw methods on
}

// SafeMath64CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMath64CallerRaw struct {
	Contract *SafeMath64Caller // Generic read-only contract binding to access the raw methods on
}

// SafeMath64TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMath64TransactorRaw struct {
	Contract *SafeMath64Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath64 creates a new instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64(address common.Address, backend bind.ContractBackend) (*SafeMath64, error) {
	contract, err := bindSafeMath64(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath64{SafeMath64Caller: SafeMath64Caller{contract: contract}, SafeMath64Transactor: SafeMath64Transactor{contract: contract}, SafeMath64Filterer: SafeMath64Filterer{contract: contract}}, nil
}

// NewSafeMath64Caller creates a new read-only instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64Caller(address common.Address, caller bind.ContractCaller) (*SafeMath64Caller, error) {
	contract, err := bindSafeMath64(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath64Caller{contract: contract}, nil
}

// NewSafeMath64Transactor creates a new write-only instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeMath64Transactor, error) {
	contract, err := bindSafeMath64(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath64Transactor{contract: contract}, nil
}

// NewSafeMath64Filterer creates a new log filterer instance of SafeMath64, bound to a specific deployed contract.
func NewSafeMath64Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeMath64Filterer, error) {
	contract, err := bindSafeMath64(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMath64Filterer{contract: contract}, nil
}

// bindSafeMath64 binds a generic wrapper to an already deployed contract.
func bindSafeMath64(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath64ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath64 *SafeMath64Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath64.Contract.SafeMath64Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath64 *SafeMath64Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath64.Contract.SafeMath64Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath64 *SafeMath64Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath64.Contract.SafeMath64Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath64 *SafeMath64CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath64.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath64 *SafeMath64TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath64.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath64 *SafeMath64TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath64.Contract.contract.Transact(opts, method, params...)
}
