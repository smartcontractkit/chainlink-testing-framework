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

// ChainlinkABI is the input ABI used to generate the binding from.
const ChainlinkABI = "[]"

// ChainlinkBin is the compiled bytecode used for deploying new contracts.
var ChainlinkBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220d35e03cf1fa89a4c6150c4c3262fc1c75dd6c0da0b7cee227011b9128fd96dda64736f6c63430006060033"

// DeployChainlink deploys a new Ethereum contract, binding an instance of Chainlink to it.
func DeployChainlink(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Chainlink, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainlinkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChainlinkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Chainlink{ChainlinkCaller: ChainlinkCaller{contract: contract}, ChainlinkTransactor: ChainlinkTransactor{contract: contract}, ChainlinkFilterer: ChainlinkFilterer{contract: contract}}, nil
}

// Chainlink is an auto generated Go binding around an Ethereum contract.
type Chainlink struct {
	ChainlinkCaller     // Read-only binding to the contract
	ChainlinkTransactor // Write-only binding to the contract
	ChainlinkFilterer   // Log filterer for contract events
}

// ChainlinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainlinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainlinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainlinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainlinkSession struct {
	Contract     *Chainlink        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainlinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainlinkCallerSession struct {
	Contract *ChainlinkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ChainlinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainlinkTransactorSession struct {
	Contract     *ChainlinkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ChainlinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainlinkRaw struct {
	Contract *Chainlink // Generic contract binding to access the raw methods on
}

// ChainlinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainlinkCallerRaw struct {
	Contract *ChainlinkCaller // Generic read-only contract binding to access the raw methods on
}

// ChainlinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainlinkTransactorRaw struct {
	Contract *ChainlinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainlink creates a new instance of Chainlink, bound to a specific deployed contract.
func NewChainlink(address common.Address, backend bind.ContractBackend) (*Chainlink, error) {
	contract, err := bindChainlink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chainlink{ChainlinkCaller: ChainlinkCaller{contract: contract}, ChainlinkTransactor: ChainlinkTransactor{contract: contract}, ChainlinkFilterer: ChainlinkFilterer{contract: contract}}, nil
}

// NewChainlinkCaller creates a new read-only instance of Chainlink, bound to a specific deployed contract.
func NewChainlinkCaller(address common.Address, caller bind.ContractCaller) (*ChainlinkCaller, error) {
	contract, err := bindChainlink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkCaller{contract: contract}, nil
}

// NewChainlinkTransactor creates a new write-only instance of Chainlink, bound to a specific deployed contract.
func NewChainlinkTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainlinkTransactor, error) {
	contract, err := bindChainlink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkTransactor{contract: contract}, nil
}

// NewChainlinkFilterer creates a new log filterer instance of Chainlink, bound to a specific deployed contract.
func NewChainlinkFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainlinkFilterer, error) {
	contract, err := bindChainlink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainlinkFilterer{contract: contract}, nil
}

// bindChainlink binds a generic wrapper to an already deployed contract.
func bindChainlink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainlinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chainlink *ChainlinkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chainlink.Contract.ChainlinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chainlink *ChainlinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chainlink.Contract.ChainlinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chainlink *ChainlinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chainlink.Contract.ChainlinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chainlink *ChainlinkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chainlink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chainlink *ChainlinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chainlink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chainlink *ChainlinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chainlink.Contract.contract.Transact(opts, method, params...)
}
