package experiments

import (
	"time"
)

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

func (e *NetworkBandwidth) SetBase(base Base) {
	e.Base = base
}

func (e *NetworkBandwidth) Resource() string {
	return "networkchaos"
}

func (e *NetworkBandwidth) Filename() string {
	return "network-bandwidth.yml"
}
