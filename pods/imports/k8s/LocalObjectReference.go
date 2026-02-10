package k8s


// LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
type LocalObjectReference struct {
	// Name of the referent.
	//
	// This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name *string `field:"optional" json:"name" yaml:"name"`
}

