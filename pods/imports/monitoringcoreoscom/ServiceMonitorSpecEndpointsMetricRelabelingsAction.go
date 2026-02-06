package monitoringcoreoscom


// Action to perform based on the regex matching.
//
// `Uppercase` and `Lowercase` actions require Prometheus >= v2.36.0. `DropEqual` and `KeepEqual` actions require Prometheus >= v2.41.0.
// Default: "Replace".
type ServiceMonitorSpecEndpointsMetricRelabelingsAction string

const (
	// replace.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_REPLACE ServiceMonitorSpecEndpointsMetricRelabelingsAction = "REPLACE"
	// keep.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_KEEP ServiceMonitorSpecEndpointsMetricRelabelingsAction = "KEEP"
	// drop.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_DROP ServiceMonitorSpecEndpointsMetricRelabelingsAction = "DROP"
	// hashmod.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_HASHMOD ServiceMonitorSpecEndpointsMetricRelabelingsAction = "HASHMOD"
	// labelmap.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_LABELMAP ServiceMonitorSpecEndpointsMetricRelabelingsAction = "LABELMAP"
	// labeldrop.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_LABELDROP ServiceMonitorSpecEndpointsMetricRelabelingsAction = "LABELDROP"
	// labelkeep.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_LABELKEEP ServiceMonitorSpecEndpointsMetricRelabelingsAction = "LABELKEEP"
	// lowercase.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_LOWERCASE ServiceMonitorSpecEndpointsMetricRelabelingsAction = "LOWERCASE"
	// uppercase.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_UPPERCASE ServiceMonitorSpecEndpointsMetricRelabelingsAction = "UPPERCASE"
	// keepequal.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_KEEPEQUAL ServiceMonitorSpecEndpointsMetricRelabelingsAction = "KEEPEQUAL"
	// dropequal.
	ServiceMonitorSpecEndpointsMetricRelabelingsAction_DROPEQUAL ServiceMonitorSpecEndpointsMetricRelabelingsAction = "DROPEQUAL"
)

