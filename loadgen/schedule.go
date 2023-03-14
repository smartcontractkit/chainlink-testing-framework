package loadgen

import (
	"math"
	"time"
)

/* Different load profile schedules definitions */

const (
	// DefaultStepChangePrecision is default amount of steps in which we split a schedule
	DefaultStepChangePrecision = 10
)

func Plain(from int64, duration time.Duration) []*Segment {
	return []*Segment{
		{
			From:         from,
			Steps:        DefaultStepChangePrecision,
			StepDuration: duration / DefaultStepChangePrecision,
		},
	}
}

func Line(from, to int64, duration time.Duration) []*Segment {
	var inc int64
	stepDur := duration / DefaultStepChangePrecision
	incFloat := (float64(to) - float64(from)) / DefaultStepChangePrecision
	if math.Signbit(incFloat) {
		inc = int64(math.Floor(incFloat))
	} else {
		inc = int64(math.Ceil(incFloat))
	}
	return []*Segment{
		{
			From:         from,
			Steps:        DefaultStepChangePrecision,
			Increase:     inc,
			StepDuration: stepDur,
		},
	}
}

func Combine(segs ...[]*Segment) []*Segment {
	acc := make([]*Segment, 0)
	for _, ss := range segs {
		acc = append(acc, ss...)
	}
	return acc
}

func CombineAndRepeat(times int, segs ...[]*Segment) []*Segment {
	if len(segs) == 0 {
		panic(ErrNoSched)
	}
	acc := make([]*Segment, 0)
	for i := 0; i < times; i++ {
		for _, ss := range segs {
			acc = append(acc, ss...)
		}
	}
	return acc
}
