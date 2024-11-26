# WASP - Testing Alerts

> [!WARNING]  
> The API used to create and check alerts is unstable, and it’s possible that this section is out of date.

With your load profile set up and running, let’s explore how you can monitor and assert on the metrics generated by the load.

This requires the use of a Grafana dashboard. For this example, we’ll generate the dashboard and alerts programmatically. However, in a real-world scenario, if you already have a dashboard with alerts, you can use them to make assertions with WASP.

We’ll use a simple HTTP `Gun` from previous examples, so its definition is skipped here for brevity. For this example, we’ll divide alerts into two groups:
* **Baseline**
* **Stress**

> [!WARNING]  
> This example assumes you have a Grafana instance set up. If you don’t, please set it up before running the test.  
> You can find more information on setting up Grafana locally [here](./how-to/start_local_observability_stack.md).  
> WASP reads Grafana configuration from environment variables. Learn more about them in the [Configuration](./configuration.md) section.

---

### Constants for Alert Groups

Let’s start by defining some constants that will be useful later:

```go
const (
    FirstGenName                 = "first API"
    SecondGenName                = "second API"
    BaselineRequirementGroupName = "baseline"
    StressRequirementGroupName   = "stress"
)
```

---

### Defining Alerts

#### Baseline Alerts

1. **99th Percentile Latency Alert**  
   This alert triggers when the 99th percentile of the response time exceeds 50ms.

```go
fiftyMsAlert := dashboard.WaspAlert{
    Name:                 "99th latency percentile is out of SLO for first API",
    AlertType:            dashboard.AlertTypeQuantile99,
    TestName:             "TestBaselineRequirements",
    GenName:              FirstGenName,
    RequirementGroupName: BaselineRequirementGroupName,
    AlertIf:              alert.IsAbove(50),
}
```

> [!NOTE]  
> Learn more about various alert types supported by WASP in the [AlertChecker](./components/alert_checker.md) documentation.

2. **Error Alert for Second API**  
   This alert fires if any errors occur while calling the second generator:

```go
anyErrorsAlert := dashboard.WaspAlert{
    Name:                 "second API has errors",
    AlertType:            dashboard.AlertTypeErrors,
    TestName:             "TestBaselineRequirements",
    GenName:              SecondGenName,
    RequirementGroupName: BaselineRequirementGroupName,
    AlertIf:              alert.IsAbove(0),
}
```

#### Stress Alerts

For the `stress` group, we’ll define a more complex custom alert:

```go
customAlert := dashboard.WaspAlert{
    RequirementGroupName: StressRequirementGroupName,
    Name:                 "MyCustomAlert",
    CustomAlert: timeseries.Alert(
        "MyCustomAlert",
        alert.For("10s"),                             // Wait 10s before considering it a firing alert
        alert.OnExecutionError(alert.ErrorAlerting),  // Set "alerting state" to "alerting" on errors
        alert.Description("My custom description"),
        alert.Tags(map[string]string{
            "service": "wasp",
            dashboard.DefaultRequirementLabelKey: StressRequirementGroupName,
        }),
        alert.WithLokiQuery(
            "MyCustomAlert",
            `
max_over_time({go_test_name="%s", test_data_type=~"stats", gen_name="%s"}
| json
| unwrap failed [10s]) by (go_test_name, gen_name)`,
        ),
        alert.If(alert.Last, "MyCustomAlert", alert.IsAbove(20)), // Trigger if ≥20 matches in the last 10s
        alert.EvaluateEvery("10s"),                               // Evaluate every 10s
    ),
}
```

---

### Building and Deploying the Dashboard

Next, we’ll build a default WASP dashboard with these alerts and deploy it to Grafana:

```go
func buildDashboard() (*dashboard.Dashboard, error) {
    d, err := dashboard.NewDashboard([]dashboard.WaspAlert{
        fiftyMsAlert,
        anyErrorsAlert,
        customAlert,
    })
    if err != nil {
        return nil, err
    }
    return d.Deploy() // Create in Grafana
}
```

---

### Writing the Test

With the dashboard deployed and the `Gun` already defined, we can create the test:

```go
func TestBaselineRequirements(t *testing.T) {
    // Start HTTP mock server
    srv := wasp.NewHTTPMockServer(
        &wasp.HTTPMockServerConfig{
            FirstAPILatency:   50 * time.Millisecond,
            FirstAPIHTTPCode:  500,
            SecondAPILatency:  50 * time.Millisecond,
            SecondAPIHTTPCode: 500,
        },
    )
    srv.Run()
    
    // Build and deploy the dashboard
    dashboard, err := buildDashboard()
    require.NoError(t, err)
    
    // Define a profile with 2 load generators
    _, err = wasp.NewProfile().
        Add(wasp.NewGenerator(&wasp.Config{
            T:          t,
            LoadType:   wasp.RPS,
            GenName:    FirstGenName,
            Schedule:   wasp.Plain(5, 20*time.Second),
            Gun:        NewExampleHTTPGun(srv.URL()),
            LokiConfig: wasp.NewEnvLokiConfig(),
        })).
        Add(wasp.NewGenerator(&wasp.Config{
            T:          t,
            LoadType:   wasp.RPS,
            GenName:    SecondGenName,
            Schedule:   wasp.Plain(5, 20*time.Second),
            Gun:        NewExampleHTTPGun(srv.URL()),
            LokiConfig: wasp.NewEnvLokiConfig(),
        })).
        Run(true)
    require.NoError(t, err)

    // Check alerts for the baseline group
    _, err = wasp.NewAlertChecker(t).AnyAlerts(dashboard.DefaultDashboardUUID, BaselineRequirementGroupName)
    require.NoError(t, err)

    // Check alerts for the stress group
    _, err = wasp.NewAlertChecker(t).AnyAlerts(dashboard.DefaultDashboardUUID, StressRequirementGroupName)
    require.NoError(t, err)
}
```

---

### Handling General Alert Checks

If you want to fail the test for any triggered alerts (not group-specific), configure your profile with the following:

```go
_, err := wasp.NewProfile().
    WithGrafana(&wasp.GrafanaOpts{
        GrafanaURL:                   os.Getenv("GRAFANA_URL"),
        GrafanaToken:                 os.Getenv("GRAFANA_TOKEN"),
        AnnotateDashboardUID:         os.Getenv("DASHBOARD_UID"),  // Annotate start/end of the test
        CheckDashboardAlertsAfterRun: os.Getenv("DASHBOARD_UID"),  // Check alerts after test ends
    }).
    // Add generator definitions
```

This setup ensures the test fails if any alerts are triggered on the dashboard after the profile finishes.

---

### Conclusion

Now you know how to write load tests and make assertions on application behavior!  
You can find the full example, including more alerts and assertions, [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/alerts).