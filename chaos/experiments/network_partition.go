package experiments

import (
	"time"
)

type NetworkPartition struct {
	Base
	AppLabel       string
	TargetAppLabel string
	Duration       time.Duration
}

func (e *NetworkPartition) SetBase(base Base) {
	e.Base = base
}

// Resource is a CRD resource that can be found in spec.names.singular
func (e *NetworkPartition) Resource() string {
	return "networkchaos"
}

func (e *NetworkPartition) Filename() string {
	return "network-partition.yml"
}
