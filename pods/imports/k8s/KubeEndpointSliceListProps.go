package k8s


// EndpointSliceList represents a list of endpoint slices.
type KubeEndpointSliceListProps struct {
	// items is the list of endpoint slices.
	Items *[]*KubeEndpointSliceProps `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

