package experiments

import "time"

type TimeShift struct {
	Base
	TargetAppLabel string
	TimeOffset     time.Duration
	Duration       time.Duration
}

func (e *TimeShift) SetBase(base Base) {
	e.Base = base
}

func (e *TimeShift) Resource() string {
	return "timechaos"
}

func (e *TimeShift) Filename() string {
	return "time-shift.yml"
}
