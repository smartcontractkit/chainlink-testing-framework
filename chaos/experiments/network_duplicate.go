package experiments

import "time"

// NetworkDuplicate struct contains objects for Network Duplication testing
type NetworkDuplicate struct {
	Base
	TargetAppLabel string
	Duplicate      int
	Correlation    int
	Duration       time.Duration
}

// SetBase sets the base
func (e *NetworkDuplicate) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *NetworkDuplicate) Resource() string {
	return "networkchaos"
}

// Filename returns the network duplicate yaml
func (e *NetworkDuplicate) Filename() string {
	return "network-duplicate.yml"
}
