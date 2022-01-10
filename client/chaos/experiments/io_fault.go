package experiments

import "time"

// IOFault struct contains objects for IO Fault testing
type IOFault struct {
	Base
	Mode       string
	LabelKey   string
	LabelValue string
	VolumePath string
	Path       string
	Errno      int
	Percent    int
	Duration   time.Duration
}

// SetBase sets the base
func (e *IOFault) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *IOFault) Resource() string {
	return "iochaos"
}

// Filename returns the io fault yaml
func (e *IOFault) Filename() string {
	return "io-fault.yml"
}
