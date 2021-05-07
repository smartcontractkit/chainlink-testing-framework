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

// CheckedMathABI is the input ABI used to generate the binding from.
const CheckedMathABI = "[]"

// CheckedMathBin is the compiled bytecode used for deploying new contracts.
var CheckedMathBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220e1183181f20547292edcbc6755d35df7377af47499632153d0179bcc0d2ce18664736f6c63430006060033"

// DeployCheckedMath deploys a new Ethereum contract, binding an instance of CheckedMath to it.
func DeployCheckedMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CheckedMath, error) {
	parsed, err := abi.JSON(strings.NewReader(CheckedMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CheckedMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CheckedMath{CheckedMathCaller: CheckedMathCaller{contract: contract}, CheckedMathTransactor: CheckedMathTransactor{contract: contract}, CheckedMathFilterer: CheckedMathFilterer{contract: contract}}, nil
}

// CheckedMath is an auto generated Go binding around an Ethereum contract.
type CheckedMath struct {
	CheckedMathCaller     // Read-only binding to the contract
	CheckedMathTransactor // Write-only binding to the contract
	CheckedMathFilterer   // Log filterer for contract events
}

// CheckedMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type CheckedMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckedMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CheckedMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckedMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CheckedMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckedMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CheckedMathSession struct {
	Contract     *CheckedMath      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CheckedMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CheckedMathCallerSession struct {
	Contract *CheckedMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// CheckedMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CheckedMathTransactorSession struct {
	Contract     *CheckedMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// CheckedMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type CheckedMathRaw struct {
	Contract *CheckedMath // Generic contract binding to access the raw methods on
}

// CheckedMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CheckedMathCallerRaw struct {
	Contract *CheckedMathCaller // Generic read-only contract binding to access the raw methods on
}

// CheckedMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CheckedMathTransactorRaw struct {
	Contract *CheckedMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCheckedMath creates a new instance of CheckedMath, bound to a specific deployed contract.
func NewCheckedMath(address common.Address, backend bind.ContractBackend) (*CheckedMath, error) {
	contract, err := bindCheckedMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CheckedMath{CheckedMathCaller: CheckedMathCaller{contract: contract}, CheckedMathTransactor: CheckedMathTransactor{contract: contract}, CheckedMathFilterer: CheckedMathFilterer{contract: contract}}, nil
}

// NewCheckedMathCaller creates a new read-only instance of CheckedMath, bound to a specific deployed contract.
func NewCheckedMathCaller(address common.Address, caller bind.ContractCaller) (*CheckedMathCaller, error) {
	contract, err := bindCheckedMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CheckedMathCaller{contract: contract}, nil
}

// NewCheckedMathTransactor creates a new write-only instance of CheckedMath, bound to a specific deployed contract.
func NewCheckedMathTransactor(address common.Address, transactor bind.ContractTransactor) (*CheckedMathTransactor, error) {
	contract, err := bindCheckedMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CheckedMathTransactor{contract: contract}, nil
}

// NewCheckedMathFilterer creates a new log filterer instance of CheckedMath, bound to a specific deployed contract.
func NewCheckedMathFilterer(address common.Address, filterer bind.ContractFilterer) (*CheckedMathFilterer, error) {
	contract, err := bindCheckedMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CheckedMathFilterer{contract: contract}, nil
}

// bindCheckedMath binds a generic wrapper to an already deployed contract.
func bindCheckedMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CheckedMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CheckedMath *CheckedMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CheckedMath.Contract.CheckedMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CheckedMath *CheckedMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CheckedMath.Contract.CheckedMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CheckedMath *CheckedMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CheckedMath.Contract.CheckedMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CheckedMath *CheckedMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CheckedMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CheckedMath *CheckedMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CheckedMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CheckedMath *CheckedMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CheckedMath.Contract.contract.Transact(opts, method, params...)
}
