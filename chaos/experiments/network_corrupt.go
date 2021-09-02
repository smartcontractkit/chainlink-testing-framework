package experiments

import (
	"time"
)

type NetworkCorrupt struct {
	Base
	TargetAppLabel string
	Corrupt        int
	Correlation    int
	Duration       time.Duration
}

func (e *NetworkCorrupt) SetBase(base Base) {
	e.Base = base
}

func (e *NetworkCorrupt) Resource() string {
	return "networkchaos"
}

func (e *NetworkCorrupt) Filename() string {
	return "network-delay.yml"
}
