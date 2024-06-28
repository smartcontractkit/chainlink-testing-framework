package ptr

// Ptr returns a pointer to the value passed in.
func Ptr[T any](value T) *T {
	return &value
}

// Value returns the value of a pointer or the zero value of T if the pointer is nil.
func Value[T any](pointer *T) T {
	if pointer == nil {
		var zero T
		return zero
	}
	return *pointer
}
