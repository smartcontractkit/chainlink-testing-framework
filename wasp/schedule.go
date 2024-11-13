package wasp

import (
	"time"
)

/* Different load profile schedules definitions */

const (
	// DefaultStepChangePrecision is default amount of steps in which we split a schedule
	DefaultStepChangePrecision = 10
)

// Plain creates a slice containing a single Segment initialized with the specified 
// starting point 'from' and 'duration'. The Segment represents a time interval 
// starting from 'from' with the given duration. 
// The function returns a pointer to the slice of Segment.
func Plain(from int64, duration time.Duration) []*Segment {
	return []*Segment{
		{
			From:     from,
			Duration: duration,
		},
	}
}

// Steps generates a slice of Segment pointers based on the provided parameters. 
// It calculates the starting point for each segment by incrementing the 'from' value 
// by 'increase' for each step, and divides the total 'duration' evenly across the 
// specified number of 'steps'. Each Segment contains the calculated starting point 
// and its corresponding duration. The function returns a slice of these Segment pointers.
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

// Combine takes multiple slices of pointers to Segment and concatenates them into a single slice. 
// It returns a new slice containing all the segments from the provided slices in the order they were received. 
// If no segments are provided, it returns an empty slice.
func Combine(segs ...[]*Segment) []*Segment {
	acc := make([]*Segment, 0)
	for _, ss := range segs {
		acc = append(acc, ss...)
	}
	return acc
}

// CombineAndRepeat takes an integer 'times' and a variadic number of slice arguments containing pointers to Segment. 
// It concatenates the provided slices of Segment pointers and repeats this concatenation 'times' times. 
// The function returns a single slice of Segment pointers containing the combined segments. 
// If no segments are provided, it will panic with ErrNoSchedule.
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
