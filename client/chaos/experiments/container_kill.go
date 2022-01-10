package experiments

// ContainerKill struct for continer kill testing
type ContainerKill struct {
	Base
	Mode       string
	LabelKey   string
	LabelValue string
	Container  string
}

// SetBase sets the base
func (e *ContainerKill) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource
func (e *ContainerKill) Resource() string {
	return "podchaos"
}

// Filename returns the filename for container kill
func (e *ContainerKill) Filename() string {
	return "container-kill.yml"
}
