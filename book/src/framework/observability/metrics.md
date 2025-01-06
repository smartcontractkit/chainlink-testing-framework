# Metrics

We use Prometheus to collect metrics of Chainlink nodes and other services.

Check [Prometheus UI](http://localhost:9090/query).

Check [Grafana](http://localhost:3000/explore?panes=%7B%22gGs%22:%7B%22datasource%22:%22PBFA97CFB590B2093%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22ctf%5C%22%7D%22,%22range%22:true,%22instant%22:true,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22PBFA97CFB590B2093%22%7D,%22editorMode%22:%22code%22,%22legendFormat%22:%22__auto%22%7D%5D,%22range%22:%7B%22from%22:%22now-5m%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1) example query.

Queries:
- All metrics of all containers
```json
{job="ctf"}
```

## Docker Resources

We are using [cadvisor](https://github.com/google/cadvisor) to monitor resources of test [containers](http://localhost:3000/d/pMEd7m0Mz/cadvisor-exporter?orgId=1).

Cadvisor UI can be found [here](http://localhost:8085/containers/)