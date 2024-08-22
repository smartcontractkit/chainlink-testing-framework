package k8schaos

func Ptr[T any](value T) *T {
	return &value
}
