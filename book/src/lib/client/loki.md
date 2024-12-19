# Loki

Loki client simplifies querying of Loki logs with `LogQL`.

The way it's designed now implies that:
* you need to create a new client instance for each query
* query results are returned as `string`

## New instance
To create a new instance you to provide the following at the very least:
* Loki URL
* query to execute
* time range

```go
// scheme is required
lokiUrl := "http://loki-host.io"
// can be empty
tenantId := "promtail"
basicAuth := LokiBasicAuth{
    Login: "admin",
    Password: "oh-so-secret",
}
queryParams := LokiQueryParams{
				Query:     "quantile_over_time(0.5, {name='my awesome app'} | json| unwrap duration [10s]) by name",
				StartTime: time.Now().Add(1 * time.Hour),
				EndTime:   time.Now(),
				Limit:     1000,
}
lokiClient := client.NewLokiClient(lokiUrl, tenantId, basicAuth, queryParams)
```

If your instance doesn't have basic auth you should use an empty string:
```go
basicAuth := LokiBasicAuth{}
```

## Executing a query
Once you have the client instance created you can execute the query with:
```go
ctx, cancelFn := context.WithTimeout(context.Background, 3 * time.Minute)
defer cancelFn()
results, err := lokiClient.QueryLogs(ctx)
if err != nil {
    panic(err)
}

for _, logEntry := range results {
    fmt.Println("At " + logEntry.Timestamp + " found following log: " + logEntry.Log)
}
```

## Log entry types
Loki can return various data types in responses to queries. We will try to convert the following ones to `string`:
* `int`
* `float64`

If it's neither of these types nor a `string` the client will return an error. Same will happen if `nil` is returned.

# Troubleshooting
If you find yourself in trouble these two environment variables might help you:
* `RESTY_DEBUG` set to `true` will enable debug mode for the underlaying HTTP client
* `LOKI_CLIENT_LOG_LEVEL` controls log level of Loki client (for supported log levels check [logging package](../logging.md) documentation)