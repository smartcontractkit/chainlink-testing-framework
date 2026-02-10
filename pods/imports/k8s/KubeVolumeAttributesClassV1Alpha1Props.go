package k8s


// VolumeAttributesClass represents a specification of mutable volume attributes defined by the CSI driver.
//
// The class can be specified during dynamic provisioning of PersistentVolumeClaims, and changed in the PersistentVolumeClaim spec after provisioning.
type KubeVolumeAttributesClassV1Alpha1Props struct {
	// Name of the CSI driver This field is immutable.
	DriverName *string `field:"required" json:"driverName" yaml:"driverName"`
	// Standard object's metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// parameters hold volume attributes defined by the CSI driver.
	//
	// These values are opaque to the Kubernetes and are passed directly to the CSI driver. The underlying storage provider supports changing these attributes on an existing volume, however the parameters field itself is immutable. To invoke a volume update, a new VolumeAttributesClass should be created with new parameters, and the PersistentVolumeClaim should be updated to reference the new VolumeAttributesClass.
	//
	// This field is required and must contain at least one key/value pair. The keys cannot be empty, and the maximum number of parameters is 512, with a cumulative max size of 256K. If the CSI driver rejects invalid parameters, the target PersistentVolumeClaim will be set to an "Infeasible" state in the modifyVolumeStatus field.
	Parameters *map[string]*string `field:"optional" json:"parameters" yaml:"parameters"`
}

