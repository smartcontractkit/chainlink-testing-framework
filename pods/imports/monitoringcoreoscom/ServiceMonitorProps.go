package monitoringcoreoscom

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// ServiceMonitor defines monitoring for a set of services.
type ServiceMonitorProps struct {
	// Specification of desired Service selection for target discovery by Prometheus.
	Spec *ServiceMonitorSpec `field:"required" json:"spec" yaml:"spec"`
	Metadata *cdk8s.ApiObjectMetadata `field:"optional" json:"metadata" yaml:"metadata"`
}

