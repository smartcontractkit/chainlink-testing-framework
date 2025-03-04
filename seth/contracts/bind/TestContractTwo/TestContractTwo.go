// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TestContractTwo

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

// UniqueEventTwoMetaData contains all meta data concerning the UniqueEventTwo contract.
var UniqueEventTwoMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"a\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"b\",\"type\":\"int256\"}],\"name\":\"NonUniqueEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"executeSecondOperation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101f2806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e31279c514610030575b600080fd5b61004a600480360381019061004591906100df565b610060565b604051610057919061012e565b60405180910390f35b600081837f192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e60405160405180910390a3818361009c9190610178565b905092915050565b600080fd5b6000819050919050565b6100bc816100a9565b81146100c757600080fd5b50565b6000813590506100d9816100b3565b92915050565b600080604083850312156100f6576100f56100a4565b5b6000610104858286016100ca565b9250506020610115858286016100ca565b9150509250929050565b610128816100a9565b82525050565b6000602082019050610143600083018461011f565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610183826100a9565b915061018e836100a9565b9250828201905082811215600083121683821260008412151617156101b6576101b5610149565b5b9291505056fea264697066735822122036206ffb0222909d65d2a872080e677681bc6601dce2b4cfebd9503a1187280964736f6c63430008130033",
}

// UniqueEventTwoABI is the input ABI used to generate the binding from.
// Deprecated: Use UniqueEventTwoMetaData.ABI instead.
var UniqueEventTwoABI = UniqueEventTwoMetaData.ABI

// UniqueEventTwoBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniqueEventTwoMetaData.Bin instead.
var UniqueEventTwoBin = UniqueEventTwoMetaData.Bin

// DeployUniqueEventTwo deploys a new Ethereum contract, binding an instance of UniqueEventTwo to it.
func DeployUniqueEventTwo(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniqueEventTwo, error) {
	parsed, err := UniqueEventTwoMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniqueEventTwoBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniqueEventTwo{UniqueEventTwoCaller: UniqueEventTwoCaller{contract: contract}, UniqueEventTwoTransactor: UniqueEventTwoTransactor{contract: contract}, UniqueEventTwoFilterer: UniqueEventTwoFilterer{contract: contract}}, nil
}

// UniqueEventTwo is an auto generated Go binding around an Ethereum contract.
type UniqueEventTwo struct {
	UniqueEventTwoCaller     // Read-only binding to the contract
	UniqueEventTwoTransactor // Write-only binding to the contract
	UniqueEventTwoFilterer   // Log filterer for contract events
}

// UniqueEventTwoCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniqueEventTwoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniqueEventTwoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniqueEventTwoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniqueEventTwoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniqueEventTwoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniqueEventTwoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniqueEventTwoSession struct {
	Contract     *UniqueEventTwo   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniqueEventTwoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniqueEventTwoCallerSession struct {
	Contract *UniqueEventTwoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// UniqueEventTwoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniqueEventTwoTransactorSession struct {
	Contract     *UniqueEventTwoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// UniqueEventTwoRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniqueEventTwoRaw struct {
	Contract *UniqueEventTwo // Generic contract binding to access the raw methods on
}

// UniqueEventTwoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniqueEventTwoCallerRaw struct {
	Contract *UniqueEventTwoCaller // Generic read-only contract binding to access the raw methods on
}

// UniqueEventTwoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniqueEventTwoTransactorRaw struct {
	Contract *UniqueEventTwoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniqueEventTwo creates a new instance of UniqueEventTwo, bound to a specific deployed contract.
func NewUniqueEventTwo(address common.Address, backend bind.ContractBackend) (*UniqueEventTwo, error) {
	contract, err := bindUniqueEventTwo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniqueEventTwo{UniqueEventTwoCaller: UniqueEventTwoCaller{contract: contract}, UniqueEventTwoTransactor: UniqueEventTwoTransactor{contract: contract}, UniqueEventTwoFilterer: UniqueEventTwoFilterer{contract: contract}}, nil
}

// NewUniqueEventTwoCaller creates a new read-only instance of UniqueEventTwo, bound to a specific deployed contract.
func NewUniqueEventTwoCaller(address common.Address, caller bind.ContractCaller) (*UniqueEventTwoCaller, error) {
	contract, err := bindUniqueEventTwo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniqueEventTwoCaller{contract: contract}, nil
}

// NewUniqueEventTwoTransactor creates a new write-only instance of UniqueEventTwo, bound to a specific deployed contract.
func NewUniqueEventTwoTransactor(address common.Address, transactor bind.ContractTransactor) (*UniqueEventTwoTransactor, error) {
	contract, err := bindUniqueEventTwo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniqueEventTwoTransactor{contract: contract}, nil
}

// NewUniqueEventTwoFilterer creates a new log filterer instance of UniqueEventTwo, bound to a specific deployed contract.
func NewUniqueEventTwoFilterer(address common.Address, filterer bind.ContractFilterer) (*UniqueEventTwoFilterer, error) {
	contract, err := bindUniqueEventTwo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniqueEventTwoFilterer{contract: contract}, nil
}

// bindUniqueEventTwo binds a generic wrapper to an already deployed contract.
func bindUniqueEventTwo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniqueEventTwoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniqueEventTwo *UniqueEventTwoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniqueEventTwo.Contract.UniqueEventTwoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniqueEventTwo *UniqueEventTwoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniqueEventTwo.Contract.UniqueEventTwoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniqueEventTwo *UniqueEventTwoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniqueEventTwo.Contract.UniqueEventTwoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniqueEventTwo *UniqueEventTwoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniqueEventTwo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniqueEventTwo *UniqueEventTwoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniqueEventTwo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniqueEventTwo *UniqueEventTwoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniqueEventTwo.Contract.contract.Transact(opts, method, params...)
}

// ExecuteSecondOperation is a paid mutator transaction binding the contract method 0xe31279c5.
//
// Solidity: function executeSecondOperation(int256 x, int256 y) returns(int256)
func (_UniqueEventTwo *UniqueEventTwoTransactor) ExecuteSecondOperation(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _UniqueEventTwo.contract.Transact(opts, "executeSecondOperation", x, y)
}

// ExecuteSecondOperation is a paid mutator transaction binding the contract method 0xe31279c5.
//
// Solidity: function executeSecondOperation(int256 x, int256 y) returns(int256)
func (_UniqueEventTwo *UniqueEventTwoSession) ExecuteSecondOperation(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _UniqueEventTwo.Contract.ExecuteSecondOperation(&_UniqueEventTwo.TransactOpts, x, y)
}

// ExecuteSecondOperation is a paid mutator transaction binding the contract method 0xe31279c5.
//
// Solidity: function executeSecondOperation(int256 x, int256 y) returns(int256)
func (_UniqueEventTwo *UniqueEventTwoTransactorSession) ExecuteSecondOperation(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _UniqueEventTwo.Contract.ExecuteSecondOperation(&_UniqueEventTwo.TransactOpts, x, y)
}

// UniqueEventTwoNonUniqueEventIterator is returned from FilterNonUniqueEvent and is used to iterate over the raw logs and unpacked data for NonUniqueEvent events raised by the UniqueEventTwo contract.
type UniqueEventTwoNonUniqueEventIterator struct {
	Event *UniqueEventTwoNonUniqueEvent // Event containing the contract specifics and raw log

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
func (it *UniqueEventTwoNonUniqueEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniqueEventTwoNonUniqueEvent)
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
		it.Event = new(UniqueEventTwoNonUniqueEvent)
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
func (it *UniqueEventTwoNonUniqueEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniqueEventTwoNonUniqueEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniqueEventTwoNonUniqueEvent represents a NonUniqueEvent event raised by the UniqueEventTwo contract.
type UniqueEventTwoNonUniqueEvent struct {
	A   *big.Int
	B   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNonUniqueEvent is a free log retrieval operation binding the contract event 0x192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e.
//
// Solidity: event NonUniqueEvent(int256 indexed a, int256 indexed b)
func (_UniqueEventTwo *UniqueEventTwoFilterer) FilterNonUniqueEvent(opts *bind.FilterOpts, a []*big.Int, b []*big.Int) (*UniqueEventTwoNonUniqueEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _UniqueEventTwo.contract.FilterLogs(opts, "NonUniqueEvent", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return &UniqueEventTwoNonUniqueEventIterator{contract: _UniqueEventTwo.contract, event: "NonUniqueEvent", logs: logs, sub: sub}, nil
}

// WatchNonUniqueEvent is a free log subscription operation binding the contract event 0x192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e.
//
// Solidity: event NonUniqueEvent(int256 indexed a, int256 indexed b)
func (_UniqueEventTwo *UniqueEventTwoFilterer) WatchNonUniqueEvent(opts *bind.WatchOpts, sink chan<- *UniqueEventTwoNonUniqueEvent, a []*big.Int, b []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _UniqueEventTwo.contract.WatchLogs(opts, "NonUniqueEvent", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniqueEventTwoNonUniqueEvent)
				if err := _UniqueEventTwo.contract.UnpackLog(event, "NonUniqueEvent", log); err != nil {
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

// ParseNonUniqueEvent is a log parse operation binding the contract event 0x192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e.
//
// Solidity: event NonUniqueEvent(int256 indexed a, int256 indexed b)
func (_UniqueEventTwo *UniqueEventTwoFilterer) ParseNonUniqueEvent(log types.Log) (*UniqueEventTwoNonUniqueEvent, error) {
	event := new(UniqueEventTwoNonUniqueEvent)
	if err := _UniqueEventTwo.contract.UnpackLog(event, "NonUniqueEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
