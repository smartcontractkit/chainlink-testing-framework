package experiments

import (
	"time"
)

// NetworkDelay stuct containing definitions for a network delay
type NetworkDelay struct {
	Base
	Mode       string
	LabelKey   string
	LabelValue string
	Latency    time.Duration
	Duration   time.Duration
}

// SetBase sets the base
func (e *NetworkDelay) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *NetworkDelay) Resource() string {
	return "networkchaos"
}

// Filename returns the file name for network delay
func (e *NetworkDelay) Filename() string {
	return "network-delay.yml"
}
