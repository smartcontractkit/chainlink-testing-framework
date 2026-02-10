package k8s


// LeaseCandidateList is a list of Lease objects.
type KubeLeaseCandidateListV1Alpha1Props struct {
	// items is a list of schema objects.
	Items *[]*KubeLeaseCandidateV1Alpha1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

