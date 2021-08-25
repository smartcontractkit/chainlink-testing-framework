package experiments

type PodKill struct {
	Base
	TargetAppLabel string
}

func (e *PodKill) SetBase(base Base) {
	e.Base = base
}

func (e *PodKill) Resource() string {
	return "podchaos"
}

func (e *PodKill) Filename() string {
	return "pod-kill.yml"
}
