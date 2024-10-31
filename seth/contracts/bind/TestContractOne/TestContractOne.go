// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TestContractOne

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

// UniqueEventOneMetaData contains all meta data concerning the UniqueEventOne contract.
var UniqueEventOneMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"a\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"b\",\"type\":\"int256\"}],\"name\":\"NonUniqueEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"executeFirstOperation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101f2806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806312a27a5e14610030575b600080fd5b61004a600480360381019061004591906100df565b610060565b604051610057919061012e565b60405180910390f35b600081837f192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e60405160405180910390a3818361009c9190610178565b905092915050565b600080fd5b6000819050919050565b6100bc816100a9565b81146100c757600080fd5b50565b6000813590506100d9816100b3565b92915050565b600080604083850312156100f6576100f56100a4565b5b6000610104858286016100ca565b9250506020610115858286016100ca565b9150509250929050565b610128816100a9565b82525050565b6000602082019050610143600083018461011f565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610183826100a9565b915061018e836100a9565b9250828201905082811215600083121683821260008412151617156101b6576101b5610149565b5b9291505056fea26469706673582212203a9131e072ba50bd7aedd7d16ac86eb26a9e83acd9a8736e900ad00a4b689ae664736f6c63430008130033",
}

// UniqueEventOneABI is the input ABI used to generate the binding from.
// Deprecated: Use UniqueEventOneMetaData.ABI instead.
var UniqueEventOneABI = UniqueEventOneMetaData.ABI

// UniqueEventOneBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniqueEventOneMetaData.Bin instead.
var UniqueEventOneBin = UniqueEventOneMetaData.Bin

// DeployUniqueEventOne deploys a new Ethereum contract, binding an instance of UniqueEventOne to it.
func DeployUniqueEventOne(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniqueEventOne, error) {
	parsed, err := UniqueEventOneMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniqueEventOneBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniqueEventOne{UniqueEventOneCaller: UniqueEventOneCaller{contract: contract}, UniqueEventOneTransactor: UniqueEventOneTransactor{contract: contract}, UniqueEventOneFilterer: UniqueEventOneFilterer{contract: contract}}, nil
}

// UniqueEventOne is an auto generated Go binding around an Ethereum contract.
type UniqueEventOne struct {
	UniqueEventOneCaller     // Read-only binding to the contract
	UniqueEventOneTransactor // Write-only binding to the contract
	UniqueEventOneFilterer   // Log filterer for contract events
}

// UniqueEventOneCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniqueEventOneCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniqueEventOneTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniqueEventOneTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniqueEventOneFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniqueEventOneFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniqueEventOneSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniqueEventOneSession struct {
	Contract     *UniqueEventOne   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniqueEventOneCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniqueEventOneCallerSession struct {
	Contract *UniqueEventOneCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// UniqueEventOneTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniqueEventOneTransactorSession struct {
	Contract     *UniqueEventOneTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// UniqueEventOneRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniqueEventOneRaw struct {
	Contract *UniqueEventOne // Generic contract binding to access the raw methods on
}

// UniqueEventOneCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniqueEventOneCallerRaw struct {
	Contract *UniqueEventOneCaller // Generic read-only contract binding to access the raw methods on
}

// UniqueEventOneTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniqueEventOneTransactorRaw struct {
	Contract *UniqueEventOneTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniqueEventOne creates a new instance of UniqueEventOne, bound to a specific deployed contract.
func NewUniqueEventOne(address common.Address, backend bind.ContractBackend) (*UniqueEventOne, error) {
	contract, err := bindUniqueEventOne(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniqueEventOne{UniqueEventOneCaller: UniqueEventOneCaller{contract: contract}, UniqueEventOneTransactor: UniqueEventOneTransactor{contract: contract}, UniqueEventOneFilterer: UniqueEventOneFilterer{contract: contract}}, nil
}

// NewUniqueEventOneCaller creates a new read-only instance of UniqueEventOne, bound to a specific deployed contract.
func NewUniqueEventOneCaller(address common.Address, caller bind.ContractCaller) (*UniqueEventOneCaller, error) {
	contract, err := bindUniqueEventOne(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniqueEventOneCaller{contract: contract}, nil
}

// NewUniqueEventOneTransactor creates a new write-only instance of UniqueEventOne, bound to a specific deployed contract.
func NewUniqueEventOneTransactor(address common.Address, transactor bind.ContractTransactor) (*UniqueEventOneTransactor, error) {
	contract, err := bindUniqueEventOne(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniqueEventOneTransactor{contract: contract}, nil
}

// NewUniqueEventOneFilterer creates a new log filterer instance of UniqueEventOne, bound to a specific deployed contract.
func NewUniqueEventOneFilterer(address common.Address, filterer bind.ContractFilterer) (*UniqueEventOneFilterer, error) {
	contract, err := bindUniqueEventOne(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniqueEventOneFilterer{contract: contract}, nil
}

// bindUniqueEventOne binds a generic wrapper to an already deployed contract.
func bindUniqueEventOne(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniqueEventOneMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniqueEventOne *UniqueEventOneRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniqueEventOne.Contract.UniqueEventOneCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniqueEventOne *UniqueEventOneRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniqueEventOne.Contract.UniqueEventOneTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniqueEventOne *UniqueEventOneRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniqueEventOne.Contract.UniqueEventOneTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniqueEventOne *UniqueEventOneCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniqueEventOne.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniqueEventOne *UniqueEventOneTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniqueEventOne.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniqueEventOne *UniqueEventOneTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniqueEventOne.Contract.contract.Transact(opts, method, params...)
}

// ExecuteFirstOperation is a paid mutator transaction binding the contract method 0x12a27a5e.
//
// Solidity: function executeFirstOperation(int256 x, int256 y) returns(int256)
func (_UniqueEventOne *UniqueEventOneTransactor) ExecuteFirstOperation(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _UniqueEventOne.contract.Transact(opts, "executeFirstOperation", x, y)
}

// ExecuteFirstOperation is a paid mutator transaction binding the contract method 0x12a27a5e.
//
// Solidity: function executeFirstOperation(int256 x, int256 y) returns(int256)
func (_UniqueEventOne *UniqueEventOneSession) ExecuteFirstOperation(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _UniqueEventOne.Contract.ExecuteFirstOperation(&_UniqueEventOne.TransactOpts, x, y)
}

// ExecuteFirstOperation is a paid mutator transaction binding the contract method 0x12a27a5e.
//
// Solidity: function executeFirstOperation(int256 x, int256 y) returns(int256)
func (_UniqueEventOne *UniqueEventOneTransactorSession) ExecuteFirstOperation(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _UniqueEventOne.Contract.ExecuteFirstOperation(&_UniqueEventOne.TransactOpts, x, y)
}

// UniqueEventOneNonUniqueEventIterator is returned from FilterNonUniqueEvent and is used to iterate over the raw logs and unpacked data for NonUniqueEvent events raised by the UniqueEventOne contract.
type UniqueEventOneNonUniqueEventIterator struct {
	Event *UniqueEventOneNonUniqueEvent // Event containing the contract specifics and raw log

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
func (it *UniqueEventOneNonUniqueEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniqueEventOneNonUniqueEvent)
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
		it.Event = new(UniqueEventOneNonUniqueEvent)
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
func (it *UniqueEventOneNonUniqueEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniqueEventOneNonUniqueEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniqueEventOneNonUniqueEvent represents a NonUniqueEvent event raised by the UniqueEventOne contract.
type UniqueEventOneNonUniqueEvent struct {
	A   *big.Int
	B   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNonUniqueEvent is a free log retrieval operation binding the contract event 0x192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e.
//
// Solidity: event NonUniqueEvent(int256 indexed a, int256 indexed b)
func (_UniqueEventOne *UniqueEventOneFilterer) FilterNonUniqueEvent(opts *bind.FilterOpts, a []*big.Int, b []*big.Int) (*UniqueEventOneNonUniqueEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _UniqueEventOne.contract.FilterLogs(opts, "NonUniqueEvent", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return &UniqueEventOneNonUniqueEventIterator{contract: _UniqueEventOne.contract, event: "NonUniqueEvent", logs: logs, sub: sub}, nil
}

// WatchNonUniqueEvent is a free log subscription operation binding the contract event 0x192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e.
//
// Solidity: event NonUniqueEvent(int256 indexed a, int256 indexed b)
func (_UniqueEventOne *UniqueEventOneFilterer) WatchNonUniqueEvent(opts *bind.WatchOpts, sink chan<- *UniqueEventOneNonUniqueEvent, a []*big.Int, b []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _UniqueEventOne.contract.WatchLogs(opts, "NonUniqueEvent", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniqueEventOneNonUniqueEvent)
				if err := _UniqueEventOne.contract.UnpackLog(event, "NonUniqueEvent", log); err != nil {
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
func (_UniqueEventOne *UniqueEventOneFilterer) ParseNonUniqueEvent(log types.Log) (*UniqueEventOneNonUniqueEvent, error) {
	event := new(UniqueEventOneNonUniqueEvent)
	if err := _UniqueEventOne.contract.UnpackLog(event, "NonUniqueEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
