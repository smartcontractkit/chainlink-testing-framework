# Defining a New Report

Each `BenchSpy` report must implement the `Reporter` interface, which handles three primary responsibilities:
- **Storage and retrieval** (`Storer` interface)
- **Data fetching** (`DataFetcher` interface)
- **Comparison** (`Comparator` interface)

---

## Reporter Interface
### Definition
```go
type Reporter interface {
	Storer
	DataFetcher
	Comparator
}
```

The comparison of actual performance data should not be part of the report itself. It should be done independently, ideally using Go's `require` and `assert` statements.

---

## Storer Interface
### Definition
```go
type Storer interface {
	// Store stores the report in persistent storage and returns the path to it, or an error
	Store() (string, error)
	// Load loads the report from persistent storage based on the test name and commit/tag, or returns an error
	Load(testName, commitOrTag string) error
	// LoadLatest loads the latest report from persistent storage for the given test name, or returns an error
	LoadLatest(testName string) error
}
```

### Usage
- If storing reports locally under Git satisfies your requirements, you can reuse the `LocalStorage` implementation of `Storer`.
- If you need to store reports in S3 or a database, you will need to implement the interface yourself.

---

## DataFetcher Interface
### Definition
```go
type DataFetcher interface {
	// FetchData populates the report with data from the test
	FetchData(ctx context.Context) error
}
```

### Purpose
This interface is solely responsible for fetching data from the data source and populating the report with results.

---

## Comparator Interface
### Definition
```go
type Comparator interface {
	// IsComparable checks whether two reports can be compared (e.g., test configuration, app resources, queries, and metrics are identical),
	// and returns an error if any differences are found
	IsComparable(otherReport Reporter) error
}
```

### Purpose
This interface ensures that both reports are comparable by verifying:
- Both use generators with identical configurations (e.g., load type, load characteristics).
- Both reports feature the same data sources and queries.

---

This design provides flexibility and composability, allowing you to store, fetch, and compare reports in a way that fits your specific requirements.