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

// FlagsInterfaceABI is the input ABI used to generate the binding from.
const FlagsInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"getFlag\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"name\":\"lowerFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"raiseFlag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"name\":\"raiseFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"setRaisingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// FlagsInterface is an auto generated Go binding around an Ethereum contract.
type FlagsInterface struct {
	FlagsInterfaceCaller     // Read-only binding to the contract
	FlagsInterfaceTransactor // Write-only binding to the contract
	FlagsInterfaceFilterer   // Log filterer for contract events
}

// FlagsInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlagsInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlagsInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlagsInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlagsInterfaceSession struct {
	Contract     *FlagsInterface   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlagsInterfaceCallerSession struct {
	Contract *FlagsInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// FlagsInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlagsInterfaceTransactorSession struct {
	Contract     *FlagsInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// FlagsInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlagsInterfaceRaw struct {
	Contract *FlagsInterface // Generic contract binding to access the raw methods on
}

// FlagsInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlagsInterfaceCallerRaw struct {
	Contract *FlagsInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// FlagsInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlagsInterfaceTransactorRaw struct {
	Contract *FlagsInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlagsInterface creates a new instance of FlagsInterface, bound to a specific deployed contract.
func NewFlagsInterface(address common.Address, backend bind.ContractBackend) (*FlagsInterface, error) {
	contract, err := bindFlagsInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FlagsInterface{FlagsInterfaceCaller: FlagsInterfaceCaller{contract: contract}, FlagsInterfaceTransactor: FlagsInterfaceTransactor{contract: contract}, FlagsInterfaceFilterer: FlagsInterfaceFilterer{contract: contract}}, nil
}

// NewFlagsInterfaceCaller creates a new read-only instance of FlagsInterface, bound to a specific deployed contract.
func NewFlagsInterfaceCaller(address common.Address, caller bind.ContractCaller) (*FlagsInterfaceCaller, error) {
	contract, err := bindFlagsInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsInterfaceCaller{contract: contract}, nil
}

// NewFlagsInterfaceTransactor creates a new write-only instance of FlagsInterface, bound to a specific deployed contract.
func NewFlagsInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*FlagsInterfaceTransactor, error) {
	contract, err := bindFlagsInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsInterfaceTransactor{contract: contract}, nil
}

// NewFlagsInterfaceFilterer creates a new log filterer instance of FlagsInterface, bound to a specific deployed contract.
func NewFlagsInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*FlagsInterfaceFilterer, error) {
	contract, err := bindFlagsInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlagsInterfaceFilterer{contract: contract}, nil
}

// bindFlagsInterface binds a generic wrapper to an already deployed contract.
func bindFlagsInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlagsInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlagsInterface *FlagsInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlagsInterface.Contract.FlagsInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlagsInterface *FlagsInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlagsInterface.Contract.FlagsInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlagsInterface *FlagsInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlagsInterface.Contract.FlagsInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlagsInterface *FlagsInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlagsInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlagsInterface *FlagsInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlagsInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlagsInterface *FlagsInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlagsInterface.Contract.contract.Transact(opts, method, params...)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address ) view returns(bool)
func (_FlagsInterface *FlagsInterfaceCaller) GetFlag(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _FlagsInterface.contract.Call(opts, &out, "getFlag", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address ) view returns(bool)
func (_FlagsInterface *FlagsInterfaceSession) GetFlag(arg0 common.Address) (bool, error) {
	return _FlagsInterface.Contract.GetFlag(&_FlagsInterface.CallOpts, arg0)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address ) view returns(bool)
func (_FlagsInterface *FlagsInterfaceCallerSession) GetFlag(arg0 common.Address) (bool, error) {
	return _FlagsInterface.Contract.GetFlag(&_FlagsInterface.CallOpts, arg0)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] ) view returns(bool[])
func (_FlagsInterface *FlagsInterfaceCaller) GetFlags(opts *bind.CallOpts, arg0 []common.Address) ([]bool, error) {
	var out []interface{}
	err := _FlagsInterface.contract.Call(opts, &out, "getFlags", arg0)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] ) view returns(bool[])
func (_FlagsInterface *FlagsInterfaceSession) GetFlags(arg0 []common.Address) ([]bool, error) {
	return _FlagsInterface.Contract.GetFlags(&_FlagsInterface.CallOpts, arg0)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] ) view returns(bool[])
func (_FlagsInterface *FlagsInterfaceCallerSession) GetFlags(arg0 []common.Address) ([]bool, error) {
	return _FlagsInterface.Contract.GetFlags(&_FlagsInterface.CallOpts, arg0)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] ) returns()
func (_FlagsInterface *FlagsInterfaceTransactor) LowerFlags(opts *bind.TransactOpts, arg0 []common.Address) (*types.Transaction, error) {
	return _FlagsInterface.contract.Transact(opts, "lowerFlags", arg0)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] ) returns()
func (_FlagsInterface *FlagsInterfaceSession) LowerFlags(arg0 []common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.LowerFlags(&_FlagsInterface.TransactOpts, arg0)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] ) returns()
func (_FlagsInterface *FlagsInterfaceTransactorSession) LowerFlags(arg0 []common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.LowerFlags(&_FlagsInterface.TransactOpts, arg0)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address ) returns()
func (_FlagsInterface *FlagsInterfaceTransactor) RaiseFlag(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _FlagsInterface.contract.Transact(opts, "raiseFlag", arg0)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address ) returns()
func (_FlagsInterface *FlagsInterfaceSession) RaiseFlag(arg0 common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.RaiseFlag(&_FlagsInterface.TransactOpts, arg0)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address ) returns()
func (_FlagsInterface *FlagsInterfaceTransactorSession) RaiseFlag(arg0 common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.RaiseFlag(&_FlagsInterface.TransactOpts, arg0)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] ) returns()
func (_FlagsInterface *FlagsInterfaceTransactor) RaiseFlags(opts *bind.TransactOpts, arg0 []common.Address) (*types.Transaction, error) {
	return _FlagsInterface.contract.Transact(opts, "raiseFlags", arg0)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] ) returns()
func (_FlagsInterface *FlagsInterfaceSession) RaiseFlags(arg0 []common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.RaiseFlags(&_FlagsInterface.TransactOpts, arg0)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] ) returns()
func (_FlagsInterface *FlagsInterfaceTransactorSession) RaiseFlags(arg0 []common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.RaiseFlags(&_FlagsInterface.TransactOpts, arg0)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address ) returns()
func (_FlagsInterface *FlagsInterfaceTransactor) SetRaisingAccessController(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _FlagsInterface.contract.Transact(opts, "setRaisingAccessController", arg0)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address ) returns()
func (_FlagsInterface *FlagsInterfaceSession) SetRaisingAccessController(arg0 common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.SetRaisingAccessController(&_FlagsInterface.TransactOpts, arg0)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address ) returns()
func (_FlagsInterface *FlagsInterfaceTransactorSession) SetRaisingAccessController(arg0 common.Address) (*types.Transaction, error) {
	return _FlagsInterface.Contract.SetRaisingAccessController(&_FlagsInterface.TransactOpts, arg0)
}
