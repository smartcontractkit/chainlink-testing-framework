package k8s

// DeviceClaim defines how to request devices with a ResourceClaim.
type DeviceClaimV1Beta1 struct {
	// This field holds configuration for multiple potential drivers which could satisfy requests in this claim.
	//
	// It is ignored while allocating the claim.
	Config *[]*DeviceClaimConfigurationV1Beta1 `field:"optional" json:"config" yaml:"config"`
	// These constraints must be satisfied by the set of devices that get allocated for the claim.
	Constraints *[]*DeviceConstraintV1Beta1 `field:"optional" json:"constraints" yaml:"constraints"`
	// Requests represent individual requests for distinct devices which must all be satisfied.
	//
	// If empty, nothing needs to be allocated.
	Requests *[]*DeviceRequestV1Beta1 `field:"optional" json:"requests" yaml:"requests"`
}
