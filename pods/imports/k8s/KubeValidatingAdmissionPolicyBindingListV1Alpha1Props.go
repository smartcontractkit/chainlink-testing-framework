package k8s


// ValidatingAdmissionPolicyBindingList is a list of ValidatingAdmissionPolicyBinding.
type KubeValidatingAdmissionPolicyBindingListV1Alpha1Props struct {
	// List of PolicyBinding.
	Items *[]*KubeValidatingAdmissionPolicyBindingV1Alpha1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

