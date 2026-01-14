package leak

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// ClNodesCheck contains thresholds which can be verified for each Chainlink node
// it is recommended to set some WarmUpDuration, 20% of overall test time
// to have more stable results
type CLNodesCheck struct {
	NumNodes        int
	Start           time.Time
	End             time.Time
	WarmUpDuration  time.Duration
	CPUThreshold    float64
	MemoryThreshold float64
}

// CLNodesLeakDetector is Chainlink node specific resource leak detector
// can be used with both local and remote Chainlink node sets (DONs)
type CLNodesLeakDetector struct {
	Mode                  string
	CPUQuery, MemoryQuery string
	c                     *ResourceLeakChecker
}

// WithCPUQuery allows to override CPU leak query (Prometheus)
func WithCPUQuery(q string) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.CPUQuery = q
	}
}

// WithCPUQuery allows to override Memory leak query (Prometheus)
func WithMemoryQuery(q string) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.MemoryQuery = q
	}
}

// NewCLNodesLeakDetector create new Chainlink node specific resource leak detector with Prometheus client
func NewCLNodesLeakDetector(c *ResourceLeakChecker, opts ...func(*CLNodesLeakDetector)) (*CLNodesLeakDetector, error) {
	cd := &CLNodesLeakDetector{
		c: c,
	}
	for _, o := range opts {
		o(cd)
	}
	if cd.Mode == "" {
		cd.Mode = "devenv"
	}
	switch cd.Mode {
	case "devenv":
		cd.CPUQuery = `sum(rate(container_cpu_usage_seconds_total{name=~"don-node%d"}[5m])) * 100`
		cd.MemoryQuery = `quantile_over_time(0.5, container_memory_rss{name="don-node%d"}[1h]) / 1024 / 1024`
	case "griddle":
		return nil, fmt.Errorf("not implemented yet")
	default:
		return nil, fmt.Errorf("invalid mode, use: 'devenv' or 'griddle'")
	}
	return cd, nil
}

// Check runs all resource leak checks and returns errors if threshold reached for any of them
func (cd *CLNodesLeakDetector) Check(t *CLNodesCheck) error {
	if t.NumNodes == 0 {
		return fmt.Errorf("cl nodes num must be > 0")
	}
	memoryDiffs := make([]float64, 0)
	cpuDiffs := make([]float64, 0)
	errs := make([]error, 0)
	for i := range t.NumNodes {
		memoryDiff, err := cd.c.MeasureLeak(&CheckConfig{
			Query:          fmt.Sprintf(cd.MemoryQuery, i),
			Start:          t.Start,
			End:            t.End,
			WarmUpDuration: t.WarmUpDuration,
		})
		if err != nil {
			return fmt.Errorf("memory leak check failed: %w", err)
		}
		memoryDiffs = append(memoryDiffs, memoryDiff)
		cpuDiff, err := cd.c.MeasureLeak(&CheckConfig{
			Query:          fmt.Sprintf(cd.CPUQuery, i),
			Start:          t.Start,
			End:            t.End,
			WarmUpDuration: t.WarmUpDuration,
		})
		if err != nil {
			return fmt.Errorf("cpu leak check failed: %w", err)
		}
		cpuDiffs = append(cpuDiffs, cpuDiff)

		if memoryDiff >= t.MemoryThreshold {
			errs = append(errs, fmt.Errorf(
				"Memory leak detected for node %d and interval: [%s -> %s], diff: %.f",
				i, t.Start, t.End, memoryDiff,
			))
		}
		if cpuDiff >= t.CPUThreshold {
			errs = append(errs, fmt.Errorf(
				"CPU leak detected for node %d and interval: [%s -> %s], diff: %.f",
				i, t.Start, t.End, cpuDiff,
			))
		}
	}
	framework.L.Info().
		Any("MemoryDiffs", memoryDiffs).
		Any("CPUDiffs", cpuDiffs).
		Msg("Leaks info")
	return errors.Join(errs...)
}
