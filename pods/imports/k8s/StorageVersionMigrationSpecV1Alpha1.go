package k8s


// Spec of the storage version migration.
type StorageVersionMigrationSpecV1Alpha1 struct {
	// The resource that is being migrated.
	//
	// The migrator sends requests to the endpoint serving the resource. Immutable.
	Resource *GroupVersionResourceV1Alpha1 `field:"required" json:"resource" yaml:"resource"`
	// The token used in the list options to get the next chunk of objects to migrate.
	//
	// When the .status.conditions indicates the migration is "Running", users can use this token to check the progress of the migration.
	ContinueToken *string `field:"optional" json:"continueToken" yaml:"continueToken"`
}

