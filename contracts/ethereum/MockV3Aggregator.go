// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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
)

// MockV3AggregatorMetaData contains all meta data concerning the MockV3Aggregator contract.
var MockV3AggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"int256\",\"name\":\"_initialAnswer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"updateAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startedAt\",\"type\":\"uint256\"}],\"name\":\"updateRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516104cf3803806104cf8339818101604052604081101561003357600080fd5b5080516020909101516000805460ff191660ff84161790556100548161005b565b50506100a2565b600181815542600281905560038054909201808355600090815260046020908152604080832095909555835482526005815284822083905592548152600690925291902055565b61041e806100b16000396000f3fe608060405234801561001057600080fd5b50600436106100a45760003560e01c8063313ce567146100a95780634aa2011f146100c757806350d25bcd1461010157806354fd4d501461011b578063668a0f02146101235780637284e4161461012b5780638205bf6a146101a85780639a6fc8f5146101b0578063a87a20ce1461021a578063b5ab58dc14610237578063b633620c14610254578063feaf968c14610271575b600080fd5b6100b1610279565b6040805160ff9092168252519081900360200190f35b6100ff600480360360808110156100dd57600080fd5b506001600160501b038135169060208101359060408101359060600135610282565b005b6101096102cc565b60408051918252519081900360200190f35b6101096102d2565b6101096102d7565b6101336102dd565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561016d578181015183820152602001610155565b50505050905090810190601f16801561019a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610109610314565b6101d6600480360360208110156101c657600080fd5b50356001600160501b031661031a565b60405180866001600160501b03168152602001858152602001848152602001838152602001826001600160501b031681526020019550505050505060405180910390f35b6100ff6004803603602081101561023057600080fd5b5035610350565b6101096004803603602081101561024d57600080fd5b5035610397565b6101096004803603602081101561026a57600080fd5b50356103a9565b6101d66103bb565b60005460ff1681565b6001600160501b0390931660038181556001849055600283905560009182526004602090815260408084209590955581548352600581528483209390935554815260069091522055565b60015481565b600081565b60035481565b60408051808201909152601f81527f76302e362f74657374732f4d6f636b563341676772656761746f722e736f6c00602082015290565b60025481565b6001600160501b038116600090815260046020908152604080832054600683528184205460059093529220549293919290918490565b600181815542600281905560038054909201808355600090815260046020908152604080832095909555835482526005815284822083905592548152600690925291902055565b60046020526000908152604090205481565b60056020526000908152604090205481565b6003546000818152600460209081526040808320546006835281842054600590935292205483909192939456fea26469706673582212202e67ecac433ebcfe150871ebdc8f25d791f60841ae48ca3364ee9e0db9fe070564736f6c63430007060033",
}

// MockV3AggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use MockV3AggregatorMetaData.ABI instead.
var MockV3AggregatorABI = MockV3AggregatorMetaData.ABI

// MockV3AggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockV3AggregatorMetaData.Bin instead.
var MockV3AggregatorBin = MockV3AggregatorMetaData.Bin

// DeployMockV3Aggregator deploys a new Ethereum contract, binding an instance of MockV3Aggregator to it.
func DeployMockV3Aggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _decimals uint8, _initialAnswer *big.Int) (common.Address, *types.Transaction, *MockV3Aggregator, error) {
	parsed, err := MockV3AggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockV3AggregatorBin), backend, _decimals, _initialAnswer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockV3Aggregator{MockV3AggregatorCaller: MockV3AggregatorCaller{contract: contract}, MockV3AggregatorTransactor: MockV3AggregatorTransactor{contract: contract}, MockV3AggregatorFilterer: MockV3AggregatorFilterer{contract: contract}}, nil
}

// MockV3Aggregator is an auto generated Go binding around an Ethereum contract.
type MockV3Aggregator struct {
	MockV3AggregatorCaller     // Read-only binding to the contract
	MockV3AggregatorTransactor // Write-only binding to the contract
	MockV3AggregatorFilterer   // Log filterer for contract events
}

// MockV3AggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockV3AggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockV3AggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockV3AggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockV3AggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockV3AggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockV3AggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockV3AggregatorSession struct {
	Contract     *MockV3Aggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockV3AggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockV3AggregatorCallerSession struct {
	Contract *MockV3AggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// MockV3AggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockV3AggregatorTransactorSession struct {
	Contract     *MockV3AggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// MockV3AggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockV3AggregatorRaw struct {
	Contract *MockV3Aggregator // Generic contract binding to access the raw methods on
}

// MockV3AggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockV3AggregatorCallerRaw struct {
	Contract *MockV3AggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// MockV3AggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockV3AggregatorTransactorRaw struct {
	Contract *MockV3AggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockV3Aggregator creates a new instance of MockV3Aggregator, bound to a specific deployed contract.
func NewMockV3Aggregator(address common.Address, backend bind.ContractBackend) (*MockV3Aggregator, error) {
	contract, err := bindMockV3Aggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockV3Aggregator{MockV3AggregatorCaller: MockV3AggregatorCaller{contract: contract}, MockV3AggregatorTransactor: MockV3AggregatorTransactor{contract: contract}, MockV3AggregatorFilterer: MockV3AggregatorFilterer{contract: contract}}, nil
}

// NewMockV3AggregatorCaller creates a new read-only instance of MockV3Aggregator, bound to a specific deployed contract.
func NewMockV3AggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockV3AggregatorCaller, error) {
	contract, err := bindMockV3Aggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorCaller{contract: contract}, nil
}

// NewMockV3AggregatorTransactor creates a new write-only instance of MockV3Aggregator, bound to a specific deployed contract.
func NewMockV3AggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockV3AggregatorTransactor, error) {
	contract, err := bindMockV3Aggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorTransactor{contract: contract}, nil
}

// NewMockV3AggregatorFilterer creates a new log filterer instance of MockV3Aggregator, bound to a specific deployed contract.
func NewMockV3AggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockV3AggregatorFilterer, error) {
	contract, err := bindMockV3Aggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorFilterer{contract: contract}, nil
}

// bindMockV3Aggregator binds a generic wrapper to an already deployed contract.
func bindMockV3Aggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockV3AggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockV3Aggregator *MockV3AggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockV3Aggregator.Contract.MockV3AggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockV3Aggregator *MockV3AggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.MockV3AggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockV3Aggregator *MockV3AggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.MockV3AggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockV3Aggregator *MockV3AggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockV3Aggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockV3Aggregator *MockV3AggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockV3Aggregator *MockV3AggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockV3Aggregator *MockV3AggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockV3Aggregator *MockV3AggregatorSession) Decimals() (uint8, error) {
	return _MockV3Aggregator.Contract.Decimals(&_MockV3Aggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) Decimals() (uint8, error) {
	return _MockV3Aggregator.Contract.Decimals(&_MockV3Aggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() pure returns(string)
func (_MockV3Aggregator *MockV3AggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() pure returns(string)
func (_MockV3Aggregator *MockV3AggregatorSession) Description() (string, error) {
	return _MockV3Aggregator.Contract.Description(&_MockV3Aggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() pure returns(string)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) Description() (string, error) {
	return _MockV3Aggregator.Contract.Description(&_MockV3Aggregator.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 ) view returns(int256)
func (_MockV3Aggregator *MockV3AggregatorCaller) GetAnswer(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "getAnswer", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 ) view returns(int256)
func (_MockV3Aggregator *MockV3AggregatorSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetAnswer(&_MockV3Aggregator.CallOpts, arg0)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 ) view returns(int256)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetAnswer(&_MockV3Aggregator.CallOpts, arg0)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3Aggregator *MockV3AggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "getRoundData", _roundId)

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
func (_MockV3Aggregator *MockV3AggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3Aggregator.Contract.GetRoundData(&_MockV3Aggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3Aggregator.Contract.GetRoundData(&_MockV3Aggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 ) view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCaller) GetTimestamp(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "getTimestamp", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 ) view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetTimestamp(&_MockV3Aggregator.CallOpts, arg0)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 ) view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _MockV3Aggregator.Contract.GetTimestamp(&_MockV3Aggregator.CallOpts, arg0)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockV3Aggregator *MockV3AggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockV3Aggregator *MockV3AggregatorSession) LatestAnswer() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestAnswer(&_MockV3Aggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestAnswer(&_MockV3Aggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorSession) LatestRound() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestRound(&_MockV3Aggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestRound(&_MockV3Aggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3Aggregator *MockV3AggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestRoundData")

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
func (_MockV3Aggregator *MockV3AggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3Aggregator.Contract.LatestRoundData(&_MockV3Aggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3Aggregator.Contract.LatestRoundData(&_MockV3Aggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestTimestamp(&_MockV3Aggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _MockV3Aggregator.Contract.LatestTimestamp(&_MockV3Aggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3Aggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorSession) Version() (*big.Int, error) {
	return _MockV3Aggregator.Contract.Version(&_MockV3Aggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockV3Aggregator *MockV3AggregatorCallerSession) Version() (*big.Int, error) {
	return _MockV3Aggregator.Contract.Version(&_MockV3Aggregator.CallOpts)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xa87a20ce.
//
// Solidity: function updateAnswer(int256 _answer) returns()
func (_MockV3Aggregator *MockV3AggregatorTransactor) UpdateAnswer(opts *bind.TransactOpts, _answer *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.contract.Transact(opts, "updateAnswer", _answer)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xa87a20ce.
//
// Solidity: function updateAnswer(int256 _answer) returns()
func (_MockV3Aggregator *MockV3AggregatorSession) UpdateAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateAnswer(&_MockV3Aggregator.TransactOpts, _answer)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xa87a20ce.
//
// Solidity: function updateAnswer(int256 _answer) returns()
func (_MockV3Aggregator *MockV3AggregatorTransactorSession) UpdateAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateAnswer(&_MockV3Aggregator.TransactOpts, _answer)
}

// UpdateRoundData is a paid mutator transaction binding the contract method 0x4aa2011f.
//
// Solidity: function updateRoundData(uint80 _roundId, int256 _answer, uint256 _timestamp, uint256 _startedAt) returns()
func (_MockV3Aggregator *MockV3AggregatorTransactor) UpdateRoundData(opts *bind.TransactOpts, _roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.contract.Transact(opts, "updateRoundData", _roundId, _answer, _timestamp, _startedAt)
}

// UpdateRoundData is a paid mutator transaction binding the contract method 0x4aa2011f.
//
// Solidity: function updateRoundData(uint80 _roundId, int256 _answer, uint256 _timestamp, uint256 _startedAt) returns()
func (_MockV3Aggregator *MockV3AggregatorSession) UpdateRoundData(_roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateRoundData(&_MockV3Aggregator.TransactOpts, _roundId, _answer, _timestamp, _startedAt)
}

// UpdateRoundData is a paid mutator transaction binding the contract method 0x4aa2011f.
//
// Solidity: function updateRoundData(uint80 _roundId, int256 _answer, uint256 _timestamp, uint256 _startedAt) returns()
func (_MockV3Aggregator *MockV3AggregatorTransactorSession) UpdateRoundData(_roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3Aggregator.Contract.UpdateRoundData(&_MockV3Aggregator.TransactOpts, _roundId, _answer, _timestamp, _startedAt)
}

// MockV3AggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the MockV3Aggregator contract.
type MockV3AggregatorAnswerUpdatedIterator struct {
	Event *MockV3AggregatorAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *MockV3AggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockV3AggregatorAnswerUpdated)
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
		it.Event = new(MockV3AggregatorAnswerUpdated)
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
func (it *MockV3AggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockV3AggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockV3AggregatorAnswerUpdated represents a AnswerUpdated event raised by the MockV3Aggregator contract.
type MockV3AggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_MockV3Aggregator *MockV3AggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*MockV3AggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorAnswerUpdatedIterator{contract: _MockV3Aggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_MockV3Aggregator *MockV3AggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockV3AggregatorAnswerUpdated)
				if err := _MockV3Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_MockV3Aggregator *MockV3AggregatorFilterer) ParseAnswerUpdated(log types.Log) (*MockV3AggregatorAnswerUpdated, error) {
	event := new(MockV3AggregatorAnswerUpdated)
	if err := _MockV3Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockV3AggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the MockV3Aggregator contract.
type MockV3AggregatorNewRoundIterator struct {
	Event *MockV3AggregatorNewRound // Event containing the contract specifics and raw log

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
func (it *MockV3AggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockV3AggregatorNewRound)
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
		it.Event = new(MockV3AggregatorNewRound)
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
func (it *MockV3AggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockV3AggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockV3AggregatorNewRound represents a NewRound event raised by the MockV3Aggregator contract.
type MockV3AggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_MockV3Aggregator *MockV3AggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*MockV3AggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorNewRoundIterator{contract: _MockV3Aggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_MockV3Aggregator *MockV3AggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _MockV3Aggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockV3AggregatorNewRound)
				if err := _MockV3Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_MockV3Aggregator *MockV3AggregatorFilterer) ParseNewRound(log types.Log) (*MockV3AggregatorNewRound, error) {
	event := new(MockV3AggregatorNewRound)
	if err := _MockV3Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
