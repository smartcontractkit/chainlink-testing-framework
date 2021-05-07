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

// AccessControlledAggregatorABI is the input ABI used to generate the binding from.
const AccessControlledAggregatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"_minSubmissionValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_maxSubmissionValue\",\"type\":\"int256\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AvailableFundsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdateRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"OracleAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"name\":\"OraclePermissionsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"delay\",\"type\":\"uint32\"}],\"name\":\"RequesterPermissionsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint128\",\"name\":\"paymentAmount\",\"type\":\"uint128\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"minSubmissionCount\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"maxSubmissionCount\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"restartDelay\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timeout\",\"type\":\"uint32\"}],\"name\":\"RoundDetailsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"submission\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"round\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"SubmissionReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"acceptAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allocatedFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableFunds\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_removed\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_added\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_addedAdmins\",\"type\":\"address[]\"},{\"internalType\":\"uint32\",\"name\":\"_minSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"}],\"name\":\"changeOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSubmissionCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSubmissionValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSubmissionCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSubmissionValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleCount\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_queriedRoundId\",\"type\":\"uint32\"}],\"name\":\"oracleRoundState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_eligibleToSubmit\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"_roundId\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"_latestSubmission\",\"type\":\"int256\"},{\"internalType\":\"uint64\",\"name\":\"_startedAt\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_timeout\",\"type\":\"uint64\"},{\"internalType\":\"uint128\",\"name\":\"_availableFunds\",\"type\":\"uint128\"},{\"internalType\":\"uint8\",\"name\":\"_oracleCount\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paymentAmount\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"restartDelay\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_authorized\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"_delay\",\"type\":\"uint32\"}],\"name\":\"setRequesterPermissions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_submission\",\"type\":\"int256\"}],\"name\":\"submit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeout\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateAvailableFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"_paymentAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"_minSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maxSubmissions\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_restartDelay\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_timeout\",\"type\":\"uint32\"}],\"name\":\"updateFutureRounds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"withdrawablePayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AccessControlledAggregatorBin is the compiled bytecode used for deploying new contracts.
var AccessControlledAggregatorBin = "0x60c06040523480156200001157600080fd5b5060405162008c3538038062008c3583398181016040526101008110156200003857600080fd5b8101908080519060200190929190805190602001909291908051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190805160405193929190846401000000008211156200009f57600080fd5b83820191506020820185811115620000b657600080fd5b8251866001820283011164010000000082111715620000d457600080fd5b8083526020830192505050908051906020019080838360005b838110156200010a578082015181840152602081019050620000ed565b50505050905090810190601f168015620001385780820380516001836020036101000a031916815260200191505b506040525050508787878787878787336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555087600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001df8760008060008a620002c760201b60201c565b620001f085620007fc60201b60201c565b83608081815250508260a0818152505081600560006101000a81548160ff021916908360ff16021790555080600690805190602001906200023392919062000b2b565b50620002548663ffffffff1642620009ba60201b62005e901790919060201c565b600960008063ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555050505050505050506001600e60006101000a81548160ff021916908315150217905550505050505050505062000bda565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146200038a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b60006200039c62000a4460201b60201c565b60ff1690508463ffffffff168463ffffffff16101562000424576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f6d6178206d75737420657175616c2f657863656564206d696e0000000000000081525060200191505060405180910390fd5b8363ffffffff168163ffffffff161015620004a7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f6d61782063616e6e6f742065786365656420746f74616c00000000000000000081525060200191505060405180910390fd5b60008163ffffffff161480620004c857508263ffffffff168163ffffffff16115b6200053b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f64656c61792063616e6e6f742065786365656420746f74616c0000000000000081525060200191505060405180910390fd5b6200055e866fffffffffffffffffffffffffffffffff1662000a5160201b60201c565b600d60000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16101562000607576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f696e73756666696369656e742066756e647320666f72207061796d656e74000081525060200191505060405180910390fd5b60006200061962000a4460201b60201c565b60ff161115620006a15760008563ffffffff1611620006a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6d696e206d7573742062652067726561746572207468616e203000000000000081525060200191505060405180910390fd5b5b85600460006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555084600460146101000a81548163ffffffff021916908363ffffffff16021790555083600460106101000a81548163ffffffff021916908363ffffffff16021790555082600460186101000a81548163ffffffff021916908363ffffffff160217905550816004601c6101000a81548163ffffffff021916908363ffffffff1602179055508363ffffffff168563ffffffff16600460009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff167f56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f8686604051808363ffffffff1663ffffffff1681526020018263ffffffff1663ffffffff1681526020019250505060405180910390a4505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614620008bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614620009b65781600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d60405160405180910390a35b5050565b60008282111562000a33576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b6000600c80549050905090565b600062000a99600262000a8562000a6d62000a4460201b60201c565b60ff168562000aa060201b6200704c1790919060201c565b62000aa060201b6200704c1790919060201c565b9050919050565b60008083141562000ab5576000905062000b25565b600082840290508284828162000ac757fe5b041462000b20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018062008c146021913960400191505060405180910390fd5b809150505b92915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1062000b6e57805160ff191683800117855562000b9f565b8280016001018555821562000b9f579182015b8281111562000b9e57825182559160200191906001019062000b81565b5b50905062000bae919062000bb2565b5090565b62000bd791905b8082111562000bd357600081600090555060010162000bb9565b5090565b90565b60805160a05161800c62000c08600039806114bf52806119285250806114295280613030525061800c6000f3fe608060405234801561001057600080fd5b50600436106102955760003560e01c806370dea79a11610167578063a4c0ed36116100ce578063d4cc54e411610087578063d4cc54e414610f6e578063dc7f012414610fb0578063e2e4031714610fd2578063e9ee6eeb1461102a578063f2fde38b1461108e578063feaf968c146110d257610295565b8063a4c0ed3614610d8d578063b5ab58dc14610e30578063b633620c14610e72578063c107532914610eb4578063c35905c614610f02578063c937450014610f4457610295565b80638823da6c116101205780638823da6c14610acc57806388aa80e714610b105780638da5cb5b14610c2f57806398e5b12a14610c795780639a6fc8f514610caf578063a118f24914610d4957610295565b806370dea79a146109cf5780637284e416146109f957806379ba509714610a7c5780637c2b0b2114610a865780638038e4a114610aa45780638205bf6a14610aae57610295565b806340884c521161020b57806358609e44116101c457806358609e44146107a8578063613d8fcc146107d2578063628806ef146107f657806364efb22b1461083a578063668a0f02146108be5780636b14daf8146108dc57610295565b806340884c521461067757806346fcff4c146106d65780634f8fc3b51461071857806350d25bcd1461072257806354fd4d501461074057806357970e931461075e57610295565b8063313ce5671161025d578063313ce5671461039e578063357ebb02146103c257806338aa4c72146103ec5780633969c20f1461046c5780633a5381b5146105bf5780633d3d77141461060957610295565b80630a7569831461029a5780631327d3d8146102a4578063202ee0ed146102e857806320ed02751461032057806323ca290314610380575b600080fd5b6102a261113c565b005b6102e6600480360360208110156102ba57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061125d565b005b61031e600480360360408110156102fe57600080fd5b810190808035906020019092919080359060200190929190505050611419565b005b61037e6004803603606081101561033657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190803563ffffffff16906020019092919050505061164d565b005b610388611926565b6040518082815260200191505060405180910390f35b6103a661194a565b604051808260ff1660ff16815260200191505060405180910390f35b6103ca61195d565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b61046a600480360360a081101561040257600080fd5b8101908080356fffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190505050611973565b005b6105bd600480360360c081101561048257600080fd5b810190808035906020019064010000000081111561049f57600080fd5b8201836020820111156104b157600080fd5b803590602001918460208302840111640100000000831117156104d357600080fd5b9091929391929390803590602001906401000000008111156104f457600080fd5b82018360208201111561050657600080fd5b8035906020019184602083028401116401000000008311171561052857600080fd5b90919293919293908035906020019064010000000081111561054957600080fd5b82018360208201111561055b57600080fd5b8035906020019184602083028401116401000000008311171561057d57600080fd5b9091929391929390803563ffffffff169060200190929190803563ffffffff169060200190929190803563ffffffff169060200190929190505050611e88565b005b6105c761216c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6106756004803603606081101561061f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050612192565b005b61067f6125b8565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156106c25780820151818401526020810190506106a7565b505050509050019250505060405180910390f35b6106de612646565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61072061266f565b005b61072a6128ab565b6040518082815260200191505060405180910390f35b61074861297b565b6040518082815260200191505060405180910390f35b610766612980565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6107b06129a6565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b6107da6129bc565b604051808260ff1660ff16815260200191505060405180910390f35b6108386004803603602081101561080c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506129c9565b005b61087c6004803603602081101561085057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612c2c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6108c6612c98565b6040518082815260200191505060405180910390f35b6109b5600480360360408110156108f257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561092f57600080fd5b82018360208201111561094157600080fd5b8035906020019184600183028401116401000000008311171561096357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050612d68565b604051808215151515815260200191505060405180910390f35b6109d7612db2565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610a01612dc8565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610a41578082015181840152602081019050610a26565b50505050905090810190601f168015610a6e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610a84612e66565b005b610a8e61302e565b6040518082815260200191505060405180910390f35b610aac613052565b005b610ab6613172565b6040518082815260200191505060405180910390f35b610b0e60048036036020811015610ae257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050613242565b005b610b6260048036036040811015610b2657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803563ffffffff169060200190929190505050613415565b60405180891515151581526020018863ffffffff1663ffffffff1681526020018781526020018667ffffffffffffffff1667ffffffffffffffff1681526020018567ffffffffffffffff1667ffffffffffffffff168152602001846fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1681526020018360ff1660ff168152602001826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1681526020019850505050505050505060405180910390f35b610c37613674565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610c81613699565b604051808269ffffffffffffffffffff1669ffffffffffffffffffff16815260200191505060405180910390f35b610ce760048036036020811015610cc557600080fd5b81019080803569ffffffffffffffffffff16906020019092919050505061386f565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b610d8b60048036036020811015610d5f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050613954565b005b610e2e60048036036060811015610da357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190640100000000811115610dea57600080fd5b820183602082011115610dfc57600080fd5b80359060200191846001830284011164010000000083111715610e1e57600080fd5b9091929391929390505050613b26565b005b610e5c60048036036020811015610e4657600080fd5b8101908080359060200190929190505050613bad565b6040518082815260200191505060405180910390f35b610e9e60048036036020811015610e8857600080fd5b8101908080359060200190929190505050613c80565b6040518082815260200191505060405180910390f35b610f0060048036036040811015610eca57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050613d53565b005b610f0a614070565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610f4c614092565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b610f766140a8565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610fb86140d1565b604051808215151515815260200191505060405180910390f35b61101460048036036020811015610fe857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506140e4565b6040518082815260200191505060405180910390f35b61108c6004803603604081101561104057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061415e565b005b6110d0600480360360208110156110a457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050614394565b005b6110da614515565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146111fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600e60009054906101000a900460ff161561125b576000600e60006101000a81548160ff0219169083151502179055507f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963860405160405180910390a15b565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461131f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146114155781600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d60405160405180910390a35b5050565b606061142533846145f7565b90507f00000000000000000000000000000000000000000000000000000000000000008212156114bd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f76616c75652062656c6f77206d696e5375626d697373696f6e56616c7565000081525060200191505060405180910390fd5b7f0000000000000000000000000000000000000000000000000000000000000000821315611553576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f76616c75652061626f7665206d61785375626d697373696f6e56616c7565000081525060200191505060405180910390fd5b600081511481906115ff576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b838110156115c45780820151818401526020810190506115a9565b50505050905090810190601f1680156115f15780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b506116098361496c565b6116138284614a86565b60008061161f85614c41565b9150915061162c85614e54565b611635856151b9565b81156116465761164585826152c1565b5b5050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461170f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b811515600b60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff161515141561176f57611921565b81156118355781600b60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a81548160ff02191690831515021790555080600b60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160016101000a81548163ffffffff021916908363ffffffff1602179055506118ba565b600b60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600080820160006101000a81549060ff02191690556000820160016101000a81549063ffffffff02191690556000820160056101000a81549063ffffffff021916905550505b8273ffffffffffffffffffffffffffffffffffffffff167fc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a838360405180831515151581526020018263ffffffff1663ffffffff1681526020019250505060405180910390a25b505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600560009054906101000a900460ff1681565b600460189054906101000a900463ffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611a35576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b6000611a3f6129bc565b60ff1690508463ffffffff168463ffffffff161015611ac6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f6d6178206d75737420657175616c2f657863656564206d696e0000000000000081525060200191505060405180910390fd5b8363ffffffff168163ffffffff161015611b48576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f6d61782063616e6e6f742065786365656420746f74616c00000000000000000081525060200191505060405180910390fd5b60008163ffffffff161480611b6857508263ffffffff168163ffffffff16115b611bda576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f64656c61792063616e6e6f742065786365656420746f74616c0000000000000081525060200191505060405180910390fd5b611bf5866fffffffffffffffffffffffffffffffff16615463565b600d60000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff161015611c9d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f696e73756666696369656e742066756e647320666f72207061796d656e74000081525060200191505060405180910390fd5b6000611ca76129bc565b60ff161115611d2d5760008563ffffffff1611611d2c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6d696e206d7573742062652067726561746572207468616e203000000000000081525060200191505060405180910390fd5b5b85600460006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555084600460146101000a81548163ffffffff021916908363ffffffff16021790555083600460106101000a81548163ffffffff021916908363ffffffff16021790555082600460186101000a81548163ffffffff021916908363ffffffff160217905550816004601c6101000a81548163ffffffff021916908363ffffffff1602179055508363ffffffff168563ffffffff16600460009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff167f56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f8686604051808363ffffffff1663ffffffff1681526020018263ffffffff1663ffffffff1681526020019250505060405180910390a4505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611f4a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b60008090505b89899050811015611f9957611f8c8a8a83818110611f6a57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1661549c565b8080600101915050611f50565b50848490508787905014612015576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f6e6565642073616d65206f7261636c6520616e642061646d696e20636f756e7481525060200191505060405180910390fd5b604d612037888890506120266129bc565b60ff166157e790919063ffffffff16565b11156120ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f6d6178206f7261636c657320616c6c6f7765640000000000000000000000000081525060200191505060405180910390fd5b60008090505b87879050811015612123576121168888838181106120cb57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff168787848181106120f457fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1661586f565b80806001019150506120b1565b50612161600460009054906101000a90046fffffffffffffffffffffffffffffffff168484846004601c9054906101000a900463ffffffff16611973565b505050505050505050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff16600860008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614612295576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6f6e6c792063616c6c61626c652062792061646d696e0000000000000000000081525060200191505060405180910390fd5b60008190506000600860008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff169050816fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff161015612397576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f696e73756666696369656e7420776974686472617761626c652066756e64730081525060200191505060405180910390fd5b6123bc82826fffffffffffffffffffffffffffffffff16615de390919063ffffffff16565b600860008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555061247a82600d60000160109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16615de390919063ffffffff16565b600d60000160106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85846fffffffffffffffffffffffffffffffff166040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561257057600080fd5b505af1158015612584573d6000803e3d6000fd5b505050506040513d602081101561259a57600080fd5b81019080805190602001909291905050506125b157fe5b5050505050565b6060600c80548060200260200160405190810160405280929190818152602001828054801561263c57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116125f2575b5050505050905090565b6000600d60000160009054906101000a90046fffffffffffffffffffffffffffffffff16905090565b612677617e2e565b600d6040518060400160405290816000820160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1681526020016000820160109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815250509050600061281e82602001516fffffffffffffffffffffffffffffffff16600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156127d557600080fd5b505afa1580156127e9573d6000803e3d6000fd5b505050506040513d60208110156127ff57600080fd5b8101908080519060200190929190505050615e9090919063ffffffff16565b90508082600001516fffffffffffffffffffffffffffffffff16146128a75780600d60000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550807ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f60405160405180910390a25b5050565b60006128fc336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b61296e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b612976615f19565b905090565b600381565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600460109054906101000a900463ffffffff1681565b6000600c80549050905090565b3373ffffffffffffffffffffffffffffffffffffffff16600860008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614612acc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f6f6e6c792063616c6c61626c652062792070656e64696e672061646d696e000081525060200191505060405180910390fd5b6000600860008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600860008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a350565b6000600860008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6000612ce9336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b612d5b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b612d63615f55565b905090565b6000612d748383615f75565b80612daa57503273ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16145b905092915050565b6004601c9054906101000a900463ffffffff1681565b60068054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015612e5e5780601f10612e3357610100808354040283529160200191612e5e565b820191906000526020600020905b815481529060010190602001808311612e4157829003601f168201915b505050505081565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612f29576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4d7573742062652070726f706f736564206f776e65720000000000000000000081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a350565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613114576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600e60009054906101000a900460ff16613170576001600e60006101000a81548160ff0219169083151502179055507faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348060405160405180910390a15b565b60006131c3336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b613235576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b61323d615fe4565b905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613304576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600f60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615613412576000600f60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d181604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15b50565b6000806000806000806000803273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146134c2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6f66662d636861696e2072656164696e67206f6e6c790000000000000000000081525060200191505060405180910390fd5b60008963ffffffff16111561364d576000600960008b63ffffffff1663ffffffff16815260200190815260200160002090506000600a60008c63ffffffff1663ffffffff16815260200190815260200160002090506135218c8c61603e565b8b600860008f73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101548460010160009054906101000a900467ffffffffffffffff168460010160089054906101000a900463ffffffff16600d60000160009054906101000a90046fffffffffffffffffffffffffffffffff166135bd6129bc565b60008960010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff161161360b57600460009054906101000a90046fffffffffffffffffffffffffffffffff1661362d565b87600101600c9054906101000a90046fffffffffffffffffffffffffffffffff165b8363ffffffff169350995099509950995099509950995099505050613667565b6136568a6160d7565b975097509750975097509750975097505b9295985092959890939650565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600b60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff1661375d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f6e6f7420617574686f72697a656420726571756573746572000000000000000081525060200191505060405180910390fd5b6000600760009054906101000a900463ffffffff1690506000600960008363ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff1611806137c857506137c781616388565b5b61383a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f7072657620726f756e64206d75737420626520737570657273656461626c650081525060200191505060405180910390fd5b600061385660018363ffffffff1661645b90919063ffffffff16565b9050613861816164ef565b8063ffffffff169250505090565b60008060008060006138c6336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b613938576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b613941866166ad565b9450945094509450945091939590929450565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613a16576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600f60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16613b23576001600f60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db481604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15b50565b60008282905014613b9f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f7472616e7366657220646f65736e2774206163636570742063616c6c6461746181525060200191505060405180910390fd5b613ba761266f565b50505050565b6000613bfe336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b613c70576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b613c79826168d1565b9050919050565b6000613cd1336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b613d43576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b613d4c82616915565b9050919050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613e15576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b6000600d60000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16905081613e99613e8a600460009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16615463565b83615e9090919063ffffffff16565b1015613f0d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f696e73756666696369656e7420726573657276652066756e647300000000000081525060200191505060405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015613fb657600080fd5b505af1158015613fca573d6000803e3d6000fd5b505050506040513d6020811015613fe057600080fd5b8101908080519060200190929190505050614063576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f746f6b656e207472616e73666572206661696c6564000000000000000000000081525060200191505060405180910390fd5b61406b61266f565b505050565b600460009054906101000a90046fffffffffffffffffffffffffffffffff1681565b600460149054906101000a900463ffffffff1681565b6000600d60000160109054906101000a90046fffffffffffffffffffffffffffffffff16905090565b600e60009054906101000a900460ff1681565b6000600860008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169050919050565b3373ffffffffffffffffffffffffffffffffffffffff16600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614614261576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6f6e6c792063616c6c61626c652062792061646d696e0000000000000000000081525060200191505060405180910390fd5b80600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff167fb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b81043383604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a25050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614614456576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b600080600080600061456c336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612d68565b6145de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6145e6616977565b945094509450945094509091929394565b60606000600860008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160109054906101000a900463ffffffff1690506000600760009054906101000a900463ffffffff16905060008263ffffffff1614156146b5576040518060400160405280601281526020017f6e6f7420656e61626c6564206f7261636c65000000000000000000000000000081525092505050614966565b8363ffffffff168263ffffffff161115614708576040518060400160405280601681526020017f6e6f742079657420656e61626c6564206f7261636c650000000000000000000081525092505050614966565b8363ffffffff16600860008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff1610156147ad576040518060400160405280601881526020017f6e6f206c6f6e67657220616c6c6f776564206f7261636c65000000000000000081525092505050614966565b8363ffffffff16600860008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160189054906101000a900463ffffffff1663ffffffff1610614851576040518060400160405280602081526020017f63616e6e6f74207265706f7274206f6e2070726576696f757320726f756e647381525092505050614966565b8063ffffffff168463ffffffff1614158015614892575061488260018263ffffffff1661645b90919063ffffffff16565b63ffffffff168463ffffffff1614155b80156148a557506148a384826169b1565b155b156148e9576040518060400160405280601781526020017f696e76616c696420726f756e6420746f207265706f727400000000000000000081525092505050614966565b60018463ffffffff161415801561491f575061491d61491860018663ffffffff16616a2f90919063ffffffff16565b616ac4565b155b15614963576040518060400160405280601f81526020017f70726576696f757320726f756e64206e6f7420737570657273656461626c650081525092505050614966565b50505b92915050565b61497581616b20565b61497e57614a83565b6000600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff169050600460189054906101000a900463ffffffff1663ffffffff1681018263ffffffff1611158015614a0c575060008114155b15614a175750614a83565b614a2082616b63565b81600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c6101000a81548163ffffffff021916908363ffffffff160217905550505b50565b614a8f81616e6c565b614b01576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f726f756e64206e6f7420616363657074696e67207375626d697373696f6e730081525060200191505060405180910390fd5b600a60008263ffffffff1663ffffffff16815260200190815260200160002060000182908060018154018082558091505060019003906000526020600020016000909190919091505580600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160186101000a81548163ffffffff021916908363ffffffff16021790555081600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101819055503373ffffffffffffffffffffffffffffffffffffffff168163ffffffff16837f92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c60405160405180910390a45050565b600080600a60008463ffffffff1663ffffffff16815260200190815260200160002060010160049054906101000a900463ffffffff1663ffffffff16600a60008563ffffffff1663ffffffff168152602001908152602001600020600001805490501015614cb85760008080905091509150614e4f565b6000614d33600a60008663ffffffff1663ffffffff168152602001908152602001600020600001805480602002602001604051908101604052809291908181526020018280548015614d2957602002820191906000526020600020905b815481526020019060010190808311614d15575b5050505050616eb1565b905080600960008663ffffffff1663ffffffff1681526020019081526020016000206000018190555042600960008663ffffffff1663ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083600960008663ffffffff1663ffffffff16815260200190815260200160002060010160106101000a81548163ffffffff021916908363ffffffff16021790555083600760046101000a81548163ffffffff021916908363ffffffff1602179055508363ffffffff16817f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a360018192509250505b915091565b6000600a60008363ffffffff1663ffffffff168152602001908152602001600020600101600c9054906101000a90046fffffffffffffffffffffffffffffffff169050614e9f617e2e565b600d6040518060400160405290816000820160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1681526020016000820160109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815250509050614f6b8282600001516fffffffffffffffffffffffffffffffff16615de390919063ffffffff16565b81600001906fffffffffffffffffffffffffffffffff1690816fffffffffffffffffffffffffffffffff1681525050614fc38282602001516fffffffffffffffffffffffffffffffff16616fa090919063ffffffff16565b81602001906fffffffffffffffffffffffffffffffff1690816fffffffffffffffffffffffffffffffff168152505080600d60008201518160000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555060208201518160000160106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055509050506150f982600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16616fa090919063ffffffff16565b600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555080600001516fffffffffffffffffffffffffffffffff167ffe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f60405160405180910390a2505050565b600a60008263ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1663ffffffff16600a60008363ffffffff1663ffffffff168152602001908152602001600020600001805490501015615223576152be565b600a60008263ffffffff1663ffffffff168152602001908152602001600020600080820160006152539190617e6c565b6001820160006101000a81549063ffffffff02191690556001820160046101000a81549063ffffffff02191690556001820160086101000a81549063ffffffff021916905560018201600c6101000a8154906fffffffffffffffffffffffffffffffff021916905550505b50565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415615323575061545f565b600061533f60018563ffffffff16616a2f90919063ffffffff16565b90506000600960008363ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff1690506000600960008463ffffffff1663ffffffff1681526020019081526020016000206000015490508373ffffffffffffffffffffffffffffffffffffffff1663beed9b51620186a084848a8a6040518663ffffffff1660e01b8152600401808563ffffffff1681526020018481526020018363ffffffff168152602001828152602001945050505050602060405180830381600088803b15801561541a57600080fd5b5087f19350505050801561544f57506040513d602081101561543b57600080fd5b810190808051906020019092919050505060015b6154585761545a565b505b505050505b5050565b600061549560026154876154756129bc565b60ff168561704c90919063ffffffff16565b61704c90919063ffffffff16565b9050919050565b6154a5816170d2565b615517576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6f7261636c65206e6f7420656e61626c6564000000000000000000000000000081525060200191505060405180910390fd5b6155436001600760009054906101000a900463ffffffff1663ffffffff1661645b90919063ffffffff16565b600860008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff1602179055506000600c6155c560016155b46129bc565b60ff16615e9090919063ffffffff16565b815481106155cf57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506000600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900461ffff16905080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81548161ffff021916908361ffff160217905550600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81549061ffff021916905581600c8261ffff168154811061571357fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600c80548061576657fe5b6001900381819060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690559055600015158373ffffffffffffffffffffffffffffffffffffffff167f18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e60405160405180910390a3505050565b600080828401905083811015615865576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b615878826170d2565b156158eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6f7261636c6520616c726561647920656e61626c65640000000000000000000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561598e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f63616e6e6f74207365742061646d696e20746f2030000000000000000000000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161480615ab857508073ffffffffffffffffffffffffffffffffffffffff16600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b615b2a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f6f776e65722063616e6e6f74206f76657277726974652061646d696e0000000081525060200191505060405180910390fd5b615b338261713c565b600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160106101000a81548163ffffffff021916908363ffffffff16021790555063ffffffff600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff160217905550600c80549050600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160006101000a81548161ffff021916908361ffff160217905550600c829080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600115158273ffffffffffffffffffffffffffffffffffffffff167f18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e60405160405180910390a38073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe90460405160405180910390a35050565b6000826fffffffffffffffffffffffffffffffff16826fffffffffffffffffffffffffffffffff161115615e7f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b600082821115615f08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b600060096000600760049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060000154905090565b6000600760049054906101000a900463ffffffff1663ffffffff16905090565b6000600f60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1680615fdc5750600e60009054906101000a900460ff16155b905092915050565b600060096000600760049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16905090565b600080600960008463ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff1611156160ae5761609182616e6c565b80156160a7575060006160a484846145f7565b51145b90506160d1565b6160b883836171f9565b80156160ce575060006160cb84846145f7565b51145b90505b92915050565b6000806000806000806000806000600960008063ffffffff16815260200190815260200160002090506000600860008c73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002090506000600760009054906101000a900463ffffffff1663ffffffff168260000160189054906101000a900463ffffffff1663ffffffff16148061619d575061619b600760009054906101000a900463ffffffff16616e6c565b155b90506161ba600760009054906101000a900463ffffffff16616ac4565b80156161c35750805b15616249576161f46001600760009054906101000a900463ffffffff1663ffffffff1661645b90919063ffffffff16565b9950600960008b63ffffffff1663ffffffff1681526020019081526020016000209250600460009054906101000a90046fffffffffffffffffffffffffffffffff1693506162428c8b6171f9565b9a506162cc565b600760009054906101000a900463ffffffff169950600960008b63ffffffff1663ffffffff1681526020019081526020016000209250600a60008b63ffffffff1663ffffffff168152602001908152602001600020600101600c9054906101000a90046fffffffffffffffffffffffffffffffff1693506162c98a616e6c565b9a505b60006162d88d8c6145f7565b51146162e35760009a505b8a8a83600101548560010160009054906101000a900467ffffffffffffffff16600a60008f63ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900463ffffffff16600d60000160009054906101000a90046fffffffffffffffffffffffffffffffff166163606129bc565b8a8363ffffffff1693509a509a509a509a509a509a509a509a50505050919395975091939597565b600080600960008463ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1690506000600a60008563ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900463ffffffff16905060008267ffffffffffffffff1611801561641b575060008163ffffffff16115b80156164525750426164468263ffffffff168467ffffffffffffffff1661728e90919063ffffffff16565b67ffffffffffffffff16105b92505050919050565b60008082840190508363ffffffff168163ffffffff1610156164e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b6164f881616b20565b616501576166aa565b6000600b60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160059054906101000a900463ffffffff1663ffffffff169050600b60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160019054906101000a900463ffffffff1663ffffffff1681018263ffffffff1611806165cc5750600081145b61663e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f6d7573742064656c61792072657175657374730000000000000000000000000081525060200191505060405180910390fd5b61664782616b63565b81600b60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160056101000a81548163ffffffff021916908363ffffffff160217905550505b50565b60008060008060006166bd617e8d565b600960008863ffffffff1663ffffffff168152602001908152602001600020604051806080016040529081600082015481526020016001820160009054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff1681526020016001820160089054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff1681526020016001820160109054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000816060015163ffffffff161180156167a857506167a78769ffffffffffffffffffff1661732a565b5b6040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090616884576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561684957808201518184015260208101905061682e565b50505050905090810190601f1680156168765780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b508681600001518260200151836040015184606001518267ffffffffffffffff1692508167ffffffffffffffff1691508063ffffffff169050955095509550955095505091939590929450565b60006168dc8261732a565b1561690b57600960008363ffffffff1663ffffffff168152602001908152602001600020600001549050616910565b600090505b919050565b60006169208261732a565b1561696d57600960008363ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff169050616972565b600090505b919050565b60008060008060006169a0600760049054906101000a900463ffffffff1663ffffffff1661386f565b945094509450945094509091929394565b60008163ffffffff166169d460018563ffffffff1661645b90919063ffffffff16565b63ffffffff16148015616a2757506000600960008463ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff16145b905092915050565b60008263ffffffff168263ffffffff161115616ab3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b600080600960008463ffffffff1663ffffffff16815260200190815260200160002060010160089054906101000a900467ffffffffffffffff1667ffffffffffffffff161180616b195750616b1882616388565b5b9050919050565b6000616b4e6001600760009054906101000a900463ffffffff1663ffffffff1661645b90919063ffffffff16565b63ffffffff168263ffffffff16149050919050565b616b85616b8060018363ffffffff16616a2f90919063ffffffff16565b61733d565b80600760006101000a81548163ffffffff021916908363ffffffff160217905550616bae617ecf565b6040518060a00160405280600067ffffffffffffffff81118015616bd157600080fd5b50604051908082528060200260200182016040528015616c005781602001602082028036833780820191505090505b508152602001600460109054906101000a900463ffffffff1663ffffffff168152602001600460149054906101000a900463ffffffff1663ffffffff1681526020016004601c9054906101000a900463ffffffff1663ffffffff168152602001600460009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815250905080600a60008463ffffffff1663ffffffff1681526020019081526020016000206000820151816000019080519060200190616cd2929190617f22565b5060208201518160010160006101000a81548163ffffffff021916908363ffffffff16021790555060408201518160010160046101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010160086101000a81548163ffffffff021916908363ffffffff160217905550608082015181600101600c6101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555090505042600960008463ffffffff1663ffffffff16815260200190815260200160002060010160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168263ffffffff167f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271600960008663ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff16604051808267ffffffffffffffff16815260200191505060405180910390a35050565b600080600a60008463ffffffff1663ffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1663ffffffff1614159050919050565b60008151600010616f2a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6c697374206d757374206e6f7420626520656d7074790000000000000000000081525060200191505060405180910390fd5b600082519050600060028281616f3c57fe5b049050600060028381616f4b57fe5b061415616f8657600080616f69866000600187036001870387617511565b8092508193505050616f7b82826175fe565b945050505050616f9b565b616f96846000600185038461769b565b925050505b919050565b6000808284019050836fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff161015617042576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b60008083141561705f57600090506170cc565b600082840290508284828161707057fe5b04146170c7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180617fb66021913960400191505060405180910390fd5b809150505b92915050565b600063ffffffff8016600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff16149050919050565b600080600760009054906101000a900463ffffffff16905060008163ffffffff16141580156171c85750600860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900463ffffffff1663ffffffff168163ffffffff16145b156171d657809150506171f4565b6171f060018263ffffffff1661645b90919063ffffffff16565b9150505b919050565b600080600860008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff169050600460189054906101000a900463ffffffff1663ffffffff1681018363ffffffff1611806172855750600081145b91505092915050565b60008082840190508367ffffffffffffffff168167ffffffffffffffff161015617320576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b600063ffffffff80168211159050919050565b61734681616388565b61734f5761750e565b600061736b60018363ffffffff16616a2f90919063ffffffff16565b9050600960008263ffffffff1663ffffffff16815260200190815260200160002060000154600960008463ffffffff1663ffffffff16815260200190815260200160002060000181905550600960008263ffffffff1663ffffffff16815260200190815260200160002060010160109054906101000a900463ffffffff16600960008463ffffffff1663ffffffff16815260200190815260200160002060010160106101000a81548163ffffffff021916908363ffffffff16021790555042600960008463ffffffff1663ffffffff16815260200190815260200160002060010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550600a60008363ffffffff1663ffffffff168152602001908152602001600020600080820160006174a29190617e6c565b6001820160006101000a81549063ffffffff02191690556001820160046101000a81549063ffffffff02191690556001820160086101000a81549063ffffffff021916905560018201600c6101000a8154906fffffffffffffffffffffffffffffffff02191690555050505b50565b60008082841061752057600080fd5b8386111580156175305750848411155b61753957600080fd5b8286111580156175495750848311155b61755257600080fd5b5b6001156175f3576007868603101561757b576175728787878787617735565b915091506175f4565b6000617588888888617ca9565b9050808411617599578095506175ed565b848110156175ac576001810196506175ec565b8085111580156175bb57508381105b6175c157fe5b6175cd8888838861769b565b92506175de8860018301888761769b565b9150828292509250506175f4565b5b50617553565b5b9550959350505050565b6000808312801561760f5750600082135b8061762657506000831380156176255750600082125b5b156176465760026176378484617da0565b8161763e57fe5b059050617695565b6000600280848161765357fe5b076002868161765e57fe5b07018161766757fe5b05905061769161768b6002868161767a57fe5b056002868161768557fe5b05617da0565b82617da0565b9150505b92915050565b6000818411156176aa57600080fd5b828211156176b757600080fd5b5b8284101561771657600784840310156176eb5760006176da8686868687617735565b80925081935050508191505061772d565b60006176f8868686617ca9565b905080831161770957809350617710565b6001810194505b506176b8565b84848151811061772257fe5b602002602001015190505b949350505050565b60008060008660018701039050600088600089018151811061775357fe5b6020026020010151905060008260011061778d577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6177a5565b8960018a018151811061779c57fe5b60200260200101515b90506000836002106177d7577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6177ef565b8a60028b01815181106177e657fe5b60200260200101515b9050600084600310617821577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff617839565b8b60038c018151811061783057fe5b60200260200101515b905060008560041061786b577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff617883565b8c60048d018151811061787a57fe5b60200260200101515b90506000866005106178b5577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6178cd565b8d60058e01815181106178c457fe5b60200260200101515b90506000876006106178ff577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff617917565b8e60068f018151811061790e57fe5b60200260200101515b90508587131561792c57858780975081985050505b8385131561793f57838580955081965050505b8183131561795257818380935081945050505b8487131561796557848780965081985050505b8386131561797857838680955081975050505b8083131561798b57808380925081945050505b8486131561799e57848680965081975050505b808213156179b157808280925081935050505b828713156179c457828780945081985050505b818613156179d757818680935081975050505b808513156179ea57808580925081965050505b828613156179fd57828680945081975050505b80841315617a1057808480925081955050505b82851315617a2357828580945081965050505b81841315617a3657818480935081955050505b82841315617a4957828480945081955050505b60008e8d0390506000811415617a6157879a50617b3b565b6001811415617a7257869a50617b3a565b6002811415617a8357859a50617b39565b6003811415617a9457849a50617b38565b6004811415617aa557839a50617b37565b6005811415617ab657829a50617b36565b6006811415617ac757819a50617b35565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f6b31206f7574206f6620626f756e64730000000000000000000000000000000081525060200191505060405180910390fd5b5b5b5b5b5b5b60008f8d0390508c8e1415617b5f578b8c9b509b5050505050505050505050617c9f565b6000811415617b7d578b899b509b5050505050505050505050617c9f565b6001811415617b9b578b889b509b5050505050505050505050617c9f565b6002811415617bb9578b879b509b5050505050505050505050617c9f565b6003811415617bd7578b869b509b5050505050505050505050617c9f565b6004811415617bf5578b859b509b5050505050505050505050617c9f565b6005811415617c13578b849b509b5050505050505050505050617c9f565b6006811415617c31578b839b509b5050505050505050505050617c9f565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f6b32206f7574206f6620626f756e64730000000000000000000000000000000081525060200191505060405180910390fd5b9550959350505050565b60008084600284860181617cb957fe5b0481518110617cc457fe5b602002602001015190506001840393506001830192505b600115617d97575b60018401935080858581518110617cf657fe5b602002602001015112617ce3575b60018303925080858481518110617d1757fe5b602002602001015113617d045782841015617d8957848381518110617d3857fe5b6020026020010151858581518110617d4c57fe5b6020026020010151868681518110617d6057fe5b60200260200101878681518110617d7357fe5b6020026020010182815250828152505050617d92565b82915050617d99565b617cdb565b505b9392505050565b600080828401905060008312158015617db95750838112155b80617dcf5750600083128015617dce57508381125b5b617e24576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180617f956021913960400191505060405180910390fd5b8091505092915050565b604051806040016040528060006fffffffffffffffffffffffffffffffff16815260200160006fffffffffffffffffffffffffffffffff1681525090565b5080546000825590600052602060002090810190617e8a9190617f6f565b50565b604051806080016040528060008152602001600067ffffffffffffffff168152602001600067ffffffffffffffff168152602001600063ffffffff1681525090565b6040518060a0016040528060608152602001600063ffffffff168152602001600063ffffffff168152602001600063ffffffff16815260200160006fffffffffffffffffffffffffffffffff1681525090565b828054828255906000526020600020908101928215617f5e579160200282015b82811115617f5d578251825591602001919060010190617f42565b5b509050617f6b9190617f6f565b5090565b617f9191905b80821115617f8d576000816000905550600101617f75565b5090565b9056fe5369676e6564536166654d6174683a206164646974696f6e206f766572666c6f77536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a2646970667358221220f27ed4cb477a87bbb608a4ff82cda9d3930fac5b4290b806c6936cd58dccb01664736f6c63430006060033536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77"

// DeployAccessControlledAggregator deploys a new Ethereum contract, binding an instance of AccessControlledAggregator to it.
func DeployAccessControlledAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _paymentAmount *big.Int, _timeout uint32, _validator common.Address, _minSubmissionValue *big.Int, _maxSubmissionValue *big.Int, _decimals uint8, _description string) (common.Address, *types.Transaction, *AccessControlledAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccessControlledAggregatorBin), backend, _link, _paymentAmount, _timeout, _validator, _minSubmissionValue, _maxSubmissionValue, _decimals, _description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessControlledAggregator{AccessControlledAggregatorCaller: AccessControlledAggregatorCaller{contract: contract}, AccessControlledAggregatorTransactor: AccessControlledAggregatorTransactor{contract: contract}, AccessControlledAggregatorFilterer: AccessControlledAggregatorFilterer{contract: contract}}, nil
}

// AccessControlledAggregator is an auto generated Go binding around an Ethereum contract.
type AccessControlledAggregator struct {
	AccessControlledAggregatorCaller     // Read-only binding to the contract
	AccessControlledAggregatorTransactor // Write-only binding to the contract
	AccessControlledAggregatorFilterer   // Log filterer for contract events
}

// AccessControlledAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessControlledAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlledAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessControlledAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlledAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessControlledAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlledAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessControlledAggregatorSession struct {
	Contract     *AccessControlledAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// AccessControlledAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessControlledAggregatorCallerSession struct {
	Contract *AccessControlledAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// AccessControlledAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessControlledAggregatorTransactorSession struct {
	Contract     *AccessControlledAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// AccessControlledAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessControlledAggregatorRaw struct {
	Contract *AccessControlledAggregator // Generic contract binding to access the raw methods on
}

// AccessControlledAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessControlledAggregatorCallerRaw struct {
	Contract *AccessControlledAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// AccessControlledAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessControlledAggregatorTransactorRaw struct {
	Contract *AccessControlledAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessControlledAggregator creates a new instance of AccessControlledAggregator, bound to a specific deployed contract.
func NewAccessControlledAggregator(address common.Address, backend bind.ContractBackend) (*AccessControlledAggregator, error) {
	contract, err := bindAccessControlledAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregator{AccessControlledAggregatorCaller: AccessControlledAggregatorCaller{contract: contract}, AccessControlledAggregatorTransactor: AccessControlledAggregatorTransactor{contract: contract}, AccessControlledAggregatorFilterer: AccessControlledAggregatorFilterer{contract: contract}}, nil
}

// NewAccessControlledAggregatorCaller creates a new read-only instance of AccessControlledAggregator, bound to a specific deployed contract.
func NewAccessControlledAggregatorCaller(address common.Address, caller bind.ContractCaller) (*AccessControlledAggregatorCaller, error) {
	contract, err := bindAccessControlledAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorCaller{contract: contract}, nil
}

// NewAccessControlledAggregatorTransactor creates a new write-only instance of AccessControlledAggregator, bound to a specific deployed contract.
func NewAccessControlledAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControlledAggregatorTransactor, error) {
	contract, err := bindAccessControlledAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorTransactor{contract: contract}, nil
}

// NewAccessControlledAggregatorFilterer creates a new log filterer instance of AccessControlledAggregator, bound to a specific deployed contract.
func NewAccessControlledAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControlledAggregatorFilterer, error) {
	contract, err := bindAccessControlledAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorFilterer{contract: contract}, nil
}

// bindAccessControlledAggregator binds a generic wrapper to an already deployed contract.
func bindAccessControlledAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControlledAggregator *AccessControlledAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledAggregator.Contract.AccessControlledAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControlledAggregator *AccessControlledAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AccessControlledAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControlledAggregator *AccessControlledAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AccessControlledAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControlledAggregator *AccessControlledAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.contract.Transact(opts, method, params...)
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) AllocatedFunds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "allocatedFunds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) AllocatedFunds() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.AllocatedFunds(&_AccessControlledAggregator.CallOpts)
}

// AllocatedFunds is a free data retrieval call binding the contract method 0xd4cc54e4.
//
// Solidity: function allocatedFunds() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) AllocatedFunds() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.AllocatedFunds(&_AccessControlledAggregator.CallOpts)
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) AvailableFunds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "availableFunds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) AvailableFunds() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.AvailableFunds(&_AccessControlledAggregator.CallOpts)
}

// AvailableFunds is a free data retrieval call binding the contract method 0x46fcff4c.
//
// Solidity: function availableFunds() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) AvailableFunds() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.AvailableFunds(&_AccessControlledAggregator.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) CheckEnabled() (bool, error) {
	return _AccessControlledAggregator.Contract.CheckEnabled(&_AccessControlledAggregator.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _AccessControlledAggregator.Contract.CheckEnabled(&_AccessControlledAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Decimals() (uint8, error) {
	return _AccessControlledAggregator.Contract.Decimals(&_AccessControlledAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) Decimals() (uint8, error) {
	return _AccessControlledAggregator.Contract.Decimals(&_AccessControlledAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Description() (string, error) {
	return _AccessControlledAggregator.Contract.Description(&_AccessControlledAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) Description() (string, error) {
	return _AccessControlledAggregator.Contract.Description(&_AccessControlledAggregator.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) GetAdmin(opts *bind.CallOpts, _oracle common.Address) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "getAdmin", _oracle)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) GetAdmin(_oracle common.Address) (common.Address, error) {
	return _AccessControlledAggregator.Contract.GetAdmin(&_AccessControlledAggregator.CallOpts, _oracle)
}

// GetAdmin is a free data retrieval call binding the contract method 0x64efb22b.
//
// Solidity: function getAdmin(address _oracle) view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) GetAdmin(_oracle common.Address) (common.Address, error) {
	return _AccessControlledAggregator.Contract.GetAdmin(&_AccessControlledAggregator.CallOpts, _oracle)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledAggregator.Contract.GetAnswer(&_AccessControlledAggregator.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledAggregator.Contract.GetAnswer(&_AccessControlledAggregator.CallOpts, _roundId)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) GetOracles(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "getOracles")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_AccessControlledAggregator *AccessControlledAggregatorSession) GetOracles() ([]common.Address, error) {
	return _AccessControlledAggregator.Contract.GetOracles(&_AccessControlledAggregator.CallOpts)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) GetOracles() ([]common.Address, error) {
	return _AccessControlledAggregator.Contract.GetOracles(&_AccessControlledAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

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
func (_AccessControlledAggregator *AccessControlledAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledAggregator.Contract.GetRoundData(&_AccessControlledAggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledAggregator.Contract.GetRoundData(&_AccessControlledAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledAggregator.Contract.GetTimestamp(&_AccessControlledAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledAggregator.Contract.GetTimestamp(&_AccessControlledAggregator.CallOpts, _roundId)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledAggregator.Contract.HasAccess(&_AccessControlledAggregator.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledAggregator.Contract.HasAccess(&_AccessControlledAggregator.CallOpts, _user, _calldata)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.LatestAnswer(&_AccessControlledAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.LatestAnswer(&_AccessControlledAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) LatestRound() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.LatestRound(&_AccessControlledAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.LatestRound(&_AccessControlledAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "latestRoundData")

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
func (_AccessControlledAggregator *AccessControlledAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledAggregator.Contract.LatestRoundData(&_AccessControlledAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledAggregator.Contract.LatestRoundData(&_AccessControlledAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.LatestTimestamp(&_AccessControlledAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.LatestTimestamp(&_AccessControlledAggregator.CallOpts)
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) LinkToken() (common.Address, error) {
	return _AccessControlledAggregator.Contract.LinkToken(&_AccessControlledAggregator.CallOpts)
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) LinkToken() (common.Address, error) {
	return _AccessControlledAggregator.Contract.LinkToken(&_AccessControlledAggregator.CallOpts)
}

// MaxSubmissionCount is a free data retrieval call binding the contract method 0x58609e44.
//
// Solidity: function maxSubmissionCount() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) MaxSubmissionCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "maxSubmissionCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MaxSubmissionCount is a free data retrieval call binding the contract method 0x58609e44.
//
// Solidity: function maxSubmissionCount() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) MaxSubmissionCount() (uint32, error) {
	return _AccessControlledAggregator.Contract.MaxSubmissionCount(&_AccessControlledAggregator.CallOpts)
}

// MaxSubmissionCount is a free data retrieval call binding the contract method 0x58609e44.
//
// Solidity: function maxSubmissionCount() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) MaxSubmissionCount() (uint32, error) {
	return _AccessControlledAggregator.Contract.MaxSubmissionCount(&_AccessControlledAggregator.CallOpts)
}

// MaxSubmissionValue is a free data retrieval call binding the contract method 0x23ca2903.
//
// Solidity: function maxSubmissionValue() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) MaxSubmissionValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "maxSubmissionValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxSubmissionValue is a free data retrieval call binding the contract method 0x23ca2903.
//
// Solidity: function maxSubmissionValue() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) MaxSubmissionValue() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.MaxSubmissionValue(&_AccessControlledAggregator.CallOpts)
}

// MaxSubmissionValue is a free data retrieval call binding the contract method 0x23ca2903.
//
// Solidity: function maxSubmissionValue() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) MaxSubmissionValue() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.MaxSubmissionValue(&_AccessControlledAggregator.CallOpts)
}

// MinSubmissionCount is a free data retrieval call binding the contract method 0xc9374500.
//
// Solidity: function minSubmissionCount() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) MinSubmissionCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "minSubmissionCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MinSubmissionCount is a free data retrieval call binding the contract method 0xc9374500.
//
// Solidity: function minSubmissionCount() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) MinSubmissionCount() (uint32, error) {
	return _AccessControlledAggregator.Contract.MinSubmissionCount(&_AccessControlledAggregator.CallOpts)
}

// MinSubmissionCount is a free data retrieval call binding the contract method 0xc9374500.
//
// Solidity: function minSubmissionCount() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) MinSubmissionCount() (uint32, error) {
	return _AccessControlledAggregator.Contract.MinSubmissionCount(&_AccessControlledAggregator.CallOpts)
}

// MinSubmissionValue is a free data retrieval call binding the contract method 0x7c2b0b21.
//
// Solidity: function minSubmissionValue() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) MinSubmissionValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "minSubmissionValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinSubmissionValue is a free data retrieval call binding the contract method 0x7c2b0b21.
//
// Solidity: function minSubmissionValue() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) MinSubmissionValue() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.MinSubmissionValue(&_AccessControlledAggregator.CallOpts)
}

// MinSubmissionValue is a free data retrieval call binding the contract method 0x7c2b0b21.
//
// Solidity: function minSubmissionValue() view returns(int256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) MinSubmissionValue() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.MinSubmissionValue(&_AccessControlledAggregator.CallOpts)
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() view returns(uint8)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) OracleCount(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "oracleCount")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() view returns(uint8)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) OracleCount() (uint8, error) {
	return _AccessControlledAggregator.Contract.OracleCount(&_AccessControlledAggregator.CallOpts)
}

// OracleCount is a free data retrieval call binding the contract method 0x613d8fcc.
//
// Solidity: function oracleCount() view returns(uint8)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) OracleCount() (uint8, error) {
	return _AccessControlledAggregator.Contract.OracleCount(&_AccessControlledAggregator.CallOpts)
}

// OracleRoundState is a free data retrieval call binding the contract method 0x88aa80e7.
//
// Solidity: function oracleRoundState(address _oracle, uint32 _queriedRoundId) view returns(bool _eligibleToSubmit, uint32 _roundId, int256 _latestSubmission, uint64 _startedAt, uint64 _timeout, uint128 _availableFunds, uint8 _oracleCount, uint128 _paymentAmount)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) OracleRoundState(opts *bind.CallOpts, _oracle common.Address, _queriedRoundId uint32) (struct {
	EligibleToSubmit bool
	RoundId          uint32
	LatestSubmission *big.Int
	StartedAt        uint64
	Timeout          uint64
	AvailableFunds   *big.Int
	OracleCount      uint8
	PaymentAmount    *big.Int
}, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "oracleRoundState", _oracle, _queriedRoundId)

	outstruct := new(struct {
		EligibleToSubmit bool
		RoundId          uint32
		LatestSubmission *big.Int
		StartedAt        uint64
		Timeout          uint64
		AvailableFunds   *big.Int
		OracleCount      uint8
		PaymentAmount    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EligibleToSubmit = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RoundId = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.LatestSubmission = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Timeout = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.AvailableFunds = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.OracleCount = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	outstruct.PaymentAmount = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// OracleRoundState is a free data retrieval call binding the contract method 0x88aa80e7.
//
// Solidity: function oracleRoundState(address _oracle, uint32 _queriedRoundId) view returns(bool _eligibleToSubmit, uint32 _roundId, int256 _latestSubmission, uint64 _startedAt, uint64 _timeout, uint128 _availableFunds, uint8 _oracleCount, uint128 _paymentAmount)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) OracleRoundState(_oracle common.Address, _queriedRoundId uint32) (struct {
	EligibleToSubmit bool
	RoundId          uint32
	LatestSubmission *big.Int
	StartedAt        uint64
	Timeout          uint64
	AvailableFunds   *big.Int
	OracleCount      uint8
	PaymentAmount    *big.Int
}, error) {
	return _AccessControlledAggregator.Contract.OracleRoundState(&_AccessControlledAggregator.CallOpts, _oracle, _queriedRoundId)
}

// OracleRoundState is a free data retrieval call binding the contract method 0x88aa80e7.
//
// Solidity: function oracleRoundState(address _oracle, uint32 _queriedRoundId) view returns(bool _eligibleToSubmit, uint32 _roundId, int256 _latestSubmission, uint64 _startedAt, uint64 _timeout, uint128 _availableFunds, uint8 _oracleCount, uint128 _paymentAmount)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) OracleRoundState(_oracle common.Address, _queriedRoundId uint32) (struct {
	EligibleToSubmit bool
	RoundId          uint32
	LatestSubmission *big.Int
	StartedAt        uint64
	Timeout          uint64
	AvailableFunds   *big.Int
	OracleCount      uint8
	PaymentAmount    *big.Int
}, error) {
	return _AccessControlledAggregator.Contract.OracleRoundState(&_AccessControlledAggregator.CallOpts, _oracle, _queriedRoundId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Owner() (common.Address, error) {
	return _AccessControlledAggregator.Contract.Owner(&_AccessControlledAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) Owner() (common.Address, error) {
	return _AccessControlledAggregator.Contract.Owner(&_AccessControlledAggregator.CallOpts)
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) PaymentAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "paymentAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) PaymentAmount() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.PaymentAmount(&_AccessControlledAggregator.CallOpts)
}

// PaymentAmount is a free data retrieval call binding the contract method 0xc35905c6.
//
// Solidity: function paymentAmount() view returns(uint128)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) PaymentAmount() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.PaymentAmount(&_AccessControlledAggregator.CallOpts)
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) RestartDelay(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "restartDelay")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) RestartDelay() (uint32, error) {
	return _AccessControlledAggregator.Contract.RestartDelay(&_AccessControlledAggregator.CallOpts)
}

// RestartDelay is a free data retrieval call binding the contract method 0x357ebb02.
//
// Solidity: function restartDelay() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) RestartDelay() (uint32, error) {
	return _AccessControlledAggregator.Contract.RestartDelay(&_AccessControlledAggregator.CallOpts)
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) Timeout(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "timeout")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Timeout() (uint32, error) {
	return _AccessControlledAggregator.Contract.Timeout(&_AccessControlledAggregator.CallOpts)
}

// Timeout is a free data retrieval call binding the contract method 0x70dea79a.
//
// Solidity: function timeout() view returns(uint32)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) Timeout() (uint32, error) {
	return _AccessControlledAggregator.Contract.Timeout(&_AccessControlledAggregator.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Validator() (common.Address, error) {
	return _AccessControlledAggregator.Contract.Validator(&_AccessControlledAggregator.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) Validator() (common.Address, error) {
	return _AccessControlledAggregator.Contract.Validator(&_AccessControlledAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Version() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.Version(&_AccessControlledAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) Version() (*big.Int, error) {
	return _AccessControlledAggregator.Contract.Version(&_AccessControlledAggregator.CallOpts)
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCaller) WithdrawablePayment(opts *bind.CallOpts, _oracle common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledAggregator.contract.Call(opts, &out, "withdrawablePayment", _oracle)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) WithdrawablePayment(_oracle common.Address) (*big.Int, error) {
	return _AccessControlledAggregator.Contract.WithdrawablePayment(&_AccessControlledAggregator.CallOpts, _oracle)
}

// WithdrawablePayment is a free data retrieval call binding the contract method 0xe2e40317.
//
// Solidity: function withdrawablePayment(address _oracle) view returns(uint256)
func (_AccessControlledAggregator *AccessControlledAggregatorCallerSession) WithdrawablePayment(_oracle common.Address) (*big.Int, error) {
	return _AccessControlledAggregator.Contract.WithdrawablePayment(&_AccessControlledAggregator.CallOpts, _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) AcceptAdmin(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "acceptAdmin", _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) AcceptAdmin(_oracle common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AcceptAdmin(&_AccessControlledAggregator.TransactOpts, _oracle)
}

// AcceptAdmin is a paid mutator transaction binding the contract method 0x628806ef.
//
// Solidity: function acceptAdmin(address _oracle) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) AcceptAdmin(_oracle common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AcceptAdmin(&_AccessControlledAggregator.TransactOpts, _oracle)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AcceptOwnership(&_AccessControlledAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AcceptOwnership(&_AccessControlledAggregator.TransactOpts)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AddAccess(&_AccessControlledAggregator.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.AddAccess(&_AccessControlledAggregator.TransactOpts, _user)
}

// ChangeOracles is a paid mutator transaction binding the contract method 0x3969c20f.
//
// Solidity: function changeOracles(address[] _removed, address[] _added, address[] _addedAdmins, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) ChangeOracles(opts *bind.TransactOpts, _removed []common.Address, _added []common.Address, _addedAdmins []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "changeOracles", _removed, _added, _addedAdmins, _minSubmissions, _maxSubmissions, _restartDelay)
}

// ChangeOracles is a paid mutator transaction binding the contract method 0x3969c20f.
//
// Solidity: function changeOracles(address[] _removed, address[] _added, address[] _addedAdmins, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) ChangeOracles(_removed []common.Address, _added []common.Address, _addedAdmins []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.ChangeOracles(&_AccessControlledAggregator.TransactOpts, _removed, _added, _addedAdmins, _minSubmissions, _maxSubmissions, _restartDelay)
}

// ChangeOracles is a paid mutator transaction binding the contract method 0x3969c20f.
//
// Solidity: function changeOracles(address[] _removed, address[] _added, address[] _addedAdmins, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) ChangeOracles(_removed []common.Address, _added []common.Address, _addedAdmins []common.Address, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.ChangeOracles(&_AccessControlledAggregator.TransactOpts, _removed, _added, _addedAdmins, _minSubmissions, _maxSubmissions, _restartDelay)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.DisableAccessCheck(&_AccessControlledAggregator.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.DisableAccessCheck(&_AccessControlledAggregator.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.EnableAccessCheck(&_AccessControlledAggregator.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.EnableAccessCheck(&_AccessControlledAggregator.TransactOpts)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes _data) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int, _data []byte) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "onTokenTransfer", arg0, arg1, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes _data) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) OnTokenTransfer(arg0 common.Address, arg1 *big.Int, _data []byte) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.OnTokenTransfer(&_AccessControlledAggregator.TransactOpts, arg0, arg1, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 , bytes _data) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) OnTokenTransfer(arg0 common.Address, arg1 *big.Int, _data []byte) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.OnTokenTransfer(&_AccessControlledAggregator.TransactOpts, arg0, arg1, _data)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.RemoveAccess(&_AccessControlledAggregator.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.RemoveAccess(&_AccessControlledAggregator.TransactOpts, _user)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_AccessControlledAggregator *AccessControlledAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.RequestNewRound(&_AccessControlledAggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.RequestNewRound(&_AccessControlledAggregator.TransactOpts)
}

// SetRequesterPermissions is a paid mutator transaction binding the contract method 0x20ed0275.
//
// Solidity: function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) SetRequesterPermissions(opts *bind.TransactOpts, _requester common.Address, _authorized bool, _delay uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "setRequesterPermissions", _requester, _authorized, _delay)
}

// SetRequesterPermissions is a paid mutator transaction binding the contract method 0x20ed0275.
//
// Solidity: function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) SetRequesterPermissions(_requester common.Address, _authorized bool, _delay uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.SetRequesterPermissions(&_AccessControlledAggregator.TransactOpts, _requester, _authorized, _delay)
}

// SetRequesterPermissions is a paid mutator transaction binding the contract method 0x20ed0275.
//
// Solidity: function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) SetRequesterPermissions(_requester common.Address, _authorized bool, _delay uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.SetRequesterPermissions(&_AccessControlledAggregator.TransactOpts, _requester, _authorized, _delay)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "setValidator", _newValidator)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.SetValidator(&_AccessControlledAggregator.TransactOpts, _newValidator)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.SetValidator(&_AccessControlledAggregator.TransactOpts, _newValidator)
}

// Submit is a paid mutator transaction binding the contract method 0x202ee0ed.
//
// Solidity: function submit(uint256 _roundId, int256 _submission) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) Submit(opts *bind.TransactOpts, _roundId *big.Int, _submission *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "submit", _roundId, _submission)
}

// Submit is a paid mutator transaction binding the contract method 0x202ee0ed.
//
// Solidity: function submit(uint256 _roundId, int256 _submission) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) Submit(_roundId *big.Int, _submission *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.Submit(&_AccessControlledAggregator.TransactOpts, _roundId, _submission)
}

// Submit is a paid mutator transaction binding the contract method 0x202ee0ed.
//
// Solidity: function submit(uint256 _roundId, int256 _submission) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) Submit(_roundId *big.Int, _submission *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.Submit(&_AccessControlledAggregator.TransactOpts, _roundId, _submission)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) TransferAdmin(opts *bind.TransactOpts, _oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "transferAdmin", _oracle, _newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) TransferAdmin(_oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.TransferAdmin(&_AccessControlledAggregator.TransactOpts, _oracle, _newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0xe9ee6eeb.
//
// Solidity: function transferAdmin(address _oracle, address _newAdmin) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) TransferAdmin(_oracle common.Address, _newAdmin common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.TransferAdmin(&_AccessControlledAggregator.TransactOpts, _oracle, _newAdmin)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.TransferOwnership(&_AccessControlledAggregator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.TransferOwnership(&_AccessControlledAggregator.TransactOpts, _to)
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) UpdateAvailableFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "updateAvailableFunds")
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) UpdateAvailableFunds() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.UpdateAvailableFunds(&_AccessControlledAggregator.TransactOpts)
}

// UpdateAvailableFunds is a paid mutator transaction binding the contract method 0x4f8fc3b5.
//
// Solidity: function updateAvailableFunds() returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) UpdateAvailableFunds() (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.UpdateAvailableFunds(&_AccessControlledAggregator.TransactOpts)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _paymentAmount, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay, uint32 _timeout) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) UpdateFutureRounds(opts *bind.TransactOpts, _paymentAmount *big.Int, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "updateFutureRounds", _paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, _timeout)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _paymentAmount, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay, uint32 _timeout) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) UpdateFutureRounds(_paymentAmount *big.Int, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.UpdateFutureRounds(&_AccessControlledAggregator.TransactOpts, _paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, _timeout)
}

// UpdateFutureRounds is a paid mutator transaction binding the contract method 0x38aa4c72.
//
// Solidity: function updateFutureRounds(uint128 _paymentAmount, uint32 _minSubmissions, uint32 _maxSubmissions, uint32 _restartDelay, uint32 _timeout) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) UpdateFutureRounds(_paymentAmount *big.Int, _minSubmissions uint32, _maxSubmissions uint32, _restartDelay uint32, _timeout uint32) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.UpdateFutureRounds(&_AccessControlledAggregator.TransactOpts, _paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, _timeout)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.WithdrawFunds(&_AccessControlledAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.WithdrawFunds(&_AccessControlledAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.contract.Transact(opts, "withdrawPayment", _oracle, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorSession) WithdrawPayment(_oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.WithdrawPayment(&_AccessControlledAggregator.TransactOpts, _oracle, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x3d3d7714.
//
// Solidity: function withdrawPayment(address _oracle, address _recipient, uint256 _amount) returns()
func (_AccessControlledAggregator *AccessControlledAggregatorTransactorSession) WithdrawPayment(_oracle common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledAggregator.Contract.WithdrawPayment(&_AccessControlledAggregator.TransactOpts, _oracle, _recipient, _amount)
}

// AccessControlledAggregatorAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorAddedAccessIterator struct {
	Event *AccessControlledAggregatorAddedAccess // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorAddedAccess)
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
		it.Event = new(AccessControlledAggregatorAddedAccess)
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
func (it *AccessControlledAggregatorAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorAddedAccess represents a AddedAccess event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*AccessControlledAggregatorAddedAccessIterator, error) {

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorAddedAccessIterator{contract: _AccessControlledAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorAddedAccess)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseAddedAccess(log types.Log) (*AccessControlledAggregatorAddedAccess, error) {
	event := new(AccessControlledAggregatorAddedAccess)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorAnswerUpdatedIterator struct {
	Event *AccessControlledAggregatorAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorAnswerUpdated)
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
		it.Event = new(AccessControlledAggregatorAnswerUpdated)
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
func (it *AccessControlledAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorAnswerUpdated represents a AnswerUpdated event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AccessControlledAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorAnswerUpdatedIterator{contract: _AccessControlledAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorAnswerUpdated)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*AccessControlledAggregatorAnswerUpdated, error) {
	event := new(AccessControlledAggregatorAnswerUpdated)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorAvailableFundsUpdatedIterator is returned from FilterAvailableFundsUpdated and is used to iterate over the raw logs and unpacked data for AvailableFundsUpdated events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorAvailableFundsUpdatedIterator struct {
	Event *AccessControlledAggregatorAvailableFundsUpdated // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorAvailableFundsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorAvailableFundsUpdated)
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
		it.Event = new(AccessControlledAggregatorAvailableFundsUpdated)
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
func (it *AccessControlledAggregatorAvailableFundsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorAvailableFundsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorAvailableFundsUpdated represents a AvailableFundsUpdated event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorAvailableFundsUpdated struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAvailableFundsUpdated is a free log retrieval operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterAvailableFundsUpdated(opts *bind.FilterOpts, amount []*big.Int) (*AccessControlledAggregatorAvailableFundsUpdatedIterator, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "AvailableFundsUpdated", amountRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorAvailableFundsUpdatedIterator{contract: _AccessControlledAggregator.contract, event: "AvailableFundsUpdated", logs: logs, sub: sub}, nil
}

// WatchAvailableFundsUpdated is a free log subscription operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchAvailableFundsUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorAvailableFundsUpdated, amount []*big.Int) (event.Subscription, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "AvailableFundsUpdated", amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorAvailableFundsUpdated)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "AvailableFundsUpdated", log); err != nil {
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

// ParseAvailableFundsUpdated is a log parse operation binding the contract event 0xfe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234f.
//
// Solidity: event AvailableFundsUpdated(uint256 indexed amount)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseAvailableFundsUpdated(log types.Log) (*AccessControlledAggregatorAvailableFundsUpdated, error) {
	event := new(AccessControlledAggregatorAvailableFundsUpdated)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "AvailableFundsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorCheckAccessDisabledIterator struct {
	Event *AccessControlledAggregatorCheckAccessDisabled // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorCheckAccessDisabled)
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
		it.Event = new(AccessControlledAggregatorCheckAccessDisabled)
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
func (it *AccessControlledAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorCheckAccessDisabled represents a CheckAccessDisabled event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*AccessControlledAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorCheckAccessDisabledIterator{contract: _AccessControlledAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorCheckAccessDisabled)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*AccessControlledAggregatorCheckAccessDisabled, error) {
	event := new(AccessControlledAggregatorCheckAccessDisabled)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorCheckAccessEnabledIterator struct {
	Event *AccessControlledAggregatorCheckAccessEnabled // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorCheckAccessEnabled)
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
		it.Event = new(AccessControlledAggregatorCheckAccessEnabled)
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
func (it *AccessControlledAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorCheckAccessEnabled represents a CheckAccessEnabled event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*AccessControlledAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorCheckAccessEnabledIterator{contract: _AccessControlledAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorCheckAccessEnabled)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*AccessControlledAggregatorCheckAccessEnabled, error) {
	event := new(AccessControlledAggregatorCheckAccessEnabled)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorNewRoundIterator struct {
	Event *AccessControlledAggregatorNewRound // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorNewRound)
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
		it.Event = new(AccessControlledAggregatorNewRound)
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
func (it *AccessControlledAggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorNewRound represents a NewRound event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AccessControlledAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorNewRoundIterator{contract: _AccessControlledAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorNewRound)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseNewRound(log types.Log) (*AccessControlledAggregatorNewRound, error) {
	event := new(AccessControlledAggregatorNewRound)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorOracleAdminUpdateRequestedIterator is returned from FilterOracleAdminUpdateRequested and is used to iterate over the raw logs and unpacked data for OracleAdminUpdateRequested events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOracleAdminUpdateRequestedIterator struct {
	Event *AccessControlledAggregatorOracleAdminUpdateRequested // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorOracleAdminUpdateRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorOracleAdminUpdateRequested)
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
		it.Event = new(AccessControlledAggregatorOracleAdminUpdateRequested)
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
func (it *AccessControlledAggregatorOracleAdminUpdateRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorOracleAdminUpdateRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorOracleAdminUpdateRequested represents a OracleAdminUpdateRequested event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOracleAdminUpdateRequested struct {
	Oracle   common.Address
	Admin    common.Address
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOracleAdminUpdateRequested is a free log retrieval operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterOracleAdminUpdateRequested(opts *bind.FilterOpts, oracle []common.Address) (*AccessControlledAggregatorOracleAdminUpdateRequestedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "OracleAdminUpdateRequested", oracleRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorOracleAdminUpdateRequestedIterator{contract: _AccessControlledAggregator.contract, event: "OracleAdminUpdateRequested", logs: logs, sub: sub}, nil
}

// WatchOracleAdminUpdateRequested is a free log subscription operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchOracleAdminUpdateRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorOracleAdminUpdateRequested, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "OracleAdminUpdateRequested", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorOracleAdminUpdateRequested)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "OracleAdminUpdateRequested", log); err != nil {
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

// ParseOracleAdminUpdateRequested is a log parse operation binding the contract event 0xb79bf2e89c2d70dde91d2991fb1ea69b7e478061ad7c04ed5b02b96bc52b8104.
//
// Solidity: event OracleAdminUpdateRequested(address indexed oracle, address admin, address newAdmin)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseOracleAdminUpdateRequested(log types.Log) (*AccessControlledAggregatorOracleAdminUpdateRequested, error) {
	event := new(AccessControlledAggregatorOracleAdminUpdateRequested)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "OracleAdminUpdateRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorOracleAdminUpdatedIterator is returned from FilterOracleAdminUpdated and is used to iterate over the raw logs and unpacked data for OracleAdminUpdated events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOracleAdminUpdatedIterator struct {
	Event *AccessControlledAggregatorOracleAdminUpdated // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorOracleAdminUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorOracleAdminUpdated)
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
		it.Event = new(AccessControlledAggregatorOracleAdminUpdated)
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
func (it *AccessControlledAggregatorOracleAdminUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorOracleAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorOracleAdminUpdated represents a OracleAdminUpdated event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOracleAdminUpdated struct {
	Oracle   common.Address
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOracleAdminUpdated is a free log retrieval operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterOracleAdminUpdated(opts *bind.FilterOpts, oracle []common.Address, newAdmin []common.Address) (*AccessControlledAggregatorOracleAdminUpdatedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "OracleAdminUpdated", oracleRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorOracleAdminUpdatedIterator{contract: _AccessControlledAggregator.contract, event: "OracleAdminUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleAdminUpdated is a free log subscription operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchOracleAdminUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorOracleAdminUpdated, oracle []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "OracleAdminUpdated", oracleRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorOracleAdminUpdated)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "OracleAdminUpdated", log); err != nil {
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

// ParseOracleAdminUpdated is a log parse operation binding the contract event 0x0c5055390645c15a4be9a21b3f8d019153dcb4a0c125685da6eb84048e2fe904.
//
// Solidity: event OracleAdminUpdated(address indexed oracle, address indexed newAdmin)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseOracleAdminUpdated(log types.Log) (*AccessControlledAggregatorOracleAdminUpdated, error) {
	event := new(AccessControlledAggregatorOracleAdminUpdated)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "OracleAdminUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorOraclePermissionsUpdatedIterator is returned from FilterOraclePermissionsUpdated and is used to iterate over the raw logs and unpacked data for OraclePermissionsUpdated events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOraclePermissionsUpdatedIterator struct {
	Event *AccessControlledAggregatorOraclePermissionsUpdated // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorOraclePermissionsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorOraclePermissionsUpdated)
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
		it.Event = new(AccessControlledAggregatorOraclePermissionsUpdated)
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
func (it *AccessControlledAggregatorOraclePermissionsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorOraclePermissionsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorOraclePermissionsUpdated represents a OraclePermissionsUpdated event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOraclePermissionsUpdated struct {
	Oracle      common.Address
	Whitelisted bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePermissionsUpdated is a free log retrieval operation binding the contract event 0x18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e.
//
// Solidity: event OraclePermissionsUpdated(address indexed oracle, bool indexed whitelisted)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterOraclePermissionsUpdated(opts *bind.FilterOpts, oracle []common.Address, whitelisted []bool) (*AccessControlledAggregatorOraclePermissionsUpdatedIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var whitelistedRule []interface{}
	for _, whitelistedItem := range whitelisted {
		whitelistedRule = append(whitelistedRule, whitelistedItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "OraclePermissionsUpdated", oracleRule, whitelistedRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorOraclePermissionsUpdatedIterator{contract: _AccessControlledAggregator.contract, event: "OraclePermissionsUpdated", logs: logs, sub: sub}, nil
}

// WatchOraclePermissionsUpdated is a free log subscription operation binding the contract event 0x18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e.
//
// Solidity: event OraclePermissionsUpdated(address indexed oracle, bool indexed whitelisted)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchOraclePermissionsUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorOraclePermissionsUpdated, oracle []common.Address, whitelisted []bool) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var whitelistedRule []interface{}
	for _, whitelistedItem := range whitelisted {
		whitelistedRule = append(whitelistedRule, whitelistedItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "OraclePermissionsUpdated", oracleRule, whitelistedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorOraclePermissionsUpdated)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "OraclePermissionsUpdated", log); err != nil {
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

// ParseOraclePermissionsUpdated is a log parse operation binding the contract event 0x18dd09695e4fbdae8d1a5edb11221eb04564269c29a089b9753a6535c54ba92e.
//
// Solidity: event OraclePermissionsUpdated(address indexed oracle, bool indexed whitelisted)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseOraclePermissionsUpdated(log types.Log) (*AccessControlledAggregatorOraclePermissionsUpdated, error) {
	event := new(AccessControlledAggregatorOraclePermissionsUpdated)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "OraclePermissionsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOwnershipTransferRequestedIterator struct {
	Event *AccessControlledAggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorOwnershipTransferRequested)
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
		it.Event = new(AccessControlledAggregatorOwnershipTransferRequested)
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
func (it *AccessControlledAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorOwnershipTransferRequestedIterator{contract: _AccessControlledAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorOwnershipTransferRequested)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*AccessControlledAggregatorOwnershipTransferRequested, error) {
	event := new(AccessControlledAggregatorOwnershipTransferRequested)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOwnershipTransferredIterator struct {
	Event *AccessControlledAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorOwnershipTransferred)
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
		it.Event = new(AccessControlledAggregatorOwnershipTransferred)
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
func (it *AccessControlledAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorOwnershipTransferredIterator{contract: _AccessControlledAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorOwnershipTransferred)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*AccessControlledAggregatorOwnershipTransferred, error) {
	event := new(AccessControlledAggregatorOwnershipTransferred)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorRemovedAccessIterator struct {
	Event *AccessControlledAggregatorRemovedAccess // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorRemovedAccess)
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
		it.Event = new(AccessControlledAggregatorRemovedAccess)
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
func (it *AccessControlledAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorRemovedAccess represents a RemovedAccess event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*AccessControlledAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorRemovedAccessIterator{contract: _AccessControlledAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorRemovedAccess)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseRemovedAccess(log types.Log) (*AccessControlledAggregatorRemovedAccess, error) {
	event := new(AccessControlledAggregatorRemovedAccess)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorRequesterPermissionsSetIterator is returned from FilterRequesterPermissionsSet and is used to iterate over the raw logs and unpacked data for RequesterPermissionsSet events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorRequesterPermissionsSetIterator struct {
	Event *AccessControlledAggregatorRequesterPermissionsSet // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorRequesterPermissionsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorRequesterPermissionsSet)
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
		it.Event = new(AccessControlledAggregatorRequesterPermissionsSet)
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
func (it *AccessControlledAggregatorRequesterPermissionsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorRequesterPermissionsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorRequesterPermissionsSet represents a RequesterPermissionsSet event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorRequesterPermissionsSet struct {
	Requester  common.Address
	Authorized bool
	Delay      uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRequesterPermissionsSet is a free log retrieval operation binding the contract event 0xc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a.
//
// Solidity: event RequesterPermissionsSet(address indexed requester, bool authorized, uint32 delay)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterRequesterPermissionsSet(opts *bind.FilterOpts, requester []common.Address) (*AccessControlledAggregatorRequesterPermissionsSetIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "RequesterPermissionsSet", requesterRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorRequesterPermissionsSetIterator{contract: _AccessControlledAggregator.contract, event: "RequesterPermissionsSet", logs: logs, sub: sub}, nil
}

// WatchRequesterPermissionsSet is a free log subscription operation binding the contract event 0xc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a.
//
// Solidity: event RequesterPermissionsSet(address indexed requester, bool authorized, uint32 delay)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchRequesterPermissionsSet(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorRequesterPermissionsSet, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "RequesterPermissionsSet", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorRequesterPermissionsSet)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "RequesterPermissionsSet", log); err != nil {
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

// ParseRequesterPermissionsSet is a log parse operation binding the contract event 0xc3df5a754e002718f2e10804b99e6605e7c701d95cec9552c7680ca2b6f2820a.
//
// Solidity: event RequesterPermissionsSet(address indexed requester, bool authorized, uint32 delay)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseRequesterPermissionsSet(log types.Log) (*AccessControlledAggregatorRequesterPermissionsSet, error) {
	event := new(AccessControlledAggregatorRequesterPermissionsSet)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "RequesterPermissionsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorRoundDetailsUpdatedIterator is returned from FilterRoundDetailsUpdated and is used to iterate over the raw logs and unpacked data for RoundDetailsUpdated events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorRoundDetailsUpdatedIterator struct {
	Event *AccessControlledAggregatorRoundDetailsUpdated // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorRoundDetailsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorRoundDetailsUpdated)
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
		it.Event = new(AccessControlledAggregatorRoundDetailsUpdated)
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
func (it *AccessControlledAggregatorRoundDetailsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorRoundDetailsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorRoundDetailsUpdated represents a RoundDetailsUpdated event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorRoundDetailsUpdated struct {
	PaymentAmount      *big.Int
	MinSubmissionCount uint32
	MaxSubmissionCount uint32
	RestartDelay       uint32
	Timeout            uint32
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRoundDetailsUpdated is a free log retrieval operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minSubmissionCount, uint32 indexed maxSubmissionCount, uint32 restartDelay, uint32 timeout)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterRoundDetailsUpdated(opts *bind.FilterOpts, paymentAmount []*big.Int, minSubmissionCount []uint32, maxSubmissionCount []uint32) (*AccessControlledAggregatorRoundDetailsUpdatedIterator, error) {

	var paymentAmountRule []interface{}
	for _, paymentAmountItem := range paymentAmount {
		paymentAmountRule = append(paymentAmountRule, paymentAmountItem)
	}
	var minSubmissionCountRule []interface{}
	for _, minSubmissionCountItem := range minSubmissionCount {
		minSubmissionCountRule = append(minSubmissionCountRule, minSubmissionCountItem)
	}
	var maxSubmissionCountRule []interface{}
	for _, maxSubmissionCountItem := range maxSubmissionCount {
		maxSubmissionCountRule = append(maxSubmissionCountRule, maxSubmissionCountItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "RoundDetailsUpdated", paymentAmountRule, minSubmissionCountRule, maxSubmissionCountRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorRoundDetailsUpdatedIterator{contract: _AccessControlledAggregator.contract, event: "RoundDetailsUpdated", logs: logs, sub: sub}, nil
}

// WatchRoundDetailsUpdated is a free log subscription operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minSubmissionCount, uint32 indexed maxSubmissionCount, uint32 restartDelay, uint32 timeout)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchRoundDetailsUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorRoundDetailsUpdated, paymentAmount []*big.Int, minSubmissionCount []uint32, maxSubmissionCount []uint32) (event.Subscription, error) {

	var paymentAmountRule []interface{}
	for _, paymentAmountItem := range paymentAmount {
		paymentAmountRule = append(paymentAmountRule, paymentAmountItem)
	}
	var minSubmissionCountRule []interface{}
	for _, minSubmissionCountItem := range minSubmissionCount {
		minSubmissionCountRule = append(minSubmissionCountRule, minSubmissionCountItem)
	}
	var maxSubmissionCountRule []interface{}
	for _, maxSubmissionCountItem := range maxSubmissionCount {
		maxSubmissionCountRule = append(maxSubmissionCountRule, maxSubmissionCountItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "RoundDetailsUpdated", paymentAmountRule, minSubmissionCountRule, maxSubmissionCountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorRoundDetailsUpdated)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "RoundDetailsUpdated", log); err != nil {
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

// ParseRoundDetailsUpdated is a log parse operation binding the contract event 0x56800c9d1ed723511246614d15e58cfcde15b6a33c245b5c961b689c1890fd8f.
//
// Solidity: event RoundDetailsUpdated(uint128 indexed paymentAmount, uint32 indexed minSubmissionCount, uint32 indexed maxSubmissionCount, uint32 restartDelay, uint32 timeout)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseRoundDetailsUpdated(log types.Log) (*AccessControlledAggregatorRoundDetailsUpdated, error) {
	event := new(AccessControlledAggregatorRoundDetailsUpdated)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "RoundDetailsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorSubmissionReceivedIterator is returned from FilterSubmissionReceived and is used to iterate over the raw logs and unpacked data for SubmissionReceived events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorSubmissionReceivedIterator struct {
	Event *AccessControlledAggregatorSubmissionReceived // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorSubmissionReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorSubmissionReceived)
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
		it.Event = new(AccessControlledAggregatorSubmissionReceived)
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
func (it *AccessControlledAggregatorSubmissionReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorSubmissionReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorSubmissionReceived represents a SubmissionReceived event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorSubmissionReceived struct {
	Submission *big.Int
	Round      uint32
	Oracle     common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSubmissionReceived is a free log retrieval operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed submission, uint32 indexed round, address indexed oracle)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterSubmissionReceived(opts *bind.FilterOpts, submission []*big.Int, round []uint32, oracle []common.Address) (*AccessControlledAggregatorSubmissionReceivedIterator, error) {

	var submissionRule []interface{}
	for _, submissionItem := range submission {
		submissionRule = append(submissionRule, submissionItem)
	}
	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "SubmissionReceived", submissionRule, roundRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorSubmissionReceivedIterator{contract: _AccessControlledAggregator.contract, event: "SubmissionReceived", logs: logs, sub: sub}, nil
}

// WatchSubmissionReceived is a free log subscription operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed submission, uint32 indexed round, address indexed oracle)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchSubmissionReceived(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorSubmissionReceived, submission []*big.Int, round []uint32, oracle []common.Address) (event.Subscription, error) {

	var submissionRule []interface{}
	for _, submissionItem := range submission {
		submissionRule = append(submissionRule, submissionItem)
	}
	var roundRule []interface{}
	for _, roundItem := range round {
		roundRule = append(roundRule, roundItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "SubmissionReceived", submissionRule, roundRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorSubmissionReceived)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "SubmissionReceived", log); err != nil {
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

// ParseSubmissionReceived is a log parse operation binding the contract event 0x92e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68c.
//
// Solidity: event SubmissionReceived(int256 indexed submission, uint32 indexed round, address indexed oracle)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseSubmissionReceived(log types.Log) (*AccessControlledAggregatorSubmissionReceived, error) {
	event := new(AccessControlledAggregatorSubmissionReceived)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "SubmissionReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledAggregatorValidatorUpdatedIterator is returned from FilterValidatorUpdated and is used to iterate over the raw logs and unpacked data for ValidatorUpdated events raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorValidatorUpdatedIterator struct {
	Event *AccessControlledAggregatorValidatorUpdated // Event containing the contract specifics and raw log

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
func (it *AccessControlledAggregatorValidatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledAggregatorValidatorUpdated)
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
		it.Event = new(AccessControlledAggregatorValidatorUpdated)
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
func (it *AccessControlledAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledAggregatorValidatorUpdated represents a ValidatorUpdated event raised by the AccessControlledAggregator contract.
type AccessControlledAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValidatorUpdated is a free log retrieval operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) FilterValidatorUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*AccessControlledAggregatorValidatorUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.FilterLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledAggregatorValidatorUpdatedIterator{contract: _AccessControlledAggregator.contract, event: "ValidatorUpdated", logs: logs, sub: sub}, nil
}

// WatchValidatorUpdated is a free log subscription operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) WatchValidatorUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledAggregatorValidatorUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AccessControlledAggregator.contract.WatchLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledAggregatorValidatorUpdated)
				if err := _AccessControlledAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
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

// ParseValidatorUpdated is a log parse operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
func (_AccessControlledAggregator *AccessControlledAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*AccessControlledAggregatorValidatorUpdated, error) {
	event := new(AccessControlledAggregatorValidatorUpdated)
	if err := _AccessControlledAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
