package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.TimeChaos",
		reflect.TypeOf((*TimeChaos)(nil)).Elem(),
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
			j := jsiiProxy_TimeChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.TimeChaosProps",
		reflect.TypeOf((*TimeChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.TimeChaosSpec",
		reflect.TypeOf((*TimeChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.TimeChaosSpecMode",
		reflect.TypeOf((*TimeChaosSpecMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": TimeChaosSpecMode_ONE,
			"ALL": TimeChaosSpecMode_ALL,
			"FIXED": TimeChaosSpecMode_FIXED,
			"FIXED_PERCENT": TimeChaosSpecMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": TimeChaosSpecMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.TimeChaosSpecSelector",
		reflect.TypeOf((*TimeChaosSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.TimeChaosSpecSelectorExpressionSelectors",
		reflect.TypeOf((*TimeChaosSpecSelectorExpressionSelectors)(nil)).Elem(),
	)
}
