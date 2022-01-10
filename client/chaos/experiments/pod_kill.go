package experiments

// PodKill struct for pod kill testing
type PodKill struct {
	Base
	Mode       string
	LabelKey   string
	LabelValue string
}

// SetBase sets the base
func (e *PodKill) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *PodKill) Resource() string {
	return "podchaos"
}

// Filename returns the file for pod kill
func (e *PodKill) Filename() string {
	return "pod-kill.yml"
}
