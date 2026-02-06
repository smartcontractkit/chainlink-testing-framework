package monitoringcoreoscom


// `oauth2` configures the OAuth2 settings to use when scraping the target.
//
// It requires Prometheus >= 2.27.0.
// Cannot be set at the same time as `authorization`, or `basicAuth`.
type ServiceMonitorSpecEndpointsOauth2 struct {
	// `clientId` specifies a key of a Secret or ConfigMap containing the OAuth2 client's ID.
	ClientId *ServiceMonitorSpecEndpointsOauth2ClientId `field:"required" json:"clientId" yaml:"clientId"`
	// `clientSecret` specifies a key of a Secret containing the OAuth2 client's secret.
	ClientSecret *ServiceMonitorSpecEndpointsOauth2ClientSecret `field:"required" json:"clientSecret" yaml:"clientSecret"`
	// `tokenURL` configures the URL to fetch the token from.
	TokenUrl *string `field:"required" json:"tokenUrl" yaml:"tokenUrl"`
	// `endpointParams` configures the HTTP parameters to append to the token URL.
	EndpointParams *map[string]*string `field:"optional" json:"endpointParams" yaml:"endpointParams"`
	// `scopes` defines the OAuth2 scopes used for the token request.
	Scopes *[]*string `field:"optional" json:"scopes" yaml:"scopes"`
}

