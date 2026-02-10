package k8s


// ServiceBackendPort is the service port being referenced.
type ServiceBackendPort struct {
	// name is the name of the port on the Service.
	//
	// This is a mutually exclusive setting with "Number".
	Name *string `field:"optional" json:"name" yaml:"name"`
	// number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with "Name".
	Number *float64 `field:"optional" json:"number" yaml:"number"`
}

