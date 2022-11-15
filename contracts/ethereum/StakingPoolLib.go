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

// StakingPoolLibMetaData contains all meta data concerning the StakingPoolLib contract.
var StakingPoolLibMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"remainingAmount\",\"type\":\"uint256\"}],\"name\":\"ExcessiveStakeAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"ExistingStakeFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentOperatorsCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minInitialOperatorsCount\",\"type\":\"uint256\"}],\"name\":\"InadequateInitialOperatorsCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"remainingPoolSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredPoolSize\",\"type\":\"uint256\"}],\"name\":\"InsufficientRemainingPoolSpace\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requiredAmount\",\"type\":\"uint256\"}],\"name\":\"InsufficientStakeAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"InvalidMaxStakeAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxPoolSize\",\"type\":\"uint256\"}],\"name\":\"InvalidPoolSize\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"currentStatus\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"requiredStatus\",\"type\":\"bool\"}],\"name\":\"InvalidPoolStatus\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorIsAssignedToFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorIsLocked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"StakeNotFound\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"feedOperators\",\"type\":\"address[]\"}],\"name\":\"FeedOperatorsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"MaxCommunityStakeAmountIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"MaxOperatorStakeAmountIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"PoolConcluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"PoolOpened\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxPoolSize\",\"type\":\"uint256\"}],\"name\":\"PoolSizeIncreased\",\"type\":\"event\"}]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212209b798aabab54b035809e86f88a28bf9c95266d3cec8e464dd973991d9d3c2c3f64736f6c63430008100033",
}

// StakingPoolLibABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingPoolLibMetaData.ABI instead.
var StakingPoolLibABI = StakingPoolLibMetaData.ABI

// StakingPoolLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingPoolLibMetaData.Bin instead.
var StakingPoolLibBin = StakingPoolLibMetaData.Bin

// DeployStakingPoolLib deploys a new Ethereum contract, binding an instance of StakingPoolLib to it.
func DeployStakingPoolLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingPoolLib, error) {
	parsed, err := StakingPoolLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingPoolLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingPoolLib{StakingPoolLibCaller: StakingPoolLibCaller{contract: contract}, StakingPoolLibTransactor: StakingPoolLibTransactor{contract: contract}, StakingPoolLibFilterer: StakingPoolLibFilterer{contract: contract}}, nil
}

// StakingPoolLib is an auto generated Go binding around an Ethereum contract.
type StakingPoolLib struct {
	StakingPoolLibCaller     // Read-only binding to the contract
	StakingPoolLibTransactor // Write-only binding to the contract
	StakingPoolLibFilterer   // Log filterer for contract events
}

// StakingPoolLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingPoolLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingPoolLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingPoolLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingPoolLibSession struct {
	Contract     *StakingPoolLib   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingPoolLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingPoolLibCallerSession struct {
	Contract *StakingPoolLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StakingPoolLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingPoolLibTransactorSession struct {
	Contract     *StakingPoolLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StakingPoolLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingPoolLibRaw struct {
	Contract *StakingPoolLib // Generic contract binding to access the raw methods on
}

// StakingPoolLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingPoolLibCallerRaw struct {
	Contract *StakingPoolLibCaller // Generic read-only contract binding to access the raw methods on
}

// StakingPoolLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingPoolLibTransactorRaw struct {
	Contract *StakingPoolLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingPoolLib creates a new instance of StakingPoolLib, bound to a specific deployed contract.
func NewStakingPoolLib(address common.Address, backend bind.ContractBackend) (*StakingPoolLib, error) {
	contract, err := bindStakingPoolLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingPoolLib{StakingPoolLibCaller: StakingPoolLibCaller{contract: contract}, StakingPoolLibTransactor: StakingPoolLibTransactor{contract: contract}, StakingPoolLibFilterer: StakingPoolLibFilterer{contract: contract}}, nil
}

// NewStakingPoolLibCaller creates a new read-only instance of StakingPoolLib, bound to a specific deployed contract.
func NewStakingPoolLibCaller(address common.Address, caller bind.ContractCaller) (*StakingPoolLibCaller, error) {
	contract, err := bindStakingPoolLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibCaller{contract: contract}, nil
}

// NewStakingPoolLibTransactor creates a new write-only instance of StakingPoolLib, bound to a specific deployed contract.
func NewStakingPoolLibTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingPoolLibTransactor, error) {
	contract, err := bindStakingPoolLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibTransactor{contract: contract}, nil
}

// NewStakingPoolLibFilterer creates a new log filterer instance of StakingPoolLib, bound to a specific deployed contract.
func NewStakingPoolLibFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingPoolLibFilterer, error) {
	contract, err := bindStakingPoolLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibFilterer{contract: contract}, nil
}

// bindStakingPoolLib binds a generic wrapper to an already deployed contract.
func bindStakingPoolLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingPoolLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPoolLib *StakingPoolLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingPoolLib.Contract.StakingPoolLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPoolLib *StakingPoolLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPoolLib.Contract.StakingPoolLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPoolLib *StakingPoolLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPoolLib.Contract.StakingPoolLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPoolLib *StakingPoolLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingPoolLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPoolLib *StakingPoolLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPoolLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPoolLib *StakingPoolLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPoolLib.Contract.contract.Transact(opts, method, params...)
}

// StakingPoolLibFeedOperatorsSetIterator is returned from FilterFeedOperatorsSet and is used to iterate over the raw logs and unpacked data for FeedOperatorsSet events raised by the StakingPoolLib contract.
type StakingPoolLibFeedOperatorsSetIterator struct {
	Event *StakingPoolLibFeedOperatorsSet // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibFeedOperatorsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibFeedOperatorsSet)
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
		it.Event = new(StakingPoolLibFeedOperatorsSet)
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
func (it *StakingPoolLibFeedOperatorsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibFeedOperatorsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibFeedOperatorsSet represents a FeedOperatorsSet event raised by the StakingPoolLib contract.
type StakingPoolLibFeedOperatorsSet struct {
	FeedOperators []common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeedOperatorsSet is a free log retrieval operation binding the contract event 0x40aed8e423b39a56b445ae160f4c071fc2cfb48ee0b6dcd5ffeb6bc5b18d10d0.
//
// Solidity: event FeedOperatorsSet(address[] feedOperators)
func (_StakingPoolLib *StakingPoolLibFilterer) FilterFeedOperatorsSet(opts *bind.FilterOpts) (*StakingPoolLibFeedOperatorsSetIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "FeedOperatorsSet")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibFeedOperatorsSetIterator{contract: _StakingPoolLib.contract, event: "FeedOperatorsSet", logs: logs, sub: sub}, nil
}

// WatchFeedOperatorsSet is a free log subscription operation binding the contract event 0x40aed8e423b39a56b445ae160f4c071fc2cfb48ee0b6dcd5ffeb6bc5b18d10d0.
//
// Solidity: event FeedOperatorsSet(address[] feedOperators)
func (_StakingPoolLib *StakingPoolLibFilterer) WatchFeedOperatorsSet(opts *bind.WatchOpts, sink chan<- *StakingPoolLibFeedOperatorsSet) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "FeedOperatorsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibFeedOperatorsSet)
				if err := _StakingPoolLib.contract.UnpackLog(event, "FeedOperatorsSet", log); err != nil {
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

// ParseFeedOperatorsSet is a log parse operation binding the contract event 0x40aed8e423b39a56b445ae160f4c071fc2cfb48ee0b6dcd5ffeb6bc5b18d10d0.
//
// Solidity: event FeedOperatorsSet(address[] feedOperators)
func (_StakingPoolLib *StakingPoolLibFilterer) ParseFeedOperatorsSet(log types.Log) (*StakingPoolLibFeedOperatorsSet, error) {
	event := new(StakingPoolLibFeedOperatorsSet)
	if err := _StakingPoolLib.contract.UnpackLog(event, "FeedOperatorsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibMaxCommunityStakeAmountIncreasedIterator is returned from FilterMaxCommunityStakeAmountIncreased and is used to iterate over the raw logs and unpacked data for MaxCommunityStakeAmountIncreased events raised by the StakingPoolLib contract.
type StakingPoolLibMaxCommunityStakeAmountIncreasedIterator struct {
	Event *StakingPoolLibMaxCommunityStakeAmountIncreased // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibMaxCommunityStakeAmountIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibMaxCommunityStakeAmountIncreased)
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
		it.Event = new(StakingPoolLibMaxCommunityStakeAmountIncreased)
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
func (it *StakingPoolLibMaxCommunityStakeAmountIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibMaxCommunityStakeAmountIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibMaxCommunityStakeAmountIncreased represents a MaxCommunityStakeAmountIncreased event raised by the StakingPoolLib contract.
type StakingPoolLibMaxCommunityStakeAmountIncreased struct {
	MaxStakeAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMaxCommunityStakeAmountIncreased is a free log retrieval operation binding the contract event 0xb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df7.
//
// Solidity: event MaxCommunityStakeAmountIncreased(uint256 maxStakeAmount)
func (_StakingPoolLib *StakingPoolLibFilterer) FilterMaxCommunityStakeAmountIncreased(opts *bind.FilterOpts) (*StakingPoolLibMaxCommunityStakeAmountIncreasedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "MaxCommunityStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibMaxCommunityStakeAmountIncreasedIterator{contract: _StakingPoolLib.contract, event: "MaxCommunityStakeAmountIncreased", logs: logs, sub: sub}, nil
}

// WatchMaxCommunityStakeAmountIncreased is a free log subscription operation binding the contract event 0xb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df7.
//
// Solidity: event MaxCommunityStakeAmountIncreased(uint256 maxStakeAmount)
func (_StakingPoolLib *StakingPoolLibFilterer) WatchMaxCommunityStakeAmountIncreased(opts *bind.WatchOpts, sink chan<- *StakingPoolLibMaxCommunityStakeAmountIncreased) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "MaxCommunityStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibMaxCommunityStakeAmountIncreased)
				if err := _StakingPoolLib.contract.UnpackLog(event, "MaxCommunityStakeAmountIncreased", log); err != nil {
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

// ParseMaxCommunityStakeAmountIncreased is a log parse operation binding the contract event 0xb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df7.
//
// Solidity: event MaxCommunityStakeAmountIncreased(uint256 maxStakeAmount)
func (_StakingPoolLib *StakingPoolLibFilterer) ParseMaxCommunityStakeAmountIncreased(log types.Log) (*StakingPoolLibMaxCommunityStakeAmountIncreased, error) {
	event := new(StakingPoolLibMaxCommunityStakeAmountIncreased)
	if err := _StakingPoolLib.contract.UnpackLog(event, "MaxCommunityStakeAmountIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibMaxOperatorStakeAmountIncreasedIterator is returned from FilterMaxOperatorStakeAmountIncreased and is used to iterate over the raw logs and unpacked data for MaxOperatorStakeAmountIncreased events raised by the StakingPoolLib contract.
type StakingPoolLibMaxOperatorStakeAmountIncreasedIterator struct {
	Event *StakingPoolLibMaxOperatorStakeAmountIncreased // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibMaxOperatorStakeAmountIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibMaxOperatorStakeAmountIncreased)
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
		it.Event = new(StakingPoolLibMaxOperatorStakeAmountIncreased)
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
func (it *StakingPoolLibMaxOperatorStakeAmountIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibMaxOperatorStakeAmountIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibMaxOperatorStakeAmountIncreased represents a MaxOperatorStakeAmountIncreased event raised by the StakingPoolLib contract.
type StakingPoolLibMaxOperatorStakeAmountIncreased struct {
	MaxStakeAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMaxOperatorStakeAmountIncreased is a free log retrieval operation binding the contract event 0x816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52.
//
// Solidity: event MaxOperatorStakeAmountIncreased(uint256 maxStakeAmount)
func (_StakingPoolLib *StakingPoolLibFilterer) FilterMaxOperatorStakeAmountIncreased(opts *bind.FilterOpts) (*StakingPoolLibMaxOperatorStakeAmountIncreasedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "MaxOperatorStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibMaxOperatorStakeAmountIncreasedIterator{contract: _StakingPoolLib.contract, event: "MaxOperatorStakeAmountIncreased", logs: logs, sub: sub}, nil
}

// WatchMaxOperatorStakeAmountIncreased is a free log subscription operation binding the contract event 0x816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52.
//
// Solidity: event MaxOperatorStakeAmountIncreased(uint256 maxStakeAmount)
func (_StakingPoolLib *StakingPoolLibFilterer) WatchMaxOperatorStakeAmountIncreased(opts *bind.WatchOpts, sink chan<- *StakingPoolLibMaxOperatorStakeAmountIncreased) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "MaxOperatorStakeAmountIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibMaxOperatorStakeAmountIncreased)
				if err := _StakingPoolLib.contract.UnpackLog(event, "MaxOperatorStakeAmountIncreased", log); err != nil {
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

// ParseMaxOperatorStakeAmountIncreased is a log parse operation binding the contract event 0x816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52.
//
// Solidity: event MaxOperatorStakeAmountIncreased(uint256 maxStakeAmount)
func (_StakingPoolLib *StakingPoolLibFilterer) ParseMaxOperatorStakeAmountIncreased(log types.Log) (*StakingPoolLibMaxOperatorStakeAmountIncreased, error) {
	event := new(StakingPoolLibMaxOperatorStakeAmountIncreased)
	if err := _StakingPoolLib.contract.UnpackLog(event, "MaxOperatorStakeAmountIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the StakingPoolLib contract.
type StakingPoolLibOperatorAddedIterator struct {
	Event *StakingPoolLibOperatorAdded // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibOperatorAdded)
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
		it.Event = new(StakingPoolLibOperatorAdded)
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
func (it *StakingPoolLibOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibOperatorAdded represents a OperatorAdded event raised by the StakingPoolLib contract.
type StakingPoolLibOperatorAdded struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d.
//
// Solidity: event OperatorAdded(address operator)
func (_StakingPoolLib *StakingPoolLibFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*StakingPoolLibOperatorAddedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibOperatorAddedIterator{contract: _StakingPoolLib.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d.
//
// Solidity: event OperatorAdded(address operator)
func (_StakingPoolLib *StakingPoolLibFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *StakingPoolLibOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibOperatorAdded)
				if err := _StakingPoolLib.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ParseOperatorAdded is a log parse operation binding the contract event 0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d.
//
// Solidity: event OperatorAdded(address operator)
func (_StakingPoolLib *StakingPoolLibFilterer) ParseOperatorAdded(log types.Log) (*StakingPoolLibOperatorAdded, error) {
	event := new(StakingPoolLibOperatorAdded)
	if err := _StakingPoolLib.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibOperatorRemovedIterator is returned from FilterOperatorRemoved and is used to iterate over the raw logs and unpacked data for OperatorRemoved events raised by the StakingPoolLib contract.
type StakingPoolLibOperatorRemovedIterator struct {
	Event *StakingPoolLibOperatorRemoved // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibOperatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibOperatorRemoved)
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
		it.Event = new(StakingPoolLibOperatorRemoved)
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
func (it *StakingPoolLibOperatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibOperatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibOperatorRemoved represents a OperatorRemoved event raised by the StakingPoolLib contract.
type StakingPoolLibOperatorRemoved struct {
	Operator common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorRemoved is a free log retrieval operation binding the contract event 0x2360404a74478febece1a14f11275f22ada88d19ef96f7d785913010bfff4479.
//
// Solidity: event OperatorRemoved(address operator, uint256 amount)
func (_StakingPoolLib *StakingPoolLibFilterer) FilterOperatorRemoved(opts *bind.FilterOpts) (*StakingPoolLibOperatorRemovedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "OperatorRemoved")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibOperatorRemovedIterator{contract: _StakingPoolLib.contract, event: "OperatorRemoved", logs: logs, sub: sub}, nil
}

// WatchOperatorRemoved is a free log subscription operation binding the contract event 0x2360404a74478febece1a14f11275f22ada88d19ef96f7d785913010bfff4479.
//
// Solidity: event OperatorRemoved(address operator, uint256 amount)
func (_StakingPoolLib *StakingPoolLibFilterer) WatchOperatorRemoved(opts *bind.WatchOpts, sink chan<- *StakingPoolLibOperatorRemoved) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "OperatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibOperatorRemoved)
				if err := _StakingPoolLib.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
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

// ParseOperatorRemoved is a log parse operation binding the contract event 0x2360404a74478febece1a14f11275f22ada88d19ef96f7d785913010bfff4479.
//
// Solidity: event OperatorRemoved(address operator, uint256 amount)
func (_StakingPoolLib *StakingPoolLibFilterer) ParseOperatorRemoved(log types.Log) (*StakingPoolLibOperatorRemoved, error) {
	event := new(StakingPoolLibOperatorRemoved)
	if err := _StakingPoolLib.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibPoolConcludedIterator is returned from FilterPoolConcluded and is used to iterate over the raw logs and unpacked data for PoolConcluded events raised by the StakingPoolLib contract.
type StakingPoolLibPoolConcludedIterator struct {
	Event *StakingPoolLibPoolConcluded // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibPoolConcludedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibPoolConcluded)
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
		it.Event = new(StakingPoolLibPoolConcluded)
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
func (it *StakingPoolLibPoolConcludedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibPoolConcludedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibPoolConcluded represents a PoolConcluded event raised by the StakingPoolLib contract.
type StakingPoolLibPoolConcluded struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPoolConcluded is a free log retrieval operation binding the contract event 0xf7d0e0f15586495da8c687328ead30fb829d9da55538cb0ef73dd229e517cdb8.
//
// Solidity: event PoolConcluded()
func (_StakingPoolLib *StakingPoolLibFilterer) FilterPoolConcluded(opts *bind.FilterOpts) (*StakingPoolLibPoolConcludedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "PoolConcluded")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibPoolConcludedIterator{contract: _StakingPoolLib.contract, event: "PoolConcluded", logs: logs, sub: sub}, nil
}

// WatchPoolConcluded is a free log subscription operation binding the contract event 0xf7d0e0f15586495da8c687328ead30fb829d9da55538cb0ef73dd229e517cdb8.
//
// Solidity: event PoolConcluded()
func (_StakingPoolLib *StakingPoolLibFilterer) WatchPoolConcluded(opts *bind.WatchOpts, sink chan<- *StakingPoolLibPoolConcluded) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "PoolConcluded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibPoolConcluded)
				if err := _StakingPoolLib.contract.UnpackLog(event, "PoolConcluded", log); err != nil {
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

// ParsePoolConcluded is a log parse operation binding the contract event 0xf7d0e0f15586495da8c687328ead30fb829d9da55538cb0ef73dd229e517cdb8.
//
// Solidity: event PoolConcluded()
func (_StakingPoolLib *StakingPoolLibFilterer) ParsePoolConcluded(log types.Log) (*StakingPoolLibPoolConcluded, error) {
	event := new(StakingPoolLibPoolConcluded)
	if err := _StakingPoolLib.contract.UnpackLog(event, "PoolConcluded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibPoolOpenedIterator is returned from FilterPoolOpened and is used to iterate over the raw logs and unpacked data for PoolOpened events raised by the StakingPoolLib contract.
type StakingPoolLibPoolOpenedIterator struct {
	Event *StakingPoolLibPoolOpened // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibPoolOpenedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibPoolOpened)
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
		it.Event = new(StakingPoolLibPoolOpened)
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
func (it *StakingPoolLibPoolOpenedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibPoolOpenedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibPoolOpened represents a PoolOpened event raised by the StakingPoolLib contract.
type StakingPoolLibPoolOpened struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPoolOpened is a free log retrieval operation binding the contract event 0xded6ebf04e261e1eb2f3e3b268a2e6aee5b478c15b341eba5cf18b9bc80c2e63.
//
// Solidity: event PoolOpened()
func (_StakingPoolLib *StakingPoolLibFilterer) FilterPoolOpened(opts *bind.FilterOpts) (*StakingPoolLibPoolOpenedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "PoolOpened")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibPoolOpenedIterator{contract: _StakingPoolLib.contract, event: "PoolOpened", logs: logs, sub: sub}, nil
}

// WatchPoolOpened is a free log subscription operation binding the contract event 0xded6ebf04e261e1eb2f3e3b268a2e6aee5b478c15b341eba5cf18b9bc80c2e63.
//
// Solidity: event PoolOpened()
func (_StakingPoolLib *StakingPoolLibFilterer) WatchPoolOpened(opts *bind.WatchOpts, sink chan<- *StakingPoolLibPoolOpened) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "PoolOpened")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibPoolOpened)
				if err := _StakingPoolLib.contract.UnpackLog(event, "PoolOpened", log); err != nil {
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

// ParsePoolOpened is a log parse operation binding the contract event 0xded6ebf04e261e1eb2f3e3b268a2e6aee5b478c15b341eba5cf18b9bc80c2e63.
//
// Solidity: event PoolOpened()
func (_StakingPoolLib *StakingPoolLibFilterer) ParsePoolOpened(log types.Log) (*StakingPoolLibPoolOpened, error) {
	event := new(StakingPoolLibPoolOpened)
	if err := _StakingPoolLib.contract.UnpackLog(event, "PoolOpened", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolLibPoolSizeIncreasedIterator is returned from FilterPoolSizeIncreased and is used to iterate over the raw logs and unpacked data for PoolSizeIncreased events raised by the StakingPoolLib contract.
type StakingPoolLibPoolSizeIncreasedIterator struct {
	Event *StakingPoolLibPoolSizeIncreased // Event containing the contract specifics and raw log

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
func (it *StakingPoolLibPoolSizeIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolLibPoolSizeIncreased)
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
		it.Event = new(StakingPoolLibPoolSizeIncreased)
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
func (it *StakingPoolLibPoolSizeIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolLibPoolSizeIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolLibPoolSizeIncreased represents a PoolSizeIncreased event raised by the StakingPoolLib contract.
type StakingPoolLibPoolSizeIncreased struct {
	MaxPoolSize *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPoolSizeIncreased is a free log retrieval operation binding the contract event 0x7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c0.
//
// Solidity: event PoolSizeIncreased(uint256 maxPoolSize)
func (_StakingPoolLib *StakingPoolLibFilterer) FilterPoolSizeIncreased(opts *bind.FilterOpts) (*StakingPoolLibPoolSizeIncreasedIterator, error) {

	logs, sub, err := _StakingPoolLib.contract.FilterLogs(opts, "PoolSizeIncreased")
	if err != nil {
		return nil, err
	}
	return &StakingPoolLibPoolSizeIncreasedIterator{contract: _StakingPoolLib.contract, event: "PoolSizeIncreased", logs: logs, sub: sub}, nil
}

// WatchPoolSizeIncreased is a free log subscription operation binding the contract event 0x7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c0.
//
// Solidity: event PoolSizeIncreased(uint256 maxPoolSize)
func (_StakingPoolLib *StakingPoolLibFilterer) WatchPoolSizeIncreased(opts *bind.WatchOpts, sink chan<- *StakingPoolLibPoolSizeIncreased) (event.Subscription, error) {

	logs, sub, err := _StakingPoolLib.contract.WatchLogs(opts, "PoolSizeIncreased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolLibPoolSizeIncreased)
				if err := _StakingPoolLib.contract.UnpackLog(event, "PoolSizeIncreased", log); err != nil {
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

// ParsePoolSizeIncreased is a log parse operation binding the contract event 0x7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c0.
//
// Solidity: event PoolSizeIncreased(uint256 maxPoolSize)
func (_StakingPoolLib *StakingPoolLibFilterer) ParsePoolSizeIncreased(log types.Log) (*StakingPoolLibPoolSizeIncreased, error) {
	event := new(StakingPoolLibPoolSizeIncreased)
	if err := _StakingPoolLib.contract.UnpackLog(event, "PoolSizeIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
