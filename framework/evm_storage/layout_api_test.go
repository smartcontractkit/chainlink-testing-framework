package evm_storage_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/evm_storage"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
)

// TestLayoutAPI that's the example of using helpers to override storage in your contracts
func TestLayoutAPI(t *testing.T) {
	//t.Skip("this test is for manual debugging and figuring out layout of custom structs")
	// load contract layout file, see testdata/layout.json
	// more docs here - https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html#
	layout, err := evm_storage.New(layoutFile)
	if err != nil {
		t.Fatalf("failed to load layout: %v", err)
	}

	encodeFunc := func(addr string, index uint8, group uint8) string {
		// huge structs can be packed differently to save space
		// cast storage 0x5FbDB2315678afecb367f032d93F642f64180aa3 0x1 --rpc-url http://localhost:8545
		addrBytes, _ := hex.DecodeString(strings.TrimPrefix(addr, "0x"))
		buf := make([]byte, evm_storage.StorageSlotSizeBytes)
		// start at the end of the slot, writing right to left
		idx := evm_storage.StorageSlotSizeBytes
		// 20 bytes for address
		idx -= 20
		copy(buf[idx:], addrBytes)
		// one byte for index
		idx--
		buf[idx] = index
		// one byte for group
		idx--
		buf[idx] = group
		return "0x" + hex.EncodeToString(buf)
	}

	{
		slot := layout.MustSlot("number_uint256")
		data := evm_storage.MustEncodeStorageSlot("uint256", big.NewInt(222))
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "number_uint256()(uint256)" --rpc-url http://localhost:8545
	}
	{
		slot := layout.MustSlot("number_uint8")
		data := evm_storage.MustEncodeStorageSlot("uint8", uint8(8))
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "number_uint8()(uint8)" --rpc-url http://localhost:8545
	}
	{
		slot := layout.MustSlot("number_int256")
		data := evm_storage.MustEncodeStorageSlot("uint256", big.NewInt(221))
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "number_int256()(int256)" --rpc-url http://localhost:8545
	}
	{
		slot := layout.MustSlot("_owner")
		data := evm_storage.MustEncodeStorageSlot("address", common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"))
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// private fields like _owner can only be verified by getter
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "getOwner()(address)" --rpc-url http://localhost:8545
	}
	{
		slot := layout.MustArraySlot("a_signers", 1)
		data := encodeFunc("0x00000000000000000000000000000000000000a5", 255, 42)
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "getASigner(uint256)(address,uint8,uint8)" --rpc-url http://localhost:8545 1
	}
	{
		slot := layout.MustMapSlot("s_signers", "0x00000000000000000000000000000000000000a5")
		data := encodeFunc("0x00000000000000000000000000000000000000a5", 254, 40)
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "getSSigner(address)(address,uint8,uint8)" 0x00000000000000000000000000000000000000a5 --rpc-url http://localhost:8545 1
	}
	{
		// offset example
		slot := layout.MustSlot("boolean")
		data := evm_storage.MustEncodeStorageSlot("bool", true)
		boolValue := evm_storage.ShiftHexByOffset(data, 1)
		uint8Value := evm_storage.MustEncodeStorageSlot("uint8", uint8(8))
		data = evm_storage.MergeHex(uint8Value, boolValue)
		fmt.Printf("setting slot: %s with data: %s\n", slot, data)
		r := rpc.New(testRPCURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
		require.NoError(t, err)
		// Contract code:
		// contract Counter {
		//    address private _owner;
		//    uint256 public number_uint256;
		//    int256 public number_int256;
		//    uint8 public number_uint8; <-- we need to change this
		//    bool public boolean;       <-- and this
		//
		// Example layout:
		// ╭----------------+-------------------------------------------+------+--------+-------+-------------------------╮
		// | Name           | Type                                      | Slot | Offset | Bytes | Contract                |
		// |----------------+-------------------------------------------+------+--------+-------+-------------------------|
		// | number_uint8   | uint8                                     | 3    | 0      | 1     | src/Counter.sol:Counter |
		// |----------------+-------------------------------------------+------+--------+-------+-------------------------|
		// | boolean        | bool                                      | 3    | 1      | 1     | src/Counter.sol:Counter |
		// |----------------+-------------------------------------------+------+--------+-------+-------------------------|
		// Resulting value with offsets: 0x0000000000000000000000000000000000000000000000000000000000000108
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "boolean()(bool)" --rpc-url http://localhost:8545
		// true
		// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "number_uint8()(uint8)" --rpc-url http://localhost:8545
		// 8
	}
}
