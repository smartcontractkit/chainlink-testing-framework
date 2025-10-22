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

	suggestions, err := estimator.Stats(t.Context(), bn, 25)
	require.NoError(t, err, "Gas estimator should not err")
	require.NotNil(t, suggestions.BaseFeePerc, "Suggested base fee percentiles should not be nil")
	require.NotNil(t, suggestions.TipCapPerc, "Suggested tip cap percentiles should not be nil")
	require.NotNil(t, suggestions.LastBaseFee, "Last base fee should not be nil")
	require.NotNil(t, suggestions.SuggestedGasTipCap, "Suggested gas tip cap should not be nil")
	require.NotNil(t, suggestions.SuggestedGasPrice, "Suggested gas price should not be nil")

	require.Greater(t, suggestions.BaseFeePerc.Perc25, float64(0), "Suggested 25th percentile gas price should be greater than 0")
	require.GreaterOrEqual(t, suggestions.BaseFeePerc.Perc50, suggestions.BaseFeePerc.Perc25, "Suggested 50th percentile gas price should be greater than or equal to 25th percentile")
	require.GreaterOrEqual(t, suggestions.BaseFeePerc.Perc75, suggestions.BaseFeePerc.Perc50, "Suggested 75th percentile gas price should be greater than or equal to 50th percentile")
	require.GreaterOrEqual(t, suggestions.BaseFeePerc.Perc99, suggestions.BaseFeePerc.Perc75, "Suggested 99th percentile gas price should be greater than or equal to 75th percentile")
	require.GreaterOrEqual(t, suggestions.BaseFeePerc.Max, suggestions.BaseFeePerc.Perc99, "Suggested max gas price should be greater than or equal to 99th percentile")

	require.GreaterOrEqual(t, suggestions.TipCapPerc.Perc25, float64(0), "Suggested 25th percentile tip cap should be greater than or equal to 0")
	require.GreaterOrEqual(t, suggestions.TipCapPerc.Perc50, suggestions.TipCapPerc.Perc25, "Suggested 50th percentile tip cap should be greater than or equal to 25th percentile")
	require.GreaterOrEqual(t, suggestions.TipCapPerc.Perc75, suggestions.TipCapPerc.Perc50, "Suggested 75th percentile tip cap should be greater than or equal to 50th percentile")
	require.GreaterOrEqual(t, suggestions.TipCapPerc.Perc99, suggestions.TipCapPerc.Perc75, "Suggested 99th percentile tip cap should be greater than or equal to 75th percentile")
	require.GreaterOrEqual(t, suggestions.TipCapPerc.Max, suggestions.TipCapPerc.Perc99, "Suggested max tip cap should be greater than or equal to 99th percentile")
}
