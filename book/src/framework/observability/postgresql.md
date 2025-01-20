# PostgreSQL

## Metrics

You can go to [this](http://localhost:3000/d/000000039/postgresql-database?orgId=1&refresh=10s&var-DS_PROMETHEUS=PBFA97CFB590B2093&var-interval=$__auto_interval_interval&var-namespace=&var-release=&var-instance=postgres_exporter_0:9187&var-datname=All&var-mode=All&from=now-5m&to=now) dashboard and select `Instance`.

We run exporter for the first 5 nodes, that should be enough to debug performance issues.

## Data sources and raw queries

You can use `PostgreSQL X` datasources, first 5 nodes to select things for your [dashboards](http://localhost:3000/explore?panes=%7B%22qrr%22:%7B%22datasource%22:%22P4DD770FAD7295D26%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22postgres%22,%22uid%22:%22P4DD770FAD7295D26%22%7D,%22format%22:%22table%22,%22rawSql%22:%22select%20%2A%20from%20sessions;%22,%22editorMode%22:%22code%22,%22sql%22:%7B%22columns%22:%5B%7B%22type%22:%22function%22,%22parameters%22:%5B%5D%7D%5D,%22groupBy%22:%5B%7B%22type%22:%22groupBy%22,%22property%22:%7B%22type%22:%22string%22%7D%7D%5D,%22limit%22:50%7D,%22rawQuery%22:true%7D%5D,%22range%22:%7B%22from%22:%22now-6h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1) or debug.

## Slow queries debug

Use `pg_stat_statements` extension, it is enabled for all databases in `NodeSet`
```
select * from pg_stat_statements;
```
