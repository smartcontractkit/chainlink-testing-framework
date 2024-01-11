package config

// GenericConfig is an interface for all product based config types to implement
type GenericConfig[T any] interface {
	Validate() error
	ApplyOverride(from T) error
	Default() error
}
