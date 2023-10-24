package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.PodIoChaos",
		reflect.TypeOf((*PodIoChaos)(nil)).Elem(),
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
			j := jsiiProxy_PodIoChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosProps",
		reflect.TypeOf((*PodIoChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpec",
		reflect.TypeOf((*PodIoChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpecActions",
		reflect.TypeOf((*PodIoChaosSpecActions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpecActionsAtime",
		reflect.TypeOf((*PodIoChaosSpecActionsAtime)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpecActionsCtime",
		reflect.TypeOf((*PodIoChaosSpecActionsCtime)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpecActionsFaults",
		reflect.TypeOf((*PodIoChaosSpecActionsFaults)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpecActionsMistake",
		reflect.TypeOf((*PodIoChaosSpecActionsMistake)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.PodIoChaosSpecActionsMistakeFilling",
		reflect.TypeOf((*PodIoChaosSpecActionsMistakeFilling)(nil)).Elem(),
		map[string]interface{}{
			"ZERO": PodIoChaosSpecActionsMistakeFilling_ZERO,
			"RANDOM": PodIoChaosSpecActionsMistakeFilling_RANDOM,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodIoChaosSpecActionsMtime",
		reflect.TypeOf((*PodIoChaosSpecActionsMtime)(nil)).Elem(),
	)
}
