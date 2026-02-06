package monitoringcoreoscom


// HTTP scheme to use for scraping.
//
// `http` and `https` are the expected values unless you rewrite the `__scheme__` label via relabeling.
// If empty, Prometheus uses the default value `http`.
type ServiceMonitorSpecEndpointsScheme string

const (
	// http.
	ServiceMonitorSpecEndpointsScheme_HTTP ServiceMonitorSpecEndpointsScheme = "HTTP"
	// https.
	ServiceMonitorSpecEndpointsScheme_HTTPS ServiceMonitorSpecEndpointsScheme = "HTTPS"
)

