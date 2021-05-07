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

// CBORChainlinkABI is the input ABI used to generate the binding from.
const CBORChainlinkABI = "[]"

// CBORChainlinkBin is the compiled bytecode used for deploying new contracts.
var CBORChainlinkBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212205e0c6de0498282a16303fb46df61437ac492720506ceadf7a889f14ca90a54f564736f6c63430006060033"

// DeployCBORChainlink deploys a new Ethereum contract, binding an instance of CBORChainlink to it.
func DeployCBORChainlink(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CBORChainlink, error) {
	parsed, err := abi.JSON(strings.NewReader(CBORChainlinkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CBORChainlinkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CBORChainlink{CBORChainlinkCaller: CBORChainlinkCaller{contract: contract}, CBORChainlinkTransactor: CBORChainlinkTransactor{contract: contract}, CBORChainlinkFilterer: CBORChainlinkFilterer{contract: contract}}, nil
}

// CBORChainlink is an auto generated Go binding around an Ethereum contract.
type CBORChainlink struct {
	CBORChainlinkCaller     // Read-only binding to the contract
	CBORChainlinkTransactor // Write-only binding to the contract
	CBORChainlinkFilterer   // Log filterer for contract events
}

// CBORChainlinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type CBORChainlinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CBORChainlinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CBORChainlinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CBORChainlinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CBORChainlinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CBORChainlinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CBORChainlinkSession struct {
	Contract     *CBORChainlink    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CBORChainlinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CBORChainlinkCallerSession struct {
	Contract *CBORChainlinkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// CBORChainlinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CBORChainlinkTransactorSession struct {
	Contract     *CBORChainlinkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// CBORChainlinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type CBORChainlinkRaw struct {
	Contract *CBORChainlink // Generic contract binding to access the raw methods on
}

// CBORChainlinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CBORChainlinkCallerRaw struct {
	Contract *CBORChainlinkCaller // Generic read-only contract binding to access the raw methods on
}

// CBORChainlinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CBORChainlinkTransactorRaw struct {
	Contract *CBORChainlinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCBORChainlink creates a new instance of CBORChainlink, bound to a specific deployed contract.
func NewCBORChainlink(address common.Address, backend bind.ContractBackend) (*CBORChainlink, error) {
	contract, err := bindCBORChainlink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CBORChainlink{CBORChainlinkCaller: CBORChainlinkCaller{contract: contract}, CBORChainlinkTransactor: CBORChainlinkTransactor{contract: contract}, CBORChainlinkFilterer: CBORChainlinkFilterer{contract: contract}}, nil
}

// NewCBORChainlinkCaller creates a new read-only instance of CBORChainlink, bound to a specific deployed contract.
func NewCBORChainlinkCaller(address common.Address, caller bind.ContractCaller) (*CBORChainlinkCaller, error) {
	contract, err := bindCBORChainlink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CBORChainlinkCaller{contract: contract}, nil
}

// NewCBORChainlinkTransactor creates a new write-only instance of CBORChainlink, bound to a specific deployed contract.
func NewCBORChainlinkTransactor(address common.Address, transactor bind.ContractTransactor) (*CBORChainlinkTransactor, error) {
	contract, err := bindCBORChainlink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CBORChainlinkTransactor{contract: contract}, nil
}

// NewCBORChainlinkFilterer creates a new log filterer instance of CBORChainlink, bound to a specific deployed contract.
func NewCBORChainlinkFilterer(address common.Address, filterer bind.ContractFilterer) (*CBORChainlinkFilterer, error) {
	contract, err := bindCBORChainlink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CBORChainlinkFilterer{contract: contract}, nil
}

// bindCBORChainlink binds a generic wrapper to an already deployed contract.
func bindCBORChainlink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CBORChainlinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CBORChainlink *CBORChainlinkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CBORChainlink.Contract.CBORChainlinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CBORChainlink *CBORChainlinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CBORChainlink.Contract.CBORChainlinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CBORChainlink *CBORChainlinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CBORChainlink.Contract.CBORChainlinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CBORChainlink *CBORChainlinkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CBORChainlink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CBORChainlink *CBORChainlinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CBORChainlink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CBORChainlink *CBORChainlinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CBORChainlink.Contract.contract.Transact(opts, method, params...)
}
