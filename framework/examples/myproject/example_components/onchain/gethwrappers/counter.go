// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gethwrappers

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// GethwrappersMetaData contains all meta data concerning the Gethwrappers contract.
var GethwrappersMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"increment\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"number\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setNumber\",\"inputs\":[{\"name\":\"newNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220b4b78a03f7df1ffcbf8115e03ede66c3aae910f19ea5b824b889e4e6056ab55c64736f6c63430008180033",
}

// GethwrappersABI is the input ABI used to generate the binding from.
// Deprecated: Use GethwrappersMetaData.ABI instead.
var GethwrappersABI = GethwrappersMetaData.ABI

// GethwrappersBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GethwrappersMetaData.Bin instead.
var GethwrappersBin = GethwrappersMetaData.Bin

// DeployGethwrappers deploys a new Ethereum contract, binding an instance of Gethwrappers to it.
func DeployGethwrappers(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Gethwrappers, error) {
	parsed, err := GethwrappersMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GethwrappersBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Gethwrappers{GethwrappersCaller: GethwrappersCaller{contract: contract}, GethwrappersTransactor: GethwrappersTransactor{contract: contract}, GethwrappersFilterer: GethwrappersFilterer{contract: contract}}, nil
}

// Gethwrappers is an auto generated Go binding around an Ethereum contract.
type Gethwrappers struct {
	GethwrappersCaller     // Read-only binding to the contract
	GethwrappersTransactor // Write-only binding to the contract
	GethwrappersFilterer   // Log filterer for contract events
}

// GethwrappersCaller is an auto generated read-only Go binding around an Ethereum contract.
type GethwrappersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GethwrappersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GethwrappersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GethwrappersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GethwrappersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GethwrappersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GethwrappersSession struct {
	Contract     *Gethwrappers     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GethwrappersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GethwrappersCallerSession struct {
	Contract *GethwrappersCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GethwrappersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GethwrappersTransactorSession struct {
	Contract     *GethwrappersTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GethwrappersRaw is an auto generated low-level Go binding around an Ethereum contract.
type GethwrappersRaw struct {
	Contract *Gethwrappers // Generic contract binding to access the raw methods on
}

// GethwrappersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GethwrappersCallerRaw struct {
	Contract *GethwrappersCaller // Generic read-only contract binding to access the raw methods on
}

// GethwrappersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GethwrappersTransactorRaw struct {
	Contract *GethwrappersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGethwrappers creates a new instance of Gethwrappers, bound to a specific deployed contract.
func NewGethwrappers(address common.Address, backend bind.ContractBackend) (*Gethwrappers, error) {
	contract, err := bindGethwrappers(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Gethwrappers{GethwrappersCaller: GethwrappersCaller{contract: contract}, GethwrappersTransactor: GethwrappersTransactor{contract: contract}, GethwrappersFilterer: GethwrappersFilterer{contract: contract}}, nil
}

// NewGethwrappersCaller creates a new read-only instance of Gethwrappers, bound to a specific deployed contract.
func NewGethwrappersCaller(address common.Address, caller bind.ContractCaller) (*GethwrappersCaller, error) {
	contract, err := bindGethwrappers(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GethwrappersCaller{contract: contract}, nil
}

// NewGethwrappersTransactor creates a new write-only instance of Gethwrappers, bound to a specific deployed contract.
func NewGethwrappersTransactor(address common.Address, transactor bind.ContractTransactor) (*GethwrappersTransactor, error) {
	contract, err := bindGethwrappers(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GethwrappersTransactor{contract: contract}, nil
}

// NewGethwrappersFilterer creates a new log filterer instance of Gethwrappers, bound to a specific deployed contract.
func NewGethwrappersFilterer(address common.Address, filterer bind.ContractFilterer) (*GethwrappersFilterer, error) {
	contract, err := bindGethwrappers(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GethwrappersFilterer{contract: contract}, nil
}

// bindGethwrappers binds a generic wrapper to an already deployed contract.
func bindGethwrappers(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GethwrappersMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gethwrappers *GethwrappersRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gethwrappers.Contract.GethwrappersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gethwrappers *GethwrappersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gethwrappers.Contract.GethwrappersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gethwrappers *GethwrappersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gethwrappers.Contract.GethwrappersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gethwrappers *GethwrappersCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gethwrappers.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gethwrappers *GethwrappersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gethwrappers.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gethwrappers *GethwrappersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gethwrappers.Contract.contract.Transact(opts, method, params...)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_Gethwrappers *GethwrappersCaller) Number(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Gethwrappers.contract.Call(opts, &out, "number")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_Gethwrappers *GethwrappersSession) Number() (*big.Int, error) {
	return _Gethwrappers.Contract.Number(&_Gethwrappers.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_Gethwrappers *GethwrappersCallerSession) Number() (*big.Int, error) {
	return _Gethwrappers.Contract.Number(&_Gethwrappers.CallOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Gethwrappers *GethwrappersTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gethwrappers.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Gethwrappers *GethwrappersSession) Increment() (*types.Transaction, error) {
	return _Gethwrappers.Contract.Increment(&_Gethwrappers.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Gethwrappers *GethwrappersTransactorSession) Increment() (*types.Transaction, error) {
	return _Gethwrappers.Contract.Increment(&_Gethwrappers.TransactOpts)
}

// SetNumber is a paid mutator transaction binding the contract method 0x3fb5c1cb.
//
// Solidity: function setNumber(uint256 newNumber) returns()
func (_Gethwrappers *GethwrappersTransactor) SetNumber(opts *bind.TransactOpts, newNumber *big.Int) (*types.Transaction, error) {
	return _Gethwrappers.contract.Transact(opts, "setNumber", newNumber)
}

// SetNumber is a paid mutator transaction binding the contract method 0x3fb5c1cb.
//
// Solidity: function setNumber(uint256 newNumber) returns()
func (_Gethwrappers *GethwrappersSession) SetNumber(newNumber *big.Int) (*types.Transaction, error) {
	return _Gethwrappers.Contract.SetNumber(&_Gethwrappers.TransactOpts, newNumber)
}

// SetNumber is a paid mutator transaction binding the contract method 0x3fb5c1cb.
//
// Solidity: function setNumber(uint256 newNumber) returns()
func (_Gethwrappers *GethwrappersTransactorSession) SetNumber(newNumber *big.Int) (*types.Transaction, error) {
	return _Gethwrappers.Contract.SetNumber(&_Gethwrappers.TransactOpts, newNumber)
}
