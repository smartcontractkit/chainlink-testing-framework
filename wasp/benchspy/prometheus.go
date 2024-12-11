package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

type PrometheusQueryExecutor struct {
	KindName           string                 `json:"kind"`
	startTime, endTime time.Time              `json:"-"`
	client             v1.API                 `json:"-"`
	Queries            map[string]string      `json:"queries"`
	QueryResults       map[string]interface{} `json:"query_results"`
	warnings           map[string]v1.Warnings `json:"-"`
}

// NewPrometheusQueryExecutor creates a new PrometheusResourceReporter, url should include basic auth if needed
func NewPrometheusQueryExecutor(url string, startTime, endTime time.Time, queries map[string]string) (*PrometheusQueryExecutor, error) {
	c, err := client.NewPrometheusClient(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Prometheus client")
	}

	return &PrometheusQueryExecutor{
		KindName:     string(StandardQueryExecutor_Prometheus),
		client:       c,
		Queries:      queries,
		startTime:    startTime,
		endTime:      endTime,
		QueryResults: make(map[string]interface{}),
	}, nil
}

func NewStandardPrometheusQueryExecutor(url string, startTime, endTime time.Time, nameRegexPattern string) (*PrometheusQueryExecutor, error) {
	c, err := client.NewPrometheusClient(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Prometheus client")
	}

	pr := &PrometheusQueryExecutor{
		client:       c,
		startTime:    startTime,
		endTime:      endTime,
		QueryResults: make(map[string]interface{}),
	}

	standardQueries, queryErr := pr.generateStandardQueries(nameRegexPattern, startTime, endTime)
	if queryErr != nil {
		return nil, errors.Wrapf(queryErr, "failed to generate standard queries for %s", nameRegexPattern)
	}

	pr.Queries = standardQueries

	return pr, nil
}

func (r *PrometheusQueryExecutor) Execute(ctx context.Context) error {
	for name, query := range r.Queries {
		result, warnings, queryErr := r.client.Query(ctx, query, r.endTime)
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

func (r *PrometheusQueryExecutor) Results() map[string]interface{} {
	return r.QueryResults
}

func (l *PrometheusQueryExecutor) Kind() string {
	return l.KindName
}

func (r *PrometheusQueryExecutor) Validate() error {
	if r.client == nil {
		return errors.New("prometheus client is nil")
	}

	if len(r.Queries) == 0 {
		return errors.New("no queries provided")
	}

	if r.startTime.IsZero() {
		return errors.New("start time is not set")
	}

	if r.endTime.IsZero() {
		return errors.New("end time is not set")
	}

	return nil
}

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

func (r *PrometheusQueryExecutor) Warnings() map[string]v1.Warnings {
	return r.warnings
}

func (r *PrometheusQueryExecutor) MustResultsAsValue() map[string]model.Value {
	results := make(map[string]model.Value)
	for name, result := range r.QueryResults {
		results[name] = result.(model.Value)
	}
	return results
}

func (r *PrometheusQueryExecutor) TimeRange(startTime, endTime time.Time) {
	// not sure if we need to set the time range for Prometheus
	// I think we should remove that method all together from all interfaces
	// and instead make sure that all segments have start/end times (i.e. they were executed), when new report is created
	return
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

func (r *PrometheusQueryExecutor) UnmarshalJSON(data []byte) error {
	// helper struct with QueryResults map[string]interface{}
	type Alias PrometheusQueryExecutor
	var raw struct {
		Alias
		QueryResults map[string]interface{} `json:"query_results"`
	}

	// unmarshal into the helper struct to populate other fields automatically
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// convert map[string]interface{} to map[string]actualType
	convertedTypes, conversionErr := convertQueryResults(raw.QueryResults)
	if conversionErr != nil {
		return conversionErr
	}

	*r = PrometheusQueryExecutor(raw.Alias)
	r.QueryResults = convertedTypes
	return nil
}
