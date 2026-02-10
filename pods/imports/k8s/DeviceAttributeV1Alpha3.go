package k8s


// DeviceAttribute must have exactly one field set.
type DeviceAttributeV1Alpha3 struct {
	// BoolValue is a true/false value.
	Bool *bool `field:"optional" json:"bool" yaml:"bool"`
	// IntValue is a number.
	Int *float64 `field:"optional" json:"int" yaml:"int"`
	// StringValue is a string.
	//
	// Must not be longer than 64 characters.
	String *string `field:"optional" json:"string" yaml:"string"`
	// VersionValue is a semantic version according to semver.org spec 2.0.0. Must not be longer than 64 characters.
	Version *string `field:"optional" json:"version" yaml:"version"`
}

