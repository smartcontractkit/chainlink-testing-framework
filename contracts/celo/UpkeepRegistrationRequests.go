// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package celo

import (
	"math/big"
	"strings"

	celo "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = celo.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// UpkeepRegistrationRequestsABI is the input ABI used to generate the binding from.
const UpkeepRegistrationRequestsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"windowStart\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"approvedInCurrentWindow\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// UpkeepRegistrationRequestsBin is the compiled bytecode used for deploying new contracts.
var UpkeepRegistrationRequestsBin = "0x60a060405234801561001057600080fd5b5060405161204f38038061204f8339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b031916331790556001600160601b031960609290921b9190911660805260025560805160601c611fb76100986000398061086a5280610cbe528061105352806117035280611a635250611fb76000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c806388b12d5511610081578063c4110e5c1161005b578063c4110e5c146103eb578063c4d252f51461057d578063f2fde38b1461059a576100c9565b806388b12d55146102f75780638da5cb5b14610351578063a4c0ed3614610359576100c9565b80635772ac92116100b25780635772ac921461022b57806379ba509714610282578063850af0cb1461028a576100c9565b8063183310b3146100ce5780631b6b6d23146101fa575b600080fd5b6101f8600480360360c08110156100e457600080fd5b8101906020810181356401000000008111156100ff57600080fd5b82018360208201111561011157600080fd5b8035906020019184600183028401116401000000008311171561013357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929573ffffffffffffffffffffffffffffffffffffffff853581169663ffffffff60208801351696604081013590921695509193509091506080810190606001356401000000008111156101b957600080fd5b8201836020820111156101cb57600080fd5b803590602001918460018302840111640100000000831117156101ed57600080fd5b9193509150356105cd565b005b610202610868565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101f8600480360360a081101561024157600080fd5b50803515159063ffffffff6020820135169061ffff6040820135169073ffffffffffffffffffffffffffffffffffffffff606082013516906080013561088c565b6101f8610a87565b610292610b89565b60408051971515885263ffffffff909616602088015261ffff9485168787015273ffffffffffffffffffffffffffffffffffffffff9093166060870152608086019190915267ffffffffffffffff1660a08501521660c0830152519081900360e00190f35b6103146004803603602081101561030d57600080fd5b5035610c24565b6040805173ffffffffffffffffffffffffffffffffffffffff90931683526bffffffffffffffffffffffff90911660208301528051918290030190f35b610202610c8a565b6101f86004803603606081101561036f57600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516916020810135918101906060810160408201356401000000008111156103ac57600080fd5b8201836020820111156103be57600080fd5b803590602001918460018302840111640100000000831117156103e057600080fd5b509092509050610ca6565b6101f8600480360361010081101561040257600080fd5b81019060208101813564010000000081111561041d57600080fd5b82018360208201111561042f57600080fd5b8035906020019184600183028401116401000000008311171561045157600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092959493602081019350359150506401000000008111156104a457600080fd5b8201836020820111156104b657600080fd5b803590602001918460018302840111640100000000831117156104d857600080fd5b9193909273ffffffffffffffffffffffffffffffffffffffff833581169363ffffffff602082013516936040820135909216929060808101906060013564010000000081111561052757600080fd5b82018360208201111561053957600080fd5b8035906020019184600183028401116401000000008311171561055b57600080fd5b919350915080356bffffffffffffffffffffffff16906020013560ff1661103b565b6101f86004803603602081101561059357600080fd5b503561151d565b6101f8600480360360208110156105b057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166117e5565b60005473ffffffffffffffffffffffffffffffffffffffff16331461065357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60008181526003602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169183019190915261072157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015290519081900360640190fd5b60008787878787604051602001808673ffffffffffffffffffffffffffffffffffffffff1681526020018563ffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001806020018281038252848482818152602001925080828437600081840152601f19601f820116905080830192505050965050505050505060405160208183030381529060405280519060200120905080831461083057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6861736820616e64207061796c6f616420646f206e6f74206d61746368000000604482015290519081900360640190fd5b600083815260036020908152604082209190915582015161085d908a908a908a908a908a908a908a6118e1565b505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60005473ffffffffffffffffffffffffffffffffffffffff16331461091257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6040805160a0808201835287151580835261ffff8716602080850182905263ffffffff8a1685870181905260006060808801829052608097880191909152600480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff166101008602177fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff1663010000008402177fffffffffffffffffffffffffffffff00000000000000000000ffffffffffffff1690556002899055600580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8c16908117909155885195865292850191909152838701929092529082015291820184905291517f421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53929181900390910190a15050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b0d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a08101825260045460ff8116151580835261ffff610100830481166020850181905263ffffffff630100000085041695850186905267ffffffffffffffff670100000000000000850416606086018190526f0100000000000000000000000000000090940490911660809094018490526005546002549296919473ffffffffffffffffffffffffffffffffffffffff9091169391565b60009081526003602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169290910182905291565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610d4a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101517fffffffff0000000000000000000000000000000000000000000000000000000081167fc4110e5c0000000000000000000000000000000000000000000000000000000014610e3557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b8484848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060e4810151828114610edf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f416d6f756e74206d69736d617463680000000000000000000000000000000000604482015290519081900360640190fd5b600254881015610f5057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e73756666696369656e74207061796d656e74000000000000000000000000604482015290519081900360640190fd5b60003073ffffffffffffffffffffffffffffffffffffffff1688886040518083838082843760405192019450600093509091505080830381855af49150503d8060008114610fba576040519150601f19603f3d011682016040523d82523d6000602084013e610fbf565b606091505b505090508061102f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f556e61626c6520746f2063726561746520726571756573740000000000000000604482015290519081900360640190fd5b50505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110df57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff851661116157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f696e76616c69642061646d696e20616464726573730000000000000000000000604482015290519081900360640190fd5b60008787878787604051602001808673ffffffffffffffffffffffffffffffffffffffff1681526020018563ffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505096505050505050506040516020818303038152906040528051906020012090508160ff168873ffffffffffffffffffffffffffffffffffffffff16827fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788e8e8e8d8d8d8d8d6040518080602001806020018863ffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff16815260200180602001856bffffffffffffffffffffffff16815260200184810384528c818151815260200191508051906020019080838360005b838110156112c85781810151838201526020016112b0565b50505050905090810190601f1680156112f55780820380516001836020036101000a031916815260200191505b5084810383528a81526020018b8b80828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528781526020019050878780828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039d50909b505050505050505050505050a46040805160a08101825260045460ff811615801580845261ffff61010084048116602086015263ffffffff63010000008504169585019590955267ffffffffffffffff67010000000000000084041660608501526f010000000000000000000000000000009092049093166080830152909161140f575061140f81611c4b565b156114325761141d81611c7f565b61142d8c8a8a8a8a8a8a896118e1565b61150f565b600082815260036020526040812054611471907401000000000000000000000000000000000000000090046bffffffffffffffffffffffff1686611db3565b60408051808201825273ffffffffffffffffffffffffffffffffffffffff8b811682526bffffffffffffffffffffffff938416602080840191825260008981526003909152939093209151825493517fffffffffffffffffffffffff00000000000000000000000000000000000000009094169082161716740100000000000000000000000000000000000000009290931691909102919091179055505b505050505050505050505050565b60008181526003602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff16918301919091523314806115a4575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61160f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792061646d696e202f206f776e65722063616e2063616e63656c000000604482015290519081900360640190fd5b805173ffffffffffffffffffffffffffffffffffffffff1661169257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015290519081900360640190fd5b60008281526003602090815260408083208390558382015181517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff9091166024820152905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb93604480850194919392918390030190829087803b15801561174a57600080fd5b505af115801561175e573d6000803e3d6000fd5b505050506040513d602081101561177457600080fd5b50516117e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4c494e4b20746f6b656e207472616e73666572206661696c6564000000000000604482015290519081900360640190fd5b5050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461186b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6005546040517fda5c674100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8981166004830190815263ffffffff8a1660248401528882166044840152608060648401908152608484018890529190931692600092849263da5c6741928d928d928d928d928d929060a401848480828437600081840152601f19601f8201169050808301925050509650505050505050602060405180830381600087803b1580156119aa57600080fd5b505af11580156119be573d6000803e3d6000fd5b505050506040513d60208110156119d457600080fd5b5051604080516020808201849052825180830382018152828401938490527f4000aea00000000000000000000000000000000000000000000000000000000090935273ffffffffffffffffffffffffffffffffffffffff868116604484019081526bffffffffffffffffffffffff8a166064850152606060848501908152855160a486015285519697506000967f000000000000000000000000000000000000000000000000000000000000000090931695634000aea0958a958d959294939260c490920191908501908083838d5b83811015611abb578181015183820152602001611aa3565b50505050905090810190601f168015611ae85780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015611b0957600080fd5b505af1158015611b1d573d6000803e3d6000fd5b505050506040513d6020811015611b3357600080fd5b5051905080611ba357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f6661696c656420746f2066756e642075706b6565700000000000000000000000604482015290519081900360640190fd5b81847fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b8d6040518080602001828103825283818151815260200191508051906020019080838360005b83811015611c04578181015183820152602001611bec565b50505050905090810190601f168015611c315780820380516001836020036101000a031916815260200191505b509250505060405180910390a35050505050505050505050565b6000611c5682611e3f565b816020015161ffff16826080015161ffff161015611c7657506001611c7a565b5060005b919050565b60808101805161ffff6001909101811691829052825160048054602086015160408701516060909701516f010000000000000000000000000000009096027fffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffff67ffffffffffffffff909716670100000000000000027fffffffffffffffffffffffffffffffffff0000000000000000ffffffffffffff63ffffffff9099166301000000027fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff93909716610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff9615157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009095169490941795909516929092171693909317949094161791909116179055565b60008282016bffffffffffffffffffffffff8085169082161015611e3857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6000816060015167ffffffffffffffff1643039050816040015163ffffffff168167ffffffffffffffff16106117e157504367ffffffffffffffff16606082018190526000608083015281516004805460208501516040909501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00909116921515929092177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff1661010061ffff90951694909402939093177fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff16630100000063ffffffff90921691909102177fffffffffffffffffffffffffffffffffff0000000000000000ffffffffffffff16670100000000000000909102177fffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffff16905556fea2646970667358221220e43784eaa4830088a2c893c1ad168455f9424f05f0440d214c9d69d18f92c90364736f6c63430007060033"

// DeployUpkeepRegistrationRequests deploys a new Ethereum contract, binding an instance of UpkeepRegistrationRequests to it.
func DeployUpkeepRegistrationRequests(auth *bind.TransactOpts, backend bind.ContractBackend, LINKAddress common.Address, minimumLINKJuels *big.Int) (common.Address, *types.Transaction, *UpkeepRegistrationRequests, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepRegistrationRequestsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UpkeepRegistrationRequestsBin), backend, LINKAddress, minimumLINKJuels)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepRegistrationRequests{UpkeepRegistrationRequestsCaller: UpkeepRegistrationRequestsCaller{contract: contract}, UpkeepRegistrationRequestsTransactor: UpkeepRegistrationRequestsTransactor{contract: contract}, UpkeepRegistrationRequestsFilterer: UpkeepRegistrationRequestsFilterer{contract: contract}}, nil
}

// UpkeepRegistrationRequests is an auto generated Go binding around an Ethereum contract.
type UpkeepRegistrationRequests struct {
	UpkeepRegistrationRequestsCaller     // Read-only binding to the contract
	UpkeepRegistrationRequestsTransactor // Write-only binding to the contract
	UpkeepRegistrationRequestsFilterer   // Log filterer for contract events
}

// UpkeepRegistrationRequestsCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpkeepRegistrationRequestsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepRegistrationRequestsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpkeepRegistrationRequestsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepRegistrationRequestsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpkeepRegistrationRequestsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepRegistrationRequestsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpkeepRegistrationRequestsSession struct {
	Contract     *UpkeepRegistrationRequests // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// UpkeepRegistrationRequestsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpkeepRegistrationRequestsCallerSession struct {
	Contract *UpkeepRegistrationRequestsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// UpkeepRegistrationRequestsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpkeepRegistrationRequestsTransactorSession struct {
	Contract     *UpkeepRegistrationRequestsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// UpkeepRegistrationRequestsRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpkeepRegistrationRequestsRaw struct {
	Contract *UpkeepRegistrationRequests // Generic contract binding to access the raw methods on
}

// UpkeepRegistrationRequestsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpkeepRegistrationRequestsCallerRaw struct {
	Contract *UpkeepRegistrationRequestsCaller // Generic read-only contract binding to access the raw methods on
}

// UpkeepRegistrationRequestsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpkeepRegistrationRequestsTransactorRaw struct {
	Contract *UpkeepRegistrationRequestsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpkeepRegistrationRequests creates a new instance of UpkeepRegistrationRequests, bound to a specific deployed contract.
func NewUpkeepRegistrationRequests(address common.Address, backend bind.ContractBackend) (*UpkeepRegistrationRequests, error) {
	contract, err := bindUpkeepRegistrationRequests(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequests{UpkeepRegistrationRequestsCaller: UpkeepRegistrationRequestsCaller{contract: contract}, UpkeepRegistrationRequestsTransactor: UpkeepRegistrationRequestsTransactor{contract: contract}, UpkeepRegistrationRequestsFilterer: UpkeepRegistrationRequestsFilterer{contract: contract}}, nil
}

// NewUpkeepRegistrationRequestsCaller creates a new read-only instance of UpkeepRegistrationRequests, bound to a specific deployed contract.
func NewUpkeepRegistrationRequestsCaller(address common.Address, caller bind.ContractCaller) (*UpkeepRegistrationRequestsCaller, error) {
	contract, err := bindUpkeepRegistrationRequests(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsCaller{contract: contract}, nil
}

// NewUpkeepRegistrationRequestsTransactor creates a new write-only instance of UpkeepRegistrationRequests, bound to a specific deployed contract.
func NewUpkeepRegistrationRequestsTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepRegistrationRequestsTransactor, error) {
	contract, err := bindUpkeepRegistrationRequests(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsTransactor{contract: contract}, nil
}

// NewUpkeepRegistrationRequestsFilterer creates a new log filterer instance of UpkeepRegistrationRequests, bound to a specific deployed contract.
func NewUpkeepRegistrationRequestsFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepRegistrationRequestsFilterer, error) {
	contract, err := bindUpkeepRegistrationRequests(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsFilterer{contract: contract}, nil
}

// bindUpkeepRegistrationRequests binds a generic wrapper to an already deployed contract.
func bindUpkeepRegistrationRequests(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepRegistrationRequestsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepRegistrationRequests.Contract.UpkeepRegistrationRequestsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.UpkeepRegistrationRequestsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.UpkeepRegistrationRequestsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepRegistrationRequests.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.contract.Transact(opts, method, params...)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) LINK() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.LINK(&_UpkeepRegistrationRequests.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) LINK() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.LINK(&_UpkeepRegistrationRequests.CallOpts)
}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) GetPendingRequest(opts *bind.CallOpts, hash [32]byte) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "getPendingRequest", hash)

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _UpkeepRegistrationRequests.Contract.GetPendingRequest(&_UpkeepRegistrationRequests.CallOpts, hash)
}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _UpkeepRegistrationRequests.Contract.GetPendingRequest(&_UpkeepRegistrationRequests.CallOpts, hash)
}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels, uint64 windowStart, uint16 approvedInCurrentWindow)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) GetRegistrationConfig(opts *bind.CallOpts) (struct {
	Enabled                 bool
	WindowSizeInBlocks      uint32
	AllowedPerWindow        uint16
	KeeperRegistry          common.Address
	MinLINKJuels            *big.Int
	WindowStart             uint64
	ApprovedInCurrentWindow uint16
}, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "getRegistrationConfig")

	outstruct := new(struct {
		Enabled                 bool
		WindowSizeInBlocks      uint32
		AllowedPerWindow        uint16
		KeeperRegistry          common.Address
		MinLINKJuels            *big.Int
		WindowStart             uint64
		ApprovedInCurrentWindow uint16
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Enabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.WindowSizeInBlocks = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.AllowedPerWindow = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.KeeperRegistry = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.MinLINKJuels = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.WindowStart = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.ApprovedInCurrentWindow = *abi.ConvertType(out[6], new(uint16)).(*uint16)

	return *outstruct, err

}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels, uint64 windowStart, uint16 approvedInCurrentWindow)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) GetRegistrationConfig() (struct {
	Enabled                 bool
	WindowSizeInBlocks      uint32
	AllowedPerWindow        uint16
	KeeperRegistry          common.Address
	MinLINKJuels            *big.Int
	WindowStart             uint64
	ApprovedInCurrentWindow uint16
}, error) {
	return _UpkeepRegistrationRequests.Contract.GetRegistrationConfig(&_UpkeepRegistrationRequests.CallOpts)
}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels, uint64 windowStart, uint16 approvedInCurrentWindow)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) GetRegistrationConfig() (struct {
	Enabled                 bool
	WindowSizeInBlocks      uint32
	AllowedPerWindow        uint16
	KeeperRegistry          common.Address
	MinLINKJuels            *big.Int
	WindowStart             uint64
	ApprovedInCurrentWindow uint16
}, error) {
	return _UpkeepRegistrationRequests.Contract.GetRegistrationConfig(&_UpkeepRegistrationRequests.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpkeepRegistrationRequests.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Owner() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.Owner(&_UpkeepRegistrationRequests.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsCallerSession) Owner() (common.Address, error) {
	return _UpkeepRegistrationRequests.Contract.Owner(&_UpkeepRegistrationRequests.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) AcceptOwnership() (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.AcceptOwnership(&_UpkeepRegistrationRequests.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.AcceptOwnership(&_UpkeepRegistrationRequests.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x183310b3.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes32 hash) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "approve", name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

// Approve is a paid mutator transaction binding the contract method 0x183310b3.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes32 hash) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Approve(&_UpkeepRegistrationRequests.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

// Approve is a paid mutator transaction binding the contract method 0x183310b3.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes32 hash) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Approve(&_UpkeepRegistrationRequests.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "cancel", hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Cancel(&_UpkeepRegistrationRequests.TransactOpts, hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Cancel(&_UpkeepRegistrationRequests.TransactOpts, hash)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 amount, bytes data) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 amount, bytes data) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.OnTokenTransfer(&_UpkeepRegistrationRequests.TransactOpts, arg0, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address , uint256 amount, bytes data) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.OnTokenTransfer(&_UpkeepRegistrationRequests.TransactOpts, arg0, amount, data)
}

// Register is a paid mutator transaction binding the contract method 0xc4110e5c.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 source) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "register", name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

// Register is a paid mutator transaction binding the contract method 0xc4110e5c.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 source) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Register(&_UpkeepRegistrationRequests.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

// Register is a paid mutator transaction binding the contract method 0xc4110e5c.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 source) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.Register(&_UpkeepRegistrationRequests.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0x5772ac92.
//
// Solidity: function setRegistrationConfig(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) SetRegistrationConfig(opts *bind.TransactOpts, enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "setRegistrationConfig", enabled, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0x5772ac92.
//
// Solidity: function setRegistrationConfig(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) SetRegistrationConfig(enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.SetRegistrationConfig(&_UpkeepRegistrationRequests.TransactOpts, enabled, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0x5772ac92.
//
// Solidity: function setRegistrationConfig(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) SetRegistrationConfig(enabled bool, windowSizeInBlocks uint32, allowedPerWindow uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.SetRegistrationConfig(&_UpkeepRegistrationRequests.TransactOpts, enabled, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.TransferOwnership(&_UpkeepRegistrationRequests.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _UpkeepRegistrationRequests.Contract.TransferOwnership(&_UpkeepRegistrationRequests.TransactOpts, _to)
}

// UpkeepRegistrationRequestsConfigChangedIterator is returned from FilterConfigChanged and is used to iterate over the raw logs and unpacked data for ConfigChanged events raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsConfigChangedIterator struct {
	Event *UpkeepRegistrationRequestsConfigChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepRegistrationRequestsConfigChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsConfigChanged)
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
		it.Event = new(UpkeepRegistrationRequestsConfigChanged)
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
func (it *UpkeepRegistrationRequestsConfigChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepRegistrationRequestsConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepRegistrationRequestsConfigChanged represents a ConfigChanged event raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsConfigChanged struct {
	Enabled            bool
	WindowSizeInBlocks uint32
	AllowedPerWindow   uint16
	KeeperRegistry     common.Address
	MinLINKJuels       *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterConfigChanged is a free log retrieval operation binding the contract event 0x421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53.
//
// Solidity: event ConfigChanged(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*UpkeepRegistrationRequestsConfigChangedIterator, error) {

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsConfigChangedIterator{contract: _UpkeepRegistrationRequests.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

// WatchConfigChanged is a free log subscription operation binding the contract event 0x421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53.
//
// Solidity: event ConfigChanged(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsConfigChanged) (event.Subscription, error) {

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpkeepRegistrationRequestsConfigChanged)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

// ParseConfigChanged is a log parse operation binding the contract event 0x421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53.
//
// Solidity: event ConfigChanged(bool enabled, uint32 windowSizeInBlocks, uint16 allowedPerWindow, address keeperRegistry, uint256 minLINKJuels)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseConfigChanged(log types.Log) (*UpkeepRegistrationRequestsConfigChanged, error) {
	event := new(UpkeepRegistrationRequestsConfigChanged)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpkeepRegistrationRequestsOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsOwnershipTransferRequestedIterator struct {
	Event *UpkeepRegistrationRequestsOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepRegistrationRequestsOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsOwnershipTransferRequested)
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
		it.Event = new(UpkeepRegistrationRequestsOwnershipTransferRequested)
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
func (it *UpkeepRegistrationRequestsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepRegistrationRequestsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepRegistrationRequestsOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpkeepRegistrationRequestsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsOwnershipTransferRequestedIterator{contract: _UpkeepRegistrationRequests.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpkeepRegistrationRequestsOwnershipTransferRequested)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseOwnershipTransferRequested(log types.Log) (*UpkeepRegistrationRequestsOwnershipTransferRequested, error) {
	event := new(UpkeepRegistrationRequestsOwnershipTransferRequested)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpkeepRegistrationRequestsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsOwnershipTransferredIterator struct {
	Event *UpkeepRegistrationRequestsOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepRegistrationRequestsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsOwnershipTransferred)
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
		it.Event = new(UpkeepRegistrationRequestsOwnershipTransferred)
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
func (it *UpkeepRegistrationRequestsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepRegistrationRequestsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepRegistrationRequestsOwnershipTransferred represents a OwnershipTransferred event raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpkeepRegistrationRequestsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsOwnershipTransferredIterator{contract: _UpkeepRegistrationRequests.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpkeepRegistrationRequestsOwnershipTransferred)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseOwnershipTransferred(log types.Log) (*UpkeepRegistrationRequestsOwnershipTransferred, error) {
	event := new(UpkeepRegistrationRequestsOwnershipTransferred)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpkeepRegistrationRequestsRegistrationApprovedIterator is returned from FilterRegistrationApproved and is used to iterate over the raw logs and unpacked data for RegistrationApproved events raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsRegistrationApprovedIterator struct {
	Event *UpkeepRegistrationRequestsRegistrationApproved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepRegistrationRequestsRegistrationApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsRegistrationApproved)
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
		it.Event = new(UpkeepRegistrationRequestsRegistrationApproved)
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
func (it *UpkeepRegistrationRequestsRegistrationApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepRegistrationRequestsRegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepRegistrationRequestsRegistrationApproved represents a RegistrationApproved event raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsRegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRegistrationApproved is a free log retrieval operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*UpkeepRegistrationRequestsRegistrationApprovedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsRegistrationApprovedIterator{contract: _UpkeepRegistrationRequests.contract, event: "RegistrationApproved", logs: logs, sub: sub}, nil
}

// WatchRegistrationApproved is a free log subscription operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpkeepRegistrationRequestsRegistrationApproved)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
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

// ParseRegistrationApproved is a log parse operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseRegistrationApproved(log types.Log) (*UpkeepRegistrationRequestsRegistrationApproved, error) {
	event := new(UpkeepRegistrationRequestsRegistrationApproved)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpkeepRegistrationRequestsRegistrationRequestedIterator is returned from FilterRegistrationRequested and is used to iterate over the raw logs and unpacked data for RegistrationRequested events raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsRegistrationRequestedIterator struct {
	Event *UpkeepRegistrationRequestsRegistrationRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  celo.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepRegistrationRequestsRegistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepRegistrationRequestsRegistrationRequested)
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
		it.Event = new(UpkeepRegistrationRequestsRegistrationRequested)
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
func (it *UpkeepRegistrationRequestsRegistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepRegistrationRequestsRegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepRegistrationRequestsRegistrationRequested represents a RegistrationRequested event raised by the UpkeepRegistrationRequests contract.
type UpkeepRegistrationRequestsRegistrationRequested struct {
	Hash           [32]byte
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	Amount         *big.Int
	Source         uint8
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRegistrationRequested is a free log retrieval operation binding the contract event 0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 indexed source)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*UpkeepRegistrationRequestsRegistrationRequestedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepRegistrationRequestsRegistrationRequestedIterator{contract: _UpkeepRegistrationRequests.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

// WatchRegistrationRequested is a free log subscription operation binding the contract event 0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 indexed source)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *UpkeepRegistrationRequestsRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _UpkeepRegistrationRequests.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpkeepRegistrationRequestsRegistrationRequested)
				if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

// ParseRegistrationRequested is a log parse operation binding the contract event 0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 indexed source)
func (_UpkeepRegistrationRequests *UpkeepRegistrationRequestsFilterer) ParseRegistrationRequested(log types.Log) (*UpkeepRegistrationRequestsRegistrationRequested, error) {
	event := new(UpkeepRegistrationRequestsRegistrationRequested)
	if err := _UpkeepRegistrationRequests.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
