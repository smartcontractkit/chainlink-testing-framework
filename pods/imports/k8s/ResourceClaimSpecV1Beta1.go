package k8s

// ResourceClaimSpec defines what is being requested in a ResourceClaim and how to configure it.
type ResourceClaimSpecV1Beta1 struct {
	// Devices defines how to request devices.
	Devices *DeviceClaimV1Beta1 `field:"optional" json:"devices" yaml:"devices"`
}
