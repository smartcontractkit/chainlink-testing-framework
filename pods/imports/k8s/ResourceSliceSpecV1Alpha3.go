package k8s


// ResourceSliceSpec contains the information published by the driver in one ResourceSlice.
type ResourceSliceSpecV1Alpha3 struct {
	// Driver identifies the DRA driver providing the capacity information.
	//
	// A field selector can be used to list only ResourceSlice objects with a certain driver name.
	//
	// Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver. This field is immutable.
	Driver *string `field:"required" json:"driver" yaml:"driver"`
	// Pool describes the pool that this ResourceSlice belongs to.
	Pool *ResourcePoolV1Alpha3 `field:"required" json:"pool" yaml:"pool"`
	// AllNodes indicates that all nodes have access to the resources in the pool.
	//
	// Exactly one of NodeName, NodeSelector and AllNodes must be set.
	AllNodes *bool `field:"optional" json:"allNodes" yaml:"allNodes"`
	// Devices lists some or all of the devices in this pool.
	//
	// Must not have more than 128 entries.
	Devices *[]*DeviceV1Alpha3 `field:"optional" json:"devices" yaml:"devices"`
	// NodeName identifies the node which provides the resources in this pool.
	//
	// A field selector can be used to list only ResourceSlice objects belonging to a certain node.
	//
	// This field can be used to limit access from nodes to ResourceSlices with the same node name. It also indicates to autoscalers that adding new nodes of the same type as some old node might also make new resources available.
	//
	// Exactly one of NodeName, NodeSelector and AllNodes must be set. This field is immutable.
	NodeName *string `field:"optional" json:"nodeName" yaml:"nodeName"`
	// NodeSelector defines which nodes have access to the resources in the pool, when that pool is not limited to a single node.
	//
	// Must use exactly one term.
	//
	// Exactly one of NodeName, NodeSelector and AllNodes must be set.
	NodeSelector *NodeSelector `field:"optional" json:"nodeSelector" yaml:"nodeSelector"`
}

