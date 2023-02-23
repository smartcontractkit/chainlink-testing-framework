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

// OnchainConfig2_0 is an auto generated low-level Go binding around an user-defined struct.
type OnchainConfig2_0 struct {
	PaymentPremiumPPB    uint32
	FlatFeeMicroLink     uint32
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	MinUpkeepSpend       *big.Int
	MaxPerformGas        uint32
	MaxCheckDataSize     uint32
	MaxPerformDataSize   uint32
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
	Transcoder           common.Address
	Registrar            common.Address
}

// State2_0 is an auto generated low-level Go binding around an user-defined struct.
type State2_0 struct {
	Nonce                   uint32
	OwnerLinkBalance        *big.Int
	ExpectedLinkBalance     *big.Int
	TotalPremium            *big.Int
	NumUpkeeps              *big.Int
	ConfigCount             uint32
	LatestConfigBlockNumber uint32
	LatestConfigDigest      [32]byte
	LatestEpoch             uint32
	Paused                  bool
}

// UpkeepInfo is an auto generated low-level Go binding around an user-defined struct.
type UpkeepInfo struct {
	Target                 common.Address
	ExecuteGas             uint32
	CheckData              []byte
	Balance                *big.Int
	Admin                  common.Address
	MaxValidBlocknumber    uint64
	LastPerformBlockNumber uint32
	AmountSpent            *big.Int
	Paused                 bool
	OffchainConfig         []byte
}

// KeeperRegistry20MetaData contains all meta data concerning the KeeperRegistry20 contract.
var KeeperRegistry20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryBase2_0\",\"name\":\"keeperRegistryLogic\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnchainConfigNonEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"checkBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getKeeperRegistryLogicAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPaymentModel\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structState\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structOnchainConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structUpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006316380380620063168339810160408190526200003591620003a7565b806001600160a01b0316634b4fd03b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200006f57600080fd5b505afa15801562000084573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000aa9190620003ce565b816001600160a01b031663ca30e6036040518163ffffffff1660e01b815260040160206040518083038186803b158015620000e457600080fd5b505afa158015620000f9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200011f9190620003a7565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200015957600080fd5b505afa1580156200016e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001949190620003a7565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b815260040160206040518083038186803b158015620001ce57600080fd5b505afa158015620001e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002099190620003a7565b3380600081620002605760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000293576200029381620002fb565b505050836002811115620002ab57620002ab620003f1565b60e0816002811115620002c257620002c2620003f1565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b909116610100525062000420565b6001600160a01b038116331415620003565760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000257565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620003ba57600080fd5b8151620003c78162000407565b9392505050565b600060208284031215620003e157600080fd5b815160038110620003c757600080fd5b634e487b7160e01b600052602160045260246000fd5b6001600160a01b03811681146200041d57600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c615e6a620004ac6000396000818161057f0152610a8201526000818161052c01528181613c5801528181614313015281816144ca015261470f0152600081816105d3015261342a015260008181610859015261351301526000818161091401526113210152615e6a6000f3fe6080604052600436106103175760003560e01c80638e86139b1161019a578063b1dc65a4116100e1578063e3d0e7121161008a578063f2fde38b11610064578063f2fde38b146109d8578063f7d334ba146109f8578063faa3e99614610a2a57610326565b8063e3d0e71214610938578063eb5dcd6c14610760578063ed56b3e11461095857610326565b8063c7c3a19a116100bb578063c7c3a19a146108d8578063c804802214610550578063ca30e6031461090557610326565b8063b1dc65a414610898578063b657bc9c146108b8578063b79550be1461047357610326565b8063aab9edd611610143578063b10b673c1161011d578063b10b673c1461084a578063b121e1471461087d578063b148ab6b1461055057610326565b8063aab9edd614610796578063aed2e929146107bd578063afcb95d7146107f457610326565b8063a4c0ed3611610174578063a4c0ed3614610740578063a710b22114610760578063a72aa27e1461077b57610326565b80638e86139b1461070a578063948108f7146107255780639fab4386146106ef57610326565b8063572e05e11161025e57806381ff7048116102075780638765ecbe116101e15780638765ecbe146105505780638da5cb5b146106c45780638dcf0fe7146106ef57610326565b806381ff70481461063a5780638456cb591461047357806385c1b0ba146106a457610326565b8063744bfe6111610238578063744bfe611461043d57806379ba5097146106255780637d9b97e01461047357610326565b8063572e05e1146105705780636709d0e5146105c45780636ded9eae146105f757610326565b80633b9cce59116102c057806348013d7b1161029a57806348013d7b146104fb5780634b4fd03b1461051d5780635165f2f51461055057610326565b80633b9cce59146104585780633f4ba83a14610473578063421d183b1461048857610326565b80631865c57d116102f15780631865c57d146103f7578063187256e81461041d5780631a2af0111461043d57610326565b806306e3b6321461032e5780630e08ae8414610364578063181f5a77146103a157610326565b3661032657610324610a7d565b005b610324610a7d565b34801561033a57600080fd5b5061034e61034936600461511e565b610aa8565b60405161035b91906154cd565b60405180910390f35b34801561037057600080fd5b5061038461037f366004615260565b610ba2565b6040516bffffffffffffffffffffffff909116815260200161035b565b3480156103ad57600080fd5b506103ea6040518060400160405280601481526020017f4b6565706572526567697374727920322e302e3200000000000000000000000081525081565b60405161035b919061557b565b34801561040357600080fd5b5061040c610ce5565b60405161035b9594939291906155b5565b34801561042957600080fd5b50610324610438366004614c02565b6110a8565b34801561044957600080fd5b506103246104383660046150ad565b34801561046457600080fd5b50610324610438366004614d38565b34801561047f57600080fd5b506103246110b4565b34801561049457600080fd5b506104a86104a3366004614bac565b6110bc565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a00161035b565b34801561050757600080fd5b50610510600081565b60405161035b91906155a8565b34801561052957600080fd5b507f0000000000000000000000000000000000000000000000000000000000000000610510565b34801561055c57600080fd5b5061032461056b366004615094565b6111da565b34801561057c57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161035b565b3480156105d057600080fd5b507f000000000000000000000000000000000000000000000000000000000000000061059f565b34801561060357600080fd5b50610617610612366004614c90565b6111e5565b60405190815260200161035b565b34801561063157600080fd5b506103246111fa565b34801561064657600080fd5b50610681601254600e5463ffffffff6c0100000000000000000000000083048116937001000000000000000000000000000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600161035b565b3480156106b057600080fd5b506103246106bf366004614efe565b6112fc565b3480156106d057600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661059f565b3480156106fb57600080fd5b506103246106bf3660046150d2565b34801561071657600080fd5b50610324610438366004614f6e565b34801561073157600080fd5b5061032461043836600461523b565b34801561074c57600080fd5b5061032461075b366004614c34565b611309565b34801561076c57600080fd5b50610324610438366004614bc9565b34801561078757600080fd5b50610324610438366004615216565b3480156107a257600080fd5b506107ab600281565b60405160ff909116815260200161035b565b3480156107c957600080fd5b506107dd6107d83660046150d2565b611524565b60408051921515835260208301919091520161035b565b34801561080057600080fd5b50600e54600f54604080516000815260208101939093527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169082015260600161035b565b34801561085657600080fd5b507f000000000000000000000000000000000000000000000000000000000000000061059f565b34801561088957600080fd5b5061032461056b366004614bac565b3480156108a457600080fd5b506103246108b3366004614e47565b61168f565b3480156108c457600080fd5b506103846108d3366004615094565b61224c565b3480156108e457600080fd5b506108f86108f3366004615094565b612270565b60405161035b91906156c2565b34801561091157600080fd5b507f000000000000000000000000000000000000000000000000000000000000000061059f565b34801561094457600080fd5b50610324610953366004614d7a565b61259b565b34801561096457600080fd5b506109bf610973366004614bac565b73ffffffffffffffffffffffffffffffffffffffff1660009081526009602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff90911660208301520161035b565b3480156109e457600080fd5b506103246109f3366004614bac565b613392565b348015610a0457600080fd5b50610a18610a13366004615094565b6133a3565b60405161035b96959493929190615511565b348015610a3657600080fd5b50610a70610a45366004614bac565b73ffffffffffffffffffffffffffffffffffffffff1660009081526016602052604090205460ff1690565b60405161035b919061558e565b610aa67f00000000000000000000000000000000000000000000000000000000000000006133c6565b565b60606000610ab660026133ea565b9050808410610af1576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82610b0357610b008482615bcf565b92505b60008367ffffffffffffffff811115610b1e57610b1e615d88565b604051908082528060200260200182016040528015610b47578160200160208202803683370190505b50905060005b84811015610b9957610b6a610b628288615a50565b6002906133f4565b828281518110610b7c57610b7c615d59565b602090810291909101015280610b9181615c93565b915050610b4d565b50949350505050565b6040805161012081018252600f5460ff808216835263ffffffff6101008084048216602086015265010000000000840482169585019590955262ffffff6901000000000000000000840416606085015261ffff6c0100000000000000000000000084041660808501526e01000000000000000000000000000083048216151560a08501526f010000000000000000000000000000008304909116151560c08401526bffffffffffffffffffffffff70010000000000000000000000000000000083041660e08401527c010000000000000000000000000000000000000000000000000000000090910416918101919091526000908180610ca183613407565b6012549193509150610cdc90849087907801000000000000000000000000000000000000000000000000900463ffffffff1685856000613603565b95945050505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101a081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101829052610180810191909152604080516101408101825260125468010000000000000000900463ffffffff1681526011546bffffffffffffffffffffffff908116602083015260155492820192909252600f54700100000000000000000000000000000000900490911660608083019190915290819060009060808101610e1a60026133ea565b815260125463ffffffff6c01000000000000000000000000808304821660208086019190915270010000000000000000000000000000000084048316604080870191909152600e54606080880191909152600f547c0100000000000000000000000000000000000000000000000000000000810486166080808a019190915260ff6e01000000000000000000000000000083048116151560a09a8b015284516101a0810186526101008085048a1682526501000000000085048a1682890152898b168288015262ffffff69010000000000000000008604169582019590955261ffff88850416928101929092526010546bffffffffffffffffffffffff81169a83019a909a526401000000008904881660c0830152740100000000000000000000000000000000000000008904881660e083015278010000000000000000000000000000000000000000000000009098049096169186019190915260135461012086015260145461014086015273ffffffffffffffffffffffffffffffffffffffff96849004871661016086015260115493909304909516610180840152600a8054865181840281018401909752808752969b509299508a958a959394600b949316929185919083018282801561102757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610ffc575b505050505092508180548060200260200160405190810160405280929190818152602001828054801561109057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611065575b50505050509150945094509450945094509091929394565b6110b0610a7d565b5050565b610aa6610a7d565b73ffffffffffffffffffffffffffffffffffffffff811660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e010000000000000000000000000000909204811660608301819052600f54849384938493849384926111689291700100000000000000000000000000000000900416615be6565b600b5490915060009061117b9083615ae7565b9050826000015183602001518285604001516111979190615aac565b6060959095015173ffffffffffffffffffffffffffffffffffffffff9b8c166000908152600c6020526040902054929c919b959a50985093169550919350505050565b6111e2610a7d565b50565b60006111ef610a7d565b979650505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611280576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611304610a7d565b505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611378576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146113b2576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006113c082840184615094565b600081815260046020526040902054909150640100000000900463ffffffff90811614611419576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546114549085906c0100000000000000000000000090046bffffffffffffffffffffffff16615aac565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556015546114bf908590615a50565b6015556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061152f61364e565b600f546e010000000000000000000000000000900460ff161561157e576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481168386015268010000000000000000820460ff16151583870152690100000000000000000090910473ffffffffffffffffffffffffffffffffffffffff1660608301526001909201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528251601f87018390048302810183019093528583529161168291839190889088908190840183828082843760009201919091525061368692505050565b9250925050935093915050565b60005a6040805161012081018252600f5460ff808216835261010080830463ffffffff90811660208601526501000000000084048116958501959095526901000000000000000000830462ffffff1660608501526c01000000000000000000000000830461ffff1660808501526e0100000000000000000000000000008304821615801560a08601526f010000000000000000000000000000008404909216151560c085015270010000000000000000000000000000000083046bffffffffffffffffffffffff1660e08501527c0100000000000000000000000000000000000000000000000000000000909204909316908201529192506117bd576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526008602052604090205460ff16611806576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006118478a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506137ec92505050565b9050600081604001515167ffffffffffffffff81111561186957611869615d88565b60405190808252806020026020018201604052801561191d57816020015b604080516101a081018252600060c0820181815260e083018290526101008301829052610120830182905261014083018290526101608301829052610180830182905282526020808301829052928201819052606082018190526080820181905260a082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816118875790505b5090506000805b836040015151811015611bc957600460008560400151838151811061194b5761194b615d59565b6020908102919091018101518252818101929092526040908101600020815160e081018352815463ffffffff8082168352640100000000820481169583019590955268010000000000000000810460ff16151593820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528351849083908110611a3557611a35615d59565b602002602001015160000181905250611a9e85848381518110611a5a57611a5a615d59565b6020026020010151600001516000015186606001518481518110611a8057611a80615d59565b60200260200101516040015151876000015188602001516001613603565b838281518110611ab057611ab0615d59565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050611b5e84604001518281518110611af757611af7615d59565b602002602001015185606001518381518110611b1557611b15615d59565b6020026020010151858481518110611b2f57611b2f615d59565b602002602001015160000151868581518110611b4d57611b4d615d59565b602002602001015160400151613898565b838281518110611b7057611b70615d59565b60200260200101516020019015159081151581525050828181518110611b9857611b98615d59565b60200260200101516020015115611bb757611bb4600183615a2a565b91505b80611bc181615c93565b915050611924565b5061ffff8116611bdd575050505050612242565b600e548d3514611c19576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8351611c26906001615a87565b60ff1689141580611c375750888714155b15611c6e576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611c7e8d8d8d8d8d8d8d8d6139e9565b60005b836040015151811015611e6857828181518110611ca057611ca0615d59565b60200260200101516020015115611e5657611cb9613c52565b63ffffffff166004600086604001518481518110611cd957611cd9615d59565b6020026020010151815260200190815260200160002060010160189054906101000a900463ffffffff1663ffffffff161415611d41576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611d89838281518110611d5657611d56615d59565b60200260200101516000015185606001518381518110611d7857611d78615d59565b602002602001015160400151613686565b848381518110611d9b57611d9b615d59565b6020026020010151606001858481518110611db857611db8615d59565b60200260200101516080018281525082151515158152505050828181518110611de357611de3615d59565b60200260200101516080015186611dfa9190615bcf565b9550611e04613c52565b6004600086604001518481518110611e1e57611e1e615d59565b6020026020010151815260200190815260200160002060010160186101000a81548163ffffffff021916908363ffffffff1602179055505b80611e6081615c93565b915050611c81565b508351611e76906001615a87565b611e859060ff1661044c615b12565b616914611e938d6010615b12565b5a611e9e9089615bcf565b611ea89190615a50565b611eb29190615a50565b611ebc9190615a50565b94506116a8611ecf61ffff831687615ad3565b611ed99190615a50565b945060008060008060005b8760400151518110156120e457868181518110611f0357611f03615d59565b602002602001015160200151156120d257611f458a89606001518381518110611f2e57611f2e615d59565b602002602001015160400151518b60000151613d17565b878281518110611f5757611f57615d59565b602002602001015160a0018181525050611fb38989604001518381518110611f8157611f81615d59565b6020026020010151898481518110611f9b57611f9b615d59565b60200260200101518b600001518c602001518b613d35565b9093509150611fc28285615aac565b9350611fce8386615aac565b9450868181518110611fe257611fe2615d59565b60200260200101516060015115158860400151828151811061200657612006615d59565b60200260200101517f29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a44208a60600151848151811061204557612045615d59565b6020026020010151600001518a858151811061206357612063615d59565b6020026020010151608001518b868151811061208157612081615d59565b602002602001015160a0015187896120999190615aac565b6040805163ffffffff90951685526020850193909352918301526bffffffffffffffffffffffff16606082015260800160405180910390a35b806120dc81615c93565b915050611ee4565b5050336000908152600860205260409020805484925060029061211c9084906201000090046bffffffffffffffffffffffff16615aac565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600f60000160108282829054906101000a90046bffffffffffffffffffffffff166121769190615aac565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106121b9576121b9615d59565b602002013560001c9050600060088264ffffffffff16901c905087610100015163ffffffff168163ffffffff16111561223857600f80547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff8416021790555b5050505050505050505b5050505050505050565b60008181526004602052604081205461226a9063ffffffff16610ba2565b92915050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815160e081018352815463ffffffff8082168352640100000000820481168387015268010000000000000000820460ff16151583860152690100000000000000000090910473ffffffffffffffffffffffffffffffffffffffff908116606084019081526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a085015278010000000000000000000000000000000000000000000000009004821660c08401528451610140810186529351168352815116828501528685526007909352928190208054929392918301916123bc90615c3f565b80601f01602080910402602001604051908101604052809291908181526020018280546123e890615c3f565b80156124355780601f1061240a57610100808354040283529160200191612435565b820191906000526020600020905b81548152906001019060200180831161241857829003601f168201915b505050505081526020018260a001516bffffffffffffffffffffffff1681526020016005600086815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001826020015163ffffffff1667ffffffffffffffff1681526020018260c0015163ffffffff16815260200182608001516bffffffffffffffffffffffff16815260200182604001511515815260200160176000868152602001908152602001600020805461251290615c3f565b80601f016020809104026020016040519081016040528092919081815260200182805461253e90615c3f565b801561258b5780601f106125605761010080835404028352916020019161258b565b820191906000526020600020905b81548152906001019060200180831161256e57829003601f168201915b5050505050815250915050919050565b6125a3613e28565b601f865111156125df576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff8416612619576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806126385750612630846003615b7b565b60ff16865111155b1561266f576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b547001000000000000000000000000000000009091046bffffffffffffffffffffffff169060005b816bffffffffffffffffffffffff16811015612704576126f1600b82815481106126c8576126c8615d59565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484613ea9565b50806126fc81615c93565b91505061269c565b5060008060005b836bffffffffffffffffffffffff1681101561280d57600a818154811061273457612734615d59565b600091825260209091200154600b805473ffffffffffffffffffffffffffffffffffffffff9092169450908290811061276f5761276f615d59565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff868116845260098352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600890925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061280581615c93565b91505061270b565b5061281a600a60006147d8565b612826600b60006147d8565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015612baa57600960008e838151811061286b5761286b615d59565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff16156128d6576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600960008f848151811061290757612907615d59565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c90829081106129af576129af615d59565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff81166000908152600883526040908190208151608081018352905460ff80821615801584526101008304909116958301959095526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015294509250612a74576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff871660009081526008909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580612ba281615c93565b91505061284c565b50508a51612bc09150600a9060208d01906147f6565b508851612bd490600b9060208c01906147f6565b50600087806020019051810190612beb9190614fa4565b60125460c082015191925063ffffffff640100000000909104811691161015612c40576040517f39abc10400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125460e082015163ffffffff74010000000000000000000000000000000000000000909204821691161015612ca2576040517f1fa9bdcb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125461010082015163ffffffff7801000000000000000000000000000000000000000000000000909204821691161015612d09576040517fd1d5faa800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061012001604052808a60ff168152602001826000015163ffffffff168152602001826020015163ffffffff168152602001826060015162ffffff168152602001826080015161ffff168152602001600015158152602001600015158152602001866bffffffffffffffffffffffff168152602001600063ffffffff16815250600f60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160056101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160096101000a81548162ffffff021916908362ffffff160217905550608082015181600001600c6101000a81548161ffff021916908361ffff16021790555060a082015181600001600e6101000a81548160ff02191690831515021790555060c082015181600001600f6101000a81548160ff02191690831515021790555060e08201518160000160106101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061010082015181600001601c6101000a81548163ffffffff021916908363ffffffff1602179055509050506040518061016001604052808260a001516bffffffffffffffffffffffff16815260200182610160015173ffffffffffffffffffffffffffffffffffffffff168152602001601060010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200182610180015173ffffffffffffffffffffffffffffffffffffffff168152602001826040015163ffffffff1681526020018260c0015163ffffffff168152602001601060020160089054906101000a900463ffffffff1663ffffffff1681526020016010600201600c9054906101000a900463ffffffff1663ffffffff168152602001601060020160109054906101000a900463ffffffff1663ffffffff1681526020018260e0015163ffffffff16815260200182610100015163ffffffff16815250601060008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160020160006101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160020160046101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600201600c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160106101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160146101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160186101000a81548163ffffffff021916908363ffffffff1602179055509050508061012001516013819055508061014001516014819055506000601060020160109054906101000a900463ffffffff16905061326f613c52565b601280547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000063ffffffff9384160217808255600192600c916132d69185916c01000000000000000000000000900416615a68565b92506101000a81548163ffffffff021916908363ffffffff16021790555061332046306010600201600c9054906101000a900463ffffffff1663ffffffff168f8f8f8f8f8f6140d0565b600e819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600e546010600201600c9054906101000a900463ffffffff168f8f8f8f8f8f60405161337c9998979695949392919061589e565b60405180910390a1505050505050505050505050565b61339a613e28565b6111e28161417a565b600060606000806000806133b561364e565b6133bd610a7d565b91939550919395565b3660008037600080366000845af43d6000803e8080156133e5573d6000f35b3d6000fd5b600061226a825490565b60006134008383614270565b9392505050565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561348e57600080fd5b505afa1580156134a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134c6919061527d565b50945090925050506000811315806134dd57508142105b806134fe57508280156134fe57506134f58242615bcf565b8463ffffffff16105b1561350d576013549550613511565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561357757600080fd5b505afa15801561358b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135af919061527d565b50945090925050506000811315806135c657508142105b806135e757508280156135e757506135de8242615bcf565b8463ffffffff16105b156135f65760145494506135fa565b8094505b50505050915091565b60008061361486896000015161429a565b905060008061362f8a8a63ffffffff16858a8a60018b6142de565b909250905061363e8183615aac565b93505050505b9695505050505050565b3215610aa6576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f5460009081906f01000000000000000000000000000000900460ff16156136db576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f010000000000000000000000000000001790555a90506000634585e33b60e01b84604051602401613733919061557b565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506137ab856000015163ffffffff168660600151836146bd565b92505a6137b89083615bcf565b915050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff16905590939092509050565b6138176040518060800160405280600081526020016000815260200160608152602001606081525090565b600080600080858060200190518101906138319190615140565b93509350935093508051825114613874576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051608081018252948552602085019390935291830152606082015292915050565b60008260c0015163ffffffff16846000015163ffffffff1610156138e95760405185907f5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a8990600090a25060006139e1565b602084015184516138ff9063ffffffff16614709565b146139375760405185907f561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc1390600090a25060006139e1565b61393f613c52565b836020015163ffffffff16116139825760405185907fd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f90600090a25060006139e1565b816bffffffffffffffffffffffff168360a001516bffffffffffffffffffffffff1610156139dd5760405185907f7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb9690600090a25060006139e1565b5060015b949350505050565b600087876040516139fb929190615496565b604051908190038120613a12918b90602001615561565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b88811015613be957600185878360208110613a7e57613a7e615d59565b613a8b91901a601b615a87565b8c8c85818110613a9d57613a9d615d59565b905060200201358b8b86818110613ab657613ab6615d59565b9050602002013560405160008152602001604052604051613af3949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015613b15573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526009602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955093509050613bc3576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b840193508080613be190615c93565b915050613a61565b50827e01010101010101010101010101010101010101010101010101010101010101841614613c44576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115613c8857613c88615d2a565b1415613d1257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613cd557600080fd5b505afa158015613ce9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d0d9190614f55565b905090565b504390565b6000613d23838361429a565b90508084101561340057509192915050565b600080613d508887608001518860a0015188888860016142de565b90925090506000613d618284615aac565b600089815260046020526040902060010180549192508291600c90613da59084906c0100000000000000000000000090046bffffffffffffffffffffffff16615be6565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092613dee91859116615aac565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610aa6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611277565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e01000000000000000000000000000090920416606082018190528290613f349086615be6565b90506000613f428583615ae7565b90508083604001818151613f569190615aac565b6bffffffffffffffffffffffff9081169091528716606085015250613f7b8582615ba4565b613f859083615be6565b60118054600090613fa59084906bffffffffffffffffffffffff16615aac565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b6000808a8a8a8a8a8a8a8a8a6040516020016140f4999897969594939291906157f9565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156141fa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611277565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061428757614287615d59565b9060005260206000200154905092915050565b60006142ad63ffffffff84166014615b12565b6142b8836001615a87565b6142c79060ff16611d4c615b12565b6142d49062011170615a50565b6134009190615a50565b6000806000896080015161ffff16876142f79190615b12565b90508380156143055750803a105b1561430d57503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561434357614343615d2a565b14156144c65760408051600081526020810190915285156143a257600036604051806080016040528060488152602001615e166048913960405160200161438c939291906154a6565b604051602081830303815290604052905061441e565b6012546143d2907801000000000000000000000000000000000000000000000000900463ffffffff166004615b4f565b63ffffffff1667ffffffffffffffff8111156143f0576143f0615d88565b6040519080825280601f01601f19166020018201604052801561441a576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9061446e90849060040161557b565b60206040518083038186803b15801561448657600080fd5b505afa15801561449a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144be9190614f55565b915050614582565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156144fa576144fa615d2a565b141561458257606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561454757600080fd5b505afa15801561455b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061457f9190614f55565b90505b8461459e57808b6080015161ffff1661459b9190615b12565b90505b6145ac61ffff871682615ad3565b9050600087826145bc8c8e615a50565b6145c69086615b12565b6145d09190615a50565b6145e290670de0b6b3a7640000615b12565b6145ec9190615ad3565b905060008c6040015163ffffffff1664e8d4a5100061460b9190615b12565b898e6020015163ffffffff16858f886146249190615b12565b61462e9190615a50565b61463c90633b9aca00615b12565b6146469190615b12565b6146509190615ad3565b61465a9190615a50565b90506b033b2e3c9fd0803ce80000006146738284615a50565b11156146ab576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60005a6113888110156146cf57600080fd5b6113888103905084604082048203116146e757600080fd5b50823b6146f357600080fd5b60008083516020850160008789f1949350505050565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111561473f5761473f615d2a565b14156147ce576040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101839052606490632b407a829060240160206040518083038186803b15801561479657600080fd5b505afa1580156147aa573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061226a9190614f55565b504090565b919050565b50805460008255906000526020600020908101906111e29190614880565b828054828255906000526020600020908101928215614870579160200282015b8281111561487057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614816565b5061487c929150614880565b5090565b5b8082111561487c5760008155600101614881565b80516147d381615dc7565b60008083601f8401126148b257600080fd5b50813567ffffffffffffffff8111156148ca57600080fd5b6020830191508360208260051b85010111156148e557600080fd5b9250929050565b600082601f8301126148fd57600080fd5b8135602061491261490d836159c0565b615971565b80838252828201915082860187848660051b890101111561493257600080fd5b60005b8581101561495a57813561494881615dc7565b84529284019290840190600101614935565b5090979650505050505050565b600082601f83011261497857600080fd5b8151602061498861490d836159c0565b80838252828201915082860187848660051b89010111156149a857600080fd5b60005b8581101561495a57815167ffffffffffffffff808211156149cb57600080fd5b818a0191506060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848e03011215614a0357600080fd5b614a0b615924565b88840151614a1881615de9565b81526040848101518a830152918401519183831115614a3657600080fd5b82850194508d603f860112614a4a57600080fd5b898501519350614a5c61490d856159e4565b92508383528d81858701011115614a7257600080fd5b614a81848b8501838801615c13565b8101919091528652505092840192908401906001016149ab565b60008083601f840112614aad57600080fd5b50813567ffffffffffffffff811115614ac557600080fd5b6020830191508360208285010111156148e557600080fd5b600082601f830112614aee57600080fd5b8135614afc61490d826159e4565b818152846020838601011115614b1157600080fd5b816020850160208301376000918101602001919091529392505050565b805161ffff811681146147d357600080fd5b805162ffffff811681146147d357600080fd5b80516147d381615de9565b803567ffffffffffffffff811681146147d357600080fd5b803560ff811681146147d357600080fd5b805169ffffffffffffffffffff811681146147d357600080fd5b80516147d381615dfb565b600060208284031215614bbe57600080fd5b813561340081615dc7565b60008060408385031215614bdc57600080fd5b8235614be781615dc7565b91506020830135614bf781615dc7565b809150509250929050565b60008060408385031215614c1557600080fd5b8235614c2081615dc7565b9150602083013560048110614bf757600080fd5b60008060008060608587031215614c4a57600080fd5b8435614c5581615dc7565b935060208501359250604085013567ffffffffffffffff811115614c7857600080fd5b614c8487828801614a9b565b95989497509550505050565b600080600080600080600060a0888a031215614cab57600080fd5b8735614cb681615dc7565b96506020880135614cc681615de9565b95506040880135614cd681615dc7565b9450606088013567ffffffffffffffff80821115614cf357600080fd5b614cff8b838c01614a9b565b909650945060808a0135915080821115614d1857600080fd5b50614d258a828b01614a9b565b989b979a50959850939692959293505050565b60008060208385031215614d4b57600080fd5b823567ffffffffffffffff811115614d6257600080fd5b614d6e858286016148a0565b90969095509350505050565b60008060008060008060c08789031215614d9357600080fd5b863567ffffffffffffffff80821115614dab57600080fd5b614db78a838b016148ec565b97506020890135915080821115614dcd57600080fd5b614dd98a838b016148ec565b9650614de760408a01614b76565b95506060890135915080821115614dfd57600080fd5b614e098a838b01614add565b9450614e1760808a01614b5e565b935060a0890135915080821115614e2d57600080fd5b50614e3a89828a01614add565b9150509295509295509295565b60008060008060008060008060e0898b031215614e6357600080fd5b606089018a811115614e7457600080fd5b8998503567ffffffffffffffff80821115614e8e57600080fd5b614e9a8c838d01614a9b565b909950975060808b0135915080821115614eb357600080fd5b614ebf8c838d016148a0565b909750955060a08b0135915080821115614ed857600080fd5b50614ee58b828c016148a0565b999c989b50969995989497949560c00135949350505050565b600080600060408486031215614f1357600080fd5b833567ffffffffffffffff811115614f2a57600080fd5b614f36868287016148a0565b9094509250506020840135614f4a81615dc7565b809150509250925092565b600060208284031215614f6757600080fd5b5051919050565b60008060208385031215614f8157600080fd5b823567ffffffffffffffff811115614f9857600080fd5b614d6e85828601614a9b565b60006101a08284031215614fb757600080fd5b614fbf61594d565b614fc883614b53565b8152614fd660208401614b53565b6020820152614fe760408401614b53565b6040820152614ff860608401614b40565b606082015261500960808401614b2e565b608082015261501a60a08401614ba1565b60a082015261502b60c08401614b53565b60c082015261503c60e08401614b53565b60e082015261010061504f818501614b53565b9082015261012083810151908201526101408084015190820152610160615077818501614895565b90820152610180615089848201614895565b908201529392505050565b6000602082840312156150a657600080fd5b5035919050565b600080604083850312156150c057600080fd5b823591506020830135614bf781615dc7565b6000806000604084860312156150e757600080fd5b83359250602084013567ffffffffffffffff81111561510557600080fd5b61511186828701614a9b565b9497909650939450505050565b6000806040838503121561513157600080fd5b50508035926020909101359150565b6000806000806080858703121561515657600080fd5b845193506020808601519350604086015167ffffffffffffffff8082111561517d57600080fd5b818801915088601f83011261519157600080fd5b815161519f61490d826159c0565b8082825285820191508585018c878560051b88010111156151bf57600080fd5b600095505b838610156151e25780518352600195909501949186019186016151c4565b5060608b015190975094505050808311156151fc57600080fd5b505061520a87828801614967565b91505092959194509250565b6000806040838503121561522957600080fd5b823591506020830135614bf781615de9565b6000806040838503121561524e57600080fd5b823591506020830135614bf781615dfb565b60006020828403121561527257600080fd5b813561340081615de9565b600080600080600060a0868803121561529557600080fd5b61529e86614b87565b94506020860151935060408601519250606086015191506152c160808701614b87565b90509295509295909350565b600081518084526020808501945080840160005b8381101561531357815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016152e1565b509495945050505050565b60008151808452615336816020860160208601615c13565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b805163ffffffff1682526020810151615389602084018263ffffffff169052565b5060408101516153a1604084018263ffffffff169052565b5060608101516153b8606084018262ffffff169052565b5060808101516153ce608084018261ffff169052565b5060a08101516153ee60a08401826bffffffffffffffffffffffff169052565b5060c081015161540660c084018263ffffffff169052565b5060e081015161541e60e084018263ffffffff169052565b506101008181015163ffffffff8116848301525050610120818101519083015261014080820151908301526101608082015173ffffffffffffffffffffffffffffffffffffffff81168285015250506101808181015173ffffffffffffffffffffffffffffffffffffffff8116848301525b50505050565b8183823760009101908152919050565b8284823760008382016000815283516154c3818360208801615c13565b0195945050505050565b6020808252825182820181905260009190848201906040850190845b81811015615505578351835292840192918401916001016154e9565b50909695505050505050565b861515815260c06020820152600061552c60c083018861531e565b90506007861061553e5761553e615d2a565b8560408301528460608301528360808301528260a0830152979650505050505050565b828152608081016060836020840137600081529392505050565b602081526000613400602083018461531e565b60208101600483106155a2576155a2615d2a565b91905290565b602081016155a283615db7565b855163ffffffff168152600061034060208801516155e360208501826bffffffffffffffffffffffff169052565b5060408801516040840152606088015161560d60608501826bffffffffffffffffffffffff169052565b506080880151608084015260a088015161562f60a085018263ffffffff169052565b5060c088015161564760c085018263ffffffff169052565b5060e088015160e08401526101008089015161566a8286018263ffffffff169052565b505061012088810151151590840152615687610140840188615368565b806102e084015261569a818401876152cd565b90508281036103008401526156af81866152cd565b91505061364461032083018460ff169052565b602081526156e960208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151615702604084018263ffffffff169052565b50604083015161014080606085015261571f61016085018361531e565b9150606085015161574060808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e08501516101006157ac818701836bffffffffffffffffffffffff169052565b86015190506101206157c18682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001838701529050613644838261531e565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526158408285018b6152cd565b91508382036080850152615854828a6152cd565b915060ff881660a085015283820360c0850152615871828861531e565b90861660e0850152838103610100850152905061588e818561531e565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526158ce8184018a6152cd565b905082810360808401526158e281896152cd565b905060ff871660a084015282810360c08401526158ff818761531e565b905067ffffffffffffffff851660e084015282810361010084015261588e818561531e565b6040516060810167ffffffffffffffff8111828210171561594757615947615d88565b60405290565b6040516101a0810167ffffffffffffffff8111828210171561594757615947615d88565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156159b8576159b8615d88565b604052919050565b600067ffffffffffffffff8211156159da576159da615d88565b5060051b60200190565b600067ffffffffffffffff8211156159fe576159fe615d88565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600061ffff808316818516808303821115615a4757615a47615ccc565b01949350505050565b60008219821115615a6357615a63615ccc565b500190565b600063ffffffff808316818516808303821115615a4757615a47615ccc565b600060ff821660ff84168060ff03821115615aa457615aa4615ccc565b019392505050565b60006bffffffffffffffffffffffff808316818516808303821115615a4757615a47615ccc565b600082615ae257615ae2615cfb565b500490565b60006bffffffffffffffffffffffff80841680615b0657615b06615cfb565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615b4a57615b4a615ccc565b500290565b600063ffffffff80831681851681830481118215151615615b7257615b72615ccc565b02949350505050565b600060ff821660ff84168160ff0481118215151615615b9c57615b9c615ccc565b029392505050565b60006bffffffffffffffffffffffff80831681851681830481118215151615615b7257615b72615ccc565b600082821015615be157615be1615ccc565b500390565b60006bffffffffffffffffffffffff83811690831681811015615c0b57615c0b615ccc565b039392505050565b60005b83811015615c2e578181015183820152602001615c16565b838111156154905750506000910152565b600181811c90821680615c5357607f821691505b60208210811415615c8d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615cc557615cc5615ccc565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600381106111e2576111e2615d2a565b73ffffffffffffffffffffffffffffffffffffffff811681146111e257600080fd5b63ffffffff811681146111e257600080fd5b6bffffffffffffffffffffffff811681146111e257600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

// KeeperRegistry20ABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistry20MetaData.ABI instead.
var KeeperRegistry20ABI = KeeperRegistry20MetaData.ABI

// KeeperRegistry20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistry20MetaData.Bin instead.
var KeeperRegistry20Bin = KeeperRegistry20MetaData.Bin

// DeployKeeperRegistry20 deploys a new Ethereum contract, binding an instance of KeeperRegistry20 to it.
func DeployKeeperRegistry20(auth *bind.TransactOpts, backend bind.ContractBackend, keeperRegistryLogic common.Address) (common.Address, *types.Transaction, *KeeperRegistry20, error) {
	parsed, err := KeeperRegistry20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistry20Bin), backend, keeperRegistryLogic)
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

// GetFastGasFeedAddress is a free data retrieval call binding the contract method 0x6709d0e5.
//
// Solidity: function getFastGasFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFastGasFeedAddress is a free data retrieval call binding the contract method 0x6709d0e5.
//
// Solidity: function getFastGasFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetFastGasFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetFastGasFeedAddress is a free data retrieval call binding the contract method 0x6709d0e5.
//
// Solidity: function getFastGasFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetFastGasFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetKeeperRegistryLogicAddress is a free data retrieval call binding the contract method 0x572e05e1.
//
// Solidity: function getKeeperRegistryLogicAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetKeeperRegistryLogicAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getKeeperRegistryLogicAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetKeeperRegistryLogicAddress is a free data retrieval call binding the contract method 0x572e05e1.
//
// Solidity: function getKeeperRegistryLogicAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetKeeperRegistryLogicAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetKeeperRegistryLogicAddress(&_KeeperRegistry20.CallOpts)
}

// GetKeeperRegistryLogicAddress is a free data retrieval call binding the contract method 0x572e05e1.
//
// Solidity: function getKeeperRegistryLogicAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetKeeperRegistryLogicAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetKeeperRegistryLogicAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkAddress is a free data retrieval call binding the contract method 0xca30e603.
//
// Solidity: function getLinkAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkAddress is a free data retrieval call binding the contract method 0xca30e603.
//
// Solidity: function getLinkAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkAddress is a free data retrieval call binding the contract method 0xca30e603.
//
// Solidity: function getLinkAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkNativeFeedAddress is a free data retrieval call binding the contract method 0xb10b673c.
//
// Solidity: function getLinkNativeFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkNativeFeedAddress is a free data retrieval call binding the contract method 0xb10b673c.
//
// Solidity: function getLinkNativeFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkNativeFeedAddress is a free data retrieval call binding the contract method 0xb10b673c.
//
// Solidity: function getLinkNativeFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x0e08ae84.
//
// Solidity: function getMaxPaymentForGas(uint32 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x0e08ae84.
//
// Solidity: function getMaxPaymentForGas(uint32 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMaxPaymentForGas(&_KeeperRegistry20.CallOpts, gasLimit)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x0e08ae84.
//
// Solidity: function getMaxPaymentForGas(uint32 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
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

// GetPaymentModel is a free data retrieval call binding the contract method 0xf1570141.
//
// Solidity: function getPaymentModel() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetPaymentModel(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getPaymentModel")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPaymentModel is a free data retrieval call binding the contract method 0xf1570141.
//
// Solidity: function getPaymentModel() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetPaymentModel() (uint8, error) {
	return _KeeperRegistry20.Contract.GetPaymentModel(&_KeeperRegistry20.CallOpts)
}

// GetPaymentModel is a free data retrieval call binding the contract method 0xf1570141.
//
// Solidity: function getPaymentModel() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetPaymentModel() (uint8, error) {
	return _KeeperRegistry20.Contract.GetPaymentModel(&_KeeperRegistry20.CallOpts)
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

// GetSignerInfo is a free data retrieval call binding the contract method 0xed56b3e1.
//
// Solidity: function getSignerInfo(address query) view returns(bool active, uint8 index)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (struct {
	Active bool
	Index  uint8
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(struct {
		Active bool
		Index  uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetSignerInfo is a free data retrieval call binding the contract method 0xed56b3e1.
//
// Solidity: function getSignerInfo(address query) view returns(bool active, uint8 index)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetSignerInfo(query common.Address) (struct {
	Active bool
	Index  uint8
}, error) {
	return _KeeperRegistry20.Contract.GetSignerInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetSignerInfo is a free data retrieval call binding the contract method 0xed56b3e1.
//
// Solidity: function getSignerInfo(address query) view returns(bool active, uint8 index)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetSignerInfo(query common.Address) (struct {
	Active bool
	Index  uint8
}, error) {
	return _KeeperRegistry20.Contract.GetSignerInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint96,uint256,uint32,uint32,bytes32,uint32,bool) state, (uint32,uint32,uint32,uint24,uint16,uint96,uint32,uint32,uint32,uint256,uint256,address,address) config, address[] signers, address[] transmitters, uint8 f)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetState(opts *bind.CallOpts) (struct {
	State        State2_0
	Config       OnchainConfig2_0
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		State        State2_0
		Config       OnchainConfig2_0
		Signers      []common.Address
		Transmitters []common.Address
		F            uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State2_0)).(*State2_0)
	outstruct.Config = *abi.ConvertType(out[1], new(OnchainConfig2_0)).(*OnchainConfig2_0)
	outstruct.Signers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.Transmitters = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	outstruct.F = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint96,uint256,uint32,uint32,bytes32,uint32,bool) state, (uint32,uint32,uint32,uint24,uint16,uint96,uint32,uint32,uint32,uint256,uint256,address,address) config, address[] signers, address[] transmitters, uint8 f)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetState() (struct {
	State        State2_0
	Config       OnchainConfig2_0
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}, error) {
	return _KeeperRegistry20.Contract.GetState(&_KeeperRegistry20.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint96,uint256,uint32,uint32,bytes32,uint32,bool) state, (uint32,uint32,uint32,uint24,uint16,uint96,uint32,uint32,uint32,uint256,uint256,address,address) config, address[] signers, address[] transmitters, uint8 f)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetState() (struct {
	State        State2_0
	Config       OnchainConfig2_0
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}, error) {
	return _KeeperRegistry20.Contract.GetState(&_KeeperRegistry20.CallOpts)
}

// GetTransmitterInfo is a free data retrieval call binding the contract method 0x421d183b.
//
// Solidity: function getTransmitterInfo(address query) view returns(bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getTransmitterInfo", query)

	outstruct := new(struct {
		Active        bool
		Index         uint8
		Balance       *big.Int
		LastCollected *big.Int
		Payee         common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastCollected = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Payee = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GetTransmitterInfo is a free data retrieval call binding the contract method 0x421d183b.
//
// Solidity: function getTransmitterInfo(address query) view returns(bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetTransmitterInfo(query common.Address) (struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}, error) {
	return _KeeperRegistry20.Contract.GetTransmitterInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetTransmitterInfo is a free data retrieval call binding the contract method 0x421d183b.
//
// Solidity: function getTransmitterInfo(address query) view returns(bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetTransmitterInfo(query common.Address) (struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}, error) {
	return _KeeperRegistry20.Contract.GetTransmitterInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns((address,uint32,bytes,uint96,address,uint64,uint32,uint96,bool,bytes) upkeepInfo)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (UpkeepInfo, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(UpkeepInfo)).(*UpkeepInfo)

	return out0, err

}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns((address,uint32,bytes,uint96,address,uint64,uint32,uint96,bool,bytes) upkeepInfo)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetUpkeep(id *big.Int) (UpkeepInfo, error) {
	return _KeeperRegistry20.Contract.GetUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns((address,uint32,bytes,uint96,address,uint64,uint32,uint96,bool,bytes) upkeepInfo)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetUpkeep(id *big.Int) (UpkeepInfo, error) {
	return _KeeperRegistry20.Contract.GetUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_KeeperRegistry20 *KeeperRegistry20Caller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_KeeperRegistry20 *KeeperRegistry20Session) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDetails(&_KeeperRegistry20.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDetails(&_KeeperRegistry20.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Caller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Session) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry20.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry20.CallOpts)
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
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptPayeeship", transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptPayeeship(&_KeeperRegistry20.TransactOpts, transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptPayeeship(&_KeeperRegistry20.TransactOpts, transmitter)
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

// CheckUpkeep is a paid mutator transaction binding the contract method 0xf7d334ba.
//
// Solidity: function checkUpkeep(uint256 id) returns(bool upkeepNeeded, bytes performData, uint8 upkeepFailureReason, uint256 gasUsed, uint256 fastGasWei, uint256 linkNative)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "checkUpkeep", id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xf7d334ba.
//
// Solidity: function checkUpkeep(uint256 id) returns(bool upkeepNeeded, bytes performData, uint8 upkeepFailureReason, uint256 gasUsed, uint256 fastGasWei, uint256 linkNative)
func (_KeeperRegistry20 *KeeperRegistry20Session) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CheckUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xf7d334ba.
//
// Solidity: function checkUpkeep(uint256 id) returns(bool upkeepNeeded, bytes performData, uint8 upkeepFailureReason, uint256 gasUsed, uint256 fastGasWei, uint256 linkNative)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CheckUpkeep(&_KeeperRegistry20.TransactOpts, id)
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

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x6ded9eae.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData, bytes offchainConfig) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData, offchainConfig)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x6ded9eae.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData, bytes offchainConfig) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RegisterUpkeep(&_KeeperRegistry20.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x6ded9eae.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData, bytes offchainConfig) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RegisterUpkeep(&_KeeperRegistry20.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetConfig(&_KeeperRegistry20.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetConfig(&_KeeperRegistry20.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetPayees is a paid mutator transaction binding the contract method 0x3b9cce59.
//
// Solidity: function setPayees(address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setPayees", payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x3b9cce59.
//
// Solidity: function setPayees(address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPayees(&_KeeperRegistry20.TransactOpts, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x3b9cce59.
//
// Solidity: function setPayees(address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPayees(&_KeeperRegistry20.TransactOpts, payees)
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

// SetUpkeepOffchainConfig is a paid mutator transaction binding the contract method 0x8dcf0fe7.
//
// Solidity: function setUpkeepOffchainConfig(uint256 id, bytes config) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

// SetUpkeepOffchainConfig is a paid mutator transaction binding the contract method 0x8dcf0fe7.
//
// Solidity: function setUpkeepOffchainConfig(uint256 id, bytes config) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepOffchainConfig(&_KeeperRegistry20.TransactOpts, id, config)
}

// SetUpkeepOffchainConfig is a paid mutator transaction binding the contract method 0x8dcf0fe7.
//
// Solidity: function setUpkeepOffchainConfig(uint256 id, bytes config) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepOffchainConfig(&_KeeperRegistry20.TransactOpts, id, config)
}

// SimulatePerformUpkeep is a paid mutator transaction binding the contract method 0xaed2e929.
//
// Solidity: function simulatePerformUpkeep(uint256 id, bytes performData) returns(bool success, uint256 gasUsed)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

// SimulatePerformUpkeep is a paid mutator transaction binding the contract method 0xaed2e929.
//
// Solidity: function simulatePerformUpkeep(uint256 id, bytes performData) returns(bool success, uint256 gasUsed)
func (_KeeperRegistry20 *KeeperRegistry20Session) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SimulatePerformUpkeep(&_KeeperRegistry20.TransactOpts, id, performData)
}

// SimulatePerformUpkeep is a paid mutator transaction binding the contract method 0xaed2e929.
//
// Solidity: function simulatePerformUpkeep(uint256 id, bytes performData) returns(bool success, uint256 gasUsed)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SimulatePerformUpkeep(&_KeeperRegistry20.TransactOpts, id, performData)
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
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferPayeeship(&_KeeperRegistry20.TransactOpts, transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferPayeeship(&_KeeperRegistry20.TransactOpts, transmitter, proposed)
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

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes rawReport, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes rawReport, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Transmit(&_KeeperRegistry20.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes rawReport, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Transmit(&_KeeperRegistry20.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
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

// KeeperRegistry20CancelledUpkeepReportIterator is returned from FilterCancelledUpkeepReport and is used to iterate over the raw logs and unpacked data for CancelledUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20CancelledUpkeepReportIterator struct {
	Event *KeeperRegistry20CancelledUpkeepReport // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20CancelledUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20CancelledUpkeepReport)
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
		it.Event = new(KeeperRegistry20CancelledUpkeepReport)
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
func (it *KeeperRegistry20CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20CancelledUpkeepReport represents a CancelledUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20CancelledUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCancelledUpkeepReport is a free log retrieval operation binding the contract event 0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f.
//
// Solidity: event CancelledUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20CancelledUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchCancelledUpkeepReport is a free log subscription operation binding the contract event 0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f.
//
// Solidity: event CancelledUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20CancelledUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

// ParseCancelledUpkeepReport is a log parse operation binding the contract event 0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f.
//
// Solidity: event CancelledUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistry20CancelledUpkeepReport, error) {
	event := new(KeeperRegistry20CancelledUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry20ConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20ConfigSetIterator{contract: _KeeperRegistry20.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
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

// KeeperRegistry20InsufficientFundsUpkeepReportIterator is returned from FilterInsufficientFundsUpkeepReport and is used to iterate over the raw logs and unpacked data for InsufficientFundsUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20InsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistry20InsufficientFundsUpkeepReport // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20InsufficientFundsUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20InsufficientFundsUpkeepReport)
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
		it.Event = new(KeeperRegistry20InsufficientFundsUpkeepReport)
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
func (it *KeeperRegistry20InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20InsufficientFundsUpkeepReport represents a InsufficientFundsUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20InsufficientFundsUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterInsufficientFundsUpkeepReport is a free log retrieval operation binding the contract event 0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96.
//
// Solidity: event InsufficientFundsUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20InsufficientFundsUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchInsufficientFundsUpkeepReport is a free log subscription operation binding the contract event 0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96.
//
// Solidity: event InsufficientFundsUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20InsufficientFundsUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

// ParseInsufficientFundsUpkeepReport is a log parse operation binding the contract event 0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96.
//
// Solidity: event InsufficientFundsUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistry20InsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistry20InsufficientFundsUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

// KeeperRegistry20PayeesUpdatedIterator is returned from FilterPayeesUpdated and is used to iterate over the raw logs and unpacked data for PayeesUpdated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeesUpdatedIterator struct {
	Event *KeeperRegistry20PayeesUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20PayeesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PayeesUpdated)
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
		it.Event = new(KeeperRegistry20PayeesUpdated)
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
func (it *KeeperRegistry20PayeesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PayeesUpdated represents a PayeesUpdated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPayeesUpdated is a free log retrieval operation binding the contract event 0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725.
//
// Solidity: event PayeesUpdated(address[] transmitters, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistry20PayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeesUpdatedIterator{contract: _KeeperRegistry20.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

// WatchPayeesUpdated is a free log subscription operation binding the contract event 0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725.
//
// Solidity: event PayeesUpdated(address[] transmitters, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PayeesUpdated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

// ParsePayeesUpdated is a log parse operation binding the contract event 0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725.
//
// Solidity: event PayeesUpdated(address[] transmitters, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistry20PayeesUpdated, error) {
	event := new(KeeperRegistry20PayeesUpdated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry20PayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeeshipTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
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
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to)
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
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry20PayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeeshipTransferredIterator{contract: _KeeperRegistry20.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
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
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to)
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
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry20PaymentWithdrawnIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PaymentWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
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
// Solidity: event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry20PaymentWithdrawn, error) {
	event := new(KeeperRegistry20PaymentWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20ReorgedUpkeepReportIterator is returned from FilterReorgedUpkeepReport and is used to iterate over the raw logs and unpacked data for ReorgedUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20ReorgedUpkeepReportIterator struct {
	Event *KeeperRegistry20ReorgedUpkeepReport // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20ReorgedUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20ReorgedUpkeepReport)
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
		it.Event = new(KeeperRegistry20ReorgedUpkeepReport)
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
func (it *KeeperRegistry20ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20ReorgedUpkeepReport represents a ReorgedUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20ReorgedUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterReorgedUpkeepReport is a free log retrieval operation binding the contract event 0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13.
//
// Solidity: event ReorgedUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20ReorgedUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchReorgedUpkeepReport is a free log subscription operation binding the contract event 0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13.
//
// Solidity: event ReorgedUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20ReorgedUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

// ParseReorgedUpkeepReport is a log parse operation binding the contract event 0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13.
//
// Solidity: event ReorgedUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistry20ReorgedUpkeepReport, error) {
	event := new(KeeperRegistry20ReorgedUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20StaleUpkeepReportIterator is returned from FilterStaleUpkeepReport and is used to iterate over the raw logs and unpacked data for StaleUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20StaleUpkeepReportIterator struct {
	Event *KeeperRegistry20StaleUpkeepReport // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20StaleUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20StaleUpkeepReport)
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
		it.Event = new(KeeperRegistry20StaleUpkeepReport)
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
func (it *KeeperRegistry20StaleUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20StaleUpkeepReport represents a StaleUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20StaleUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterStaleUpkeepReport is a free log retrieval operation binding the contract event 0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89.
//
// Solidity: event StaleUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20StaleUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchStaleUpkeepReport is a free log subscription operation binding the contract event 0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89.
//
// Solidity: event StaleUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20StaleUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

// ParseStaleUpkeepReport is a log parse operation binding the contract event 0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89.
//
// Solidity: event StaleUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistry20StaleUpkeepReport, error) {
	event := new(KeeperRegistry20StaleUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20TransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the KeeperRegistry20 contract.
type KeeperRegistry20TransmittedIterator struct {
	Event *KeeperRegistry20Transmitted // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20TransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20Transmitted)
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
		it.Event = new(KeeperRegistry20Transmitted)
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
func (it *KeeperRegistry20TransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20TransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20Transmitted represents a Transmitted event raised by the KeeperRegistry20 contract.
type KeeperRegistry20Transmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterTransmitted(opts *bind.FilterOpts) (*KeeperRegistry20TransmittedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20TransmittedIterator{contract: _KeeperRegistry20.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20Transmitted) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20Transmitted)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseTransmitted(log types.Log) (*KeeperRegistry20Transmitted, error) {
	event := new(KeeperRegistry20Transmitted)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// KeeperRegistry20UpkeepOffchainConfigSetIterator is returned from FilterUpkeepOffchainConfigSet and is used to iterate over the raw logs and unpacked data for UpkeepOffchainConfigSet events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistry20UpkeepOffchainConfigSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry20UpkeepOffchainConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepOffchainConfigSet)
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
		it.Event = new(KeeperRegistry20UpkeepOffchainConfigSet)
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
func (it *KeeperRegistry20UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepOffchainConfigSet represents a UpkeepOffchainConfigSet event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpkeepOffchainConfigSet is a free log retrieval operation binding the contract event 0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850.
//
// Solidity: event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepOffchainConfigSetIterator{contract: _KeeperRegistry20.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepOffchainConfigSet is a free log subscription operation binding the contract event 0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850.
//
// Solidity: event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepOffchainConfigSet)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

// ParseUpkeepOffchainConfigSet is a log parse operation binding the contract event 0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850.
//
// Solidity: event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistry20UpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistry20UpkeepOffchainConfigSet)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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
	Id               *big.Int
	Success          bool
	CheckBlockNumber uint32
	GasUsed          *big.Int
	GasOverhead      *big.Int
	TotalPayment     *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPerformed is a free log retrieval operation binding the contract event 0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, uint32 checkBlockNumber, uint256 gasUsed, uint256 gasOverhead, uint96 totalPayment)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistry20UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepPerformedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, uint32 checkBlockNumber, uint256 gasUsed, uint256 gasOverhead, uint96 totalPayment)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
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

// ParseUpkeepPerformed is a log parse operation binding the contract event 0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, uint32 checkBlockNumber, uint256 gasUsed, uint256 gasOverhead, uint96 totalPayment)
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
