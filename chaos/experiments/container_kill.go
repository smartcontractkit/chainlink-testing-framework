package experiments

type ContainerKill struct {
	Base
	TargetAppLabel string
	Container      string
}

func (e *ContainerKill) SetBase(base Base) {
	e.Base = base
}

func (e *ContainerKill) Resource() string {
	return "podchaos"
}

func (e *ContainerKill) Filename() string {
	return "container-kill.yml"
}
