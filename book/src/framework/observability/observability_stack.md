# Local Observability Stack

You can use a local observability stack, framework is connected to it by default

```bash
ctf obs up
```

To remove it use

```bash
ctf obs down
```

Read more about how to check [logs](logs.md) and [profiles](profiling.md)

## Helpful links for your local setup

- [Loki](http://localhost:3000/explore?panes=%7B%22qZw%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-6h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1)
- [Prometheus](http://localhost:3000/explore?panes=%7B%22qZw%22:%7B%22datasource%22:%22PBFA97CFB590B2093%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%22,%22range%22:true,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22PBFA97CFB590B2093%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-6h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1)
- [PostgreSQL](http://localhost:3000/d/000000039/postgresql-database?orgId=1&refresh=10s&var-DS_PROMETHEUS=PBFA97CFB590B2093&var-interval=$__auto_interval_interval&var-namespace=&var-release=&var-instance=postgres_exporter_0:9187&var-datname=All&var-mode=All&from=now-5m&to=now)
- [Pyroscope](http://localhost:4040)
