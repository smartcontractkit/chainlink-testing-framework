package monitoringcoreoscom


// Selector to select which namespaces the Kubernetes `Endpoints` objects are discovered from.
type ServiceMonitorSpecNamespaceSelector struct {
	// Boolean describing whether all namespaces are selected in contrast to a list restricting them.
	Any *bool `field:"optional" json:"any" yaml:"any"`
	// List of namespace names to select from.
	MatchNames *[]*string `field:"optional" json:"matchNames" yaml:"matchNames"`
}

