package experiments

// Base base experiment data, name and namespace where to store CRD entities
type Base struct {
	Name string
	// Namespace is a namespace where experiment entity will be stored
	Namespace string
}
