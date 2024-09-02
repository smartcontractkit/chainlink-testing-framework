package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.PodNetworkChaos",
		reflect.TypeOf((*PodNetworkChaos)(nil)).Elem(),
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
			j := jsiiProxy_PodNetworkChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosProps",
		reflect.TypeOf((*PodNetworkChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpec",
		reflect.TypeOf((*PodNetworkChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecIpsets",
		reflect.TypeOf((*PodNetworkChaosSpecIpsets)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecIptables",
		reflect.TypeOf((*PodNetworkChaosSpecIptables)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcs",
		reflect.TypeOf((*PodNetworkChaosSpecTcs)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcsBandwidth",
		reflect.TypeOf((*PodNetworkChaosSpecTcsBandwidth)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcsCorrupt",
		reflect.TypeOf((*PodNetworkChaosSpecTcsCorrupt)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcsDelay",
		reflect.TypeOf((*PodNetworkChaosSpecTcsDelay)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcsDelayReorder",
		reflect.TypeOf((*PodNetworkChaosSpecTcsDelayReorder)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcsDuplicate",
		reflect.TypeOf((*PodNetworkChaosSpecTcsDuplicate)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodNetworkChaosSpecTcsLoss",
		reflect.TypeOf((*PodNetworkChaosSpecTcsLoss)(nil)).Elem(),
	)
}
