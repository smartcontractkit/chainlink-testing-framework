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

// FlagsABI is the input ABI used to generate the binding from.
const FlagsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"racAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"FlagLowered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"FlagRaised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RaisingAccessControllerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"getFlag\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"lowerFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"raiseFlag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"raiseFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"raisingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"racAddress\",\"type\":\"address\"}],\"name\":\"setRaisingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// FlagsBin is the compiled bytecode used for deploying new contracts.
var FlagsBin = "0x608060405234801561001057600080fd5b50604051611d87380380611d878339818101604052602081101561003357600080fd5b8101908080519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060018060146101000a81548160ff0219169083151502179055506100ad816100b360201b60201c565b5061026f565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610175576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461026b5781600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167fbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac60405160405180910390a35b5050565b611b098061027e6000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80637d723cac11610097578063a118f24911610066578063a118f2491461054e578063d74af26314610592578063dc7f0124146105d6578063f2fde38b146105f857610100565b80637d723cac146103e85780638038e4a1146104b65780638823da6c146104c05780638da5cb5b1461050457610100565b8063517e89fe116100d3578063517e89fe1461022e5780636b14daf814610272578063760bc82d1461036557806379ba5097146103de57610100565b80630a75698314610105578063282865961461010f5780632e1d859c14610188578063357e47fe146101d2575b600080fd5b61010d61063c565b005b6101866004803603602081101561012557600080fd5b810190808035906020019064010000000081111561014257600080fd5b82018360208201111561015457600080fd5b8035906020019184602083028401116401000000008311171561017657600080fd5b909192939192939050505061075d565b005b61019061095e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610214600480360360208110156101e857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610984565b604051808215151515815260200191505060405180910390f35b6102706004803603602081101561024457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a9b565b005b61034b6004803603604081101561028857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156102c557600080fd5b8201836020820111156102d757600080fd5b803590602001918460018302840111640100000000831117156102f957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610c57565b604051808215151515815260200191505060405180910390f35b6103dc6004803603602081101561037b57600080fd5b810190808035906020019064010000000081111561039857600080fd5b8201836020820111156103aa57600080fd5b803590602001918460208302840111640100000000831117156103cc57600080fd5b9091929391929390505050610ca1565b005b6103e6610d6f565b005b61045f600480360360208110156103fe57600080fd5b810190808035906020019064010000000081111561041b57600080fd5b82018360208201111561042d57600080fd5b8035906020019184602083028401116401000000008311171561044f57600080fd5b9091929391929390505050610f37565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156104a2578082015181840152602081019050610487565b505050509050019250505060405180910390f35b6104be611104565b005b610502600480360360208110156104d657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611223565b005b61050c6113f6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6105906004803603602081101561056457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061141b565b005b6105d4600480360360208110156105a857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506115ed565b005b6105de611673565b604051808215151515815260200191505060405180910390f35b61063a6004803603602081101561060e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611686565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600160149054906101000a900460ff161561075b576000600160146101000a81548160ff0219169083151502179055507f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963860405160405180910390a15b565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461081f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b60008090505b8282905081101561095957600083838381811061083e57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff169050600460008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161561094b576000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508073ffffffffffffffffffffffffffffffffffffffff167fd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc860405160405180910390a25b508080600101915050610825565b505050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006109d5336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610c57565b610a47576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610b5d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610c535781600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167fbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac60405160405180910390a35b5050565b6000610c638383611807565b80610c9957503273ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16145b905092915050565b610ca9611876565b610d1b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f4e6f7420616c6c6f77656420746f20726169736520666c61677300000000000081525060200191505060405180910390fd5b60008090505b82829050811015610d6a57610d5d838383818110610d3b57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff166119e3565b8080600101915050610d21565b505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610e32576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4d7573742062652070726f706f736564206f776e65720000000000000000000081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a350565b6060610f88336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610c57565b610ffa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f4e6f20616363657373000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b60608383905067ffffffffffffffff8111801561101657600080fd5b506040519080825280602002602001820160405280156110455781602001602082028036833780820191505090505b50905060008090505b848490508110156110f9576004600086868481811061106957fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff168282815181106110da57fe5b602002602001019015159081151581525050808060010191505061104e565b508091505092915050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146111c6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600160149054906101000a900460ff166112215760018060146101000a81548160ff0219169083151502179055507faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348060405160405180910390a15b565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146112e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600260008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156113f3576000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d181604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15b50565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146114dd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b600260008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166115ea576001600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db481604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15b50565b6115f5611876565b611667576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f4e6f7420616c6c6f77656420746f20726169736520666c61677300000000000081525060200191505060405180910390fd5b611670816119e3565b50565b600160149054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611748576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000081525060200191505060405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b6000600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff168061186e5750600160149054906101000a900460ff16155b905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806119de5750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636b14daf8336000366040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001806020018281038252848482818152602001925080828437600081840152601f19601f82011690508083019250505094505050505060206040518083038186803b1580156119a257600080fd5b505afa1580156119b6573d6000803e3d6000fd5b505050506040513d60208110156119cc57600080fd5b81019080805190602001909291905050505b905090565b600460008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16611ad0576001600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508073ffffffffffffffffffffffffffffffffffffffff167f881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c00660405160405180910390a25b5056fea2646970667358221220c783a91cbc8e98d80cb146668525bedbb3c8ad06f2293efc24701ca3fd545b4764736f6c63430006060033"

// DeployFlags deploys a new Ethereum contract, binding an instance of Flags to it.
func DeployFlags(auth *bind.TransactOpts, backend bind.ContractBackend, racAddress common.Address) (common.Address, *types.Transaction, *Flags, error) {
	parsed, err := abi.JSON(strings.NewReader(FlagsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FlagsBin), backend, racAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Flags{FlagsCaller: FlagsCaller{contract: contract}, FlagsTransactor: FlagsTransactor{contract: contract}, FlagsFilterer: FlagsFilterer{contract: contract}}, nil
}

// Flags is an auto generated Go binding around an Ethereum contract.
type Flags struct {
	FlagsCaller     // Read-only binding to the contract
	FlagsTransactor // Write-only binding to the contract
	FlagsFilterer   // Log filterer for contract events
}

// FlagsCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlagsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlagsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlagsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlagsSession struct {
	Contract     *Flags            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlagsCallerSession struct {
	Contract *FlagsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FlagsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlagsTransactorSession struct {
	Contract     *FlagsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlagsRaw struct {
	Contract *Flags // Generic contract binding to access the raw methods on
}

// FlagsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlagsCallerRaw struct {
	Contract *FlagsCaller // Generic read-only contract binding to access the raw methods on
}

// FlagsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlagsTransactorRaw struct {
	Contract *FlagsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlags creates a new instance of Flags, bound to a specific deployed contract.
func NewFlags(address common.Address, backend bind.ContractBackend) (*Flags, error) {
	contract, err := bindFlags(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Flags{FlagsCaller: FlagsCaller{contract: contract}, FlagsTransactor: FlagsTransactor{contract: contract}, FlagsFilterer: FlagsFilterer{contract: contract}}, nil
}

// NewFlagsCaller creates a new read-only instance of Flags, bound to a specific deployed contract.
func NewFlagsCaller(address common.Address, caller bind.ContractCaller) (*FlagsCaller, error) {
	contract, err := bindFlags(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsCaller{contract: contract}, nil
}

// NewFlagsTransactor creates a new write-only instance of Flags, bound to a specific deployed contract.
func NewFlagsTransactor(address common.Address, transactor bind.ContractTransactor) (*FlagsTransactor, error) {
	contract, err := bindFlags(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsTransactor{contract: contract}, nil
}

// NewFlagsFilterer creates a new log filterer instance of Flags, bound to a specific deployed contract.
func NewFlagsFilterer(address common.Address, filterer bind.ContractFilterer) (*FlagsFilterer, error) {
	contract, err := bindFlags(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlagsFilterer{contract: contract}, nil
}

// bindFlags binds a generic wrapper to an already deployed contract.
func bindFlags(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlagsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flags *FlagsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flags.Contract.FlagsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flags *FlagsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.Contract.FlagsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flags *FlagsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flags.Contract.FlagsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flags *FlagsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flags.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flags *FlagsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flags *FlagsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flags.Contract.contract.Transact(opts, method, params...)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsSession) CheckEnabled() (bool, error) {
	return _Flags.Contract.CheckEnabled(&_Flags.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsCallerSession) CheckEnabled() (bool, error) {
	return _Flags.Contract.CheckEnabled(&_Flags.CallOpts)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsCaller) GetFlag(opts *bind.CallOpts, subject common.Address) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "getFlag", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsSession) GetFlag(subject common.Address) (bool, error) {
	return _Flags.Contract.GetFlag(&_Flags.CallOpts, subject)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsCallerSession) GetFlag(subject common.Address) (bool, error) {
	return _Flags.Contract.GetFlag(&_Flags.CallOpts, subject)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsCaller) GetFlags(opts *bind.CallOpts, subjects []common.Address) ([]bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "getFlags", subjects)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsSession) GetFlags(subjects []common.Address) ([]bool, error) {
	return _Flags.Contract.GetFlags(&_Flags.CallOpts, subjects)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsCallerSession) GetFlags(subjects []common.Address) ([]bool, error) {
	return _Flags.Contract.GetFlags(&_Flags.CallOpts, subjects)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _Flags.Contract.HasAccess(&_Flags.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _Flags.Contract.HasAccess(&_Flags.CallOpts, _user, _calldata)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsSession) Owner() (common.Address, error) {
	return _Flags.Contract.Owner(&_Flags.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsCallerSession) Owner() (common.Address, error) {
	return _Flags.Contract.Owner(&_Flags.CallOpts)
}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsCaller) RaisingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "raisingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsSession) RaisingAccessController() (common.Address, error) {
	return _Flags.Contract.RaisingAccessController(&_Flags.CallOpts)
}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsCallerSession) RaisingAccessController() (common.Address, error) {
	return _Flags.Contract.RaisingAccessController(&_Flags.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsSession) AcceptOwnership() (*types.Transaction, error) {
	return _Flags.Contract.AcceptOwnership(&_Flags.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Flags.Contract.AcceptOwnership(&_Flags.TransactOpts)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.AddAccess(&_Flags.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.AddAccess(&_Flags.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsSession) DisableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.DisableAccessCheck(&_Flags.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.DisableAccessCheck(&_Flags.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsSession) EnableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.EnableAccessCheck(&_Flags.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.EnableAccessCheck(&_Flags.TransactOpts)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsTransactor) LowerFlags(opts *bind.TransactOpts, subjects []common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "lowerFlags", subjects)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsSession) LowerFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.LowerFlags(&_Flags.TransactOpts, subjects)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsTransactorSession) LowerFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.LowerFlags(&_Flags.TransactOpts, subjects)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsTransactor) RaiseFlag(opts *bind.TransactOpts, subject common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "raiseFlag", subject)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsSession) RaiseFlag(subject common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlag(&_Flags.TransactOpts, subject)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsTransactorSession) RaiseFlag(subject common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlag(&_Flags.TransactOpts, subject)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsTransactor) RaiseFlags(opts *bind.TransactOpts, subjects []common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "raiseFlags", subjects)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsSession) RaiseFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlags(&_Flags.TransactOpts, subjects)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsTransactorSession) RaiseFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlags(&_Flags.TransactOpts, subjects)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RemoveAccess(&_Flags.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RemoveAccess(&_Flags.TransactOpts, _user)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsTransactor) SetRaisingAccessController(opts *bind.TransactOpts, racAddress common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "setRaisingAccessController", racAddress)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsSession) SetRaisingAccessController(racAddress common.Address) (*types.Transaction, error) {
	return _Flags.Contract.SetRaisingAccessController(&_Flags.TransactOpts, racAddress)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsTransactorSession) SetRaisingAccessController(racAddress common.Address) (*types.Transaction, error) {
	return _Flags.Contract.SetRaisingAccessController(&_Flags.TransactOpts, racAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Flags *FlagsTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Flags *FlagsSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Flags.Contract.TransferOwnership(&_Flags.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Flags *FlagsTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Flags.Contract.TransferOwnership(&_Flags.TransactOpts, _to)
}

// FlagsAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the Flags contract.
type FlagsAddedAccessIterator struct {
	Event *FlagsAddedAccess // Event containing the contract specifics and raw log

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
func (it *FlagsAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsAddedAccess)
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
		it.Event = new(FlagsAddedAccess)
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
func (it *FlagsAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsAddedAccess represents a AddedAccess event raised by the Flags contract.
type FlagsAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*FlagsAddedAccessIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &FlagsAddedAccessIterator{contract: _Flags.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *FlagsAddedAccess) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsAddedAccess)
				if err := _Flags.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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
func (_Flags *FlagsFilterer) ParseAddedAccess(log types.Log) (*FlagsAddedAccess, error) {
	event := new(FlagsAddedAccess)
	if err := _Flags.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the Flags contract.
type FlagsCheckAccessDisabledIterator struct {
	Event *FlagsCheckAccessDisabled // Event containing the contract specifics and raw log

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
func (it *FlagsCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsCheckAccessDisabled)
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
		it.Event = new(FlagsCheckAccessDisabled)
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
func (it *FlagsCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsCheckAccessDisabled represents a CheckAccessDisabled event raised by the Flags contract.
type FlagsCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*FlagsCheckAccessDisabledIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &FlagsCheckAccessDisabledIterator{contract: _Flags.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *FlagsCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsCheckAccessDisabled)
				if err := _Flags.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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
func (_Flags *FlagsFilterer) ParseCheckAccessDisabled(log types.Log) (*FlagsCheckAccessDisabled, error) {
	event := new(FlagsCheckAccessDisabled)
	if err := _Flags.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the Flags contract.
type FlagsCheckAccessEnabledIterator struct {
	Event *FlagsCheckAccessEnabled // Event containing the contract specifics and raw log

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
func (it *FlagsCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsCheckAccessEnabled)
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
		it.Event = new(FlagsCheckAccessEnabled)
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
func (it *FlagsCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsCheckAccessEnabled represents a CheckAccessEnabled event raised by the Flags contract.
type FlagsCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*FlagsCheckAccessEnabledIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &FlagsCheckAccessEnabledIterator{contract: _Flags.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *FlagsCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsCheckAccessEnabled)
				if err := _Flags.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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
func (_Flags *FlagsFilterer) ParseCheckAccessEnabled(log types.Log) (*FlagsCheckAccessEnabled, error) {
	event := new(FlagsCheckAccessEnabled)
	if err := _Flags.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsFlagLoweredIterator is returned from FilterFlagLowered and is used to iterate over the raw logs and unpacked data for FlagLowered events raised by the Flags contract.
type FlagsFlagLoweredIterator struct {
	Event *FlagsFlagLowered // Event containing the contract specifics and raw log

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
func (it *FlagsFlagLoweredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsFlagLowered)
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
		it.Event = new(FlagsFlagLowered)
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
func (it *FlagsFlagLoweredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsFlagLoweredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsFlagLowered represents a FlagLowered event raised by the Flags contract.
type FlagsFlagLowered struct {
	Subject common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFlagLowered is a free log retrieval operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) FilterFlagLowered(opts *bind.FilterOpts, subject []common.Address) (*FlagsFlagLoweredIterator, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "FlagLowered", subjectRule)
	if err != nil {
		return nil, err
	}
	return &FlagsFlagLoweredIterator{contract: _Flags.contract, event: "FlagLowered", logs: logs, sub: sub}, nil
}

// WatchFlagLowered is a free log subscription operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) WatchFlagLowered(opts *bind.WatchOpts, sink chan<- *FlagsFlagLowered, subject []common.Address) (event.Subscription, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "FlagLowered", subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsFlagLowered)
				if err := _Flags.contract.UnpackLog(event, "FlagLowered", log); err != nil {
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

// ParseFlagLowered is a log parse operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) ParseFlagLowered(log types.Log) (*FlagsFlagLowered, error) {
	event := new(FlagsFlagLowered)
	if err := _Flags.contract.UnpackLog(event, "FlagLowered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsFlagRaisedIterator is returned from FilterFlagRaised and is used to iterate over the raw logs and unpacked data for FlagRaised events raised by the Flags contract.
type FlagsFlagRaisedIterator struct {
	Event *FlagsFlagRaised // Event containing the contract specifics and raw log

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
func (it *FlagsFlagRaisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsFlagRaised)
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
		it.Event = new(FlagsFlagRaised)
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
func (it *FlagsFlagRaisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsFlagRaisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsFlagRaised represents a FlagRaised event raised by the Flags contract.
type FlagsFlagRaised struct {
	Subject common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFlagRaised is a free log retrieval operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) FilterFlagRaised(opts *bind.FilterOpts, subject []common.Address) (*FlagsFlagRaisedIterator, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "FlagRaised", subjectRule)
	if err != nil {
		return nil, err
	}
	return &FlagsFlagRaisedIterator{contract: _Flags.contract, event: "FlagRaised", logs: logs, sub: sub}, nil
}

// WatchFlagRaised is a free log subscription operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) WatchFlagRaised(opts *bind.WatchOpts, sink chan<- *FlagsFlagRaised, subject []common.Address) (event.Subscription, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "FlagRaised", subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsFlagRaised)
				if err := _Flags.contract.UnpackLog(event, "FlagRaised", log); err != nil {
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

// ParseFlagRaised is a log parse operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) ParseFlagRaised(log types.Log) (*FlagsFlagRaised, error) {
	event := new(FlagsFlagRaised)
	if err := _Flags.contract.UnpackLog(event, "FlagRaised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Flags contract.
type FlagsOwnershipTransferRequestedIterator struct {
	Event *FlagsOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *FlagsOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsOwnershipTransferRequested)
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
		it.Event = new(FlagsOwnershipTransferRequested)
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
func (it *FlagsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Flags contract.
type FlagsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FlagsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FlagsOwnershipTransferRequestedIterator{contract: _Flags.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FlagsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsOwnershipTransferRequested)
				if err := _Flags.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_Flags *FlagsFilterer) ParseOwnershipTransferRequested(log types.Log) (*FlagsOwnershipTransferRequested, error) {
	event := new(FlagsOwnershipTransferRequested)
	if err := _Flags.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Flags contract.
type FlagsOwnershipTransferredIterator struct {
	Event *FlagsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FlagsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsOwnershipTransferred)
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
		it.Event = new(FlagsOwnershipTransferred)
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
func (it *FlagsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsOwnershipTransferred represents a OwnershipTransferred event raised by the Flags contract.
type FlagsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FlagsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FlagsOwnershipTransferredIterator{contract: _Flags.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FlagsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsOwnershipTransferred)
				if err := _Flags.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Flags *FlagsFilterer) ParseOwnershipTransferred(log types.Log) (*FlagsOwnershipTransferred, error) {
	event := new(FlagsOwnershipTransferred)
	if err := _Flags.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsRaisingAccessControllerUpdatedIterator is returned from FilterRaisingAccessControllerUpdated and is used to iterate over the raw logs and unpacked data for RaisingAccessControllerUpdated events raised by the Flags contract.
type FlagsRaisingAccessControllerUpdatedIterator struct {
	Event *FlagsRaisingAccessControllerUpdated // Event containing the contract specifics and raw log

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
func (it *FlagsRaisingAccessControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsRaisingAccessControllerUpdated)
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
		it.Event = new(FlagsRaisingAccessControllerUpdated)
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
func (it *FlagsRaisingAccessControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsRaisingAccessControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsRaisingAccessControllerUpdated represents a RaisingAccessControllerUpdated event raised by the Flags contract.
type FlagsRaisingAccessControllerUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRaisingAccessControllerUpdated is a free log retrieval operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) FilterRaisingAccessControllerUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*FlagsRaisingAccessControllerUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "RaisingAccessControllerUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &FlagsRaisingAccessControllerUpdatedIterator{contract: _Flags.contract, event: "RaisingAccessControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchRaisingAccessControllerUpdated is a free log subscription operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) WatchRaisingAccessControllerUpdated(opts *bind.WatchOpts, sink chan<- *FlagsRaisingAccessControllerUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "RaisingAccessControllerUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsRaisingAccessControllerUpdated)
				if err := _Flags.contract.UnpackLog(event, "RaisingAccessControllerUpdated", log); err != nil {
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

// ParseRaisingAccessControllerUpdated is a log parse operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) ParseRaisingAccessControllerUpdated(log types.Log) (*FlagsRaisingAccessControllerUpdated, error) {
	event := new(FlagsRaisingAccessControllerUpdated)
	if err := _Flags.contract.UnpackLog(event, "RaisingAccessControllerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlagsRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the Flags contract.
type FlagsRemovedAccessIterator struct {
	Event *FlagsRemovedAccess // Event containing the contract specifics and raw log

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
func (it *FlagsRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsRemovedAccess)
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
		it.Event = new(FlagsRemovedAccess)
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
func (it *FlagsRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsRemovedAccess represents a RemovedAccess event raised by the Flags contract.
type FlagsRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*FlagsRemovedAccessIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &FlagsRemovedAccessIterator{contract: _Flags.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *FlagsRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsRemovedAccess)
				if err := _Flags.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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
func (_Flags *FlagsFilterer) ParseRemovedAccess(log types.Log) (*FlagsRemovedAccess, error) {
	event := new(FlagsRemovedAccess)
	if err := _Flags.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
