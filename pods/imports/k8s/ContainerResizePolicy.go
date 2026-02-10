package k8s


// ContainerResizePolicy represents resource resize policy for the container.
type ContainerResizePolicy struct {
	// Name of the resource to which this resource resize policy applies.
	//
	// Supported values: cpu, memory.
	ResourceName *string `field:"required" json:"resourceName" yaml:"resourceName"`
	// Restart policy to apply when specified resource is resized.
	//
	// If not specified, it defaults to NotRequired.
	RestartPolicy *string `field:"required" json:"restartPolicy" yaml:"restartPolicy"`
}

