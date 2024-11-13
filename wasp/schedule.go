package wasp

import (
	"time"
)

/* Different load profile schedules definitions */

const (
	// DefaultStepChangePrecision is default amount of steps in which we split a schedule
	DefaultStepChangePrecision = 10
)

// Plain creates a slice containing a single Segment with the specified start time and duration.
// It returns a slice of pointers to Segment, initialized with the provided 'from' time and 'duration'.
func Plain(from int64, duration time.Duration) []*Segment {
	return []*Segment{
		{
			From:     from,
			Duration: duration,
		},
	}
}

// Steps generates a slice of Segment pointers, each representing a step in a sequence.
// It starts from the given 'from' value, increasing by 'increase' for each step.
// The total number of steps is specified by 'steps', and each step lasts for 'duration' divided by 'steps'.
// It returns a slice of Segment pointers, each containing the starting value and duration for that step.
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

// Combine concatenates multiple slices of Segment pointers into a single slice.
// It takes a variadic number of slices as input and returns a new slice containing
// all the elements from the input slices in the order they were provided.
func Combine(segs ...[]*Segment) []*Segment {
	acc := make([]*Segment, 0)
	for _, ss := range segs {
		acc = append(acc, ss...)
	}
	return acc
}

// CombineAndRepeat concatenates multiple slices of Segment pointers and repeats the concatenation a specified number of times.
// It panics if no slices are provided. The function returns a single slice containing the repeated concatenations.
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
