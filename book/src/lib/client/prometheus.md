# Prometheus

This client is basically a wrapper over the official Prometheus client that gas three usesful functions:
* new client creation
* fetching of all firing alerts
* executing an arbitrary query with time range of `(-infinity, now)`

## New instance
```go
c, err := NewPrometheusClient("prometheus.io")
if err != nil {
    panic(err)
}
```

## Get all firing alerts
```go
alerts, err := c.GetAlerts()
if err != nil {
    panic(err)
}

for _, alert := range alerts {
    fmt.Println("Found alert: " + alert.Value + "in state: " + alert.AlertState)
}
```

## Execute a query
```go
queryResult, err := c.GetQuery(`100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[2m])) * 100)`)
if err != nil {
    panic(err)
}

if asV, ok := queryResult.(.model.Vector); ok {
    for _, v := range asV {
        fmt.Println("Metric data: " +v.Metric)
        fmt.Println("Value: " + v.Value)
    }
} else {
    panic(fmt.Sprintf("Result wasn't a model.Vector, but %T", queryResult))
}

```