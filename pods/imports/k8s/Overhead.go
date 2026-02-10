package k8s


// Overhead structure represents the resource overhead associated with running a pod.
type Overhead struct {
	// podFixed represents the fixed resource overhead associated with running a pod.
	PodFixed *map[string]Quantity `field:"optional" json:"podFixed" yaml:"podFixed"`
}

