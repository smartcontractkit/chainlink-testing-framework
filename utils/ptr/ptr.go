package ptr

// Ptr returns a pointer to the value passed in.
func Ptr[T any](value T) *T {
	return &value
}

// Val returns the value of the pointer passed in.
func Val[T any](value *T) T {
	return *value
}
