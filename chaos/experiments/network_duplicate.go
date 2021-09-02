package experiments

import "time"

type NetworkDuplicate struct {
	Base
	TargetAppLabel string
	Duplicate      int
	Correlation    int
	Duration       time.Duration
}

func (e *NetworkDuplicate) SetBase(base Base) {
	e.Base = base
}

func (e *NetworkDuplicate) Resource() string {
	return "networkchaos"
}

func (e *NetworkDuplicate) Filename() string {
	return "network-duplicate.yml"
}
