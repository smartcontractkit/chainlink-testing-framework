package monitoringcoreoscom


// Action to perform based on the regex matching.
//
// `Uppercase` and `Lowercase` actions require Prometheus >= v2.36.0. `DropEqual` and `KeepEqual` actions require Prometheus >= v2.41.0.
// Default: "Replace".
type ServiceMonitorSpecEndpointsRelabelingsAction string

const (
	// replace.
	ServiceMonitorSpecEndpointsRelabelingsAction_REPLACE ServiceMonitorSpecEndpointsRelabelingsAction = "REPLACE"
	// keep.
	ServiceMonitorSpecEndpointsRelabelingsAction_KEEP ServiceMonitorSpecEndpointsRelabelingsAction = "KEEP"
	// drop.
	ServiceMonitorSpecEndpointsRelabelingsAction_DROP ServiceMonitorSpecEndpointsRelabelingsAction = "DROP"
	// hashmod.
	ServiceMonitorSpecEndpointsRelabelingsAction_HASHMOD ServiceMonitorSpecEndpointsRelabelingsAction = "HASHMOD"
	// labelmap.
	ServiceMonitorSpecEndpointsRelabelingsAction_LABELMAP ServiceMonitorSpecEndpointsRelabelingsAction = "LABELMAP"
	// labeldrop.
	ServiceMonitorSpecEndpointsRelabelingsAction_LABELDROP ServiceMonitorSpecEndpointsRelabelingsAction = "LABELDROP"
	// labelkeep.
	ServiceMonitorSpecEndpointsRelabelingsAction_LABELKEEP ServiceMonitorSpecEndpointsRelabelingsAction = "LABELKEEP"
	// lowercase.
	ServiceMonitorSpecEndpointsRelabelingsAction_LOWERCASE ServiceMonitorSpecEndpointsRelabelingsAction = "LOWERCASE"
	// uppercase.
	ServiceMonitorSpecEndpointsRelabelingsAction_UPPERCASE ServiceMonitorSpecEndpointsRelabelingsAction = "UPPERCASE"
	// keepequal.
	ServiceMonitorSpecEndpointsRelabelingsAction_KEEPEQUAL ServiceMonitorSpecEndpointsRelabelingsAction = "KEEPEQUAL"
	// dropequal.
	ServiceMonitorSpecEndpointsRelabelingsAction_DROPEQUAL ServiceMonitorSpecEndpointsRelabelingsAction = "DROPEQUAL"
)

