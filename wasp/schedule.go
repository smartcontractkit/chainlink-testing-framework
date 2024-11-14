package wasp

import (
	"time"
)

/* Different load profile schedules definitions */

const (
	// DefaultStepChangePrecision is default amount of steps in which we split a schedule
	DefaultStepChangePrecision = 10
)

// Plain returns a slice of Segment containing a single Segment initialized with the provided from time and duration.
func Plain(from int64, duration time.Duration) []*Segment {
	return []*Segment{
		{
			From:     from,
			Duration: duration,
		},
	}
}

// Steps generates a slice of Segment pointers, each starting from 'from' and incremented by 'increase'.
// It creates 'steps' number of segments, dividing the total 'duration' equally across each segment.
// The function returns the resulting slice of segments.
func Steps(from, increase int64, steps int, duration time.Duration) []*Segment {
	segments := make([]*Segment, 0)
	perStepDuration := duration / time.Duration(steps)
	for i := 0; i < steps; i++ {
		newFrom := from + int64(i)*increase
		segments = append(segments, &Segment{
			From:     newFrom,
			Duration: perStepDuration,
		})
	}
	return segments
}

// Combine merges multiple slices of *Segment into a single slice.
// It accepts a variable number of []*Segment and appends them in order.
// The returned slice contains all segments from the provided slices.
func Combine(segs ...[]*Segment) []*Segment {
	acc := make([]*Segment, 0)
	for _, ss := range segs {
		acc = append(acc, ss...)
	}
	return acc
}

// CombineAndRepeat combines multiple Segment slices and repeats the combined sequence the specified number of times.
// It panics with ErrNoSchedule if no segment slices are provided.
// The returned slice contains all segments from each input slice, repeated `times` times.
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
