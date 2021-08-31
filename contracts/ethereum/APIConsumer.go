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

// APIConsumerABI is the input ABI used to generate the binding from.
const APIConsumerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_jobId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_url\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_path\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"_times\",\"type\":\"int256\"}],\"name\":\"createRequestTo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_data\",\"type\":\"uint256\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selector\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// APIConsumerBin is the compiled bytecode used for deploying new contracts.
var APIConsumerBin = "0x608060405260016004553480156200001657600080fd5b50604051620016f6380380620016f6833981810160405260208110156200003c57600080fd5b5051600680546001600160a01b0319163317908190556040516001600160a01b0391909116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36001600160a01b038116620000b257620000ac6001600160e01b03620000cd16565b620000c6565b620000c6816001600160e01b036200015e16565b5062000180565b6200015c73c89bd4e1632d3a43cb03aaad5262cbe4038bc5716001600160a01b03166338cc48316040518163ffffffff1660e01b815260040160206040518083038186803b1580156200011f57600080fd5b505afa15801562000134573d6000803e3d6000fd5b505050506040513d60208110156200014b57600080fd5b50516001600160e01b036200015e16565b565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b61156680620001906000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638dc654a211610076578063ea3d508a1161005b578063ea3d508a146102b7578063ec65d0f8146102f4578063f2fde38b14610345576100be565b80638dc654a2146102935780638f32d59b1461029b576100be565b80634357855e116100a75780634357855e1461025e57806373d4a13a146102835780638da5cb5b1461028b576100be565b8063165d35e1146100c357806316ef7f1a146100f4575b600080fd5b6100cb610378565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61024c600480360360c081101561010a57600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516916020810135916040820135919081019060808101606082013564010000000081111561014e57600080fd5b82018360208201111561016057600080fd5b8035906020019184600183028401116401000000008311171561018257600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092959493602081019350359150506401000000008111156101d557600080fd5b8201836020820111156101e757600080fd5b8035906020019184600183028401116401000000008311171561020957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505091359250610387915050565b60408051918252519081900360200190f35b6102816004803603604081101561027457600080fd5b5080359060200135610549565b005b61024c61054f565b6100cb610555565b610281610571565b6102a36107a1565b604080519115158252519081900360200190f35b6102bf6107bf565b604080517fffffffff000000000000000000000000000000000000000000000000000000009092168252519081900360200190f35b6102816004803603608081101561030a57600080fd5b508035906020810135907fffffffff0000000000000000000000000000000000000000000000000000000060408201351690606001356107c8565b6102816004803603602081101561035b57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661084d565b60006103826108c9565b905090565b60006103916107a1565b6103fc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000016634357855e1790556104326114c1565b61045d87307f4357855e000000000000000000000000000000000000000000000000000000006108e5565b60408051808201909152600381527f676574000000000000000000000000000000000000000000000000000000000060208201529091506104a69082908763ffffffff61091016565b60408051808201909152600481527f706174680000000000000000000000000000000000000000000000000000000060208201526104ec9082908663ffffffff61091016565b60408051808201909152600581527f74696d657300000000000000000000000000000000000000000000000000000060208201526105329082908563ffffffff61093f16565b61053d888288610969565b98975050505050505050565b60075550565b60075481565b60065473ffffffffffffffffffffffffffffffffffffffff1690565b6105796107a1565b6105e457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b60006105ee6108c9565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925073ffffffffffffffffffffffffffffffffffffffff83169163a9059cbb91339184916370a08231916024808301926020929190829003018186803b15801561066757600080fd5b505afa15801561067b573d6000803e3d6000fd5b505050506040513d602081101561069157600080fd5b5051604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff909316600484015260248301919091525160448083019260209291908290030181600087803b15801561070757600080fd5b505af115801561071b573d6000803e3d6000fd5b505050506040513d602081101561073157600080fd5b505161079e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015290519081900360640190fd5b50565b60065473ffffffffffffffffffffffffffffffffffffffff16331490565b60085460e01b81565b6107d06107a1565b61083b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b61084784848484610ba6565b50505050565b6108556107a1565b6108c057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b61079e81610ce1565b60025473ffffffffffffffffffffffffffffffffffffffff1690565b6108ed6114c1565b6108f56114c1565b6109078186868663ffffffff610ddb16565b95945050505050565b6080830151610925908363ffffffff610e3d16565b608083015161093a908263ffffffff610e3d16565b505050565b6080830151610954908363ffffffff610e3d16565b608083015161093a908263ffffffff610e5a16565b6004546040805130606090811b60208084019190915260348084018690528451808503909101815260549093018452825192810192909220908601939093526000838152600590915281812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8816179055905182917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af991a260025473ffffffffffffffffffffffffffffffffffffffff16634000aea08584610a4387610ed0565b6040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610ac7578181015183820152602001610aaf565b50505050905090810190601f168015610af45780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610b1557600080fd5b505af1158015610b29573d6000803e3d6000fd5b505050506040513d6020811015610b3f57600080fd5b5051610b96576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806115376023913960400191505060405180910390fd5b6004805460010190559392505050565b60008481526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000008116909155905173ffffffffffffffffffffffffffffffffffffffff9091169186917fe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c59190a2604080517f6ee4d55300000000000000000000000000000000000000000000000000000000815260048101879052602481018690527fffffffff000000000000000000000000000000000000000000000000000000008516604482015260648101849052905173ffffffffffffffffffffffffffffffffffffffff831691636ee4d55391608480830192600092919082900301818387803b158015610cc257600080fd5b505af1158015610cd6573d6000803e3d6000fd5b505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116610d4d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806115116026913960400191505060405180910390fd5b60065460405173ffffffffffffffffffffffffffffffffffffffff8084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610de36114c1565b610df385608001516101006110b9565b505091835273ffffffffffffffffffffffffffffffffffffffff1660208301527fffffffff0000000000000000000000000000000000000000000000000000000016604082015290565b610e4a82600383516110f9565b61093a828263ffffffff61120316565b7fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000811215610e9157610e8c8282611224565b610ecc565b67ffffffffffffffff811315610eab57610e8c8282611281565b60008112610ebf57610e8c826000836110f9565b610ecc82600183196110f9565b5050565b6060634042994660e01b60008084600001518560200151866040015187606001516001896080015160000151604051602401808973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018881526020018781526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610ffc578181015183820152602001610fe4565b50505050905090810190601f1680156110295780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909d169c909c17909b5250989950505050505050505050919050565b6110c16114f6565b60208206156110d65760208206602003820191505b506020808301829052604080518085526000815283019091019052815b92915050565b601781116111205761111a8360e0600585901b16831763ffffffff6112bc16565b5061093a565b60ff811161115657611143836018611fe0600586901b161763ffffffff6112bc16565b5061111a8382600163ffffffff6112d416565b61ffff811161118d5761117a836019611fe0600586901b161763ffffffff6112bc16565b5061111a8382600263ffffffff6112d416565b63ffffffff81116111c6576111b383601a611fe0600586901b161763ffffffff6112bc16565b5061111a8382600463ffffffff6112d416565b67ffffffffffffffff811161093a576111f083601b611fe0600586901b161763ffffffff6112bc16565b506108478382600863ffffffff6112d416565b61120b6114f6565b61121d838460000151518485516112f5565b9392505050565b6112358260c363ffffffff6112bc16565b50610ecc82827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03604051602001808281526020019150506040516020818303038152906040526113dd565b6112928260c263ffffffff6112bc16565b50610ecc8282604051602001808281526020019150506040516020818303038152906040526113dd565b6112c46114f6565b61121d83846000015151846113ea565b6112dc6114f6565b6112ed848560000151518585611435565b949350505050565b6112fd6114f6565b825182111561130b57600080fd5b84602001518285011115611335576113358561132d8760200151878601611493565b6002026114aa565b6000808651805187602083010193508088870111156113545787860182525b505050602084015b6020841061139957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909301926020918201910161135c565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b610e4a82600283516110f9565b6113f26114f6565b8360200151831061140e5761140e8485602001516002026114aa565b83518051602085830101848153508085141561142b576001810182525b5093949350505050565b61143d6114f6565b8460200151848301111561145a5761145a858584016002026114aa565b60006001836101000a0390508551838682010185831982511617815250805184870111156114885783860181525b509495945050505050565b6000818311156114a45750816110f3565b50919050565b81516114b683836110b9565b506108478382611203565b6040805160a0810182526000808252602082018190529181018290526060810191909152608081016114f16114f6565b905290565b60405180604001604052806060815260200160008152509056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c65a164736f6c6343000606000a"

// DeployAPIConsumer deploys a new Ethereum contract, binding an instance of APIConsumer to it.
func DeployAPIConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address) (common.Address, *types.Transaction, *APIConsumer, error) {
	parsed, err := abi.JSON(strings.NewReader(APIConsumerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(APIConsumerBin), backend, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &APIConsumer{APIConsumerCaller: APIConsumerCaller{contract: contract}, APIConsumerTransactor: APIConsumerTransactor{contract: contract}, APIConsumerFilterer: APIConsumerFilterer{contract: contract}}, nil
}

// APIConsumer is an auto generated Go binding around an Ethereum contract.
type APIConsumer struct {
	APIConsumerCaller     // Read-only binding to the contract
	APIConsumerTransactor // Write-only binding to the contract
	APIConsumerFilterer   // Log filterer for contract events
}

// APIConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type APIConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// APIConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type APIConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// APIConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type APIConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// APIConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type APIConsumerSession struct {
	Contract     *APIConsumer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// APIConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type APIConsumerCallerSession struct {
	Contract *APIConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// APIConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type APIConsumerTransactorSession struct {
	Contract     *APIConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// APIConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type APIConsumerRaw struct {
	Contract *APIConsumer // Generic contract binding to access the raw methods on
}

// APIConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type APIConsumerCallerRaw struct {
	Contract *APIConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// APIConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type APIConsumerTransactorRaw struct {
	Contract *APIConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAPIConsumer creates a new instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumer(address common.Address, backend bind.ContractBackend) (*APIConsumer, error) {
	contract, err := bindAPIConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &APIConsumer{APIConsumerCaller: APIConsumerCaller{contract: contract}, APIConsumerTransactor: APIConsumerTransactor{contract: contract}, APIConsumerFilterer: APIConsumerFilterer{contract: contract}}, nil
}

// NewAPIConsumerCaller creates a new read-only instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumerCaller(address common.Address, caller bind.ContractCaller) (*APIConsumerCaller, error) {
	contract, err := bindAPIConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &APIConsumerCaller{contract: contract}, nil
}

// NewAPIConsumerTransactor creates a new write-only instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*APIConsumerTransactor, error) {
	contract, err := bindAPIConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &APIConsumerTransactor{contract: contract}, nil
}

// NewAPIConsumerFilterer creates a new log filterer instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*APIConsumerFilterer, error) {
	contract, err := bindAPIConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &APIConsumerFilterer{contract: contract}, nil
}

// bindAPIConsumer binds a generic wrapper to an already deployed contract.
func bindAPIConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(APIConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_APIConsumer *APIConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _APIConsumer.Contract.APIConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_APIConsumer *APIConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _APIConsumer.Contract.APIConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_APIConsumer *APIConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _APIConsumer.Contract.APIConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_APIConsumer *APIConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _APIConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_APIConsumer *APIConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _APIConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_APIConsumer *APIConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _APIConsumer.Contract.contract.Transact(opts, method, params...)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint256)
func (_APIConsumer *APIConsumerCaller) Data(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "data")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint256)
func (_APIConsumer *APIConsumerSession) Data() (*big.Int, error) {
	return _APIConsumer.Contract.Data(&_APIConsumer.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint256)
func (_APIConsumer *APIConsumerCallerSession) Data() (*big.Int, error) {
	return _APIConsumer.Contract.Data(&_APIConsumer.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_APIConsumer *APIConsumerCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_APIConsumer *APIConsumerSession) GetChainlinkToken() (common.Address, error) {
	return _APIConsumer.Contract.GetChainlinkToken(&_APIConsumer.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_APIConsumer *APIConsumerCallerSession) GetChainlinkToken() (common.Address, error) {
	return _APIConsumer.Contract.GetChainlinkToken(&_APIConsumer.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_APIConsumer *APIConsumerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_APIConsumer *APIConsumerSession) IsOwner() (bool, error) {
	return _APIConsumer.Contract.IsOwner(&_APIConsumer.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_APIConsumer *APIConsumerCallerSession) IsOwner() (bool, error) {
	return _APIConsumer.Contract.IsOwner(&_APIConsumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_APIConsumer *APIConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_APIConsumer *APIConsumerSession) Owner() (common.Address, error) {
	return _APIConsumer.Contract.Owner(&_APIConsumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_APIConsumer *APIConsumerCallerSession) Owner() (common.Address, error) {
	return _APIConsumer.Contract.Owner(&_APIConsumer.CallOpts)
}

// Selector is a free data retrieval call binding the contract method 0xea3d508a.
//
// Solidity: function selector() view returns(bytes4)
func (_APIConsumer *APIConsumerCaller) Selector(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "selector")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// Selector is a free data retrieval call binding the contract method 0xea3d508a.
//
// Solidity: function selector() view returns(bytes4)
func (_APIConsumer *APIConsumerSession) Selector() ([4]byte, error) {
	return _APIConsumer.Contract.Selector(&_APIConsumer.CallOpts)
}

// Selector is a free data retrieval call binding the contract method 0xea3d508a.
//
// Solidity: function selector() view returns(bytes4)
func (_APIConsumer *APIConsumerCallerSession) Selector() ([4]byte, error) {
	return _APIConsumer.Contract.Selector(&_APIConsumer.CallOpts)
}

// CancelRequest is a paid mutator transaction binding the contract method 0xec65d0f8.
//
// Solidity: function cancelRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_APIConsumer *APIConsumerTransactor) CancelRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "cancelRequest", _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0xec65d0f8.
//
// Solidity: function cancelRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_APIConsumer *APIConsumerSession) CancelRequest(_requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CancelRequest(&_APIConsumer.TransactOpts, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0xec65d0f8.
//
// Solidity: function cancelRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_APIConsumer *APIConsumerTransactorSession) CancelRequest(_requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CancelRequest(&_APIConsumer.TransactOpts, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CreateRequestTo is a paid mutator transaction binding the contract method 0x16ef7f1a.
//
// Solidity: function createRequestTo(address _oracle, bytes32 _jobId, uint256 _payment, string _url, string _path, int256 _times) returns(bytes32 requestId)
func (_APIConsumer *APIConsumerTransactor) CreateRequestTo(opts *bind.TransactOpts, _oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "createRequestTo", _oracle, _jobId, _payment, _url, _path, _times)
}

// CreateRequestTo is a paid mutator transaction binding the contract method 0x16ef7f1a.
//
// Solidity: function createRequestTo(address _oracle, bytes32 _jobId, uint256 _payment, string _url, string _path, int256 _times) returns(bytes32 requestId)
func (_APIConsumer *APIConsumerSession) CreateRequestTo(_oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CreateRequestTo(&_APIConsumer.TransactOpts, _oracle, _jobId, _payment, _url, _path, _times)
}

// CreateRequestTo is a paid mutator transaction binding the contract method 0x16ef7f1a.
//
// Solidity: function createRequestTo(address _oracle, bytes32 _jobId, uint256 _payment, string _url, string _path, int256 _times) returns(bytes32 requestId)
func (_APIConsumer *APIConsumerTransactorSession) CreateRequestTo(_oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CreateRequestTo(&_APIConsumer.TransactOpts, _oracle, _jobId, _payment, _url, _path, _times)
}

// Fulfill is a paid mutator transaction binding the contract method 0x4357855e.
//
// Solidity: function fulfill(bytes32 _requestId, uint256 _data) returns()
func (_APIConsumer *APIConsumerTransactor) Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "fulfill", _requestId, _data)
}

// Fulfill is a paid mutator transaction binding the contract method 0x4357855e.
//
// Solidity: function fulfill(bytes32 _requestId, uint256 _data) returns()
func (_APIConsumer *APIConsumerSession) Fulfill(_requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.Fulfill(&_APIConsumer.TransactOpts, _requestId, _data)
}

// Fulfill is a paid mutator transaction binding the contract method 0x4357855e.
//
// Solidity: function fulfill(bytes32 _requestId, uint256 _data) returns()
func (_APIConsumer *APIConsumerTransactorSession) Fulfill(_requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.Fulfill(&_APIConsumer.TransactOpts, _requestId, _data)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_APIConsumer *APIConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_APIConsumer *APIConsumerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _APIConsumer.Contract.TransferOwnership(&_APIConsumer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_APIConsumer *APIConsumerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _APIConsumer.Contract.TransferOwnership(&_APIConsumer.TransactOpts, newOwner)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_APIConsumer *APIConsumerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "withdrawLink")
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_APIConsumer *APIConsumerSession) WithdrawLink() (*types.Transaction, error) {
	return _APIConsumer.Contract.WithdrawLink(&_APIConsumer.TransactOpts)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_APIConsumer *APIConsumerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _APIConsumer.Contract.WithdrawLink(&_APIConsumer.TransactOpts)
}

// APIConsumerChainlinkCancelledIterator is returned from FilterChainlinkCancelled and is used to iterate over the raw logs and unpacked data for ChainlinkCancelled events raised by the APIConsumer contract.
type APIConsumerChainlinkCancelledIterator struct {
	Event *APIConsumerChainlinkCancelled // Event containing the contract specifics and raw log

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
func (it *APIConsumerChainlinkCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerChainlinkCancelled)
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
		it.Event = new(APIConsumerChainlinkCancelled)
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
func (it *APIConsumerChainlinkCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerChainlinkCancelled represents a ChainlinkCancelled event raised by the APIConsumer contract.
type APIConsumerChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkCancelled is a free log retrieval operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*APIConsumerChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerChainlinkCancelledIterator{contract: _APIConsumer.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

// WatchChainlinkCancelled is a free log subscription operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *APIConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerChainlinkCancelled)
				if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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
func (_APIConsumer *APIConsumerFilterer) ParseChainlinkCancelled(log types.Log) (*APIConsumerChainlinkCancelled, error) {
	event := new(APIConsumerChainlinkCancelled)
	if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerChainlinkFulfilledIterator is returned from FilterChainlinkFulfilled and is used to iterate over the raw logs and unpacked data for ChainlinkFulfilled events raised by the APIConsumer contract.
type APIConsumerChainlinkFulfilledIterator struct {
	Event *APIConsumerChainlinkFulfilled // Event containing the contract specifics and raw log

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
func (it *APIConsumerChainlinkFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerChainlinkFulfilled)
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
		it.Event = new(APIConsumerChainlinkFulfilled)
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
func (it *APIConsumerChainlinkFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerChainlinkFulfilled represents a ChainlinkFulfilled event raised by the APIConsumer contract.
type APIConsumerChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkFulfilled is a free log retrieval operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*APIConsumerChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerChainlinkFulfilledIterator{contract: _APIConsumer.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

// WatchChainlinkFulfilled is a free log subscription operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *APIConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerChainlinkFulfilled)
				if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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
func (_APIConsumer *APIConsumerFilterer) ParseChainlinkFulfilled(log types.Log) (*APIConsumerChainlinkFulfilled, error) {
	event := new(APIConsumerChainlinkFulfilled)
	if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerChainlinkRequestedIterator is returned from FilterChainlinkRequested and is used to iterate over the raw logs and unpacked data for ChainlinkRequested events raised by the APIConsumer contract.
type APIConsumerChainlinkRequestedIterator struct {
	Event *APIConsumerChainlinkRequested // Event containing the contract specifics and raw log

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
func (it *APIConsumerChainlinkRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerChainlinkRequested)
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
		it.Event = new(APIConsumerChainlinkRequested)
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
func (it *APIConsumerChainlinkRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerChainlinkRequested represents a ChainlinkRequested event raised by the APIConsumer contract.
type APIConsumerChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkRequested is a free log retrieval operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*APIConsumerChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerChainlinkRequestedIterator{contract: _APIConsumer.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

// WatchChainlinkRequested is a free log subscription operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *APIConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerChainlinkRequested)
				if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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
func (_APIConsumer *APIConsumerFilterer) ParseChainlinkRequested(log types.Log) (*APIConsumerChainlinkRequested, error) {
	event := new(APIConsumerChainlinkRequested)
	if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the APIConsumer contract.
type APIConsumerOwnershipTransferredIterator struct {
	Event *APIConsumerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *APIConsumerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerOwnershipTransferred)
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
		it.Event = new(APIConsumerOwnershipTransferred)
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
func (it *APIConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerOwnershipTransferred represents a OwnershipTransferred event raised by the APIConsumer contract.
type APIConsumerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_APIConsumer *APIConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*APIConsumerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerOwnershipTransferredIterator{contract: _APIConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_APIConsumer *APIConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *APIConsumerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerOwnershipTransferred)
				if err := _APIConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_APIConsumer *APIConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*APIConsumerOwnershipTransferred, error) {
	event := new(APIConsumerOwnershipTransferred)
	if err := _APIConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
