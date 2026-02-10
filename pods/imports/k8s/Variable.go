package k8s


// Variable is the definition of a variable that is used for composition.
//
// A variable is defined as a named expression.
type Variable struct {
	// Expression is the expression that will be evaluated as the value of the variable.
	//
	// The CEL expression has access to the same identifiers as the CEL expressions in Validation.
	Expression *string `field:"required" json:"expression" yaml:"expression"`
	// Name is the name of the variable.
	//
	// The name must be a valid CEL identifier and unique among all variables. The variable can be accessed in other expressions through `variables` For example, if name is "foo", the variable will be available as `variables.foo`
	Name *string `field:"required" json:"name" yaml:"name"`
}

