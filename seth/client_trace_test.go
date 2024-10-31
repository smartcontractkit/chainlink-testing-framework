package seth_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	network_debug_contract "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/NetworkDebugContract"
	"github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/TestContractOne"
	"github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/TestContractTwo"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"
	"github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link_token_interface"
	"github.com/smartcontractkit/chainlink-testing-framework/seth/test_utils"
)

const (
	NoAnvilSupport = "Anvil doesn't support tracing"
	FailedToDecode = "failed to decode transaction"
)

func SkipAnvil(t *testing.T, c *seth.Client) {
	if c.Cfg.Network.Name == "Anvil" {
		t.Skip(NoAnvilSupport)
	}
}

// since we uploaded the contracts via Seth, we have the contract address in the map
// and we can trace the calls correctly even though both calls have the same signature
func TestTraceContractTracingSameMethodSignatures_UploadedViaSeth(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	var x int64 = 2
	var y int64 = 4
	tx, err := c.Decode(TestEnv.DebugContract.Trace(c.NewTXOpts(), big.NewInt(x), big.NewInt(y)))
	require.NoError(t, err, FailedToDecode)

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")
	require.Equal(t, 2, len(c.Tracer.GetDecodedCalls(tx.Hash)), "expected 2 decoded calls for this transaction")

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())

	firstExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "3e41f135",
			CallType:  "CALL",
			Method:    "trace(int256,int256)",
			Input:     map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y)},
			Output:    map[string]interface{}{"0": big.NewInt(y + 2)},
		},
		Events: []seth.DecodedCommonLog{
			{
				Signature: "TwoIndexEvent(uint256,address)",
				EventData: map[string]interface{}{"roundId": big.NewInt(y), "startedBy": c.Addresses[0]},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5",
					"0x0000000000000000000000000000000000000000000000000000000000000004",
					"0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				},
			},
		},
		Comment: "",
	}

	require.EqualValues(t, firstExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "first decoded call does not match")

	secondExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugSubContractAddress.Hex()),
		From:        "NetworkDebugContract",
		To:          "NetworkDebugSubContract",
		CommonData: seth.CommonData{
			Signature:       "3e41f135",
			ParentSignature: "3e41f135",
			CallType:        "CALL",
			NestingLevel:    1,
			Method:          "trace(int256,int256)",
			Input:           map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y)},
			Output:          map[string]interface{}{"0": big.NewInt(y + 4)},
		},
		Comment: "",
	}

	actualSecondEvents := c.Tracer.GetDecodedCalls(tx.Hash)[1].Events
	c.Tracer.GetDecodedCalls(tx.Hash)[1].Events = nil

	require.EqualValues(t, secondExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[1], "second decoded call does not match")
	require.Equal(t, 1, len(actualSecondEvents), "second decoded call events count does not match")
	require.Equal(t, 3, len(actualSecondEvents[0].Topics), "second decoded event topics count does not match")

	expectedSecondEvents := []seth.DecodedCommonLog{
		{
			Signature: "TwoIndexEvent(uint256,address)",
			EventData: map[string]interface{}{"roundId": big.NewInt(6), "startedBy": TestEnv.DebugContractAddress},
			Address:   TestEnv.DebugSubContractAddress,
			Topics: []string{
				"0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5",
				"0x0000000000000000000000000000000000000000000000000000000000000006",
				// this one changes dynamically depending on sender address
				// "0x000000000000000000000000c351628eb244ec633d5f21fbd6621e1a683b1181",
			},
		},
	}
	actualSecondEvents[0].Topics = actualSecondEvents[0].Topics[0:2]
	require.EqualValues(t, expectedSecondEvents, actualSecondEvents, "second decoded call events do not match")
}

// we test a scenario, where because two contracts have the same method signature, both addresses
// were mapped to the same contract name (it doesn't happen always, it all depends on how data is ordered
// in the maps and that depends on addresses generated). We show that even if the initial mapping is incorrect,
// once we trace a transaction with different method signature, the mapping is corrected and the second transaction
// is traced correctly.
func TestTraceContractTracingSameMethodSignatures_UploadedManually(t *testing.T) {
	c := newClient(t)
	SkipAnvil(t, c)

	for k := range c.ContractAddressToNameMap.GetContractMap() {
		delete(c.ContractAddressToNameMap.GetContractMap(), k)
	}

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	// let's simulate this case, because it doesn't happen always, it all depends on the order of the
	// contract map, which is non-deterministic (hash map with keys being dynamically generated addresses)
	c.ContractAddressToNameMap.AddContract(TestEnv.DebugContractAddress.Hex(), "NetworkDebugContract")
	c.ContractAddressToNameMap.AddContract(TestEnv.DebugSubContractAddress.Hex(), "NetworkDebugContract")

	var x int64 = 2
	var y int64 = 4

	diffSigTx, txErr := c.Decode(TestEnv.DebugContract.TraceDifferent(c.NewTXOpts(), big.NewInt(x), big.NewInt(y)))
	require.NoError(t, txErr, FailedToDecode)
	sameSigTx, txErr := c.Decode(TestEnv.DebugContract.Trace(c.NewTXOpts(), big.NewInt(x), big.NewInt(y)))
	require.NoError(t, txErr, FailedToDecode)

	require.NotNil(t, c.Tracer.GetDecodedCalls(diffSigTx.Hash), "expected decoded calls to contain the diffSig transaction hash")
	require.Equal(t, 2, len(c.Tracer.GetDecodedCalls(diffSigTx.Hash)), "expected 2 decoded calls for diffSig transaction")

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())

	firstDiffSigCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "30985bcc",
			CallType:  "CALL",
			Method:    "traceDifferent(int256,int256)",
			Input:     map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y)},
			Output:    map[string]interface{}{"0": big.NewInt(y + 2)},
		},
		Events: []seth.DecodedCommonLog{
			{
				Signature: "OneIndexEvent(uint256)",
				EventData: map[string]interface{}{"a": big.NewInt(x)},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35",
					"0x0000000000000000000000000000000000000000000000000000000000000002",
				},
			},
		},
		Comment: "",
	}

	require.EqualValues(t, firstDiffSigCall, c.Tracer.GetDecodedCalls(diffSigTx.Hash)[0], "first diffSig decoded call does not match")

	secondDiffSigCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugSubContractAddress.Hex()),
		From:        "NetworkDebugContract",
		To:          "NetworkDebugSubContract",
		CommonData: seth.CommonData{
			Signature:       "047c4425",
			ParentSignature: "30985bcc",
			NestingLevel:    1,
			CallType:        "CALL",
			Method:          "traceOneInt(int256)",
			Input:           map[string]interface{}{"x": big.NewInt(x + 2)},
			Output:          map[string]interface{}{"r": big.NewInt(y + 3)},
		},
		Comment: "",
	}

	c.Tracer.GetDecodedCalls(diffSigTx.Hash)[1].Events = nil
	require.EqualValues(t, secondDiffSigCall, c.Tracer.GetDecodedCalls(diffSigTx.Hash)[1], "second diffSig decoded call does not match")

	require.Equal(t, 2, len(c.Tracer.GetAllDecodedCalls()), "expected 2 decoded transactons")
	require.NotNil(t, c.Tracer.GetDecodedCalls(sameSigTx.Hash), "expected decoded calls to contain the sameSig transaction hash")
	require.Equal(t, 2, len(c.Tracer.GetDecodedCalls(sameSigTx.Hash)), "expected 2 decoded calls for sameSig transaction")

	firstSameSigCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "3e41f135",
			CallType:  "CALL",
			Method:    "trace(int256,int256)",
			Input:     map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y)},
			Output:    map[string]interface{}{"0": big.NewInt(y + 2)},
		},
		Events: []seth.DecodedCommonLog{
			{
				Signature: "TwoIndexEvent(uint256,address)",
				EventData: map[string]interface{}{"roundId": big.NewInt(y), "startedBy": c.Addresses[0]},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5",
					"0x0000000000000000000000000000000000000000000000000000000000000004",
					"0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				},
			},
		},
		Comment: "",
	}

	require.EqualValues(t, firstSameSigCall, c.Tracer.GetDecodedCalls(sameSigTx.Hash)[0], "first sameSig decoded call does not match")

	secondSameSigCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugSubContractAddress.Hex()),
		From:        "NetworkDebugContract",
		To:          "NetworkDebugSubContract",
		CommonData: seth.CommonData{
			Signature:       "3e41f135",
			ParentSignature: "3e41f135",
			NestingLevel:    1,
			CallType:        "CALL",
			Method:          "trace(int256,int256)",
			Input:           map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y)},
			Output:          map[string]interface{}{"0": big.NewInt(y + 4)},
		},
	}

	actualSecondEvents := c.Tracer.GetDecodedCalls(sameSigTx.Hash)[1].Events
	c.Tracer.GetDecodedCalls(sameSigTx.Hash)[1].Events = nil

	require.EqualValues(t, secondSameSigCall, c.Tracer.GetDecodedCalls(sameSigTx.Hash)[1], "second sameSig decoded call does not match")
	require.Equal(t, 1, len(actualSecondEvents), "second sameSig decoded call events count does not match")
	require.Equal(t, 3, len(actualSecondEvents[0].Topics), "second sameSig decoded event topics count does not match")

	expectedSecondEvents := []seth.DecodedCommonLog{
		{
			Signature: "TwoIndexEvent(uint256,address)",
			EventData: map[string]interface{}{"roundId": big.NewInt(6), "startedBy": TestEnv.DebugContractAddress},
			Address:   TestEnv.DebugSubContractAddress,
			Topics: []string{
				"0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5",
				"0x0000000000000000000000000000000000000000000000000000000000000006",
				// third topic changes dynamically depending on sender address
				// "0x000000000000000000000000c351628eb244ec633d5f21fbd6621e1a683b1181",
			},
		},
	}
	actualSecondEvents[0].Topics = actualSecondEvents[0].Topics[0:2]
	require.EqualValues(t, expectedSecondEvents, actualSecondEvents, "second sameSig decoded call events do not match")
}

func TestTraceContractTracingSameMethodSignaturesWarningInComment_UploadedManually(t *testing.T) {
	c := newClient(t)
	SkipAnvil(t, c)

	c.ContractAddressToNameMap = seth.NewEmptyContractMap()

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	sameSigTx, err := c.Decode(TestEnv.DebugContract.Trace(c.NewTXOpts(), big.NewInt(2), big.NewInt(2)))
	require.NoError(t, err, "failed to send transaction")

	require.NotNil(t, c.Tracer.GetDecodedCalls(sameSigTx.Hash), "expected decoded calls to contain the transaction hash")
	require.Equal(t, 2, len(c.Tracer.GetDecodedCalls(sameSigTx.Hash)), "expected 2 decoded calls for transaction")
	require.Equal(t, "potentially inaccurate - method present in 1 other contracts", c.Tracer.GetDecodedCalls(sameSigTx.Hash)[1].Comment, "expected comment to be set")
}

func TestTraceContractTracingWithCallback_UploadedViaSeth(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	// As this test might fail if run multiple times due to non-deterministic addressed in contract mapping
	// which sometime causes the call to be traced and sometimes not (it all depends on the order of
	// addresses in the map), I just remove potentially problematic ABI.
	delete(c.ContractStore.ABIs, "DebugContractCallback.abi")

	var x int64 = 2
	var y int64 = 4
	tx, txErr := c.Decode(TestEnv.DebugContract.TraceSubWithCallback(c.NewTXOpts(), big.NewInt(x), big.NewInt(y)))
	require.NoError(t, txErr, FailedToDecode)

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")
	require.Equal(t, 3, len(c.Tracer.GetDecodedCalls(tx.Hash)), "expected 2 decoded calls for test transaction")

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())

	firstExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			CallType:  "CALL",
			Signature: "3837a75e",
			Method:    "traceSubWithCallback(int256,int256)",
			Input:     map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y)},
			Output:    map[string]interface{}{"0": big.NewInt(y + 4)},
		},
		Events: []seth.DecodedCommonLog{
			{
				Signature: "TwoIndexEvent(uint256,address)",
				EventData: map[string]interface{}{"roundId": big.NewInt(1), "startedBy": c.Addresses[0]},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5",
					"0x0000000000000000000000000000000000000000000000000000000000000001",
					"0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				},
			},
		},
		Comment: "",
	}

	require.EqualValues(t, firstExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "first decoded call does not match")

	require.Equal(t, 2, len(c.Tracer.GetDecodedCalls(tx.Hash)[1].Events), "second decoded call events count does not match")
	require.Equal(t, 3, len(c.Tracer.GetDecodedCalls(tx.Hash)[1].Events[0].Topics), "second decoded first event topics count does not match")

	separatedTopcis := c.Tracer.GetDecodedCalls(tx.Hash)[1].Events[0].Topics
	separatedTopcis = separatedTopcis[0:2]
	c.Tracer.GetDecodedCalls(tx.Hash)[1].Events[0].Topics = nil

	secondExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugSubContractAddress.Hex()),
		From:        "NetworkDebugContract",
		To:          "NetworkDebugSubContract",
		CommonData: seth.CommonData{
			Signature:       "fa8fca7a",
			CallType:        "CALL",
			Method:          "traceWithCallback(int256,int256)",
			NestingLevel:    1,
			ParentSignature: "3837a75e",
			Input:           map[string]interface{}{"x": big.NewInt(x), "y": big.NewInt(y + 2)},
			Output:          map[string]interface{}{"0": big.NewInt(y + 2)},
		},
		Events: []seth.DecodedCommonLog{
			{
				Signature: "TwoIndexEvent(uint256,address)",
				EventData: map[string]interface{}{"roundId": big.NewInt(6), "startedBy": TestEnv.DebugContractAddress},
				Address:   TestEnv.DebugSubContractAddress,
			},
			{
				Signature: "OneIndexEvent(uint256)",
				EventData: map[string]interface{}{"a": big.NewInt(y + 2)},
				Address:   TestEnv.DebugSubContractAddress,
				Topics: []string{
					"0xeace1be0b97ec11f959499c07b9f60f0cc47bf610b28fda8fb0e970339cf3b35",
					"0x0000000000000000000000000000000000000000000000000000000000000006",
				},
			},
		},
		Comment: "",
	}

	require.EqualValues(t, secondExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[1], "second decoded call does not match")

	expectedTopics := []string{
		"0x33b47a1cd66813164ec00800d74296f57415217c22505ee380594a712936a0b5",
		"0x0000000000000000000000000000000000000000000000000000000000000006",
		// third topic is dynamic (sender address), skip it
		// "0x00000000000000000000000056fc17a65ccfec6b7ad0ade9bd9416cb365b9be8",
	}

	require.EqualValues(t, expectedTopics, separatedTopcis, "second decoded first event topics do not match")

	thirdExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(TestEnv.DebugSubContractAddress.Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "NetworkDebugSubContract",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature:       "fbcb8d07",
			CallType:        "CALL",
			Method:          "callbackMethod(int256)",
			ParentSignature: "fa8fca7a",
			NestingLevel:    2,
			Input:           map[string]interface{}{"x": big.NewInt(x + y)},
			Output:          map[string]interface{}{"0": big.NewInt(y + x)},
		},
		Events: []seth.DecodedCommonLog{
			{
				Signature: "CallbackEvent(int256)",
				EventData: map[string]interface{}{"a": big.NewInt(y + 2)},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0xb16dba9242e1aa07ccc47228094628f72c8cc9699ee45d5bc8d67b84d3038c68",
					"0x0000000000000000000000000000000000000000000000000000000000000006",
				},
			},
		},
	}
	require.EqualValues(t, thirdExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[2], "third decoded call does not match")
}

// Here we show that partial tracing works even if we don't have the ABI for the contract.
// We still try to decode what we can even without ABI and that we can decode the other call
// for which we do have ABI.
func TestTraceContractTracingUnknownAbi(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	// simulate missing ABI
	delete(c.ContractAddressToNameMap.GetContractMap(), strings.ToLower(TestEnv.DebugContractAddress.Hex()))
	delete(c.ContractStore.ABIs, "NetworkDebugContract.abi")

	var x int64 = 2
	var y int64 = 4
	tx, txErr := c.Decode(TestEnv.DebugContract.TraceDifferent(c.NewTXOpts(), big.NewInt(x), big.NewInt(y)))
	require.NoError(t, txErr, FailedToDecode)

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")
	require.Equal(t, 2, len(c.Tracer.GetDecodedCalls(tx.Hash)), "expected 2 decoded calls for test transaction")

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())

	firstExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          seth.UNKNOWN,
		CommonData: seth.CommonData{
			Signature: "30985bcc",
			CallType:  "CALL",
			Method:    seth.UNKNOWN,
			Input:     make(map[string]interface{}),
			Output:    make(map[string]interface{}),
		},
		Events:  []seth.DecodedCommonLog{},
		Comment: seth.CommentMissingABI,
	}

	require.EqualValues(t, firstExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "first decoded call does not match")

	secondExpectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugSubContractAddress.Hex()),
		From:        seth.UNKNOWN,
		To:          "NetworkDebugSubContract",
		CommonData: seth.CommonData{
			Signature:       "047c4425",
			CallType:        "CALL",
			ParentSignature: "30985bcc",
			NestingLevel:    1,
			Method:          "traceOneInt(int256)",
			Input:           map[string]interface{}{"x": big.NewInt(x + 2)},
			Output:          map[string]interface{}{"r": big.NewInt(y + 3)},
		},
		Comment: "",
	}

	c.Tracer.GetDecodedCalls(tx.Hash)[1].Events = nil
	require.EqualValues(t, secondExpectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[1], "second decoded call does not match")
}

func TestTraceContractTracingNamedInputsAndOutputs(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	x := big.NewInt(1000)
	var testString = "string"
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitNamedInputsOutputs(c.NewTXOpts(), x, testString))
	require.NoError(t, txErr, FailedToDecode)

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "45f0c9e6",
			CallType:  "CALL",
			Method:    "emitNamedInputsOutputs(uint256,string)",
			Input:     map[string]interface{}{"inputVal1": x, "inputVal2": testString},
			Output:    map[string]interface{}{"outputVal1": x, "outputVal2": testString},
		},
		Comment: "",
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingNamedInputsAnonymousOutputs(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	x := big.NewInt(1001)
	var testString = "string"
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitInputsOutputs(c.NewTXOpts(), x, testString))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "d7a80205",
			CallType:  "CALL",
			Method:    "emitInputsOutputs(uint256,string)",
			Input:     map[string]interface{}{"inputVal1": x, "inputVal2": testString},
			Output:    map[string]interface{}{"0": x, "1": testString},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

// Shows that when output mixes named and unnamed parameters, we can still decode the transaction,
// but that named outputs become unnamed and referenced by their index.
func TestTraceContractTracingIntInputsWithoutLength(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	x := big.NewInt(1001)
	y := big.NewInt(2)
	z := big.NewInt(26)
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitInts(c.NewTXOpts(), x, y, z))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "9e099652",
			CallType:  "CALL",
			Method:    "emitInts(int256,int128,uint256)",
			Input:     map[string]interface{}{"first": x, "second": y, "third": z},
			Output:    map[string]interface{}{"0": x, "1": y, "2": z},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingAddressInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	address := c.Addresses[0]
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitAddress(c.NewTXOpts(), address))
	require.NoError(t, txErr, "failed to send transaction")

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "ec5c3ede",
			CallType:  "CALL",
			Method:    "emitAddress(address)",
			Input:     map[string]interface{}{"addr": address},
			Output:    map[string]interface{}{"0": address},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingBytes32InputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	addrAsBytes := c.Addresses[0].Bytes()
	addrAsBytes = append(addrAsBytes, c.Addresses[0].Bytes()...)
	var bytes32 = [32]byte(addrAsBytes)
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitBytes32(c.NewTXOpts(), bytes32))
	require.NoError(t, txErr, FailedToDecode)

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "33311ef3",
			CallType:  "CALL",
			Method:    "emitBytes32(bytes32)",
			Input:     map[string]interface{}{"input": bytes32},
			Output:    map[string]interface{}{"output": bytes32},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingUint256ArrayInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	uint256Array := []*big.Int{big.NewInt(1), big.NewInt(19271), big.NewInt(261), big.NewInt(271911), big.NewInt(821762721)}
	tx, txErr := c.Decode(TestEnv.DebugContract.ProcessUintArray(c.NewTXOpts(), uint256Array))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	var output []*big.Int
	for _, x := range uint256Array {
		output = append(output, big.NewInt(0).Add(x, big.NewInt(1)))
	}

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "12d91233",
			CallType:  "CALL",
			Method:    "processUintArray(uint256[])",
			Input:     map[string]interface{}{"input": uint256Array},
			Output:    map[string]interface{}{"0": output},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingAddressArrayInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	addressArray := []common.Address{c.Addresses[0], TestEnv.DebugSubContractAddress}
	tx, txErr := c.Decode(TestEnv.DebugContract.ProcessAddressArray(c.NewTXOpts(), addressArray))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "e1111f79",
			CallType:  "CALL",
			Method:    "processAddressArray(address[])",
			Input:     map[string]interface{}{"input": addressArray},
			Output:    map[string]interface{}{"0": addressArray},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingStructWithDynamicFieldsInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	data := network_debug_contract.NetworkDebugContractData{
		Name:   "my awesome name",
		Values: []*big.Int{big.NewInt(2), big.NewInt(266810), big.NewInt(473878233)},
	}
	tx, txErr := c.Decode(TestEnv.DebugContract.ProcessDynamicData(c.NewTXOpts(), data))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expected := struct {
		Name   string     `json:"name"`
		Values []*big.Int `json:"values"`
	}{
		Name:   data.Name,
		Values: data.Values,
	}

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "7fdc8fe1",
			CallType:  "CALL",
			Method:    "processDynamicData((string,uint256[]))",
			Input:     map[string]interface{}{"data": expected},
			Output:    map[string]interface{}{"0": expected},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingStructArrayWithDynamicFieldsInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	data := network_debug_contract.NetworkDebugContractData{
		Name:   "my awesome name",
		Values: []*big.Int{big.NewInt(2), big.NewInt(266810), big.NewInt(473878233)},
	}
	dataArray := [3]network_debug_contract.NetworkDebugContractData{data, data, data}
	tx, txErr := c.Decode(TestEnv.DebugContract.ProcessFixedDataArray(c.NewTXOpts(), dataArray))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	input := [3]struct {
		Name   string     `json:"name"`
		Values []*big.Int `json:"values"`
	}{
		{
			Name:   data.Name,
			Values: data.Values,
		},
		{
			Name:   data.Name,
			Values: data.Values,
		},
		{
			Name:   data.Name,
			Values: data.Values,
		},
	}

	output := [2]struct {
		Name   string     `json:"name"`
		Values []*big.Int `json:"values"`
	}{
		{
			Name:   data.Name,
			Values: data.Values,
		},
		{
			Name:   data.Name,
			Values: data.Values,
		},
	}

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "99adad2e",
			CallType:  "CALL",
			Method:    "processFixedDataArray((string,uint256[])[3])",
			Input:     map[string]interface{}{"data": input},
			Output:    map[string]interface{}{"0": output},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingNestedStructsWithDynamicFieldsInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	data := network_debug_contract.NetworkDebugContractNestedData{
		Data: network_debug_contract.NetworkDebugContractData{
			Name:   "my awesome name",
			Values: []*big.Int{big.NewInt(2), big.NewInt(266810), big.NewInt(473878233)},
		},
		DynamicBytes: []byte("dynamic bytes"),
	}
	tx, txErr := c.Decode(TestEnv.DebugContract.ProcessNestedData(c.NewTXOpts(), data))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	input := struct {
		Data struct {
			Name   string     `json:"name"`
			Values []*big.Int `json:"values"`
		} `json:"data"`
		DynamicBytes []byte `json:"dynamicBytes"`
	}{
		struct {
			Name   string     `json:"name"`
			Values []*big.Int `json:"values"`
		}{
			Name:   data.Data.Name,
			Values: data.Data.Values,
		},
		data.DynamicBytes,
	}

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "7f12881c",
			CallType:  "CALL",
			Method:    "processNestedData(((string,uint256[]),bytes))",
			Input:     map[string]interface{}{"data": input},
			Output:    map[string]interface{}{"0": input},
		},
		Comment: "",
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingNestedStructsWithDynamicFieldsInputAndStructOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	data := network_debug_contract.NetworkDebugContractData{
		Name:   "my awesome name",
		Values: []*big.Int{big.NewInt(2), big.NewInt(266810), big.NewInt(473878233)},
	}
	tx, txErr := c.Decode(TestEnv.DebugContract.ProcessNestedData0(c.NewTXOpts(), data))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	input := struct {
		Name   string     `json:"name"`
		Values []*big.Int `json:"values"`
	}{
		Name:   data.Name,
		Values: data.Values,
	}

	hash := crypto.Keccak256Hash([]byte(input.Name))

	output := struct {
		Data struct {
			Name   string     `json:"name"`
			Values []*big.Int `json:"values"`
		} `json:"data"`
		DynamicBytes []byte `json:"dynamicBytes"`
	}{
		struct {
			Name   string     `json:"name"`
			Values []*big.Int `json:"values"`
		}{
			Name:   data.Name,
			Values: data.Values,
		},
		hash.Bytes(),
	}

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "f499af2a",
			CallType:  "CALL",
			Method:    "processNestedData((string,uint256[]))",
			Input:     map[string]interface{}{"data": input},
			Output:    map[string]interface{}{"0": output},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingPayable(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	var value int64 = 1000
	tx, txErr := c.Decode(TestEnv.DebugContract.Pay(c.NewTXOpts(seth.WithValue(big.NewInt(value)))))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "1b9265b8",
			CallType:  "CALL",
			Method:    "pay()",
			Output:    map[string]interface{}{},
		},
		Value: value,
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingFallback(t *testing.T) {
	t.Skip("Need to investigate further how to support it, the call succeeds, but we fail to decode it")
	// our ABIFinder doesn't know anything about fallback, but maybe we should use it, when everything else fails?
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	tx, txErr := c.Decode(TestEnv.DebugContractRaw.RawTransact(c.NewTXOpts(), []byte("iDontExist")))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transaction")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "1b9265b8",
			CallType:  "CALL",
			Method:    "pay()",
			Output:    map[string]interface{}{},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingReceive(t *testing.T) {
	t.Skip("Need to investigate further how to support it, the call succreds, but we fail to match the signature as input is 0x")
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	value := big.NewInt(29121)
	tx, txErr := c.Decode(TestEnv.DebugContract.Receive(c.NewTXOpts(seth.WithValue(value))))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transaction")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "1b9265b8",
			Method:    "pay()",
			CallType:  "CALL",
			Output:    map[string]interface{}{},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingEnumInputAndOutput(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	var status uint8 = 1 // Active
	tx, txErr := c.Decode(TestEnv.DebugContract.SetStatus(c.NewTXOpts(), status))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "2e49d78b",
			CallType:  "CALL",
			Method:    "setStatus(uint8)",
			Input:     map[string]interface{}{"status": status},
			Output:    map[string]interface{}{"0": status},
		},
		Comment: "",
		Events: []seth.DecodedCommonLog{
			{
				Signature: "CurrentStatus(uint8)",
				EventData: map[string]interface{}{"status": status},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0xbea054406fdf249b05d1aef1b5f848d62d902d94389fca702b2d8337677c359a",
					"0x0000000000000000000000000000000000000000000000000000000000000001",
				},
			},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingNonIndexedEventParameter(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	tx, txErr := c.Decode(TestEnv.DebugContract.EmitNoIndexEventString(c.NewTXOpts()))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "788c4772",
			CallType:  "CALL",
			Method:    "emitNoIndexEventString()",
			Input:     nil,
			Output:    map[string]interface{}{},
		},
		Comment: "",
		Events: []seth.DecodedCommonLog{
			{
				Signature: "NoIndexEventString(string)",
				EventData: map[string]interface{}{"str": "myString"},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0x25b7adba1b046a19379db4bc06aa1f2e71604d7b599a0ee8783d58110f00e16a",
				},
			},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingEventThreeIndexedParameters(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	tx, txErr := c.Decode(TestEnv.DebugContract.EmitThreeIndexEvent(c.NewTXOpts()))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "aa3fdcf4",
			CallType:  "CALL",
			Method:    "emitThreeIndexEvent()",
			Input:     nil,
			Output:    map[string]interface{}{},
		},
		Comment: "",
		Events: []seth.DecodedCommonLog{
			{
				Signature: "ThreeIndexEvent(uint256,address,uint256)",
				EventData: map[string]interface{}{"roundId": big.NewInt(1), "startedAt": big.NewInt(3), "startedBy": c.Addresses[0]},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0x5660e8f93f0146f45abcd659e026b75995db50053cbbca4d7f365934ade68bf3",
					"0x0000000000000000000000000000000000000000000000000000000000000001",
					"0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
					"0x0000000000000000000000000000000000000000000000000000000000000003",
				},
			},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingEventFourMixedParameters(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	tx, txErr := c.Decode(TestEnv.DebugContract.EmitFourParamMixedEvent(c.NewTXOpts()))
	require.NoError(t, txErr, FailedToDecode)
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "c2124b22",
			CallType:  "CALL",
			Method:    "emitFourParamMixedEvent()",
			Input:     nil,
			Output:    map[string]interface{}{},
		},
		Comment: "",
		Events: []seth.DecodedCommonLog{
			{
				Signature: "ThreeIndexAndOneNonIndexedEvent(uint256,address,uint256,string)",
				EventData: map[string]interface{}{"roundId": big.NewInt(2), "startedAt": big.NewInt(3), "startedBy": c.Addresses[0], "dataId": "some id"},
				Address:   TestEnv.DebugContractAddress,
				Topics: []string{
					"0x56c2ea44ba516098cee0c181dd9d8db262657368b6e911e83ae0ccfae806c73d",
					"0x0000000000000000000000000000000000000000000000000000000000000002",
					"0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
					"0x0000000000000000000000000000000000000000000000000000000000000003",
				},
			},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractAll(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_All
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	revertedTx, txErr := TestEnv.DebugContract.AlwaysRevertsCustomError(c.NewTXOpts())
	require.NoError(t, txErr, "transaction sending should not fail")
	_, decodeErr := c.Decode(revertedTx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [12 21]", decodeErr.Error(), "expected error message to contain the reverted error type and values")

	okTx, txErr := TestEnv.DebugContract.AddCounter(c.NewTXOpts(), big.NewInt(1), big.NewInt(2))
	require.NoError(t, txErr, "transaction should not have reverted")
	_, decodeErr = c.Decode(okTx, txErr)
	require.NoError(t, decodeErr, "transaction decoding should not err")
	require.Equal(t, 2, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "5e9c80d6",
			CallType:  "CALL",
			Method:    "alwaysRevertsCustomError()",
			Output:    map[string]interface{}{},
			Error:     "execution reverted",
		},
	}

	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(revertedTx.Hash().Hex())[0], "reverted decoded call does not match")

	expectedCall = &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "23515760",
			CallType:  "CALL",
			Method:    "addCounter(int256,int256)",
			Output:    map[string]interface{}{"value": big.NewInt(2)},
			Input:     map[string]interface{}{"idx": big.NewInt(1), "x": big.NewInt(2)},
		},
	}

	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(okTx.Hash().Hex())[0], "successful decoded call does not match")
}

func TestTraceContractOnlyReverted(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_Reverted
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	revertedTx, txErr := TestEnv.DebugContract.AlwaysRevertsCustomError(c.NewTXOpts())
	require.NoError(t, txErr, "transaction sending should not fail")
	_, decodeErr := c.Decode(revertedTx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [12 21]", decodeErr.Error(), "expected error message to contain the reverted error type and values")

	okTx, txErr := TestEnv.DebugContract.AddCounter(c.NewTXOpts(), big.NewInt(1), big.NewInt(2))
	require.NoError(t, txErr, "transaction should not have reverted")
	_, decodeErr = c.Decode(okTx, txErr)
	require.NoError(t, decodeErr, "transaction decoding should not err")

	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "5e9c80d6",
			CallType:  "CALL",
			Method:    "alwaysRevertsCustomError()",
			Output:    map[string]interface{}{},
			Error:     "execution reverted",
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(revertedTx.Hash().Hex())[0], "decoded call does not match")
}

func TestTraceContractNone(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	// when this level nothing is ever traced or debugged
	c.Cfg.TracingLevel = seth.TracingLevel_None

	revertedTx, txErr := TestEnv.DebugContract.AlwaysRevertsCustomError(c.NewTXOpts())
	require.NoError(t, txErr, "transaction sending should not fail")
	_, decodeErr := c.Decode(revertedTx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [12 21]", decodeErr.Error(), "expected error message to contain the reverted error type and values")

	okTx, txErr := TestEnv.DebugContract.AddCounter(c.NewTXOpts(), big.NewInt(1), big.NewInt(2))
	require.NoError(t, txErr, "transaction should not have reverted")
	_, decodeErr = c.Decode(okTx, txErr)
	require.NoError(t, decodeErr, "transaction decoding should not err")

	require.Empty(t, c.Tracer.GetAllDecodedCalls(), "expected 1 decoded transaction")
}

func TestTraceContractRevertedErrNoValues(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_Reverted
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	tx, txErr := TestEnv.DebugContract.AlwaysRevertsCustomErrorNoValues(c.NewTXOpts())
	require.NoError(t, txErr, "transaction should have reverted")
	_, decodeErr := c.Decode(tx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErrNoValues, error values: []", decodeErr.Error(), "expected error message to contain the reverted error type and values")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "b600141f",
			CallType:  "CALL",
			Method:    "alwaysRevertsCustomErrorNoValues()",
			Output:    map[string]interface{}{},
			Error:     "execution reverted",
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash().Hex())[0], "decoded call does not match")
}

func TestTraceCallRevertFunctionInTheContract(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_Reverted
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	tx, txErr := TestEnv.DebugContract.CallRevertFunctionInTheContract(c.NewTXOpts())
	require.NoError(t, txErr, "transaction should have reverted")
	_, decodeErr := c.Decode(tx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [12 21]", decodeErr.Error(), "expected error message to contain the reverted error type and values")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "9349d00b",
			CallType:  "CALL",
			Method:    "callRevertFunctionInTheContract()",
			Output:    map[string]interface{}{},
			Error:     "execution reverted",
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash().Hex())[0], "decoded call does not match")
}

func TestTraceCallRevertFunctionInSubContract(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_Reverted
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	x := big.NewInt(1001)
	y := big.NewInt(2)
	tx, txErr := TestEnv.DebugContract.CallRevertFunctionInSubContract(c.NewTXOpts(), x, y)
	require.NoError(t, txErr, "transaction should have reverted")
	_, decodeErr := c.Decode(tx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [1001 2]", decodeErr.Error(), "expected error message to contain the reverted error type and values")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transaction")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "11b3c478",
			CallType:  "CALL",
			Method:    "callRevertFunctionInSubContract(uint256,uint256)",
			Input:     map[string]interface{}{"x": x, "y": y},
			Output:    map[string]interface{}{},
			Error:     "execution reverted",
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash().Hex())[0], "decoded call does not match")
}

func TestTraceCallRevertInCallback(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_Reverted
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")
	c.ContractStore.AddABI("LinkToken", *linkAbi)
	c.ContractStore.AddBIN("LinkToken", common.FromHex(link_token.LinkTokenMetaData.Bin))

	amount := big.NewInt(0)
	tx, txErr := TestEnv.LinkTokenContract.TransferAndCall(c.NewTXOpts(), TestEnv.DebugContractAddress, amount, []byte{})
	require.NoError(t, txErr, "transaction should have reverted")
	_, decodeErr := c.Decode(tx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [99 101]", decodeErr.Error(), "expected error message to contain the reverted error type and values")
}

func TestTraceOldPragmaNoRevertReason(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TracingLevel = seth.TracingLevel_Reverted
	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	// this is old Link contract used on Ethereum Mainnet that's in pragma 0.4
	linkAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	data, err := c.DeployContract(c.NewTXOpts(), "LinkTokenInterface", *linkAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy contract")

	instance, err := link_token_interface.NewLinkToken(data.Address, c.Client)
	require.NoError(t, err, "failed to create contract instance")

	amount := big.NewInt(0)
	tx, txErr := instance.TransferAndCall(c.NewTXOpts(), TestEnv.DebugContractAddress, amount, []byte{})
	require.NoError(t, txErr, "transaction should have reverted")
	_, decodeErr := c.Decode(tx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "execution reverted", decodeErr.Error(), "expected error message to contain the reverted error type and values")
}

func TestTraceeRevertReasonNonRootSender(t *testing.T) {
	cBeta := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, cBeta)

	cfg := cBeta.Cfg
	one := int64(1)
	cfg.EphemeralAddrs = &one
	cfg.TracingLevel = seth.TracingLevel_Reverted
	cfg.TraceOutputs = []string{seth.TraceOutput_Console}

	c, err := seth.NewClientWithConfig(cfg)
	require.NoError(t, err, "failed to create client")

	x := big.NewInt(1001)
	y := big.NewInt(2)
	tx, txErr := TestEnv.DebugContract.CallRevertFunctionInSubContract(c.NewTXKeyOpts(1), x, y)
	require.NoError(t, txErr, "transaction should have reverted")
	_, decodeErr := c.Decode(tx, txErr)
	require.Error(t, decodeErr, "transaction should have reverted")
	require.Equal(t, "error type: CustomErr, error values: [1001 2]", decodeErr.Error(), "expected error message to contain the reverted error type and values")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transaction")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[1].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "11b3c478",
			CallType:  "CALL",
			Method:    "callRevertFunctionInSubContract(uint256,uint256)",
			Input:     map[string]interface{}{"x": x, "y": y},
			Output:    map[string]interface{}{},
			Error:     "execution reverted",
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash().Hex())[0], "decoded call does not match")
}

func TestTraceContractTracingClientInitialisesTracerIfTracingIsEnabled(t *testing.T) {
	cfg, err := test_utils.CopyConfig(TestEnv.Client.Cfg)
	require.NoError(t, err, "failed to copy config")

	as, err := seth.NewContractStore(filepath.Join(cfg.ConfigDir, cfg.ABIDir), filepath.Join(cfg.ConfigDir, cfg.BINDir), nil)
	require.NoError(t, err, "failed to create contract store")

	nm, err := seth.NewNonceManager(cfg, TestEnv.Client.Addresses, TestEnv.Client.PrivateKeys)
	require.NoError(t, err, "failed to create nonce manager")

	cfg.TracingLevel = seth.TracingLevel_All
	cfg.TraceOutputs = []string{seth.TraceOutput_Console}
	cfg.Network.TxnTimeout = seth.MustMakeDuration(time.Duration(5 * time.Second))

	c, err := seth.NewClientRaw(
		cfg,
		TestEnv.Client.Addresses,
		TestEnv.Client.PrivateKeys,
		seth.WithContractStore(as),
		seth.WithNonceManager(nm),
	)
	require.NoError(t, err, "failed to create client")
	SkipAnvil(t, c)

	x := big.NewInt(1001)
	y := big.NewInt(2)
	z := big.NewInt(26)
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitInts(c.NewTXOpts(), x, y, z))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "9e099652",
			CallType:  "CALL",
			Method:    "emitInts(int256,int128,uint256)",
			Input:     map[string]interface{}{"first": x, "second": y, "third": z},
			Output:    map[string]interface{}{"0": x, "1": y, "2": z},
		},
	}

	removeGasDataFromDecodedCalls(c.Tracer.GetAllDecodedCalls())
	require.EqualValues(t, expectedCall, c.Tracer.GetDecodedCalls(tx.Hash)[0], "decoded call does not match")
}

func TestTraceContractTracingSaveToJson(t *testing.T) {
	cfg, err := test_utils.CopyConfig(TestEnv.Client.Cfg)
	require.NoError(t, err, "failed to copy config")

	as, err := seth.NewContractStore(filepath.Join(cfg.ConfigDir, cfg.ABIDir), filepath.Join(cfg.ConfigDir, cfg.BINDir), nil)
	require.NoError(t, err, "failed to create contract store")

	nm, err := seth.NewNonceManager(cfg, TestEnv.Client.Addresses, TestEnv.Client.PrivateKeys)
	require.NoError(t, err, "failed to create nonce manager")

	cfg.TracingLevel = seth.TracingLevel_All
	cfg.TraceOutputs = []string{seth.TraceOutput_JSON}
	cfg.Network.TxnTimeout = seth.MustMakeDuration(time.Duration(5 * time.Second))

	c, err := seth.NewClientRaw(
		cfg,
		TestEnv.Client.Addresses,
		TestEnv.Client.PrivateKeys,
		seth.WithContractStore(as),
		seth.WithNonceManager(nm),
	)
	require.NoError(t, err, "failed to create client")
	SkipAnvil(t, c)

	x := big.NewInt(1001)
	y := big.NewInt(2)
	z := big.NewInt(26)
	tx, txErr := c.Decode(TestEnv.DebugContract.EmitInts(c.NewTXOpts(), x, y, z))
	require.NoError(t, txErr, "failed to send transaction")
	require.Equal(t, 1, len(c.Tracer.GetAllDecodedCalls()), "expected 1 decoded transacton")
	require.NotNil(t, c.Tracer.GetDecodedCalls(tx.Hash), "expected decoded calls to contain the transaction hash")

	fileName := filepath.Join(c.Cfg.ArtifactsDir, "traces", fmt.Sprintf("%s.json", tx.Hash))
	t.Cleanup(func() {
		_ = os.Remove(fileName)
	})

	expectedCall := &seth.DecodedCall{
		FromAddress: strings.ToLower(c.Addresses[0].Hex()),
		ToAddress:   strings.ToLower(TestEnv.DebugContractAddress.Hex()),
		From:        "you",
		To:          "NetworkDebugContract",
		CommonData: seth.CommonData{
			Signature: "9e099652",
			CallType:  "CALL",
			Method:    "emitInts(int256,int128,uint256)",
			Input:     map[string]interface{}{"first": 1001.0, "second": 2.0, "third": 26.0},
			Output:    map[string]interface{}{"0": 1001.0, "1": 2.0, "2": 26.0},
		},
		Comment: "",
	}

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	require.NoError(t, err, "expected trace file to exist")

	var readCall []seth.DecodedCall

	defer func() { _ = f.Close() }()
	b, _ := io.ReadAll(f)
	err = json.Unmarshal(b, &readCall)
	require.NoError(t, err, "failed to unmarshal trace file")

	removeGasDataFromDecodedCalls(map[string][]*seth.DecodedCall{tx.Hash: {&readCall[0]}})

	require.Equal(t, 1, len(readCall), "expected 1 decoded transaction")
	require.EqualValues(t, expectedCall, &readCall[0], "decoded call does not match one read from file")
}

func TestTraceContractTracingSaveToDot(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TraceOutputs = []string{seth.TraceOutput_DOT}
	c.Cfg.TracingLevel = seth.TracingLevel_All

	linkTokenAbi, err := link_token.LinkTokenMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	linkDeploymentData, err := c.DeployContract(c.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	if err != nil {
		panic(err)
	}
	linkToken, err := link_token.NewLinkToken(linkDeploymentData.Address, c.Client)
	if err != nil {
		panic(err)
	}

	_, err = c.Decode(linkToken.GrantMintRole(c.NewTXOpts(), c.MustGetRootKeyAddress()))
	if err != nil {
		fmt.Println("failed to grant mint LINK role")
		os.Exit(1)
	}

	_, err = c.Decode(linkToken.Mint(c.NewTXOpts(), c.MustGetRootKeyAddress(), big.NewInt(1000000000000000000)))
	if err != nil {
		fmt.Println("failed to mint LINK")
		os.Exit(1)
	}

	debugAbi, err := abi.JSON(strings.NewReader(network_debug_contract.NetworkDebugContractMetaData.ABI))
	if err != nil {
		fmt.Println("failed to get debug contract ABI")
		os.Exit(1)
	}

	var x int64 = 6
	var y int64 = 5

	req, err := debugAbi.Pack(
		"traceWithValidate",
		big.NewInt(x),
		big.NewInt(y),
	)

	if err != nil {
		fmt.Println("failed to pack arguments")
		os.Exit(1)
	}

	amount := big.NewInt(10)
	decodedTx, decodeErr := c.Decode(linkToken.TransferAndCall(c.NewTXOpts(), TestEnv.DebugContractAddress, amount, req))
	require.NoError(t, decodeErr, "transaction should not have reverted")

	fileName := filepath.Join(c.Cfg.ArtifactsDir, "dot_graphs", fmt.Sprintf("%s.dot", decodedTx.Hash))
	t.Cleanup(func() {
		_ = os.Remove(fileName)
	})

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	require.NoError(t, err, "expected trace file to exist")
	require.NotNil(t, f, "expected file to exist")
	s, err := f.Stat()
	require.NoError(t, err, "expected file to exist")
	require.Greater(t, s.Size(), int64(0), "expected file to have content")
}

func TestTraceVariousCallTypesAndNestingLevels(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	c.Cfg.TraceOutputs = []string{seth.TraceOutput_Console}
	c.Cfg.TracingLevel = seth.TracingLevel_All

	linkTokenAbi, err := link_token.LinkTokenMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	linkDeploymentData, err := c.DeployContract(c.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	if err != nil {
		panic(err)
	}
	linkToken, err := link_token.NewLinkToken(linkDeploymentData.Address, c.Client)
	if err != nil {
		panic(err)
	}

	_, err = c.Decode(linkToken.GrantMintRole(c.NewTXOpts(), c.MustGetRootKeyAddress()))
	if err != nil {
		fmt.Println("failed to grant mint LINK role")
		os.Exit(1)
	}

	_, err = c.Decode(linkToken.Mint(c.NewTXOpts(), c.MustGetRootKeyAddress(), big.NewInt(1000000000000000000)))
	if err != nil {
		fmt.Println("failed to mint LINK")
		os.Exit(1)
	}

	debugAbi, err := abi.JSON(strings.NewReader(network_debug_contract.NetworkDebugContractMetaData.ABI))
	if err != nil {
		fmt.Println("failed to get debug contract ABI")
		os.Exit(1)
	}

	var x int64 = 6
	var y int64 = 5

	req, err := debugAbi.Pack(
		"traceWithValidate",
		big.NewInt(x),
		big.NewInt(y),
	)

	if err != nil {
		fmt.Println("failed to pack arguments")
		os.Exit(1)
	}

	amount := big.NewInt(10)
	decodedTx, decodeErr := c.Decode(linkToken.TransferAndCall(c.NewTXOpts(), TestEnv.DebugContractAddress, amount, req))
	require.NoError(t, decodeErr, "transaction should not have reverted")
	require.Equal(t, 3, len(c.Tracer.GetAllDecodedCalls()), "expected 3 decoded transaction")
	require.Equal(t, 9, len(c.Tracer.GetDecodedCalls(decodedTx.Hash)), "expected 9 decoded transaction for tx: "+decodedTx.Hash)

	require.Equal(t, "CALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[0].CallType, "expected call type to be CALL")
	require.Equal(t, 0, c.Tracer.GetDecodedCalls(decodedTx.Hash)[0].NestingLevel, "expected nesting level to be 0")

	require.Equal(t, "CALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[1].CallType, "expected call type to be CALL")
	require.Equal(t, 1, c.Tracer.GetDecodedCalls(decodedTx.Hash)[1].NestingLevel, "expected nesting level to be 1")

	require.Equal(t, "DELEGATECALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[2].CallType, "expected call type to be DELEGATECALL")
	require.Equal(t, 2, c.Tracer.GetDecodedCalls(decodedTx.Hash)[2].NestingLevel, "expected nesting level to be 2")

	require.Equal(t, "CALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[3].CallType, "expected call type to be CALL")
	require.Equal(t, 3, c.Tracer.GetDecodedCalls(decodedTx.Hash)[3].NestingLevel, "expected nesting level to be 3")

	require.Equal(t, "STATICCALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[4].CallType, "expected call type to be STATICCALL")
	require.Equal(t, 2, c.Tracer.GetDecodedCalls(decodedTx.Hash)[4].NestingLevel, "expected nesting level to be 2")

	require.Equal(t, "STATICCALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[5].CallType, "expected call type to be STATICCALL")
	require.Equal(t, 3, c.Tracer.GetDecodedCalls(decodedTx.Hash)[5].NestingLevel, "expected nesting level to be 3")

	require.Equal(t, "CALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[6].CallType, "expected call type to be CALL")
	require.Equal(t, 2, c.Tracer.GetDecodedCalls(decodedTx.Hash)[6].NestingLevel, "expected nesting level to be 2")

	require.Equal(t, "CALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[7].CallType, "expected call type to be CALL")
	require.Equal(t, 3, c.Tracer.GetDecodedCalls(decodedTx.Hash)[7].NestingLevel, "expected nesting level to be 3")

	require.Equal(t, "CALL", c.Tracer.GetDecodedCalls(decodedTx.Hash)[8].CallType, "expected call type to be CALL")
	require.Equal(t, 4, c.Tracer.GetDecodedCalls(decodedTx.Hash)[8].NestingLevel, "expected nesting level to be 4")
}

func TestNestedEvents(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	tx, txErr := TestEnv.DebugContract.TraceNestedEvents(c.NewTXOpts())
	require.NoError(t, txErr, "transaction should have succeeded")
	decoded, decodeErr := c.Decode(tx, txErr)
	require.NoError(t, decodeErr, "transaction should have succeeded")

	expectedLogs := []seth.DecodedCommonLog{
		{
			Signature: "UniqueSubDebugEvent()",
			Address:   TestEnv.DebugSubContractAddress,
			EventData: map[string]interface{}{},
			Topics:    []string{"0xe0b03c5e88196d907268b0babc690e041bdc7fcc1abf4bbf1e363e28c17e6b9b"},
		},
		{
			Signature: "UniqueDebugEvent()",
			Address:   TestEnv.DebugContractAddress,
			EventData: map[string]interface{}{},
			Topics:    []string{"0xa0f7c7c1fff15178b5db3e56860767f0889c56b591bd2d9ba3121b491347d74c"},
		},
	}

	require.Equal(t, 2, len(decoded.Events), "expected 2 events")
	var actualEvents []seth.DecodedCommonLog
	for _, event := range decoded.Events {
		actualEvents = append(actualEvents, event.DecodedCommonLog)
	}

	require.EqualValues(t, expectedLogs, actualEvents, "decoded events do not match")
}

func TestSameEventTwoABIs(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	contractAbi, err := TestContractOne.UniqueEventOneMetaData.GetAbi()
	require.NoError(t, err, "failed to get contract ABI")
	oneData, err := c.DeployContract(c.NewTXOpts(), "TestContractOne", *contractAbi, common.FromHex(TestContractOne.UniqueEventOneMetaData.Bin))
	require.NoError(t, err, "failed to deploy contract")

	contractAbi, err = TestContractTwo.UniqueEventTwoMetaData.GetAbi()
	require.NoError(t, err, "failed to get contract ABI")
	_, err = c.DeployContract(c.NewTXOpts(), "TestContractTwo", *contractAbi, common.FromHex(TestContractTwo.UniqueEventTwoMetaData.Bin))
	require.NoError(t, err, "failed to deploy contract")

	oneInstance, err := TestContractOne.NewUniqueEventOne(oneData.Address, c.Client)
	require.NoError(t, err, "failed to create contract instance")
	decoded, txErr := c.Decode(oneInstance.ExecuteFirstOperation(c.NewTXOpts(), big.NewInt(1), big.NewInt(2)))
	require.NoError(t, txErr, "transaction should have succeeded")

	expectedLogs := []seth.DecodedCommonLog{
		{
			Signature: "NonUniqueEvent(int256,int256)",
			Address:   oneData.Address,
			EventData: map[string]interface{}{
				"a": big.NewInt(1),
				"b": big.NewInt(2),
			},
			Topics: []string{"0x192aedde7837c0cbfb2275e082ba2391de36cf5a893681e9dac2cced6947614e", "0x0000000000000000000000000000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000000000000000000000000000002"},
		},
	}

	require.Equal(t, 1, len(decoded.Events), "expected 1 event")
	var actualEvents []seth.DecodedCommonLog
	for _, event := range decoded.Events {
		actualEvents = append(actualEvents, event.DecodedCommonLog)
	}

	require.EqualValues(t, expectedLogs, actualEvents, "decoded events do not match")
}

func removeGasDataFromDecodedCalls(decodedCall map[string][]*seth.DecodedCall) {
	for _, decodedCalls := range decodedCall {
		for _, call := range decodedCalls {
			call.GasUsed = 0
			call.GasLimit = 0
		}
	}
}
