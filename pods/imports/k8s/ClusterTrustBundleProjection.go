package k8s


// ClusterTrustBundleProjection describes how to select a set of ClusterTrustBundle objects and project their contents into the pod filesystem.
type ClusterTrustBundleProjection struct {
	// Relative path from the volume root to write the bundle.
	Path *string `field:"required" json:"path" yaml:"path"`
	// Select all ClusterTrustBundles that match this label selector.
	//
	// Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as "match nothing".  If set but empty, interpreted as "match everything".
	LabelSelector *LabelSelector `field:"optional" json:"labelSelector" yaml:"labelSelector"`
	// Select a single ClusterTrustBundle by object name.
	//
	// Mutually-exclusive with signerName and labelSelector.
	Name *string `field:"optional" json:"name" yaml:"name"`
	// If true, don't block pod startup if the referenced ClusterTrustBundle(s) aren't available.
	//
	// If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
	Optional *bool `field:"optional" json:"optional" yaml:"optional"`
	// Select all ClusterTrustBundles that match this signer name.
	//
	// Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
	SignerName *string `field:"optional" json:"signerName" yaml:"signerName"`
}

