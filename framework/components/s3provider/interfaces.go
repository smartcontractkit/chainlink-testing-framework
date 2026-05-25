package s3provider

import "context"

// Provider is the interface that wraps S3 interaction methods.
type Provider interface {
	GetEndpoint() string
	GetBaseEndpoint() string
	GetConsoleURL() string
	GetConsoleBaseURL() string
	GetSecretKey() string
	GetAccessKey() string
	GetBucket() string
	GetRegion() string
	Output() *Output
}

// ProviderFactory is the interface that standardizes S3 providers constructors.
type ProviderFactory interface {
	New(...Option) (Provider, error)
	NewFrom(*Input) (*Output, error)
	NewWithContext(ctx context.Context, options ...Option) (Provider, error)
	NewWithContextFrom(ctx context.Context, input *Input) (*Output, error)
}
