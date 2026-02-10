package k8s


// StorageVersionMigration represents a migration of stored data to the latest storage version.
type KubeStorageVersionMigrationV1Alpha1Props struct {
	// Standard object metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// Specification of the migration.
	Spec *StorageVersionMigrationSpecV1Alpha1 `field:"optional" json:"spec" yaml:"spec"`
}

