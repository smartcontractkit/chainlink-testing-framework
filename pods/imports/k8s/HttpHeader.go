package k8s


// HTTPHeader describes a custom header to be used in HTTP probes.
type HttpHeader struct {
	// The header field name.
	//
	// This will be canonicalized upon output, so case-variant names will be understood as the same header.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The header field value.
	Value *string `field:"required" json:"value" yaml:"value"`
}

