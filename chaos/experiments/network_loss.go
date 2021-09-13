package experiments

import "time"

// NetworkLoss struct with objects for Network Loss testing
type NetworkLoss struct {
	Base
	TargetAppLabel string
	Loss           int
	Correlation    int
	Duration       time.Duration
}

// SetBase sets the base
func (e *NetworkLoss) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *NetworkLoss) Resource() string {
	return "networkchaos"
}

// Filename returns the network loss yaml
func (e *NetworkLoss) Filename() string {
	return "network-loss.yml"
}
