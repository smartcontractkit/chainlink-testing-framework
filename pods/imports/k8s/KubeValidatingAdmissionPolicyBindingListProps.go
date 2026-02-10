package k8s


// ValidatingAdmissionPolicyBindingList is a list of ValidatingAdmissionPolicyBinding.
type KubeValidatingAdmissionPolicyBindingListProps struct {
	// List of PolicyBinding.
	Items *[]*KubeValidatingAdmissionPolicyBindingProps `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

