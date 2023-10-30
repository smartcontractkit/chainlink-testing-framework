package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.StressChaos",
		reflect.TypeOf((*StressChaos)(nil)).Elem(),
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
			j := jsiiProxy_StressChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosProps",
		reflect.TypeOf((*StressChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosSpec",
		reflect.TypeOf((*StressChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.StressChaosSpecMode",
		reflect.TypeOf((*StressChaosSpecMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": StressChaosSpecMode_ONE,
			"ALL": StressChaosSpecMode_ALL,
			"FIXED": StressChaosSpecMode_FIXED,
			"FIXED_PERCENT": StressChaosSpecMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": StressChaosSpecMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosSpecSelector",
		reflect.TypeOf((*StressChaosSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosSpecSelectorExpressionSelectors",
		reflect.TypeOf((*StressChaosSpecSelectorExpressionSelectors)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosSpecStressors",
		reflect.TypeOf((*StressChaosSpecStressors)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosSpecStressorsCpu",
		reflect.TypeOf((*StressChaosSpecStressorsCpu)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.StressChaosSpecStressorsMemory",
		reflect.TypeOf((*StressChaosSpecStressorsMemory)(nil)).Elem(),
	)
}
