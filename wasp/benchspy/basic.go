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

	for name2 := range otherData.GeneratorConfigs {
		if _, ok := b.GeneratorConfigs[name2]; !ok {
			return fmt.Errorf("generator config %s is missing from the current report", name2)
		}
	}

	// TODO: would be good to be able to check if Gun and VU are the same, but idk yet how we could do that easily [hash the code?]

	return nil
}

func compareGeneratorConfigs(cfg1, cfg2 *wasp.Config) error {
	if cfg1.LoadType != cfg2.LoadType {
		return fmt.Errorf("load types are different. Expected %s, got %s", cfg1.LoadType, cfg2.LoadType)
	}

	if len(cfg1.Schedule) != len(cfg2.Schedule) {
		return fmt.Errorf("schedules are different. Expected %d, got %d", len(cfg1.Schedule), len(cfg2.Schedule))
	}

	for i, segment1 := range cfg1.Schedule {
		segment2 := cfg2.Schedule[i]
		if segment1 == nil {
			return fmt.Errorf("schedule at index %d is nil in the current report", i)
		}
		if segment2 == nil {
			return fmt.Errorf("schedule at index %d is nil in the other report", i)
		}
		if *segment1 != *segment2 {
			return fmt.Errorf("schedules at index %d are different. Expected %s, got %s", i, mustMarshallSegment(segment1), mustMarshallSegment(segment2))
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
