package k8s

// PersistentVolumeSpec is the specification of a persistent volume.
type PersistentVolumeSpec struct {
	// accessModes contains all ways the volume can be mounted.
	//
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes
	AccessModes *[]*string `field:"optional" json:"accessModes" yaml:"accessModes"`
	// awsElasticBlockStore represents an AWS Disk resource that is attached to a kubelet's host machine and then exposed to the pod.
	//
	// Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	AwsElasticBlockStore *AwsElasticBlockStoreVolumeSource `field:"optional" json:"awsElasticBlockStore" yaml:"awsElasticBlockStore"`
	// azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
	//
	// Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type are redirected to the disk.csi.azure.com CSI driver.
	AzureDisk *AzureDiskVolumeSource `field:"optional" json:"azureDisk" yaml:"azureDisk"`
	// azureFile represents an Azure File Service mount on the host and bind mount to the pod.
	//
	// Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type are redirected to the file.csi.azure.com CSI driver.
	AzureFile *AzureFilePersistentVolumeSource `field:"optional" json:"azureFile" yaml:"azureFile"`
	// capacity is the description of the persistent volume's resources and capacity.
	//
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
	Capacity *map[string]Quantity `field:"optional" json:"capacity" yaml:"capacity"`
	// cephFS represents a Ceph FS mount on the host that shares a pod's lifetime.
	//
	// Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.
	Cephfs *CephFsPersistentVolumeSource `field:"optional" json:"cephfs" yaml:"cephfs"`
	// cinder represents a cinder volume attached and mounted on kubelets host machine.
	//
	// Deprecated: Cinder is deprecated. All operations for the in-tree cinder type are redirected to the cinder.csi.openstack.org CSI driver. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	Cinder *CinderPersistentVolumeSource `field:"optional" json:"cinder" yaml:"cinder"`
	// claimRef is part of a bi-directional binding between PersistentVolume and PersistentVolumeClaim.
	//
	// Expected to be non-nil when bound. claim.VolumeName is the authoritative bind between PV and PVC. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#binding
	ClaimRef *ObjectReference `field:"optional" json:"claimRef" yaml:"claimRef"`
	// csi represents storage that is handled by an external CSI driver.
	Csi *CsiPersistentVolumeSource `field:"optional" json:"csi" yaml:"csi"`
	// fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod.
	Fc *FcVolumeSource `field:"optional" json:"fc" yaml:"fc"`
	// flexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin.
	//
	// Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.
	FlexVolume *FlexPersistentVolumeSource `field:"optional" json:"flexVolume" yaml:"flexVolume"`
	// flocker represents a Flocker volume attached to a kubelet's host machine and exposed to the pod for its usage.
	//
	// This depends on the Flocker control service being running. Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.
	Flocker *FlockerVolumeSource `field:"optional" json:"flocker" yaml:"flocker"`
	// gcePersistentDisk represents a GCE Disk resource that is attached to a kubelet's host machine and then exposed to the pod.
	//
	// Provisioned by an admin. Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	GcePersistentDisk *GcePersistentDiskVolumeSource `field:"optional" json:"gcePersistentDisk" yaml:"gcePersistentDisk"`
	// glusterfs represents a Glusterfs volume that is attached to a host and exposed to the pod.
	//
	// Provisioned by an admin. Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported. More info: https://examples.k8s.io/volumes/glusterfs/README.md
	Glusterfs *GlusterfsPersistentVolumeSource `field:"optional" json:"glusterfs" yaml:"glusterfs"`
	// hostPath represents a directory on the host.
	//
	// Provisioned by a developer or tester. This is useful for single-node development and testing only! On-host storage is not supported in any way and WILL NOT WORK in a multi-node cluster. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	HostPath *HostPathVolumeSource `field:"optional" json:"hostPath" yaml:"hostPath"`
	// iscsi represents an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod.
	//
	// Provisioned by an admin.
	Iscsi *IscsiPersistentVolumeSource `field:"optional" json:"iscsi" yaml:"iscsi"`
	// local represents directly-attached storage with node affinity.
	Local *LocalVolumeSource `field:"optional" json:"local" yaml:"local"`
	// mountOptions is the list of mount options, e.g. ["ro", "soft"]. Not validated - mount will simply fail if one is invalid. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options.
	MountOptions *[]*string `field:"optional" json:"mountOptions" yaml:"mountOptions"`
	// nfs represents an NFS mount on the host.
	//
	// Provisioned by an admin. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	Nfs *NfsVolumeSource `field:"optional" json:"nfs" yaml:"nfs"`
	// nodeAffinity defines constraints that limit what nodes this volume can be accessed from.
	//
	// This field influences the scheduling of pods that use this volume.
	NodeAffinity *VolumeNodeAffinity `field:"optional" json:"nodeAffinity" yaml:"nodeAffinity"`
	// persistentVolumeReclaimPolicy defines what happens to a persistent volume when released from its claim.
	//
	// Valid options are Retain (default for manually created PersistentVolumes), Delete (default for dynamically provisioned PersistentVolumes), and Recycle (deprecated). Recycle must be supported by the volume plugin underlying this PersistentVolume. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#reclaiming
	PersistentVolumeReclaimPolicy *string `field:"optional" json:"persistentVolumeReclaimPolicy" yaml:"persistentVolumeReclaimPolicy"`
	// photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine.
	//
	// Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.
	PhotonPersistentDisk *PhotonPersistentDiskVolumeSource `field:"optional" json:"photonPersistentDisk" yaml:"photonPersistentDisk"`
	// portworxVolume represents a portworx volume attached and mounted on kubelets host machine.
	//
	// Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate is on.
	PortworxVolume *PortworxVolumeSource `field:"optional" json:"portworxVolume" yaml:"portworxVolume"`
	// quobyte represents a Quobyte mount on the host that shares a pod's lifetime.
	//
	// Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.
	Quobyte *QuobyteVolumeSource `field:"optional" json:"quobyte" yaml:"quobyte"`
	// rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.
	//
	// Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported. More info: https://examples.k8s.io/volumes/rbd/README.md
	Rbd *RbdPersistentVolumeSource `field:"optional" json:"rbd" yaml:"rbd"`
	// scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.
	//
	// Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.
	ScaleIo *ScaleIoPersistentVolumeSource `field:"optional" json:"scaleIo" yaml:"scaleIo"`
	// storageClassName is the name of StorageClass to which this persistent volume belongs.
	//
	// Empty value means that this volume does not belong to any StorageClass.
	StorageClassName *string `field:"optional" json:"storageClassName" yaml:"storageClassName"`
	// storageOS represents a StorageOS volume that is attached to the kubelet's host machine and mounted into the pod.
	//
	// Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported. More info: https://examples.k8s.io/volumes/storageos/README.md
	Storageos *StorageOsPersistentVolumeSource `field:"optional" json:"storageos" yaml:"storageos"`
	// Name of VolumeAttributesClass to which this persistent volume belongs.
	//
	// Empty value is not allowed. When this field is not set, it indicates that this volume does not belong to any VolumeAttributesClass. This field is mutable and can be changed by the CSI driver after a volume has been updated successfully to a new class. For an unbound PersistentVolume, the volumeAttributesClassName will be matched with unbound PersistentVolumeClaims during the binding process. This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
	VolumeAttributesClassName *string `field:"optional" json:"volumeAttributesClassName" yaml:"volumeAttributesClassName"`
	// volumeMode defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state.
	//
	// Value of Filesystem is implied when not included in spec.
	VolumeMode *string `field:"optional" json:"volumeMode" yaml:"volumeMode"`
	// vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine.
	//
	// Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type are redirected to the csi.vsphere.vmware.com CSI driver.
	VsphereVolume *VsphereVirtualDiskVolumeSource `field:"optional" json:"vsphereVolume" yaml:"vsphereVolume"`
}
