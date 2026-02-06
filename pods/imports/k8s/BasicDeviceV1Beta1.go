package k8s

// BasicDevice defines one device instance.
type BasicDeviceV1Beta1 struct {
	// Attributes defines the set of attributes for this device.
	//
	// The name of each attribute must be unique in that set.
	//
	// The maximum number of attributes and capacities combined is 32.
	Attributes *map[string]*DeviceAttributeV1Beta1 `field:"optional" json:"attributes" yaml:"attributes"`
	// Capacity defines the set of capacities for this device.
	//
	// The name of each capacity must be unique in that set.
	//
	// The maximum number of attributes and capacities combined is 32.
	Capacity *map[string]*DeviceCapacityV1Beta1 `field:"optional" json:"capacity" yaml:"capacity"`
}
