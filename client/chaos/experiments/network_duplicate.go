package experiments

import "time"

// NetworkDuplicate struct contains objects for NetworkConfig Duplication testing
type NetworkDuplicate struct {
	Base
	Mode        string
	LabelKey    string
	LabelValue  string
	Duplicate   int
	Correlation int
	Duration    time.Duration
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
