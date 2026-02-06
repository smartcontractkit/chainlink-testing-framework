package k8s

// MutatingAdmissionPolicy describes the definition of an admission mutation policy that mutates the object coming into admission chain.
type KubeMutatingAdmissionPolicyV1Alpha1Props struct {
	// Standard object metadata;
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// Specification of the desired behavior of the MutatingAdmissionPolicy.
	Spec *MutatingAdmissionPolicySpecV1Alpha1 `field:"optional" json:"spec" yaml:"spec"`
}
