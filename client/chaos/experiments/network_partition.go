package experiments

// NetworkPartition struct with objects for NetworkConfig Partition testing
type NetworkPartition struct {
	Base
	FromMode       string
	FromLabelKey   string
	FromLabelValue string
	ToMode         string
	ToLabelKey     string
	ToLabelValue   string
}

// SetBase sets the base
func (e *NetworkPartition) SetBase(base Base) {
	e.Base = base
}

// Resource is a CRD resource that can be found in spec.names.singular
func (e *NetworkPartition) Resource() string {
	return "networkchaos"
}

// Filename returns the network partition yaml
func (e *NetworkPartition) Filename() string {
	return "network-partition.yml"
}
