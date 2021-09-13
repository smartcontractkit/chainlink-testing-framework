package experiments

import "time"

// CPUHog struct for cpu hog testing
type CPUHog struct {
	Base
	TargetAppLabel string
	Workers        int
	Load           int
	OptsCPU        int
	OptsTimeout    int
	OptsHDD        int
	Duration       time.Duration
}

// SetBase sets the base
func (e *CPUHog) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *CPUHog) Resource() string {
	return "stresschaos"
}

// Filename returns the cpu chaos yaml
func (e *CPUHog) Filename() string {
	return "cpu-chaos.yml"
}
