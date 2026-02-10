package k8s


// IPBlock describes a particular CIDR (Ex.
//
// "192.168.1.0/24","2001:db8::/64") that is allowed to the pods matched by a NetworkPolicySpec's podSelector. The except entry describes CIDRs that should not be included within this rule.
type IpBlock struct {
	// cidr is a string representing the IPBlock Valid examples are "192.168.1.0/24" or "2001:db8::/64".
	Cidr *string `field:"required" json:"cidr" yaml:"cidr"`
	// except is a slice of CIDRs that should not be included within an IPBlock Valid examples are "192.168.1.0/24" or "2001:db8::/64" Except values will be rejected if they are outside the cidr range.
	Except *[]*string `field:"optional" json:"except" yaml:"except"`
}

