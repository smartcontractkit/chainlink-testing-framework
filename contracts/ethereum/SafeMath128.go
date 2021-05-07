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

// SafeMath128ABI is the input ABI used to generate the binding from.
const SafeMath128ABI = "[]"

// SafeMath128Bin is the compiled bytecode used for deploying new contracts.
var SafeMath128Bin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220732c9427201297b1f872704b6c74711ffca0d360a39ea4b44ad467938b4e669264736f6c63430006060033"

// DeploySafeMath128 deploys a new Ethereum contract, binding an instance of SafeMath128 to it.
func DeploySafeMath128(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath128, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath128ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMath128Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath128{SafeMath128Caller: SafeMath128Caller{contract: contract}, SafeMath128Transactor: SafeMath128Transactor{contract: contract}, SafeMath128Filterer: SafeMath128Filterer{contract: contract}}, nil
}

// SafeMath128 is an auto generated Go binding around an Ethereum contract.
type SafeMath128 struct {
	SafeMath128Caller     // Read-only binding to the contract
	SafeMath128Transactor // Write-only binding to the contract
	SafeMath128Filterer   // Log filterer for contract events
}

// SafeMath128Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMath128Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath128Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMath128Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath128Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMath128Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMath128Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMath128Session struct {
	Contract     *SafeMath128      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMath128CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMath128CallerSession struct {
	Contract *SafeMath128Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SafeMath128TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMath128TransactorSession struct {
	Contract     *SafeMath128Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SafeMath128Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMath128Raw struct {
	Contract *SafeMath128 // Generic contract binding to access the raw methods on
}

// SafeMath128CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMath128CallerRaw struct {
	Contract *SafeMath128Caller // Generic read-only contract binding to access the raw methods on
}

// SafeMath128TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMath128TransactorRaw struct {
	Contract *SafeMath128Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath128 creates a new instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128(address common.Address, backend bind.ContractBackend) (*SafeMath128, error) {
	contract, err := bindSafeMath128(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath128{SafeMath128Caller: SafeMath128Caller{contract: contract}, SafeMath128Transactor: SafeMath128Transactor{contract: contract}, SafeMath128Filterer: SafeMath128Filterer{contract: contract}}, nil
}

// NewSafeMath128Caller creates a new read-only instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128Caller(address common.Address, caller bind.ContractCaller) (*SafeMath128Caller, error) {
	contract, err := bindSafeMath128(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath128Caller{contract: contract}, nil
}

// NewSafeMath128Transactor creates a new write-only instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeMath128Transactor, error) {
	contract, err := bindSafeMath128(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMath128Transactor{contract: contract}, nil
}

// NewSafeMath128Filterer creates a new log filterer instance of SafeMath128, bound to a specific deployed contract.
func NewSafeMath128Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeMath128Filterer, error) {
	contract, err := bindSafeMath128(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMath128Filterer{contract: contract}, nil
}

// bindSafeMath128 binds a generic wrapper to an already deployed contract.
func bindSafeMath128(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMath128ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath128 *SafeMath128Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath128.Contract.SafeMath128Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath128 *SafeMath128Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath128.Contract.SafeMath128Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath128 *SafeMath128Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath128.Contract.SafeMath128Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath128 *SafeMath128CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath128.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath128 *SafeMath128TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath128.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath128 *SafeMath128TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath128.Contract.contract.Transact(opts, method, params...)
}
