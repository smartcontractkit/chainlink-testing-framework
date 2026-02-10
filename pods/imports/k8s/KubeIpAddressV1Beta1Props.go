package k8s


// IPAddress represents a single IP of a single IP Family.
//
// The object is designed to be used by APIs that operate on IP addresses. The object is used by the Service core API for allocation of IP addresses. An IP address can be represented in different formats, to guarantee the uniqueness of the IP, the name of the object is the IP address in canonical format, four decimal digits separated by dots suppressing leading zeros for IPv4 and the representation defined by RFC 5952 for IPv6. Valid: 192.168.1.5 or 2001:db8::1 or 2001:db8:aaaa:bbbb:cccc:dddd:eeee:1 Invalid: 10.01.2.3 or 2001:db8:0:0:0::1
type KubeIpAddressV1Beta1Props struct {
	// Standard object's metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// spec is the desired state of the IPAddress.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *IpAddressSpecV1Beta1 `field:"optional" json:"spec" yaml:"spec"`
}

