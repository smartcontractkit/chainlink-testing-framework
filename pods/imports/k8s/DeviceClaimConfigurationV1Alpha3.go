package k8s


// DeviceClaimConfiguration is used for configuration parameters in DeviceClaim.
type DeviceClaimConfigurationV1Alpha3 struct {
	// Opaque provides driver-specific configuration parameters.
	Opaque *OpaqueDeviceConfigurationV1Alpha3 `field:"optional" json:"opaque" yaml:"opaque"`
	// Requests lists the names of requests where the configuration applies.
	//
	// If empty, it applies to all requests.
	Requests *[]*string `field:"optional" json:"requests" yaml:"requests"`
}

