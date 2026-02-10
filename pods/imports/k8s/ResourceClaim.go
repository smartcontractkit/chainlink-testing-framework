package k8s


// ResourceClaim references one entry in PodSpec.ResourceClaims.
type ResourceClaim struct {
	// Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
	Name *string `field:"required" json:"name" yaml:"name"`
	// Request is the name chosen for a request in the referenced claim.
	//
	// If empty, everything from the claim is made available, otherwise only the result of this request.
	Request *string `field:"optional" json:"request" yaml:"request"`
}

