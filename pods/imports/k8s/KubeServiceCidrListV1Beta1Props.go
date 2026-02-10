package k8s


// ServiceCIDRList contains a list of ServiceCIDR objects.
type KubeServiceCidrListV1Beta1Props struct {
	// items is the list of ServiceCIDRs.
	Items *[]*KubeServiceCidrv1Beta1Props `field:"required" json:"items" yaml:"items"`
	// Standard object's metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

