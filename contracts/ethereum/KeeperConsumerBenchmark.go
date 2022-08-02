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

// KeeperConsumerBenchmarkMetaData contains all meta data concerning the KeeperConsumerBenchmark contract.
var KeeperConsumerBenchmarkMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_averageEligibilityCadence\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_performGasToBurn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"averageEligibilityCadence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"firstEligibleBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newTestRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newAverageEligibilityCadence\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155600060085534801561001e57600080fd5b506040516106d43803806106d48339818101604052608081101561004157600080fd5b5080516020820151604083015160609093015160028390556003829055600484905560058190559192909182610075610091565b8161007c57fe5b064301600101600681905550505050506100c2565b6040805160001943014060208083019190915230828401528251808303840181526060909201909252805191012090565b610603806100d16000396000f3fe608060405234801561001057600080fd5b50600436106100eb5760003560e01c80637145f11b116100925780637145f11b146102bb5780637f407edf146102ec578063926f086e1461030f578063a9a4c57c14610317578063b30566b41461031f578063c228a98e14610327578063c4da244d1461032f578063d826f88f14610337578063e303666f1461033f576100eb565b806306661abd146100f057806313bda75b1461010a5780632555d2cf146101295780632ff3617d146101465780634585e33b1461014e578063523d9b8a146101bc5780636250a13a146101c45780636e04ff0d146101cc575b600080fd5b6100f8610347565b60408051918252519081900360200190f35b6101276004803603602081101561012057600080fd5b503561034d565b005b6101276004803603602081101561013f57600080fd5b5035610352565b6100f8610357565b6101276004803603602081101561016457600080fd5b810190602081018135600160201b81111561017e57600080fd5b82018360208201111561019057600080fd5b803590602001918460018302840111600160201b831117156101b157600080fd5b50909250905061035d565b6100f861044e565b6100f8610454565b61023a600480360360208110156101e257600080fd5b810190602081018135600160201b8111156101fc57600080fd5b82018360208201111561020e57600080fd5b803590602001918460018302840111600160201b8311171561022f57600080fd5b50909250905061045a565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561027f578181015183820152602001610267565b50505050905090810190601f1680156102ac5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6102d8600480360360208110156102d157600080fd5b50356104f9565b604080519115158252519081900360200190f35b6101276004803603604081101561030257600080fd5b508035906020013561050e565b6100f8610519565b6100f861051f565b6100f8610525565b6102d861052b565b6100f861053a565b610127610540565b6100f8610568565b60085481565b600455565b600555565b60045481565b61036561056e565b61036e57600080fd5b60005a90506000546000141561038357436000555b600354439081016001818155600880549091019055600054604080513281526020810192909252818101929092526060810192909252517f1313be6f6d6263f115d3e986c9622f868fcda43c8b8e7ef193e7a53d75a4d27c9181900360800190a160001943014060005b6005545a8403101561044757808015610414575060008281526007602052604090205460ff165b604080516020808201959095523081830152815180820383018152606090910190915280519301929092209190506103ed565b5050505050565b60015481565b60025481565b6000606060005a905060001943014060005b6004545a840310156104c657808015610493575060008281526007602052604090205460ff165b6040805160208082019590955230818301528151808203830181526060909101909152805193019290922091905061046c565b6104ce61056e565b6040805192151560208085019190915281518085039091018152928101905297909650945050505050565b60076020526000908152604090205460ff1681565b600291909155600355565b60005481565b60035481565b60055481565b600061053561056e565b905090565b60065481565b600080805560085560035461055361059c565b8161055a57fe5b064301600101600681905550565b60085490565b600080541561059257600254600054430310801561058d575060015443115b610535565b5060065443101590565b604080516000194301406020808301919091523082840152825180830384018152606090920190925280519101209056fea26469706673582212201e010f07b1e726f98e38b548f87af748753ce70c37b5d461603fc63ee4cd837d64736f6c63430007060033",
}

// KeeperConsumerBenchmarkABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperConsumerBenchmarkMetaData.ABI instead.
var KeeperConsumerBenchmarkABI = KeeperConsumerBenchmarkMetaData.ABI

// KeeperConsumerBenchmarkBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperConsumerBenchmarkMetaData.Bin instead.
var KeeperConsumerBenchmarkBin = KeeperConsumerBenchmarkMetaData.Bin

// DeployKeeperConsumerBenchmark deploys a new Ethereum contract, binding an instance of KeeperConsumerBenchmark to it.
func DeployKeeperConsumerBenchmark(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _averageEligibilityCadence *big.Int, _checkGasToBurn *big.Int, _performGasToBurn *big.Int) (common.Address, *types.Transaction, *KeeperConsumerBenchmark, error) {
	parsed, err := KeeperConsumerBenchmarkMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperConsumerBenchmarkBin), backend, _testRange, _averageEligibilityCadence, _checkGasToBurn, _performGasToBurn)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperConsumerBenchmark{KeeperConsumerBenchmarkCaller: KeeperConsumerBenchmarkCaller{contract: contract}, KeeperConsumerBenchmarkTransactor: KeeperConsumerBenchmarkTransactor{contract: contract}, KeeperConsumerBenchmarkFilterer: KeeperConsumerBenchmarkFilterer{contract: contract}}, nil
}

// KeeperConsumerBenchmark is an auto generated Go binding around an Ethereum contract.
type KeeperConsumerBenchmark struct {
	KeeperConsumerBenchmarkCaller     // Read-only binding to the contract
	KeeperConsumerBenchmarkTransactor // Write-only binding to the contract
	KeeperConsumerBenchmarkFilterer   // Log filterer for contract events
}

// KeeperConsumerBenchmarkCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperConsumerBenchmarkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerBenchmarkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperConsumerBenchmarkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerBenchmarkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperConsumerBenchmarkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerBenchmarkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperConsumerBenchmarkSession struct {
	Contract     *KeeperConsumerBenchmark // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// KeeperConsumerBenchmarkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperConsumerBenchmarkCallerSession struct {
	Contract *KeeperConsumerBenchmarkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// KeeperConsumerBenchmarkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperConsumerBenchmarkTransactorSession struct {
	Contract     *KeeperConsumerBenchmarkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// KeeperConsumerBenchmarkRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperConsumerBenchmarkRaw struct {
	Contract *KeeperConsumerBenchmark // Generic contract binding to access the raw methods on
}

// KeeperConsumerBenchmarkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperConsumerBenchmarkCallerRaw struct {
	Contract *KeeperConsumerBenchmarkCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperConsumerBenchmarkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperConsumerBenchmarkTransactorRaw struct {
	Contract *KeeperConsumerBenchmarkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperConsumerBenchmark creates a new instance of KeeperConsumerBenchmark, bound to a specific deployed contract.
func NewKeeperConsumerBenchmark(address common.Address, backend bind.ContractBackend) (*KeeperConsumerBenchmark, error) {
	contract, err := bindKeeperConsumerBenchmark(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerBenchmark{KeeperConsumerBenchmarkCaller: KeeperConsumerBenchmarkCaller{contract: contract}, KeeperConsumerBenchmarkTransactor: KeeperConsumerBenchmarkTransactor{contract: contract}, KeeperConsumerBenchmarkFilterer: KeeperConsumerBenchmarkFilterer{contract: contract}}, nil
}

// NewKeeperConsumerBenchmarkCaller creates a new read-only instance of KeeperConsumerBenchmark, bound to a specific deployed contract.
func NewKeeperConsumerBenchmarkCaller(address common.Address, caller bind.ContractCaller) (*KeeperConsumerBenchmarkCaller, error) {
	contract, err := bindKeeperConsumerBenchmark(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerBenchmarkCaller{contract: contract}, nil
}

// NewKeeperConsumerBenchmarkTransactor creates a new write-only instance of KeeperConsumerBenchmark, bound to a specific deployed contract.
func NewKeeperConsumerBenchmarkTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperConsumerBenchmarkTransactor, error) {
	contract, err := bindKeeperConsumerBenchmark(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerBenchmarkTransactor{contract: contract}, nil
}

// NewKeeperConsumerBenchmarkFilterer creates a new log filterer instance of KeeperConsumerBenchmark, bound to a specific deployed contract.
func NewKeeperConsumerBenchmarkFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperConsumerBenchmarkFilterer, error) {
	contract, err := bindKeeperConsumerBenchmark(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerBenchmarkFilterer{contract: contract}, nil
}

// bindKeeperConsumerBenchmark binds a generic wrapper to an already deployed contract.
func bindKeeperConsumerBenchmark(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperConsumerBenchmarkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumerBenchmark.Contract.KeeperConsumerBenchmarkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.KeeperConsumerBenchmarkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.KeeperConsumerBenchmarkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumerBenchmark.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.contract.Transact(opts, method, params...)
}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "averageEligibilityCadence")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) AverageEligibilityCadence() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.AverageEligibilityCadence(&_KeeperConsumerBenchmark.CallOpts)
}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) AverageEligibilityCadence() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.AverageEligibilityCadence(&_KeeperConsumerBenchmark.CallOpts)
}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) CheckEligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "checkEligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) CheckEligible() (bool, error) {
	return _KeeperConsumerBenchmark.Contract.CheckEligible(&_KeeperConsumerBenchmark.CallOpts)
}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) CheckEligible() (bool, error) {
	return _KeeperConsumerBenchmark.Contract.CheckEligible(&_KeeperConsumerBenchmark.CallOpts)
}

// CheckGasToBurn is a free data retrieval call binding the contract method 0x2ff3617d.
//
// Solidity: function checkGasToBurn() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) CheckGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "checkGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CheckGasToBurn is a free data retrieval call binding the contract method 0x2ff3617d.
//
// Solidity: function checkGasToBurn() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) CheckGasToBurn() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.CheckGasToBurn(&_KeeperConsumerBenchmark.CallOpts)
}

// CheckGasToBurn is a free data retrieval call binding the contract method 0x2ff3617d.
//
// Solidity: function checkGasToBurn() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) CheckGasToBurn() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.CheckGasToBurn(&_KeeperConsumerBenchmark.CallOpts)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _KeeperConsumerBenchmark.Contract.CheckUpkeep(&_KeeperConsumerBenchmark.CallOpts, data)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _KeeperConsumerBenchmark.Contract.CheckUpkeep(&_KeeperConsumerBenchmark.CallOpts, data)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) Count() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.Count(&_KeeperConsumerBenchmark.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) Count() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.Count(&_KeeperConsumerBenchmark.CallOpts)
}

// DummyMap is a free data retrieval call binding the contract method 0x7145f11b.
//
// Solidity: function dummyMap(bytes32 ) view returns(bool)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DummyMap is a free data retrieval call binding the contract method 0x7145f11b.
//
// Solidity: function dummyMap(bytes32 ) view returns(bool)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _KeeperConsumerBenchmark.Contract.DummyMap(&_KeeperConsumerBenchmark.CallOpts, arg0)
}

// DummyMap is a free data retrieval call binding the contract method 0x7145f11b.
//
// Solidity: function dummyMap(bytes32 ) view returns(bool)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _KeeperConsumerBenchmark.Contract.DummyMap(&_KeeperConsumerBenchmark.CallOpts, arg0)
}

// FirstEligibleBlock is a free data retrieval call binding the contract method 0xc4da244d.
//
// Solidity: function firstEligibleBlock() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) FirstEligibleBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "firstEligibleBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FirstEligibleBlock is a free data retrieval call binding the contract method 0xc4da244d.
//
// Solidity: function firstEligibleBlock() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) FirstEligibleBlock() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.FirstEligibleBlock(&_KeeperConsumerBenchmark.CallOpts)
}

// FirstEligibleBlock is a free data retrieval call binding the contract method 0xc4da244d.
//
// Solidity: function firstEligibleBlock() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) FirstEligibleBlock() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.FirstEligibleBlock(&_KeeperConsumerBenchmark.CallOpts)
}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) GetCountPerforms(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "getCountPerforms")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) GetCountPerforms() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.GetCountPerforms(&_KeeperConsumerBenchmark.CallOpts)
}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) GetCountPerforms() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.GetCountPerforms(&_KeeperConsumerBenchmark.CallOpts)
}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) InitialCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "initialCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) InitialCall() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.InitialCall(&_KeeperConsumerBenchmark.CallOpts)
}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) InitialCall() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.InitialCall(&_KeeperConsumerBenchmark.CallOpts)
}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) NextEligible(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "nextEligible")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) NextEligible() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.NextEligible(&_KeeperConsumerBenchmark.CallOpts)
}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) NextEligible() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.NextEligible(&_KeeperConsumerBenchmark.CallOpts)
}

// PerformGasToBurn is a free data retrieval call binding the contract method 0xb30566b4.
//
// Solidity: function performGasToBurn() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) PerformGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "performGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerformGasToBurn is a free data retrieval call binding the contract method 0xb30566b4.
//
// Solidity: function performGasToBurn() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) PerformGasToBurn() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.PerformGasToBurn(&_KeeperConsumerBenchmark.CallOpts)
}

// PerformGasToBurn is a free data retrieval call binding the contract method 0xb30566b4.
//
// Solidity: function performGasToBurn() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) PerformGasToBurn() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.PerformGasToBurn(&_KeeperConsumerBenchmark.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerBenchmark.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) TestRange() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.TestRange(&_KeeperConsumerBenchmark.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkCallerSession) TestRange() (*big.Int, error) {
	return _KeeperConsumerBenchmark.Contract.TestRange(&_KeeperConsumerBenchmark.CallOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes ) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactor) PerformUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.contract.Transact(opts, "performUpkeep", arg0)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes ) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.PerformUpkeep(&_KeeperConsumerBenchmark.TransactOpts, arg0)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes ) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.PerformUpkeep(&_KeeperConsumerBenchmark.TransactOpts, arg0)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) Reset() (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.Reset(&_KeeperConsumerBenchmark.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorSession) Reset() (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.Reset(&_KeeperConsumerBenchmark.TransactOpts)
}

// SetCheckGasToBurn is a paid mutator transaction binding the contract method 0x13bda75b.
//
// Solidity: function setCheckGasToBurn(uint256 value) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.contract.Transact(opts, "setCheckGasToBurn", value)
}

// SetCheckGasToBurn is a paid mutator transaction binding the contract method 0x13bda75b.
//
// Solidity: function setCheckGasToBurn(uint256 value) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.SetCheckGasToBurn(&_KeeperConsumerBenchmark.TransactOpts, value)
}

// SetCheckGasToBurn is a paid mutator transaction binding the contract method 0x13bda75b.
//
// Solidity: function setCheckGasToBurn(uint256 value) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.SetCheckGasToBurn(&_KeeperConsumerBenchmark.TransactOpts, value)
}

// SetPerformGasToBurn is a paid mutator transaction binding the contract method 0x2555d2cf.
//
// Solidity: function setPerformGasToBurn(uint256 value) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.contract.Transact(opts, "setPerformGasToBurn", value)
}

// SetPerformGasToBurn is a paid mutator transaction binding the contract method 0x2555d2cf.
//
// Solidity: function setPerformGasToBurn(uint256 value) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.SetPerformGasToBurn(&_KeeperConsumerBenchmark.TransactOpts, value)
}

// SetPerformGasToBurn is a paid mutator transaction binding the contract method 0x2555d2cf.
//
// Solidity: function setPerformGasToBurn(uint256 value) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.SetPerformGasToBurn(&_KeeperConsumerBenchmark.TransactOpts, value)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactor) SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.contract.Transact(opts, "setSpread", _newTestRange, _newAverageEligibilityCadence)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.SetSpread(&_KeeperConsumerBenchmark.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkTransactorSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerBenchmark.Contract.SetSpread(&_KeeperConsumerBenchmark.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

// KeeperConsumerBenchmarkPerformingUpkeepIterator is returned from FilterPerformingUpkeep and is used to iterate over the raw logs and unpacked data for PerformingUpkeep events raised by the KeeperConsumerBenchmark contract.
type KeeperConsumerBenchmarkPerformingUpkeepIterator struct {
	Event *KeeperConsumerBenchmarkPerformingUpkeep // Event containing the contract specifics and raw log

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
func (it *KeeperConsumerBenchmarkPerformingUpkeepIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperConsumerBenchmarkPerformingUpkeep)
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
		it.Event = new(KeeperConsumerBenchmarkPerformingUpkeep)
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
func (it *KeeperConsumerBenchmarkPerformingUpkeepIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperConsumerBenchmarkPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperConsumerBenchmarkPerformingUpkeep represents a PerformingUpkeep event raised by the KeeperConsumerBenchmark contract.
type KeeperConsumerBenchmarkPerformingUpkeep struct {
	From         common.Address
	InitialCall  *big.Int
	NextEligible *big.Int
	BlockNumber  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPerformingUpkeep is a free log retrieval operation binding the contract event 0x1313be6f6d6263f115d3e986c9622f868fcda43c8b8e7ef193e7a53d75a4d27c.
//
// Solidity: event PerformingUpkeep(address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*KeeperConsumerBenchmarkPerformingUpkeepIterator, error) {

	logs, sub, err := _KeeperConsumerBenchmark.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerBenchmarkPerformingUpkeepIterator{contract: _KeeperConsumerBenchmark.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

// WatchPerformingUpkeep is a free log subscription operation binding the contract event 0x1313be6f6d6263f115d3e986c9622f868fcda43c8b8e7ef193e7a53d75a4d27c.
//
// Solidity: event PerformingUpkeep(address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *KeeperConsumerBenchmarkPerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _KeeperConsumerBenchmark.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperConsumerBenchmarkPerformingUpkeep)
				if err := _KeeperConsumerBenchmark.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

// ParsePerformingUpkeep is a log parse operation binding the contract event 0x1313be6f6d6263f115d3e986c9622f868fcda43c8b8e7ef193e7a53d75a4d27c.
//
// Solidity: event PerformingUpkeep(address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_KeeperConsumerBenchmark *KeeperConsumerBenchmarkFilterer) ParsePerformingUpkeep(log types.Log) (*KeeperConsumerBenchmarkPerformingUpkeep, error) {
	event := new(KeeperConsumerBenchmarkPerformingUpkeep)
	if err := _KeeperConsumerBenchmark.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
