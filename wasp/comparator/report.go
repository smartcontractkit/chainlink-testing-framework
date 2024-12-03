package comparator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type Report interface {
	// Store stores the report in a persistent storage and returns the path to it, or an error
	Store() (string, error)
	// Load loads the report from a persistent storage and returns it, or an error
	Load() error
	// Fetch populates the report with the data from the test
	Fetch() error
	// IsComparable checks whether both reports can be compared (e.g. test config is the same, app's resources are the same, queries or metrics used are the same, etc.), and returns a map of the differences and an error (if any difference is found)
	IsComparable(otherReport Report) (bool, map[string]string, error)
}

var directory = "performance_reports"

type BasicData struct {
	TestName    string `json:"test_name"`
	CommitOrTag string `json:"commit_or_tag"`

	// Test metrics
	TestStart time.Time `json:"test_start_timestamp"`
	TestEnd   time.Time `json:"test_end_timestamp"`

	// all, generator settings, including segments
	GeneratorConfigs map[string]*wasp.Config `json:"generator_configs"`
}

type BasicReport struct {
	BasicData
	LocalReportStorage
	ResourceReporter

	// Performance queries
	// a map of name to query template, ex: "average cpu usage": "avg(rate(cpu_usage_seconds_total[5m]))"
	LokiQueries map[string]string `json:"loki_queries"`
	// Performance queries results
	// can be anything, avg RPS, amount of errors, 95th percentile of CPU utilization, etc
	Results map[string][]string `json:"results"`
	// In case something went wrong
	Errors []error `json:"errors"`

	LokiConfig *wasp.LokiConfig `json:"-"`
}

func (b *BasicReport) Store() (string, error) {
	return b.LocalReportStorage.Store(b.TestName, b.CommitOrTag, b)
}

func (b *BasicReport) Load() error {
	return b.LocalReportStorage.Load(b.TestName, b.CommitOrTag, b)
}

func (b *BasicReport) Fetch() error {
	if len(b.LokiQueries) == 0 {
		return errors.New("there are no Loki queries, there's nothing to fetch. Please set them and try again")
	}
	if b.LokiConfig == nil {
		return errors.New("loki config is missing. Please set it and try again")
	}
	if b.TestStart.IsZero() {
		return errors.New("test start time is missing. We cannot query Loki without a time range. Please set it and try again")
	}
	if b.TestEnd.IsZero() {
		return errors.New("test end time is missing. We cannot query Loki without a time range. Please set it and try again")
	}

	splitAuth := strings.Split(b.LokiConfig.BasicAuth, ":")
	var basicAuth client.LokiBasicAuth
	if len(splitAuth) == 2 {
		basicAuth = client.LokiBasicAuth{
			Login:    splitAuth[0],
			Password: splitAuth[1],
		}
	}

	b.Results = make(map[string][]string)

	for name, query := range b.LokiQueries {
		queryParams := client.LokiQueryParams{
			Query:     query,
			StartTime: b.TestStart,
			EndTime:   b.TestEnd,
			Limit:     1000, //TODO make this configurable
		}

		parsedLokiUrl, err := url.Parse(b.LokiConfig.URL)
		if err != nil {
			return errors.Wrapf(err, "failed to parse Loki URL %s", b.LokiConfig.URL)
		}

		lokiUrl := parsedLokiUrl.Scheme + "://" + parsedLokiUrl.Host
		lokiClient := client.NewLokiClient(lokiUrl, b.LokiConfig.TenantID, basicAuth, queryParams)

		ctx, cancelFn := context.WithTimeout(context.Background(), b.LokiConfig.Timeout)
		rawLogs, err := lokiClient.QueryLogs(ctx)
		if err != nil {
			b.Errors = append(b.Errors, err)
			cancelFn()
			continue
		}

		cancelFn()
		b.Results[name] = []string{}
		for _, log := range rawLogs {
			b.Results[name] = append(b.Results[name], log.Log)
		}
	}

	if len(b.Errors) > 0 {
		return errors.New("there were errors while fetching the results. Please check the errors and try again")
	}

	resourceErr := b.FetchResources()
	if resourceErr != nil {
		return resourceErr
	}

	return nil
}

func (b *BasicReport) IsComparable(otherReport BasicReport) (bool, []error) {
	// check if generator configs are the same
	// are all configs present? do they have the same schedule type? do they have the same segments?
	// is call timeout the same?
	// is rate limit timeout the same?
	// would be good to be able to check if Gun and VU are the same, but idk yet how we could do that easily [hash the code?]

	if len(b.GeneratorConfigs) != len(otherReport.GeneratorConfigs) {
		return false, []error{fmt.Errorf("generator configs count is different. Expected %d, got %d", len(b.GeneratorConfigs), len(otherReport.GeneratorConfigs))}
	}

	for name1, cfg1 := range b.GeneratorConfigs {
		if cfg2, ok := otherReport.GeneratorConfigs[name1]; !ok {
			return false, []error{fmt.Errorf("generator config %s is missing from the other report", name1)}
		} else {
			if err := compareGeneratorConfigs(cfg1, cfg2); err != nil {
				return false, []error{err}
			}
		}
	}

	for name2 := range otherReport.GeneratorConfigs {
		if _, ok := b.GeneratorConfigs[name2]; !ok {
			return false, []error{fmt.Errorf("generator config %s is missing from the current report", name2)}
		}
	}

	if b.ExecutionEnvironment != otherReport.ExecutionEnvironment {
		return false, []error{fmt.Errorf("execution environments are different. Expected %s, got %s", b.ExecutionEnvironment, otherReport.ExecutionEnvironment)}
	}

	// check if pods resources are the same
	// are all pods present? do they have the same resources?
	if resourceErr := b.CompareResources(&otherReport.ResourceReporter); resourceErr != nil {
		return false, []error{resourceErr}
	}

	// check if queries are the same
	// are all queries present? do they have the same template?
	lokiQueriesErr := compareLokiQueries(b.LokiQueries, otherReport.LokiQueries)
	if lokiQueriesErr != nil {
		return false, []error{lokiQueriesErr}
	}

	return true, nil
}

func compareLokiQueries(this, other map[string]string) error {
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

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("query %s is missing from the current report", name2)
		}
	}

	return nil
}

func compareGeneratorConfigs(cfg1, cfg2 *wasp.Config) error {
	if cfg1.LoadType != cfg2.LoadType {
		return fmt.Errorf("load types are different. Expected %s, got %s", cfg1.LoadType, cfg2.LoadType)
	}

	if len(cfg1.Schedule) != len(cfg2.Schedule) {
		return fmt.Errorf("schedules are different. Expected %d, got %d", len(cfg1.Schedule), len(cfg2.Schedule))
	}

	for i, segment1 := range cfg1.Schedule {
		segment2 := cfg2.Schedule[i]
		if segment1 == nil {
			return fmt.Errorf("schedule at index %d is nil in the current report", i)
		}
		if segment2 == nil {
			return fmt.Errorf("schedule at index %d is nil in the other report", i)
		}
		if *segment1 != *segment2 {
			return fmt.Errorf("schedules at index %d are different. Expected %s, got %s", i, mustMarshallSegment(segment1), mustMarshallSegment(segment2))
		}
	}

	if cfg1.CallTimeout != cfg2.CallTimeout {
		return fmt.Errorf("call timeouts are different. Expected %s, got %s", cfg1.CallTimeout, cfg2.CallTimeout)
	}

	if cfg1.RateLimitUnitDuration != cfg2.RateLimitUnitDuration {
		return fmt.Errorf("rate limit unit durations are different. Expected %s, got %s", cfg1.RateLimitUnitDuration, cfg2.RateLimitUnitDuration)
	}

	return nil
}

func mustMarshallSegment(segment *wasp.Segment) string {
	segmentBytes, err := json.MarshalIndent(segment, "", " ")
	if err != nil {
		panic(err)
	}

	return string(segmentBytes)
}
