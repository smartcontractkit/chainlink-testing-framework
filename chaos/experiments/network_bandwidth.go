package experiments

import (
	"time"
)

// NetworkBandwidth struct with objects for Network Bandwidth testing
type NetworkBandwidth struct {
	Base
	TargetAppLabel string
	// kbps
	Rate     string
	Limit    int
	Buffer   int
	PeakRate int
	MinBurst int
	Duration time.Duration
}

// SetBase sets the base
func (e *NetworkBandwidth) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *NetworkBandwidth) Resource() string {
	return "networkchaos"
}

// Filename returns the file name of the network bandwidth yaml
func (e *NetworkBandwidth) Filename() string {
	return "network-bandwidth.yml"
}
