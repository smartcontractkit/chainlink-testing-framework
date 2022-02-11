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

// KeeperConsumerPerformanceMetaData contains all meta data concerning the KeeperConsumerPerformance contract.
var KeeperConsumerPerformanceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_averageEligibilityCadence\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"eligible\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"averageEligibilityCadence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newTestRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newAverageEligibilityCadence\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155600060045534801561001e57600080fd5b5060405161042d38038061042d8339818101604052604081101561004157600080fd5b5080516020909101516002919091556003556103cb806100626000396000f3fe608060405234801561001057600080fd5b506004361061008e5760003560e01c80634585e33b14610093578063523d9b8a146101035780636250a13a1461011d5780636e04ff0d146101255780637f407edf14610214578063926f086e14610237578063a9a4c57c1461023f578063c228a98e14610247578063d826f88f14610263578063e303666f1461026b575b600080fd5b610101600480360360208110156100a957600080fd5b810190602081018135600160201b8111156100c357600080fd5b8201836020820111156100d557600080fd5b803590602001918460018302840111600160201b831117156100f657600080fd5b509092509050610273565b005b61010b610307565b60408051918252519081900360200190f35b61010b61030d565b6101936004803603602081101561013b57600080fd5b810190602081018135600160201b81111561015557600080fd5b82018360208201111561016757600080fd5b803590602001918460018302840111600160201b8311171561018857600080fd5b509092509050610313565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156101d85781810151838201526020016101c0565b50505050905090810190601f1680156102055780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101016004803603604081101561022a57600080fd5b508035906020013561033a565b61010b610345565b61010b61034b565b61024f610351565b604080519115158252519081900360200190f35b610101610360565b61010b61036a565b600061027d610370565b60005460015460408051841515815232602082015280820193909352606083019190915243608083018190529051929350917fbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc09181900360a00190a1816102e357600080fd5b6000546102f05760008190555b600354016001908155600480549091019055505050565b60015481565b60025481565b6000606061031f610370565b60405180602001604052806000815250915091509250929050565b600291909155600355565b60005481565b60035481565b600061035b610370565b905090565b6000808055600455565b60045490565b60008054158061035b5750600254600054430310801561035b5750506001544310159056fea2646970667358221220ea467b763c458ff66dc749aed84e7e57b50539bb5ff0c6b6cb57c2d353915fcc64736f6c63430007060033",
}

// KeeperConsumerPerformanceABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperConsumerPerformanceMetaData.ABI instead.
var KeeperConsumerPerformanceABI = KeeperConsumerPerformanceMetaData.ABI

// KeeperConsumerPerformanceBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperConsumerPerformanceMetaData.Bin instead.
var KeeperConsumerPerformanceBin = KeeperConsumerPerformanceMetaData.Bin

// DeployKeeperConsumerPerformance deploys a new Ethereum contract, binding an instance of KeeperConsumerPerformance to it.
func DeployKeeperConsumerPerformance(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _averageEligibilityCadence *big.Int) (common.Address, *types.Transaction, *KeeperConsumerPerformance, error) {
	parsed, err := KeeperConsumerPerformanceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperConsumerPerformanceBin), backend, _testRange, _averageEligibilityCadence)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperConsumerPerformance{KeeperConsumerPerformanceCaller: KeeperConsumerPerformanceCaller{contract: contract}, KeeperConsumerPerformanceTransactor: KeeperConsumerPerformanceTransactor{contract: contract}, KeeperConsumerPerformanceFilterer: KeeperConsumerPerformanceFilterer{contract: contract}}, nil
}

// KeeperConsumerPerformance is an auto generated Go binding around an Ethereum contract.
type KeeperConsumerPerformance struct {
	KeeperConsumerPerformanceCaller     // Read-only binding to the contract
	KeeperConsumerPerformanceTransactor // Write-only binding to the contract
	KeeperConsumerPerformanceFilterer   // Log filterer for contract events
}

// KeeperConsumerPerformanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperConsumerPerformanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerPerformanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperConsumerPerformanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerPerformanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperConsumerPerformanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerPerformanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperConsumerPerformanceSession struct {
	Contract     *KeeperConsumerPerformance // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// KeeperConsumerPerformanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperConsumerPerformanceCallerSession struct {
	Contract *KeeperConsumerPerformanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// KeeperConsumerPerformanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperConsumerPerformanceTransactorSession struct {
	Contract     *KeeperConsumerPerformanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// KeeperConsumerPerformanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperConsumerPerformanceRaw struct {
	Contract *KeeperConsumerPerformance // Generic contract binding to access the raw methods on
}

// KeeperConsumerPerformanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperConsumerPerformanceCallerRaw struct {
	Contract *KeeperConsumerPerformanceCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperConsumerPerformanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperConsumerPerformanceTransactorRaw struct {
	Contract *KeeperConsumerPerformanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperConsumerPerformance creates a new instance of KeeperConsumerPerformance, bound to a specific deployed contract.
func NewKeeperConsumerPerformance(address common.Address, backend bind.ContractBackend) (*KeeperConsumerPerformance, error) {
	contract, err := bindKeeperConsumerPerformance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformance{KeeperConsumerPerformanceCaller: KeeperConsumerPerformanceCaller{contract: contract}, KeeperConsumerPerformanceTransactor: KeeperConsumerPerformanceTransactor{contract: contract}, KeeperConsumerPerformanceFilterer: KeeperConsumerPerformanceFilterer{contract: contract}}, nil
}

// NewKeeperConsumerPerformanceCaller creates a new read-only instance of KeeperConsumerPerformance, bound to a specific deployed contract.
func NewKeeperConsumerPerformanceCaller(address common.Address, caller bind.ContractCaller) (*KeeperConsumerPerformanceCaller, error) {
	contract, err := bindKeeperConsumerPerformance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformanceCaller{contract: contract}, nil
}

// NewKeeperConsumerPerformanceTransactor creates a new write-only instance of KeeperConsumerPerformance, bound to a specific deployed contract.
func NewKeeperConsumerPerformanceTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperConsumerPerformanceTransactor, error) {
	contract, err := bindKeeperConsumerPerformance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformanceTransactor{contract: contract}, nil
}

// NewKeeperConsumerPerformanceFilterer creates a new log filterer instance of KeeperConsumerPerformance, bound to a specific deployed contract.
func NewKeeperConsumerPerformanceFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperConsumerPerformanceFilterer, error) {
	contract, err := bindKeeperConsumerPerformance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformanceFilterer{contract: contract}, nil
}

// bindKeeperConsumerPerformance binds a generic wrapper to an already deployed contract.
func bindKeeperConsumerPerformance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperConsumerPerformanceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumerPerformance.Contract.KeeperConsumerPerformanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.KeeperConsumerPerformanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.KeeperConsumerPerformanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumerPerformance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.contract.Transact(opts, method, params...)
}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "averageEligibilityCadence")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) AverageEligibilityCadence() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.AverageEligibilityCadence(&_KeeperConsumerPerformance.CallOpts)
}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) AverageEligibilityCadence() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.AverageEligibilityCadence(&_KeeperConsumerPerformance.CallOpts)
}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) CheckEligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "checkEligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) CheckEligible() (bool, error) {
	return _KeeperConsumerPerformance.Contract.CheckEligible(&_KeeperConsumerPerformance.CallOpts)
}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) CheckEligible() (bool, error) {
	return _KeeperConsumerPerformance.Contract.CheckEligible(&_KeeperConsumerPerformance.CallOpts)
}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) GetCountPerforms(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "getCountPerforms")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) GetCountPerforms() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.GetCountPerforms(&_KeeperConsumerPerformance.CallOpts)
}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) GetCountPerforms() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.GetCountPerforms(&_KeeperConsumerPerformance.CallOpts)
}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) InitialCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "initialCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) InitialCall() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.InitialCall(&_KeeperConsumerPerformance.CallOpts)
}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) InitialCall() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.InitialCall(&_KeeperConsumerPerformance.CallOpts)
}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) NextEligible(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "nextEligible")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) NextEligible() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.NextEligible(&_KeeperConsumerPerformance.CallOpts)
}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) NextEligible() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.NextEligible(&_KeeperConsumerPerformance.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) TestRange() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.TestRange(&_KeeperConsumerPerformance.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) TestRange() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.TestRange(&_KeeperConsumerPerformance.CallOpts)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) returns(bool, bytes)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) CheckUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "checkUpkeep", data)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) returns(bool, bytes)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) CheckUpkeep(data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.CheckUpkeep(&_KeeperConsumerPerformance.TransactOpts, data)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) returns(bool, bytes)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) CheckUpkeep(data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.CheckUpkeep(&_KeeperConsumerPerformance.TransactOpts, data)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes data) returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) PerformUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "performUpkeep", data)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes data) returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.PerformUpkeep(&_KeeperConsumerPerformance.TransactOpts, data)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes data) returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.PerformUpkeep(&_KeeperConsumerPerformance.TransactOpts, data)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) Reset() (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.Reset(&_KeeperConsumerPerformance.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) Reset() (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.Reset(&_KeeperConsumerPerformance.TransactOpts)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "setSpread", _newTestRange, _newAverageEligibilityCadence)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetSpread(&_KeeperConsumerPerformance.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetSpread(&_KeeperConsumerPerformance.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

// KeeperConsumerPerformancePerformingUpkeepIterator is returned from FilterPerformingUpkeep and is used to iterate over the raw logs and unpacked data for PerformingUpkeep events raised by the KeeperConsumerPerformance contract.
type KeeperConsumerPerformancePerformingUpkeepIterator struct {
	Event *KeeperConsumerPerformancePerformingUpkeep // Event containing the contract specifics and raw log

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
func (it *KeeperConsumerPerformancePerformingUpkeepIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperConsumerPerformancePerformingUpkeep)
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
		it.Event = new(KeeperConsumerPerformancePerformingUpkeep)
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
func (it *KeeperConsumerPerformancePerformingUpkeepIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperConsumerPerformancePerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperConsumerPerformancePerformingUpkeep represents a PerformingUpkeep event raised by the KeeperConsumerPerformance contract.
type KeeperConsumerPerformancePerformingUpkeep struct {
	Eligible     bool
	From         common.Address
	InitialCall  *big.Int
	NextEligible *big.Int
	BlockNumber  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPerformingUpkeep is a free log retrieval operation binding the contract event 0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0.
//
// Solidity: event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*KeeperConsumerPerformancePerformingUpkeepIterator, error) {

	logs, sub, err := _KeeperConsumerPerformance.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformancePerformingUpkeepIterator{contract: _KeeperConsumerPerformance.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

// WatchPerformingUpkeep is a free log subscription operation binding the contract event 0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0.
//
// Solidity: event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *KeeperConsumerPerformancePerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _KeeperConsumerPerformance.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperConsumerPerformancePerformingUpkeep)
				if err := _KeeperConsumerPerformance.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

// ParsePerformingUpkeep is a log parse operation binding the contract event 0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0.
//
// Solidity: event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_KeeperConsumerPerformance *KeeperConsumerPerformanceFilterer) ParsePerformingUpkeep(log types.Log) (*KeeperConsumerPerformancePerformingUpkeep, error) {
	event := new(KeeperConsumerPerformancePerformingUpkeep)
	if err := _KeeperConsumerPerformance.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
