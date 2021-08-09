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

// VRFCoordinatorABI is the input ABI used to generate the binding from.
const VRFCoordinatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_blockHashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"NewServiceAgreement\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestID\",\"type\":\"bytes32\"}],\"name\":\"RandomnessRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"output\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequestFulfilled\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PRESEED_OFFSET\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROOF_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PUBLIC_KEY_OFFSET\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackContract\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"randomnessFee\",\"type\":\"uint96\"},{\"internalType\":\"bytes32\",\"name\":\"seedAndBlockNum\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"fulfillRandomnessRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"_publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"_publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes32\",\"name\":\"_jobID\",\"type\":\"bytes32\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"serviceAgreements\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"vRFOracle\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"fee\",\"type\":\"uint96\"},{\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// VRFCoordinatorBin is the compiled bytecode used for deploying new contracts.
var VRFCoordinatorBin = "0x608060405234801561001057600080fd5b506040516121fa3803806121fa8339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b03191633178082556040516001600160a01b039190911691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3600180546001600160a01b039384166001600160a01b0319918216179091556002805492909316911617905561213a806100c06000396000f3fe608060405234801561001057600080fd5b50600436106100e95760003560e01c8063a4c0ed361161008c578063d834020911610066578063d834020914610393578063e911439c146103ca578063f2fde38b146103d2578063f3fef3a3146103f8576100e9565b8063a4c0ed3614610285578063b415f4f514610340578063caf70c4a14610348576100e9565b806375d35070116100c857806375d35070146102205780638aa7927b1461023d5780638da5cb5b146102455780638f32d59b14610269576100e9565b80626f6ad0146100ee57806321f36509146101265780635e1c105914610178575b600080fd5b6101146004803603602081101561010457600080fd5b50356001600160a01b0316610424565b60408051918252519081900360200190f35b6101436004803603602081101561013c57600080fd5b5035610436565b604080516001600160a01b0390941684526bffffffffffffffffffffffff909216602084015282820152519081900360600190f35b61021e6004803603602081101561018e57600080fd5b8101906020810181356401000000008111156101a957600080fd5b8201836020820111156101bb57600080fd5b803590602001918460018302840111640100000000831117156101dd57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610470945050505050565b005b6101436004803603602081101561023657600080fd5b503561055e565b610114610598565b61024d61059d565b604080516001600160a01b039092168252519081900360200190f35b6102716105ac565b604080519115158252519081900360200190f35b61021e6004803603606081101561029b57600080fd5b6001600160a01b03823516916020810135918101906060810160408201356401000000008111156102cb57600080fd5b8201836020820111156102dd57600080fd5b803590602001918460018302840111640100000000831117156102ff57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506105bd945050505050565b610114610655565b6101146004803603604081101561035e57600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525091945061065a9350505050565b61021e600480360360a08110156103a957600080fd5b508035906001600160a01b03602082013516906040810190608001356106b0565b6101146108fb565b61021e600480360360208110156103e857600080fd5b50356001600160a01b0316610901565b61021e6004803603604081101561040e57600080fd5b506001600160a01b038135169060200135610966565b60056020526000908152604090205481565b600360205260009081526040902080546001909101546001600160a01b03821691600160a01b90046bffffffffffffffffffffffff169083565b600061047a611fc7565b60008061048685610aa2565b600084815260046020908152604080832054828701516001600160a01b039091168085526005909352922054959950939750919550935090916104dc916bffffffffffffffffffffffff1663ffffffff610d9616565b6001600160a01b0382166000908152600560209081526040808320939093558582526003905290812081815560010155835161051b9084908490610df9565b604080518481526020810184905281517fa2e7a402243ebda4a69ceeb3dfb682943b7a9b3ac66d6eefa8db65894009611c929181900390910190a1505050505050565b600460205260009081526040902080546001909101546001600160a01b03821691600160a01b90046bffffffffffffffffffffffff169083565b602081565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b6001546001600160a01b0316331461061c576040805162461bcd60e51b815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b60008082806020019051604081101561063457600080fd5b508051602090910151909250905061064e82828688610f71565b5050505050565b60e081565b6000816040516020018082600260200280838360005b83811015610688578181015183820152602001610670565b505050509050019150506040516020818303038152906040528051906020012090505b919050565b6106b86105ac565b610709576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b60408051808201825260009161073891908590600290839083908082843760009201919091525061065a915050565b6000818152600460205260409020549091506001600160a01b031680156107a6576040805162461bcd60e51b815260206004820152601960248201527f706c656173652072656769737465722061206e6577206b657900000000000000604482015290519081900360640190fd5b6001600160a01b038516610801576040805162461bcd60e51b815260206004820152601760248201527f5f6f7261636c65206d757374206e6f7420626520307830000000000000000000604482015290519081900360640190fd5b6000828152600460205260409020805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0387161781556001018390556b033b2e3c9fd0803ce80000008611156108865760405162461bcd60e51b815260040180806020018281038252603c815260200180612086603c913960400191505060405180910390fd5b60008281526004602090815260409182902080546001600160a01b0316600160a01b6bffffffffffffffffffffffff8b1602179055815184815290810188905281517fae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe929181900390910190a1505050505050565b6101a081565b6109096105ac565b61095a576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b610963816111bd565b50565b3360009081526005602052604090205481908111156109cc576040805162461bcd60e51b815260206004820181905260248201527f63616e2774207769746864726177206d6f7265207468616e2062616c616e6365604482015290519081900360640190fd5b336000908152600560205260409020546109ec908363ffffffff61126a16565b3360009081526005602090815260408083209390935560015483517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b038881166004830152602482018890529451949091169363a9059cbb93604480840194938390030190829087803b158015610a6b57600080fd5b505af1158015610a7f573d6000803e3d6000fd5b505050506040513d6020811015610a9557600080fd5b5051610a9d57fe5b505050565b6000610aac611fc7565b825160009081906101c0908114610b0a576040805162461bcd60e51b815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b610b12611fe7565b5060e086015181870151602088019190610b2b8361065a565b9750610b3788836112c7565b600081815260036020908152604091829020825160608101845281546001600160a01b038116808352600160a01b9091046bffffffffffffffffffffffff169382019390935260019091015492810192909252909850909650610be1576040805162461bcd60e51b815260206004820152601860248201527f6e6f20636f72726573706f6e64696e6720726571756573740000000000000000604482015290519081900360640190fd5b6040805160208082018590528183018490528251808303840181526060909201835281519101209088015114610c5e576040805162461bcd60e51b815260206004820152601a60248201527f77726f6e672070726553656564206f7220626c6f636b206e756d000000000000604482015290519081900360640190fd5b804080610d4a57600254604080517fe9413d380000000000000000000000000000000000000000000000000000000081526004810185905290516001600160a01b039092169163e9413d3891602480820192602092909190829003018186803b158015610cca57600080fd5b505afa158015610cde573d6000803e3d6000fd5b505050506040513d6020811015610cf457600080fd5b5051905080610d4a576040805162461bcd60e51b815260206004820152601660248201527f706c656173652070726f766520626c6f636b6861736800000000000000000000604482015290519081900360640190fd5b6040805160208082018690528183018490528251808303840181526060909201909252805191012060e08b018190526101a08b52610d878b6112f3565b96505050505050509193509193565b600082820183811015610df0576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b90505b92915050565b604080516024810185905260448082018590528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f94985ddd00000000000000000000000000000000000000000000000000000000179052600090620324b0805a1015610ec2576040805162461bcd60e51b815260206004820152601b60248201527f6e6f7420656e6f7567682067617320666f7220636f6e73756d65720000000000604482015290519081900360640190fd5b6000846001600160a01b0316836040518082805190602001908083835b60208310610efe5780518252601f199092019160209182019101610edf565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610f60576040519150601f19603f3d011682016040523d82523d6000602084013e610f65565b606091505b50505050505050505050565b60008481526004602052604090205482908590600160a01b90046bffffffffffffffffffffffff16821015610fed576040805162461bcd60e51b815260206004820152601460248201527f42656c6f7720616772656564207061796d656e74000000000000000000000000604482015290519081900360640190fd5b60008681526006602090815260408083206001600160a01b03871684529091528120549061101d88888785611447565b9050600061102b89836112c7565b6000818152600360205260409020549091506001600160a01b03161561104d57fe5b6000818152600360205260409020805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0388161790556b033b2e3c9fd0803ce8000000871061109657fe5b600081815260036020908152604080832080546bffffffffffffffffffffffff8c16600160a01b026001600160a01b0391821617825582518085018890524381850152835180820385018152606082018086528151918701919091206001948501558f875260049095529483902090910154928d905260808401869052891660a084015260c083018a905260e083018490525190917f56bd374744a66d531874338def36c906e3a6cf31176eb1e9afd9f1de69725d5191908190036101000190a260008981526006602090815260408083206001600160a01b038a16845290915290205461118b90600163ffffffff610d9616565b6000998a52600660209081526040808c206001600160a01b039099168c52979052959098209490945550505050505050565b6001600160a01b0381166112025760405162461bcd60e51b81526004018080602001828103825260268152602001806120606026913960400191505060405180910390fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b6000828211156112c1576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60006101a082511461134c576040805162461bcd60e51b815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b611354611fe7565b61135c611fe7565b611364612005565b600061136e611fe7565b611376611fe7565b6000888060200190516101a081101561138e57600080fd5b5060e08101516101808201519198506040890197506080890196509450610100880193506101408801925090506113e187878760006020020151886001602002015189600260200201518989898961148e565b6003866040516020018083815260200182600260200280838360005b838110156114155781810151838201526020016113fd565b50505050905001925050506040516020818303038152906040528051906020012060001c975050505050505050919050565b60408051602080820196909652808201949094526001600160a01b039290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b611497896116f3565b6114e8576040805162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015290519081900360640190fd5b6114f1886116f3565b611542576040805162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015290519081900360640190fd5b61154b836116f3565b61159c576040805162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015290519081900360640190fd5b6115a5826116f3565b6115f6576040805162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015290519081900360640190fd5b611602878a888761171d565b611653576040805162461bcd60e51b815260206004820152601a60248201527f6164647228632a706b2b732a6729e289a05f755769746e657373000000000000604482015290519081900360640190fd5b61165b611fe7565b6116658a8761185d565b905061166f611fe7565b61167e898b878b868989611900565b9050600061168f838d8d8a86611a0b565b9050808a146116e5576040805162461bcd60e51b815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015290519081900360640190fd5b505050505050505050505050565b60208101516000906401000003d0199080096117168360005b6020020151611b14565b1492915050565b60006001600160a01b03821661177a576040805162461bcd60e51b815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015290519081900360640190fd5b60208401516000906001161561179157601c611794565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe1991820392506000919089098751604080516000808252602082810180855288905260ff8916838501526060830194909452608082018590529151939450909260019260a0808401939192601f1981019281900390910190855afa158015611835573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b611865611fe7565b6118c3600184846040516020018084815260200183600260200280838360005b8381101561189d578181015183820152602001611885565b505050509050018281526020019350505050604051602081830303815290604052611b38565b90505b6118cf816116f3565b610df35780516040805160208181019390935281518082039093018352810190526118f990611b38565b90506118c6565b611908611fe7565b825186516401000003d01991900306611968576040805162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015290519081900360640190fd5b611973878988611b86565b6119ae5760405162461bcd60e51b81526004018080602001828103825260218152602001806120c26021913960400191505060405180910390fd5b6119b9848685611b86565b6119f45760405162461bcd60e51b81526004018080602001828103825260228152602001806120e36022913960400191505060405180910390fd5b6119ff868484611ca6565b98975050505050505050565b6000600286868685876040516020018087815260200186600260200280838360005b83811015611a45578181015183820152602001611a2d565b5050505090500185600260200280838360005b83811015611a70578181015183820152602001611a58565b5050505090500184600260200280838360005b83811015611a9b578181015183820152602001611a83565b5050505090500183600260200280838360005b83811015611ac6578181015183820152602001611aae565b50505050905001826001600160a01b03166001600160a01b031660601b815260140196505050505050506040516020818303038152906040528051906020012060001c905095945050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b611b40611fe7565b611b4982611d6c565b8152611b5e611b5982600061170c565b611da7565b6020820181905260029006600114156106ab576020810180516401000003d019039052919050565b600082611b9257600080fd5b8351602085015160009060011615611bab57601c611bae565b601b5b9050600070014551231950b75fc4402da1732fc9bebe19838709604080516000808252602080830180855282905260ff871683850152606083018890526080830185905292519394509260019260a0808401939192601f1981019281900390910190855afa158015611c24573d6000803e3d6000fd5b5050506020604051035190506000866040516020018082600260200280838360005b83811015611c5e578181015183820152602001611c46565b505050509050019150506040516020818303038152906040528051906020012060001c9050806001600160a01b0316826001600160a01b031614955050505050509392505050565b611cae611fe7565b835160208086015185519186015160009384938493611ccf93909190611dd3565b919450925090506401000003d019858209600114611d34576040805162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a00000000000000604482015290519081900360640190fd5b60405180604001604052806401000003d01980611d4d57fe5b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d01981106106ab57604080516020808201939093528151808203840181529082019091528051910120611d74565b6000610df3827f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c611eb3565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000611e1383838585611f5a565b9098509050611e2488828e88611f7e565b9098509050611e3588828c87611f7e565b90985090506000611e488d878b85611f7e565b9098509050611e5988828686611f5a565b9098509050611e6a88828e89611f7e565b9098509050818114611e9f576401000003d019818a0998506401000003d01982890997506401000003d0198183099650611ea3565b8196505b5050505050509450945094915050565b600080611ebe612023565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152611ef0612041565b60208160c0846005600019fa925082611f50576040805162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015290519081900360640190fd5b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b604080516060810182526000808252602082018190529181019190915290565b60405180604001604052806002906020820280368337509192915050565b60405180606001604052806003906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b6040518060200160405280600190602082028036833750919291505056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373796f752063616e277420636861726765206d6f7265207468616e20616c6c20746865204c494e4b20696e2074686520776f726c642c206772656564794669727374206d756c7469706c69636174696f6e20636865636b206661696c65645365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c6564a26469706673582212209df11cae09f0d92fc956e3096bd5448cdd94115af89490206f202fcfdf17510964736f6c63430006060033"

// DeployVRFCoordinator deploys a new Ethereum contract, binding an instance of VRFCoordinator to it.
func DeployVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _blockHashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinator, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFCoordinatorBin), backend, _link, _blockHashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

// VRFCoordinator is an auto generated Go binding around an Ethereum contract.
type VRFCoordinator struct {
	VRFCoordinatorCaller     // Read-only binding to the contract
	VRFCoordinatorTransactor // Write-only binding to the contract
	VRFCoordinatorFilterer   // Log filterer for contract events
}

// VRFCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFCoordinatorSession struct {
	Contract     *VRFCoordinator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFCoordinatorCallerSession struct {
	Contract *VRFCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// VRFCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFCoordinatorTransactorSession struct {
	Contract     *VRFCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// VRFCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFCoordinatorRaw struct {
	Contract *VRFCoordinator // Generic contract binding to access the raw methods on
}

// VRFCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFCoordinatorCallerRaw struct {
	Contract *VRFCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// VRFCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFCoordinatorTransactorRaw struct {
	Contract *VRFCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFCoordinator creates a new instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinator(address common.Address, backend bind.ContractBackend) (*VRFCoordinator, error) {
	contract, err := bindVRFCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

// NewVRFCoordinatorCaller creates a new read-only instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorCaller, error) {
	contract, err := bindVRFCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCaller{contract: contract}, nil
}

// NewVRFCoordinatorTransactor creates a new write-only instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTransactor, error) {
	contract, err := bindVRFCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTransactor{contract: contract}, nil
}

// NewVRFCoordinatorFilterer creates a new log filterer instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorFilterer, error) {
	contract, err := bindVRFCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorFilterer{contract: contract}, nil
}

// bindVRFCoordinator binds a generic wrapper to an already deployed contract.
func bindVRFCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinator *VRFCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.VRFCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinator *VRFCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinator *VRFCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinator *VRFCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transact(opts, method, params...)
}

// PRESEEDOFFSET is a free data retrieval call binding the contract method 0xb415f4f5.
//
// Solidity: function PRESEED_OFFSET() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCaller) PRESEEDOFFSET(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "PRESEED_OFFSET")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRESEEDOFFSET is a free data retrieval call binding the contract method 0xb415f4f5.
//
// Solidity: function PRESEED_OFFSET() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorSession) PRESEEDOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PRESEEDOFFSET(&_VRFCoordinator.CallOpts)
}

// PRESEEDOFFSET is a free data retrieval call binding the contract method 0xb415f4f5.
//
// Solidity: function PRESEED_OFFSET() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCallerSession) PRESEEDOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PRESEEDOFFSET(&_VRFCoordinator.CallOpts)
}

// PROOFLENGTH is a free data retrieval call binding the contract method 0xe911439c.
//
// Solidity: function PROOF_LENGTH() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCaller) PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "PROOF_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PROOFLENGTH is a free data retrieval call binding the contract method 0xe911439c.
//
// Solidity: function PROOF_LENGTH() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorSession) PROOFLENGTH() (*big.Int, error) {
	return _VRFCoordinator.Contract.PROOFLENGTH(&_VRFCoordinator.CallOpts)
}

// PROOFLENGTH is a free data retrieval call binding the contract method 0xe911439c.
//
// Solidity: function PROOF_LENGTH() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCallerSession) PROOFLENGTH() (*big.Int, error) {
	return _VRFCoordinator.Contract.PROOFLENGTH(&_VRFCoordinator.CallOpts)
}

// PUBLICKEYOFFSET is a free data retrieval call binding the contract method 0x8aa7927b.
//
// Solidity: function PUBLIC_KEY_OFFSET() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCaller) PUBLICKEYOFFSET(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "PUBLIC_KEY_OFFSET")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PUBLICKEYOFFSET is a free data retrieval call binding the contract method 0x8aa7927b.
//
// Solidity: function PUBLIC_KEY_OFFSET() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorSession) PUBLICKEYOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PUBLICKEYOFFSET(&_VRFCoordinator.CallOpts)
}

// PUBLICKEYOFFSET is a free data retrieval call binding the contract method 0x8aa7927b.
//
// Solidity: function PUBLIC_KEY_OFFSET() view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCallerSession) PUBLICKEYOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PUBLICKEYOFFSET(&_VRFCoordinator.CallOpts)
}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) view returns(address callbackContract, uint96 randomnessFee, bytes32 seedAndBlockNum)
func (_VRFCoordinator *VRFCoordinatorCaller) Callbacks(opts *bind.CallOpts, arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	SeedAndBlockNum  [32]byte
}, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "callbacks", arg0)

	outstruct := new(struct {
		CallbackContract common.Address
		RandomnessFee    *big.Int
		SeedAndBlockNum  [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CallbackContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.RandomnessFee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SeedAndBlockNum = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) view returns(address callbackContract, uint96 randomnessFee, bytes32 seedAndBlockNum)
func (_VRFCoordinator *VRFCoordinatorSession) Callbacks(arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	SeedAndBlockNum  [32]byte
}, error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) view returns(address callbackContract, uint96 randomnessFee, bytes32 seedAndBlockNum)
func (_VRFCoordinator *VRFCoordinatorCallerSession) Callbacks(arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	SeedAndBlockNum  [32]byte
}, error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) pure returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCaller) HashOfKey(opts *bind.CallOpts, _publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "hashOfKey", _publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) pure returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) pure returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCallerSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_VRFCoordinator *VRFCoordinatorCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_VRFCoordinator *VRFCoordinatorSession) IsOwner() (bool, error) {
	return _VRFCoordinator.Contract.IsOwner(&_VRFCoordinator.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_VRFCoordinator *VRFCoordinatorCallerSession) IsOwner() (bool, error) {
	return _VRFCoordinator.Contract.IsOwner(&_VRFCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFCoordinator *VRFCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFCoordinator *VRFCoordinatorSession) Owner() (common.Address, error) {
	return _VRFCoordinator.Contract.Owner(&_VRFCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFCoordinator *VRFCoordinatorCallerSession) Owner() (common.Address, error) {
	return _VRFCoordinator.Contract.Owner(&_VRFCoordinator.CallOpts)
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) view returns(address vRFOracle, uint96 fee, bytes32 jobID)
func (_VRFCoordinator *VRFCoordinatorCaller) ServiceAgreements(opts *bind.CallOpts, arg0 [32]byte) (struct {
	VRFOracle common.Address
	Fee       *big.Int
	JobID     [32]byte
}, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "serviceAgreements", arg0)

	outstruct := new(struct {
		VRFOracle common.Address
		Fee       *big.Int
		JobID     [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.VRFOracle = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Fee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.JobID = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) view returns(address vRFOracle, uint96 fee, bytes32 jobID)
func (_VRFCoordinator *VRFCoordinatorSession) ServiceAgreements(arg0 [32]byte) (struct {
	VRFOracle common.Address
	Fee       *big.Int
	JobID     [32]byte
}, error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) view returns(address vRFOracle, uint96 fee, bytes32 jobID)
func (_VRFCoordinator *VRFCoordinatorCallerSession) ServiceAgreements(arg0 [32]byte) (struct {
	VRFOracle common.Address
	Fee       *big.Int
	JobID     [32]byte
}, error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCaller) WithdrawableTokens(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "withdrawableTokens", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCallerSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

// FulfillRandomnessRequest is a paid mutator transaction binding the contract method 0x5e1c1059.
//
// Solidity: function fulfillRandomnessRequest(bytes _proof) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) FulfillRandomnessRequest(opts *bind.TransactOpts, _proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "fulfillRandomnessRequest", _proof)
}

// FulfillRandomnessRequest is a paid mutator transaction binding the contract method 0x5e1c1059.
//
// Solidity: function fulfillRandomnessRequest(bytes _proof) returns()
func (_VRFCoordinator *VRFCoordinatorSession) FulfillRandomnessRequest(_proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomnessRequest(&_VRFCoordinator.TransactOpts, _proof)
}

// FulfillRandomnessRequest is a paid mutator transaction binding the contract method 0x5e1c1059.
//
// Solidity: function fulfillRandomnessRequest(bytes _proof) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) FulfillRandomnessRequest(_proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomnessRequest(&_VRFCoordinator.TransactOpts, _proof)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _fee, bytes _data) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "onTokenTransfer", _sender, _fee, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _fee, bytes _data) returns()
func (_VRFCoordinator *VRFCoordinatorSession) OnTokenTransfer(_sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, _sender, _fee, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _fee, bytes _data) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) OnTokenTransfer(_sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, _sender, _fee, _data)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0xd8340209.
//
// Solidity: function registerProvingKey(uint256 _fee, address _oracle, uint256[2] _publicProvingKey, bytes32 _jobID) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) RegisterProvingKey(opts *bind.TransactOpts, _fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "registerProvingKey", _fee, _oracle, _publicProvingKey, _jobID)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0xd8340209.
//
// Solidity: function registerProvingKey(uint256 _fee, address _oracle, uint256[2] _publicProvingKey, bytes32 _jobID) returns()
func (_VRFCoordinator *VRFCoordinatorSession) RegisterProvingKey(_fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _oracle, _publicProvingKey, _jobID)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0xd8340209.
//
// Solidity: function registerProvingKey(uint256 _fee, address _oracle, uint256[2] _publicProvingKey, bytes32 _jobID) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) RegisterProvingKey(_fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _oracle, _publicProvingKey, _jobID)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VRFCoordinator *VRFCoordinatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.TransferOwnership(&_VRFCoordinator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.TransferOwnership(&_VRFCoordinator.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "withdraw", _recipient, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_VRFCoordinator *VRFCoordinatorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Withdraw(&_VRFCoordinator.TransactOpts, _recipient, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Withdraw(&_VRFCoordinator.TransactOpts, _recipient, _amount)
}

// VRFCoordinatorNewServiceAgreementIterator is returned from FilterNewServiceAgreement and is used to iterate over the raw logs and unpacked data for NewServiceAgreement events raised by the VRFCoordinator contract.
type VRFCoordinatorNewServiceAgreementIterator struct {
	Event *VRFCoordinatorNewServiceAgreement // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorNewServiceAgreementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorNewServiceAgreement)
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
		it.Event = new(VRFCoordinatorNewServiceAgreement)
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
func (it *VRFCoordinatorNewServiceAgreementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorNewServiceAgreementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorNewServiceAgreement represents a NewServiceAgreement event raised by the VRFCoordinator contract.
type VRFCoordinatorNewServiceAgreement struct {
	KeyHash [32]byte
	Fee     *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNewServiceAgreement is a free log retrieval operation binding the contract event 0xae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe.
//
// Solidity: event NewServiceAgreement(bytes32 keyHash, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterNewServiceAgreement(opts *bind.FilterOpts) (*VRFCoordinatorNewServiceAgreementIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "NewServiceAgreement")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorNewServiceAgreementIterator{contract: _VRFCoordinator.contract, event: "NewServiceAgreement", logs: logs, sub: sub}, nil
}

// WatchNewServiceAgreement is a free log subscription operation binding the contract event 0xae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe.
//
// Solidity: event NewServiceAgreement(bytes32 keyHash, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchNewServiceAgreement(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNewServiceAgreement) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "NewServiceAgreement")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorNewServiceAgreement)
				if err := _VRFCoordinator.contract.UnpackLog(event, "NewServiceAgreement", log); err != nil {
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

// ParseNewServiceAgreement is a log parse operation binding the contract event 0xae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe.
//
// Solidity: event NewServiceAgreement(bytes32 keyHash, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseNewServiceAgreement(log types.Log) (*VRFCoordinatorNewServiceAgreement, error) {
	event := new(VRFCoordinatorNewServiceAgreement)
	if err := _VRFCoordinator.contract.UnpackLog(event, "NewServiceAgreement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the VRFCoordinator contract.
type VRFCoordinatorOwnershipTransferredIterator struct {
	Event *VRFCoordinatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorOwnershipTransferred)
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
		it.Event = new(VRFCoordinatorOwnershipTransferred)
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
func (it *VRFCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorOwnershipTransferred represents a OwnershipTransferred event raised by the VRFCoordinator contract.
type VRFCoordinatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VRFCoordinatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorOwnershipTransferredIterator{contract: _VRFCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorOwnershipTransferred)
				if err := _VRFCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorOwnershipTransferred, error) {
	event := new(VRFCoordinatorOwnershipTransferred)
	if err := _VRFCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorRandomnessRequestIterator is returned from FilterRandomnessRequest and is used to iterate over the raw logs and unpacked data for RandomnessRequest events raised by the VRFCoordinator contract.
type VRFCoordinatorRandomnessRequestIterator struct {
	Event *VRFCoordinatorRandomnessRequest // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorRandomnessRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessRequest)
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
		it.Event = new(VRFCoordinatorRandomnessRequest)
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
func (it *VRFCoordinatorRandomnessRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorRandomnessRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorRandomnessRequest represents a RandomnessRequest event raised by the VRFCoordinator contract.
type VRFCoordinatorRandomnessRequest struct {
	KeyHash   [32]byte
	Seed      *big.Int
	JobID     [32]byte
	Sender    common.Address
	Fee       *big.Int
	RequestID [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRandomnessRequest is a free log retrieval operation binding the contract event 0x56bd374744a66d531874338def36c906e3a6cf31176eb1e9afd9f1de69725d51.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 indexed jobID, address sender, uint256 fee, bytes32 requestID)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequest(opts *bind.FilterOpts, jobID [][32]byte) (*VRFCoordinatorRandomnessRequestIterator, error) {

	var jobIDRule []interface{}
	for _, jobIDItem := range jobID {
		jobIDRule = append(jobIDRule, jobIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequest", jobIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequest", logs: logs, sub: sub}, nil
}

// WatchRandomnessRequest is a free log subscription operation binding the contract event 0x56bd374744a66d531874338def36c906e3a6cf31176eb1e9afd9f1de69725d51.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 indexed jobID, address sender, uint256 fee, bytes32 requestID)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequest(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequest, jobID [][32]byte) (event.Subscription, error) {

	var jobIDRule []interface{}
	for _, jobIDItem := range jobID {
		jobIDRule = append(jobIDRule, jobIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequest", jobIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorRandomnessRequest)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
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

// ParseRandomnessRequest is a log parse operation binding the contract event 0x56bd374744a66d531874338def36c906e3a6cf31176eb1e9afd9f1de69725d51.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 indexed jobID, address sender, uint256 fee, bytes32 requestID)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequest(log types.Log) (*VRFCoordinatorRandomnessRequest, error) {
	event := new(VRFCoordinatorRandomnessRequest)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFCoordinatorRandomnessRequestFulfilledIterator is returned from FilterRandomnessRequestFulfilled and is used to iterate over the raw logs and unpacked data for RandomnessRequestFulfilled events raised by the VRFCoordinator contract.
type VRFCoordinatorRandomnessRequestFulfilledIterator struct {
	Event *VRFCoordinatorRandomnessRequestFulfilled // Event containing the contract specifics and raw log

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
func (it *VRFCoordinatorRandomnessRequestFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessRequestFulfilled)
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
		it.Event = new(VRFCoordinatorRandomnessRequestFulfilled)
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
func (it *VRFCoordinatorRandomnessRequestFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorRandomnessRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorRandomnessRequestFulfilled represents a RandomnessRequestFulfilled event raised by the VRFCoordinator contract.
type VRFCoordinatorRandomnessRequestFulfilled struct {
	RequestId [32]byte
	Output    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRandomnessRequestFulfilled is a free log retrieval operation binding the contract event 0xa2e7a402243ebda4a69ceeb3dfb682943b7a9b3ac66d6eefa8db65894009611c.
//
// Solidity: event RandomnessRequestFulfilled(bytes32 requestId, uint256 output)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequestFulfilled(opts *bind.FilterOpts) (*VRFCoordinatorRandomnessRequestFulfilledIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestFulfilledIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequestFulfilled", logs: logs, sub: sub}, nil
}

// WatchRandomnessRequestFulfilled is a free log subscription operation binding the contract event 0xa2e7a402243ebda4a69ceeb3dfb682943b7a9b3ac66d6eefa8db65894009611c.
//
// Solidity: event RandomnessRequestFulfilled(bytes32 requestId, uint256 output)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorRandomnessRequestFulfilled)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequestFulfilled", log); err != nil {
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

// ParseRandomnessRequestFulfilled is a log parse operation binding the contract event 0xa2e7a402243ebda4a69ceeb3dfb682943b7a9b3ac66d6eefa8db65894009611c.
//
// Solidity: event RandomnessRequestFulfilled(bytes32 requestId, uint256 output)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequestFulfilled(log types.Log) (*VRFCoordinatorRandomnessRequestFulfilled, error) {
	event := new(VRFCoordinatorRandomnessRequestFulfilled)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
