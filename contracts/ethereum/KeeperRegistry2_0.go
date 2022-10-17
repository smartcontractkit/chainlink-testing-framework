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

// KeeperRegistry20MetaData contains all meta data concerning the KeeperRegistry20 contract.
var KeeperRegistry20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.PaymentModel\",\"name\":\"paymentModel\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"registryGasOverhead\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"keeperRegistryLogic\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"ARB_NITRO_ORACLE\",\"outputs\":[{\"internalType\":\"contractArbGasInfo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"KEEPER_REGISTRY_LOGIC\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_FEE_DATA_PADDING\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_INPUT_DATA\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPTIMISM_ORACLE\",\"outputs\":[{\"internalType\":\"contractOVM_GasPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAYMENT_MODEL\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REGISTRY_GAS_OVERHEAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getKeeperInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"}],\"internalType\":\"structState\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"lastKeeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setKeepers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x61020060405260486101808181529062003e546101a03980516200002c91600491602090910190620005e6565b50604051806101400160405280610110815260200162003e9c610110913980516200006091600591602090910190620005e6565b507f420000000000000000000000000000000000000f00000000000000000000000060e0526c6c00000000000000000000000061010052348015620000a457600080fd5b5060405162003fac38038062003fac833981016040819052620000c791620006fd565b86868686863380600081620001235760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000156576200015681620001e7565b5050600160029081556003805460ff191690558691508111156200017e576200017e620009f6565b610120816002811115620001965762000196620009f6565b60f81b905250610140939093526001600160601b0319606092831b811660805290821b811660a05291811b821660c05284901b166101605250620001da8162000293565b5050505050505062000a0c565b6001600160a01b038116331415620002425760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200011a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6200029d62000588565b60105460e082015163ffffffff91821691161015620002cf57604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600f60010160049054906101000a900463ffffffff1663ffffffff16815250600f60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151601181905550806101200151601281905550806101400151601560006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601660006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325816040516200057d919062000882565b60405180910390a150565b6000546001600160a01b03163314620005e45760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016200011a565b565b828054620005f490620009b9565b90600052602060002090601f01602090048101928262000618576000855562000663565b82601f106200063357805160ff191683800117855562000663565b8280016001018555821562000663579182015b828111156200066357825182559160200191906001019062000646565b506200067192915062000675565b5090565b5b8082111562000671576000815560010162000676565b80516001600160a01b0381168114620006a457600080fd5b919050565b805161ffff81168114620006a457600080fd5b805162ffffff81168114620006a457600080fd5b805163ffffffff81168114620006a457600080fd5b80516001600160601b0381168114620006a457600080fd5b60008060008060008060008789036102408112156200071b57600080fd5b8851600381106200072b57600080fd5b60208a015190985096506200074360408a016200068c565b95506200075360608a016200068c565b94506200076360808a016200068c565b93506200077360a08a016200068c565b92506101808060bf19830112156200078a57600080fd5b6200079462000981565b9150620007a460c08b01620006d0565b8252620007b460e08b01620006d0565b6020830152610100620007c9818c01620006bc565b6040840152610120620007de818d01620006d0565b6060850152610140620007f3818e01620006bc565b608086015261016062000808818f01620006a9565b60a08701526200081a858f01620006e5565b60c08701526200082e6101a08f01620006d0565b60e08701526101c08e0151848701526101e08e015183870152620008566102008f016200068c565b82870152620008696102208f016200068c565b8187015250505050508091505092959891949750929550565b815163ffffffff16815261018081016020830151620008a9602084018263ffffffff169052565b506040830151620008c1604084018262ffffff169052565b506060830151620008da606084018263ffffffff169052565b506080830151620008f2608084018262ffffff169052565b5060a08301516200090960a084018261ffff169052565b5060c08301516200092560c08401826001600160601b03169052565b5060e08301516200093e60e084018263ffffffff169052565b5061010083810151908301526101208084015190830152610140808401516001600160a01b03908116918401919091526101609384015116929091019190915290565b60405161018081016001600160401b0381118282101715620009b357634e487b7160e01b600052604160045260246000fd5b60405290565b600181811c90821680620009ce57607f821691505b60208210811415620009f057634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160601c6101005160601c6101205160f81c610140516101605160601c61338e62000ac6600039600081816102b6015261096f0152600081816104f101526120310152600081816106540152818161208401526122320152600081816105ac015261226a0152600081816105e001526121a101526000818161049b0152611dbe0152600081816107ac0152611e920152600081816103ae01526111f8015261338e6000f3fe6080604052600436106102575760003560e01c80638da5cb5b1161013a578063b657bc9c116100b1578063b657bc9c14610809578063b79550be14610474578063b7fdb43614610829578063be30817714610849578063c41b813a1461085e578063c7c3a19a1461088f578063c8048022146107ee578063da5c6741146108c4578063eb5dcd6c14610764578063ef47a0ce146108e4578063f2fde38b14610904578063faa3e9961461092457610266565b80638da5cb5b146106835780638e86139b146106a157806393f0c1fc146106bc578063948108f7146106f45780639fab43861461070f578063a0ad00cf1461072f578063a4c0ed3614610744578063a710b22114610764578063a72aa27e1461077f578063ad1783611461079a578063b121e147146107ce578063b148ab6b146107ee57610266565b80635165f2f5116101ce5780635165f2f5146105215780635c975abb14610541578063744bfe611461038157806379ba5097146105655780637bbaf1ea1461057a5780637d9b97e0146104745780637f37618e1461059a5780638456cb5914610474578063850cce34146105ce57806385c1b0ba146106025780638765ecbe146106225780638811cbe81461064257610266565b806306e3b6321461026e5780630852c7c9146102a4578063181f5a77146102f05780631865c57d1461033d578063187256e8146103615780631a2af011146103815780631b6b6d231461039c5780631e12b8a5146103d05780633f4ba83a146104745780634584a4191461048957806348013d7b146104bd5780635077b210146104df57610266565b366102665761026461096a565b005b61026461096a565b34801561027a57600080fd5b5061028e610289366004612cab565b610995565b60405161029b9190612fd4565b60405180910390f35b3480156102b057600080fd5b506102d87f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161029b565b3480156102fc57600080fd5b506103306040518060400160405280601481526020017304b6565706572526567697374727920312e332e360641b81525081565b60405161029b9190613047565b34801561034957600080fd5b50610352610a77565b60405161029b939291906130e2565b34801561036d57600080fd5b5061026461037c366004612931565b610cd1565b34801561038d57600080fd5b5061026461037c366004612c3d565b3480156103a857600080fd5b506102d87f000000000000000000000000000000000000000000000000000000000000000081565b3480156103dc57600080fd5b506104466103eb3660046128e3565b6001600160a01b039081166000908152600a602090815260409182902082516060810184528154948516808252600160a01b9095046001600160601b031692810183905260019091015460ff16151592018290529192909190565b604080516001600160a01b03909416845291151560208401526001600160601b03169082015260600161029b565b34801561048057600080fd5b50610264610cdd565b34801561049557600080fd5b506102d87f000000000000000000000000000000000000000000000000000000000000000081565b3480156104c957600080fd5b506104d2600181565b60405161029b91906130bf565b3480156104eb57600080fd5b506105137f000000000000000000000000000000000000000000000000000000000000000081565b60405190815260200161029b565b34801561052d57600080fd5b5061026461053c366004612c0b565b610ce5565b34801561054d57600080fd5b5060035460ff165b604051901515815260200161029b565b34801561057157600080fd5b50610264610e03565b34801561058657600080fd5b50610555610595366004612c60565b610eb2565b3480156105a657600080fd5b506102d87f000000000000000000000000000000000000000000000000000000000000000081565b3480156105da57600080fd5b506102d87f000000000000000000000000000000000000000000000000000000000000000081565b34801561060e57600080fd5b5061026461061d366004612a99565b610f10565b34801561062e57600080fd5b5061026461063d366004612c0b565b610f1d565b34801561064e57600080fd5b506106767f000000000000000000000000000000000000000000000000000000000000000081565b60405161029b91906130ab565b34801561068f57600080fd5b506000546001600160a01b03166102d8565b3480156106ad57600080fd5b5061026461037c366004612aec565b3480156106c857600080fd5b506106dc6106d7366004612c0b565b611042565b6040516001600160601b03909116815260200161029b565b34801561070057600080fd5b5061026461037c366004612cf0565b34801561071b57600080fd5b5061026461072a366004612c60565b611060565b34801561073b57600080fd5b5061033061115f565b34801561075057600080fd5b5061026461075f36600461296c565b6111ed565b34801561077057600080fd5b5061026461037c3660046128fe565b34801561078b57600080fd5b5061026461037c366004612ccd565b3480156107a657600080fd5b506102d87f000000000000000000000000000000000000000000000000000000000000000081565b3480156107da57600080fd5b506102646107e93660046128e3565b611357565b3480156107fa57600080fd5b506102646107e9366004612c0b565b34801561081557600080fd5b506106dc610824366004612c0b565b611362565b34801561083557600080fd5b50610264610844366004612a3a565b611383565b34801561085557600080fd5b50610330611391565b34801561086a57600080fd5b5061087e610879366004612c3d565b61139e565b60405161029b95949392919061305a565b34801561089b57600080fd5b506108af6108aa366004612c0b565b6113c0565b60405161029b99989796959493929190612f53565b3480156108d057600080fd5b506105136108df3660046129c5565b611556565b3480156108f057600080fd5b506102646108ff366004612b2d565b611569565b34801561091057600080fd5b5061026461091f3660046128e3565b611859565b34801561093057600080fd5b5061095d61093f3660046128e3565b6001600160a01b03166000908152600e602052604090205460ff1690565b60405161029b9190613091565b6109937f000000000000000000000000000000000000000000000000000000000000000061186a565b565b606060006109a3600761188e565b90508084106109c557604051631390f2a160e01b815260040160405180910390fd5b826109d7576109d48482613255565b92505b6000836001600160401b038111156109f1576109f1613342565b604051908082528060200260200182016040528015610a1a578160200160208202803683370190505b50905060005b84811015610a6c57610a3d610a3582886131d1565b600790611898565b828281518110610a4f57610a4f61332c565b602090810291909101015280610a64816132cf565b915050610a20565b509150505b92915050565b6040805160808101825260008082526020820181905291810182905260608101919091526040805161018081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101919091526040805161012081018252600f5463ffffffff8082168352600160201b808304821660208086019190915262ffffff600160401b8504811686880152600160581b85048416606087810191909152600160781b8604909116608087015261ffff600160901b86041660a08701526001600160601b03600160a01b909504851660c087015260105480851660e088015292909204909216610100850181905287526013549092169086015260145492850192909252610bba600761188e565b606080860191909152815163ffffffff908116855260208084015182168187015260408085015162ffffff90811682890152858501518416948801949094526080808601519094169387019390935260a08085015161ffff169087015260c0808501516001600160601b03169087015260e080850151909216918601919091526011546101008601526012546101208601526015546001600160a01b0390811661014087015260165416610160860152600680548351818402810184019094528084528793879390918391830182828015610cbe57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610ca0575b5050505050905093509350935050909192565b610cd961096a565b5050565b61099361096a565b60008181526009602090815260409182902082516101008101845281546001600160601b0380821683526001600160a01b03600160601b92839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff8082166080850152600160201b82041660a0840152600160401b810490911660c083015260ff600160e01b90910416151560e0820152610d8a816118ab565b8060e00151610dac576040516306e229e160e21b815260040160405180910390fd5b6000828152600960205260409020600201805460ff60e01b19169055610dd360078361190c565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b6001546001600160a01b03163314610e5b5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610ebc611918565b610f08610f03338686868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506001925061195e915050565b611a29565b949350505050565b610f1861096a565b505050565b60008181526009602090815260409182902082516101008101845281546001600160601b0380821683526001600160a01b03600160601b92839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff8082166080850152600160201b82041660a0840152600160401b810490911660c083015260ff600160e01b90910416151560e0820152610fc2816118ab565b8060e0015115610fe557604051631452db0960e21b815260040160405180910390fd5b6000828152600960205260409020600201805460ff60e01b1916600160e01b179055611012600783611d7f565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b600080600061104f611d8b565b91509150610f088483836000611f6c565b60008381526009602090815260409182902082516101008101845281546001600160601b0380821683526001600160a01b03600160601b92839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff8082166080850152600160201b82041660a0840152600160401b810490911660c083015260ff600160e01b90910416151560e0820152611105816118ab565b6000848152600d6020526040902061111e908484612738565b50837f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf8484604051611151929190613018565b60405180910390a250505050565b6005805461116c90613294565b80601f016020809104026020016040519081016040528092919081815260200182805461119890613294565b80156111e55780601f106111ba576101008083540402835291602001916111e5565b820191906000526020600020905b8154815290600101906020018083116111c857829003601f168201915b505050505081565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146112365760405163c8bad78d60e01b815260040160405180910390fd5b6020811461125757604051630dfe930960e41b815260040160405180910390fd5b600061126582840184612c0b565b600081815260096020526040902060020154909150600160201b900463ffffffff908116146112a757604051634e0041d160e11b815260040160405180910390fd5b6000818152600960205260409020546112ca9085906001600160601b03166131e9565b600082815260096020526040902080546001600160601b0319166001600160601b03929092169190911790556014546113049085906131d1565b6014556040516001600160601b03851681526001600160a01b0386169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b61135f61096a565b50565b600081815260096020526040812060020154610a719063ffffffff16611042565b61138b61096a565b50505050565b6004805461116c90613294565b60606000806000806113ae6123aa565b6113b661096a565b9295509295909350565b600081815260096020908152604080832081516101008101835281546001600160601b038082168352600160601b918290046001600160a01b039081168488019081526001860154928316858801908152939092048116606080860191825260029096015463ffffffff80821660808801819052600160201b830490911660a08801908152600160401b830490941660c08801819052600160e01b90920460ff16151560e088019081528c8c52600d909a52978a2086519451925193519551995181548c9b999a8c9a8b9a8b9a8b9a8b9a8b9a93999895979596919593949093909187906114ad90613294565b80601f01602080910402602001604051908101604052809291908181526020018280546114d990613294565b80156115265780601f106114fb57610100808354040283529160200191611526565b820191906000526020600020905b81548152906001019060200180831161150957829003601f168201915b505050505096508263ffffffff169250995099509950995099509950995099509950509193959799909294969850565b600061156061096a565b95945050505050565b6115716123c9565b60105460e082015163ffffffff918216911610156115a257604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600f60010160049054906101000a900463ffffffff1663ffffffff16815250600f60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151601181905550806101200151601281905550806101400151601560006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601660006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de3258160405161184e91906130d3565b60405180910390a150565b6118616123c9565b61135f8161241c565b3660008037600080366000845af43d6000803e808015611889573d6000f35b3d6000fd5b6000610a71825490565b60006118a483836124c0565b9392505050565b80606001516001600160a01b0316336001600160a01b0316146118e15760405163523e0b8360e11b815260040160405180910390fd5b60a081015163ffffffff9081161461135f57604051634e0041d160e11b815260040160405180910390fd5b60006118a483836124ea565b60035460ff16156109935760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b6044820152606401610e52565b6119a76040518060e0016040528060006001600160a01b031681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526009602052604081206002015463ffffffff1690806119c9611d8b565b9150915060006119db84848489611f6c565b6040805160e0810182526001600160a01b03909b168b5260208b0199909952978901969096526001600160601b039096166060880152608087019190915260a086015250505060c082015290565b6000600280541415611a7d5760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610e52565b600280805560208084015160009081526009825260409081902081516101008101835281546001600160601b0380821683526001600160a01b03600160601b92839004811696840196909652600184015490811694830194909452909204831660608301529092015463ffffffff8082166080850152600160201b82041660a08401819052600160401b820490921660c084015260ff600160e01b90910416151560e08301524310611b4257604051634e0041d160e11b815260040160405180910390fd5b611b558184600001518560600151612539565b60005a90506000634585e33b60e01b8560400151604051602401611b799190613047565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b0383818316178352505050509050611bc185608001518460c00151836125f9565b93505a611bce9083613255565b91506000611be7838760a001518860c001516001611f6c565b602080880151600090815260099091526040902054909150611c139082906001600160601b031661326c565b6020878101805160009081526009909252604080832080546001600160601b0319166001600160601b0395861617905590518252902060010154611c59918391166131e9565b60208781018051600090815260098352604080822060010180546001600160601b0319166001600160601b039687161790558a519251825280822080548616600160601b6001600160a01b03958616021790558a519092168152600a909252902054611cce918391600160a01b9004166131e9565b600a600088600001516001600160a01b03166001600160a01b0316815260200190815260200160002060000160146101000a8154816001600160601b0302191690836001600160601b0316021790555085600001516001600160a01b031685151587602001517fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6848a60400151604051611d69929190613176565b60405180910390a4505050506001600255919050565b60006118a48383612645565b6000806000600f600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015611e1557600080fd5b505afa158015611e29573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e4d9190612d13565b509450909250849150508015611e715750611e688242613255565b8463ffffffff16105b80611e7d575060008113155b15611e8c576011549550611e90565b8095505b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015611ee957600080fd5b505afa158015611efd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f219190612d13565b509450909250849150508015611f455750611f3c8242613255565b8463ffffffff16105b80611f51575060008113155b15611f60576012549450611f64565b8094505b505050509091565b6040805161012081018252600f5463ffffffff8082168352600160201b80830482166020850152600160401b830462ffffff90811695850195909552600160581b830482166060850152600160781b83049094166080840152600160901b820461ffff1660a08401819052600160a01b9092046001600160601b031660c084015260105480821660e08501529390930490921661010082015260009182906120149087613236565b90508380156120225750803a105b1561202a57503a5b60006120567f0000000000000000000000000000000000000000000000000000000000000000896131d1565b6120609083613236565b835190915060009061207c9063ffffffff16633b9aca006131d1565b9050600060027f000000000000000000000000000000000000000000000000000000000000000060028111156120b4576120b4613300565b141561222e5760408051600081526020810190915287156120fc5760003660046040516020016120e693929190612ea8565b604051602081830303815290604052905061218a565b6005805461210990613294565b80601f016020809104026020016040519081016040528092919081815260200182805461213590613294565b80156121825780601f1061215757610100808354040283529160200191612182565b820191906000526020600020905b81548152906001019060200180831161216557829003601f168201915b505050505090505b6040516324ca470760e11b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906349948e0e906121d6908490600401613047565b60206040518083038186803b1580156121ee57600080fd5b505afa158015612202573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122269190612c24565b9150506122fc565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561226257612262613300565b14156122fc577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156122c157600080fd5b505afa1580156122d5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122f99190612c24565b90505b8661231857808560a0015161ffff166123159190613236565b90505b6000856020015163ffffffff1664e8d4a510006123359190613236565b898461234185886131d1565b61234f90633b9aca00613236565b6123599190613236565b6123639190613214565b61236d91906131d1565b90506b033b2e3c9fd0803ce800000081111561239c5760405163156baa3d60e11b815260040160405180910390fd5b9a9950505050505050505050565b32156109935760405163b60ac5db60e01b815260040160405180910390fd5b6000546001600160a01b031633146109935760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610e52565b6001600160a01b03811633141561246f5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610e52565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106124d7576124d761332c565b9060005260206000200154905092915050565b600081815260018301602052604081205461253157508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610a71565b506000610a71565b8260e001511561255c57604051631452db0960e21b815260040160405180910390fd5b6001600160a01b0382166000908152600a602052604090206001015460ff16612598576040516319f759fb60e31b815260040160405180910390fd5b82516001600160601b03168111156125c35760405163356680b760e01b815260040160405180910390fd5b816001600160a01b031683602001516001600160a01b03161415610f1857604051621af04160e61b815260040160405180910390fd5b60005a61138881101561260b57600080fd5b61138881039050846040820482031161262357600080fd5b50823b61262f57600080fd5b60008083516020850160008789f1949350505050565b6000818152600183016020526040812054801561272e576000612669600183613255565b855490915060009061267d90600190613255565b90508181146126e257600086600001828154811061269d5761269d61332c565b90600052602060002001549050808760000184815481106126c0576126c061332c565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806126f3576126f3613316565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610a71565b6000915050610a71565b82805461274490613294565b90600052602060002090601f01602090048101928261276657600085556127ac565b82601f1061277f5782800160ff198235161785556127ac565b828001600101855582156127ac579182015b828111156127ac578235825591602001919060010190612791565b506127b89291506127bc565b5090565b5b808211156127b857600081556001016127bd565b80356001600160a01b03811681146127e857600080fd5b919050565b60008083601f8401126127ff57600080fd5b5081356001600160401b0381111561281657600080fd5b6020830191508360208260051b850101111561283157600080fd5b9250929050565b60008083601f84011261284a57600080fd5b5081356001600160401b0381111561286157600080fd5b60208301915083602082850101111561283157600080fd5b803561ffff811681146127e857600080fd5b803562ffffff811681146127e857600080fd5b803563ffffffff811681146127e857600080fd5b805169ffffffffffffffffffff811681146127e857600080fd5b80356001600160601b03811681146127e857600080fd5b6000602082840312156128f557600080fd5b6118a4826127d1565b6000806040838503121561291157600080fd5b61291a836127d1565b9150612928602084016127d1565b90509250929050565b6000806040838503121561294457600080fd5b61294d836127d1565b915060208301356004811061296157600080fd5b809150509250929050565b6000806000806060858703121561298257600080fd5b61298b856127d1565b93506020850135925060408501356001600160401b038111156129ad57600080fd5b6129b987828801612838565b95989497509550505050565b6000806000806000608086880312156129dd57600080fd5b6129e6866127d1565b94506129f46020870161289e565b9350612a02604087016127d1565b925060608601356001600160401b03811115612a1d57600080fd5b612a2988828901612838565b969995985093965092949392505050565b60008060008060408587031215612a5057600080fd5b84356001600160401b0380821115612a6757600080fd5b612a73888389016127ed565b90965094506020870135915080821115612a8c57600080fd5b506129b9878288016127ed565b600080600060408486031215612aae57600080fd5b83356001600160401b03811115612ac457600080fd5b612ad0868287016127ed565b9094509250612ae39050602085016127d1565b90509250925092565b60008060208385031215612aff57600080fd5b82356001600160401b03811115612b1557600080fd5b612b2185828601612838565b90969095509350505050565b60006101808284031215612b4057600080fd5b612b4861319a565b612b518361289e565b8152612b5f6020840161289e565b6020820152612b706040840161288b565b6040820152612b816060840161289e565b6060820152612b926080840161288b565b6080820152612ba360a08401612879565b60a0820152612bb460c084016128cc565b60c0820152612bc560e0840161289e565b60e082015261010083810135908201526101208084013590820152610140612bee8185016127d1565b90820152610160612c008482016127d1565b908201529392505050565b600060208284031215612c1d57600080fd5b5035919050565b600060208284031215612c3657600080fd5b5051919050565b60008060408385031215612c5057600080fd5b82359150612928602084016127d1565b600080600060408486031215612c7557600080fd5b8335925060208401356001600160401b03811115612c9257600080fd5b612c9e86828701612838565b9497909650939450505050565b60008060408385031215612cbe57600080fd5b50508035926020909101359150565b60008060408385031215612ce057600080fd5b823591506129286020840161289e565b60008060408385031215612d0357600080fd5b82359150612928602084016128cc565b600080600080600060a08688031215612d2b57600080fd5b612d34866128b2565b9450602086015193506040860151925060608601519150612d57608087016128b2565b90509295509295909350565b6000815180845260005b81811015612d8957602081850181015186830182015201612d6d565b81811115612d9b576000602083870101525b50601f01601f19169290920160200192915050565b805163ffffffff1682526020810151612dd1602084018263ffffffff169052565b506040810151612de8604084018262ffffff169052565b506060810151612e00606084018263ffffffff169052565b506080810151612e17608084018262ffffff169052565b5060a0810151612e2d60a084018261ffff169052565b5060c0810151612e4860c08401826001600160601b03169052565b5060e0810151612e6060e084018263ffffffff169052565b5061010081810151908301526101208082015190830152610140808201516001600160a01b038116828501525050610160818101516001600160a01b0381168483015261138b565b828482376000838201600080825280855482600182811c915080831680612ed057607f831692505b6020808410821415612ef057634e487b7160e01b87526022600452602487fd5b818015612f045760018114612f1557612f41565b60ff19861689528489019650612f41565b60008c815260209020885b86811015612f395781548b820152908501908301612f20565b505084890196505b50949c9b505050505050505050505050565b6001600160a01b038a8116825263ffffffff8a16602083015261012060408301819052600091612f858483018c612d63565b6001600160601b039a8b16606086015298811660808501529690961660a0830152506001600160401b039390931660c0840152941660e082015292151561010090930192909252949350505050565b6020808252825182820181905260009190848201906040850190845b8181101561300c57835183529284019291840191600101612ff0565b50909695505050505050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b6020815260006118a46020830184612d63565b60a08152600061306d60a0830188612d63565b90508560208301528460408301528360608301528260808301529695505050505050565b60208101600483106130a5576130a5613300565b91905290565b60208101600383106130a5576130a5613300565b60208101600283106130a5576130a5613300565b6101808101610a718284612db0565b600061022080830163ffffffff8751168452602060018060601b0381890151168186015260408801516040860152606088015160608601526131276080860188612db0565b6102008501929092528451908190526102408401918086019160005b818110156131685783516001600160a01b031685529382019392820192600101613143565b509298975050505050505050565b6001600160601b0383168152604060208201819052600090610f0890830184612d63565b60405161018081016001600160401b03811182821017156131cb57634e487b7160e01b600052604160045260246000fd5b60405290565b600082198211156131e4576131e46132ea565b500190565b60006001600160601b0382811684821680830382111561320b5761320b6132ea565b01949350505050565b60008261323157634e487b7160e01b600052601260045260246000fd5b500490565b6000816000190483118215151615613250576132506132ea565b500290565b600082821015613267576132676132ea565b500390565b60006001600160601b038381169083168181101561328c5761328c6132ea565b039392505050565b600181811c908216806132a857607f821691505b602082108114156132c957634e487b7160e01b600052602260045260246000fd5b50919050565b60006000198214156132e3576132e36132ea565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220d847111b812cd574e6dad1f58f4fd5c8b8115a1ba3734352da8ba6672cc6da1264736f6c634300080600333078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666663078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666",
}

// KeeperRegistry20ABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistry20MetaData.ABI instead.
var KeeperRegistry20ABI = KeeperRegistry20MetaData.ABI

// KeeperRegistry20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistry20MetaData.Bin instead.
var KeeperRegistry20Bin = KeeperRegistry20MetaData.Bin

// DeployKeeperRegistry20 deploys a new Ethereum contract, binding an instance of KeeperRegistry20 to it.
func DeployKeeperRegistry20(auth *bind.TransactOpts, backend bind.ContractBackend, paymentModel uint8, registryGasOverhead *big.Int, link common.Address, linkEthFeed common.Address, fastGasFeed common.Address, keeperRegistryLogic common.Address, config Config) (common.Address, *types.Transaction, *KeeperRegistry20, error) {
	parsed, err := KeeperRegistry20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistry20Bin), backend, paymentModel, registryGasOverhead, link, linkEthFeed, fastGasFeed, keeperRegistryLogic, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry20{KeeperRegistry20Caller: KeeperRegistry20Caller{contract: contract}, KeeperRegistry20Transactor: KeeperRegistry20Transactor{contract: contract}, KeeperRegistry20Filterer: KeeperRegistry20Filterer{contract: contract}}, nil
}

// KeeperRegistry20 is an auto generated Go binding around an Ethereum contract.
type KeeperRegistry20 struct {
	KeeperRegistry20Caller     // Read-only binding to the contract
	KeeperRegistry20Transactor // Write-only binding to the contract
	KeeperRegistry20Filterer   // Log filterer for contract events
}

// KeeperRegistry20Caller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistry20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistry20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistry20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistry20Session struct {
	Contract     *KeeperRegistry20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRegistry20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistry20CallerSession struct {
	Contract *KeeperRegistry20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// KeeperRegistry20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistry20TransactorSession struct {
	Contract     *KeeperRegistry20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// KeeperRegistry20Raw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistry20Raw struct {
	Contract *KeeperRegistry20 // Generic contract binding to access the raw methods on
}

// KeeperRegistry20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistry20CallerRaw struct {
	Contract *KeeperRegistry20Caller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistry20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistry20TransactorRaw struct {
	Contract *KeeperRegistry20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistry20 creates a new instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20(address common.Address, backend bind.ContractBackend) (*KeeperRegistry20, error) {
	contract, err := bindKeeperRegistry20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20{KeeperRegistry20Caller: KeeperRegistry20Caller{contract: contract}, KeeperRegistry20Transactor: KeeperRegistry20Transactor{contract: contract}, KeeperRegistry20Filterer: KeeperRegistry20Filterer{contract: contract}}, nil
}

// NewKeeperRegistry20Caller creates a new read-only instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistry20Caller, error) {
	contract, err := bindKeeperRegistry20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20Caller{contract: contract}, nil
}

// NewKeeperRegistry20Transactor creates a new write-only instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistry20Transactor, error) {
	contract, err := bindKeeperRegistry20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20Transactor{contract: contract}, nil
}

// NewKeeperRegistry20Filterer creates a new log filterer instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistry20Filterer, error) {
	contract, err := bindKeeperRegistry20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20Filterer{contract: contract}, nil
}

// bindKeeperRegistry20 binds a generic wrapper to an already deployed contract.
func bindKeeperRegistry20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistry20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry20 *KeeperRegistry20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry20.Contract.KeeperRegistry20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry20 *KeeperRegistry20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.KeeperRegistry20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry20 *KeeperRegistry20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.KeeperRegistry20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry20 *KeeperRegistry20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry20 *KeeperRegistry20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry20 *KeeperRegistry20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.contract.Transact(opts, method, params...)
}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "ARB_NITRO_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistry20.Contract.ARBNITROORACLE(&_KeeperRegistry20.CallOpts)
}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistry20.Contract.ARBNITROORACLE(&_KeeperRegistry20.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry20.Contract.FASTGASFEED(&_KeeperRegistry20.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry20.Contract.FASTGASFEED(&_KeeperRegistry20.CallOpts)
}

// KEEPERREGISTRYLOGIC is a free data retrieval call binding the contract method 0x0852c7c9.
//
// Solidity: function KEEPER_REGISTRY_LOGIC() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) KEEPERREGISTRYLOGIC(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "KEEPER_REGISTRY_LOGIC")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KEEPERREGISTRYLOGIC is a free data retrieval call binding the contract method 0x0852c7c9.
//
// Solidity: function KEEPER_REGISTRY_LOGIC() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) KEEPERREGISTRYLOGIC() (common.Address, error) {
	return _KeeperRegistry20.Contract.KEEPERREGISTRYLOGIC(&_KeeperRegistry20.CallOpts)
}

// KEEPERREGISTRYLOGIC is a free data retrieval call binding the contract method 0x0852c7c9.
//
// Solidity: function KEEPER_REGISTRY_LOGIC() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) KEEPERREGISTRYLOGIC() (common.Address, error) {
	return _KeeperRegistry20.Contract.KEEPERREGISTRYLOGIC(&_KeeperRegistry20.CallOpts)
}

// L1FEEDATAPADDING is a free data retrieval call binding the contract method 0xbe308177.
//
// Solidity: function L1_FEE_DATA_PADDING() view returns(bytes)
func (_KeeperRegistry20 *KeeperRegistry20Caller) L1FEEDATAPADDING(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "L1_FEE_DATA_PADDING")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// L1FEEDATAPADDING is a free data retrieval call binding the contract method 0xbe308177.
//
// Solidity: function L1_FEE_DATA_PADDING() view returns(bytes)
func (_KeeperRegistry20 *KeeperRegistry20Session) L1FEEDATAPADDING() ([]byte, error) {
	return _KeeperRegistry20.Contract.L1FEEDATAPADDING(&_KeeperRegistry20.CallOpts)
}

// L1FEEDATAPADDING is a free data retrieval call binding the contract method 0xbe308177.
//
// Solidity: function L1_FEE_DATA_PADDING() view returns(bytes)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) L1FEEDATAPADDING() ([]byte, error) {
	return _KeeperRegistry20.Contract.L1FEEDATAPADDING(&_KeeperRegistry20.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) LINK() (common.Address, error) {
	return _KeeperRegistry20.Contract.LINK(&_KeeperRegistry20.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) LINK() (common.Address, error) {
	return _KeeperRegistry20.Contract.LINK(&_KeeperRegistry20.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry20.Contract.LINKETHFEED(&_KeeperRegistry20.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry20.Contract.LINKETHFEED(&_KeeperRegistry20.CallOpts)
}

// MAXINPUTDATA is a free data retrieval call binding the contract method 0xa0ad00cf.
//
// Solidity: function MAX_INPUT_DATA() view returns(bytes)
func (_KeeperRegistry20 *KeeperRegistry20Caller) MAXINPUTDATA(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "MAX_INPUT_DATA")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MAXINPUTDATA is a free data retrieval call binding the contract method 0xa0ad00cf.
//
// Solidity: function MAX_INPUT_DATA() view returns(bytes)
func (_KeeperRegistry20 *KeeperRegistry20Session) MAXINPUTDATA() ([]byte, error) {
	return _KeeperRegistry20.Contract.MAXINPUTDATA(&_KeeperRegistry20.CallOpts)
}

// MAXINPUTDATA is a free data retrieval call binding the contract method 0xa0ad00cf.
//
// Solidity: function MAX_INPUT_DATA() view returns(bytes)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) MAXINPUTDATA() ([]byte, error) {
	return _KeeperRegistry20.Contract.MAXINPUTDATA(&_KeeperRegistry20.CallOpts)
}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "OPTIMISM_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistry20.Contract.OPTIMISMORACLE(&_KeeperRegistry20.CallOpts)
}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistry20.Contract.OPTIMISMORACLE(&_KeeperRegistry20.CallOpts)
}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) PAYMENTMODEL(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "PAYMENT_MODEL")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistry20.Contract.PAYMENTMODEL(&_KeeperRegistry20.CallOpts)
}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistry20.Contract.PAYMENTMODEL(&_KeeperRegistry20.CallOpts)
}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistry20 *KeeperRegistry20Caller) REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "REGISTRY_GAS_OVERHEAD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistry20 *KeeperRegistry20Session) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistry20.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistry20.CallOpts)
}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistry20.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistry20.CallOpts)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry20 *KeeperRegistry20Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry20.Contract.GetActiveUpkeepIDs(&_KeeperRegistry20.CallOpts, startIndex, maxCount)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry20.Contract.GetActiveUpkeepIDs(&_KeeperRegistry20.CallOpts, startIndex, maxCount)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetKeeperInfo(opts *bind.CallOpts, query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getKeeperInfo", query)

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
func (_KeeperRegistry20 *KeeperRegistry20Session) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry20.Contract.GetKeeperInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry20.Contract.GetKeeperInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMaxPaymentForGas(&_KeeperRegistry20.CallOpts, gasLimit)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMaxPaymentForGas(&_KeeperRegistry20.CallOpts, gasLimit)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry20.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry20.CallOpts, peer)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry20.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry20.CallOpts, peer)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetState(opts *bind.CallOpts) (struct {
	State   State
	Config  Config
	Keepers []common.Address
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getState")

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
func (_KeeperRegistry20 *KeeperRegistry20Session) GetState() (struct {
	State   State
	Config  Config
	Keepers []common.Address
}, error) {
	return _KeeperRegistry20.Contract.GetState(&_KeeperRegistry20.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetState() (struct {
	State   State
	Config  Config
	Keepers []common.Address
}, error) {
	return _KeeperRegistry20.Contract.GetState(&_KeeperRegistry20.CallOpts)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent, bool paused)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
	Paused              bool
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getUpkeep", id)

	outstruct := new(struct {
		Target              common.Address
		ExecuteGas          uint32
		CheckData           []byte
		Balance             *big.Int
		LastKeeper          common.Address
		Admin               common.Address
		MaxValidBlocknumber uint64
		AmountSpent         *big.Int
		Paused              bool
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
	outstruct.Paused = *abi.ConvertType(out[8], new(bool)).(*bool)

	return *outstruct, err

}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent, bool paused)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
	Paused              bool
}, error) {
	return _KeeperRegistry20.Contract.GetUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent, bool paused)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
	Paused              bool
}, error) {
	return _KeeperRegistry20.Contract.GetUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) Owner() (common.Address, error) {
	return _KeeperRegistry20.Contract.Owner(&_KeeperRegistry20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry20.Contract.Owner(&_KeeperRegistry20.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry20 *KeeperRegistry20Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry20 *KeeperRegistry20Session) Paused() (bool, error) {
	return _KeeperRegistry20.Contract.Paused(&_KeeperRegistry20.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) Paused() (bool, error) {
	return _KeeperRegistry20.Contract.Paused(&_KeeperRegistry20.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry20 *KeeperRegistry20Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry20 *KeeperRegistry20Session) TypeAndVersion() (string, error) {
	return _KeeperRegistry20.Contract.TypeAndVersion(&_KeeperRegistry20.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry20.Contract.TypeAndVersion(&_KeeperRegistry20.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry20.Contract.UpkeepTranscoderVersion(&_KeeperRegistry20.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry20.Contract.UpkeepTranscoderVersion(&_KeeperRegistry20.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptOwnership(&_KeeperRegistry20.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptOwnership(&_KeeperRegistry20.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptPayeeship", keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptPayeeship(&_KeeperRegistry20.TransactOpts, keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptPayeeship(&_KeeperRegistry20.TransactOpts, keeper)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "addFunds", id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AddFunds(&_KeeperRegistry20.TransactOpts, id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AddFunds(&_KeeperRegistry20.TransactOpts, id, amount)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "cancelUpkeep", id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CancelUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CancelUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "checkUpkeep", id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry20 *KeeperRegistry20Session) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CheckUpkeep(&_KeeperRegistry20.TransactOpts, id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CheckUpkeep(&_KeeperRegistry20.TransactOpts, id, from)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.MigrateUpkeeps(&_KeeperRegistry20.TransactOpts, ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.MigrateUpkeeps(&_KeeperRegistry20.TransactOpts, ids, destination)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.OnTokenTransfer(&_KeeperRegistry20.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.OnTokenTransfer(&_KeeperRegistry20.TransactOpts, sender, amount, data)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Pause(&_KeeperRegistry20.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Pause(&_KeeperRegistry20.TransactOpts)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "pauseUpkeep", id)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.PauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.PauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "performUpkeep", id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry20 *KeeperRegistry20Session) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.PerformUpkeep(&_KeeperRegistry20.TransactOpts, id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.PerformUpkeep(&_KeeperRegistry20.TransactOpts, id, performData)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.ReceiveUpkeeps(&_KeeperRegistry20.TransactOpts, encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.ReceiveUpkeeps(&_KeeperRegistry20.TransactOpts, encodedUpkeeps)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "recoverFunds")
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RecoverFunds(&_KeeperRegistry20.TransactOpts)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RecoverFunds(&_KeeperRegistry20.TransactOpts)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RegisterUpkeep(&_KeeperRegistry20.TransactOpts, target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RegisterUpkeep(&_KeeperRegistry20.TransactOpts, target, gasLimit, admin, checkData)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetConfig(opts *bind.TransactOpts, config Config) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setConfig", config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetConfig(&_KeeperRegistry20.TransactOpts, config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetConfig(&_KeeperRegistry20.TransactOpts, config)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setKeepers", keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetKeepers(&_KeeperRegistry20.TransactOpts, keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetKeepers(&_KeeperRegistry20.TransactOpts, keepers, payees)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry20.TransactOpts, peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry20.TransactOpts, peer, permission)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepGasLimit(&_KeeperRegistry20.TransactOpts, id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepGasLimit(&_KeeperRegistry20.TransactOpts, id, gasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferOwnership(&_KeeperRegistry20.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferOwnership(&_KeeperRegistry20.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferPayeeship(&_KeeperRegistry20.TransactOpts, keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferPayeeship(&_KeeperRegistry20.TransactOpts, keeper, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id, proposed)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Unpause(&_KeeperRegistry20.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Unpause(&_KeeperRegistry20.TransactOpts)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "unpauseUpkeep", id)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UnpauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UnpauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UpdateCheckData(&_KeeperRegistry20.TransactOpts, id, newCheckData)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UpdateCheckData(&_KeeperRegistry20.TransactOpts, id, newCheckData)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "withdrawFunds", id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawFunds(&_KeeperRegistry20.TransactOpts, id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawFunds(&_KeeperRegistry20.TransactOpts, id, to)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "withdrawOwnerFunds")
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawOwnerFunds(&_KeeperRegistry20.TransactOpts)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawOwnerFunds(&_KeeperRegistry20.TransactOpts)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "withdrawPayment", from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawPayment(&_KeeperRegistry20.TransactOpts, from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawPayment(&_KeeperRegistry20.TransactOpts, from, to)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Fallback(&_KeeperRegistry20.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Fallback(&_KeeperRegistry20.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Receive() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Receive(&_KeeperRegistry20.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Receive() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Receive(&_KeeperRegistry20.TransactOpts)
}

// KeeperRegistry20ConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the KeeperRegistry20 contract.
type KeeperRegistry20ConfigSetIterator struct {
	Event *KeeperRegistry20ConfigSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20ConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20ConfigSet)
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
		it.Event = new(KeeperRegistry20ConfigSet)
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
func (it *KeeperRegistry20ConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20ConfigSet represents a ConfigSet event raised by the KeeperRegistry20 contract.
type KeeperRegistry20ConfigSet struct {
	Config Config
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry20ConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20ConfigSetIterator{contract: _KeeperRegistry20.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20ConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20ConfigSet)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseConfigSet(log types.Log) (*KeeperRegistry20ConfigSet, error) {
	event := new(KeeperRegistry20ConfigSet)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20FundsAddedIterator is returned from FilterFundsAdded and is used to iterate over the raw logs and unpacked data for FundsAdded events raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsAddedIterator struct {
	Event *KeeperRegistry20FundsAdded // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20FundsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20FundsAdded)
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
		it.Event = new(KeeperRegistry20FundsAdded)
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
func (it *KeeperRegistry20FundsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20FundsAdded represents a FundsAdded event raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsAdded is a free log retrieval operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistry20FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20FundsAddedIterator{contract: _KeeperRegistry20.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

// WatchFundsAdded is a free log subscription operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20FundsAdded)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistry20FundsAdded, error) {
	event := new(KeeperRegistry20FundsAdded)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20FundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsWithdrawnIterator struct {
	Event *KeeperRegistry20FundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20FundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20FundsWithdrawn)
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
		it.Event = new(KeeperRegistry20FundsWithdrawn)
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
func (it *KeeperRegistry20FundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20FundsWithdrawn represents a FundsWithdrawn event raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20FundsWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20FundsWithdrawn)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistry20FundsWithdrawn, error) {
	event := new(KeeperRegistry20FundsWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20KeepersUpdatedIterator is returned from FilterKeepersUpdated and is used to iterate over the raw logs and unpacked data for KeepersUpdated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20KeepersUpdatedIterator struct {
	Event *KeeperRegistry20KeepersUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20KeepersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20KeepersUpdated)
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
		it.Event = new(KeeperRegistry20KeepersUpdated)
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
func (it *KeeperRegistry20KeepersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20KeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20KeepersUpdated represents a KeepersUpdated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20KeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKeepersUpdated is a free log retrieval operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistry20KeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20KeepersUpdatedIterator{contract: _KeeperRegistry20.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

// WatchKeepersUpdated is a free log subscription operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20KeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20KeepersUpdated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistry20KeepersUpdated, error) {
	event := new(KeeperRegistry20KeepersUpdated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20OwnerFundsWithdrawnIterator is returned from FilterOwnerFundsWithdrawn and is used to iterate over the raw logs and unpacked data for OwnerFundsWithdrawn events raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistry20OwnerFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20OwnerFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20OwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistry20OwnerFundsWithdrawn)
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
func (it *KeeperRegistry20OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20OwnerFundsWithdrawn represents a OwnerFundsWithdrawn event raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerFundsWithdrawn is a free log retrieval operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistry20OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20OwnerFundsWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchOwnerFundsWithdrawn is a free log subscription operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20OwnerFundsWithdrawn)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistry20OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistry20OwnerFundsWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20OwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistry20OwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20OwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20OwnershipTransferRequested)
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
		it.Event = new(KeeperRegistry20OwnershipTransferRequested)
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
func (it *KeeperRegistry20OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20OwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry20OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20OwnershipTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20OwnershipTransferRequested)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistry20OwnershipTransferRequested, error) {
	event := new(KeeperRegistry20OwnershipTransferRequested)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferredIterator struct {
	Event *KeeperRegistry20OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20OwnershipTransferred)
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
		it.Event = new(KeeperRegistry20OwnershipTransferred)
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
func (it *KeeperRegistry20OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20OwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry20OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20OwnershipTransferredIterator{contract: _KeeperRegistry20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20OwnershipTransferred)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistry20OwnershipTransferred, error) {
	event := new(KeeperRegistry20OwnershipTransferred)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PausedIterator struct {
	Event *KeeperRegistry20Paused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20Paused)
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
		it.Event = new(KeeperRegistry20Paused)
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
func (it *KeeperRegistry20PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20Paused represents a Paused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistry20PausedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PausedIterator{contract: _KeeperRegistry20.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20Paused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePaused(log types.Log) (*KeeperRegistry20Paused, error) {
	event := new(KeeperRegistry20Paused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistry20PayeeshipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20PayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistry20PayeeshipTransferRequested)
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
func (it *KeeperRegistry20PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry20PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeeshipTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PayeeshipTransferRequested)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistry20PayeeshipTransferRequested, error) {
	event := new(KeeperRegistry20PayeeshipTransferRequested)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferredIterator struct {
	Event *KeeperRegistry20PayeeshipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20PayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PayeeshipTransferred)
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
		it.Event = new(KeeperRegistry20PayeeshipTransferred)
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
func (it *KeeperRegistry20PayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PayeeshipTransferred represents a PayeeshipTransferred event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry20PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeeshipTransferredIterator{contract: _KeeperRegistry20.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PayeeshipTransferred)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistry20PayeeshipTransferred, error) {
	event := new(KeeperRegistry20PayeeshipTransferred)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PaymentWithdrawnIterator is returned from FilterPaymentWithdrawn and is used to iterate over the raw logs and unpacked data for PaymentWithdrawn events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PaymentWithdrawnIterator struct {
	Event *KeeperRegistry20PaymentWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20PaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PaymentWithdrawn)
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
		it.Event = new(KeeperRegistry20PaymentWithdrawn)
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
func (it *KeeperRegistry20PaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PaymentWithdrawn represents a PaymentWithdrawn event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry20PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PaymentWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PaymentWithdrawn)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry20PaymentWithdrawn, error) {
	event := new(KeeperRegistry20PaymentWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UnpausedIterator struct {
	Event *KeeperRegistry20Unpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20Unpaused)
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
		it.Event = new(KeeperRegistry20Unpaused)
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
func (it *KeeperRegistry20UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20Unpaused represents a Unpaused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistry20UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UnpausedIterator{contract: _KeeperRegistry20.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20Unpaused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUnpaused(log types.Log) (*KeeperRegistry20Unpaused, error) {
	event := new(KeeperRegistry20Unpaused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepAdminTransferRequestedIterator is returned from FilterUpkeepAdminTransferRequested and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferRequested events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistry20UpkeepAdminTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepAdminTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepAdminTransferRequested)
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
		it.Event = new(KeeperRegistry20UpkeepAdminTransferRequested)
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
func (it *KeeperRegistry20UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepAdminTransferRequested represents a UpkeepAdminTransferRequested event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferRequested is a free log retrieval operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry20UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferRequested is a free log subscription operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepAdminTransferRequested)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistry20UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistry20UpkeepAdminTransferRequested)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepAdminTransferredIterator is returned from FilterUpkeepAdminTransferred and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferred events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistry20UpkeepAdminTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepAdminTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepAdminTransferred)
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
		it.Event = new(KeeperRegistry20UpkeepAdminTransferred)
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
func (it *KeeperRegistry20UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepAdminTransferred represents a UpkeepAdminTransferred event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferred is a free log retrieval operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry20UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepAdminTransferredIterator{contract: _KeeperRegistry20.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferred is a free log subscription operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepAdminTransferred)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistry20UpkeepAdminTransferred, error) {
	event := new(KeeperRegistry20UpkeepAdminTransferred)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepCanceledIterator is returned from FilterUpkeepCanceled and is used to iterate over the raw logs and unpacked data for UpkeepCanceled events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCanceledIterator struct {
	Event *KeeperRegistry20UpkeepCanceled // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepCanceled)
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
		it.Event = new(KeeperRegistry20UpkeepCanceled)
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
func (it *KeeperRegistry20UpkeepCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepCanceled represents a UpkeepCanceled event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCanceled is a free log retrieval operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistry20UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepCanceledIterator{contract: _KeeperRegistry20.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

// WatchUpkeepCanceled is a free log subscription operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepCanceled)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistry20UpkeepCanceled, error) {
	event := new(KeeperRegistry20UpkeepCanceled)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepCheckDataUpdatedIterator is returned from FilterUpkeepCheckDataUpdated and is used to iterate over the raw logs and unpacked data for UpkeepCheckDataUpdated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistry20UpkeepCheckDataUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepCheckDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepCheckDataUpdated)
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
		it.Event = new(KeeperRegistry20UpkeepCheckDataUpdated)
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
func (it *KeeperRegistry20UpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepCheckDataUpdated represents a UpkeepCheckDataUpdated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCheckDataUpdated is a free log retrieval operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepCheckDataUpdatedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

// WatchUpkeepCheckDataUpdated is a free log subscription operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepCheckDataUpdated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistry20UpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistry20UpkeepCheckDataUpdated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepGasLimitSetIterator is returned from FilterUpkeepGasLimitSet and is used to iterate over the raw logs and unpacked data for UpkeepGasLimitSet events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistry20UpkeepGasLimitSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistry20UpkeepGasLimitSet)
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
func (it *KeeperRegistry20UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepGasLimitSet represents a UpkeepGasLimitSet event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpkeepGasLimitSet is a free log retrieval operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepGasLimitSetIterator{contract: _KeeperRegistry20.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepGasLimitSet is a free log subscription operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepGasLimitSet)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistry20UpkeepGasLimitSet, error) {
	event := new(KeeperRegistry20UpkeepGasLimitSet)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepMigratedIterator is returned from FilterUpkeepMigrated and is used to iterate over the raw logs and unpacked data for UpkeepMigrated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepMigratedIterator struct {
	Event *KeeperRegistry20UpkeepMigrated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepMigrated)
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
		it.Event = new(KeeperRegistry20UpkeepMigrated)
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
func (it *KeeperRegistry20UpkeepMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepMigrated represents a UpkeepMigrated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepMigrated is a free log retrieval operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepMigratedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

// WatchUpkeepMigrated is a free log subscription operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepMigrated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistry20UpkeepMigrated, error) {
	event := new(KeeperRegistry20UpkeepMigrated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepPausedIterator is returned from FilterUpkeepPaused and is used to iterate over the raw logs and unpacked data for UpkeepPaused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPausedIterator struct {
	Event *KeeperRegistry20UpkeepPaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepPaused)
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
		it.Event = new(KeeperRegistry20UpkeepPaused)
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
func (it *KeeperRegistry20UpkeepPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepPaused represents a UpkeepPaused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPaused is a free log retrieval operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepPausedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepPaused is a free log subscription operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepPaused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistry20UpkeepPaused, error) {
	event := new(KeeperRegistry20UpkeepPaused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepPerformedIterator is returned from FilterUpkeepPerformed and is used to iterate over the raw logs and unpacked data for UpkeepPerformed events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPerformedIterator struct {
	Event *KeeperRegistry20UpkeepPerformed // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepPerformedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepPerformed)
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
		it.Event = new(KeeperRegistry20UpkeepPerformed)
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
func (it *KeeperRegistry20UpkeepPerformedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepPerformed represents a UpkeepPerformed event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPerformed struct {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistry20UpkeepPerformedIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepPerformedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepPerformed)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistry20UpkeepPerformed, error) {
	event := new(KeeperRegistry20UpkeepPerformed)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepReceivedIterator is returned from FilterUpkeepReceived and is used to iterate over the raw logs and unpacked data for UpkeepReceived events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepReceivedIterator struct {
	Event *KeeperRegistry20UpkeepReceived // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepReceived)
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
		it.Event = new(KeeperRegistry20UpkeepReceived)
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
func (it *KeeperRegistry20UpkeepReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepReceived represents a UpkeepReceived event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpkeepReceived is a free log retrieval operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepReceivedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

// WatchUpkeepReceived is a free log subscription operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepReceived)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistry20UpkeepReceived, error) {
	event := new(KeeperRegistry20UpkeepReceived)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepRegisteredIterator is returned from FilterUpkeepRegistered and is used to iterate over the raw logs and unpacked data for UpkeepRegistered events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepRegisteredIterator struct {
	Event *KeeperRegistry20UpkeepRegistered // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepRegistered)
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
		it.Event = new(KeeperRegistry20UpkeepRegistered)
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
func (it *KeeperRegistry20UpkeepRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepRegistered represents a UpkeepRegistered event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpkeepRegistered is a free log retrieval operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepRegisteredIterator{contract: _KeeperRegistry20.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

// WatchUpkeepRegistered is a free log subscription operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepRegistered)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistry20UpkeepRegistered, error) {
	event := new(KeeperRegistry20UpkeepRegistered)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepUnpausedIterator is returned from FilterUpkeepUnpaused and is used to iterate over the raw logs and unpacked data for UpkeepUnpaused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepUnpausedIterator struct {
	Event *KeeperRegistry20UpkeepUnpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepUnpaused)
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
		it.Event = new(KeeperRegistry20UpkeepUnpaused)
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
func (it *KeeperRegistry20UpkeepUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepUnpaused represents a UpkeepUnpaused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepUnpaused is a free log retrieval operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepUnpausedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepUnpaused is a free log subscription operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepUnpaused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistry20UpkeepUnpaused, error) {
	event := new(KeeperRegistry20UpkeepUnpaused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
