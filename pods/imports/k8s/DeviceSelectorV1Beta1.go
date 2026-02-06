package k8s

// DeviceSelector must have exactly one field set.
type DeviceSelectorV1Beta1 struct {
	// CEL contains a CEL expression for selecting a device.
	Cel *CelDeviceSelectorV1Beta1 `field:"optional" json:"cel" yaml:"cel"`
}
