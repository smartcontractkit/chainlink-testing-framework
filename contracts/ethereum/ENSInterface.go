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

// ENSInterfaceABI is the input ABI used to generate the binding from.
const ENSInterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"label\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NewOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"NewResolver\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"ttl\",\"type\":\"uint64\"}],\"name\":\"NewTTL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"resolver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_resolver\",\"type\":\"address\"}],\"name\":\"setResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"label\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"setSubnodeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_ttl\",\"type\":\"uint64\"}],\"name\":\"setTTL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"ttl\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ENSInterface is an auto generated Go binding around an Ethereum contract.
type ENSInterface struct {
	ENSInterfaceCaller     // Read-only binding to the contract
	ENSInterfaceTransactor // Write-only binding to the contract
	ENSInterfaceFilterer   // Log filterer for contract events
}

// ENSInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type ENSInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ENSInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ENSInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ENSInterfaceSession struct {
	Contract     *ENSInterface     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ENSInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ENSInterfaceCallerSession struct {
	Contract *ENSInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ENSInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ENSInterfaceTransactorSession struct {
	Contract     *ENSInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ENSInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type ENSInterfaceRaw struct {
	Contract *ENSInterface // Generic contract binding to access the raw methods on
}

// ENSInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ENSInterfaceCallerRaw struct {
	Contract *ENSInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// ENSInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ENSInterfaceTransactorRaw struct {
	Contract *ENSInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewENSInterface creates a new instance of ENSInterface, bound to a specific deployed contract.
func NewENSInterface(address common.Address, backend bind.ContractBackend) (*ENSInterface, error) {
	contract, err := bindENSInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ENSInterface{ENSInterfaceCaller: ENSInterfaceCaller{contract: contract}, ENSInterfaceTransactor: ENSInterfaceTransactor{contract: contract}, ENSInterfaceFilterer: ENSInterfaceFilterer{contract: contract}}, nil
}

// NewENSInterfaceCaller creates a new read-only instance of ENSInterface, bound to a specific deployed contract.
func NewENSInterfaceCaller(address common.Address, caller bind.ContractCaller) (*ENSInterfaceCaller, error) {
	contract, err := bindENSInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceCaller{contract: contract}, nil
}

// NewENSInterfaceTransactor creates a new write-only instance of ENSInterface, bound to a specific deployed contract.
func NewENSInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*ENSInterfaceTransactor, error) {
	contract, err := bindENSInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceTransactor{contract: contract}, nil
}

// NewENSInterfaceFilterer creates a new log filterer instance of ENSInterface, bound to a specific deployed contract.
func NewENSInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*ENSInterfaceFilterer, error) {
	contract, err := bindENSInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceFilterer{contract: contract}, nil
}

// bindENSInterface binds a generic wrapper to an already deployed contract.
func bindENSInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ENSInterface *ENSInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ENSInterface.Contract.ENSInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ENSInterface *ENSInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENSInterface.Contract.ENSInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ENSInterface *ENSInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ENSInterface.Contract.ENSInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ENSInterface *ENSInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ENSInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ENSInterface *ENSInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENSInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ENSInterface *ENSInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ENSInterface.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(bytes32 node) view returns(address)
func (_ENSInterface *ENSInterfaceCaller) Owner(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ENSInterface.contract.Call(opts, &out, "owner", node)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(bytes32 node) view returns(address)
func (_ENSInterface *ENSInterfaceSession) Owner(node [32]byte) (common.Address, error) {
	return _ENSInterface.Contract.Owner(&_ENSInterface.CallOpts, node)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(bytes32 node) view returns(address)
func (_ENSInterface *ENSInterfaceCallerSession) Owner(node [32]byte) (common.Address, error) {
	return _ENSInterface.Contract.Owner(&_ENSInterface.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 node) view returns(address)
func (_ENSInterface *ENSInterfaceCaller) Resolver(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ENSInterface.contract.Call(opts, &out, "resolver", node)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 node) view returns(address)
func (_ENSInterface *ENSInterfaceSession) Resolver(node [32]byte) (common.Address, error) {
	return _ENSInterface.Contract.Resolver(&_ENSInterface.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 node) view returns(address)
func (_ENSInterface *ENSInterfaceCallerSession) Resolver(node [32]byte) (common.Address, error) {
	return _ENSInterface.Contract.Resolver(&_ENSInterface.CallOpts, node)
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// Solidity: function ttl(bytes32 node) view returns(uint64)
func (_ENSInterface *ENSInterfaceCaller) Ttl(opts *bind.CallOpts, node [32]byte) (uint64, error) {
	var out []interface{}
	err := _ENSInterface.contract.Call(opts, &out, "ttl", node)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// Solidity: function ttl(bytes32 node) view returns(uint64)
func (_ENSInterface *ENSInterfaceSession) Ttl(node [32]byte) (uint64, error) {
	return _ENSInterface.Contract.Ttl(&_ENSInterface.CallOpts, node)
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// Solidity: function ttl(bytes32 node) view returns(uint64)
func (_ENSInterface *ENSInterfaceCallerSession) Ttl(node [32]byte) (uint64, error) {
	return _ENSInterface.Contract.Ttl(&_ENSInterface.CallOpts, node)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// Solidity: function setOwner(bytes32 node, address _owner) returns()
func (_ENSInterface *ENSInterfaceTransactor) SetOwner(opts *bind.TransactOpts, node [32]byte, _owner common.Address) (*types.Transaction, error) {
	return _ENSInterface.contract.Transact(opts, "setOwner", node, _owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// Solidity: function setOwner(bytes32 node, address _owner) returns()
func (_ENSInterface *ENSInterfaceSession) SetOwner(node [32]byte, _owner common.Address) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetOwner(&_ENSInterface.TransactOpts, node, _owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// Solidity: function setOwner(bytes32 node, address _owner) returns()
func (_ENSInterface *ENSInterfaceTransactorSession) SetOwner(node [32]byte, _owner common.Address) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetOwner(&_ENSInterface.TransactOpts, node, _owner)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 node, address _resolver) returns()
func (_ENSInterface *ENSInterfaceTransactor) SetResolver(opts *bind.TransactOpts, node [32]byte, _resolver common.Address) (*types.Transaction, error) {
	return _ENSInterface.contract.Transact(opts, "setResolver", node, _resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 node, address _resolver) returns()
func (_ENSInterface *ENSInterfaceSession) SetResolver(node [32]byte, _resolver common.Address) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetResolver(&_ENSInterface.TransactOpts, node, _resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 node, address _resolver) returns()
func (_ENSInterface *ENSInterfaceTransactorSession) SetResolver(node [32]byte, _resolver common.Address) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetResolver(&_ENSInterface.TransactOpts, node, _resolver)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// Solidity: function setSubnodeOwner(bytes32 node, bytes32 label, address _owner) returns()
func (_ENSInterface *ENSInterfaceTransactor) SetSubnodeOwner(opts *bind.TransactOpts, node [32]byte, label [32]byte, _owner common.Address) (*types.Transaction, error) {
	return _ENSInterface.contract.Transact(opts, "setSubnodeOwner", node, label, _owner)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// Solidity: function setSubnodeOwner(bytes32 node, bytes32 label, address _owner) returns()
func (_ENSInterface *ENSInterfaceSession) SetSubnodeOwner(node [32]byte, label [32]byte, _owner common.Address) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetSubnodeOwner(&_ENSInterface.TransactOpts, node, label, _owner)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// Solidity: function setSubnodeOwner(bytes32 node, bytes32 label, address _owner) returns()
func (_ENSInterface *ENSInterfaceTransactorSession) SetSubnodeOwner(node [32]byte, label [32]byte, _owner common.Address) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetSubnodeOwner(&_ENSInterface.TransactOpts, node, label, _owner)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(bytes32 node, uint64 _ttl) returns()
func (_ENSInterface *ENSInterfaceTransactor) SetTTL(opts *bind.TransactOpts, node [32]byte, _ttl uint64) (*types.Transaction, error) {
	return _ENSInterface.contract.Transact(opts, "setTTL", node, _ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(bytes32 node, uint64 _ttl) returns()
func (_ENSInterface *ENSInterfaceSession) SetTTL(node [32]byte, _ttl uint64) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetTTL(&_ENSInterface.TransactOpts, node, _ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(bytes32 node, uint64 _ttl) returns()
func (_ENSInterface *ENSInterfaceTransactorSession) SetTTL(node [32]byte, _ttl uint64) (*types.Transaction, error) {
	return _ENSInterface.Contract.SetTTL(&_ENSInterface.TransactOpts, node, _ttl)
}

// ENSInterfaceNewOwnerIterator is returned from FilterNewOwner and is used to iterate over the raw logs and unpacked data for NewOwner events raised by the ENSInterface contract.
type ENSInterfaceNewOwnerIterator struct {
	Event *ENSInterfaceNewOwner // Event containing the contract specifics and raw log

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
func (it *ENSInterfaceNewOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSInterfaceNewOwner)
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
		it.Event = new(ENSInterfaceNewOwner)
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
func (it *ENSInterfaceNewOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSInterfaceNewOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSInterfaceNewOwner represents a NewOwner event raised by the ENSInterface contract.
type ENSInterfaceNewOwner struct {
	Node  [32]byte
	Label [32]byte
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNewOwner is a free log retrieval operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// Solidity: event NewOwner(bytes32 indexed node, bytes32 indexed label, address owner)
func (_ENSInterface *ENSInterfaceFilterer) FilterNewOwner(opts *bind.FilterOpts, node [][32]byte, label [][32]byte) (*ENSInterfaceNewOwnerIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	logs, sub, err := _ENSInterface.contract.FilterLogs(opts, "NewOwner", nodeRule, labelRule)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceNewOwnerIterator{contract: _ENSInterface.contract, event: "NewOwner", logs: logs, sub: sub}, nil
}

// WatchNewOwner is a free log subscription operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// Solidity: event NewOwner(bytes32 indexed node, bytes32 indexed label, address owner)
func (_ENSInterface *ENSInterfaceFilterer) WatchNewOwner(opts *bind.WatchOpts, sink chan<- *ENSInterfaceNewOwner, node [][32]byte, label [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	logs, sub, err := _ENSInterface.contract.WatchLogs(opts, "NewOwner", nodeRule, labelRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSInterfaceNewOwner)
				if err := _ENSInterface.contract.UnpackLog(event, "NewOwner", log); err != nil {
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

// ParseNewOwner is a log parse operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// Solidity: event NewOwner(bytes32 indexed node, bytes32 indexed label, address owner)
func (_ENSInterface *ENSInterfaceFilterer) ParseNewOwner(log types.Log) (*ENSInterfaceNewOwner, error) {
	event := new(ENSInterfaceNewOwner)
	if err := _ENSInterface.contract.UnpackLog(event, "NewOwner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ENSInterfaceNewResolverIterator is returned from FilterNewResolver and is used to iterate over the raw logs and unpacked data for NewResolver events raised by the ENSInterface contract.
type ENSInterfaceNewResolverIterator struct {
	Event *ENSInterfaceNewResolver // Event containing the contract specifics and raw log

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
func (it *ENSInterfaceNewResolverIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSInterfaceNewResolver)
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
		it.Event = new(ENSInterfaceNewResolver)
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
func (it *ENSInterfaceNewResolverIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSInterfaceNewResolverIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSInterfaceNewResolver represents a NewResolver event raised by the ENSInterface contract.
type ENSInterfaceNewResolver struct {
	Node     [32]byte
	Resolver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewResolver is a free log retrieval operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: event NewResolver(bytes32 indexed node, address resolver)
func (_ENSInterface *ENSInterfaceFilterer) FilterNewResolver(opts *bind.FilterOpts, node [][32]byte) (*ENSInterfaceNewResolverIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENSInterface.contract.FilterLogs(opts, "NewResolver", nodeRule)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceNewResolverIterator{contract: _ENSInterface.contract, event: "NewResolver", logs: logs, sub: sub}, nil
}

// WatchNewResolver is a free log subscription operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: event NewResolver(bytes32 indexed node, address resolver)
func (_ENSInterface *ENSInterfaceFilterer) WatchNewResolver(opts *bind.WatchOpts, sink chan<- *ENSInterfaceNewResolver, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENSInterface.contract.WatchLogs(opts, "NewResolver", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSInterfaceNewResolver)
				if err := _ENSInterface.contract.UnpackLog(event, "NewResolver", log); err != nil {
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

// ParseNewResolver is a log parse operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: event NewResolver(bytes32 indexed node, address resolver)
func (_ENSInterface *ENSInterfaceFilterer) ParseNewResolver(log types.Log) (*ENSInterfaceNewResolver, error) {
	event := new(ENSInterfaceNewResolver)
	if err := _ENSInterface.contract.UnpackLog(event, "NewResolver", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ENSInterfaceNewTTLIterator is returned from FilterNewTTL and is used to iterate over the raw logs and unpacked data for NewTTL events raised by the ENSInterface contract.
type ENSInterfaceNewTTLIterator struct {
	Event *ENSInterfaceNewTTL // Event containing the contract specifics and raw log

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
func (it *ENSInterfaceNewTTLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSInterfaceNewTTL)
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
		it.Event = new(ENSInterfaceNewTTL)
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
func (it *ENSInterfaceNewTTLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSInterfaceNewTTLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSInterfaceNewTTL represents a NewTTL event raised by the ENSInterface contract.
type ENSInterfaceNewTTL struct {
	Node [32]byte
	Ttl  uint64
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNewTTL is a free log retrieval operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// Solidity: event NewTTL(bytes32 indexed node, uint64 ttl)
func (_ENSInterface *ENSInterfaceFilterer) FilterNewTTL(opts *bind.FilterOpts, node [][32]byte) (*ENSInterfaceNewTTLIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENSInterface.contract.FilterLogs(opts, "NewTTL", nodeRule)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceNewTTLIterator{contract: _ENSInterface.contract, event: "NewTTL", logs: logs, sub: sub}, nil
}

// WatchNewTTL is a free log subscription operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// Solidity: event NewTTL(bytes32 indexed node, uint64 ttl)
func (_ENSInterface *ENSInterfaceFilterer) WatchNewTTL(opts *bind.WatchOpts, sink chan<- *ENSInterfaceNewTTL, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENSInterface.contract.WatchLogs(opts, "NewTTL", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSInterfaceNewTTL)
				if err := _ENSInterface.contract.UnpackLog(event, "NewTTL", log); err != nil {
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

// ParseNewTTL is a log parse operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// Solidity: event NewTTL(bytes32 indexed node, uint64 ttl)
func (_ENSInterface *ENSInterfaceFilterer) ParseNewTTL(log types.Log) (*ENSInterfaceNewTTL, error) {
	event := new(ENSInterfaceNewTTL)
	if err := _ENSInterface.contract.UnpackLog(event, "NewTTL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ENSInterfaceTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ENSInterface contract.
type ENSInterfaceTransferIterator struct {
	Event *ENSInterfaceTransfer // Event containing the contract specifics and raw log

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
func (it *ENSInterfaceTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSInterfaceTransfer)
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
		it.Event = new(ENSInterfaceTransfer)
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
func (it *ENSInterfaceTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSInterfaceTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSInterfaceTransfer represents a Transfer event raised by the ENSInterface contract.
type ENSInterfaceTransfer struct {
	Node  [32]byte
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// Solidity: event Transfer(bytes32 indexed node, address owner)
func (_ENSInterface *ENSInterfaceFilterer) FilterTransfer(opts *bind.FilterOpts, node [][32]byte) (*ENSInterfaceTransferIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENSInterface.contract.FilterLogs(opts, "Transfer", nodeRule)
	if err != nil {
		return nil, err
	}
	return &ENSInterfaceTransferIterator{contract: _ENSInterface.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// Solidity: event Transfer(bytes32 indexed node, address owner)
func (_ENSInterface *ENSInterfaceFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ENSInterfaceTransfer, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENSInterface.contract.WatchLogs(opts, "Transfer", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSInterfaceTransfer)
				if err := _ENSInterface.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// Solidity: event Transfer(bytes32 indexed node, address owner)
func (_ENSInterface *ENSInterfaceFilterer) ParseTransfer(log types.Log) (*ENSInterfaceTransfer, error) {
	event := new(ENSInterfaceTransfer)
	if err := _ENSInterface.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
