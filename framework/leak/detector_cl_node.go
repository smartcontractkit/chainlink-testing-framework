package leak

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// ClNodesCheck contains thresholds which can be verified for each Chainlink node
// it is recommended to set some WarmUpDuration, 20% of overall test time
// to have more stable results
type CLNodesCheck struct {
	CheckMode       string
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
	EnvironmentType                            string
	CPUQuery, MemoryQuery, ContainerAliveQuery string
	c                                          *ResourceLeakChecker
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
	if cd.EnvironmentType == "" {
		cd.EnvironmentType = "devenv"
	}
	switch cd.EnvironmentType {
	case "devenv":
		cd.ContainerAliveQuery = `time() - container_start_time_seconds{name=~"don-node%d"}`
		// avg from intervals of 1h with 30m step to mitigate spikes
		cd.CPUQuery = `avg_over_time((sum(rate(container_cpu_usage_seconds_total{name="don-node%d"}[1h])))[1h:30m]) * 100`
		cd.MemoryQuery = `avg_over_time(container_memory_rss{name="don-node%d"}[1h:30m]) / 1024 / 1024`
	case "griddle":
		return nil, fmt.Errorf("not implemented yet")
	default:
		return nil, fmt.Errorf("invalid mode, use: 'devenv' or 'griddle'")
	}
	return cd, nil
}

func (cd *CLNodesLeakDetector) checkContainerUptime(t *CLNodesCheck, nodeIdx int) (float64, error) {
	uptimeResp, err := cd.c.c.Query(fmt.Sprintf(cd.ContainerAliveQuery, nodeIdx), t.End)
	if err != nil {
		return 0, fmt.Errorf("failed to execute container alive query: %w", err)
	}
	uptimeResult := uptimeResp.Data.Result
	if len(uptimeResult) == 0 {
		return 0, fmt.Errorf("no results for end timestamp: %s", t.End)
	}

	uptimeResultValue, resOk := uptimeResult[0].Value[1].(string)
	if !resOk {
		return 0, fmt.Errorf("invalid Prometheus response value for timestamp: %s, value: %v", t.End, uptimeResult[0].Value[1])
	}

	uptimeResultValueFloat, err := strconv.ParseFloat(uptimeResultValue, 64)
	if err != nil {
		return 0, fmt.Errorf("uptime can't be parsed from string: %w", err)
	}
	if uptimeResultValueFloat <= float64(t.End.Unix())-float64(t.Start.Unix()) {
		return uptimeResultValueFloat, fmt.Errorf("container hasn't lived long enough and was killed while the test was running")
	}
	return uptimeResultValueFloat, nil
}

// Check runs all resource leak checks and returns errors if threshold reached for any of them
func (cd *CLNodesLeakDetector) Check(t *CLNodesCheck) error {
	if t.NumNodes == 0 {
		return fmt.Errorf("cl nodes num must be > 0")
	}
	memMeasurements := make([]*Measurement, 0)
	cpuMeasurements := make([]*Measurement, 0)
	uptimes := make([]float64, 0)
	errs := make([]error, 0)
	for i := range t.NumNodes {
		memMeasurement, err := cd.c.MeasureDelta(&CheckConfig{
			CheckMode:           t.CheckMode,
			Query:          fmt.Sprintf(cd.MemoryQuery, i),
			Start:          t.Start,
			End:            t.End,
			WarmUpDuration: t.WarmUpDuration,
		})
		if err != nil {
			return fmt.Errorf("memory leak check failed: %w", err)
		}
		memMeasurements = append(memMeasurements, memMeasurement)

		cpuMeasurement, err := cd.c.MeasureDelta(&CheckConfig{
			CheckMode:           t.CheckMode,
			Query:          fmt.Sprintf(cd.CPUQuery, i),
			Start:          t.Start,
			End:            t.End,
			WarmUpDuration: t.WarmUpDuration,
		})
		if err != nil {
			return fmt.Errorf("cpu leak check failed: %w", err)
		}
		cpuMeasurements = append(cpuMeasurements, cpuMeasurement)

		switch t.CheckMode {
		case CheckModePercentage:
			fallthrough
		case CheckModeDiff:
			if memMeasurement.Delta >= t.MemoryThreshold {
				errs = append(errs, fmt.Errorf(
					"Memory leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, memMeasurement, t.CheckMode,
				))
			}
			if cpuMeasurement.Delta >= t.CPUThreshold {
				errs = append(errs, fmt.Errorf(
					"CPU leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, cpuMeasurement, t.CheckMode,
				))
			}
		case CheckModeAbsolute:
			if memMeasurement.End >= t.MemoryThreshold {
				errs = append(errs, fmt.Errorf(
					"Memory leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, memMeasurement, t.CheckMode,
				))
			}
			if cpuMeasurement.End >= t.CPUThreshold {
				errs = append(errs, fmt.Errorf(
					"CPU leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, cpuMeasurement, t.CheckMode,
				))
			}
		}

		uptime, err := cd.checkContainerUptime(t, i)
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"Container uptime issue for node %d and interval: [%s -> %s], uptime: %.f, err: %w",
				i, t.Start, t.End, uptime, err,
			))
		}
		uptimes = append(uptimes, uptime)
	}
	framework.L.Info().
		Any("MemoryDiffs", memMeasurements).
		Any("CPUDiffs", cpuMeasurements).
		Any("Uptimes", uptimes).
		Str("TestDuration", t.End.Sub(t.Start).String()).
		Float64("TestDurationSec", t.End.Sub(t.Start).Seconds()).
		Msg("Leaks info")
	framework.L.Info().Msg("Downloading pprof profile..")
	dumper := NewProfileDumper(framework.LocalPyroscopeBaseURL)
	profilePath, err := dumper.MemoryProfile(&ProfileDumperConfig{
		ServiceName: "chainlink-node",
	})
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to download Pyroscopt profile: %w", err))
		return errors.Join(errs...)
	}
	framework.L.Info().Str("Path", profilePath).Msg("Saved pprof profile")
	return errors.Join(errs...)
}
