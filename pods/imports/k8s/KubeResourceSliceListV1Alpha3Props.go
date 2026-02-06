package k8s

// ResourceSliceList is a collection of ResourceSlices.
type KubeResourceSliceListV1Alpha3Props struct {
	// Items is the list of resource ResourceSlices.
	Items *[]*KubeResourceSliceV1Alpha3Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}
