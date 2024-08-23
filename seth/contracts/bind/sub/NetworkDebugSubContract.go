// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package network_debug_sub_contract

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
	_ = abi.ConvertType
)

// NetworkDebugSubContractAccount is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugSubContractAccount struct {
	Name       string
	Balance    uint64
	DailyLimit *big.Int
}

// NetworkDebugSubContractMetaData contains all meta data concerning the NetworkDebugSubContract contract.
var NetworkDebugSubContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"required\",\"type\":\"uint256\"}],\"name\":\"CustomErr\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NoIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"NoIndexEventString\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"balance\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"dailyLimit\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structNetworkDebugSubContract.Account\",\"name\":\"a\",\"type\":\"tuple\"}],\"name\":\"NoIndexStructEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"OneIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"ThreeIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"}],\"name\":\"TwoIndexEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"alwaysRevertsCustomError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pay\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"trace\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"traceOneInt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"r\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceWithCallback\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506105f4806100206000396000f3fe60806040526004361061004a5760003560e01c8063047c44251461004f57806311abb0021461008c5780631b9265b8146100b55780633e41f135146100bf578063fa8fca7a146100fc575b600080fd5b34801561005b57600080fd5b5061007660048036038101906100719190610368565b610139565b60405161008391906103a4565b60405180910390f35b34801561009857600080fd5b506100b360048036038101906100ae91906103f5565b610186565b005b6100bd6101c5565b005b3480156100cb57600080fd5b506100e660048036038101906100e19190610435565b6101c7565b6040516100f391906103a4565b60405180910390f35b34801561010857600080fd5b50610123600480360381019061011e9190610435565b610230565b60405161013091906103a4565b60405180910390f35b60007f33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c3360405161016a91906104b6565b60405180910390a160038261017f9190610500565b9050919050565b81816040517f4a2eaf7e0000000000000000000000000000000000000000000000000000000081526004016101bc929190610553565b60405180910390fd5b565b60006002826101d69190610500565b91503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a381836102289190610500565b905092915050565b60003373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a360003373ffffffffffffffffffffffffffffffffffffffff1663fbcb8d07846040518263ffffffff1660e01b81526004016102b191906103a4565b6020604051808303816000875af11580156102d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f49190610591565b9050807feace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b3560405160405180910390a28091505092915050565b600080fd5b6000819050919050565b61034581610332565b811461035057600080fd5b50565b6000813590506103628161033c565b92915050565b60006020828403121561037e5761037d61032d565b5b600061038c84828501610353565b91505092915050565b61039e81610332565b82525050565b60006020820190506103b96000830184610395565b92915050565b6000819050919050565b6103d2816103bf565b81146103dd57600080fd5b50565b6000813590506103ef816103c9565b92915050565b6000806040838503121561040c5761040b61032d565b5b600061041a858286016103e0565b925050602061042b858286016103e0565b9150509250929050565b6000806040838503121561044c5761044b61032d565b5b600061045a85828601610353565b925050602061046b85828601610353565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006104a082610475565b9050919050565b6104b081610495565b82525050565b60006020820190506104cb60008301846104a7565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061050b82610332565b915061051683610332565b92508282019050828112156000831216838212600084121516171561053e5761053d6104d1565b5b92915050565b61054d816103bf565b82525050565b60006040820190506105686000830185610544565b6105756020830184610544565b9392505050565b60008151905061058b8161033c565b92915050565b6000602082840312156105a7576105a661032d565b5b60006105b58482850161057c565b9150509291505056fea2646970667358221220e1f967a055585aa7f00afcc65e7374aca7ee992f6da265b4b82d7c2b38f461c864736f6c63430008130033",
}

// NetworkDebugSubContractABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkDebugSubContractMetaData.ABI instead.
var NetworkDebugSubContractABI = NetworkDebugSubContractMetaData.ABI

// NetworkDebugSubContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkDebugSubContractMetaData.Bin instead.
var NetworkDebugSubContractBin = NetworkDebugSubContractMetaData.Bin

// DeployNetworkDebugSubContract deploys a new Ethereum contract, binding an instance of NetworkDebugSubContract to it.
func DeployNetworkDebugSubContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NetworkDebugSubContract, error) {
	parsed, err := NetworkDebugSubContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkDebugSubContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkDebugSubContract{NetworkDebugSubContractCaller: NetworkDebugSubContractCaller{contract: contract}, NetworkDebugSubContractTransactor: NetworkDebugSubContractTransactor{contract: contract}, NetworkDebugSubContractFilterer: NetworkDebugSubContractFilterer{contract: contract}}, nil
}

// NetworkDebugSubContract is an auto generated Go binding around an Ethereum contract.
type NetworkDebugSubContract struct {
	NetworkDebugSubContractCaller     // Read-only binding to the contract
	NetworkDebugSubContractTransactor // Write-only binding to the contract
	NetworkDebugSubContractFilterer   // Log filterer for contract events
}

// NetworkDebugSubContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkDebugSubContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugSubContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkDebugSubContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugSubContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkDebugSubContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugSubContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkDebugSubContractSession struct {
	Contract     *NetworkDebugSubContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// NetworkDebugSubContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkDebugSubContractCallerSession struct {
	Contract *NetworkDebugSubContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// NetworkDebugSubContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkDebugSubContractTransactorSession struct {
	Contract     *NetworkDebugSubContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// NetworkDebugSubContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkDebugSubContractRaw struct {
	Contract *NetworkDebugSubContract // Generic contract binding to access the raw methods on
}

// NetworkDebugSubContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkDebugSubContractCallerRaw struct {
	Contract *NetworkDebugSubContractCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkDebugSubContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkDebugSubContractTransactorRaw struct {
	Contract *NetworkDebugSubContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkDebugSubContract creates a new instance of NetworkDebugSubContract, bound to a specific deployed contract.
func NewNetworkDebugSubContract(address common.Address, backend bind.ContractBackend) (*NetworkDebugSubContract, error) {
	contract, err := bindNetworkDebugSubContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContract{NetworkDebugSubContractCaller: NetworkDebugSubContractCaller{contract: contract}, NetworkDebugSubContractTransactor: NetworkDebugSubContractTransactor{contract: contract}, NetworkDebugSubContractFilterer: NetworkDebugSubContractFilterer{contract: contract}}, nil
}

// NewNetworkDebugSubContractCaller creates a new read-only instance of NetworkDebugSubContract, bound to a specific deployed contract.
func NewNetworkDebugSubContractCaller(address common.Address, caller bind.ContractCaller) (*NetworkDebugSubContractCaller, error) {
	contract, err := bindNetworkDebugSubContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractCaller{contract: contract}, nil
}

// NewNetworkDebugSubContractTransactor creates a new write-only instance of NetworkDebugSubContract, bound to a specific deployed contract.
func NewNetworkDebugSubContractTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkDebugSubContractTransactor, error) {
	contract, err := bindNetworkDebugSubContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractTransactor{contract: contract}, nil
}

// NewNetworkDebugSubContractFilterer creates a new log filterer instance of NetworkDebugSubContract, bound to a specific deployed contract.
func NewNetworkDebugSubContractFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkDebugSubContractFilterer, error) {
	contract, err := bindNetworkDebugSubContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractFilterer{contract: contract}, nil
}

// bindNetworkDebugSubContract binds a generic wrapper to an already deployed contract.
func bindNetworkDebugSubContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkDebugSubContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkDebugSubContract *NetworkDebugSubContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkDebugSubContract.Contract.NetworkDebugSubContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkDebugSubContract *NetworkDebugSubContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.NetworkDebugSubContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkDebugSubContract *NetworkDebugSubContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.NetworkDebugSubContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkDebugSubContract *NetworkDebugSubContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkDebugSubContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.contract.Transact(opts, method, params...)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x11abb002.
//
// Solidity: function alwaysRevertsCustomError(uint256 x, uint256 y) returns()
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactor) AlwaysRevertsCustomError(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.contract.Transact(opts, "alwaysRevertsCustomError", x, y)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x11abb002.
//
// Solidity: function alwaysRevertsCustomError(uint256 x, uint256 y) returns()
func (_NetworkDebugSubContract *NetworkDebugSubContractSession) AlwaysRevertsCustomError(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.AlwaysRevertsCustomError(&_NetworkDebugSubContract.TransactOpts, x, y)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x11abb002.
//
// Solidity: function alwaysRevertsCustomError(uint256 x, uint256 y) returns()
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorSession) AlwaysRevertsCustomError(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.AlwaysRevertsCustomError(&_NetworkDebugSubContract.TransactOpts, x, y)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactor) Pay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugSubContract.contract.Transact(opts, "pay")
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugSubContract *NetworkDebugSubContractSession) Pay() (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.Pay(&_NetworkDebugSubContract.TransactOpts)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorSession) Pay() (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.Pay(&_NetworkDebugSubContract.TransactOpts)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactor) Trace(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.contract.Transact(opts, "trace", x, y)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugSubContract *NetworkDebugSubContractSession) Trace(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.Trace(&_NetworkDebugSubContract.TransactOpts, x, y)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorSession) Trace(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.Trace(&_NetworkDebugSubContract.TransactOpts, x, y)
}

// TraceOneInt is a paid mutator transaction binding the contract method 0x047c4425.
//
// Solidity: function traceOneInt(int256 x) returns(int256 r)
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactor) TraceOneInt(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.contract.Transact(opts, "traceOneInt", x)
}

// TraceOneInt is a paid mutator transaction binding the contract method 0x047c4425.
//
// Solidity: function traceOneInt(int256 x) returns(int256 r)
func (_NetworkDebugSubContract *NetworkDebugSubContractSession) TraceOneInt(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.TraceOneInt(&_NetworkDebugSubContract.TransactOpts, x)
}

// TraceOneInt is a paid mutator transaction binding the contract method 0x047c4425.
//
// Solidity: function traceOneInt(int256 x) returns(int256 r)
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorSession) TraceOneInt(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.TraceOneInt(&_NetworkDebugSubContract.TransactOpts, x)
}

// TraceWithCallback is a paid mutator transaction binding the contract method 0xfa8fca7a.
//
// Solidity: function traceWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactor) TraceWithCallback(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.contract.Transact(opts, "traceWithCallback", x, y)
}

// TraceWithCallback is a paid mutator transaction binding the contract method 0xfa8fca7a.
//
// Solidity: function traceWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugSubContract *NetworkDebugSubContractSession) TraceWithCallback(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.TraceWithCallback(&_NetworkDebugSubContract.TransactOpts, x, y)
}

// TraceWithCallback is a paid mutator transaction binding the contract method 0xfa8fca7a.
//
// Solidity: function traceWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugSubContract *NetworkDebugSubContractTransactorSession) TraceWithCallback(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugSubContract.Contract.TraceWithCallback(&_NetworkDebugSubContract.TransactOpts, x, y)
}

// NetworkDebugSubContractNoIndexEventIterator is returned from FilterNoIndexEvent and is used to iterate over the raw logs and unpacked data for NoIndexEvent events raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractNoIndexEventIterator struct {
	Event *NetworkDebugSubContractNoIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugSubContractNoIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugSubContractNoIndexEvent)
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
		it.Event = new(NetworkDebugSubContractNoIndexEvent)
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
func (it *NetworkDebugSubContractNoIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugSubContractNoIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugSubContractNoIndexEvent represents a NoIndexEvent event raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractNoIndexEvent struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNoIndexEvent is a free log retrieval operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) FilterNoIndexEvent(opts *bind.FilterOpts) (*NetworkDebugSubContractNoIndexEventIterator, error) {

	logs, sub, err := _NetworkDebugSubContract.contract.FilterLogs(opts, "NoIndexEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractNoIndexEventIterator{contract: _NetworkDebugSubContract.contract, event: "NoIndexEvent", logs: logs, sub: sub}, nil
}

// WatchNoIndexEvent is a free log subscription operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) WatchNoIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugSubContractNoIndexEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugSubContract.contract.WatchLogs(opts, "NoIndexEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugSubContractNoIndexEvent)
				if err := _NetworkDebugSubContract.contract.UnpackLog(event, "NoIndexEvent", log); err != nil {
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

// ParseNoIndexEvent is a log parse operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) ParseNoIndexEvent(log types.Log) (*NetworkDebugSubContractNoIndexEvent, error) {
	event := new(NetworkDebugSubContractNoIndexEvent)
	if err := _NetworkDebugSubContract.contract.UnpackLog(event, "NoIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugSubContractNoIndexEventStringIterator is returned from FilterNoIndexEventString and is used to iterate over the raw logs and unpacked data for NoIndexEventString events raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractNoIndexEventStringIterator struct {
	Event *NetworkDebugSubContractNoIndexEventString // Event containing the contract specifics and raw log

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
func (it *NetworkDebugSubContractNoIndexEventStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugSubContractNoIndexEventString)
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
		it.Event = new(NetworkDebugSubContractNoIndexEventString)
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
func (it *NetworkDebugSubContractNoIndexEventStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugSubContractNoIndexEventStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugSubContractNoIndexEventString represents a NoIndexEventString event raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractNoIndexEventString struct {
	Str string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoIndexEventString is a free log retrieval operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) FilterNoIndexEventString(opts *bind.FilterOpts) (*NetworkDebugSubContractNoIndexEventStringIterator, error) {

	logs, sub, err := _NetworkDebugSubContract.contract.FilterLogs(opts, "NoIndexEventString")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractNoIndexEventStringIterator{contract: _NetworkDebugSubContract.contract, event: "NoIndexEventString", logs: logs, sub: sub}, nil
}

// WatchNoIndexEventString is a free log subscription operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) WatchNoIndexEventString(opts *bind.WatchOpts, sink chan<- *NetworkDebugSubContractNoIndexEventString) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugSubContract.contract.WatchLogs(opts, "NoIndexEventString")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugSubContractNoIndexEventString)
				if err := _NetworkDebugSubContract.contract.UnpackLog(event, "NoIndexEventString", log); err != nil {
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

// ParseNoIndexEventString is a log parse operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) ParseNoIndexEventString(log types.Log) (*NetworkDebugSubContractNoIndexEventString, error) {
	event := new(NetworkDebugSubContractNoIndexEventString)
	if err := _NetworkDebugSubContract.contract.UnpackLog(event, "NoIndexEventString", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugSubContractNoIndexStructEventIterator is returned from FilterNoIndexStructEvent and is used to iterate over the raw logs and unpacked data for NoIndexStructEvent events raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractNoIndexStructEventIterator struct {
	Event *NetworkDebugSubContractNoIndexStructEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugSubContractNoIndexStructEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugSubContractNoIndexStructEvent)
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
		it.Event = new(NetworkDebugSubContractNoIndexStructEvent)
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
func (it *NetworkDebugSubContractNoIndexStructEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugSubContractNoIndexStructEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugSubContractNoIndexStructEvent represents a NoIndexStructEvent event raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractNoIndexStructEvent struct {
	A   NetworkDebugSubContractAccount
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoIndexStructEvent is a free log retrieval operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) FilterNoIndexStructEvent(opts *bind.FilterOpts) (*NetworkDebugSubContractNoIndexStructEventIterator, error) {

	logs, sub, err := _NetworkDebugSubContract.contract.FilterLogs(opts, "NoIndexStructEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractNoIndexStructEventIterator{contract: _NetworkDebugSubContract.contract, event: "NoIndexStructEvent", logs: logs, sub: sub}, nil
}

// WatchNoIndexStructEvent is a free log subscription operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) WatchNoIndexStructEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugSubContractNoIndexStructEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugSubContract.contract.WatchLogs(opts, "NoIndexStructEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugSubContractNoIndexStructEvent)
				if err := _NetworkDebugSubContract.contract.UnpackLog(event, "NoIndexStructEvent", log); err != nil {
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

// ParseNoIndexStructEvent is a log parse operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) ParseNoIndexStructEvent(log types.Log) (*NetworkDebugSubContractNoIndexStructEvent, error) {
	event := new(NetworkDebugSubContractNoIndexStructEvent)
	if err := _NetworkDebugSubContract.contract.UnpackLog(event, "NoIndexStructEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugSubContractOneIndexEventIterator is returned from FilterOneIndexEvent and is used to iterate over the raw logs and unpacked data for OneIndexEvent events raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractOneIndexEventIterator struct {
	Event *NetworkDebugSubContractOneIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugSubContractOneIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugSubContractOneIndexEvent)
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
		it.Event = new(NetworkDebugSubContractOneIndexEvent)
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
func (it *NetworkDebugSubContractOneIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugSubContractOneIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugSubContractOneIndexEvent represents a OneIndexEvent event raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractOneIndexEvent struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterOneIndexEvent is a free log retrieval operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) FilterOneIndexEvent(opts *bind.FilterOpts, a []*big.Int) (*NetworkDebugSubContractOneIndexEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugSubContract.contract.FilterLogs(opts, "OneIndexEvent", aRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractOneIndexEventIterator{contract: _NetworkDebugSubContract.contract, event: "OneIndexEvent", logs: logs, sub: sub}, nil
}

// WatchOneIndexEvent is a free log subscription operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) WatchOneIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugSubContractOneIndexEvent, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugSubContract.contract.WatchLogs(opts, "OneIndexEvent", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugSubContractOneIndexEvent)
				if err := _NetworkDebugSubContract.contract.UnpackLog(event, "OneIndexEvent", log); err != nil {
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

// ParseOneIndexEvent is a log parse operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) ParseOneIndexEvent(log types.Log) (*NetworkDebugSubContractOneIndexEvent, error) {
	event := new(NetworkDebugSubContractOneIndexEvent)
	if err := _NetworkDebugSubContract.contract.UnpackLog(event, "OneIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugSubContractThreeIndexEventIterator is returned from FilterThreeIndexEvent and is used to iterate over the raw logs and unpacked data for ThreeIndexEvent events raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractThreeIndexEventIterator struct {
	Event *NetworkDebugSubContractThreeIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugSubContractThreeIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugSubContractThreeIndexEvent)
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
		it.Event = new(NetworkDebugSubContractThreeIndexEvent)
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
func (it *NetworkDebugSubContractThreeIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugSubContractThreeIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugSubContractThreeIndexEvent represents a ThreeIndexEvent event raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractThreeIndexEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterThreeIndexEvent is a free log retrieval operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) FilterThreeIndexEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*NetworkDebugSubContractThreeIndexEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugSubContract.contract.FilterLogs(opts, "ThreeIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractThreeIndexEventIterator{contract: _NetworkDebugSubContract.contract, event: "ThreeIndexEvent", logs: logs, sub: sub}, nil
}

// WatchThreeIndexEvent is a free log subscription operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) WatchThreeIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugSubContractThreeIndexEvent, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugSubContract.contract.WatchLogs(opts, "ThreeIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugSubContractThreeIndexEvent)
				if err := _NetworkDebugSubContract.contract.UnpackLog(event, "ThreeIndexEvent", log); err != nil {
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

// ParseThreeIndexEvent is a log parse operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) ParseThreeIndexEvent(log types.Log) (*NetworkDebugSubContractThreeIndexEvent, error) {
	event := new(NetworkDebugSubContractThreeIndexEvent)
	if err := _NetworkDebugSubContract.contract.UnpackLog(event, "ThreeIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugSubContractTwoIndexEventIterator is returned from FilterTwoIndexEvent and is used to iterate over the raw logs and unpacked data for TwoIndexEvent events raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractTwoIndexEventIterator struct {
	Event *NetworkDebugSubContractTwoIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugSubContractTwoIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugSubContractTwoIndexEvent)
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
		it.Event = new(NetworkDebugSubContractTwoIndexEvent)
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
func (it *NetworkDebugSubContractTwoIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugSubContractTwoIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugSubContractTwoIndexEvent represents a TwoIndexEvent event raised by the NetworkDebugSubContract contract.
type NetworkDebugSubContractTwoIndexEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTwoIndexEvent is a free log retrieval operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) FilterTwoIndexEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*NetworkDebugSubContractTwoIndexEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugSubContract.contract.FilterLogs(opts, "TwoIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugSubContractTwoIndexEventIterator{contract: _NetworkDebugSubContract.contract, event: "TwoIndexEvent", logs: logs, sub: sub}, nil
}

// WatchTwoIndexEvent is a free log subscription operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) WatchTwoIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugSubContractTwoIndexEvent, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugSubContract.contract.WatchLogs(opts, "TwoIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugSubContractTwoIndexEvent)
				if err := _NetworkDebugSubContract.contract.UnpackLog(event, "TwoIndexEvent", log); err != nil {
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

// ParseTwoIndexEvent is a log parse operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugSubContract *NetworkDebugSubContractFilterer) ParseTwoIndexEvent(log types.Log) (*NetworkDebugSubContractTwoIndexEvent, error) {
	event := new(NetworkDebugSubContractTwoIndexEvent)
	if err := _NetworkDebugSubContract.contract.UnpackLog(event, "TwoIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
