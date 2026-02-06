package k8s

// DeviceClassConfiguration is used in DeviceClass.
type DeviceClassConfigurationV1Beta1 struct {
	// Opaque provides driver-specific configuration parameters.
	Opaque *OpaqueDeviceConfigurationV1Beta1 `field:"optional" json:"opaque" yaml:"opaque"`
}
