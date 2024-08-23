// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package network_debug_contract

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
	_ = abi.ConvertType
)

// NetworkDebugContractAccount is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugContractAccount struct {
	Name       string
	Balance    uint64
	DailyLimit *big.Int
}

// NetworkDebugContractData is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugContractData struct {
	Name   string
	Values []*big.Int
}

// NetworkDebugContractNestedData is an auto generated low-level Go binding around an user-defined struct.
type NetworkDebugContractNestedData struct {
	Data         NetworkDebugContractData
	DynamicBytes []byte
}

// NetworkDebugContractMetaData contains all meta data concerning the NetworkDebugContract contract.
var NetworkDebugContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"required\",\"type\":\"uint256\"}],\"name\":\"CustomErr\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CustomErrNoValues\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"CustomErrWithMessage\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"CallDataLength\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"a\",\"type\":\"int256\"}],\"name\":\"CallbackEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"CurrentStatus\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EtherReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"IsValidEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NoIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"NoIndexEventString\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"balance\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"dailyLimit\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structNetworkDebugContract.Account\",\"name\":\"a\",\"type\":\"tuple\"}],\"name\":\"NoIndexStructEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"OneIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"dataId\",\"type\":\"string\"}],\"name\":\"ThreeIndexAndOneNonIndexedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"ThreeIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"}],\"name\":\"TwoIndexEvent\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"idx\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"addCounter\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsAssert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsCustomError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsCustomErrorNoValues\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsRequire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"callRevertFunctionInSubContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"callRevertFunctionInTheContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"callbackMethod\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"counterMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStatus\",\"outputs\":[{\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"emitAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"input\",\"type\":\"bytes32\"}],\"name\":\"emitBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"output\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitFourParamMixedEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"inputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inputVal2\",\"type\":\"string\"}],\"name\":\"emitInputs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"inputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inputVal2\",\"type\":\"string\"}],\"name\":\"emitInputsOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"first\",\"type\":\"int256\"},{\"internalType\":\"int128\",\"name\":\"second\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"third\",\"type\":\"uint256\"}],\"name\":\"emitInts\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"int128\",\"name\":\"outputVal1\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"outputVal2\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"inputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inputVal2\",\"type\":\"string\"}],\"name\":\"emitNamedInputsOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"outputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"outputVal2\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNamedOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"outputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"outputVal2\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNoIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNoIndexEventString\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNoIndexStructEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitOneIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitThreeIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitTwoIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"data\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"idx\",\"type\":\"int256\"}],\"name\":\"getCounter\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"data\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"data\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pay\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performStaticCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"input\",\"type\":\"address[]\"}],\"name\":\"processAddressArray\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"processDynamicData\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data[3]\",\"name\":\"data\",\"type\":\"tuple[3]\"}],\"name\":\"processFixedDataArray\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data[2]\",\"name\":\"\",\"type\":\"tuple[2]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"dynamicBytes\",\"type\":\"bytes\"}],\"internalType\":\"structNetworkDebugContract.NestedData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"processNestedData\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"dynamicBytes\",\"type\":\"bytes\"}],\"internalType\":\"structNetworkDebugContract.NestedData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"processNestedData\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"dynamicBytes\",\"type\":\"bytes\"}],\"internalType\":\"structNetworkDebugContract.NestedData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"input\",\"type\":\"uint256[]\"}],\"name\":\"processUintArray\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"idx\",\"type\":\"int256\"}],\"name\":\"resetCounter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"set\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"setMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"setStatus\",\"outputs\":[{\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storedData\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"storedDataMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subContract\",\"outputs\":[{\"internalType\":\"contractNetworkDebugSubContract\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"trace\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceDifferent\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceSubWithCallback\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceWithValidate\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceYetDifferent\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162004015380380620040158339818101604052810190620000379190620000f2565b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101006004819055505062000124565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000ba826200008d565b9050919050565b620000cc81620000ad565b8114620000d857600080fd5b50565b600081519050620000ec81620000c1565b92915050565b6000602082840312156200010b576200010a62000088565b5b60006200011b84828501620000db565b91505092915050565b613ee180620001346000396000f3fe6080604052600436106102e85760003560e01c8063788c477211610190578063b1ae9d85116100dc578063e5c19b2d11610095578063ef8a92351161006f578063ef8a923514610bc8578063f3396bd914610bf3578063f499af2a14610c1c578063fbcb8d0714610c5957610328565b8063e5c19b2d14610b11578063e8116e2814610b4e578063ec5c3ede14610b8b57610328565b8063b1ae9d8514610a0d578063b600141f14610a3d578063c0d06d8914610a54578063c2124b2214610a7f578063d7a8020514610a96578063e1111f7914610ad457610328565b80639349d00b116101495780639e099652116101235780639e09965214610963578063a4c0ed36146109a2578063aa3fdcf4146109cb578063ad3de14c146109e257610328565b80639349d00b146108f857806395a81a4c1461090f57806399adad2e1461092657610328565b8063788c4772146107fb5780637f12881c146108125780637fdc8fe11461084f57806381b375a01461088c5780638db611be146108b55780638f856296146108e157610328565b806333311ef31161024f57806358379d71116102085780636284117d116101e25780636284117d1461075057806362c270e11461078d5780636d4ce63c146107a45780637014c81d146107cf57610328565b806358379d71146106bf5780635921483f146106fc5780635e9c80d61461073957610328565b806333311ef3146105625780633837a75e1461059f5780633bc5de30146105dc5780633e41f1351461060757806345f0c9e61461064457806348ad9fe81461068257610328565b806323515760116102a1578063235157601461043e578063256560d51461047b5780632a1afcd9146104925780632e49d78b146104bd57806330985bcc146104fa5780633170428e1461053757610328565b806304d8215b1461036357806306595f75146103a057806311b3c478146103b757806312d91233146103e05780631b9265b81461041d5780631e31d0a81461042757610328565b36610328577f59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832333460405161031e9291906120c3565b60405180910390a1005b7f1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b33346040516103599291906120ff565b60405180910390a1005b34801561036f57600080fd5b5061038a60048036038101906103859190612172565b610c96565b60405161039791906121cd565b60405180910390f35b3480156103ac57600080fd5b506103b5610cdc565b005b3480156103c357600080fd5b506103de60048036038101906103d99190612214565b610d1f565b005b3480156103ec57600080fd5b50610407600480360381019061040291906123ad565b610db2565b60405161041491906124b4565b60405180910390f35b610425610e71565b005b34801561043357600080fd5b5061043c610e73565b005b34801561044a57600080fd5b5061046560048036038101906104609190612172565b610eba565b60405161047291906124e5565b60405180910390f35b34801561048757600080fd5b50610490610eef565b005b34801561049e57600080fd5b506104a7610f00565b6040516104b491906124e5565b60405180910390f35b3480156104c957600080fd5b506104e460048036038101906104df9190612525565b610f06565b6040516104f191906125c9565b60405180910390f35b34801561050657600080fd5b50610521600480360381019061051c9190612172565b610f88565b60405161052e91906124e5565b60405180910390f35b34801561054357600080fd5b5061054c61106a565b60405161055991906125e4565b60405180910390f35b34801561056e57600080fd5b5061058960048036038101906105849190612635565b6111a9565b6040516105969190612671565b60405180910390f35b3480156105ab57600080fd5b506105c660048036038101906105c19190612172565b6111b3565b6040516105d391906124e5565b60405180910390f35b3480156105e857600080fd5b506105f16112be565b6040516105fe91906125e4565b60405180910390f35b34801561061357600080fd5b5061062e60048036038101906106299190612172565b6112c8565b60405161063b91906124e5565b60405180910390f35b34801561065057600080fd5b5061066b60048036038101906106669190612741565b6113c3565b60405161067992919061280b565b60405180910390f35b34801561068e57600080fd5b506106a960048036038101906106a49190612867565b6113d4565b6040516106b691906124e5565b60405180910390f35b3480156106cb57600080fd5b506106e660048036038101906106e19190612172565b6113ec565b6040516106f391906124e5565b60405180910390f35b34801561070857600080fd5b50610723600480360381019061071e9190612894565b6114e7565b60405161073091906124e5565b60405180910390f35b34801561074557600080fd5b5061074e611504565b005b34801561075c57600080fd5b5061077760048036038101906107729190612894565b611545565b60405161078491906124e5565b60405180910390f35b34801561079957600080fd5b506107a261155d565b005b3480156107b057600080fd5b506107b96115f1565b6040516107c691906124e5565b60405180910390f35b3480156107db57600080fd5b506107e46115fa565b6040516107f292919061280b565b60405180910390f35b34801561080757600080fd5b5061081061163f565b005b34801561081e57600080fd5b50610839600480360381019061083491906128e5565b611676565b6040516108469190612ac4565b60405180910390f35b34801561085b57600080fd5b5061087660048036038101906108719190612b05565b61168f565b6040516108839190612d12565b60405180910390f35b34801561089857600080fd5b506108b360048036038101906108ae9190612741565b611698565b005b3480156108c157600080fd5b506108ca61169c565b6040516108d892919061280b565b60405180910390f35b3480156108ed57600080fd5b506108f66116e1565b005b34801561090457600080fd5b5061090d611711565b005b34801561091b57600080fd5b5061092461171b565b005b34801561093257600080fd5b5061094d60048036038101906109489190612d56565b611754565b60405161095a9190612e55565b60405180910390f35b34801561096f57600080fd5b5061098a60048036038101906109859190612eb0565b6117ff565b60405161099993929190612f12565b60405180910390f35b3480156109ae57600080fd5b506109c960048036038101906109c49190612fa4565b611816565b005b3480156109d757600080fd5b506109e0611b02565b005b3480156109ee57600080fd5b506109f7611b4b565b604051610a0491906124e5565b60405180910390f35b610a276004803603810190610a229190612172565b611b92565b604051610a3491906124e5565b60405180910390f35b348015610a4957600080fd5b50610a52611cdc565b005b348015610a6057600080fd5b50610a69611d0e565b604051610a769190613077565b60405180910390f35b348015610a8b57600080fd5b50610a94611d34565b005b348015610aa257600080fd5b50610abd6004803603810190610ab89190612741565b611d86565b604051610acb92919061280b565b60405180910390f35b348015610ae057600080fd5b50610afb6004803603810190610af69190613155565b611d97565b604051610b08919061325c565b60405180910390f35b348015610b1d57600080fd5b50610b386004803603810190610b339190612894565b611da1565b604051610b4591906124e5565b60405180910390f35b348015610b5a57600080fd5b50610b756004803603810190610b709190612894565b611db2565b604051610b8291906124e5565b60405180910390f35b348015610b9757600080fd5b50610bb26004803603810190610bad9190612867565b611e00565b604051610bbf919061327e565b60405180910390f35b348015610bd457600080fd5b50610bdd611e0a565b604051610bea91906125c9565b60405180910390f35b348015610bff57600080fd5b50610c1a6004803603810190610c159190612894565b611e1d565b005b348015610c2857600080fd5b50610c436004803603810190610c3e9190612b05565b611e39565b604051610c509190612ac4565b60405180910390f35b348015610c6557600080fd5b50610c806004803603810190610c7b9190612894565b611f6e565b604051610c8d91906124e5565b60405180910390f35b60007fdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b828413604051610cc991906121cd565b60405180910390a1818313905092915050565b6000610d1d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d14906132e5565b60405180910390fd5b565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166311abb00283836040518363ffffffff1660e01b8152600401610d7c929190613305565b600060405180830381600087803b158015610d9657600080fd5b505af1158015610daa573d6000803e3d6000fd5b505050505050565b60606000825167ffffffffffffffff811115610dd157610dd061226a565b5b604051908082528060200260200182016040528015610dff5781602001602082028036833780820191505090505b50905060005b8351811015610e67576001848281518110610e2357610e2261332e565b5b6020026020010151610e35919061338c565b828281518110610e4857610e4761332e565b5b6020026020010181815250508080610e5f906133c0565b915050610e05565b5080915050919050565b565b3373ffffffffffffffffffffffffffffffffffffffff1660017f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a3565b600081600260008581526020019081526020016000206000828254610edf9190613408565b9250508190555081905092915050565b6000610efe57610efd61344c565b5b565b60005481565b600081600560006101000a81548160ff02191690836003811115610f2d57610f2c612552565b5b0217905550600560009054906101000a900460ff166003811115610f5457610f53612552565b5b7fbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a60405160405180910390a2819050919050565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663047c4425836040518263ffffffff1660e01b8152600401610fe591906124e5565b6020604051808303816000875af1158015611004573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110289190613490565b50827feace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b3560405160405180910390a281836110629190613408565b905092915050565b6000803090506000808273ffffffffffffffffffffffffffffffffffffffff16633bc5de3060e01b604051602401604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161110291906134f9565b600060405180830381855afa9150503d806000811461113d576040519150601f19603f3d011682016040523d82523d6000602084013e611142565b606091505b509150915081611187576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161117e9061355c565b60405180910390fd5b60008180602001905181019061119d9190613591565b90508094505050505090565b6000819050919050565b60006002826111c29190613408565b9150600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fa8fca7a84846040518363ffffffff1660e01b81526004016112219291906135be565b6020604051808303816000875af1158015611240573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112649190613490565b503373ffffffffffffffffffffffffffffffffffffffff1660017f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a381836112b69190613408565b905092915050565b6000600454905090565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633e41f13584846040518363ffffffff1660e01b81526004016113279291906135be565b6020604051808303816000875af1158015611346573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061136a9190613490565b503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a381836113bb9190613408565b905092915050565b600060608383915091509250929050565b60016020528060005260406000206000915090505481565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633e41f13584846040518363ffffffff1660e01b815260040161144b9291906135be565b6020604051808303816000875af115801561146a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061148e9190613490565b503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a381836114df9190613408565b905092915050565b600060026000838152602001908152602001600020549050919050565b600c60156040517f4a2eaf7e00000000000000000000000000000000000000000000000000000000815260040161153c92919061365d565b60405180910390fd5b60026020528060005260406000206000915090505481565b7febe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc04418860405180606001604052806040518060400160405280600481526020017f4a6f686e000000000000000000000000000000000000000000000000000000008152508152602001600567ffffffffffffffff168152602001600a8152506040516115e791906136f9565b60405180910390a1565b60008054905090565b60006060617a696040518060400160405280600a81526020017f6f757470757456616c3100000000000000000000000000000000000000000000815250915091509091565b7f25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a60405161166c90613767565b60405180910390a1565b61167e611fa5565b8161168890613942565b9050919050565b36819050919050565b5050565b60006060617a696040518060400160405280600a81526020017f6f757470757456616c3100000000000000000000000000000000000000000000815250915091509091565b60537feace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b3560405160405180910390a2565b611719611504565b565b7f33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c3360405161174a919061327e565b60405180910390a1565b61175c611fc5565b611764611fc5565b826000600381106117785761177761332e565b5b6020028101906117889190613964565b6117919061398c565b816000600281106117a5576117a461332e565b5b6020020181905250826001600381106117c1576117c061332e565b5b6020028101906117d19190613964565b6117da9061398c565b816001600281106117ee576117ed61332e565b5b602002018190525080915050919050565b600080600085858592509250925093509350939050565b7f962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed8282905060405161184891906125e4565b60405180910390a16000828290500361189c57606360656040517f4a2eaf7e000000000000000000000000000000000000000000000000000000008152600401611893929190613a15565b60405180910390fd5b6000803073ffffffffffffffffffffffffffffffffffffffff1684846040516118c6929190613a63565b600060405180830381855af49150503d8060008114611901576040519150601f19603f3d011682016040523d82523d6000602084013e611906565b606091505b50915091508161195e576000815111156119235780518082602001fd5b6040517f2350eb5200000000000000000000000000000000000000000000000000000000815260040161195590613aee565b60405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff16633170428e6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156119a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119cd9190613591565b50600084846000906004926119e493929190613b18565b906119ef9190613b97565b90506358379d7160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611a78576040517f2350eb52000000000000000000000000000000000000000000000000000000008152600401611a6f90613c42565b60405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff16633837a75e600160026040518363ffffffff1660e01b8152600401611ab5929190613cd8565b6020604051808303816000875af1158015611ad4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611af89190613490565b5050505050505050565b60033373ffffffffffffffffffffffffffffffffffffffff1660017f5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf360405160405180910390a4565b6000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905090565b6000611b9e8383610c96565b15611c9b57600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633e41f13584846040518363ffffffff1660e01b8152600401611c009291906135be565b6020604051808303816000875af1158015611c1f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c439190613490565b503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a38183611c949190613408565b9050611cd6565b6040517f2350eb52000000000000000000000000000000000000000000000000000000008152600401611ccd90613d73565b60405180910390fd5b92915050565b6040517fa0c2d2db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60033373ffffffffffffffffffffffffffffffffffffffff1660027f56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d604051611d7c90613ddf565b60405180910390a4565b600060608383915091509250929050565b6060819050919050565b600081600081905550819050919050565b600081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550819050919050565b6000819050919050565b600560009054906101000a900460ff1681565b6000600260008381526020019081526020016000208190555050565b611e41611fa5565b6000828060000190611e539190613dff565b604051602001611e64929190613e92565b6040516020818303038152906040528051906020012090506000602067ffffffffffffffff811115611e9957611e9861226a565b5b6040519080825280601f01601f191660200182016040528015611ecb5781602001600182028036833780820191505090505b50905060005b6020811015611f4657828160208110611eed57611eec61332e565b5b1a60f81b828281518110611f0457611f0361332e565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508080611f3e906133c0565b915050611ed1565b50604051806040016040528085611f5c9061398c565b81526020018281525092505050919050565b6000817fb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c6860405160405180910390a2819050919050565b6040518060400160405280611fb8611ff2565b8152602001606081525090565b60405180604001604052806002905b611fdc611ff2565b815260200190600190039081611fd45790505090565b604051806040016040528060608152602001606081525090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006120378261200c565b9050919050565b6120478161202c565b82525050565b6000819050919050565b6120608161204d565b82525050565b600082825260208201905092915050565b7f5265636569766564204574686572000000000000000000000000000000000000600082015250565b60006120ad600e83612066565b91506120b882612077565b602082019050919050565b60006060820190506120d8600083018561203e565b6120e56020830184612057565b81810360408301526120f6816120a0565b90509392505050565b6000604082019050612114600083018561203e565b6121216020830184612057565b9392505050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b61214f8161213c565b811461215a57600080fd5b50565b60008135905061216c81612146565b92915050565b6000806040838503121561218957612188612132565b5b60006121978582860161215d565b92505060206121a88582860161215d565b9150509250929050565b60008115159050919050565b6121c7816121b2565b82525050565b60006020820190506121e260008301846121be565b92915050565b6121f18161204d565b81146121fc57600080fd5b50565b60008135905061220e816121e8565b92915050565b6000806040838503121561222b5761222a612132565b5b6000612239858286016121ff565b925050602061224a858286016121ff565b9150509250929050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6122a282612259565b810181811067ffffffffffffffff821117156122c1576122c061226a565b5b80604052505050565b60006122d4612128565b90506122e08282612299565b919050565b600067ffffffffffffffff821115612300576122ff61226a565b5b602082029050602081019050919050565b600080fd5b6000612329612324846122e5565b6122ca565b9050808382526020820190506020840283018581111561234c5761234b612311565b5b835b81811015612375578061236188826121ff565b84526020840193505060208101905061234e565b5050509392505050565b600082601f83011261239457612393612254565b5b81356123a4848260208601612316565b91505092915050565b6000602082840312156123c3576123c2612132565b5b600082013567ffffffffffffffff8111156123e1576123e0612137565b5b6123ed8482850161237f565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61242b8161204d565b82525050565b600061243d8383612422565b60208301905092915050565b6000602082019050919050565b6000612461826123f6565b61246b8185612401565b935061247683612412565b8060005b838110156124a757815161248e8882612431565b975061249983612449565b92505060018101905061247a565b5085935050505092915050565b600060208201905081810360008301526124ce8184612456565b905092915050565b6124df8161213c565b82525050565b60006020820190506124fa60008301846124d6565b92915050565b6004811061250d57600080fd5b50565b60008135905061251f81612500565b92915050565b60006020828403121561253b5761253a612132565b5b600061254984828501612510565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6004811061259257612591612552565b5b50565b60008190506125a382612581565b919050565b60006125b382612595565b9050919050565b6125c3816125a8565b82525050565b60006020820190506125de60008301846125ba565b92915050565b60006020820190506125f96000830184612057565b92915050565b6000819050919050565b612612816125ff565b811461261d57600080fd5b50565b60008135905061262f81612609565b92915050565b60006020828403121561264b5761264a612132565b5b600061265984828501612620565b91505092915050565b61266b816125ff565b82525050565b60006020820190506126866000830184612662565b92915050565b600080fd5b600067ffffffffffffffff8211156126ac576126ab61226a565b5b6126b582612259565b9050602081019050919050565b82818337600083830152505050565b60006126e46126df84612691565b6122ca565b905082815260208101848484011115612700576126ff61268c565b5b61270b8482856126c2565b509392505050565b600082601f83011261272857612727612254565b5b81356127388482602086016126d1565b91505092915050565b6000806040838503121561275857612757612132565b5b6000612766858286016121ff565b925050602083013567ffffffffffffffff81111561278757612786612137565b5b61279385828601612713565b9150509250929050565b600081519050919050565b60005b838110156127c65780820151818401526020810190506127ab565b60008484015250505050565b60006127dd8261279d565b6127e78185612066565b93506127f78185602086016127a8565b61280081612259565b840191505092915050565b60006040820190506128206000830185612057565b818103602083015261283281846127d2565b90509392505050565b6128448161202c565b811461284f57600080fd5b50565b6000813590506128618161283b565b92915050565b60006020828403121561287d5761287c612132565b5b600061288b84828501612852565b91505092915050565b6000602082840312156128aa576128a9612132565b5b60006128b88482850161215d565b91505092915050565b600080fd5b6000604082840312156128dc576128db6128c1565b5b81905092915050565b6000602082840312156128fb576128fa612132565b5b600082013567ffffffffffffffff81111561291957612918612137565b5b612925848285016128c6565b91505092915050565b600082825260208201905092915050565b600061294a8261279d565b612954818561292e565b93506129648185602086016127a8565b61296d81612259565b840191505092915050565b600082825260208201905092915050565b6000612994826123f6565b61299e8185612978565b93506129a983612412565b8060005b838110156129da5781516129c18882612431565b97506129cc83612449565b9250506001810190506129ad565b5085935050505092915050565b60006040830160008301518482036000860152612a04828261293f565b91505060208301518482036020860152612a1e8282612989565b9150508091505092915050565b600081519050919050565b600082825260208201905092915050565b6000612a5282612a2b565b612a5c8185612a36565b9350612a6c8185602086016127a8565b612a7581612259565b840191505092915050565b60006040830160008301518482036000860152612a9d82826129e7565b91505060208301518482036020860152612ab78282612a47565b9150508091505092915050565b60006020820190508181036000830152612ade8184612a80565b905092915050565b600060408284031215612afc57612afb6128c1565b5b81905092915050565b600060208284031215612b1b57612b1a612132565b5b600082013567ffffffffffffffff811115612b3957612b38612137565b5b612b4584828501612ae6565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112612b7a57612b79612b58565b5b83810192508235915060208301925067ffffffffffffffff821115612ba257612ba1612b4e565b5b600182023603831315612bb857612bb7612b53565b5b509250929050565b6000612bcc838561292e565b9350612bd98385846126c2565b612be283612259565b840190509392505050565b60008083356001602003843603038112612c0a57612c09612b58565b5b83810192508235915060208301925067ffffffffffffffff821115612c3257612c31612b4e565b5b602082023603831315612c4857612c47612b53565b5b509250929050565b600080fd5b82818337505050565b6000612c6a8385612978565b93507f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115612c9d57612c9c612c50565b5b602083029250612cae838584612c55565b82840190509392505050565b600060408301612ccd6000840184612b5d565b8583036000870152612ce0838284612bc0565b92505050612cf16020840184612bed565b8583036020870152612d04838284612c5e565b925050508091505092915050565b60006020820190508181036000830152612d2c8184612cba565b905092915050565b600081905082602060030282011115612d5057612d4f612311565b5b92915050565b600060208284031215612d6c57612d6b612132565b5b600082013567ffffffffffffffff811115612d8a57612d89612137565b5b612d9684828501612d34565b91505092915050565b600060029050919050565b600081905092915050565b6000819050919050565b6000612dcb83836129e7565b905092915050565b6000602082019050919050565b6000612deb82612d9f565b612df58185612daa565b935083602082028501612e0785612db5565b8060005b85811015612e435784840389528151612e248582612dbf565b9450612e2f83612dd3565b925060208a01995050600181019050612e0b565b50829750879550505050505092915050565b60006020820190508181036000830152612e6f8184612de0565b905092915050565b600081600f0b9050919050565b612e8d81612e77565b8114612e9857600080fd5b50565b600081359050612eaa81612e84565b92915050565b600080600060608486031215612ec957612ec8612132565b5b6000612ed78682870161215d565b9350506020612ee886828701612e9b565b9250506040612ef9868287016121ff565b9150509250925092565b612f0c81612e77565b82525050565b6000606082019050612f2760008301866124d6565b612f346020830185612f03565b612f416040830184612057565b949350505050565b600080fd5b60008083601f840112612f6457612f63612254565b5b8235905067ffffffffffffffff811115612f8157612f80612f49565b5b602083019150836001820283011115612f9d57612f9c612311565b5b9250929050565b60008060008060608587031215612fbe57612fbd612132565b5b6000612fcc87828801612852565b9450506020612fdd878288016121ff565b935050604085013567ffffffffffffffff811115612ffe57612ffd612137565b5b61300a87828801612f4e565b925092505092959194509250565b6000819050919050565b600061303d6130386130338461200c565b613018565b61200c565b9050919050565b600061304f82613022565b9050919050565b600061306182613044565b9050919050565b61307181613056565b82525050565b600060208201905061308c6000830184613068565b92915050565b600067ffffffffffffffff8211156130ad576130ac61226a565b5b602082029050602081019050919050565b60006130d16130cc84613092565b6122ca565b905080838252602082019050602084028301858111156130f4576130f3612311565b5b835b8181101561311d57806131098882612852565b8452602084019350506020810190506130f6565b5050509392505050565b600082601f83011261313c5761313b612254565b5b813561314c8482602086016130be565b91505092915050565b60006020828403121561316b5761316a612132565b5b600082013567ffffffffffffffff81111561318957613188612137565b5b61319584828501613127565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6131d38161202c565b82525050565b60006131e583836131ca565b60208301905092915050565b6000602082019050919050565b60006132098261319e565b61321381856131a9565b935061321e836131ba565b8060005b8381101561324f57815161323688826131d9565b9750613241836131f1565b925050600181019050613222565b5085935050505092915050565b6000602082019050818103600083015261327681846131fe565b905092915050565b6000602082019050613293600083018461203e565b92915050565b7f616c7761797320726576657274206572726f7200000000000000000000000000600082015250565b60006132cf601383612066565b91506132da82613299565b602082019050919050565b600060208201905081810360008301526132fe816132c2565b9050919050565b600060408201905061331a6000830185612057565b6133276020830184612057565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006133978261204d565b91506133a28361204d565b92508282019050808211156133ba576133b961335d565b5b92915050565b60006133cb8261204d565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036133fd576133fc61335d565b5b600182019050919050565b60006134138261213c565b915061341e8361213c565b9250828201905082811215600083121683821260008412151617156134465761344561335d565b5b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b60008151905061348a81612146565b92915050565b6000602082840312156134a6576134a5612132565b5b60006134b48482850161347b565b91505092915050565b600081905092915050565b60006134d382612a2b565b6134dd81856134bd565b93506134ed8185602086016127a8565b80840191505092915050565b600061350582846134c8565b915081905092915050565b7f5374617469632063616c6c206661696c65640000000000000000000000000000600082015250565b6000613546601283612066565b915061355182613510565b602082019050919050565b6000602082019050818103600083015261357581613539565b9050919050565b60008151905061358b816121e8565b92915050565b6000602082840312156135a7576135a6612132565b5b60006135b58482850161357c565b91505092915050565b60006040820190506135d360008301856124d6565b6135e060208301846124d6565b9392505050565b6000819050919050565b600061360c613607613602846135e7565b613018565b61204d565b9050919050565b61361c816135f1565b82525050565b6000819050919050565b600061364761364261363d84613622565b613018565b61204d565b9050919050565b6136578161362c565b82525050565b60006040820190506136726000830185613613565b61367f602083018461364e565b9392505050565b600067ffffffffffffffff82169050919050565b6136a381613686565b82525050565b600060608301600083015184820360008601526136c6828261293f565b91505060208301516136db602086018261369a565b5060408301516136ee6040860182612422565b508091505092915050565b6000602082019050818103600083015261371381846136a9565b905092915050565b7f6d79537472696e67000000000000000000000000000000000000000000000000600082015250565b6000613751600883612066565b915061375c8261371b565b602082019050919050565b6000602082019050818103600083015261378081613744565b9050919050565b600080fd5b600080fd5b6000604082840312156137a7576137a6613787565b5b6137b160406122ca565b9050600082013567ffffffffffffffff8111156137d1576137d061378c565b5b6137dd84828501612713565b600083015250602082013567ffffffffffffffff8111156138015761380061378c565b5b61380d8482850161237f565b60208301525092915050565b600067ffffffffffffffff8211156138345761383361226a565b5b61383d82612259565b9050602081019050919050565b600061385d61385884613819565b6122ca565b9050828152602081018484840111156138795761387861268c565b5b6138848482856126c2565b509392505050565b600082601f8301126138a1576138a0612254565b5b81356138b184826020860161384a565b91505092915050565b6000604082840312156138d0576138cf613787565b5b6138da60406122ca565b9050600082013567ffffffffffffffff8111156138fa576138f961378c565b5b61390684828501613791565b600083015250602082013567ffffffffffffffff81111561392a5761392961378c565b5b6139368482850161388c565b60208301525092915050565b600061394e36836138ba565b9050919050565b600080fd5b600080fd5b600080fd5b6000823560016040038336030381126139805761397f613955565b5b80830191505092915050565b60006139983683613791565b9050919050565b6000819050919050565b60006139c46139bf6139ba8461399f565b613018565b61204d565b9050919050565b6139d4816139a9565b82525050565b6000819050919050565b60006139ff6139fa6139f5846139da565b613018565b61204d565b9050919050565b613a0f816139e4565b82525050565b6000604082019050613a2a60008301856139cb565b613a376020830184613a06565b9392505050565b6000613a4a83856134bd565b9350613a578385846126c2565b82840190509392505050565b6000613a70828486613a3e565b91508190509392505050565b7f64656c656761746563616c6c206661696c65642077697468206e6f207265617360008201527f6f6e000000000000000000000000000000000000000000000000000000000000602082015250565b6000613ad8602283612066565b9150613ae382613a7c565b604082019050919050565b60006020820190508181036000830152613b0781613acb565b9050919050565b600080fd5b600080fd5b60008085851115613b2c57613b2b613b0e565b5b83861115613b3d57613b3c613b13565b5b6001850283019150848603905094509492505050565b600082905092915050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600082821b905092915050565b6000613ba38383613b53565b82613bae8135613b5e565b92506004821015613bee57613be97fffffffff0000000000000000000000000000000000000000000000000000000083600403600802613b8a565b831692505b505092915050565b7f6f68206f68206f682069742773206d6167696321000000000000000000000000600082015250565b6000613c2c601483612066565b9150613c3782613bf6565b602082019050919050565b60006020820190508181036000830152613c5b81613c1f565b9050919050565b6000819050919050565b6000613c87613c82613c7d84613c62565b613018565b61213c565b9050919050565b613c9781613c6c565b82525050565b6000819050919050565b6000613cc2613cbd613cb884613c9d565b613018565b61213c565b9050919050565b613cd281613ca7565b82525050565b6000604082019050613ced6000830185613c8e565b613cfa6020830184613cc9565b9392505050565b7f666972737420696e7420776173206e6f742067726561746572207468616e207360008201527f65636f6e6420696e740000000000000000000000000000000000000000000000602082015250565b6000613d5d602983612066565b9150613d6882613d01565b604082019050919050565b60006020820190508181036000830152613d8c81613d50565b9050919050565b7f736f6d6520696400000000000000000000000000000000000000000000000000600082015250565b6000613dc9600783612066565b9150613dd482613d93565b602082019050919050565b60006020820190508181036000830152613df881613dbc565b9050919050565b60008083356001602003843603038112613e1c57613e1b613955565b5b80840192508235915067ffffffffffffffff821115613e3e57613e3d61395a565b5b602083019250600182023603831315613e5a57613e5961395f565b5b509250929050565b600081905092915050565b6000613e798385613e62565b9350613e868385846126c2565b82840190509392505050565b6000613e9f828486613e6d565b9150819050939250505056fea26469706673582212201d3e2a5deeb452c5209f93da6f6ac4039dcbe2db9159584b889c06892fed4a3764736f6c63430008130033",
}

// NetworkDebugContractABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkDebugContractMetaData.ABI instead.
var NetworkDebugContractABI = NetworkDebugContractMetaData.ABI

// NetworkDebugContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkDebugContractMetaData.Bin instead.
var NetworkDebugContractBin = NetworkDebugContractMetaData.Bin

// DeployNetworkDebugContract deploys a new Ethereum contract, binding an instance of NetworkDebugContract to it.
func DeployNetworkDebugContract(auth *bind.TransactOpts, backend bind.ContractBackend, subAddr common.Address) (common.Address, *types.Transaction, *NetworkDebugContract, error) {
	parsed, err := NetworkDebugContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkDebugContractBin), backend, subAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkDebugContract{NetworkDebugContractCaller: NetworkDebugContractCaller{contract: contract}, NetworkDebugContractTransactor: NetworkDebugContractTransactor{contract: contract}, NetworkDebugContractFilterer: NetworkDebugContractFilterer{contract: contract}}, nil
}

// NetworkDebugContract is an auto generated Go binding around an Ethereum contract.
type NetworkDebugContract struct {
	NetworkDebugContractCaller     // Read-only binding to the contract
	NetworkDebugContractTransactor // Write-only binding to the contract
	NetworkDebugContractFilterer   // Log filterer for contract events
}

// NetworkDebugContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkDebugContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkDebugContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkDebugContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkDebugContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkDebugContractSession struct {
	Contract     *NetworkDebugContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// NetworkDebugContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkDebugContractCallerSession struct {
	Contract *NetworkDebugContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// NetworkDebugContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkDebugContractTransactorSession struct {
	Contract     *NetworkDebugContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// NetworkDebugContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkDebugContractRaw struct {
	Contract *NetworkDebugContract // Generic contract binding to access the raw methods on
}

// NetworkDebugContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkDebugContractCallerRaw struct {
	Contract *NetworkDebugContractCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkDebugContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkDebugContractTransactorRaw struct {
	Contract *NetworkDebugContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkDebugContract creates a new instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContract(address common.Address, backend bind.ContractBackend) (*NetworkDebugContract, error) {
	contract, err := bindNetworkDebugContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContract{NetworkDebugContractCaller: NetworkDebugContractCaller{contract: contract}, NetworkDebugContractTransactor: NetworkDebugContractTransactor{contract: contract}, NetworkDebugContractFilterer: NetworkDebugContractFilterer{contract: contract}}, nil
}

// NewNetworkDebugContractCaller creates a new read-only instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContractCaller(address common.Address, caller bind.ContractCaller) (*NetworkDebugContractCaller, error) {
	contract, err := bindNetworkDebugContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCaller{contract: contract}, nil
}

// NewNetworkDebugContractTransactor creates a new write-only instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContractTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkDebugContractTransactor, error) {
	contract, err := bindNetworkDebugContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractTransactor{contract: contract}, nil
}

// NewNetworkDebugContractFilterer creates a new log filterer instance of NetworkDebugContract, bound to a specific deployed contract.
func NewNetworkDebugContractFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkDebugContractFilterer, error) {
	contract, err := bindNetworkDebugContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractFilterer{contract: contract}, nil
}

// bindNetworkDebugContract binds a generic wrapper to an already deployed contract.
func bindNetworkDebugContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkDebugContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkDebugContract *NetworkDebugContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkDebugContract.Contract.NetworkDebugContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkDebugContract *NetworkDebugContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.NetworkDebugContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkDebugContract *NetworkDebugContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.NetworkDebugContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkDebugContract *NetworkDebugContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkDebugContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkDebugContract *NetworkDebugContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkDebugContract *NetworkDebugContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.contract.Transact(opts, method, params...)
}

// CounterMap is a free data retrieval call binding the contract method 0x6284117d.
//
// Solidity: function counterMap(int256 ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCaller) CounterMap(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "counterMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CounterMap is a free data retrieval call binding the contract method 0x6284117d.
//
// Solidity: function counterMap(int256 ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) CounterMap(arg0 *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.CounterMap(&_NetworkDebugContract.CallOpts, arg0)
}

// CounterMap is a free data retrieval call binding the contract method 0x6284117d.
//
// Solidity: function counterMap(int256 ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) CounterMap(arg0 *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.CounterMap(&_NetworkDebugContract.CallOpts, arg0)
}

// CurrentStatus is a free data retrieval call binding the contract method 0xef8a9235.
//
// Solidity: function currentStatus() view returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractCaller) CurrentStatus(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "currentStatus")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CurrentStatus is a free data retrieval call binding the contract method 0xef8a9235.
//
// Solidity: function currentStatus() view returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractSession) CurrentStatus() (uint8, error) {
	return _NetworkDebugContract.Contract.CurrentStatus(&_NetworkDebugContract.CallOpts)
}

// CurrentStatus is a free data retrieval call binding the contract method 0xef8a9235.
//
// Solidity: function currentStatus() view returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) CurrentStatus() (uint8, error) {
	return _NetworkDebugContract.Contract.CurrentStatus(&_NetworkDebugContract.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCaller) Get(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "get")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractSession) Get() (*big.Int, error) {
	return _NetworkDebugContract.Contract.Get(&_NetworkDebugContract.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) Get() (*big.Int, error) {
	return _NetworkDebugContract.Contract.Get(&_NetworkDebugContract.CallOpts)
}

// GetCounter is a free data retrieval call binding the contract method 0x5921483f.
//
// Solidity: function getCounter(int256 idx) view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCaller) GetCounter(opts *bind.CallOpts, idx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "getCounter", idx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCounter is a free data retrieval call binding the contract method 0x5921483f.
//
// Solidity: function getCounter(int256 idx) view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractSession) GetCounter(idx *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetCounter(&_NetworkDebugContract.CallOpts, idx)
}

// GetCounter is a free data retrieval call binding the contract method 0x5921483f.
//
// Solidity: function getCounter(int256 idx) view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) GetCounter(idx *big.Int) (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetCounter(&_NetworkDebugContract.CallOpts, idx)
}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCaller) GetData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "getData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractSession) GetData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetData(&_NetworkDebugContract.CallOpts)
}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) GetData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetData(&_NetworkDebugContract.CallOpts)
}

// GetMap is a free data retrieval call binding the contract method 0xad3de14c.
//
// Solidity: function getMap() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCaller) GetMap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "getMap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMap is a free data retrieval call binding the contract method 0xad3de14c.
//
// Solidity: function getMap() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractSession) GetMap() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetMap(&_NetworkDebugContract.CallOpts)
}

// GetMap is a free data retrieval call binding the contract method 0xad3de14c.
//
// Solidity: function getMap() view returns(int256 data)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) GetMap() (*big.Int, error) {
	return _NetworkDebugContract.Contract.GetMap(&_NetworkDebugContract.CallOpts)
}

// PerformStaticCall is a free data retrieval call binding the contract method 0x3170428e.
//
// Solidity: function performStaticCall() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCaller) PerformStaticCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "performStaticCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerformStaticCall is a free data retrieval call binding the contract method 0x3170428e.
//
// Solidity: function performStaticCall() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractSession) PerformStaticCall() (*big.Int, error) {
	return _NetworkDebugContract.Contract.PerformStaticCall(&_NetworkDebugContract.CallOpts)
}

// PerformStaticCall is a free data retrieval call binding the contract method 0x3170428e.
//
// Solidity: function performStaticCall() view returns(uint256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) PerformStaticCall() (*big.Int, error) {
	return _NetworkDebugContract.Contract.PerformStaticCall(&_NetworkDebugContract.CallOpts)
}

// StoredData is a free data retrieval call binding the contract method 0x2a1afcd9.
//
// Solidity: function storedData() view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCaller) StoredData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "storedData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StoredData is a free data retrieval call binding the contract method 0x2a1afcd9.
//
// Solidity: function storedData() view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) StoredData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredData(&_NetworkDebugContract.CallOpts)
}

// StoredData is a free data retrieval call binding the contract method 0x2a1afcd9.
//
// Solidity: function storedData() view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) StoredData() (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredData(&_NetworkDebugContract.CallOpts)
}

// StoredDataMap is a free data retrieval call binding the contract method 0x48ad9fe8.
//
// Solidity: function storedDataMap(address ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCaller) StoredDataMap(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "storedDataMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StoredDataMap is a free data retrieval call binding the contract method 0x48ad9fe8.
//
// Solidity: function storedDataMap(address ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) StoredDataMap(arg0 common.Address) (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredDataMap(&_NetworkDebugContract.CallOpts, arg0)
}

// StoredDataMap is a free data retrieval call binding the contract method 0x48ad9fe8.
//
// Solidity: function storedDataMap(address ) view returns(int256)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) StoredDataMap(arg0 common.Address) (*big.Int, error) {
	return _NetworkDebugContract.Contract.StoredDataMap(&_NetworkDebugContract.CallOpts, arg0)
}

// SubContract is a free data retrieval call binding the contract method 0xc0d06d89.
//
// Solidity: function subContract() view returns(address)
func (_NetworkDebugContract *NetworkDebugContractCaller) SubContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkDebugContract.contract.Call(opts, &out, "subContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SubContract is a free data retrieval call binding the contract method 0xc0d06d89.
//
// Solidity: function subContract() view returns(address)
func (_NetworkDebugContract *NetworkDebugContractSession) SubContract() (common.Address, error) {
	return _NetworkDebugContract.Contract.SubContract(&_NetworkDebugContract.CallOpts)
}

// SubContract is a free data retrieval call binding the contract method 0xc0d06d89.
//
// Solidity: function subContract() view returns(address)
func (_NetworkDebugContract *NetworkDebugContractCallerSession) SubContract() (common.Address, error) {
	return _NetworkDebugContract.Contract.SubContract(&_NetworkDebugContract.CallOpts)
}

// AddCounter is a paid mutator transaction binding the contract method 0x23515760.
//
// Solidity: function addCounter(int256 idx, int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactor) AddCounter(opts *bind.TransactOpts, idx *big.Int, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "addCounter", idx, x)
}

// AddCounter is a paid mutator transaction binding the contract method 0x23515760.
//
// Solidity: function addCounter(int256 idx, int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractSession) AddCounter(idx *big.Int, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AddCounter(&_NetworkDebugContract.TransactOpts, idx, x)
}

// AddCounter is a paid mutator transaction binding the contract method 0x23515760.
//
// Solidity: function addCounter(int256 idx, int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AddCounter(idx *big.Int, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AddCounter(&_NetworkDebugContract.TransactOpts, idx, x)
}

// AlwaysRevertsAssert is a paid mutator transaction binding the contract method 0x256560d5.
//
// Solidity: function alwaysRevertsAssert() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsAssert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsAssert")
}

// AlwaysRevertsAssert is a paid mutator transaction binding the contract method 0x256560d5.
//
// Solidity: function alwaysRevertsAssert() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsAssert() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsAssert(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsAssert is a paid mutator transaction binding the contract method 0x256560d5.
//
// Solidity: function alwaysRevertsAssert() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsAssert() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsAssert(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x5e9c80d6.
//
// Solidity: function alwaysRevertsCustomError() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsCustomError(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsCustomError")
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x5e9c80d6.
//
// Solidity: function alwaysRevertsCustomError() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsCustomError() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomError(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomError is a paid mutator transaction binding the contract method 0x5e9c80d6.
//
// Solidity: function alwaysRevertsCustomError() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsCustomError() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomError(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomErrorNoValues is a paid mutator transaction binding the contract method 0xb600141f.
//
// Solidity: function alwaysRevertsCustomErrorNoValues() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsCustomErrorNoValues(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsCustomErrorNoValues")
}

// AlwaysRevertsCustomErrorNoValues is a paid mutator transaction binding the contract method 0xb600141f.
//
// Solidity: function alwaysRevertsCustomErrorNoValues() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsCustomErrorNoValues() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomErrorNoValues(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsCustomErrorNoValues is a paid mutator transaction binding the contract method 0xb600141f.
//
// Solidity: function alwaysRevertsCustomErrorNoValues() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsCustomErrorNoValues() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsCustomErrorNoValues(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsRequire is a paid mutator transaction binding the contract method 0x06595f75.
//
// Solidity: function alwaysRevertsRequire() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) AlwaysRevertsRequire(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "alwaysRevertsRequire")
}

// AlwaysRevertsRequire is a paid mutator transaction binding the contract method 0x06595f75.
//
// Solidity: function alwaysRevertsRequire() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) AlwaysRevertsRequire() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsRequire(&_NetworkDebugContract.TransactOpts)
}

// AlwaysRevertsRequire is a paid mutator transaction binding the contract method 0x06595f75.
//
// Solidity: function alwaysRevertsRequire() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) AlwaysRevertsRequire() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.AlwaysRevertsRequire(&_NetworkDebugContract.TransactOpts)
}

// CallRevertFunctionInSubContract is a paid mutator transaction binding the contract method 0x11b3c478.
//
// Solidity: function callRevertFunctionInSubContract(uint256 x, uint256 y) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) CallRevertFunctionInSubContract(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "callRevertFunctionInSubContract", x, y)
}

// CallRevertFunctionInSubContract is a paid mutator transaction binding the contract method 0x11b3c478.
//
// Solidity: function callRevertFunctionInSubContract(uint256 x, uint256 y) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) CallRevertFunctionInSubContract(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInSubContract(&_NetworkDebugContract.TransactOpts, x, y)
}

// CallRevertFunctionInSubContract is a paid mutator transaction binding the contract method 0x11b3c478.
//
// Solidity: function callRevertFunctionInSubContract(uint256 x, uint256 y) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) CallRevertFunctionInSubContract(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInSubContract(&_NetworkDebugContract.TransactOpts, x, y)
}

// CallRevertFunctionInTheContract is a paid mutator transaction binding the contract method 0x9349d00b.
//
// Solidity: function callRevertFunctionInTheContract() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) CallRevertFunctionInTheContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "callRevertFunctionInTheContract")
}

// CallRevertFunctionInTheContract is a paid mutator transaction binding the contract method 0x9349d00b.
//
// Solidity: function callRevertFunctionInTheContract() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) CallRevertFunctionInTheContract() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInTheContract(&_NetworkDebugContract.TransactOpts)
}

// CallRevertFunctionInTheContract is a paid mutator transaction binding the contract method 0x9349d00b.
//
// Solidity: function callRevertFunctionInTheContract() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) CallRevertFunctionInTheContract() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallRevertFunctionInTheContract(&_NetworkDebugContract.TransactOpts)
}

// CallbackMethod is a paid mutator transaction binding the contract method 0xfbcb8d07.
//
// Solidity: function callbackMethod(int256 x) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) CallbackMethod(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "callbackMethod", x)
}

// CallbackMethod is a paid mutator transaction binding the contract method 0xfbcb8d07.
//
// Solidity: function callbackMethod(int256 x) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) CallbackMethod(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallbackMethod(&_NetworkDebugContract.TransactOpts, x)
}

// CallbackMethod is a paid mutator transaction binding the contract method 0xfbcb8d07.
//
// Solidity: function callbackMethod(int256 x) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) CallbackMethod(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.CallbackMethod(&_NetworkDebugContract.TransactOpts, x)
}

// EmitAddress is a paid mutator transaction binding the contract method 0xec5c3ede.
//
// Solidity: function emitAddress(address addr) returns(address)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitAddress", addr)
}

// EmitAddress is a paid mutator transaction binding the contract method 0xec5c3ede.
//
// Solidity: function emitAddress(address addr) returns(address)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitAddress(addr common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitAddress(&_NetworkDebugContract.TransactOpts, addr)
}

// EmitAddress is a paid mutator transaction binding the contract method 0xec5c3ede.
//
// Solidity: function emitAddress(address addr) returns(address)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitAddress(addr common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitAddress(&_NetworkDebugContract.TransactOpts, addr)
}

// EmitBytes32 is a paid mutator transaction binding the contract method 0x33311ef3.
//
// Solidity: function emitBytes32(bytes32 input) returns(bytes32 output)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitBytes32(opts *bind.TransactOpts, input [32]byte) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitBytes32", input)
}

// EmitBytes32 is a paid mutator transaction binding the contract method 0x33311ef3.
//
// Solidity: function emitBytes32(bytes32 input) returns(bytes32 output)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitBytes32(input [32]byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitBytes32(&_NetworkDebugContract.TransactOpts, input)
}

// EmitBytes32 is a paid mutator transaction binding the contract method 0x33311ef3.
//
// Solidity: function emitBytes32(bytes32 input) returns(bytes32 output)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitBytes32(input [32]byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitBytes32(&_NetworkDebugContract.TransactOpts, input)
}

// EmitFourParamMixedEvent is a paid mutator transaction binding the contract method 0xc2124b22.
//
// Solidity: function emitFourParamMixedEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitFourParamMixedEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitFourParamMixedEvent")
}

// EmitFourParamMixedEvent is a paid mutator transaction binding the contract method 0xc2124b22.
//
// Solidity: function emitFourParamMixedEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitFourParamMixedEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitFourParamMixedEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitFourParamMixedEvent is a paid mutator transaction binding the contract method 0xc2124b22.
//
// Solidity: function emitFourParamMixedEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitFourParamMixedEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitFourParamMixedEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitInputs is a paid mutator transaction binding the contract method 0x81b375a0.
//
// Solidity: function emitInputs(uint256 inputVal1, string inputVal2) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitInputs(opts *bind.TransactOpts, inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitInputs", inputVal1, inputVal2)
}

// EmitInputs is a paid mutator transaction binding the contract method 0x81b375a0.
//
// Solidity: function emitInputs(uint256 inputVal1, string inputVal2) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitInputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInputs is a paid mutator transaction binding the contract method 0x81b375a0.
//
// Solidity: function emitInputs(uint256 inputVal1, string inputVal2) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitInputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInputsOutputs is a paid mutator transaction binding the contract method 0xd7a80205.
//
// Solidity: function emitInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitInputsOutputs(opts *bind.TransactOpts, inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitInputsOutputs", inputVal1, inputVal2)
}

// EmitInputsOutputs is a paid mutator transaction binding the contract method 0xd7a80205.
//
// Solidity: function emitInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInputsOutputs is a paid mutator transaction binding the contract method 0xd7a80205.
//
// Solidity: function emitInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitInts is a paid mutator transaction binding the contract method 0x9e099652.
//
// Solidity: function emitInts(int256 first, int128 second, uint256 third) returns(int256, int128 outputVal1, uint256 outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitInts(opts *bind.TransactOpts, first *big.Int, second *big.Int, third *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitInts", first, second, third)
}

// EmitInts is a paid mutator transaction binding the contract method 0x9e099652.
//
// Solidity: function emitInts(int256 first, int128 second, uint256 third) returns(int256, int128 outputVal1, uint256 outputVal2)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitInts(first *big.Int, second *big.Int, third *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInts(&_NetworkDebugContract.TransactOpts, first, second, third)
}

// EmitInts is a paid mutator transaction binding the contract method 0x9e099652.
//
// Solidity: function emitInts(int256 first, int128 second, uint256 third) returns(int256, int128 outputVal1, uint256 outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitInts(first *big.Int, second *big.Int, third *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitInts(&_NetworkDebugContract.TransactOpts, first, second, third)
}

// EmitNamedInputsOutputs is a paid mutator transaction binding the contract method 0x45f0c9e6.
//
// Solidity: function emitNamedInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNamedInputsOutputs(opts *bind.TransactOpts, inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNamedInputsOutputs", inputVal1, inputVal2)
}

// EmitNamedInputsOutputs is a paid mutator transaction binding the contract method 0x45f0c9e6.
//
// Solidity: function emitNamedInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNamedInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitNamedInputsOutputs is a paid mutator transaction binding the contract method 0x45f0c9e6.
//
// Solidity: function emitNamedInputsOutputs(uint256 inputVal1, string inputVal2) returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNamedInputsOutputs(inputVal1 *big.Int, inputVal2 string) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedInputsOutputs(&_NetworkDebugContract.TransactOpts, inputVal1, inputVal2)
}

// EmitNamedOutputs is a paid mutator transaction binding the contract method 0x7014c81d.
//
// Solidity: function emitNamedOutputs() returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNamedOutputs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNamedOutputs")
}

// EmitNamedOutputs is a paid mutator transaction binding the contract method 0x7014c81d.
//
// Solidity: function emitNamedOutputs() returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNamedOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitNamedOutputs is a paid mutator transaction binding the contract method 0x7014c81d.
//
// Solidity: function emitNamedOutputs() returns(uint256 outputVal1, string outputVal2)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNamedOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNamedOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEvent is a paid mutator transaction binding the contract method 0x95a81a4c.
//
// Solidity: function emitNoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNoIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNoIndexEvent")
}

// EmitNoIndexEvent is a paid mutator transaction binding the contract method 0x95a81a4c.
//
// Solidity: function emitNoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEvent is a paid mutator transaction binding the contract method 0x95a81a4c.
//
// Solidity: function emitNoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEventString is a paid mutator transaction binding the contract method 0x788c4772.
//
// Solidity: function emitNoIndexEventString() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNoIndexEventString(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNoIndexEventString")
}

// EmitNoIndexEventString is a paid mutator transaction binding the contract method 0x788c4772.
//
// Solidity: function emitNoIndexEventString() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNoIndexEventString() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEventString(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexEventString is a paid mutator transaction binding the contract method 0x788c4772.
//
// Solidity: function emitNoIndexEventString() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNoIndexEventString() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexEventString(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexStructEvent is a paid mutator transaction binding the contract method 0x62c270e1.
//
// Solidity: function emitNoIndexStructEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitNoIndexStructEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitNoIndexStructEvent")
}

// EmitNoIndexStructEvent is a paid mutator transaction binding the contract method 0x62c270e1.
//
// Solidity: function emitNoIndexStructEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitNoIndexStructEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexStructEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitNoIndexStructEvent is a paid mutator transaction binding the contract method 0x62c270e1.
//
// Solidity: function emitNoIndexStructEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitNoIndexStructEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitNoIndexStructEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitOneIndexEvent is a paid mutator transaction binding the contract method 0x8f856296.
//
// Solidity: function emitOneIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitOneIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitOneIndexEvent")
}

// EmitOneIndexEvent is a paid mutator transaction binding the contract method 0x8f856296.
//
// Solidity: function emitOneIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitOneIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOneIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitOneIndexEvent is a paid mutator transaction binding the contract method 0x8f856296.
//
// Solidity: function emitOneIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitOneIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOneIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitOutputs is a paid mutator transaction binding the contract method 0x8db611be.
//
// Solidity: function emitOutputs() returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitOutputs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitOutputs")
}

// EmitOutputs is a paid mutator transaction binding the contract method 0x8db611be.
//
// Solidity: function emitOutputs() returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractSession) EmitOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitOutputs is a paid mutator transaction binding the contract method 0x8db611be.
//
// Solidity: function emitOutputs() returns(uint256, string)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitOutputs() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitOutputs(&_NetworkDebugContract.TransactOpts)
}

// EmitThreeIndexEvent is a paid mutator transaction binding the contract method 0xaa3fdcf4.
//
// Solidity: function emitThreeIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitThreeIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitThreeIndexEvent")
}

// EmitThreeIndexEvent is a paid mutator transaction binding the contract method 0xaa3fdcf4.
//
// Solidity: function emitThreeIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitThreeIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitThreeIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitThreeIndexEvent is a paid mutator transaction binding the contract method 0xaa3fdcf4.
//
// Solidity: function emitThreeIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitThreeIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitThreeIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitTwoIndexEvent is a paid mutator transaction binding the contract method 0x1e31d0a8.
//
// Solidity: function emitTwoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) EmitTwoIndexEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "emitTwoIndexEvent")
}

// EmitTwoIndexEvent is a paid mutator transaction binding the contract method 0x1e31d0a8.
//
// Solidity: function emitTwoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) EmitTwoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitTwoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// EmitTwoIndexEvent is a paid mutator transaction binding the contract method 0x1e31d0a8.
//
// Solidity: function emitTwoIndexEvent() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) EmitTwoIndexEvent() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.EmitTwoIndexEvent(&_NetworkDebugContract.TransactOpts)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.OnTokenTransfer(&_NetworkDebugContract.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.OnTokenTransfer(&_NetworkDebugContract.TransactOpts, sender, amount, data)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) Pay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "pay")
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugContract *NetworkDebugContractSession) Pay() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Pay(&_NetworkDebugContract.TransactOpts)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Pay() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Pay(&_NetworkDebugContract.TransactOpts)
}

// ProcessAddressArray is a paid mutator transaction binding the contract method 0xe1111f79.
//
// Solidity: function processAddressArray(address[] input) returns(address[])
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessAddressArray(opts *bind.TransactOpts, input []common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processAddressArray", input)
}

// ProcessAddressArray is a paid mutator transaction binding the contract method 0xe1111f79.
//
// Solidity: function processAddressArray(address[] input) returns(address[])
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessAddressArray(input []common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessAddressArray(&_NetworkDebugContract.TransactOpts, input)
}

// ProcessAddressArray is a paid mutator transaction binding the contract method 0xe1111f79.
//
// Solidity: function processAddressArray(address[] input) returns(address[])
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessAddressArray(input []common.Address) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessAddressArray(&_NetworkDebugContract.TransactOpts, input)
}

// ProcessDynamicData is a paid mutator transaction binding the contract method 0x7fdc8fe1.
//
// Solidity: function processDynamicData((string,uint256[]) data) returns((string,uint256[]))
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessDynamicData(opts *bind.TransactOpts, data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processDynamicData", data)
}

// ProcessDynamicData is a paid mutator transaction binding the contract method 0x7fdc8fe1.
//
// Solidity: function processDynamicData((string,uint256[]) data) returns((string,uint256[]))
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessDynamicData(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessDynamicData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessDynamicData is a paid mutator transaction binding the contract method 0x7fdc8fe1.
//
// Solidity: function processDynamicData((string,uint256[]) data) returns((string,uint256[]))
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessDynamicData(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessDynamicData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessFixedDataArray is a paid mutator transaction binding the contract method 0x99adad2e.
//
// Solidity: function processFixedDataArray((string,uint256[])[3] data) returns((string,uint256[])[2])
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessFixedDataArray(opts *bind.TransactOpts, data [3]NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processFixedDataArray", data)
}

// ProcessFixedDataArray is a paid mutator transaction binding the contract method 0x99adad2e.
//
// Solidity: function processFixedDataArray((string,uint256[])[3] data) returns((string,uint256[])[2])
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessFixedDataArray(data [3]NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessFixedDataArray(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessFixedDataArray is a paid mutator transaction binding the contract method 0x99adad2e.
//
// Solidity: function processFixedDataArray((string,uint256[])[3] data) returns((string,uint256[])[2])
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessFixedDataArray(data [3]NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessFixedDataArray(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData is a paid mutator transaction binding the contract method 0x7f12881c.
//
// Solidity: function processNestedData(((string,uint256[]),bytes) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessNestedData(opts *bind.TransactOpts, data NetworkDebugContractNestedData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processNestedData", data)
}

// ProcessNestedData is a paid mutator transaction binding the contract method 0x7f12881c.
//
// Solidity: function processNestedData(((string,uint256[]),bytes) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessNestedData(data NetworkDebugContractNestedData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData is a paid mutator transaction binding the contract method 0x7f12881c.
//
// Solidity: function processNestedData(((string,uint256[]),bytes) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessNestedData(data NetworkDebugContractNestedData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData0 is a paid mutator transaction binding the contract method 0xf499af2a.
//
// Solidity: function processNestedData((string,uint256[]) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessNestedData0(opts *bind.TransactOpts, data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processNestedData0", data)
}

// ProcessNestedData0 is a paid mutator transaction binding the contract method 0xf499af2a.
//
// Solidity: function processNestedData((string,uint256[]) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessNestedData0(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData0(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessNestedData0 is a paid mutator transaction binding the contract method 0xf499af2a.
//
// Solidity: function processNestedData((string,uint256[]) data) returns(((string,uint256[]),bytes))
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessNestedData0(data NetworkDebugContractData) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessNestedData0(&_NetworkDebugContract.TransactOpts, data)
}

// ProcessUintArray is a paid mutator transaction binding the contract method 0x12d91233.
//
// Solidity: function processUintArray(uint256[] input) returns(uint256[])
func (_NetworkDebugContract *NetworkDebugContractTransactor) ProcessUintArray(opts *bind.TransactOpts, input []*big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "processUintArray", input)
}

// ProcessUintArray is a paid mutator transaction binding the contract method 0x12d91233.
//
// Solidity: function processUintArray(uint256[] input) returns(uint256[])
func (_NetworkDebugContract *NetworkDebugContractSession) ProcessUintArray(input []*big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessUintArray(&_NetworkDebugContract.TransactOpts, input)
}

// ProcessUintArray is a paid mutator transaction binding the contract method 0x12d91233.
//
// Solidity: function processUintArray(uint256[] input) returns(uint256[])
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ProcessUintArray(input []*big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ProcessUintArray(&_NetworkDebugContract.TransactOpts, input)
}

// ResetCounter is a paid mutator transaction binding the contract method 0xf3396bd9.
//
// Solidity: function resetCounter(int256 idx) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) ResetCounter(opts *bind.TransactOpts, idx *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "resetCounter", idx)
}

// ResetCounter is a paid mutator transaction binding the contract method 0xf3396bd9.
//
// Solidity: function resetCounter(int256 idx) returns()
func (_NetworkDebugContract *NetworkDebugContractSession) ResetCounter(idx *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ResetCounter(&_NetworkDebugContract.TransactOpts, idx)
}

// ResetCounter is a paid mutator transaction binding the contract method 0xf3396bd9.
//
// Solidity: function resetCounter(int256 idx) returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) ResetCounter(idx *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.ResetCounter(&_NetworkDebugContract.TransactOpts, idx)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactor) Set(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "set", x)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractSession) Set(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Set(&_NetworkDebugContract.TransactOpts, x)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Set(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Set(&_NetworkDebugContract.TransactOpts, x)
}

// SetMap is a paid mutator transaction binding the contract method 0xe8116e28.
//
// Solidity: function setMap(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactor) SetMap(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "setMap", x)
}

// SetMap is a paid mutator transaction binding the contract method 0xe8116e28.
//
// Solidity: function setMap(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractSession) SetMap(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetMap(&_NetworkDebugContract.TransactOpts, x)
}

// SetMap is a paid mutator transaction binding the contract method 0xe8116e28.
//
// Solidity: function setMap(int256 x) returns(int256 value)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) SetMap(x *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetMap(&_NetworkDebugContract.TransactOpts, x)
}

// SetStatus is a paid mutator transaction binding the contract method 0x2e49d78b.
//
// Solidity: function setStatus(uint8 status) returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractTransactor) SetStatus(opts *bind.TransactOpts, status uint8) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "setStatus", status)
}

// SetStatus is a paid mutator transaction binding the contract method 0x2e49d78b.
//
// Solidity: function setStatus(uint8 status) returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractSession) SetStatus(status uint8) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetStatus(&_NetworkDebugContract.TransactOpts, status)
}

// SetStatus is a paid mutator transaction binding the contract method 0x2e49d78b.
//
// Solidity: function setStatus(uint8 status) returns(uint8)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) SetStatus(status uint8) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.SetStatus(&_NetworkDebugContract.TransactOpts, status)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) Trace(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "trace", x, y)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) Trace(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Trace(&_NetworkDebugContract.TransactOpts, x, y)
}

// Trace is a paid mutator transaction binding the contract method 0x3e41f135.
//
// Solidity: function trace(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Trace(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Trace(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceDifferent is a paid mutator transaction binding the contract method 0x30985bcc.
//
// Solidity: function traceDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceDifferent(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceDifferent", x, y)
}

// TraceDifferent is a paid mutator transaction binding the contract method 0x30985bcc.
//
// Solidity: function traceDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceDifferent is a paid mutator transaction binding the contract method 0x30985bcc.
//
// Solidity: function traceDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceSubWithCallback is a paid mutator transaction binding the contract method 0x3837a75e.
//
// Solidity: function traceSubWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceSubWithCallback(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceSubWithCallback", x, y)
}

// TraceSubWithCallback is a paid mutator transaction binding the contract method 0x3837a75e.
//
// Solidity: function traceSubWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceSubWithCallback(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceSubWithCallback(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceSubWithCallback is a paid mutator transaction binding the contract method 0x3837a75e.
//
// Solidity: function traceSubWithCallback(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceSubWithCallback(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceSubWithCallback(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceWithValidate is a paid mutator transaction binding the contract method 0xb1ae9d85.
//
// Solidity: function traceWithValidate(int256 x, int256 y) payable returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceWithValidate(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceWithValidate", x, y)
}

// TraceWithValidate is a paid mutator transaction binding the contract method 0xb1ae9d85.
//
// Solidity: function traceWithValidate(int256 x, int256 y) payable returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceWithValidate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceWithValidate(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceWithValidate is a paid mutator transaction binding the contract method 0xb1ae9d85.
//
// Solidity: function traceWithValidate(int256 x, int256 y) payable returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceWithValidate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceWithValidate(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceYetDifferent is a paid mutator transaction binding the contract method 0x58379d71.
//
// Solidity: function traceYetDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceYetDifferent(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceYetDifferent", x, y)
}

// TraceYetDifferent is a paid mutator transaction binding the contract method 0x58379d71.
//
// Solidity: function traceYetDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractSession) TraceYetDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceYetDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// TraceYetDifferent is a paid mutator transaction binding the contract method 0x58379d71.
//
// Solidity: function traceYetDifferent(int256 x, int256 y) returns(int256)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceYetDifferent(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceYetDifferent(&_NetworkDebugContract.TransactOpts, x, y)
}

// Validate is a paid mutator transaction binding the contract method 0x04d8215b.
//
// Solidity: function validate(int256 x, int256 y) returns(bool)
func (_NetworkDebugContract *NetworkDebugContractTransactor) Validate(opts *bind.TransactOpts, x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "validate", x, y)
}

// Validate is a paid mutator transaction binding the contract method 0x04d8215b.
//
// Solidity: function validate(int256 x, int256 y) returns(bool)
func (_NetworkDebugContract *NetworkDebugContractSession) Validate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Validate(&_NetworkDebugContract.TransactOpts, x, y)
}

// Validate is a paid mutator transaction binding the contract method 0x04d8215b.
//
// Solidity: function validate(int256 x, int256 y) returns(bool)
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Validate(x *big.Int, y *big.Int) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Validate(&_NetworkDebugContract.TransactOpts, x, y)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_NetworkDebugContract *NetworkDebugContractSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Fallback(&_NetworkDebugContract.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Fallback(&_NetworkDebugContract.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkDebugContract *NetworkDebugContractSession) Receive() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Receive(&_NetworkDebugContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) Receive() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.Receive(&_NetworkDebugContract.TransactOpts)
}

// NetworkDebugContractCallDataLengthIterator is returned from FilterCallDataLength and is used to iterate over the raw logs and unpacked data for CallDataLength events raised by the NetworkDebugContract contract.
type NetworkDebugContractCallDataLengthIterator struct {
	Event *NetworkDebugContractCallDataLength // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractCallDataLengthIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractCallDataLength)
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
		it.Event = new(NetworkDebugContractCallDataLength)
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
func (it *NetworkDebugContractCallDataLengthIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractCallDataLengthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractCallDataLength represents a CallDataLength event raised by the NetworkDebugContract contract.
type NetworkDebugContractCallDataLength struct {
	Length *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCallDataLength is a free log retrieval operation binding the contract event 0x962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed.
//
// Solidity: event CallDataLength(uint256 length)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterCallDataLength(opts *bind.FilterOpts) (*NetworkDebugContractCallDataLengthIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "CallDataLength")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCallDataLengthIterator{contract: _NetworkDebugContract.contract, event: "CallDataLength", logs: logs, sub: sub}, nil
}

// WatchCallDataLength is a free log subscription operation binding the contract event 0x962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed.
//
// Solidity: event CallDataLength(uint256 length)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchCallDataLength(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractCallDataLength) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "CallDataLength")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractCallDataLength)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "CallDataLength", log); err != nil {
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

// ParseCallDataLength is a log parse operation binding the contract event 0x962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed.
//
// Solidity: event CallDataLength(uint256 length)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseCallDataLength(log types.Log) (*NetworkDebugContractCallDataLength, error) {
	event := new(NetworkDebugContractCallDataLength)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "CallDataLength", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractCallbackEventIterator is returned from FilterCallbackEvent and is used to iterate over the raw logs and unpacked data for CallbackEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractCallbackEventIterator struct {
	Event *NetworkDebugContractCallbackEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractCallbackEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractCallbackEvent)
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
		it.Event = new(NetworkDebugContractCallbackEvent)
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
func (it *NetworkDebugContractCallbackEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractCallbackEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractCallbackEvent represents a CallbackEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractCallbackEvent struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCallbackEvent is a free log retrieval operation binding the contract event 0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68.
//
// Solidity: event CallbackEvent(int256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterCallbackEvent(opts *bind.FilterOpts, a []*big.Int) (*NetworkDebugContractCallbackEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "CallbackEvent", aRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCallbackEventIterator{contract: _NetworkDebugContract.contract, event: "CallbackEvent", logs: logs, sub: sub}, nil
}

// WatchCallbackEvent is a free log subscription operation binding the contract event 0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68.
//
// Solidity: event CallbackEvent(int256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchCallbackEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractCallbackEvent, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "CallbackEvent", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractCallbackEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "CallbackEvent", log); err != nil {
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

// ParseCallbackEvent is a log parse operation binding the contract event 0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68.
//
// Solidity: event CallbackEvent(int256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseCallbackEvent(log types.Log) (*NetworkDebugContractCallbackEvent, error) {
	event := new(NetworkDebugContractCallbackEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "CallbackEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractCurrentStatusIterator is returned from FilterCurrentStatus and is used to iterate over the raw logs and unpacked data for CurrentStatus events raised by the NetworkDebugContract contract.
type NetworkDebugContractCurrentStatusIterator struct {
	Event *NetworkDebugContractCurrentStatus // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractCurrentStatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractCurrentStatus)
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
		it.Event = new(NetworkDebugContractCurrentStatus)
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
func (it *NetworkDebugContractCurrentStatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractCurrentStatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractCurrentStatus represents a CurrentStatus event raised by the NetworkDebugContract contract.
type NetworkDebugContractCurrentStatus struct {
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCurrentStatus is a free log retrieval operation binding the contract event 0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a.
//
// Solidity: event CurrentStatus(uint8 indexed status)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterCurrentStatus(opts *bind.FilterOpts, status []uint8) (*NetworkDebugContractCurrentStatusIterator, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "CurrentStatus", statusRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractCurrentStatusIterator{contract: _NetworkDebugContract.contract, event: "CurrentStatus", logs: logs, sub: sub}, nil
}

// WatchCurrentStatus is a free log subscription operation binding the contract event 0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a.
//
// Solidity: event CurrentStatus(uint8 indexed status)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchCurrentStatus(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractCurrentStatus, status []uint8) (event.Subscription, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "CurrentStatus", statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractCurrentStatus)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "CurrentStatus", log); err != nil {
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

// ParseCurrentStatus is a log parse operation binding the contract event 0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a.
//
// Solidity: event CurrentStatus(uint8 indexed status)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseCurrentStatus(log types.Log) (*NetworkDebugContractCurrentStatus, error) {
	event := new(NetworkDebugContractCurrentStatus)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "CurrentStatus", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractEtherReceivedIterator is returned from FilterEtherReceived and is used to iterate over the raw logs and unpacked data for EtherReceived events raised by the NetworkDebugContract contract.
type NetworkDebugContractEtherReceivedIterator struct {
	Event *NetworkDebugContractEtherReceived // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractEtherReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractEtherReceived)
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
		it.Event = new(NetworkDebugContractEtherReceived)
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
func (it *NetworkDebugContractEtherReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractEtherReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractEtherReceived represents a EtherReceived event raised by the NetworkDebugContract contract.
type NetworkDebugContractEtherReceived struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherReceived is a free log retrieval operation binding the contract event 0x1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b.
//
// Solidity: event EtherReceived(address sender, uint256 amount)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterEtherReceived(opts *bind.FilterOpts) (*NetworkDebugContractEtherReceivedIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "EtherReceived")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractEtherReceivedIterator{contract: _NetworkDebugContract.contract, event: "EtherReceived", logs: logs, sub: sub}, nil
}

// WatchEtherReceived is a free log subscription operation binding the contract event 0x1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b.
//
// Solidity: event EtherReceived(address sender, uint256 amount)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchEtherReceived(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractEtherReceived) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "EtherReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractEtherReceived)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "EtherReceived", log); err != nil {
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

// ParseEtherReceived is a log parse operation binding the contract event 0x1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b.
//
// Solidity: event EtherReceived(address sender, uint256 amount)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseEtherReceived(log types.Log) (*NetworkDebugContractEtherReceived, error) {
	event := new(NetworkDebugContractEtherReceived)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "EtherReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractIsValidEventIterator is returned from FilterIsValidEvent and is used to iterate over the raw logs and unpacked data for IsValidEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractIsValidEventIterator struct {
	Event *NetworkDebugContractIsValidEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractIsValidEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractIsValidEvent)
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
		it.Event = new(NetworkDebugContractIsValidEvent)
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
func (it *NetworkDebugContractIsValidEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractIsValidEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractIsValidEvent represents a IsValidEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractIsValidEvent struct {
	Success bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterIsValidEvent is a free log retrieval operation binding the contract event 0xdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b.
//
// Solidity: event IsValidEvent(bool success)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterIsValidEvent(opts *bind.FilterOpts) (*NetworkDebugContractIsValidEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "IsValidEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractIsValidEventIterator{contract: _NetworkDebugContract.contract, event: "IsValidEvent", logs: logs, sub: sub}, nil
}

// WatchIsValidEvent is a free log subscription operation binding the contract event 0xdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b.
//
// Solidity: event IsValidEvent(bool success)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchIsValidEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractIsValidEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "IsValidEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractIsValidEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "IsValidEvent", log); err != nil {
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

// ParseIsValidEvent is a log parse operation binding the contract event 0xdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b.
//
// Solidity: event IsValidEvent(bool success)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseIsValidEvent(log types.Log) (*NetworkDebugContractIsValidEvent, error) {
	event := new(NetworkDebugContractIsValidEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "IsValidEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractNoIndexEventIterator is returned from FilterNoIndexEvent and is used to iterate over the raw logs and unpacked data for NoIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEventIterator struct {
	Event *NetworkDebugContractNoIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractNoIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractNoIndexEvent)
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
		it.Event = new(NetworkDebugContractNoIndexEvent)
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
func (it *NetworkDebugContractNoIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractNoIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractNoIndexEvent represents a NoIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEvent struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNoIndexEvent is a free log retrieval operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterNoIndexEvent(opts *bind.FilterOpts) (*NetworkDebugContractNoIndexEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "NoIndexEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractNoIndexEventIterator{contract: _NetworkDebugContract.contract, event: "NoIndexEvent", logs: logs, sub: sub}, nil
}

// WatchNoIndexEvent is a free log subscription operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchNoIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractNoIndexEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "NoIndexEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractNoIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEvent", log); err != nil {
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

// ParseNoIndexEvent is a log parse operation binding the contract event 0x33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c.
//
// Solidity: event NoIndexEvent(address sender)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseNoIndexEvent(log types.Log) (*NetworkDebugContractNoIndexEvent, error) {
	event := new(NetworkDebugContractNoIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractNoIndexEventStringIterator is returned from FilterNoIndexEventString and is used to iterate over the raw logs and unpacked data for NoIndexEventString events raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEventStringIterator struct {
	Event *NetworkDebugContractNoIndexEventString // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractNoIndexEventStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractNoIndexEventString)
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
		it.Event = new(NetworkDebugContractNoIndexEventString)
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
func (it *NetworkDebugContractNoIndexEventStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractNoIndexEventStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractNoIndexEventString represents a NoIndexEventString event raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexEventString struct {
	Str string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoIndexEventString is a free log retrieval operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterNoIndexEventString(opts *bind.FilterOpts) (*NetworkDebugContractNoIndexEventStringIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "NoIndexEventString")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractNoIndexEventStringIterator{contract: _NetworkDebugContract.contract, event: "NoIndexEventString", logs: logs, sub: sub}, nil
}

// WatchNoIndexEventString is a free log subscription operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchNoIndexEventString(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractNoIndexEventString) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "NoIndexEventString")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractNoIndexEventString)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEventString", log); err != nil {
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

// ParseNoIndexEventString is a log parse operation binding the contract event 0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a.
//
// Solidity: event NoIndexEventString(string str)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseNoIndexEventString(log types.Log) (*NetworkDebugContractNoIndexEventString, error) {
	event := new(NetworkDebugContractNoIndexEventString)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexEventString", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractNoIndexStructEventIterator is returned from FilterNoIndexStructEvent and is used to iterate over the raw logs and unpacked data for NoIndexStructEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexStructEventIterator struct {
	Event *NetworkDebugContractNoIndexStructEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractNoIndexStructEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractNoIndexStructEvent)
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
		it.Event = new(NetworkDebugContractNoIndexStructEvent)
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
func (it *NetworkDebugContractNoIndexStructEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractNoIndexStructEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractNoIndexStructEvent represents a NoIndexStructEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractNoIndexStructEvent struct {
	A   NetworkDebugContractAccount
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoIndexStructEvent is a free log retrieval operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterNoIndexStructEvent(opts *bind.FilterOpts) (*NetworkDebugContractNoIndexStructEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "NoIndexStructEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractNoIndexStructEventIterator{contract: _NetworkDebugContract.contract, event: "NoIndexStructEvent", logs: logs, sub: sub}, nil
}

// WatchNoIndexStructEvent is a free log subscription operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchNoIndexStructEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractNoIndexStructEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "NoIndexStructEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractNoIndexStructEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexStructEvent", log); err != nil {
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

// ParseNoIndexStructEvent is a log parse operation binding the contract event 0xebe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc044188.
//
// Solidity: event NoIndexStructEvent((string,uint64,uint256) a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseNoIndexStructEvent(log types.Log) (*NetworkDebugContractNoIndexStructEvent, error) {
	event := new(NetworkDebugContractNoIndexStructEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "NoIndexStructEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractOneIndexEventIterator is returned from FilterOneIndexEvent and is used to iterate over the raw logs and unpacked data for OneIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractOneIndexEventIterator struct {
	Event *NetworkDebugContractOneIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractOneIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractOneIndexEvent)
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
		it.Event = new(NetworkDebugContractOneIndexEvent)
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
func (it *NetworkDebugContractOneIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractOneIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractOneIndexEvent represents a OneIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractOneIndexEvent struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterOneIndexEvent is a free log retrieval operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterOneIndexEvent(opts *bind.FilterOpts, a []*big.Int) (*NetworkDebugContractOneIndexEventIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "OneIndexEvent", aRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractOneIndexEventIterator{contract: _NetworkDebugContract.contract, event: "OneIndexEvent", logs: logs, sub: sub}, nil
}

// WatchOneIndexEvent is a free log subscription operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchOneIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractOneIndexEvent, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "OneIndexEvent", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractOneIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "OneIndexEvent", log); err != nil {
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

// ParseOneIndexEvent is a log parse operation binding the contract event 0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35.
//
// Solidity: event OneIndexEvent(uint256 indexed a)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseOneIndexEvent(log types.Log) (*NetworkDebugContractOneIndexEvent, error) {
	event := new(NetworkDebugContractOneIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "OneIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractReceivedIterator is returned from FilterReceived and is used to iterate over the raw logs and unpacked data for Received events raised by the NetworkDebugContract contract.
type NetworkDebugContractReceivedIterator struct {
	Event *NetworkDebugContractReceived // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractReceived)
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
		it.Event = new(NetworkDebugContractReceived)
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
func (it *NetworkDebugContractReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractReceived represents a Received event raised by the NetworkDebugContract contract.
type NetworkDebugContractReceived struct {
	Caller  common.Address
	Amount  *big.Int
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReceived is a free log retrieval operation binding the contract event 0x59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832.
//
// Solidity: event Received(address caller, uint256 amount, string message)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterReceived(opts *bind.FilterOpts) (*NetworkDebugContractReceivedIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractReceivedIterator{contract: _NetworkDebugContract.contract, event: "Received", logs: logs, sub: sub}, nil
}

// WatchReceived is a free log subscription operation binding the contract event 0x59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832.
//
// Solidity: event Received(address caller, uint256 amount, string message)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchReceived(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractReceived) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractReceived)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "Received", log); err != nil {
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

// ParseReceived is a log parse operation binding the contract event 0x59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b5832.
//
// Solidity: event Received(address caller, uint256 amount, string message)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseReceived(log types.Log) (*NetworkDebugContractReceived, error) {
	event := new(NetworkDebugContractReceived)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "Received", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator is returned from FilterThreeIndexAndOneNonIndexedEvent and is used to iterate over the raw logs and unpacked data for ThreeIndexAndOneNonIndexedEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator struct {
	Event *NetworkDebugContractThreeIndexAndOneNonIndexedEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
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
		it.Event = new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
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
func (it *NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractThreeIndexAndOneNonIndexedEvent represents a ThreeIndexAndOneNonIndexedEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexAndOneNonIndexedEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	DataId    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterThreeIndexAndOneNonIndexedEvent is a free log retrieval operation binding the contract event 0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d.
//
// Solidity: event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterThreeIndexAndOneNonIndexedEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (*NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "ThreeIndexAndOneNonIndexedEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractThreeIndexAndOneNonIndexedEventIterator{contract: _NetworkDebugContract.contract, event: "ThreeIndexAndOneNonIndexedEvent", logs: logs, sub: sub}, nil
}

// WatchThreeIndexAndOneNonIndexedEvent is a free log subscription operation binding the contract event 0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d.
//
// Solidity: event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchThreeIndexAndOneNonIndexedEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractThreeIndexAndOneNonIndexedEvent, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "ThreeIndexAndOneNonIndexedEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexAndOneNonIndexedEvent", log); err != nil {
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

// ParseThreeIndexAndOneNonIndexedEvent is a log parse operation binding the contract event 0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d.
//
// Solidity: event ThreeIndexAndOneNonIndexedEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt, string dataId)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseThreeIndexAndOneNonIndexedEvent(log types.Log) (*NetworkDebugContractThreeIndexAndOneNonIndexedEvent, error) {
	event := new(NetworkDebugContractThreeIndexAndOneNonIndexedEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexAndOneNonIndexedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractThreeIndexEventIterator is returned from FilterThreeIndexEvent and is used to iterate over the raw logs and unpacked data for ThreeIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexEventIterator struct {
	Event *NetworkDebugContractThreeIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractThreeIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractThreeIndexEvent)
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
		it.Event = new(NetworkDebugContractThreeIndexEvent)
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
func (it *NetworkDebugContractThreeIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractThreeIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractThreeIndexEvent represents a ThreeIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractThreeIndexEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterThreeIndexEvent is a free log retrieval operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterThreeIndexEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (*NetworkDebugContractThreeIndexEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "ThreeIndexEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractThreeIndexEventIterator{contract: _NetworkDebugContract.contract, event: "ThreeIndexEvent", logs: logs, sub: sub}, nil
}

// WatchThreeIndexEvent is a free log subscription operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchThreeIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractThreeIndexEvent, roundId []*big.Int, startedBy []common.Address, startedAt []*big.Int) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}
	var startedAtRule []interface{}
	for _, startedAtItem := range startedAt {
		startedAtRule = append(startedAtRule, startedAtItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "ThreeIndexEvent", roundIdRule, startedByRule, startedAtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractThreeIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexEvent", log); err != nil {
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

// ParseThreeIndexEvent is a log parse operation binding the contract event 0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3.
//
// Solidity: event ThreeIndexEvent(uint256 indexed roundId, address indexed startedBy, uint256 indexed startedAt)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseThreeIndexEvent(log types.Log) (*NetworkDebugContractThreeIndexEvent, error) {
	event := new(NetworkDebugContractThreeIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "ThreeIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkDebugContractTwoIndexEventIterator is returned from FilterTwoIndexEvent and is used to iterate over the raw logs and unpacked data for TwoIndexEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractTwoIndexEventIterator struct {
	Event *NetworkDebugContractTwoIndexEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractTwoIndexEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractTwoIndexEvent)
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
		it.Event = new(NetworkDebugContractTwoIndexEvent)
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
func (it *NetworkDebugContractTwoIndexEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractTwoIndexEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractTwoIndexEvent represents a TwoIndexEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractTwoIndexEvent struct {
	RoundId   *big.Int
	StartedBy common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTwoIndexEvent is a free log retrieval operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterTwoIndexEvent(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*NetworkDebugContractTwoIndexEventIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "TwoIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractTwoIndexEventIterator{contract: _NetworkDebugContract.contract, event: "TwoIndexEvent", logs: logs, sub: sub}, nil
}

// WatchTwoIndexEvent is a free log subscription operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchTwoIndexEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractTwoIndexEvent, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "TwoIndexEvent", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractTwoIndexEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "TwoIndexEvent", log); err != nil {
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

// ParseTwoIndexEvent is a log parse operation binding the contract event 0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5.
//
// Solidity: event TwoIndexEvent(uint256 indexed roundId, address indexed startedBy)
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseTwoIndexEvent(log types.Log) (*NetworkDebugContractTwoIndexEvent, error) {
	event := new(NetworkDebugContractTwoIndexEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "TwoIndexEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
