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

// VRFConsumerMetaData contains all meta data concerning the VRFConsumer contract.
var VRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"randomnessOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b5060405161047738038061047783398101604081905261002f91610062565b6001600160a01b0391821660a05216608052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a0516103b76100c06000396000818160ba01526101640152600061013501526103b76000f3fe608060405234801561001057600080fd5b506004361061004b5760003560e01c80626d6cae146100505780632f47fd861461006b578063866ee7481461007457806394985ddd14610087575b600080fd5b61005960025481565b60405190815260200160405180910390f35b61005960015481565b6100596100823660046102ab565b61009c565b61009a6100953660046102ab565b6100af565b005b60006100a88383610131565b9392505050565b336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461012b5760405162461bcd60e51b815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c00604482015260640160405180910390fd5b60015550565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316634000aea07f0000000000000000000000000000000000000000000000000000000000000000848660006040516020016101a1929190918252602082015260400190565b6040516020818303038152906040526040518463ffffffff1660e01b81526004016101ce939291906102cd565b6020604051808303816000875af11580156101ed573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102119190610339565b5060008381526020818152604080832054815180840188905280830185905230606082015260808082018390528351808303909101815260a09091019092528151918301919091208684529290915261026b90600161035b565b6000858152602081815260409182902092909255805180830187905280820184905281518082038301815260609091019091528051910120949350505050565b600080604083850312156102be57600080fd5b50508035926020909101359150565b60018060a01b038416815260006020848184015260606040840152835180606085015260005b8181101561030f578581018301518582016080015282016102f3565b81811115610321576000608083870101525b50601f01601f19169290920160800195945050505050565b60006020828403121561034b57600080fd5b815180151581146100a857600080fd5b6000821982111561037c57634e487b7160e01b600052601160045260246000fd5b50019056fea26469706673582212205b43ce73c510b47d8ddc6008fe7a20ffeff82df7b3179da0175e39ff8464b37a64736f6c634300080d0033",
}

// VRFConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFConsumerMetaData.ABI instead.
var VRFConsumerABI = VRFConsumerMetaData.ABI

// VRFConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VRFConsumerMetaData.Bin instead.
var VRFConsumerBin = VRFConsumerMetaData.Bin

// DeployVRFConsumer deploys a new Ethereum contract, binding an instance of VRFConsumer to it.
func DeployVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFConsumer, error) {
	parsed, err := VRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumer{VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
}

// VRFConsumer is an auto generated Go binding around an Ethereum contract.
type VRFConsumer struct {
	VRFConsumerCaller     // Read-only binding to the contract
	VRFConsumerTransactor // Write-only binding to the contract
	VRFConsumerFilterer   // Log filterer for contract events
}

// VRFConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFConsumerSession struct {
	Contract     *VRFConsumer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFConsumerCallerSession struct {
	Contract *VRFConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// VRFConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFConsumerTransactorSession struct {
	Contract     *VRFConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// VRFConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFConsumerRaw struct {
	Contract *VRFConsumer // Generic contract binding to access the raw methods on
}

// VRFConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFConsumerCallerRaw struct {
	Contract *VRFConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// VRFConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFConsumerTransactorRaw struct {
	Contract *VRFConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFConsumer creates a new instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumer(address common.Address, backend bind.ContractBackend) (*VRFConsumer, error) {
	contract, err := bindVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumer{VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
}

// NewVRFConsumerCaller creates a new read-only instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerCaller, error) {
	contract, err := bindVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerCaller{contract: contract}, nil
}

// NewVRFConsumerTransactor creates a new write-only instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerTransactor, error) {
	contract, err := bindVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerTransactor{contract: contract}, nil
}

// NewVRFConsumerFilterer creates a new log filterer instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerFilterer, error) {
	contract, err := bindVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerFilterer{contract: contract}, nil
}

// bindVRFConsumer binds a generic wrapper to an already deployed contract.
func bindVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumer *VRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumer.Contract.VRFConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumer *VRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumer.Contract.VRFConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumer *VRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumer.Contract.VRFConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumer *VRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumer *VRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumer *VRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumer.Contract.contract.Transact(opts, method, params...)
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() view returns(uint256)
func (_VRFConsumer *VRFConsumerCaller) RandomnessOutput(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "randomnessOutput")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() view returns(uint256)
func (_VRFConsumer *VRFConsumerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() view returns(uint256)
func (_VRFConsumer *VRFConsumerCallerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() view returns(bytes32)
func (_VRFConsumer *VRFConsumerCaller) RequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "requestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() view returns(bytes32)
func (_VRFConsumer *VRFConsumerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() view returns(bytes32)
func (_VRFConsumer *VRFConsumerCallerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumer *VRFConsumerTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumer *VRFConsumerSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RawFulfillRandomness(&_VRFConsumer.TransactOpts, requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumer *VRFConsumerTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RawFulfillRandomness(&_VRFConsumer.TransactOpts, requestId, randomness)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x866ee748.
//
// Solidity: function testRequestRandomness(bytes32 keyHash, uint256 fee) returns(bytes32)
func (_VRFConsumer *VRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "testRequestRandomness", keyHash, fee)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x866ee748.
//
// Solidity: function testRequestRandomness(bytes32 keyHash, uint256 fee) returns(bytes32)
func (_VRFConsumer *VRFConsumerSession) TestRequestRandomness(keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, keyHash, fee)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x866ee748.
//
// Solidity: function testRequestRandomness(bytes32 keyHash, uint256 fee) returns(bytes32)
func (_VRFConsumer *VRFConsumerTransactorSession) TestRequestRandomness(keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, keyHash, fee)
}
