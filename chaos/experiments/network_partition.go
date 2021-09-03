package experiments

import (
	"time"
)

type NetworkPartition struct {
	Base
	FromMode       string
	FromLabelKey   string
	FromLabelValue string
	ToMode         string
	ToLabelKey     string
	ToLabelValue   string
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
