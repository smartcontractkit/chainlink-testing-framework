package dashboard

import (
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
)

// RPSPanel creates a time series panel for displaying responses per second.
// It configures the panel with a specified data source and query parameters,
// setting up the axis, legend, and Prometheus targets for visualizing response
// rates by generator and call group. The panel is styled with transparency,
// a specific span, and height, and includes legends for detailed insights.
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

// RPSVUPerScheduleSegmentsPanel creates a time series panel for visualizing
// requests per second (RPS) and virtual users (VUs) per schedule segments.
// It uses the provided data source and query parameters to configure Prometheus
// targets for displaying maximum and total RPS and VUs. The panel is styled
// with a transparent background, a span of 6, and a height of 300px.
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
