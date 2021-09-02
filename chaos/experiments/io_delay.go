package experiments

import "time"

type IODelay struct {
	Base
	TargetAppLabel string
	VolumePath     string
	Path           string
	Delay          time.Duration
	Percent        int
	Duration       time.Duration
}

func (e *IODelay) SetBase(base Base) {
	e.Base = base
}

func (e *IODelay) Resource() string {
	return "iochaos"
}

func (e *IODelay) Filename() string {
	return "io-delay.yml"
}
