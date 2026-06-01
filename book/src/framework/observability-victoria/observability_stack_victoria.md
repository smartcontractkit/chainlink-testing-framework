# Local Observability Stack (VictoriaMetrics)

Minimal setup for experimenting with OTEL metrics & logs in Grafana.

## Components

| Service | Port | Purpose |
|---|---|---|
| Grafana | 3000 | UI, anonymous admin enabled |
| OTEL Collector | 4317 (gRPC), 4318 (HTTP) | Receives OTLP from your app |
| VictoriaMetrics | 8428 | Metrics TSDB (Prom remote_write in, MetricsQL out) |
| VictoriaLogs | 9428 | Logs DB (OTLP in, LogsQL out) |

## Run

```bash
docker compose up -d
```

Grafana: <http://localhost:3000> (no login required, anonymous admin).

Both datasources (`VictoriaMetrics`, `VictoriaLogs`) are auto-provisioned.

## Point your Go app at it

In the OTEL exporter, use endpoint `localhost:4317` (gRPC) or `localhost:4318` (HTTP). Both metrics and logs go to the same collector — it fans them out to VM/VL.

## Tear down

```bash
# start the observability stack
ctf obs up -vm
# remove the stack with all the data (volumes)
ctf obs d -vm
# restart the stack removing all the data (volumes)
ctf obs r -vm
```

## Developing

Change compose files under `framework/cmd/observability` and restart the stack (removing volumes too)
```
just reload-cli && ctf obs r
```

## Local Dashboards (Docker)

You can create a dashboard using [UI](http://localhost:3000) and put them under `$pwd/dashboards` folder then commit, they'll be loaded automatically on start and you can find them [here](http://localhost:3000/dashboards) under `local` directory.

`$pwd` is you current working directory from which you call `ctf obs u`
