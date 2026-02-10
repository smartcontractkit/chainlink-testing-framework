package k8s


// FlowSchemaList is a list of FlowSchema objects.
type KubeFlowSchemaListProps struct {
	// `items` is a list of FlowSchemas.
	Items *[]*KubeFlowSchemaProps `field:"required" json:"items" yaml:"items"`
	// `metadata` is the standard list metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

