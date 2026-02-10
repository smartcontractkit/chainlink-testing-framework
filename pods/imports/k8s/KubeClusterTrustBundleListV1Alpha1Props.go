package k8s


// ClusterTrustBundleList is a collection of ClusterTrustBundle objects.
type KubeClusterTrustBundleListV1Alpha1Props struct {
	// items is a collection of ClusterTrustBundle objects.
	Items *[]*KubeClusterTrustBundleV1Alpha1Props `field:"required" json:"items" yaml:"items"`
	// metadata contains the list metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

