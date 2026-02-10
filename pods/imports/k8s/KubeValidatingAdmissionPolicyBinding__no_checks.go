//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubeValidatingAdmissionPolicyBinding_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubeValidatingAdmissionPolicyBinding_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubeValidatingAdmissionPolicyBinding_ManifestParameters(props *KubeValidatingAdmissionPolicyBindingProps) error {
	return nil
}

func validateKubeValidatingAdmissionPolicyBinding_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubeValidatingAdmissionPolicyBindingParameters(scope constructs.Construct, id *string, props *KubeValidatingAdmissionPolicyBindingProps) error {
	return nil
}

