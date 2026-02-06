package k8s

// ResourceClaimSpec defines what is being requested in a ResourceClaim and how to configure it.
type ResourceClaimSpecV1Alpha3 struct {
	// Devices defines how to request devices.
	Devices *DeviceClaimV1Alpha3 `field:"optional" json:"devices" yaml:"devices"`
}
