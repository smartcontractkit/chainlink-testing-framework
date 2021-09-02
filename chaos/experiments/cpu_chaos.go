package experiments

import "time"

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

func (e *CPUHog) SetBase(base Base) {
	e.Base = base
}

func (e *CPUHog) Resource() string {
	return "stresschaos"
}

func (e *CPUHog) Filename() string {
	return "cpu-chaos.yml"
}
