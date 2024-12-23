// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package NetworkDebugContract

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"required\",\"type\":\"uint256\"}],\"name\":\"CustomErr\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CustomErrNoValues\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"CustomErrWithMessage\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"CallDataLength\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"a\",\"type\":\"int256\"}],\"name\":\"CallbackEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"CurrentStatus\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EtherReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"IsValidEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NoIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"NoIndexEventString\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"balance\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"dailyLimit\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structNetworkDebugContract.Account\",\"name\":\"a\",\"type\":\"tuple\"}],\"name\":\"NoIndexStructEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"OneIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"dataId\",\"type\":\"string\"}],\"name\":\"ThreeIndexAndOneNonIndexedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"ThreeIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"}],\"name\":\"TwoIndexEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"UniqueDebugEvent\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"idx\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"addCounter\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsAssert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsCustomError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsCustomErrorNoValues\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alwaysRevertsRequire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"callRevertFunctionInSubContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"callRevertFunctionInTheContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"callbackMethod\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"counterMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStatus\",\"outputs\":[{\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"emitAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"input\",\"type\":\"bytes32\"}],\"name\":\"emitBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"output\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitFourParamMixedEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"inputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inputVal2\",\"type\":\"string\"}],\"name\":\"emitInputs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"inputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inputVal2\",\"type\":\"string\"}],\"name\":\"emitInputsOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"first\",\"type\":\"int256\"},{\"internalType\":\"int128\",\"name\":\"second\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"third\",\"type\":\"uint256\"}],\"name\":\"emitInts\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"int128\",\"name\":\"outputVal1\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"outputVal2\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"inputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inputVal2\",\"type\":\"string\"}],\"name\":\"emitNamedInputsOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"outputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"outputVal2\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNamedOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"outputVal1\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"outputVal2\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNoIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNoIndexEventString\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitNoIndexStructEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitOneIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitOutputs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitThreeIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emitTwoIndexEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"data\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"idx\",\"type\":\"int256\"}],\"name\":\"getCounter\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"data\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"data\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pay\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performStaticCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"input\",\"type\":\"address[]\"}],\"name\":\"processAddressArray\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"processDynamicData\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data[3]\",\"name\":\"data\",\"type\":\"tuple[3]\"}],\"name\":\"processFixedDataArray\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data[2]\",\"name\":\"\",\"type\":\"tuple[2]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"dynamicBytes\",\"type\":\"bytes\"}],\"internalType\":\"structNetworkDebugContract.NestedData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"processNestedData\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"dynamicBytes\",\"type\":\"bytes\"}],\"internalType\":\"structNetworkDebugContract.NestedData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"processNestedData\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"internalType\":\"structNetworkDebugContract.Data\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"dynamicBytes\",\"type\":\"bytes\"}],\"internalType\":\"structNetworkDebugContract.NestedData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"input\",\"type\":\"uint256[]\"}],\"name\":\"processUintArray\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"idx\",\"type\":\"int256\"}],\"name\":\"resetCounter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"set\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"setMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"setStatus\",\"outputs\":[{\"internalType\":\"enumNetworkDebugContract.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storedData\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"storedDataMap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subContract\",\"outputs\":[{\"internalType\":\"contractNetworkDebugSubContract\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"trace\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceDifferent\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"traceNestedEvents\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceSubWithCallback\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceWithValidate\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"traceYetDifferent\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620040f7380380620040f78339818101604052810190620000379190620000f2565b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101006004819055505062000124565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000ba826200008d565b9050919050565b620000cc81620000ad565b8114620000d857600080fd5b50565b600081519050620000ec81620000c1565b92915050565b6000602082840312156200010b576200010a62000088565b5b60006200011b84828501620000db565b91505092915050565b613fc380620001346000396000f3fe6080604052600436106103035760003560e01c80637f12881c11610190578063b3f8d1b2116100dc578063e5c19b2d11610095578063ef8a92351161006f578063ef8a923514610bfa578063f3396bd914610c25578063f499af2a14610c4e578063fbcb8d0714610c8b57610343565b8063e5c19b2d14610b43578063e8116e2814610b80578063ec5c3ede14610bbd57610343565b8063b3f8d1b214610a58578063b600141f14610a6f578063c0d06d8914610a86578063c2124b2214610ab1578063d7a8020514610ac8578063e1111f7914610b0657610343565b806395a81a4c11610149578063a4c0ed3611610123578063a4c0ed36146109bd578063aa3fdcf4146109e6578063ad3de14c146109fd578063b1ae9d8514610a2857610343565b806395a81a4c1461092a57806399adad2e146109415780639e0996521461097e57610343565b80637f12881c1461082d5780637fdc8fe11461086a57806381b375a0146108a75780638db611be146108d05780638f856296146108fc5780639349d00b1461091357610343565b80633837a75e1161024f5780635921483f1161020857806362c270e1116101e257806362c270e1146107a85780636d4ce63c146107bf5780637014c81d146107ea578063788c47721461081657610343565b80635921483f146107175780635e9c80d6146107545780636284117d1461076b57610343565b80633837a75e146105ba5780633bc5de30146105f75780633e41f1351461062257806345f0c9e61461065f57806348ad9fe81461069d57806358379d71146106da57610343565b806323515760116102bc5780632e49d78b116102965780632e49d78b146104d857806330985bcc146105155780633170428e1461055257806333311ef31461057d57610343565b80632351576014610459578063256560d5146104965780632a1afcd9146104ad57610343565b806304d8215b1461037e57806306595f75146103bb57806311b3c478146103d257806312d91233146103fb5780631b9265b8146104385780631e31d0a81461044257610343565b36610343577f59e04c3f0d44b7caf6e8ef854b61d9a51cf1960d7a88ff6356cc5e946b4b583233346040516103399291906121a5565b60405180910390a1005b7f1e57e3bb474320be3d2c77138f75b7c3941292d647f5f9634e33a8e94e0e069b33346040516103749291906121e1565b60405180910390a1005b34801561038a57600080fd5b506103a560048036038101906103a09190612254565b610cc8565b6040516103b291906122af565b60405180910390f35b3480156103c757600080fd5b506103d0610d0e565b005b3480156103de57600080fd5b506103f960048036038101906103f491906122f6565b610d51565b005b34801561040757600080fd5b50610422600480360381019061041d919061248f565b610de4565b60405161042f9190612596565b60405180910390f35b610440610ea3565b005b34801561044e57600080fd5b50610457610ea5565b005b34801561046557600080fd5b50610480600480360381019061047b9190612254565b610eec565b60405161048d91906125c7565b60405180910390f35b3480156104a257600080fd5b506104ab610f21565b005b3480156104b957600080fd5b506104c2610f32565b6040516104cf91906125c7565b60405180910390f35b3480156104e457600080fd5b506104ff60048036038101906104fa9190612607565b610f38565b60405161050c91906126ab565b60405180910390f35b34801561052157600080fd5b5061053c60048036038101906105379190612254565b610fba565b60405161054991906125c7565b60405180910390f35b34801561055e57600080fd5b5061056761109c565b60405161057491906126c6565b60405180910390f35b34801561058957600080fd5b506105a4600480360381019061059f9190612717565b6111db565b6040516105b19190612753565b60405180910390f35b3480156105c657600080fd5b506105e160048036038101906105dc9190612254565b6111e5565b6040516105ee91906125c7565b60405180910390f35b34801561060357600080fd5b5061060c6112f0565b60405161061991906126c6565b60405180910390f35b34801561062e57600080fd5b5061064960048036038101906106449190612254565b6112fa565b60405161065691906125c7565b60405180910390f35b34801561066b57600080fd5b5061068660048036038101906106819190612823565b6113f5565b6040516106949291906128ed565b60405180910390f35b3480156106a957600080fd5b506106c460048036038101906106bf9190612949565b611406565b6040516106d191906125c7565b60405180910390f35b3480156106e657600080fd5b5061070160048036038101906106fc9190612254565b61141e565b60405161070e91906125c7565b60405180910390f35b34801561072357600080fd5b5061073e60048036038101906107399190612976565b611519565b60405161074b91906125c7565b60405180910390f35b34801561076057600080fd5b50610769611536565b005b34801561077757600080fd5b50610792600480360381019061078d9190612976565b611577565b60405161079f91906125c7565b60405180910390f35b3480156107b457600080fd5b506107bd61158f565b005b3480156107cb57600080fd5b506107d4611623565b6040516107e191906125c7565b60405180910390f35b3480156107f657600080fd5b506107ff61162c565b60405161080d9291906128ed565b60405180910390f35b34801561082257600080fd5b5061082b611671565b005b34801561083957600080fd5b50610854600480360381019061084f91906129c7565b6116a8565b6040516108619190612ba6565b60405180910390f35b34801561087657600080fd5b50610891600480360381019061088c9190612be7565b6116c1565b60405161089e9190612df4565b60405180910390f35b3480156108b357600080fd5b506108ce60048036038101906108c99190612823565b6116ca565b005b3480156108dc57600080fd5b506108e56116ce565b6040516108f39291906128ed565b60405180910390f35b34801561090857600080fd5b50610911611713565b005b34801561091f57600080fd5b50610928611743565b005b34801561093657600080fd5b5061093f61174d565b005b34801561094d57600080fd5b5061096860048036038101906109639190612e38565b611786565b6040516109759190612f37565b60405180910390f35b34801561098a57600080fd5b506109a560048036038101906109a09190612f92565b611831565b6040516109b493929190612ff4565b60405180910390f35b3480156109c957600080fd5b506109e460048036038101906109df9190613086565b611848565b005b3480156109f257600080fd5b506109fb611b34565b005b348015610a0957600080fd5b50610a12611b7d565b604051610a1f91906125c7565b60405180910390f35b610a426004803603810190610a3d9190612254565b611bc4565b604051610a4f91906125c7565b60405180910390f35b348015610a6457600080fd5b50610a6d611d0e565b005b348015610a7b57600080fd5b50610a84611dbe565b005b348015610a9257600080fd5b50610a9b611df0565b604051610aa89190613159565b60405180910390f35b348015610abd57600080fd5b50610ac6611e16565b005b348015610ad457600080fd5b50610aef6004803603810190610aea9190612823565b611e68565b604051610afd9291906128ed565b60405180910390f35b348015610b1257600080fd5b50610b2d6004803603810190610b289190613237565b611e79565b604051610b3a919061333e565b60405180910390f35b348015610b4f57600080fd5b50610b6a6004803603810190610b659190612976565b611e83565b604051610b7791906125c7565b60405180910390f35b348015610b8c57600080fd5b50610ba76004803603810190610ba29190612976565b611e94565b604051610bb491906125c7565b60405180910390f35b348015610bc957600080fd5b50610be46004803603810190610bdf9190612949565b611ee2565b604051610bf19190613360565b60405180910390f35b348015610c0657600080fd5b50610c0f611eec565b604051610c1c91906126ab565b60405180910390f35b348015610c3157600080fd5b50610c4c6004803603810190610c479190612976565b611eff565b005b348015610c5a57600080fd5b50610c756004803603810190610c709190612be7565b611f1b565b604051610c829190612ba6565b60405180910390f35b348015610c9757600080fd5b50610cb26004803603810190610cad9190612976565b612050565b604051610cbf91906125c7565b60405180910390f35b60007fdfac7500004753b91139af55816e7eade36d96faec68b343f77ed66b89912a7b828413604051610cfb91906122af565b60405180910390a1818313905092915050565b6000610d4f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d46906133c7565b60405180910390fd5b565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166311abb00283836040518363ffffffff1660e01b8152600401610dae9291906133e7565b600060405180830381600087803b158015610dc857600080fd5b505af1158015610ddc573d6000803e3d6000fd5b505050505050565b60606000825167ffffffffffffffff811115610e0357610e0261234c565b5b604051908082528060200260200182016040528015610e315781602001602082028036833780820191505090505b50905060005b8351811015610e99576001848281518110610e5557610e54613410565b5b6020026020010151610e67919061346e565b828281518110610e7a57610e79613410565b5b6020026020010181815250508080610e91906134a2565b915050610e37565b5080915050919050565b565b3373ffffffffffffffffffffffffffffffffffffffff1660017f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a3565b600081600260008581526020019081526020016000206000828254610f1191906134ea565b9250508190555081905092915050565b6000610f3057610f2f61352e565b5b565b60005481565b600081600560006101000a81548160ff02191690836003811115610f5f57610f5e612634565b5b0217905550600560009054906101000a900460ff166003811115610f8657610f85612634565b5b7fbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a60405160405180910390a2819050919050565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663047c4425836040518263ffffffff1660e01b815260040161101791906125c7565b6020604051808303816000875af1158015611036573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061105a9190613572565b50827feace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b3560405160405180910390a2818361109491906134ea565b905092915050565b6000803090506000808273ffffffffffffffffffffffffffffffffffffffff16633bc5de3060e01b604051602401604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161113491906135db565b600060405180830381855afa9150503d806000811461116f576040519150601f19603f3d011682016040523d82523d6000602084013e611174565b606091505b5091509150816111b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111b09061363e565b60405180910390fd5b6000818060200190518101906111cf9190613673565b90508094505050505090565b6000819050919050565b60006002826111f491906134ea565b9150600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fa8fca7a84846040518363ffffffff1660e01b81526004016112539291906136a0565b6020604051808303816000875af1158015611272573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112969190613572565b503373ffffffffffffffffffffffffffffffffffffffff1660017f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a381836112e891906134ea565b905092915050565b6000600454905090565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633e41f13584846040518363ffffffff1660e01b81526004016113599291906136a0565b6020604051808303816000875af1158015611378573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061139c9190613572565b503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a381836113ed91906134ea565b905092915050565b600060608383915091509250929050565b60016020528060005260406000206000915090505481565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633e41f13584846040518363ffffffff1660e01b815260040161147d9291906136a0565b6020604051808303816000875af115801561149c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114c09190613572565b503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a3818361151191906134ea565b905092915050565b600060026000838152602001908152602001600020549050919050565b600c60156040517f4a2eaf7e00000000000000000000000000000000000000000000000000000000815260040161156e92919061373f565b60405180910390fd5b60026020528060005260406000206000915090505481565b7febe3ff7e2071d351bf2e65b4fccd24e3ae99485f02468f1feecf7d64dc04418860405180606001604052806040518060400160405280600481526020017f4a6f686e000000000000000000000000000000000000000000000000000000008152508152602001600567ffffffffffffffff168152602001600a81525060405161161991906137db565b60405180910390a1565b60008054905090565b60006060617a696040518060400160405280600a81526020017f6f757470757456616c3100000000000000000000000000000000000000000000815250915091509091565b7f25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a60405161169e90613849565b60405180910390a1565b6116b0612087565b816116ba90613a24565b9050919050565b36819050919050565b5050565b60006060617a696040518060400160405280600a81526020017f6f757470757456616c3100000000000000000000000000000000000000000000815250915091509091565b60537feace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b3560405160405180910390a2565b61174b611536565b565b7f33bc9bae48dbe1e057f264b3fc6a1dacdcceacb3ba28d937231c70e068a02f1c3360405161177c9190613360565b60405180910390a1565b61178e6120a7565b6117966120a7565b826000600381106117aa576117a9613410565b5b6020028101906117ba9190613a46565b6117c390613a6e565b816000600281106117d7576117d6613410565b5b6020020181905250826001600381106117f3576117f2613410565b5b6020028101906118039190613a46565b61180c90613a6e565b816001600281106118205761181f613410565b5b602002018190525080915050919050565b600080600085858592509250925093509350939050565b7f962c5df4c8ad201a4f54a88f47715bb2cf291d019e350e2dff50ca6fc0f5d0ed8282905060405161187a91906126c6565b60405180910390a1600082829050036118ce57606360656040517f4a2eaf7e0000000000000000000000000000000000000000000000000000000081526004016118c5929190613af7565b60405180910390fd5b6000803073ffffffffffffffffffffffffffffffffffffffff1684846040516118f8929190613b45565b600060405180830381855af49150503d8060008114611933576040519150601f19603f3d011682016040523d82523d6000602084013e611938565b606091505b509150915081611990576000815111156119555780518082602001fd5b6040517f2350eb5200000000000000000000000000000000000000000000000000000000815260040161198790613bd0565b60405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff16633170428e6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156119db573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119ff9190613673565b5060008484600090600492611a1693929190613bfa565b90611a219190613c79565b90506358379d7160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611aaa576040517f2350eb52000000000000000000000000000000000000000000000000000000008152600401611aa190613d24565b60405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff16633837a75e600160026040518363ffffffff1660e01b8152600401611ae7929190613dba565b6020604051808303816000875af1158015611b06573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b2a9190613572565b5050505050505050565b60033373ffffffffffffffffffffffffffffffffffffffff1660017f5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf360405160405180910390a4565b6000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905090565b6000611bd08383610cc8565b15611ccd57600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633e41f13584846040518363ffffffff1660e01b8152600401611c329291906136a0565b6020604051808303816000875af1158015611c51573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c759190613572565b503373ffffffffffffffffffffffffffffffffffffffff16827f33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b560405160405180910390a38183611cc691906134ea565b9050611d08565b6040517f2350eb52000000000000000000000000000000000000000000000000000000008152600401611cff90613e55565b60405180910390fd5b92915050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632d948d896040518163ffffffff1660e01b8152600401600060405180830381600087803b158015611d7857600080fd5b505af1158015611d8c573d6000803e3d6000fd5b505050507fa0f7c7c1fff15178b5db3e56860767f0889c56b591bd2d9ba3121b491347d74c60405160405180910390a1565b6040517fa0c2d2db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60033373ffffffffffffffffffffffffffffffffffffffff1660027f56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d604051611e5e90613ec1565b60405180910390a4565b600060608383915091509250929050565b6060819050919050565b600081600081905550819050919050565b600081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550819050919050565b6000819050919050565b600560009054906101000a900460ff1681565b6000600260008381526020019081526020016000208190555050565b611f23612087565b6000828060000190611f359190613ee1565b604051602001611f46929190613f74565b6040516020818303038152906040528051906020012090506000602067ffffffffffffffff811115611f7b57611f7a61234c565b5b6040519080825280601f01601f191660200182016040528015611fad5781602001600182028036833780820191505090505b50905060005b602081101561202857828160208110611fcf57611fce613410565b5b1a60f81b828281518110611fe657611fe5613410565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508080612020906134a2565b915050611fb3565b5060405180604001604052808561203e90613a6e565b81526020018281525092505050919050565b6000817fb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c6860405160405180910390a2819050919050565b604051806040016040528061209a6120d4565b8152602001606081525090565b60405180604001604052806002905b6120be6120d4565b8152602001906001900390816120b65790505090565b604051806040016040528060608152602001606081525090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000612119826120ee565b9050919050565b6121298161210e565b82525050565b6000819050919050565b6121428161212f565b82525050565b600082825260208201905092915050565b7f5265636569766564204574686572000000000000000000000000000000000000600082015250565b600061218f600e83612148565b915061219a82612159565b602082019050919050565b60006060820190506121ba6000830185612120565b6121c76020830184612139565b81810360408301526121d881612182565b90509392505050565b60006040820190506121f66000830185612120565b6122036020830184612139565b9392505050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b6122318161221e565b811461223c57600080fd5b50565b60008135905061224e81612228565b92915050565b6000806040838503121561226b5761226a612214565b5b60006122798582860161223f565b925050602061228a8582860161223f565b9150509250929050565b60008115159050919050565b6122a981612294565b82525050565b60006020820190506122c460008301846122a0565b92915050565b6122d38161212f565b81146122de57600080fd5b50565b6000813590506122f0816122ca565b92915050565b6000806040838503121561230d5761230c612214565b5b600061231b858286016122e1565b925050602061232c858286016122e1565b9150509250929050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6123848261233b565b810181811067ffffffffffffffff821117156123a3576123a261234c565b5b80604052505050565b60006123b661220a565b90506123c2828261237b565b919050565b600067ffffffffffffffff8211156123e2576123e161234c565b5b602082029050602081019050919050565b600080fd5b600061240b612406846123c7565b6123ac565b9050808382526020820190506020840283018581111561242e5761242d6123f3565b5b835b81811015612457578061244388826122e1565b845260208401935050602081019050612430565b5050509392505050565b600082601f83011261247657612475612336565b5b81356124868482602086016123f8565b91505092915050565b6000602082840312156124a5576124a4612214565b5b600082013567ffffffffffffffff8111156124c3576124c2612219565b5b6124cf84828501612461565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61250d8161212f565b82525050565b600061251f8383612504565b60208301905092915050565b6000602082019050919050565b6000612543826124d8565b61254d81856124e3565b9350612558836124f4565b8060005b838110156125895781516125708882612513565b975061257b8361252b565b92505060018101905061255c565b5085935050505092915050565b600060208201905081810360008301526125b08184612538565b905092915050565b6125c18161221e565b82525050565b60006020820190506125dc60008301846125b8565b92915050565b600481106125ef57600080fd5b50565b600081359050612601816125e2565b92915050565b60006020828403121561261d5761261c612214565b5b600061262b848285016125f2565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6004811061267457612673612634565b5b50565b600081905061268582612663565b919050565b600061269582612677565b9050919050565b6126a58161268a565b82525050565b60006020820190506126c0600083018461269c565b92915050565b60006020820190506126db6000830184612139565b92915050565b6000819050919050565b6126f4816126e1565b81146126ff57600080fd5b50565b600081359050612711816126eb565b92915050565b60006020828403121561272d5761272c612214565b5b600061273b84828501612702565b91505092915050565b61274d816126e1565b82525050565b60006020820190506127686000830184612744565b92915050565b600080fd5b600067ffffffffffffffff82111561278e5761278d61234c565b5b6127978261233b565b9050602081019050919050565b82818337600083830152505050565b60006127c66127c184612773565b6123ac565b9050828152602081018484840111156127e2576127e161276e565b5b6127ed8482856127a4565b509392505050565b600082601f83011261280a57612809612336565b5b813561281a8482602086016127b3565b91505092915050565b6000806040838503121561283a57612839612214565b5b6000612848858286016122e1565b925050602083013567ffffffffffffffff81111561286957612868612219565b5b612875858286016127f5565b9150509250929050565b600081519050919050565b60005b838110156128a857808201518184015260208101905061288d565b60008484015250505050565b60006128bf8261287f565b6128c98185612148565b93506128d981856020860161288a565b6128e28161233b565b840191505092915050565b60006040820190506129026000830185612139565b818103602083015261291481846128b4565b90509392505050565b6129268161210e565b811461293157600080fd5b50565b6000813590506129438161291d565b92915050565b60006020828403121561295f5761295e612214565b5b600061296d84828501612934565b91505092915050565b60006020828403121561298c5761298b612214565b5b600061299a8482850161223f565b91505092915050565b600080fd5b6000604082840312156129be576129bd6129a3565b5b81905092915050565b6000602082840312156129dd576129dc612214565b5b600082013567ffffffffffffffff8111156129fb576129fa612219565b5b612a07848285016129a8565b91505092915050565b600082825260208201905092915050565b6000612a2c8261287f565b612a368185612a10565b9350612a4681856020860161288a565b612a4f8161233b565b840191505092915050565b600082825260208201905092915050565b6000612a76826124d8565b612a808185612a5a565b9350612a8b836124f4565b8060005b83811015612abc578151612aa38882612513565b9750612aae8361252b565b925050600181019050612a8f565b5085935050505092915050565b60006040830160008301518482036000860152612ae68282612a21565b91505060208301518482036020860152612b008282612a6b565b9150508091505092915050565b600081519050919050565b600082825260208201905092915050565b6000612b3482612b0d565b612b3e8185612b18565b9350612b4e81856020860161288a565b612b578161233b565b840191505092915050565b60006040830160008301518482036000860152612b7f8282612ac9565b91505060208301518482036020860152612b998282612b29565b9150508091505092915050565b60006020820190508181036000830152612bc08184612b62565b905092915050565b600060408284031215612bde57612bdd6129a3565b5b81905092915050565b600060208284031215612bfd57612bfc612214565b5b600082013567ffffffffffffffff811115612c1b57612c1a612219565b5b612c2784828501612bc8565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112612c5c57612c5b612c3a565b5b83810192508235915060208301925067ffffffffffffffff821115612c8457612c83612c30565b5b600182023603831315612c9a57612c99612c35565b5b509250929050565b6000612cae8385612a10565b9350612cbb8385846127a4565b612cc48361233b565b840190509392505050565b60008083356001602003843603038112612cec57612ceb612c3a565b5b83810192508235915060208301925067ffffffffffffffff821115612d1457612d13612c30565b5b602082023603831315612d2a57612d29612c35565b5b509250929050565b600080fd5b82818337505050565b6000612d4c8385612a5a565b93507f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115612d7f57612d7e612d32565b5b602083029250612d90838584612d37565b82840190509392505050565b600060408301612daf6000840184612c3f565b8583036000870152612dc2838284612ca2565b92505050612dd36020840184612ccf565b8583036020870152612de6838284612d40565b925050508091505092915050565b60006020820190508181036000830152612e0e8184612d9c565b905092915050565b600081905082602060030282011115612e3257612e316123f3565b5b92915050565b600060208284031215612e4e57612e4d612214565b5b600082013567ffffffffffffffff811115612e6c57612e6b612219565b5b612e7884828501612e16565b91505092915050565b600060029050919050565b600081905092915050565b6000819050919050565b6000612ead8383612ac9565b905092915050565b6000602082019050919050565b6000612ecd82612e81565b612ed78185612e8c565b935083602082028501612ee985612e97565b8060005b85811015612f255784840389528151612f068582612ea1565b9450612f1183612eb5565b925060208a01995050600181019050612eed565b50829750879550505050505092915050565b60006020820190508181036000830152612f518184612ec2565b905092915050565b600081600f0b9050919050565b612f6f81612f59565b8114612f7a57600080fd5b50565b600081359050612f8c81612f66565b92915050565b600080600060608486031215612fab57612faa612214565b5b6000612fb98682870161223f565b9350506020612fca86828701612f7d565b9250506040612fdb868287016122e1565b9150509250925092565b612fee81612f59565b82525050565b600060608201905061300960008301866125b8565b6130166020830185612fe5565b6130236040830184612139565b949350505050565b600080fd5b60008083601f84011261304657613045612336565b5b8235905067ffffffffffffffff8111156130635761306261302b565b5b60208301915083600182028301111561307f5761307e6123f3565b5b9250929050565b600080600080606085870312156130a05761309f612214565b5b60006130ae87828801612934565b94505060206130bf878288016122e1565b935050604085013567ffffffffffffffff8111156130e0576130df612219565b5b6130ec87828801613030565b925092505092959194509250565b6000819050919050565b600061311f61311a613115846120ee565b6130fa565b6120ee565b9050919050565b600061313182613104565b9050919050565b600061314382613126565b9050919050565b61315381613138565b82525050565b600060208201905061316e600083018461314a565b92915050565b600067ffffffffffffffff82111561318f5761318e61234c565b5b602082029050602081019050919050565b60006131b36131ae84613174565b6123ac565b905080838252602082019050602084028301858111156131d6576131d56123f3565b5b835b818110156131ff57806131eb8882612934565b8452602084019350506020810190506131d8565b5050509392505050565b600082601f83011261321e5761321d612336565b5b813561322e8482602086016131a0565b91505092915050565b60006020828403121561324d5761324c612214565b5b600082013567ffffffffffffffff81111561326b5761326a612219565b5b61327784828501613209565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6132b58161210e565b82525050565b60006132c783836132ac565b60208301905092915050565b6000602082019050919050565b60006132eb82613280565b6132f5818561328b565b93506133008361329c565b8060005b8381101561333157815161331888826132bb565b9750613323836132d3565b925050600181019050613304565b5085935050505092915050565b6000602082019050818103600083015261335881846132e0565b905092915050565b60006020820190506133756000830184612120565b92915050565b7f616c7761797320726576657274206572726f7200000000000000000000000000600082015250565b60006133b1601383612148565b91506133bc8261337b565b602082019050919050565b600060208201905081810360008301526133e0816133a4565b9050919050565b60006040820190506133fc6000830185612139565b6134096020830184612139565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006134798261212f565b91506134848361212f565b925082820190508082111561349c5761349b61343f565b5b92915050565b60006134ad8261212f565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036134df576134de61343f565b5b600182019050919050565b60006134f58261221e565b91506135008361221e565b9250828201905082811215600083121683821260008412151617156135285761352761343f565b5b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b60008151905061356c81612228565b92915050565b60006020828403121561358857613587612214565b5b60006135968482850161355d565b91505092915050565b600081905092915050565b60006135b582612b0d565b6135bf818561359f565b93506135cf81856020860161288a565b80840191505092915050565b60006135e782846135aa565b915081905092915050565b7f5374617469632063616c6c206661696c65640000000000000000000000000000600082015250565b6000613628601283612148565b9150613633826135f2565b602082019050919050565b600060208201905081810360008301526136578161361b565b9050919050565b60008151905061366d816122ca565b92915050565b60006020828403121561368957613688612214565b5b60006136978482850161365e565b91505092915050565b60006040820190506136b560008301856125b8565b6136c260208301846125b8565b9392505050565b6000819050919050565b60006136ee6136e96136e4846136c9565b6130fa565b61212f565b9050919050565b6136fe816136d3565b82525050565b6000819050919050565b600061372961372461371f84613704565b6130fa565b61212f565b9050919050565b6137398161370e565b82525050565b600060408201905061375460008301856136f5565b6137616020830184613730565b9392505050565b600067ffffffffffffffff82169050919050565b61378581613768565b82525050565b600060608301600083015184820360008601526137a88282612a21565b91505060208301516137bd602086018261377c565b5060408301516137d06040860182612504565b508091505092915050565b600060208201905081810360008301526137f5818461378b565b905092915050565b7f6d79537472696e67000000000000000000000000000000000000000000000000600082015250565b6000613833600883612148565b915061383e826137fd565b602082019050919050565b6000602082019050818103600083015261386281613826565b9050919050565b600080fd5b600080fd5b60006040828403121561388957613888613869565b5b61389360406123ac565b9050600082013567ffffffffffffffff8111156138b3576138b261386e565b5b6138bf848285016127f5565b600083015250602082013567ffffffffffffffff8111156138e3576138e261386e565b5b6138ef84828501612461565b60208301525092915050565b600067ffffffffffffffff8211156139165761391561234c565b5b61391f8261233b565b9050602081019050919050565b600061393f61393a846138fb565b6123ac565b90508281526020810184848401111561395b5761395a61276e565b5b6139668482856127a4565b509392505050565b600082601f83011261398357613982612336565b5b813561399384826020860161392c565b91505092915050565b6000604082840312156139b2576139b1613869565b5b6139bc60406123ac565b9050600082013567ffffffffffffffff8111156139dc576139db61386e565b5b6139e884828501613873565b600083015250602082013567ffffffffffffffff811115613a0c57613a0b61386e565b5b613a188482850161396e565b60208301525092915050565b6000613a30368361399c565b9050919050565b600080fd5b600080fd5b600080fd5b600082356001604003833603038112613a6257613a61613a37565b5b80830191505092915050565b6000613a7a3683613873565b9050919050565b6000819050919050565b6000613aa6613aa1613a9c84613a81565b6130fa565b61212f565b9050919050565b613ab681613a8b565b82525050565b6000819050919050565b6000613ae1613adc613ad784613abc565b6130fa565b61212f565b9050919050565b613af181613ac6565b82525050565b6000604082019050613b0c6000830185613aad565b613b196020830184613ae8565b9392505050565b6000613b2c838561359f565b9350613b398385846127a4565b82840190509392505050565b6000613b52828486613b20565b91508190509392505050565b7f64656c656761746563616c6c206661696c65642077697468206e6f207265617360008201527f6f6e000000000000000000000000000000000000000000000000000000000000602082015250565b6000613bba602283612148565b9150613bc582613b5e565b604082019050919050565b60006020820190508181036000830152613be981613bad565b9050919050565b600080fd5b600080fd5b60008085851115613c0e57613c0d613bf0565b5b83861115613c1f57613c1e613bf5565b5b6001850283019150848603905094509492505050565b600082905092915050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600082821b905092915050565b6000613c858383613c35565b82613c908135613c40565b92506004821015613cd057613ccb7fffffffff0000000000000000000000000000000000000000000000000000000083600403600802613c6c565b831692505b505092915050565b7f6f68206f68206f682069742773206d6167696321000000000000000000000000600082015250565b6000613d0e601483612148565b9150613d1982613cd8565b602082019050919050565b60006020820190508181036000830152613d3d81613d01565b9050919050565b6000819050919050565b6000613d69613d64613d5f84613d44565b6130fa565b61221e565b9050919050565b613d7981613d4e565b82525050565b6000819050919050565b6000613da4613d9f613d9a84613d7f565b6130fa565b61221e565b9050919050565b613db481613d89565b82525050565b6000604082019050613dcf6000830185613d70565b613ddc6020830184613dab565b9392505050565b7f666972737420696e7420776173206e6f742067726561746572207468616e207360008201527f65636f6e6420696e740000000000000000000000000000000000000000000000602082015250565b6000613e3f602983612148565b9150613e4a82613de3565b604082019050919050565b60006020820190508181036000830152613e6e81613e32565b9050919050565b7f736f6d6520696400000000000000000000000000000000000000000000000000600082015250565b6000613eab600783612148565b9150613eb682613e75565b602082019050919050565b60006020820190508181036000830152613eda81613e9e565b9050919050565b60008083356001602003843603038112613efe57613efd613a37565b5b80840192508235915067ffffffffffffffff821115613f2057613f1f613a3c565b5b602083019250600182023603831315613f3c57613f3b613a41565b5b509250929050565b600081905092915050565b6000613f5b8385613f44565b9350613f688385846127a4565b82840190509392505050565b6000613f81828486613f4f565b9150819050939250505056fea2646970667358221220811fec16fe5ce732327b99430471548e9245369ec054d5ac5f20da3b4b1c9bb264736f6c63430008130033",
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

// TraceNestedEvents is a paid mutator transaction binding the contract method 0xb3f8d1b2.
//
// Solidity: function traceNestedEvents() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactor) TraceNestedEvents(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkDebugContract.contract.Transact(opts, "traceNestedEvents")
}

// TraceNestedEvents is a paid mutator transaction binding the contract method 0xb3f8d1b2.
//
// Solidity: function traceNestedEvents() returns()
func (_NetworkDebugContract *NetworkDebugContractSession) TraceNestedEvents() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceNestedEvents(&_NetworkDebugContract.TransactOpts)
}

// TraceNestedEvents is a paid mutator transaction binding the contract method 0xb3f8d1b2.
//
// Solidity: function traceNestedEvents() returns()
func (_NetworkDebugContract *NetworkDebugContractTransactorSession) TraceNestedEvents() (*types.Transaction, error) {
	return _NetworkDebugContract.Contract.TraceNestedEvents(&_NetworkDebugContract.TransactOpts)
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

// NetworkDebugContractUniqueDebugEventIterator is returned from FilterUniqueDebugEvent and is used to iterate over the raw logs and unpacked data for UniqueDebugEvent events raised by the NetworkDebugContract contract.
type NetworkDebugContractUniqueDebugEventIterator struct {
	Event *NetworkDebugContractUniqueDebugEvent // Event containing the contract specifics and raw log

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
func (it *NetworkDebugContractUniqueDebugEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkDebugContractUniqueDebugEvent)
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
		it.Event = new(NetworkDebugContractUniqueDebugEvent)
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
func (it *NetworkDebugContractUniqueDebugEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkDebugContractUniqueDebugEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkDebugContractUniqueDebugEvent represents a UniqueDebugEvent event raised by the NetworkDebugContract contract.
type NetworkDebugContractUniqueDebugEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUniqueDebugEvent is a free log retrieval operation binding the contract event 0xa0f7c7c1fff15178b5db3e56860767f0889c56b591bd2d9ba3121b491347d74c.
//
// Solidity: event UniqueDebugEvent()
func (_NetworkDebugContract *NetworkDebugContractFilterer) FilterUniqueDebugEvent(opts *bind.FilterOpts) (*NetworkDebugContractUniqueDebugEventIterator, error) {

	logs, sub, err := _NetworkDebugContract.contract.FilterLogs(opts, "UniqueDebugEvent")
	if err != nil {
		return nil, err
	}
	return &NetworkDebugContractUniqueDebugEventIterator{contract: _NetworkDebugContract.contract, event: "UniqueDebugEvent", logs: logs, sub: sub}, nil
}

// WatchUniqueDebugEvent is a free log subscription operation binding the contract event 0xa0f7c7c1fff15178b5db3e56860767f0889c56b591bd2d9ba3121b491347d74c.
//
// Solidity: event UniqueDebugEvent()
func (_NetworkDebugContract *NetworkDebugContractFilterer) WatchUniqueDebugEvent(opts *bind.WatchOpts, sink chan<- *NetworkDebugContractUniqueDebugEvent) (event.Subscription, error) {

	logs, sub, err := _NetworkDebugContract.contract.WatchLogs(opts, "UniqueDebugEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkDebugContractUniqueDebugEvent)
				if err := _NetworkDebugContract.contract.UnpackLog(event, "UniqueDebugEvent", log); err != nil {
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

// ParseUniqueDebugEvent is a log parse operation binding the contract event 0xa0f7c7c1fff15178b5db3e56860767f0889c56b591bd2d9ba3121b491347d74c.
//
// Solidity: event UniqueDebugEvent()
func (_NetworkDebugContract *NetworkDebugContractFilterer) ParseUniqueDebugEvent(log types.Log) (*NetworkDebugContractUniqueDebugEvent, error) {
	event := new(NetworkDebugContractUniqueDebugEvent)
	if err := _NetworkDebugContract.contract.UnpackLog(event, "UniqueDebugEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
