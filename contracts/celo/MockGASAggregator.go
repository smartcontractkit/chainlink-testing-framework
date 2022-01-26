// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package celo

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = celo.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MockGASAggregatorMetaData contains all meta data concerning the MockGASAggregator contract.
var MockGASAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161029738038061029783398101604081905261002f91610037565b60005561004f565b600060208284031215610048578081fd5b5051919050565b6102398061005e6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063313ce5671461006757806354fd4d50146100855780637284e4161461009a57806385bb7d69146100af5780639a6fc8f5146100b7578063feaf968c146100db575b600080fd5b61006f6100e3565b60405161007c91906101f5565b60405180910390f35b61008d6100e8565b60405161007c9190610169565b6100a26100ed565b60405161007c9190610172565b61008d610118565b6100ca6100c536600461013b565b61011e565b60405161007c9594939291906101c5565b6100ca61012d565b601290565b600190565b60408051808201909152601181527026b7b1b5a3a0a9a0b3b3b932b3b0ba37b960791b602082015290565b60005481565b50600190600090429081908490565b600160004280839091929394565b60006020828403121561014c578081fd5b81356001600160501b0381168114610162578182fd5b9392505050565b90815260200190565b6000602080835283518082850152825b8181101561019e57858101830151858201604001528201610182565b818111156101af5783604083870101525b50601f01601f1916929092016040019392505050565b6001600160501b039586168152602081019490945260408401929092526060830152909116608082015260a00190565b60ff9190911681526020019056fea2646970667358221220ccc1e41a4368e05f078bb88943fcc2c5fcf58dc23a2fd036f26c13e140bae64964736f6c63430008000033",
}

// MockGASAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use MockGASAggregatorMetaData.ABI instead.
var MockGASAggregatorABI = MockGASAggregatorMetaData.ABI

// MockGASAggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockGASAggregatorMetaData.Bin instead.
var MockGASAggregatorBin = MockGASAggregatorMetaData.Bin

// DeployMockGASAggregator deploys a new Ethereum contract, binding an instance of MockGASAggregator to it.
func DeployMockGASAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *MockGASAggregator, error) {
	parsed, err := MockGASAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockGASAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockGASAggregator{MockGASAggregatorCaller: MockGASAggregatorCaller{contract: contract}, MockGASAggregatorTransactor: MockGASAggregatorTransactor{contract: contract}, MockGASAggregatorFilterer: MockGASAggregatorFilterer{contract: contract}}, nil
}

// MockGASAggregator is an auto generated Go binding around an Ethereum contract.
type MockGASAggregator struct {
	MockGASAggregatorCaller     // Read-only binding to the contract
	MockGASAggregatorTransactor // Write-only binding to the contract
	MockGASAggregatorFilterer   // Log filterer for contract events
}

// MockGASAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockGASAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockGASAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockGASAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockGASAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockGASAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockGASAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockGASAggregatorSession struct {
	Contract     *MockGASAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MockGASAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockGASAggregatorCallerSession struct {
	Contract *MockGASAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MockGASAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockGASAggregatorTransactorSession struct {
	Contract     *MockGASAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MockGASAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockGASAggregatorRaw struct {
	Contract *MockGASAggregator // Generic contract binding to access the raw methods on
}

// MockGASAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockGASAggregatorCallerRaw struct {
	Contract *MockGASAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// MockGASAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockGASAggregatorTransactorRaw struct {
	Contract *MockGASAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockGASAggregator creates a new instance of MockGASAggregator, bound to a specific deployed contract.
func NewMockGASAggregator(address common.Address, backend bind.ContractBackend) (*MockGASAggregator, error) {
	contract, err := bindMockGASAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregator{MockGASAggregatorCaller: MockGASAggregatorCaller{contract: contract}, MockGASAggregatorTransactor: MockGASAggregatorTransactor{contract: contract}, MockGASAggregatorFilterer: MockGASAggregatorFilterer{contract: contract}}, nil
}

// NewMockGASAggregatorCaller creates a new read-only instance of MockGASAggregator, bound to a specific deployed contract.
func NewMockGASAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockGASAggregatorCaller, error) {
	contract, err := bindMockGASAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregatorCaller{contract: contract}, nil
}

// NewMockGASAggregatorTransactor creates a new write-only instance of MockGASAggregator, bound to a specific deployed contract.
func NewMockGASAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockGASAggregatorTransactor, error) {
	contract, err := bindMockGASAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregatorTransactor{contract: contract}, nil
}

// NewMockGASAggregatorFilterer creates a new log filterer instance of MockGASAggregator, bound to a specific deployed contract.
func NewMockGASAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockGASAggregatorFilterer, error) {
	contract, err := bindMockGASAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockGASAggregatorFilterer{contract: contract}, nil
}

// bindMockGASAggregator binds a generic wrapper to an already deployed contract.
func bindMockGASAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockGASAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockGASAggregator *MockGASAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockGASAggregator.Contract.MockGASAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockGASAggregator *MockGASAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.MockGASAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockGASAggregator *MockGASAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.MockGASAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockGASAggregator *MockGASAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockGASAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockGASAggregator *MockGASAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockGASAggregator *MockGASAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockGASAggregator.Contract.contract.Transact(opts, method, params...)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockGASAggregator *MockGASAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockGASAggregator *MockGASAggregatorSession) Answer() (*big.Int, error) {
	return _MockGASAggregator.Contract.Answer(&_MockGASAggregator.CallOpts)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockGASAggregator *MockGASAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockGASAggregator.Contract.Answer(&_MockGASAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockGASAggregator *MockGASAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockGASAggregator *MockGASAggregatorSession) Decimals() (uint8, error) {
	return _MockGASAggregator.Contract.Decimals(&_MockGASAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockGASAggregator *MockGASAggregatorCallerSession) Decimals() (uint8, error) {
	return _MockGASAggregator.Contract.Decimals(&_MockGASAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_MockGASAggregator *MockGASAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_MockGASAggregator *MockGASAggregatorSession) Description() (string, error) {
	return _MockGASAggregator.Contract.Description(&_MockGASAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_MockGASAggregator *MockGASAggregatorCallerSession) Description() (string, error) {
	return _MockGASAggregator.Contract.Description(&_MockGASAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockGASAggregator *MockGASAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockGASAggregator *MockGASAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockGASAggregator.Contract.GetRoundData(&_MockGASAggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockGASAggregator *MockGASAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockGASAggregator.Contract.GetRoundData(&_MockGASAggregator.CallOpts, _roundId)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockGASAggregator *MockGASAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockGASAggregator *MockGASAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockGASAggregator.Contract.LatestRoundData(&_MockGASAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockGASAggregator *MockGASAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockGASAggregator.Contract.LatestRoundData(&_MockGASAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockGASAggregator *MockGASAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockGASAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockGASAggregator *MockGASAggregatorSession) Version() (*big.Int, error) {
	return _MockGASAggregator.Contract.Version(&_MockGASAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockGASAggregator *MockGASAggregatorCallerSession) Version() (*big.Int, error) {
	return _MockGASAggregator.Contract.Version(&_MockGASAggregator.CallOpts)
}
