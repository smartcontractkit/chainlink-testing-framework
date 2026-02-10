package k8s


// ValidatingAdmissionPolicy describes the definition of an admission validation policy that accepts or rejects an object without changing it.
type KubeValidatingAdmissionPolicyProps struct {
	// Standard object metadata;
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// Specification of the desired behavior of the ValidatingAdmissionPolicy.
	Spec *ValidatingAdmissionPolicySpec `field:"optional" json:"spec" yaml:"spec"`
}

