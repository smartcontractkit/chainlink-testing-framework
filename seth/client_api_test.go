package seth_test

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/seth/test_utils"
)

func TestAPI(t *testing.T) {
	c := newClientWithEphemeralAddresses(t)

	t.Cleanup(func() {
		err := c.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(c, c.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
	})

	type test struct {
		name            string
		EIP1559Enabled  bool
		transactionOpts []seth.TransactOpt
		callOpts        []seth.CallOpt
	}

	bn, err := c.Client.BlockNumber(context.Background())
	require.NoError(t, err)
	weiValue := big.NewInt(1)
	overriddenGasPrice := big.NewInt(c.Cfg.Network.GasPrice + 1)
	overriddenGasFeeCap := big.NewInt(c.Cfg.Network.GasFeeCap + 1)
	overriddenGasTipCap := big.NewInt(c.Cfg.Network.GasTipCap + 1)
	overriddenGasLimit := uint64(c.Cfg.Network.GasLimit) + 1

	tests := []test{
		{
			name: "default tx gas opts from cfg",
		},
		{
			name: "custom legacy tx opts override",
			transactionOpts: []seth.TransactOpt{
				seth.WithGasPrice(overriddenGasPrice),
				seth.WithGasLimit(overriddenGasLimit),
			},
		},
		{
			name: "custom EIP-1559 tx opts override",
			transactionOpts: []seth.TransactOpt{
				seth.WithGasFeeCap(overriddenGasFeeCap),
				seth.WithGasTipCap(overriddenGasTipCap),
			},
			EIP1559Enabled: true,
		},
		{
			name: "with value override",
			transactionOpts: []seth.TransactOpt{
				seth.WithValue(weiValue),
			},
		},
		{
			name: "custom call opts",
			callOpts: []seth.CallOpt{
				seth.WithPending(true),
				seth.WithBlockNumber(bn),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var dtx *seth.DecodedTransaction
			c.Cfg.Network.EIP1559DynamicFees = tc.EIP1559Enabled
			if tc.name == "with value override" {
				dtx, err = c.Decode(
					TestEnv.DebugContract.Pay(c.NewTXOpts(tc.transactionOpts...)),
				)
				require.NoError(t, err)
				require.Equal(t, weiValue, dtx.Transaction.Value())
			} else {
				dtx, err = c.Decode(
					TestEnv.DebugContract.Set(c.NewTXOpts(tc.transactionOpts...), big.NewInt(1)),
				)
			}
			require.NoError(t, err)
			require.NotEmpty(t, dtx.Transaction)
			require.NotEmpty(t, dtx.Receipt)
			val, err := TestEnv.DebugContract.Get(c.NewCallOpts(tc.callOpts...))
			require.NoError(t, err)
			require.Equal(t, big.NewInt(1), val)
		})
	}
}

func TestAPINonces(t *testing.T) {
	c := newClientWithEphemeralAddresses(t)

	t.Cleanup(func() {
		err := c.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(c, c.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
	})

	type test struct {
		name            string
		EIP1559Enabled  bool
		transactionOpts []seth.TransactOpt
		callOpts        []seth.CallOpt
	}

	pnonce, err := c.Client.PendingNonceAt(context.Background(), c.Addresses[0])
	require.NoError(t, err)

	tests := []test{
		{
			name: "with nonce override",
			transactionOpts: []seth.TransactOpt{
				//nolint
				seth.WithNonce(big.NewInt(int64(pnonce))),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.Cfg.Network.EIP1559DynamicFees = tc.EIP1559Enabled
			_, err := c.Decode(
				TestEnv.DebugContract.Set(c.NewTXOpts(tc.transactionOpts...), big.NewInt(1)),
			)
			require.NoError(t, err)
			val, err := TestEnv.DebugContract.Get(c.NewCallOpts(tc.callOpts...))
			require.NoError(t, err)
			require.Equal(t, big.NewInt(1), val)
		})
	}
}

func TestAPISeqErrors(t *testing.T) {
	c := newClientWithEphemeralAddresses(t)

	t.Cleanup(func() {
		err := c.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(c, c.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
	})

	type test struct {
		name string
	}

	tests := []test{
		{
			name: "raise previous call error first",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.Errors = append(c.Errors, errors.New("previous call error"))
			_, err := c.Decode(
				TestEnv.DebugContract.Set(c.NewTXOpts(), big.NewInt(1)),
			)
			require.Error(t, err)
		})
	}
}

func TestAPIConfig(t *testing.T) {
	cfg, err := seth.ReadConfig()
	require.NoError(t, err)
	addrs, pkeys, err := cfg.ParseKeys()
	require.NoError(t, err)
	c, err := seth.NewClientRaw(cfg, addrs, pkeys)
	require.NoError(t, err)

	type test struct {
		name string
	}

	tests := []test{
		{
			name: "can run without ABI",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.Decode(
				TestEnv.DebugContract.Set(c.NewTXOpts(), big.NewInt(1)),
			)
			require.NoError(t, err)
			val, err := TestEnv.DebugContract.Get(c.NewCallOpts())
			require.NoError(t, err)
			require.Equal(t, big.NewInt(1), val)
		})
	}
}

func TestAPIKeys(t *testing.T) {
	type test struct {
		name string
	}

	keyCount := 60
	c := test_utils.NewClientWithAddresses(t, keyCount, nil)

	t.Cleanup(func() {
		err := c.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(c, c.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
	})

	tests := []test{
		{
			name: "multiple separate keys used",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			for i := 1; i < 61; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					_, err := c.Decode(
						TestEnv.DebugContract.AddCounter(c.NewTXKeyOpts(i), big.NewInt(0), big.NewInt(1)),
					)
					require.NoError(t, err)
				}(i)
			}
			wg.Wait()
			for i := 1; i < 61; i++ {
				val, err := TestEnv.DebugContract.GetCounter(c.NewCallOpts(), big.NewInt(0))
				require.NoError(t, err)
				require.Equal(t, big.NewInt(60), val)
			}
		})
	}
}

func TestManualAPIReconnect(t *testing.T) {
	c := newClientWithEphemeralAddresses(t)

	t.Cleanup(func() {
		err := c.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(c, c.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
	})

	type test struct {
		name            string
		transactionOpts []seth.TransactOpt
	}

	tests := []test{
		{
			name: "can reconnect",
		},
	}

	for _, tc := range tests {
		for i := 0; i < 20; i++ {
			_, err := c.RetryTxAndDecode(func() (*types.Transaction, error) {
				return TestEnv.DebugContract.Set(c.NewTXOpts(tc.transactionOpts...), big.NewInt(1))
			})
			require.NoError(t, err)
			time.Sleep(1 * time.Second)
		}
	}
}
