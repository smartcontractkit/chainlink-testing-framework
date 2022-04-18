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

// FlagsMetaData contains all meta data concerning the Flags contract.
var FlagsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"racAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"FlagLowered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"FlagRaised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RaisingAccessControllerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"getFlag\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"lowerFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"raiseFlag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"raiseFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"raisingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"racAddress\",\"type\":\"address\"}],\"name\":\"setRaisingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200101838038062001018833981016040819052620000349162000263565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e6565b50506001805460ff60a01b1916600160a01b17905550620000df8162000191565b5062000295565b336001600160a01b03821603620001405760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6200019b62000205565b6003546001600160a01b0390811690821681146200020157600380546001600160a01b0319166001600160a01b0384811691821790925560405190918316907fbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac90600090a35b5050565b6000546001600160a01b03163314620002615760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b6000602082840312156200027657600080fd5b81516001600160a01b03811681146200028e57600080fd5b9392505050565b610d7380620002a56000396000f3fe608060405234801561001057600080fd5b50600436106100e05760003560e01c80637d723cac116100875780637d723cac1461018f5780638038e4a1146101af5780638823da6c146101b75780638da5cb5b146101ca578063a118f249146101db578063d74af263146101ee578063dc7f012414610201578063f2fde38b1461021557600080fd5b80630a756983146100e557806328286596146100ef5780632e1d859c14610102578063357e47fe1461012b578063517e89fe1461014e5780636b14daf814610161578063760bc82d1461017457806379ba509714610187575b600080fd5b6100ed610228565b005b6100ed6100fd366004610a6e565b61027b565b600354610115906001600160a01b031681565b6040516101229190610ae3565b60405180910390f35b61013e610139366004610b0e565b61033a565b6040519015158152602001610122565b6100ed61015c366004610b0e565b6103c5565b61013e61016f366004610b3f565b610436565b6100ed610182366004610a6e565b61045c565b6100ed6104cc565b6101a261019d366004610a6e565b610576565b6040516101229190610c01565b6100ed6106b0565b6100ed6101c5366004610b0e565b610707565b6000546001600160a01b0316610115565b6100ed6101e9366004610b0e565b61078c565b6100ed6101fc366004610b0e565b610807565b60015461013e90600160a01b900460ff1681565b6100ed610223366004610b0e565b610834565b610230610845565b600154600160a01b900460ff1615610279576001805460ff60a01b191690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b610283610845565b60005b818110156103355760008383838181106102a2576102a2610c47565b90506020020160208101906102b79190610b0e565b6001600160a01b03811660009081526004602052604090205490915060ff1615610322576001600160a01b038116600081815260046020526040808220805460ff19169055517fd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc89190a25b508061032d81610c5d565b915050610286565b505050565b600061037d336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061043692505050565b6103a25760405162461bcd60e51b815260040161039990610c84565b60405180910390fd5b506001600160a01b03811660009081526004602052604090205460ff165b919050565b6103cd610845565b6003546001600160a01b03908116908216811461043257600380546001600160a01b0319166001600160a01b0384811691821790925560405190918316907fbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac90600090a35b5050565b60006104428383610898565b8061045557506001600160a01b03831632145b9392505050565b6104646108cf565b6104805760405162461bcd60e51b815260040161039990610ca7565b60005b81811015610335576104ba8383838181106104a0576104a0610c47565b90506020020160208101906104b59190610b0e565b61095f565b806104c481610c5d565b915050610483565b6001546001600160a01b0316331461051f5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610399565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60606105b9336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061043692505050565b6105d55760405162461bcd60e51b815260040161039990610c84565b60008267ffffffffffffffff8111156105f0576105f0610b29565b604051908082528060200260200182016040528015610619578160200160208202803683370190505b50905060005b838110156106a8576004600086868481811061063d5761063d610c47565b90506020020160208101906106529190610b0e565b6001600160a01b03168152602081019190915260400160002054825160ff9091169083908390811061068657610686610c47565b91151560209283029190910190910152806106a081610c5d565b91505061061f565b509392505050565b6106b8610845565b600154600160a01b900460ff16610279576001805460ff60a01b1916600160a01b1790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b61070f610845565b6001600160a01b03811660009081526002602052604090205460ff1615610789576001600160a01b03811660009081526002602052604090819020805460ff19169055517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d190610780908390610ae3565b60405180910390a15b50565b610794610845565b6001600160a01b03811660009081526002602052604090205460ff16610789576001600160a01b03811660009081526002602052604090819020805460ff19166001179055517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db490610780908390610ae3565b61080f6108cf565b61082b5760405162461bcd60e51b815260040161039990610ca7565b6107898161095f565b61083c610845565b610789816109cb565b6000546001600160a01b031633146102795760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610399565b6001600160a01b03821660009081526002602052604081205460ff1680610455575050600154600160a01b900460ff161592915050565b600080546001600160a01b031633148061095a5750600354604051630d629b5f60e31b81526001600160a01b0390911690636b14daf8906109199033906000903690600401610cdb565b602060405180830381865afa158015610936573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061095a9190610d1b565b905090565b6001600160a01b03811660009081526004602052604090205460ff16610789576001600160a01b038116600081815260046020526040808220805460ff19166001179055517f881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c0069190a250565b336001600160a01b03821603610a1d5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610399565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060208385031215610a8157600080fd5b823567ffffffffffffffff80821115610a9957600080fd5b818501915085601f830112610aad57600080fd5b813581811115610abc57600080fd5b8660208260051b8501011115610ad157600080fd5b60209290920196919550909350505050565b6001600160a01b0391909116815260200190565b80356001600160a01b03811681146103c057600080fd5b600060208284031215610b2057600080fd5b61045582610af7565b634e487b7160e01b600052604160045260246000fd5b60008060408385031215610b5257600080fd5b610b5b83610af7565b9150602083013567ffffffffffffffff80821115610b7857600080fd5b818501915085601f830112610b8c57600080fd5b813581811115610b9e57610b9e610b29565b604051601f8201601f19908116603f01168101908382118183101715610bc657610bc6610b29565b81604052828152886020848701011115610bdf57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6020808252825182820181905260009190848201906040850190845b81811015610c3b578351151583529284019291840191600101610c1d565b50909695505050505050565b634e487b7160e01b600052603260045260246000fd5b600060018201610c7d57634e487b7160e01b600052601160045260246000fd5b5060010190565b6020808252600990820152684e6f2061636365737360b81b604082015260600190565b6020808252601a90820152794e6f7420616c6c6f77656420746f20726169736520666c61677360301b604082015260600190565b6001600160a01b03841681526040602082018190528101829052818360608301376000818301606090810191909152601f909201601f1916010192915050565b600060208284031215610d2d57600080fd5b8151801515811461045557600080fdfea26469706673582212207a80c72a55d52a53a1d8e4e8f31dee71ff7da558dc68138037d9b7a4a8d8dad264736f6c634300080d0033",
}

// FlagsABI is the input ABI used to generate the binding from.
// Deprecated: Use FlagsMetaData.ABI instead.
var FlagsABI = FlagsMetaData.ABI

// FlagsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FlagsMetaData.Bin instead.
var FlagsBin = FlagsMetaData.Bin

// DeployFlags deploys a new Ethereum contract, binding an instance of Flags to it.
func DeployFlags(auth *bind.TransactOpts, backend bind.ContractBackend, racAddress common.Address) (common.Address, *types.Transaction, *Flags, error) {
	parsed, err := FlagsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FlagsBin), backend, racAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Flags{FlagsCaller: FlagsCaller{contract: contract}, FlagsTransactor: FlagsTransactor{contract: contract}, FlagsFilterer: FlagsFilterer{contract: contract}}, nil
}

// Flags is an auto generated Go binding around an Ethereum contract.
type Flags struct {
	FlagsCaller     // Read-only binding to the contract
	FlagsTransactor // Write-only binding to the contract
	FlagsFilterer   // Log filterer for contract events
}

// FlagsCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlagsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlagsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlagsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlagsSession struct {
	Contract     *Flags            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlagsCallerSession struct {
	Contract *FlagsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FlagsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlagsTransactorSession struct {
	Contract     *FlagsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlagsRaw struct {
	Contract *Flags // Generic contract binding to access the raw methods on
}

// FlagsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlagsCallerRaw struct {
	Contract *FlagsCaller // Generic read-only contract binding to access the raw methods on
}

// FlagsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlagsTransactorRaw struct {
	Contract *FlagsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlags creates a new instance of Flags, bound to a specific deployed contract.
func NewFlags(address common.Address, backend bind.ContractBackend) (*Flags, error) {
	contract, err := bindFlags(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Flags{FlagsCaller: FlagsCaller{contract: contract}, FlagsTransactor: FlagsTransactor{contract: contract}, FlagsFilterer: FlagsFilterer{contract: contract}}, nil
}

// NewFlagsCaller creates a new read-only instance of Flags, bound to a specific deployed contract.
func NewFlagsCaller(address common.Address, caller bind.ContractCaller) (*FlagsCaller, error) {
	contract, err := bindFlags(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsCaller{contract: contract}, nil
}

// NewFlagsTransactor creates a new write-only instance of Flags, bound to a specific deployed contract.
func NewFlagsTransactor(address common.Address, transactor bind.ContractTransactor) (*FlagsTransactor, error) {
	contract, err := bindFlags(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsTransactor{contract: contract}, nil
}

// NewFlagsFilterer creates a new log filterer instance of Flags, bound to a specific deployed contract.
func NewFlagsFilterer(address common.Address, filterer bind.ContractFilterer) (*FlagsFilterer, error) {
	contract, err := bindFlags(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlagsFilterer{contract: contract}, nil
}

// bindFlags binds a generic wrapper to an already deployed contract.
func bindFlags(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlagsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flags *FlagsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flags.Contract.FlagsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flags *FlagsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.Contract.FlagsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flags *FlagsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flags.Contract.FlagsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flags *FlagsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flags.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flags *FlagsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flags *FlagsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flags.Contract.contract.Transact(opts, method, params...)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsSession) CheckEnabled() (bool, error) {
	return _Flags.Contract.CheckEnabled(&_Flags.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsCallerSession) CheckEnabled() (bool, error) {
	return _Flags.Contract.CheckEnabled(&_Flags.CallOpts)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsCaller) GetFlag(opts *bind.CallOpts, subject common.Address) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "getFlag", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsSession) GetFlag(subject common.Address) (bool, error) {
	return _Flags.Contract.GetFlag(&_Flags.CallOpts, subject)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsCallerSession) GetFlag(subject common.Address) (bool, error) {
	return _Flags.Contract.GetFlag(&_Flags.CallOpts, subject)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsCaller) GetFlags(opts *bind.CallOpts, subjects []common.Address) ([]bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "getFlags", subjects)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsSession) GetFlags(subjects []common.Address) ([]bool, error) {
	return _Flags.Contract.GetFlags(&_Flags.CallOpts, subjects)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsCallerSession) GetFlags(subjects []common.Address) ([]bool, error) {
	return _Flags.Contract.GetFlags(&_Flags.CallOpts, subjects)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _Flags.Contract.HasAccess(&_Flags.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _Flags.Contract.HasAccess(&_Flags.CallOpts, _user, _calldata)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsSession) Owner() (common.Address, error) {
	return _Flags.Contract.Owner(&_Flags.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsCallerSession) Owner() (common.Address, error) {
	return _Flags.Contract.Owner(&_Flags.CallOpts)
}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsCaller) RaisingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "raisingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsSession) RaisingAccessController() (common.Address, error) {
	return _Flags.Contract.RaisingAccessController(&_Flags.CallOpts)
}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsCallerSession) RaisingAccessController() (common.Address, error) {
	return _Flags.Contract.RaisingAccessController(&_Flags.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsSession) AcceptOwnership() (*types.Transaction, error) {
	return _Flags.Contract.AcceptOwnership(&_Flags.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Flags.Contract.AcceptOwnership(&_Flags.TransactOpts)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.AddAccess(&_Flags.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.AddAccess(&_Flags.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsSession) DisableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.DisableAccessCheck(&_Flags.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.DisableAccessCheck(&_Flags.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsSession) EnableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.EnableAccessCheck(&_Flags.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.EnableAccessCheck(&_Flags.TransactOpts)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsTransactor) LowerFlags(opts *bind.TransactOpts, subjects []common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "lowerFlags", subjects)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsSession) LowerFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.LowerFlags(&_Flags.TransactOpts, subjects)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsTransactorSession) LowerFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.LowerFlags(&_Flags.TransactOpts, subjects)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsTransactor) RaiseFlag(opts *bind.TransactOpts, subject common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "raiseFlag", subject)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsSession) RaiseFlag(subject common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlag(&_Flags.TransactOpts, subject)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsTransactorSession) RaiseFlag(subject common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlag(&_Flags.TransactOpts, subject)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsTransactor) RaiseFlags(opts *bind.TransactOpts, subjects []common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "raiseFlags", subjects)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsSession) RaiseFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlags(&_Flags.TransactOpts, subjects)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsTransactorSession) RaiseFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlags(&_Flags.TransactOpts, subjects)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RemoveAccess(&_Flags.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RemoveAccess(&_Flags.TransactOpts, _user)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsTransactor) SetRaisingAccessController(opts *bind.TransactOpts, racAddress common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "setRaisingAccessController", racAddress)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsSession) SetRaisingAccessController(racAddress common.Address) (*types.Transaction, error) {
	return _Flags.Contract.SetRaisingAccessController(&_Flags.TransactOpts, racAddress)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsTransactorSession) SetRaisingAccessController(racAddress common.Address) (*types.Transaction, error) {
	return _Flags.Contract.SetRaisingAccessController(&_Flags.TransactOpts, racAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Flags *FlagsTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Flags *FlagsSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Flags.Contract.TransferOwnership(&_Flags.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Flags *FlagsTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Flags.Contract.TransferOwnership(&_Flags.TransactOpts, to)
}

// FlagsAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the Flags contract.
type FlagsAddedAccessIterator struct {
	Event *FlagsAddedAccess // Event containing the contract specifics and raw log

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
func (it *FlagsAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsAddedAccess)
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
		it.Event = new(FlagsAddedAccess)
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
func (it *FlagsAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsAddedAccess represents a AddedAccess event raised by the Flags contract.
type FlagsAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*FlagsAddedAccessIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &FlagsAddedAccessIterator{contract: _Flags.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *FlagsAddedAccess) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsAddedAccess)
				if err := _Flags.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

// ParseAddedAccess is a log parse operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) ParseAddedAccess(log types.Log) (*FlagsAddedAccess, error) {
	event := new(FlagsAddedAccess)
	if err := _Flags.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the Flags contract.
type FlagsCheckAccessDisabledIterator struct {
	Event *FlagsCheckAccessDisabled // Event containing the contract specifics and raw log

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
func (it *FlagsCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsCheckAccessDisabled)
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
		it.Event = new(FlagsCheckAccessDisabled)
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
func (it *FlagsCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsCheckAccessDisabled represents a CheckAccessDisabled event raised by the Flags contract.
type FlagsCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*FlagsCheckAccessDisabledIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &FlagsCheckAccessDisabledIterator{contract: _Flags.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *FlagsCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsCheckAccessDisabled)
				if err := _Flags.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

// ParseCheckAccessDisabled is a log parse operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) ParseCheckAccessDisabled(log types.Log) (*FlagsCheckAccessDisabled, error) {
	event := new(FlagsCheckAccessDisabled)
	if err := _Flags.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the Flags contract.
type FlagsCheckAccessEnabledIterator struct {
	Event *FlagsCheckAccessEnabled // Event containing the contract specifics and raw log

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
func (it *FlagsCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsCheckAccessEnabled)
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
		it.Event = new(FlagsCheckAccessEnabled)
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
func (it *FlagsCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsCheckAccessEnabled represents a CheckAccessEnabled event raised by the Flags contract.
type FlagsCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*FlagsCheckAccessEnabledIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &FlagsCheckAccessEnabledIterator{contract: _Flags.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *FlagsCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsCheckAccessEnabled)
				if err := _Flags.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

// ParseCheckAccessEnabled is a log parse operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) ParseCheckAccessEnabled(log types.Log) (*FlagsCheckAccessEnabled, error) {
	event := new(FlagsCheckAccessEnabled)
	if err := _Flags.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsFlagLoweredIterator is returned from FilterFlagLowered and is used to iterate over the raw logs and unpacked data for FlagLowered events raised by the Flags contract.
type FlagsFlagLoweredIterator struct {
	Event *FlagsFlagLowered // Event containing the contract specifics and raw log

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
func (it *FlagsFlagLoweredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsFlagLowered)
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
		it.Event = new(FlagsFlagLowered)
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
func (it *FlagsFlagLoweredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsFlagLoweredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsFlagLowered represents a FlagLowered event raised by the Flags contract.
type FlagsFlagLowered struct {
	Subject common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFlagLowered is a free log retrieval operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) FilterFlagLowered(opts *bind.FilterOpts, subject []common.Address) (*FlagsFlagLoweredIterator, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "FlagLowered", subjectRule)
	if err != nil {
		return nil, err
	}
	return &FlagsFlagLoweredIterator{contract: _Flags.contract, event: "FlagLowered", logs: logs, sub: sub}, nil
}

// WatchFlagLowered is a free log subscription operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) WatchFlagLowered(opts *bind.WatchOpts, sink chan<- *FlagsFlagLowered, subject []common.Address) (event.Subscription, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "FlagLowered", subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsFlagLowered)
				if err := _Flags.contract.UnpackLog(event, "FlagLowered", log); err != nil {
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

// ParseFlagLowered is a log parse operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) ParseFlagLowered(log types.Log) (*FlagsFlagLowered, error) {
	event := new(FlagsFlagLowered)
	if err := _Flags.contract.UnpackLog(event, "FlagLowered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsFlagRaisedIterator is returned from FilterFlagRaised and is used to iterate over the raw logs and unpacked data for FlagRaised events raised by the Flags contract.
type FlagsFlagRaisedIterator struct {
	Event *FlagsFlagRaised // Event containing the contract specifics and raw log

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
func (it *FlagsFlagRaisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsFlagRaised)
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
		it.Event = new(FlagsFlagRaised)
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
func (it *FlagsFlagRaisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsFlagRaisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsFlagRaised represents a FlagRaised event raised by the Flags contract.
type FlagsFlagRaised struct {
	Subject common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFlagRaised is a free log retrieval operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) FilterFlagRaised(opts *bind.FilterOpts, subject []common.Address) (*FlagsFlagRaisedIterator, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "FlagRaised", subjectRule)
	if err != nil {
		return nil, err
	}
	return &FlagsFlagRaisedIterator{contract: _Flags.contract, event: "FlagRaised", logs: logs, sub: sub}, nil
}

// WatchFlagRaised is a free log subscription operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) WatchFlagRaised(opts *bind.WatchOpts, sink chan<- *FlagsFlagRaised, subject []common.Address) (event.Subscription, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "FlagRaised", subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsFlagRaised)
				if err := _Flags.contract.UnpackLog(event, "FlagRaised", log); err != nil {
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

// ParseFlagRaised is a log parse operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) ParseFlagRaised(log types.Log) (*FlagsFlagRaised, error) {
	event := new(FlagsFlagRaised)
	if err := _Flags.contract.UnpackLog(event, "FlagRaised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Flags contract.
type FlagsOwnershipTransferRequestedIterator struct {
	Event *FlagsOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *FlagsOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsOwnershipTransferRequested)
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
		it.Event = new(FlagsOwnershipTransferRequested)
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
func (it *FlagsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Flags contract.
type FlagsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FlagsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FlagsOwnershipTransferRequestedIterator{contract: _Flags.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FlagsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsOwnershipTransferRequested)
				if err := _Flags.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) ParseOwnershipTransferRequested(log types.Log) (*FlagsOwnershipTransferRequested, error) {
	event := new(FlagsOwnershipTransferRequested)
	if err := _Flags.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Flags contract.
type FlagsOwnershipTransferredIterator struct {
	Event *FlagsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FlagsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsOwnershipTransferred)
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
		it.Event = new(FlagsOwnershipTransferred)
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
func (it *FlagsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsOwnershipTransferred represents a OwnershipTransferred event raised by the Flags contract.
type FlagsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FlagsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FlagsOwnershipTransferredIterator{contract: _Flags.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FlagsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsOwnershipTransferred)
				if err := _Flags.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) ParseOwnershipTransferred(log types.Log) (*FlagsOwnershipTransferred, error) {
	event := new(FlagsOwnershipTransferred)
	if err := _Flags.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsRaisingAccessControllerUpdatedIterator is returned from FilterRaisingAccessControllerUpdated and is used to iterate over the raw logs and unpacked data for RaisingAccessControllerUpdated events raised by the Flags contract.
type FlagsRaisingAccessControllerUpdatedIterator struct {
	Event *FlagsRaisingAccessControllerUpdated // Event containing the contract specifics and raw log

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
func (it *FlagsRaisingAccessControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsRaisingAccessControllerUpdated)
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
		it.Event = new(FlagsRaisingAccessControllerUpdated)
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
func (it *FlagsRaisingAccessControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsRaisingAccessControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsRaisingAccessControllerUpdated represents a RaisingAccessControllerUpdated event raised by the Flags contract.
type FlagsRaisingAccessControllerUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRaisingAccessControllerUpdated is a free log retrieval operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) FilterRaisingAccessControllerUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*FlagsRaisingAccessControllerUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "RaisingAccessControllerUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &FlagsRaisingAccessControllerUpdatedIterator{contract: _Flags.contract, event: "RaisingAccessControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchRaisingAccessControllerUpdated is a free log subscription operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) WatchRaisingAccessControllerUpdated(opts *bind.WatchOpts, sink chan<- *FlagsRaisingAccessControllerUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "RaisingAccessControllerUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsRaisingAccessControllerUpdated)
				if err := _Flags.contract.UnpackLog(event, "RaisingAccessControllerUpdated", log); err != nil {
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

// ParseRaisingAccessControllerUpdated is a log parse operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) ParseRaisingAccessControllerUpdated(log types.Log) (*FlagsRaisingAccessControllerUpdated, error) {
	event := new(FlagsRaisingAccessControllerUpdated)
	if err := _Flags.contract.UnpackLog(event, "RaisingAccessControllerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the Flags contract.
type FlagsRemovedAccessIterator struct {
	Event *FlagsRemovedAccess // Event containing the contract specifics and raw log

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
func (it *FlagsRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsRemovedAccess)
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
		it.Event = new(FlagsRemovedAccess)
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
func (it *FlagsRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsRemovedAccess represents a RemovedAccess event raised by the Flags contract.
type FlagsRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*FlagsRemovedAccessIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &FlagsRemovedAccessIterator{contract: _Flags.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *FlagsRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsRemovedAccess)
				if err := _Flags.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

// ParseRemovedAccess is a log parse operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) ParseRemovedAccess(log types.Log) (*FlagsRemovedAccess, error) {
	event := new(FlagsRemovedAccess)
	if err := _Flags.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
