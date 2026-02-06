package monitoringcoreoscom


// Certificate authority used when verifying server certificates.
type ServiceMonitorSpecEndpointsTlsConfigCa struct {
	// ConfigMap containing data to use for the targets.
	ConfigMap *ServiceMonitorSpecEndpointsTlsConfigCaConfigMap `field:"optional" json:"configMap" yaml:"configMap"`
	// Secret containing data to use for the targets.
	Secret *ServiceMonitorSpecEndpointsTlsConfigCaSecret `field:"optional" json:"secret" yaml:"secret"`
}

