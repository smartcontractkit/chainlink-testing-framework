package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.IoChaos",
		reflect.TypeOf((*IoChaos)(nil)).Elem(),
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
			j := jsiiProxy_IoChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosProps",
		reflect.TypeOf((*IoChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpec",
		reflect.TypeOf((*IoChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.IoChaosSpecAction",
		reflect.TypeOf((*IoChaosSpecAction)(nil)).Elem(),
		map[string]interface{}{
			"LATENCY": IoChaosSpecAction_LATENCY,
			"FAULT": IoChaosSpecAction_FAULT,
			"ATTR_OVERRIDE": IoChaosSpecAction_ATTR_OVERRIDE,
			"MISTAKE": IoChaosSpecAction_MISTAKE,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecAttr",
		reflect.TypeOf((*IoChaosSpecAttr)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecAttrAtime",
		reflect.TypeOf((*IoChaosSpecAttrAtime)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecAttrCtime",
		reflect.TypeOf((*IoChaosSpecAttrCtime)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecAttrMtime",
		reflect.TypeOf((*IoChaosSpecAttrMtime)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecMistake",
		reflect.TypeOf((*IoChaosSpecMistake)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.IoChaosSpecMistakeFilling",
		reflect.TypeOf((*IoChaosSpecMistakeFilling)(nil)).Elem(),
		map[string]interface{}{
			"ZERO": IoChaosSpecMistakeFilling_ZERO,
			"RANDOM": IoChaosSpecMistakeFilling_RANDOM,
		},
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.IoChaosSpecMode",
		reflect.TypeOf((*IoChaosSpecMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": IoChaosSpecMode_ONE,
			"ALL": IoChaosSpecMode_ALL,
			"FIXED": IoChaosSpecMode_FIXED,
			"FIXED_PERCENT": IoChaosSpecMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": IoChaosSpecMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecSelector",
		reflect.TypeOf((*IoChaosSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.IoChaosSpecSelectorExpressionSelectors",
		reflect.TypeOf((*IoChaosSpecSelectorExpressionSelectors)(nil)).Elem(),
	)
}
