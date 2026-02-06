package monitoringcoreoscom


// TLS configuration to use when scraping the target.
type ServiceMonitorSpecEndpointsTlsConfig struct {
	// Certificate authority used when verifying server certificates.
	Ca *ServiceMonitorSpecEndpointsTlsConfigCa `field:"optional" json:"ca" yaml:"ca"`
	// Path to the CA cert in the Prometheus container to use for the targets.
	CaFile *string `field:"optional" json:"caFile" yaml:"caFile"`
	// Client certificate to present when doing client-authentication.
	Cert *ServiceMonitorSpecEndpointsTlsConfigCert `field:"optional" json:"cert" yaml:"cert"`
	// Path to the client cert file in the Prometheus container for the targets.
	CertFile *string `field:"optional" json:"certFile" yaml:"certFile"`
	// Disable target certificate validation.
	InsecureSkipVerify *bool `field:"optional" json:"insecureSkipVerify" yaml:"insecureSkipVerify"`
	// Path to the client key file in the Prometheus container for the targets.
	KeyFile *string `field:"optional" json:"keyFile" yaml:"keyFile"`
	// Secret containing the client key file for the targets.
	KeySecret *ServiceMonitorSpecEndpointsTlsConfigKeySecret `field:"optional" json:"keySecret" yaml:"keySecret"`
	// Used to verify the hostname for the targets.
	ServerName *string `field:"optional" json:"serverName" yaml:"serverName"`
}

