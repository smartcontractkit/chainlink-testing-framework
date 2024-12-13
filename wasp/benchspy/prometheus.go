package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/pkg/errors"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

type PrometheusConfig struct {
	Url               string
	NameRegexPatterns []string
}

const PrometheusUrlEnvVar = "PROMETHEUS_URL"

// NewPrometheusConfig creates a new PrometheusConfig instance with the specified name regex patterns.
// It retrieves the Prometheus URL from the environment and is used to configure query execution for Prometheus data sources.
func NewPrometheusConfig(nameRegexPatterns ...string) *PrometheusConfig {
	return &PrometheusConfig{
		Url:               os.Getenv(PrometheusUrlEnvVar),
		NameRegexPatterns: nameRegexPatterns,
	}
}

var WithoutPrometheus *PrometheusConfig = nil

type PrometheusQueryExecutor struct {
	KindName     string                 `json:"kind"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	client       v1.API                 `json:"-"`
	Queries      map[string]string      `json:"queries"`
	QueryResults map[string]interface{} `json:"query_results"`
	warnings     map[string]v1.Warnings `json:"-"`
}

// NewPrometheusQueryExecutor creates a new PrometheusResourceReporter, url should include basic auth if needed
func NewPrometheusQueryExecutor(queries map[string]string, config *PrometheusConfig) (*PrometheusQueryExecutor, error) {
	c, err := client.NewPrometheusClient(config.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Prometheus client")
	}

	return &PrometheusQueryExecutor{
		KindName:     string(StandardQueryExecutor_Prometheus),
		client:       c,
		Queries:      queries,
		QueryResults: make(map[string]interface{}),
	}, nil
}

// NewStandardPrometheusQueryExecutor creates a PrometheusQueryExecutor with standard queries 
// based on the provided time range and configuration. It simplifies the process of generating 
// queries for Prometheus, making it easier to integrate Prometheus data into reports.
func NewStandardPrometheusQueryExecutor(startTime, endTime time.Time, config *PrometheusConfig) (*PrometheusQueryExecutor, error) {
	p := &PrometheusQueryExecutor{}

	standardQueries := make(map[string]string)
	for _, nameRegexPattern := range config.NameRegexPatterns {
		queries, queryErr := p.generateStandardQueries(nameRegexPattern, startTime, endTime)
		if queryErr != nil {
			return nil, errors.Wrapf(queryErr, "failed to generate standard queries for %s", nameRegexPattern)
		}

		for name, query := range queries {
			standardQueries[name] = query
		}
	}

	return NewPrometheusQueryExecutor(standardQueries, config)
}

// Execute runs the defined Prometheus queries concurrently, collecting results and warnings.
// It returns an error if any query fails, allowing for efficient data retrieval in reporting tasks.
func (r *PrometheusQueryExecutor) Execute(ctx context.Context) error {
	for name, query := range r.Queries {
		result, warnings, queryErr := r.client.Query(ctx, query, r.EndTime)
		if queryErr != nil {
			return errors.Wrapf(queryErr, "failed to query Prometheus for %s", name)
		}

		if len(warnings) > 0 {
			r.warnings[name] = warnings
		}

		r.QueryResults[name] = result
	}

	return nil
}

// Results returns the query results as a map of string to interface{}.
// It allows users to access the results of executed queries, facilitating data retrieval and manipulation.
func (r *PrometheusQueryExecutor) Results() map[string]interface{} {
	return r.QueryResults
}

// Kind returns the type of the query executor as a string.
// It is used to identify the specific kind of executor in a collection of query executors.
func (l *PrometheusQueryExecutor) Kind() string {
	return l.KindName
}

// Validate checks the PrometheusQueryExecutor for a valid client and ensures that at least one query is provided. 
// It returns an error if the client is nil or no queries are specified, helping to ensure proper configuration before execution.
func (r *PrometheusQueryExecutor) Validate() error {
	if r.client == nil {
		return errors.New("prometheus client is nil")
	}

	if len(r.Queries) == 0 {
		return errors.New("no queries provided")
	}

	return nil
}

// IsComparable checks if the provided QueryExecutor is of the same type as the receiver.
// It returns an error if the types do not match, ensuring type safety for query comparisons.
func (r *PrometheusQueryExecutor) IsComparable(other QueryExecutor) error {
	otherType := reflect.TypeOf(other)
	if otherType != reflect.TypeOf(r) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(r), otherType)
	}

	asPrometheusResourceReporter := other.(*PrometheusQueryExecutor)

	return r.compareQueries(asPrometheusResourceReporter.Queries)
}

func (r *PrometheusQueryExecutor) compareQueries(other map[string]string) error {
	this := r.Queries
	if len(this) != len(other) {
		return fmt.Errorf("queries count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, query1 := range this {
		if query2, ok := other[name1]; !ok {
			return fmt.Errorf("query %s is missing from the other report", name1)
		} else {
			if query1 != query2 {
				return fmt.Errorf("query %s is different. Expected %s, got %s", name1, query1, query2)
			}
		}
	}

	return nil
}

// Warnings returns a map of warnings encountered during query execution.
// This function is useful for retrieving any issues that may have arisen, 
// allowing users to handle or log them appropriately.
func (r *PrometheusQueryExecutor) Warnings() map[string]v1.Warnings {
	return r.warnings
}

// MustResultsAsValue retrieves the query results as a map of metric names to their corresponding values.
// It ensures that the results are in a consistent format, making it easier to work with metrics in subsequent operations.
func (r *PrometheusQueryExecutor) MustResultsAsValue() map[string]model.Value {
	results := make(map[string]model.Value)
	for name, result := range r.QueryResults {
		var val model.Value
		switch v := result.(type) {
		case model.Matrix:
			// model.Matrix implements model.Value with value receivers
			val = v
		case *model.Matrix:
			val = v
		case model.Vector:
			// model.Vector implements model.Value with value receivers
			val = v
		case *model.Vector:
			val = v
		case model.Scalar:
			scalar := v
			// *model.Scalar implements model.Value
			val = &scalar
		case *model.Scalar:
			val = v
		case model.String:
			str := v
			// *model.String implements model.Value
			val = &str
		case *model.String:
			val = v
		default:
			panic(fmt.Sprintf("Unknown result type: %T", result))
		}
		results[name] = val
	}
	return results
}

// TimeRange sets the start and end time for the Prometheus query execution.
// This function is essential for defining the time window for data retrieval, ensuring accurate and relevant results.
func (r *PrometheusQueryExecutor) TimeRange(startTime, endTime time.Time) {
	r.StartTime = startTime
	r.EndTime = endTime
}

func (r *PrometheusQueryExecutor) standardQuery(metric StandardResourceMetric, nameRegexPattern string, startTime, endTime time.Time) (string, error) {
	duration := calculateTimeRange(startTime, endTime)
	switch metric {
	case MedianCPUUsage:
		return fmt.Sprintf("quantile_over_time(0.5, rate(container_cpu_usage_seconds_total{name=~\"%s\"}[5m])[%s:10s]) * 100", nameRegexPattern, duration), nil
	case P95CPUUsage:
		return fmt.Sprintf("quantile_over_time(0.95, rate(container_cpu_usage_seconds_total{name=~\"%s\"}[5m])[%s:10s]) * 100", nameRegexPattern, duration), nil
	case MedianMemUsage:
		return fmt.Sprintf("quantile_over_time(0.5, rate(container_memory_usage_bytes{name=~\"%s\"}[5m])[%s:10s]) * 100", nameRegexPattern, duration), nil
	case P95MemUsage:
		// this becomes problematic if we want to only consider plain segments, because each might have a different length and thus should have a different range window for accurate calculation
		// unless... we will are only interested in comparing the differences between reports, not the actual values, then it won't matter that error rate is skewed (calculated over ranges longer than query interval)
		return fmt.Sprintf("quantile_over_time(0.95, rate(container_memory_usage_bytes{name=~\"%s\"}[5m])[%s:10s]) * 100", nameRegexPattern, duration), nil
	default:
		return "", fmt.Errorf("unsupported standard metric %s", metric)
	}
}

func (r *PrometheusQueryExecutor) generateStandardQueries(nameRegexPattern string, startTime, endTime time.Time) (map[string]string, error) {
	standardQueries := make(map[string]string)

	for _, metric := range standardResourceMetrics {
		query, err := r.standardQuery(metric, nameRegexPattern, startTime, endTime)
		if err != nil {
			return nil, err
		}
		standardQueries[string(metric)] = query
	}

	return standardQueries, nil
}

type TypedMetric struct {
	Value      model.Value `json:"value"`
	MetricType string      `json:"metric_type"`
}

// MarshalJSON customizes the JSON representation of PrometheusQueryExecutor.
// It includes only essential fields: Kind, Queries, and simplified QueryResults.
// This function is useful for serializing the executor's state in a concise format.
func (g *PrometheusQueryExecutor) MarshalJSON() ([]byte, error) {
	// we need custom marshalling to only include some parts of the metrics
	type QueryExecutor struct {
		Kind         string                 `json:"kind"`
		Queries      map[string]string      `json:"queries"`
		QueryResults map[string]TypedMetric `json:"query_results"`
	}

	q := &QueryExecutor{
		Kind:    g.KindName,
		Queries: g.Queries,
		QueryResults: func() map[string]TypedMetric {
			simplifiedMetrics := make(map[string]TypedMetric)
			for name, value := range g.MustResultsAsValue() {
				simplifiedMetrics[name] = TypedMetric{
					MetricType: value.Type().String(),
					Value:      value,
				}
			}
			return simplifiedMetrics
		}(),
	}

	return json.Marshal(q)
}

// UnmarshalJSON decodes JSON data into a PrometheusQueryExecutor instance.
// It populates the QueryResults field with appropriately typed metrics,
// enabling easy access to the results of Prometheus queries.
func (r *PrometheusQueryExecutor) UnmarshalJSON(data []byte) error {
	// helper struct with QueryResults map[string]interface{}
	type Alias PrometheusQueryExecutor
	var raw struct {
		Alias
		QueryResults map[string]json.RawMessage `json:"query_results"`
	}

	// unmarshal into the helper struct to populate other fields automatically
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var convertedQueryResults = make(map[string]interface{})

	for name, rawResult := range raw.QueryResults {
		var rawTypedMetric struct {
			MetricType string          `json:"metric_type"`
			Value      json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(rawResult, &rawTypedMetric); err != nil {
			return errors.Wrapf(err, "failed to unmarshal query result for %s", name)
		}

		switch rawTypedMetric.MetricType {
		case "scalar":
			var scalar model.Scalar
			if err := json.Unmarshal(rawTypedMetric.Value, &scalar); err != nil {
				return errors.Wrapf(err, "failed to unmarshal scalar value for %s", name)
			}
			convertedQueryResults[name] = &scalar
		case "vector":
			var vector model.Vector
			if err := json.Unmarshal(rawTypedMetric.Value, &vector); err != nil {
				return errors.Wrapf(err, "failed to unmarshal vector value for %s", name)
			}
			convertedQueryResults[name] = vector
		case "matrix":
			var matrix model.Matrix
			if err := json.Unmarshal(rawTypedMetric.Value, &matrix); err != nil {
				return errors.Wrapf(err, "failed to unmarshal matrix value for %s", name)
			}
			convertedQueryResults[name] = matrix
		case "string":
			var str model.String
			if err := json.Unmarshal(rawTypedMetric.Value, &str); err != nil {
				return errors.Wrapf(err, "failed to unmarshal string value for %s", name)
			}
			convertedQueryResults[name] = &str
		default:
			return fmt.Errorf("unknown metric type %s", rawTypedMetric.MetricType)
		}
	}

	*r = PrometheusQueryExecutor(raw.Alias)
	r.QueryResults = convertedQueryResults
	return nil
}
