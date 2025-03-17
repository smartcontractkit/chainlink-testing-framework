package benchspy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

// BasicData is the basic data that is required for a report, common to all reports
type BasicData struct {
	TestName    string `json:"test_name"`
	CommitOrTag string `json:"commit_or_tag"`

	// Test metrics
	TestStart time.Time `json:"test_start_timestamp"`
	TestEnd   time.Time `json:"test_end_timestamp"`

	// all, generator settings, including segments
	GeneratorConfigs map[string]*wasp.Config `json:"generator_configs"`
}

// MustNewBasicData creates a new BasicData instance from a commit or tag.
// It panics if the creation fails, ensuring that the caller receives a valid instance.
func MustNewBasicData(commitOrTag string, generators ...*wasp.Generator) BasicData {
	b, err := NewBasicData(commitOrTag, generators...)
	if err != nil {
		panic(err)
	}

	return *b
}

// NewBasicData creates a new BasicData instance using the provided commit or tag and a list of generators.
// It ensures that at least one generator is provided and that it is associated with a testing.T instance.
// This function is essential for initializing test data configurations in a structured manner.
func NewBasicData(commitOrTag string, generators ...*wasp.Generator) (*BasicData, error) {
	L.Debug().
		Msg("Creating new basic data instance")

	if len(generators) == 0 {
		return nil, errors.New("at least one generator is required")
	}

	if generators[0].Cfg.T == nil {
		return nil, errors.New("generators are not associated with a testing.T instance. Please set it as generator.Cfg.T and try again")
	}

	b := &BasicData{
		TestName:         generators[0].Cfg.T.Name(),
		CommitOrTag:      commitOrTag,
		GeneratorConfigs: make(map[string]*wasp.Config),
	}

	for _, g := range generators {
		b.GeneratorConfigs[g.Cfg.GenName] = g.Cfg
	}

	timeErr := b.FillStartEndTimes()
	if timeErr != nil {
		return nil, timeErr
	}

	L.Debug().
		Msg("Basic data instance created successfully")

	return b, nil
}

// FillStartEndTimes calculates the earliest start time and latest end time from generator schedules.
// It updates the BasicData instance with these times, ensuring all segments have valid start and end times.
func (b *BasicData) FillStartEndTimes() error {
	L.Debug().
		Msg("Filling test start and end times for basic data instance based on generator schedules")

	earliestTime := time.Now()
	var latestTime time.Time

	for _, cfg := range b.GeneratorConfigs {
		if len(cfg.Schedule) == 0 {
			return fmt.Errorf("schedule is empty for generator %s", cfg.GenName)
		}

		for _, segment := range cfg.Schedule {
			if segment.StartTime.IsZero() {
				return fmt.Errorf("start time is missing in one of the segments belonging to generator %s. Did that generator run?", cfg.GenName)
			}
			if segment.StartTime.Before(earliestTime) {
				earliestTime = segment.StartTime
			}
			if segment.EndTime.IsZero() {
				return fmt.Errorf("end time is missing in one of the segments belonging to generator %s. Did that generator finish running?", cfg.GenName)
			}
			if segment.EndTime.After(latestTime) {
				latestTime = segment.EndTime
			}
		}
	}

	b.TestStart = earliestTime
	b.TestEnd = latestTime
	L.Debug().
		Str("Test start time", earliestTime.Format(time.RFC3339)).
		Str("Test end time", latestTime.Format(time.RFC3339)).
		Msg("Start and end times filled successfully")

	return nil
}

// Validate checks the integrity of the BasicData fields, ensuring that the test start and end times are set,
// and that at least one generator configuration is provided. It returns an error if any of these conditions are not met.
func (b *BasicData) Validate() error {
	L.Debug().
		Msg("Validating basic data instance")
	if b.TestStart.IsZero() {
		return errors.New("test start time is missing. We cannot query Loki without a time range. Please set it and try again")
	}
	if b.TestEnd.IsZero() {
		return errors.New("test end time is missing. We cannot query Loki without a time range. Please set it and try again")
	}

	if b.TestEnd.Before(b.TestStart) {
		return errors.New("test end time is before test start time. Please set valid times and try again")
	}
	if b.TestEnd.Sub(b.TestStart) < time.Second {
		return errors.New("test duration is less than a second. Please set a valid time range and try again")
	}

	if len(b.GeneratorConfigs) == 0 {
		return errors.New("generator configs are missing. At least one is required. Please set them and try again")
	}

	L.Debug().
		Msg("Basic data instance is valid")

	return nil
}

// IsComparable checks if two BasicData instances have the same configuration settings.
// It validates the count, presence, and equivalence of generator configurations,
// returning an error if any discrepancies are found. This function is useful for ensuring
// consistency between data reports before processing or comparison.
func (b *BasicData) IsComparable(otherData BasicData) error {
	L.Debug().
		Msg("Checking if basic data instances are comparable")

	if len(b.GeneratorConfigs) != len(otherData.GeneratorConfigs) {
		return fmt.Errorf("generator configs count is different. Expected %d, got %d", len(b.GeneratorConfigs), len(otherData.GeneratorConfigs))
	}

	for name1, cfg1 := range b.GeneratorConfigs {
		cfg2, ok := otherData.GeneratorConfigs[name1]
		if !ok {
			return fmt.Errorf("generator config %s is missing from the other report", name1)
		}
		if err := compareGeneratorConfigs(cfg1, cfg2); err != nil {
			return err
		}
	}

	L.Debug().
		Msg("Basic data instances are comparable")

	return nil
}

func compareGeneratorConfigs(cfg1, cfg2 *wasp.Config) error {
	if cfg1.GenName != cfg2.GenName {
		return fmt.Errorf("generator names are different. Expected %s, got %s", cfg1.GenName, cfg2.GenName)
	}
	if cfg1.LoadType != cfg2.LoadType {
		return fmt.Errorf("load types are different. Expected %s, got %s", cfg1.LoadType, cfg2.LoadType)
	}

	if len(cfg1.Schedule) != len(cfg2.Schedule) {
		return fmt.Errorf("schedules are different. Expected %d, got %d", len(cfg1.Schedule), len(cfg2.Schedule))
	}

	var areSegmentsEqual = func(segment1, segment2 *wasp.Segment) bool {
		return segment1.From == segment2.From && segment1.Duration == segment2.Duration && segment1.Type == segment2.Type
	}

	for i, segment1 := range cfg1.Schedule {
		segment2 := cfg2.Schedule[i]
		if segment1 == nil {
			return fmt.Errorf("segment at index %d is nil in the current report", i)
		}
		if segment2 == nil {
			return fmt.Errorf("segment at index %d is nil in the other report", i)
		}
		if !areSegmentsEqual(segment1, segment2) {
			return fmt.Errorf("segments at index %d are different. Expected %s segment(s), got %s segment(s)", i, mustMarshallSegment(segment1), mustMarshallSegment(segment2))
		}
	}

	if cfg1.CallTimeout != cfg2.CallTimeout {
		return fmt.Errorf("call timeouts are different. Expected %s, got %s", cfg1.CallTimeout, cfg2.CallTimeout)
	}

	if cfg1.RateLimitUnitDuration != cfg2.RateLimitUnitDuration {
		return fmt.Errorf("rate limit unit durations are different. Expected %s, got %s", cfg1.RateLimitUnitDuration, cfg2.RateLimitUnitDuration)
	}

	return nil
}

func mustMarshallSegment(segment *wasp.Segment) string {
	segmentBytes, err := json.MarshalIndent(segment, "", " ")
	if err != nil {
		panic(err)
	}

	return string(segmentBytes)
}
