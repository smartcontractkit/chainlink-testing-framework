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

// SafeMath32ABI is the input ABI used to generate the binding from.
const SafeMath32ABI = "[]"

// SafeMath32Bin is the compiled bytecode used for deploying new contracts.
var SafeMath32Bin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220d56dd678fccb45b2f4bd9a7d4ccba4db22281f90ea595c5bbc9455d1ecdc2a7b64736f6c63430006060033"

// DeploySafeMath32 deploys a new Ethereum contract, binding an instance of SafeMath32 to it.
func DeploySafeMath32(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath32, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath32ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMath32Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath32{SafeMath32Caller: SafeMath32Caller{contract: contract}, SafeMath32Transactor: SafeMath32Transactor{contract: contract}, SafeMath32Filterer: SafeMath32Filterer{contract: contract}}, nil
}

// SafeMath32 is an auto generated Go binding around an Ethereum contract.
type SafeMath32 struct {
	SafeMath32Caller     // Read-only binding to the contract
	SafeMath32Transactor // Write-only binding to the contract
	SafeMath32Filterer   // Log filterer for contract events
}

// SafeMath32Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMath32Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath32Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMath32Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath32Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMath32Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath32Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMath32Session struct {
	Contract     *SafeMath32       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMath32CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMath32CallerSession struct {
	Contract *SafeMath32Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SafeMath32TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMath32TransactorSession struct {
	Contract     *SafeMath32Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SafeMath32Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMath32Raw struct {
	Contract *SafeMath32 // Generic contract binding to access the raw methods on
}

// SafeMath32CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMath32CallerRaw struct {
	Contract *SafeMath32Caller // Generic read-only contract binding to access the raw methods on
}

// SafeMath32TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMath32TransactorRaw struct {
	Contract *SafeMath32Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath32 creates a new instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32(address common.Address, backend bind.ContractBackend) (*SafeMath32, error) {
	contract, err := bindSafeMath32(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath32{SafeMath32Caller: SafeMath32Caller{contract: contract}, SafeMath32Transactor: SafeMath32Transactor{contract: contract}, SafeMath32Filterer: SafeMath32Filterer{contract: contract}}, nil
}

// NewSafeMath32Caller creates a new read-only instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32Caller(address common.Address, caller bind.ContractCaller) (*SafeMath32Caller, error) {
	contract, err := bindSafeMath32(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath32Caller{contract: contract}, nil
}

// NewSafeMath32Transactor creates a new write-only instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeMath32Transactor, error) {
	contract, err := bindSafeMath32(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath32Transactor{contract: contract}, nil
}

// NewSafeMath32Filterer creates a new log filterer instance of SafeMath32, bound to a specific deployed contract.
func NewSafeMath32Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeMath32Filterer, error) {
	contract, err := bindSafeMath32(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMath32Filterer{contract: contract}, nil
}

// bindSafeMath32 binds a generic wrapper to an already deployed contract.
func bindSafeMath32(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath32ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath32 *SafeMath32Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath32.Contract.SafeMath32Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath32 *SafeMath32Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath32.Contract.SafeMath32Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath32 *SafeMath32Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath32.Contract.SafeMath32Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath32 *SafeMath32CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath32.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath32 *SafeMath32TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath32.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath32 *SafeMath32TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath32.Contract.contract.Transact(opts, method, params...)
}
