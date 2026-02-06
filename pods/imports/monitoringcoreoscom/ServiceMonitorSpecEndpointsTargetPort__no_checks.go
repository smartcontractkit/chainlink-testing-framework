//go:build no_runtime_type_checking

package monitoringcoreoscom

// Building without runtime type checking enabled, so all the below just return nil

func validateServiceMonitorSpecEndpointsTargetPort_FromNumberParameters(value *float64) error {
	return nil
}

func validateServiceMonitorSpecEndpointsTargetPort_FromStringParameters(value *string) error {
	return nil
}

