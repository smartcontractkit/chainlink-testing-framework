package common

import (
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"sync"
	"time"
)

// TestSpec common test spec
type TestSpec struct {
	EnvSetup                *actions.DefaultSuiteSetup
	Nodes                   []client.Chainlink
	NodesAddresses          []common.Address
	Adapter                 environment.ExternalAdapter
	OnChainCheckAttemptsOpt func(config *retry.Config)
}

// InstanceDeployment common deploy data for one contract deployed concurrently
type InstanceDeployment struct {
	Wg            *sync.WaitGroup
	Index         int
	Suite         *actions.DefaultSuiteSetup
	NodeAddresses []common.Address
	Nodes         []client.Chainlink
	Adapter       environment.ExternalAdapter
}

// Test common test runtime instance
type Test struct {
	DefaultSetup            *actions.DefaultSuiteSetup
	OnChainCheckAttemptsOpt func(config *retry.Config)
	Nodes                   []client.Chainlink
	Adapter                 environment.ExternalAdapter
	Prom                    *tools.PromChecker
}

// PercentileReport common percentile report
type PercentileReport struct {
	StdDev float64
	Max    float64
	Min    float64
	P99    float64
	P95    float64
	P90    float64
	P50    float64
}

// PrintPercentileMetrics prints percentile metrics
func (t *Test) PrintPercentileMetrics(m *PercentileReport) {
	log.Info().Float64("Round duration ms MAX", m.Max).Send()
	log.Info().Float64("Round duration ms P99", m.P99).Send()
	log.Info().Float64("Round duration ms P95", m.P95).Send()
	log.Info().Float64("Round duration ms P90", m.P90).Send()
	log.Info().Float64("Round duration ms P50", m.P50).Send()
	log.Info().Float64("Round duration ms MIN", m.Min).Send()
	log.Info().Float64("Round duration ms standard deviation", m.StdDev).Send()
}

// CalculatePercentiles calculates percentiles for arbitrary float64 data
func (t *Test) CalculatePercentiles(data []time.Duration) (*PercentileReport, error) {
	dataFloat64 := make([]float64, 0)
	for _, d := range data {
		dataFloat64 = append(dataFloat64, float64(d.Milliseconds()))
	}
	perc99, err := stats.Percentile(dataFloat64, 99)
	if err != nil {
		return nil, err
	}
	perc95, err := stats.Percentile(dataFloat64, 95)
	if err != nil {
		return nil, err
	}
	perc90, err := stats.Percentile(dataFloat64, 90)
	if err != nil {
		return nil, err
	}
	perc50, err := stats.Percentile(dataFloat64, 50)
	if err != nil {
		return nil, err
	}
	max, err := stats.Max(dataFloat64)
	if err != nil {
		return nil, err
	}
	min, err := stats.Min(dataFloat64)
	if err != nil {
		return nil, err
	}
	stdDev, err := stats.StandardDeviation(dataFloat64)
	if err != nil {
		return nil, err
	}
	return &PercentileReport{P99: perc99, P95: perc95, P90: perc90, P50: perc50, Max: max, Min: min, StdDev: stdDev}, nil
}
