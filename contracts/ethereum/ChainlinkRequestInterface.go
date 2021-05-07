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

// ChainlinkRequestInterfaceABI is the input ABI used to generate the binding from.
const ChainlinkRequestInterfaceABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestPrice\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"serviceAgreementID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"oracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ChainlinkRequestInterface is an auto generated Go binding around an Ethereum contract.
type ChainlinkRequestInterface struct {
	ChainlinkRequestInterfaceCaller     // Read-only binding to the contract
	ChainlinkRequestInterfaceTransactor // Write-only binding to the contract
	ChainlinkRequestInterfaceFilterer   // Log filterer for contract events
}

// ChainlinkRequestInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainlinkRequestInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkRequestInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainlinkRequestInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkRequestInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainlinkRequestInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkRequestInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainlinkRequestInterfaceSession struct {
	Contract     *ChainlinkRequestInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainlinkRequestInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainlinkRequestInterfaceCallerSession struct {
	Contract *ChainlinkRequestInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ChainlinkRequestInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainlinkRequestInterfaceTransactorSession struct {
	Contract     *ChainlinkRequestInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ChainlinkRequestInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainlinkRequestInterfaceRaw struct {
	Contract *ChainlinkRequestInterface // Generic contract binding to access the raw methods on
}

// ChainlinkRequestInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainlinkRequestInterfaceCallerRaw struct {
	Contract *ChainlinkRequestInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// ChainlinkRequestInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainlinkRequestInterfaceTransactorRaw struct {
	Contract *ChainlinkRequestInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainlinkRequestInterface creates a new instance of ChainlinkRequestInterface, bound to a specific deployed contract.
func NewChainlinkRequestInterface(address common.Address, backend bind.ContractBackend) (*ChainlinkRequestInterface, error) {
	contract, err := bindChainlinkRequestInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainlinkRequestInterface{ChainlinkRequestInterfaceCaller: ChainlinkRequestInterfaceCaller{contract: contract}, ChainlinkRequestInterfaceTransactor: ChainlinkRequestInterfaceTransactor{contract: contract}, ChainlinkRequestInterfaceFilterer: ChainlinkRequestInterfaceFilterer{contract: contract}}, nil
}

// NewChainlinkRequestInterfaceCaller creates a new read-only instance of ChainlinkRequestInterface, bound to a specific deployed contract.
func NewChainlinkRequestInterfaceCaller(address common.Address, caller bind.ContractCaller) (*ChainlinkRequestInterfaceCaller, error) {
	contract, err := bindChainlinkRequestInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkRequestInterfaceCaller{contract: contract}, nil
}

// NewChainlinkRequestInterfaceTransactor creates a new write-only instance of ChainlinkRequestInterface, bound to a specific deployed contract.
func NewChainlinkRequestInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainlinkRequestInterfaceTransactor, error) {
	contract, err := bindChainlinkRequestInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkRequestInterfaceTransactor{contract: contract}, nil
}

// NewChainlinkRequestInterfaceFilterer creates a new log filterer instance of ChainlinkRequestInterface, bound to a specific deployed contract.
func NewChainlinkRequestInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainlinkRequestInterfaceFilterer, error) {
	contract, err := bindChainlinkRequestInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainlinkRequestInterfaceFilterer{contract: contract}, nil
}

// bindChainlinkRequestInterface binds a generic wrapper to an already deployed contract.
func bindChainlinkRequestInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainlinkRequestInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainlinkRequestInterface.Contract.ChainlinkRequestInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.ChainlinkRequestInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.ChainlinkRequestInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainlinkRequestInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.contract.Transact(opts, method, params...)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 requestId, uint256 payment, bytes4 callbackFunctionId, uint256 expiration) returns()
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceTransactor) CancelOracleRequest(opts *bind.TransactOpts, requestId [32]byte, payment *big.Int, callbackFunctionId [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.contract.Transact(opts, "cancelOracleRequest", requestId, payment, callbackFunctionId, expiration)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 requestId, uint256 payment, bytes4 callbackFunctionId, uint256 expiration) returns()
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceSession) CancelOracleRequest(requestId [32]byte, payment *big.Int, callbackFunctionId [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.CancelOracleRequest(&_ChainlinkRequestInterface.TransactOpts, requestId, payment, callbackFunctionId, expiration)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 requestId, uint256 payment, bytes4 callbackFunctionId, uint256 expiration) returns()
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceTransactorSession) CancelOracleRequest(requestId [32]byte, payment *big.Int, callbackFunctionId [4]byte, expiration *big.Int) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.CancelOracleRequest(&_ChainlinkRequestInterface.TransactOpts, requestId, payment, callbackFunctionId, expiration)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address sender, uint256 requestPrice, bytes32 serviceAgreementID, address callbackAddress, bytes4 callbackFunctionId, uint256 nonce, uint256 dataVersion, bytes data) returns()
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceTransactor) OracleRequest(opts *bind.TransactOpts, sender common.Address, requestPrice *big.Int, serviceAgreementID [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.contract.Transact(opts, "oracleRequest", sender, requestPrice, serviceAgreementID, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address sender, uint256 requestPrice, bytes32 serviceAgreementID, address callbackAddress, bytes4 callbackFunctionId, uint256 nonce, uint256 dataVersion, bytes data) returns()
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceSession) OracleRequest(sender common.Address, requestPrice *big.Int, serviceAgreementID [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.OracleRequest(&_ChainlinkRequestInterface.TransactOpts, sender, requestPrice, serviceAgreementID, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address sender, uint256 requestPrice, bytes32 serviceAgreementID, address callbackAddress, bytes4 callbackFunctionId, uint256 nonce, uint256 dataVersion, bytes data) returns()
func (_ChainlinkRequestInterface *ChainlinkRequestInterfaceTransactorSession) OracleRequest(sender common.Address, requestPrice *big.Int, serviceAgreementID [32]byte, callbackAddress common.Address, callbackFunctionId [4]byte, nonce *big.Int, dataVersion *big.Int, data []byte) (*types.Transaction, error) {
	return _ChainlinkRequestInterface.Contract.OracleRequest(&_ChainlinkRequestInterface.TransactOpts, sender, requestPrice, serviceAgreementID, callbackAddress, callbackFunctionId, nonce, dataVersion, data)
}
