package config

// GenericConfig is an interface for all product based config types to implement
type GenericConfig interface {
	Validate() error
	ApplyOverrides(from interface{}) error
	Default() error
}
