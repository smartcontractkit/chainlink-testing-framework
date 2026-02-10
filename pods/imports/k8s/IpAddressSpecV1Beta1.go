package k8s


// IPAddressSpec describe the attributes in an IP Address.
type IpAddressSpecV1Beta1 struct {
	// ParentRef references the resource that an IPAddress is attached to.
	//
	// An IPAddress must reference a parent object.
	ParentRef *ParentReferenceV1Beta1 `field:"required" json:"parentRef" yaml:"parentRef"`
}

