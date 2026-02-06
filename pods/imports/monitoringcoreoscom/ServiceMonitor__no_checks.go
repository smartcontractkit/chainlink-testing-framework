//go:build no_runtime_type_checking

package monitoringcoreoscom

// Building without runtime type checking enabled, so all the below just return nil

func validateServiceMonitor_IsApiObjectParameters(o interface{}) error {
	return nil
}

func validateServiceMonitor_IsConstructParameters(x interface{}) error {
	return nil
}

func validateServiceMonitor_ManifestParameters(props *ServiceMonitorProps) error {
	return nil
}

func validateServiceMonitor_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewServiceMonitorParameters(scope constructs.Construct, id *string, props *ServiceMonitorProps) error {
	return nil
}

