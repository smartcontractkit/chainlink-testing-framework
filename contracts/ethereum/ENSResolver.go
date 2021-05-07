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

// ENSResolverABI is the input ABI used to generate the binding from.
const ENSResolverABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"addr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ENSResolver is an auto generated Go binding around an Ethereum contract.
type ENSResolver struct {
	ENSResolverCaller     // Read-only binding to the contract
	ENSResolverTransactor // Write-only binding to the contract
	ENSResolverFilterer   // Log filterer for contract events
}

// ENSResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type ENSResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ENSResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ENSResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ENSResolverSession struct {
	Contract     *ENSResolver      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ENSResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ENSResolverCallerSession struct {
	Contract *ENSResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ENSResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ENSResolverTransactorSession struct {
	Contract     *ENSResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ENSResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type ENSResolverRaw struct {
	Contract *ENSResolver // Generic contract binding to access the raw methods on
}

// ENSResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ENSResolverCallerRaw struct {
	Contract *ENSResolverCaller // Generic read-only contract binding to access the raw methods on
}

// ENSResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ENSResolverTransactorRaw struct {
	Contract *ENSResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewENSResolver creates a new instance of ENSResolver, bound to a specific deployed contract.
func NewENSResolver(address common.Address, backend bind.ContractBackend) (*ENSResolver, error) {
	contract, err := bindENSResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ENSResolver{ENSResolverCaller: ENSResolverCaller{contract: contract}, ENSResolverTransactor: ENSResolverTransactor{contract: contract}, ENSResolverFilterer: ENSResolverFilterer{contract: contract}}, nil
}

// NewENSResolverCaller creates a new read-only instance of ENSResolver, bound to a specific deployed contract.
func NewENSResolverCaller(address common.Address, caller bind.ContractCaller) (*ENSResolverCaller, error) {
	contract, err := bindENSResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ENSResolverCaller{contract: contract}, nil
}

// NewENSResolverTransactor creates a new write-only instance of ENSResolver, bound to a specific deployed contract.
func NewENSResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*ENSResolverTransactor, error) {
	contract, err := bindENSResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ENSResolverTransactor{contract: contract}, nil
}

// NewENSResolverFilterer creates a new log filterer instance of ENSResolver, bound to a specific deployed contract.
func NewENSResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*ENSResolverFilterer, error) {
	contract, err := bindENSResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ENSResolverFilterer{contract: contract}, nil
}

// bindENSResolver binds a generic wrapper to an already deployed contract.
func bindENSResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSResolverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ENSResolver *ENSResolverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ENSResolver.Contract.ENSResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ENSResolver *ENSResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENSResolver.Contract.ENSResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ENSResolver *ENSResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ENSResolver.Contract.ENSResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ENSResolver *ENSResolverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ENSResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ENSResolver *ENSResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENSResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ENSResolver *ENSResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ENSResolver.Contract.contract.Transact(opts, method, params...)
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(bytes32 node) view returns(address)
func (_ENSResolver *ENSResolverCaller) Addr(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ENSResolver.contract.Call(opts, &out, "addr", node)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(bytes32 node) view returns(address)
func (_ENSResolver *ENSResolverSession) Addr(node [32]byte) (common.Address, error) {
	return _ENSResolver.Contract.Addr(&_ENSResolver.CallOpts, node)
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(bytes32 node) view returns(address)
func (_ENSResolver *ENSResolverCallerSession) Addr(node [32]byte) (common.Address, error) {
	return _ENSResolver.Contract.Addr(&_ENSResolver.CallOpts, node)
}
