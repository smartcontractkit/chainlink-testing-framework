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

// AggregatorProxyABI is the input ABI used to generate the binding from.
const AggregatorProxyABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"aggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"confirmAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"name\":\"phaseAggregators\",\"outputs\":[{\"internalType\":\"contractAggregatorV2V3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"phaseId\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"proposeAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposedAggregator\",\"outputs\":[{\"internalType\":\"contractAggregatorV2V3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"proposedGetRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposedLatestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AggregatorProxyBin is the compiled bytecode used for deploying new contracts.
var AggregatorProxyBin = "0x608060405234801561001057600080fd5b50604051611ff1380380611ff18339818101604052602081101561003357600080fd5b8101908080519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506100938161009960201b60201c565b506101b0565b60006001600260000160009054906101000a900461ffff1601905060405180604001604052808261ffff1681526020018373ffffffffffffffffffffffffffffffffffffffff16815250600260008201518160000160006101000a81548161ffff021916908361ffff16021790555060208201518160000160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505081600460008361ffff1661ffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b611e32806101bf6000396000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c80638da5cb5b116100b8578063b633620c1161007c578063b633620c14610543578063c159730414610585578063e8c4be30146105f7578063f2fde38b14610641578063f8a2abd314610685578063feaf968c146106c957610137565b80638da5cb5b1461036f5780638f6b4d91146103b95780639a6fc8f514610423578063a928c096146104bd578063b5ab58dc1461050157610137565b80636001ac53116100ff5780636001ac531461020c578063668a0f02146102a65780637284e416146102c457806379ba5097146103475780638205bf6a1461035157610137565b8063245a7bfc1461013c578063313ce5671461018657806350d25bcd146101aa57806354fd4d50146101c857806358303b10146101e6575b600080fd5b610144610733565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61018e610760565b604051808260ff1660ff16815260200191505060405180910390f35b6101b261080d565b6040518082815260200191505060405180910390f35b6101d06108ba565b6040518082815260200191505060405180910390f35b6101ee610967565b604051808261ffff1661ffff16815260200191505060405180910390f35b6102446004803603602081101561022257600080fd5b81019080803569ffffffffffffffffffff169060200190929190505050610982565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102ae610b4f565b6040518082815260200191505060405180910390f35b6102cc610c82565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561030c5780820151818401526020810190506102f1565b50505050905090810190601f1680156103395780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61034f610ded565b005b610359610fb5565b6040518082815260200191505060405180910390f35b610377611062565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6103c1611087565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b61045b6004803603602081101561043957600080fd5b81019080803569ffffffffffffffffffff16906020019092919050505061122f565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b6104ff600480360360208110156104d357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061138f565b005b61052d6004803603602081101561051757600080fd5b8101908080359060200190929190505050611545565b6040518082815260200191505060405180910390f35b61056f6004803603602081101561055957600080fd5b8101908080359060200190929190505050611694565b6040518082815260200191505060405180910390f35b6105b56004803603602081101561059b57600080fd5b81019080803561ffff1690602001909291905050506117e3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6105ff611816565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6106836004803603602081101561065757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061183c565b005b6106c76004803603602081101561069b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506119bd565b005b6106d1611ac3565b604051808669ffffffffffffffffffff1669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1669ffffffffffffffffffff1681526020019550505050505060405180910390f35b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b1580156107cd57600080fd5b505afa1580156107e1573d6000803e3d6000fd5b505050506040513d60208110156107f757600080fd5b8101908080519060200190929190505050905090565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b815260040160206040518083038186803b15801561087a57600080fd5b505afa15801561088e573d6000803e3d6000fd5b505050506040513d60208110156108a457600080fd5b8101908080519060200190929190505050905090565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166354fd4d506040518163ffffffff1660e01b815260040160206040518083038186803b15801561092757600080fd5b505afa15801561093b573d6000803e3d6000fd5b505050506040513d602081101561095157600080fd5b8101908080519060200190929190505050905090565b6000600260000160009054906101000a900461ffff16905090565b60008060008060008073ffffffffffffffffffffffffffffffffffffffff16600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415610a4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f4e6f2070726f706f7365642061676772656761746f722070726573656e74000081525060200191505060405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639a6fc8f5876040518263ffffffff1660e01b8152600401808269ffffffffffffffffffff1669ffffffffffffffffffff16815260200191505060a06040518083038186803b158015610ad957600080fd5b505afa158015610aed573d6000803e3d6000fd5b505050506040513d60a0811015610b0357600080fd5b8101908080519060200190929190805190602001909291908051906020019092919080519060200190929190805190602001909291905050509450945094509450945091939590929450565b6000610b59611dc8565b60026040518060400160405290816000820160009054906101000a900461ffff1661ffff1661ffff1681526020016000820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815250509050610c708160000151826020015173ffffffffffffffffffffffffffffffffffffffff1663668a0f026040518163ffffffff1660e01b815260040160206040518083038186803b158015610c3057600080fd5b505afa158015610c44573d6000803e3d6000fd5b505050506040513d6020811015610c5a57600080fd5b8101908080519060200190929190505050611c3d565b69ffffffffffffffffffff1691505090565b6060600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16637284e4166040518163ffffffff1660e01b815260040160006040518083038186803b158015610cef57600080fd5b505afa158015610d03573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506020811015610d2d57600080fd5b8101908080516040519392919084640100000000821115610d4d57600080fd5b83820191506020820185811115610d6357600080fd5b8251866001820283011164010000000082111715610d8057600080fd5b8083526020830192505050908051906020019080838360005b83811015610db4578082015181840152602081019050610d99565b50505050905090810190601f168015610de15780820380516001836020036101000a031916815260200191505b50604052505050905090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610eb0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4d7573742062652070726f706f736564206f776e65720000000000000000000081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a350565b6000600260000160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638205bf6a6040518163ffffffff1660e01b815260040160206040518083038186803b15801561102257600080fd5b505afa158015611036573d6000803e3d6000fd5b505050506040513d602081101561104c57600080fd5b8101908080519060200190929190505050905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008060008073ffffffffffffffffffffffffffffffffffffffff16600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415611153576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f4e6f2070726f706f7365642061676772656761746f722070726573656e74000081525060200191505060405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156111bb57600080fd5b505afa1580156111cf573d6000803e3d6000fd5b505050506040513d60a08110156111e557600080fd5b810190808051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050945094509450945094509091929394565b600080600080600080600061124f8869ffffffffffffffffffff16611c5c565b91509150600460008361ffff1661ffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639a6fc8f5826040518263ffffffff1660e01b8152600401808267ffffffffffffffff1669ffffffffffffffffffff16815260200191505060a06040518083038186803b1580156112f557600080fd5b505afa158015611309573d6000803e3d6000fd5b505050506040513d60a081101561131f57600080fd5b810190808051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050809750819850829950839a50849b50505050505061137a878787878787611c7a565b96509650965096509650505091939590929450565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611451576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614611514576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f496e76616c69642070726f706f7365642061676772656761746f72000000000081525060200191505060405180910390fd5b600360006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905561154281611cb1565b50565b60006001601060400160020a03821115611562576000905061168f565b60008061156e84611c5c565b915091506000600460008461ffff1661ffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156115f3576000935050505061168f565b8073ffffffffffffffffffffffffffffffffffffffff1663b5ab58dc836040518263ffffffff1660e01b8152600401808267ffffffffffffffff16815260200191505060206040518083038186803b15801561164e57600080fd5b505afa158015611662573d6000803e3d6000fd5b505050506040513d602081101561167857600080fd5b810190808051906020019092919050505093505050505b919050565b60006001601060400160020a038211156116b157600090506117de565b6000806116bd84611c5c565b915091506000600460008461ffff1661ffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561174257600093505050506117de565b8073ffffffffffffffffffffffffffffffffffffffff1663b633620c836040518263ffffffff1660e01b8152600401808267ffffffffffffffff16815260200191505060206040518083038186803b15801561179d57600080fd5b505afa1580156117b1573d6000803e3d6000fd5b505050506040513d60208110156117c757600080fd5b810190808051906020019092919050505093505050505b919050565b60046020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146118fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611a7f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000806000806000611ad3611dc8565b60026040518060400160405290816000820160009054906101000a900461ffff1661ffff1661ffff1681526020016000820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815250509050806020015173ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015611ba257600080fd5b505afa158015611bb6573d6000803e3d6000fd5b505050506040513d60a0811015611bcc57600080fd5b810190808051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050809650819750829850839950849a505050505050611c2b86868686868660000151611c7a565b95509550955095509550509091929394565b60008167ffffffffffffffff1660408461ffff16901b17905092915050565b6000806000604084901c905060008490508181935093505050915091565b6000806000806000611c8c868c611c3d565b8a8a8a611c998a8c611c3d565b94509450945094509450965096509650965096915050565b60006001600260000160009054906101000a900461ffff1601905060405180604001604052808261ffff1681526020018373ffffffffffffffffffffffffffffffffffffffff16815250600260008201518160000160006101000a81548161ffff021916908361ffff16021790555060208201518160000160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505081600460008361ffff1661ffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b6040518060400160405280600061ffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152509056fea2646970667358221220447a1ac2b249a058936277eb4ecb10ca1be14bbab05e07d4a4bdbbe656b5a0d464736f6c63430006060033"

// DeployAggregatorProxy deploys a new Ethereum contract, binding an instance of AggregatorProxy to it.
func DeployAggregatorProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _aggregator common.Address) (common.Address, *types.Transaction, *AggregatorProxy, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorProxyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AggregatorProxyBin), backend, _aggregator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AggregatorProxy{AggregatorProxyCaller: AggregatorProxyCaller{contract: contract}, AggregatorProxyTransactor: AggregatorProxyTransactor{contract: contract}, AggregatorProxyFilterer: AggregatorProxyFilterer{contract: contract}}, nil
}

// AggregatorProxy is an auto generated Go binding around an Ethereum contract.
type AggregatorProxy struct {
	AggregatorProxyCaller     // Read-only binding to the contract
	AggregatorProxyTransactor // Write-only binding to the contract
	AggregatorProxyFilterer   // Log filterer for contract events
}

// AggregatorProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorProxySession struct {
	Contract     *AggregatorProxy  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AggregatorProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorProxyCallerSession struct {
	Contract *AggregatorProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AggregatorProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorProxyTransactorSession struct {
	Contract     *AggregatorProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AggregatorProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorProxyRaw struct {
	Contract *AggregatorProxy // Generic contract binding to access the raw methods on
}

// AggregatorProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorProxyCallerRaw struct {
	Contract *AggregatorProxyCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorProxyTransactorRaw struct {
	Contract *AggregatorProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorProxy creates a new instance of AggregatorProxy, bound to a specific deployed contract.
func NewAggregatorProxy(address common.Address, backend bind.ContractBackend) (*AggregatorProxy, error) {
	contract, err := bindAggregatorProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxy{AggregatorProxyCaller: AggregatorProxyCaller{contract: contract}, AggregatorProxyTransactor: AggregatorProxyTransactor{contract: contract}, AggregatorProxyFilterer: AggregatorProxyFilterer{contract: contract}}, nil
}

// NewAggregatorProxyCaller creates a new read-only instance of AggregatorProxy, bound to a specific deployed contract.
func NewAggregatorProxyCaller(address common.Address, caller bind.ContractCaller) (*AggregatorProxyCaller, error) {
	contract, err := bindAggregatorProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyCaller{contract: contract}, nil
}

// NewAggregatorProxyTransactor creates a new write-only instance of AggregatorProxy, bound to a specific deployed contract.
func NewAggregatorProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorProxyTransactor, error) {
	contract, err := bindAggregatorProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyTransactor{contract: contract}, nil
}

// NewAggregatorProxyFilterer creates a new log filterer instance of AggregatorProxy, bound to a specific deployed contract.
func NewAggregatorProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorProxyFilterer, error) {
	contract, err := bindAggregatorProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyFilterer{contract: contract}, nil
}

// bindAggregatorProxy binds a generic wrapper to an already deployed contract.
func bindAggregatorProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorProxy *AggregatorProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorProxy.Contract.AggregatorProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorProxy *AggregatorProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AggregatorProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorProxy *AggregatorProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AggregatorProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorProxy *AggregatorProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorProxy *AggregatorProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorProxy *AggregatorProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.contract.Transact(opts, method, params...)
}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_AggregatorProxy *AggregatorProxyCaller) Aggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "aggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_AggregatorProxy *AggregatorProxySession) Aggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.Aggregator(&_AggregatorProxy.CallOpts)
}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_AggregatorProxy *AggregatorProxyCallerSession) Aggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.Aggregator(&_AggregatorProxy.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorProxy *AggregatorProxyCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorProxy *AggregatorProxySession) Decimals() (uint8, error) {
	return _AggregatorProxy.Contract.Decimals(&_AggregatorProxy.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorProxy *AggregatorProxyCallerSession) Decimals() (uint8, error) {
	return _AggregatorProxy.Contract.Decimals(&_AggregatorProxy.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorProxy *AggregatorProxyCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorProxy *AggregatorProxySession) Description() (string, error) {
	return _AggregatorProxy.Contract.Description(&_AggregatorProxy.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorProxy *AggregatorProxyCallerSession) Description() (string, error) {
	return _AggregatorProxy.Contract.Description(&_AggregatorProxy.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256 answer)
func (_AggregatorProxy *AggregatorProxyCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256 answer)
func (_AggregatorProxy *AggregatorProxySession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetAnswer(&_AggregatorProxy.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256 answer)
func (_AggregatorProxy *AggregatorProxyCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetAnswer(&_AggregatorProxy.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "getRoundData", _roundId)

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
func (_AggregatorProxy *AggregatorProxySession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.GetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.GetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxyCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxySession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetTimestamp(&_AggregatorProxy.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxyCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AggregatorProxy.Contract.GetTimestamp(&_AggregatorProxy.CallOpts, _roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256 answer)
func (_AggregatorProxy *AggregatorProxyCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256 answer)
func (_AggregatorProxy *AggregatorProxySession) LatestAnswer() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestAnswer(&_AggregatorProxy.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256 answer)
func (_AggregatorProxy *AggregatorProxyCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestAnswer(&_AggregatorProxy.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256 roundId)
func (_AggregatorProxy *AggregatorProxyCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256 roundId)
func (_AggregatorProxy *AggregatorProxySession) LatestRound() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestRound(&_AggregatorProxy.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256 roundId)
func (_AggregatorProxy *AggregatorProxyCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestRound(&_AggregatorProxy.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestRoundData")

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
func (_AggregatorProxy *AggregatorProxySession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.LatestRoundData(&_AggregatorProxy.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.LatestRoundData(&_AggregatorProxy.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxyCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxySession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestTimestamp(&_AggregatorProxy.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxyCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorProxy.Contract.LatestTimestamp(&_AggregatorProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AggregatorProxy *AggregatorProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AggregatorProxy *AggregatorProxySession) Owner() (common.Address, error) {
	return _AggregatorProxy.Contract.Owner(&_AggregatorProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AggregatorProxy *AggregatorProxyCallerSession) Owner() (common.Address, error) {
	return _AggregatorProxy.Contract.Owner(&_AggregatorProxy.CallOpts)
}

// PhaseAggregators is a free data retrieval call binding the contract method 0xc1597304.
//
// Solidity: function phaseAggregators(uint16 ) view returns(address)
func (_AggregatorProxy *AggregatorProxyCaller) PhaseAggregators(opts *bind.CallOpts, arg0 uint16) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "phaseAggregators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PhaseAggregators is a free data retrieval call binding the contract method 0xc1597304.
//
// Solidity: function phaseAggregators(uint16 ) view returns(address)
func (_AggregatorProxy *AggregatorProxySession) PhaseAggregators(arg0 uint16) (common.Address, error) {
	return _AggregatorProxy.Contract.PhaseAggregators(&_AggregatorProxy.CallOpts, arg0)
}

// PhaseAggregators is a free data retrieval call binding the contract method 0xc1597304.
//
// Solidity: function phaseAggregators(uint16 ) view returns(address)
func (_AggregatorProxy *AggregatorProxyCallerSession) PhaseAggregators(arg0 uint16) (common.Address, error) {
	return _AggregatorProxy.Contract.PhaseAggregators(&_AggregatorProxy.CallOpts, arg0)
}

// PhaseId is a free data retrieval call binding the contract method 0x58303b10.
//
// Solidity: function phaseId() view returns(uint16)
func (_AggregatorProxy *AggregatorProxyCaller) PhaseId(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "phaseId")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// PhaseId is a free data retrieval call binding the contract method 0x58303b10.
//
// Solidity: function phaseId() view returns(uint16)
func (_AggregatorProxy *AggregatorProxySession) PhaseId() (uint16, error) {
	return _AggregatorProxy.Contract.PhaseId(&_AggregatorProxy.CallOpts)
}

// PhaseId is a free data retrieval call binding the contract method 0x58303b10.
//
// Solidity: function phaseId() view returns(uint16)
func (_AggregatorProxy *AggregatorProxyCallerSession) PhaseId() (uint16, error) {
	return _AggregatorProxy.Contract.PhaseId(&_AggregatorProxy.CallOpts)
}

// ProposedAggregator is a free data retrieval call binding the contract method 0xe8c4be30.
//
// Solidity: function proposedAggregator() view returns(address)
func (_AggregatorProxy *AggregatorProxyCaller) ProposedAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "proposedAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProposedAggregator is a free data retrieval call binding the contract method 0xe8c4be30.
//
// Solidity: function proposedAggregator() view returns(address)
func (_AggregatorProxy *AggregatorProxySession) ProposedAggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.ProposedAggregator(&_AggregatorProxy.CallOpts)
}

// ProposedAggregator is a free data retrieval call binding the contract method 0xe8c4be30.
//
// Solidity: function proposedAggregator() view returns(address)
func (_AggregatorProxy *AggregatorProxyCallerSession) ProposedAggregator() (common.Address, error) {
	return _AggregatorProxy.Contract.ProposedAggregator(&_AggregatorProxy.CallOpts)
}

// ProposedGetRoundData is a free data retrieval call binding the contract method 0x6001ac53.
//
// Solidity: function proposedGetRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCaller) ProposedGetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "proposedGetRoundData", _roundId)

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

// ProposedGetRoundData is a free data retrieval call binding the contract method 0x6001ac53.
//
// Solidity: function proposedGetRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxySession) ProposedGetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.ProposedGetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

// ProposedGetRoundData is a free data retrieval call binding the contract method 0x6001ac53.
//
// Solidity: function proposedGetRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCallerSession) ProposedGetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.ProposedGetRoundData(&_AggregatorProxy.CallOpts, _roundId)
}

// ProposedLatestRoundData is a free data retrieval call binding the contract method 0x8f6b4d91.
//
// Solidity: function proposedLatestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCaller) ProposedLatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "proposedLatestRoundData")

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

// ProposedLatestRoundData is a free data retrieval call binding the contract method 0x8f6b4d91.
//
// Solidity: function proposedLatestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxySession) ProposedLatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.ProposedLatestRoundData(&_AggregatorProxy.CallOpts)
}

// ProposedLatestRoundData is a free data retrieval call binding the contract method 0x8f6b4d91.
//
// Solidity: function proposedLatestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorProxy *AggregatorProxyCallerSession) ProposedLatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorProxy.Contract.ProposedLatestRoundData(&_AggregatorProxy.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorProxy *AggregatorProxyCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorProxy.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorProxy *AggregatorProxySession) Version() (*big.Int, error) {
	return _AggregatorProxy.Contract.Version(&_AggregatorProxy.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorProxy *AggregatorProxyCallerSession) Version() (*big.Int, error) {
	return _AggregatorProxy.Contract.Version(&_AggregatorProxy.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AggregatorProxy *AggregatorProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AggregatorProxy *AggregatorProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AcceptOwnership(&_AggregatorProxy.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AggregatorProxy *AggregatorProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AggregatorProxy.Contract.AcceptOwnership(&_AggregatorProxy.TransactOpts)
}

// ConfirmAggregator is a paid mutator transaction binding the contract method 0xa928c096.
//
// Solidity: function confirmAggregator(address _aggregator) returns()
func (_AggregatorProxy *AggregatorProxyTransactor) ConfirmAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "confirmAggregator", _aggregator)
}

// ConfirmAggregator is a paid mutator transaction binding the contract method 0xa928c096.
//
// Solidity: function confirmAggregator(address _aggregator) returns()
func (_AggregatorProxy *AggregatorProxySession) ConfirmAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ConfirmAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

// ConfirmAggregator is a paid mutator transaction binding the contract method 0xa928c096.
//
// Solidity: function confirmAggregator(address _aggregator) returns()
func (_AggregatorProxy *AggregatorProxyTransactorSession) ConfirmAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ConfirmAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

// ProposeAggregator is a paid mutator transaction binding the contract method 0xf8a2abd3.
//
// Solidity: function proposeAggregator(address _aggregator) returns()
func (_AggregatorProxy *AggregatorProxyTransactor) ProposeAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "proposeAggregator", _aggregator)
}

// ProposeAggregator is a paid mutator transaction binding the contract method 0xf8a2abd3.
//
// Solidity: function proposeAggregator(address _aggregator) returns()
func (_AggregatorProxy *AggregatorProxySession) ProposeAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ProposeAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

// ProposeAggregator is a paid mutator transaction binding the contract method 0xf8a2abd3.
//
// Solidity: function proposeAggregator(address _aggregator) returns()
func (_AggregatorProxy *AggregatorProxyTransactorSession) ProposeAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.ProposeAggregator(&_AggregatorProxy.TransactOpts, _aggregator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AggregatorProxy *AggregatorProxyTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AggregatorProxy *AggregatorProxySession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.TransferOwnership(&_AggregatorProxy.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AggregatorProxy *AggregatorProxyTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AggregatorProxy.Contract.TransferOwnership(&_AggregatorProxy.TransactOpts, _to)
}

// AggregatorProxyAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AggregatorProxy contract.
type AggregatorProxyAnswerUpdatedIterator struct {
	Event *AggregatorProxyAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *AggregatorProxyAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyAnswerUpdated)
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
		it.Event = new(AggregatorProxyAnswerUpdated)
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
func (it *AggregatorProxyAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorProxyAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorProxyAnswerUpdated represents a AnswerUpdated event raised by the AggregatorProxy contract.
type AggregatorProxyAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxyFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorProxyAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyAnswerUpdatedIterator{contract: _AggregatorProxy.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorProxy *AggregatorProxyFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorProxyAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorProxyAnswerUpdated)
				if err := _AggregatorProxy.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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
func (_AggregatorProxy *AggregatorProxyFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorProxyAnswerUpdated, error) {
	event := new(AggregatorProxyAnswerUpdated)
	if err := _AggregatorProxy.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorProxyNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AggregatorProxy contract.
type AggregatorProxyNewRoundIterator struct {
	Event *AggregatorProxyNewRound // Event containing the contract specifics and raw log

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
func (it *AggregatorProxyNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyNewRound)
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
		it.Event = new(AggregatorProxyNewRound)
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
func (it *AggregatorProxyNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorProxyNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorProxyNewRound represents a NewRound event raised by the AggregatorProxy contract.
type AggregatorProxyNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorProxy *AggregatorProxyFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorProxyNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyNewRoundIterator{contract: _AggregatorProxy.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorProxy *AggregatorProxyFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorProxyNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorProxyNewRound)
				if err := _AggregatorProxy.contract.UnpackLog(event, "NewRound", log); err != nil {
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
func (_AggregatorProxy *AggregatorProxyFilterer) ParseNewRound(log types.Log) (*AggregatorProxyNewRound, error) {
	event := new(AggregatorProxyNewRound)
	if err := _AggregatorProxy.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorProxyOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the AggregatorProxy contract.
type AggregatorProxyOwnershipTransferRequestedIterator struct {
	Event *AggregatorProxyOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *AggregatorProxyOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyOwnershipTransferRequested)
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
		it.Event = new(AggregatorProxyOwnershipTransferRequested)
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
func (it *AggregatorProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorProxyOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the AggregatorProxy contract.
type AggregatorProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AggregatorProxy *AggregatorProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AggregatorProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyOwnershipTransferRequestedIterator{contract: _AggregatorProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AggregatorProxy *AggregatorProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AggregatorProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorProxyOwnershipTransferRequested)
				if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_AggregatorProxy *AggregatorProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*AggregatorProxyOwnershipTransferRequested, error) {
	event := new(AggregatorProxyOwnershipTransferRequested)
	if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorProxyOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AggregatorProxy contract.
type AggregatorProxyOwnershipTransferredIterator struct {
	Event *AggregatorProxyOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AggregatorProxyOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorProxyOwnershipTransferred)
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
		it.Event = new(AggregatorProxyOwnershipTransferred)
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
func (it *AggregatorProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorProxyOwnershipTransferred represents a OwnershipTransferred event raised by the AggregatorProxy contract.
type AggregatorProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AggregatorProxy *AggregatorProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AggregatorProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorProxyOwnershipTransferredIterator{contract: _AggregatorProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AggregatorProxy *AggregatorProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AggregatorProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AggregatorProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AggregatorProxyOwnershipTransferred)
				if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AggregatorProxy *AggregatorProxyFilterer) ParseOwnershipTransferred(log types.Log) (*AggregatorProxyOwnershipTransferred, error) {
	event := new(AggregatorProxyOwnershipTransferred)
	if err := _AggregatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
