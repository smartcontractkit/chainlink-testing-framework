package k8s


// ServiceCIDR defines a range of IP addresses using CIDR format (e.g. 192.168.0.0/24 or 2001:db2::/64). This range is used to allocate ClusterIPs to Service objects.
type KubeServiceCidrv1Beta1Props struct {
	// Standard object's metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// spec is the desired state of the ServiceCIDR.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *ServiceCidrSpecV1Beta1 `field:"optional" json:"spec" yaml:"spec"`
}

