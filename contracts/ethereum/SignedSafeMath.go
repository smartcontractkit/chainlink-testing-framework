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

// SignedSafeMathABI is the input ABI used to generate the binding from.
const SignedSafeMathABI = "[]"

// SignedSafeMathBin is the compiled bytecode used for deploying new contracts.
var SignedSafeMathBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220d56b6fee9e4a85ae7e811984c5baa98c381136f80f35598102559b3d9031534b64736f6c63430006060033"

// DeploySignedSafeMath deploys a new Ethereum contract, binding an instance of SignedSafeMath to it.
func DeploySignedSafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SignedSafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SignedSafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SignedSafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SignedSafeMath{SignedSafeMathCaller: SignedSafeMathCaller{contract: contract}, SignedSafeMathTransactor: SignedSafeMathTransactor{contract: contract}, SignedSafeMathFilterer: SignedSafeMathFilterer{contract: contract}}, nil
}

// SignedSafeMath is an auto generated Go binding around an Ethereum contract.
type SignedSafeMath struct {
	SignedSafeMathCaller     // Read-only binding to the contract
	SignedSafeMathTransactor // Write-only binding to the contract
	SignedSafeMathFilterer   // Log filterer for contract events
}

// SignedSafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignedSafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedSafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignedSafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedSafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignedSafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedSafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignedSafeMathSession struct {
	Contract     *SignedSafeMath   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignedSafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignedSafeMathCallerSession struct {
	Contract *SignedSafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SignedSafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignedSafeMathTransactorSession struct {
	Contract     *SignedSafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SignedSafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignedSafeMathRaw struct {
	Contract *SignedSafeMath // Generic contract binding to access the raw methods on
}

// SignedSafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignedSafeMathCallerRaw struct {
	Contract *SignedSafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SignedSafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignedSafeMathTransactorRaw struct {
	Contract *SignedSafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignedSafeMath creates a new instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMath(address common.Address, backend bind.ContractBackend) (*SignedSafeMath, error) {
	contract, err := bindSignedSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMath{SignedSafeMathCaller: SignedSafeMathCaller{contract: contract}, SignedSafeMathTransactor: SignedSafeMathTransactor{contract: contract}, SignedSafeMathFilterer: SignedSafeMathFilterer{contract: contract}}, nil
}

// NewSignedSafeMathCaller creates a new read-only instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SignedSafeMathCaller, error) {
	contract, err := bindSignedSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMathCaller{contract: contract}, nil
}

// NewSignedSafeMathTransactor creates a new write-only instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SignedSafeMathTransactor, error) {
	contract, err := bindSignedSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMathTransactor{contract: contract}, nil
}

// NewSignedSafeMathFilterer creates a new log filterer instance of SignedSafeMath, bound to a specific deployed contract.
func NewSignedSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SignedSafeMathFilterer, error) {
	contract, err := bindSignedSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignedSafeMathFilterer{contract: contract}, nil
}

// bindSignedSafeMath binds a generic wrapper to an already deployed contract.
func bindSignedSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SignedSafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignedSafeMath *SignedSafeMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignedSafeMath.Contract.SignedSafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignedSafeMath *SignedSafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.SignedSafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignedSafeMath *SignedSafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.SignedSafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignedSafeMath *SignedSafeMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignedSafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignedSafeMath *SignedSafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignedSafeMath *SignedSafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignedSafeMath.Contract.contract.Transact(opts, method, params...)
}
