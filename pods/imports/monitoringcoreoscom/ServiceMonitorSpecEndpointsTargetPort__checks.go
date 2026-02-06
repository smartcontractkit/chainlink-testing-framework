//go:build !no_runtime_type_checking

package monitoringcoreoscom

import (
	"fmt"
)

func validateServiceMonitorSpecEndpointsTargetPort_FromNumberParameters(value *float64) error {
	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func validateServiceMonitorSpecEndpointsTargetPort_FromStringParameters(value *string) error {
	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

