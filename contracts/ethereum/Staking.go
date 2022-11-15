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

// StakingPoolConstructorParams is an auto generated low-level Go binding around an user-defined struct.
type StakingPoolConstructorParams struct {
	LINKAddress                    common.Address
	MonitoredFeed                  common.Address
	InitialMaxPoolSize             *big.Int
	InitialMaxCommunityStakeAmount *big.Int
	InitialMaxOperatorStakeAmount  *big.Int
	MinCommunityStakeAmount        *big.Int
	MinOperatorStakeAmount         *big.Int
	PriorityPeriodThreshold        *big.Int
	RegularPeriodThreshold         *big.Int
	MaxAlertingRewardAmount        *big.Int
	MinInitialOperatorCount        *big.Int
	MinRewardDuration              *big.Int
	SlashableDuration              *big.Int
	DelegationRateDenominator      *big.Int
}

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"monitoredFeed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"initialMaxPoolSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialMaxCommunityStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialMaxOperatorStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minCommunityStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minOperatorStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priorityPeriodThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"regularPeriodThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxAlertingRewardAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minInitialOperatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minRewardDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashableDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRateDenominator\",\"type\":\"uint256\"}],\"internalType\":\"structStaking.PoolConstructorParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"AlertAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlertInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CastError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"remainingAmount\",\"type\":\"uint256\"}],\"name\":\"ExcessiveStakeAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"ExistingStakeFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentOperatorsCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minInitialOperatorsCount\",\"type\":\"uint256\"}],\"name\":\"InadequateInitialOperatorsCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"remainingPoolSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredPoolSize\",\"type\":\"uint256\"}],\"name\":\"InsufficientRemainingPoolSpace\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requiredAmount\",\"type\":\"uint256\"}],\"name\":\"InsufficientStakeAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDelegationRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxStakeAmount\",\"type\":\"uint256\"}],\"name\":\"InvalidMaxStakeAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMigrationTarget\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxPoolSize\",\"type\":\"uint256\"}],\"name\":\"InvalidPoolSize\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"currentStatus\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"requiredStatus\",\"type\":\"bool\"}],\"name\":\"InvalidPoolStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRegularPeriodThreshold\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MerkleRootNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorIsAssignedToFeed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorIsLocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RewardDurationTooShort\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderNotLinkToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"StakeNotFound\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardAmount\",\"type\":\"uint256\"}],\"name\":\"AlertRaised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newMerkleRoot\",\"type\":\"bytes32\"}],\"name\":\"MerkleRootChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"delegationReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Migrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"MigrationTargetAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"MigrationTargetProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newStake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"delegationReward\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptMigrationTarget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"addOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"canAlert\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newRate\",\"type\":\"uint256\"}],\"name\":\"changeRewardRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emergencyPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emergencyUnpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAvailableReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"getBaseReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommunityStakerLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDelegatesCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDelegationRateDenominator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"getDelegationReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEarnedBaseRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEarnedDelegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeedOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxPoolSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMigrationTarget\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMonitoredFeed\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperatorLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardTimestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"getStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalDelegatedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalRemovedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalStakedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"migrationTarget\",\"type\":\"address\"}],\"name\":\"proposeMigrationTarget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"raiseAlert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"removeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"setFeedOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newMerkleRoot\",\"type\":\"bytes32\"}],\"name\":\"setMerkleRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxPoolSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCommunityStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxOperatorStakeAmount\",\"type\":\"uint256\"}],\"name\":\"setPoolConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialRewardRate\",\"type\":\"uint256\"}],\"name\":\"start\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawRemovedStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawUnusedReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x34620006275762006518388190036101e0601f8201601f19168101906001600160401b0382119082101762000565576101c09282916040526101e0391262000627576040516101c081016001600160401b0381118282101762000565576040526101e0516001600160a01b038116810362000627578152610200516001600160a01b0381168103620006275760208201526102205160408201526102405160608201526102605160808201526102805160a08201526102a05160c08201526102c05160e08201526102e05161010082015261030051610120820152610320516101408201526103405161016082015261036051610180820152610380516101a08201523315620005e257600080546001600160a01b031916331790556001805460ff60a01b1916905580516001600160a01b031615620005d05760208101516001600160a01b031615620005d0576101a081015115620005be5761010081015160e08201511015620005ac5780516001600160a01b0316608090815260408201516060830151918301519092918382116200057b57600554936001600160601b03851681106200059457606085901c6001600160501b031682106200057b57828560b01c1162000536576040516001600160401b0360808201908111908211176200056557608081016040526004549060ff82161515815260ff8260081c169081602082015260606040820193600180831b038160101c168552600180831b039060701c1691015280600019048511811515166200054f5790516001600160601b03169084029081019081106200054f57811062000536576001600160601b038516819003620004c4575b506005546101a09450606081901c6001600160501b031682900362000462575b5050600554818160b01c0362000409575b505060018060a01b0360208201511660a05260e081015160c05261010081015160e0526101208101516101005260c08101516101205260a0810151610140526101408101516101605261016081015161018052610180810151825201516101c052604051615ed690816200064282396080518181816104a70152818161075401528181611069015281816116e401528181611c3101528181612c7c01528181614e4a01526150fb015260a0518181816114cc01528181612b0e0152614c98015260c051818181612b5d0152614ce7015260e051818181612b8e0152614d18015261010051818181615b850152615bb001526101205181818161153301528181612ce1015261569a01526101405181818161040d015261534501526101605181611bca01526101805181818161110f0152818161123d0152611c9a01526101a05181612cc001526101c051818181610c5d01528181614f65015281816150b9015281816153dc0152615ab50152f35b7f816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52916020916001600160b01b031962000442836200062c565b60b01b169060018060b01b031617600555604051908152a138806200029b565b7fb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df791602091600160601b600160b01b036200049d836200062c565b60601b16600160601b600160b01b03199190911617600555604051908152a138806200028a565b6001600160601b03811162000524576001600160601b03199094166001600160601b038516176005556040519384526101a0937f7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c090602090a1386200026a565b60405163408ba96f60e11b8152600490fd5b60405163bc91aa3360e01b815260048101849052602490fd5b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b60405163bc91aa3360e01b815260048101839052602490fd5b60249060405190630f9e1c3b60e11b82526004820152fd5b6040516310919fb960e11b8152600490fd5b60405163027953ef60e61b8152600490fd5b60405163f6b2911f60e01b8152600490fd5b60405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f00000000000000006044820152606490fd5b600080fd5b6001600160501b03908181116200052457169056fe60806040526004361015610013575b600080fd5b60003560e01c80630641bdd8146103f25780630fbc8f5b146103e9578063165d35e1146103e0578063181f5a77146103d75780631a9d4c7c146103ce5780631ddb5552146103c557806322f3e2d4146103bc5780632def6620146103b357806332e28850146103aa57806338adb6f0146103a157806349590657146103985780634a4e3bd51461038f57806351858e271461038657806359f018791461037d5780635aa6e013146103745780635c975abb146102935780635e8b40d71461036b5780635fec60f81461036257806363b2c85a146103595780636d70f7ae14610350578063741040021461034757806374de4ec41461033e57806374f237c41461033557806379ba50971461032c5780637a766460146103235780637cb647591461031a5780637e1a3786146103115780638019e7d01461030857806383db28a0146102ff57806387e900b1146102f65780638856398f146102ed5780638932a90d146102e45780638a44f337146102db5780638da5cb5b146102d25780638fb4b573146102c95780639a109bc2146102c05780639d0a3864146102b7578063a07aea1c146102ae578063a4c0ed36146102a5578063a7a2f5aa1461029c578063b187bd2614610293578063bfbd9b1b1461028a578063c1852f5814610281578063d365a37714610278578063da9c732f1461026f578063e0974ea514610266578063e5f929731461025d578063e937fdaa14610254578063ebdb56f31461024b5763f2fde38b1461024357600080fd5b61000e613173565b5061000e6130e1565b5061000e612fe2565b5061000e612edf565b5061000e612ec3565b5061000e612a71565b5061000e6125c7565b5061000e61259f565b5061000e612350565b5061000e610c1d565b5061000e612334565b5061000e6122ae565b5061000e611f32565b5061000e611e39565b5061000e611d18565b5061000e611b9c565b5061000e611b67565b5061000e611770565b5061000e611562565b5061000e611518565b5061000e6114f0565b5061000e61149e565b5061000e61146d565b5061000e611442565b5061000e6113f5565b5061000e611398565b5061000e61126b565b5061000e6111a7565b5061000e610fd6565b5061000e610f91565b5061000e610f40565b5061000e610d19565b5061000e610cca565b5061000e610c44565b5061000e610a0a565b5061000e6109c8565b5061000e610933565b5061000e61085b565b5061000e61083c565b5061000e610820565b5061000e6107fe565b5061000e610680565b5061000e61065a565b5061000e610625565b5061000e610588565b5061000e610529565b5061000e610479565b5061000e61044c565b503461000e57600060031936011261000e57600554604080517f0000000000000000000000000000000000000000000000000000000000000000815260609290921c69ffffffffffffffffffff16602083015290f35b0390f35b503461000e57600060031936011261000e5760206bffffffffffffffffffffffff60055416604051908152f35b503461000e57600060031936011261000e57602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b919082519283825260005b8481106105155750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b6020818301810151848301820152016104d6565b503461000e57600060031936011261000e5761044860405161054a81611db5565b600d81527f5374616b696e6720302e312e300000000000000000000000000000000000000060208201526040519182916020835260208301906104cb565b503461000e57600060031936011261000e5760206105a4615c70565b64e8d4a510006105f86105b56150a3565b926bffffffffffffffffffffffff600a5416938103908111610618575b6105f36105dd61396e565b9169ffffffffffffffffffff60095416906134d9565b6134d9565b04810390811161060b575b604051908152f35b6106136133c2565b610603565b6106206133c2565b6105d2565b503461000e57600060031936011261000e57602073ffffffffffffffffffffffffffffffffffffffff600e5416604051908152f35b503461000e57600060031936011261000e576020610676614dd8565b6040519015158152f35b503461000e576000806003193601126107fb5761069b614dd8565b6107c35761073a60206106fa6106f56106b33361586c565b92917f204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de0060806040513381528389820152846040820152866060820152a1613560565b613560565b6040517fa9059cbb000000000000000000000000000000000000000000000000000000008152336004820152602481019190915291829081906044820190565b03818573ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af180156107b6575b610788575b50604051f35b6107a89060203d81116107af575b6107a08183611dd1565b810190613a3d565b5038610782565b503d610796565b6107be613a55565b61077d565b604490604051907fa30a70c2000000000000000000000000000000000000000000000000000000008252600160048301526024820152fd5b80fd5b503461000e57600060031936011261000e57602060ff60085416604051908152f35b503461000e57600060031936011261000e576020610603615c70565b503461000e57600060031936011261000e576020601054604051908152f35b503461000e57600060031936011261000e576108756132a6565b60015460ff8160a01c16156108d5577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff166001557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152fd5b503461000e57600060031936011261000e5761094d6132a6565b740100000000000000000000000000000000000000007fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff60015461099760ff8260a01c161561423e565b16176001557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b503461000e57600060031936011261000e57600b546040805163ffffffff608084901c1681526fffffffffffffffffffffffffffffffff909216602083015290f35b503461000e576000806003193601126107fb57610a25614dd8565b6107c357610a7f610a6e610a593373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b5460701c6bffffffffffffffffffffffff1690565b6bffffffffffffffffffffffff1690565b8015610be857602081610af1610ab4610a9a61073a95613a20565b6006546fffffffffffffffffffffffffffffffff16614db6565b6fffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffffff000000000000000000000000000000006006541617600655565b610b46610b1e3373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b7fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff8154169055565b7f204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de0060405180610ba5843383606090600092949373ffffffffffffffffffffffffffffffffffffffff608083019616825260208201528260408201520152565b0390a16040517fa9059cbb000000000000000000000000000000000000000000000000000000008152336004820152602481019190915291829081906044820190565b6040517fe4adde72000000000000000000000000000000000000000000000000000000008152336004820152602490fd5b0390fd5b503461000e57600060031936011261000e57602060ff60015460a01c166040519015158152f35b503461000e57600060031936011261000e5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b90815180825260208080930193019160005b828110610ca0575050505090565b835173ffffffffffffffffffffffffffffffffffffffff1685529381019392810192600101610c92565b503461000e57600060031936011261000e57610448610ce76142a3565b604051918291602083526020830190610c80565b73ffffffffffffffffffffffffffffffffffffffff81160361000e57565b503461000e57602060031936011261000e57600435610d3781610cfb565b610d3f6132a6565b803b158015610f21575b8015610efd575b8015610ed9575b8015610e2b575b610e0157610dfc81610dcc7f5c74c441be501340b2713817a6c6975e6f3d4a4ae39fa1ac0bf75d3c54a0cad39373ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff0000000000000000000000000000000000000000600c541617600c55565b610dd542600d55565b60405173ffffffffffffffffffffffffffffffffffffffff90911681529081906020820190565b0390a1005b60046040517f367a1038000000000000000000000000000000000000000000000000000000008152fd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527fa4c0ed3600000000000000000000000000000000000000000000000000000000600482015260208160248173ffffffffffffffffffffffffffffffffffffffff86165afa908115610ecc575b600091610eae575b5015610d5e565b610ec6915060203d81116107af576107a08183611dd1565b38610ea7565b610ed4613a55565b610e9f565b50600e5473ffffffffffffffffffffffffffffffffffffffff828116911614610d57565b50600c5473ffffffffffffffffffffffffffffffffffffffff828116911614610d50565b503073ffffffffffffffffffffffffffffffffffffffff821614610d49565b503461000e57602060031936011261000e5773ffffffffffffffffffffffffffffffffffffffff600435610f7381610cfb565b166000526002602052602060ff604060002054166040519015158152f35b503461000e57600060031936011261000e576020610fad6150a3565b64e8d4a510006105f86bffffffffffffffffffffffff600a5460601c16926105f36105dd61396e565b503461000e57602060031936011261000e57600435610ff36132a6565b610ffb614dd8565b1561116f576040517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018290527fde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d91610dfc9160208160648160007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff165af18015611162575b611144575b506111346110ca610a6e6005546bffffffffffffffffffffffff1690565b6110d2615c70565b906110fa6110eb60095469ffffffffffffffffffff1690565b69ffffffffffffffffffff1690565b90611103614e01565b9161110c6150a3565b937f00000000000000000000000000000000000000000000000000000000000000009261356d565b6040519081529081906020820190565b61115b9060203d81116107af576107a08183611dd1565b50386110ac565b61116a613a55565b6110a7565b60446040517fa30a70c20000000000000000000000000000000000000000000000000000000081526000600482015260016024820152fd5b503461000e57602060031936011261000e576004356111c46132a6565b6111cc614dd8565b1561116f57801561000e5760207f1e3be2efa25bca5bff2215c7b30b31086e703d6aa7d9b9a1f8ba62c5291219ad916112626112066150a3565b61120f81613d04565b611217613c41565b6bffffffffffffffffffffffff60055416611230615c70565b908461123a614e01565b927f00000000000000000000000000000000000000000000000000000000000000009261356d565b604051908152a1005b503461000e576000806003193601126107fb5773ffffffffffffffffffffffffffffffffffffffff8060015416330361133a57815473ffffffffffffffffffffffffffffffffffffffff16600080547fffffffffffffffffffffffff0000000000000000000000000000000000000000163317905561130d7fffffffffffffffffffffffff000000000000000000000000000000000000000060015416600155565b604051913391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08484a3f35b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152fd5b503461000e57602060031936011261000e5773ffffffffffffffffffffffffffffffffffffffff6004356113cb81610cfb565b16600052600260205260206bffffffffffffffffffffffff60406000205460101c16604051908152f35b503461000e57602060031936011261000e577f1b930366dfeaa7eb3b325021e4ae81e36527063452ee55b86c95f85b36f4c31c60206004356114356132a6565b80601055604051908152a1005b503461000e57600060031936011261000e57602069ffffffffffffffffffff60095416604051908152f35b503461000e57600060031936011261000e5760206fffffffffffffffffffffffffffffffff60065416604051908152f35b503461000e57600060031936011261000e57602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b503461000e57602060031936011261000e57602061060360043561151381610cfb565b61501a565b503461000e57600060031936011261000e57600554604080517f0000000000000000000000000000000000000000000000000000000000000000815260b09290921c602083015290f35b503461000e57602060031936011261000e5767ffffffffffffffff60043581811161000e573660238201121561000e57806004013591821161000e576024810190602483369201011161000e576115b7614dd8565b6117385773ffffffffffffffffffffffffffffffffffffffff91826115f1600e5473ffffffffffffffffffffffffffffffffffffffff1690565b1615610e01576116e09260209260006106f56116aa61166b6116123361586c565b9194907f667838b33bdc898470de09e0e746990f2adc11b965b7fe6828e502ebc39e0434604051806116498d8c88878d3387614365565b0390a1600e5473ffffffffffffffffffffffffffffffffffffffff1695613560565b9361167e604051978892338b85016143a4565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101875286611dd1565b604051968795869485937f4000aea0000000000000000000000000000000000000000000000000000000008552600485016143d1565b03927f0000000000000000000000000000000000000000000000000000000000000000165af1801561172b575b611714575b005b6117129060203d81116107af576107a08183611dd1565b611733613a55565b61170d565b60446040517fa30a70c20000000000000000000000000000000000000000000000000000000081526001600482015260006024820152fd5b503461000e57606060031936011261000e57604435600480356024356117946132a6565b61179c614dd8565b15611b2f57818411611af7576005546bffffffffffffffffffffffff811690838211611ac2578269ffffffffffffffffffff8260601c1611611a8d57859060b01c11611a54576117ea615823565b9361182c611809876105f361180360208a015160ff1690565b60ff1690565b611826610a6e6040809901516bffffffffffffffffffffffff1690565b90613560565b8410611a20575093829161188693600096036119ab575b806118606110eb60055469ffffffffffffffffffff9060601c1690565b0361191a575b50806118776110eb60055460b01c90565b03611889575b506110d2615c70565b51f35b611911816119026118ba7f816587cb2e773af4f3689a03d7520fabff3462605ded374b485b13994c0d7b52946139c6565b75ffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffff000000000000000000000000000000000000000000006005549260b01b16911617600555565b85519081529081906020820190565b0390a13861187d565b6119a28161199361194b7fb5f554e5ef00806bace1edbb84186512ebcefa2af7706085143f501f29314df7946139c6565b7fffffffffffffffffffff00000000000000000000ffffffffffffffffffffffff75ffffffffffffffffffff0000000000000000000000006005549260601b16911617600555565b86519081529081906020820190565b0390a138611866565b6119f06119b784613a07565b6bffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffff0000000000000000000000006005541617600555565b84518381527f7f4f497e086b2eb55f8a9885ba00d33399bbe0ebcb92ea092834386435a1b9c090602090a1611843565b84517fbc91aa3300000000000000000000000000000000000000000000000000000000815290810186815281906020010390fd5b50506040517fbc91aa33000000000000000000000000000000000000000000000000000000008152918201928352509081906020010390fd5b6040517fbc91aa3300000000000000000000000000000000000000000000000000000000815280860184815281906020010390fd5b6040517f1f3c387600000000000000000000000000000000000000000000000000000000815280860185815281906020010390fd5b50506040517fbc91aa330000000000000000000000000000000000000000000000000000000081529081019182529081906020010390fd5b6044836000604051917fa30a70c200000000000000000000000000000000000000000000000000000000835282015260016024820152fd5b503461000e57600060031936011261000e57602073ffffffffffffffffffffffffffffffffffffffff60005416604051908152f35b503461000e57604060031936011261000e57600435611bb96132a6565b60105415611cee5761171290611bee7f0000000000000000000000000000000000000000000000000000000000000000615bd2565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526044810182905260208160648160007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff165af18015611ce1575b611cc3575b50611c8f610a6e6005546bffffffffffffffffffffffff1690565b611c97614e01565b917f00000000000000000000000000000000000000000000000000000000000000009160243590613a62565b611cda9060203d81116107af576107a08183611dd1565b5038611c74565b611ce9613a55565b611c6f565b60046040517f9f8a28f2000000000000000000000000000000000000000000000000000000008152fd5b503461000e57602060031936011261000e576020610603600435611d3b81610cfb565b614efd565b507f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6080810190811067ffffffffffffffff821117611d8c57604052565b611d94611d40565b604052565b6060810190811067ffffffffffffffff821117611d8c57604052565b6040810190811067ffffffffffffffff821117611d8c57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117611d8c57604052565b60209067ffffffffffffffff8111611e2c575b60051b0190565b611e34611d40565b611e25565b503461000e57604060031936011261000e57600435611e5781610cfb565b6024359067ffffffffffffffff821161000e573660238301121561000e578160040135611e8381611e12565b92611e916040519485611dd1565b81845260209160248386019160051b8301019136831161000e57602401905b828210611ed657610448611ec48787613325565b60405190151581529081906020820190565b81358152908301908301611eb0565b90602060031983011261000e5760043567ffffffffffffffff9283821161000e578060238301121561000e57816004013593841161000e5760248460051b8301011161000e576024019190565b503461000e57611f4136611ee5565b90611f4a6132a6565b63ffffffff600b5460801c16151580612256575b61116f57611f7a611f746110eb60055460b01c90565b836134d9565b611f82615c9e565b80821161221b57505060005b828110611feb57611712611fb7611fb2856106f561180360045460ff9060081c1690565b614230565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff006004549260081b16911617600455565b61203161202a612004611fff84878761415b565b614173565b73ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b5460ff1690565b6121be57612061610a6e61204c612004611fff85888861415b565b5460101c6bffffffffffffffffffffffff1690565b6121605761207c610a6e610a59612004611fff85888861415b565b61210257806120c3612098612004611fff6120fd95888861415b565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b7fac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d6120f5610dd5611fff84888861415b565b0390a16133f2565b611f8e565b611fff90610c19936121139361415b565b6040517f7a378b9c00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201529081906024820190565b611fff90610c19936121719361415b565b6040517f602d4d1100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201529081906024820190565b611fff90610c19936121cf9361415b565b6040517ea5216600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201529081906024820190565b6040517f35cf446b00000000000000000000000000000000000000000000000000000000815260048101919091526024810191909152604490fd5b5061225f614dd8565b15611f5e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f60209267ffffffffffffffff81116122a1575b01160190565b6122a9611d40565b61229b565b503461000e57606060031936011261000e576004356122cc81610cfb565b6044359067ffffffffffffffff821161000e573660238301121561000e578160040135906122f982612265565b916123076040519384611dd1565b808352366024828601011161000e57602081600092602461171297018387013784010152602435906150e3565b503461000e57600060031936011261000e5760206106036150a3565b503461000e5761235f36611ee5565b6123676132a6565b60005b6003548110156123dc57806123d26123aa61200461238a6123d795615cf9565b905473ffffffffffffffffffffffffffffffffffffffff9160031b1c1690565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff8154169055565b6133f2565b61236a565b506123e5615d3e565b60005b81811061242b57508161241c827f40aed8e423b39a56b445ae160f4c071fc2cfb48ee0b6dcd5ffeb6bc5b18d10d094615d91565b610dfc60405192839283615e6e565b612439611fff82848661415b565b61246d61246961202a8373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b1590565b612556576124a861249e8273ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b5460081c60ff1690565b61250e57906123d26124dd6125099373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b6101007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff825416179055565b6123e8565b6040517ea5216600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff919091166004820152602490fd5b6040517feac13dcd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff919091166004820152602490fd5b503461000e57602060031936011261000e5760206106766004356125c281610cfb565b614c20565b503461000e576125d636611ee5565b906125df6132a6565b6125e7614dd8565b1561116f576125fc6125f76150a3565b613d04565b60005b82811061262357611712611fb761261585614230565b60045460081c60ff166141f8565b612631611fff82858561415b565b61266361265e8273ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b61417d565b906126716124698351151590565b612556576020820151612a285761274e92917f2360404a74478febece1a14f11275f22ada88d19ef96f7d785913010bfff4479916120f56126c5610a6e6040809501516bffffffffffffffffffffffff1690565b9283612753575b6127216126f98473ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008154169055565b51928392836020909392919373ffffffffffffffffffffffffffffffffffffffff60408201951681520152565b6125ff565b61286561281b61276a61276586614efd565b613a07565b6127c761278e600a9261278984546bffffffffffffffffffffffff1690565b613e7b565b6bffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffff000000000000000000000000600a541617600a55565b6127f661278e6127e16127656127db61396e565b8b613906565b83546bffffffffffffffffffffffff16613e7b565b6127896128056127658861501a565b915460601c6bffffffffffffffffffffffff1690565b7fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff77ffffffffffffffffffffffff000000000000000000000000600a549260601b16911617600a55565b6128aa61287c61287760085460ff1690565b6141c7565b60ff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff006008541617600855565b6128ff6128d78473ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b7fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff8154169055565b612a23610ab46bffffffffffffffffffffffff61291b87613a07565b61298961293d826127896004546bffffffffffffffffffffffff9060701c1690565b7fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff79ffffffffffffffffffffffff00000000000000000000000000006004549260701b16911617600455565b612a02816129b78973ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b907fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff79ffffffffffffffffffffffff000000000000000000000000000083549260701b169116179055565b16612a1e6006546fffffffffffffffffffffffffffffffff1690565b61420c565b6126cc565b6040517fded6031900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff919091166004820152602490fd5b503461000e576000806003193601126107fb57612a8c614dd8565b15612e8b57612ac1610a6e61204c3373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b8015612e615760409081517ffeaf968c00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9160a082600481867f0000000000000000000000000000000000000000000000000000000000000000165afa918215612e54575b85908693612e14575b5069ffffffffffffffffffff1691600f548314612de357612b827f000000000000000000000000000000000000000000000000000000000000000082613560565b4210612d8457612bb3907f000000000000000000000000000000000000000000000000000000000000000090613560565b4210908180612dad575b612d845792827fd2720e8f454493f612cc97499fe8cbce7fa4d4c18d346fe7104e9042df1c1edd612c1d612c78948997612c18612c13612bfe60209a613a20565b6fffffffffffffffffffffffffffffffff1690565b600f55565b615b7a565b875133815260208101939093526040830181905291606090a185517fa9059cbb000000000000000000000000000000000000000000000000000000008152336004820152602481019190915293849283919082906044820190565b03927f0000000000000000000000000000000000000000000000000000000000000000165af18015612d77575b612d59575b50612d05612cb66150a3565b612cbe6142a3565b7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000614506565b611886612d22610a6e6005546bffffffffffffffffffffffff1690565b612d2a615c70565b612d426110eb60095469ffffffffffffffffffff1690565b612d4a614e01565b91612d536150a3565b936136bf565b612d709060203d81116107af576107a08183611dd1565b5038612caa565b612d7f613a55565b612ca5565b600485517ffc53c50a000000000000000000000000000000000000000000000000000000008152fd5b50612dde61246961202a3373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b612bbd565b84517ff3553c2200000000000000000000000000000000000000000000000000000000815260048101849052602490fd5b69ffffffffffffffffffff9350612e42915060a03d8111612e4d575b612e3a8183611dd1565b81019061441a565b50949150612b419050565b503d612e30565b612e5c613a55565b612b38565b60046040517fef67f5d8000000000000000000000000000000000000000000000000000000008152fd5b604490604051907fa30a70c2000000000000000000000000000000000000000000000000000000008252600482015260016024820152fd5b503461000e57600060031936011261000e576020610603614e01565b503461000e576000806003193601126107fb57612efa6132a6565b612f02614dd8565b15612e8b57612f3b612f12615c70565b612f1a6150a3565b90612f2482613d04565b612f2c613c41565b818103908111612fd557613eb7565b612f84612f4742613a20565b6fffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffffff00000000000000000000000000000000600b541617600b55565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600454166004556040517ff7d0e0f15586495da8c687328ead30fb829d9da55538cb0ef73dd229e517cdb88282a1f35b612fdd6133c2565b613eb7565b503461000e57600060031936011261000e57612ffc6132a6565b600c5473ffffffffffffffffffffffffffffffffffffffff168015610e0157613026600d54613550565b4210612e61576130719073ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff0000000000000000000000000000000000000000600e541617600e55565b61309e7fffffffffffffffffffffffff0000000000000000000000000000000000000000600c5416600c55565b7ffa33c052bbee754f3c0482a89962daffe749191fa33c696a61e947fbfd68bd84610dfc610dd5600e5473ffffffffffffffffffffffffffffffffffffffff1690565b503461000e576000806003193601126107fb576130fc6132a6565b613104614dd8565b6107c35761073a6020613141613118614e01565b600a54906131366bffffffffffffffffffffffff918284169061347e565b9160601c169061347e565b6040518181527f150a6ec0e6f4e9ddcaaaa1674f157d91165a42d60653016f87a9fc870a39f050908060208101610ba5565b503461000e57602060031936011261000e5773ffffffffffffffffffffffffffffffffffffffff6004356131a681610cfb565b6131ae6132a6565b1633811461324857806000917fffffffffffffffffffffffff0000000000000000000000000000000000000000600154161760015561321d613204835473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1690565b90604051917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788484a3f35b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152fd5b73ffffffffffffffffffffffffffffffffffffffff6000541633036132c757565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152fd5b90613364916010549060405173ffffffffffffffffffffffffffffffffffffffff60208201921682526020815261335b81611db5565b51902091613367565b90565b929091906000915b84518310156133ba57613382838661345c565b51908181116133a55760005260205261339f6040600020926133f2565b9161336f565b9060005260205261339f6040600020926133f2565b915092501490565b507f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114613420570190565b6134286133c2565b0190565b507f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020918151811015613471575b60051b010190565b61347961342c565b613469565b9190820391821161348b57565b6134936133c2565b565b64e8d4a5100090807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211811515166134cd570290565b6134d56133c2565b0290565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211811515166134cd570290565b507f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b8115613544570490565b61354c61350a565b0490565b9062093a80820180921161348b57565b9190820180921161348b57565b909392916135d2906135cc6135828888613929565b936135c66135a36135928b61399d565b8988829b039081116136b25761347e565b6135bb6110eb60095469ffffffffffffffffffff1690565b850361366a57613495565b926134d9565b9061353a565b9182106136405761281b61276561363b9464e8d4a5100061362e866105f38b61361961278e6127656134939f9b612f479f876136146136359f6118269361347e565b613906565b69ffffffffffffffffffff60095416906134d9565b0490613560565b42613560565b613a20565b60046040517fda056d00000000000000000000000000000000000000000000000000000000008152fd5b6136ad613676866139c6565b69ffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffffffff000000000000000000006009541617600955565b613495565b6136ba6133c2565b61347e565b928491926136cd8385613929565b906137106136da8561399d565b9683850394851161386b575b87850394851161385e575b6135c669ffffffffffffffffffff958660095416850361366a57613495565b908115613851575b046137258197829661347e565b9061372f91613906565b61373891613560565b61374190613a07565b61377e906bffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffff000000000000000000000000600a541617600a55565b6009541661378b916134d9565b90613795916134d9565b64e8d4a5100090046137a691613560565b6137af90613a07565b6137fd907fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff77ffffffffffffffffffffffff000000000000000000000000600a549260601b16911617600a55565b6138079042613560565b61381090613a20565b613493906fffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffffff00000000000000000000000000000000600b541617600b55565b61385961350a565b613718565b6138666133c2565b6136f1565b6138736133c2565b6136e6565b906138e2909392936135cc600a54936135c66bffffffffffffffffffffffff9569ffffffffffffffffffff6138d5888316986138b261396e565b506009549360601c16809a6138c561396e565b508a81039081116136b25761347e565b9116850361366a57613495565b928310613640576134939261363561281b612f479461276561278e61363b96613a07565b64e8d4a51000916105f361354c9269ffffffffffffffffffff60095416906134d9565b61395964e8d4a51000916bffffffffffffffffffffffff600a5416938103908111610618576105f36105dd61396e565b0481039081116139665790565b6133646133c2565b600b546fffffffffffffffffffffffffffffffff164281116139905750600090565b4281039081116139665790565b64e8d4a510006139596bffffffffffffffffffffffff600a5460601c16926105f36105dd61396e565b69ffffffffffffffffffff908181116139dd571690565b60046040517f811752de000000000000000000000000000000000000000000000000000000008152fd5b6bffffffffffffffffffffffff908181116139dd571690565b6fffffffffffffffffffffffffffffffff908181116139dd571690565b9081602091031261000e5751801515810361000e5790565b506040513d6000823e3d90fd5b9093919293600b549163ffffffff93848460801c1661000e577f125fc8494f786b470e3c39d0932a62e9e09e291ebd81ea19c57604f6d2b1d1679683613bc493613bbf613c0b9769ffffffffffffffffffff613abd856139c6565b167fffffffffffffffffffffffffffffffffffffffffffff000000000000000000006009541617600955613af042613c10565b907fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff73ffffffff000000000000000000000000000000008360801b16911617600b556008547fffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffff70ffffffff000000000000000000000000008360681b169116176008557fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff79ffffffff000000000000000000000000000000000000000000006009549260b01b16911617600955565b613878565b600b5460408051928352602083019590955263ffffffff608082811c90951616948201949094526fffffffffffffffffffffffffffffffff90931660608401528291820190565b0390a1565b63ffffffff908181116139dd571690565b9190916bffffffffffffffffffffffff8080941691160191821161348b57565b60095475ffffffffffffffffffffffff00000000000000000000613ca1613c8a61276563ffffffff8560b01c164203428111613cf7575b69ffffffffffffffffffff86166134d9565b6bffffffffffffffffffffffff8460501c16613c21565b60501b167fffffffffffff00000000000000000000000000000000ffffffffffffffffffff79ffffffff00000000000000000000000000000000000000000000613cea42613c10565b60b01b1692161717600955565b613cff6133c2565b613c78565b6008549060ff8216613d5f575b5050613d1c42613c10565b7fffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffff70ffffffff000000000000000000000000006008549260681b16911617600855565b6cffffffffffffffffffffffff00613d9a6127657fffffffffffffffffffffffffffffffffffffff000000000000000000000000ff93613dab565b60081b169116176008553880613d11565b600b5461336491906fffffffffffffffffffffffffffffffff16428111613e465764e8d4a51000613e216bffffffffffffffffffffffff9263ffffffff60085460681c168103908111613e39575b6105f36008549560ff87168015613e2c575b69ffffffffffffffffffff6009541691046134d9565b049160081c16613560565b613e3461350a565b613e0b565b613e416133c2565b613df9565b506bffffffffffffffffffffffff64e8d4a51000613e2163ffffffff60085460681c16420342811115613df957613e416133c2565b6bffffffffffffffffffffffff918216908216039190821161348b57565b6001906bffffffffffffffffffffffff809116908114613420570190565b7fffffffffffffffff00000000000000000000000000000000000000000000000077ffffffffffffffffffffffff000000000000000000000000613f8c613f40613f3a613f34613f5c966105f3613f0c61396e565b9169ffffffffffffffffffff60095416948591613f34856105f364e8d4a51000998a946134d9565b04613a07565b9a6134d9565b600a54956bffffffffffffffffffffffff958691828916613e7b565b1694857fffffffffffffffffffffffffffffffffffffffff00000000000000000000000088161760601c16613e7b565b60601b1692161717600a55565b61281b9061349392613fa961396e565b69ffffffffffffffffffff60095416916105f3613fe784613f34856105f3613fe08a613f3464e8d4a510009a8b998a9889946134d9565b98876134d9565b940661405b575b06614046575b614035907fffffffffffffffffffffffffffffffffffffffff000000000000000000000000600a54916bffffffffffffffffffffffff938491828516613c21565b1691161780600a5560601c16613c21565b9061405361403591613e99565b919050613ff4565b9161406590613e99565b91613fee565b6140e36134939164e8d4a5100061409f81613f3461408761396e565b6105f369ffffffffffffffffffff60095416876134d9565b910661414d575b600a547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000006bffffffffffffffffffffffff80948194828516613c21565b1691161780600a5560601c16908111614140577fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff77ffffffffffffffffffffffff000000000000000000000000600a549260601b16911617600a55565b6141486133c2565b61281b565b61415690613e99565b6140a6565b919081101561416b5760051b0190565b611e3461342c565b3561336481610cfb565b9060405161418a81611d70565b606081935460ff81161515835260ff8160081c16151560208401526bffffffffffffffffffffffff90818160101c16604085015260701c16910152565b60ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9116019060ff821161348b57565b9060ff8091169116039060ff821161348b57565b9190916fffffffffffffffffffffffffffffffff8080941691160191821161348b57565b60ff81116139dd5760ff1690565b1561424557565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152fd5b60405190600354808352826020918282019060036000527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b936000905b8282106142f65750505061349392500383611dd1565b855473ffffffffffffffffffffffffffffffffffffffff16845260019586019588955093810193909101906142e0565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b919260a09373ffffffffffffffffffffffffffffffffffffffff6133649896931684526020840152604083015260608201528160808201520191614326565b60409073ffffffffffffffffffffffffffffffffffffffff61336495931681528160208201520191614326565b613364939273ffffffffffffffffffffffffffffffffffffffff606093168252602082015281604082015201906104cb565b519069ffffffffffffffffffff8216820361000e57565b908160a091031261000e5761442e81614403565b91602082015191604081015191613364608060608401519301614403565b9061445682611e12565b6144636040519182611dd1565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06144918294611e12565b0190602036910137565b90815180825260208080930193019160005b8281106144bb575050505090565b8351855293810193928101926001016144ad565b916144f8906144ea6133649593606086526060860190610c80565b90848203602086015261449b565b91604081840391015261449b565b939291909360ff61451960085460ff1690565b16156146c2578264e8d4a51000614546876105f361454d9569ffffffffffffffffffff60095416906134d9565b04956146c9565b916000918261455c825161444c565b94614567835161444c565b92855b81518710156146415761459a614580888461345c565b5173ffffffffffffffffffffffffffffffffffffffff1690565b886145cb610a6e61204c8473ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b8c811561463257898461461a946145f78e6145f161462b9a976146259a6145fe9861473d565b9261345c565b52886148f4565b6146088b8a61345c565b526146138a8c61345c565b5190613560565b94614613898861345c565b966133f2565b959261456a565b50505050929561462b906133f2565b7e635ea9da6e262e92bb713d71840af7c567807ff35bf73e927490c61283248098995061469d919650613c0b955061281b92509261276561278e6146876146b696613a07565b600a546bffffffffffffffffffffffff16613e7b565b600a5460601c6bffffffffffffffffffffffff16613e7b565b604051938493846144cf565b5050509050565b90604051916146d783611d99565b60085492604063ffffffff60ff8616958684526bffffffffffffffffffffffff8160081c16602085015260681c16910152821560001461471957505050600090565b6105f361354c9264e8d4a510009469ffffffffffffffffffff6009541691046134d9565b9161486f61483261336493946147676fffffffffffffffffffffffffffffffff600b541642101590565b156148b3576147f66147eb6147af6fffffffffffffffffffffffffffffffff600b54166147a96147a060095463ffffffff9060b01c1690565b63ffffffff1690565b9061347e565b955b6147e56147cd6009549869ffffffffffffffffffff8a166134d9565b6bffffffffffffffffffffffff809960501c16613560565b906134d9565b64e8d4a51000900490565b90846148228873ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b541682039182116148a6576148e2565b9361486661483f86613a07565b9173ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b92835416613c21565b6bffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffff000000000000000000000000825416179055565b6148ae6133c2565b6148e2565b6147f66147eb63ffffffff60095460b01c1642034281116148d5575b956147b1565b6148dd6133c2565b6148cf565b90808210156148ef575090565b905090565b916149da6149c1613364939461491e6fffffffffffffffffffffffffffffffff600b541642101590565b15614a225761498261496a6149596fffffffffffffffffffffffffffffffff600b54166147a963ffffffff60085460681c1663ffffffff1690565b6136146008549760ff89169061353a565b6bffffffffffffffffffffffff809660081c16613560565b90846149ae8873ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b5460601c1682039182116148a6576148e2565b936149ce61483f86613a07565b92835460601c16613c21565b7fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff77ffffffffffffffffffffffff00000000000000000000000083549260601b169116179055565b61498261496a63ffffffff60085460681c1642034281111561495957614a466133c2565b614959565b600b549091906fffffffffffffffffffffffffffffffff16428111614af257614ad56147eb614a9f73ffffffffffffffffffffffffffffffffffffffff936147a96147a060095463ffffffff9060b01c1690565b935b6147e5614abd6009549669ffffffffffffffffffff88166134d9565b6bffffffffffffffffffffffff809760501c16613560565b921660005260076020526040600020541681039081116139665790565b5073ffffffffffffffffffffffffffffffffffffffff614ad56147eb63ffffffff60095460b01c164203428111614b2a575b93614aa1565b614b326133c2565b614b24565b600b54613364916147eb916fffffffffffffffffffffffffffffffff16428111614bb2576147e59063ffffffff60095460b01c168103908111614ba5575b6bffffffffffffffffffffffff614b9b6009549269ffffffffffffffffffff84166134d9565b9160501c16613560565b614bad6133c2565b614b75565b506147e563ffffffff60095460b01c16420342811115614b7557614bad6133c2565b90614bf373ffffffffffffffffffffffffffffffffffffffff91613dab565b911660005260076020526bffffffffffffffffffffffff60406000205460601c1681039081116139665790565b614c50610a6e61204c8373ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b15614db0576040517ffeaf968c00000000000000000000000000000000000000000000000000000000815260a08160048173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa908115614da3575b6000908192614d7e575b5069ffffffffffffffffffff600f54911614614d7757614d0c7f000000000000000000000000000000000000000000000000000000000000000082613560565b4210614d7757614d3d907f000000000000000000000000000000000000000000000000000000000000000090613560565b421015614d715761202a6133649173ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b50600190565b5050600090565b9050614d98915060a03d8111612e4d57612e3a8183611dd1565b509291505038614ccc565b614dab613a55565b614cc2565b50600090565b6fffffffffffffffffffffffffffffffff918216908216039190821161348b57565b60ff6004541680614de65790565b506fffffffffffffffffffffffffffffffff600b5416421090565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015614ef0575b600090614ebd575b6133649150614e88615c70565b8103908111614eb0575b6147a9612bfe6006546fffffffffffffffffffffffffffffffff1690565b614eb86133c2565b614e92565b6020823d8211614ee8575b81614ed560209383611dd1565b810103126107fb57506133649051614e7b565b3d9150614ec8565b614ef8613a55565b614e73565b73ffffffffffffffffffffffffffffffffffffffff811660005260026020526bffffffffffffffffffffffff908160406000205460101c16918215614fef5760ff6040600020541615614f5557509061336491614a4b565b6147a990614fc2614f9b613364957f00000000000000000000000000000000000000000000000000000000000000008015614fe2575b81048103908111614fd557614b37565b9373ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b54166bffffffffffffffffffffffff1690565b614fdd6133c2565b614b37565b614fea61350a565b614f8b565b505050600090565b90801561500d575b810481039081116139665790565b61501561350a565b614fff565b73ffffffffffffffffffffffffffffffffffffffff8116600052600260205260406000206040519061504b82611d70565b549060ff821615908115815260ff8360081c161515602082015260606bffffffffffffffffffffffff808560101c169485604085015260701c16910152614d775715614db0576133649061509d6150a3565b90614bd4565b6bffffffffffffffffffffffff60045460101c167f0000000000000000000000000000000000000000000000000000000000000000908115613544570490565b9173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633036152205761513360ff60015460a01c161561423e565b61513b614dd8565b1561116f5761517061246961202a8573ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b1561521657601054908161518a575b505061349391615305565b805115612e615761520b916151ab826020806124699551830101910161524a565b90604051602081019061335b816151df8a8591909173ffffffffffffffffffffffffffffffffffffffff6020820193169052565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282611dd1565b612e6157388061517f565b5061349391615642565b60046040517f4d695438000000000000000000000000000000000000000000000000000000008152fd5b602090818184031261000e5780519067ffffffffffffffff821161000e57019180601f8401121561000e57825161528081611e12565b9361528e6040519586611dd1565b818552838086019260051b82010192831161000e578301905b8282106152b5575050505090565b815181529083019083016152a7565b604051906152d182611d99565b8160406005546bffffffffffffffffffffffff8116835269ffffffffffffffffffff8160601c16602084015260b01c910152565b9190615337610a6e61204c8573ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b926153428285613560565b937f00000000000000000000000000000000000000000000000000000000000000008086106155fd575061538d6110eb602061537c6152c4565b015169ffffffffffffffffffff1690565b8086116155b957505061539e615c9e565b8083116155865750613c0b7f1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee9093946153d76125f76150a3565b6154767f000000000000000000000000000000000000000000000000000000000000000061547061541161540b838961353a565b8861347e565b9161546a61542161276585614b37565b61486f61544e8a73ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b9161546583546bffffffffffffffffffffffff1690565b613c21565b8761353a565b90613f99565b6154de61549e61548586613a07565b60045460101c6bffffffffffffffffffffffff16613c21565b7fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff6dffffffffffffffffffffffff00006004549260101b16911617600455565b6155536154ea82613a07565b6155148573ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b907fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff6dffffffffffffffffffffffff000083549260101b169116179055565b6040519384938460409194939273ffffffffffffffffffffffffffffffffffffffff606083019616825260208201520152565b6040517fb94339d80000000000000000000000000000000000000000000000000000000081526004810191909152602490fd5b610c19916155c69161347e565b6040519182917fb94339d8000000000000000000000000000000000000000000000000000000008352600483019190602083019252565b6040517f1d820b170000000000000000000000000000000000000000000000000000000081526004810191909152602490fd5b60ff6001911660ff8114613420570190565b919061568c610a6e604061567961265e8773ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b01516bffffffffffffffffffffffff1690565b926156978285613560565b937f00000000000000000000000000000000000000000000000000000000000000008086106155fd57506156d06110eb60055460b01c90565b8086116155b957507f1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee909394613c0b9115615772575b61574161571461276586614b37565b61486f61544e8673ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b61576961293d61575086613a07565b60045460701c6bffffffffffffffffffffffff16613c21565b6154de8461406b565b61577d6125f76150a3565b61579461287c61578f60085460ff1690565b615630565b6008805461581e911c6bffffffffffffffffffffffff166157d58573ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b907fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff77ffffffffffffffffffffffff00000000000000000000000083549260601b169116179055565b615705565b6040519061583082611d70565b81606060045460ff81161515835260ff8160081c1660208401526bffffffffffffffffffffffff90818160101c16604085015260701c16910152565b61589961265e8273ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b91604083016bffffffffffffffffffffffff806158c283516bffffffffffffffffffffffff1690565b1615615adb576158da6158d3615823565b9551151590565b15615a2e5750610a6e8161591c61293d61590361598c95516bffffffffffffffffffffffff1690565b60045460701c6bffffffffffffffffffffffff16613e7b565b855115615a035761593761592f86614efd565b965b51151590565b156159915761597a6128d761595361594d6150a3565b88614bd4565b9673ffffffffffffffffffffffffffffffffffffffff166000526002602052604060002090565b516bffffffffffffffffffffffff1690565b929190565b61597a6128d76159fd6159b8610a6e6008546bffffffffffffffffffffffff9060081c1690565b6147a9610a6e6159e88b73ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b5460601c6bffffffffffffffffffffffff1690565b96612004565b615937615a2886615a23610a6e85516bffffffffffffffffffffffff1690565b615b22565b96615931565b90610a6e9184615a72615a8895969761593161549e615a5987516bffffffffffffffffffffffff1690565b60045460101c6bffffffffffffffffffffffff16613e7b565b15615a8f5761597a91506159536128d791614efd565b9190600090565b6159fd6128d791615a2361597a94615ab387516bffffffffffffffffffffffff1690565b7f00000000000000000000000000000000000000000000000000000000000000009116614ff7565b6040517fe4adde7200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602490fd5b73ffffffffffffffffffffffffffffffffffffffff64e8d4a51000615b5c6bffffffffffffffffffffffff938460095460501c16906134d9565b04921660005260076020526040600020541681039081116139665790565b90615bad57613364907f00000000000000000000000000000000000000000000000000000000000000009060011c6148e2565b507f000000000000000000000000000000000000000000000000000000000000000090565b6004549060ff8260081c1690808210615c395750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060019116176004557fded6ebf04e261e1eb2f3e3b268a2e6aee5b478c15b341eba5cf18b9bc80c2e636000604051a1565b60449250604051917fe709379900000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b615c78615823565b6bffffffffffffffffffffffff6060816040840151169201511681018091116139665790565b615ca6615823565b6005549060406bffffffffffffffffffffffff91615cd38385169460ff6020840151169060b01c906134d9565b8403938411615cec575b01511681039081116139665790565b615cf46133c2565b615cdd565b600354811015615d31575b60036000527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b0190600090565b615d3961342c565b615d04565b60035460008060035581615d50575050565b600381527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b918201915b828110615d8657505050565b818155600101615d7a565b67ffffffffffffffff8211615e61575b680100000000000000008211615e54575b60035482600355808310615e13575b50600360005260005b828110615dd657505050565b6001906020833593615de785610cfb565b0192817fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b015501615dca565b827fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b91820191015b818110615e485750615dc1565b60008155600101615e3b565b615e5c611d40565b615db2565b615e69611d40565b615da1565b90916040602092828482018583525201929160005b828110615e91575050505090565b90919293828060019273ffffffffffffffffffffffffffffffffffffffff8835615eba81610cfb565b16815201950193929101615e8356fea164736f6c6343000810000a",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// StakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingMetaData.Bin instead.
var StakingBin = StakingMetaData.Bin

// DeployStaking deploys a new Ethereum contract, binding an instance of Staking to it.
func DeployStaking(auth *bind.TransactOpts, backend bind.ContractBackend, params StakingPoolConstructorParams) (common.Address, *types.Transaction, *Staking, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingBin), backend, params)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// CanAlert is a free data retrieval call binding the contract method 0xc1852f58.
//
// Solidity: function canAlert(address alerter) view returns(bool)
func (_Staking *StakingCaller) CanAlert(opts *bind.CallOpts, alerter common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "canAlert", alerter)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanAlert is a free data retrieval call binding the contract method 0xc1852f58.
//
// Solidity: function canAlert(address alerter) view returns(bool)
func (_Staking *StakingSession) CanAlert(alerter common.Address) (bool, error) {
	return _Staking.Contract.CanAlert(&_Staking.CallOpts, alerter)
}

// CanAlert is a free data retrieval call binding the contract method 0xc1852f58.
//
// Solidity: function canAlert(address alerter) view returns(bool)
func (_Staking *StakingCallerSession) CanAlert(alerter common.Address) (bool, error) {
	return _Staking.Contract.CanAlert(&_Staking.CallOpts, alerter)
}

// GetAvailableReward is a free data retrieval call binding the contract method 0xe0974ea5.
//
// Solidity: function getAvailableReward() view returns(uint256)
func (_Staking *StakingCaller) GetAvailableReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getAvailableReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAvailableReward is a free data retrieval call binding the contract method 0xe0974ea5.
//
// Solidity: function getAvailableReward() view returns(uint256)
func (_Staking *StakingSession) GetAvailableReward() (*big.Int, error) {
	return _Staking.Contract.GetAvailableReward(&_Staking.CallOpts)
}

// GetAvailableReward is a free data retrieval call binding the contract method 0xe0974ea5.
//
// Solidity: function getAvailableReward() view returns(uint256)
func (_Staking *StakingCallerSession) GetAvailableReward() (*big.Int, error) {
	return _Staking.Contract.GetAvailableReward(&_Staking.CallOpts)
}

// GetBaseReward is a free data retrieval call binding the contract method 0x9a109bc2.
//
// Solidity: function getBaseReward(address staker) view returns(uint256)
func (_Staking *StakingCaller) GetBaseReward(opts *bind.CallOpts, staker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getBaseReward", staker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBaseReward is a free data retrieval call binding the contract method 0x9a109bc2.
//
// Solidity: function getBaseReward(address staker) view returns(uint256)
func (_Staking *StakingSession) GetBaseReward(staker common.Address) (*big.Int, error) {
	return _Staking.Contract.GetBaseReward(&_Staking.CallOpts, staker)
}

// GetBaseReward is a free data retrieval call binding the contract method 0x9a109bc2.
//
// Solidity: function getBaseReward(address staker) view returns(uint256)
func (_Staking *StakingCallerSession) GetBaseReward(staker common.Address) (*big.Int, error) {
	return _Staking.Contract.GetBaseReward(&_Staking.CallOpts, staker)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Staking *StakingCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Staking *StakingSession) GetChainlinkToken() (common.Address, error) {
	return _Staking.Contract.GetChainlinkToken(&_Staking.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_Staking *StakingCallerSession) GetChainlinkToken() (common.Address, error) {
	return _Staking.Contract.GetChainlinkToken(&_Staking.CallOpts)
}

// GetCommunityStakerLimits is a free data retrieval call binding the contract method 0x0641bdd8.
//
// Solidity: function getCommunityStakerLimits() view returns(uint256, uint256)
func (_Staking *StakingCaller) GetCommunityStakerLimits(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getCommunityStakerLimits")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetCommunityStakerLimits is a free data retrieval call binding the contract method 0x0641bdd8.
//
// Solidity: function getCommunityStakerLimits() view returns(uint256, uint256)
func (_Staking *StakingSession) GetCommunityStakerLimits() (*big.Int, *big.Int, error) {
	return _Staking.Contract.GetCommunityStakerLimits(&_Staking.CallOpts)
}

// GetCommunityStakerLimits is a free data retrieval call binding the contract method 0x0641bdd8.
//
// Solidity: function getCommunityStakerLimits() view returns(uint256, uint256)
func (_Staking *StakingCallerSession) GetCommunityStakerLimits() (*big.Int, *big.Int, error) {
	return _Staking.Contract.GetCommunityStakerLimits(&_Staking.CallOpts)
}

// GetDelegatesCount is a free data retrieval call binding the contract method 0x32e28850.
//
// Solidity: function getDelegatesCount() view returns(uint256)
func (_Staking *StakingCaller) GetDelegatesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDelegatesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDelegatesCount is a free data retrieval call binding the contract method 0x32e28850.
//
// Solidity: function getDelegatesCount() view returns(uint256)
func (_Staking *StakingSession) GetDelegatesCount() (*big.Int, error) {
	return _Staking.Contract.GetDelegatesCount(&_Staking.CallOpts)
}

// GetDelegatesCount is a free data retrieval call binding the contract method 0x32e28850.
//
// Solidity: function getDelegatesCount() view returns(uint256)
func (_Staking *StakingCallerSession) GetDelegatesCount() (*big.Int, error) {
	return _Staking.Contract.GetDelegatesCount(&_Staking.CallOpts)
}

// GetDelegationRateDenominator is a free data retrieval call binding the contract method 0x5e8b40d7.
//
// Solidity: function getDelegationRateDenominator() view returns(uint256)
func (_Staking *StakingCaller) GetDelegationRateDenominator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDelegationRateDenominator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDelegationRateDenominator is a free data retrieval call binding the contract method 0x5e8b40d7.
//
// Solidity: function getDelegationRateDenominator() view returns(uint256)
func (_Staking *StakingSession) GetDelegationRateDenominator() (*big.Int, error) {
	return _Staking.Contract.GetDelegationRateDenominator(&_Staking.CallOpts)
}

// GetDelegationRateDenominator is a free data retrieval call binding the contract method 0x5e8b40d7.
//
// Solidity: function getDelegationRateDenominator() view returns(uint256)
func (_Staking *StakingCallerSession) GetDelegationRateDenominator() (*big.Int, error) {
	return _Staking.Contract.GetDelegationRateDenominator(&_Staking.CallOpts)
}

// GetDelegationReward is a free data retrieval call binding the contract method 0x87e900b1.
//
// Solidity: function getDelegationReward(address staker) view returns(uint256)
func (_Staking *StakingCaller) GetDelegationReward(opts *bind.CallOpts, staker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDelegationReward", staker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDelegationReward is a free data retrieval call binding the contract method 0x87e900b1.
//
// Solidity: function getDelegationReward(address staker) view returns(uint256)
func (_Staking *StakingSession) GetDelegationReward(staker common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDelegationReward(&_Staking.CallOpts, staker)
}

// GetDelegationReward is a free data retrieval call binding the contract method 0x87e900b1.
//
// Solidity: function getDelegationReward(address staker) view returns(uint256)
func (_Staking *StakingCallerSession) GetDelegationReward(staker common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDelegationReward(&_Staking.CallOpts, staker)
}

// GetEarnedBaseRewards is a free data retrieval call binding the contract method 0x1a9d4c7c.
//
// Solidity: function getEarnedBaseRewards() view returns(uint256)
func (_Staking *StakingCaller) GetEarnedBaseRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getEarnedBaseRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEarnedBaseRewards is a free data retrieval call binding the contract method 0x1a9d4c7c.
//
// Solidity: function getEarnedBaseRewards() view returns(uint256)
func (_Staking *StakingSession) GetEarnedBaseRewards() (*big.Int, error) {
	return _Staking.Contract.GetEarnedBaseRewards(&_Staking.CallOpts)
}

// GetEarnedBaseRewards is a free data retrieval call binding the contract method 0x1a9d4c7c.
//
// Solidity: function getEarnedBaseRewards() view returns(uint256)
func (_Staking *StakingCallerSession) GetEarnedBaseRewards() (*big.Int, error) {
	return _Staking.Contract.GetEarnedBaseRewards(&_Staking.CallOpts)
}

// GetEarnedDelegationRewards is a free data retrieval call binding the contract method 0x74104002.
//
// Solidity: function getEarnedDelegationRewards() view returns(uint256)
func (_Staking *StakingCaller) GetEarnedDelegationRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getEarnedDelegationRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEarnedDelegationRewards is a free data retrieval call binding the contract method 0x74104002.
//
// Solidity: function getEarnedDelegationRewards() view returns(uint256)
func (_Staking *StakingSession) GetEarnedDelegationRewards() (*big.Int, error) {
	return _Staking.Contract.GetEarnedDelegationRewards(&_Staking.CallOpts)
}

// GetEarnedDelegationRewards is a free data retrieval call binding the contract method 0x74104002.
//
// Solidity: function getEarnedDelegationRewards() view returns(uint256)
func (_Staking *StakingCallerSession) GetEarnedDelegationRewards() (*big.Int, error) {
	return _Staking.Contract.GetEarnedDelegationRewards(&_Staking.CallOpts)
}

// GetFeedOperators is a free data retrieval call binding the contract method 0x5fec60f8.
//
// Solidity: function getFeedOperators() view returns(address[])
func (_Staking *StakingCaller) GetFeedOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getFeedOperators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetFeedOperators is a free data retrieval call binding the contract method 0x5fec60f8.
//
// Solidity: function getFeedOperators() view returns(address[])
func (_Staking *StakingSession) GetFeedOperators() ([]common.Address, error) {
	return _Staking.Contract.GetFeedOperators(&_Staking.CallOpts)
}

// GetFeedOperators is a free data retrieval call binding the contract method 0x5fec60f8.
//
// Solidity: function getFeedOperators() view returns(address[])
func (_Staking *StakingCallerSession) GetFeedOperators() ([]common.Address, error) {
	return _Staking.Contract.GetFeedOperators(&_Staking.CallOpts)
}

// GetMaxPoolSize is a free data retrieval call binding the contract method 0x0fbc8f5b.
//
// Solidity: function getMaxPoolSize() view returns(uint256)
func (_Staking *StakingCaller) GetMaxPoolSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getMaxPoolSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPoolSize is a free data retrieval call binding the contract method 0x0fbc8f5b.
//
// Solidity: function getMaxPoolSize() view returns(uint256)
func (_Staking *StakingSession) GetMaxPoolSize() (*big.Int, error) {
	return _Staking.Contract.GetMaxPoolSize(&_Staking.CallOpts)
}

// GetMaxPoolSize is a free data retrieval call binding the contract method 0x0fbc8f5b.
//
// Solidity: function getMaxPoolSize() view returns(uint256)
func (_Staking *StakingCallerSession) GetMaxPoolSize() (*big.Int, error) {
	return _Staking.Contract.GetMaxPoolSize(&_Staking.CallOpts)
}

// GetMerkleRoot is a free data retrieval call binding the contract method 0x49590657.
//
// Solidity: function getMerkleRoot() view returns(bytes32)
func (_Staking *StakingCaller) GetMerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getMerkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetMerkleRoot is a free data retrieval call binding the contract method 0x49590657.
//
// Solidity: function getMerkleRoot() view returns(bytes32)
func (_Staking *StakingSession) GetMerkleRoot() ([32]byte, error) {
	return _Staking.Contract.GetMerkleRoot(&_Staking.CallOpts)
}

// GetMerkleRoot is a free data retrieval call binding the contract method 0x49590657.
//
// Solidity: function getMerkleRoot() view returns(bytes32)
func (_Staking *StakingCallerSession) GetMerkleRoot() ([32]byte, error) {
	return _Staking.Contract.GetMerkleRoot(&_Staking.CallOpts)
}

// GetMigrationTarget is a free data retrieval call binding the contract method 0x1ddb5552.
//
// Solidity: function getMigrationTarget() view returns(address)
func (_Staking *StakingCaller) GetMigrationTarget(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getMigrationTarget")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMigrationTarget is a free data retrieval call binding the contract method 0x1ddb5552.
//
// Solidity: function getMigrationTarget() view returns(address)
func (_Staking *StakingSession) GetMigrationTarget() (common.Address, error) {
	return _Staking.Contract.GetMigrationTarget(&_Staking.CallOpts)
}

// GetMigrationTarget is a free data retrieval call binding the contract method 0x1ddb5552.
//
// Solidity: function getMigrationTarget() view returns(address)
func (_Staking *StakingCallerSession) GetMigrationTarget() (common.Address, error) {
	return _Staking.Contract.GetMigrationTarget(&_Staking.CallOpts)
}

// GetMonitoredFeed is a free data retrieval call binding the contract method 0x83db28a0.
//
// Solidity: function getMonitoredFeed() view returns(address)
func (_Staking *StakingCaller) GetMonitoredFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getMonitoredFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMonitoredFeed is a free data retrieval call binding the contract method 0x83db28a0.
//
// Solidity: function getMonitoredFeed() view returns(address)
func (_Staking *StakingSession) GetMonitoredFeed() (common.Address, error) {
	return _Staking.Contract.GetMonitoredFeed(&_Staking.CallOpts)
}

// GetMonitoredFeed is a free data retrieval call binding the contract method 0x83db28a0.
//
// Solidity: function getMonitoredFeed() view returns(address)
func (_Staking *StakingCallerSession) GetMonitoredFeed() (common.Address, error) {
	return _Staking.Contract.GetMonitoredFeed(&_Staking.CallOpts)
}

// GetOperatorLimits is a free data retrieval call binding the contract method 0x8856398f.
//
// Solidity: function getOperatorLimits() view returns(uint256, uint256)
func (_Staking *StakingCaller) GetOperatorLimits(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getOperatorLimits")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetOperatorLimits is a free data retrieval call binding the contract method 0x8856398f.
//
// Solidity: function getOperatorLimits() view returns(uint256, uint256)
func (_Staking *StakingSession) GetOperatorLimits() (*big.Int, *big.Int, error) {
	return _Staking.Contract.GetOperatorLimits(&_Staking.CallOpts)
}

// GetOperatorLimits is a free data retrieval call binding the contract method 0x8856398f.
//
// Solidity: function getOperatorLimits() view returns(uint256, uint256)
func (_Staking *StakingCallerSession) GetOperatorLimits() (*big.Int, *big.Int, error) {
	return _Staking.Contract.GetOperatorLimits(&_Staking.CallOpts)
}

// GetRewardRate is a free data retrieval call binding the contract method 0x7e1a3786.
//
// Solidity: function getRewardRate() view returns(uint256)
func (_Staking *StakingCaller) GetRewardRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRewardRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardRate is a free data retrieval call binding the contract method 0x7e1a3786.
//
// Solidity: function getRewardRate() view returns(uint256)
func (_Staking *StakingSession) GetRewardRate() (*big.Int, error) {
	return _Staking.Contract.GetRewardRate(&_Staking.CallOpts)
}

// GetRewardRate is a free data retrieval call binding the contract method 0x7e1a3786.
//
// Solidity: function getRewardRate() view returns(uint256)
func (_Staking *StakingCallerSession) GetRewardRate() (*big.Int, error) {
	return _Staking.Contract.GetRewardRate(&_Staking.CallOpts)
}

// GetRewardTimestamps is a free data retrieval call binding the contract method 0x59f01879.
//
// Solidity: function getRewardTimestamps() view returns(uint256, uint256)
func (_Staking *StakingCaller) GetRewardTimestamps(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRewardTimestamps")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetRewardTimestamps is a free data retrieval call binding the contract method 0x59f01879.
//
// Solidity: function getRewardTimestamps() view returns(uint256, uint256)
func (_Staking *StakingSession) GetRewardTimestamps() (*big.Int, *big.Int, error) {
	return _Staking.Contract.GetRewardTimestamps(&_Staking.CallOpts)
}

// GetRewardTimestamps is a free data retrieval call binding the contract method 0x59f01879.
//
// Solidity: function getRewardTimestamps() view returns(uint256, uint256)
func (_Staking *StakingCallerSession) GetRewardTimestamps() (*big.Int, *big.Int, error) {
	return _Staking.Contract.GetRewardTimestamps(&_Staking.CallOpts)
}

// GetStake is a free data retrieval call binding the contract method 0x7a766460.
//
// Solidity: function getStake(address staker) view returns(uint256)
func (_Staking *StakingCaller) GetStake(opts *bind.CallOpts, staker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getStake", staker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStake is a free data retrieval call binding the contract method 0x7a766460.
//
// Solidity: function getStake(address staker) view returns(uint256)
func (_Staking *StakingSession) GetStake(staker common.Address) (*big.Int, error) {
	return _Staking.Contract.GetStake(&_Staking.CallOpts, staker)
}

// GetStake is a free data retrieval call binding the contract method 0x7a766460.
//
// Solidity: function getStake(address staker) view returns(uint256)
func (_Staking *StakingCallerSession) GetStake(staker common.Address) (*big.Int, error) {
	return _Staking.Contract.GetStake(&_Staking.CallOpts, staker)
}

// GetTotalDelegatedAmount is a free data retrieval call binding the contract method 0xa7a2f5aa.
//
// Solidity: function getTotalDelegatedAmount() view returns(uint256)
func (_Staking *StakingCaller) GetTotalDelegatedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getTotalDelegatedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalDelegatedAmount is a free data retrieval call binding the contract method 0xa7a2f5aa.
//
// Solidity: function getTotalDelegatedAmount() view returns(uint256)
func (_Staking *StakingSession) GetTotalDelegatedAmount() (*big.Int, error) {
	return _Staking.Contract.GetTotalDelegatedAmount(&_Staking.CallOpts)
}

// GetTotalDelegatedAmount is a free data retrieval call binding the contract method 0xa7a2f5aa.
//
// Solidity: function getTotalDelegatedAmount() view returns(uint256)
func (_Staking *StakingCallerSession) GetTotalDelegatedAmount() (*big.Int, error) {
	return _Staking.Contract.GetTotalDelegatedAmount(&_Staking.CallOpts)
}

// GetTotalRemovedAmount is a free data retrieval call binding the contract method 0x8019e7d0.
//
// Solidity: function getTotalRemovedAmount() view returns(uint256)
func (_Staking *StakingCaller) GetTotalRemovedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getTotalRemovedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalRemovedAmount is a free data retrieval call binding the contract method 0x8019e7d0.
//
// Solidity: function getTotalRemovedAmount() view returns(uint256)
func (_Staking *StakingSession) GetTotalRemovedAmount() (*big.Int, error) {
	return _Staking.Contract.GetTotalRemovedAmount(&_Staking.CallOpts)
}

// GetTotalRemovedAmount is a free data retrieval call binding the contract method 0x8019e7d0.
//
// Solidity: function getTotalRemovedAmount() view returns(uint256)
func (_Staking *StakingCallerSession) GetTotalRemovedAmount() (*big.Int, error) {
	return _Staking.Contract.GetTotalRemovedAmount(&_Staking.CallOpts)
}

// GetTotalStakedAmount is a free data retrieval call binding the contract method 0x38adb6f0.
//
// Solidity: function getTotalStakedAmount() view returns(uint256)
func (_Staking *StakingCaller) GetTotalStakedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getTotalStakedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalStakedAmount is a free data retrieval call binding the contract method 0x38adb6f0.
//
// Solidity: function getTotalStakedAmount() view returns(uint256)
func (_Staking *StakingSession) GetTotalStakedAmount() (*big.Int, error) {
	return _Staking.Contract.GetTotalStakedAmount(&_Staking.CallOpts)
}

// GetTotalStakedAmount is a free data retrieval call binding the contract method 0x38adb6f0.
//
// Solidity: function getTotalStakedAmount() view returns(uint256)
func (_Staking *StakingCallerSession) GetTotalStakedAmount() (*big.Int, error) {
	return _Staking.Contract.GetTotalStakedAmount(&_Staking.CallOpts)
}

// HasAccess is a free data retrieval call binding the contract method 0x9d0a3864.
//
// Solidity: function hasAccess(address staker, bytes32[] proof) view returns(bool)
func (_Staking *StakingCaller) HasAccess(opts *bind.CallOpts, staker common.Address, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "hasAccess", staker, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x9d0a3864.
//
// Solidity: function hasAccess(address staker, bytes32[] proof) view returns(bool)
func (_Staking *StakingSession) HasAccess(staker common.Address, proof [][32]byte) (bool, error) {
	return _Staking.Contract.HasAccess(&_Staking.CallOpts, staker, proof)
}

// HasAccess is a free data retrieval call binding the contract method 0x9d0a3864.
//
// Solidity: function hasAccess(address staker, bytes32[] proof) view returns(bool)
func (_Staking *StakingCallerSession) HasAccess(staker common.Address, proof [][32]byte) (bool, error) {
	return _Staking.Contract.HasAccess(&_Staking.CallOpts, staker, proof)
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() view returns(bool)
func (_Staking *StakingCaller) IsActive(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "isActive")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() view returns(bool)
func (_Staking *StakingSession) IsActive() (bool, error) {
	return _Staking.Contract.IsActive(&_Staking.CallOpts)
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() view returns(bool)
func (_Staking *StakingCallerSession) IsActive() (bool, error) {
	return _Staking.Contract.IsActive(&_Staking.CallOpts)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address staker) view returns(bool)
func (_Staking *StakingCaller) IsOperator(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "isOperator", staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address staker) view returns(bool)
func (_Staking *StakingSession) IsOperator(staker common.Address) (bool, error) {
	return _Staking.Contract.IsOperator(&_Staking.CallOpts, staker)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address staker) view returns(bool)
func (_Staking *StakingCallerSession) IsOperator(staker common.Address) (bool, error) {
	return _Staking.Contract.IsOperator(&_Staking.CallOpts, staker)
}

// IsPaused is a free data retrieval call binding the contract method 0xb187bd26.
//
// Solidity: function isPaused() view returns(bool)
func (_Staking *StakingCaller) IsPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "isPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPaused is a free data retrieval call binding the contract method 0xb187bd26.
//
// Solidity: function isPaused() view returns(bool)
func (_Staking *StakingSession) IsPaused() (bool, error) {
	return _Staking.Contract.IsPaused(&_Staking.CallOpts)
}

// IsPaused is a free data retrieval call binding the contract method 0xb187bd26.
//
// Solidity: function isPaused() view returns(bool)
func (_Staking *StakingCallerSession) IsPaused() (bool, error) {
	return _Staking.Contract.IsPaused(&_Staking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Staking *StakingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Staking *StakingSession) Owner() (common.Address, error) {
	return _Staking.Contract.Owner(&_Staking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Staking *StakingCallerSession) Owner() (common.Address, error) {
	return _Staking.Contract.Owner(&_Staking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Staking *StakingCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Staking *StakingSession) Paused() (bool, error) {
	return _Staking.Contract.Paused(&_Staking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Staking *StakingCallerSession) Paused() (bool, error) {
	return _Staking.Contract.Paused(&_Staking.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Staking *StakingCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Staking *StakingSession) TypeAndVersion() (string, error) {
	return _Staking.Contract.TypeAndVersion(&_Staking.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Staking *StakingCallerSession) TypeAndVersion() (string, error) {
	return _Staking.Contract.TypeAndVersion(&_Staking.CallOpts)
}

// AcceptMigrationTarget is a paid mutator transaction binding the contract method 0xe937fdaa.
//
// Solidity: function acceptMigrationTarget() returns()
func (_Staking *StakingTransactor) AcceptMigrationTarget(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "acceptMigrationTarget")
}

// AcceptMigrationTarget is a paid mutator transaction binding the contract method 0xe937fdaa.
//
// Solidity: function acceptMigrationTarget() returns()
func (_Staking *StakingSession) AcceptMigrationTarget() (*types.Transaction, error) {
	return _Staking.Contract.AcceptMigrationTarget(&_Staking.TransactOpts)
}

// AcceptMigrationTarget is a paid mutator transaction binding the contract method 0xe937fdaa.
//
// Solidity: function acceptMigrationTarget() returns()
func (_Staking *StakingTransactorSession) AcceptMigrationTarget() (*types.Transaction, error) {
	return _Staking.Contract.AcceptMigrationTarget(&_Staking.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Staking *StakingTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Staking *StakingSession) AcceptOwnership() (*types.Transaction, error) {
	return _Staking.Contract.AcceptOwnership(&_Staking.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Staking *StakingTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Staking.Contract.AcceptOwnership(&_Staking.TransactOpts)
}

// AddOperators is a paid mutator transaction binding the contract method 0xa07aea1c.
//
// Solidity: function addOperators(address[] operators) returns()
func (_Staking *StakingTransactor) AddOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "addOperators", operators)
}

// AddOperators is a paid mutator transaction binding the contract method 0xa07aea1c.
//
// Solidity: function addOperators(address[] operators) returns()
func (_Staking *StakingSession) AddOperators(operators []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.AddOperators(&_Staking.TransactOpts, operators)
}

// AddOperators is a paid mutator transaction binding the contract method 0xa07aea1c.
//
// Solidity: function addOperators(address[] operators) returns()
func (_Staking *StakingTransactorSession) AddOperators(operators []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.AddOperators(&_Staking.TransactOpts, operators)
}

// AddReward is a paid mutator transaction binding the contract method 0x74de4ec4.
//
// Solidity: function addReward(uint256 amount) returns()
func (_Staking *StakingTransactor) AddReward(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "addReward", amount)
}

// AddReward is a paid mutator transaction binding the contract method 0x74de4ec4.
//
// Solidity: function addReward(uint256 amount) returns()
func (_Staking *StakingSession) AddReward(amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddReward(&_Staking.TransactOpts, amount)
}

// AddReward is a paid mutator transaction binding the contract method 0x74de4ec4.
//
// Solidity: function addReward(uint256 amount) returns()
func (_Staking *StakingTransactorSession) AddReward(amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.AddReward(&_Staking.TransactOpts, amount)
}

// ChangeRewardRate is a paid mutator transaction binding the contract method 0x74f237c4.
//
// Solidity: function changeRewardRate(uint256 newRate) returns()
func (_Staking *StakingTransactor) ChangeRewardRate(opts *bind.TransactOpts, newRate *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "changeRewardRate", newRate)
}

// ChangeRewardRate is a paid mutator transaction binding the contract method 0x74f237c4.
//
// Solidity: function changeRewardRate(uint256 newRate) returns()
func (_Staking *StakingSession) ChangeRewardRate(newRate *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.ChangeRewardRate(&_Staking.TransactOpts, newRate)
}

// ChangeRewardRate is a paid mutator transaction binding the contract method 0x74f237c4.
//
// Solidity: function changeRewardRate(uint256 newRate) returns()
func (_Staking *StakingTransactorSession) ChangeRewardRate(newRate *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.ChangeRewardRate(&_Staking.TransactOpts, newRate)
}

// Conclude is a paid mutator transaction binding the contract method 0xe5f92973.
//
// Solidity: function conclude() returns()
func (_Staking *StakingTransactor) Conclude(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "conclude")
}

// Conclude is a paid mutator transaction binding the contract method 0xe5f92973.
//
// Solidity: function conclude() returns()
func (_Staking *StakingSession) Conclude() (*types.Transaction, error) {
	return _Staking.Contract.Conclude(&_Staking.TransactOpts)
}

// Conclude is a paid mutator transaction binding the contract method 0xe5f92973.
//
// Solidity: function conclude() returns()
func (_Staking *StakingTransactorSession) Conclude() (*types.Transaction, error) {
	return _Staking.Contract.Conclude(&_Staking.TransactOpts)
}

// EmergencyPause is a paid mutator transaction binding the contract method 0x51858e27.
//
// Solidity: function emergencyPause() returns()
func (_Staking *StakingTransactor) EmergencyPause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "emergencyPause")
}

// EmergencyPause is a paid mutator transaction binding the contract method 0x51858e27.
//
// Solidity: function emergencyPause() returns()
func (_Staking *StakingSession) EmergencyPause() (*types.Transaction, error) {
	return _Staking.Contract.EmergencyPause(&_Staking.TransactOpts)
}

// EmergencyPause is a paid mutator transaction binding the contract method 0x51858e27.
//
// Solidity: function emergencyPause() returns()
func (_Staking *StakingTransactorSession) EmergencyPause() (*types.Transaction, error) {
	return _Staking.Contract.EmergencyPause(&_Staking.TransactOpts)
}

// EmergencyUnpause is a paid mutator transaction binding the contract method 0x4a4e3bd5.
//
// Solidity: function emergencyUnpause() returns()
func (_Staking *StakingTransactor) EmergencyUnpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "emergencyUnpause")
}

// EmergencyUnpause is a paid mutator transaction binding the contract method 0x4a4e3bd5.
//
// Solidity: function emergencyUnpause() returns()
func (_Staking *StakingSession) EmergencyUnpause() (*types.Transaction, error) {
	return _Staking.Contract.EmergencyUnpause(&_Staking.TransactOpts)
}

// EmergencyUnpause is a paid mutator transaction binding the contract method 0x4a4e3bd5.
//
// Solidity: function emergencyUnpause() returns()
func (_Staking *StakingTransactorSession) EmergencyUnpause() (*types.Transaction, error) {
	return _Staking.Contract.EmergencyUnpause(&_Staking.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8932a90d.
//
// Solidity: function migrate(bytes data) returns()
func (_Staking *StakingTransactor) Migrate(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "migrate", data)
}

// Migrate is a paid mutator transaction binding the contract method 0x8932a90d.
//
// Solidity: function migrate(bytes data) returns()
func (_Staking *StakingSession) Migrate(data []byte) (*types.Transaction, error) {
	return _Staking.Contract.Migrate(&_Staking.TransactOpts, data)
}

// Migrate is a paid mutator transaction binding the contract method 0x8932a90d.
//
// Solidity: function migrate(bytes data) returns()
func (_Staking *StakingTransactorSession) Migrate(data []byte) (*types.Transaction, error) {
	return _Staking.Contract.Migrate(&_Staking.TransactOpts, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_Staking *StakingTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_Staking *StakingSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Staking.Contract.OnTokenTransfer(&_Staking.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_Staking *StakingTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Staking.Contract.OnTokenTransfer(&_Staking.TransactOpts, sender, amount, data)
}

// ProposeMigrationTarget is a paid mutator transaction binding the contract method 0x63b2c85a.
//
// Solidity: function proposeMigrationTarget(address migrationTarget) returns()
func (_Staking *StakingTransactor) ProposeMigrationTarget(opts *bind.TransactOpts, migrationTarget common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "proposeMigrationTarget", migrationTarget)
}

// ProposeMigrationTarget is a paid mutator transaction binding the contract method 0x63b2c85a.
//
// Solidity: function proposeMigrationTarget(address migrationTarget) returns()
func (_Staking *StakingSession) ProposeMigrationTarget(migrationTarget common.Address) (*types.Transaction, error) {
	return _Staking.Contract.ProposeMigrationTarget(&_Staking.TransactOpts, migrationTarget)
}

// ProposeMigrationTarget is a paid mutator transaction binding the contract method 0x63b2c85a.
//
// Solidity: function proposeMigrationTarget(address migrationTarget) returns()
func (_Staking *StakingTransactorSession) ProposeMigrationTarget(migrationTarget common.Address) (*types.Transaction, error) {
	return _Staking.Contract.ProposeMigrationTarget(&_Staking.TransactOpts, migrationTarget)
}

// RaiseAlert is a paid mutator transaction binding the contract method 0xda9c732f.
//
// Solidity: function raiseAlert() returns()
func (_Staking *StakingTransactor) RaiseAlert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "raiseAlert")
}

// RaiseAlert is a paid mutator transaction binding the contract method 0xda9c732f.
//
// Solidity: function raiseAlert() returns()
func (_Staking *StakingSession) RaiseAlert() (*types.Transaction, error) {
	return _Staking.Contract.RaiseAlert(&_Staking.TransactOpts)
}

// RaiseAlert is a paid mutator transaction binding the contract method 0xda9c732f.
//
// Solidity: function raiseAlert() returns()
func (_Staking *StakingTransactorSession) RaiseAlert() (*types.Transaction, error) {
	return _Staking.Contract.RaiseAlert(&_Staking.TransactOpts)
}

// RemoveOperators is a paid mutator transaction binding the contract method 0xd365a377.
//
// Solidity: function removeOperators(address[] operators) returns()
func (_Staking *StakingTransactor) RemoveOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "removeOperators", operators)
}

// RemoveOperators is a paid mutator transaction binding the contract method 0xd365a377.
//
// Solidity: function removeOperators(address[] operators) returns()
func (_Staking *StakingSession) RemoveOperators(operators []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.RemoveOperators(&_Staking.TransactOpts, operators)
}

// RemoveOperators is a paid mutator transaction binding the contract method 0xd365a377.
//
// Solidity: function removeOperators(address[] operators) returns()
func (_Staking *StakingTransactorSession) RemoveOperators(operators []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.RemoveOperators(&_Staking.TransactOpts, operators)
}

// SetFeedOperators is a paid mutator transaction binding the contract method 0xbfbd9b1b.
//
// Solidity: function setFeedOperators(address[] operators) returns()
func (_Staking *StakingTransactor) SetFeedOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setFeedOperators", operators)
}

// SetFeedOperators is a paid mutator transaction binding the contract method 0xbfbd9b1b.
//
// Solidity: function setFeedOperators(address[] operators) returns()
func (_Staking *StakingSession) SetFeedOperators(operators []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.SetFeedOperators(&_Staking.TransactOpts, operators)
}

// SetFeedOperators is a paid mutator transaction binding the contract method 0xbfbd9b1b.
//
// Solidity: function setFeedOperators(address[] operators) returns()
func (_Staking *StakingTransactorSession) SetFeedOperators(operators []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.SetFeedOperators(&_Staking.TransactOpts, operators)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 newMerkleRoot) returns()
func (_Staking *StakingTransactor) SetMerkleRoot(opts *bind.TransactOpts, newMerkleRoot [32]byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setMerkleRoot", newMerkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 newMerkleRoot) returns()
func (_Staking *StakingSession) SetMerkleRoot(newMerkleRoot [32]byte) (*types.Transaction, error) {
	return _Staking.Contract.SetMerkleRoot(&_Staking.TransactOpts, newMerkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 newMerkleRoot) returns()
func (_Staking *StakingTransactorSession) SetMerkleRoot(newMerkleRoot [32]byte) (*types.Transaction, error) {
	return _Staking.Contract.SetMerkleRoot(&_Staking.TransactOpts, newMerkleRoot)
}

// SetPoolConfig is a paid mutator transaction binding the contract method 0x8a44f337.
//
// Solidity: function setPoolConfig(uint256 maxPoolSize, uint256 maxCommunityStakeAmount, uint256 maxOperatorStakeAmount) returns()
func (_Staking *StakingTransactor) SetPoolConfig(opts *bind.TransactOpts, maxPoolSize *big.Int, maxCommunityStakeAmount *big.Int, maxOperatorStakeAmount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setPoolConfig", maxPoolSize, maxCommunityStakeAmount, maxOperatorStakeAmount)
}

// SetPoolConfig is a paid mutator transaction binding the contract method 0x8a44f337.
//
// Solidity: function setPoolConfig(uint256 maxPoolSize, uint256 maxCommunityStakeAmount, uint256 maxOperatorStakeAmount) returns()
func (_Staking *StakingSession) SetPoolConfig(maxPoolSize *big.Int, maxCommunityStakeAmount *big.Int, maxOperatorStakeAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.SetPoolConfig(&_Staking.TransactOpts, maxPoolSize, maxCommunityStakeAmount, maxOperatorStakeAmount)
}

// SetPoolConfig is a paid mutator transaction binding the contract method 0x8a44f337.
//
// Solidity: function setPoolConfig(uint256 maxPoolSize, uint256 maxCommunityStakeAmount, uint256 maxOperatorStakeAmount) returns()
func (_Staking *StakingTransactorSession) SetPoolConfig(maxPoolSize *big.Int, maxCommunityStakeAmount *big.Int, maxOperatorStakeAmount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.SetPoolConfig(&_Staking.TransactOpts, maxPoolSize, maxCommunityStakeAmount, maxOperatorStakeAmount)
}

// Start is a paid mutator transaction binding the contract method 0x8fb4b573.
//
// Solidity: function start(uint256 amount, uint256 initialRewardRate) returns()
func (_Staking *StakingTransactor) Start(opts *bind.TransactOpts, amount *big.Int, initialRewardRate *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "start", amount, initialRewardRate)
}

// Start is a paid mutator transaction binding the contract method 0x8fb4b573.
//
// Solidity: function start(uint256 amount, uint256 initialRewardRate) returns()
func (_Staking *StakingSession) Start(amount *big.Int, initialRewardRate *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Start(&_Staking.TransactOpts, amount, initialRewardRate)
}

// Start is a paid mutator transaction binding the contract method 0x8fb4b573.
//
// Solidity: function start(uint256 amount, uint256 initialRewardRate) returns()
func (_Staking *StakingTransactorSession) Start(amount *big.Int, initialRewardRate *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Start(&_Staking.TransactOpts, amount, initialRewardRate)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Staking *StakingTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Staking *StakingSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Staking.Contract.TransferOwnership(&_Staking.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Staking *StakingTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Staking.Contract.TransferOwnership(&_Staking.TransactOpts, to)
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Staking *StakingTransactor) Unstake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "unstake")
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Staking *StakingSession) Unstake() (*types.Transaction, error) {
	return _Staking.Contract.Unstake(&_Staking.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Staking *StakingTransactorSession) Unstake() (*types.Transaction, error) {
	return _Staking.Contract.Unstake(&_Staking.TransactOpts)
}

// WithdrawRemovedStake is a paid mutator transaction binding the contract method 0x5aa6e013.
//
// Solidity: function withdrawRemovedStake() returns()
func (_Staking *StakingTransactor) WithdrawRemovedStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdrawRemovedStake")
}

// WithdrawRemovedStake is a paid mutator transaction binding the contract method 0x5aa6e013.
//
// Solidity: function withdrawRemovedStake() returns()
func (_Staking *StakingSession) WithdrawRemovedStake() (*types.Transaction, error) {
	return _Staking.Contract.WithdrawRemovedStake(&_Staking.TransactOpts)
}

// WithdrawRemovedStake is a paid mutator transaction binding the contract method 0x5aa6e013.
//
// Solidity: function withdrawRemovedStake() returns()
func (_Staking *StakingTransactorSession) WithdrawRemovedStake() (*types.Transaction, error) {
	return _Staking.Contract.WithdrawRemovedStake(&_Staking.TransactOpts)
}

// WithdrawUnusedReward is a paid mutator transaction binding the contract method 0xebdb56f3.
//
// Solidity: function withdrawUnusedReward() returns()
func (_Staking *StakingTransactor) WithdrawUnusedReward(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdrawUnusedReward")
}

// WithdrawUnusedReward is a paid mutator transaction binding the contract method 0xebdb56f3.
//
// Solidity: function withdrawUnusedReward() returns()
func (_Staking *StakingSession) WithdrawUnusedReward() (*types.Transaction, error) {
	return _Staking.Contract.WithdrawUnusedReward(&_Staking.TransactOpts)
}

// WithdrawUnusedReward is a paid mutator transaction binding the contract method 0xebdb56f3.
//
// Solidity: function withdrawUnusedReward() returns()
func (_Staking *StakingTransactorSession) WithdrawUnusedReward() (*types.Transaction, error) {
	return _Staking.Contract.WithdrawUnusedReward(&_Staking.TransactOpts)
}

// StakingAlertRaisedIterator is returned from FilterAlertRaised and is used to iterate over the raw logs and unpacked data for AlertRaised events raised by the Staking contract.
type StakingAlertRaisedIterator struct {
	Event *StakingAlertRaised // Event containing the contract specifics and raw log

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
func (it *StakingAlertRaisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingAlertRaised)
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
		it.Event = new(StakingAlertRaised)
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
func (it *StakingAlertRaisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingAlertRaisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingAlertRaised represents a AlertRaised event raised by the Staking contract.
type StakingAlertRaised struct {
	Alerter      common.Address
	RoundId      *big.Int
	RewardAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAlertRaised is a free log retrieval operation binding the contract event 0xd2720e8f454493f612cc97499fe8cbce7fa4d4c18d346fe7104e9042df1c1edd.
//
// Solidity: event AlertRaised(address alerter, uint256 roundId, uint256 rewardAmount)
func (_Staking *StakingFilterer) FilterAlertRaised(opts *bind.FilterOpts) (*StakingAlertRaisedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "AlertRaised")
	if err != nil {
		return nil, err
	}
	return &StakingAlertRaisedIterator{contract: _Staking.contract, event: "AlertRaised", logs: logs, sub: sub}, nil
}

// WatchAlertRaised is a free log subscription operation binding the contract event 0xd2720e8f454493f612cc97499fe8cbce7fa4d4c18d346fe7104e9042df1c1edd.
//
// Solidity: event AlertRaised(address alerter, uint256 roundId, uint256 rewardAmount)
func (_Staking *StakingFilterer) WatchAlertRaised(opts *bind.WatchOpts, sink chan<- *StakingAlertRaised) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "AlertRaised")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingAlertRaised)
				if err := _Staking.contract.UnpackLog(event, "AlertRaised", log); err != nil {
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

// ParseAlertRaised is a log parse operation binding the contract event 0xd2720e8f454493f612cc97499fe8cbce7fa4d4c18d346fe7104e9042df1c1edd.
//
// Solidity: event AlertRaised(address alerter, uint256 roundId, uint256 rewardAmount)
func (_Staking *StakingFilterer) ParseAlertRaised(log types.Log) (*StakingAlertRaised, error) {
	event := new(StakingAlertRaised)
	if err := _Staking.contract.UnpackLog(event, "AlertRaised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingMerkleRootChangedIterator is returned from FilterMerkleRootChanged and is used to iterate over the raw logs and unpacked data for MerkleRootChanged events raised by the Staking contract.
type StakingMerkleRootChangedIterator struct {
	Event *StakingMerkleRootChanged // Event containing the contract specifics and raw log

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
func (it *StakingMerkleRootChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingMerkleRootChanged)
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
		it.Event = new(StakingMerkleRootChanged)
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
func (it *StakingMerkleRootChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingMerkleRootChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingMerkleRootChanged represents a MerkleRootChanged event raised by the Staking contract.
type StakingMerkleRootChanged struct {
	NewMerkleRoot [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootChanged is a free log retrieval operation binding the contract event 0x1b930366dfeaa7eb3b325021e4ae81e36527063452ee55b86c95f85b36f4c31c.
//
// Solidity: event MerkleRootChanged(bytes32 newMerkleRoot)
func (_Staking *StakingFilterer) FilterMerkleRootChanged(opts *bind.FilterOpts) (*StakingMerkleRootChangedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "MerkleRootChanged")
	if err != nil {
		return nil, err
	}
	return &StakingMerkleRootChangedIterator{contract: _Staking.contract, event: "MerkleRootChanged", logs: logs, sub: sub}, nil
}

// WatchMerkleRootChanged is a free log subscription operation binding the contract event 0x1b930366dfeaa7eb3b325021e4ae81e36527063452ee55b86c95f85b36f4c31c.
//
// Solidity: event MerkleRootChanged(bytes32 newMerkleRoot)
func (_Staking *StakingFilterer) WatchMerkleRootChanged(opts *bind.WatchOpts, sink chan<- *StakingMerkleRootChanged) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "MerkleRootChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingMerkleRootChanged)
				if err := _Staking.contract.UnpackLog(event, "MerkleRootChanged", log); err != nil {
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

// ParseMerkleRootChanged is a log parse operation binding the contract event 0x1b930366dfeaa7eb3b325021e4ae81e36527063452ee55b86c95f85b36f4c31c.
//
// Solidity: event MerkleRootChanged(bytes32 newMerkleRoot)
func (_Staking *StakingFilterer) ParseMerkleRootChanged(log types.Log) (*StakingMerkleRootChanged, error) {
	event := new(StakingMerkleRootChanged)
	if err := _Staking.contract.UnpackLog(event, "MerkleRootChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the Staking contract.
type StakingMigratedIterator struct {
	Event *StakingMigrated // Event containing the contract specifics and raw log

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
func (it *StakingMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingMigrated)
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
		it.Event = new(StakingMigrated)
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
func (it *StakingMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingMigrated represents a Migrated event raised by the Staking contract.
type StakingMigrated struct {
	Staker           common.Address
	Principal        *big.Int
	BaseReward       *big.Int
	DelegationReward *big.Int
	Data             []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0x667838b33bdc898470de09e0e746990f2adc11b965b7fe6828e502ebc39e0434.
//
// Solidity: event Migrated(address staker, uint256 principal, uint256 baseReward, uint256 delegationReward, bytes data)
func (_Staking *StakingFilterer) FilterMigrated(opts *bind.FilterOpts) (*StakingMigratedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Migrated")
	if err != nil {
		return nil, err
	}
	return &StakingMigratedIterator{contract: _Staking.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0x667838b33bdc898470de09e0e746990f2adc11b965b7fe6828e502ebc39e0434.
//
// Solidity: event Migrated(address staker, uint256 principal, uint256 baseReward, uint256 delegationReward, bytes data)
func (_Staking *StakingFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *StakingMigrated) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Migrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingMigrated)
				if err := _Staking.contract.UnpackLog(event, "Migrated", log); err != nil {
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

// ParseMigrated is a log parse operation binding the contract event 0x667838b33bdc898470de09e0e746990f2adc11b965b7fe6828e502ebc39e0434.
//
// Solidity: event Migrated(address staker, uint256 principal, uint256 baseReward, uint256 delegationReward, bytes data)
func (_Staking *StakingFilterer) ParseMigrated(log types.Log) (*StakingMigrated, error) {
	event := new(StakingMigrated)
	if err := _Staking.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingMigrationTargetAcceptedIterator is returned from FilterMigrationTargetAccepted and is used to iterate over the raw logs and unpacked data for MigrationTargetAccepted events raised by the Staking contract.
type StakingMigrationTargetAcceptedIterator struct {
	Event *StakingMigrationTargetAccepted // Event containing the contract specifics and raw log

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
func (it *StakingMigrationTargetAcceptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingMigrationTargetAccepted)
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
		it.Event = new(StakingMigrationTargetAccepted)
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
func (it *StakingMigrationTargetAcceptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingMigrationTargetAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingMigrationTargetAccepted represents a MigrationTargetAccepted event raised by the Staking contract.
type StakingMigrationTargetAccepted struct {
	MigrationTarget common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMigrationTargetAccepted is a free log retrieval operation binding the contract event 0xfa33c052bbee754f3c0482a89962daffe749191fa33c696a61e947fbfd68bd84.
//
// Solidity: event MigrationTargetAccepted(address migrationTarget)
func (_Staking *StakingFilterer) FilterMigrationTargetAccepted(opts *bind.FilterOpts) (*StakingMigrationTargetAcceptedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "MigrationTargetAccepted")
	if err != nil {
		return nil, err
	}
	return &StakingMigrationTargetAcceptedIterator{contract: _Staking.contract, event: "MigrationTargetAccepted", logs: logs, sub: sub}, nil
}

// WatchMigrationTargetAccepted is a free log subscription operation binding the contract event 0xfa33c052bbee754f3c0482a89962daffe749191fa33c696a61e947fbfd68bd84.
//
// Solidity: event MigrationTargetAccepted(address migrationTarget)
func (_Staking *StakingFilterer) WatchMigrationTargetAccepted(opts *bind.WatchOpts, sink chan<- *StakingMigrationTargetAccepted) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "MigrationTargetAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingMigrationTargetAccepted)
				if err := _Staking.contract.UnpackLog(event, "MigrationTargetAccepted", log); err != nil {
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

// ParseMigrationTargetAccepted is a log parse operation binding the contract event 0xfa33c052bbee754f3c0482a89962daffe749191fa33c696a61e947fbfd68bd84.
//
// Solidity: event MigrationTargetAccepted(address migrationTarget)
func (_Staking *StakingFilterer) ParseMigrationTargetAccepted(log types.Log) (*StakingMigrationTargetAccepted, error) {
	event := new(StakingMigrationTargetAccepted)
	if err := _Staking.contract.UnpackLog(event, "MigrationTargetAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingMigrationTargetProposedIterator is returned from FilterMigrationTargetProposed and is used to iterate over the raw logs and unpacked data for MigrationTargetProposed events raised by the Staking contract.
type StakingMigrationTargetProposedIterator struct {
	Event *StakingMigrationTargetProposed // Event containing the contract specifics and raw log

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
func (it *StakingMigrationTargetProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingMigrationTargetProposed)
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
		it.Event = new(StakingMigrationTargetProposed)
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
func (it *StakingMigrationTargetProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingMigrationTargetProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingMigrationTargetProposed represents a MigrationTargetProposed event raised by the Staking contract.
type StakingMigrationTargetProposed struct {
	MigrationTarget common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMigrationTargetProposed is a free log retrieval operation binding the contract event 0x5c74c441be501340b2713817a6c6975e6f3d4a4ae39fa1ac0bf75d3c54a0cad3.
//
// Solidity: event MigrationTargetProposed(address migrationTarget)
func (_Staking *StakingFilterer) FilterMigrationTargetProposed(opts *bind.FilterOpts) (*StakingMigrationTargetProposedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "MigrationTargetProposed")
	if err != nil {
		return nil, err
	}
	return &StakingMigrationTargetProposedIterator{contract: _Staking.contract, event: "MigrationTargetProposed", logs: logs, sub: sub}, nil
}

// WatchMigrationTargetProposed is a free log subscription operation binding the contract event 0x5c74c441be501340b2713817a6c6975e6f3d4a4ae39fa1ac0bf75d3c54a0cad3.
//
// Solidity: event MigrationTargetProposed(address migrationTarget)
func (_Staking *StakingFilterer) WatchMigrationTargetProposed(opts *bind.WatchOpts, sink chan<- *StakingMigrationTargetProposed) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "MigrationTargetProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingMigrationTargetProposed)
				if err := _Staking.contract.UnpackLog(event, "MigrationTargetProposed", log); err != nil {
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

// ParseMigrationTargetProposed is a log parse operation binding the contract event 0x5c74c441be501340b2713817a6c6975e6f3d4a4ae39fa1ac0bf75d3c54a0cad3.
//
// Solidity: event MigrationTargetProposed(address migrationTarget)
func (_Staking *StakingFilterer) ParseMigrationTargetProposed(log types.Log) (*StakingMigrationTargetProposed, error) {
	event := new(StakingMigrationTargetProposed)
	if err := _Staking.contract.UnpackLog(event, "MigrationTargetProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Staking contract.
type StakingOwnershipTransferRequestedIterator struct {
	Event *StakingOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *StakingOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOwnershipTransferRequested)
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
		it.Event = new(StakingOwnershipTransferRequested)
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
func (it *StakingOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Staking contract.
type StakingOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Staking *StakingFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingOwnershipTransferRequestedIterator{contract: _Staking.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Staking *StakingFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *StakingOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOwnershipTransferRequested)
				if err := _Staking.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_Staking *StakingFilterer) ParseOwnershipTransferRequested(log types.Log) (*StakingOwnershipTransferRequested, error) {
	event := new(StakingOwnershipTransferRequested)
	if err := _Staking.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Staking contract.
type StakingOwnershipTransferredIterator struct {
	Event *StakingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOwnershipTransferred)
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
		it.Event = new(StakingOwnershipTransferred)
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
func (it *StakingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOwnershipTransferred represents a OwnershipTransferred event raised by the Staking contract.
type StakingOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Staking *StakingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingOwnershipTransferredIterator{contract: _Staking.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Staking *StakingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOwnershipTransferred)
				if err := _Staking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Staking *StakingFilterer) ParseOwnershipTransferred(log types.Log) (*StakingOwnershipTransferred, error) {
	event := new(StakingOwnershipTransferred)
	if err := _Staking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Staking contract.
type StakingPausedIterator struct {
	Event *StakingPaused // Event containing the contract specifics and raw log

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
func (it *StakingPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPaused)
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
		it.Event = new(StakingPaused)
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
func (it *StakingPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPaused represents a Paused event raised by the Staking contract.
type StakingPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Staking *StakingFilterer) FilterPaused(opts *bind.FilterOpts) (*StakingPausedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &StakingPausedIterator{contract: _Staking.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Staking *StakingFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *StakingPaused) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPaused)
				if err := _Staking.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Staking *StakingFilterer) ParsePaused(log types.Log) (*StakingPaused, error) {
	event := new(StakingPaused)
	if err := _Staking.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Staking contract.
type StakingStakedIterator struct {
	Event *StakingStaked // Event containing the contract specifics and raw log

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
func (it *StakingStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingStaked)
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
		it.Event = new(StakingStaked)
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
func (it *StakingStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingStaked represents a Staked event raised by the Staking contract.
type StakingStaked struct {
	Staker     common.Address
	NewStake   *big.Int
	TotalStake *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address staker, uint256 newStake, uint256 totalStake)
func (_Staking *StakingFilterer) FilterStaked(opts *bind.FilterOpts) (*StakingStakedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Staked")
	if err != nil {
		return nil, err
	}
	return &StakingStakedIterator{contract: _Staking.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address staker, uint256 newStake, uint256 totalStake)
func (_Staking *StakingFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakingStaked) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Staked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingStaked)
				if err := _Staking.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address staker, uint256 newStake, uint256 totalStake)
func (_Staking *StakingFilterer) ParseStaked(log types.Log) (*StakingStaked, error) {
	event := new(StakingStaked)
	if err := _Staking.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Staking contract.
type StakingUnpausedIterator struct {
	Event *StakingUnpaused // Event containing the contract specifics and raw log

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
func (it *StakingUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUnpaused)
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
		it.Event = new(StakingUnpaused)
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
func (it *StakingUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUnpaused represents a Unpaused event raised by the Staking contract.
type StakingUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Staking *StakingFilterer) FilterUnpaused(opts *bind.FilterOpts) (*StakingUnpausedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &StakingUnpausedIterator{contract: _Staking.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Staking *StakingFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StakingUnpaused) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUnpaused)
				if err := _Staking.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Staking *StakingFilterer) ParseUnpaused(log types.Log) (*StakingUnpaused, error) {
	event := new(StakingUnpaused)
	if err := _Staking.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Staking contract.
type StakingUnstakedIterator struct {
	Event *StakingUnstaked // Event containing the contract specifics and raw log

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
func (it *StakingUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUnstaked)
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
		it.Event = new(StakingUnstaked)
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
func (it *StakingUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUnstaked represents a Unstaked event raised by the Staking contract.
type StakingUnstaked struct {
	Staker           common.Address
	Principal        *big.Int
	BaseReward       *big.Int
	DelegationReward *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address staker, uint256 principal, uint256 baseReward, uint256 delegationReward)
func (_Staking *StakingFilterer) FilterUnstaked(opts *bind.FilterOpts) (*StakingUnstakedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Unstaked")
	if err != nil {
		return nil, err
	}
	return &StakingUnstakedIterator{contract: _Staking.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address staker, uint256 principal, uint256 baseReward, uint256 delegationReward)
func (_Staking *StakingFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakingUnstaked) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Unstaked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUnstaked)
				if err := _Staking.contract.UnpackLog(event, "Unstaked", log); err != nil {
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

// ParseUnstaked is a log parse operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address staker, uint256 principal, uint256 baseReward, uint256 delegationReward)
func (_Staking *StakingFilterer) ParseUnstaked(log types.Log) (*StakingUnstaked, error) {
	event := new(StakingUnstaked)
	if err := _Staking.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
