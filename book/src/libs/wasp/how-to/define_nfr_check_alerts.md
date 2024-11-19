# WASP - How to Define NFRs and Check Alerts

WASP allows you to define and monitor **Non-Functional Requirements (NFRs)** by grouping them into categories for independent checks. Ideally, these NFRs should be defined as dashboard-as-code and committed to your Git repository. However, you can also create them programmatically on an existing or new dashboard.

---

### Defining Alerts

WASP supports two types of alerts:

#### Built-in Alerts

These alerts are automatically supported by every generator and are based on the following metrics:
* **99th quantile (p99)**  
* **Errors**  
* **Timeouts**  

You can specify simple alert conditions such as:
* Value above or below
* Average
* Median  

These conditions use Grabana's [ConditionEvaluator](https://pkg.go.dev/github.com/K-Phoen/grabana@v0.21.18/alert#ConditionEvaluator).

---

#### Custom Alerts

Custom alerts can be much more complex and can:
* Combine multiple simple conditions.
* Execute Loki queries.

Custom alerts use Grabana's [timeseries.Alert](https://pkg.go.dev/github.com/K-Phoen/grabana@v0.21.18/timeseries#Alert) and must be timeseries-based.

> [!NOTE]  
> For a programmatic example, check the [alerts example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/wasp/examples/alerts/main_test.go).

---

### Checking Alerts

Once you have a dashboard with alerts, you can choose between two approaches:

#### 1. Automatic Checking with `Profile` and `GrafanaOpts`

You can automatically check alerts and add dashboard annotations by utilizing the `WithGrafana()` function with a `wasp.Profile`. This approach integrates dashboard annotations and evaluates alerts after the test run.

**Example**:

```go
_, err = wasp.NewProfile().
    WithGrafana(grafanaOpts).
    Add(wasp.NewGenerator(getLatestReportByTimestampCfg)).
    Run(true)
require.NoError(t, err)
```

Where `GrafanaOpts` is defined as:

```go
type GrafanaOpts struct {
	GrafanaURL                   string        `toml:"grafana_url"`
	GrafanaToken                 string        `toml:"grafana_token_secret"`
	WaitBeforeAlertCheck         time.Duration `toml:"grafana_wait_before_alert_check"`                  // Cooldown period before checking for alerts
	AnnotateDashboardUIDs        []string      `toml:"grafana_annotate_dashboard_uids"`                  // Dashboard UIDs to annotate the start and end of the run
	CheckDashboardAlertsAfterRun []string      `toml:"grafana_check_alerts_after_run_on_dashboard_uids"` // Dashboard UIDs to check for alerts after the run
}
```

#### 2. Manual Checking with `AlertChecker`

You can manually check alerts using the [AlertChecker](../components/alert_checker.md) component.

---

### Summary

WASP provides flexibility in defining and monitoring NFRs:
* Use **built-in alerts** for standard metrics like p99, errors, and timeouts.
* Use **custom alerts** for complex conditions or Loki queries.
* Automate alert checks using `Profile` with `GrafanaOpts` or manually verify them with `AlertChecker`.

By combining these tools, you can ensure your application's performance and reliability align with your NFRs.