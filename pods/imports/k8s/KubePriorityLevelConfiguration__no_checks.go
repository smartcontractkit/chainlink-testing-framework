//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubePriorityLevelConfiguration_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubePriorityLevelConfiguration_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubePriorityLevelConfiguration_ManifestParameters(props *KubePriorityLevelConfigurationProps) error {
	return nil
}

func validateKubePriorityLevelConfiguration_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubePriorityLevelConfigurationParameters(scope constructs.Construct, id *string, props *KubePriorityLevelConfigurationProps) error {
	return nil
}

