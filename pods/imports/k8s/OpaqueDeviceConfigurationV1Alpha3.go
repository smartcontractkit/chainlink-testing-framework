package k8s


// OpaqueDeviceConfiguration contains configuration parameters for a driver in a format defined by the driver vendor.
type OpaqueDeviceConfigurationV1Alpha3 struct {
	// Driver is used to determine which kubelet plugin needs to be passed these configuration parameters.
	//
	// An admission policy provided by the driver developer could use this to decide whether it needs to validate them.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
	Driver *string `field:"required" json:"driver" yaml:"driver"`
	// Parameters can contain arbitrary data.
	//
	// It is the responsibility of the driver developer to handle validation and versioning. Typically this includes self-identification and a version ("kind" + "apiVersion" for Kubernetes types), with conversion between different versions.
	Parameters interface{} `field:"required" json:"parameters" yaml:"parameters"`
}

