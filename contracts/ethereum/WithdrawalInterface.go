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

// WithdrawalInterfaceABI is the input ABI used to generate the binding from.
const WithdrawalInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// WithdrawalInterface is an auto generated Go binding around an Ethereum contract.
type WithdrawalInterface struct {
	WithdrawalInterfaceCaller     // Read-only binding to the contract
	WithdrawalInterfaceTransactor // Write-only binding to the contract
	WithdrawalInterfaceFilterer   // Log filterer for contract events
}

// WithdrawalInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithdrawalInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawalInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithdrawalInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawalInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithdrawalInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawalInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithdrawalInterfaceSession struct {
	Contract     *WithdrawalInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// WithdrawalInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithdrawalInterfaceCallerSession struct {
	Contract *WithdrawalInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// WithdrawalInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithdrawalInterfaceTransactorSession struct {
	Contract     *WithdrawalInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// WithdrawalInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithdrawalInterfaceRaw struct {
	Contract *WithdrawalInterface // Generic contract binding to access the raw methods on
}

// WithdrawalInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithdrawalInterfaceCallerRaw struct {
	Contract *WithdrawalInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// WithdrawalInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithdrawalInterfaceTransactorRaw struct {
	Contract *WithdrawalInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithdrawalInterface creates a new instance of WithdrawalInterface, bound to a specific deployed contract.
func NewWithdrawalInterface(address common.Address, backend bind.ContractBackend) (*WithdrawalInterface, error) {
	contract, err := bindWithdrawalInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WithdrawalInterface{WithdrawalInterfaceCaller: WithdrawalInterfaceCaller{contract: contract}, WithdrawalInterfaceTransactor: WithdrawalInterfaceTransactor{contract: contract}, WithdrawalInterfaceFilterer: WithdrawalInterfaceFilterer{contract: contract}}, nil
}

// NewWithdrawalInterfaceCaller creates a new read-only instance of WithdrawalInterface, bound to a specific deployed contract.
func NewWithdrawalInterfaceCaller(address common.Address, caller bind.ContractCaller) (*WithdrawalInterfaceCaller, error) {
	contract, err := bindWithdrawalInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawalInterfaceCaller{contract: contract}, nil
}

// NewWithdrawalInterfaceTransactor creates a new write-only instance of WithdrawalInterface, bound to a specific deployed contract.
func NewWithdrawalInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*WithdrawalInterfaceTransactor, error) {
	contract, err := bindWithdrawalInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawalInterfaceTransactor{contract: contract}, nil
}

// NewWithdrawalInterfaceFilterer creates a new log filterer instance of WithdrawalInterface, bound to a specific deployed contract.
func NewWithdrawalInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*WithdrawalInterfaceFilterer, error) {
	contract, err := bindWithdrawalInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithdrawalInterfaceFilterer{contract: contract}, nil
}

// bindWithdrawalInterface binds a generic wrapper to an already deployed contract.
func bindWithdrawalInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WithdrawalInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithdrawalInterface *WithdrawalInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WithdrawalInterface.Contract.WithdrawalInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithdrawalInterface *WithdrawalInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithdrawalInterface.Contract.WithdrawalInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithdrawalInterface *WithdrawalInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithdrawalInterface.Contract.WithdrawalInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithdrawalInterface *WithdrawalInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WithdrawalInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithdrawalInterface *WithdrawalInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithdrawalInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithdrawalInterface *WithdrawalInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithdrawalInterface.Contract.contract.Transact(opts, method, params...)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_WithdrawalInterface *WithdrawalInterfaceCaller) Withdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WithdrawalInterface.contract.Call(opts, &out, "withdrawable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_WithdrawalInterface *WithdrawalInterfaceSession) Withdrawable() (*big.Int, error) {
	return _WithdrawalInterface.Contract.Withdrawable(&_WithdrawalInterface.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_WithdrawalInterface *WithdrawalInterfaceCallerSession) Withdrawable() (*big.Int, error) {
	return _WithdrawalInterface.Contract.Withdrawable(&_WithdrawalInterface.CallOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_WithdrawalInterface *WithdrawalInterfaceTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WithdrawalInterface.contract.Transact(opts, "withdraw", recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_WithdrawalInterface *WithdrawalInterfaceSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WithdrawalInterface.Contract.Withdraw(&_WithdrawalInterface.TransactOpts, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_WithdrawalInterface *WithdrawalInterfaceTransactorSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WithdrawalInterface.Contract.Withdraw(&_WithdrawalInterface.TransactOpts, recipient, amount)
}
