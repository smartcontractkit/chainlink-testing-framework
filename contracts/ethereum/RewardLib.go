// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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
)

// RewardLibMetaData contains all meta data concerning the RewardLib contract.
var RewardLibMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"RewardDurationTooShort\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountAdded\",\"type\":\"uint256\"}],\"name\":\"RewardAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"name\":\"RewardInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"RewardRateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"operator\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"slashedBaseRewards\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"slashedDelegatedRewards\",\"type\":\"uint256[]\"}],\"name\":\"RewardSlashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardWithdrawn\",\"type\":\"event\"}]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212203a1d042dfd6b700be27ad9cd04c9df01cbd70225b35d383cb78cfd44b717a46164736f6c63430008100033",
}

// RewardLibABI is the input ABI used to generate the binding from.
// Deprecated: Use RewardLibMetaData.ABI instead.
var RewardLibABI = RewardLibMetaData.ABI

// RewardLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RewardLibMetaData.Bin instead.
var RewardLibBin = RewardLibMetaData.Bin

// DeployRewardLib deploys a new Ethereum contract, binding an instance of RewardLib to it.
func DeployRewardLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RewardLib, error) {
	parsed, err := RewardLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RewardLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RewardLib{RewardLibCaller: RewardLibCaller{contract: contract}, RewardLibTransactor: RewardLibTransactor{contract: contract}, RewardLibFilterer: RewardLibFilterer{contract: contract}}, nil
}

// RewardLib is an auto generated Go binding around an Ethereum contract.
type RewardLib struct {
	RewardLibCaller     // Read-only binding to the contract
	RewardLibTransactor // Write-only binding to the contract
	RewardLibFilterer   // Log filterer for contract events
}

// RewardLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardLibSession struct {
	Contract     *RewardLib        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardLibCallerSession struct {
	Contract *RewardLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RewardLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardLibTransactorSession struct {
	Contract     *RewardLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RewardLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardLibRaw struct {
	Contract *RewardLib // Generic contract binding to access the raw methods on
}

// RewardLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardLibCallerRaw struct {
	Contract *RewardLibCaller // Generic read-only contract binding to access the raw methods on
}

// RewardLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardLibTransactorRaw struct {
	Contract *RewardLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewardLib creates a new instance of RewardLib, bound to a specific deployed contract.
func NewRewardLib(address common.Address, backend bind.ContractBackend) (*RewardLib, error) {
	contract, err := bindRewardLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RewardLib{RewardLibCaller: RewardLibCaller{contract: contract}, RewardLibTransactor: RewardLibTransactor{contract: contract}, RewardLibFilterer: RewardLibFilterer{contract: contract}}, nil
}

// NewRewardLibCaller creates a new read-only instance of RewardLib, bound to a specific deployed contract.
func NewRewardLibCaller(address common.Address, caller bind.ContractCaller) (*RewardLibCaller, error) {
	contract, err := bindRewardLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardLibCaller{contract: contract}, nil
}

// NewRewardLibTransactor creates a new write-only instance of RewardLib, bound to a specific deployed contract.
func NewRewardLibTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardLibTransactor, error) {
	contract, err := bindRewardLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardLibTransactor{contract: contract}, nil
}

// NewRewardLibFilterer creates a new log filterer instance of RewardLib, bound to a specific deployed contract.
func NewRewardLibFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardLibFilterer, error) {
	contract, err := bindRewardLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardLibFilterer{contract: contract}, nil
}

// bindRewardLib binds a generic wrapper to an already deployed contract.
func bindRewardLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RewardLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardLib *RewardLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardLib.Contract.RewardLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardLib *RewardLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardLib.Contract.RewardLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardLib *RewardLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardLib.Contract.RewardLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardLib *RewardLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardLib *RewardLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardLib *RewardLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardLib.Contract.contract.Transact(opts, method, params...)
}

// RewardLibRewardAddedIterator is returned from FilterRewardAdded and is used to iterate over the raw logs and unpacked data for RewardAdded events raised by the RewardLib contract.
type RewardLibRewardAddedIterator struct {
	Event *RewardLibRewardAdded // Event containing the contract specifics and raw log

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
func (it *RewardLibRewardAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardLibRewardAdded)
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
		it.Event = new(RewardLibRewardAdded)
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
func (it *RewardLibRewardAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardLibRewardAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardLibRewardAdded represents a RewardAdded event raised by the RewardLib contract.
type RewardLibRewardAdded struct {
	AmountAdded *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRewardAdded is a free log retrieval operation binding the contract event 0xde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d.
//
// Solidity: event RewardAdded(uint256 amountAdded)
func (_RewardLib *RewardLibFilterer) FilterRewardAdded(opts *bind.FilterOpts) (*RewardLibRewardAddedIterator, error) {

	logs, sub, err := _RewardLib.contract.FilterLogs(opts, "RewardAdded")
	if err != nil {
		return nil, err
	}
	return &RewardLibRewardAddedIterator{contract: _RewardLib.contract, event: "RewardAdded", logs: logs, sub: sub}, nil
}

// WatchRewardAdded is a free log subscription operation binding the contract event 0xde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d.
//
// Solidity: event RewardAdded(uint256 amountAdded)
func (_RewardLib *RewardLibFilterer) WatchRewardAdded(opts *bind.WatchOpts, sink chan<- *RewardLibRewardAdded) (event.Subscription, error) {

	logs, sub, err := _RewardLib.contract.WatchLogs(opts, "RewardAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardLibRewardAdded)
				if err := _RewardLib.contract.UnpackLog(event, "RewardAdded", log); err != nil {
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

// ParseRewardAdded is a log parse operation binding the contract event 0xde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d.
//
// Solidity: event RewardAdded(uint256 amountAdded)
func (_RewardLib *RewardLibFilterer) ParseRewardAdded(log types.Log) (*RewardLibRewardAdded, error) {
	event := new(RewardLibRewardAdded)
	if err := _RewardLib.contract.UnpackLog(event, "RewardAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardLibRewardInitializedIterator is returned from FilterRewardInitialized and is used to iterate over the raw logs and unpacked data for RewardInitialized events raised by the RewardLib contract.
type RewardLibRewardInitializedIterator struct {
	Event *RewardLibRewardInitialized // Event containing the contract specifics and raw log

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
func (it *RewardLibRewardInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardLibRewardInitialized)
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
		it.Event = new(RewardLibRewardInitialized)
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
func (it *RewardLibRewardInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardLibRewardInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardLibRewardInitialized represents a RewardInitialized event raised by the RewardLib contract.
type RewardLibRewardInitialized struct {
	Rate           *big.Int
	Available      *big.Int
	StartTimestamp *big.Int
	EndTimestamp   *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRewardInitialized is a free log retrieval operation binding the contract event 0x125fc8494f786b470e3c39d0932a62e9e09e291ebd81ea19c57604f6d2b1d167.
//
// Solidity: event RewardInitialized(uint256 rate, uint256 available, uint256 startTimestamp, uint256 endTimestamp)
func (_RewardLib *RewardLibFilterer) FilterRewardInitialized(opts *bind.FilterOpts) (*RewardLibRewardInitializedIterator, error) {

	logs, sub, err := _RewardLib.contract.FilterLogs(opts, "RewardInitialized")
	if err != nil {
		return nil, err
	}
	return &RewardLibRewardInitializedIterator{contract: _RewardLib.contract, event: "RewardInitialized", logs: logs, sub: sub}, nil
}

// WatchRewardInitialized is a free log subscription operation binding the contract event 0x125fc8494f786b470e3c39d0932a62e9e09e291ebd81ea19c57604f6d2b1d167.
//
// Solidity: event RewardInitialized(uint256 rate, uint256 available, uint256 startTimestamp, uint256 endTimestamp)
func (_RewardLib *RewardLibFilterer) WatchRewardInitialized(opts *bind.WatchOpts, sink chan<- *RewardLibRewardInitialized) (event.Subscription, error) {

	logs, sub, err := _RewardLib.contract.WatchLogs(opts, "RewardInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardLibRewardInitialized)
				if err := _RewardLib.contract.UnpackLog(event, "RewardInitialized", log); err != nil {
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

// ParseRewardInitialized is a log parse operation binding the contract event 0x125fc8494f786b470e3c39d0932a62e9e09e291ebd81ea19c57604f6d2b1d167.
//
// Solidity: event RewardInitialized(uint256 rate, uint256 available, uint256 startTimestamp, uint256 endTimestamp)
func (_RewardLib *RewardLibFilterer) ParseRewardInitialized(log types.Log) (*RewardLibRewardInitialized, error) {
	event := new(RewardLibRewardInitialized)
	if err := _RewardLib.contract.UnpackLog(event, "RewardInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardLibRewardRateChangedIterator is returned from FilterRewardRateChanged and is used to iterate over the raw logs and unpacked data for RewardRateChanged events raised by the RewardLib contract.
type RewardLibRewardRateChangedIterator struct {
	Event *RewardLibRewardRateChanged // Event containing the contract specifics and raw log

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
func (it *RewardLibRewardRateChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardLibRewardRateChanged)
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
		it.Event = new(RewardLibRewardRateChanged)
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
func (it *RewardLibRewardRateChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardLibRewardRateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardLibRewardRateChanged represents a RewardRateChanged event raised by the RewardLib contract.
type RewardLibRewardRateChanged struct {
	Rate *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRewardRateChanged is a free log retrieval operation binding the contract event 0x1e3be2efa25bca5bff2215c7b30b31086e703d6aa7d9b9a1f8ba62c5291219ad.
//
// Solidity: event RewardRateChanged(uint256 rate)
func (_RewardLib *RewardLibFilterer) FilterRewardRateChanged(opts *bind.FilterOpts) (*RewardLibRewardRateChangedIterator, error) {

	logs, sub, err := _RewardLib.contract.FilterLogs(opts, "RewardRateChanged")
	if err != nil {
		return nil, err
	}
	return &RewardLibRewardRateChangedIterator{contract: _RewardLib.contract, event: "RewardRateChanged", logs: logs, sub: sub}, nil
}

// WatchRewardRateChanged is a free log subscription operation binding the contract event 0x1e3be2efa25bca5bff2215c7b30b31086e703d6aa7d9b9a1f8ba62c5291219ad.
//
// Solidity: event RewardRateChanged(uint256 rate)
func (_RewardLib *RewardLibFilterer) WatchRewardRateChanged(opts *bind.WatchOpts, sink chan<- *RewardLibRewardRateChanged) (event.Subscription, error) {

	logs, sub, err := _RewardLib.contract.WatchLogs(opts, "RewardRateChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardLibRewardRateChanged)
				if err := _RewardLib.contract.UnpackLog(event, "RewardRateChanged", log); err != nil {
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

// ParseRewardRateChanged is a log parse operation binding the contract event 0x1e3be2efa25bca5bff2215c7b30b31086e703d6aa7d9b9a1f8ba62c5291219ad.
//
// Solidity: event RewardRateChanged(uint256 rate)
func (_RewardLib *RewardLibFilterer) ParseRewardRateChanged(log types.Log) (*RewardLibRewardRateChanged, error) {
	event := new(RewardLibRewardRateChanged)
	if err := _RewardLib.contract.UnpackLog(event, "RewardRateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardLibRewardSlashedIterator is returned from FilterRewardSlashed and is used to iterate over the raw logs and unpacked data for RewardSlashed events raised by the RewardLib contract.
type RewardLibRewardSlashedIterator struct {
	Event *RewardLibRewardSlashed // Event containing the contract specifics and raw log

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
func (it *RewardLibRewardSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardLibRewardSlashed)
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
		it.Event = new(RewardLibRewardSlashed)
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
func (it *RewardLibRewardSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardLibRewardSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardLibRewardSlashed represents a RewardSlashed event raised by the RewardLib contract.
type RewardLibRewardSlashed struct {
	Operator                []common.Address
	SlashedBaseRewards      []*big.Int
	SlashedDelegatedRewards []*big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterRewardSlashed is a free log retrieval operation binding the contract event 0x00635ea9da6e262e92bb713d71840af7c567807ff35bf73e927490c612832480.
//
// Solidity: event RewardSlashed(address[] operator, uint256[] slashedBaseRewards, uint256[] slashedDelegatedRewards)
func (_RewardLib *RewardLibFilterer) FilterRewardSlashed(opts *bind.FilterOpts) (*RewardLibRewardSlashedIterator, error) {

	logs, sub, err := _RewardLib.contract.FilterLogs(opts, "RewardSlashed")
	if err != nil {
		return nil, err
	}
	return &RewardLibRewardSlashedIterator{contract: _RewardLib.contract, event: "RewardSlashed", logs: logs, sub: sub}, nil
}

// WatchRewardSlashed is a free log subscription operation binding the contract event 0x00635ea9da6e262e92bb713d71840af7c567807ff35bf73e927490c612832480.
//
// Solidity: event RewardSlashed(address[] operator, uint256[] slashedBaseRewards, uint256[] slashedDelegatedRewards)
func (_RewardLib *RewardLibFilterer) WatchRewardSlashed(opts *bind.WatchOpts, sink chan<- *RewardLibRewardSlashed) (event.Subscription, error) {

	logs, sub, err := _RewardLib.contract.WatchLogs(opts, "RewardSlashed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardLibRewardSlashed)
				if err := _RewardLib.contract.UnpackLog(event, "RewardSlashed", log); err != nil {
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

// ParseRewardSlashed is a log parse operation binding the contract event 0x00635ea9da6e262e92bb713d71840af7c567807ff35bf73e927490c612832480.
//
// Solidity: event RewardSlashed(address[] operator, uint256[] slashedBaseRewards, uint256[] slashedDelegatedRewards)
func (_RewardLib *RewardLibFilterer) ParseRewardSlashed(log types.Log) (*RewardLibRewardSlashed, error) {
	event := new(RewardLibRewardSlashed)
	if err := _RewardLib.contract.UnpackLog(event, "RewardSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardLibRewardWithdrawnIterator is returned from FilterRewardWithdrawn and is used to iterate over the raw logs and unpacked data for RewardWithdrawn events raised by the RewardLib contract.
type RewardLibRewardWithdrawnIterator struct {
	Event *RewardLibRewardWithdrawn // Event containing the contract specifics and raw log

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
func (it *RewardLibRewardWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardLibRewardWithdrawn)
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
		it.Event = new(RewardLibRewardWithdrawn)
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
func (it *RewardLibRewardWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardLibRewardWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardLibRewardWithdrawn represents a RewardWithdrawn event raised by the RewardLib contract.
type RewardLibRewardWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardWithdrawn is a free log retrieval operation binding the contract event 0x150a6ec0e6f4e9ddcaaaa1674f157d91165a42d60653016f87a9fc870a39f050.
//
// Solidity: event RewardWithdrawn(uint256 amount)
func (_RewardLib *RewardLibFilterer) FilterRewardWithdrawn(opts *bind.FilterOpts) (*RewardLibRewardWithdrawnIterator, error) {

	logs, sub, err := _RewardLib.contract.FilterLogs(opts, "RewardWithdrawn")
	if err != nil {
		return nil, err
	}
	return &RewardLibRewardWithdrawnIterator{contract: _RewardLib.contract, event: "RewardWithdrawn", logs: logs, sub: sub}, nil
}

// WatchRewardWithdrawn is a free log subscription operation binding the contract event 0x150a6ec0e6f4e9ddcaaaa1674f157d91165a42d60653016f87a9fc870a39f050.
//
// Solidity: event RewardWithdrawn(uint256 amount)
func (_RewardLib *RewardLibFilterer) WatchRewardWithdrawn(opts *bind.WatchOpts, sink chan<- *RewardLibRewardWithdrawn) (event.Subscription, error) {

	logs, sub, err := _RewardLib.contract.WatchLogs(opts, "RewardWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardLibRewardWithdrawn)
				if err := _RewardLib.contract.UnpackLog(event, "RewardWithdrawn", log); err != nil {
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

// ParseRewardWithdrawn is a log parse operation binding the contract event 0x150a6ec0e6f4e9ddcaaaa1674f157d91165a42d60653016f87a9fc870a39f050.
//
// Solidity: event RewardWithdrawn(uint256 amount)
func (_RewardLib *RewardLibFilterer) ParseRewardWithdrawn(log types.Log) (*RewardLibRewardWithdrawn, error) {
	event := new(RewardLibRewardWithdrawn)
	if err := _RewardLib.contract.UnpackLog(event, "RewardWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
