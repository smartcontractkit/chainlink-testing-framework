//go:build no_runtime_type_checking

package k8s

// Building without runtime type checking enabled, so all the below just return nil

func validateKubeResourceClaimTemplateListV1Alpha3_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateKubeResourceClaimTemplateListV1Alpha3_IsConstructParameters(x interface{}) error {
	return nil
}

func validateKubeResourceClaimTemplateListV1Alpha3_ManifestParameters(props *KubeResourceClaimTemplateListV1Alpha3Props) error {
	return nil
}

func validateKubeResourceClaimTemplateListV1Alpha3_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewKubeResourceClaimTemplateListV1Alpha3Parameters(scope constructs.Construct, id *string, props *KubeResourceClaimTemplateListV1Alpha3Props) error {
	return nil
}

