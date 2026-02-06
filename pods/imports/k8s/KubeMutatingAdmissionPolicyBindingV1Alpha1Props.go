package k8s

// MutatingAdmissionPolicyBinding binds the MutatingAdmissionPolicy with parametrized resources.
//
// MutatingAdmissionPolicyBinding and the optional parameter resource together define how cluster administrators configure policies for clusters.
//
// For a given admission request, each binding will cause its policy to be evaluated N times, where N is 1 for policies/bindings that don't use params, otherwise N is the number of parameters selected by the binding. Each evaluation is constrained by a [runtime cost budget](https://kubernetes.io/docs/reference/using-api/cel/#runtime-cost-budget).
//
// Adding/removing policies, bindings, or params can not affect whether a given (policy, binding, param) combination is within its own CEL budget.
type KubeMutatingAdmissionPolicyBindingV1Alpha1Props struct {
	// Standard object metadata;
	//
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *ObjectMeta `field:"optional" json:"metadata" yaml:"metadata"`
	// Specification of the desired behavior of the MutatingAdmissionPolicyBinding.
	Spec *MutatingAdmissionPolicyBindingSpecV1Alpha1 `field:"optional" json:"spec" yaml:"spec"`
}
