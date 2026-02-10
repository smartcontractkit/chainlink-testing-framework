package k8s


// ResourceClaim describes a request for access to resources in the cluster, for use by workloads.
//
// For example, if a workload needs an accelerator device with specific properties, this is how that request is expressed. The status stanza tracks whether this claim has been satisfied and what specific resources have been allocated.
//
// This is an alpha type and requires enabling the DynamicResourceAllocation feature gate.
type KubeResourceClaimV1Alpha3Props struct {
	// Spec describes what is being requested and how to configure it.
	//
	// The spec is immutable.
	Spec *ResourceClaimSpecV1Alpha3 `field:"required" json:"spec" yaml:"spec"`
	// Standard object metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

