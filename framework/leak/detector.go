package leak

/*
Resource leak detector
This module provides a Prometheus-based leak detector for long-running soak tests. It detects leaks by comparing the median resource usage at the start and end of a test and flags any increases that breach configured thresholds.

Usage Note: Set the WarmUpDuration to at least 20% of your test length for reliable metrics.
It is also recommend to use it with 3h+ soak tests for less false-positives.
*/

import (
	"fmt"
	"strconv"
	"time"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// ResourceLeakCheckerConfig is resource leak checker config with Prometheus base URL
type ResourceLeakCheckerConfig struct {
	PrometheusBaseURL string
}

// ResourceLeakChecker is resource leak cheker instance
type ResourceLeakChecker struct {
	PrometheusURL string
	c             PromQuerier
}

// WithPrometheusBaseURL sets Prometheus base URL, example http://localhost:9099
func WithPrometheusBaseURL(url string) func(*ResourceLeakChecker) {
	return func(rlc *ResourceLeakChecker) {
		rlc.PrometheusURL = url
	}
}

// WithQueryClient sets Prometheus query client
func WithQueryClient(c PromQuerier) func(*ResourceLeakChecker) {
	return func(rlc *ResourceLeakChecker) {
		rlc.c = c
	}
}

// PromQueries is an interface for querying Prometheus containing only methods we need for detecting resource leaks
type PromQuerier interface {
	Query(query string, timestamp time.Time) (*f.PrometheusQueryResponse, error)
}

// NewResourceLeakChecker creates a new resource leak checker
func NewResourceLeakChecker(opts ...func(*ResourceLeakChecker)) *ResourceLeakChecker {
	lc := &ResourceLeakChecker{}
	for _, o := range opts {
		o(lc)
	}
	if lc.c == nil {
		lc.c = f.NewPrometheusQueryClient(f.LocalPrometheusBaseURL)
	}
	return lc
}

// CheckConfig describes leak check configuration
type CheckConfig struct {
	// ComparisonMode how do we compaore start/end values: percentage, diff or absolute
	ComparisonMode string
	// Query prometheus query
	Query string
	// Start start time
	Start time.Time
	// End end time
	End time.Time
	// WarmUpDuration duration that will be excluded from comparison, load/soak test warmup duration
	WarmUpDuration time.Duration
}

const (
	// ComparisonModePercentage compares start and end values in percentage
	ComparisonModePercentage = "percentage"
	// ComparisonModeDiff compares start and end values by subtracting start from end
	ComparisonModeDiff = "diff"
	// ComparisonModeAbsolute compares only end values with threshold
	ComparisonModeAbsolute = "absolute"
)

type Measurement struct {
	Start float64
	End   float64
	Delta float64
}

// MeasureDelta measures resource leak delta between start and end timestamps
// WarmUpDuration is used to ignore warm up interval results for more stable comparison
func (rc *ResourceLeakChecker) MeasureDelta(
	c *CheckConfig,
) (*Measurement, error) {
	if c.Start.After(c.End) {
		return nil, fmt.Errorf("start time is greated than end time: %s -> %s", c.Start, c.End)
	}
	if c.WarmUpDuration > c.End.Sub(c.Start)/2 {
		return nil, fmt.Errorf("warm up duration can't be more than 50 percent of test interval between start and end timestamps: %s", c.WarmUpDuration)
	}
	startWithWarmUp := c.Start.Add(c.WarmUpDuration)
	memStart, err := rc.c.Query(c.Query, startWithWarmUp)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory for the test start: %w", err)
	}

	memEnd, err := rc.c.Query(c.Query, c.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory for the test end: %w", err)
	}

	resStart := memStart.Data.Result
	resEnd := memEnd.Data.Result
	if len(resStart) == 0 {
		return nil, fmt.Errorf("no results for start timestamp: %s, query: %s", startWithWarmUp, c.Query)
	}
	if len(resEnd) == 0 {
		return nil, fmt.Errorf("no results for end timestamp: %s, query: %s", c.End, c.Query)
	}

	if len(resStart[0].Value) < 2 {
		return nil, fmt.Errorf("invalid Prometheus response for start timestamp, should have timestamp and value: %s", c.Start)
	}
	if len(resEnd[0].Value) < 2 {
		return nil, fmt.Errorf("invalid Prometheus response for end timestamp, should have timestamp and value: %s", c.End)
	}

	memStartVal, startOk := memStart.Data.Result[0].Value[1].(string)
	if !startOk {
		return nil, fmt.Errorf("invalid Prometheus response value for timestamp: %s, value: %v", c.Start, memStart.Data.Result[0].Value[1])
	}
	memEndVal, endOk := memEnd.Data.Result[0].Value[1].(string)
	if !endOk {
		return nil, fmt.Errorf("invalid Prometheus response value for timestamp: %s, value: %v", c.End, memEnd.Data.Result[0].Value[1])
	}

	startValFloat, err := strconv.ParseFloat(memStartVal, 64)
	if err != nil {
		return nil, fmt.Errorf("start quantile can't be parsed from string: %w", err)
	}
	endValFloat, err := strconv.ParseFloat(memEndVal, 64)
	if err != nil {
		return nil, fmt.Errorf("start quantile can't be parsed from string: %w", err)
	}

	var delta float64
	switch c.ComparisonMode {
	case ComparisonModePercentage:
		delta = (endValFloat / startValFloat * 100) - 100
	case ComparisonModeDiff:
		delta = endValFloat - startValFloat
	}

	f.L.Info().
		Str("Mode", c.ComparisonMode).
		Float64("Start", startValFloat).
		Float64("End", endValFloat).
		Float64("Increase", delta).
		Msg("Increase total")
	return &Measurement{Start: startValFloat, End: endValFloat, Delta: delta}, nil
}
