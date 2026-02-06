package monitoringcoreoscom

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/smartcontractkit/pods/imports/monitoringcoreoscom/jsii"
)

// Name or number of the target port of the `Pod` object behind the Service, the port must be specified with container port property.
//
// Deprecated: use `port` instead.
type ServiceMonitorSpecEndpointsTargetPort interface {
	Value() interface{}
}

// The jsii proxy struct for ServiceMonitorSpecEndpointsTargetPort
type jsiiProxy_ServiceMonitorSpecEndpointsTargetPort struct {
	_ byte // padding
}

func (j *jsiiProxy_ServiceMonitorSpecEndpointsTargetPort) Value() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


func ServiceMonitorSpecEndpointsTargetPort_FromNumber(value *float64) ServiceMonitorSpecEndpointsTargetPort {
	_init_.Initialize()

	if err := validateServiceMonitorSpecEndpointsTargetPort_FromNumberParameters(value); err != nil {
		panic(err)
	}
	var returns ServiceMonitorSpecEndpointsTargetPort

	_jsii_.StaticInvoke(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTargetPort",
		"fromNumber",
		[]interface{}{value},
		&returns,
	)

	return returns
}

func ServiceMonitorSpecEndpointsTargetPort_FromString(value *string) ServiceMonitorSpecEndpointsTargetPort {
	_init_.Initialize()

	if err := validateServiceMonitorSpecEndpointsTargetPort_FromStringParameters(value); err != nil {
		panic(err)
	}
	var returns ServiceMonitorSpecEndpointsTargetPort

	_jsii_.StaticInvoke(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTargetPort",
		"fromString",
		[]interface{}{value},
		&returns,
	)

	return returns
}

