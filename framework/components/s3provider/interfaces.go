package s3provider

// Provider is the interface that wraps S3 interaction methods.
type Provider interface {
	GetEndpoint() string
	GetDockerEndpoint() string
	GetConsoleURL() string
	GetSecretKey() string
	GetAccessKey() string
	GetBucket() string
	GetRegion() string
	GetOutput() *Output
}

// ProviderFactory is the interface that standardizes S3 providers constructors.
type ProviderFactory interface {
	New(...Option) (Provider, error)
	NewFrom(*Input) (*Output, error)
}
