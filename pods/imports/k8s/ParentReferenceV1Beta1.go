package k8s


// ParentReference describes a reference to a parent object.
type ParentReferenceV1Beta1 struct {
	// Name is the name of the object being referenced.
	Name *string `field:"required" json:"name" yaml:"name"`
	// Resource is the resource of the object being referenced.
	Resource *string `field:"required" json:"resource" yaml:"resource"`
	// Group is the group of the object being referenced.
	Group *string `field:"optional" json:"group" yaml:"group"`
	// Namespace is the namespace of the object being referenced.
	Namespace *string `field:"optional" json:"namespace" yaml:"namespace"`
}

