package k8s


// ResourceClaimTemplateSpec contains the metadata and fields for a ResourceClaim.
type ResourceClaimTemplateSpecV1Alpha3 struct {
	// Spec for the ResourceClaim.
	//
	// The entire content is copied unchanged into the ResourceClaim that gets created from this template. The same fields as in a ResourceClaim are also valid here.
	Spec *ResourceClaimSpecV1Alpha3 `field:"required" json:"spec" yaml:"spec"`
	// ObjectMeta may contain labels and annotations that will be copied into the PVC when creating it.
	//
	// No other fields are allowed and will be rejected during validation.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

