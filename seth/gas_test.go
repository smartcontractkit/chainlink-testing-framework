package seth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

func TestGasEstimator(t *testing.T) {
	c := newClient(t)
	bn, err := c.Client.BlockNumber(context.Background())
	require.NoError(t, err, "BlockNumber should not error")
	for range 10 {
		_, err := c.DeployContractFromContractStore(c.NewTXOpts(), "NetworkDebugSubContract")
		require.NoError(t, err, "Deploying contract should not error")
	}
	estimator := seth.NewGasEstimator(c)

	suggestions, err := estimator.Stats(bn, 25)
	require.NoError(t, err, "Gas estimator should not err")
	require.NotNil(t, suggestions.GasPrice, "Suggested gas price should not be nil")
	require.NotNil(t, suggestions.TipCap, "Suggested tip cap should not be nil")

	require.Greater(t, suggestions.GasPrice.Perc25, float64(0), "Suggested 25th percentile gas price should be greater than 0")
	require.GreaterOrEqual(t, suggestions.GasPrice.Perc50, suggestions.GasPrice.Perc25, "Suggested 50th percentile gas price should be greater than or equal to 25th percentile")
	require.GreaterOrEqual(t, suggestions.GasPrice.Perc75, suggestions.GasPrice.Perc50, "Suggested 75th percentile gas price should be greater than or equal to 50th percentile")
	require.GreaterOrEqual(t, suggestions.GasPrice.Perc99, suggestions.GasPrice.Perc75, "Suggested 99th percentile gas price should be greater than or equal to 75th percentile")
	require.GreaterOrEqual(t, suggestions.GasPrice.Max, suggestions.GasPrice.Perc99, "Suggested max gas price should be greater than or equal to 99th percentile")

	require.GreaterOrEqual(t, suggestions.TipCap.Perc25, float64(0), "Suggested 25th percentile tip cap should be greater than or equal to 0")
	require.GreaterOrEqual(t, suggestions.TipCap.Perc50, suggestions.TipCap.Perc25, "Suggested 50th percentile tip cap should be greater than or equal to 25th percentile")
	require.GreaterOrEqual(t, suggestions.TipCap.Perc75, suggestions.TipCap.Perc50, "Suggested 75th percentile tip cap should be greater than or equal to 50th percentile")
	require.GreaterOrEqual(t, suggestions.TipCap.Perc99, suggestions.TipCap.Perc75, "Suggested 99th percentile tip cap should be greater than or equal to 75th percentile")
	require.GreaterOrEqual(t, suggestions.TipCap.Max, suggestions.TipCap.Perc99, "Suggested max tip cap should be greater than or equal to 99th percentile")
}
