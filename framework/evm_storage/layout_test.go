package evm_storage_test

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/evm_storage"
)

const (
	contractAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	testRPCURL   = "http://localhost:8545"
	layoutFile   = "testdata/layout.json"
)

type testCase struct {
	Name                  string
	Slot                  string
	ValueHex              string
	ExpectValue           string
	AssertMethodSignature string
	AssertMethodArgs      []string
}

func prettyStructResult(output []byte) string {
	result := strings.TrimSpace(string(output))
	lines := strings.Split(result, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	joined := strings.Join(lines, " ")
	return joined
}

func setupAnvil(t *testing.T) {
	out, err := exec.Command("./setup.sh").CombinedOutput()
	fmt.Println(string(out))
	require.NoError(t, err)
	t.Cleanup(func() {
		out, err := exec.Command("./teardown.sh").CombinedOutput()
		fmt.Println(string(out))
		require.NoError(t, err)
	})
}

func TestStorageMutations(t *testing.T) {
	setupAnvil(t)
	// load contract layout file, see testdata/layout.json
	// more docs here - https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html#
	layout, err := evm_storage.New(layoutFile)
	if err != nil {
		t.Fatalf("failed to load layout: %v", err)
	}

	structPackingFunc := func(addr string, index uint8, group uint8) string {
		// huge structs can be packed differently to save space
		// comment the teardown and use this command to understand the layout
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

	testCases := []testCase{
		{
			Name:                  "Set number = 224",
			Slot:                  layout.Slot("number"),
			ValueHex:              "0x00000000000000000000000000000000000000000000000000000000000000E0",
			ExpectValue:           "224",
			AssertMethodSignature: "number()(uint256)",
		},
		{
			Name:                  "Set values[0] = 123",
			Slot:                  layout.ArraySlot("values", 0),
			ValueHex:              "0x000000000000000000000000000000000000000000000000000000000000007b",
			ExpectValue:           "123",
			AssertMethodSignature: "values(uint256)(uint256)",
			AssertMethodArgs:      []string{"0"},
		},
		{
			Name:                  "Set values[2] = 777",
			Slot:                  layout.ArraySlot("values", 2),
			ValueHex:              "0x0000000000000000000000000000000000000000000000000000000000000309",
			ExpectValue:           "777",
			AssertMethodSignature: "values(uint256)(uint256)",
			AssertMethodArgs:      []string{"2"},
		},
		{
			Name:                  "Set scores[dead] = 456",
			Slot:                  layout.MapSlot("scores", "0x000000000000000000000000000000000000dead"),
			ValueHex:              "0x00000000000000000000000000000000000000000000000000000000000001c8",
			ExpectValue:           "456",
			AssertMethodSignature: "scores(address)(uint256)",
			AssertMethodArgs:      []string{"0x000000000000000000000000000000000000dead"},
		},
		{
			Name:                  "Overwrite a_signers[0] with addr=a5, index=255, group=42",
			Slot:                  layout.ArraySlot("a_signers", 0),
			ValueHex:              structPackingFunc("0x00000000000000000000000000000000000000a5", 255, 42),
			AssertMethodSignature: "getASigner(uint256)(address,uint8,uint8)",
			AssertMethodArgs:      []string{"0"},
			ExpectValue:           "0x00000000000000000000000000000000000000A5 255 42",
		},
		{
			Name:                  "Overwrite a_signers[1] with addr=a5, index=255, group=42",
			Slot:                  layout.ArraySlot("a_signers", 1),
			ValueHex:              structPackingFunc("0x00000000000000000000000000000000000000a5", 255, 42),
			AssertMethodSignature: "getASigner(uint256)(address,uint8,uint8)",
			AssertMethodArgs:      []string{"1"},
			ExpectValue:           "0x00000000000000000000000000000000000000A5 255 42",
		},
		{
			Name:                  "Overwrite s_signers[0x5FbDB2315678afecb367f032d93F642f64180aa3] with addr=a6, index=12, group=34",
			Slot:                  layout.MapSlot("s_signers", "0x5cf8c07638e3be26449806d3dc21b622a946f877"),
			ValueHex:              structPackingFunc("0x00000000000000000000000000000000000000a6", 12, 34),
			AssertMethodSignature: "getSSigner(address)(address,uint8,uint8)",
			AssertMethodArgs:      []string{"0x5cf8c07638e3be26449806d3dc21b622a946f877"},
			ExpectValue:           "0x00000000000000000000000000000000000000a6 12 34",
		},
		{
			Name:                  "Overwrite s_signers[0x5FbDB2315678afecb367f032d93F642f64180aa4] with addr=a6, index=12, group=34",
			Slot:                  layout.MapSlot("s_signers", "0x5FbDB2315678afecb367f032d93F642f64180aa4"),
			ValueHex:              structPackingFunc("0x00000000000000000000000000000000000000a6", 14, 38),
			AssertMethodSignature: "getSSigner(address)(address,uint8,uint8)",
			AssertMethodArgs:      []string{"0x5FbDB2315678afecb367f032d93F642f64180aa4"},
			ExpectValue:           "0x00000000000000000000000000000000000000a6 14 38",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fmt.Println(tc.Name)
			callArgs := append([]string{
				"call", contractAddr, tc.AssertMethodSignature,
				"--rpc-url", testRPCURL,
			}, tc.AssertMethodArgs...)
			fmt.Println(callArgs)
			output, _ := exec.Command("cast", callArgs...).CombinedOutput()
			res := prettyStructResult(output)
			fmt.Println("Before:", res)

			fmt.Printf("Setting slot %s to %s\n", tc.Slot, tc.ValueHex)
			rpcArgs := []string{"rpc", "anvil_setStorageAt", contractAddr, tc.Slot, tc.ValueHex, "--rpc-url", testRPCURL}
			_, err := exec.Command("cast", rpcArgs...).CombinedOutput()
			if err != nil {
				t.Fatalf("set slot failed: %v", err)
			}

			if tc.AssertMethodArgs != nil {
				callArgs = append([]string{
					"call", contractAddr, tc.AssertMethodSignature,
					"--rpc-url", testRPCURL,
				}, tc.AssertMethodArgs...)
			} else {
				callArgs = []string{
					"call", contractAddr, tc.AssertMethodSignature,
					"--rpc-url", testRPCURL,
				}
			}
			fmt.Println(callArgs)
			output, err = exec.Command("cast", callArgs...).CombinedOutput()
			require.NoError(t, err)
			res = prettyStructResult(output)
			fmt.Println("After:", res)
			require.Equal(t, tc.ExpectValue, res)
		})
	}
}
