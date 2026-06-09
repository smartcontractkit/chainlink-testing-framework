# Logs

WASP load data is stored in `VictoriaLogs` and queried with [LogsQL](https://docs.victoriametrics.com/victorialogs/logsql/). The bundled **WASP (VictoriaLogs)** dashboard is loaded automatically, but you can also explore the raw data in [localhost:3000](http://localhost:3000/explore) (`VictoriaLogs` datasource) or directly via the VictoriaLogs UI at [localhost:9428](http://localhost:9428/select/vmui).

Every WASP record carries these fields: `go_test_name`, `gen_name`, `call_group`, `branch`, `commit` and `test_data_type` (`stats` or `responses`). Numeric values (`current_rps`, `current_instances`, `duration`, ...) live in the JSON `_msg`, so add `| unpack_json` before aggregating them.

[Explore](http://localhost:3000/explore?schemaVersion=1&panes=%7B%22il3%22:%7B%22datasource%22:%22victorialogs%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22victoriametrics-logs-datasource%22,%22uid%22:%22victorialogs%22%7D,%22editorMode%22:%22code%22,%22expr%22:%22%2A%22,%22queryType%22:%22instant%22%7D%5D,%22range%22:%7B%22from%22:%22now-5m%22,%22to%22:%22now%22%7D%7D%7D&orgId=1)

## Example WASP Log Queries

Queries:
- All data for a single test
```sql
go_test_name:="TestMyLoad"
```
- Periodic generator stats (RPS, VUs, sampling)
```sql
go_test_name:="TestMyLoad" AND test_data_type:stats
```
- Individual responses (one log line per call)
```sql
go_test_name:="TestMyLoad" AND test_data_type:responses
```
- Current RPS per generator
```sql
test_data_type:stats | unpack_json | stats by (go_test_name, gen_name) max(current_rps) value
```
- Current VUs (virtual users) per generator
```sql
test_data_type:stats | unpack_json | stats by (go_test_name, gen_name) max(current_instances) value
```
- Failed responses
```sql
test_data_type:responses AND _msg:"\"failed\":true"
```
- Timed out responses
```sql
test_data_type:responses AND _msg:"\"timeout\":true"
```
- Latency quantiles (p99/p95/p50, in nanoseconds) per call group
```sql
test_data_type:responses | unpack_json | stats by (gen_name, call_group) quantile(0.99, duration) value
```
