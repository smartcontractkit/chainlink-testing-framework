package benchspy

import (
	"context"
	"time"
)

type Storer interface {
	// Store stores the report in a persistent storage and returns the path to it, or an error
	Store() (string, error)
	// Load loads the report from a persistent storage and returns it, or an error
	Load(testName, commitOrTag string) error
	// LoadLatest loads the latest report from a persistent storage and returns it, or an error
	LoadLatest(testName string) error
}

type DataFetcher interface {
	// Fetch populates the report with the data from the test
	FetchData(ctx context.Context) error
}

type Comparator interface {
	// IsComparable checks whether both reports can be compared (e.g. test config is the same, app's resources are the same, queries or metrics used are the same, etc.), and an error if any difference is found
	IsComparable(otherReport Reporter) error
}

type ResourceFetcher interface {
	// FetchResources fetches the resources used by the AUT (e.g. CPU, memory, etc.)
	FetchResources(ctx context.Context) error
}

type Reporter interface {
	Storer
	DataFetcher
	ResourceFetcher
	Comparator
}

type QueryExecutor interface {
	// Validate checks if the QueryExecutor has all the necessary data and configuration to execute the queries
	Validate() error
	// Execute executes the queries and populates the QueryExecutor with the results
	Execute(ctx context.Context) error
	// Results returns the results of the queries, where key is the name of the query and value is the result
	Results() map[string][]string
	// IsComparable checks whether both QueryExecutors can be compared (e.g. they have the same type, queries are the same, etc.), and returns an error (if any difference is found)
	IsComparable(other QueryExecutor) error
	// TimeRange sets the time range for the queries
	TimeRange(time.Time, time.Time)
}
