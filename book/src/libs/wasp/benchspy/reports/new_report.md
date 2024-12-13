# Defining a new report

Each `BenchSpy` report should implement the `Reporter` interface, which handles 3 responsibilities:
* storage and retrival (`Storer` interface)
* data fetching (`DataFetcher` interface)
* comparator (`Comparator` interface)

### Definition
```go
type Reporter interface {
	Storer
	DataFetcher
	Comparator
}

```

Comparison of actual performance data should not be part of the report and should be done independently from it,
ideally using simple Go's `require` and `assert` statements.

# Storer interface
## Definition
```go
type Storer interface {
	// Store stores the report in a persistent storage and returns the path to it, or an error
	Store() (string, error)
	// Load loads the report from a persistent storage and returns it, or an error
	Load(testName, commitOrTag string) error
	// LoadLatest loads the latest report from a persistent storage and returns it, or an error
	LoadLatest(testName string) error
}
```

If storing the reports on the local filesystem under Git fulfills your requirements you can reuse the `LocalStorage`
implementation of `Storer`.

If you would like to store them in S3 or a database, you will need to implement the interface yourself.

# DataFetcher interface
## Definition
```go
type DataFetcher interface {
	// Fetch populates the report with the data from the test
	FetchData(ctx context.Context) error
}
```
This interface is only concerned with fetching the data from the data source and populating the results.

# Comparator interface
## Definition
```go
type Comparator interface {
	// IsComparable checks whether both reports can be compared (e.g. test config is the same, app's resources are the same, queries or metrics used are the same, etc.), and an error if any difference is found
	IsComparable(otherReport Reporter) error
}

```

This interface is only concerned with making sure that both report are comparable, for example by checking:
* whether both use generators with identical configurations (such as load type, load characteristics)
* whether both report feature the same data sources and queries