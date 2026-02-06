package monitoringcoreoscom


// `authorization` configures the Authorization header credentials to use when scraping the target.
//
// Cannot be set at the same time as `basicAuth`, or `oauth2`.
type ServiceMonitorSpecEndpointsAuthorization struct {
	// Selects a key of a Secret in the namespace that contains the credentials for authentication.
	Credentials *ServiceMonitorSpecEndpointsAuthorizationCredentials `field:"optional" json:"credentials" yaml:"credentials"`
	// Defines the authentication type.
	//
	// The value is case-insensitive.
	// "Basic" is not a supported value.
	// Default: "Bearer".
	Type *string `field:"optional" json:"type" yaml:"type"`
}

