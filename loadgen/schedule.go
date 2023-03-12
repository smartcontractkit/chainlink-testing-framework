package loadgen

import (
	"time"

	"github.com/rs/zerolog/log"
)

/* Different load profile schedules definitions */

const (
	DefaultStepChangePrecision = 10
)

type SawScheduleProfile struct {
	From         int64
	Increase     int64
	Steps        int64
	StepDuration time.Duration
	Length       int
}

func HorizontalLine(from int64) []*Segment {
	return []*Segment{
		{
			From: from,
		},
	}
}

func Line(from, to int64, duration time.Duration) []*Segment {
	stepDur := duration / DefaultStepChangePrecision
	inc := (to - from) / DefaultStepChangePrecision
	log.Info().
		Dur("StepDur", stepDur).
		Int64("Increase", inc).
		Msg("Stats generated")
	return []*Segment{
		{
			From:         from,
			Steps:        DefaultStepChangePrecision,
			Increase:     inc,
			StepDuration: stepDur,
		},
	}
}

func Saw(prof SawScheduleProfile) []*Segment {
	segs := make([]*Segment, 0)
	for i := 0; i < prof.Length; i++ {
		s := &Segment{
			From:         prof.From,
			Steps:        prof.Steps,
			StepDuration: prof.StepDuration,
		}
		if i%2 == 0 {
			s.Increase = prof.Increase
		} else {
			s.From = prof.From + (prof.Increase * prof.Steps)
			s.Increase = -prof.Increase
		}
		segs = append(segs, s)
	}
	return segs
}
