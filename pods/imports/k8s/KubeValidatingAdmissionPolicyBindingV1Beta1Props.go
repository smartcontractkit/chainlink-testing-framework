package k8s


// ValidatingAdmissionPolicyBinding binds the ValidatingAdmissionPolicy with paramerized resources.
//
// ValidatingAdmissionPolicyBinding and parameter CRDs together define how cluster administrators configure policies for clusters.
//
// For a given admission request, each binding will cause its policy to be evaluated N times, where N is 1 for policies/bindings that don't use params, otherwise N is the number of parameters selected by the binding.
//
// The CEL expressions of a policy must have a computed CEL cost below the maximum CEL budget. Each evaluation of the policy is given an independent CEL cost budget. Adding/removing policies, bindings, or params can not affect whether a given (policy, binding, param) combination is within its own CEL budget.
type KubeValidatingAdmissionPolicyBindingV1Beta1Props struct {
	// Standard object metadata;
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// Specification of the desired behavior of the ValidatingAdmissionPolicyBinding.
	Spec *ValidatingAdmissionPolicyBindingSpecV1Beta1 `field:"optional" json:"spec" yaml:"spec"`
}

