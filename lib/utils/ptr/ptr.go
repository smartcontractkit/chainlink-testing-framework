package ptr

// Ptr returns a pointer to the value passed in.
func Ptr[T any](value T) *T {
	return &value
}
