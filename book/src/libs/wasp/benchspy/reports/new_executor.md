# BenchSpy - Adding new QueryExecutor

As mentioned previously the `StandardReport` comes with support of three different data types:
* `Direct`
* `Loki`
* `Prometheus`

Each of them implements the `QueryExecutor` interface:
```go
type QueryExecutor interface {
	// Kind returns the type of the QueryExecutor
	Kind() string
	// Validate checks if the QueryExecutor has all the necessary data and configuration to execute the queries
	Validate() error
	// Execute executes the queries and populates the QueryExecutor with the results
	Execute(ctx context.Context) error
	// Results returns the results of the queries, where key is the name of the query and value is the result
	Results() map[string]interface{}
	// IsComparable checks whether both QueryExecutors can be compared (e.g. they have the same type, queries are the same, etc.), and returns an error (if any difference is found)
	IsComparable(other QueryExecutor) error
	// TimeRange sets the time range for the queries
	TimeRange(startTime, endTime time.Time)
}
```

Most of the functions that your new `QueryExecutor` should implement are self-explanatory and I will skip them and focus on two that might not be obvious.

## Kind
Kind should return a name as string of your `QueryExecutor`. It needs to be unique, because `StandardReport` uses it, when unmarshalling JSON files with
stored reports. That also means that you need to add support for your new executor in the `StandardReport.UnmarshallJSON` function.

> [!NOTE]
> If your new `QueryExecutor` uses interfaces or `interface{}` or `any` types, or has some fields that should not/cannot be serialized,
> remember to add custom `MarshallJSON` and `UnmarshallJSON` functions to it. Existing executors can serve as a good example.

## TimeRange
`StandardReport` calls this method just before calling `Execute()` for each executor. It is used to set the time range for the query (required
for Loki and Prometheus). This is done primarly to avoid the need for manual calculation of the time range for the test, because `StandardReport` does it automatically by
analysing schedules of all generators and finding earliest start time and latest end time.