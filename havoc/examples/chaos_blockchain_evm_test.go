package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestBlockchainChaos(t *testing.T) {
	srcChain := os.Getenv("CTF_CHAOS_SRC_CHAIN_RPC_HTTP_URL")
	require.NotEmpty(t, srcChain, "source chain RPC must be set")
	dstChain := os.Getenv("CTF_CHAOS_DST_CHAIN_RPC_HTTP_URL")
	require.NotEmpty(t, dstChain, "destination chain RPC must be set")

	recoveryIntervalDuration := 120 * time.Second

	testCases := []struct {
		name       string
		chainURL   string
		reorgDepth int
	}{
		{
			name:       "Reorg src with depth: 1",
			chainURL:   srcChain,
			reorgDepth: 1,
		},
		{
			name:       "Reorg dst with depth: 1",
			chainURL:   dstChain,
			reorgDepth: 1,
		},
		{
			name:       "Reorg src with depth: 5",
			chainURL:   srcChain,
			reorgDepth: 5,
		},
		{
			name:       "Reorg dst with depth: 5",
			chainURL:   dstChain,
			reorgDepth: 5,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			r := rpc.New(tc.chainURL, nil)
			err := r.GethSetHead(tc.reorgDepth)
			require.NoError(t, err)
			t.Logf("Awaiting chaos recovery: %s", tc.name)
			time.Sleep(recoveryIntervalDuration)

			// Validate Chainlink product here, use load test assertions
		})
	}
}
