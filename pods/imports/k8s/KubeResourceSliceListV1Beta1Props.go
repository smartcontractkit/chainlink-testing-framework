package k8s

// ResourceSliceList is a collection of ResourceSlices.
type KubeResourceSliceListV1Beta1Props struct {
	// Items is the list of resource ResourceSlices.
	Items *[]*KubeResourceSliceV1Beta1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}
