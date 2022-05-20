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

// State is an auto generated low-level Go binding around an user-defined struct.
type State struct {
	Nonce               uint32
	OwnerLinkBalance    *big.Int
	ExpectedLinkBalance *big.Int
	NumUpkeeps          *big.Int
}

// KeeperRegistryMetaData contains all meta data concerning the KeeperRegistry contract.
var KeeperRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getKeeperInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistry.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"}],\"internalType\":\"structState\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"lastKeeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setKeepers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistry.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620055f6380380620055f68339810160408190526200003491620005a5565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fe565b50506001600255506003805460ff191690556001600160a01b0380851660805283811660a052821660c052620000f481620001a9565b50505050620007f0565b336001600160a01b03821603620001585760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001b36200049e565b600d5460e082015163ffffffff91821691161015620001e557604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600c60010160049054906101000a900463ffffffff1663ffffffff16815250600c60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600e81905550806101200151600f81905550806101400151601260006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601360006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051620004939190620006f1565b60405180910390a150565b6000546001600160a01b03163314620004fa5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200051457600080fd5b919050565b60405161018081016001600160401b03811182821017156200054b57634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff811681146200051457600080fd5b805162ffffff811681146200051457600080fd5b805161ffff811681146200051457600080fd5b80516001600160601b03811681146200051457600080fd5b6000806000808486036101e0811215620005be57600080fd5b620005c986620004fc565b9450620005d960208701620004fc565b9350620005e960408701620004fc565b925061018080605f19830112156200060057600080fd5b6200060a62000519565b91506200061a6060880162000551565b82526200062a6080880162000551565b60208301526200063d60a0880162000566565b60408301526200065060c0880162000551565b60608301526200066360e0880162000566565b6080830152610100620006788189016200057a565b60a08401526101206200068d818a016200058d565b60c0850152610140620006a2818b0162000551565b60e0860152610160808b015184870152848b015183870152620006c96101a08c01620004fc565b82870152620006dc6101c08c01620004fc565b90860152509699959850939650909450505050565b815163ffffffff1681526101808101602083015162000718602084018263ffffffff169052565b50604083015162000730604084018262ffffff169052565b50606083015162000749606084018263ffffffff169052565b50608083015162000761608084018262ffffff169052565b5060a08301516200077860a084018261ffff169052565b5060c08301516200079460c08401826001600160601b03169052565b5060e0830151620007ad60e084018263ffffffff169052565b5061010083810151908301526101208084015190830152610140808401516001600160a01b03908116918401919091526101609384015116929091019190915290565b60805160a05160c051614d90620008666000396000818161034e01526133b0015260008181610479015261347501526000818161027b01528181610b3c01528181610d770152818161143c015281816117aa0152818161187201528181611b2501528181611dbd0152611e370152614d906000f3fe608060405234801561001057600080fd5b50600436106101dc5760003560e01c806393f0c1fc11610105578063b7fdb4361161009d578063b7fdb436146104c9578063c41b813a146104dc578063c7c3a19a14610500578063c804802214610527578063da5c67411461053a578063eb5dcd6c1461055b578063ef47a0ce1461056e578063f2fde38b14610581578063faa3e9961461059457600080fd5b806393f0c1fc14610408578063948108f714610428578063a4c0ed361461043b578063a710b2211461044e578063a72aa27e14610461578063ad17836114610474578063b121e1471461049b578063b657bc9c146104ae578063b79550be146104c157600080fd5b80635c975abb116101785780635c975abb14610385578063744bfe611461039c57806379ba5097146103af5780637bbaf1ea146103b75780637d9b97e0146103ca5780638456cb59146103d257806385c1b0ba146103da5780638da5cb5b146103ed5780638e86139b146103f557600080fd5b806306e3b632146101e1578063181f5a771461020a5780631865c57d1461024a578063187256e8146102615780631b6b6d23146102765780631e12b8a5146102aa5780633f4ba83a146103415780634584a4191461034957806348013d7b14610370575b600080fd5b6101f46101ef366004613ac0565b6105cd565b6040516102019190613ae2565b60405180910390f35b61023d6040518060400160405280601481526020017304b6565706572526567697374727920312e322e360641b81525081565b6040516102019190613b82565b6102526106b2565b60405161020193929190613c91565b61027461026f366004613d41565b61090c565b005b61029d7f000000000000000000000000000000000000000000000000000000000000000081565b6040516102019190613d7c565b6103136102b8366004613d90565b6001600160a01b0390811660009081526008602090815260409182902082516060810184528154948516808252600160a01b9095046001600160601b031692810183905260019091015460ff16151592018290529192909190565b604080516001600160a01b03909416845291151560208401526001600160601b031690820152606001610201565b610274610952565b61029d7f000000000000000000000000000000000000000000000000000000000000000081565b610378600081565b6040516102019190613dd5565b60035460ff165b6040519015158152602001610201565b6102746103aa366004613de3565b610964565b610274610bc2565b61038c6103c5366004613e57565b610c71565b610274610cf0565b610274610df5565b6102746103e8366004613ee6565b610e05565b61029d6114c6565b610274610403366004613f39565b6114d5565b61041b610416366004613f7a565b6116ab565b6040516102019190613f93565b610274610436366004613fbe565b6116df565b610274610449366004613fe1565b611867565b61027461045c36600461403a565b6119d0565b61027461046f366004614078565b611ba4565b61029d7f000000000000000000000000000000000000000000000000000000000000000081565b6102746104a9366004613d90565b611ccd565b61041b6104bc366004613f7a565b611d7a565b610274611d9b565b6102746104d736600461409b565b611e90565b6104ef6104ea366004613de3565b6120dd565b6040516102019594939291906140fa565b61051361050e366004613f7a565b6122f3565b604051610201989796959493929190614131565b610274610535366004613f7a565b612461565b61054d6105483660046141aa565b6125c5565b604051908152602001610201565b61027461056936600461403a565b612759565b61027461057c3660046142d5565b612839565b61027461058f366004613d90565b612b29565b6105c06105a2366004613d90565b6001600160a01b03166000908152600b602052604090205460ff1690565b60405161020191906143b3565b606060006105db6005612b3d565b90508084106105fd57604051631390f2a160e01b815260040160405180910390fd5b826000036106125761060f84826143e3565b92505b6000836001600160401b0381111561062c5761062c61421f565b604051908082528060200260200182016040528015610655578160200160208202803683370190505b50905060005b848110156106a75761067861067082886143fa565b600590612b47565b82828151811061068a5761068a614412565b60209081029190910101528061069f81614428565b91505061065b565b509150505b92915050565b6040805160808101825260008082526020820181905291810182905260608101919091526040805161018081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101919091526040805161012081018252600c5463ffffffff8082168352600160201b808304821660208086019190915262ffffff600160401b8504811686880152600160581b85048416606087810191909152600160781b8604909116608087015261ffff600160901b86041660a08701526001600160601b03600160a01b909504851660c0870152600d5480851660e0880152929092049092166101008501819052875260105490921690860152601154928501929092526107f56005612b3d565b606080860191909152815163ffffffff908116855260208084015182168187015260408085015162ffffff90811682890152858501518416948801949094526080808601519094169387019390935260a08085015161ffff169087015260c0808501516001600160601b03169087015260e08085015190921691860191909152600e54610100860152600f546101208601526012546001600160a01b03908116610140870152601354166101608601526004805483518184028101840190945280845287938793909183918301828280156108f957602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116108db575b5050505050905093509350935050909192565b610914612b5a565b6001600160a01b0382166000908152600b60205260409020805482919060ff1916600183600381111561094957610949613dab565b02179055505050565b61095a612b5a565b610962612bad565b565b806001600160a01b03811661098c57604051634e46966960e11b815260040160405180910390fd5b6000838152600760205260409020600201548390600160601b90046001600160a01b031633146109cf5760405163523e0b8360e11b815260040160405180910390fd5b60008481526007602052604090206001015443600160201b9091046001600160401b03161115610a15576040516001627b1a2360e01b0319815260040160405180910390fd5b600c54600085815260076020526040812080546002909101546001600160601b03600160a01b9094048416939182169291169083821015610a7957610a5a8285614441565b9050826001600160601b0316816001600160601b03161115610a795750815b6000610a858285614441565b60008a815260076020526040902080546001600160601b0319169055601054909150610abb9083906001600160601b0316614469565b601080546001600160601b0319166001600160601b03928316179055601154610ae6918316906143e3565b60115560405189907ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a6831890610b1d9084908c90614494565b60405180910390a260405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90610b73908b9085906004016144b6565b6020604051808303816000875af1158015610b92573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bb691906144e8565b50505050505050505050565b6001546001600160a01b03163314610c1a5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610c7f60035460ff1690565b15610c9c5760405162461bcd60e51b8152600401610c1190614503565b610ce8610ce3338686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250612c3a915050565b612d15565b949350505050565b610cf8612b5a565b6010546011546001600160601b0390911690610d159082906143e3565b601155601080546001600160601b03191690556040517f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f190610d58908390613f93565b60405180910390a160405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90610dae90339085906004016144b6565b6020604051808303816000875af1158015610dcd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610df191906144e8565b5050565b610dfd612b5a565b610962613085565b60016001600160a01b0382166000908152600b602052604090205460ff166003811115610e3457610e34613dab565b14158015610e6f575060036001600160a01b0382166000908152600b602052604090205460ff166003811115610e6c57610e6c613dab565b14155b15610e8d576040516303afbb0f60e21b815260040160405180910390fd5b6012546001600160a01b0316610eb65760405163d12d7d8d60e01b815260040160405180910390fd5b6000829003610ed857604051632c2fc94160e01b815260040160405180910390fd5b60006007600085856000818110610ef157610ef1614412565b905060200201358152602001908152602001600020600201600c9054906101000a90046001600160a01b031690506000610f296114c6565b336001600160a01b039182168114925090831614801590610f48575080155b15610f6657604051637dedc72b60e11b815260040160405180910390fd5b6000610f7061395e565b600080876001600160401b03811115610f8b57610f8b61421f565b604051908082528060200260200182016040528015610fbe57816020015b6060815260200190600190039081610fa95790505b5090506000886001600160401b03811115610fdb57610fdb61421f565b60405190808252806020026020018201604052801561101457816020015b61100161395e565b815260200190600190039081610ff95790505b50905060005b8981101561129f578a8a8281811061103457611034614412565b60209081029290920135600081815260078452604090819020815160e08101835281546001600160601b0380821683526001600160a01b03600160601b92839004811698840198909852600184015463ffffffff8116958401959095526001600160401b03600160201b8604166060840152938190048716608083015260029092015492831660a08201529104841660c082018190529199509750918a1690911490506110f45760405163523e0b8360e11b815260040160405180910390fd5b60608501516001600160401b039081161461112257604051633425886760e21b815260040160405180910390fd5b8482828151811061113557611135614412565b6020026020010181905250600a6000878152602001908152602001600020805461115e9061452d565b80601f016020809104026020016040519081016040528092919081815260200182805461118a9061452d565b80156111d75780601f106111ac576101008083540402835291602001916111d7565b820191906000526020600020905b8154815290600101906020018083116111ba57829003601f168201915b50505050508382815181106111ee576111ee614412565b6020908102919091010152845161120e906001600160601b0316856143fa565b600087815260076020908152604080832083815560018101849055600201839055600a9091528120919550611243919061399a565b61124e6005876130dd565b50857fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff86600001518b604051611285929190614494565b60405180910390a28061129781614428565b91505061101a565b50826011546112ae91906143e3565b6011556040516000906112cb908c908c90859087906020016145bc565b6040516020818303038152906040529050886001600160a01b0316638e86139b601260009054906101000a90046001600160a01b03166001600160a01b031663c71249ab60008d6001600160a01b03166348013d7b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015611351573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061137591906146d1565b866040518463ffffffff1660e01b8152600401611394939291906146f2565b600060405180830381865afa1580156113b1573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526113d99190810190614790565b6040518263ffffffff1660e01b81526004016113f59190613b82565b600060405180830381600087803b15801561140f57600080fd5b505af1158015611423573d6000803e3d6000fd5b505060405163a9059cbb60e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150611475908c9088906004016147c4565b6020604051808303816000875af1158015611494573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114b891906144e8565b505050505050505050505050565b6000546001600160a01b031690565b6002336000908152600b602052604090205460ff1660038111156114fb576114fb613dab565b1415801561152d57506003336000908152600b602052604090205460ff16600381111561152a5761152a613dab565b14155b1561154b576040516303afbb0f60e21b815260040160405180910390fd5b6000808061155b848601866149c5565b92509250925060005b83518110156116a35761162184828151811061158257611582614412565b602002602001015184838151811061159c5761159c614412565b6020026020010151608001518584815181106115ba576115ba614412565b6020026020010151604001518685815181106115d8576115d8614412565b602002602001015160c001518786815181106115f6576115f6614412565b60200260200101516000015187878151811061161457611614614412565b60200260200101516130e9565b83818151811061163357611633614412565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7184838151811061166e5761166e614412565b60200260200101516000015133604051611689929190614494565b60405180910390a28061169b81614428565b915050611564565b505050505050565b60008060006116b861337d565b9150915060006116c9836000613540565b90506116d6858284613576565b95945050505050565b6000828152600760205260409020600101548290600160201b90046001600160401b039081161461172357604051633425886760e21b815260040160405180910390fd5b6000838152600760205260409020546117469083906001600160601b0316614469565b600084815260076020526040902080546001600160601b0319166001600160601b0392831617905560115461177d918416906143fa565b6011556040516323b872dd60e01b81523360048201523060248201526001600160601b03831660448201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906323b872dd906064016020604051808303816000875af11580156117fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061181f91906144e8565b50336001600160a01b0316837fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062038460405161185a9190613f93565b60405180910390a3505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146118b05760405163c8bad78d60e01b815260040160405180910390fd5b602081146118d157604051630dfe930960e41b815260040160405180910390fd5b60006118df82840184613f7a565b600081815260076020526040902060010154909150600160201b90046001600160401b039081161461192457604051633425886760e21b815260040160405180910390fd5b6000818152600760205260409020546119479085906001600160601b0316614469565b600082815260076020526040902080546001600160601b0319166001600160601b03929092169190911790556011546119819085906143fa565b601181905550846001600160a01b0316817fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203866040516119c19190613f93565b60405180910390a35050505050565b806001600160a01b0381166119f857604051634e46966960e11b815260040160405180910390fd5b6001600160a01b0383811660009081526008602090815260409182902082516060810184528154948516808252600160a01b9095046001600160601b0316928101929092526001015460ff16151591810191909152903314611a6d5760405163cebf515b60e01b815260040160405180910390fd5b6001600160a01b0380851660009081526008602090815260409091208054909216909155810151601154611aaa916001600160601b0316906143e3565b601181905550826001600160a01b031681602001516001600160601b0316856001600160a01b03167f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f4069833604051611b019190613d7c565b60405180910390a4602081015160405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163a9059cbb91611b5a9187916004016144b6565b6020604051808303816000875af1158015611b79573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b9d91906144e8565b5050505050565b6000828152600760205260409020600101548290600160201b90046001600160401b0390811614611be857604051633425886760e21b815260040160405180910390fd5b6000838152600760205260409020600201548390600160601b90046001600160a01b03163314611c2b5760405163523e0b8360e11b815260040160405180910390fd5b6108fc8363ffffffff161080611c4c5750600d5463ffffffff908116908416115b15611c6a576040516314c237fb60e01b815260040160405180910390fd5b600084815260076020908152604091829020600101805463ffffffff191663ffffffff8716908117909155915191825285917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a250505050565b6001600160a01b03818116600090815260096020526040902054163314611d07576040516333a973d560e11b815260040160405180910390fd5b6001600160a01b0381811660008181526008602090815260408083208054336001600160a01b031980831682179093556009909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000818152600760205260408120600101546106ac9063ffffffff166116ab565b611da3612b5a565b6040516370a0823160e01b81526000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190611df2903090600401613d7c565b602060405180830381865afa158015611e0f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e339190614aa1565b90507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb3360115484611e7391906143e3565b6040518363ffffffff1660e01b8152600401610dae9291906147c4565b611e98612b5a565b8281141580611ea75750600283105b15611ec5576040516367aa603560e11b815260040160405180910390fd5b60005b600454811015611f2657600060048281548110611ee757611ee7614412565b60009182526020808320909101546001600160a01b031682526008905260409020600101805460ff191690555080611f1e81614428565b915050611ec8565b5060005b8381101561208c576000858583818110611f4657611f46614412565b9050602002016020810190611f5b9190613d90565b6001600160a01b03808216600090815260086020526040812080549394509290911690868686818110611f9057611f90614412565b9050602002016020810190611fa59190613d90565b90506001600160a01b0381161580611ff757506001600160a01b03821615801590611fe25750806001600160a01b0316826001600160a01b031614155b8015611ff757506001600160a01b0381811614155b1561201557604051631670f44760e31b815260040160405180910390fd5b600183015460ff161561203b57604051630d5f433160e21b815260040160405180910390fd5b6001838101805460ff191690911790556001600160a01b03818116146120755782546001600160a01b0319166001600160a01b0382161783555b50505050808061208490614428565b915050611f2a565b50612099600485856139d4565b507f056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f848484846040516120cf9493929190614b01565b60405180910390a150505050565b60606000806000806120ed613639565b6000878152600760209081526040808320815160e08101835281546001600160601b0380821683526001600160a01b03600160601b92839004811684880152600185015463ffffffff8116858801526001600160401b03600160201b82041660608601528390048116608085015260029094015490811660a08401520490911660c08201528a8452600a9092528083209051919291636e04ff0d60e01b9161219791602401614b33565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b038381831617835250505050905060008083608001516001600160a01b0316600c600001600b9054906101000a900463ffffffff1663ffffffff16846040516122079190614bda565b60006040518083038160008787f1925050503d8060008114612245576040519150601f19603f3d011682016040523d82523d6000602084013e61224a565b606091505b50915091508161226f57806040516396c3623560e01b8152600401610c119190613b82565b808060200190518101906122839190614bf6565b99509150816122a55760405163865676e360e01b815260040160405180910390fd5b60006122b48b8d8c6000612c3a565b90506122c98582600001518360600151613658565b6060810151608082015160a083015160c0909301519b9e919d509b50909998509650505050505050565b6000818152600760209081526040808320815160e08101835281546001600160601b0380821683526001600160a01b03600160601b928390048116848801908152600186015463ffffffff81168689018190526001600160401b03600160201b83041660608881019182529287900485166080890181905260029099015495861660a089019081529690950490931660c087019081528b8b52600a9099529689208551915198519351945181548b9a8b998a998a998a998a99929893979296929593949390929086906123c59061452d565b80601f01602080910402602001604051908101604052809291908181526020018280546123f19061452d565b801561243e5780601f106124135761010080835404028352916020019161243e565b820191906000526020600020905b81548152906001019060200180831161242157829003601f168201915b505050505095509850985098509850985098509850985050919395975091939597565b6000818152600760205260408120600101546001600160401b03600160201b90910481169190821415906124936114c6565b6001600160a01b0316336001600160a01b03161490508180156124c857508080156124c6575043836001600160401b0316115b155b156124e657604051631f7806af60e31b815260040160405180910390fd5b801580156125155750600084815260076020526040902060020154600160601b90046001600160a01b03163314155b1561253357604051637dedc72b60e11b815260040160405180910390fd5b4381612547576125446032826143fa565b90505b600085815260076020526040902060010180546bffffffffffffffff000000001916600160201b6001600160401b038416021790556125876005866130dd565b506040516001600160401b0382169086907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a35050505050565b60006125cf6114c6565b6001600160a01b0316336001600160a01b0316141580156125fb57506013546001600160a01b03163314155b156126195760405163d48b678b60e01b815260040160405180910390fd5b6126246001436143e3565b600d5460408051924060208401523060601b6001600160601b03191690830152600160201b900460e01b6001600160e01b03191660548201526058016040516020818303038152906040528051906020012060001c90506126c081878787600088888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506130e992505050565b600d8054600160201b900463ffffffff169060046126dd83614c43565b91906101000a81548163ffffffff021916908363ffffffff16021790555050807fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012868660405161274892919063ffffffff9290921682526001600160a01b0316602082015260400190565b60405180910390a295945050505050565b6001600160a01b038281166000908152600860205260409020541633146127935760405163cebf515b60e01b815260040160405180910390fd5b336001600160a01b038216036127bc57604051638c8728c760e01b815260040160405180910390fd5b6001600160a01b03828116600090815260096020526040902054811690821614610df1576001600160a01b0382811660008181526009602052604080822080546001600160a01b0319169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b612841612b5a565b600d5460e082015163ffffffff9182169116101561287257604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600c60010160049054906101000a900463ffffffff1663ffffffff16815250600c60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600e81905550806101200151600f81905550806101400151601260006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601360006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051612b1e9190614c66565b60405180910390a150565b612b31612b5a565b612b3a816136f9565b50565b60006106ac825490565b6000612b53838361379c565b9392505050565b6000546001600160a01b031633146109625760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610c11565b60035460ff16612bf65760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b6044820152606401610c11565b6003805460ff191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b604051612c309190613d7c565b60405180910390a1565b612c836040518060e0016040528060006001600160a01b031681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526007602052604081206001015463ffffffff169080612ca561337d565b915091506000612cb58387613540565b90506000612cc4858385613576565b6040805160e0810182526001600160a01b038d168152602081018c90529081018a90526001600160601b03909116606082015260808101959095525060a084015260c0830152509050949350505050565b60006002805403612d685760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610c11565b600280556020828101516000818152600790925260409091206001015443600160201b9091046001600160401b031611612db557604051633425886760e21b815260040160405180910390fd5b602080840151600090815260078252604090819020815160e08101835281546001600160601b0380821683526001600160a01b03600160601b92839004811696840196909652600184015463ffffffff8116958401959095526001600160401b03600160201b860416606080850191909152948290048616608084015260029093015492831660a083015290910490921660c0830152845190850151612e5c918391613658565b60005a90506000634585e33b60e01b8660400151604051602401612e809190613b82565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b0383818316178352505050509050612ec886608001518460800151836137c6565b94505a612ed590836143e3565b91506000612eec838860a001518960c00151613576565b602080890151600090815260079091526040902054909150612f189082906001600160601b0316614441565b6020888101805160009081526007909252604080832080546001600160601b0319166001600160601b0395861617905590518252902060020154612f5e91839116614469565b60208881018051600090815260078352604080822060020180546001600160601b0319166001600160601b039687161790558b519251825280822080548616600160601b6001600160a01b03958616021790558b5190921681526008909252902054612fd3918391600160a01b900416614469565b6008600089600001516001600160a01b03166001600160a01b0316815260200190815260200160002060000160146101000a8154816001600160601b0302191690836001600160601b0316021790555086600001516001600160a01b031686151588602001517fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6848b6040015160405161306e929190614c75565b60405180910390a450505050506001600255919050565b60035460ff16156130a85760405162461bcd60e51b8152600401610c1190614503565b6003805460ff191660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258612c233390565b6000612b538383613812565b60035460ff161561310c5760405162461bcd60e51b8152600401610c1190614503565b6001600160a01b0385163b613134576040516309ee12d560e01b815260040160405180910390fd5b6108fc8463ffffffff1610806131555750600d5463ffffffff908116908516115b15613173576040516314c237fb60e01b815260040160405180910390fd5b6040518060e00160405280836001600160601b0316815260200160006001600160a01b031681526020018563ffffffff1681526020016001600160401b0380168152602001866001600160a01b0316815260200160006001600160601b03168152602001846001600160a01b03168152506007600088815260200190815260200160002060008201518160000160006101000a8154816001600160601b0302191690836001600160601b03160217905550602082015181600001600c6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160010160006101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010160046101000a8154816001600160401b0302191690836001600160401b03160217905550608082015181600101600c6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060a08201518160020160006101000a8154816001600160601b0302191690836001600160601b0316021790555060c082015181600201600c6101000a8154816001600160a01b0302191690836001600160a01b03160217905550905050816001600160601b031660115461334691906143fa565b6011556000868152600a60209081526040909120825161336892840190613a37565b50613374600587613905565b50505050505050565b6000806000600c600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561340c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134309190614cb3565b509450909250849150508015613454575061344b82426143e3565b8463ffffffff16105b80613460575060008113155b1561346f57600e549550613473565b8095505b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156134d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134f59190614cb3565b509450909250849150508015613519575061351082426143e3565b8463ffffffff16105b80613525575060008113155b1561353457600f549450613538565b8094505b505050509091565b600c5460009061355b90600160901b900461ffff1684614d03565b90508180156135695750803a105b156106ac57503a92915050565b60008061358662013880866143fa565b6135909085614d03565b600c549091506000906135ad9063ffffffff16633b9aca006143fa565b600c549091506000906135d290600160201b900463ffffffff1664e8d4a51000614d03565b85836135e286633b9aca00614d03565b6135ec9190614d03565b6135f69190614d22565b61360091906143fa565b90506b033b2e3c9fd0803ce800000081111561362f5760405163156baa3d60e11b815260040160405180910390fd5b9695505050505050565b32156109625760405163b60ac5db60e01b815260040160405180910390fd5b6001600160a01b03821660009081526008602052604090206001015460ff16613694576040516319f759fb60e31b815260040160405180910390fd5b82516001600160601b03168111156136bf5760405163356680b760e01b815260040160405180910390fd5b816001600160a01b031683602001516001600160a01b0316036136f457604051621af04160e61b815260040160405180910390fd5b505050565b336001600160a01b0382160361374b5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610c11565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106137b3576137b3614412565b9060005260206000200154905092915050565b60005a6113888110156137d857600080fd5b6113888103905084604082048203116137f057600080fd5b50823b6137fc57600080fd5b60008083516020850160008789f1949350505050565b600081815260018301602052604081205480156138fb5760006138366001836143e3565b855490915060009061384a906001906143e3565b90508181146138af57600086600001828154811061386a5761386a614412565b906000526020600020015490508087600001848154811061388d5761388d614412565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806138c0576138c0614d44565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106ac565b60009150506106ac565b6000818152600183016020526040812054612b5390849084908490613956575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106ac565b5060006106ac565b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081019190915290565b5080546139a69061452d565b6000825580601f106139b6575050565b601f016020900490600052602060002090810190612b3a9190613aab565b828054828255906000526020600020908101928215613a27579160200282015b82811115613a275781546001600160a01b0319166001600160a01b038435161782556020909201916001909101906139f4565b50613a33929150613aab565b5090565b828054613a439061452d565b90600052602060002090601f016020900481019282613a655760008555613a27565b82601f10613a7e57805160ff1916838001178555613a27565b82800160010185558215613a27579182015b82811115613a27578251825591602001919060010190613a90565b5b80821115613a335760008155600101613aac565b60008060408385031215613ad357600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613b1a57835183529284019291840191600101613afe565b50909695505050505050565b60005b83811015613b41578181015183820152602001613b29565b83811115613b50576000848401525b50505050565b60008151808452613b6e816020860160208601613b26565b601f01601f19169290920160200192915050565b602081526000612b536020830184613b56565b6001600160a01b03169052565b805163ffffffff1682526020810151613bc3602084018263ffffffff169052565b506040810151613bda604084018262ffffff169052565b506060810151613bf2606084018263ffffffff169052565b506080810151613c09608084018262ffffff169052565b5060a0810151613c1f60a084018261ffff169052565b5060c0810151613c3a60c08401826001600160601b03169052565b5060e0810151613c5260e084018263ffffffff169052565b506101008181015190830152610120808201519083015261014080820151613c7c82850182613b95565b505061016080820151613b5082850182613b95565b600061022080830163ffffffff8751168452602060018060601b038189015116818601526040880151604086015260608801516060860152613cd66080860188613ba2565b6102008501929092528451908190526102408401918086019160005b81811015613d175783516001600160a01b031685529382019392820192600101613cf2565b509298975050505050505050565b80356001600160a01b0381168114613d3c57600080fd5b919050565b60008060408385031215613d5457600080fd5b613d5d83613d25565b9150602083013560048110613d7157600080fd5b809150509250929050565b6001600160a01b0391909116815260200190565b600060208284031215613da257600080fd5b612b5382613d25565b634e487b7160e01b600052602160045260246000fd5b60018110613dd157613dd1613dab565b9052565b602081016106ac8284613dc1565b60008060408385031215613df657600080fd5b82359150613e0660208401613d25565b90509250929050565b60008083601f840112613e2157600080fd5b5081356001600160401b03811115613e3857600080fd5b602083019150836020828501011115613e5057600080fd5b9250929050565b600080600060408486031215613e6c57600080fd5b8335925060208401356001600160401b03811115613e8957600080fd5b613e9586828701613e0f565b9497909650939450505050565b60008083601f840112613eb457600080fd5b5081356001600160401b03811115613ecb57600080fd5b6020830191508360208260051b8501011115613e5057600080fd5b600080600060408486031215613efb57600080fd5b83356001600160401b03811115613f1157600080fd5b613f1d86828701613ea2565b9094509250613f30905060208501613d25565b90509250925092565b60008060208385031215613f4c57600080fd5b82356001600160401b03811115613f6257600080fd5b613f6e85828601613e0f565b90969095509350505050565b600060208284031215613f8c57600080fd5b5035919050565b6001600160601b0391909116815260200190565b80356001600160601b0381168114613d3c57600080fd5b60008060408385031215613fd157600080fd5b82359150613e0660208401613fa7565b60008060008060608587031215613ff757600080fd5b61400085613d25565b93506020850135925060408501356001600160401b0381111561402257600080fd5b61402e87828801613e0f565b95989497509550505050565b6000806040838503121561404d57600080fd5b61405683613d25565b9150613e0660208401613d25565b803563ffffffff81168114613d3c57600080fd5b6000806040838503121561408b57600080fd5b82359150613e0660208401614064565b600080600080604085870312156140b157600080fd5b84356001600160401b03808211156140c857600080fd5b6140d488838901613ea2565b909650945060208701359150808211156140ed57600080fd5b5061402e87828801613ea2565b60a08152600061410d60a0830188613b56565b90508560208301528460408301528360608301528260808301529695505050505050565b6001600160a01b03898116825263ffffffff89166020830152610100604083018190526000916141638483018b613b56565b6001600160601b03998a16606086015297811660808501529590951660a0830152506001600160401b039290921660c083015290931660e090930192909252949350505050565b6000806000806000608086880312156141c257600080fd5b6141cb86613d25565b94506141d960208701614064565b93506141e760408701613d25565b925060608601356001600160401b0381111561420257600080fd5b61420e88828901613e0f565b969995985093965092949392505050565b634e487b7160e01b600052604160045260246000fd5b60405161018081016001600160401b03811182821017156142585761425861421f565b60405290565b60405160e081016001600160401b03811182821017156142585761425861421f565b604051601f8201601f191681016001600160401b03811182821017156142a8576142a861421f565b604052919050565b803562ffffff81168114613d3c57600080fd5b803561ffff81168114613d3c57600080fd5b600061018082840312156142e857600080fd5b6142f0614235565b6142f983614064565b815261430760208401614064565b6020820152614318604084016142b0565b604082015261432960608401614064565b606082015261433a608084016142b0565b608082015261434b60a084016142c3565b60a082015261435c60c08401613fa7565b60c082015261436d60e08401614064565b60e082015261010083810135908201526101208084013590820152610140614396818501613d25565b908201526101606143a8848201613d25565b908201529392505050565b60208101600483106143c7576143c7613dab565b91905290565b634e487b7160e01b600052601160045260246000fd5b6000828210156143f5576143f56143cd565b500390565b6000821982111561440d5761440d6143cd565b500190565b634e487b7160e01b600052603260045260246000fd5b60006001820161443a5761443a6143cd565b5060010190565b60006001600160601b0383811690831681811015614461576144616143cd565b039392505050565b60006001600160601b0382811684821680830382111561448b5761448b6143cd565b01949350505050565b6001600160601b039290921682526001600160a01b0316602082015260400190565b6001600160a01b039290921682526001600160601b0316602082015260400190565b80518015158114613d3c57600080fd5b6000602082840312156144fa57600080fd5b612b53826144d8565b60208082526010908201526f14185d5cd8589b194e881c185d5cd95960821b604082015260600190565b600181811c9082168061454157607f821691505b60208210810361456157634e487b7160e01b600052602260045260246000fd5b50919050565b600081518084526020808501808196508360051b8101915082860160005b858110156145af57828403895261459d848351613b56565b98850198935090840190600101614585565b5091979650505050505050565b606080825281810185905260009060806001600160fb1b038711156145e057600080fd5b8660051b808983870137808501905081810160008152602083878403018188015281895180845260a093508385019150828b01945060005b818110156146ad57855180516001600160601b03168452848101516001600160a01b039081168686015260408083015163ffffffff1690860152898201516001600160401b03168a86015288820151168885015285810151614684878601826001600160601b03169052565b5060c0908101519061469885820183613b95565b50509483019460e09290920191600101614618565b505087810360408901526146c1818a614567565b9c9b505050505050505050505050565b6000602082840312156146e357600080fd5b815160018110612b5357600080fd5b6146fc8185613dc1565b6147096020820184613dc1565b6060604082015260006116d66060830184613b56565b60006001600160401b038211156147385761473861421f565b50601f01601f191660200190565b600082601f83011261475757600080fd5b815161476a6147658261471f565b614280565b81815284602083860101111561477f57600080fd5b610ce8826020830160208701613b26565b6000602082840312156147a257600080fd5b81516001600160401b038111156147b857600080fd5b610ce884828501614746565b6001600160a01b03929092168252602082015260400190565b60006001600160401b038211156147f6576147f661421f565b5060051b60200190565b600082601f83011261481157600080fd5b81356020614821614765836147dd565b82815260e0928302850182019282820191908785111561484057600080fd5b8387015b858110156148ef5781818a03121561485c5760008081fd5b61486461425e565b61486d82613fa7565b815261487a868301613d25565b86820152604061488b818401614064565b908201526060828101356001600160401b03811681146148ab5760008081fd5b9082015260806148bc838201613d25565b9082015260a06148cd838201613fa7565b9082015260c06148de838201613d25565b908201528452928401928101614844565b5090979650505050505050565b600082601f83011261490d57600080fd5b8135602061491d614765836147dd565b82815260059290921b8401810191818101908684111561493c57600080fd5b8286015b848110156149ba5780356001600160401b0381111561495f5760008081fd5b8701603f810189136149715760008081fd5b8481013560406149836147658361471f565b8281528b828486010111156149985760008081fd5b8282850189830137600092810188019290925250845250918301918301614940565b509695505050505050565b6000806000606084860312156149da57600080fd5b83356001600160401b03808211156149f157600080fd5b818601915086601f830112614a0557600080fd5b81356020614a15614765836147dd565b82815260059290921b8401810191818101908a841115614a3457600080fd5b948201945b83861015614a5257853582529482019490820190614a39565b97505087013592505080821115614a6857600080fd5b614a7487838801614800565b93506040860135915080821115614a8a57600080fd5b50614a97868287016148fc565b9150509250925092565b600060208284031215614ab357600080fd5b5051919050565b8183526000602080850194508260005b85811015614af6576001600160a01b03614ae383613d25565b1687529582019590820190600101614aca565b509495945050505050565b604081526000614b15604083018688614aba565b8281036020840152614b28818587614aba565b979650505050505050565b600060208083526000845481600182811c915080831680614b5557607f831692505b8583108103614b7257634e487b7160e01b85526022600452602485fd5b878601838152602001818015614b8f5760018114614ba057614bcb565b60ff19861682528782019650614bcb565b60008b81526020902060005b86811015614bc557815484820152908501908901614bac565b83019750505b50949998505050505050505050565b60008251614bec818460208701613b26565b9190910192915050565b60008060408385031215614c0957600080fd5b614c12836144d8565b915060208301516001600160401b03811115614c2d57600080fd5b614c3985828601614746565b9150509250929050565b600063ffffffff808316818103614c5c57614c5c6143cd565b6001019392505050565b61018081016106ac8284613ba2565b6001600160601b0383168152604060208201819052600090610ce890830184613b56565b805169ffffffffffffffffffff81168114613d3c57600080fd5b600080600080600060a08688031215614ccb57600080fd5b614cd486614c99565b9450602086015193506040860151925060608601519150614cf760808701614c99565b90509295509295909350565b6000816000190483118215151615614d1d57614d1d6143cd565b500290565b600082614d3f57634e487b7160e01b600052601260045260246000fd5b500490565b634e487b7160e01b600052603160045260246000fdfea2646970667358221220be603f5a0a5092046a5f16d5d4e84a057f6772597853c7e76218cc44c52994f264736f6c634300080d0033",
}

// KeeperRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistryMetaData.ABI instead.
var KeeperRegistryABI = KeeperRegistryMetaData.ABI

// KeeperRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistryMetaData.Bin instead.
var KeeperRegistryBin = KeeperRegistryMetaData.Bin

// DeployKeeperRegistry deploys a new Ethereum contract, binding an instance of KeeperRegistry to it.
func DeployKeeperRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkEthFeed common.Address, fastGasFeed common.Address, config Config) (common.Address, *types.Transaction, *KeeperRegistry, error) {
	parsed, err := KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryBin), backend, link, linkEthFeed, fastGasFeed, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry{KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
}

// KeeperRegistry is an auto generated Go binding around an Ethereum contract.
type KeeperRegistry struct {
	KeeperRegistryCaller     // Read-only binding to the contract
	KeeperRegistryTransactor // Write-only binding to the contract
	KeeperRegistryFilterer   // Log filterer for contract events
}

// KeeperRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistrySession struct {
	Contract     *KeeperRegistry   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistryCallerSession struct {
	Contract *KeeperRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// KeeperRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistryTransactorSession struct {
	Contract     *KeeperRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// KeeperRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistryRaw struct {
	Contract *KeeperRegistry // Generic contract binding to access the raw methods on
}

// KeeperRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistryCallerRaw struct {
	Contract *KeeperRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistryTransactorRaw struct {
	Contract *KeeperRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistry creates a new instance of KeeperRegistry, bound to a specific deployed contract.
func NewKeeperRegistry(address common.Address, backend bind.ContractBackend) (*KeeperRegistry, error) {
	contract, err := bindKeeperRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry{KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
}

// NewKeeperRegistryCaller creates a new read-only instance of KeeperRegistry, bound to a specific deployed contract.
func NewKeeperRegistryCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCaller, error) {
	contract, err := bindKeeperRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCaller{contract: contract}, nil
}

// NewKeeperRegistryTransactor creates a new write-only instance of KeeperRegistry, bound to a specific deployed contract.
func NewKeeperRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryTransactor, error) {
	contract, err := bindKeeperRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryTransactor{contract: contract}, nil
}

// NewKeeperRegistryFilterer creates a new log filterer instance of KeeperRegistry, bound to a specific deployed contract.
func NewKeeperRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryFilterer, error) {
	contract, err := bindKeeperRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFilterer{contract: contract}, nil
}

// bindKeeperRegistry binds a generic wrapper to an already deployed contract.
func bindKeeperRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry *KeeperRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry.Contract.KeeperRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry *KeeperRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.KeeperRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry *KeeperRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.KeeperRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry *KeeperRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry *KeeperRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry *KeeperRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.contract.Transact(opts, method, params...)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry *KeeperRegistryCaller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry *KeeperRegistrySession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.FASTGASFEED(&_KeeperRegistry.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry *KeeperRegistryCallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.FASTGASFEED(&_KeeperRegistry.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry *KeeperRegistryCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry *KeeperRegistrySession) LINK() (common.Address, error) {
	return _KeeperRegistry.Contract.LINK(&_KeeperRegistry.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry *KeeperRegistryCallerSession) LINK() (common.Address, error) {
	return _KeeperRegistry.Contract.LINK(&_KeeperRegistry.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry *KeeperRegistryCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry *KeeperRegistrySession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.LINKETHFEED(&_KeeperRegistry.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry *KeeperRegistryCallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.LINKETHFEED(&_KeeperRegistry.CallOpts)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry *KeeperRegistryCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry *KeeperRegistrySession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDs(&_KeeperRegistry.CallOpts, startIndex, maxCount)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry *KeeperRegistryCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDs(&_KeeperRegistry.CallOpts, startIndex, maxCount)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry *KeeperRegistryCaller) GetKeeperInfo(opts *bind.CallOpts, query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getKeeperInfo", query)

	outstruct := new(struct {
		Payee   common.Address
		Active  bool
		Balance *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Payee = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Active = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry *KeeperRegistrySession) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry.Contract.GetKeeperInfo(&_KeeperRegistry.CallOpts, query)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry *KeeperRegistryCallerSession) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry.Contract.GetKeeperInfo(&_KeeperRegistry.CallOpts, query)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry *KeeperRegistryCaller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry *KeeperRegistrySession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMaxPaymentForGas(&_KeeperRegistry.CallOpts, gasLimit)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry *KeeperRegistryCallerSession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMaxPaymentForGas(&_KeeperRegistry.CallOpts, gasLimit)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry *KeeperRegistryCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry *KeeperRegistrySession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry.CallOpts, id)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry *KeeperRegistryCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry.CallOpts, id)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry *KeeperRegistryCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry *KeeperRegistrySession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry.CallOpts, peer)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry *KeeperRegistryCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry.CallOpts, peer)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry *KeeperRegistryCaller) GetState(opts *bind.CallOpts) (struct {
	State   State
	Config  Config
	Keepers []common.Address
}, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		State   State
		Config  Config
		Keepers []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State)).(*State)
	outstruct.Config = *abi.ConvertType(out[1], new(Config)).(*Config)
	outstruct.Keepers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry *KeeperRegistrySession) GetState() (struct {
	State   State
	Config  Config
	Keepers []common.Address
}, error) {
	return _KeeperRegistry.Contract.GetState(&_KeeperRegistry.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry *KeeperRegistryCallerSession) GetState() (struct {
	State   State
	Config  Config
	Keepers []common.Address
}, error) {
	return _KeeperRegistry.Contract.GetState(&_KeeperRegistry.CallOpts)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent)
func (_KeeperRegistry *KeeperRegistryCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getUpkeep", id)

	outstruct := new(struct {
		Target              common.Address
		ExecuteGas          uint32
		CheckData           []byte
		Balance             *big.Int
		LastKeeper          common.Address
		Admin               common.Address
		MaxValidBlocknumber uint64
		AmountSpent         *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ExecuteGas = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.CheckData = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.Balance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LastKeeper = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Admin = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.MaxValidBlocknumber = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.AmountSpent = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent)
func (_KeeperRegistry *KeeperRegistrySession) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}, error) {
	return _KeeperRegistry.Contract.GetUpkeep(&_KeeperRegistry.CallOpts, id)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent)
func (_KeeperRegistry *KeeperRegistryCallerSession) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}, error) {
	return _KeeperRegistry.Contract.GetUpkeep(&_KeeperRegistry.CallOpts, id)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry *KeeperRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry *KeeperRegistrySession) Owner() (common.Address, error) {
	return _KeeperRegistry.Contract.Owner(&_KeeperRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry *KeeperRegistryCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry.Contract.Owner(&_KeeperRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry *KeeperRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry *KeeperRegistrySession) Paused() (bool, error) {
	return _KeeperRegistry.Contract.Paused(&_KeeperRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry *KeeperRegistryCallerSession) Paused() (bool, error) {
	return _KeeperRegistry.Contract.Paused(&_KeeperRegistry.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry *KeeperRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry *KeeperRegistrySession) TypeAndVersion() (string, error) {
	return _KeeperRegistry.Contract.TypeAndVersion(&_KeeperRegistry.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry *KeeperRegistryCallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry.Contract.TypeAndVersion(&_KeeperRegistry.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry *KeeperRegistryCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry *KeeperRegistrySession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepTranscoderVersion(&_KeeperRegistry.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry *KeeperRegistryCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepTranscoderVersion(&_KeeperRegistry.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry *KeeperRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry *KeeperRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptOwnership(&_KeeperRegistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptOwnership(&_KeeperRegistry.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptPayeeship", keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry *KeeperRegistrySession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptPayeeship(&_KeeperRegistry.TransactOpts, keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptPayeeship(&_KeeperRegistry.TransactOpts, keeper)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "addFunds", id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry *KeeperRegistrySession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AddFunds(&_KeeperRegistry.TransactOpts, id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AddFunds(&_KeeperRegistry.TransactOpts, id, amount)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "cancelUpkeep", id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry *KeeperRegistrySession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CancelUpkeep(&_KeeperRegistry.TransactOpts, id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CancelUpkeep(&_KeeperRegistry.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry *KeeperRegistryTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "checkUpkeep", id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry *KeeperRegistrySession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep(&_KeeperRegistry.TransactOpts, id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry *KeeperRegistryTransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep(&_KeeperRegistry.TransactOpts, id, from)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry *KeeperRegistrySession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.MigrateUpkeeps(&_KeeperRegistry.TransactOpts, ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.MigrateUpkeeps(&_KeeperRegistry.TransactOpts, ids, destination)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry *KeeperRegistrySession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.OnTokenTransfer(&_KeeperRegistry.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.OnTokenTransfer(&_KeeperRegistry.TransactOpts, sender, amount, data)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry *KeeperRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry *KeeperRegistrySession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Pause(&_KeeperRegistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Pause(&_KeeperRegistry.TransactOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry *KeeperRegistryTransactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "performUpkeep", id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry *KeeperRegistrySession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.PerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry *KeeperRegistryTransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.PerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry *KeeperRegistrySession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.ReceiveUpkeeps(&_KeeperRegistry.TransactOpts, encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.ReceiveUpkeeps(&_KeeperRegistry.TransactOpts, encodedUpkeeps)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry *KeeperRegistryTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "recoverFunds")
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry *KeeperRegistrySession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RecoverFunds(&_KeeperRegistry.TransactOpts)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RecoverFunds(&_KeeperRegistry.TransactOpts)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry *KeeperRegistryTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry *KeeperRegistrySession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RegisterUpkeep(&_KeeperRegistry.TransactOpts, target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry *KeeperRegistryTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RegisterUpkeep(&_KeeperRegistry.TransactOpts, target, gasLimit, admin, checkData)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) SetConfig(opts *bind.TransactOpts, config Config) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setConfig", config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry *KeeperRegistrySession) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, config)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setKeepers", keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry *KeeperRegistrySession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetKeepers(&_KeeperRegistry.TransactOpts, keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetKeepers(&_KeeperRegistry.TransactOpts, keepers, payees)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry *KeeperRegistrySession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry.TransactOpts, peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry.TransactOpts, peer, permission)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry *KeeperRegistrySession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepGasLimit(&_KeeperRegistry.TransactOpts, id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepGasLimit(&_KeeperRegistry.TransactOpts, id, gasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry *KeeperRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferOwnership(&_KeeperRegistry.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferOwnership(&_KeeperRegistry.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry *KeeperRegistrySession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferPayeeship(&_KeeperRegistry.TransactOpts, keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferPayeeship(&_KeeperRegistry.TransactOpts, keeper, proposed)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry *KeeperRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry *KeeperRegistrySession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Unpause(&_KeeperRegistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Unpause(&_KeeperRegistry.TransactOpts)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawFunds", id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry *KeeperRegistrySession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawFunds(&_KeeperRegistry.TransactOpts, id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawFunds(&_KeeperRegistry.TransactOpts, id, to)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawOwnerFunds")
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry *KeeperRegistrySession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawOwnerFunds(&_KeeperRegistry.TransactOpts)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawOwnerFunds(&_KeeperRegistry.TransactOpts)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawPayment", from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry *KeeperRegistrySession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawPayment(&_KeeperRegistry.TransactOpts, from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawPayment(&_KeeperRegistry.TransactOpts, from, to)
}

// KeeperRegistryConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the KeeperRegistry contract.
type KeeperRegistryConfigSetIterator struct {
	Event *KeeperRegistryConfigSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryConfigSet)
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
		it.Event = new(KeeperRegistryConfigSet)
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
func (it *KeeperRegistryConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryConfigSet represents a ConfigSet event raised by the KeeperRegistry contract.
type KeeperRegistryConfigSet struct {
	Config Config
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryConfigSetIterator{contract: _KeeperRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error) {
	event := new(KeeperRegistryConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryFundsAddedIterator is returned from FilterFundsAdded and is used to iterate over the raw logs and unpacked data for FundsAdded events raised by the KeeperRegistry contract.
type KeeperRegistryFundsAddedIterator struct {
	Event *KeeperRegistryFundsAdded // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryFundsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryFundsAdded)
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
		it.Event = new(KeeperRegistryFundsAdded)
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
func (it *KeeperRegistryFundsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryFundsAdded represents a FundsAdded event raised by the KeeperRegistry contract.
type KeeperRegistryFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsAdded is a free log retrieval operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFundsAddedIterator{contract: _KeeperRegistry.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

// WatchFundsAdded is a free log subscription operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryFundsAdded)
				if err := _KeeperRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryFundsAdded, error) {
	event := new(KeeperRegistryFundsAdded)
	if err := _KeeperRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryFundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the KeeperRegistry contract.
type KeeperRegistryFundsWithdrawnIterator struct {
	Event *KeeperRegistryFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryFundsWithdrawn)
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
		it.Event = new(KeeperRegistryFundsWithdrawn)
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
func (it *KeeperRegistryFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryFundsWithdrawn represents a FundsWithdrawn event raised by the KeeperRegistry contract.
type KeeperRegistryFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFundsWithdrawnIterator{contract: _KeeperRegistry.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryFundsWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryFundsWithdrawn, error) {
	event := new(KeeperRegistryFundsWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryKeepersUpdatedIterator is returned from FilterKeepersUpdated and is used to iterate over the raw logs and unpacked data for KeepersUpdated events raised by the KeeperRegistry contract.
type KeeperRegistryKeepersUpdatedIterator struct {
	Event *KeeperRegistryKeepersUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryKeepersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryKeepersUpdated)
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
		it.Event = new(KeeperRegistryKeepersUpdated)
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
func (it *KeeperRegistryKeepersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryKeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryKeepersUpdated represents a KeepersUpdated event raised by the KeeperRegistry contract.
type KeeperRegistryKeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKeepersUpdated is a free log retrieval operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryKeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryKeepersUpdatedIterator{contract: _KeeperRegistry.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

// WatchKeepersUpdated is a free log subscription operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryKeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryKeepersUpdated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistryKeepersUpdated, error) {
	event := new(KeeperRegistryKeepersUpdated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryOwnerFundsWithdrawnIterator is returned from FilterOwnerFundsWithdrawn and is used to iterate over the raw logs and unpacked data for OwnerFundsWithdrawn events raised by the KeeperRegistry contract.
type KeeperRegistryOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryOwnerFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistryOwnerFundsWithdrawn)
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
func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryOwnerFundsWithdrawn represents a OwnerFundsWithdrawn event raised by the KeeperRegistry contract.
type KeeperRegistryOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerFundsWithdrawn is a free log retrieval operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnerFundsWithdrawnIterator{contract: _KeeperRegistry.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchOwnerFundsWithdrawn is a free log subscription operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryOwnerFundsWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryOwnerFundsWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistry contract.
type KeeperRegistryOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryOwnershipTransferRequested)
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
func (it *KeeperRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistry contract.
type KeeperRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnershipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryOwnershipTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryOwnershipTransferRequested, error) {
	event := new(KeeperRegistryOwnershipTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistry contract.
type KeeperRegistryOwnershipTransferredIterator struct {
	Event *KeeperRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnershipTransferred)
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
		it.Event = new(KeeperRegistryOwnershipTransferred)
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
func (it *KeeperRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistry contract.
type KeeperRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnershipTransferredIterator{contract: _KeeperRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryOwnershipTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryOwnershipTransferred, error) {
	event := new(KeeperRegistryOwnershipTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the KeeperRegistry contract.
type KeeperRegistryPausedIterator struct {
	Event *KeeperRegistryPaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPaused)
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
		it.Event = new(KeeperRegistryPaused)
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
func (it *KeeperRegistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryPaused represents a Paused event raised by the KeeperRegistry contract.
type KeeperRegistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryPausedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPausedIterator{contract: _KeeperRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryPaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParsePaused(log types.Log) (*KeeperRegistryPaused, error) {
	event := new(KeeperRegistryPaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the KeeperRegistry contract.
type KeeperRegistryPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryPayeeshipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistryPayeeshipTransferRequested)
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
func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the KeeperRegistry contract.
type KeeperRegistryPayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryPayeeshipTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryPayeeshipTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the KeeperRegistry contract.
type KeeperRegistryPayeeshipTransferredIterator struct {
	Event *KeeperRegistryPayeeshipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeeshipTransferred)
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
		it.Event = new(KeeperRegistryPayeeshipTransferred)
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
func (it *KeeperRegistryPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryPayeeshipTransferred represents a PayeeshipTransferred event raised by the KeeperRegistry contract.
type KeeperRegistryPayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferredIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryPayeeshipTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryPayeeshipTransferred, error) {
	event := new(KeeperRegistryPayeeshipTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryPaymentWithdrawnIterator is returned from FilterPaymentWithdrawn and is used to iterate over the raw logs and unpacked data for PaymentWithdrawn events raised by the KeeperRegistry contract.
type KeeperRegistryPaymentWithdrawnIterator struct {
	Event *KeeperRegistryPaymentWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryPaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPaymentWithdrawn)
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
		it.Event = new(KeeperRegistryPaymentWithdrawn)
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
func (it *KeeperRegistryPaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryPaymentWithdrawn represents a PaymentWithdrawn event raised by the KeeperRegistry contract.
type KeeperRegistryPaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPaymentWithdrawnIterator{contract: _KeeperRegistry.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryPaymentWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryPaymentWithdrawn, error) {
	event := new(KeeperRegistryPaymentWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the KeeperRegistry contract.
type KeeperRegistryUnpausedIterator struct {
	Event *KeeperRegistryUnpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUnpaused)
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
		it.Event = new(KeeperRegistryUnpaused)
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
func (it *KeeperRegistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUnpaused represents a Unpaused event raised by the KeeperRegistry contract.
type KeeperRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUnpausedIterator{contract: _KeeperRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUnpaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryUnpaused, error) {
	event := new(KeeperRegistryUnpaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUpkeepCanceledIterator is returned from FilterUpkeepCanceled and is used to iterate over the raw logs and unpacked data for UpkeepCanceled events raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepCanceledIterator struct {
	Event *KeeperRegistryUpkeepCanceled // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUpkeepCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepCanceled)
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
		it.Event = new(KeeperRegistryUpkeepCanceled)
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
func (it *KeeperRegistryUpkeepCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUpkeepCanceled represents a UpkeepCanceled event raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCanceled is a free log retrieval operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepCanceledIterator{contract: _KeeperRegistry.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

// WatchUpkeepCanceled is a free log subscription operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUpkeepCanceled)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryUpkeepCanceled, error) {
	event := new(KeeperRegistryUpkeepCanceled)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUpkeepGasLimitSetIterator is returned from FilterUpkeepGasLimitSet and is used to iterate over the raw logs and unpacked data for UpkeepGasLimitSet events raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryUpkeepGasLimitSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUpkeepGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistryUpkeepGasLimitSet)
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
func (it *KeeperRegistryUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUpkeepGasLimitSet represents a UpkeepGasLimitSet event raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpkeepGasLimitSet is a free log retrieval operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepGasLimitSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepGasLimitSet is a free log subscription operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUpkeepGasLimitSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryUpkeepGasLimitSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUpkeepMigratedIterator is returned from FilterUpkeepMigrated and is used to iterate over the raw logs and unpacked data for UpkeepMigrated events raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepMigratedIterator struct {
	Event *KeeperRegistryUpkeepMigrated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUpkeepMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepMigrated)
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
		it.Event = new(KeeperRegistryUpkeepMigrated)
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
func (it *KeeperRegistryUpkeepMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUpkeepMigrated represents a UpkeepMigrated event raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepMigrated is a free log retrieval operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepMigratedIterator{contract: _KeeperRegistry.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

// WatchUpkeepMigrated is a free log subscription operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUpkeepMigrated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryUpkeepMigrated, error) {
	event := new(KeeperRegistryUpkeepMigrated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUpkeepPerformedIterator is returned from FilterUpkeepPerformed and is used to iterate over the raw logs and unpacked data for UpkeepPerformed events raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepPerformedIterator struct {
	Event *KeeperRegistryUpkeepPerformed // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUpkeepPerformedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepPerformed)
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
		it.Event = new(KeeperRegistryUpkeepPerformed)
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
func (it *KeeperRegistryUpkeepPerformedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUpkeepPerformed represents a UpkeepPerformed event raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepPerformed struct {
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
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryUpkeepPerformedIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPerformedIterator{contract: _KeeperRegistry.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUpkeepPerformed)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryUpkeepPerformed, error) {
	event := new(KeeperRegistryUpkeepPerformed)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUpkeepReceivedIterator is returned from FilterUpkeepReceived and is used to iterate over the raw logs and unpacked data for UpkeepReceived events raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepReceivedIterator struct {
	Event *KeeperRegistryUpkeepReceived // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUpkeepReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepReceived)
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
		it.Event = new(KeeperRegistryUpkeepReceived)
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
func (it *KeeperRegistryUpkeepReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUpkeepReceived represents a UpkeepReceived event raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpkeepReceived is a free log retrieval operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepReceivedIterator{contract: _KeeperRegistry.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

// WatchUpkeepReceived is a free log subscription operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUpkeepReceived)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryUpkeepReceived, error) {
	event := new(KeeperRegistryUpkeepReceived)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistryUpkeepRegisteredIterator is returned from FilterUpkeepRegistered and is used to iterate over the raw logs and unpacked data for UpkeepRegistered events raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepRegisteredIterator struct {
	Event *KeeperRegistryUpkeepRegistered // Event containing the contract specifics and raw log

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
func (it *KeeperRegistryUpkeepRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepRegistered)
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
		it.Event = new(KeeperRegistryUpkeepRegistered)
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
func (it *KeeperRegistryUpkeepRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistryUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistryUpkeepRegistered represents a UpkeepRegistered event raised by the KeeperRegistry contract.
type KeeperRegistryUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpkeepRegistered is a free log retrieval operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepRegisteredIterator{contract: _KeeperRegistry.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

// WatchUpkeepRegistered is a free log subscription operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistryUpkeepRegistered)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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
func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryUpkeepRegistered, error) {
	event := new(KeeperRegistryUpkeepRegistered)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
