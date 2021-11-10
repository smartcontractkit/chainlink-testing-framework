// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package celo

import (
	"math/big"
	"strings"

	celo "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = celo.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// KeeperConsumerABI is the input ABI used to generate the binding from.
const KeeperConsumerABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"updateInterval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimeStamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KeeperConsumerBin is the compiled bytecode used for deploying new contracts.
var KeeperConsumerBin = "0x60a060405234801561001057600080fd5b506040516102e93803806102e98339818101604052602081101561003357600080fd5b5051608052426001556000805560805161028f61005a60003980610260525061028f6000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806361bc221a1161005057806361bc221a146100f85780636e04ff0d14610100578063947a36fb146101f157610067565b80633f3b3b271461006c5780634585e33b14610086575b600080fd5b6100746101f9565b60408051918252519081900360200190f35b6100f66004803603602081101561009c57600080fd5b8101906020810181356401000000008111156100b757600080fd5b8201836020820111156100c957600080fd5b803590602001918460018302840111640100000000831117156100eb57600080fd5b5090925090506101ff565b005b61007461020c565b6101706004803603602081101561011657600080fd5b81019060208101813564010000000081111561013157600080fd5b82018360208201111561014357600080fd5b8035906020019184600183028401116401000000008311171561016557600080fd5b509092509050610212565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156101b557818101518382015260200161019d565b50505050905090810190601f1680156101e25780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b61007461025e565b60015481565b5050600080546001019055565b60005481565b600060606001848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b7f00000000000000000000000000000000000000000000000000000000000000008156fea164736f6c6343000706000a"

// DeployKeeperConsumer deploys a new Ethereum contract, binding an instance of KeeperConsumer to it.
func DeployKeeperConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, updateInterval *big.Int) (common.Address, *types.Transaction, *KeeperConsumer, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperConsumerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(KeeperConsumerBin), backend, updateInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperConsumer{KeeperConsumerCaller: KeeperConsumerCaller{contract: contract}, KeeperConsumerTransactor: KeeperConsumerTransactor{contract: contract}, KeeperConsumerFilterer: KeeperConsumerFilterer{contract: contract}}, nil
}

// KeeperConsumer is an auto generated Go binding around an Ethereum contract.
type KeeperConsumer struct {
	KeeperConsumerCaller     // Read-only binding to the contract
	KeeperConsumerTransactor // Write-only binding to the contract
	KeeperConsumerFilterer   // Log filterer for contract events
}

// KeeperConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperConsumerSession struct {
	Contract     *KeeperConsumer   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperConsumerCallerSession struct {
	Contract *KeeperConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// KeeperConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperConsumerTransactorSession struct {
	Contract     *KeeperConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// KeeperConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperConsumerRaw struct {
	Contract *KeeperConsumer // Generic contract binding to access the raw methods on
}

// KeeperConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperConsumerCallerRaw struct {
	Contract *KeeperConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperConsumerTransactorRaw struct {
	Contract *KeeperConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperConsumer creates a new instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumer(address common.Address, backend bind.ContractBackend) (*KeeperConsumer, error) {
	contract, err := bindKeeperConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumer{KeeperConsumerCaller: KeeperConsumerCaller{contract: contract}, KeeperConsumerTransactor: KeeperConsumerTransactor{contract: contract}, KeeperConsumerFilterer: KeeperConsumerFilterer{contract: contract}}, nil
}

// NewKeeperConsumerCaller creates a new read-only instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumerCaller(address common.Address, caller bind.ContractCaller) (*KeeperConsumerCaller, error) {
	contract, err := bindKeeperConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerCaller{contract: contract}, nil
}

// NewKeeperConsumerTransactor creates a new write-only instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperConsumerTransactor, error) {
	contract, err := bindKeeperConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerTransactor{contract: contract}, nil
}

// NewKeeperConsumerFilterer creates a new log filterer instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperConsumerFilterer, error) {
	contract, err := bindKeeperConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerFilterer{contract: contract}, nil
}

// bindKeeperConsumer binds a generic wrapper to an already deployed contract.
func bindKeeperConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumer *KeeperConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumer.Contract.KeeperConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumer *KeeperConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.KeeperConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumer *KeeperConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.KeeperConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumer *KeeperConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumer *KeeperConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumer *KeeperConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.contract.Transact(opts, method, params...)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerSession) Counter() (*big.Int, error) {
	return _KeeperConsumer.Contract.Counter(&_KeeperConsumer.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCallerSession) Counter() (*big.Int, error) {
	return _KeeperConsumer.Contract.Counter(&_KeeperConsumer.CallOpts)
}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerSession) Interval() (*big.Int, error) {
	return _KeeperConsumer.Contract.Interval(&_KeeperConsumer.CallOpts)
}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCallerSession) Interval() (*big.Int, error) {
	return _KeeperConsumer.Contract.Interval(&_KeeperConsumer.CallOpts)
}

// LastTimeStamp is a free data retrieval call binding the contract method 0x3f3b3b27.
//
// Solidity: function lastTimeStamp() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCaller) LastTimeStamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "lastTimeStamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastTimeStamp is a free data retrieval call binding the contract method 0x3f3b3b27.
//
// Solidity: function lastTimeStamp() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerSession) LastTimeStamp() (*big.Int, error) {
	return _KeeperConsumer.Contract.LastTimeStamp(&_KeeperConsumer.CallOpts)
}

// LastTimeStamp is a free data retrieval call binding the contract method 0x3f3b3b27.
//
// Solidity: function lastTimeStamp() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCallerSession) LastTimeStamp() (*big.Int, error) {
	return _KeeperConsumer.Contract.LastTimeStamp(&_KeeperConsumer.CallOpts)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes checkData) returns(bool upkeepNeeded, bytes performData)
func (_KeeperConsumer *KeeperConsumerTransactor) CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.contract.Transact(opts, "checkUpkeep", checkData)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes checkData) returns(bool upkeepNeeded, bytes performData)
func (_KeeperConsumer *KeeperConsumerSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.CheckUpkeep(&_KeeperConsumer.TransactOpts, checkData)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes checkData) returns(bool upkeepNeeded, bytes performData)
func (_KeeperConsumer *KeeperConsumerTransactorSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.CheckUpkeep(&_KeeperConsumer.TransactOpts, checkData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_KeeperConsumer *KeeperConsumerTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.contract.Transact(opts, "performUpkeep", performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_KeeperConsumer *KeeperConsumerSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.PerformUpkeep(&_KeeperConsumer.TransactOpts, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_KeeperConsumer *KeeperConsumerTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.PerformUpkeep(&_KeeperConsumer.TransactOpts, performData)
}
