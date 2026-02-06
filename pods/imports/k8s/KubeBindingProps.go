package k8s

// Binding ties one object to another;
//
// for example, a pod is bound to a node by a scheduler.
type KubeBindingProps struct {
	// The target object that you want to bind to the standard object.
	Target *ObjectReference `field:"required" json:"target" yaml:"target"`
	// Standard object's metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
}
