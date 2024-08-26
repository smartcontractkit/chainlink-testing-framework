package seth

import (
	"context"
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

// Stats prints gas stats
func (m *GasEstimator) Stats(fromNumber uint64, priorityPerc float64) (GasSuggestions, error) {
	bn, err := m.Client.Client.BlockNumber(context.Background())
	if err != nil {
		return GasSuggestions{}, err
	}
	hist, err := m.Client.Client.FeeHistory(context.Background(), fromNumber, big.NewInt(int64(bn)), []float64{priorityPerc})
	if err != nil {
		return GasSuggestions{}, err
	}
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
		return GasSuggestions{}, err
	}
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
		return GasSuggestions{}, err
	}
	suggestedGasPrice, err := m.Client.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return GasSuggestions{}, err
	}
	suggestedGasTipCap, err := m.Client.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return GasSuggestions{}, err
	}
	L.Trace().
		Interface("History", hist).
		Msg("Fee history")
	return GasSuggestions{
		GasPrice:           gasPercs,
		TipCap:             tipPercs,
		SuggestedGasPrice:  suggestedGasPrice,
		SuggestedGasTipCap: suggestedGasTipCap,
	}, nil
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
	GasPrice           *GasPercentiles
	TipCap             *GasPercentiles
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
