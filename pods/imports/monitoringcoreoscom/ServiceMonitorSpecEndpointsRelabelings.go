package monitoringcoreoscom


// RelabelConfig allows dynamic rewriting of the label set for targets, alerts, scraped samples and remote write samples.
//
// More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
type ServiceMonitorSpecEndpointsRelabelings struct {
	// Action to perform based on the regex matching.
	//
	// `Uppercase` and `Lowercase` actions require Prometheus >= v2.36.0. `DropEqual` and `KeepEqual` actions require Prometheus >= v2.41.0.
	// Default: "Replace".
	Action ServiceMonitorSpecEndpointsRelabelingsAction `field:"optional" json:"action" yaml:"action"`
	// Modulus to take of the hash of the source label values.
	//
	// Only applicable when the action is `HashMod`.
	Modulus *float64 `field:"optional" json:"modulus" yaml:"modulus"`
	// Regular expression against which the extracted value is matched.
	Regex *string `field:"optional" json:"regex" yaml:"regex"`
	// Replacement value against which a Replace action is performed if the regular expression matches.
	//
	// Regex capture groups are available.
	Replacement *string `field:"optional" json:"replacement" yaml:"replacement"`
	// Separator is the string between concatenated SourceLabels.
	Separator *string `field:"optional" json:"separator" yaml:"separator"`
	// The source labels select values from existing labels.
	//
	// Their content is concatenated using the configured Separator and matched against the configured regular expression.
	SourceLabels *[]*string `field:"optional" json:"sourceLabels" yaml:"sourceLabels"`
	// Label to which the resulting string is written in a replacement.
	//
	// It is mandatory for `Replace`, `HashMod`, `Lowercase`, `Uppercase`, `KeepEqual` and `DropEqual` actions.
	// Regex capture groups are available.
	TargetLabel *string `field:"optional" json:"targetLabel" yaml:"targetLabel"`
}

