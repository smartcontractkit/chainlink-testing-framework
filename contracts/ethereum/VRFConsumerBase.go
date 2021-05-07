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

// VRFConsumerBaseABI is the input ABI used to generate the binding from.
const VRFConsumerBaseABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// VRFConsumerBase is an auto generated Go binding around an Ethereum contract.
type VRFConsumerBase struct {
	VRFConsumerBaseCaller     // Read-only binding to the contract
	VRFConsumerBaseTransactor // Write-only binding to the contract
	VRFConsumerBaseFilterer   // Log filterer for contract events
}

// VRFConsumerBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFConsumerBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFConsumerBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFConsumerBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFConsumerBaseSession struct {
	Contract     *VRFConsumerBase  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFConsumerBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFConsumerBaseCallerSession struct {
	Contract *VRFConsumerBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// VRFConsumerBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFConsumerBaseTransactorSession struct {
	Contract     *VRFConsumerBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// VRFConsumerBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFConsumerBaseRaw struct {
	Contract *VRFConsumerBase // Generic contract binding to access the raw methods on
}

// VRFConsumerBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFConsumerBaseCallerRaw struct {
	Contract *VRFConsumerBaseCaller // Generic read-only contract binding to access the raw methods on
}

// VRFConsumerBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFConsumerBaseTransactorRaw struct {
	Contract *VRFConsumerBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFConsumerBase creates a new instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBase(address common.Address, backend bind.ContractBackend) (*VRFConsumerBase, error) {
	contract, err := bindVRFConsumerBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBase{VRFConsumerBaseCaller: VRFConsumerBaseCaller{contract: contract}, VRFConsumerBaseTransactor: VRFConsumerBaseTransactor{contract: contract}, VRFConsumerBaseFilterer: VRFConsumerBaseFilterer{contract: contract}}, nil
}

// NewVRFConsumerBaseCaller creates a new read-only instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBaseCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerBaseCaller, error) {
	contract, err := bindVRFConsumerBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBaseCaller{contract: contract}, nil
}

// NewVRFConsumerBaseTransactor creates a new write-only instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerBaseTransactor, error) {
	contract, err := bindVRFConsumerBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBaseTransactor{contract: contract}, nil
}

// NewVRFConsumerBaseFilterer creates a new log filterer instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerBaseFilterer, error) {
	contract, err := bindVRFConsumerBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBaseFilterer{contract: contract}, nil
}

// bindVRFConsumerBase binds a generic wrapper to an already deployed contract.
func bindVRFConsumerBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumerBase *VRFConsumerBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerBase.Contract.VRFConsumerBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumerBase *VRFConsumerBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.VRFConsumerBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumerBase *VRFConsumerBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.VRFConsumerBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumerBase *VRFConsumerBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumerBase *VRFConsumerBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumerBase *VRFConsumerBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.contract.Transact(opts, method, params...)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumerBase *VRFConsumerBaseTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumerBase.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumerBase *VRFConsumerBaseSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.RawFulfillRandomness(&_VRFConsumerBase.TransactOpts, requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumerBase *VRFConsumerBaseTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.RawFulfillRandomness(&_VRFConsumerBase.TransactOpts, requestId, randomness)
}
