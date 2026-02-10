package k8s


// StorageVersionMigrationList is a collection of storage version migrations.
type KubeStorageVersionMigrationListV1Alpha1Props struct {
	// Items is the list of StorageVersionMigration.
	Items *[]*KubeStorageVersionMigrationV1Alpha1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

