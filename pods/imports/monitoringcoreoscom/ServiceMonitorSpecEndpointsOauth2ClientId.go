package monitoringcoreoscom


// `clientId` specifies a key of a Secret or ConfigMap containing the OAuth2 client's ID.
type ServiceMonitorSpecEndpointsOauth2ClientId struct {
	// ConfigMap containing data to use for the targets.
	ConfigMap *ServiceMonitorSpecEndpointsOauth2ClientIdConfigMap `field:"optional" json:"configMap" yaml:"configMap"`
	// Secret containing data to use for the targets.
	Secret *ServiceMonitorSpecEndpointsOauth2ClientIdSecret `field:"optional" json:"secret" yaml:"secret"`
}

