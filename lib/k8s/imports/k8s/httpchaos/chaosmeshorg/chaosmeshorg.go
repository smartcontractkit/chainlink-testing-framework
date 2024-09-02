// chaos-meshorg
package chaosmeshorg

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/imports/k8s/httpchaos/chaosmeshorg/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/imports/k8s/httpchaos/chaosmeshorg/internal"
)

// HTTPChaos is the Schema for the HTTPchaos API.
type HttpChaos interface {
	cdk8s.ApiObject
	// The group portion of the API version (e.g. `authorization.k8s.io`).
	ApiGroup() *string
	// The object's API version (e.g. `authorization.k8s.io/v1`).
	ApiVersion() *string
	// The chart in which this object is defined.
	Chart() cdk8s.Chart
	// The object kind.
	Kind() *string
	// Metadata associated with this API object.
	Metadata() cdk8s.ApiObjectMetadataDefinition
	// The name of the API object.
	//
	// If a name is specified in `metadata.name` this will be the name returned.
	// Otherwise, a name will be generated by calling
	// `Chart.of(this).generatedObjectName(this)`, which by default uses the
	// construct path to generate a DNS-compatible name for the resource.
	Name() *string
	// The tree node.
	Node() constructs.Node
	// Create a dependency between this ApiObject and other constructs.
	//
	// These can be other ApiObjects, Charts, or custom.
	AddDependency(dependencies ...constructs.IConstruct)
	// Applies a set of RFC-6902 JSON-Patch operations to the manifest synthesized for this API object.
	//
	// Example:
	//     kubePod.addJsonPatch(JsonPatch.replace('/spec/enableServiceLinks', true));
	//
	AddJsonPatch(ops ...cdk8s.JsonPatch)
	// Renders the object to Kubernetes JSON.
	ToJson() interface{}
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for HttpChaos
type jsiiProxy_HttpChaos struct {
	internal.Type__cdk8sApiObject
}

func (j *jsiiProxy_HttpChaos) ApiGroup() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HttpChaos) ApiVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HttpChaos) Chart() cdk8s.Chart {
	var returns cdk8s.Chart
	_jsii_.Get(
		j,
		"chart",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HttpChaos) Kind() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kind",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HttpChaos) Metadata() cdk8s.ApiObjectMetadataDefinition {
	var returns cdk8s.ApiObjectMetadataDefinition
	_jsii_.Get(
		j,
		"metadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HttpChaos) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HttpChaos) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

// Defines a "HTTPChaos" API object.
func NewHttpChaos(scope constructs.Construct, id *string, props *HttpChaosProps) HttpChaos {
	_init_.Initialize()

	j := jsiiProxy_HttpChaos{}

	_jsii_.Create(
		"chaos-meshorg.HttpChaos",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Defines a "HTTPChaos" API object.
func NewHttpChaos_Override(h HttpChaos, scope constructs.Construct, id *string, props *HttpChaosProps) {
	_init_.Initialize()

	_jsii_.Create(
		"chaos-meshorg.HttpChaos",
		[]interface{}{scope, id, props},
		h,
	)
}

// Checks if `x` is a construct.
//
// Use this method instead of `instanceof` to properly detect `Construct`
// instances, even when the construct library is symlinked.
//
// Explanation: in JavaScript, multiple copies of the `constructs` library on
// disk are seen as independent, completely different libraries. As a
// consequence, the class `Construct` in each copy of the `constructs` library
// is seen as a different class, and an instance of one class will not test as
// `instanceof` the other class. `npm install` will not create installations
// like this, but users may manually symlink construct libraries together or
// use a monorepo tool: in those cases, multiple copies of the `constructs`
// library can be accidentally installed, and `instanceof` will behave
// unpredictably. It is safest to avoid using `instanceof`, and using
// this type-testing method instead.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
func HttpChaos_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"chaos-meshorg.HttpChaos",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Renders a Kubernetes manifest for "HTTPChaos".
//
// This can be used to inline resource manifests inside other objects (e.g. as templates).
func HttpChaos_Manifest(props *HttpChaosProps) interface{} {
	_init_.Initialize()

	var returns interface{}

	_jsii_.StaticInvoke(
		"chaos-meshorg.HttpChaos",
		"manifest",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Returns the `ApiObject` named `Resource` which is a child of the given construct.
//
// If `c` is an `ApiObject`, it is returned directly. Throws an
// exception if the construct does not have a child named `Default` _or_ if
// this child is not an `ApiObject`.
func HttpChaos_Of(c constructs.IConstruct) cdk8s.ApiObject {
	_init_.Initialize()

	var returns cdk8s.ApiObject

	_jsii_.StaticInvoke(
		"chaos-meshorg.HttpChaos",
		"of",
		[]interface{}{c},
		&returns,
	)

	return returns
}

func HttpChaos_GVK() *cdk8s.GroupVersionKind {
	_init_.Initialize()
	var returns *cdk8s.GroupVersionKind
	_jsii_.StaticGet(
		"chaos-meshorg.HttpChaos",
		"GVK",
		&returns,
	)
	return returns
}

func (h *jsiiProxy_HttpChaos) AddDependency(dependencies ...constructs.IConstruct) {
	args := []interface{}{}
	for _, a := range dependencies {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		h,
		"addDependency",
		args,
	)
}

func (h *jsiiProxy_HttpChaos) AddJsonPatch(ops ...cdk8s.JsonPatch) {
	args := []interface{}{}
	for _, a := range ops {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		h,
		"addJsonPatch",
		args,
	)
}

func (h *jsiiProxy_HttpChaos) ToJson() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		h,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (h *jsiiProxy_HttpChaos) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		h,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// HTTPChaos is the Schema for the HTTPchaos API.
type HttpChaosProps struct {
	Metadata *cdk8s.ApiObjectMetadata `field:"optional" json:"metadata" yaml:"metadata"`
	Spec     *HttpChaosSpec           `field:"optional" json:"spec" yaml:"spec"`
}

type HttpChaosSpec struct {
	// Mode defines the mode to run chaos action.
	//
	// Supported mode: one / all / fixed / fixed-percent / random-max-percent.
	Mode HttpChaosSpecMode `field:"required" json:"mode" yaml:"mode"`
	// Selector is used to select pods that are used to inject chaos action.
	Selector *HttpChaosSpecSelector `field:"required" json:"selector" yaml:"selector"`
	// Target is the object to be selected and injected.
	Target HttpChaosSpecTarget `field:"required" json:"target" yaml:"target"`
	// Abort is a rule to abort a http session.
	Abort *bool `field:"optional" json:"abort" yaml:"abort"`
	// Code is a rule to select target by http status code in response.
	Code *float64 `field:"optional" json:"code" yaml:"code"`
	// Delay represents the delay of the target request/response.
	//
	// A duration string is a possibly unsigned sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	Delay *string `field:"optional" json:"delay" yaml:"delay"`
	// Duration represents the duration of the chaos action.
	Duration *string `field:"optional" json:"duration" yaml:"duration"`
	// Method is a rule to select target by http method in request.
	Method *string `field:"optional" json:"method" yaml:"method"`
	// Patch is a rule to patch some contents in target.
	Patch *HttpChaosSpecPatch `field:"optional" json:"patch" yaml:"patch"`
	// Path is a rule to select target by uri path in http request.
	Path *string `field:"optional" json:"path" yaml:"path"`
	// Port represents the target port to be proxy of.
	Port *float64 `field:"optional" json:"port" yaml:"port"`
	// Replace is a rule to replace some contents in target.
	Replace *HttpChaosSpecReplace `field:"optional" json:"replace" yaml:"replace"`
	// RequestHeaders is a rule to select target by http headers in request.
	//
	// The key-value pairs represent header name and header value pairs.
	RequestHeaders *map[string]*string `field:"optional" json:"requestHeaders" yaml:"requestHeaders"`
	// ResponseHeaders is a rule to select target by http headers in response.
	//
	// The key-value pairs represent header name and header value pairs.
	ResponseHeaders *map[string]*string `field:"optional" json:"responseHeaders" yaml:"responseHeaders"`
	// Value is required when the mode is set to `FixedPodMode` / `FixedPercentPodMod` / `RandomMaxPercentPodMod`.
	//
	// If `FixedPodMode`, provide an integer of pods to do chaos action. If `FixedPercentPodMod`, provide a number from 0-100 to specify the percent of pods the server can do chaos action. IF `RandomMaxPercentPodMod`,  provide a number from 0-100 to specify the max percent of pods to do chaos action
	Value *string `field:"optional" json:"value" yaml:"value"`
}

// Mode defines the mode to run chaos action.
//
// Supported mode: one / all / fixed / fixed-percent / random-max-percent.
type HttpChaosSpecMode string

const (
	// one.
	HttpChaosSpecMode_ONE HttpChaosSpecMode = "ONE"
	// all.
	HttpChaosSpecMode_ALL HttpChaosSpecMode = "ALL"
	// fixed.
	HttpChaosSpecMode_FIXED HttpChaosSpecMode = "FIXED"
	// fixed-percent.
	HttpChaosSpecMode_FIXED_PERCENT HttpChaosSpecMode = "FIXED_PERCENT"
	// random-max-percent.
	HttpChaosSpecMode_RANDOM_MAX_PERCENT HttpChaosSpecMode = "RANDOM_MAX_PERCENT"
)

// Patch is a rule to patch some contents in target.
type HttpChaosSpecPatch struct {
	// Body is a rule to patch message body of target.
	Body *HttpChaosSpecPatchBody `field:"optional" json:"body" yaml:"body"`
	// Headers is a rule to append http headers of target.
	//
	// For example: `[["Set-Cookie", "<one cookie>"], ["Set-Cookie", "<another cookie>"]]`.
	Headers *[]*[]*string `field:"optional" json:"headers" yaml:"headers"`
	// Queries is a rule to append uri queries of target(Request only).
	//
	// For example: `[["foo", "bar"], ["foo", "unknown"]]`.
	Queries *[]*[]*string `field:"optional" json:"queries" yaml:"queries"`
}

// Body is a rule to patch message body of target.
type HttpChaosSpecPatchBody struct {
	// Type represents the patch type, only support `JSON` as [merge patch json](https://tools.ietf.org/html/rfc7396) currently.
	Type *string `field:"required" json:"type" yaml:"type"`
	// Value is the patch contents.
	Value *string `field:"required" json:"value" yaml:"value"`
}

// Replace is a rule to replace some contents in target.
type HttpChaosSpecReplace struct {
	// Body is a rule to replace http message body in target.
	Body *string `field:"optional" json:"body" yaml:"body"`
	// Code is a rule to replace http status code in response.
	Code *float64 `field:"optional" json:"code" yaml:"code"`
	// Headers is a rule to replace http headers of target.
	//
	// The key-value pairs represent header name and header value pairs.
	Headers *map[string]*string `field:"optional" json:"headers" yaml:"headers"`
	// Method is a rule to replace http method in request.
	Method *string `field:"optional" json:"method" yaml:"method"`
	// Path is rule to to replace uri path in http request.
	Path *string `field:"optional" json:"path" yaml:"path"`
	// Queries is a rule to replace uri queries in http request.
	//
	// For example, with value `{ "foo": "unknown" }`, the `/?foo=bar` will be altered to `/?foo=unknown`,
	Queries *map[string]*string `field:"optional" json:"queries" yaml:"queries"`
}

// Selector is used to select pods that are used to inject chaos action.
type HttpChaosSpecSelector struct {
	// Map of string keys and values that can be used to select objects.
	//
	// A selector based on annotations.
	AnnotationSelectors *map[string]*string `field:"optional" json:"annotationSelectors" yaml:"annotationSelectors"`
	// a slice of label selector expressions that can be used to select objects.
	//
	// A list of selectors based on set-based label expressions.
	ExpressionSelectors *[]*HttpChaosSpecSelectorExpressionSelectors `field:"optional" json:"expressionSelectors" yaml:"expressionSelectors"`
	// Map of string keys and values that can be used to select objects.
	//
	// A selector based on fields.
	FieldSelectors *map[string]*string `field:"optional" json:"fieldSelectors" yaml:"fieldSelectors"`
	// Map of string keys and values that can be used to select objects.
	//
	// A selector based on labels.
	LabelSelectors *map[string]*string `field:"optional" json:"labelSelectors" yaml:"labelSelectors"`
	// Namespaces is a set of namespace to which objects belong.
	Namespaces *[]*string `field:"optional" json:"namespaces" yaml:"namespaces"`
	// Nodes is a set of node name and objects must belong to these nodes.
	Nodes *[]*string `field:"optional" json:"nodes" yaml:"nodes"`
	// Map of string keys and values that can be used to select nodes.
	//
	// Selector which must match a node's labels, and objects must belong to these selected nodes.
	NodeSelectors *map[string]*string `field:"optional" json:"nodeSelectors" yaml:"nodeSelectors"`
	// PodPhaseSelectors is a set of condition of a pod at the current time.
	//
	// supported value: Pending / Running / Succeeded / Failed / Unknown.
	PodPhaseSelectors *[]*string `field:"optional" json:"podPhaseSelectors" yaml:"podPhaseSelectors"`
	// Pods is a map of string keys and a set values that used to select pods.
	//
	// The key defines the namespace which pods belong, and the each values is a set of pod names.
	Pods *map[string]*[]*string `field:"optional" json:"pods" yaml:"pods"`
}

// A label selector requirement is a selector that contains values, a key, and an operator that relates the key and values.
type HttpChaosSpecSelectorExpressionSelectors struct {
	// key is the label key that the selector applies to.
	Key *string `field:"required" json:"key" yaml:"key"`
	// operator represents a key's relationship to a set of values.
	//
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator *string `field:"required" json:"operator" yaml:"operator"`
	// values is an array of string values.
	//
	// If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
	Values *[]*string `field:"optional" json:"values" yaml:"values"`
}

// Target is the object to be selected and injected.
type HttpChaosSpecTarget string

const (
	// Request.
	HttpChaosSpecTarget_REQUEST HttpChaosSpecTarget = "REQUEST"
	// Response.
	HttpChaosSpecTarget_RESPONSE HttpChaosSpecTarget = "RESPONSE"
)
