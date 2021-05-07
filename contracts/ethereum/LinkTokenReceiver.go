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

// LinkTokenReceiverABI is the input ABI used to generate the binding from.
const LinkTokenReceiverABI = "[{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// LinkTokenReceiver is an auto generated Go binding around an Ethereum contract.
type LinkTokenReceiver struct {
	LinkTokenReceiverCaller     // Read-only binding to the contract
	LinkTokenReceiverTransactor // Write-only binding to the contract
	LinkTokenReceiverFilterer   // Log filterer for contract events
}

// LinkTokenReceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type LinkTokenReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenReceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LinkTokenReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenReceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LinkTokenReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenReceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LinkTokenReceiverSession struct {
	Contract     *LinkTokenReceiver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LinkTokenReceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LinkTokenReceiverCallerSession struct {
	Contract *LinkTokenReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// LinkTokenReceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LinkTokenReceiverTransactorSession struct {
	Contract     *LinkTokenReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// LinkTokenReceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type LinkTokenReceiverRaw struct {
	Contract *LinkTokenReceiver // Generic contract binding to access the raw methods on
}

// LinkTokenReceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LinkTokenReceiverCallerRaw struct {
	Contract *LinkTokenReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// LinkTokenReceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LinkTokenReceiverTransactorRaw struct {
	Contract *LinkTokenReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLinkTokenReceiver creates a new instance of LinkTokenReceiver, bound to a specific deployed contract.
func NewLinkTokenReceiver(address common.Address, backend bind.ContractBackend) (*LinkTokenReceiver, error) {
	contract, err := bindLinkTokenReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiver{LinkTokenReceiverCaller: LinkTokenReceiverCaller{contract: contract}, LinkTokenReceiverTransactor: LinkTokenReceiverTransactor{contract: contract}, LinkTokenReceiverFilterer: LinkTokenReceiverFilterer{contract: contract}}, nil
}

// NewLinkTokenReceiverCaller creates a new read-only instance of LinkTokenReceiver, bound to a specific deployed contract.
func NewLinkTokenReceiverCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenReceiverCaller, error) {
	contract, err := bindLinkTokenReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiverCaller{contract: contract}, nil
}

// NewLinkTokenReceiverTransactor creates a new write-only instance of LinkTokenReceiver, bound to a specific deployed contract.
func NewLinkTokenReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenReceiverTransactor, error) {
	contract, err := bindLinkTokenReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiverTransactor{contract: contract}, nil
}

// NewLinkTokenReceiverFilterer creates a new log filterer instance of LinkTokenReceiver, bound to a specific deployed contract.
func NewLinkTokenReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenReceiverFilterer, error) {
	contract, err := bindLinkTokenReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiverFilterer{contract: contract}, nil
}

// bindLinkTokenReceiver binds a generic wrapper to an already deployed contract.
func bindLinkTokenReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LinkTokenReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenReceiver *LinkTokenReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenReceiver.Contract.LinkTokenReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenReceiver *LinkTokenReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.LinkTokenReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenReceiver *LinkTokenReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.LinkTokenReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenReceiver *LinkTokenReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenReceiver *LinkTokenReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenReceiver *LinkTokenReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.contract.Transact(opts, method, params...)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_LinkTokenReceiver *LinkTokenReceiverCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LinkTokenReceiver.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_LinkTokenReceiver *LinkTokenReceiverSession) GetChainlinkToken() (common.Address, error) {
	return _LinkTokenReceiver.Contract.GetChainlinkToken(&_LinkTokenReceiver.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_LinkTokenReceiver *LinkTokenReceiverCallerSession) GetChainlinkToken() (common.Address, error) {
	return _LinkTokenReceiver.Contract.GetChainlinkToken(&_LinkTokenReceiver.CallOpts)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_LinkTokenReceiver *LinkTokenReceiverTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _LinkTokenReceiver.contract.Transact(opts, "onTokenTransfer", _sender, _amount, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_LinkTokenReceiver *LinkTokenReceiverSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.OnTokenTransfer(&_LinkTokenReceiver.TransactOpts, _sender, _amount, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_LinkTokenReceiver *LinkTokenReceiverTransactorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.OnTokenTransfer(&_LinkTokenReceiver.TransactOpts, _sender, _amount, _data)
}
