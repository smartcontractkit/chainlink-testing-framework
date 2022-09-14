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

// Config is an auto generated low-level Go binding around an user-defined struct.
type Config struct {
	PaymentPremiumPPB    uint32
	FlatFeeMicroLink     uint32
	BlockCountPerTurn    *big.Int
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	MinUpkeepSpend       *big.Int
	MaxPerformGas        uint32
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
	Transcoder           common.Address
	Registrar            common.Address
}

// KeeperRegistryBaseMetaData contains all meta data concerning the KeeperRegistryBase contract.
var KeeperRegistryBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ARB_NITRO_ORACLE\",\"outputs\":[{\"internalType\":\"contractArbGasInfo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPTIMISM_ORACLE\",\"outputs\":[{\"internalType\":\"contractOVM_GasPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAYMENT_MODEL\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REGISTRY_GAS_OVERHEAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// KeeperRegistryBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistryBaseMetaData.ABI instead.
var KeeperRegistryBaseABI = KeeperRegistryBaseMetaData.ABI

// KeeperRegistryBase is an auto generated Go binding around an Ethereum contract.
type KeeperRegistryBase struct {
	KeeperRegistryBaseCaller     // Read-only binding to the contract
	KeeperRegistryBaseTransactor // Write-only binding to the contract
	KeeperRegistryBaseFilterer   // Log filterer for contract events
}

// KeeperRegistryBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistryBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistryBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistryBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistryBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistryBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistryBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistryBaseSession struct {
	Contract     *KeeperRegistryBase // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// KeeperRegistryBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistryBaseCallerSession struct {
	Contract *KeeperRegistryBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// KeeperRegistryBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistryBaseTransactorSession struct {
	Contract     *KeeperRegistryBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// KeeperRegistryBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistryBaseRaw struct {
	Contract *KeeperRegistryBase // Generic contract binding to access the raw methods on
}

// KeeperRegistryBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistryBaseCallerRaw struct {
	Contract *KeeperRegistryBaseCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistryBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistryBaseTransactorRaw struct {
	Contract *KeeperRegistryBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistryBase creates a new instance of KeeperRegistryBase, bound to a specific deployed contract.
func NewKeeperRegistryBase(address common.Address, backend bind.ContractBackend) (*KeeperRegistryBase, error) {
	contract, err := bindKeeperRegistryBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBase{KeeperRegistryBaseCaller: KeeperRegistryBaseCaller{contract: contract}, KeeperRegistryBaseTransactor: KeeperRegistryBaseTransactor{contract: contract}, KeeperRegistryBaseFilterer: KeeperRegistryBaseFilterer{contract: contract}}, nil
}

// NewKeeperRegistryBaseCaller creates a new read-only instance of KeeperRegistryBase, bound to a specific deployed contract.
func NewKeeperRegistryBaseCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryBaseCaller, error) {
	contract, err := bindKeeperRegistryBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseCaller{contract: contract}, nil
}

// NewKeeperRegistryBaseTransactor creates a new write-only instance of KeeperRegistryBase, bound to a specific deployed contract.
func NewKeeperRegistryBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryBaseTransactor, error) {
	contract, err := bindKeeperRegistryBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseTransactor{contract: contract}, nil
}

// NewKeeperRegistryBaseFilterer creates a new log filterer instance of KeeperRegistryBase, bound to a specific deployed contract.
func NewKeeperRegistryBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryBaseFilterer, error) {
	contract, err := bindKeeperRegistryBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseFilterer{contract: contract}, nil
}

// bindKeeperRegistryBase binds a generic wrapper to an already deployed contract.
func bindKeeperRegistryBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistryBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistryBase *KeeperRegistryBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryBase.Contract.KeeperRegistryBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistryBase *KeeperRegistryBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.KeeperRegistryBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistryBase *KeeperRegistryBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.KeeperRegistryBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistryBase *KeeperRegistryBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistryBase *KeeperRegistryBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistryBase *KeeperRegistryBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.contract.Transact(opts, method, params...)
}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "ARB_NITRO_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistryBase.Contract.ARBNITROORACLE(&_KeeperRegistryBase.CallOpts)
}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistryBase.Contract.ARBNITROORACLE(&_KeeperRegistryBase.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistryBase.Contract.FASTGASFEED(&_KeeperRegistryBase.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistryBase.Contract.FASTGASFEED(&_KeeperRegistryBase.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) LINK() (common.Address, error) {
	return _KeeperRegistryBase.Contract.LINK(&_KeeperRegistryBase.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) LINK() (common.Address, error) {
	return _KeeperRegistryBase.Contract.LINK(&_KeeperRegistryBase.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistryBase.Contract.LINKETHFEED(&_KeeperRegistryBase.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistryBase.Contract.LINKETHFEED(&_KeeperRegistryBase.CallOpts)
}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "OPTIMISM_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistryBase.Contract.OPTIMISMORACLE(&_KeeperRegistryBase.CallOpts)
}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistryBase.Contract.OPTIMISMORACLE(&_KeeperRegistryBase.CallOpts)
}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) PAYMENTMODEL(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "PAYMENT_MODEL")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistryBase.Contract.PAYMENTMODEL(&_KeeperRegistryBase.CallOpts)
}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistryBase.Contract.PAYMENTMODEL(&_KeeperRegistryBase.CallOpts)
}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "REGISTRY_GAS_OVERHEAD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistryBase.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistryBase.CallOpts)
}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistryBase.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistryBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) Owner() (common.Address, error) {
	return _KeeperRegistryBase.Contract.Owner(&_KeeperRegistryBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryBase.Contract.Owner(&_KeeperRegistryBase.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistryBase *KeeperRegistryBaseCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistryBase.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistryBase *KeeperRegistryBaseSession) Paused() (bool, error) {
	return _KeeperRegistryBase.Contract.Paused(&_KeeperRegistryBase.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistryBase *KeeperRegistryBaseCallerSession) Paused() (bool, error) {
	return _KeeperRegistryBase.Contract.Paused(&_KeeperRegistryBase.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistryBase *KeeperRegistryBaseTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryBase.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistryBase *KeeperRegistryBaseSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.AcceptOwnership(&_KeeperRegistryBase.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistryBase *KeeperRegistryBaseTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.AcceptOwnership(&_KeeperRegistryBase.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistryBase *KeeperRegistryBaseTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryBase.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistryBase *KeeperRegistryBaseSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.TransferOwnership(&_KeeperRegistryBase.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistryBase *KeeperRegistryBaseTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryBase.Contract.TransferOwnership(&_KeeperRegistryBase.TransactOpts, to)
}

// KeeperRegistryBaseConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseConfigSetIterator struct {
	Event *KeeperRegistryBaseConfigSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseConfigSet)
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
		it.Event = new(KeeperRegistryBaseConfigSet)
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
func (it *KeeperRegistryBaseConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseConfigSet represents a ConfigSet event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseConfigSet struct {
	Config Config
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryBaseConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseConfigSetIterator{contract: _KeeperRegistryBase.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseConfigSet)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseConfigSet(log types.Log) (*KeeperRegistryBaseConfigSet, error) {
	event := new(KeeperRegistryBaseConfigSet)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseFundsAddedIterator is returned from FilterFundsAdded and is used to iterate over the raw logs and unpacked data for FundsAdded events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseFundsAddedIterator struct {
	Event *KeeperRegistryBaseFundsAdded // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseFundsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseFundsAdded)
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
		it.Event = new(KeeperRegistryBaseFundsAdded)
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
func (it *KeeperRegistryBaseFundsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseFundsAdded represents a FundsAdded event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsAdded is a free log retrieval operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryBaseFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseFundsAddedIterator{contract: _KeeperRegistryBase.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

// WatchFundsAdded is a free log subscription operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseFundsAdded)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

// ParseFundsAdded is a log parse operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryBaseFundsAdded, error) {
	event := new(KeeperRegistryBaseFundsAdded)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseFundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseFundsWithdrawnIterator struct {
	Event *KeeperRegistryBaseFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseFundsWithdrawn)
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
		it.Event = new(KeeperRegistryBaseFundsWithdrawn)
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
func (it *KeeperRegistryBaseFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseFundsWithdrawn represents a FundsWithdrawn event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseFundsWithdrawnIterator{contract: _KeeperRegistryBase.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseFundsWithdrawn)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

// ParseFundsWithdrawn is a log parse operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryBaseFundsWithdrawn, error) {
	event := new(KeeperRegistryBaseFundsWithdrawn)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseKeepersUpdatedIterator is returned from FilterKeepersUpdated and is used to iterate over the raw logs and unpacked data for KeepersUpdated events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseKeepersUpdatedIterator struct {
	Event *KeeperRegistryBaseKeepersUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseKeepersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseKeepersUpdated)
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
		it.Event = new(KeeperRegistryBaseKeepersUpdated)
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
func (it *KeeperRegistryBaseKeepersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseKeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseKeepersUpdated represents a KeepersUpdated event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseKeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKeepersUpdated is a free log retrieval operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryBaseKeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseKeepersUpdatedIterator{contract: _KeeperRegistryBase.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

// WatchKeepersUpdated is a free log subscription operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseKeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseKeepersUpdated)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
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

// ParseKeepersUpdated is a log parse operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistryBaseKeepersUpdated, error) {
	event := new(KeeperRegistryBaseKeepersUpdated)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseOwnerFundsWithdrawnIterator is returned from FilterOwnerFundsWithdrawn and is used to iterate over the raw logs and unpacked data for OwnerFundsWithdrawn events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryBaseOwnerFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseOwnerFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseOwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistryBaseOwnerFundsWithdrawn)
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
func (it *KeeperRegistryBaseOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseOwnerFundsWithdrawn represents a OwnerFundsWithdrawn event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerFundsWithdrawn is a free log retrieval operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryBaseOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseOwnerFundsWithdrawnIterator{contract: _KeeperRegistryBase.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchOwnerFundsWithdrawn is a free log subscription operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseOwnerFundsWithdrawn)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

// ParseOwnerFundsWithdrawn is a log parse operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryBaseOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryBaseOwnerFundsWithdrawn)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryBaseOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryBaseOwnershipTransferRequested)
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
func (it *KeeperRegistryBaseOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryBaseOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseOwnershipTransferRequestedIterator{contract: _KeeperRegistryBase.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseOwnershipTransferRequested)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryBaseOwnershipTransferRequested, error) {
	event := new(KeeperRegistryBaseOwnershipTransferRequested)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseOwnershipTransferredIterator struct {
	Event *KeeperRegistryBaseOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseOwnershipTransferred)
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
		it.Event = new(KeeperRegistryBaseOwnershipTransferred)
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
func (it *KeeperRegistryBaseOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseOwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryBaseOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseOwnershipTransferredIterator{contract: _KeeperRegistryBase.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseOwnershipTransferred)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryBaseOwnershipTransferred, error) {
	event := new(KeeperRegistryBaseOwnershipTransferred)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBasePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePausedIterator struct {
	Event *KeeperRegistryBasePaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBasePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBasePaused)
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
		it.Event = new(KeeperRegistryBasePaused)
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
func (it *KeeperRegistryBasePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBasePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBasePaused represents a Paused event raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryBasePausedIterator, error) {

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBasePausedIterator{contract: _KeeperRegistryBase.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBasePaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBasePaused)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParsePaused(log types.Log) (*KeeperRegistryBasePaused, error) {
	event := new(KeeperRegistryBasePaused)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBasePayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryBasePayeeshipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBasePayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBasePayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistryBasePayeeshipTransferRequested)
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
func (it *KeeperRegistryBasePayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBasePayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBasePayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryBasePayeeshipTransferRequestedIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBasePayeeshipTransferRequestedIterator{contract: _KeeperRegistryBase.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBasePayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBasePayeeshipTransferRequested)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryBasePayeeshipTransferRequested, error) {
	event := new(KeeperRegistryBasePayeeshipTransferRequested)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBasePayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePayeeshipTransferredIterator struct {
	Event *KeeperRegistryBasePayeeshipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBasePayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBasePayeeshipTransferred)
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
		it.Event = new(KeeperRegistryBasePayeeshipTransferred)
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
func (it *KeeperRegistryBasePayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBasePayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBasePayeeshipTransferred represents a PayeeshipTransferred event raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryBasePayeeshipTransferredIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBasePayeeshipTransferredIterator{contract: _KeeperRegistryBase.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBasePayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBasePayeeshipTransferred)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryBasePayeeshipTransferred, error) {
	event := new(KeeperRegistryBasePayeeshipTransferred)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBasePaymentWithdrawnIterator is returned from FilterPaymentWithdrawn and is used to iterate over the raw logs and unpacked data for PaymentWithdrawn events raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePaymentWithdrawnIterator struct {
	Event *KeeperRegistryBasePaymentWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBasePaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBasePaymentWithdrawn)
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
		it.Event = new(KeeperRegistryBasePaymentWithdrawn)
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
func (it *KeeperRegistryBasePaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBasePaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBasePaymentWithdrawn represents a PaymentWithdrawn event raised by the KeeperRegistryBase contract.
type KeeperRegistryBasePaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryBasePaymentWithdrawnIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBasePaymentWithdrawnIterator{contract: _KeeperRegistryBase.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBasePaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBasePaymentWithdrawn)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

// ParsePaymentWithdrawn is a log parse operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryBasePaymentWithdrawn, error) {
	event := new(KeeperRegistryBasePaymentWithdrawn)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUnpausedIterator struct {
	Event *KeeperRegistryBaseUnpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUnpaused)
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
		it.Event = new(KeeperRegistryBaseUnpaused)
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
func (it *KeeperRegistryBaseUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUnpaused represents a Unpaused event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryBaseUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUnpausedIterator{contract: _KeeperRegistryBase.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUnpaused)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryBaseUnpaused, error) {
	event := new(KeeperRegistryBaseUnpaused)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepAdminTransferRequestedIterator is returned from FilterUpkeepAdminTransferRequested and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferRequested events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryBaseUpkeepAdminTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepAdminTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepAdminTransferRequested)
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
		it.Event = new(KeeperRegistryBaseUpkeepAdminTransferRequested)
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
func (it *KeeperRegistryBaseUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepAdminTransferRequested represents a UpkeepAdminTransferRequested event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferRequested is a free log retrieval operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryBaseUpkeepAdminTransferRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferRequested is a free log subscription operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepAdminTransferRequested)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

// ParseUpkeepAdminTransferRequested is a log parse operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryBaseUpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryBaseUpkeepAdminTransferRequested)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepAdminTransferredIterator is returned from FilterUpkeepAdminTransferred and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferred events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryBaseUpkeepAdminTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepAdminTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepAdminTransferred)
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
		it.Event = new(KeeperRegistryBaseUpkeepAdminTransferred)
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
func (it *KeeperRegistryBaseUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepAdminTransferred represents a UpkeepAdminTransferred event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferred is a free log retrieval operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryBaseUpkeepAdminTransferredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepAdminTransferredIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferred is a free log subscription operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepAdminTransferred)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

// ParseUpkeepAdminTransferred is a log parse operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryBaseUpkeepAdminTransferred, error) {
	event := new(KeeperRegistryBaseUpkeepAdminTransferred)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepCanceledIterator is returned from FilterUpkeepCanceled and is used to iterate over the raw logs and unpacked data for UpkeepCanceled events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepCanceledIterator struct {
	Event *KeeperRegistryBaseUpkeepCanceled // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepCanceled)
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
		it.Event = new(KeeperRegistryBaseUpkeepCanceled)
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
func (it *KeeperRegistryBaseUpkeepCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepCanceled represents a UpkeepCanceled event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCanceled is a free log retrieval operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryBaseUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepCanceledIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

// WatchUpkeepCanceled is a free log subscription operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepCanceled)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

// ParseUpkeepCanceled is a log parse operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryBaseUpkeepCanceled, error) {
	event := new(KeeperRegistryBaseUpkeepCanceled)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepCheckDataUpdatedIterator is returned from FilterUpkeepCheckDataUpdated and is used to iterate over the raw logs and unpacked data for UpkeepCheckDataUpdated events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistryBaseUpkeepCheckDataUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepCheckDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepCheckDataUpdated)
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
		it.Event = new(KeeperRegistryBaseUpkeepCheckDataUpdated)
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
func (it *KeeperRegistryBaseUpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepCheckDataUpdated represents a UpkeepCheckDataUpdated event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCheckDataUpdated is a free log retrieval operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepCheckDataUpdatedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

// WatchUpkeepCheckDataUpdated is a free log subscription operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepCheckDataUpdated)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
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

// ParseUpkeepCheckDataUpdated is a log parse operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryBaseUpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistryBaseUpkeepCheckDataUpdated)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepGasLimitSetIterator is returned from FilterUpkeepGasLimitSet and is used to iterate over the raw logs and unpacked data for UpkeepGasLimitSet events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryBaseUpkeepGasLimitSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistryBaseUpkeepGasLimitSet)
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
func (it *KeeperRegistryBaseUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepGasLimitSet represents a UpkeepGasLimitSet event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpkeepGasLimitSet is a free log retrieval operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepGasLimitSetIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepGasLimitSet is a free log subscription operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepGasLimitSet)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

// ParseUpkeepGasLimitSet is a log parse operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryBaseUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryBaseUpkeepGasLimitSet)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepMigratedIterator is returned from FilterUpkeepMigrated and is used to iterate over the raw logs and unpacked data for UpkeepMigrated events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepMigratedIterator struct {
	Event *KeeperRegistryBaseUpkeepMigrated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepMigrated)
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
		it.Event = new(KeeperRegistryBaseUpkeepMigrated)
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
func (it *KeeperRegistryBaseUpkeepMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepMigrated represents a UpkeepMigrated event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepMigrated is a free log retrieval operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepMigratedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

// WatchUpkeepMigrated is a free log subscription operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepMigrated)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

// ParseUpkeepMigrated is a log parse operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryBaseUpkeepMigrated, error) {
	event := new(KeeperRegistryBaseUpkeepMigrated)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepPausedIterator is returned from FilterUpkeepPaused and is used to iterate over the raw logs and unpacked data for UpkeepPaused events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepPausedIterator struct {
	Event *KeeperRegistryBaseUpkeepPaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepPaused)
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
		it.Event = new(KeeperRegistryBaseUpkeepPaused)
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
func (it *KeeperRegistryBaseUpkeepPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepPaused represents a UpkeepPaused event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPaused is a free log retrieval operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepPausedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepPaused is a free log subscription operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepPaused)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

// ParseUpkeepPaused is a log parse operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryBaseUpkeepPaused, error) {
	event := new(KeeperRegistryBaseUpkeepPaused)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepPerformedIterator is returned from FilterUpkeepPerformed and is used to iterate over the raw logs and unpacked data for UpkeepPerformed events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepPerformedIterator struct {
	Event *KeeperRegistryBaseUpkeepPerformed // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepPerformedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepPerformed)
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
		it.Event = new(KeeperRegistryBaseUpkeepPerformed)
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
func (it *KeeperRegistryBaseUpkeepPerformedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepPerformed represents a UpkeepPerformed event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepPerformed struct {
	Id          *big.Int
	Success     bool
	From        common.Address
	Payment     *big.Int
	PerformData []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPerformed is a free log retrieval operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryBaseUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepPerformedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepPerformed)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

// ParseUpkeepPerformed is a log parse operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryBaseUpkeepPerformed, error) {
	event := new(KeeperRegistryBaseUpkeepPerformed)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepReceivedIterator is returned from FilterUpkeepReceived and is used to iterate over the raw logs and unpacked data for UpkeepReceived events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepReceivedIterator struct {
	Event *KeeperRegistryBaseUpkeepReceived // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepReceived)
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
		it.Event = new(KeeperRegistryBaseUpkeepReceived)
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
func (it *KeeperRegistryBaseUpkeepReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepReceived represents a UpkeepReceived event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpkeepReceived is a free log retrieval operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepReceivedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

// WatchUpkeepReceived is a free log subscription operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepReceived)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

// ParseUpkeepReceived is a log parse operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryBaseUpkeepReceived, error) {
	event := new(KeeperRegistryBaseUpkeepReceived)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepRegisteredIterator is returned from FilterUpkeepRegistered and is used to iterate over the raw logs and unpacked data for UpkeepRegistered events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepRegisteredIterator struct {
	Event *KeeperRegistryBaseUpkeepRegistered // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepRegistered)
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
		it.Event = new(KeeperRegistryBaseUpkeepRegistered)
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
func (it *KeeperRegistryBaseUpkeepRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepRegistered represents a UpkeepRegistered event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpkeepRegistered is a free log retrieval operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepRegisteredIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

// WatchUpkeepRegistered is a free log subscription operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepRegistered)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

// ParseUpkeepRegistered is a log parse operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryBaseUpkeepRegistered, error) {
	event := new(KeeperRegistryBaseUpkeepRegistered)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryBaseUpkeepUnpausedIterator is returned from FilterUpkeepUnpaused and is used to iterate over the raw logs and unpacked data for UpkeepUnpaused events raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepUnpausedIterator struct {
	Event *KeeperRegistryBaseUpkeepUnpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryBaseUpkeepUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryBaseUpkeepUnpaused)
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
		it.Event = new(KeeperRegistryBaseUpkeepUnpaused)
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
func (it *KeeperRegistryBaseUpkeepUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryBaseUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryBaseUpkeepUnpaused represents a UpkeepUnpaused event raised by the KeeperRegistryBase contract.
type KeeperRegistryBaseUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepUnpaused is a free log retrieval operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryBaseUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryBaseUpkeepUnpausedIterator{contract: _KeeperRegistryBase.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepUnpaused is a free log subscription operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryBaseUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryBase.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryBaseUpkeepUnpaused)
				if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

// ParseUpkeepUnpaused is a log parse operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistryBase *KeeperRegistryBaseFilterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryBaseUpkeepUnpaused, error) {
	event := new(KeeperRegistryBaseUpkeepUnpaused)
	if err := _KeeperRegistryBase.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
