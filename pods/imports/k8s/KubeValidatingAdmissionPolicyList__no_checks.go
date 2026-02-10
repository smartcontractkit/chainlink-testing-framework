//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubeValidatingAdmissionPolicyList_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubeValidatingAdmissionPolicyList_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubeValidatingAdmissionPolicyList_ManifestParameters(props *KubeValidatingAdmissionPolicyListProps) error {
	return nil
}

func validateKubeValidatingAdmissionPolicyList_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubeValidatingAdmissionPolicyListParameters(scope constructs.Construct, id *string, props *KubeValidatingAdmissionPolicyListProps) error {
	return nil
}

