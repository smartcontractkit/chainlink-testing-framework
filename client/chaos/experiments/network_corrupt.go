package experiments

import (
	"time"
)

// NetworkCorrupt struct for network corruption
type NetworkCorrupt struct {
	Base
	TargetAppLabel string
	Corrupt        int
	Correlation    int
	Duration       time.Duration
}

// SetBase sets the base
func (e *NetworkCorrupt) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *NetworkCorrupt) Resource() string {
	return "networkchaos"
}

// Filename returns the filename for a network corruption
func (e *NetworkCorrupt) Filename() string {
	return "network-corrupt.yml"
}
