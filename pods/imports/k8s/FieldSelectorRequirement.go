package k8s


// FieldSelectorRequirement is a selector that contains values, a key, and an operator that relates the key and values.
type FieldSelectorRequirement struct {
	// key is the field selector key that the requirement applies to.
	Key *string `field:"required" json:"key" yaml:"key"`
	// operator represents a key's relationship to a set of values.
	//
	// Valid operators are In, NotIn, Exists, DoesNotExist. The list of operators may grow in the future.
	Operator *string `field:"required" json:"operator" yaml:"operator"`
	// values is an array of string values.
	//
	// If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty.
	Values *[]*string `field:"optional" json:"values" yaml:"values"`
}

