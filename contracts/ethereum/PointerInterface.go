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

// PointerInterfaceABI is the input ABI used to generate the binding from.
const PointerInterfaceABI = "[{\"inputs\":[],\"name\":\"getAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PointerInterface is an auto generated Go binding around an Ethereum contract.
type PointerInterface struct {
	PointerInterfaceCaller     // Read-only binding to the contract
	PointerInterfaceTransactor // Write-only binding to the contract
	PointerInterfaceFilterer   // Log filterer for contract events
}

// PointerInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type PointerInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PointerInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PointerInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PointerInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PointerInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PointerInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PointerInterfaceSession struct {
	Contract     *PointerInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PointerInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PointerInterfaceCallerSession struct {
	Contract *PointerInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// PointerInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PointerInterfaceTransactorSession struct {
	Contract     *PointerInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// PointerInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type PointerInterfaceRaw struct {
	Contract *PointerInterface // Generic contract binding to access the raw methods on
}

// PointerInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PointerInterfaceCallerRaw struct {
	Contract *PointerInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// PointerInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PointerInterfaceTransactorRaw struct {
	Contract *PointerInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPointerInterface creates a new instance of PointerInterface, bound to a specific deployed contract.
func NewPointerInterface(address common.Address, backend bind.ContractBackend) (*PointerInterface, error) {
	contract, err := bindPointerInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PointerInterface{PointerInterfaceCaller: PointerInterfaceCaller{contract: contract}, PointerInterfaceTransactor: PointerInterfaceTransactor{contract: contract}, PointerInterfaceFilterer: PointerInterfaceFilterer{contract: contract}}, nil
}

// NewPointerInterfaceCaller creates a new read-only instance of PointerInterface, bound to a specific deployed contract.
func NewPointerInterfaceCaller(address common.Address, caller bind.ContractCaller) (*PointerInterfaceCaller, error) {
	contract, err := bindPointerInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PointerInterfaceCaller{contract: contract}, nil
}

// NewPointerInterfaceTransactor creates a new write-only instance of PointerInterface, bound to a specific deployed contract.
func NewPointerInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*PointerInterfaceTransactor, error) {
	contract, err := bindPointerInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PointerInterfaceTransactor{contract: contract}, nil
}

// NewPointerInterfaceFilterer creates a new log filterer instance of PointerInterface, bound to a specific deployed contract.
func NewPointerInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*PointerInterfaceFilterer, error) {
	contract, err := bindPointerInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PointerInterfaceFilterer{contract: contract}, nil
}

// bindPointerInterface binds a generic wrapper to an already deployed contract.
func bindPointerInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PointerInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PointerInterface *PointerInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PointerInterface.Contract.PointerInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PointerInterface *PointerInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PointerInterface.Contract.PointerInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PointerInterface *PointerInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PointerInterface.Contract.PointerInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PointerInterface *PointerInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PointerInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PointerInterface *PointerInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PointerInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PointerInterface *PointerInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PointerInterface.Contract.contract.Transact(opts, method, params...)
}

// GetAddress is a free data retrieval call binding the contract method 0x38cc4831.
//
// Solidity: function getAddress() view returns(address)
func (_PointerInterface *PointerInterfaceCaller) GetAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PointerInterface.contract.Call(opts, &out, "getAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress is a free data retrieval call binding the contract method 0x38cc4831.
//
// Solidity: function getAddress() view returns(address)
func (_PointerInterface *PointerInterfaceSession) GetAddress() (common.Address, error) {
	return _PointerInterface.Contract.GetAddress(&_PointerInterface.CallOpts)
}

// GetAddress is a free data retrieval call binding the contract method 0x38cc4831.
//
// Solidity: function getAddress() view returns(address)
func (_PointerInterface *PointerInterfaceCallerSession) GetAddress() (common.Address, error) {
	return _PointerInterface.Contract.GetAddress(&_PointerInterface.CallOpts)
}
