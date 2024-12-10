package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"golang.org/x/sync/errgroup"
)

// StandardReport is a report that contains all the necessary data for a performance test
type StandardReport struct {
	BasicData
	LocalStorage
	QueryExecutors []QueryExecutor `json:"query_executors"`
}

func (b *StandardReport) Store() (string, error) {
	return b.LocalStorage.Store(b.TestName, b.CommitOrTag, b)
}

func (b *StandardReport) Load(testName, commitOrTag string) error {
	return b.LocalStorage.Load(testName, commitOrTag, b)
}

func (b *StandardReport) LoadLatest(testName string) error {
	return b.LocalStorage.Load(testName, "", b)
}

func ResultsAs[Type any](newType Type, queryExecutors []QueryExecutor, queryExecutorType StandardQueryExecutorType, queryNames ...string) map[string]Type {
	results := make(map[string]Type)
	for _, queryExecutor := range queryExecutors {
		if strings.EqualFold(queryExecutor.Kind(), string(queryExecutorType)) {
			if len(queryNames) > 0 {
				for _, queryName := range queryNames {
					if result, ok := queryExecutor.Results()[queryName]; ok {
						results[queryName] = result.(Type)
					}
				}
			} else {
				for queryName, result := range queryExecutor.Results() {
					results[queryName] = result.(Type)
				}
			}
		}
	}

	return results
}

func (b *StandardReport) FetchData(ctx context.Context) error {
	// if b.TestStart.IsZero() || b.TestEnd.IsZero() {
	// 	fillErr := b.BasicData.FillStartEndTimes()
	// 	if fillErr != nil {
	// 		return fillErr
	// 	}
	// }

	errGroup, errCtx := errgroup.WithContext(ctx)
	for _, queryExecutor := range b.QueryExecutors {
		errGroup.Go(func() error {
			// feature: PLAIN SEGEMENT ONLY
			// go over all schedules and execute the code below only for ones with type "plain"
			// and then concatenate that data and return that; if parallelizing then we should first
			// create a slice of plain segments and then, when sending results over channel include the index,
			// so that we can concatenate them in the right order
			// queryExecutor.TimeRange(b.TestStart, b.TestEnd)

			// if validateErr := queryExecutor.ValidateQueries(); validateErr != nil {
			// 	return validateErr
			// }

			if execErr := queryExecutor.Execute(errCtx); execErr != nil {
				return execErr
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}

func (b *StandardReport) IsComparable(otherReport Reporter) error {
	if _, ok := otherReport.(*StandardReport); !ok {
		return fmt.Errorf("expected type %s, got %T", "*StandardReport", otherReport)
	}

	asStandardReport := otherReport.(*StandardReport)

	basicErr := b.BasicData.IsComparable(asStandardReport.BasicData)
	if basicErr != nil {
		return basicErr
	}

	for _, queryExecutor := range b.QueryExecutors {
		queryErr := queryExecutor.IsComparable(queryExecutor)
		if queryErr != nil {
			return queryErr
		}
	}

	return nil
}

type PrometheusConfig struct {
	url               string
	nameRegexPatterns []string
}

var WithoutPrometheus *PrometheusConfig = nil

type standardReportConfig struct {
	executorType     StandardQueryExecutorType
	generators       []*wasp.Generator
	prometheusConfig *PrometheusConfig
}

type StandardReportOption func(*standardReportConfig)

func WithStandardQueryExecutorType(executorType StandardQueryExecutorType) StandardReportOption {
	return func(c *standardReportConfig) {
		c.executorType = executorType
	}
}

func WithGenerators(generators ...*wasp.Generator) StandardReportOption {
	return func(c *standardReportConfig) {
		c.generators = generators
	}
}

func WithPrometheus(prometheusConfig *PrometheusConfig) StandardReportOption {
	return func(c *standardReportConfig) {
		c.prometheusConfig = prometheusConfig
	}
}

func (c *standardReportConfig) validate() error {
	if c.executorType == "" {
		return errors.New("executor type is not set")
	}

	if c.executorType == StandardQueryExecutor_Prometheus {
		return errors.New("prometheus as query executor is not supported currently. Use either Loki or Generator")
	}

	if len(c.generators) == 0 {
		return errors.New("generators are not set")
	}

	if c.prometheusConfig != WithoutPrometheus {
		if c.prometheusConfig.url == "" {
			return errors.New("prometheus url is not set")
		}
		if len(c.prometheusConfig.nameRegexPatterns) == 0 {
			return errors.New("prometheus name regex patterns are not set. At least one pattern is needed to match containers by name")
		}
	}

	return nil
}

// func NewStandardReport(commitOrTag string, standardQueryExecutorType StandardQueryExecutorType, generators []*wasp.Generator, prometheusConfig *PrometheusConfig) (*StandardReport, error) {
// 	basicData, basicErr := NewBasicData(commitOrTag, generators...)
// 	if basicErr != nil {
// 		return nil, errors.Wrapf(basicErr, "failed to create basic data for generators %v", generators)
// 	}

// 	startEndErr := basicData.FillStartEndTimes()
// 	if startEndErr != nil {
// 		return nil, startEndErr
// 	}

// 	basicValidateErr := basicData.Validate()
// 	if basicValidateErr != nil {
// 		return nil, basicValidateErr
// 	}

// 	var queryExecutors []QueryExecutor
// 	for _, g := range generators {
// 		executor, executorErr := initStandardQueryExecutor(standardQueryExecutorType, basicData, g)
// 		if executorErr != nil {
// 			return nil, errors.Wrapf(executorErr, "failed to create standard %s query executor for generator %s", standardQueryExecutorType, g.Cfg.GenName)
// 		}

// 		validateErr := executor.Validate()
// 		if validateErr != nil {
// 			return nil, errors.Wrapf(validateErr, "failed to validate queries for generator %s", g.Cfg.GenName)
// 		}

// 		queryExecutors = append(queryExecutors, executor)
// 	}

// 	var resourceFetchers []ResourceFetcher
// 	if prometheusConfig != nil {
// 		for _, nameRegexPattern := range prometheusConfig.nameRegexPatterns {
// 			resourceFetcher, prometheusErr := NewStandardPrometheusResourceReporter(prometheusConfig.url, basicData.TestStart, basicData.TestEnd, nameRegexPattern)
// 			if prometheusErr != nil {
// 				return nil, errors.Wrapf(prometheusErr, "failed to create Prometheus resource reporter for name %s", nameRegexPattern)
// 			}
// 			validateErr := resourceFetcher.Validate()
// 			if validateErr != nil {
// 				return nil, errors.Wrapf(validateErr, "failed to validate resources for name %s", nameRegexPattern)
// 			}
// 			resourceFetchers = append(resourceFetchers, resourceFetcher)
// 		}
// 	}

// 	return &StandardReport{
// 		BasicData:        *basicData,
// 		QueryExecutors:   queryExecutors,
// 		ResourceFetchers: resourceFetchers,
// 	}, nil
// }

func NewStandardReport(commitOrTag string, opts ...StandardReportOption) (*StandardReport, error) {
	config := standardReportConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	basicData, basicErr := NewBasicData(commitOrTag, config.generators...)
	if basicErr != nil {
		var generatorNames string
		for _, g := range config.generators {
			generatorNames += g.Cfg.GenName + ", "
		}
		return nil, errors.Wrapf(basicErr, "failed to create basic data for generators %s", generatorNames)
	}

	configErr := config.validate()
	if configErr != nil {
		return nil, configErr
	}

	startEndErr := basicData.FillStartEndTimes()
	if startEndErr != nil {
		return nil, startEndErr
	}

	basicValidateErr := basicData.Validate()
	if basicValidateErr != nil {
		return nil, basicValidateErr
	}

	var queryExecutors []QueryExecutor
	for _, g := range config.generators {
		executor, executorErr := initStandardQueryExecutor(config.executorType, basicData, g)
		if executorErr != nil {
			return nil, errors.Wrapf(executorErr, "failed to create standard %s query executor for generator %s", config.executorType, g.Cfg.GenName)
		}

		validateErr := executor.Validate()
		if validateErr != nil {
			return nil, errors.Wrapf(validateErr, "failed to validate queries for generator %s", g.Cfg.GenName)
		}

		queryExecutors = append(queryExecutors, executor)
	}

	if config.prometheusConfig != WithoutPrometheus {
		for _, nameRegexPattern := range config.prometheusConfig.nameRegexPatterns {
			resourceFetcher, prometheusErr := NewStandardPrometheusQueryExecutor(config.prometheusConfig.url, basicData.TestStart, basicData.TestEnd, nameRegexPattern)
			if prometheusErr != nil {
				return nil, errors.Wrapf(prometheusErr, "failed to create Prometheus resource reporter for name %s", nameRegexPattern)
			}
			validateErr := resourceFetcher.Validate()
			if validateErr != nil {
				return nil, errors.Wrapf(validateErr, "failed to validate resources for name %s", nameRegexPattern)
			}
			queryExecutors = append(queryExecutors, resourceFetcher)
		}
	}

	return &StandardReport{
		BasicData:      *basicData,
		QueryExecutors: queryExecutors,
	}, nil
}

func initStandardQueryExecutor(kind StandardQueryExecutorType, basicData *BasicData, g *wasp.Generator) (QueryExecutor, error) {
	switch kind {
	case StandardQueryExecutor_Loki:
		if !generatorHasLabels(g) {
			return nil, fmt.Errorf("generator %s is missing branch or commit labels", g.Cfg.GenName)
		}
		executor, executorErr := NewStandardMetricsLokiExecutor(g.Cfg.LokiConfig, basicData.TestName, g.Cfg.GenName, g.Cfg.Labels["branch"], g.Cfg.Labels["commit"], basicData.TestStart, basicData.TestEnd)
		if executorErr != nil {
			return nil, errors.Wrapf(executorErr, "failed to create standard Loki query executor for generator %s", g.Cfg.GenName)
		}
		return executor, nil
	case StandardQueryExecutor_Generator:
		executor, executorErr := NewGeneratorQueryExecutor(g)
		if executorErr != nil {
			return nil, errors.Wrapf(executorErr, "failed to create standard generator query executor for generator %s", g.Cfg.GenName)
		}
		return executor, nil
	default:
		return nil, fmt.Errorf("unknown standard query executor type: %s", kind)
	}
}

func generatorHasLabels(g *wasp.Generator) bool {
	return g.Cfg.Labels["branch"] != "" && g.Cfg.Labels["commit"] != ""
}

func (s *StandardReport) UnmarshalJSON(data []byte) error {
	// helper struct with QueryExecutors as json.RawMessage
	type Alias StandardReport
	var raw struct {
		Alias
		QueryExecutors   []json.RawMessage `json:"query_executors"`
		ResourceFetchers []json.RawMessage `json:"resource_fetchers"`
	}

	// unmarshal into the helper struct to populate other fields automatically
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// var queryExecutors []QueryExecutor

	// // manually decide, which QueryExecutor implementation to use based on the "kind" field
	// for _, rawExecutor := range raw.QueryExecutors {
	// 	var typeIndicator struct {
	// 		Kind string `json:"kind"`
	// 	}
	// 	if err := json.Unmarshal(rawExecutor, &typeIndicator); err != nil {
	// 		return err
	// 	}

	// 	var executor QueryExecutor
	// 	switch typeIndicator.Kind {
	// 	case "loki":
	// 		executor = &LokiQueryExecutor{}
	// 	case "generator":
	// 		executor = &GeneratorQueryExecutor{}
	// 	default:
	// 		return fmt.Errorf("unknown query executor type: %s\nIf you added a new query executor make sure to add a custom JSON unmarshaller to StandardReport.UnmarshalJSON()", typeIndicator.Kind)
	// 	}

	// 	if err := json.Unmarshal(rawExecutor, executor); err != nil {
	// 		return err
	// 	}

	// 	queryExecutors = append(s.QueryExecutors, executor)
	// }

	queryExecutors, queryErr := unmarshallQueryExecutors(raw.QueryExecutors)
	if queryErr != nil {
		return queryErr
	}

	*s = StandardReport(raw.Alias)
	s.QueryExecutors = queryExecutors
	return nil
}

func unmarshallQueryExecutors(raw []json.RawMessage) ([]QueryExecutor, error) {
	var queryExecutors []QueryExecutor

	// manually decide, which QueryExecutor implementation to use based on the "kind" field
	for _, rawExecutor := range raw {
		var typeIndicator struct {
			Kind string `json:"kind"`
		}
		if err := json.Unmarshal(rawExecutor, &typeIndicator); err != nil {
			return nil, err
		}

		var executor QueryExecutor
		switch typeIndicator.Kind {
		case "loki":
			executor = &LokiQueryExecutor{}
		case "generator":
			executor = &GeneratorQueryExecutor{}
		case "prometheus":
			executor = &PrometheusQueryExecutor{}
		default:
			return nil, fmt.Errorf("unknown query executor type: %s\nIf you added a new query executor make sure to add a custom JSON unmarshaller to StandardReport.UnmarshalJSON()", typeIndicator.Kind)
		}

		if err := json.Unmarshal(rawExecutor, executor); err != nil {
			return nil, err
		}

		queryExecutors = append(queryExecutors, executor)
	}

	return queryExecutors, nil
}
