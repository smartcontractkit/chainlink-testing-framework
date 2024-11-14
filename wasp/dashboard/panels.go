package dashboard

import (
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
)

// RPSPanel generates a responses per second time series panel for the dashboard using the provided data source and query parameters.
// It configures visualization settings including legend placement, transparency, span, height, and axis units.
// The panel aggregates response counts from Prometheus, grouping them by node ID, test name, generator name, and call group.
// This allows monitoring of response rates segmented by generator and call group within the specified time frame.
// It returns a row.Option that can be integrated into a dashboard layout.
func RPSPanel(dataSource string, query map[string]string) row.Option {
	queryString := ""
	for key, value := range query {
		queryString += key + value + ", "
	}
	return row.WithTimeSeries(
		"Responses/sec (Generator, CallGroup)",
		timeseries.Legend(timeseries.Hide),
		timeseries.Transparent(),
		timeseries.Span(6),
		timeseries.Height("300px"),
		timeseries.DataSource(dataSource),
		timeseries.Axis(
			axis.Unit("Responses"),
			axis.Label("Responses"),
		),
		timeseries.Legend(timeseries.Bottom),
		timeseries.WithPrometheusTarget(
			`sum(count_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}", call_group=~"${call_group:pipe}"} [1s])) by (node_id, go_test_name, gen_name, call_group)`,
			prometheus.Legend("{{go_test_name}} {{gen_name}} {{call_group}} responses/sec"),
		),
		timeseries.WithPrometheusTarget(
			`sum(count_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"} [1s])) by (node_id, go_test_name, gen_name)`,
			prometheus.Legend("{{go_test_name}} Total responses/sec"),
		),
	)
}

// RPSVUPerScheduleSegmentsPanel creates a dashboard row panel that visualizes Requests Per Second (RPS) and Virtual Users (VUs) segmented by schedule. It utilizes the provided data source and query parameters to configure multiple time series widgets, fetching and displaying relevant performance metrics from Prometheus. This panel integrates seamlessly with other dashboard components to offer comprehensive insights into system performance across different schedule segments.
func RPSVUPerScheduleSegmentsPanel(dataSource string, query map[string]string) row.Option {
	queryString := ""
	for key, value := range query {
		queryString += key + value + ", "
	}
	return row.WithTimeSeries(
		"RPS/VUs per schedule segments",
		timeseries.Transparent(),
		timeseries.Span(6),
		timeseries.Height("300px"),
		timeseries.DataSource(dataSource),
		timeseries.WithPrometheusTarget(
			`
			max_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap current_rps [$__interval]) by (node_id, go_test_name, gen_name)
			`, prometheus.Legend("{{go_test_name}} {{gen_name}} RPS"),
		),
		timeseries.WithPrometheusTarget(
			`
			sum(last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap current_rps [$__interval]) by (node_id, go_test_name, gen_name))
			`,
			prometheus.Legend("{{go_test_name}} Total RPS"),
		),
		timeseries.WithPrometheusTarget(
			`
			max_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap current_instances [$__interval]) by (node_id, go_test_name, gen_name)
			`, prometheus.Legend("{{go_test_name}} {{gen_name}} VUs"),
		),
		timeseries.WithPrometheusTarget(
			`
			sum(last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap current_instances [$__interval]) by (node_id, go_test_name, gen_name))
			`,
			prometheus.Legend("{{go_test_name}} Total VUs"),
		),
	)
}
