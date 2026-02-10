package k8s


// ValidatingAdmissionPolicyList is a list of ValidatingAdmissionPolicy.
type KubeValidatingAdmissionPolicyListV1Alpha1Props struct {
	// List of ValidatingAdmissionPolicy.
	Items *[]*KubeValidatingAdmissionPolicyV1Alpha1Props `field:"required" json:"items" yaml:"items"`
	// Standard list metadata.
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Metadata *ListMeta `field:"optional" json:"metadata" yaml:"metadata"`
}

