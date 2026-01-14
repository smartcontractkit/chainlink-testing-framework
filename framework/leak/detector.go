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
	Query          string
	Start          time.Time
	End            time.Time
	WarmUpDuration time.Duration
}

// MeasureLeak measures resource leak between start and end timestamps
// WarmUpDuration is used to ignore warm up interval results for more stable comparison
func (rc *ResourceLeakChecker) MeasureLeak(
	c *CheckConfig,
) (float64, error) {
	if c.Start.After(c.End) {
		return 0, fmt.Errorf("start time is greated than end time: %s -> %s", c.Start, c.End)
	}
	if c.WarmUpDuration > c.End.Sub(c.Start)/2 {
		return 0, fmt.Errorf("warm up duration can't be more than 50 percent of test interval between start and end timestamps: %s", c.WarmUpDuration)
	}
	startWithWarmUp := c.Start.Add(c.WarmUpDuration)
	memStart, err := rc.c.Query(c.Query, startWithWarmUp)
	if err != nil {
		return 0, fmt.Errorf("failed to get memory for the test start: %w", err)
	}

	memEnd, err := rc.c.Query(c.Query, c.End)
	if err != nil {
		return 0, fmt.Errorf("failed to get memory for the test end: %w", err)
	}

	resStart := memStart.Data.Result
	resEnd := memEnd.Data.Result
	if len(resStart) == 0 {
		return 0, fmt.Errorf("no results for start timestamp: %s", c.Start)
	}
	if len(resEnd) == 0 {
		return 0, fmt.Errorf("no results for end timestamp: %s", c.End)
	}

	if len(resStart[0].Value) < 2 {
		return 0, fmt.Errorf("invalid Prometheus response for start timestamp, should have timestamp and value: %s", c.Start)
	}
	if len(resEnd[0].Value) < 2 {
		return 0, fmt.Errorf("invalid Prometheus response for end timestamp, should have timestamp and value: %s", c.End)
	}

	memStartVal, startOk := memStart.Data.Result[0].Value[1].(string)
	if !startOk {
		return 0, fmt.Errorf("invalid Prometheus response value for timestamp: %s, value: %v", c.Start, memStart.Data.Result[0].Value[1])
	}
	memEndVal, endOk := memEnd.Data.Result[0].Value[1].(string)
	if !endOk {
		return 0, fmt.Errorf("invalid Prometheus response value for timestamp: %s, value: %v", c.End, memEnd.Data.Result[0].Value[1])
	}

	memStartValFloat, err := strconv.ParseFloat(memStartVal, 64)
	if err != nil {
		return 0, fmt.Errorf("start quantile can't be parsed from string: %w", err)
	}
	memEndValFloat, err := strconv.ParseFloat(memEndVal, 64)
	if err != nil {
		return 0, fmt.Errorf("start quantile can't be parsed from string: %w", err)
	}

	totalIncreasePercentage := (memEndValFloat / memStartValFloat * 100) - 100

	f.L.Debug().
		Float64("Start", memStartValFloat).
		Float64("End", memEndValFloat).
		Float64("Increase", totalIncreasePercentage).
		Msg("Memory increase total (percentage)")
	return totalIncreasePercentage, nil
}
