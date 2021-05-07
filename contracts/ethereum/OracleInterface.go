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

// OracleInterfaceABI is the input ABI used to generate the binding from.
const OracleInterfaceABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"fulfillOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"getAuthorizationStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setFulfillmentPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OracleInterface is an auto generated Go binding around an Ethereum contract.
type OracleInterface struct {
	OracleInterfaceCaller     // Read-only binding to the contract
	OracleInterfaceTransactor // Write-only binding to the contract
	OracleInterfaceFilterer   // Log filterer for contract events
}

// OracleInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleInterfaceSession struct {
	Contract     *OracleInterface  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleInterfaceCallerSession struct {
	Contract *OracleInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// OracleInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleInterfaceTransactorSession struct {
	Contract     *OracleInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// OracleInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleInterfaceRaw struct {
	Contract *OracleInterface // Generic contract binding to access the raw methods on
}

// OracleInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleInterfaceCallerRaw struct {
	Contract *OracleInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// OracleInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleInterfaceTransactorRaw struct {
	Contract *OracleInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleInterface creates a new instance of OracleInterface, bound to a specific deployed contract.
func NewOracleInterface(address common.Address, backend bind.ContractBackend) (*OracleInterface, error) {
	contract, err := bindOracleInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleInterface{OracleInterfaceCaller: OracleInterfaceCaller{contract: contract}, OracleInterfaceTransactor: OracleInterfaceTransactor{contract: contract}, OracleInterfaceFilterer: OracleInterfaceFilterer{contract: contract}}, nil
}

// NewOracleInterfaceCaller creates a new read-only instance of OracleInterface, bound to a specific deployed contract.
func NewOracleInterfaceCaller(address common.Address, caller bind.ContractCaller) (*OracleInterfaceCaller, error) {
	contract, err := bindOracleInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleInterfaceCaller{contract: contract}, nil
}

// NewOracleInterfaceTransactor creates a new write-only instance of OracleInterface, bound to a specific deployed contract.
func NewOracleInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleInterfaceTransactor, error) {
	contract, err := bindOracleInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleInterfaceTransactor{contract: contract}, nil
}

// NewOracleInterfaceFilterer creates a new log filterer instance of OracleInterface, bound to a specific deployed contract.
func NewOracleInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleInterfaceFilterer, error) {
	contract, err := bindOracleInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleInterfaceFilterer{contract: contract}, nil
}

// bindOracleInterface binds a generic wrapper to an already deployed contract.
func bindOracleInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleInterface *OracleInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleInterface.Contract.OracleInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleInterface *OracleInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleInterface.Contract.OracleInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleInterface *OracleInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleInterface.Contract.OracleInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleInterface *OracleInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleInterface *OracleInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleInterface *OracleInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleInterface.Contract.contract.Transact(opts, method, params...)
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address node) view returns(bool)
func (_OracleInterface *OracleInterfaceCaller) GetAuthorizationStatus(opts *bind.CallOpts, node common.Address) (bool, error) {
	var out []interface{}
	err := _OracleInterface.contract.Call(opts, &out, "getAuthorizationStatus", node)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address node) view returns(bool)
func (_OracleInterface *OracleInterfaceSession) GetAuthorizationStatus(node common.Address) (bool, error) {
	return _OracleInterface.Contract.GetAuthorizationStatus(&_OracleInterface.CallOpts, node)
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address node) view returns(bool)
func (_OracleInterface *OracleInterfaceCallerSession) GetAuthorizationStatus(node common.Address) (bool, error) {
	return _OracleInterface.Contract.GetAuthorizationStatus(&_OracleInterface.CallOpts, node)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_OracleInterface *OracleInterfaceCaller) Withdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleInterface.contract.Call(opts, &out, "withdrawable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_OracleInterface *OracleInterfaceSession) Withdrawable() (*big.Int, error) {
	return _OracleInterface.Contract.Withdrawable(&_OracleInterface.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_OracleInterface *OracleInterfaceCallerSession) Withdrawable() (*big.Int, error) {
	return _OracleInterface.Contract.Withdrawable(&_OracleInterface.CallOpts)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data) returns(bool)
func (_OracleInterface *OracleInterfaceTransactor) FulfillOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _OracleInterface.contract.Transact(opts, "fulfillOracleRequest", requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data) returns(bool)
func (_OracleInterface *OracleInterfaceSession) FulfillOracleRequest(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _OracleInterface.Contract.FulfillOracleRequest(&_OracleInterface.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 data) returns(bool)
func (_OracleInterface *OracleInterfaceTransactorSession) FulfillOracleRequest(requestId [32]byte, payment *big.Int, callbackAddress common.Address, callbackFunctionId [4]byte, expiration *big.Int, data [32]byte) (*types.Transaction, error) {
	return _OracleInterface.Contract.FulfillOracleRequest(&_OracleInterface.TransactOpts, requestId, payment, callbackAddress, callbackFunctionId, expiration, data)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address node, bool allowed) returns()
func (_OracleInterface *OracleInterfaceTransactor) SetFulfillmentPermission(opts *bind.TransactOpts, node common.Address, allowed bool) (*types.Transaction, error) {
	return _OracleInterface.contract.Transact(opts, "setFulfillmentPermission", node, allowed)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address node, bool allowed) returns()
func (_OracleInterface *OracleInterfaceSession) SetFulfillmentPermission(node common.Address, allowed bool) (*types.Transaction, error) {
	return _OracleInterface.Contract.SetFulfillmentPermission(&_OracleInterface.TransactOpts, node, allowed)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address node, bool allowed) returns()
func (_OracleInterface *OracleInterfaceTransactorSession) SetFulfillmentPermission(node common.Address, allowed bool) (*types.Transaction, error) {
	return _OracleInterface.Contract.SetFulfillmentPermission(&_OracleInterface.TransactOpts, node, allowed)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_OracleInterface *OracleInterfaceTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OracleInterface.contract.Transact(opts, "withdraw", recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_OracleInterface *OracleInterfaceSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OracleInterface.Contract.Withdraw(&_OracleInterface.TransactOpts, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address recipient, uint256 amount) returns()
func (_OracleInterface *OracleInterfaceTransactorSession) Withdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OracleInterface.Contract.Withdraw(&_OracleInterface.TransactOpts, recipient, amount)
}
