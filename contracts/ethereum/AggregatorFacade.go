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

// AggregatorFacadeABI is the input ABI used to generate the binding from.
const AggregatorFacadeABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"aggregator\",\"outputs\":[{\"internalType\":\"contractAggregatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AggregatorFacadeBin is the compiled bytecode used for deploying new contracts.
var AggregatorFacadeBin = "0x608060405234801561001057600080fd5b5060405162000db238038062000db28339818101604052606081101561003557600080fd5b8101908080519060200190929190805190602001909291908051604051939291908464010000000082111561006957600080fd5b8382019150602082018581111561007f57600080fd5b825186600182028301116401000000008211171561009c57600080fd5b8083526020830192505050908051906020019080838360005b838110156100d05780820151818401526020810190506100b5565b50505050905090810190601f1680156100fd5780820380516001836020036101000a031916815260200191505b50604052505050826000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600060146101000a81548160ff021916908360ff160217905550806001908051906020019061017592919061017e565b50505050610223565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101bf57805160ff19168380011785556101ed565b828001600101855582156101ed579182015b828111156101ec5782518255916020019190600101906101d1565b5b5090506101fa91906101fe565b5090565b61022091905b8082111561021c576000816000905550600101610204565b5090565b90565b610b7f80620002336000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80637284e416116100715780637284e416146101765780638205bf6a146101f95780639a6fc8f514610217578063b5ab58dc146102b1578063b633620c146102f3578063feaf968c14610335576100a9565b8063245a7bfc146100ae578063313ce567146100f857806350d25bcd1461011c57806354fd4d501461013a578063668a0f0214610158575b600080fd5b6100b661039f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101006103c4565b604051808260ff1660ff16815260200191505060405180910390f35b6101246103d7565b6040518082815260200191505060405180910390f35b610142610480565b6040518082815260200191505060405180910390f35b610160610485565b6040518082815260200191505060405180910390f35b61017e61052e565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101be5780820151818401526020810190506101a3565b50505050905090810190601f1680156101eb5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102016105cc565b6040518082815260200191505060405180910390f35b61024f6004803603602081101561022d57600080fd5b81019080803569ffffffffffffffffffff169060200190929190505050610675565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102dd600480360360208110156102c757600080fd5b8101908080359060200190929190505050610699565b6040518082815260200191505060405180910390f35b61031f6004803603602081101561030957600080fd5b810190808035906020019092919050505061074f565b6040518082815260200191505060405180910390f35b61033d610805565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600060149054906101000a900460ff1681565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b815260040160206040518083038186803b15801561044057600080fd5b505afa158015610454573d6000803e3d6000fd5b505050506040513d602081101561046a57600080fd5b8101908080519060200190929190505050905090565b600281565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663668a0f026040518163ffffffff1660e01b815260040160206040518083038186803b1580156104ee57600080fd5b505afa158015610502573d6000803e3d6000fd5b505050506040513d602081101561051857600080fd5b8101908080519060200190929190505050905090565b60018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105c45780601f10610599576101008083540402835291602001916105c4565b820191906000526020600020905b8154815290600101906020018083116105a757829003601f168201915b505050505081565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638205bf6a6040518163ffffffff1660e01b815260040160206040518083038186803b15801561063557600080fd5b505afa158015610649573d6000803e3d6000fd5b505050506040513d602081101561065f57600080fd5b8101908080519060200190929190505050905090565b6000806000806000610686866108c8565b9450945094509450945091939590929450565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b5ab58dc836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561070d57600080fd5b505afa158015610721573d6000803e3d6000fd5b505050506040513d602081101561073757600080fd5b81019080805190602001909291905050509050919050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b633620c836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b1580156107c357600080fd5b505afa1580156107d7573d6000803e3d6000fd5b505050506040513d60208110156107ed57600080fd5b81019080805190602001909291905050509050919050565b60008060008060006108b76000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663668a0f026040518163ffffffff1660e01b815260040160206040518083038186803b15801561087757600080fd5b505afa15801561088b573d6000803e3d6000fd5b505050506040513d60208110156108a157600080fd5b81019080805190602001909291905050506108c8565b945094509450945094509091929394565b60008060008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b5ab58dc876040518263ffffffff1660e01b8152600401808269ffffffffffffffffffff16815260200191505060206040518083038186803b15801561094e57600080fd5b505afa158015610962573d6000803e3d6000fd5b505050506040513d602081101561097857600080fd5b810190808051906020019092919050505093506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b633620c876040518263ffffffff1660e01b8152600401808269ffffffffffffffffffff16815260200191505060206040518083038186803b158015610a0957600080fd5b505afa158015610a1d573d6000803e3d6000fd5b505050506040513d6020811015610a3357600080fd5b810190808051906020019092919050505067ffffffffffffffff169150600082116040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090610b30576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610af5578082015181840152602081019050610ada565b50505050905090810190601f168015610b225780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b508584838489945094509450945094509193959092945056fea2646970667358221220a822653af921158ca9b15557d4814111b34541df93bc60601648a3d086980fa064736f6c63430006060033"

// DeployAggregatorFacade deploys a new Ethereum contract, binding an instance of AggregatorFacade to it.
func DeployAggregatorFacade(auth *bind.TransactOpts, backend bind.ContractBackend, _aggregator common.Address, _decimals uint8, _description string) (common.Address, *types.Transaction, *AggregatorFacade, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorFacadeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AggregatorFacadeBin), backend, _aggregator, _decimals, _description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AggregatorFacade{AggregatorFacadeCaller: AggregatorFacadeCaller{contract: contract}, AggregatorFacadeTransactor: AggregatorFacadeTransactor{contract: contract}, AggregatorFacadeFilterer: AggregatorFacadeFilterer{contract: contract}}, nil
}

// AggregatorFacade is an auto generated Go binding around an Ethereum contract.
type AggregatorFacade struct {
	AggregatorFacadeCaller     // Read-only binding to the contract
	AggregatorFacadeTransactor // Write-only binding to the contract
	AggregatorFacadeFilterer   // Log filterer for contract events
}

// AggregatorFacadeCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorFacadeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorFacadeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorFacadeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorFacadeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorFacadeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorFacadeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorFacadeSession struct {
	Contract     *AggregatorFacade // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AggregatorFacadeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorFacadeCallerSession struct {
	Contract *AggregatorFacadeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// AggregatorFacadeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorFacadeTransactorSession struct {
	Contract     *AggregatorFacadeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// AggregatorFacadeRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorFacadeRaw struct {
	Contract *AggregatorFacade // Generic contract binding to access the raw methods on
}

// AggregatorFacadeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorFacadeCallerRaw struct {
	Contract *AggregatorFacadeCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorFacadeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorFacadeTransactorRaw struct {
	Contract *AggregatorFacadeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorFacade creates a new instance of AggregatorFacade, bound to a specific deployed contract.
func NewAggregatorFacade(address common.Address, backend bind.ContractBackend) (*AggregatorFacade, error) {
	contract, err := bindAggregatorFacade(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorFacade{AggregatorFacadeCaller: AggregatorFacadeCaller{contract: contract}, AggregatorFacadeTransactor: AggregatorFacadeTransactor{contract: contract}, AggregatorFacadeFilterer: AggregatorFacadeFilterer{contract: contract}}, nil
}

// NewAggregatorFacadeCaller creates a new read-only instance of AggregatorFacade, bound to a specific deployed contract.
func NewAggregatorFacadeCaller(address common.Address, caller bind.ContractCaller) (*AggregatorFacadeCaller, error) {
	contract, err := bindAggregatorFacade(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorFacadeCaller{contract: contract}, nil
}

// NewAggregatorFacadeTransactor creates a new write-only instance of AggregatorFacade, bound to a specific deployed contract.
func NewAggregatorFacadeTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorFacadeTransactor, error) {
	contract, err := bindAggregatorFacade(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorFacadeTransactor{contract: contract}, nil
}

// NewAggregatorFacadeFilterer creates a new log filterer instance of AggregatorFacade, bound to a specific deployed contract.
func NewAggregatorFacadeFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorFacadeFilterer, error) {
	contract, err := bindAggregatorFacade(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorFacadeFilterer{contract: contract}, nil
}

// bindAggregatorFacade binds a generic wrapper to an already deployed contract.
func bindAggregatorFacade(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorFacadeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorFacade *AggregatorFacadeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorFacade.Contract.AggregatorFacadeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorFacade *AggregatorFacadeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorFacade.Contract.AggregatorFacadeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorFacade *AggregatorFacadeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorFacade.Contract.AggregatorFacadeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorFacade *AggregatorFacadeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorFacade.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorFacade *AggregatorFacadeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorFacade.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorFacade *AggregatorFacadeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorFacade.Contract.contract.Transact(opts, method, params...)
}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_AggregatorFacade *AggregatorFacadeCaller) Aggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "aggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_AggregatorFacade *AggregatorFacadeSession) Aggregator() (common.Address, error) {
	return _AggregatorFacade.Contract.Aggregator(&_AggregatorFacade.CallOpts)
}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_AggregatorFacade *AggregatorFacadeCallerSession) Aggregator() (common.Address, error) {
	return _AggregatorFacade.Contract.Aggregator(&_AggregatorFacade.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorFacade *AggregatorFacadeCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorFacade *AggregatorFacadeSession) Decimals() (uint8, error) {
	return _AggregatorFacade.Contract.Decimals(&_AggregatorFacade.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorFacade *AggregatorFacadeCallerSession) Decimals() (uint8, error) {
	return _AggregatorFacade.Contract.Decimals(&_AggregatorFacade.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorFacade *AggregatorFacadeCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorFacade *AggregatorFacadeSession) Description() (string, error) {
	return _AggregatorFacade.Contract.Description(&_AggregatorFacade.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorFacade *AggregatorFacadeCallerSession) Description() (string, error) {
	return _AggregatorFacade.Contract.Description(&_AggregatorFacade.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AggregatorFacade *AggregatorFacadeCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AggregatorFacade *AggregatorFacadeSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorFacade.Contract.GetAnswer(&_AggregatorFacade.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AggregatorFacade *AggregatorFacadeCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorFacade.Contract.GetAnswer(&_AggregatorFacade.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorFacade *AggregatorFacadeCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorFacade *AggregatorFacadeSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorFacade.Contract.GetRoundData(&_AggregatorFacade.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorFacade *AggregatorFacadeCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorFacade.Contract.GetRoundData(&_AggregatorFacade.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorFacade.Contract.GetTimestamp(&_AggregatorFacade.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorFacade.Contract.GetTimestamp(&_AggregatorFacade.CallOpts, _roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorFacade *AggregatorFacadeCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorFacade *AggregatorFacadeSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorFacade.Contract.LatestAnswer(&_AggregatorFacade.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorFacade *AggregatorFacadeCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorFacade.Contract.LatestAnswer(&_AggregatorFacade.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeSession) LatestRound() (*big.Int, error) {
	return _AggregatorFacade.Contract.LatestRound(&_AggregatorFacade.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorFacade.Contract.LatestRound(&_AggregatorFacade.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorFacade *AggregatorFacadeCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorFacade *AggregatorFacadeSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorFacade.Contract.LatestRoundData(&_AggregatorFacade.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorFacade *AggregatorFacadeCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorFacade.Contract.LatestRoundData(&_AggregatorFacade.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorFacade.Contract.LatestTimestamp(&_AggregatorFacade.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorFacade.Contract.LatestTimestamp(&_AggregatorFacade.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorFacade.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeSession) Version() (*big.Int, error) {
	return _AggregatorFacade.Contract.Version(&_AggregatorFacade.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorFacade *AggregatorFacadeCallerSession) Version() (*big.Int, error) {
	return _AggregatorFacade.Contract.Version(&_AggregatorFacade.CallOpts)
}

// AggregatorFacadeAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AggregatorFacade contract.
type AggregatorFacadeAnswerUpdatedIterator struct {
	Event *AggregatorFacadeAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *AggregatorFacadeAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorFacadeAnswerUpdated)
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
		it.Event = new(AggregatorFacadeAnswerUpdated)
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
func (it *AggregatorFacadeAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorFacadeAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorFacadeAnswerUpdated represents a AnswerUpdated event raised by the AggregatorFacade contract.
type AggregatorFacadeAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorFacade *AggregatorFacadeFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorFacadeAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorFacade.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorFacadeAnswerUpdatedIterator{contract: _AggregatorFacade.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorFacade *AggregatorFacadeFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorFacadeAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorFacade.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorFacadeAnswerUpdated)
				if err := _AggregatorFacade.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorFacade *AggregatorFacadeFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorFacadeAnswerUpdated, error) {
	event := new(AggregatorFacadeAnswerUpdated)
	if err := _AggregatorFacade.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorFacadeNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AggregatorFacade contract.
type AggregatorFacadeNewRoundIterator struct {
	Event *AggregatorFacadeNewRound // Event containing the contract specifics and raw log

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
func (it *AggregatorFacadeNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorFacadeNewRound)
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
		it.Event = new(AggregatorFacadeNewRound)
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
func (it *AggregatorFacadeNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorFacadeNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorFacadeNewRound represents a NewRound event raised by the AggregatorFacade contract.
type AggregatorFacadeNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorFacade *AggregatorFacadeFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorFacadeNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorFacade.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorFacadeNewRoundIterator{contract: _AggregatorFacade.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorFacade *AggregatorFacadeFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorFacadeNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorFacade.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorFacadeNewRound)
				if err := _AggregatorFacade.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorFacade *AggregatorFacadeFilterer) ParseNewRound(log types.Log) (*AggregatorFacadeNewRound, error) {
	event := new(AggregatorFacadeNewRound)
	if err := _AggregatorFacade.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
