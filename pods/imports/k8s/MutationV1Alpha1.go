package k8s

// Mutation specifies the CEL expression which is used to apply the Mutation.
type MutationV1Alpha1 struct {
	// patchType indicates the patch strategy used.
	//
	// Allowed values are "ApplyConfiguration" and "JSONPatch". Required.
	PatchType *string `field:"required" json:"patchType" yaml:"patchType"`
	// applyConfiguration defines the desired configuration values of an object.
	//
	// The configuration is applied to the admission object using [structured merge diff](https://github.com/kubernetes-sigs/structured-merge-diff). A CEL expression is used to create apply configuration.
	ApplyConfiguration *ApplyConfigurationV1Alpha1 `field:"optional" json:"applyConfiguration" yaml:"applyConfiguration"`
	// jsonPatch defines a [JSON patch](https://jsonpatch.com/) operation to perform a mutation to the object. A CEL expression is used to create the JSON patch.
	JsonPatch *JsonPatchV1Alpha1 `field:"optional" json:"jsonPatch" yaml:"jsonPatch"`
}
