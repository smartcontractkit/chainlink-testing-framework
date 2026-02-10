package k8s


// StatefulSetOrdinals describes the policy used for replica ordinal assignment in this StatefulSet.
type StatefulSetOrdinals struct {
	// start is the number representing the first replica's index.
	//
	// It may be used to number replicas from an alternate index (eg: 1-indexed) over the default 0-indexed names, or to orchestrate progressive movement of replicas from one StatefulSet to another. If set, replica indices will be in the range:
	// [.spec.ordinals.start, .spec.ordinals.start + .spec.replicas).
	// If unset, defaults to 0. Replica indices will be in the range:
	// [0, .spec.replicas).
	Start *float64 `field:"optional" json:"start" yaml:"start"`
}

