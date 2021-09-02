package experiments

import "time"

type NetworkLoss struct {
	Base
	TargetAppLabel string
	Loss           int
	Correlation    int
	Duration       time.Duration
}

func (e *NetworkLoss) SetBase(base Base) {
	e.Base = base
}

func (e *NetworkLoss) Resource() string {
	return "networkchaos"
}

func (e *NetworkLoss) Filename() string {
	return "network-loss.yml"
}
