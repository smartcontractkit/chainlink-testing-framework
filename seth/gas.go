package seth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// GasEstimator estimates gas prices
type GasEstimator struct {
	Client              *Client
	BlockGasLimits      []uint64
	TransactionGasPrice []uint64
}

func (m *GasEstimator) logger() *zerolog.Logger {
	if m == nil || m.Client == nil {
		l := newLogger()
		return &l
	}
	return m.Client.Logger()
}

// NewGasEstimator creates a new gas estimator
func NewGasEstimator(c *Client) *GasEstimator {
	return &GasEstimator{Client: c}
}

// Stats calculates gas price and tip cap suggestions based on historical fee data over a specified number of blocks.
// It computes quantiles for base fees and tip caps and provides suggested gas price and tip cap values.
func (m *GasEstimator) Stats(ctx context.Context, blockCount uint64, priorityPerc float64) (GasSuggestions, error) {
	estimations := GasSuggestions{}
	logger := m.logger()

	if blockCount == 0 {
		return estimations, errors.New("block count must be greater than zero")
	}

	currentBlock, err := m.Client.Client.BlockNumber(ctx)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get current block number: %w", err)
	}
	if currentBlock == 0 {
		return GasSuggestions{}, errors.New("current block number is zero. No fee history available")
	}
	if blockCount >= currentBlock {
		blockCount = currentBlock - 1
	}

	hist, err := m.Client.Client.FeeHistory(ctx, blockCount, big.NewInt(mustSafeInt64(currentBlock)), []float64{priorityPerc})
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get fee history: %w", err)
	}
	logger.Trace().
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
		return GasSuggestions{}, fmt.Errorf("failed to calculate quantiles from fee history for base fee: %w", err)
	}
	estimations.BaseFeePerc = gasPercs

	logger.Trace().
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
		return GasSuggestions{}, fmt.Errorf("failed to calculate quantiles from fee history for tip cap: %w", err)
	}
	estimations.TipCapPerc = tipPercs
	logger.Trace().
		Interface("Gas percentiles ", tipPercs).
		Msg("Tip caps")

	suggestedGasPrice, err := m.Client.Client.SuggestGasPrice(ctx)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get suggested gas price: %w", err)
	}
	estimations.SuggestedGasPrice = suggestedGasPrice

	suggestedGasTipCap, err := m.Client.Client.SuggestGasTipCap(ctx)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get suggested gas tip cap: %w", err)
	}

	estimations.SuggestedGasTipCap = suggestedGasTipCap

	header, err := m.Client.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return GasSuggestions{}, fmt.Errorf("failed to get latest block header: %w", err)
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
