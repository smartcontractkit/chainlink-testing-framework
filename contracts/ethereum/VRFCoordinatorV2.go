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

// VRFCoordinatorV2FeeConfig is an auto generated low-level Go binding around an user-defined struct.
type VRFCoordinatorV2FeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}

// VRFCoordinatorV2RequestCommitment is an auto generated low-level Go binding around an user-defined struct.
type VRFCoordinatorV2RequestCommitment struct {
	BlockNum         uint64
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
}

// VRFProof is an auto generated low-level Go binding around an user-defined struct.
type VRFProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

// VRFCoordinatorV2MetaData contains all meta data concerning the VRFCoordinatorV2 contract.
var VRFCoordinatorV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structVRFCoordinatorV2.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"}],\"name\":\"getFeeTier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004942380380620049428339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000de565b5050506001600160a01b03928316608052821660a0521660c052620001f0565b336001600160a01b03821603620001385760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001a157600080fd5b919050565b600080600060608486031215620001bc57600080fd5b620001c78462000189565b9250620001d76020850162000189565b9150620001e76040850162000189565b90509250925092565b60805160a05160c0516146f162000251600039600081816104090152612bef0152600081816104f30152612f440152600081816102da01528181610fdb01528181611b11015281816122fa015281816123df01526128ed01526146f16000f3fe608060405234801561001057600080fd5b50600436106101db5760003560e01c80636f64f03f11610105578063ad1783611161009d578063ad178361146104ee578063af198b9714610515578063c3f909d414610540578063caf70c4a14610586578063d2f9f9a714610599578063d7ae1d30146105ac578063e72f6e30146105bf578063e82ad7d4146105d2578063f2fde38b146105f557600080fd5b80636f64f03f1461044b5780637341c10c1461045e57806379ba50971461047157806382359740146104795780638da5cb5b1461048c5780639f87fad71461049d578063a21a23e4146104b0578063a47c7696146104b8578063a4c0ed36146104db57600080fd5b8063356dac7111610178578063356dac711461030957806340d6bb82146103115780634cb48a541461032f5780635d3b1d30146103425780635fbbc0d21461035557806364d51a2a146103e957806366316d8d146103f1578063689c45171461040457806369bcdb7d1461042b57600080fd5b8062012291146101e057806302bcc5b61461020057806304c357cb1461021557806306bfa6371461022857806308821d581461024d57806312b583491461026057806315c48b8414610282578063181f5a771461029d5780631b6b6d23146102d5575b600080fd5b6101e8610608565b6040516101f793929190613b3e565b60405180910390f35b61021361020e366004613bb4565b610684565b005b610213610223366004613be6565b6106fb565b6005546001600160401b03165b6040516001600160401b0390911681526020016101f7565b61021361025b366004613c2a565b61083f565b600554600160401b90046001600160601b03165b6040519081526020016101f7565b61028a60c881565b60405161ffff90911681526020016101f7565b60408051808201825260168152750565246436f6f7264696e61746f72563220312e302e360541b602082015290516101f79190613c46565b6102fc7f000000000000000000000000000000000000000000000000000000000000000081565b6040516101f79190613c9b565b600a54610274565b61031a6101f481565b60405163ffffffff90911681526020016101f7565b61021361033d366004613d27565b6109d2565b610274610350366004613e46565b610c1e565b600c546040805163ffffffff8084168252600160201b840481166020830152600160401b8404811692820192909252600160601b830482166060820152600160801b8304909116608082015262ffffff600160a01b8304811660a0830152600160b81b8304811660c0830152600160d01b8304811660e0830152600160e81b909204909116610100820152610120016101f7565b61028a606481565b6102136103ff366004613ea4565b610ede565b6102fc7f000000000000000000000000000000000000000000000000000000000000000081565b610274610439366004613ee7565b60009081526009602052604090205490565b610213610459366004613f00565b61108b565b61021361046c366004613be6565b61118a565b61021361132f565b610213610487366004613bb4565b6113d9565b6000546001600160a01b03166102fc565b6102136104ab366004613be6565b611531565b610235611840565b6104cb6104c6366004613bb4565b6119e1565b6040516101f79493929190613f2b565b6102136104e9366004613fa9565b611adb565b6102fc7f000000000000000000000000000000000000000000000000000000000000000081565b61052861052336600461412e565b611caa565b6040516001600160601b0390911681526020016101f7565b600b546040805161ffff8316815263ffffffff62010000840481166020830152600160381b8404811692820192909252600160581b9092041660608201526080016101f7565b6102746105943660046141f6565b612063565b61031a6105a7366004613bb4565b612093565b6102136105ba366004613be6565b61220f565b6102136105cd366004614212565b6122d8565b6105e56105e0366004613bb4565b61249a565b60405190151581526020016101f7565b610213610603366004614212565b612645565b600b546007805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff1693919283919083018282801561067257602002820191906000526020600020905b81548152602001906001019080831161065e575b50505050509050925092509250909192565b61068c612656565b6001600160401b0381166000908152600360205260409020546001600160a01b03166106cb57604051630fb532db60e11b815260040160405180910390fd5b6001600160401b0381166000908152600360205260409020546106f89082906001600160a01b03166126ab565b50565b6001600160401b03821660009081526003602052604090205482906001600160a01b03168061073d57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b038216146107715780604051636c51fda960e11b81526004016107689190613c9b565b60405180910390fd5b600b54600160301b900460ff161561079c5760405163769dd35360e11b815260040160405180910390fd5b6001600160401b0384166000908152600360205260409020600101546001600160a01b03848116911614610839576001600160401b0384166000818152600360205260409081902060010180546001600160a01b0319166001600160a01b038716179055517f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be90610830903390879061422d565b60405180910390a25b50505050565b610847612656565b604080518082018252600091610876919084906002908390839080828437600092019190915250612063915050565b6000818152600660205260409020549091506001600160a01b0316806108b257604051631dfd6e1360e21b815260048101839052602401610768565b600082815260066020526040812080546001600160a01b03191690555b6007548110156109895782600782815481106108ed576108ed614247565b90600052602060002001540361097757600780546000919061091190600190614273565b8154811061092157610921614247565b90600052602060002001549050806007838154811061094257610942614247565b600091825260209091200155600780548061095f5761095f61428a565b60019003818190600052602060002001600090559055505b80610981816142a0565b9150506108cf565b50806001600160a01b03167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d836040516109c591815260200190565b60405180910390a2505050565b6109da612656565b60c861ffff87161115610a145760405163539c34bb60e11b815261ffff871660048201819052602482015260c86044820152606401610768565b60008213610a38576040516321ea67b360e11b815260048101839052602401610768565b6040805160a0808201835261ffff891680835263ffffffff89811660208086018290526000868801528a831660608088018290528b85166080988901819052600b805465ffffffffffff1916909717620100009095029490941768ffffffffffffffffff60301b1916600160381b90920263ffffffff60581b191691909117600160581b909302929092179093558651600c80549489015189890151938a0151978a0151968a015160c08b015160e08c01516101008d01519588166001600160401b031990991698909817600160201b938816939093029290921767ffffffffffffffff60401b1916600160401b9587169590950263ffffffff60601b191694909417600160601b988616989098029790971766ffffffffffffff60801b1916600160801b969094169590950262ffffff60a01b191692909217600160a01b62ffffff928316021765ffffffffffff60b81b1916600160b81b9582169590950262ffffff60d01b191694909417600160d01b92851692909202919091176001600160e81b0316600160e81b9390911692909202919091178155600a84905590517fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb291610c0e9189918991899189918991906142b9565b60405180910390a1505050505050565b600b54600090600160301b900460ff1615610c4c5760405163769dd35360e11b815260040160405180910390fd5b6001600160401b0385166000908152600360205260409020546001600160a01b0316610c8b57604051630fb532db60e11b815260040160405180910390fd5b3360009081526002602090815260408083206001600160401b03808a16855292528220541690819003610cd5578533604051637800cff360e11b8152600401610768929190614395565b600b5461ffff9081169086161080610cf1575060c861ffff8616115b15610d2857600b5460405163539c34bb60e11b815261ffff8088166004830152909116602482015260c86044820152606401610768565b600b5463ffffffff6201000090910481169085161115610d7657600b54604051637aebf00f60e11b815263ffffffff8087166004830152620100009092049091166024820152604401610768565b6101f463ffffffff84161115610daf576040516311ce1afb60e21b815263ffffffff841660048201526101f46024820152604401610768565b6000610dbc8260016143b7565b9050600080610dcd8a338b866129ec565b604080516020810184905243918101919091526001600160401b038c16606082015263ffffffff808b166080830152891660a08201523360c0820152919350915060e00160408051808303601f19018152828252805160209182012060008681526009835283902055848352820183905261ffff8a169082015263ffffffff80891660608301528716608082015233906001600160401b038b16908c907f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a97729060a00160405180910390a4503360009081526002602090815260408083206001600160401b03808d16855292529091208054919093166001600160401b03199091161790915591505095945050505050565b600b54600160301b900460ff1615610f095760405163769dd35360e11b815260040160405180910390fd5b336000908152600860205260409020546001600160601b0380831691161015610f4557604051631e9acf1760e31b815260040160405180910390fd5b3360009081526008602052604081208054839290610f6d9084906001600160601b03166143e2565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600560088282829054906101000a90046001600160601b0316610fb591906143e2565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb83836040518363ffffffff1660e01b815260040161102792919061440a565b6020604051808303816000875af1158015611046573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061106a919061442c565b61108757604051631e9acf1760e31b815260040160405180910390fd5b5050565b611093612656565b6040805180820182526000916110c2919084906002908390839080828437600092019190915250612063915050565b6000818152600660205260409020549091506001600160a01b0316156110fe57604051634a0b8fa760e01b815260048101829052602401610768565b600081815260066020908152604080832080546001600160a01b0319166001600160a01b0388169081179091556007805460018101825594527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b891016109c5565b6001600160401b03821660009081526003602052604090205482906001600160a01b0316806111cc57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b038216146111f75780604051636c51fda960e11b81526004016107689190613c9b565b600b54600160301b900460ff16156112225760405163769dd35360e11b815260040160405180910390fd5b6001600160401b0384166000908152600360205260409020600201546063190161125f576040516305a48e0f60e01b815260040160405180910390fd5b6001600160a01b03831660009081526002602090815260408083206001600160401b0380891685529252822054169003610839576001600160a01b03831660008181526002602081815260408084206001600160401b038a1680865290835281852080546001600160401b0319166001908117909155600384528286209094018054948501815585529190932090910180546001600160a01b031916909317909255517f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e090610830908690613c9b565b6001546001600160a01b031633146113825760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610768565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600b54600160301b900460ff16156114045760405163769dd35360e11b815260040160405180910390fd5b6001600160401b0381166000908152600360205260409020546001600160a01b031661144357604051630fb532db60e11b815260040160405180910390fd5b6001600160401b0381166000908152600360205260409020600101546001600160a01b031633146114ae576001600160401b0381166000908152600360205260409081902060010154905163d084e97560e01b8152610768916001600160a01b031690600401613c9b565b6001600160401b0381166000818152600360205260409081902080546001600160a01b031980821633908117845560019093018054909116905591516001600160a01b0390921692917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0916115259185919061422d565b60405180910390a25050565b6001600160401b03821660009081526003602052604090205482906001600160a01b03168061157357604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461159e5780604051636c51fda960e11b81526004016107689190613c9b565b600b54600160301b900460ff16156115c95760405163769dd35360e11b815260040160405180910390fd5b6001600160a01b03831660009081526002602090815260408083206001600160401b038089168552925282205416900361161a578383604051637800cff360e11b8152600401610768929190614395565b6001600160401b03841660009081526003602090815260408083206002018054825181850281018501909352808352919290919083018282801561168757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611669575b5050505050905060006001825161169e9190614273565b905060005b82518110156117c457856001600160a01b03168382815181106116c8576116c8614247565b60200260200101516001600160a01b0316036117b25760008383815181106116f2576116f2614247565b6020026020010151905080600360008a6001600160401b03166001600160401b03168152602001908152602001600020600201838154811061173657611736614247565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092556001600160401b038a16815260039091526040902060020180548061178a5761178a61428a565b600082815260209020810160001990810180546001600160a01b0319169055019055506117c4565b806117bc816142a0565b9150506116a3565b506001600160a01b03851660009081526002602090815260408083206001600160401b038a1680855292529182902080546001600160401b031916905590517f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b90611830908890613c9b565b60405180910390a2505050505050565b600b54600090600160301b900460ff161561186e5760405163769dd35360e11b815260040160405180910390fd5b600580546001600160401b03169060006118878361444e565b82546101009290920a6001600160401b038181021990931691831602179091556005541690506000806040519080825280602002602001820160405280156118d9578160200160208202803683370190505b50604080518082018252600080825260208083018281526001600160401b03888116808552600484528685209551865493516001600160601b039091166001600160a01b031994851617600160601b919093160291909117909455845160608101865233815280830184815281870188815295855260038452959093208351815483166001600160a01b03918216178255955160018201805490931696169590951790559151805194955090936119969260028501920190613a4c565b50905050816001600160401b03167f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf336040516119d39190613c9b565b60405180910390a250905090565b6001600160401b038116600090815260036020526040812054819081906060906001600160a01b0316611a2757604051630fb532db60e11b815260040160405180910390fd5b6001600160401b0380861660009081526004602090815260408083205460038352928190208054600290910180548351818602810186019094528084526001600160601b03861696600160601b909604909516946001600160a01b03909216939092918391830182828015611ac557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611aa7575b5050505050905093509350935093509193509193565b600b54600160301b900460ff1615611b065760405163769dd35360e11b815260040160405180910390fd5b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611b4f576040516344b0e3c360e01b815260040160405180910390fd5b60208114611b7057604051638129bbcd60e01b815260040160405180910390fd5b6000611b7e82840184613bb4565b6001600160401b0381166000908152600360205260409020549091506001600160a01b0316611bc057604051630fb532db60e11b815260040160405180910390fd5b6001600160401b038116600090815260046020526040812080546001600160601b031691869190611bf18385614474565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600560088282829054906101000a90046001600160601b0316611c399190614474565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550816001600160401b03167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8828784611c959190614496565b60408051928352602083019190915201611830565b600b54600090600160301b900460ff1615611cd85760405163769dd35360e11b815260040160405180910390fd5b60005a90506000806000611cec8787612a66565b9250925092506000866060015163ffffffff166001600160401b03811115611d1657611d16613cd5565b604051908082528060200260200182016040528015611d3f578160200160208202803683370190505b50905060005b876060015163ffffffff16811015611db35760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c828281518110611d9657611d96614247565b602090810291909101015280611dab816142a0565b915050611d45565b50600083815260096020526040808220829055518190631fe543e360e01b90611de290879086906024016144ae565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600b805460ff60301b1916600160301b179055908a015160808b0151919250600091611e479163ffffffff169084612ce4565b600b805460ff60301b191690556020808c0180516001600160401b03908116600090815260049093526040808420549251821684529092208054939450600160601b918290048316936001939192600c92611ea69286929004166143b7565b92506101000a8154816001600160401b0302191690836001600160401b031602179055506000611efb8a600b600001600b9054906101000a900463ffffffff1663ffffffff16611ef585612093565b3a612d32565b6020808e01516001600160401b03166000908152600490915260409020549091506001600160601b0380831691161015611f4857604051631e9acf1760e31b815260040160405180910390fd5b6020808d01516001600160401b031660009081526004909152604081208054839290611f7e9084906001600160601b03166143e2565b82546101009290920a6001600160601b0381810219909316918316021790915560008b8152600660209081526040808320546001600160a01b031683526008909152812080548594509092611fd591859116614474565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550877f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4888386604051612049939291909283526001600160601b039190911660208301521515604082015260600190565b60405180910390a299505050505050505050505b92915050565b600081604051602001612076919061451f565b604051602081830303815290604052805190602001209050919050565b6040805161012081018252600c5463ffffffff8082168352600160201b820481166020840152600160401b8204811693830193909352600160601b810483166060830152600160801b8104909216608082015262ffffff600160a01b8304811660a08301819052600160b81b8404821660c0840152600160d01b8404821660e0840152600160e81b909304166101008201526000916001600160401b0384161161213e575192915050565b826001600160401b03168160a0015162ffffff1610801561217157508060c0015162ffffff16836001600160401b031611155b15612180576020015192915050565b826001600160401b03168160c0015162ffffff161080156121b357508060e0015162ffffff16836001600160401b031611155b156121c2576040015192915050565b826001600160401b03168160e0015162ffffff161080156121f6575080610100015162ffffff16836001600160401b031611155b15612205576060015192915050565b6080015192915050565b6001600160401b03821660009081526003602052604090205482906001600160a01b03168061225157604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461227c5780604051636c51fda960e11b81526004016107689190613c9b565b600b54600160301b900460ff16156122a75760405163769dd35360e11b815260040160405180910390fd5b6122b08461249a565b156122ce57604051631685ecdd60e31b815260040160405180910390fd5b61083984846126ab565b6122e0612656565b6040516370a0823160e01b81526000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a082319061232f903090600401613c9b565b602060405180830381865afa15801561234c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612370919061452d565b600554909150600160401b90046001600160601b0316818111156123b1576040516354ced18160e11b81526004810182905260248101839052604401610768565b818110156124955760006123c58284614273565b60405163a9059cbb60e01b81529091506001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906124169087908590600401614546565b6020604051808303816000875af1158015612435573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612459919061442c565b507f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600848260405161248b929190614546565b60405180910390a1505b505050565b6001600160401b0381166000908152600360209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561252e57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612510575b505050505081525050905060005b81604001515181101561263b5760005b6007548110156126285760006125f16007838154811061256e5761256e614247565b90600052602060002001548560400151858151811061258f5761258f614247565b60200260200101518860026000896040015189815181106125b2576125b2614247565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081206001600160401b03808f16835293522054166129ec565b50600081815260096020526040902054909150156126155750600195945050505050565b5080612620816142a0565b91505061254c565b5080612633816142a0565b91505061253c565b5060009392505050565b61264d612656565b6106f881612e08565b6000546001600160a01b031633146126a95760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610768565b565b600b54600160301b900460ff16156126d65760405163769dd35360e11b815260040160405180910390fd5b6001600160401b0382166000908152600360209081526040808320815160608101835281546001600160a01b0390811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561276657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612748575b505050919092525050506001600160401b0380851660009081526004602090815260408083208151808301909252546001600160601b038116808352600160601b909104909416918101919091529293505b8360400151518110156128385760026000856040015183815181106127df576127df614247565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081206001600160401b038a168252909252902080546001600160401b031916905580612830816142a0565b9150506127b8565b506001600160401b038516600090815260036020526040812080546001600160a01b0319908116825560018201805490911690559061287a6002830182613ab1565b50506001600160401b038516600090815260046020526040902080546001600160a01b0319169055600580548291906008906128c7908490600160401b90046001600160601b03166143e2565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb85836001600160601b03166040518363ffffffff1660e01b8152600401612942929190614546565b6020604051808303816000875af1158015612961573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612985919061442c565b6129a257604051631e9acf1760e31b815260040160405180910390fd5b846001600160401b03167fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd4981585836040516129dd92919061440a565b60405180910390a25050505050565b6040805160208082018790526001600160a01b0395909516818301526001600160401b039384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b6000806000612a788560000151612063565b6000818152600660205260409020549093506001600160a01b031680612ab457604051631dfd6e1360e21b815260048101859052602401610768565b6080860151604051612ad3918691602001918252602082015260400190565b60408051601f198184030181529181528151602092830120600081815260099093529082205490945090819003612b1d57604051631b44092560e11b815260040160405180910390fd5b85516020808801516040808a015160608b015160808c01519251612b88968b9690959491019586526001600160401b03948516602087015292909316604085015263ffffffff90811660608501529190911660808301526001600160a01b031660a082015260c00190565b604051602081830303815290604052805190602001208114612bbd5760405163354a450b60e21b815260040160405180910390fd5b85516001600160401b03164080612c90578651604051631d2827a760e31b81526001600160401b0390911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d3890602401602060405180830381865afa158015612c3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c62919061452d565b905080612c9057865160405163175dadad60e01b81526001600160401b039091166004820152602401610768565b6000886080015182604051602001612cb2929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050612cd78982612eab565b9450505050509250925092565b60005a611388811015612cf657600080fd5b611388810390508460408204820311612d0e57600080fd5b50823b612d1a57600080fd5b60008083516020850160008789f190505b9392505050565b600080612d3d612f16565b905060008113612d63576040516321ea67b360e11b815260048101829052602401610768565b6000815a612d718989614496565b612d7b9190614273565b612d8d86670de0b6b3a764000061455f565b612d97919061455f565b612da19190614594565b90506000612dba63ffffffff871664e8d4a5100061455f565b9050612dd2816b033b2e3c9fd0803ce8000000614273565b821115612df25760405163e80fa38160e01b815260040160405180910390fd5b612dfc8183614496565b98975050505050505050565b336001600160a01b03821603612e5a5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610768565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612edf8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151612ff1565b60038360200151604051602001612ef79291906145a8565b60408051601f1981840301815291905280516020909101209392505050565b600b5460408051633fabe5a360e21b81529051600092600160381b900463ffffffff169182151591849182917f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169163feaf968c9160048083019260a09291908290030181865afa158015612f97573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612fbb91906145d6565b509450909250849150508015612fdf5750612fd68242614273565b8463ffffffff16105b15612fe95750600a545b949350505050565b612ffa8961320d565b6130435760405162461bcd60e51b815260206004820152601a6024820152797075626c6963206b6579206973206e6f74206f6e20637572766560301b6044820152606401610768565b61304c8861320d565b6130905760405162461bcd60e51b815260206004820152601560248201527467616d6d61206973206e6f74206f6e20637572766560581b6044820152606401610768565b6130998361320d565b6130e55760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610768565b6130ee8261320d565b61313a5760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610768565b613146878a88876132d0565b61318e5760405162461bcd60e51b81526020600482015260196024820152786164647228632a706b2b732a6729213d5f755769746e65737360381b6044820152606401610768565b600061319a8a876133f3565b905060006131ad898b878b868989613457565b905060006131be838d8d8a8661356a565b9050808a146131ff5760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610768565b505050505050505050505050565b80516000906401000003d0191161325b5760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420782d6f7264696e61746560701b6044820152606401610768565b60208201516401000003d019116132a95760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420792d6f7264696e61746560701b6044820152606401610768565b60208201516401000003d0199080096132c98360005b60200201516135aa565b1492915050565b60006001600160a01b0382166133165760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610768565b60208401516000906001161561332d57601c613330565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156133cb573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6133fb613acf565b613428600184846040516020016134149392919061461a565b6040516020818303038152906040526135ce565b90505b6134348161320d565b61205d5780516040805160208101929092526134509101613414565b905061342b565b61345f613acf565b825186516401000003d01991829006919006036134be5760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610768565b6134c987898861361c565b61350e5760405162461bcd60e51b8152602060048201526016602482015275119a5c9cdd081b5d5b0818da1958dac819985a5b195960521b6044820152606401610768565b61351984868561361c565b61355f5760405162461bcd60e51b815260206004820152601760248201527614d958dbdb99081b5d5b0818da1958dac819985a5b1959604a1b6044820152606401610768565b612dfc868484613747565b6000600286868685876040516020016135889695949392919061463b565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b6135d6613acf565b6135df8261380a565b81526135f46135ef8260006132bf565b613845565b6020820181905260029006600103613617576020810180516401000003d0190390525b919050565b60008260000361365c5760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610768565b8351602085015160009061367290600290614695565b1561367e57601c613681565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa1580156136f3573d6000803e3d6000fd5b50505060206040510351905060008660405160200161371291906146a9565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b61374f613acf565b83516020808601518551918601516000938493849361377093909190613865565b919450925090506401000003d0198582096001146137cc5760405162461bcd60e51b815260206004820152601960248201527834b73b2d1036bab9ba1031329034b73b32b939b29037b3103d60391b6044820152606401610768565b60405180604001604052806401000003d019806137eb576137eb61457e565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d019811061361757604080516020808201939093528151808203840181529082019091528051910120613812565b600061205d82600261385e6401000003d0196001614496565b901c613945565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a08905060006138a5838385856139df565b90985090506138b688828e88613a03565b90985090506138c788828c87613a03565b909850905060006138da8d878b85613a03565b90985090506138eb888286866139df565b90985090506138fc88828e89613a03565b9098509050818114613931576401000003d019818a0998506401000003d01982890997506401000003d0198183099650613935565b8196505b5050505050509450945094915050565b600080613950613aed565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152613982613b0b565b60208160c0846005600019fa9250826000036139d55760405162461bcd60e51b81526020600482015260126024820152716269674d6f64457870206661696c7572652160701b6044820152606401610768565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215613aa1579160200282015b82811115613aa157825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613a6c565b50613aad929150613b29565b5090565b50805460008255906000526020600020908101906106f89190613b29565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115613aad5760008155600101613b2a565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015613b8f57845183529383019391830191600101613b73565b509098975050505050505050565b80356001600160401b038116811461361757600080fd5b600060208284031215613bc657600080fd5b612d2b82613b9d565b80356001600160a01b038116811461361757600080fd5b60008060408385031215613bf957600080fd5b613c0283613b9d565b9150613c1060208401613bcf565b90509250929050565b806040810183101561205d57600080fd5b600060408284031215613c3c57600080fd5b612d2b8383613c19565b600060208083528351808285015260005b81811015613c7357858101830151858201604001528201613c57565b81811115613c85576000604083870101525b50601f01601f1916929092016040019392505050565b6001600160a01b0391909116815260200190565b803561ffff8116811461361757600080fd5b803563ffffffff8116811461361757600080fd5b634e487b7160e01b600052604160045260246000fd5b60405161012081016001600160401b0381118282101715613d0e57613d0e613cd5565b60405290565b803562ffffff8116811461361757600080fd5b6000806000806000808688036101c0811215613d4257600080fd5b613d4b88613caf565b9650613d5960208901613cc1565b9550613d6760408901613cc1565b9450613d7560608901613cc1565b93506080880135925061012080609f1983011215613d9257600080fd5b613d9a613ceb565b9150613da860a08a01613cc1565b8252613db660c08a01613cc1565b6020830152613dc760e08a01613cc1565b6040830152610100613dda818b01613cc1565b6060840152613dea828b01613cc1565b6080840152613dfc6101408b01613d14565b60a0840152613e0e6101608b01613d14565b60c0840152613e206101808b01613d14565b60e0840152613e326101a08b01613d14565b818401525050809150509295509295509295565b600080600080600060a08688031215613e5e57600080fd5b85359450613e6e60208701613b9d565b9350613e7c60408701613caf565b9250613e8a60608701613cc1565b9150613e9860808701613cc1565b90509295509295909350565b60008060408385031215613eb757600080fd5b613ec083613bcf565b915060208301356001600160601b0381168114613edc57600080fd5b809150509250929050565b600060208284031215613ef957600080fd5b5035919050565b60008060608385031215613f1357600080fd5b613f1c83613bcf565b9150613c108460208501613c19565b6001600160601b03851681526001600160401b0384166020808301919091526001600160a01b038481166040840152608060608401819052845190840181905260009285810192909160a0860190855b81811015613f99578551841683529484019491840191600101613f7b565b50909a9950505050505050505050565b60008060008060608587031215613fbf57600080fd5b613fc885613bcf565b93506020850135925060408501356001600160401b0380821115613feb57600080fd5b818701915087601f830112613fff57600080fd5b81358181111561400e57600080fd5b88602082850101111561402057600080fd5b95989497505060200194505050565b600082601f83011261404057600080fd5b604051604081018181106001600160401b038211171561406257614062613cd5565b806040525080604084018581111561407957600080fd5b845b8181101561409357803583526020928301920161407b565b509195945050505050565b600060a082840312156140b057600080fd5b60405160a081018181106001600160401b03821117156140d2576140d2613cd5565b6040529050806140e183613b9d565b81526140ef60208401613b9d565b602082015261410060408401613cc1565b604082015261411160608401613cc1565b606082015261412260808401613bcf565b60808201525092915050565b60008082840361024081121561414357600080fd5b6101a08082121561415357600080fd5b61415b613ceb565b9150614167868661402f565b8252614176866040870161402f565b60208301526080850135604083015260a0850135606083015260c085013560808301526141a560e08601613bcf565b60a08301526101006141b98782880161402f565b60c08401526141cc87610140880161402f565b60e084015261018086013581840152508193506141eb8682870161409e565b925050509250929050565b60006040828403121561420857600080fd5b612d2b838361402f565b60006020828403121561422457600080fd5b612d2b82613bcf565b6001600160a01b0392831681529116602082015260400190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000828210156142855761428561425d565b500390565b634e487b7160e01b600052603160045260246000fd5b6000600182016142b2576142b261425d565b5060010190565b61ffff8716815263ffffffff86811660208301528581166040830152848116606083015260808201849052825480821660a08401526101c08301919061430c60c08501838360201c1663ffffffff169052565b61432360e08501838360401c1663ffffffff169052565b61433b6101008501838360601c1663ffffffff169052565b6143536101208501838360801c1663ffffffff169052565b62ffffff60a082901c811661014086015260b882901c811661016086015260d082901c1661018085015260e81c6101a090930192909252979650505050505050565b6001600160401b039290921682526001600160a01b0316602082015260400190565b60006001600160401b038083168185168083038211156143d9576143d961425d565b01949350505050565b60006001600160601b03838116908316818110156144025761440261425d565b039392505050565b6001600160a01b039290921682526001600160601b0316602082015260400190565b60006020828403121561443e57600080fd5b81518015158114612d2b57600080fd5b60006001600160401b0380831681810361446a5761446a61425d565b6001019392505050565b60006001600160601b038281168482168083038211156143d9576143d961425d565b600082198211156144a9576144a961425d565b500190565b6000604082018483526020604081850152818551808452606086019150828701935060005b818110156144ef578451835293830193918301916001016144d3565b5090979650505050505050565b8060005b6002811015610839578151845260209384019390910190600101614500565b6040810161205d82846144fc565b60006020828403121561453f57600080fd5b5051919050565b6001600160a01b03929092168252602082015260400190565b60008160001904831182151516156145795761457961425d565b500290565b634e487b7160e01b600052601260045260246000fd5b6000826145a3576145a361457e565b500490565b82815260608101612d2b60208301846144fc565b805169ffffffffffffffffffff8116811461361757600080fd5b600080600080600060a086880312156145ee57600080fd5b6145f7866145bc565b9450602086015193506040860151925060608601519150613e98608087016145bc565b83815261462a60208201846144fc565b606081019190915260800192915050565b86815261464b60208201876144fc565b61465860608201866144fc565b61466560a08201856144fc565b61467260e08201846144fc565b60609190911b6001600160601b0319166101208201526101340195945050505050565b6000826146a4576146a461457e565b500690565b6146b381836144fc565b60400191905056fea26469706673582212207dacc4df1823a6207ad382421a4c8874c014f085a040bc6f2f052eb5eda86f9764736f6c634300080d0033",
}

// VRFCoordinatorV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFCoordinatorV2MetaData.ABI instead.
var VRFCoordinatorV2ABI = VRFCoordinatorV2MetaData.ABI

// VRFCoordinatorV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VRFCoordinatorV2MetaData.Bin instead.
var VRFCoordinatorV2Bin = VRFCoordinatorV2MetaData.Bin

// DeployVRFCoordinatorV2 deploys a new Ethereum contract, binding an instance of VRFCoordinatorV2 to it.
func DeployVRFCoordinatorV2(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, blockhashStore common.Address, linkEthFeed common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV2, error) {
	parsed, err := VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2Bin), backend, link, blockhashStore, linkEthFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2{VRFCoordinatorV2Caller: VRFCoordinatorV2Caller{contract: contract}, VRFCoordinatorV2Transactor: VRFCoordinatorV2Transactor{contract: contract}, VRFCoordinatorV2Filterer: VRFCoordinatorV2Filterer{contract: contract}}, nil
}

// VRFCoordinatorV2 is an auto generated Go binding around an Ethereum contract.
type VRFCoordinatorV2 struct {
	VRFCoordinatorV2Caller     // Read-only binding to the contract
	VRFCoordinatorV2Transactor // Write-only binding to the contract
	VRFCoordinatorV2Filterer   // Log filterer for contract events
}

// VRFCoordinatorV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type VRFCoordinatorV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFCoordinatorV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFCoordinatorV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFCoordinatorV2Session struct {
	Contract     *VRFCoordinatorV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFCoordinatorV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFCoordinatorV2CallerSession struct {
	Contract *VRFCoordinatorV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// VRFCoordinatorV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFCoordinatorV2TransactorSession struct {
	Contract     *VRFCoordinatorV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// VRFCoordinatorV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type VRFCoordinatorV2Raw struct {
	Contract *VRFCoordinatorV2 // Generic contract binding to access the raw methods on
}

// VRFCoordinatorV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFCoordinatorV2CallerRaw struct {
	Contract *VRFCoordinatorV2Caller // Generic read-only contract binding to access the raw methods on
}

// VRFCoordinatorV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFCoordinatorV2TransactorRaw struct {
	Contract *VRFCoordinatorV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFCoordinatorV2 creates a new instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2, error) {
	contract, err := bindVRFCoordinatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2{VRFCoordinatorV2Caller: VRFCoordinatorV2Caller{contract: contract}, VRFCoordinatorV2Transactor: VRFCoordinatorV2Transactor{contract: contract}, VRFCoordinatorV2Filterer: VRFCoordinatorV2Filterer{contract: contract}}, nil
}

// NewVRFCoordinatorV2Caller creates a new read-only instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2Caller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2Caller, error) {
	contract, err := bindVRFCoordinatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Caller{contract: contract}, nil
}

// NewVRFCoordinatorV2Transactor creates a new write-only instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2Transactor, error) {
	contract, err := bindVRFCoordinatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Transactor{contract: contract}, nil
}

// NewVRFCoordinatorV2Filterer creates a new log filterer instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2Filterer, error) {
	contract, err := bindVRFCoordinatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Filterer{contract: contract}, nil
}

// bindVRFCoordinatorV2 binds a generic wrapper to an already deployed contract.
func bindVRFCoordinatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.contract.Transact(opts, method, params...)
}

// BLOCKHASHSTORE is a free data retrieval call binding the contract method 0x689c4517.
//
// Solidity: function BLOCKHASH_STORE() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BLOCKHASHSTORE is a free data retrieval call binding the contract method 0x689c4517.
//
// Solidity: function BLOCKHASH_STORE() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2.CallOpts)
}

// BLOCKHASHSTORE is a free data retrieval call binding the contract method 0x689c4517.
//
// Solidity: function BLOCKHASH_STORE() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) LINK() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINK(&_VRFCoordinatorV2.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINK(&_VRFCoordinatorV2.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINKETHFEED(&_VRFCoordinatorV2.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINKETHFEED(&_VRFCoordinatorV2.CallOpts)
}

// MAXCONSUMERS is a free data retrieval call binding the contract method 0x64d51a2a.
//
// Solidity: function MAX_CONSUMERS() view returns(uint16)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXCONSUMERS is a free data retrieval call binding the contract method 0x64d51a2a.
//
// Solidity: function MAX_CONSUMERS() view returns(uint16)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXCONSUMERS(&_VRFCoordinatorV2.CallOpts)
}

// MAXCONSUMERS is a free data retrieval call binding the contract method 0x64d51a2a.
//
// Solidity: function MAX_CONSUMERS() view returns(uint16)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXCONSUMERS(&_VRFCoordinatorV2.CallOpts)
}

// MAXNUMWORDS is a free data retrieval call binding the contract method 0x40d6bb82.
//
// Solidity: function MAX_NUM_WORDS() view returns(uint32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MAXNUMWORDS is a free data retrieval call binding the contract method 0x40d6bb82.
//
// Solidity: function MAX_NUM_WORDS() view returns(uint32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2.Contract.MAXNUMWORDS(&_VRFCoordinatorV2.CallOpts)
}

// MAXNUMWORDS is a free data retrieval call binding the contract method 0x40d6bb82.
//
// Solidity: function MAX_NUM_WORDS() view returns(uint32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2.Contract.MAXNUMWORDS(&_VRFCoordinatorV2.CallOpts)
}

// MAXREQUESTCONFIRMATIONS is a free data retrieval call binding the contract method 0x15c48b84.
//
// Solidity: function MAX_REQUEST_CONFIRMATIONS() view returns(uint16)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXREQUESTCONFIRMATIONS is a free data retrieval call binding the contract method 0x15c48b84.
//
// Solidity: function MAX_REQUEST_CONFIRMATIONS() view returns(uint16)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2.CallOpts)
}

// MAXREQUESTCONFIRMATIONS is a free data retrieval call binding the contract method 0x15c48b84.
//
// Solidity: function MAX_REQUEST_CONFIRMATIONS() view returns(uint16)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2.CallOpts)
}

// GetCommitment is a free data retrieval call binding the contract method 0x69bcdb7d.
//
// Solidity: function getCommitment(uint256 requestId) view returns(bytes32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getCommitment", requestId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x69bcdb7d.
//
// Solidity: function getCommitment(uint256 requestId) view returns(bytes32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetCommitment(&_VRFCoordinatorV2.CallOpts, requestId)
}

// GetCommitment is a free data retrieval call binding the contract method 0x69bcdb7d.
//
// Solidity: function getCommitment(uint256 requestId) view returns(bytes32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetCommitment(&_VRFCoordinatorV2.CallOpts, requestId)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetConfig(opts *bind.CallOpts) (struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getConfig")

	outstruct := new(struct {
		MinimumRequestConfirmations uint16
		MaxGasLimit                 uint32
		StalenessSeconds            uint32
		GasAfterPaymentCalculation  uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimumRequestConfirmations = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.MaxGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.StalenessSeconds = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetConfig() (struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}, error) {
	return _VRFCoordinatorV2.Contract.GetConfig(&_VRFCoordinatorV2.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetConfig() (struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}, error) {
	return _VRFCoordinatorV2.Contract.GetConfig(&_VRFCoordinatorV2.CallOpts)
}

// GetCurrentSubId is a free data retrieval call binding the contract method 0x06bfa637.
//
// Solidity: function getCurrentSubId() view returns(uint64)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getCurrentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetCurrentSubId is a free data retrieval call binding the contract method 0x06bfa637.
//
// Solidity: function getCurrentSubId() view returns(uint64)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetCurrentSubId() (uint64, error) {
	return _VRFCoordinatorV2.Contract.GetCurrentSubId(&_VRFCoordinatorV2.CallOpts)
}

// GetCurrentSubId is a free data retrieval call binding the contract method 0x06bfa637.
//
// Solidity: function getCurrentSubId() view returns(uint64)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetCurrentSubId() (uint64, error) {
	return _VRFCoordinatorV2.Contract.GetCurrentSubId(&_VRFCoordinatorV2.CallOpts)
}

// GetFallbackWeiPerUnitLink is a free data retrieval call binding the contract method 0x356dac71.
//
// Solidity: function getFallbackWeiPerUnitLink() view returns(int256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getFallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFallbackWeiPerUnitLink is a free data retrieval call binding the contract method 0x356dac71.
//
// Solidity: function getFallbackWeiPerUnitLink() view returns(int256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetFallbackWeiPerUnitLink(&_VRFCoordinatorV2.CallOpts)
}

// GetFallbackWeiPerUnitLink is a free data retrieval call binding the contract method 0x356dac71.
//
// Solidity: function getFallbackWeiPerUnitLink() view returns(int256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetFallbackWeiPerUnitLink(&_VRFCoordinatorV2.CallOpts)
}

// GetFeeConfig is a free data retrieval call binding the contract method 0x5fbbc0d2.
//
// Solidity: function getFeeConfig() view returns(uint32 fulfillmentFlatFeeLinkPPMTier1, uint32 fulfillmentFlatFeeLinkPPMTier2, uint32 fulfillmentFlatFeeLinkPPMTier3, uint32 fulfillmentFlatFeeLinkPPMTier4, uint32 fulfillmentFlatFeeLinkPPMTier5, uint24 reqsForTier2, uint24 reqsForTier3, uint24 reqsForTier4, uint24 reqsForTier5)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetFeeConfig(opts *bind.CallOpts) (struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getFeeConfig")

	outstruct := new(struct {
		FulfillmentFlatFeeLinkPPMTier1 uint32
		FulfillmentFlatFeeLinkPPMTier2 uint32
		FulfillmentFlatFeeLinkPPMTier3 uint32
		FulfillmentFlatFeeLinkPPMTier4 uint32
		FulfillmentFlatFeeLinkPPMTier5 uint32
		ReqsForTier2                   *big.Int
		ReqsForTier3                   *big.Int
		ReqsForTier4                   *big.Int
		ReqsForTier5                   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FulfillmentFlatFeeLinkPPMTier1 = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier2 = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier3 = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier4 = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier5 = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.ReqsForTier2 = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier3 = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier4 = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier5 = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetFeeConfig is a free data retrieval call binding the contract method 0x5fbbc0d2.
//
// Solidity: function getFeeConfig() view returns(uint32 fulfillmentFlatFeeLinkPPMTier1, uint32 fulfillmentFlatFeeLinkPPMTier2, uint32 fulfillmentFlatFeeLinkPPMTier3, uint32 fulfillmentFlatFeeLinkPPMTier4, uint32 fulfillmentFlatFeeLinkPPMTier5, uint24 reqsForTier2, uint24 reqsForTier3, uint24 reqsForTier4, uint24 reqsForTier5)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetFeeConfig() (struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}, error) {
	return _VRFCoordinatorV2.Contract.GetFeeConfig(&_VRFCoordinatorV2.CallOpts)
}

// GetFeeConfig is a free data retrieval call binding the contract method 0x5fbbc0d2.
//
// Solidity: function getFeeConfig() view returns(uint32 fulfillmentFlatFeeLinkPPMTier1, uint32 fulfillmentFlatFeeLinkPPMTier2, uint32 fulfillmentFlatFeeLinkPPMTier3, uint32 fulfillmentFlatFeeLinkPPMTier4, uint32 fulfillmentFlatFeeLinkPPMTier5, uint24 reqsForTier2, uint24 reqsForTier3, uint24 reqsForTier4, uint24 reqsForTier5)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetFeeConfig() (struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}, error) {
	return _VRFCoordinatorV2.Contract.GetFeeConfig(&_VRFCoordinatorV2.CallOpts)
}

// GetFeeTier is a free data retrieval call binding the contract method 0xd2f9f9a7.
//
// Solidity: function getFeeTier(uint64 reqCount) view returns(uint32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getFeeTier", reqCount)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetFeeTier is a free data retrieval call binding the contract method 0xd2f9f9a7.
//
// Solidity: function getFeeTier(uint64 reqCount) view returns(uint32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorV2.Contract.GetFeeTier(&_VRFCoordinatorV2.CallOpts, reqCount)
}

// GetFeeTier is a free data retrieval call binding the contract method 0xd2f9f9a7.
//
// Solidity: function getFeeTier(uint64 reqCount) view returns(uint32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorV2.Contract.GetFeeTier(&_VRFCoordinatorV2.CallOpts, reqCount)
}

// GetRequestConfig is a free data retrieval call binding the contract method 0x00012291.
//
// Solidity: function getRequestConfig() view returns(uint16, uint32, bytes32[])
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

// GetRequestConfig is a free data retrieval call binding the contract method 0x00012291.
//
// Solidity: function getRequestConfig() view returns(uint16, uint32, bytes32[])
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetRequestConfig(&_VRFCoordinatorV2.CallOpts)
}

// GetRequestConfig is a free data retrieval call binding the contract method 0x00012291.
//
// Solidity: function getRequestConfig() view returns(uint16, uint32, bytes32[])
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetRequestConfig(&_VRFCoordinatorV2.CallOpts)
}

// GetSubscription is a free data retrieval call binding the contract method 0xa47c7696.
//
// Solidity: function getSubscription(uint64 subId) view returns(uint96 balance, uint64 reqCount, address owner, address[] consumers)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetSubscription(opts *bind.CallOpts, subId uint64) (struct {
	Balance   *big.Int
	ReqCount  uint64
	Owner     common.Address
	Consumers []common.Address
}, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(struct {
		Balance   *big.Int
		ReqCount  uint64
		Owner     common.Address
		Consumers []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetSubscription is a free data retrieval call binding the contract method 0xa47c7696.
//
// Solidity: function getSubscription(uint64 subId) view returns(uint96 balance, uint64 reqCount, address owner, address[] consumers)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetSubscription(subId uint64) (struct {
	Balance   *big.Int
	ReqCount  uint64
	Owner     common.Address
	Consumers []common.Address
}, error) {
	return _VRFCoordinatorV2.Contract.GetSubscription(&_VRFCoordinatorV2.CallOpts, subId)
}

// GetSubscription is a free data retrieval call binding the contract method 0xa47c7696.
//
// Solidity: function getSubscription(uint64 subId) view returns(uint96 balance, uint64 reqCount, address owner, address[] consumers)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetSubscription(subId uint64) (struct {
	Balance   *big.Int
	ReqCount  uint64
	Owner     common.Address
	Consumers []common.Address
}, error) {
	return _VRFCoordinatorV2.Contract.GetSubscription(&_VRFCoordinatorV2.CallOpts, subId)
}

// GetTotalBalance is a free data retrieval call binding the contract method 0x12b58349.
//
// Solidity: function getTotalBalance() view returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalBalance is a free data retrieval call binding the contract method 0x12b58349.
//
// Solidity: function getTotalBalance() view returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetTotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetTotalBalance(&_VRFCoordinatorV2.CallOpts)
}

// GetTotalBalance is a free data retrieval call binding the contract method 0x12b58349.
//
// Solidity: function getTotalBalance() view returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetTotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetTotalBalance(&_VRFCoordinatorV2.CallOpts)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] publicKey) pure returns(bytes32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] publicKey) pure returns(bytes32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.HashOfKey(&_VRFCoordinatorV2.CallOpts, publicKey)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] publicKey) pure returns(bytes32)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.HashOfKey(&_VRFCoordinatorV2.CallOpts, publicKey)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) Owner() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.Owner(&_VRFCoordinatorV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.Owner(&_VRFCoordinatorV2.CallOpts)
}

// PendingRequestExists is a free data retrieval call binding the contract method 0xe82ad7d4.
//
// Solidity: function pendingRequestExists(uint64 subId) view returns(bool)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PendingRequestExists is a free data retrieval call binding the contract method 0xe82ad7d4.
//
// Solidity: function pendingRequestExists(uint64 subId) view returns(bool)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorV2.Contract.PendingRequestExists(&_VRFCoordinatorV2.CallOpts, subId)
}

// PendingRequestExists is a free data retrieval call binding the contract method 0xe82ad7d4.
//
// Solidity: function pendingRequestExists(uint64 subId) view returns(bool)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorV2.Contract.PendingRequestExists(&_VRFCoordinatorV2.CallOpts, subId)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) TypeAndVersion() (string, error) {
	return _VRFCoordinatorV2.Contract.TypeAndVersion(&_VRFCoordinatorV2.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) TypeAndVersion() (string, error) {
	return _VRFCoordinatorV2.Contract.TypeAndVersion(&_VRFCoordinatorV2.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptOwnership(&_VRFCoordinatorV2.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptOwnership(&_VRFCoordinatorV2.TransactOpts)
}

// AcceptSubscriptionOwnerTransfer is a paid mutator transaction binding the contract method 0x82359740.
//
// Solidity: function acceptSubscriptionOwnerTransfer(uint64 subId) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

// AcceptSubscriptionOwnerTransfer is a paid mutator transaction binding the contract method 0x82359740.
//
// Solidity: function acceptSubscriptionOwnerTransfer(uint64 subId) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId)
}

// AcceptSubscriptionOwnerTransfer is a paid mutator transaction binding the contract method 0x82359740.
//
// Solidity: function acceptSubscriptionOwnerTransfer(uint64 subId) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId)
}

// AddConsumer is a paid mutator transaction binding the contract method 0x7341c10c.
//
// Solidity: function addConsumer(uint64 subId, address consumer) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "addConsumer", subId, consumer)
}

// AddConsumer is a paid mutator transaction binding the contract method 0x7341c10c.
//
// Solidity: function addConsumer(uint64 subId, address consumer) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AddConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

// AddConsumer is a paid mutator transaction binding the contract method 0x7341c10c.
//
// Solidity: function addConsumer(uint64 subId, address consumer) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AddConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

// CancelSubscription is a paid mutator transaction binding the contract method 0xd7ae1d30.
//
// Solidity: function cancelSubscription(uint64 subId, address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "cancelSubscription", subId, to)
}

// CancelSubscription is a paid mutator transaction binding the contract method 0xd7ae1d30.
//
// Solidity: function cancelSubscription(uint64 subId, address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId, to)
}

// CancelSubscription is a paid mutator transaction binding the contract method 0xd7ae1d30.
//
// Solidity: function cancelSubscription(uint64 subId, address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId, to)
}

// CreateSubscription is a paid mutator transaction binding the contract method 0xa21a23e4.
//
// Solidity: function createSubscription() returns(uint64)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "createSubscription")
}

// CreateSubscription is a paid mutator transaction binding the contract method 0xa21a23e4.
//
// Solidity: function createSubscription() returns(uint64)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CreateSubscription(&_VRFCoordinatorV2.TransactOpts)
}

// CreateSubscription is a paid mutator transaction binding the contract method 0xa21a23e4.
//
// Solidity: function createSubscription() returns(uint64)
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CreateSubscription(&_VRFCoordinatorV2.TransactOpts)
}

// DeregisterProvingKey is a paid mutator transaction binding the contract method 0x08821d58.
//
// Solidity: function deregisterProvingKey(uint256[2] publicProvingKey) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

// DeregisterProvingKey is a paid mutator transaction binding the contract method 0x08821d58.
//
// Solidity: function deregisterProvingKey(uint256[2] publicProvingKey) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.DeregisterProvingKey(&_VRFCoordinatorV2.TransactOpts, publicProvingKey)
}

// DeregisterProvingKey is a paid mutator transaction binding the contract method 0x08821d58.
//
// Solidity: function deregisterProvingKey(uint256[2] publicProvingKey) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.DeregisterProvingKey(&_VRFCoordinatorV2.TransactOpts, publicProvingKey)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0xaf198b97.
//
// Solidity: function fulfillRandomWords((uint256[2],uint256[2],uint256,uint256,uint256,address,uint256[2],uint256[2],uint256) proof, (uint64,uint64,uint32,uint32,address) rc) returns(uint96)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0xaf198b97.
//
// Solidity: function fulfillRandomWords((uint256[2],uint256[2],uint256,uint256,uint256,address,uint256[2],uint256[2],uint256) proof, (uint64,uint64,uint32,uint32,address) rc) returns(uint96)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.FulfillRandomWords(&_VRFCoordinatorV2.TransactOpts, proof, rc)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0xaf198b97.
//
// Solidity: function fulfillRandomWords((uint256[2],uint256[2],uint256,uint256,uint256,address,uint256[2],uint256[2],uint256) proof, (uint64,uint64,uint32,uint32,address) rc) returns(uint96)
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.FulfillRandomWords(&_VRFCoordinatorV2.TransactOpts, proof, rc)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 amount, bytes data) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 amount, bytes data) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OnTokenTransfer(&_VRFCoordinatorV2.TransactOpts, arg0, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 amount, bytes data) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OnTokenTransfer(&_VRFCoordinatorV2.TransactOpts, arg0, amount, data)
}

// OracleWithdraw is a paid mutator transaction binding the contract method 0x66316d8d.
//
// Solidity: function oracleWithdraw(address recipient, uint96 amount) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

// OracleWithdraw is a paid mutator transaction binding the contract method 0x66316d8d.
//
// Solidity: function oracleWithdraw(address recipient, uint96 amount) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OracleWithdraw(&_VRFCoordinatorV2.TransactOpts, recipient, amount)
}

// OracleWithdraw is a paid mutator transaction binding the contract method 0x66316d8d.
//
// Solidity: function oracleWithdraw(address recipient, uint96 amount) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OracleWithdraw(&_VRFCoordinatorV2.TransactOpts, recipient, amount)
}

// OwnerCancelSubscription is a paid mutator transaction binding the contract method 0x02bcc5b6.
//
// Solidity: function ownerCancelSubscription(uint64 subId) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "ownerCancelSubscription", subId)
}

// OwnerCancelSubscription is a paid mutator transaction binding the contract method 0x02bcc5b6.
//
// Solidity: function ownerCancelSubscription(uint64 subId) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId)
}

// OwnerCancelSubscription is a paid mutator transaction binding the contract method 0x02bcc5b6.
//
// Solidity: function ownerCancelSubscription(uint64 subId) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xe72f6e30.
//
// Solidity: function recoverFunds(address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "recoverFunds", to)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xe72f6e30.
//
// Solidity: function recoverFunds(address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RecoverFunds(&_VRFCoordinatorV2.TransactOpts, to)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xe72f6e30.
//
// Solidity: function recoverFunds(address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RecoverFunds(&_VRFCoordinatorV2.TransactOpts, to)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0x6f64f03f.
//
// Solidity: function registerProvingKey(address oracle, uint256[2] publicProvingKey) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0x6f64f03f.
//
// Solidity: function registerProvingKey(address oracle, uint256[2] publicProvingKey) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RegisterProvingKey(&_VRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0x6f64f03f.
//
// Solidity: function registerProvingKey(address oracle, uint256[2] publicProvingKey) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RegisterProvingKey(&_VRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

// RemoveConsumer is a paid mutator transaction binding the contract method 0x9f87fad7.
//
// Solidity: function removeConsumer(uint64 subId, address consumer) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "removeConsumer", subId, consumer)
}

// RemoveConsumer is a paid mutator transaction binding the contract method 0x9f87fad7.
//
// Solidity: function removeConsumer(uint64 subId, address consumer) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RemoveConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

// RemoveConsumer is a paid mutator transaction binding the contract method 0x9f87fad7.
//
// Solidity: function removeConsumer(uint64 subId, address consumer) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RemoveConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x5d3b1d30.
//
// Solidity: function requestRandomWords(bytes32 keyHash, uint64 subId, uint16 requestConfirmations, uint32 callbackGasLimit, uint32 numWords) returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "requestRandomWords", keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x5d3b1d30.
//
// Solidity: function requestRandomWords(bytes32 keyHash, uint64 subId, uint16 requestConfirmations, uint32 callbackGasLimit, uint32 numWords) returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestRandomWords(&_VRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x5d3b1d30.
//
// Solidity: function requestRandomWords(bytes32 keyHash, uint64 subId, uint16 requestConfirmations, uint32 callbackGasLimit, uint32 numWords) returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestRandomWords(&_VRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

// RequestSubscriptionOwnerTransfer is a paid mutator transaction binding the contract method 0x04c357cb.
//
// Solidity: function requestSubscriptionOwnerTransfer(uint64 subId, address newOwner) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

// RequestSubscriptionOwnerTransfer is a paid mutator transaction binding the contract method 0x04c357cb.
//
// Solidity: function requestSubscriptionOwnerTransfer(uint64 subId, address newOwner) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId, newOwner)
}

// RequestSubscriptionOwnerTransfer is a paid mutator transaction binding the contract method 0x04c357cb.
//
// Solidity: function requestSubscriptionOwnerTransfer(uint64 subId, address newOwner) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId, newOwner)
}

// SetConfig is a paid mutator transaction binding the contract method 0x4cb48a54.
//
// Solidity: function setConfig(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation, int256 fallbackWeiPerUnitLink, (uint32,uint32,uint32,uint32,uint32,uint24,uint24,uint24,uint24) feeConfig) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0x4cb48a54.
//
// Solidity: function setConfig(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation, int256 fallbackWeiPerUnitLink, (uint32,uint32,uint32,uint32,uint32,uint24,uint24,uint24,uint24) feeConfig) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.SetConfig(&_VRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0x4cb48a54.
//
// Solidity: function setConfig(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation, int256 fallbackWeiPerUnitLink, (uint32,uint32,uint32,uint32,uint32,uint24,uint24,uint24,uint24) feeConfig) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.SetConfig(&_VRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.TransferOwnership(&_VRFCoordinatorV2.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.TransferOwnership(&_VRFCoordinatorV2.TransactOpts, to)
}

// VRFCoordinatorV2ConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2ConfigSetIterator struct {
	Event *VRFCoordinatorV2ConfigSet // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2ConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2ConfigSet)
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
		it.Event = new(VRFCoordinatorV2ConfigSet)
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
func (it *VRFCoordinatorV2ConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2ConfigSet represents a ConfigSet event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2ConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FeeConfig                   VRFCoordinatorV2FeeConfig
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2.
//
// Solidity: event ConfigSet(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation, int256 fallbackWeiPerUnitLink, (uint32,uint32,uint32,uint32,uint32,uint24,uint24,uint24,uint24) feeConfig)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2ConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2ConfigSetIterator{contract: _VRFCoordinatorV2.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2.
//
// Solidity: event ConfigSet(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation, int256 fallbackWeiPerUnitLink, (uint32,uint32,uint32,uint32,uint32,uint24,uint24,uint24,uint24) feeConfig)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2ConfigSet)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2.
//
// Solidity: event ConfigSet(uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 stalenessSeconds, uint32 gasAfterPaymentCalculation, int256 fallbackWeiPerUnitLink, (uint32,uint32,uint32,uint32,uint32,uint24,uint24,uint24,uint24) feeConfig)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV2ConfigSet, error) {
	event := new(VRFCoordinatorV2ConfigSet)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2FundsRecoveredIterator is returned from FilterFundsRecovered and is used to iterate over the raw logs and unpacked data for FundsRecovered events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2FundsRecoveredIterator struct {
	Event *VRFCoordinatorV2FundsRecovered // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2FundsRecoveredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2FundsRecovered)
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
		it.Event = new(VRFCoordinatorV2FundsRecovered)
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
func (it *VRFCoordinatorV2FundsRecoveredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2FundsRecovered represents a FundsRecovered event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsRecovered is a free log retrieval operation binding the contract event 0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600.
//
// Solidity: event FundsRecovered(address to, uint256 amount)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2FundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2FundsRecoveredIterator{contract: _VRFCoordinatorV2.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

// WatchFundsRecovered is a free log subscription operation binding the contract event 0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600.
//
// Solidity: event FundsRecovered(address to, uint256 amount)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2FundsRecovered)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

// ParseFundsRecovered is a log parse operation binding the contract event 0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600.
//
// Solidity: event FundsRecovered(address to, uint256 amount)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2FundsRecovered, error) {
	event := new(VRFCoordinatorV2FundsRecovered)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2OwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2OwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV2OwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2OwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2OwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV2OwnershipTransferRequested)
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
func (it *VRFCoordinatorV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2OwnershipTransferRequested represents a OwnershipTransferRequested event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2OwnershipTransferRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2OwnershipTransferRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2OwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV2OwnershipTransferRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2OwnershipTransferredIterator struct {
	Event *VRFCoordinatorV2OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2OwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV2OwnershipTransferred)
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
func (it *VRFCoordinatorV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2OwnershipTransferred represents a OwnershipTransferred event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2OwnershipTransferredIterator{contract: _VRFCoordinatorV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2OwnershipTransferred)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2OwnershipTransferred, error) {
	event := new(VRFCoordinatorV2OwnershipTransferred)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2ProvingKeyDeregisteredIterator is returned from FilterProvingKeyDeregistered and is used to iterate over the raw logs and unpacked data for ProvingKeyDeregistered events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2ProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV2ProvingKeyDeregistered // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2ProvingKeyDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2ProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorV2ProvingKeyDeregistered)
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
func (it *VRFCoordinatorV2ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2ProvingKeyDeregistered represents a ProvingKeyDeregistered event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2ProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterProvingKeyDeregistered is a free log retrieval operation binding the contract event 0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d.
//
// Solidity: event ProvingKeyDeregistered(bytes32 keyHash, address indexed oracle)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2ProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2ProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV2.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

// WatchProvingKeyDeregistered is a free log subscription operation binding the contract event 0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d.
//
// Solidity: event ProvingKeyDeregistered(bytes32 keyHash, address indexed oracle)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2ProvingKeyDeregistered)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

// ParseProvingKeyDeregistered is a log parse operation binding the contract event 0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d.
//
// Solidity: event ProvingKeyDeregistered(bytes32 keyHash, address indexed oracle)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2ProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV2ProvingKeyDeregistered)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2ProvingKeyRegisteredIterator is returned from FilterProvingKeyRegistered and is used to iterate over the raw logs and unpacked data for ProvingKeyRegistered events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2ProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV2ProvingKeyRegistered // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2ProvingKeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2ProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV2ProvingKeyRegistered)
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
func (it *VRFCoordinatorV2ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2ProvingKeyRegistered represents a ProvingKeyRegistered event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2ProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterProvingKeyRegistered is a free log retrieval operation binding the contract event 0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8.
//
// Solidity: event ProvingKeyRegistered(bytes32 keyHash, address indexed oracle)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2ProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2ProvingKeyRegisteredIterator{contract: _VRFCoordinatorV2.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

// WatchProvingKeyRegistered is a free log subscription operation binding the contract event 0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8.
//
// Solidity: event ProvingKeyRegistered(bytes32 keyHash, address indexed oracle)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2ProvingKeyRegistered)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

// ParseProvingKeyRegistered is a log parse operation binding the contract event 0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8.
//
// Solidity: event ProvingKeyRegistered(bytes32 keyHash, address indexed oracle)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2ProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV2ProvingKeyRegistered)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2RandomWordsFulfilledIterator is returned from FilterRandomWordsFulfilled and is used to iterate over the raw logs and unpacked data for RandomWordsFulfilled events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2RandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV2RandomWordsFulfilled // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2RandomWordsFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2RandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV2RandomWordsFulfilled)
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
func (it *VRFCoordinatorV2RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2RandomWordsFulfilled represents a RandomWordsFulfilled event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2RandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRandomWordsFulfilled is a free log retrieval operation binding the contract event 0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4.
//
// Solidity: event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bool success)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV2RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2RandomWordsFulfilledIterator{contract: _VRFCoordinatorV2.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

// WatchRandomWordsFulfilled is a free log subscription operation binding the contract event 0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4.
//
// Solidity: event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bool success)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2RandomWordsFulfilled)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

// ParseRandomWordsFulfilled is a log parse operation binding the contract event 0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4.
//
// Solidity: event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bool success)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2RandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV2RandomWordsFulfilled)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2RandomWordsRequestedIterator is returned from FilterRandomWordsRequested and is used to iterate over the raw logs and unpacked data for RandomWordsRequested events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2RandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV2RandomWordsRequested // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2RandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV2RandomWordsRequested)
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
func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2RandomWordsRequested represents a RandomWordsRequested event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2RandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       uint64
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	Sender                      common.Address
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterRandomWordsRequested is a free log retrieval operation binding the contract event 0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772.
//
// Solidity: event RandomWordsRequested(bytes32 indexed keyHash, uint256 requestId, uint256 preSeed, uint64 indexed subId, uint16 minimumRequestConfirmations, uint32 callbackGasLimit, uint32 numWords, address indexed sender)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV2RandomWordsRequestedIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2RandomWordsRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

// WatchRandomWordsRequested is a free log subscription operation binding the contract event 0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772.
//
// Solidity: event RandomWordsRequested(bytes32 indexed keyHash, uint256 requestId, uint256 preSeed, uint64 indexed subId, uint16 minimumRequestConfirmations, uint32 callbackGasLimit, uint32 numWords, address indexed sender)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2RandomWordsRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

// ParseRandomWordsRequested is a log parse operation binding the contract event 0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772.
//
// Solidity: event RandomWordsRequested(bytes32 indexed keyHash, uint256 requestId, uint256 preSeed, uint64 indexed subId, uint16 minimumRequestConfirmations, uint32 callbackGasLimit, uint32 numWords, address indexed sender)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2RandomWordsRequested, error) {
	event := new(VRFCoordinatorV2RandomWordsRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionCanceledIterator is returned from FilterSubscriptionCanceled and is used to iterate over the raw logs and unpacked data for SubscriptionCanceled events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV2SubscriptionCanceled // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV2SubscriptionCanceled)
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
func (it *VRFCoordinatorV2SubscriptionCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionCanceled represents a SubscriptionCanceled event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionCanceled struct {
	SubId  uint64
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionCanceled is a free log retrieval operation binding the contract event 0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815.
//
// Solidity: event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionCanceledIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

// WatchSubscriptionCanceled is a free log subscription operation binding the contract event 0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815.
//
// Solidity: event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionCanceled)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

// ParseSubscriptionCanceled is a log parse operation binding the contract event 0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815.
//
// Solidity: event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2SubscriptionCanceled, error) {
	event := new(VRFCoordinatorV2SubscriptionCanceled)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionConsumerAddedIterator is returned from FilterSubscriptionConsumerAdded and is used to iterate over the raw logs and unpacked data for SubscriptionConsumerAdded events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV2SubscriptionConsumerAdded // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionConsumerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV2SubscriptionConsumerAdded)
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
func (it *VRFCoordinatorV2SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionConsumerAdded represents a SubscriptionConsumerAdded event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionConsumerAdded is a free log retrieval operation binding the contract event 0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0.
//
// Solidity: event SubscriptionConsumerAdded(uint64 indexed subId, address consumer)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

// WatchSubscriptionConsumerAdded is a free log subscription operation binding the contract event 0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0.
//
// Solidity: event SubscriptionConsumerAdded(uint64 indexed subId, address consumer)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionConsumerAdded)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

// ParseSubscriptionConsumerAdded is a log parse operation binding the contract event 0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0.
//
// Solidity: event SubscriptionConsumerAdded(uint64 indexed subId, address consumer)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2SubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV2SubscriptionConsumerAdded)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionConsumerRemovedIterator is returned from FilterSubscriptionConsumerRemoved and is used to iterate over the raw logs and unpacked data for SubscriptionConsumerRemoved events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV2SubscriptionConsumerRemoved // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionConsumerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV2SubscriptionConsumerRemoved)
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
func (it *VRFCoordinatorV2SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionConsumerRemoved represents a SubscriptionConsumerRemoved event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionConsumerRemoved is a free log retrieval operation binding the contract event 0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b.
//
// Solidity: event SubscriptionConsumerRemoved(uint64 indexed subId, address consumer)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

// WatchSubscriptionConsumerRemoved is a free log subscription operation binding the contract event 0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b.
//
// Solidity: event SubscriptionConsumerRemoved(uint64 indexed subId, address consumer)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

// ParseSubscriptionConsumerRemoved is a log parse operation binding the contract event 0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b.
//
// Solidity: event SubscriptionConsumerRemoved(uint64 indexed subId, address consumer)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2SubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV2SubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionCreatedIterator is returned from FilterSubscriptionCreated and is used to iterate over the raw logs and unpacked data for SubscriptionCreated events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV2SubscriptionCreated // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV2SubscriptionCreated)
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
func (it *VRFCoordinatorV2SubscriptionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionCreated represents a SubscriptionCreated event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionCreated is a free log retrieval operation binding the contract event 0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf.
//
// Solidity: event SubscriptionCreated(uint64 indexed subId, address owner)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionCreatedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

// WatchSubscriptionCreated is a free log subscription operation binding the contract event 0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf.
//
// Solidity: event SubscriptionCreated(uint64 indexed subId, address owner)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionCreated)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

// ParseSubscriptionCreated is a log parse operation binding the contract event 0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf.
//
// Solidity: event SubscriptionCreated(uint64 indexed subId, address owner)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2SubscriptionCreated, error) {
	event := new(VRFCoordinatorV2SubscriptionCreated)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionFundedIterator is returned from FilterSubscriptionFunded and is used to iterate over the raw logs and unpacked data for SubscriptionFunded events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionFundedIterator struct {
	Event *VRFCoordinatorV2SubscriptionFunded // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV2SubscriptionFunded)
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
func (it *VRFCoordinatorV2SubscriptionFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionFunded represents a SubscriptionFunded event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionFunded is a free log retrieval operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionFundedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

// WatchSubscriptionFunded is a free log subscription operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionFunded)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

// ParseSubscriptionFunded is a log parse operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2SubscriptionFunded, error) {
	event := new(VRFCoordinatorV2SubscriptionFunded)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator is returned from FilterSubscriptionOwnerTransferRequested and is used to iterate over the raw logs and unpacked data for SubscriptionOwnerTransferRequested events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV2SubscriptionOwnerTransferRequested // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
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
func (it *VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionOwnerTransferRequested represents a SubscriptionOwnerTransferRequested event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionOwnerTransferRequested is a free log retrieval operation binding the contract event 0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be.
//
// Solidity: event SubscriptionOwnerTransferRequested(uint64 indexed subId, address from, address to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

// WatchSubscriptionOwnerTransferRequested is a free log subscription operation binding the contract event 0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be.
//
// Solidity: event SubscriptionOwnerTransferRequested(uint64 indexed subId, address from, address to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

// ParseSubscriptionOwnerTransferRequested is a log parse operation binding the contract event 0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be.
//
// Solidity: event SubscriptionOwnerTransferRequested(uint64 indexed subId, address from, address to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2SubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorV2SubscriptionOwnerTransferredIterator is returned from FilterSubscriptionOwnerTransferred and is used to iterate over the raw logs and unpacked data for SubscriptionOwnerTransferred events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV2SubscriptionOwnerTransferred // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorV2SubscriptionOwnerTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferred)
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
func (it *VRFCoordinatorV2SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2SubscriptionOwnerTransferred represents a SubscriptionOwnerTransferred event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2SubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionOwnerTransferred is a free log retrieval operation binding the contract event 0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0.
//
// Solidity: event SubscriptionOwnerTransferred(uint64 indexed subId, address from, address to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

// WatchSubscriptionOwnerTransferred is a free log subscription operation binding the contract event 0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0.
//
// Solidity: event SubscriptionOwnerTransferred(uint64 indexed subId, address from, address to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2SubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

// ParseSubscriptionOwnerTransferred is a log parse operation binding the contract event 0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0.
//
// Solidity: event SubscriptionOwnerTransferred(uint64 indexed subId, address from, address to)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2SubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV2SubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
