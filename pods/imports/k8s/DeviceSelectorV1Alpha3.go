package k8s


// DeviceSelector must have exactly one field set.
type DeviceSelectorV1Alpha3 struct {
	// CEL contains a CEL expression for selecting a device.
	Cel *CelDeviceSelectorV1Alpha3 `field:"optional" json:"cel" yaml:"cel"`
}

