package k8s


// The names of the group, the version, and the resource.
type GroupVersionResourceV1Alpha1 struct {
	// The name of the group.
	Group *string `field:"optional" json:"group" yaml:"group"`
	// The name of the resource.
	Resource *string `field:"optional" json:"resource" yaml:"resource"`
	// The name of the version.
	Version *string `field:"optional" json:"version" yaml:"version"`
}

