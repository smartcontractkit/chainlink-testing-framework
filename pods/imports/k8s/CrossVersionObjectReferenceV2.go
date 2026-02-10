package k8s


// CrossVersionObjectReference contains enough information to let you identify the referred resource.
type CrossVersionObjectReferenceV2 struct {
	// kind is the kind of the referent;
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `field:"required" json:"kind" yaml:"kind"`
	// name is the name of the referent;
	//
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name *string `field:"required" json:"name" yaml:"name"`
	// apiVersion is the API version of the referent.
	ApiVersion *string `field:"optional" json:"apiVersion" yaml:"apiVersion"`
}

