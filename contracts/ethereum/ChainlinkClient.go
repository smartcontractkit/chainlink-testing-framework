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

// ChainlinkClientABI is the input ABI used to generate the binding from.
const ChainlinkClientABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"}]"

// ChainlinkClientBin is the compiled bytecode used for deploying new contracts.
var ChainlinkClientBin = "0x60806040526001600455348015601457600080fd5b50603f8060226000396000f3fe6080604052600080fdfea26469706673582212204ea398bb6cc8cf31c687d7a3dc34f2a27c3dbab682f789777cc5f5e98af9ff7f64736f6c63430006060033"

// DeployChainlinkClient deploys a new Ethereum contract, binding an instance of ChainlinkClient to it.
func DeployChainlinkClient(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainlinkClient, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainlinkClientABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChainlinkClientBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainlinkClient{ChainlinkClientCaller: ChainlinkClientCaller{contract: contract}, ChainlinkClientTransactor: ChainlinkClientTransactor{contract: contract}, ChainlinkClientFilterer: ChainlinkClientFilterer{contract: contract}}, nil
}

// ChainlinkClient is an auto generated Go binding around an Ethereum contract.
type ChainlinkClient struct {
	ChainlinkClientCaller     // Read-only binding to the contract
	ChainlinkClientTransactor // Write-only binding to the contract
	ChainlinkClientFilterer   // Log filterer for contract events
}

// ChainlinkClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainlinkClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainlinkClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainlinkClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainlinkClientSession struct {
	Contract     *ChainlinkClient  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainlinkClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainlinkClientCallerSession struct {
	Contract *ChainlinkClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ChainlinkClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainlinkClientTransactorSession struct {
	Contract     *ChainlinkClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainlinkClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainlinkClientRaw struct {
	Contract *ChainlinkClient // Generic contract binding to access the raw methods on
}

// ChainlinkClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainlinkClientCallerRaw struct {
	Contract *ChainlinkClientCaller // Generic read-only contract binding to access the raw methods on
}

// ChainlinkClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainlinkClientTransactorRaw struct {
	Contract *ChainlinkClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainlinkClient creates a new instance of ChainlinkClient, bound to a specific deployed contract.
func NewChainlinkClient(address common.Address, backend bind.ContractBackend) (*ChainlinkClient, error) {
	contract, err := bindChainlinkClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClient{ChainlinkClientCaller: ChainlinkClientCaller{contract: contract}, ChainlinkClientTransactor: ChainlinkClientTransactor{contract: contract}, ChainlinkClientFilterer: ChainlinkClientFilterer{contract: contract}}, nil
}

// NewChainlinkClientCaller creates a new read-only instance of ChainlinkClient, bound to a specific deployed contract.
func NewChainlinkClientCaller(address common.Address, caller bind.ContractCaller) (*ChainlinkClientCaller, error) {
	contract, err := bindChainlinkClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClientCaller{contract: contract}, nil
}

// NewChainlinkClientTransactor creates a new write-only instance of ChainlinkClient, bound to a specific deployed contract.
func NewChainlinkClientTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainlinkClientTransactor, error) {
	contract, err := bindChainlinkClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClientTransactor{contract: contract}, nil
}

// NewChainlinkClientFilterer creates a new log filterer instance of ChainlinkClient, bound to a specific deployed contract.
func NewChainlinkClientFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainlinkClientFilterer, error) {
	contract, err := bindChainlinkClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClientFilterer{contract: contract}, nil
}

// bindChainlinkClient binds a generic wrapper to an already deployed contract.
func bindChainlinkClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainlinkClientABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainlinkClient *ChainlinkClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainlinkClient.Contract.ChainlinkClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainlinkClient *ChainlinkClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainlinkClient.Contract.ChainlinkClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainlinkClient *ChainlinkClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainlinkClient.Contract.ChainlinkClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainlinkClient *ChainlinkClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainlinkClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainlinkClient *ChainlinkClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainlinkClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainlinkClient *ChainlinkClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainlinkClient.Contract.contract.Transact(opts, method, params...)
}

// ChainlinkClientChainlinkCancelledIterator is returned from FilterChainlinkCancelled and is used to iterate over the raw logs and unpacked data for ChainlinkCancelled events raised by the ChainlinkClient contract.
type ChainlinkClientChainlinkCancelledIterator struct {
	Event *ChainlinkClientChainlinkCancelled // Event containing the contract specifics and raw log

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
func (it *ChainlinkClientChainlinkCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainlinkClientChainlinkCancelled)
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
		it.Event = new(ChainlinkClientChainlinkCancelled)
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
func (it *ChainlinkClientChainlinkCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainlinkClientChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainlinkClientChainlinkCancelled represents a ChainlinkCancelled event raised by the ChainlinkClient contract.
type ChainlinkClientChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkCancelled is a free log retrieval operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*ChainlinkClientChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ChainlinkClient.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClientChainlinkCancelledIterator{contract: _ChainlinkClient.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

// WatchChainlinkCancelled is a free log subscription operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *ChainlinkClientChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ChainlinkClient.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainlinkClientChainlinkCancelled)
				if err := _ChainlinkClient.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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

// ParseChainlinkCancelled is a log parse operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) ParseChainlinkCancelled(log types.Log) (*ChainlinkClientChainlinkCancelled, error) {
	event := new(ChainlinkClientChainlinkCancelled)
	if err := _ChainlinkClient.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChainlinkClientChainlinkFulfilledIterator is returned from FilterChainlinkFulfilled and is used to iterate over the raw logs and unpacked data for ChainlinkFulfilled events raised by the ChainlinkClient contract.
type ChainlinkClientChainlinkFulfilledIterator struct {
	Event *ChainlinkClientChainlinkFulfilled // Event containing the contract specifics and raw log

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
func (it *ChainlinkClientChainlinkFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainlinkClientChainlinkFulfilled)
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
		it.Event = new(ChainlinkClientChainlinkFulfilled)
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
func (it *ChainlinkClientChainlinkFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainlinkClientChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainlinkClientChainlinkFulfilled represents a ChainlinkFulfilled event raised by the ChainlinkClient contract.
type ChainlinkClientChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkFulfilled is a free log retrieval operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*ChainlinkClientChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ChainlinkClient.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClientChainlinkFulfilledIterator{contract: _ChainlinkClient.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

// WatchChainlinkFulfilled is a free log subscription operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *ChainlinkClientChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ChainlinkClient.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainlinkClientChainlinkFulfilled)
				if err := _ChainlinkClient.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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

// ParseChainlinkFulfilled is a log parse operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) ParseChainlinkFulfilled(log types.Log) (*ChainlinkClientChainlinkFulfilled, error) {
	event := new(ChainlinkClientChainlinkFulfilled)
	if err := _ChainlinkClient.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChainlinkClientChainlinkRequestedIterator is returned from FilterChainlinkRequested and is used to iterate over the raw logs and unpacked data for ChainlinkRequested events raised by the ChainlinkClient contract.
type ChainlinkClientChainlinkRequestedIterator struct {
	Event *ChainlinkClientChainlinkRequested // Event containing the contract specifics and raw log

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
func (it *ChainlinkClientChainlinkRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainlinkClientChainlinkRequested)
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
		it.Event = new(ChainlinkClientChainlinkRequested)
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
func (it *ChainlinkClientChainlinkRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainlinkClientChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainlinkClientChainlinkRequested represents a ChainlinkRequested event raised by the ChainlinkClient contract.
type ChainlinkClientChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkRequested is a free log retrieval operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*ChainlinkClientChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ChainlinkClient.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClientChainlinkRequestedIterator{contract: _ChainlinkClient.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

// WatchChainlinkRequested is a free log subscription operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *ChainlinkClientChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ChainlinkClient.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainlinkClientChainlinkRequested)
				if err := _ChainlinkClient.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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

// ParseChainlinkRequested is a log parse operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_ChainlinkClient *ChainlinkClientFilterer) ParseChainlinkRequested(log types.Log) (*ChainlinkClientChainlinkRequested, error) {
	event := new(ChainlinkClientChainlinkRequested)
	if err := _ChainlinkClient.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
