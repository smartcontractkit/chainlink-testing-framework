package monitoringcoreoscom


// Specification of desired Service selection for target discovery by Prometheus.
type ServiceMonitorSpec struct {
	// Label selector to select the Kubernetes `Endpoints` objects.
	Selector *ServiceMonitorSpecSelector `field:"required" json:"selector" yaml:"selector"`
	// `attachMetadata` defines additional metadata which is added to the discovered targets.
	//
	// It requires Prometheus >= v2.37.0.
	AttachMetadata *ServiceMonitorSpecAttachMetadata `field:"optional" json:"attachMetadata" yaml:"attachMetadata"`
	// List of endpoints part of this ServiceMonitor.
	Endpoints *[]*ServiceMonitorSpecEndpoints `field:"optional" json:"endpoints" yaml:"endpoints"`
	// `jobLabel` selects the label from the associated Kubernetes `Service` object which will be used as the `job` label for all metrics.
	//
	// For example if `jobLabel` is set to `foo` and the Kubernetes `Service` object is labeled with `foo: bar`, then Prometheus adds the `job="bar"` label to all ingested metrics.
	// If the value of this field is empty or if the label doesn't exist for the given Service, the `job` label of the metrics defaults to the name of the associated Kubernetes `Service`.
	JobLabel *string `field:"optional" json:"jobLabel" yaml:"jobLabel"`
	// Per-scrape limit on the number of targets dropped by relabeling that will be kept in memory.
	//
	// 0 means no limit.
	// It requires Prometheus >= v2.47.0.
	KeepDroppedTargets *float64 `field:"optional" json:"keepDroppedTargets" yaml:"keepDroppedTargets"`
	// Per-scrape limit on number of labels that will be accepted for a sample.
	//
	// It requires Prometheus >= v2.27.0.
	LabelLimit *float64 `field:"optional" json:"labelLimit" yaml:"labelLimit"`
	// Per-scrape limit on length of labels name that will be accepted for a sample.
	//
	// It requires Prometheus >= v2.27.0.
	LabelNameLengthLimit *float64 `field:"optional" json:"labelNameLengthLimit" yaml:"labelNameLengthLimit"`
	// Per-scrape limit on length of labels value that will be accepted for a sample.
	//
	// It requires Prometheus >= v2.27.0.
	LabelValueLengthLimit *float64 `field:"optional" json:"labelValueLengthLimit" yaml:"labelValueLengthLimit"`
	// Selector to select which namespaces the Kubernetes `Endpoints` objects are discovered from.
	NamespaceSelector *ServiceMonitorSpecNamespaceSelector `field:"optional" json:"namespaceSelector" yaml:"namespaceSelector"`
	// `podTargetLabels` defines the labels which are transferred from the associated Kubernetes `Pod` object onto the ingested metrics.
	PodTargetLabels *[]*string `field:"optional" json:"podTargetLabels" yaml:"podTargetLabels"`
	// `sampleLimit` defines a per-scrape limit on the number of scraped samples that will be accepted.
	SampleLimit *float64 `field:"optional" json:"sampleLimit" yaml:"sampleLimit"`
	// `targetLabels` defines the labels which are transferred from the associated Kubernetes `Service` object onto the ingested metrics.
	TargetLabels *[]*string `field:"optional" json:"targetLabels" yaml:"targetLabels"`
	// `targetLimit` defines a limit on the number of scraped targets that will be accepted.
	TargetLimit *float64 `field:"optional" json:"targetLimit" yaml:"targetLimit"`
}

