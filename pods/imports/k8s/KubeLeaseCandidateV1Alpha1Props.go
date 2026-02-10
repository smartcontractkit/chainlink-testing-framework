package k8s


// LeaseCandidate defines a candidate for a Lease object.
//
// Candidates are created such that coordinated leader election will pick the best leader from the list of candidates.
type KubeLeaseCandidateV1Alpha1Props struct {
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// spec contains the specification of the Lease.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *LeaseCandidateSpecV1Alpha1 `field:"optional" json:"spec" yaml:"spec"`
}

