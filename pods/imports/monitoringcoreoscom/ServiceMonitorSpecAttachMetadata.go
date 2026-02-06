package monitoringcoreoscom


// `attachMetadata` defines additional metadata which is added to the discovered targets.
//
// It requires Prometheus >= v2.37.0.
type ServiceMonitorSpecAttachMetadata struct {
	// When set to true, Prometheus must have the `get` permission on the `Nodes` objects.
	Node *bool `field:"optional" json:"node" yaml:"node"`
}

