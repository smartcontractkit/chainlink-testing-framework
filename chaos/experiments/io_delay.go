package experiments

import "time"

// IODelay struct contains objects for IODelay testing
type IODelay struct {
	Base
	TargetAppLabel string
	VolumePath     string
	Path           string
	Delay          time.Duration
	Percent        int
	Duration       time.Duration
}

// SetBase sets the base
func (e *IODelay) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *IODelay) Resource() string {
	return "iochaos"
}

// Filename returns the io delay yaml
func (e *IODelay) Filename() string {
	return "io-delay.yml"
}
