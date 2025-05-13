package evm_storage_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/evm_storage"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
)

// TestLayoutAPI that's the example of using helpers to override storage in your contracts
func TestLayoutAPI(t *testing.T) {
	t.Skip("this test is for manual debugging and figuring out layout of custom structs")
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

	slot := layout.ArraySlot("a_signers", 1)
	data := encodeFunc("0x00000000000000000000000000000000000000a5", 255, 42)
	fmt.Printf("setting slot: %s with data: %s\n", slot, data)
	r := rpc.New(testRPCURL, nil)
	err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
	require.NoError(t, err)
	// verify it manually
	// cast call 0x5FbDB2315678afecb367f032d93F642f64180aa3 "getASigner(uint256)(address,uint8,uint8)" --rpc-url http://localhost:8545 1
}
