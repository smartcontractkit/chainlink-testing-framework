package k8s


// DeviceClass is a vendor- or admin-provided resource that contains device configuration and selectors.
//
// It can be referenced in the device requests of a claim to apply these presets. Cluster scoped.
//
// This is an alpha type and requires enabling the DynamicResourceAllocation feature gate.
type KubeDeviceClassV1Alpha3Props struct {
	// Spec defines what can be allocated and how to configure it.
	//
	// This is mutable. Consumers have to be prepared for classes changing at any time, either because they get updated or replaced. Claim allocations are done once based on whatever was set in classes at the time of allocation.
	//
	// Changing the spec automatically increments the metadata.generation number.
	Spec *DeviceClassSpecV1Alpha3 `field:"required" json:"spec" yaml:"spec"`
	// Standard object metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

