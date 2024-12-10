package benchspy

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

type PrometheusResourceReporter struct {
	Kind               StandardResourceMonitorType `json:"kind"`
	startTime, endTime time.Time                   `json:"-"`
	client             *client.Prometheus          `json:"-"`
	Queries            map[string]string           `json:"queries"`
	ResourceResults    map[string]interface{}      `json:"resources"`
	warnings           map[string]v1.Warnings      `json:"-"`
}

// NewPrometheusResourceReporter creates a new PrometheusResourceReporter, url should include basic auth if needed
func NewPrometheusResourceReporter(url string, startTime, endTime time.Time, queries map[string]string) (*PrometheusResourceReporter, error) {
	c, err := client.NewPrometheusClient(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Prometheus client")
	}

	return &PrometheusResourceReporter{
		Kind:            StandardResourceMonitor_Prometheus,
		client:          c,
		Queries:         queries,
		startTime:       startTime,
		endTime:         endTime,
		ResourceResults: make(map[string]interface{}),
	}, nil
}

func NewStandardPrometheusResourceReporter(url string, startTime, endTime time.Time, nameRegexPattern string) (*PrometheusResourceReporter, error) {
	c, err := client.NewPrometheusClient(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Prometheus client")
	}

	standardQueries, queryErr := (&PrometheusResourceReporter{client: c}).generateStandardQueries(nameRegexPattern, startTime, endTime)
	if queryErr != nil {
		return nil, errors.Wrapf(queryErr, "failed to generate standard queries for %s", nameRegexPattern)
	}

	return &PrometheusResourceReporter{
		client:          c,
		Queries:         standardQueries,
		startTime:       startTime,
		endTime:         endTime,
		ResourceResults: make(map[string]interface{}),
	}, nil
}

func (r *PrometheusResourceReporter) Fetch(ctx context.Context) error {
	for name, query := range r.Queries {
		result, warnings, queryErr := r.client.Query(ctx, query, r.endTime)
		if queryErr != nil {
			return errors.Wrapf(queryErr, "failed to query Prometheus for %s", name)
		}

		if len(warnings) > 0 {
			r.warnings[name] = warnings
		}

		r.ResourceResults[name] = result
	}

	return nil
}

func (r *PrometheusResourceReporter) Resources() map[string]interface{} {
	return r.ResourceResults
}

func (r *PrometheusResourceReporter) Validate() error {
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

func (r *PrometheusResourceReporter) IsComparable(other ResourceMonitor) error {
	otherType := reflect.TypeOf(other)
	if otherType != reflect.TypeOf(r) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(r), otherType)
	}

	asPrometheusResourceReporter := other.(*PrometheusResourceReporter)

	return r.compareQueries(asPrometheusResourceReporter.Queries)
}

func (r *PrometheusResourceReporter) compareQueries(other map[string]string) error {
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

func (r *PrometheusResourceReporter) Warnings() map[string]v1.Warnings {
	return r.warnings
}

func (r *PrometheusResourceReporter) MustResourcesAsValue() map[string]model.Value {
	resources := make(map[string]model.Value)
	for name, resource := range r.ResourceResults {
		resources[name] = resource.(model.Value)
	}
	return resources
}

func (r *PrometheusResourceReporter) TimeRange(startTime, endTime time.Time) {
	// not sure if we need to set the time range for Prometheus
	// I think we should remove that method all together from all interfaces
	// and instead make sure that all segments have start/end times (i.e. they were executed), when new report is created
	return
}

func (r *PrometheusResourceReporter) standardQuery(metric StandardResourceMetric, nameRegexPattern string, startTime, endTime time.Time) (string, error) {
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

func (r *PrometheusResourceReporter) generateStandardQueries(nameRegexPattern string, startTime, endTime time.Time) (map[string]string, error) {
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
