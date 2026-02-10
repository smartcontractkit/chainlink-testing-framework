package k8s


// Device represents one individual hardware instance that can be selected based on its attributes.
//
// Besides the name, exactly one field must be set.
type DeviceV1Alpha3 struct {
	// Name is unique identifier among all devices managed by the driver in the pool.
	//
	// It must be a DNS label.
	Name *string `field:"required" json:"name" yaml:"name"`
	// Basic defines one device instance.
	Basic *BasicDeviceV1Alpha3 `field:"optional" json:"basic" yaml:"basic"`
}

