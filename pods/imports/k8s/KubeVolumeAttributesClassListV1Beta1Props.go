package k8s


// VolumeAttributesClassList is a collection of VolumeAttributesClass objects.
type KubeVolumeAttributesClassListV1Beta1Props struct {
	// items is the list of VolumeAttributesClass objects.
	Items *[]*KubeVolumeAttributesClassV1Beta1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

