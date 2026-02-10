package k8s


// ResourceClaimTemplate is used to produce ResourceClaim objects.
//
// This is an alpha type and requires enabling the DynamicResourceAllocation feature gate.
type KubeResourceClaimTemplateV1Alpha3Props struct {
	// Describes the ResourceClaim that is to be generated.
	//
	// This field is immutable. A ResourceClaim will get created by the control plane for a Pod when needed and then not get updated anymore.
	Spec *ResourceClaimTemplateSpecV1Alpha3 `field:"required" json:"spec" yaml:"spec"`
	// Standard object metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

