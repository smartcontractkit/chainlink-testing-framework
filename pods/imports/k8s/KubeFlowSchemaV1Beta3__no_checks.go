//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubeFlowSchemaV1Beta3_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubeFlowSchemaV1Beta3_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubeFlowSchemaV1Beta3_ManifestParameters(props *KubeFlowSchemaV1Beta3Props) error {
	return nil
}

func validateKubeFlowSchemaV1Beta3_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubeFlowSchemaV1Beta3Parameters(scope constructs.Construct, id *string, props *KubeFlowSchemaV1Beta3Props) error {
	return nil
}

