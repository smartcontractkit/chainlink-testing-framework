package s3provider

type Provider interface {
	GetEndpoint() string
	GetConsoleURL() string
	GetSecretKey() string
	GetAccessKey() string
	GetBucket() string
	GetRegion() string
}

type ProviderFactory interface {
	New(...Option) (Provider, error)
}
