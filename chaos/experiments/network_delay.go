package experiments

import (
	"time"
)

type NetworkDelay struct {
	Base
	TargetAppLabel string
	Latency        time.Duration
	Duration       time.Duration
}

func (e *NetworkDelay) SetBase(base Base) {
	e.Base = base
}

func (e *NetworkDelay) Resource() string {
	return "networkchaos"
}

func (e *NetworkDelay) Filename() string {
	return "network-delay.yml"
}
