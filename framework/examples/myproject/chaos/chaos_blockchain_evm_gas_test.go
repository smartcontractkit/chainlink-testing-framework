package chaos

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
)

type CfgGas struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestBlockchainGasChaos(t *testing.T) {
	in, err := f.Load[CfgGas](t)
	require.NoError(t, err)

	// Can replace deployments with CRIB here

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	c, err := clclient.New(out.CLNodes)
	require.NoError(t, err)

	// !! This value must match anvil block speed to set gas values for every block !!
	blockEvery := 1 * time.Second
	waitBetweenTests := 1 * time.Minute

	gasControlFunc := func(t *testing.T, r *rpc.RPCClient) {
		startGasPrice := big.NewInt(2e9)
		// ramp
		for i := 0; i < 10; i++ {
			err := r.PrintBlockBaseFee()
			require.NoError(t, err)
			t.Logf("Setting block base fee: %d", startGasPrice)
			err = r.AnvilSetNextBlockBaseFeePerGas(startGasPrice)
			require.NoError(t, err)
			startGasPrice = startGasPrice.Add(startGasPrice, big.NewInt(1e9))
			time.Sleep(blockEvery)
		}
		// hold
		for i := 0; i < 10; i++ {
			err := r.PrintBlockBaseFee()
			require.NoError(t, err)
			time.Sleep(blockEvery)
			t.Logf("Setting block base fee: %d", startGasPrice)
			err = r.AnvilSetNextBlockBaseFeePerGas(startGasPrice)
			require.NoError(t, err)
		}
		// release
		for i := 0; i < 10; i++ {
			err := r.PrintBlockBaseFee()
			require.NoError(t, err)
			time.Sleep(blockEvery)
		}
	}

	testCases := []struct {
		name             string
		chainURL         string
		increase         *big.Int
		waitBetweenTests time.Duration
		gasFunc          func(t *testing.T, r *rpc.RPCClient)
		validate         func(t *testing.T, c []*clclient.ChainlinkClient)
	}{
		{
			name:             "Slow and low",
			chainURL:         bc.Nodes[0].ExternalHTTPUrl,
			waitBetweenTests: 30 * time.Second,
			increase:         big.NewInt(1e9),
			gasFunc:          gasControlFunc,
			validate: func(t *testing.T, c []*clclient.ChainlinkClient) {
				// add more clients and validate
			},
		},
		{
			name:             "Fast and degen",
			chainURL:         bc.Nodes[0].ExternalHTTPUrl,
			waitBetweenTests: 30 * time.Second,
			increase:         big.NewInt(5e9),
			gasFunc:          gasControlFunc,
			validate: func(t *testing.T, c []*clclient.ChainlinkClient) {
				// add more clients and validate
			},
		},
	}

	// Start WASP load test here, apply average load profile that you expect in production!
	// Configure timeouts and validate all the test cases until the test ends

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			r := rpc.New(tc.chainURL, nil)
			tc.gasFunc(t, r)
			tc.validate(t, c)
			time.Sleep(waitBetweenTests)
		})
	}
}
