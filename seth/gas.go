package seth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/montanaflynn/stats"
)

// GasEstimator estimates gas prices
type GasEstimator struct {
	Client              *Client
	BlockGasLimits      []uint64
	TransactionGasPrice []uint64
}

// NewGasEstimator creates a new gas estimator
func NewGasEstimator(c *Client) *GasEstimator {
	return &GasEstimator{Client: c}
}

// Stats calculates gas price and tip cap suggestions based on historical fee data over a specified number of blocks.
// It computes quantiles for base fees and tip caps and provides suggested gas price and tip cap values.
func (m *GasEstimator) Stats(ctx context.Context, blockCount uint64, priorityPerc float64) (GasSuggestions, error) {
	estimations := GasSuggestions{}

	if blockCount == 0 {
		return estimations, fmt.Errorf("block count must be greater than zero for gas estimation. "+
			"Check 'gas_price_estimation_blocks' in your config (seth.toml or ClientBuilder) - current value: %d", blockCount)
	}

	currentBlock, err := m.Client.Client.BlockNumber(ctx)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get current block number: %w\n"+
			"Ensure RPC endpoint is accessible and synced. "+
			"Block history-based gas estimation requires access to recent block data. "+
			"Alternatively, set 'gas_price_estimation_blocks = 0' to disable block-based estimation",
			err)
	}

	if currentBlock == 0 {
		return GasSuggestions{}, fmt.Errorf("current block number is zero, which indicates either:\n" +
			"  1. The network hasn't produced any blocks yet (check if network is running)\n" +
			"  2. RPC node is not synced\n" +
			"  3. Connection to RPC node failed\n" +
			"Block history-based gas estimation is not possible without block history. " +
			"You can set 'gas_price_estimation_blocks = 0' to disable block-based estimation")
	}
	if blockCount >= currentBlock {
		blockCount = max(currentBlock-1, 1) // avoid a case, when we ask for more blocks than exist and when currentBlock = 1
	}

	hist, err := m.Client.Client.FeeHistory(ctx, blockCount, big.NewInt(mustSafeInt64(currentBlock)), []float64{priorityPerc})
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get fee history for %d blocks: %w\n"+
			"Possible causes:\n"+
			"  1. RPC node doesn't support eth_feeHistory\n"+
			"  2. Not enough blocks available (current block: %d)\n"+
			"  3. Network connection issues\n"+
			"Try reducing 'gas_price_estimation_blocks' in config",
			blockCount, err, currentBlock)
	}
	L.Trace().
		Interface("History", hist).
		Msg("Fee history")

	baseFees := make([]float64, 0)
	for _, bf := range hist.BaseFee {
		if bf == nil {
			bf = big.NewInt(0)
		}
		f := new(big.Float).SetInt(bf)
		ff, _ := f.Float64()
		baseFees = append(baseFees, ff)
	}
	gasPercs, err := quantilesFromFloatArray(baseFees)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to calculate gas price quantiles from %d blocks of fee history: %w\n"+
			"This might indicate insufficient or invalid fee data. "+
			"Try reducing 'gas_price_estimation_blocks' in config",
			len(baseFees), err)
	}
	estimations.BaseFeePerc = gasPercs

	L.Trace().
		Interface("Gas percentiles ", gasPercs).
		Msg("Base fees")

	tips := make([]float64, 0)
	for _, bf := range hist.Reward {
		if len(bf) == 0 {
			continue
		}
		if bf[0] == nil {
			bf[0] = big.NewInt(0)
		}
		f := new(big.Float).SetInt(bf[0])
		ff, _ := f.Float64()
		tips = append(tips, ff)
	}
	tipPercs, err := quantilesFromFloatArray(tips)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to calculate tip cap quantiles from %d blocks of fee history: %w\n"+
			"This might indicate insufficient or invalid tip data. "+
			"Try reducing 'gas_price_estimation_blocks' in config",
			len(tips), err)
	}
	estimations.TipCapPerc = tipPercs
	L.Trace().
		Interface("Gas percentiles ", tipPercs).
		Msg("Tip caps")

	suggestedGasPrice, err := m.Client.Client.SuggestGasPrice(ctx)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get suggested gas price from RPC: %w\n"+
			"Possible solutions:\n"+
			"  1. Disable gas estimation and set explicit 'gas_price' in config (gas_price_estimation_enabled = false)\n"+
			"  2. Check RPC node capabilities and accessibility\n"+
			"  3. Verify the network supports gas price queries",
			err)
	}
	estimations.SuggestedGasPrice = suggestedGasPrice

	suggestedGasTipCap, err := m.Client.Client.SuggestGasTipCap(ctx)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get suggested gas tip cap from RPC: %w\n"+
			"Possible solutions:\n"+
			"  1. Disable gas estimation and set explicit 'gas_tip_cap' in config (gas_price_estimation_enabled = false)\n"+
			"  2. Check if network supports EIP-1559\n"+
			"  3. Verify RPC node capabilities",
			err)
	}

	estimations.SuggestedGasTipCap = suggestedGasTipCap

	header, err := m.Client.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get latest block header: %w\n"+
			"Cannot determine current base fee. Check RPC connection",
			err)
	}
	estimations.LastBaseFee = header.BaseFee

	return estimations, nil
}

// GasPercentiles contains gas percentiles
type GasPercentiles struct {
	Max    float64
	Perc99 float64
	Perc75 float64
	Perc50 float64
	Perc25 float64
}

type GasSuggestions struct {
	BaseFeePerc        *GasPercentiles
	TipCapPerc         *GasPercentiles
	LastBaseFee        *big.Int
	SuggestedGasPrice  *big.Int
	SuggestedGasTipCap *big.Int
}

// quantilesFromFloatArray calculates quantiles from a float array
func quantilesFromFloatArray(fa []float64) (*GasPercentiles, error) {
	perMax, err := stats.Max(fa)
	if err != nil {
		return nil, err
	}
	perc99, err := stats.Percentile(fa, 99)
	if err != nil {
		return nil, err
	}
	perc75, err := stats.Percentile(fa, 75)
	if err != nil {
		return nil, err
	}
	perc50, err := stats.Percentile(fa, 50)
	if err != nil {
		return nil, err
	}
	perc25, err := stats.Percentile(fa, 25)
	if err != nil {
		return nil, err
	}
	return &GasPercentiles{
		Max:    perMax,
		Perc99: perc99,
		Perc75: perc75,
		Perc50: perc50,
		Perc25: perc25,
	}, nil
}
