<p align="center">
    <img alt="wasp" src="./images/wasp-2.png">
</p>
<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/wasp)](https://goreportcard.com/report/github.com/smartcontractkit/wasp)
[![Component Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test.yml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test.yml)
[![E2E tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test-e2e.yml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test-e2e.yml)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-80%25-brightgreen.svg?longCache=true&style=flat)</a>

Scalable protocol-agnostic load testing library for `Go`

</div>

## Goals
- Easy to reuse any custom client `Go` code
- Easy to grasp
- Have a slim codebase (500-1k loc)
- No test harness or CLI, easy to integrate and run with plain `go test`
- Have a predictable performance footprint
- Easy to create synthetic or user-based scenarios
- Scalable in `k8s` without complicated configuration or vendored UI interfaces
- Non-opinionated reporting, push any data to `Loki`

## Setup
We are using `nix` for deps, see [installation](https://nixos.org/manual/nix/stable/installation/installation.html) guide
```bash
nix develop
```


## Run example tests with Grafana + Loki
```bash
make start
```
Insert `GRAFANA_TOKEN` created in previous command
```bash
export LOKI_TOKEN=
export LOKI_URL=http://localhost:3030/loki/api/v1/push
export GRAFANA_URL=http://localhost:3000
export GRAFANA_TOKEN=
export DATA_SOURCE_NAME=Loki
export DASHBOARD_FOLDER=LoadTests
export DASHBOARD_NAME=Wasp

make dashboard
```
Run some tests:
```
make test_loki
```
Open your [Grafana dashboard](http://localhost:3000/d/wasp/wasp-load-generator?orgId=1&refresh=5s)

In case you deploy to your own Grafana check `DASHBOARD_FOLDER` and `DASHBOARD_NAME`, defaults are `LoadTests` dir and dashboard is called `Wasp`

Remove environment:
```bash
make stop
```

## Test Layout and examples
Check [examples](../../../wasp/examples/README.md) to understand what is the easiest way to structure your tests, run them both locally and remotely, at scale, inside `k8s`

## Run pyroscope test
```bash
make pyro_start
make test_pyro_rps
make test_pyro_vu
make pyro_stop
```
Open [pyroscope](http://localhost:4040/)

You can also use `trace.out` in the root folder with `Go` default tracing UI

## How it works
![img.png](../../../wasp/images/how-it-works.png)

Check this [doc](../../../wasp/HOW_IT_WORKS.md) for more examples and project overview

## Loki debug
You can check all the messages the tool sends with env var `WASP_LOG_LEVEL=trace`

If Loki client fail to deliver a batch test will proceed, if you experience Loki issues, consider setting `Timeout` in `LokiConfig` or set `MaxErrors: 10` to return an error after N Loki errors

`MaxErrors: -1` can be used to ignore all the errors

Default Promtail settings are:
```golang
&LokiConfig{
    TenantID:                os.Getenv("LOKI_TENANT_ID"),
    URL:                     os.Getenv("LOKI_URL"),
    Token:                   os.Getenv("LOKI_TOKEN"),
    BasicAuth:               os.Getenv("LOKI_BASIC_AUTH"),
    MaxErrors:               10,
    BatchWait:               5 * time.Second,
    BatchSize:               500 * 1024,
    Timeout:                 20 * time.Second,
    DropRateLimitedBatches:  false,
    ExposePrometheusMetrics: false,
    MaxStreams:              600,
    MaxLineSize:             999999,
    MaxLineSizeTruncate:     false,
}
```
If you see errors like
```
ERR Malformed promtail log message, skipping Line=["level",{},"component","client","host","...","msg","batch add err","tenant","","error",{}]
```
Try to increase `MaxStreams` even more or check your `Loki` configuration


## WASP Dashboard

Basic [dashboard](../../../wasp/dashboard/dashboard.go):

![dashboard_img](./images/dashboard_basic.png)

### Reusing Dashboard Components

You can integrate components from the WASP dashboard into your custom dashboards.

Example:

```golang
import (
    waspdashboard "github.com/smartcontractkit/wasp/dashboard"
)

func BuildCustomLoadTestDashboard(dashboardName string) (dashboard.Builder, error) {
    // Custom key,value used to query for panels
    panelQuery := map[string]string{
		"branch": `=~"${branch:pipe}"`,
		"commit": `=~"${commit:pipe}"`,
        "network_type": `="testnet"`,
	}

	return dashboard.New(
		dashboardName,
        waspdashboard.WASPLoadStatsRow("Loki", panelQuery),
		waspdashboard.WASPDebugDataRow("Loki", panelQuery, true),
        # other options
    )
}
```

## Annotate Dashboards and Monitor Alerts

To enable dashboard annotations and alert monitoring, utilize the `WithGrafana()` function in conjunction with `wasp.Profile`. This approach allows for the integration of dashboard annotations and the evaluation of dashboard alerts.

Example:

```golang
_, err = wasp.NewProfile().
    WithGrafana(grafanaOpts).
    Add(wasp.NewGenerator(getLatestReportByTimestampCfg)).
    Run(true)
require.NoError(t, err)
```

Where:

```golang
type GrafanaOpts struct {
	GrafanaURL                   string        `toml:"grafana_url"`
	GrafanaToken                 string        `toml:"grafana_token_secret"`
	WaitBeforeAlertCheck         time.Duration `toml:"grafana_wait_before_alert_check"`                  // Cooldown period to wait before checking for alerts
	AnnotateDashboardUIDs        []string      `toml:"grafana_annotate_dashboard_uids"`                  // Grafana dashboardUIDs to annotate start and end of the run
	CheckDashboardAlertsAfterRun []string      `toml:"grafana_check_alerts_after_run_on_dashboard_uids"` // Grafana dashboardIds to check for alerts after run
}

```
