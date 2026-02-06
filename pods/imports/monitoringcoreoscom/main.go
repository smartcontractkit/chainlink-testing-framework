// monitoringcoreoscom
package monitoringcoreoscom

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monitoringcoreoscom.ServiceMonitor",
		reflect.TypeOf((*ServiceMonitor)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addJsonPatch", GoMethod: "AddJsonPatch"},
			_jsii_.MemberProperty{JsiiProperty: "apiGroup", GoGetter: "ApiGroup"},
			_jsii_.MemberProperty{JsiiProperty: "apiVersion", GoGetter: "ApiVersion"},
			_jsii_.MemberProperty{JsiiProperty: "chart", GoGetter: "Chart"},
			_jsii_.MemberProperty{JsiiProperty: "kind", GoGetter: "Kind"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_ServiceMonitor{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorProps",
		reflect.TypeOf((*ServiceMonitorProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpec",
		reflect.TypeOf((*ServiceMonitorSpec)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecAttachMetadata",
		reflect.TypeOf((*ServiceMonitorSpecAttachMetadata)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpoints",
		reflect.TypeOf((*ServiceMonitorSpecEndpoints)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsAuthorization",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsAuthorization)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsAuthorizationCredentials",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsAuthorizationCredentials)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsBasicAuth",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsBasicAuth)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsBasicAuthPassword",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsBasicAuthPassword)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsBasicAuthUsername",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsBasicAuthUsername)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsBearerTokenSecret",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsBearerTokenSecret)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsMetricRelabelings",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsMetricRelabelings)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsMetricRelabelingsAction",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsMetricRelabelingsAction)(nil)).Elem(),
		map[string]interface{}{
			"REPLACE": ServiceMonitorSpecEndpointsMetricRelabelingsAction_REPLACE,
			"KEEP": ServiceMonitorSpecEndpointsMetricRelabelingsAction_KEEP,
			"DROP": ServiceMonitorSpecEndpointsMetricRelabelingsAction_DROP,
			"HASHMOD": ServiceMonitorSpecEndpointsMetricRelabelingsAction_HASHMOD,
			"LABELMAP": ServiceMonitorSpecEndpointsMetricRelabelingsAction_LABELMAP,
			"LABELDROP": ServiceMonitorSpecEndpointsMetricRelabelingsAction_LABELDROP,
			"LABELKEEP": ServiceMonitorSpecEndpointsMetricRelabelingsAction_LABELKEEP,
			"LOWERCASE": ServiceMonitorSpecEndpointsMetricRelabelingsAction_LOWERCASE,
			"UPPERCASE": ServiceMonitorSpecEndpointsMetricRelabelingsAction_UPPERCASE,
			"KEEPEQUAL": ServiceMonitorSpecEndpointsMetricRelabelingsAction_KEEPEQUAL,
			"DROPEQUAL": ServiceMonitorSpecEndpointsMetricRelabelingsAction_DROPEQUAL,
		},
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsOauth2",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsOauth2)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsOauth2ClientId",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsOauth2ClientId)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsOauth2ClientIdConfigMap",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsOauth2ClientIdConfigMap)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsOauth2ClientIdSecret",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsOauth2ClientIdSecret)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsOauth2ClientSecret",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsOauth2ClientSecret)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsRelabelings",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsRelabelings)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsRelabelingsAction",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsRelabelingsAction)(nil)).Elem(),
		map[string]interface{}{
			"REPLACE": ServiceMonitorSpecEndpointsRelabelingsAction_REPLACE,
			"KEEP": ServiceMonitorSpecEndpointsRelabelingsAction_KEEP,
			"DROP": ServiceMonitorSpecEndpointsRelabelingsAction_DROP,
			"HASHMOD": ServiceMonitorSpecEndpointsRelabelingsAction_HASHMOD,
			"LABELMAP": ServiceMonitorSpecEndpointsRelabelingsAction_LABELMAP,
			"LABELDROP": ServiceMonitorSpecEndpointsRelabelingsAction_LABELDROP,
			"LABELKEEP": ServiceMonitorSpecEndpointsRelabelingsAction_LABELKEEP,
			"LOWERCASE": ServiceMonitorSpecEndpointsRelabelingsAction_LOWERCASE,
			"UPPERCASE": ServiceMonitorSpecEndpointsRelabelingsAction_UPPERCASE,
			"KEEPEQUAL": ServiceMonitorSpecEndpointsRelabelingsAction_KEEPEQUAL,
			"DROPEQUAL": ServiceMonitorSpecEndpointsRelabelingsAction_DROPEQUAL,
		},
	)
	_jsii_.RegisterEnum(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsScheme",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsScheme)(nil)).Elem(),
		map[string]interface{}{
			"HTTP": ServiceMonitorSpecEndpointsScheme_HTTP,
			"HTTPS": ServiceMonitorSpecEndpointsScheme_HTTPS,
		},
	)
	_jsii_.RegisterClass(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTargetPort",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTargetPort)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			return &jsiiProxy_ServiceMonitorSpecEndpointsTargetPort{}
		},
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfig",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigCa",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigCa)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigCaConfigMap",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigCaConfigMap)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigCaSecret",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigCaSecret)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigCert",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigCert)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigCertConfigMap",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigCertConfigMap)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigCertSecret",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigCertSecret)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecEndpointsTlsConfigKeySecret",
		reflect.TypeOf((*ServiceMonitorSpecEndpointsTlsConfigKeySecret)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecNamespaceSelector",
		reflect.TypeOf((*ServiceMonitorSpecNamespaceSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecSelector",
		reflect.TypeOf((*ServiceMonitorSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monitoringcoreoscom.ServiceMonitorSpecSelectorMatchExpressions",
		reflect.TypeOf((*ServiceMonitorSpecSelectorMatchExpressions)(nil)).Elem(),
	)
}
