package k8s


// CSIStorageCapacityList is a collection of CSIStorageCapacity objects.
type KubeCsiStorageCapacityListProps struct {
	// items is the list of CSIStorageCapacity objects.
	Items *[]*KubeCsiStorageCapacityProps `field:"required" json:"items" yaml:"items"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

