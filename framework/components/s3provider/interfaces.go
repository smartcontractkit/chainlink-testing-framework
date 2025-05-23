package s3provider

type Provider interface {
	GetURL() string
	GetEndpoint() string
	GetConsoleURL() string
	GetSecretKey() string
	GetAccessKey() string
	GetBucket() string
	GetRegion() string
}

type ProviderFactory interface {
	NewProvider(...Option) (Provider, error)
}
