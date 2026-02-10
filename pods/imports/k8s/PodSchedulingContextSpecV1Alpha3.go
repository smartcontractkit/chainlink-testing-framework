package k8s


// PodSchedulingContextSpec describes where resources for the Pod are needed.
type PodSchedulingContextSpecV1Alpha3 struct {
	// PotentialNodes lists nodes where the Pod might be able to run.
	//
	// The size of this field is limited to 128. This is large enough for many clusters. Larger clusters may need more attempts to find a node that suits all pending resources. This may get increased in the future, but not reduced.
	PotentialNodes *[]*string `field:"optional" json:"potentialNodes" yaml:"potentialNodes"`
	// SelectedNode is the node for which allocation of ResourceClaims that are referenced by the Pod and that use "WaitForFirstConsumer" allocation is to be attempted.
	SelectedNode *string `field:"optional" json:"selectedNode" yaml:"selectedNode"`
}

