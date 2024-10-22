package main

import (
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/dashboard"
	"os"
	"testing"
	"time"

	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/stretchr/testify/require"
)

const (
	FirstGenName                 = "first API"
	SecondGenName                = "second API"
	BaselineRequirementGroupName = "baseline"
	StressRequirementGroupName   = "stress"
)

var (
	commonLabels = map[string]string{
		"branch": "generator_healthcheck",
		"commit": "generator_healthcheck",
	}
)

func TestMain(m *testing.M) {
	srv := wasp.NewHTTPMockServer(
		&wasp.HTTPMockServerConfig{
			FirstAPILatency:   50 * time.Millisecond,
			FirstAPIHTTPCode:  500,
			SecondAPILatency:  50 * time.Millisecond,
			SecondAPIHTTPCode: 500,
		},
	)
	srv.Run()

	// we define 2 NFRs groups
	// - baseline - basic latency requirements for 99th percentiles and no errors
	// - stress - another custom NFRs for stress
	// WaspAlert can be defined on per Generator level
	// usually, you define it once per project, generate your dashboard and upload it, it's here only for example purposes
	d, err := dashboard.NewDashboard([]dashboard.WaspAlert{
		// baseline group alerts
		{
			Name:                 "99th latency percentile is out of SLO for first API",
			AlertType:            dashboard.AlertTypeQuantile99,
			TestName:             "TestBaselineRequirements",
			GenName:              FirstGenName,
			RequirementGroupName: BaselineRequirementGroupName,
			AlertIf:              alert.IsAbove(50),
		},
		{
			Name:                 "first API has errors",
			AlertType:            dashboard.AlertTypeErrors,
			TestName:             "TestBaselineRequirements",
			GenName:              FirstGenName,
			RequirementGroupName: BaselineRequirementGroupName,
			AlertIf:              alert.IsAbove(0),
		},
		{
			Name:                 "99th latency percentile is out of SLO for second API",
			AlertType:            dashboard.AlertTypeQuantile99,
			TestName:             "TestBaselineRequirements",
			GenName:              FirstGenName,
			RequirementGroupName: BaselineRequirementGroupName,
			AlertIf:              alert.IsAbove(50),
		},
		{
			Name:                 "second API has errors",
			AlertType:            dashboard.AlertTypeErrors,
			TestName:             "TestBaselineRequirements",
			GenName:              FirstGenName,
			RequirementGroupName: BaselineRequirementGroupName,
			AlertIf:              alert.IsAbove(0),
		},
		// stress group alerts
		{
			Name:                 "first API has errors > threshold",
			AlertType:            dashboard.AlertTypeErrors,
			TestName:             "TestStressRequirements",
			GenName:              FirstGenName,
			RequirementGroupName: StressRequirementGroupName,
			AlertIf:              alert.IsAbove(10),
		},
		// custom alert if you don't have some metrics on wasp dashboard, but you need those alerts
		{
			RequirementGroupName: StressRequirementGroupName,
			Name:                 "MyCustomALert",
			CustomAlert: timeseries.Alert(
				"MyCustomAlert",
				alert.For("10s"),
				alert.OnExecutionError(alert.ErrorAlerting),
				alert.Description("My custom description"),
				alert.Tags(map[string]string{
					"service": "wasp",
					// set group label so it can be filtered
					dashboard.DefaultRequirementLabelKey: StressRequirementGroupName,
				}),
				alert.WithLokiQuery(
					"MyCustomAlert",
					`
max_over_time({go_test_name="%s", test_data_type=~"stats", gen_name="%s"}
| json
| unwrap failed [10s]) by (go_test_name, gen_name)`,
				),
				alert.If(alert.Last, "MyCustomAlert", alert.IsAbove(20)),
				alert.EvaluateEvery("10s"),
			),
		},
	},
		// here you can add additional dashboard fields for your project, example:
		// []dashboard.Option{
		//			dashboard.Row("my_custom_logs", row.WithLogs(
		//				"Timed out responses",
		//				logs.DataSource("Loki"),
		//				logs.Span(6),
		//				logs.Height("300px"),
		//				logs.Transparent(),
		//				logs.WithLokiTarget(`
		//	{go_test_name=~"${go_test_name:pipe}", test_data_type=~"responses", branch=~"${branch:pipe}", commit=~"${commit:pipe}", gen_name=~"${gen_name:pipe}"} |~ "timeout\":true"`),
		//			)),
		//		}
		nil)
	if err != nil {
		panic(err)
	}
	if _, err := d.Deploy(); err != nil {
		panic(err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestBaselineRequirements(t *testing.T) {
	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    FirstGenName,
			Schedule:   wasp.Plain(5, 20*time.Second),
			Gun:        NewExampleHTTPGun("http://localhost:8080/1"),
			Labels:     commonLabels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    SecondGenName,
			Schedule:   wasp.Plain(5, 20*time.Second),
			Gun:        NewExampleHTTPGun("http://localhost:8080/2"),
			Labels:     commonLabels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	require.NoError(t, err)

	// we are checking all active alerts for dashboard with UUID = "wasp" which have label "requirement_name" = "baseline"
	// if any alerts of particular group, for example "baseline" were raised - we fail the test
	// change some data in NewHTTPMockServer to make alerts disappear
	_, err = wasp.NewAlertChecker(t).AnyAlerts(dashboard.DefaultDashboardUUID, BaselineRequirementGroupName)
	require.NoError(t, err)
}

func TestStressRequirements(t *testing.T) {
	// we are testing the same APIs but for different NFRs group
	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    FirstGenName,
			Schedule:   wasp.Plain(5, 20*time.Second),
			Gun:        NewExampleHTTPGun("http://localhost:8080/1"),
			Labels:     commonLabels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    SecondGenName,
			Schedule:   wasp.Plain(5, 20*time.Second),
			Gun:        NewExampleHTTPGun("http://localhost:8080/2"),
			Labels:     commonLabels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	require.NoError(t, err)

	// we are checking all active alerts for dashboard with UUID = "wasp" which have label "requirement_name" = "stress"
	// if any alerts of particular group, for example "stress" were raised - we fail the test
	// change some data in NewHTTPMockServer to make alerts disappear
	_, err = wasp.NewAlertChecker(t).AnyAlerts(dashboard.DefaultDashboardUUID, StressRequirementGroupName)
	require.NoError(t, err)
}
