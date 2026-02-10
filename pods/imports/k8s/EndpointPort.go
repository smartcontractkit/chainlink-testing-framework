package k8s


// EndpointPort is a tuple that describes a single port.
type EndpointPort struct {
	// The port number of the endpoint.
	Port *float64 `field:"required" json:"port" yaml:"port"`
	// The application protocol for this port.
	//
	// This is used as a hint for implementations to offer richer behavior for protocols that they understand. This field follows standard Kubernetes label syntax. Valid values are either:
	//
	// * Un-prefixed protocol names - reserved for IANA standard service names (as per RFC-6335 and https://www.iana.org/assignments/service-names).
	//
	// * Kubernetes-defined prefixed names:
	// * 'kubernetes.io/h2c' - HTTP/2 prior knowledge over cleartext as described in https://www.rfc-editor.org/rfc/rfc9113.html#name-starting-http-2-with-prior-
	// * 'kubernetes.io/ws'  - WebSocket over cleartext as described in https://www.rfc-editor.org/rfc/rfc6455
	// * 'kubernetes.io/wss' - WebSocket over TLS as described in https://www.rfc-editor.org/rfc/rfc6455
	//
	// * Other protocols should use implementation-defined prefixed names such as mycompany.com/my-custom-protocol.
	AppProtocol *string `field:"optional" json:"appProtocol" yaml:"appProtocol"`
	// The name of this port.
	//
	// This must match the 'name' field in the corresponding ServicePort. Must be a DNS_LABEL. Optional only if one port is defined.
	Name *string `field:"optional" json:"name" yaml:"name"`
	// The IP protocol for this port.
	//
	// Must be UDP, TCP, or SCTP. Default is TCP.
	// Default: TCP.
	//
	Protocol *string `field:"optional" json:"protocol" yaml:"protocol"`
}

