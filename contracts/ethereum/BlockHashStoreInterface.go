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

// BlockHashStoreInterfaceABI is the input ABI used to generate the binding from.
const BlockHashStoreInterfaceABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// BlockHashStoreInterface is an auto generated Go binding around an Ethereum contract.
type BlockHashStoreInterface struct {
	BlockHashStoreInterfaceCaller     // Read-only binding to the contract
	BlockHashStoreInterfaceTransactor // Write-only binding to the contract
	BlockHashStoreInterfaceFilterer   // Log filterer for contract events
}

// BlockHashStoreInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockHashStoreInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockHashStoreInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockHashStoreInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockHashStoreInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlockHashStoreInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockHashStoreInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockHashStoreInterfaceSession struct {
	Contract     *BlockHashStoreInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BlockHashStoreInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockHashStoreInterfaceCallerSession struct {
	Contract *BlockHashStoreInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// BlockHashStoreInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockHashStoreInterfaceTransactorSession struct {
	Contract     *BlockHashStoreInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// BlockHashStoreInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockHashStoreInterfaceRaw struct {
	Contract *BlockHashStoreInterface // Generic contract binding to access the raw methods on
}

// BlockHashStoreInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockHashStoreInterfaceCallerRaw struct {
	Contract *BlockHashStoreInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// BlockHashStoreInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockHashStoreInterfaceTransactorRaw struct {
	Contract *BlockHashStoreInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockHashStoreInterface creates a new instance of BlockHashStoreInterface, bound to a specific deployed contract.
func NewBlockHashStoreInterface(address common.Address, backend bind.ContractBackend) (*BlockHashStoreInterface, error) {
	contract, err := bindBlockHashStoreInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BlockHashStoreInterface{BlockHashStoreInterfaceCaller: BlockHashStoreInterfaceCaller{contract: contract}, BlockHashStoreInterfaceTransactor: BlockHashStoreInterfaceTransactor{contract: contract}, BlockHashStoreInterfaceFilterer: BlockHashStoreInterfaceFilterer{contract: contract}}, nil
}

// NewBlockHashStoreInterfaceCaller creates a new read-only instance of BlockHashStoreInterface, bound to a specific deployed contract.
func NewBlockHashStoreInterfaceCaller(address common.Address, caller bind.ContractCaller) (*BlockHashStoreInterfaceCaller, error) {
	contract, err := bindBlockHashStoreInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockHashStoreInterfaceCaller{contract: contract}, nil
}

// NewBlockHashStoreInterfaceTransactor creates a new write-only instance of BlockHashStoreInterface, bound to a specific deployed contract.
func NewBlockHashStoreInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockHashStoreInterfaceTransactor, error) {
	contract, err := bindBlockHashStoreInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockHashStoreInterfaceTransactor{contract: contract}, nil
}

// NewBlockHashStoreInterfaceFilterer creates a new log filterer instance of BlockHashStoreInterface, bound to a specific deployed contract.
func NewBlockHashStoreInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockHashStoreInterfaceFilterer, error) {
	contract, err := bindBlockHashStoreInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockHashStoreInterfaceFilterer{contract: contract}, nil
}

// bindBlockHashStoreInterface binds a generic wrapper to an already deployed contract.
func bindBlockHashStoreInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BlockHashStoreInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockHashStoreInterface *BlockHashStoreInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockHashStoreInterface.Contract.BlockHashStoreInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockHashStoreInterface *BlockHashStoreInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockHashStoreInterface.Contract.BlockHashStoreInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockHashStoreInterface *BlockHashStoreInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockHashStoreInterface.Contract.BlockHashStoreInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockHashStoreInterface *BlockHashStoreInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockHashStoreInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockHashStoreInterface *BlockHashStoreInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockHashStoreInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockHashStoreInterface *BlockHashStoreInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockHashStoreInterface.Contract.contract.Transact(opts, method, params...)
}

// GetBlockhash is a free data retrieval call binding the contract method 0xe9413d38.
//
// Solidity: function getBlockhash(uint256 number) view returns(bytes32)
func (_BlockHashStoreInterface *BlockHashStoreInterfaceCaller) GetBlockhash(opts *bind.CallOpts, number *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BlockHashStoreInterface.contract.Call(opts, &out, "getBlockhash", number)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockhash is a free data retrieval call binding the contract method 0xe9413d38.
//
// Solidity: function getBlockhash(uint256 number) view returns(bytes32)
func (_BlockHashStoreInterface *BlockHashStoreInterfaceSession) GetBlockhash(number *big.Int) ([32]byte, error) {
	return _BlockHashStoreInterface.Contract.GetBlockhash(&_BlockHashStoreInterface.CallOpts, number)
}

// GetBlockhash is a free data retrieval call binding the contract method 0xe9413d38.
//
// Solidity: function getBlockhash(uint256 number) view returns(bytes32)
func (_BlockHashStoreInterface *BlockHashStoreInterfaceCallerSession) GetBlockhash(number *big.Int) ([32]byte, error) {
	return _BlockHashStoreInterface.Contract.GetBlockhash(&_BlockHashStoreInterface.CallOpts, number)
}
