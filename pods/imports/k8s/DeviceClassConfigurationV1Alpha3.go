package k8s


// DeviceClassConfiguration is used in DeviceClass.
type DeviceClassConfigurationV1Alpha3 struct {
	// Opaque provides driver-specific configuration parameters.
	Opaque *OpaqueDeviceConfigurationV1Alpha3 `field:"optional" json:"opaque" yaml:"opaque"`
}

