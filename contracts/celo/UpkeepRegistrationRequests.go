// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package celo

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = celo.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// UpkeepRegistrationRequestsMetaData contains all meta data concerning the UpkeepRegistrationRequests contract.
var UpkeepRegistrationRequestsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"windowStart\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"approvedInCurrentWindow\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"windowSizeInBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"allowedPerWindow\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516118ac3803806118ac8339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b031916331790556001600160601b031960609290921b9190911660805260025560805160601c6118146100986000398061071052806109fc5280610c9e528061117c528061141b52506118146000f3fe608060405234801561001057600080fd5b50600436106100995760003560e01c8063183310b31461009e5780631b6b6d23146101b95780635772ac92146101dd57806379ba509714610227578063850af0cb1461022f57806388b12d551461028e5780638da5cb5b146102d6578063a4c0ed36146102de578063c4110e5c14610361578063c4d252f5146104db578063f2fde38b146104f8575b600080fd5b6101b7600480360360c08110156100b457600080fd5b810190602081018135600160201b8111156100ce57600080fd5b8201836020820111156100e057600080fd5b803590602001918460018302840111600160201b8311171561010157600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092956001600160a01b03853581169663ffffffff6020880135169660408101359092169550919350909150608081019060600135600160201b81111561017957600080fd5b82018360208201111561018b57600080fd5b803590602001918460018302840111600160201b831117156101ac57600080fd5b91935091503561051e565b005b6101c161070e565b604080516001600160a01b039092168252519081900360200190f35b6101b7600480360360a08110156101f357600080fd5b50803515159063ffffffff6020820135169061ffff604082013516906001600160a01b036060820135169060800135610732565b6101b7610873565b610237610922565b60408051971515885263ffffffff909616602088015261ffff948516878701526001600160a01b03909316606087015260808601919091526001600160401b031660a08501521660c0830152519081900360e00190f35b6102ab600480360360208110156102a457600080fd5b503561099f565b604080516001600160a01b0390931683526001600160601b0390911660208301528051918290030190f35b6101c16109e2565b6101b7600480360360608110156102f457600080fd5b6001600160a01b0382351691602081013591810190606081016040820135600160201b81111561032357600080fd5b82018360208201111561033557600080fd5b803590602001918460018302840111600160201b8311171561035657600080fd5b5090925090506109f1565b6101b7600480360361010081101561037857600080fd5b810190602081018135600160201b81111561039257600080fd5b8201836020820111156103a457600080fd5b803590602001918460018302840111600160201b831117156103c557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561041757600080fd5b82018360208201111561042957600080fd5b803590602001918460018302840111600160201b8311171561044a57600080fd5b919390926001600160a01b03833581169363ffffffff6020820135169360408201359092169290608081019060600135600160201b81111561048b57600080fd5b82018360208201111561049d57600080fd5b803590602001918460018302840111600160201b831117156104be57600080fd5b919350915080356001600160601b0316906020013560ff16610c93565b6101b7600480360360208110156104f157600080fd5b503561103e565b6101b76004803603602081101561050e57600080fd5b50356001600160a01b0316611241565b6000546001600160a01b03163314610576576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6000818152600360209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b0316918301919091526105fb576040805162461bcd60e51b81526020600482015260116024820152701c995c5d595cdd081b9bdd08199bdd5b99607a1b604482015290519081900360640190fd5b6000878787878760405160200180866001600160a01b031681526020018563ffffffff168152602001846001600160a01b03168152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505096505050505050506040516020818303038152906040528051906020012090508083146106d6576040805162461bcd60e51b815260206004820152601d60248201527f6861736820616e64207061796c6f616420646f206e6f74206d61746368000000604482015290519081900360640190fd5b6000838152600360209081526040822091909155820151610703908a908a908a908a908a908a908a6112ea565b505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000546001600160a01b0316331461078a576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6040805160a0808201835287151580835261ffff8716602080850182905263ffffffff8a16858701819052600060608088018290526080978801919091526004805460ff1916861762ffff00191661010086021766ffffffff00000019166301000000840217600160381b600160881b03191690556002899055600580546001600160a01b0319166001600160a01b038c16908117909155885195865292850191909152838701929092529082015291820184905291517f421e8abed178b5e0b94e3f39d2eaa021143b1c90449f70e0f404c911098a1d53929181900390910190a15050505050565b6001546001600160a01b031633146108cb576040805162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015290519081900360640190fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a08101825260045460ff8116151580835261ffff610100830481166020850181905263ffffffff63010000008504169585018690526001600160401b03600160381b85041660608601819052600160781b9094049091166080909401849052600554600254929691946001600160a01b039091169391565b6000908152600360209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b03169290910182905291565b6000546001600160a01b031681565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610a64576040805162461bcd60e51b815260206004820152601360248201527226bab9ba103ab9b2902624a725903a37b5b2b760691b604482015290519081900360640190fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101516001600160e01b03198116633104439760e21b14610b04576040805162461bcd60e51b815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b8484848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060e4810151828114610b86576040805162461bcd60e51b815260206004820152600f60248201526e082dadeeadce840dad2e6dac2e8c6d608b1b604482015290519081900360640190fd5b600254881015610bd4576040805162461bcd60e51b8152602060048201526014602482015273125b9cdd59999a58da595b9d081c185e5b595b9d60621b604482015290519081900360640190fd5b6000306001600160a01b031688886040518083838082843760405192019450600093509091505080830381855af49150503d8060008114610c31576040519150601f19603f3d011682016040523d82523d6000602084013e610c36565b606091505b5050905080610c87576040805162461bcd60e51b8152602060048201526018602482015277155b98589b19481d1bc818dc99585d19481c995c5d595cdd60421b604482015290519081900360640190fd5b50505050505050505050565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610d06576040805162461bcd60e51b815260206004820152601360248201527226bab9ba103ab9b2902624a725903a37b5b2b760691b604482015290519081900360640190fd5b6001600160a01b038516610d59576040805162461bcd60e51b8152602060048201526015602482015274696e76616c69642061646d696e206164647265737360581b604482015290519081900360640190fd5b6000878787878760405160200180866001600160a01b031681526020018563ffffffff168152602001846001600160a01b03168152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505096505050505050506040516020818303038152906040528051906020012090508160ff16886001600160a01b0316827fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788e8e8e8d8d8d8d8d6040518080602001806020018863ffffffff168152602001876001600160a01b0316815260200180602001856001600160601b0316815260200184810384528c818151815260200191508051906020019080838360005b83811015610e87578181015183820152602001610e6f565b50505050905090810190601f168015610eb45780820380516001836020036101000a031916815260200191505b5084810383528a81526020018b8b80828437600083820152601f01601f191690910185810383528781526020019050878780828437600083820152604051601f909101601f19169092018290039d50909b505050505050505050505050a46040805160a08101825260045460ff811615801580845261ffff61010084048116602086015263ffffffff6301000000850416958501959095526001600160401b03600160381b8404166060850152600160781b90920490931660808301529091610f815750610f81816115e1565b15610fa457610f8f81611615565b610f9f8c8a8a8a8a8a8a896112ea565b611030565b600082815260036020526040812054610fcd90600160a01b90046001600160601b0316866116b8565b6040805180820182526001600160a01b038b811682526001600160601b03938416602080840191825260008981526003909152939093209151825493516001600160a01b03199094169082161716600160a01b9290931691909102919091179055505b505050505050505050505050565b6000818152600360209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b03169183019190915233148061109557506000546001600160a01b031633145b6110e6576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792061646d696e202f206f776e65722063616e2063616e63656c000000604482015290519081900360640190fd5b80516001600160a01b0316611136576040805162461bcd60e51b81526020600482015260116024820152701c995c5d595cdd081b9bdd08199bdd5b99607a1b604482015290519081900360640190fd5b600082815260036020908152604080832083905583820151815163a9059cbb60e01b81523360048201526001600160601b03909116602482015290516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb93604480850194919392918390030190829087803b1580156111c357600080fd5b505af11580156111d7573d6000803e3d6000fd5b505050506040513d60208110156111ed57600080fd5b505161123d576040805162461bcd60e51b815260206004820152601a60248201527913125392c81d1bdad95b881d1c985b9cd9995c8819985a5b195960321b604482015290519081900360640190fd5b5050565b6000546001600160a01b03163314611299576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60055460405163da5c674160e01b81526001600160a01b038981166004830190815263ffffffff8a1660248401528882166044840152608060648401908152608484018890529190931692600092849263da5c6741928d928d928d928d928d929060a401848480828437600081840152601f19601f8201169050808301925050509650505050505050602060405180830381600087803b15801561138d57600080fd5b505af11580156113a1573d6000803e3d6000fd5b505050506040513d60208110156113b757600080fd5b505160408051602080820184905282518083038201815282840193849052630200057560e51b9093526001600160a01b03868116604484019081526001600160601b038a166064850152606060848501908152855160a486015285519697506000967f000000000000000000000000000000000000000000000000000000000000000090931695634000aea0958a958d959294939260c490920191908501908083838d5b8381101561147357818101518382015260200161145b565b50505050905090810190601f1680156114a05780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b1580156114c157600080fd5b505af11580156114d5573d6000803e3d6000fd5b505050506040513d60208110156114eb57600080fd5b5051905080611539576040805162461bcd60e51b815260206004820152601560248201527406661696c656420746f2066756e642075706b65657605c1b604482015290519081900360640190fd5b81847fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b8d6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561159a578181015183820152602001611582565b50505050905090810190601f1680156115c75780820380516001836020036101000a031916815260200191505b509250505060405180910390a35050505050505050505050565b60006115ec82611723565b816020015161ffff16826080015161ffff16101561160c57506001611610565b5060005b919050565b60808101805161ffff600190910181169182905282516004805460208601516040870151606090970151600160781b90960261ffff60781b196001600160401b03909716600160381b0267ffffffffffffffff60381b1963ffffffff90991663010000000266ffffffff00000019939097166101000262ffff001996151560ff199095169490941795909516929092171693909317949094161791909116179055565b60008282016001600160601b03808516908216101561171c576040805162461bcd60e51b815260206004820152601b60248201527a536166654d6174683a206164646974696f6e206f766572666c6f7760281b604482015290519081900360640190fd5b9392505050565b600081606001516001600160401b031643039050816040015163ffffffff16816001600160401b03161061123d5750436001600160401b03166060820181905260006080830152815160048054602085015160409095015160ff199091169215159290921762ffff00191661010061ffff909516949094029390931766ffffffff0000001916630100000063ffffffff909216919091021767ffffffffffffffff60381b1916600160381b9091021761ffff60781b1916905556fea26469706673582212204409d297842855696cc7bd1d219497ab6aad8bd9b87c766499d6bd2d913a14c064736f6c63430007060033",
}

// UpkeepRegistrationRequestsABI is the input ABI used to generate the binding from.
// Deprecated: Use UpkeepRegistrationRequestsMetaData.ABI instead.
var UpkeepRegistrationRequestsABI = UpkeepRegistrationRequestsMetaData.ABI

// UpkeepRegistrationRequestsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpkeepRegistrationRequestsMetaData.Bin instead.
var UpkeepRegistrationRequestsBin = UpkeepRegistrationRequestsMetaData.Bin

// DeployUpkeepRegistrationRequests deploys a new Ethereum contract, binding an instance of UpkeepRegistrationRequests to it.
func DeployUpkeepRegistrationRequests(auth *bind.TransactOpts, backend bind.ContractBackend, LINKAddress common.Address, minimumLINKJuels *big.Int) (common.Address, *types.Transaction, *UpkeepRegistrationRequests, error) {
	parsed, err := UpkeepRegistrationRequestsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepRegistrationRequestsBin), backend, LINKAddress, minimumLINKJuels)
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
