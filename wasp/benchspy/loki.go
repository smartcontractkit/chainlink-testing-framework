package benchspy

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"golang.org/x/sync/errgroup"
)

func NewLokiQueryExecutor(queries map[string]string, lokiConfig *wasp.LokiConfig) *LokiQueryExecutor {
	return &LokiQueryExecutor{
		Kind:         "loki",
		Queries:      queries,
		Config:       lokiConfig,
		QueryResults: make(map[string][]string),
	}
}

type LokiQueryExecutor struct {
	Kind string `json:"kind"`
	// Test metrics
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	// Performance queries
	// a map of name to query template, ex: "average cpu usage": "avg(rate(cpu_usage_seconds_total[5m]))"
	Queries map[string]string `json:"queries"`
	// Performance queries results
	// can be anything, avg RPS, amount of errors, 95th percentile of CPU utilization, etc
	QueryResults map[string][]string `json:"query_results"`

	Config *wasp.LokiConfig `json:"-"`
}

func (l *LokiQueryExecutor) Results() map[string][]string {
	return l.QueryResults
}

func (l *LokiQueryExecutor) IsComparable(otherQueryExecutor QueryExecutor) error {
	otherType := reflect.TypeOf(otherQueryExecutor)

	if otherType != reflect.TypeOf(l) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(l), otherType)
	}

	return l.compareLokiQueries(otherQueryExecutor.(*LokiQueryExecutor).Queries)
}

func (l *LokiQueryExecutor) Validate() error {
	if len(l.Queries) == 0 {
		return errors.New("there are no Loki queries, there's nothing to fetch. Please set them and try again")
	}
	if l.Config == nil {
		return errors.New("loki config is missing. Please set it and try again")
	}

	return nil
}

func (l *LokiQueryExecutor) Execute(ctx context.Context) error {
	splitAuth := strings.Split(l.Config.BasicAuth, ":")
	var basicAuth client.LokiBasicAuth
	if len(splitAuth) == 2 {
		basicAuth = client.LokiBasicAuth{
			Login:    splitAuth[0],
			Password: splitAuth[1],
		}
	}

	l.QueryResults = make(map[string][]string)
	resultCh := make(chan map[string][]string, len(l.Queries))
	errGroup, errCtx := errgroup.WithContext(ctx)

	for name, query := range l.Queries {
		errGroup.Go(func() error {
			queryParams := client.LokiQueryParams{
				Query:     query,
				StartTime: l.StartTime,
				EndTime:   l.EndTime,
				Limit:     1000, //TODO make this configurable
			}

			parsedLokiUrl, err := url.Parse(l.Config.URL)
			if err != nil {
				return errors.Wrapf(err, "failed to parse Loki URL %s", l.Config.URL)
			}

			lokiUrl := parsedLokiUrl.Scheme + "://" + parsedLokiUrl.Host
			lokiClient := client.NewLokiClient(lokiUrl, l.Config.TenantID, basicAuth, queryParams)

			rawLogs, err := lokiClient.QueryLogs(errCtx)
			if err != nil {
				return errors.Wrapf(err, "failed to query logs for %s", name)
			}

			resultMap := make(map[string][]string)
			for _, log := range rawLogs {
				resultMap[name] = append(resultMap[name], log.Log)
			}

			select {
			case resultCh <- resultMap:
				return nil
			case <-errCtx.Done():
				return errCtx.Err() // Allows goroutine to exit if timeout occurs
			}
		})
	}

	if err := errGroup.Wait(); err != nil {
		return errors.Wrap(err, "failed to execute Loki queries")
	}

	for i := 0; i < len(l.Queries); i++ {
		result := <-resultCh
		for name, logs := range result {
			l.QueryResults[name] = logs
		}
	}

	return nil
}

func (l *LokiQueryExecutor) compareLokiQueries(other map[string]string) error {
	this := l.Queries
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

func (l *LokiQueryExecutor) TimeRange(start, end time.Time) {
	l.StartTime = start
	l.EndTime = end
}

type StandardMetric string

const (
	MedianLatency       StandardMetric = "median_latency"
	Percentile95Latency StandardMetric = "95th_percentile_latency"
	ErrorRate           StandardMetric = "error_rate"
)

var standardMetrics = []StandardMetric{MedianLatency, Percentile95Latency, ErrorRate}

func NewStandardMetricsLokiExecutor(lokiConfig *wasp.LokiConfig, testName, generatorName, branch, commit string, startTime, endTime time.Time) (*LokiQueryExecutor, error) {
	standardQueries, queryErr := generateStandardLokiQueries(testName, generatorName, branch, commit, startTime, endTime)
	if queryErr != nil {
		return nil, queryErr
	}

	return &LokiQueryExecutor{
		Kind:         "loki",
		Queries:      standardQueries,
		Config:       lokiConfig,
		QueryResults: make(map[string][]string),
	}, nil
}

func standardQuery(standardMetric StandardMetric, testName, generatorName, branch, commit string, startTime, endTime time.Time) (string, error) {
	switch standardMetric {
	case MedianLatency:
		return fmt.Sprintf("quantile_over_time(0.5, {branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} | json| unwrap duration [10s]) by (go_test_name, gen_name) / 1e6", branch, commit, testName, generatorName), nil
	case Percentile95Latency:
		return fmt.Sprintf("quantile_over_time(0.95, {branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} | json| unwrap duration [10s]) by (go_test_name, gen_name) / 1e6", branch, commit, testName, generatorName), nil
	case ErrorRate:
		queryRange := calculateTimeRange(startTime, endTime)
		// this becomes problematic if we want to only consider plain segments, because each might have a different length and thus should have a different range window for accurate calculation
		// unless... we will are only interested in comparing the differences between reports, not the actual values, then it won't matter that error rate is skewed (calculated over ranges longer than query interval)
		return fmt.Sprintf("sum(max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json| unwrap failed [%s]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)", branch, commit, testName, generatorName, queryRange), nil
	default:
		return "", fmt.Errorf("unsupported standard metric %s", standardMetric)
	}
}

func generateStandardLokiQueries(testName, generatorName, branch, commit string, startTime, endTime time.Time) (map[string]string, error) {
	standardQueries := make(map[string]string)

	for _, metric := range standardMetrics {
		query, err := standardQuery(metric, testName, generatorName, branch, commit, startTime, endTime)
		if err != nil {
			return nil, err
		}
		standardQueries[string(metric)] = query
	}

	return standardQueries, nil
}

func calculateTimeRange(startTime, endTime time.Time) string {
	totalSeconds := int(endTime.Sub(startTime).Seconds())

	var rangeStr string
	if totalSeconds%3600 == 0 { // Exact hours
		rangeStr = fmt.Sprintf("%dh", totalSeconds/3600)
	} else if totalSeconds%60 == 0 { // Exact minutes
		rangeStr = fmt.Sprintf("%dm", totalSeconds/60)
	} else { // Use seconds for uneven intervals
		rangeStr = fmt.Sprintf("%ds", totalSeconds)
	}

	return rangeStr
}
