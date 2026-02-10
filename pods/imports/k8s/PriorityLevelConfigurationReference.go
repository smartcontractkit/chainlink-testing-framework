package k8s


// PriorityLevelConfigurationReference contains information that points to the "request-priority" being used.
type PriorityLevelConfigurationReference struct {
	// `name` is the name of the priority level configuration being referenced Required.
	Name *string `field:"required" json:"name" yaml:"name"`
}

