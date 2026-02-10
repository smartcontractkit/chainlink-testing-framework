//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubePriorityLevelConfigurationList_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubePriorityLevelConfigurationList_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubePriorityLevelConfigurationList_ManifestParameters(props *KubePriorityLevelConfigurationListProps) error {
	return nil
}

func validateKubePriorityLevelConfigurationList_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubePriorityLevelConfigurationListParameters(scope constructs.Construct, id *string, props *KubePriorityLevelConfigurationListProps) error {
	return nil
}

