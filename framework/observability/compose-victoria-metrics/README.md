# Local VictoriaMetrics + VictoriaLogs + OTEL stack

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
docker compose down          # keep data
docker compose down -v       # wipe volumes
```

## Quick sanity checks

```bash
# Metrics ingest
curl http://localhost:8428/api/v1/query?query=up

# Logs ingest
curl 'http://localhost:9428/select/logsql/query?query=*&limit=10'
```
