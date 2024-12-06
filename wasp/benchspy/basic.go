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

func MustNewBasicData(commitOrTag string, generators ...*wasp.Generator) BasicData {
	b, err := NewBasicData(commitOrTag, generators...)
	if err != nil {
		panic(err)
	}

	return *b
}

func NewBasicData(commitOrTag string, generators ...*wasp.Generator) (*BasicData, error) {
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

	return b, nil
}

func (b *BasicData) FillStartEndTimes() error {
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

	return nil
}

func (b *BasicData) Validate() error {
	if b.TestStart.IsZero() {
		return errors.New("test start time is missing. We cannot query Loki without a time range. Please set it and try again")
	}
	if b.TestEnd.IsZero() {
		return errors.New("test end time is missing. We cannot query Loki without a time range. Please set it and try again")
	}

	if len(b.GeneratorConfigs) == 0 {
		return errors.New("generator configs are missing. At least one is required. Please set them and try again")
	}

	return nil
}

func (b *BasicData) IsComparable(otherData BasicData) error {
	// are all configs present? do they have the same schedule type? do they have the same segments? is call timeout the same? is rate limit timeout the same?
	if len(b.GeneratorConfigs) != len(otherData.GeneratorConfigs) {
		return fmt.Errorf("generator configs count is different. Expected %d, got %d", len(b.GeneratorConfigs), len(otherData.GeneratorConfigs))
	}

	for name1, cfg1 := range b.GeneratorConfigs {
		if cfg2, ok := otherData.GeneratorConfigs[name1]; !ok {
			return fmt.Errorf("generator config %s is missing from the other report", name1)
		} else {
			if err := compareGeneratorConfigs(cfg1, cfg2); err != nil {
				return err
			}
		}
	}

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
