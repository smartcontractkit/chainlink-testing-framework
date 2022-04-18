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

// VRFConsumerV2ABI is the input ABI used to generate the binding from.
const VRFConsumerV2ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// VRFConsumerV2Bin is the compiled bytecode used for deploying new contracts.
var VRFConsumerV2Bin = "0x60a060405234801561001057600080fd5b50604051610b1a380380610b1a83398101604081905261002f91610082565b6001600160a01b039182166080819052600280546001600160a01b03199081169092179055600380549290931691161790556100b5565b80516001600160a01b038116811461007d57600080fd5b919050565b6000806040838503121561009557600080fd5b61009e83610066565b91506100ac60208401610066565b90509250929050565b608051610a4d6100cd60003960005050610a4d6000f3fe608060405234801561001057600080fd5b50600436106100835760003560e01c80631fe543e31461008857806327784fad1461009d5780632fa4e442146100c357806336bfffed146100d65780636802f726146100e9578063706da1ca146100fc578063e89e106a1461012e578063f08c5daa14610137578063f6eaffc814610140575b600080fd5b61009b6100963660046106c2565b610153565b005b6100b06100ab366004610794565b610161565b6040519081526020015b60405180910390f35b61009b6100d13660046107fd565b61020d565b61009b6100e436600461082d565b610309565b61009b6100f73660046107fd565b61040e565b60035461011690600160a01b90046001600160401b031681565b6040516001600160401b0390911681526020016100ba565b6100b060015481565b6100b060045481565b6100b061014e3660046108d3565b610571565b61015d8282610592565b5050565b6002546040516305d3b1d360e41b8152600481018790526001600160401b038616602482015261ffff8516604482015263ffffffff8085166064830152831660848201526000916001600160a01b031690635d3b1d309060a4016020604051808303816000875af11580156101da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101fe91906108ec565b60018190559695505050505050565b600354600160a01b90046001600160401b03166000036102625760405162461bcd60e51b815260206004820152600b60248201526a1cdd58881b9bdd081cd95d60aa1b60448201526064015b60405180910390fd5b60035460025460408051600160a01b84046001600160401b031660208201526001600160a01b0393841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b81526004016102c693929190610905565b6020604051808303816000875af11580156102e5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061015d9190610979565b600354600160a01b90046001600160401b031660000361035b5760405162461bcd60e51b815260206004820152600d60248201526c1cdd589251081b9bdd081cd95d609a1b6044820152606401610259565b60005b815181101561015d5760025460035483516001600160a01b0390921691637341c10c91600160a01b90046001600160401b0316908590859081106103a4576103a461099b565b60200260200101516040518363ffffffff1660e01b81526004016103c99291906109b1565b600060405180830381600087803b1580156103e357600080fd5b505af11580156103f7573d6000803e3d6000fd5b505050508080610406906109d3565b91505061035e565b600354600160a01b90046001600160401b031660000361026257600260009054906101000a90046001600160a01b03166001600160a01b031663a21a23e46040518163ffffffff1660e01b81526004016020604051808303816000875af115801561047d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104a191906109fa565b6003805467ffffffffffffffff60a01b1916600160a01b6001600160401b0393841681029190911791829055600254604051631cd0704360e21b81526001600160a01b0390911693637341c10c93610501939004169030906004016109b1565b600060405180830381600087803b15801561051b57600080fd5b505af115801561052f573d6000803e3d6000fd5b505060035460025460408051600160a01b84046001600160401b031660208201526001600160a01b039384169550634000aea094509290911691859101610299565b6000818154811061058157600080fd5b600091825260209091200154905081565b60015482146105dd5760405162461bcd60e51b81526020600482015260176024820152761c995c5d595cdd081251081a5cc81a5b98dbdc9c9958dd604a1b6044820152606401610259565b5a60045580516105f49060009060208401906105f9565b505050565b828054828255906000526020600020908101928215610634579160200282015b82811115610634578251825591602001919060010190610619565b50610640929150610644565b5090565b5b808211156106405760008155600101610645565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171561069757610697610659565b604052919050565b60006001600160401b038211156106b8576106b8610659565b5060051b60200190565b600080604083850312156106d557600080fd5b823591506020808401356001600160401b038111156106f357600080fd5b8401601f8101861361070457600080fd5b80356107176107128261069f565b61066f565b81815260059190911b8201830190838101908883111561073657600080fd5b928401925b828410156107545783358252928401929084019061073b565b80955050505050509250929050565b6001600160401b038116811461077857600080fd5b50565b803563ffffffff8116811461078f57600080fd5b919050565b600080600080600060a086880312156107ac57600080fd5b8535945060208601356107be81610763565b9350604086013561ffff811681146107d557600080fd5b92506107e36060870161077b565b91506107f16080870161077b565b90509295509295909350565b60006020828403121561080f57600080fd5b81356001600160601b038116811461082657600080fd5b9392505050565b6000602080838503121561084057600080fd5b82356001600160401b0381111561085657600080fd5b8301601f8101851361086757600080fd5b80356108756107128261069f565b81815260059190911b8201830190838101908783111561089457600080fd5b928401925b828410156108c85783356001600160a01b03811681146108b95760008081fd5b82529284019290840190610899565b979650505050505050565b6000602082840312156108e557600080fd5b5035919050565b6000602082840312156108fe57600080fd5b5051919050565b60018060a01b03841681526000602060018060601b0385168184015260606040840152835180606085015260005b8181101561094f57858101830151858201608001528201610933565b81811115610961576000608083870101525b50601f01601f19169290920160800195945050505050565b60006020828403121561098b57600080fd5b8151801515811461082657600080fd5b634e487b7160e01b600052603260045260246000fd5b6001600160401b039290921682526001600160a01b0316602082015260400190565b6000600182016109f357634e487b7160e01b600052601160045260246000fd5b5060010190565b600060208284031215610a0c57600080fd5b81516108268161076356fea2646970667358221220440f61541a58900e2bd56c661fe15f704d34b703df3331d0429228a7f95dcadd64736f6c634300080d0033"

// DeployVRFConsumerV2 deploys a new Ethereum contract, binding an instance of VRFConsumerV2 to it.
func DeployVRFConsumerV2(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFConsumerV2, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFConsumerV2Bin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerV2{VRFConsumerV2Caller: VRFConsumerV2Caller{contract: contract}, VRFConsumerV2Transactor: VRFConsumerV2Transactor{contract: contract}, VRFConsumerV2Filterer: VRFConsumerV2Filterer{contract: contract}}, nil
}

// VRFConsumerV2 is an auto generated Go binding around an Ethereum contract.
type VRFConsumerV2 struct {
	VRFConsumerV2Caller     // Read-only binding to the contract
	VRFConsumerV2Transactor // Write-only binding to the contract
	VRFConsumerV2Filterer   // Log filterer for contract events
}

// VRFConsumerV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type VRFConsumerV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFConsumerV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFConsumerV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFConsumerV2Session struct {
	Contract     *VRFConsumerV2    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFConsumerV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFConsumerV2CallerSession struct {
	Contract *VRFConsumerV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VRFConsumerV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFConsumerV2TransactorSession struct {
	Contract     *VRFConsumerV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VRFConsumerV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type VRFConsumerV2Raw struct {
	Contract *VRFConsumerV2 // Generic contract binding to access the raw methods on
}

// VRFConsumerV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFConsumerV2CallerRaw struct {
	Contract *VRFConsumerV2Caller // Generic read-only contract binding to access the raw methods on
}

// VRFConsumerV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFConsumerV2TransactorRaw struct {
	Contract *VRFConsumerV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFConsumerV2 creates a new instance of VRFConsumerV2, bound to a specific deployed contract.
func NewVRFConsumerV2(address common.Address, backend bind.ContractBackend) (*VRFConsumerV2, error) {
	contract, err := bindVRFConsumerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2{VRFConsumerV2Caller: VRFConsumerV2Caller{contract: contract}, VRFConsumerV2Transactor: VRFConsumerV2Transactor{contract: contract}, VRFConsumerV2Filterer: VRFConsumerV2Filterer{contract: contract}}, nil
}

// NewVRFConsumerV2Caller creates a new read-only instance of VRFConsumerV2, bound to a specific deployed contract.
func NewVRFConsumerV2Caller(address common.Address, caller bind.ContractCaller) (*VRFConsumerV2Caller, error) {
	contract, err := bindVRFConsumerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Caller{contract: contract}, nil
}

// NewVRFConsumerV2Transactor creates a new write-only instance of VRFConsumerV2, bound to a specific deployed contract.
func NewVRFConsumerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerV2Transactor, error) {
	contract, err := bindVRFConsumerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Transactor{contract: contract}, nil
}

// NewVRFConsumerV2Filterer creates a new log filterer instance of VRFConsumerV2, bound to a specific deployed contract.
func NewVRFConsumerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerV2Filterer, error) {
	contract, err := bindVRFConsumerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Filterer{contract: contract}, nil
}

// bindVRFConsumerV2 binds a generic wrapper to an already deployed contract.
func bindVRFConsumerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumerV2 *VRFConsumerV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2.Contract.VRFConsumerV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumerV2 *VRFConsumerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.VRFConsumerV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumerV2 *VRFConsumerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.VRFConsumerV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumerV2 *VRFConsumerV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumerV2 *VRFConsumerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumerV2 *VRFConsumerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.contract.Transact(opts, method, params...)
}

// SGasAvailable is a free data retrieval call binding the contract method 0xf08c5daa.
//
// Solidity: function s_gasAvailable() view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Caller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SGasAvailable is a free data retrieval call binding the contract method 0xf08c5daa.
//
// Solidity: function s_gasAvailable() view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Session) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SGasAvailable(&_VRFConsumerV2.CallOpts)
}

// SGasAvailable is a free data retrieval call binding the contract method 0xf08c5daa.
//
// Solidity: function s_gasAvailable() view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SGasAvailable(&_VRFConsumerV2.CallOpts)
}

// SRandomWords is a free data retrieval call binding the contract method 0xf6eaffc8.
//
// Solidity: function s_randomWords(uint256 ) view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Caller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRandomWords is a free data retrieval call binding the contract method 0xf6eaffc8.
//
// Solidity: function s_randomWords(uint256 ) view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Session) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

// SRandomWords is a free data retrieval call binding the contract method 0xf6eaffc8.
//
// Solidity: function s_randomWords(uint256 ) view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

// SRequestId is a free data retrieval call binding the contract method 0xe89e106a.
//
// Solidity: function s_requestId() view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Caller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SRequestId is a free data retrieval call binding the contract method 0xe89e106a.
//
// Solidity: function s_requestId() view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Session) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRequestId(&_VRFConsumerV2.CallOpts)
}

// SRequestId is a free data retrieval call binding the contract method 0xe89e106a.
//
// Solidity: function s_requestId() view returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRequestId(&_VRFConsumerV2.CallOpts)
}

// SSubId is a free data retrieval call binding the contract method 0x706da1ca.
//
// Solidity: function s_subId() view returns(uint64)
func (_VRFConsumerV2 *VRFConsumerV2Caller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SSubId is a free data retrieval call binding the contract method 0x706da1ca.
//
// Solidity: function s_subId() view returns(uint64)
func (_VRFConsumerV2 *VRFConsumerV2Session) SSubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SSubId(&_VRFConsumerV2.CallOpts)
}

// SSubId is a free data retrieval call binding the contract method 0x706da1ca.
//
// Solidity: function s_subId() view returns(uint64)
func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SSubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SSubId(&_VRFConsumerV2.CallOpts)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFConsumerV2 *VRFConsumerV2Transactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFConsumerV2 *VRFConsumerV2Session) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.RawFulfillRandomWords(&_VRFConsumerV2.TransactOpts, requestId, randomWords)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.RawFulfillRandomWords(&_VRFConsumerV2.TransactOpts, requestId, randomWords)
}

// TestCreateSubscriptionAndFund is a paid mutator transaction binding the contract method 0x6802f726.
//
// Solidity: function testCreateSubscriptionAndFund(uint96 amount) returns()
func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testCreateSubscriptionAndFund", amount)
}

// TestCreateSubscriptionAndFund is a paid mutator transaction binding the contract method 0x6802f726.
//
// Solidity: function testCreateSubscriptionAndFund(uint96 amount) returns()
func (_VRFConsumerV2 *VRFConsumerV2Session) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2.TransactOpts, amount)
}

// TestCreateSubscriptionAndFund is a paid mutator transaction binding the contract method 0x6802f726.
//
// Solidity: function testCreateSubscriptionAndFund(uint96 amount) returns()
func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2.TransactOpts, amount)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x27784fad.
//
// Solidity: function testRequestRandomness(bytes32 keyHash, uint64 subId, uint16 minReqConfs, uint32 callbackGasLimit, uint32 numWords) returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testRequestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x27784fad.
//
// Solidity: function testRequestRandomness(bytes32 keyHash, uint64 subId, uint16 minReqConfs, uint32 callbackGasLimit, uint32 numWords) returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2Session) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x27784fad.
//
// Solidity: function testRequestRandomness(bytes32 keyHash, uint64 subId, uint16 minReqConfs, uint32 callbackGasLimit, uint32 numWords) returns(uint256)
func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

// TopUpSubscription is a paid mutator transaction binding the contract method 0x2fa4e442.
//
// Solidity: function topUpSubscription(uint96 amount) returns()
func (_VRFConsumerV2 *VRFConsumerV2Transactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "topUpSubscription", amount)
}

// TopUpSubscription is a paid mutator transaction binding the contract method 0x2fa4e442.
//
// Solidity: function topUpSubscription(uint96 amount) returns()
func (_VRFConsumerV2 *VRFConsumerV2Session) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TopUpSubscription(&_VRFConsumerV2.TransactOpts, amount)
}

// TopUpSubscription is a paid mutator transaction binding the contract method 0x2fa4e442.
//
// Solidity: function topUpSubscription(uint96 amount) returns()
func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TopUpSubscription(&_VRFConsumerV2.TransactOpts, amount)
}

// UpdateSubscription is a paid mutator transaction binding the contract method 0x36bfffed.
//
// Solidity: function updateSubscription(address[] consumers) returns()
func (_VRFConsumerV2 *VRFConsumerV2Transactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "updateSubscription", consumers)
}

// UpdateSubscription is a paid mutator transaction binding the contract method 0x36bfffed.
//
// Solidity: function updateSubscription(address[] consumers) returns()
func (_VRFConsumerV2 *VRFConsumerV2Session) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.UpdateSubscription(&_VRFConsumerV2.TransactOpts, consumers)
}

// UpdateSubscription is a paid mutator transaction binding the contract method 0x36bfffed.
//
// Solidity: function updateSubscription(address[] consumers) returns()
func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.UpdateSubscription(&_VRFConsumerV2.TransactOpts, consumers)
}
