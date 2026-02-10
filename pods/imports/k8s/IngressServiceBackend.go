package k8s


// IngressServiceBackend references a Kubernetes Service as a Backend.
type IngressServiceBackend struct {
	// name is the referenced service.
	//
	// The service must exist in the same namespace as the Ingress object.
	Name *string `field:"required" json:"name" yaml:"name"`
	// port of the referenced service.
	//
	// A port name or port number is required for a IngressServiceBackend.
	Port *ServiceBackendPort `field:"optional" json:"port" yaml:"port"`
}

