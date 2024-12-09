package wasp

import (
	"time"
)

/* Different load profile schedules definitions */

const (
	// DefaultStepChangePrecision is default amount of steps in which we split a schedule
	DefaultStepChangePrecision = 10
)

// Plain creates a slice containing a single Segment starting at `from` with the specified `duration`.
// It is used to initialize basic segments with defined timing.
func Plain(from int64, duration time.Duration) []*Segment {
	return []*Segment{
		{
			From:     from,
			Duration: duration,
			Type:     SegmentType_Plain,
		},
	}
}

// Steps generates a slice of Segment pointers starting from 'from', incremented by 'increase' for each of 'steps' steps.
// Each Segment has a duration equal to the total duration divided by the number of steps.
// Use this function to create uniformly distributed segments over a specified time period.
func Steps(from, increase int64, steps int, duration time.Duration) []*Segment {
	segments := make([]*Segment, 0)
	perStepDuration := duration / time.Duration(steps)
	for i := 0; i < steps; i++ {
		newFrom := from + int64(i)*increase
		segments = append(segments, &Segment{
			From:     newFrom,
			Duration: perStepDuration,
			Type:     SegmentType_Steps,
		})
	}
	return segments
}

// Combine merges multiple slices of Segment pointers into a single slice.
// It is useful for aggregating segment data from various sources.
func Combine(segs ...[]*Segment) []*Segment {
	acc := make([]*Segment, 0)
	for _, ss := range segs {
		acc = append(acc, ss...)
	}
	return acc
}

// CombineAndRepeat concatenates multiple Segment slices and repeats the combined sequence the specified number of times.
// It returns a single slice containing the repeated segments.
// Panics with ErrNoSchedule if no segments are provided.
func CombineAndRepeat(times int, segs ...[]*Segment) []*Segment {
	if len(segs) == 0 {
		panic(ErrNoSchedule)
	}
	acc := make([]*Segment, 0)
	for i := 0; i < times; i++ {
		for _, ss := range segs {
			acc = append(acc, ss...)
		}
	}
	return acc
}
