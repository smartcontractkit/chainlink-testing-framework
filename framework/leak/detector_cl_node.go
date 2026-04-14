package leak

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// ClNodesCheck contains thresholds which can be verified for each Chainlink node
// it is recommended to set some WarmUpDuration, 20% of overall test time
// to have more stable results
type CLNodesCheck struct {
	// ComparisonMode how do we compaore start/end values: percentage, diff or absolute
	ComparisonMode string
	// NumNodes number of Chainlink nodes
	NumNodes int
	// Start start time
	Start time.Time
	// End end time
	End time.Time
	// WarmUpDuration duration that will be excluded from comparison, load/soak test warmup duration
	WarmUpDuration time.Duration
	// CPUThreshold CPU threshold as a float: 200.0 means full 2 CPU cores
	CPUThreshold float64
	// MemoryThreshold memory threshold in Megabytes
	MemoryThreshold float64
}

// CLNodesLeakDetector is Chainlink node specific resource leak detector
// can be used with both local and remote Chainlink node sets (DONs)
type CLNodesLeakDetector struct {
	// CPUQuery Prometheus query for CPU
	CPUQuery string
	// MemoryQuery Prometheus query for memory
	MemoryQuery         string
	CPUQueryAbsolute    string
	MemoryQueryAbsolute string
	// ContainerAliveQuery Prometheus memory for checking if container was alive the whole time
	ContainerAliveQuery string
	c                   *ResourceLeakChecker

	nodesetName           string
	dumpPyroscopeProfiles bool
	dumpAdminProfiles     bool
}

// WithCPUQuery allows to override CPU leak query (Prometheus)
func WithCPUQuery(q string) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.CPUQuery = q
	}
}

// WithMemoryQuery allows to override Memory leak query (Prometheus)
func WithMemoryQuery(q string) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.MemoryQuery = q
	}
}

// WithNodesetName overrides the default nodeset name "don" in all Prometheus queries.
// The name is used to build container labels like "name-node0", "name-node1", etc.
// Name should be alphanumeric with hyphens/underscores; characters that could break
// format strings (% or PromQL literals like ", \) are escaped for safety.
func WithNodesetName(name string) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.nodesetName = sanitizeNodesetName(name)
	}
}

// WithDumpPyroscopeProfiles allows to dump Pyroscope profiles for each node at the end of the test.
// Dumped profiles are aggragate (cumulative) profiles from the whole test duration.
func WithDumpPyroscopeProfiles(dump bool) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.dumpPyroscopeProfiles = dump
	}
}

// WithDumpAdminProfiles allows to dump admin profiles for each node at the end of the test.
// Uses CL node's debug endpoint to fetch pprof snapshots.
func WithDumpAdminProfiles(dump bool) func(*CLNodesLeakDetector) {
	return func(cd *CLNodesLeakDetector) {
		cd.dumpAdminProfiles = dump
	}
}

// sanitizeNodesetName escapes characters that would corrupt fmt.Sprintf format strings
// or invalidate PromQL double-quoted label literals.
func sanitizeNodesetName(name string) string {
	s := strings.ReplaceAll(name, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "%", "%%")
	return s
}

// NewCLNodesLeakDetector create new Chainlink node specific resource leak detector with Prometheus client
func NewCLNodesLeakDetector(c *ResourceLeakChecker, opts ...func(*CLNodesLeakDetector)) (*CLNodesLeakDetector, error) {
	cd := &CLNodesLeakDetector{
		c:                   c,
		ContainerAliveQuery: `time() - container_start_time_seconds{name=~"don-node%d"}`,
		// avg from intervals of 1h with 30m step to mitigate spikes
		// these queries is for estimating soak tests which runs for 2h+
		CPUQuery:    `avg_over_time((sum(rate(container_cpu_usage_seconds_total{name="don-node%d"}[1h])))[1h:30m]) * 100`,
		MemoryQuery: `avg_over_time(container_memory_rss{name="don-node%d"}[1h:30m]) / 1024 / 1024`,
		// these are for very stable soak tests where we want to catch even a small deviation
		// by measuring only the end value
		CPUQueryAbsolute:    `sum(rate(container_cpu_usage_seconds_total{name="don-node%d"}[5m])) * 100`,
		MemoryQueryAbsolute: `avg_over_time(container_memory_rss{name="don-node%d"}[5m:5m]) / 1024 / 1024`,
	}
	for _, o := range opts {
		o(cd)
	}

	if cd.nodesetName != "" {
		replaceNodeset := func(s string) string {
			return strings.ReplaceAll(s, "don-node%d", cd.nodesetName+"-node%d")
		}
		cd.ContainerAliveQuery = replaceNodeset(cd.ContainerAliveQuery)
		cd.CPUQuery = replaceNodeset(cd.CPUQuery)
		cd.MemoryQuery = replaceNodeset(cd.MemoryQuery)
		cd.CPUQueryAbsolute = replaceNodeset(cd.CPUQueryAbsolute)
		cd.MemoryQueryAbsolute = replaceNodeset(cd.MemoryQueryAbsolute)
	}

	if cd.dumpPyroscopeProfiles == true && cd.dumpAdminProfiles == true {
		return nil, fmt.Errorf("both Pyroscope and admin profile dumping enabled, please choose only one. Dumping admin profiles will fail if Pyroscope is enabled.")
	}

	if cd.dumpAdminProfiles == false && cd.dumpPyroscopeProfiles == false {
		// default to dumping admin profiles since that's what engineers prefer
		cd.dumpAdminProfiles = true
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

		switch t.ComparisonMode {
		case ComparisonModePercentage:
			fallthrough
		case ComparisonModeDiff:
			memMeasurement, err := cd.c.MeasureDelta(&CheckConfig{
				ComparisonMode: t.ComparisonMode,
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
				ComparisonMode: t.ComparisonMode,
				Query:          fmt.Sprintf(cd.CPUQuery, i),
				Start:          t.Start,
				End:            t.End,
				WarmUpDuration: t.WarmUpDuration,
			})
			if err != nil {
				return fmt.Errorf("cpu leak check failed: %w", err)
			}
			cpuMeasurements = append(cpuMeasurements, cpuMeasurement)

			if memMeasurement.Delta >= t.MemoryThreshold {
				errs = append(errs, fmt.Errorf(
					"Memory leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, memMeasurement.Delta, t.ComparisonMode,
				))
			}
			if cpuMeasurement.Delta >= t.CPUThreshold {
				errs = append(errs, fmt.Errorf(
					"CPU leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, cpuMeasurement.Delta, t.ComparisonMode,
				))
			}
		case ComparisonModeAbsolute:
			memMeasurement, err := cd.c.MeasureDelta(&CheckConfig{
				ComparisonMode: t.ComparisonMode,
				Query:          fmt.Sprintf(cd.MemoryQueryAbsolute, i),
				Start:          t.Start,
				End:            t.End,
				WarmUpDuration: t.WarmUpDuration,
			})
			if err != nil {
				return fmt.Errorf("memory leak check failed: %w", err)
			}
			memMeasurements = append(memMeasurements, memMeasurement)

			cpuMeasurement, err := cd.c.MeasureDelta(&CheckConfig{
				ComparisonMode: t.ComparisonMode,
				Query:          fmt.Sprintf(cd.CPUQueryAbsolute, i),
				Start:          t.Start,
				End:            t.End,
				WarmUpDuration: t.WarmUpDuration,
			})
			if err != nil {
				return fmt.Errorf("cpu leak check failed: %w", err)
			}
			cpuMeasurements = append(cpuMeasurements, cpuMeasurement)
			if memMeasurement.End >= t.MemoryThreshold {
				errs = append(errs, fmt.Errorf(
					"Memory leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, memMeasurement.End, t.ComparisonMode,
				))
			}
			if cpuMeasurement.End >= t.CPUThreshold {
				errs = append(errs, fmt.Errorf(
					"CPU leak detected for node %d and interval: [%s -> %s], diff: %.f, comparison mode: %s",
					i, t.Start, t.End, cpuMeasurement.End, t.ComparisonMode,
				))
			}
		default:
			return fmt.Errorf("comparison mode is incorrect: %s, see available leak.ComparisonMode constants", t.ComparisonMode)
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

	if cd.dumpPyroscopeProfiles {
		profilesToDump := []string{DefaultProfileType, "memory:inuse_space:bytes:space:bytes"}
		framework.L.Info().Msgf("Downloading %d pprof profiles..", len(profilesToDump))
		dumper := NewProfileDumper(framework.LocalPyroscopeBaseURL)

		for _, profileType := range profilesToDump {
			profileSplit := strings.Split(profileType, ":")
			outputPath := DefaultOutputPath
			if len(profileSplit) > 1 {
				// e.g. for "memory:inuse_space:bytes:space:bytes" we want to have output file "memory-inuse_space.pprof"
				outputPath = fmt.Sprintf("%s-%s.pprof", profileSplit[0], profileSplit[1])
			}
			profilePath, err := dumper.MemoryProfile(&ProfileDumperConfig{
				ServiceName: "chainlink-node",
				ProfileType: profileType,
				OutputPath:  outputPath,
			})
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to download Pyroscope profile %s: %w", profileType, err))
				return errors.Join(errs...)
			}
			framework.L.Info().Str("Path", profilePath).Str("ProfileType", profileType).Msg("Saved pprof profile")
		}
	}

	if cd.dumpAdminProfiles {
		framework.L.Info().Msg("Dumping admin profiles..")
		ctx, cancel := context.WithTimeout(context.Background(), DefaultNodeProfileDumpTimeout)
		defer cancel()
		if err := DumpNodeProfiles(ctx, cd.nodesetName+"-node", DefaultAdminProfilesDir); err != nil {
			framework.L.Error().Err(err).Msg("Failed to dump node profiles")
			errs = append(errs, fmt.Errorf("failed to dump node profiles: %w", err))
		}
		framework.L.Info().Str("Path", DefaultAdminProfilesDir).Msg("Admin profiles dumped successfully")
	}

	return errors.Join(errs...)
}
