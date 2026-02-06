package k8s

// DeviceClassList is a collection of classes.
type KubeDeviceClassListV1Beta1Props struct {
	// Items is the list of resource classes.
	Items *[]*KubeDeviceClassV1Beta1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}
