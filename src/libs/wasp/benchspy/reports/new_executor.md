# BenchSpy - Adding a New QueryExecutor

# QueryExecutor interface

As mentioned earlier, the `StandardReport` supports three different data sources:
- `Direct`
- `Loki`
- `Prometheus`

Each of these implements the `QueryExecutor` interface:

```go
type QueryExecutor interface {
	// Kind returns the type of the QueryExecutor
	Kind() string
	// Validate checks if the QueryExecutor has all the necessary data and configuration to execute the queries
	Validate() error
	// Execute executes the queries and populates the QueryExecutor with the results
	Execute(ctx context.Context) error
	// Results returns the results of the queries, where the key is the query name and the value is the result
	Results() map[string]interface{}
	// IsComparable checks whether both QueryExecutors can be compared (e.g., they have the same type, queries, etc.),
	// and returns an error if any differences are found
	IsComparable(other QueryExecutor) error
	// TimeRange sets the time range for the queries
	TimeRange(startTime, endTime time.Time)
}
```

When creating a new `QueryExecutor`, most functions in this interface are straightforward. Below, we focus on two that may require additional explanation.

---

## `Kind`

The `Kind` function should return a unique string identifier for your `QueryExecutor`. This identifier is crucial because `StandardReport` uses it when unmarshalling JSON files with stored reports.

Additionally, you need to extend the `StandardReport.UnmarshalJSON` function to support your new executor.

> [!NOTE]
> If your `QueryExecutor` includes interfaces, `interface{}` or `any` types, or fields that should not or cannot be serialized, ensure you implement custom `MarshalJSON` and `UnmarshalJSON` functions. Existing executors can provide useful examples.

---

## `TimeRange`

The `TimeRange` method is called by `StandardReport` just before invoking `Execute()` for each executor. This method sets the time range for queries (required for Loki and Prometheus).

By default, `StandardReport` calculates the test time range automatically by analyzing the schedules of all generators and determining the earliest start time and latest end time. This eliminates the need for manual calculations.

---

With these details in mind, you should have a clear path to implementing your own `QueryExecutor` and integrating it seamlessly with `BenchSpy`'s `StandardReport`.

# NamedGenerator interface

Executors that query load generation metrics should also implement this simple interface:
```go
type NamedGenerator interface {
	// GeneratorName returns the name of the generator
	GeneratorName() string
}
```

It is used primarly, when casting results from `map[string]interface{}` to target type, while splitting them between different generators.

Currently, this interface is implemented by `Direct` and `Loki` exectors, but not by `Prometheus`.