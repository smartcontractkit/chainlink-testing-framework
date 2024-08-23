package seth_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"
)

func TestSmokeDebugReverts(t *testing.T) {
	c := newClient(t)

	type test struct {
		name   string
		method string
		output map[string]string
	}

	tests := []test{
		{
			name:   "revert with require",
			method: "alwaysRevertsRequire",
			output: map[string]string{
				seth.GETH:  "execution reverted: always revert error",
				seth.ANVIL: "execution reverted: revert: always revert error",
			},
		},
		{
			name:   "revert with assert(panic)",
			method: "alwaysRevertsAssert",
			output: map[string]string{
				seth.GETH:  "execution reverted: assert(false)",
				seth.ANVIL: "execution reverted: panic: assertion failed (0x01)",
			},
		},
		{
			name:   "revert with a custom err",
			method: "alwaysRevertsCustomError",
			output: map[string]string{
				seth.GETH: "error type: CustomErr, error values: [12 21]",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "revert with a custom err" && (c.Cfg.Network.Name == "Mumbai" || c.Cfg.Network.Name == "Fuji") {
				t.Skip("typed errors are not supported\nnodes payload of rpc.DataError is empty and tx fails on send, not on execution")
			}
			_, err := c.Decode(TestEnv.DebugContractRaw.Transact(c.NewTXOpts(), tc.method))
			require.Error(t, err)
			var expectedOutput = tc.output[seth.GETH]
			if c.Cfg.Network.Name != seth.GETH {
				if eo, ok := tc.output[c.Cfg.Network.Name]; ok {
					expectedOutput = eo
				}
			}
			require.Equal(t, expectedOutput, err.Error())
		})
	}
}

func TestSmokeDebugData(t *testing.T) {
	c := newClient(t)
	c.Cfg.TracingLevel = seth.TracingLevel_All

	type test struct {
		name   string
		method string
		params []interface{}
		output seth.DecodedTransaction
		write  bool
	}

	tests := []test{
		{
			name:   "test named inputs/outputs",
			method: "emitNamedInputsOutputs",
			params: []interface{}{big.NewInt(1337), "test"},
			write:  true,
			output: seth.DecodedTransaction{
				CommonData: seth.CommonData{
					Input: map[string]interface{}{
						"inputVal1": big.NewInt(1337),
						"inputVal2": "test",
					},
				},
			},
		},
		// TODO: https://docs.soliditylang.org/en/v0.8.19/control-structures.html read and figure out if
		// decoding anynymous + named is heavily used and needed, usually people name params and omit output names
		{
			name:   "test anonymous inputs/outputs",
			method: "emitInputsOutputs",
			params: []interface{}{big.NewInt(1337), "test"},
			write:  true,
			output: seth.DecodedTransaction{
				CommonData: seth.CommonData{
					Input: map[string]interface{}{
						"inputVal1": big.NewInt(1337),
						"inputVal2": "test",
					},
				},
			},
		},
		{
			name:   "test one log no index",
			method: "emitNoIndexEvent",
			write:  true,
			output: seth.DecodedTransaction{
				Events: []seth.DecodedTransactionLog{
					{
						DecodedCommonLog: seth.DecodedCommonLog{
							EventData: map[string]interface{}{
								"sender": c.Addresses[0],
							},
						},
					},
				},
			},
		},
		{
			name:   "test one log index",
			method: "emitOneIndexEvent",
			write:  true,
			output: seth.DecodedTransaction{
				Events: []seth.DecodedTransactionLog{
					{
						DecodedCommonLog: seth.DecodedCommonLog{
							EventData: map[string]interface{}{
								"a": big.NewInt(83),
							},
						},
					},
				},
			},
		},
		{
			name:   "test two log index",
			method: "emitTwoIndexEvent",
			write:  true,
			output: seth.DecodedTransaction{
				Events: []seth.DecodedTransactionLog{
					{
						DecodedCommonLog: seth.DecodedCommonLog{
							EventData: map[string]interface{}{
								"roundId":   big.NewInt(1),
								"startedBy": c.Addresses[0],
							},
						},
					},
				},
			},
		},
		{
			name:   "test three log index",
			method: "emitThreeIndexEvent",
			write:  true,
			output: seth.DecodedTransaction{
				Events: []seth.DecodedTransactionLog{
					{
						DecodedCommonLog: seth.DecodedCommonLog{
							EventData: map[string]interface{}{
								"roundId":   big.NewInt(1),
								"startedAt": big.NewInt(3),
								"startedBy": c.Addresses[0],
							},
						},
					},
				},
			},
		},
		{
			name:   "test log no index string",
			method: "emitNoIndexEventString",
			write:  true,
			output: seth.DecodedTransaction{
				Events: []seth.DecodedTransactionLog{
					{
						DecodedCommonLog: seth.DecodedCommonLog{
							EventData: map[string]interface{}{
								"str": "myString",
							},
						},
					},
				},
			},
		},
		// emitNoIndexStructEvent
		{
			name:   "test log struct",
			method: "emitNoIndexStructEvent",
			write:  true,
			output: seth.DecodedTransaction{
				Events: []seth.DecodedTransactionLog{
					{
						DecodedCommonLog: seth.DecodedCommonLog{
							EventData: map[string]interface{}{
								"a": struct {
									Name       string   `json:"name"`
									Balance    uint64   `json:"balance"`
									DailyLimit *big.Int `json:"dailyLimit"`
								}{
									Name:       "John",
									Balance:    5,
									DailyLimit: big.NewInt(10),
								},
							},
						},
					},
				},
			},
		},
		// TODO: another case - figure out if indexed strings are used by anyone in events
		// https://ethereum.stackexchange.com/questions/6840/indexed-event-with-string-not-getting-logged
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.write {
				dtx, err := c.Decode(
					TestEnv.DebugContractRaw.Transact(c.NewTXOpts(), tc.method, tc.params...),
				)
				require.NoError(t, err)
				require.Equal(t, dtx.Input, tc.output.Input)
				require.Equal(t, dtx.Output, tc.output.Output)
				for i, e := range tc.output.Events {
					require.NotNil(t, dtx.Events[i])
					require.Equal(t, dtx.Events[i].EventData, e.EventData)
				}
			}
		})
	}
}
