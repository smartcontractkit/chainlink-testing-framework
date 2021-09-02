package experiments

import (
	"time"
)

type PodFailure struct {
	Base
	TargetAppLabel string
	Duration       time.Duration
}

func (e *PodFailure) SetBase(base Base) {
	e.Base = base
}

func (e *PodFailure) Resource() string {
	return "podchaos"
}

func (e *PodFailure) Filename() string {
	return "pod-failure.yml"
}
