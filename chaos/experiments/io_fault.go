package experiments

import "time"

type IOFault struct {
	Base
	TargetAppLabel string
	VolumePath     string
	Path           string
	Errno          int
	Percent        int
	Duration       time.Duration
}

func (e *IOFault) SetBase(base Base) {
	e.Base = base
}

func (e *IOFault) Resource() string {
	return "iochaos"
}

func (e *IOFault) Filename() string {
	return "io-fault.yml"
}
