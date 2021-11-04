package experiments

import (
	"time"
)

// PodFailure struct contains objects for Pod Failure testing
type PodFailure struct {
	Base
	LabelKey   string
	LabelValue string
	Duration   time.Duration
}

// SetBase sets the base
func (e *PodFailure) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *PodFailure) Resource() string {
	return "podchaos"
}

// Filename returns the pod failure yaml
func (e *PodFailure) Filename() string {
	return "pod-failure.yml"
}
