// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package celo

import (
	"math/big"
	"strings"

	celo "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/event"
	"github.com/smartcontractkit/integrations-framework/celoextended"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = celo.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OracleABI is the input ABI used to generate the binding from.
const OracleABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CancelOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cancelExpiration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXPIRY_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunc\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_data\",\"type\":\"bytes32\"}],\"name\":\"fulfillOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node\",\"type\":\"address\"}],\"name\":\"getAuthorizationStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_dataVersion\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"oracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_allowed\",\"type\":\"bool\"}],\"name\":\"setFulfillmentPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OracleBin is the compiled bytecode used for deploying new contracts.
var OracleBin = "0x6080604052600160045534801561001557600080fd5b5060405161130d38038061130d8339818101604052602081101561003857600080fd5b5051600080546001600160a01b03191633178082556040516001600160a01b039190911691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3600180546001600160a01b0319166001600160a01b039290921691909117905561125b806100b26000396000f3fe608060405234801561001057600080fd5b50600436106100bf5760003560e01c80637fcd56db1161007c5780637fcd56db146102565780638da5cb5b146102845780638f32d59b1461028c578063a4c0ed3614610294578063d3e9c3141461034d578063f2fde38b14610373578063f3fef3a314610399576100bf565b8063165d35e1146100c457806340429946146100e85780634ab0d190146101995780634b602282146101fb57806350188301146102155780636ee4d5531461021d575b600080fd5b6100cc6103c5565b604080516001600160a01b039092168252519081900360200190f35b61019760048036036101008110156100ff57600080fd5b6001600160a01b038235811692602081013592604082013592606083013516916001600160e01b03196080820135169160a08201359160c081013591810190610100810160e0820135600160201b81111561015957600080fd5b82018360208201111561016b57600080fd5b803590602001918460018302840111600160201b8311171561018c57600080fd5b5090925090506103d4565b005b6101e7600480360360c08110156101af57600080fd5b508035906020810135906001600160a01b03604082013516906001600160e01b03196060820135169060808101359060a0013561069f565b604080519115158252519081900360200190f35b610203610994565b60408051918252519081900360200190f35b61020361099a565b6101976004803603608081101561023357600080fd5b508035906020810135906001600160e01b031960408201351690606001356109fc565b6101976004803603604081101561026c57600080fd5b506001600160a01b0381351690602001351515610bb6565b6100cc610c28565b6101e7610c37565b610197600480360360608110156102aa57600080fd5b6001600160a01b0382351691602081013591810190606081016040820135600160201b8111156102d957600080fd5b8201836020820111156102eb57600080fd5b803590602001918460018302840111600160201b8311171561030c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610c48945050505050565b6101e76004803603602081101561036357600080fd5b50356001600160a01b0316610e70565b6101976004803603602081101561038957600080fd5b50356001600160a01b0316610e8e565b610197600480360360408110156103af57600080fd5b506001600160a01b038135169060200135610ee1565b6001546001600160a01b031690565b6103dc6103c5565b6001600160a01b0316336001600160a01b031614610437576040805162461bcd60e51b815260206004820152601360248201527226bab9ba103ab9b2902624a725903a37b5b2b760691b604482015290519081900360640190fd5b60015486906001600160a01b0380831691161415610496576040805162461bcd60e51b815260206004820152601760248201527643616e6e6f742063616c6c6261636b20746f204c494e4b60481b604482015290519081900360640190fd5b604080516001600160601b031960608d901b166020808301919091526034808301899052835180840390910181526054909201835281519181019190912060008181526002909252919020541561052b576040805162461bcd60e51b8152602060048201526014602482015273135d5cdd081d5cd94818481d5b9a5c5d5948125160621b604482015290519081900360640190fd5b600061053f4261012c63ffffffff61102216565b90508a89898360405160200180858152602001846001600160a01b03166001600160a01b031660601b8152601401836001600160e01b0319166001600160e01b0319168152600401828152602001945050505050604051602081830303815290604052805190602001206002600084815260200190815260200160002081905550897fd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c658d848e8d8d878d8d8d604051808a6001600160a01b03166001600160a01b03168152602001898152602001888152602001876001600160a01b03166001600160a01b03168152602001866001600160e01b0319166001600160e01b0319168152602001858152602001848152602001806020018281038252848482818152602001925080828437600083820152604051601f909101601f19169092018290039c50909a5050505050505050505050a2505050505050505050505050565b3360009081526003602052604081205460ff16806106d557506106c0610c28565b6001600160a01b0316336001600160a01b0316145b6107105760405162461bcd60e51b815260040180806020018281038252602a8152602001806111dc602a913960400191505060405180910390fd5b6000878152600260205260409020548790610772576040805162461bcd60e51b815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b6040805160208082018a90526001600160601b031960608a901b16828401526001600160e01b0319881660548301526058808301889052835180840390910181526078909201835281519181019190912060008b81526002909252919020548114610824576040805162461bcd60e51b815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b600454610837908963ffffffff61102216565b60045560008981526002602052604081205562061a805a10156108a1576040805162461bcd60e51b815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60408051602481018b9052604480820187905282518083039091018152606490910182526020810180516001600160e01b03166001600160e01b03198a16178152915181516000936001600160a01b038c169392918291908083835b6020831061091c5780518252601f1990920191602091820191016108fd565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461097e576040519150601f19603f3d011682016040523d82523d6000602084013e610983565b606091505b50909b9a5050505050505050505050565b61012c81565b60006109a4610c37565b6109e3576040805162461bcd60e51b81526020600482018190526024820152600080516020611206833981519152604482015290519081900360640190fd5b6004546109f790600163ffffffff61108316565b905090565b6040805160208082018690523360601b828401526001600160e01b0319851660548301526058808301859052835180840390910181526078909201835281519181019190912060008781526002909252919020548114610aa3576040805162461bcd60e51b815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b42821115610af1576040805162461bcd60e51b815260206004820152601660248201527514995c5d595cdd081a5cc81b9bdd08195e1c1a5c995960521b604482015290519081900360640190fd5b6000858152600260205260408082208290555186917fa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e9391a26001546040805163a9059cbb60e01b81523360048201526024810187905290516001600160a01b039092169163a9059cbb916044808201926020929091908290030181600087803b158015610b7d57600080fd5b505af1158015610b91573d6000803e3d6000fd5b505050506040513d6020811015610ba757600080fd5b5051610baf57fe5b5050505050565b610bbe610c37565b610bfd576040805162461bcd60e51b81526020600482018190526024820152600080516020611206833981519152604482015290519081900360640190fd5b6001600160a01b03919091166000908152600360205260409020805460ff1916911515919091179055565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b610c506103c5565b6001600160a01b0316336001600160a01b031614610cab576040805162461bcd60e51b815260206004820152601360248201527226bab9ba103ab9b2902624a725903a37b5b2b760691b604482015290519081900360640190fd5b8051819060441115610cfd576040805162461bcd60e51b8152602060048201526016602482015275092dcecc2d8d2c840e4cae2eacae6e840d8cadccee8d60531b604482015290519081900360640190fd5b602082015182906001600160e01b031981166320214ca360e11b14610d69576040805162461bcd60e51b815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b8560248501528460448501526000306001600160a01b0316856040518082805190602001908083835b60208310610db15780518252601f199092019160209182019101610d92565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d8060008114610e11576040519150601f19603f3d011682016040523d82523d6000602084013e610e16565b606091505b5050905080610e67576040805162461bcd60e51b8152602060048201526018602482015277155b98589b19481d1bc818dc99585d19481c995c5d595cdd60421b604482015290519081900360640190fd5b50505050505050565b6001600160a01b031660009081526003602052604090205460ff1690565b610e96610c37565b610ed5576040805162461bcd60e51b81526020600482018190526024820152600080516020611206833981519152604482015290519081900360640190fd5b610ede816110e0565b50565b610ee9610c37565b610f28576040805162461bcd60e51b81526020600482018190526024820152600080516020611206833981519152604482015290519081900360640190fd5b80610f3a81600163ffffffff61102216565b6004541015610f7a5760405162461bcd60e51b81526004018080602001828103825260358152602001806111a76035913960400191505060405180910390fd5b600454610f8d908363ffffffff61108316565b60049081556001546040805163a9059cbb60e01b81526001600160a01b0387811694820194909452602481018690529051929091169163a9059cbb916044808201926020929091908290030181600087803b158015610feb57600080fd5b505af1158015610fff573d6000803e3d6000fd5b505050506040513d602081101561101557600080fd5b505161101d57fe5b505050565b60008282018381101561107c576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6000828211156110da576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6001600160a01b0381166111255760405162461bcd60e51b81526004018080602001828103825260268152602001806111816026913960400191505060405180910390fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b039290921691909117905556fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373416d6f756e74207265717565737465642069732067726561746572207468616e20776974686472617761626c652062616c616e63654e6f7420616e20617574686f72697a6564206e6f646520746f2066756c66696c6c2072657175657374734f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572a2646970667358221220c2ac70a484a973f724797613a8dae84bdcc9f03fb39700ee0f65d6564a4cc56264736f6c63430006060033"

// DeployOracle deploys a new Ethereum contract, binding an instance of Oracle to it.
func DeployOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address) (common.Address, *types.Transaction, *Oracle, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OracleBin), backend, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// EXPIRYTIME is a free data retrieval call binding the contract method 0x4b602282.
//
// Solidity: function EXPIRY_TIME() view returns(uint256)
func (_Oracle *OracleCaller) EXPIRYTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "EXPIRY_TIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *celoextended.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EXPIRYTIME is a free data retrieval call binding the contract method 0x4b602282.
//
// Solidity: function EXPIRY_TIME() view returns(uint256)
func (_Oracle *OracleSession) EXPIRYTIME() (*big.Int, error) {
	return _Oracle.Contract.EXPIRYTIME(&_Oracle.CallOpts)
}

// EXPIRYTIME is a free data retrieval call binding the contract method 0x4b602282.
//
// Solidity: function EXPIRY_TIME() view returns(uint256)
func (_Oracle *OracleCallerSession) EXPIRYTIME() (*big.Int, error) {
	return _Oracle.Contract.EXPIRYTIME(&_Oracle.CallOpts)
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address _node) view returns(bool)
func (_Oracle *OracleCaller) GetAuthorizationStatus(opts *bind.CallOpts, _node common.Address) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getAuthorizationStatus", _node)

	if err != nil {
		return *new(bool), err
	}

	out0 := *celoextended.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address _node) view returns(bool)
func (_Oracle *OracleSession) GetAuthorizationStatus(_node common.Address) (bool, error) {
	return _Oracle.Contract.GetAuthorizationStatus(&_Oracle.CallOpts, _node)
}

// GetAuthorizationStatus is a free data retrieval call binding the contract method 0xd3e9c314.
//
// Solidity: function getAuthorizationStatus(address _node) view returns(bool)
func (_Oracle *OracleCallerSession) GetAuthorizationStatus(_node common.Address) (bool, error) {
	return _Oracle.Contract.GetAuthorizationStatus(&_Oracle.CallOpts, _node)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Oracle *OracleCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *celoextended.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Oracle *OracleSession) GetChainlinkToken() (common.Address, error) {
	return _Oracle.Contract.GetChainlinkToken(&_Oracle.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Oracle *OracleCallerSession) GetChainlinkToken() (common.Address, error) {
	return _Oracle.Contract.GetChainlinkToken(&_Oracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Oracle *OracleCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *celoextended.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Oracle *OracleSession) IsOwner() (bool, error) {
	return _Oracle.Contract.IsOwner(&_Oracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Oracle *OracleCallerSession) IsOwner() (bool, error) {
	return _Oracle.Contract.IsOwner(&_Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *celoextended.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleCallerSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_Oracle *OracleCaller) Withdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "withdrawable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *celoextended.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_Oracle *OracleSession) Withdrawable() (*big.Int, error) {
	return _Oracle.Contract.Withdrawable(&_Oracle.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0x50188301.
//
// Solidity: function withdrawable() view returns(uint256)
func (_Oracle *OracleCallerSession) Withdrawable() (*big.Int, error) {
	return _Oracle.Contract.Withdrawable(&_Oracle.CallOpts)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunc, uint256 _expiration) returns()
func (_Oracle *OracleTransactor) CancelOracleRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "cancelOracleRequest", _requestId, _payment, _callbackFunc, _expiration)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunc, uint256 _expiration) returns()
func (_Oracle *OracleSession) CancelOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.CancelOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackFunc, _expiration)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x6ee4d553.
//
// Solidity: function cancelOracleRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunc, uint256 _expiration) returns()
func (_Oracle *OracleTransactorSession) CancelOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.CancelOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackFunc, _expiration)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 _requestId, uint256 _payment, address _callbackAddress, bytes4 _callbackFunctionId, uint256 _expiration, bytes32 _data) returns(bool)
func (_Oracle *OracleTransactor) FulfillOracleRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "fulfillOracleRequest", _requestId, _payment, _callbackAddress, _callbackFunctionId, _expiration, _data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 _requestId, uint256 _payment, address _callbackAddress, bytes4 _callbackFunctionId, uint256 _expiration, bytes32 _data) returns(bool)
func (_Oracle *OracleSession) FulfillOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.FulfillOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackAddress, _callbackFunctionId, _expiration, _data)
}

// FulfillOracleRequest is a paid mutator transaction binding the contract method 0x4ab0d190.
//
// Solidity: function fulfillOracleRequest(bytes32 _requestId, uint256 _payment, address _callbackAddress, bytes4 _callbackFunctionId, uint256 _expiration, bytes32 _data) returns(bool)
func (_Oracle *OracleTransactorSession) FulfillOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.FulfillOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackAddress, _callbackFunctionId, _expiration, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_Oracle *OracleTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "onTokenTransfer", _sender, _amount, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_Oracle *OracleSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OnTokenTransfer(&_Oracle.TransactOpts, _sender, _amount, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _amount, bytes _data) returns()
func (_Oracle *OracleTransactorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OnTokenTransfer(&_Oracle.TransactOpts, _sender, _amount, _data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address _sender, uint256 _payment, bytes32 _specId, address _callbackAddress, bytes4 _callbackFunctionId, uint256 _nonce, uint256 _dataVersion, bytes _data) returns()
func (_Oracle *OracleTransactor) OracleRequest(opts *bind.TransactOpts, _sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "oracleRequest", _sender, _payment, _specId, _callbackAddress, _callbackFunctionId, _nonce, _dataVersion, _data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address _sender, uint256 _payment, bytes32 _specId, address _callbackAddress, bytes4 _callbackFunctionId, uint256 _nonce, uint256 _dataVersion, bytes _data) returns()
func (_Oracle *OracleSession) OracleRequest(_sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OracleRequest(&_Oracle.TransactOpts, _sender, _payment, _specId, _callbackAddress, _callbackFunctionId, _nonce, _dataVersion, _data)
}

// OracleRequest is a paid mutator transaction binding the contract method 0x40429946.
//
// Solidity: function oracleRequest(address _sender, uint256 _payment, bytes32 _specId, address _callbackAddress, bytes4 _callbackFunctionId, uint256 _nonce, uint256 _dataVersion, bytes _data) returns()
func (_Oracle *OracleTransactorSession) OracleRequest(_sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OracleRequest(&_Oracle.TransactOpts, _sender, _payment, _specId, _callbackAddress, _callbackFunctionId, _nonce, _dataVersion, _data)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address _node, bool _allowed) returns()
func (_Oracle *OracleTransactor) SetFulfillmentPermission(opts *bind.TransactOpts, _node common.Address, _allowed bool) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setFulfillmentPermission", _node, _allowed)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address _node, bool _allowed) returns()
func (_Oracle *OracleSession) SetFulfillmentPermission(_node common.Address, _allowed bool) (*types.Transaction, error) {
	return _Oracle.Contract.SetFulfillmentPermission(&_Oracle.TransactOpts, _node, _allowed)
}

// SetFulfillmentPermission is a paid mutator transaction binding the contract method 0x7fcd56db.
//
// Solidity: function setFulfillmentPermission(address _node, bool _allowed) returns()
func (_Oracle *OracleTransactorSession) SetFulfillmentPermission(_node common.Address, _allowed bool) (*types.Transaction, error) {
	return _Oracle.Contract.SetFulfillmentPermission(&_Oracle.TransactOpts, _node, _allowed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracle *OracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracle *OracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.TransferOwnership(&_Oracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracle *OracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.TransferOwnership(&_Oracle.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_Oracle *OracleTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "withdraw", _recipient, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_Oracle *OracleSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Withdraw(&_Oracle.TransactOpts, _recipient, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_Oracle *OracleTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Withdraw(&_Oracle.TransactOpts, _recipient, _amount)
}

// OracleCancelOracleRequestIterator is returned from FilterCancelOracleRequest and is used to iterate over the raw logs and unpacked data for CancelOracleRequest events raised by the Oracle contract.
type OracleCancelOracleRequestIterator struct {
	Event *OracleCancelOracleRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleCancelOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleCancelOracleRequest)
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
		it.Event = new(OracleCancelOracleRequest)
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
func (it *OracleCancelOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleCancelOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleCancelOracleRequest represents a CancelOracleRequest event raised by the Oracle contract.
type OracleCancelOracleRequest struct {
	RequestId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCancelOracleRequest is a free log retrieval operation binding the contract event 0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93.
//
// Solidity: event CancelOracleRequest(bytes32 indexed requestId)
func (_Oracle *OracleFilterer) FilterCancelOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*OracleCancelOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "CancelOracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleCancelOracleRequestIterator{contract: _Oracle.contract, event: "CancelOracleRequest", logs: logs, sub: sub}, nil
}

// WatchCancelOracleRequest is a free log subscription operation binding the contract event 0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93.
//
// Solidity: event CancelOracleRequest(bytes32 indexed requestId)
func (_Oracle *OracleFilterer) WatchCancelOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleCancelOracleRequest, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "CancelOracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleCancelOracleRequest)
				if err := _Oracle.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
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

// ParseCancelOracleRequest is a log parse operation binding the contract event 0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93.
//
// Solidity: event CancelOracleRequest(bytes32 indexed requestId)
func (_Oracle *OracleFilterer) ParseCancelOracleRequest(log types.Log) (*OracleCancelOracleRequest, error) {
	event := new(OracleCancelOracleRequest)
	if err := _Oracle.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleOracleRequestIterator is returned from FilterOracleRequest and is used to iterate over the raw logs and unpacked data for OracleRequest events raised by the Oracle contract.
type OracleOracleRequestIterator struct {
	Event *OracleOracleRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOracleRequest)
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
		it.Event = new(OracleOracleRequest)
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
func (it *OracleOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleOracleRequest represents a OracleRequest event raised by the Oracle contract.
type OracleOracleRequest struct {
	SpecId             [32]byte
	Requester          common.Address
	RequestId          [32]byte
	Payment            *big.Int
	CallbackAddr       common.Address
	CallbackFunctionId [4]byte
	CancelExpiration   *big.Int
	DataVersion        *big.Int
	Data               []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOracleRequest is a free log retrieval operation binding the contract event 0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65.
//
// Solidity: event OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)
func (_Oracle *OracleFilterer) FilterOracleRequest(opts *bind.FilterOpts, specId [][32]byte) (*OracleOracleRequestIterator, error) {

	var specIdRule []interface{}
	for _, specIdItem := range specId {
		specIdRule = append(specIdRule, specIdItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OracleRequest", specIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleOracleRequestIterator{contract: _Oracle.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

// WatchOracleRequest is a free log subscription operation binding the contract event 0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65.
//
// Solidity: event OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)
func (_Oracle *OracleFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleOracleRequest, specId [][32]byte) (event.Subscription, error) {

	var specIdRule []interface{}
	for _, specIdItem := range specId {
		specIdRule = append(specIdRule, specIdItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OracleRequest", specIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleOracleRequest)
				if err := _Oracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

// ParseOracleRequest is a log parse operation binding the contract event 0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65.
//
// Solidity: event OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)
func (_Oracle *OracleFilterer) ParseOracleRequest(log types.Log) (*OracleOracleRequest, error) {
	event := new(OracleOracleRequest)
	if err := _Oracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Oracle contract.
type OracleOwnershipTransferredIterator struct {
	Event *OracleOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOwnershipTransferred)
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
		it.Event = new(OracleOwnershipTransferred)
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
func (it *OracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleOwnershipTransferred represents a OwnershipTransferred event raised by the Oracle contract.
type OracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleOwnershipTransferredIterator{contract: _Oracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleOwnershipTransferred)
				if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) ParseOwnershipTransferred(log types.Log) (*OracleOwnershipTransferred, error) {
	event := new(OracleOwnershipTransferred)
	if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
