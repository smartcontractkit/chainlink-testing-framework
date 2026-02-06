package monitoringcoreoscom


// Client certificate to present when doing client-authentication.
type ServiceMonitorSpecEndpointsTlsConfigCert struct {
	// ConfigMap containing data to use for the targets.
	ConfigMap *ServiceMonitorSpecEndpointsTlsConfigCertConfigMap `field:"optional" json:"configMap" yaml:"configMap"`
	// Secret containing data to use for the targets.
	Secret *ServiceMonitorSpecEndpointsTlsConfigCertSecret `field:"optional" json:"secret" yaml:"secret"`
}

