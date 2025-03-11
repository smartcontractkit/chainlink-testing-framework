package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

// all metrics, but error rate are calculated over a 10s interval
var (
	Loki_MedianQuery = `quantile_over_time(0.5, {branch=~"%s", commit=~"%s", go_test_name=~"%s", test_data_type=~"responses", gen_name=~"%s"} | json| unwrap duration [10s]) by (go_test_name, gen_name) / 1e6`
	Loki_95thQuery   = `quantile_over_time(0.95, {branch=~"%s", commit=~"%s", go_test_name=~"%s", test_data_type=~"responses", gen_name=~"%s"} | json| unwrap duration [10s]) by (go_test_name, gen_name) / 1e6`
	Loki_MaxQuery    = `max(max_over_time({branch=~"%s", commit=~"%s", go_test_name=~"%s", test_data_type=~"responses", gen_name=~"%s"} | json| unwrap duration [10s]) by (go_test_name, gen_name) / 1e6)`
	Loki_ErrorRate   = `sum(max_over_time({branch=~"%s", commit=~"%s", go_test_name=~"%s", test_data_type=~"stats", gen_name=~"%s"} | json| unwrap failed [%s]) by (node_id, go_test_name, gen_name)) by (__stream_shard__)`
)

// NewLokiQueryExecutor creates a new LokiQueryExecutor instance.
// It initializes the executor with the specified generator name, queries, and Loki configuration.
// This function is useful for setting up a query executor to interact with Loki for log data retrieval.
func NewLokiQueryExecutor(generatorName string, queries map[string]string, lokiConfig *wasp.LokiConfig) *LokiQueryExecutor {
	L.Debug().
		Str("Generator", generatorName).
		Int("Queries", len(queries)).
		Msg("Creating new Loki query executor")

	return &LokiQueryExecutor{
		KindName:            string(StandardQueryExecutor_Loki),
		GeneratorNameString: generatorName,
		Queries:             queries,
		Config:              lokiConfig,
		QueryResults:        make(map[string]interface{}),
	}
}

type LokiQueryExecutor struct {
	KindName            string `json:"kind"`
	GeneratorNameString string `json:"generator_name"`

	// Test metrics
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	// a map of name to query template, ex: "average cpu usage": "avg(rate(cpu_usage_seconds_total[5m]))"
	Queries map[string]string `json:"queries"`
	// can be anything, avg RPS, amount of errors, 95th percentile of CPU utilization, etc
	QueryResults map[string]interface{} `json:"query_results"`

	Config *wasp.LokiConfig `json:"-"`
}

// GeneratorName returns the name of the generator associated with the LokiQueryExecutor.
// It is useful for identifying the source of results in reports or logs.
func (l *LokiQueryExecutor) GeneratorName() string {
	return l.GeneratorNameString
}

// Results returns the query results as a map of string to interface{}.
// It allows users to access the outcomes of executed queries, facilitating further processing or type assertions.
func (l *LokiQueryExecutor) Results() map[string]interface{} {
	return l.QueryResults
}

// Kind returns the type of the query executor as a string.
// It is used to identify the specific kind of query executor in various operations.
func (l *LokiQueryExecutor) Kind() string {
	return l.KindName
}

// IsComparable checks if the given QueryExecutor is of the same type as the current instance.
// It compares the queries of both executors to ensure they are equivalent in structure and content.
// This function is useful for validating compatibility between different query executors.
func (l *LokiQueryExecutor) IsComparable(otherQueryExecutor QueryExecutor) error {
	L.Debug().
		Str("Expected kind", l.KindName).
		Msg("Checking if query executors are comparable")

	otherType := reflect.TypeOf(otherQueryExecutor)

	if otherType != reflect.TypeOf(l) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(l), otherType)
	}

	otherAsLoki := otherQueryExecutor.(*LokiQueryExecutor)
	if l.GeneratorNameString != otherAsLoki.GeneratorNameString {
		return fmt.Errorf("generator name is different. Expected %s, got %s", l.GeneratorNameString, otherAsLoki.GeneratorNameString)
	}

	queryErr := l.compareQueries(otherQueryExecutor.(*LokiQueryExecutor).Queries)
	if queryErr != nil {
		return queryErr
	}

	L.Debug().
		Str("Kind", l.KindName).
		Msg("Query executors are comparable")

	return nil
}

// Validate checks if the LokiQueryExecutor has valid queries and configuration.
// It returns an error if no queries are set or if the configuration is missing,
// ensuring that the executor is ready for execution.
func (l *LokiQueryExecutor) Validate() error {
	L.Debug().
		Msg("Validating Loki query executor")

	if len(l.Queries) == 0 {
		return errors.New("there are no Loki queries, there's nothing to fetch. Please set them and try again")
	}
	if l.Config == nil {
		return errors.New("loki config is missing. Please set it and try again")
	}
	if l.GeneratorNameString == "" {
		return errors.New("generator name is missing. Please set it and try again")
	}

	L.Debug().
		Msg("Loki query executor is valid")

	return nil
}

// Execute runs the configured Loki queries concurrently and collects the results.
// It requires a valid configuration and handles basic authentication if provided.
// The function returns an error if any query execution fails or if the configuration is missing.
func (l *LokiQueryExecutor) Execute(ctx context.Context) error {
	L.Info().
		Str("Generator", l.GeneratorNameString).
		Int("Queries", len(l.Queries)).
		Msg("Executing Loki queries")

	var basicAuth client.LokiBasicAuth

	if l.Config == nil {
		return errors.New("loki config is missing. Please set it and try again")
	}

	if l.Config.BasicAuth != "" {
		splitAuth := strings.Split(l.Config.BasicAuth, ":")
		if len(splitAuth) == 2 {
			basicAuth = client.LokiBasicAuth{
				Login:    splitAuth[0],
				Password: splitAuth[1],
			}
		}
	}

	l.QueryResults = make(map[string]interface{})
	resultCh := make(chan map[string][]string, len(l.Queries))
	errGroup, errCtx := errgroup.WithContext(ctx)

	for name, query := range l.Queries {
		L.Debug().
			Str("Generator", l.GeneratorNameString).
			Str("Query name", name).
			Str("Query", query).
			Msg("Executing Loki query")

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
				L.Debug().
					Str("Generator", l.GeneratorNameString).
					Str("Query name", name).
					Msg("Loki query executed successfully")
				return nil
			case <-errCtx.Done():
				L.Debug().
					Str("Generator", l.GeneratorNameString).
					Str("Query name", name).
					Str("Upstream error", errCtx.Err().Error()).
					Msg("Loki query execution cancelled")
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

	L.Info().
		Str("Generator", l.GeneratorNameString).
		Int("Queries", len(l.Queries)).
		Msg("Loki queries executed successfully")

	return nil
}

func (l *LokiQueryExecutor) compareQueries(other map[string]string) error {
	this := l.Queries
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

// TimeRange sets the start and end time for the Loki query execution.
// This function is essential for defining the time window of the data to be fetched.
func (l *LokiQueryExecutor) TimeRange(start, end time.Time) {
	l.StartTime = start
	l.EndTime = end
}

// UnmarshalJSON parses the JSON-encoded data and populates the LokiQueryExecutor fields.
// It converts the query results from a generic map to a specific type map, enabling type-safe access to the results.
func (l *LokiQueryExecutor) UnmarshalJSON(data []byte) error {
	// helper struct with QueryResults map[string]interface{}
	type Alias LokiQueryExecutor
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

	*l = LokiQueryExecutor(raw.Alias)
	l.QueryResults = convertedTypes
	return nil
}

// NewStandardMetricsLokiExecutor creates a LokiQueryExecutor configured with standard metrics queries.
// It generates queries based on provided test parameters and time range, returning the executor or an error if query generation fails.
func NewStandardMetricsLokiExecutor(lokiConfig *wasp.LokiConfig, testName, generatorName, branch, commit string, startTime, endTime time.Time) (*LokiQueryExecutor, error) {
	lq := &LokiQueryExecutor{
		KindName:            string(StandardQueryExecutor_Loki),
		GeneratorNameString: generatorName,
		Config:              lokiConfig,
		QueryResults:        make(map[string]interface{}),
	}

	standardQueries, queryErr := lq.generateStandardQueries(testName, generatorName, branch, commit, startTime, endTime)
	if queryErr != nil {
		return nil, queryErr
	}

	lq.Queries = standardQueries

	return lq, nil
}

func (l *LokiQueryExecutor) standardQuery(standardMetric StandardLoadMetric, testName, generatorName, branch, commit string, startTime, endTime time.Time) (string, error) {
	// if we decide to include only plain segments for the calculation, we we will need to execute this function for each of them
	// and then aggregate the results
	switch standardMetric {
	case MedianLatency:
		return fmt.Sprintf(Loki_MedianQuery, branch, commit, testName, generatorName), nil
	case Percentile95Latency:
		return fmt.Sprintf(Loki_95thQuery, branch, commit, testName, generatorName), nil
	case MaxLatency:
		return fmt.Sprintf(Loki_MaxQuery, branch, commit, testName, generatorName), nil
	case ErrorRate:
		queryRange := calculateTimeRange(startTime, endTime)
		return fmt.Sprintf(Loki_ErrorRate, branch, commit, testName, generatorName, queryRange), nil
	default:
		return "", fmt.Errorf("unsupported standard metric %s", standardMetric)
	}
}

func (l *LokiQueryExecutor) generateStandardQueries(testName, generatorName, branch, commit string, startTime, endTime time.Time) (map[string]string, error) {
	L.Debug().
		Msg("Generating standard Loki queries")

	standardQueries := make(map[string]string)

	for _, metric := range StandardLoadMetrics {
		query, err := l.standardQuery(metric, testName, generatorName, branch, commit, startTime, endTime)
		if err != nil {
			return nil, err
		}
		standardQueries[string(metric)] = query
	}

	L.Debug().
		Int("Queries", len(standardQueries)).
		Msg("Standard Loki queries generated")

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
