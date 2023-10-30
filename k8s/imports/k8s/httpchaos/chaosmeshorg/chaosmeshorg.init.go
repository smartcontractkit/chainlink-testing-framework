package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.HttpChaos",
		reflect.TypeOf((*HttpChaos)(nil)).Elem(),
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
			j := jsiiProxy_HttpChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosProps",
		reflect.TypeOf((*HttpChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosSpec",
		reflect.TypeOf((*HttpChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.HttpChaosSpecMode",
		reflect.TypeOf((*HttpChaosSpecMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": HttpChaosSpecMode_ONE,
			"ALL": HttpChaosSpecMode_ALL,
			"FIXED": HttpChaosSpecMode_FIXED,
			"FIXED_PERCENT": HttpChaosSpecMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": HttpChaosSpecMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosSpecPatch",
		reflect.TypeOf((*HttpChaosSpecPatch)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosSpecPatchBody",
		reflect.TypeOf((*HttpChaosSpecPatchBody)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosSpecReplace",
		reflect.TypeOf((*HttpChaosSpecReplace)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosSpecSelector",
		reflect.TypeOf((*HttpChaosSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.HttpChaosSpecSelectorExpressionSelectors",
		reflect.TypeOf((*HttpChaosSpecSelectorExpressionSelectors)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.HttpChaosSpecTarget",
		reflect.TypeOf((*HttpChaosSpecTarget)(nil)).Elem(),
		map[string]interface{}{
			"REQUEST": HttpChaosSpecTarget_REQUEST,
			"RESPONSE": HttpChaosSpecTarget_RESPONSE,
		},
	)
}
