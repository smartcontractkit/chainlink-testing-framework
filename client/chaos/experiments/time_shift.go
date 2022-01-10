package experiments

import "time"

// TimeShift struct to contain info needed for TimeShift testing
type TimeShift struct {
	Base
	Mode       string
	LabelKey   string
	LabelValue string
	TimeOffset time.Duration
	Duration   time.Duration
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
