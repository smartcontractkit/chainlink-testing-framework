package experiments

import "time"

// TimeShift stuct to contain info needed for TimeShift testing
type TimeShift struct {
	Base
	TargetAppLabel string
	TimeOffset     time.Duration
	Duration       time.Duration
}

// SetBase sets the base
func (e *TimeShift) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *TimeShift) Resource() string {
	return "timechaos"
}

// Filename returns the filename for a time shift
func (e *TimeShift) Filename() string {
	return "time-shift.yml"
}
