package dashboard

import (
	"fmt"
	"net/http"
	"os"

	"context"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/query"
)

const (
	DefaultStatTextSize       = 12
	DefaultStatValueSize      = 20
	DefaultAlertEvaluateEvery = "10s"
	DefaultAlertFor           = "10s"
	DefaultDashboardUUID      = "Wasp"

	DefaultRequirementLabelKey = "requirement_name"
)

const (
	AlertTypeQuantile99 = "quantile_99"
	AlertTypeErrors     = "errors"
	AlertTypeTimeouts   = "timeouts"
)

type WaspAlert struct {
	Name                 string
	AlertType            string
	TestName             string
	GenName              string
	RequirementGroupName string
	AlertIf              alert.ConditionEvaluator
	CustomAlert          timeseries.Option
}

// Dashboard is a Wasp dashboard
type Dashboard struct {
	Name           string
	DataSourceName string
	Folder         string
	GrafanaURL     string
	GrafanaToken   string
	extendedOpts   []dashboard.Option
	builder        dashboard.Builder
}

// NewDashboard creates new dashboard
func NewDashboard(reqs []WaspAlert, opts []dashboard.Option) (*Dashboard, error) {
	name := os.Getenv("DASHBOARD_NAME")
	if name == "" {
		return nil, fmt.Errorf("DASHBOARD_NAME must be provided")
	}
	dsn := os.Getenv("DATA_SOURCE_NAME")
	if dsn == "" {
		return nil, fmt.Errorf("DATA_SOURCE_NAME must be provided")
	}
	dbf := os.Getenv("DASHBOARD_FOLDER")
	if dbf == "" {
		return nil, fmt.Errorf("DASHBOARD_FOLDER must be provided")
	}
	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		return nil, fmt.Errorf("GRAFANA_URL must be provided")
	}
	grafanaToken := os.Getenv("GRAFANA_TOKEN")
	if grafanaToken == "" {
		return nil, fmt.Errorf("GRAFANA_TOKEN must be provided")
	}
	dash := &Dashboard{
		Name:           name,
		DataSourceName: dsn,
		Folder:         dbf,
		GrafanaURL:     grafanaURL,
		GrafanaToken:   grafanaToken,
		extendedOpts:   opts,
	}
	err := dash.Build(name, dsn, reqs)
	if err != nil {
		return nil, fmt.Errorf("failed to build dashboard: %s", err)
	}
	return dash, nil
}

// Deploy deploys this dashboard to some Grafana folder
func (m *Dashboard) Deploy() (*grabana.Dashboard, error) {
	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, m.GrafanaURL, grabana.WithAPIToken(m.GrafanaToken))
	fo, err := client.FindOrCreateFolder(ctx, m.Folder)
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}
	return client.UpsertDashboard(ctx, fo, m.builder)
}

// defaultStatWidget creates default Stat widget
func defaultStatWidget(name, datasourceName, target, legend string) row.Option {
	return row.WithStat(
		name,
		stat.Transparent(),
		stat.DataSource(datasourceName),
		stat.Text(stat.TextValueAndName),
		stat.Orientation(stat.OrientationHorizontal),
		stat.TitleFontSize(DefaultStatTextSize),
		stat.ValueFontSize(DefaultStatValueSize),
		stat.Span(2),
		stat.WithPrometheusTarget(target, prometheus.Legend(legend)),
	)
}

// defaultLastValueAlertWidget creates default last value alert
func defaultLastValueAlertWidget(a WaspAlert) timeseries.Option {
	if a.CustomAlert != nil {
		return a.CustomAlert
	}
	return timeseries.Alert(
		a.Name,
		alert.For(DefaultAlertFor),
		alert.OnExecutionError(alert.ErrorKO),
		alert.Description(a.Name),
		alert.Tags(map[string]string{
			"service":                  "wasp",
			DefaultRequirementLabelKey: a.RequirementGroupName,
		}),
		alert.WithLokiQuery(
			a.Name,
			InlineLokiAlertParams(a.AlertType, a.TestName, a.GenName),
		),
		alert.If(alert.Last, a.Name, a.AlertIf),
		alert.EvaluateEvery(DefaultAlertEvaluateEvery),
	)
}

// defaultLabelValuesVar creates a dashboard variable with All/Multiple options
func defaultLabelValuesVar(name, datasourceName string) dashboard.Option {
	return dashboard.VariableAsQuery(
		name,
		query.DataSource(datasourceName),
		query.Multiple(),
		query.IncludeAll(),
		query.Request(fmt.Sprintf("label_values(%s)", name)),
		query.Sort(query.NumericalAsc),
	)
}

// timeSeriesWithAlerts creates timeseries graphs per alert + definition of alert
func timeSeriesWithAlerts(datasourceName string, alertDefs []WaspAlert) []dashboard.Option {
	dashboardOpts := make([]dashboard.Option, 0)
	for _, a := range alertDefs {
		// for wasp metrics we also create additional row per alert
		tsOpts := []timeseries.Option{
			timeseries.Transparent(),
			timeseries.Span(12),
			timeseries.Height("200px"),
			timeseries.DataSource(datasourceName),
			timeseries.Legend(timeseries.Bottom),
		}
		tsOpts = append(tsOpts, defaultLastValueAlertWidget(a))

		var rowTitle string
		// for wasp metrics we also create additional row per alert
		if a.CustomAlert == nil {
			rowTitle = fmt.Sprintf("Alert: %s, Requirement: %s", a.Name, a.RequirementGroupName)
			tsOpts = append(tsOpts, timeseries.WithPrometheusTarget(InlineLokiAlertParams(a.AlertType, a.TestName, a.GenName)))
		} else {
			rowTitle = fmt.Sprintf("External alert: %s, Requirement: %s", a.Name, a.RequirementGroupName)
		}
		// all the other custom alerts may burden the dashboard,
		dashboardOpts = append(dashboardOpts,
			dashboard.Row(
				rowTitle,
				row.Collapse(),
				row.HideTitle(),
				row.WithTimeSeries(a.Name, tsOpts...),
			))
	}
	return dashboardOpts
}

func AddVariables(datasourceName string) []dashboard.Option {
	opts := []dashboard.Option{
		defaultLabelValuesVar("go_test_name", datasourceName),
		defaultLabelValuesVar("gen_name", datasourceName),
		defaultLabelValuesVar("branch", datasourceName),
		defaultLabelValuesVar("commit", datasourceName),
		defaultLabelValuesVar("call_group", datasourceName),
	}
	return opts
}

// dashboard is internal appendable representation of all Dashboard widgets
func (m *Dashboard) dashboard(datasourceName string, requirements []WaspAlert) []dashboard.Option {
	panelQuery := map[string]string{
		"branch": `=~"${branch:pipe}"`,
		"commit": `=~"${commit:pipe}"`,
	}

	defaultOpts := []dashboard.Option{
		dashboard.UID(m.Name),
		dashboard.AutoRefresh("5"),
		dashboard.Time("now-30m", "now"),
		dashboard.Tags([]string{"generated", "load-test"}),
	}
	defaultOpts = append(defaultOpts, AddVariables(datasourceName)...)
	defaultOpts = append(defaultOpts, WASPLoadStatsRow(datasourceName, panelQuery))
	defaultOpts = append(defaultOpts, WASPDebugDataRow(datasourceName, panelQuery, false))
	defaultOpts = append(defaultOpts, timeSeriesWithAlerts(datasourceName, requirements)...)
	defaultOpts = append(defaultOpts, m.extendedOpts...)
	return defaultOpts
}

// Build creates dashboard instance
func (m *Dashboard) Build(dashboardName, datasourceName string, requirements []WaspAlert) error {
	b, err := dashboard.New(
		dashboardName,
		m.dashboard(datasourceName, requirements)...,
	)
	if err != nil {
		return fmt.Errorf("failed to create a dashboard builder: %s", err)
	}
	m.builder = b
	return nil
}

// JSON render dashboard as JSON
func (m *Dashboard) JSON() ([]byte, error) {
	return m.builder.MarshalIndentJSON()
}

// InlineLokiAlertParams is specific params for predefined alerts for wasp dashboard
func InlineLokiAlertParams(queryType, testName, genName string) string {
	switch queryType {
	case AlertTypeQuantile99:
		return fmt.Sprintf(`
avg(quantile_over_time(0.99, {go_test_name="%s", test_data_type=~"responses", gen_name="%s"}
| json
| unwrap duration [10s]) / 1e6)`, testName, genName)
	case AlertTypeErrors:
		return fmt.Sprintf(`
max_over_time({go_test_name="%s", test_data_type=~"stats", gen_name="%s"}
| json
| unwrap failed [10s]) by (go_test_name, gen_name)`, testName, genName)
	case AlertTypeTimeouts:
		return fmt.Sprintf(`
max_over_time({go_test_name="%s", test_data_type=~"stats", gen_name="%s"}
| json
| unwrap callTimeout [10s]) by (go_test_name, gen_name)`, testName, genName)
	default:
		return ""
	}
}

func WASPLoadStatsRow(dataSource string, query map[string]string) dashboard.Option {
	queryString := ""
	for key, value := range query {
		queryString += key + value + ", "
	}

	return dashboard.Row(
		"WASP Load Stats",
		defaultStatWidget(
			"RPS (Now)",
			dataSource,
			`sum(last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap current_rps [1s]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)`,
			`{{go_test_name}} {{gen_name}} RPS`,
		),
		defaultStatWidget(
			"VUs (Now)",
			dataSource,
			`sum(max_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap current_instances [$__range]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)`,
			`{{go_test_name}} {{gen_name}} VUs`,
		),
		defaultStatWidget(
			"Responses/sec (Now)",
			dataSource,
			`sum(count_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"} [1s])) by (node_id, go_test_name, gen_name)`,
			`{{go_test_name}} {{gen_name}} Responses/sec`,
		),
		defaultStatWidget(
			"Successful requests (Total)",
			dataSource,
			`
			sum(max_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap success [$__range]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)
			`,
			`{{go_test_name}} {{gen_name}} Successful requests`,
		),
		defaultStatWidget(
			"Errored requests (Total)",
			dataSource,
			`
			sum(max_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap failed [$__range]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)
			`,
			`{{go_test_name}} {{gen_name}} Errored requests`,
		),
		defaultStatWidget(
			"Timed out requests (Total)",
			dataSource,
			`
			sum(max_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
			| json
			| unwrap callTimeout [$__range]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)
			`,
			`{{go_test_name}} {{gen_name}} Timed out requests`,
		),
		RPSVUPerScheduleSegmentsPanel(dataSource, query),
		RPSPanel(dataSource, query),
		row.WithTimeSeries(
			"Latency quantiles over groups (99, 95, 50)",
			timeseries.Legend(timeseries.Hide),
			timeseries.Transparent(),
			timeseries.Span(6),
			timeseries.Height("300px"),
			timeseries.DataSource(dataSource),
			timeseries.Legend(timeseries.Bottom),
			timeseries.Axis(
				axis.Unit("ms"),
				axis.Label("ms"),
			),
			timeseries.WithPrometheusTarget(`
				quantile_over_time(0.99, {`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"}
				| json
				| unwrap duration [$__interval]) by (go_test_name, gen_name) / 1e6`,
				prometheus.Legend("{{go_test_name}} {{gen_name}} Q 99 - {{error}}"),
			),
			timeseries.WithPrometheusTarget(`
				quantile_over_time(0.95, {`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"}
				| json
				| unwrap duration [$__interval]) by (go_test_name, gen_name) / 1e6`,
				prometheus.Legend("{{go_test_name}} {{gen_name}} Q 95 - {{error}}"),
			),
			timeseries.WithPrometheusTarget(
				`
				quantile_over_time(0.50, {`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"}
				| json
				| unwrap duration [$__interval]) by (go_test_name, gen_name) / 1e6`,
				prometheus.Legend("{{go_test_name}} {{gen_name}} Q 50 - {{error}}"),
			),
		),
		row.WithTimeSeries(
			"Responses latencies by types over time (Generator, CallGroup)",
			timeseries.Legend(timeseries.Hide),
			timeseries.Transparent(),
			timeseries.Span(6),
			timeseries.Height("300px"),
			timeseries.DataSource(dataSource),
			timeseries.Axis(
				axis.Unit("ms"),
				axis.Label("ms"),
			),
			timeseries.Legend(timeseries.Bottom),
			timeseries.WithPrometheusTarget(`
				last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}", call_group=~"${call_group}"}
				| json
				| unwrap duration [$__interval]) / 1e6`,
				prometheus.Legend("{{go_test_name}} {{gen_name}} {{call_group}} T: {{timeout}} E: {{error}}"),
			),
			timeseries.WithPrometheusTarget(`
				last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"}
				| json
				| unwrap duration [$__interval]) / 1e6`,
				prometheus.Legend("{{go_test_name}} {{gen_name}} all groups T: {{timeout}} E: {{error}}"),
			),
		),
	)
}

func WASPDebugDataRow(dataSource string, query map[string]string, collapse bool) dashboard.Option {
	queryString := ""
	for key, value := range query {
		queryString += key + value + ", "
	}

	defaultRowOpts := []row.Option{}
	if collapse {
		defaultRowOpts = append(defaultRowOpts, row.Collapse())
	}

	return dashboard.Row(
		"WASP Debug",
		append(defaultRowOpts,
			row.WithStat(
				"Latest segment stats",
				stat.Transparent(),
				stat.DataSource(dataSource),
				stat.Text(stat.TextValueAndName),
				stat.SparkLine(),
				stat.Span(12),
				stat.Height("100px"),
				stat.ColorValue(),
				stat.WithPrometheusTarget(`
                sum(bytes_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", gen_name=~"${gen_name:pipe}"} [$__range]) * 1e-6)
                `, prometheus.Legend("Overall logs size"),
				),
				stat.WithPrometheusTarget(`
                sum(bytes_rate({`+queryString+`go_test_name=~"${go_test_name:pipe}", gen_name=~"${gen_name:pipe}"} [$__interval]) * 1e-6)
                `, prometheus.Legend("Logs size per second"),
				),
			),
			row.WithTimeSeries(
				"CallResult sampling (successful results)",
				timeseries.Transparent(),
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(dataSource),
				timeseries.Axis(
					axis.Label("CallResults"),
				),
				timeseries.WithPrometheusTarget(`
                sum(last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
                | json
                | unwrap samples_recorded [$__interval])) by (go_test_name, gen_name)
                `, prometheus.Legend("{{go_test_name}} {{gen_name}} recorded"),
				),
				timeseries.WithPrometheusTarget(`
                sum(last_over_time({`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}
                | json
                | unwrap samples_skipped [$__interval])) by (go_test_name, gen_name)
                `, prometheus.Legend("{{go_test_name}} {{gen_name}} skipped"),
				),
			),
			row.WithLogs(
				"Stats logs",
				logs.DataSource(dataSource),
				logs.Span(12),
				logs.Height("300px"),
				logs.Transparent(),
				logs.WithLokiTarget(`{`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"stats", gen_name=~"${gen_name:pipe}"}`),
			),
			row.WithLogs(
				"Failed responses",
				logs.DataSource(dataSource),
				logs.Span(6),
				logs.Height("300px"),
				logs.Transparent(),
				logs.WithLokiTarget(`{`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"} |~ "failed\":true"`),
			),
			row.WithLogs(
				"Timed out responses",
				logs.DataSource(dataSource),
				logs.Span(6),
				logs.Height("300px"),
				logs.Transparent(),
				logs.WithLokiTarget(`{`+queryString+`go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", gen_name=~"${gen_name:pipe}"} |~ "timeout\":true"`),
			),
		)...,
	)
}
