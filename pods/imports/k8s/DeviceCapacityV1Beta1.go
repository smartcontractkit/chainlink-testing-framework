package k8s

// DeviceCapacity describes a quantity associated with a device.
type DeviceCapacityV1Beta1 struct {
	// Value defines how much of a certain device capacity is available.
	Value Quantity `field:"required" json:"value" yaml:"value"`
}
