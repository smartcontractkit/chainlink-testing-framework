package monitoringcoreoscom


// Endpoint defines an endpoint serving Prometheus metrics to be scraped by Prometheus.
type ServiceMonitorSpecEndpoints struct {
	// `authorization` configures the Authorization header credentials to use when scraping the target.
	//
	// Cannot be set at the same time as `basicAuth`, or `oauth2`.
	Authorization *ServiceMonitorSpecEndpointsAuthorization `field:"optional" json:"authorization" yaml:"authorization"`
	// `basicAuth` configures the Basic Authentication credentials to use when scraping the target.
	//
	// Cannot be set at the same time as `authorization`, or `oauth2`.
	BasicAuth *ServiceMonitorSpecEndpointsBasicAuth `field:"optional" json:"basicAuth" yaml:"basicAuth"`
	// File to read bearer token for scraping the target.
	//
	// Deprecated: use `authorization` instead.
	BearerTokenFile *string `field:"optional" json:"bearerTokenFile" yaml:"bearerTokenFile"`
	// `bearerTokenSecret` specifies a key of a Secret containing the bearer token for scraping targets.
	//
	// The secret needs to be in the same namespace as the ServiceMonitor object and readable by the Prometheus Operator.
	// Deprecated: use `authorization` instead.
	BearerTokenSecret *ServiceMonitorSpecEndpointsBearerTokenSecret `field:"optional" json:"bearerTokenSecret" yaml:"bearerTokenSecret"`
	// `enableHttp2` can be used to disable HTTP2 when scraping the target.
	EnableHttp2 *bool `field:"optional" json:"enableHttp2" yaml:"enableHttp2"`
	// When true, the pods which are not running (e.g. either in Failed or Succeeded state) are dropped during the target discovery. If unset, the filtering is enabled. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase.
	FilterRunning *bool `field:"optional" json:"filterRunning" yaml:"filterRunning"`
	// `followRedirects` defines whether the scrape requests should follow HTTP 3xx redirects.
	FollowRedirects *bool `field:"optional" json:"followRedirects" yaml:"followRedirects"`
	// When true, `honorLabels` preserves the metric's labels when they collide with the target's labels.
	HonorLabels *bool `field:"optional" json:"honorLabels" yaml:"honorLabels"`
	// `honorTimestamps` controls whether Prometheus preserves the timestamps when exposed by the target.
	HonorTimestamps *bool `field:"optional" json:"honorTimestamps" yaml:"honorTimestamps"`
	// Interval at which Prometheus scrapes the metrics from the target.
	//
	// If empty, Prometheus uses the global scrape interval.
	Interval *string `field:"optional" json:"interval" yaml:"interval"`
	// `metricRelabelings` configures the relabeling rules to apply to the samples before ingestion.
	MetricRelabelings *[]*ServiceMonitorSpecEndpointsMetricRelabelings `field:"optional" json:"metricRelabelings" yaml:"metricRelabelings"`
	// `oauth2` configures the OAuth2 settings to use when scraping the target.
	//
	// It requires Prometheus >= 2.27.0.
	// Cannot be set at the same time as `authorization`, or `basicAuth`.
	Oauth2 *ServiceMonitorSpecEndpointsOauth2 `field:"optional" json:"oauth2" yaml:"oauth2"`
	// params define optional HTTP URL parameters.
	Params *map[string]*[]*string `field:"optional" json:"params" yaml:"params"`
	// HTTP path from which to scrape for metrics.
	//
	// If empty, Prometheus uses the default value (e.g. `/metrics`).
	Path *string `field:"optional" json:"path" yaml:"path"`
	// Name of the Service port which this endpoint refers to.
	//
	// It takes precedence over `targetPort`.
	Port *string `field:"optional" json:"port" yaml:"port"`
	// `proxyURL` configures the HTTP Proxy URL (e.g. "http://proxyserver:2195") to go through when scraping the target.
	ProxyUrl *string `field:"optional" json:"proxyUrl" yaml:"proxyUrl"`
	// `relabelings` configures the relabeling rules to apply the target's metadata labels.
	//
	// The Operator automatically adds relabelings for a few standard Kubernetes fields.
	// The original scrape job's name is available via the `__tmp_prometheus_job_name` label.
	// More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
	Relabelings *[]*ServiceMonitorSpecEndpointsRelabelings `field:"optional" json:"relabelings" yaml:"relabelings"`
	// HTTP scheme to use for scraping.
	//
	// `http` and `https` are the expected values unless you rewrite the `__scheme__` label via relabeling.
	// If empty, Prometheus uses the default value `http`.
	Scheme ServiceMonitorSpecEndpointsScheme `field:"optional" json:"scheme" yaml:"scheme"`
	// Timeout after which Prometheus considers the scrape to be failed.
	//
	// If empty, Prometheus uses the global scrape timeout unless it is less than the target's scrape interval value in which the latter is used.
	ScrapeTimeout *string `field:"optional" json:"scrapeTimeout" yaml:"scrapeTimeout"`
	// Name or number of the target port of the `Pod` object behind the Service, the port must be specified with container port property.
	//
	// Deprecated: use `port` instead.
	TargetPort ServiceMonitorSpecEndpointsTargetPort `field:"optional" json:"targetPort" yaml:"targetPort"`
	// TLS configuration to use when scraping the target.
	TlsConfig *ServiceMonitorSpecEndpointsTlsConfig `field:"optional" json:"tlsConfig" yaml:"tlsConfig"`
	// `trackTimestampsStaleness` defines whether Prometheus tracks staleness of the metrics that have an explicit timestamp present in scraped data.
	//
	// Has no effect if `honorTimestamps` is false.
	// It requires Prometheus >= v2.48.0.
	TrackTimestampsStaleness *bool `field:"optional" json:"trackTimestampsStaleness" yaml:"trackTimestampsStaleness"`
}

