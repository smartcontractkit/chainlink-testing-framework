package k8s


// EndpointAddress is a tuple that describes single IP address.
type EndpointAddress struct {
	// The IP of this endpoint.
	//
	// May not be loopback (127.0.0.0/8 or ::1), link-local (169.254.0.0/16 or fe80::/10), or link-local multicast (224.0.0.0/24 or ff02::/16).
	Ip *string `field:"required" json:"ip" yaml:"ip"`
	// The Hostname of this endpoint.
	Hostname *string `field:"optional" json:"hostname" yaml:"hostname"`
	// Optional: Node hosting this endpoint.
	//
	// This can be used to determine endpoints local to a node.
	NodeName *string `field:"optional" json:"nodeName" yaml:"nodeName"`
	// Reference to object providing the endpoint.
	TargetRef *ObjectReference `field:"optional" json:"targetRef" yaml:"targetRef"`
}

