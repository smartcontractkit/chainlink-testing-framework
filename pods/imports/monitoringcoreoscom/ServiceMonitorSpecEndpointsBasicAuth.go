package monitoringcoreoscom


// `basicAuth` configures the Basic Authentication credentials to use when scraping the target.
//
// Cannot be set at the same time as `authorization`, or `oauth2`.
type ServiceMonitorSpecEndpointsBasicAuth struct {
	// `password` specifies a key of a Secret containing the password for authentication.
	Password *ServiceMonitorSpecEndpointsBasicAuthPassword `field:"optional" json:"password" yaml:"password"`
	// `username` specifies a key of a Secret containing the username for authentication.
	Username *ServiceMonitorSpecEndpointsBasicAuthUsername `field:"optional" json:"username" yaml:"username"`
}

