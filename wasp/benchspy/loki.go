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

func NewLokiQuery(queries map[string]string, lokiConfig *wasp.LokiConfig) *LokiQuery {
	return &LokiQuery{
		Kind:         "loki",
		Queries:      queries,
		LokiConfig:   lokiConfig,
		QueryResults: make(map[string][]string),
	}
}

type LokiQuery struct {
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
	// In case something went wrong
	Errors []error `json:"errors"`

	LokiConfig *wasp.LokiConfig `json:"-"`
}

func (l *LokiQuery) Results() map[string][]string {
	return l.QueryResults
}

func (l *LokiQuery) IsComparable(otherQueryExecutor QueryExecutor) error {
	otherType := reflect.TypeOf(otherQueryExecutor)

	if otherType != reflect.TypeOf(l) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(l), otherType)
	}

	return l.compareLokiQueries(otherQueryExecutor.(*LokiQuery).Queries)
}

func (l *LokiQuery) Validate() error {
	if len(l.Queries) == 0 {
		return errors.New("there are no Loki queries, there's nothing to fetch. Please set them and try again")
	}
	if l.LokiConfig == nil {
		return errors.New("loki config is missing. Please set it and try again")
	}

	return nil
}

func (l *LokiQuery) Execute(ctx context.Context) error {
	splitAuth := strings.Split(l.LokiConfig.BasicAuth, ":")
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

			parsedLokiUrl, err := url.Parse(l.LokiConfig.URL)
			if err != nil {
				return errors.Wrapf(err, "failed to parse Loki URL %s", l.LokiConfig.URL)
			}

			lokiUrl := parsedLokiUrl.Scheme + "://" + parsedLokiUrl.Host
			lokiClient := client.NewLokiClient(lokiUrl, l.LokiConfig.TenantID, basicAuth, queryParams)

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

func (l *LokiQuery) compareLokiQueries(other map[string]string) error {
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

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("query %s is missing from the current report", name2)
		}
	}

	return nil
}

func (l *LokiQuery) TimeRange(start, end time.Time) {
	l.StartTime = start
	l.EndTime = end
}
