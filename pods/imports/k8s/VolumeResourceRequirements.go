package k8s


// VolumeResourceRequirements describes the storage resource requirements for a volume.
type VolumeResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	//
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Limits *map[string]Quantity `field:"optional" json:"limits" yaml:"limits"`
	// Requests describes the minimum amount of compute resources required.
	//
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Requests *map[string]Quantity `field:"optional" json:"requests" yaml:"requests"`
}

