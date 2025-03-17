package benchspy

import (
	"context"
	"encoding/json"
	goerrors "errors"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

// StandardReport is a report that contains all the necessary data for a performance test
type StandardReport struct {
	BasicData
	LocalStorage
	QueryExecutors []QueryExecutor `json:"query_executors"`
}

// Store saves the report to local storage as a JSON file.
// It returns the absolute path of the stored file and any error encountered.
func (sr *StandardReport) Store() (string, error) {
	return sr.LocalStorage.Store(sr.TestName, sr.CommitOrTag, sr)
}

// Load retrieves a report based on the specified test name and commit or tag.
// It utilizes local storage to find and decode the corresponding report file,
// ensuring that the report is available for further processing or analysis.
func (sr *StandardReport) Load(testName, commitOrTag string) error {
	return sr.LocalStorage.Load(testName, commitOrTag, sr)
}

// LoadLatest retrieves the most recent report for the specified test name from local storage.
// It returns an error if the report cannot be loaded, enabling users to access historical test data efficiently.
func (sr *StandardReport) LoadLatest(testName string) error {
	return sr.LocalStorage.Load(testName, "", sr)
}

// ResultsAs retrieves and casts results from a query executor to a specified type.
// It returns a map of query names to their corresponding results, or an error if casting fails.
func ResultsAs[Type any](newType Type, queryExecutor QueryExecutor, queryNames ...string) (map[string]Type, error) {
	L.Debug().
		Str("Executor kind", queryExecutor.Kind()).
		Str("Query names", strings.Join(queryNames, ", ")).
		Str("New type", fmt.Sprintf("%T", newType)).
		Msg("Casting query results to new type")

	results := make(map[string]Type)

	var toTypeOrErr = func(result interface{}, queryName string) error {
		if asType, ok := result.(Type); ok {
			results[queryName] = asType
		} else {
			return fmt.Errorf("failed to cast result to type %T. It's actual type is: %T", newType, result)
		}

		return nil
	}

	if len(queryNames) > 0 {
		for _, queryName := range queryNames {
			if result, ok := queryExecutor.Results()[queryName]; ok {
				if err := toTypeOrErr(result, queryName); err != nil {
					return nil, err
				}
			}
		}
	} else {
		for queryName, result := range queryExecutor.Results() {
			if err := toTypeOrErr(result, queryName); err != nil {
				return nil, err
			}
		}
	}

	L.Debug().
		Str("Executor kind", queryExecutor.Kind()).
		Str("Query names", strings.Join(queryNames, ", ")).
		Str("New type", fmt.Sprintf("%T", newType)).
		Msg("Successfully casted query results to new type")

	return results, nil
}

type LokiResultsByGenerator map[string]map[string][]string

// MustAllLokiResults retrieves and aggregates results from all Loki query executors in a StandardReport.
// It panics if any query execution fails, ensuring that only successful results are returned.
func MustAllLokiResults(sr *StandardReport) LokiResultsByGenerator {
	results := make(LokiResultsByGenerator)

	for _, queryExecutor := range sr.QueryExecutors {
		if strings.EqualFold(queryExecutor.Kind(), string(StandardQueryExecutor_Loki)) {
			singleResult, err := ResultsAs([]string{}, queryExecutor)
			if err != nil {
				panic(err)
			}

			asNamedGenerator := queryExecutor.(NamedGenerator)
			results[asNamedGenerator.GeneratorName()] = singleResult
		}
	}

	return results
}

type DirectResultsByGenerator map[string]map[string]float64

// MustAllDirectResults extracts and returns all direct results from a given StandardReport.
// It panics if any result extraction fails, ensuring that only valid results are processed.
func MustAllDirectResults(sr *StandardReport) DirectResultsByGenerator {
	results := make(DirectResultsByGenerator)

	for _, queryExecutor := range sr.QueryExecutors {
		if strings.EqualFold(queryExecutor.Kind(), string(StandardQueryExecutor_Direct)) {
			singleResult, err := ResultsAs(0.0, queryExecutor)
			if err != nil {
				panic(err)
			}

			asNamedGenerator := queryExecutor.(NamedGenerator)
			results[asNamedGenerator.GeneratorName()] = singleResult
		}
	}

	return results
}

// MustAllPrometheusResults retrieves all Prometheus query results from a StandardReport.
// It returns a map of query names to their corresponding model.Values, ensuring type safety.
// This function is useful for aggregating and accessing Prometheus metrics efficiently.
func MustAllPrometheusResults(sr *StandardReport) map[string]model.Value {
	results := make(map[string]model.Value)

	for _, queryExecutor := range sr.QueryExecutors {
		if strings.EqualFold(queryExecutor.Kind(), string(StandardQueryExecutor_Prometheus)) {
			for queryName, result := range queryExecutor.Results() {
				if asValue, ok := result.(model.Value); ok {
					results[queryName] = asValue
				}
			}
		}
	}

	return results
}

func calculateDiffPercentage(current, previous float64) float64 {
	if previous == 0.0 {
		if current == 0.0 {
			return 0.0
		}
		return 999.0 // Convention for infinite change when previous is 0
	}

	if current == 0.0 {
		return -100.0 // Complete improvement when current is 0
	}

	return (current - previous) / previous * 100
}

// CompareDirectWithThresholds evaluates the current and previous reports against specified thresholds.
// It checks for significant differences in metrics and returns any discrepancies found, aiding in performance analysis.
func CompareDirectWithThresholds(medianThreshold, p95Threshold, maxThreshold, errorRateThreshold float64, currentReport, previousReport *StandardReport) (bool, error) {
	if currentReport == nil || previousReport == nil {
		return true, errors.New("one or both reports are nil")
	}

	L.Info().
		Str("Current report", currentReport.CommitOrTag).
		Str("Previous report", previousReport.CommitOrTag).
		Float64("Median threshold", medianThreshold).
		Float64("P95 threshold", p95Threshold).
		Float64("Max threshold", maxThreshold).
		Float64("Error rate threshold", errorRateThreshold).
		Msg("Comparing Direct metrics with thresholds")

	if thresholdsErr := validateThresholds(medianThreshold, p95Threshold, maxThreshold, errorRateThreshold); thresholdsErr != nil {
		return true, thresholdsErr
	}

	allCurrentResults := MustAllDirectResults(currentReport)
	allPreviousResults := MustAllDirectResults(previousReport)

	var compareValues = func(
		metricName, generatorName string,
		maxDiffPercentage float64,
	) error {
		if _, ok := allCurrentResults[generatorName]; !ok {
			return fmt.Errorf("generator %s results were missing from current report", generatorName)
		}

		if _, ok := allPreviousResults[generatorName]; !ok {
			return fmt.Errorf("generator %s results were missing from previous report", generatorName)
		}

		currentForGenerator := allCurrentResults[generatorName]
		previousForGenerator := allPreviousResults[generatorName]

		if _, ok := currentForGenerator[metricName]; !ok {
			return fmt.Errorf("%s metric results were missing from current report for generator %s", metricName, generatorName)
		}

		if _, ok := previousForGenerator[metricName]; !ok {
			return fmt.Errorf("%s metric results were missing from previous report for generator %s", metricName, generatorName)
		}

		currentMetric := currentForGenerator[metricName]
		previousMetric := previousForGenerator[metricName]

		diffPercentage := calculateDiffPercentage(currentMetric, previousMetric)
		if diffPercentage > maxDiffPercentage {
			return fmt.Errorf("%s is %.4f%% different, which is higher than the threshold %.4f%%", metricName, diffPercentage, maxDiffPercentage)
		}

		return nil
	}

	errors := make(map[string][]error)

	for _, genCfg := range currentReport.GeneratorConfigs {
		if err := compareValues(string(MedianLatency), genCfg.GenName, medianThreshold); err != nil {
			errors[genCfg.GenName] = append(errors[genCfg.GenName], err)
		}

		if err := compareValues(string(Percentile95Latency), genCfg.GenName, p95Threshold); err != nil {
			errors[genCfg.GenName] = append(errors[genCfg.GenName], err)
		}

		if err := compareValues(string(MaxLatency), genCfg.GenName, maxThreshold); err != nil {
			errors[genCfg.GenName] = append(errors[genCfg.GenName], err)
		}

		if err := compareValues(string(ErrorRate), genCfg.GenName, errorRateThreshold); err != nil {
			errors[genCfg.GenName] = append(errors[genCfg.GenName], err)
		}
	}

	PrintStandardDirectMetrics(currentReport, previousReport)

	L.Info().
		Str("Current report", currentReport.CommitOrTag).
		Str("Previous report", previousReport.CommitOrTag).
		Int("Number of meaningful differences", len(errors)).
		Msg("Finished comparing Direct metrics with thresholds")

	return len(errors) > 0, concatenateGeneratorErrors(errors)
}

func concatenateGeneratorErrors(errors map[string][]error) error {
	var errs []error
	for generatorName, errors := range errors {
		for _, err := range errors {
			errs = append(errs, fmt.Errorf("[%s] %w", generatorName, err))
		}
	}
	return goerrors.Join(errs...)
}

func validateThresholds(medianThreshold, p95Threshold, maxThreshold, errorRateThreshold float64) error {
	var errs []error

	var validateThreshold = func(name string, threshold float64) error {
		if threshold < 0 || threshold > 100 {
			return fmt.Errorf("%s threshold %.4f is not in the range [0, 100]", name, threshold)
		}
		return nil
	}

	if err := validateThreshold("median", medianThreshold); err != nil {
		errs = append(errs, err)
	}

	if err := validateThreshold("p95", p95Threshold); err != nil {
		errs = append(errs, err)
	}

	if err := validateThreshold("max", maxThreshold); err != nil {
		errs = append(errs, err)
	}

	if err := validateThreshold("error rate", errorRateThreshold); err != nil {
		errs = append(errs, err)
	}

	return goerrors.Join(errs...)
}

// PrintStandardDirectMetrics outputs a comparison of direct metrics between two reports.
// It displays the current and previous values along with the percentage difference for each metric,
// helping users to quickly assess performance changes across different generator configurations.
func PrintStandardDirectMetrics(currentReport, previousReport *StandardReport) {
	currentResults := MustAllDirectResults(currentReport)
	previousResults := MustAllDirectResults(previousReport)

	for _, genCfg := range currentReport.GeneratorConfigs {
		generatorName := genCfg.GenName
		table := tablewriter.NewWriter(os.Stderr)
		table.SetHeader([]string{"Metric", previousReport.CommitOrTag, currentReport.CommitOrTag, "Diff %"})

		for _, metricName := range StandardLoadMetrics {
			metricString := string(metricName)
			diff := calculateDiffPercentage(currentResults[generatorName][metricString], previousResults[generatorName][metricString])
			table.Append([]string{metricString, fmt.Sprintf("%.4f", previousResults[genCfg.GenName][metricString]), fmt.Sprintf("%.4f", currentResults[genCfg.GenName][metricString]), fmt.Sprintf("%.4f", diff)})
		}

		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		title := "Generator: " + generatorName
		fmt.Println(title)
		fmt.Println(strings.Repeat("=", len(title)))

		table.Render()
	}
}

// FetchData retrieves data for the report within the specified time range.
// It validates the time range and executes queries in parallel, returning any errors encountered during execution.
func (sr *StandardReport) FetchData(ctx context.Context) error {
	L.Info().
		Str("Test name", sr.TestName).
		Str("Reference", sr.CommitOrTag).
		Msg("Fetching data for standard report")

	if sr.TestStart.IsZero() || sr.TestEnd.IsZero() {
		return errors.New("start and end times are not set")
	}

	errGroup, errCtx := errgroup.WithContext(ctx)
	for _, queryExecutor := range sr.QueryExecutors {
		errGroup.Go(func() error {
			// feature: PLAIN SEGMENT ONLY
			// go over all schedules and execute the code below only for ones with type "plain"
			// and then concatenate that data and return that; if parallelizing then we should first
			// create a slice of plain segments and then, when sending results over channel include the index,
			// so that we can concatenate them in the right order
			queryExecutor.TimeRange(sr.TestStart, sr.TestEnd)

			// in case someone skipped helper functions and didn't set the start and end times
			if validateErr := queryExecutor.Validate(); validateErr != nil {
				return validateErr
			}
			return queryExecutor.Execute(errCtx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	L.Info().
		Str("Test name", sr.TestName).
		Str("Reference", sr.CommitOrTag).
		Msg("Finished fetching data for standard report")

	return nil
}

// IsComparable checks if the current report can be compared with another report.
// It validates the type of the other report and ensures that their basic data and query executors are comparable.
// This function is useful for verifying report consistency before performing further analysis.
func (sr *StandardReport) IsComparable(otherReport Reporter) error {
	L.Debug().
		Str("Expected type", "*StandardReport").
		Msg("Checking if reports are comparable")

	if _, ok := otherReport.(*StandardReport); !ok {
		return fmt.Errorf("expected type %s, got %T", "*StandardReport", otherReport)
	}

	asStandardReport := otherReport.(*StandardReport)

	basicErr := sr.BasicData.IsComparable(asStandardReport.BasicData)
	if basicErr != nil {
		return basicErr
	}

	for _, queryExecutor := range sr.QueryExecutors {
		queryErr := queryExecutor.IsComparable(queryExecutor)
		if queryErr != nil {
			return queryErr
		}
	}

	L.Debug().
		Msg("Reports are comparable")

	return nil
}

type standardReportConfig struct {
	executorTypes    []StandardQueryExecutorType
	generators       []*wasp.Generator
	prometheusConfig *PrometheusConfig
	queryExecutors   []QueryExecutor
	reportDirectory  string
}

type StandardReportOption func(*standardReportConfig)

// WithStandardQueries sets the executor types for a standard report configuration.
// It allows users to specify which types of query executors to use, enabling customization
// of report generation based on their requirements.
func WithStandardQueries(executorTypes ...StandardQueryExecutorType) StandardReportOption {
	return func(c *standardReportConfig) {
		c.executorTypes = executorTypes
	}
}

// WithGenerators sets the generators for the standard report configuration.
// It allows users to specify custom generator instances to be included in the report.
func WithGenerators(generators ...*wasp.Generator) StandardReportOption {
	return func(c *standardReportConfig) {
		c.generators = generators
	}
}

// WithPrometheusConfig sets the Prometheus configuration for the standard report.
// It returns a StandardReportOption that can be used to customize report generation.
func WithPrometheusConfig(prometheusConfig *PrometheusConfig) StandardReportOption {
	return func(c *standardReportConfig) {
		c.prometheusConfig = prometheusConfig
	}
}

// WithReportDirectory sets the directory for storing report files.
// This function is useful for configuring the output location of reports
// generated by the standard reporting system.
func WithReportDirectory(reportDirectory string) StandardReportOption {
	return func(c *standardReportConfig) {
		c.reportDirectory = reportDirectory
	}
}

// WithQueryExecutors sets the query executors for a standard report configuration.
// It allows customization of how queries are executed, enhancing report generation flexibility.
func WithQueryExecutors(queryExecutors ...QueryExecutor) StandardReportOption {
	return func(c *standardReportConfig) {
		c.queryExecutors = queryExecutors
	}
}

func (c *standardReportConfig) validate() error {
	L.Debug().
		Msg("Validating standard report configuration")

	if len(c.executorTypes) == 0 && len(c.queryExecutors) == 0 {
		return errors.New("no standard executor types and no custom query executors are provided. At least one is needed")
	}

	hasPrometehus := false
	for _, t := range c.executorTypes {
		if t == StandardQueryExecutor_Prometheus {
			hasPrometehus = true
			if c.prometheusConfig == WithoutPrometheus {
				return errors.New("prometheus as query executor type is set, but prometheus config is not provided")
			}
		}
	}

	if len(c.generators) == 0 {
		return errors.New("generators are not set, at least one is required")
	}

	if c.prometheusConfig != WithoutPrometheus {
		if !hasPrometehus {
			return errors.New("prometheus config is set, but query executor type is not set to prometheus")
		}

		if c.prometheusConfig.Url == "" {
			return errors.New("prometheus url is not set")
		}
		if len(c.prometheusConfig.NameRegexPatterns) == 0 {
			return errors.New("prometheus name regex patterns are not set. At least one pattern is needed to match containers by name")
		}
	}

	L.Debug().
		Msg("Standard report configuration is valid")

	return nil
}

// NewStandardReport creates a new StandardReport based on the provided commit or tag and options.
// It initializes necessary data and query executors, ensuring all configurations are validated.
// This function is essential for generating reports that require specific data sources and execution strategies.
func NewStandardReport(commitOrTag string, opts ...StandardReportOption) (*StandardReport, error) {
	L.Info().
		Str("Reference", commitOrTag).
		Msg("Creating new standard report")

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

	basicValidateErr := basicData.Validate()
	if basicValidateErr != nil {
		return nil, basicValidateErr
	}

	queryExecutors, initErr := initStandardLoadExecutors(config, basicData)
	if initErr != nil {
		return nil, errors.Wrap(initErr, "failed to initialize standard query executors")
	}

	if len(config.queryExecutors) > 0 {
		queryExecutors = append(queryExecutors, config.queryExecutors...)
	}

	prometheusExecutors, promErr := initPrometheusQueryExecutor(config, basicData)
	if promErr != nil {
		return nil, errors.Wrap(promErr, "failed to initialize prometheus query executor")
	}

	queryExecutors = append(queryExecutors, prometheusExecutors...)

	sr := &StandardReport{
		BasicData:      *basicData,
		QueryExecutors: queryExecutors,
	}

	if config.reportDirectory != "" {
		sr.LocalStorage.Directory = config.reportDirectory
	}

	L.Info().
		Str("Reference", commitOrTag).
		Str("Test name", sr.TestName).
		Int("Number of query executors", len(sr.QueryExecutors)).
		Msg("New standard report created")

	return sr, nil
}

func initPrometheusQueryExecutor(config standardReportConfig, basicData *BasicData) ([]QueryExecutor, error) {
	var queryExecutors []QueryExecutor
	if config.prometheusConfig != WithoutPrometheus {
		// not ideal, but we want to follow the same pattern as with other executors
		for _, n := range config.prometheusConfig.NameRegexPatterns {
			prometheusExecutor, prometheusErr := NewStandardPrometheusQueryExecutor(basicData.TestStart, basicData.TestEnd, NewPrometheusConfig(config.prometheusConfig.Url, n))
			if prometheusErr != nil {
				return nil, errors.Wrapf(prometheusErr, "failed to create Prometheus executor for name patterns: %s", strings.Join(config.prometheusConfig.NameRegexPatterns, ", "))
			}
			validateErr := prometheusExecutor.Validate()
			if validateErr != nil {
				return nil, errors.Wrapf(validateErr, "failed to Prometheus executor for for name patterns: %s", strings.Join(config.prometheusConfig.NameRegexPatterns, ", "))
			}
			queryExecutors = append(queryExecutors, prometheusExecutor)
		}
	}

	return queryExecutors, nil
}

func initStandardLoadExecutors(config standardReportConfig, basicData *BasicData) ([]QueryExecutor, error) {
	var queryExecutors []QueryExecutor
	if len(config.executorTypes) != 0 {
		for _, g := range config.generators {
			for _, exType := range config.executorTypes {
				if exType != StandardQueryExecutor_Prometheus {
					executor, executorErr := initStandardQueryExecutor(exType, basicData, g)
					if executorErr != nil {
						return nil, errors.Wrapf(executorErr, "failed to create standard %s query executor for generator %s", exType, g.Cfg.GenName)
					}

					validateErr := executor.Validate()
					if validateErr != nil {
						return nil, errors.Wrapf(validateErr, "failed to validate queries for generator %s", g.Cfg.GenName)
					}

					queryExecutors = append(queryExecutors, executor)
				}
			}
		}
	}

	return queryExecutors, nil
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
	case StandardQueryExecutor_Direct:
		executor, executorErr := NewStandardDirectQueryExecutor(g)
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

// UnmarshalJSON decodes JSON data into a StandardReport struct.
// It populates the QueryExecutors and ResourceFetchers fields,
// allowing for dynamic handling of JSON structures in reports.
func (sr *StandardReport) UnmarshalJSON(data []byte) error {
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

	queryExecutors, queryErr := unmarshallQueryExecutors(raw.QueryExecutors)
	if queryErr != nil {
		return queryErr
	}

	*sr = StandardReport(raw.Alias)
	sr.QueryExecutors = queryExecutors
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

		// each new implementation of QueryExecutor might need a custom JSON unmarshaller
		// especially if it's using interface{} fields and when unmarhsalling you would like them
		// to have actual types (e.g. []string instead of []interface{}) as that will help
		// with type safety and readability
		var executor QueryExecutor
		switch typeIndicator.Kind {
		case "loki":
			executor = &LokiQueryExecutor{}
		case "direct":
			executor = &DirectQueryExecutor{}
		case "prometheus":
			executor = &PrometheusQueryExecutor{}
		default:
			return nil, fmt.Errorf("unknown query executor type: %s\nIf you added a new query executor make sure to add a custom JSON unmarshaller to StandardReport.UnmarshalJSON()", typeIndicator.Kind)
		}

		if unmarshalErr := json.Unmarshal(rawExecutor, executor); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		queryExecutors = append(queryExecutors, executor)
	}

	return queryExecutors, nil
}

func convertQueryResults(results map[string]interface{}) (map[string]interface{}, error) {
	converted := make(map[string]interface{})

	var convertToStringSlice = func(v []interface{}, key string) {
		strSlice := make([]string, len(v))
		allConverted := true
		for i, elem := range v {
			str, ok := elem.(string)
			if !ok {
				// return original slice if we can't convert, because its composed of different types
				converted[key] = v
				allConverted = false
				break
			}
			strSlice[i] = str
		}
		if allConverted {
			converted[key] = strSlice
		}
	}

	var convertToIntSlice = func(v []interface{}, key string) {
		intSlice := make([]int, len(v))
		allConverted := true
		for i, elem := range v {
			num, ok := elem.(int)
			if !ok {
				// return original slice if we can't convert, because its composed of different types
				converted[key] = v
				allConverted = false
				break
			}
			intSlice[i] = num
		}
		if allConverted {
			converted[key] = intSlice
		}
	}

	var convertToFloatSlice = func(v []interface{}, key string) {
		floatSlice := make([]float64, len(v))
		allConverted := true
		for i, elem := range v {
			f, ok := elem.(float64)
			if !ok {
				// return original slice if we can't convert, because its composed of different types
				converted[key] = v
				allConverted = false
				break
			}
			floatSlice[i] = f
		}
		if allConverted {
			converted[key] = floatSlice
		}
	}

	for key, value := range results {
		switch v := value.(type) {
		case string, int, float64:
			converted[key] = v
		case []interface{}:
			if len(v) == 0 {
				converted[key] = v
				continue
			}
			// Convert first element to determine slice type
			switch v[0].(type) {
			case string:
				convertToStringSlice(v, key)
			case int:
				convertToIntSlice(v, key)
			case float64:
				convertToFloatSlice(v, key)
			default:
				// do nothing if it's not a type we can convert
				converted[key] = v
			}
		default:
			// do nothing if it's not a type we can convert
			converted[key] = v
		}
	}
	return converted, nil
}

// FetchNewStandardReportAndLoadLatestPrevious creates a new standard report for a given commit or tag,
// loads the latest previous report, and checks their comparability.
// It returns the new report, the previous report, and any error encountered during the process.
func FetchNewStandardReportAndLoadLatestPrevious(ctx context.Context, newCommitOrTag string, newReportOpts ...StandardReportOption) (*StandardReport, *StandardReport, error) {
	newReport, err := NewStandardReport(newCommitOrTag, newReportOpts...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create new report for commit or tag %s", newCommitOrTag)
	}

	config := standardReportConfig{}
	for _, opt := range newReportOpts {
		opt(&config)
	}

	var localStorage LocalStorage

	if config.reportDirectory != "" {
		localStorage.Directory = config.reportDirectory
	}

	previousReport := &StandardReport{
		LocalStorage: localStorage,
	}

	if err = previousReport.LoadLatest(newReport.TestName); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to load latest report for test %s", newReport.TestName)
	}

	if err = newReport.FetchData(ctx); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to fetch data for new report")
	}

	err = newReport.IsComparable(previousReport)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "new report is not comparable to previous report")
	}

	return newReport, previousReport, nil
}
