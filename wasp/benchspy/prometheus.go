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

// all of them are calculated over 5 minutes intervals (rate query), that are later sampled every 10 seconds over %s duration (quantile_over_time query)
var (
	Prometheus_MedianCPU = `quantile_over_time(0.5, rate(container_cpu_usage_seconds_total{name=~"%s"}[5m])[%s:10s]) * 100`
	Prometheus_P95CPU    = `quantile_over_time(0.95, rate(container_cpu_usage_seconds_total{name=~"%s"}[5m])[%s:10s]) * 100`
	Prometheus_MaxCPU    = `max(max_over_time(rate(container_cpu_usage_seconds_total{name=~"%s"}[5m])[%s:10s]) * 100)`
	Prometheus_MedianMem = `quantile_over_time(0.5, rate(container_memory_usage_bytes{name=~"%s"}[5m])[%s:10s]) * 100`
	Prometheus_P95Mem    = `quantile_over_time(0.95, rate(container_memory_usage_bytes{name=~"%s"}[5m])[%s:10s]) * 100`
	Prometheus_MaxMem    = `max(max_over_time(rate(container_memory_usage_bytes{name=~"%s"}[5m])[%s:10s]) * 100)`
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

var WithoutPrometheus *PrometheusConfig

type PrometheusQueryExecutor struct {
	KindName     string                 `json:"kind"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Queries      map[string]string      `json:"queries"`
	QueryResults map[string]interface{} `json:"query_results"`
	client       v1.API
	warnings     map[string]v1.Warnings
}

// NewPrometheusQueryExecutor creates a new PrometheusResourceReporter, url should include basic auth if needed
func NewPrometheusQueryExecutor(queries map[string]string, config *PrometheusConfig) (*PrometheusQueryExecutor, error) {
	L.Debug().
		Int("Queries", len(queries)).
		Msg("Creating new Prometheus query executor")

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
func (pe *PrometheusQueryExecutor) Execute(ctx context.Context) error {
	L.Info().
		Int("Queries", len(pe.Queries)).
		Msg("Executing Prometheus queries")

	for name, query := range pe.Queries {
		L.Debug().
			Str("Query name", name).
			Str("Query", query).
			Msg("Executing Prometheus query")

		result, warnings, queryErr := pe.client.Query(ctx, query, pe.EndTime)
		if queryErr != nil {
			return errors.Wrapf(queryErr, "failed to query Prometheus for %s", name)
		}

		if len(warnings) > 0 {
			pe.warnings[name] = warnings
		}

		pe.QueryResults[name] = result
	}

	L.Info().
		Int("Queries", len(pe.Queries)).
		Msg("Prometheus queries executed successfully")

	return nil
}

// Results returns the query results as a map of string to interface{}.
// It allows users to access the results of executed queries, facilitating data retrieval and manipulation.
func (pe *PrometheusQueryExecutor) Results() map[string]interface{} {
	return pe.QueryResults
}

// Kind returns the type of the query executor as a string.
// It is used to identify the specific kind of executor in a collection of query executors.
func (pe *PrometheusQueryExecutor) Kind() string {
	return pe.KindName
}

// Validate checks the PrometheusQueryExecutor for a valid client and ensures that at least one query is provided.
// It returns an error if the client is nil or no queries are specified, helping to ensure proper configuration before execution.
func (pe *PrometheusQueryExecutor) Validate() error {
	L.Debug().
		Msg("Validating Prometheus query executor")

	if pe.client == nil {
		return errors.New("prometheus client is nil")
	}

	if len(pe.Queries) == 0 {
		return errors.New("no queries provided")
	}

	L.Debug().
		Msg("Prometheus query executor is valid")

	return nil
}

// IsComparable checks if the provided QueryExecutor is of the same type as the receiver.
// It returns an error if the types do not match, ensuring type safety for query comparisons.
func (pe *PrometheusQueryExecutor) IsComparable(other QueryExecutor) error {
	L.Debug().
		Str("Expected kind", pe.KindName).
		Msg("Checking if query executors are comparable")

	otherType := reflect.TypeOf(other)
	if otherType != reflect.TypeOf(pe) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(pe), otherType)
	}

	asPrometheusResourceReporter := other.(*PrometheusQueryExecutor)

	queryErr := pe.compareQueries(asPrometheusResourceReporter.Queries)
	if queryErr != nil {
		return queryErr
	}

	L.Debug().
		Str("Kind", pe.KindName).
		Msg("Query executors are comparable")

	return nil
}

func (pe *PrometheusQueryExecutor) compareQueries(other map[string]string) error {
	this := pe.Queries
	if len(this) != len(other) {
		return fmt.Errorf("queries count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, query1 := range this {
		query2, ok := other[name1]
		if !ok {
			return fmt.Errorf("query %s is missing from the other report", name1)
		}
		if query1 != query2 {
			return fmt.Errorf("query %s is different. Expected %s, got %s", name1, query1, query2)
		}
	}

	return nil
}

// Warnings returns a map of warnings encountered during query execution.
// This function is useful for retrieving any issues that may have arisen,
// allowing users to handle or log them appropriately.
func (pe *PrometheusQueryExecutor) Warnings() map[string]v1.Warnings {
	return pe.warnings
}

// MustResultsAsValue retrieves the query results as a map of metric names to their corresponding values.
// It ensures that the results are in a consistent format, making it easier to work with metrics in subsequent operations.
func (pe *PrometheusQueryExecutor) MustResultsAsValue() map[string]model.Value {
	L.Debug().
		Msg("Casting query results to expected types")

	results := make(map[string]model.Value)
	for name, result := range pe.QueryResults {
		L.Debug().
			Str("Query name", name).
			Msg("Casting query result to expected type")
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

		L.Debug().
			Str("Query name", name).
			Str("Type", val.Type().String()).
			Msg("Query result casted to expected type")
	}
	return results
}

// TimeRange sets the start and end time for the Prometheus query execution.
// This function is essential for defining the time window for data retrieval, ensuring accurate and relevant results.
func (pe *PrometheusQueryExecutor) TimeRange(startTime, endTime time.Time) {
	pe.StartTime = startTime
	pe.EndTime = endTime
}

func (pe *PrometheusQueryExecutor) standardQuery(metric StandardResourceMetric, nameRegexPattern string, startTime, endTime time.Time) (string, error) {
	duration := calculateTimeRange(startTime, endTime)
	switch metric {
	case MedianCPUUsage:
		return fmt.Sprintf(Prometheus_MedianCPU, nameRegexPattern, duration), nil
	case P95CPUUsage:
		return fmt.Sprintf(Prometheus_P95CPU, nameRegexPattern, duration), nil
	case MaxCPUUsage:
		return fmt.Sprintf(Prometheus_MaxCPU, nameRegexPattern, duration), nil
	case MedianMemUsage:
		return fmt.Sprintf(Prometheus_MedianMem, nameRegexPattern, duration), nil
	case P95MemUsage:
		return fmt.Sprintf(Prometheus_P95Mem, nameRegexPattern, duration), nil
	case MaxMemUsage:
		return fmt.Sprintf(Prometheus_MaxMem, nameRegexPattern, duration), nil
	default:
		return "", fmt.Errorf("unsupported standard metric %s", metric)
	}
}

func (pe *PrometheusQueryExecutor) generateStandardQueries(nameRegexPattern string, startTime, endTime time.Time) (map[string]string, error) {
	L.Debug().
		Msg("Generating standard Prometheus queries")

	standardQueries := make(map[string]string)

	for _, metric := range StandardResourceMetrics {
		query, err := pe.standardQuery(metric, nameRegexPattern, startTime, endTime)
		if err != nil {
			return nil, err
		}
		standardQueries[string(metric)] = query
	}

	L.Debug().
		Int("Queries", len(standardQueries)).
		Msg("Standard Prometheus queries generated")

	return standardQueries, nil
}

type TypedMetric struct {
	Value      model.Value `json:"value"`
	MetricType string      `json:"metric_type"`
}

// MarshalJSON customizes the JSON representation of PrometheusQueryExecutor.
// It includes only essential fields: Kind, Queries, and simplified QueryResults.
// This function is useful for serializing the executor's state in a concise format.
func (pe *PrometheusQueryExecutor) MarshalJSON() ([]byte, error) {
	// we need custom marshalling to only include some parts of the metrics
	type QueryExecutor struct {
		Kind         string                 `json:"kind"`
		Queries      map[string]string      `json:"queries"`
		QueryResults map[string]TypedMetric `json:"query_results"`
	}

	q := &QueryExecutor{
		Kind:    pe.KindName,
		Queries: pe.Queries,
		QueryResults: func() map[string]TypedMetric {
			simplifiedMetrics := make(map[string]TypedMetric)
			for name, value := range pe.MustResultsAsValue() {
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
func (pe *PrometheusQueryExecutor) UnmarshalJSON(data []byte) error {
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

	*pe = PrometheusQueryExecutor(raw.Alias)
	pe.QueryResults = convertedQueryResults
	return nil
}
