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

type Reporter interface {
	Storer
	DataFetcher
	Comparator
}

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

type NamedGenerator interface {
	// GeneratorName returns the name of the generator
	GeneratorName() string
}

type StandardQueryExecutorType string

const (
	StandardQueryExecutor_Loki       StandardQueryExecutorType = "loki"
	StandardQueryExecutor_Direct     StandardQueryExecutorType = "direct"
	StandardQueryExecutor_Prometheus StandardQueryExecutorType = "prometheus"
)

type StandardLoadMetric string

const (
	MedianLatency       StandardLoadMetric = "median_latency"
	Percentile95Latency StandardLoadMetric = "95th_percentile_latency"
	MaxLatency          StandardLoadMetric = "max_latency"
	ErrorRate           StandardLoadMetric = "error_rate"
)

var StandardLoadMetrics = []StandardLoadMetric{MedianLatency, Percentile95Latency, MaxLatency, ErrorRate}

type StandardResourceMetric string

const (
	MedianCPUUsage StandardResourceMetric = "median_cpu_usage"
	MedianMemUsage StandardResourceMetric = "median_mem_usage"
	P95CPUUsage    StandardResourceMetric = "p95_cpu_usage"
	MaxCPUUsage    StandardResourceMetric = "max_cpu_usage"
	P95MemUsage    StandardResourceMetric = "p95_mem_usage"
	MaxMemUsage    StandardResourceMetric = "max_mem_usage"
)

var StandardResourceMetrics = []StandardResourceMetric{MedianCPUUsage, MedianMemUsage, P95CPUUsage, P95MemUsage, MaxCPUUsage, MaxMemUsage}
