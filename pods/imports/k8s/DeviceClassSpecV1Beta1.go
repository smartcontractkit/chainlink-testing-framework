package k8s

// DeviceClassSpec is used in a [DeviceClass] to define what can be allocated and how to configure it.
type DeviceClassSpecV1Beta1 struct {
	// Config defines configuration parameters that apply to each device that is claimed via this class.
	//
	// Some classses may potentially be satisfied by multiple drivers, so each instance of a vendor configuration applies to exactly one driver.
	//
	// They are passed to the driver, but are not considered while allocating the claim.
	Config *[]*DeviceClassConfigurationV1Beta1 `field:"optional" json:"config" yaml:"config"`
	// Each selector must be satisfied by a device which is claimed via this class.
	Selectors *[]*DeviceSelectorV1Beta1 `field:"optional" json:"selectors" yaml:"selectors"`
}
