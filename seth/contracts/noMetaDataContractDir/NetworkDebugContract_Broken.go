// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package network_debug_contract

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

// NetworkDebugContractAccount is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugContractAccount struct {
	Name       string
	Balance    uint64
	DailyLimit *big.Int
}

// NetworkDebugContractData is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugContractData struct {
	Name   string
	Values []*big.Int
}

// NetworkDebugContractNestedData is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugContractNestedData struct {
	Data         NetworkDebugContractData
	DynamicBytes []byte
}


// NetworkDebugContractABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkDebugContractMetaData.ABI instead.
var NetworkDebugContractABI = abi.ABI{}

// NetworkDebugContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkDebugContractMetaData.Bin instead.
var NetworkDebugContractBin = ""

// DeployNetworkDebugContract deploys a new Ethereum contract, binding an instance of NetworkDebugContract to it.
func DeployNetworkDebugContract(auth *bind.TransactOpts, backend bind.ContractBackend, subAddr common.Address) (common.Address, *types.Transaction, *NetworkDebugContract, error) {
	parsed, err := func() (*abi.ABI, error) {
		return &abi.ABI{}, errors.New("This is a broken contract. Don't try to deploy it")
	}()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkDebugContractBin), backend, subAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkDebugContract{NetworkDebugContractCaller: NetworkDebugContractCaller{contract: contract}, NetworkDebugContractTransactor: NetworkDebugContractTransactor{contract: contract}, NetworkDebugContractFilterer: NetworkDebugContractFilterer{contract: contract}}, nil
}

// NetworkDebugContract is an auto generated Go binding around an Ethereum contract.
type NetworkDebugContract struct {
	NetworkDebugContractCaller     // Read-only binding to the contract
	NetworkDebugContractTransactor // Write-only binding to the contract
	NetworkDebugContractFilterer   // Log filterer for contract events
}

// NetworkDebugContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkDebugContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkDebugContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkDebugContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkDebugContractSession struct {
	Contract     *NetworkDebugContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// NetworkDebugContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkDebugContractCallerSession struct {
	Contract *NetworkDebugContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// NetworkDebugContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkDebugContractTransactorSession struct {
	Contract     *NetworkDebugContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// NetworkDebugContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkDebugContractRaw struct {
	Contract *NetworkDebugContract // Generic contract binding to access the raw methods on
}

// NetworkDebugContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkDebugContractCallerRaw struct {
	Contract *NetworkDebugContractCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkDebugContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkDebugContractTransactorRaw struct {
	Contract *NetworkDebugContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkDebugContract creates a new instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContract(address common.Address, backend bind.ContractBackend) (*NetworkDebugContract, error) {
	contract, err := bindNetworkDebugContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContract{NetworkDebugContractCaller: NetworkDebugContractCaller{contract: contract}, NetworkDebugContractTransactor: NetworkDebugContractTransactor{contract: contract}, NetworkDebugContractFilterer: NetworkDebugContractFilterer{contract: contract}}, nil
}

// NewNetworkDebugContractCaller creates a new read-only instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContractCaller(address common.Address, caller bind.ContractCaller) (*NetworkDebugContractCaller, error) {
	contract, err := bindNetworkDebugContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCaller{contract: contract}, nil
}

// NewNetworkDebugContractTransactor creates a new write-only instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContractTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkDebugContractTransactor, error) {
	contract, err := bindNetworkDebugContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractTransactor{contract: contract}, nil
}

// NewNetworkDebugContractFilterer creates a new log filterer instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContractFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkDebugContractFilterer, error) {
	contract, err := bindNetworkDebugContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractFilterer{contract: contract}, nil
}

// bindNetworkDebugContract binds a generic wrapper to an already deployed contract.
func bindNetworkDebugContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := func() (*abi.ABI, error) {
		return &abi.ABI{}, errors.New("This is a broken contract. Don't try to deploy it")
	}()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkDebugContract *NetworkDebugContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkDebugContract.Contract.NetworkDebugContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkDebugContract *NetworkDebugContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.NetworkDebugContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkDebugContract *NetworkDebugContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.NetworkDebugContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkDebugContract *NetworkDebugContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkDebugContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkDebugContract *NetworkDebugContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkDebugContract *NetworkDebugContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.contract.Transact(opts, method, params...)
}

// CounterMap is a free data retrieval call binding the contract method 0x6284117d.
//
// Solidity: function counterMap(int256 ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCaller) CounterMap(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "counterMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CounterMap is a free data retrieval call binding the contract method 0x6284117d.
//
// Solidity: function counterMap(int256 ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) CounterMap(arg0 *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.CounterMap(&_NetworkDebugContract.CallOpts, arg0)
}

// CounterMap is a free data retrieval call binding the contract method 0x6284117d.
//
// Solidity: function counterMap(int256 ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) CounterMap(arg0 *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.CounterMap(&_NetworkDebugContract.CallOpts, arg0)
}

// CurrentStatus is a free data retrieval call binding the contract method 0xef8a9235.
//
// Solidity: function currentStatus() view returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractCaller) CurrentStatus(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "currentStatus")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CurrentStatus is a free data retrieval call binding the contract method 0xef8a9235.
//
// Solidity: function currentStatus() view returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractSession) CurrentStatus() (uint8, error) {
	return _NetworkDebugContract.Contract.CurrentStatus(&_NetworkDebugContract.CallOpts)
}

// CurrentStatus is a free data retrieval call binding the contract method 0xef8a9235.
//
// Solidity: function currentStatus() view returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) CurrentStatus() (uint8, error) {
	return _NetworkDebugContract.Contract.CurrentStatus(&_NetworkDebugContract.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCaller) Get(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "get")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractSession) Get() (*big.Int, error) {
	return _NetworkDebugContract.Contract.Get(&_NetworkDebugContract.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) Get() (*big.Int, error) {
	return _NetworkDebugContract.Contract.Get(&_NetworkDebugContract.CallOpts)
}

// GetCounter is a free data retrieval call binding the contract method 0x5921483f.
//
// Solidity: function getCounter(int256 idx) view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCaller) GetCounter(opts *bind.CallOpts, idx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "getCounter", idx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCounter is a free data retrieval call binding the contract method 0x5921483f.
//
// Solidity: function getCounter(int256 idx) view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractSession) GetCounter(idx *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetCounter(&_NetworkDebugContract.CallOpts, idx)
}

// GetCounter is a free data retrieval call binding the contract method 0x5921483f.
//
// Solidity: function getCounter(int256 idx) view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) GetCounter(idx *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetCounter(&_NetworkDebugContract.CallOpts, idx)
}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCaller) GetData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "getData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractSession) GetData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetData(&_NetworkDebugContract.CallOpts)
}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) GetData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetData(&_NetworkDebugContract.CallOpts)
}

// GetMap is a free data retrieval call binding the contract method 0xad3de14c.
//
// Solidity: function getMap() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCaller) GetMap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "getMap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMap is a free data retrieval call binding the contract method 0xad3de14c.
//
// Solidity: function getMap() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractSession) GetMap() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetMap(&_NetworkDebugContract.CallOpts)
}

// GetMap is a free data retrieval call binding the contract method 0xad3de14c.
//
// Solidity: function getMap() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) GetMap() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetMap(&_NetworkDebugContract.CallOpts)
}

// PerformStaticCall is a free data retrieval call binding the contract method 0x3170428e.
//
// Solidity: function performStaticCall() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCaller) PerformStaticCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "performStaticCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerformStaticCall is a free data retrieval call binding the contract method 0x3170428e.
//
// Solidity: function performStaticCall() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractSession) PerformStaticCall() (*big.Int, error) {
	return _NetworkDebugContract.Contract.PerformStaticCall(&_NetworkDebugContract.CallOpts)
}

// PerformStaticCall is a free data retrieval call binding the contract method 0x3170428e.
//
// Solidity: function performStaticCall() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) PerformStaticCall() (*big.Int, error) {
	return _NetworkDebugContract.Contract.PerformStaticCall(&_NetworkDebugContract.CallOpts)
}

// StoredData is a free data retrieval call binding the contract method 0x2a1afcd9.
//
// Solidity: function storedData() view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCaller) StoredData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "storedData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StoredData is a free data retrieval call binding the contract method 0x2a1afcd9.
//
// Solidity: function storedData() view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) StoredData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredData(&_NetworkDebugContract.CallOpts)
}

// StoredData is a free data retrieval call binding the contract method 0x2a1afcd9.
//
// Solidity: function storedData() view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) StoredData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredData(&_NetworkDebugContract.CallOpts)
}

// StoredDataMap is a free data retrieval call binding the contract method 0x48ad9fe8.
//
// Solidity: function storedDataMap(address ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCaller) StoredDataMap(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "storedDataMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StoredDataMap is a free data retrieval call binding the contract method 0x48ad9fe8.
//
// Solidity: function storedDataMap(address ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) StoredDataMap(arg0 common.Address) (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredDataMap(&_NetworkDebugContract.CallOpts, arg0)
}

// StoredDataMap is a free data retrieval call binding the contract method 0x48ad9fe8.
//
// Solidity: function storedDataMap(address ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) StoredDataMap(arg0 common.Address) (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredDataMap(&_NetworkDebugContract.CallOpts, arg0)
}

// SubContract is a free data retrieval call binding the contract method 0xc0d06d89.
//
// Solidity: function subContract() view returns(address)
func (_NetworkDebugContract *NetworkDebugContractCaller) SubContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "subContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SubContract is a free data retrieval call binding the contract method 0xc0d06d89.
//
// Solidity: function subContract() view returns(address)
func (_NetworkDebugContract *NetworkDebugContractSession) SubContract() (common.Address, error) {
	return _NetworkDebugContract.Contract.SubContract(&_NetworkDebugContract.CallOpts)
}

// SubContract is a free data retrieval call binding the contract method 0xc0d06d89.
//
// Solidity: function subContract() view returns(address)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) SubContract() (common.Address, error) {
	return _NetworkDebugContract.Contract.SubContract(&_NetworkDebugContract.CallOpts)
}

// AddCounter is a paid mutator transaction binding the contract method 0x23515760.
//
// Solidity: function addCounter(int256 idx, int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactor) AddCounter(opts *bind.TransactOpts, idx *big.Int, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "addCounter", idx, x)
}

// AddCounter is a paid mutator transaction binding the contract method 0x23515760.
//
// Solidity: function addCounter(int256 idx, int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractSession) AddCounter(idx *big.Int, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AddCounter(&_NetworkDebugContract.TransactOpts, idx, x)
}

// AddCounter is a paid mutator transaction binding the contract method 0x23515760.
//
// Solidity: function addCounter(int256 idx, int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AddCounter(idx *big.Int, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AddCounter(&_NetworkDebugContract.TransactOpts, idx, x)
}

// AlwaysRevertsAssert is a paid mutator transaction binding the contract method 0x256560d5.
//
// Solidity: function alwaysRevertsAssert() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsAssert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsAssert")
}

// AlwaysRevertsAssert is a paid mutator transaction binding the contract method 0x256560d5.
//
// Solidity: function alwaysRevertsAssert() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsAssert() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsAssert(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsAssert is a paid mutator transaction binding the contract method 0x256560d5.
//
// Solidity: function alwaysRevertsAssert() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsAssert() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsAssert(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x5e9c80d6.
//
// Solidity: function alwaysRevertsCustomError() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsCustomError(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsCustomError")
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x5e9c80d6.
//
// Solidity: function alwaysRevertsCustomError() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsCustomError() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomError(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x5e9c80d6.
//
// Solidity: function alwaysRevertsCustomError() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsCustomError() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomError(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomErrorNoValues is a paid mutator transaction binding the contract method 0xb600141f.
//
// Solidity: function alwaysRevertsCustomErrorNoValues() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsCustomErrorNoValues(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsCustomErrorNoValues")
}

// AlwaysRevertsCustomErrorNoValues is a paid mutator transaction binding the contract method 0xb600141f.
//
// Solidity: function alwaysRevertsCustomErrorNoValues() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsCustomErrorNoValues() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomErrorNoValues(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomErrorNoValues is a paid mutator transaction binding the contract method 0xb600141f.
//
// Solidity: function alwaysRevertsCustomErrorNoValues() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsCustomErrorNoValues() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomErrorNoValues(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsRequire is a paid mutator transaction binding the contract method 0x06595f75.
//
// Solidity: function alwaysRevertsRequire() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsRequire(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsRequire")
}

// AlwaysRevertsRequire is a paid mutator transaction binding the contract method 0x06595f75.
//
// Solidity: function alwaysRevertsRequire() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsRequire() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsRequire(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsRequire is a paid mutator transaction binding the contract method 0x06595f75.
//
// Solidity: function alwaysRevertsRequire() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsRequire() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsRequire(&_NetworkDebugContract.TransactOpts)
}

// CallRevertFunctionInSubContract is a paid mutator transaction binding the contract method 0x11b3c478.
//
// Solidity: function callRevertFunctionInSubContract(uint256 x, uint256 y) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) CallRevertFunctionInSubContract(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "callRevertFunctionInSubContract", x, y)
}

// CallRevertFunctionInSubContract is a paid mutator transaction binding the contract method 0x11b3c478.
//
// Solidity: function callRevertFunctionInSubContract(uint256 x, uint256 y) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) CallRevertFunctionInSubContract(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInSubContract(&_NetworkDebugContract.TransactOpts, x, y)
}

// CallRevertFunctionInSubContract is a paid mutator transaction binding the contract method 0x11b3c478.
//
// Solidity: function callRevertFunctionInSubContract(uint256 x, uint256 y) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) CallRevertFunctionInSubContract(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInSubContract(&_NetworkDebugContract.TransactOpts, x, y)
}

// CallRevertFunctionInTheContract is a paid mutator transaction binding the contract method 0x9349d00b.
//
// Solidity: function callRevertFunctionInTheContract() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) CallRevertFunctionInTheContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "callRevertFunctionInTheContract")
}

// CallRevertFunctionInTheContract is a paid mutator transaction binding the contract method 0x9349d00b.
//
// Solidity: function callRevertFunctionInTheContract() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) CallRevertFunctionInTheContract() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInTheContract(&_NetworkDebugContract.TransactOpts)
}

// CallRevertFunctionInTheContract is a paid mutator transaction binding the contract method 0x9349d00b.
//
// Solidity: function callRevertFunctionInTheContract() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) CallRevertFunctionInTheContract() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInTheContract(&_NetworkDebugContract.TransactOpts)
}

// CallbackMethod is a paid mutator transaction binding the contract method 0xfbcb8d07.
//
// Solidity: function callbackMethod(int256 x) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) CallbackMethod(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "callbackMethod", x)
}

// CallbackMethod is a paid mutator transaction binding the contract method 0xfbcb8d07.
//
// Solidity: function callbackMethod(int256 x) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) CallbackMethod(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallbackMethod(&_NetworkDebugContract.TransactOpts, x)
}

// CallbackMethod is a paid mutator transaction binding the contract method 0xfbcb8d07.
//
// Solidity: function callbackMethod(int256 x) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) CallbackMethod(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallbackMethod(&_NetworkDebugContract.TransactOpts, x)
}

// EmitAddress is a paid mutator transaction binding the contract method 0xec5c3ede.
//
// Solidity: function emitAddress(address addr) returns(address)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitAddress", addr)
}

// EmitAddress is a paid mutator transaction binding the contract method 0xec5c3ede.
//
// Solidity: function emitAddress(address addr) returns(address)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitAddress(addr common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitAddress(&_NetworkDebugContract.TransactOpts, addr)
}

// EmitAddress is a paid mutator transaction binding the contract method 0xec5c3ede.
//
// Solidity: function emitAddress(address addr) returns(address)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitAddress(addr common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitAddress(&_NetworkDebugContract.TransactOpts, addr)
}

// EmitBytes32 is a paid mutator transaction binding the contract method 0x33311ef3.
//
// Solidity: function emitBytes32(bytes32 input) returns(bytes32 output)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitBytes32(opts *bind.TransactOpts, input [32]byte) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitBytes32", input)
}

// EmitBytes32 is a paid mutator transaction binding the contract method 0x33311ef3.
//
// Solidity: function emitBytes32(bytes32 input) returns(bytes32 output)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitBytes32(input [32]byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitBytes32(&_NetworkDebugContract.TransactOpts, input)
}

// EmitBytes32 is a paid mutator transaction binding the contract method 0x33311ef3.
//
// Solidity: function emitBytes32(bytes32 input) returns(bytes32 output)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitBytes32(input [32]byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitBytes32(&_NetworkDebugContract.TransactOpts, input)
}

// EmitFourParamMixedEvent is a paid mutator transaction binding the contract method 0xc2124b22.
//
// Solidity: function emitFourParamMixedEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitFourParamMixedEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitFourParamMixedEvent")
}

// EmitFourParamMixedEvent is a paid mutator transaction binding the contract method 0xc2124b22.
//
// Solidity: function emitFourParamMixedEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitFourParamMixedEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitFourParamMixedEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitFourParamMixedEvent is a paid mutator transaction binding the contract method 0xc2124b22.
//
// Solidity: function emitFourParamMixedEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitFourParamMixedEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitFourParamMixedEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitInputs is a paid mutator transaction binding the contract method 0x81b375a0.
//
// Solidity: function emitInputs(uint256 inputVal1, string inputVal2) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitInputs(opts *bind.TransactOpts, inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitInputs", inputVal1, inputVal2)
}

// EmitInputs is a paid mutator transaction binding the contract method 0x81b375a0.
//
// Solidity: function emitInputs(uint256 inputVal1, string inputVal2) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitInputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInputs is a paid mutator transaction binding the contract method 0x81b375a0.
//
// Solidity: function emitInputs(uint256 inputVal1, string inputVal2) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitInputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInputsOutputs is a paid mutator transaction binding the contract method 0xd7a80205.
//
// Solidity: function emitInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitInputsOutputs(opts *bind.TransactOpts, inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitInputsOutputs", inputVal1, inputVal2)
}

// EmitInputsOutputs is a paid mutator transaction binding the contract method 0xd7a80205.
//
// Solidity: function emitInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInputsOutputs is a paid mutator transaction binding the contract method 0xd7a80205.
//
// Solidity: function emitInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInts is a paid mutator transaction binding the contract method 0x9e099652.
//
// Solidity: function emitInts(int256 first, int128 second, uint256 third) returns(int256, int128 outputVal1, uint256 outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitInts(opts *bind.TransactOpts, first *big.Int, second *big.Int, third *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitInts", first, second, third)
}

// EmitInts is a paid mutator transaction binding the contract method 0x9e099652.
//
// Solidity: function emitInts(int256 first, int128 second, uint256 third) returns(int256, int128 outputVal1, uint256 outputVal2)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitInts(first *big.Int, second *big.Int, third *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInts(&_NetworkDebugContract.TransactOpts, first, second, third)
}

// EmitInts is a paid mutator transaction binding the contract method 0x9e099652.
//
// Solidity: function emitInts(int256 first, int128 second, uint256 third) returns(int256, int128 outputVal1, uint256 outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitInts(first *big.Int, second *big.Int, third *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInts(&_NetworkDebugContract.TransactOpts, first, second, third)
}

// EmitNamedInputsOutputs is a paid mutator transaction binding the contract method 0x45f0c9e6.
//
// Solidity: function emitNamedInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNamedInputsOutputs(opts *bind.TransactOpts, inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNamedInputsOutputs", inputVal1, inputVal2)
}

// EmitNamedInputsOutputs is a paid mutator transaction binding the contract method 0x45f0c9e6.
//
// Solidity: function emitNamedInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNamedInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitNamedInputsOutputs is a paid mutator transaction binding the contract method 0x45f0c9e6.
//
// Solidity: function emitNamedInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNamedInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitNamedOutputs is a paid mutator transaction binding the contract method 0x7014c81d.
//
// Solidity: function emitNamedOutputs() returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNamedOutputs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNamedOutputs")
}

// EmitNamedOutputs is a paid mutator transaction binding the contract method 0x7014c81d.
//
// Solidity: function emitNamedOutputs() returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNamedOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitNamedOutputs is a paid mutator transaction binding the contract method 0x7014c81d.
//
// Solidity: function emitNamedOutputs() returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNamedOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEvent is a paid mutator transaction binding the contract method 0x95a81a4c.
//
// Solidity: function emitNoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNoIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNoIndexEvent")
}

// EmitNoIndexEvent is a paid mutator transaction binding the contract method 0x95a81a4c.
//
// Solidity: function emitNoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEvent is a paid mutator transaction binding the contract method 0x95a81a4c.
//
// Solidity: function emitNoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEventString is a paid mutator transaction binding the contract method 0x788c4772.
//
// Solidity: function emitNoIndexEventString() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNoIndexEventString(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNoIndexEventString")
}

// EmitNoIndexEventString is a paid mutator transaction binding the contract method 0x788c4772.
//
// Solidity: function emitNoIndexEventString() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNoIndexEventString() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEventString(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEventString is a paid mutator transaction binding the contract method 0x788c4772.
//
// Solidity: function emitNoIndexEventString() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNoIndexEventString() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEventString(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexStructEvent is a paid mutator transaction binding the contract method 0x62c270e1.
//
// Solidity: function emitNoIndexStructEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNoIndexStructEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNoIndexStructEvent")
}

// EmitNoIndexStructEvent is a paid mutator transaction binding the contract method 0x62c270e1.
//
// Solidity: function emitNoIndexStructEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNoIndexStructEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexStructEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexStructEvent is a paid mutator transaction binding the contract method 0x62c270e1.
//
// Solidity: function emitNoIndexStructEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNoIndexStructEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexStructEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitOneIndexEvent is a paid mutator transaction binding the contract method 0x8f856296.
//
// Solidity: function emitOneIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitOneIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitOneIndexEvent")
}

// EmitOneIndexEvent is a paid mutator transaction binding the contract method 0x8f856296.
//
// Solidity: function emitOneIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitOneIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOneIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitOneIndexEvent is a paid mutator transaction binding the contract method 0x8f856296.
//
// Solidity: function emitOneIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitOneIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOneIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitOutputs is a paid mutator transaction binding the contract method 0x8db611be.
//
// Solidity: function emitOutputs() returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitOutputs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitOutputs")
}

// EmitOutputs is a paid mutator transaction binding the contract method 0x8db611be.
//
// Solidity: function emitOutputs() returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitOutputs is a paid mutator transaction binding the contract method 0x8db611be.
//
// Solidity: function emitOutputs() returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitThreeIndexEvent is a paid mutator transaction binding the contract method 0xaa3fdcf4.
//
// Solidity: function emitThreeIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitThreeIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitThreeIndexEvent")
}

// EmitThreeIndexEvent is a paid mutator transaction binding the contract method 0xaa3fdcf4.
//
// Solidity: function emitThreeIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitThreeIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitThreeIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitThreeIndexEvent is a paid mutator transaction binding the contract method 0xaa3fdcf4.
//
// Solidity: function emitThreeIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitThreeIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitThreeIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitTwoIndexEvent is a paid mutator transaction binding the contract method 0x1e31d0a8.
//
// Solidity: function emitTwoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitTwoIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitTwoIndexEvent")
}

// EmitTwoIndexEvent is a paid mutator transaction binding the contract method 0x1e31d0a8.
//
// Solidity: function emitTwoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitTwoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitTwoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitTwoIndexEvent is a paid mutator transaction binding the contract method 0x1e31d0a8.
//
// Solidity: function emitTwoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitTwoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitTwoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.OnTokenTransfer(&_NetworkDebugContract.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.OnTokenTransfer(&_NetworkDebugContract.TransactOpts, sender, amount, data)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) Pay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "pay")
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugContract *NetworkDebugContractSession) Pay() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Pay(&_NetworkDebugContract.TransactOpts)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Pay() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Pay(&_NetworkDebugContract.TransactOpts)
}

// ProcessAddressArray is a paid mutator transaction binding the contract method 0xe1111f79.
//
// Solidity: function processAddressArray(address[] input) returns(address[])
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessAddressArray(opts *bind.TransactOpts, input []common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processAddressArray", input)
}

// ProcessAddressArray is a paid mutator transaction binding the contract method 0xe1111f79.
//
// Solidity: function processAddressArray(address[] input) returns(address[])
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessAddressArray(input []common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessAddressArray(&_NetworkDebugContract.TransactOpts, input)
}

// ProcessAddressArray is a paid mutator transaction binding the contract method 0xe1111f79.
//
// Solidity: function processAddressArray(address[] input) returns(address[])
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessAddressArray(input []common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessAddressArray(&_NetworkDebugContract.TransactOpts, input)
}

// ProcessDynamicData is a paid mutator transaction binding the contract method 0x7fdc8fe1.
//
// Solidity: function processDynamicData((string,uint256[]) data) returns((string,uint256[]))
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessDynamicData(opts *bind.TransactOpts, data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processDynamicData", data)
}

// ProcessDynamicData is a paid mutator transaction binding the contract method 0x7fdc8fe1.
//
// Solidity: function processDynamicData((string,uint256[]) data) returns((string,uint256[]))
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessDynamicData(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessDynamicData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessDynamicData is a paid mutator transaction binding the contract method 0x7fdc8fe1.
//
// Solidity: function processDynamicData((string,uint256[]) data) returns((string,uint256[]))
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessDynamicData(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessDynamicData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessFixedDataArray is a paid mutator transaction binding the contract method 0x99adad2e.
//
// Solidity: function processFixedDataArray((string,uint256[])[3] data) returns((string,uint256[])[2])
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessFixedDataArray(opts *bind.TransactOpts, data [3]NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processFixedDataArray", data)
}

// ProcessFixedDataArray is a paid mutator transaction binding the contract method 0x99adad2e.
//
// Solidity: function processFixedDataArray((string,uint256[])[3] data) returns((string,uint256[])[2])
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessFixedDataArray(data [3]NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessFixedDataArray(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessFixedDataArray is a paid mutator transaction binding the contract method 0x99adad2e.
//
// Solidity: function processFixedDataArray((string,uint256[])[3] data) returns((string,uint256[])[2])
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessFixedDataArray(data [3]NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessFixedDataArray(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData is a paid mutator transaction binding the contract method 0x7f12881c.
//
// Solidity: function processNestedData(((string,uint256[]),bytes) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessNestedData(opts *bind.TransactOpts, data NetworkDebugContractNestedData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processNestedData", data)
}

// ProcessNestedData is a paid mutator transaction binding the contract method 0x7f12881c.
//
// Solidity: function processNestedData(((string,uint256[]),bytes) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessNestedData(data NetworkDebugContractNestedData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData is a paid mutator transaction binding the contract method 0x7f12881c.
//
// Solidity: function processNestedData(((string,uint256[]),bytes) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessNestedData(data NetworkDebugContractNestedData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData0 is a paid mutator transaction binding the contract method 0xf499af2a.
//
// Solidity: function processNestedData((string,uint256[]) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessNestedData0(opts *bind.TransactOpts, data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processNestedData0", data)
}

// ProcessNestedData0 is a paid mutator transaction binding the contract method 0xf499af2a.
//
// Solidity: function processNestedData((string,uint256[]) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessNestedData0(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData0(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData0 is a paid mutator transaction binding the contract method 0xf499af2a.
//
// Solidity: function processNestedData((string,uint256[]) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessNestedData0(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData0(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessUintArray is a paid mutator transaction binding the contract method 0x12d91233.
//
// Solidity: function processUintArray(uint256[] input) returns(uint256[])
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessUintArray(opts *bind.TransactOpts, input []*big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processUintArray", input)
}

// ProcessUintArray is a paid mutator transaction binding the contract method 0x12d91233.
//
// Solidity: function processUintArray(uint256[] input) returns(uint256[])
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessUintArray(input []*big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessUintArray(&_NetworkDebugContract.TransactOpts, input)
}

// ProcessUintArray is a paid mutator transaction binding the contract method 0x12d91233.
//
// Solidity: function processUintArray(uint256[] input) returns(uint256[])
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessUintArray(input []*big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessUintArray(&_NetworkDebugContract.TransactOpts, input)
}

// ResetCounter is a paid mutator transaction binding the contract method 0xf3396bd9.
//
// Solidity: function resetCounter(int256 idx) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) ResetCounter(opts *bind.TransactOpts, idx *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "resetCounter", idx)
}

// ResetCounter is a paid mutator transaction binding the contract method 0xf3396bd9.
//
// Solidity: function resetCounter(int256 idx) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) ResetCounter(idx *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ResetCounter(&_NetworkDebugContract.TransactOpts, idx)
}

// ResetCounter is a paid mutator transaction binding the contract method 0xf3396bd9.
//
// Solidity: function resetCounter(int256 idx) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ResetCounter(idx *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ResetCounter(&_NetworkDebugContract.TransactOpts, idx)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactor) Set(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "set", x)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractSession) Set(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Set(&_NetworkDebugContract.TransactOpts, x)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Set(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Set(&_NetworkDebugContract.TransactOpts, x)
}

// SetMap is a paid mutator transaction binding the contract method 0xe8116e28.
//
// Solidity: function setMap(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactor) SetMap(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "setMap", x)
}

// SetMap is a paid mutator transaction binding the contract method 0xe8116e28.
//
// Solidity: function setMap(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractSession) SetMap(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetMap(&_NetworkDebugContract.TransactOpts, x)
}

// SetMap is a paid mutator transaction binding the contract method 0xe8116e28.
//
// Solidity: function setMap(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) SetMap(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetMap(&_NetworkDebugContract.TransactOpts, x)
}

// SetStatus is a paid mutator transaction binding the contract method 0x2e49d78b.
//
// Solidity: function setStatus(uint8 status) returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractTransactor) SetStatus(opts *bind.TransactOpts, status uint8) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "setStatus", status)
}

// SetStatus is a paid mutator transaction binding the contract method 0x2e49d78b.
//
// Solidity: function setStatus(uint8 status) returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractSession) SetStatus(status uint8) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetStatus(&_NetworkDebugContract.TransactOpts, status)
}

// SetStatus is a paid mutator transaction binding the contract method 0x2e49d78b.
//
// Solidity: function setStatus(uint8 status) returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) SetStatus(status uint8) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetStatus(&_NetworkDebugContract.TransactOpts, status)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) Trace(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "trace", x, y)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) Trace(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Trace(&_NetworkDebugContract.TransactOpts, x, y)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Trace(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Trace(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceDifferent is a paid mutator transaction binding the contract method 0x30985bcc.
//
// Solidity: function traceDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceDifferent(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceDifferent", x, y)
}

// TraceDifferent is a paid mutator transaction binding the contract method 0x30985bcc.
//
// Solidity: function traceDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceDifferent is a paid mutator transaction binding the contract method 0x30985bcc.
//
// Solidity: function traceDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceSubWithCallback is a paid mutator transaction binding the contract method 0x3837a75e.
//
// Solidity: function traceSubWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceSubWithCallback(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceSubWithCallback", x, y)
}

// TraceSubWithCallback is a paid mutator transaction binding the contract method 0x3837a75e.
//
// Solidity: function traceSubWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceSubWithCallback(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceSubWithCallback(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceSubWithCallback is a paid mutator transaction binding the contract method 0x3837a75e.
//
// Solidity: function traceSubWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceSubWithCallback(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceSubWithCallback(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceWithValidate is a paid mutator transaction binding the contract method 0xb1ae9d85.
//
// Solidity: function traceWithValidate(int256 x, int256 y) payable returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceWithValidate(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceWithValidate", x, y)
}

// TraceWithValidate is a paid mutator transaction binding the contract method 0xb1ae9d85.
//
// Solidity: function traceWithValidate(int256 x, int256 y) payable returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceWithValidate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceWithValidate(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceWithValidate is a paid mutator transaction binding the contract method 0xb1ae9d85.
//
// Solidity: function traceWithValidate(int256 x, int256 y) payable returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceWithValidate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceWithValidate(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceYetDifferent is a paid mutator transaction binding the contract method 0x58379d71.
//
// Solidity: function traceYetDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceYetDifferent(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceYetDifferent", x, y)
}

// TraceYetDifferent is a paid mutator transaction binding the contract method 0x58379d71.
//
// Solidity: function traceYetDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceYetDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceYetDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceYetDifferent is a paid mutator transaction binding the contract method 0x58379d71.
//
// Solidity: function traceYetDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceYetDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceYetDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// Validate is a paid mutator transaction binding the contract method 0x04d8215b.
//
// Solidity: function validate(int256 x, int256 y) returns(bool)
func (_NetworkDebugContract *NetworkDebugContractTransactor) Validate(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "validate", x, y)
}

// Validate is a paid mutator transaction binding the contract method 0x04d8215b.
//
// Solidity: function validate(int256 x, int256 y) returns(bool)
func (_NetworkDebugContract *NetworkDebugContractSession) Validate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Validate(&_NetworkDebugContract.TransactOpts, x, y)
}

// Validate is a paid mutator transaction binding the contract method 0x04d8215b.
//
// Solidity: function validate(int256 x, int256 y) returns(bool)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Validate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Validate(&_NetworkDebugContract.TransactOpts, x, y)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_NetworkDebugContract *NetworkDebugContractSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Fallback(&_NetworkDebugContract.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Fallback(&_NetworkDebugContract.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkDebugContract *NetworkDebugContractSession) Receive() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Receive(&_NetworkDebugContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Receive() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Receive(&_NetworkDebugContract.TransactOpts)
}

// NetworkDebugContractCallDataLengthIterator is returned from FilterCallDataLength and is used to iterate over the raw logs and unpacked data for CallDataLength events raised by the NetworkDebugContract contract.
type NetworkDebugContractCallDataLengthIterator struct {
	Event *NetworkDebugContractCallDataLength // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractCallDataLengthIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractCallDataLength)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractCallDataLength)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractCallDataLengthIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractCallDataLengthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractCallDataLength represents a CallDataLength event raised by the NetworkDebugContract contract.
type NetworkDebugContractCallDataLength struct {
	Length *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCallDataLength is a free log retrieval operation binding the contract event 0x962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed.
//
// Solidity: event CallDataLength(uint256 length)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterCallDataLength(opts *bind.FilterOpts) (*NetworkDebugContractCallDataLengthIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "CallDataLength")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCallDataLengthIterator{contract: _NetworkDebugContract.contract, event: "CallDataLength", logs: logs, sub: sub}, nil
}

// WatchCallDataLength is a free log subscription operation binding the contract event 0x962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed.
//
// Solidity: event CallDataLength(uint256 length)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchCallDataLength(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractCallDataLength) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "CallDataLength")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractCallDataLength)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "CallDataLength", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCallDataLength is a log parse operation binding the contract event 0x962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed.
//
// Solidity: event CallDataLength(uint256 length)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseCallDataLength(log types.Log) (*NetworkDebugContractCallDataLength, error) {
	event := new(NetworkDebugContractCallDataLength)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "CallDataLength", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractCallbackEventIterator is returned from FilterCallbackEvent and is used to iterate over the raw logs and unpacked data for CallbackEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractCallbackEventIterator struct {
	Event *NetworkDebugContractCallbackEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractCallbackEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractCallbackEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractCallbackEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractCallbackEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractCallbackEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractCallbackEvent represents a CallbackEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractCallbackEvent struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCallbackEvent is a free log retrieval operation binding the contract event 0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68.
//
// Solidity: event CallbackEvent(int256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterCallbackEvent(opts *bind.FilterOpts, a []*big.Int) (*NetworkDebugContractCallbackEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "CallbackEvent", aRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCallbackEventIterator{contract: _NetworkDebugContract.contract, event: "CallbackEvent", logs: logs, sub: sub}, nil
}

// WatchCallbackEvent is a free log subscription operation binding the contract event 0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68.
//
// Solidity: event CallbackEvent(int256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchCallbackEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractCallbackEvent, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "CallbackEvent", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractCallbackEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "CallbackEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCallbackEvent is a log parse operation binding the contract event 0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68.
//
// Solidity: event CallbackEvent(int256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseCallbackEvent(log types.Log) (*NetworkDebugContractCallbackEvent, error) {
	event := new(NetworkDebugContractCallbackEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "CallbackEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractCurrentStatusIterator is returned from FilterCurrentStatus and is used to iterate over the raw logs and unpacked data for CurrentStatus events raised by the NetworkDebugContract contract.
type NetworkDebugContractCurrentStatusIterator struct {
	Event *NetworkDebugContractCurrentStatus // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractCurrentStatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractCurrentStatus)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractCurrentStatus)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractCurrentStatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractCurrentStatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractCurrentStatus represents a CurrentStatus event raised by the NetworkDebugContract contract.
type NetworkDebugContractCurrentStatus struct {
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCurrentStatus is a free log retrieval operation binding the contract event 0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a.
//
// Solidity: event CurrentStatus(uint8 indexed status)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterCurrentStatus(opts *bind.FilterOpts, status []uint8) (*NetworkDebugContractCurrentStatusIterator, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "CurrentStatus", statusRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCurrentStatusIterator{contract: _NetworkDebugContract.contract, event: "CurrentStatus", logs: logs, sub: sub}, nil
}

// WatchCurrentStatus is a free log subscription operation binding the contract event 0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a.
//
// Solidity: event CurrentStatus(uint8 indexed status)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchCurrentStatus(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractCurrentStatus, status []uint8) (event.Subscription, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "CurrentStatus", statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractCurrentStatus)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "CurrentStatus", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCurrentStatus is a log parse operation binding the contract event 0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a.
//
// Solidity: event CurrentStatus(uint8 indexed status)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseCurrentStatus(log types.Log) (*NetworkDebugContractCurrentStatus, error) {
	event := new(NetworkDebugContractCurrentStatus)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "CurrentStatus", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractEtherReceivedIterator is returned from FilterEtherReceived and is used to iterate over the raw logs and unpacked data for EtherReceived events raised by the NetworkDebugContract contract.
type NetworkDebugContractEtherReceivedIterator struct {
	Event *NetworkDebugContractEtherReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractEtherReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractEtherReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractEtherReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractEtherReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractEtherReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractEtherReceived represents a EtherReceived event raised by the NetworkDebugContract contract.
type NetworkDebugContractEtherReceived struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherReceived is a free log retrieval operation binding the contract event 0x1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b.
//
// Solidity: event EtherReceived(address sender, uint256 amount)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterEtherReceived(opts *bind.FilterOpts) (*NetworkDebugContractEtherReceivedIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "EtherReceived")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractEtherReceivedIterator{contract: _NetworkDebugContract.contract, event: "EtherReceived", logs: logs, sub: sub}, nil
}

// WatchEtherReceived is a free log subscription operation binding the contract event 0x1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b.
//
// Solidity: event EtherReceived(address sender, uint256 amount)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchEtherReceived(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractEtherReceived) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "EtherReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractEtherReceived)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "EtherReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEtherReceived is a log parse operation binding the contract event 0x1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b.
//
// Solidity: event EtherReceived(address sender, uint256 amount)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseEtherReceived(log types.Log) (*NetworkDebugContractEtherReceived, error) {
	event := new(NetworkDebugContractEtherReceived)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "EtherReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractIsValidEventIterator is returned from FilterIsValidEvent and is used to iterate over the raw logs and unpacked data for IsValidEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractIsValidEventIterator struct {
	Event *NetworkDebugContractIsValidEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractIsValidEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractIsValidEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractIsValidEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractIsValidEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractIsValidEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractIsValidEvent represents a IsValidEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractIsValidEvent struct {
	Success bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterIsValidEvent is a free log retrieval operation binding the contract event 0xdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b.
//
// Solidity: event IsValidEvent(bool success)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterIsValidEvent(opts *bind.FilterOpts) (*NetworkDebugContractIsValidEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "IsValidEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractIsValidEventIterator{contract: _NetworkDebugContract.contract, event: "IsValidEvent", logs: logs, sub: sub}, nil
}

// WatchIsValidEvent is a free log subscription operation binding the contract event 0xdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b.
//
// Solidity: event IsValidEvent(bool success)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchIsValidEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractIsValidEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "IsValidEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractIsValidEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "IsValidEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIsValidEvent is a log parse operation binding the contract event 0xdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b.
//
// Solidity: event IsValidEvent(bool success)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseIsValidEvent(log types.Log) (*NetworkDebugContractIsValidEvent, error) {
	event := new(NetworkDebugContractIsValidEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "IsValidEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractNoIndexEventIterator is returned from FilterNoIndexEvent and is used to iterate over the raw logs and unpacked data for NoIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEventIterator struct {
	Event *NetworkDebugContractNoIndexEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractNoIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractNoIndexEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractNoIndexEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractNoIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractNoIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractNoIndexEvent represents a NoIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEvent struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNoIndexEvent is a free log retrieval operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterNoIndexEvent(opts *bind.FilterOpts) (*NetworkDebugContractNoIndexEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "NoIndexEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractNoIndexEventIterator{contract: _NetworkDebugContract.contract, event: "NoIndexEvent", logs: logs, sub: sub}, nil
}

// WatchNoIndexEvent is a free log subscription operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchNoIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractNoIndexEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "NoIndexEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractNoIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNoIndexEvent is a log parse operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseNoIndexEvent(log types.Log) (*NetworkDebugContractNoIndexEvent, error) {
	event := new(NetworkDebugContractNoIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractNoIndexEventStringIterator is returned from FilterNoIndexEventString and is used to iterate over the raw logs and unpacked data for NoIndexEventString events raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEventStringIterator struct {
	Event *NetworkDebugContractNoIndexEventString // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractNoIndexEventStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractNoIndexEventString)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractNoIndexEventString)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractNoIndexEventStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractNoIndexEventStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractNoIndexEventString represents a NoIndexEventString event raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEventString struct {
	Str string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoIndexEventString is a free log retrieval operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterNoIndexEventString(opts *bind.FilterOpts) (*NetworkDebugContractNoIndexEventStringIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "NoIndexEventString")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractNoIndexEventStringIterator{contract: _NetworkDebugContract.contract, event: "NoIndexEventString", logs: logs, sub: sub}, nil
}

// WatchNoIndexEventString is a free log subscription operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchNoIndexEventString(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractNoIndexEventString) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "NoIndexEventString")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractNoIndexEventString)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEventString", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNoIndexEventString is a log parse operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseNoIndexEventString(log types.Log) (*NetworkDebugContractNoIndexEventString, error) {
	event := new(NetworkDebugContractNoIndexEventString)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEventString", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractNoIndexStructEventIterator is returned from FilterNoIndexStructEvent and is used to iterate over the raw logs and unpacked data for NoIndexStructEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexStructEventIterator struct {
	Event *NetworkDebugContractNoIndexStructEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractNoIndexStructEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractNoIndexStructEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractNoIndexStructEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractNoIndexStructEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractNoIndexStructEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractNoIndexStructEvent represents a NoIndexStructEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexStructEvent struct {
	A   NetworkDebugContractAccount
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoIndexStructEvent is a free log retrieval operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterNoIndexStructEvent(opts *bind.FilterOpts) (*NetworkDebugContractNoIndexStructEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "NoIndexStructEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractNoIndexStructEventIterator{contract: _NetworkDebugContract.contract, event: "NoIndexStructEvent", logs: logs, sub: sub}, nil
}

// WatchNoIndexStructEvent is a free log subscription operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchNoIndexStructEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractNoIndexStructEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "NoIndexStructEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractNoIndexStructEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexStructEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNoIndexStructEvent is a log parse operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseNoIndexStructEvent(log types.Log) (*NetworkDebugContractNoIndexStructEvent, error) {
	event := new(NetworkDebugContractNoIndexStructEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexStructEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractOneIndexEventIterator is returned from FilterOneIndexEvent and is used to iterate over the raw logs and unpacked data for OneIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractOneIndexEventIterator struct {
	Event *NetworkDebugContractOneIndexEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractOneIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractOneIndexEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractOneIndexEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractOneIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractOneIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractOneIndexEvent represents a OneIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractOneIndexEvent struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterOneIndexEvent is a free log retrieval operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterOneIndexEvent(opts *bind.FilterOpts, a []*big.Int) (*NetworkDebugContractOneIndexEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "OneIndexEvent", aRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractOneIndexEventIterator{contract: _NetworkDebugContract.contract, event: "OneIndexEvent", logs: logs, sub: sub}, nil
}

// WatchOneIndexEvent is a free log subscription operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchOneIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractOneIndexEvent, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "OneIndexEvent", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractOneIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "OneIndexEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOneIndexEvent is a log parse operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseOneIndexEvent(log types.Log) (*NetworkDebugContractOneIndexEvent, error) {
	event := new(NetworkDebugContractOneIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "OneIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractReceivedIterator is returned from FilterReceived and is used to iterate over the raw logs and unpacked data for Received events raised by the NetworkDebugContract contract.
type NetworkDebugContractReceivedIterator struct {
	Event *NetworkDebugContractReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractReceived represents a Received event raised by the NetworkDebugContract contract.
type NetworkDebugContractReceived struct {
	Caller  common.Address
	Amount  *big.Int
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReceived is a free log retrieval operation binding the contract event 0x59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832.
//
// Solidity: event Received(address caller, uint256 amount, string message)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterReceived(opts *bind.FilterOpts) (*NetworkDebugContractReceivedIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractReceivedIterator{contract: _NetworkDebugContract.contract, event: "Received", logs: logs, sub: sub}, nil
}

// WatchReceived is a free log subscription operation binding the contract event 0x59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832.
//
// Solidity: event Received(address caller, uint256 amount, string message)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchReceived(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractReceived) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractReceived)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "Received", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReceived is a log parse operation binding the contract event 0x59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832.
//
// Solidity: event Received(address caller, uint256 amount, string message)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseReceived(log types.Log) (*NetworkDebugContractReceived, error) {
	event := new(NetworkDebugContractReceived)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "Received", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator is returned from FilterThreeIndexAndOneNonIndexedEvent and is used to iterate over the raw logs and unpacked data for ThreeIndexAndOneNonIndexedEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator struct {
	Event *NetworkDebugContractThreeIndexAndOneNonIndexedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractThreeIndexAndOneNonIndexedEvent represents a ThreeIndexAndOneNonIndexedEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexAndOneNonIndexedEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	DataId    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterThreeIndexAndOneNonIndexedEvent is a free log retrieval operation binding the contract event 0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d.
//
// Solidity: event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterThreeIndexAndOneNonIndexedEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (*NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "ThreeIndexAndOneNonIndexedEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator{contract: _NetworkDebugContract.contract, event: "ThreeIndexAndOneNonIndexedEvent", logs: logs, sub: sub}, nil
}

// WatchThreeIndexAndOneNonIndexedEvent is a free log subscription operation binding the contract event 0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d.
//
// Solidity: event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchThreeIndexAndOneNonIndexedEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractThreeIndexAndOneNonIndexedEvent, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "ThreeIndexAndOneNonIndexedEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexAndOneNonIndexedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseThreeIndexAndOneNonIndexedEvent is a log parse operation binding the contract event 0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d.
//
// Solidity: event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseThreeIndexAndOneNonIndexedEvent(log types.Log) (*NetworkDebugContractThreeIndexAndOneNonIndexedEvent, error) {
	event := new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexAndOneNonIndexedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractThreeIndexEventIterator is returned from FilterThreeIndexEvent and is used to iterate over the raw logs and unpacked data for ThreeIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexEventIterator struct {
	Event *NetworkDebugContractThreeIndexEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractThreeIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractThreeIndexEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractThreeIndexEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractThreeIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractThreeIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractThreeIndexEvent represents a ThreeIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterThreeIndexEvent is a free log retrieval operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterThreeIndexEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (*NetworkDebugContractThreeIndexEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "ThreeIndexEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractThreeIndexEventIterator{contract: _NetworkDebugContract.contract, event: "ThreeIndexEvent", logs: logs, sub: sub}, nil
}

// WatchThreeIndexEvent is a free log subscription operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchThreeIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractThreeIndexEvent, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "ThreeIndexEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractThreeIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseThreeIndexEvent is a log parse operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseThreeIndexEvent(log types.Log) (*NetworkDebugContractThreeIndexEvent, error) {
	event := new(NetworkDebugContractThreeIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractTwoIndexEventIterator is returned from FilterTwoIndexEvent and is used to iterate over the raw logs and unpacked data for TwoIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractTwoIndexEventIterator struct {
	Event *NetworkDebugContractTwoIndexEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkDebugContractTwoIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractTwoIndexEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkDebugContractTwoIndexEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkDebugContractTwoIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractTwoIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractTwoIndexEvent represents a TwoIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractTwoIndexEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTwoIndexEvent is a free log retrieval operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterTwoIndexEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*NetworkDebugContractTwoIndexEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "TwoIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractTwoIndexEventIterator{contract: _NetworkDebugContract.contract, event: "TwoIndexEvent", logs: logs, sub: sub}, nil
}

// WatchTwoIndexEvent is a free log subscription operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchTwoIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractTwoIndexEvent, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "TwoIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractTwoIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "TwoIndexEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTwoIndexEvent is a log parse operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseTwoIndexEvent(log types.Log) (*NetworkDebugContractTwoIndexEvent, error) {
	event := new(NetworkDebugContractTwoIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "TwoIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
