//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubeFlowSchemaList_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubeFlowSchemaList_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubeFlowSchemaList_ManifestParameters(props *KubeFlowSchemaListProps) error {
	return nil
}

func validateKubeFlowSchemaList_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubeFlowSchemaListParameters(scope constructs.Construct, id *string, props *KubeFlowSchemaListProps) error {
	return nil
}

