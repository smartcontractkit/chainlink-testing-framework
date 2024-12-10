package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type GeneratorQueryFn = func(responses *wasp.SliceBuffer[wasp.Response]) (string, error)

type GeneratorQueryExecutor struct {
	KindName     string                      `json:"kind"`
	Generator    *wasp.Generator             `json:"generator_config"`
	Queries      map[string]GeneratorQueryFn `json:"queries"`
	QueryResults map[string]interface{}      `json:"query_results"`
}

func NewGeneratorQueryExecutor(generator *wasp.Generator) (*GeneratorQueryExecutor, error) {
	g := &GeneratorQueryExecutor{
		KindName:  string(StandardQueryExecutor_Generator),
		Generator: generator,
	}

	queries, err := g.generateStandardQueries()
	if err != nil {
		return nil, err
	}

	g.Queries = queries
	g.QueryResults = make(map[string]interface{})

	return g, nil
}

func (g *GeneratorQueryExecutor) Results() map[string]interface{} {
	return g.QueryResults
}

func (l *GeneratorQueryExecutor) Kind() string {
	return l.KindName
}

func (g *GeneratorQueryExecutor) IsComparable(otherQueryExecutor QueryExecutor) error {
	otherType := reflect.TypeOf(otherQueryExecutor)

	if otherType != reflect.TypeOf(g) {
		return fmt.Errorf("expected type %s, got %s", reflect.TypeOf(g), otherType)
	}

	otherGeneratorQueryExecutor := otherQueryExecutor.(*GeneratorQueryExecutor)

	if compareGeneratorConfigs(g.Generator.Cfg, otherGeneratorQueryExecutor.Generator.Cfg) != nil {
		return errors.New("generators are not comparable")
	}

	return g.compareQueries(otherGeneratorQueryExecutor.Queries)
}

func (l *GeneratorQueryExecutor) compareQueries(other map[string]GeneratorQueryFn) error {
	this := l.Queries
	if len(this) != len(other) {
		return fmt.Errorf("queries count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1 := range this {
		if _, ok := other[name1]; !ok {
			return fmt.Errorf("query %s is missing from the other report", name1)
		}
	}

	return nil
}

func (g *GeneratorQueryExecutor) Validate() error {
	if g.Generator == nil {
		return errors.New("generator is not set")
	}

	if len(g.Queries) == 0 {
		return errors.New("at least one query is needed")
	}

	return nil
}

func (g *GeneratorQueryExecutor) Execute(_ context.Context) error {
	if g.Generator == nil {
		return errors.New("generator is not set")
	}

	for queryName, queryFunction := range g.Queries {
		if g.Generator.GetData() == nil {
			return fmt.Errorf("generator %s has no data", g.Generator.Cfg.GenName)
		}
		length := len(g.Generator.GetData().FailResponses.Data) + len(g.Generator.GetData().OKData.Data)
		allResponses := wasp.NewSliceBuffer[wasp.Response](length)

		for _, response := range g.Generator.GetData().OKResponses.Data {
			allResponses.Append(*response)
		}

		for _, response := range g.Generator.GetData().FailResponses.Data {
			allResponses.Append(*response)
		}

		if len(allResponses.Data) == 0 {
			return fmt.Errorf("no responses found for generator %s", g.Generator.Cfg.GenName)
		}

		results, queryErr := queryFunction(allResponses)
		if queryErr != nil {
			return queryErr
		}

		g.QueryResults[queryName] = results
	}

	return nil
}

func (g *GeneratorQueryExecutor) TimeRange(_, _ time.Time) {
	// nothing to do here, since all responses stored in the generator are already in the right time range
}

func (g *GeneratorQueryExecutor) generateStandardQueries() (map[string]GeneratorQueryFn, error) {
	standardQueries := make(map[string]GeneratorQueryFn)

	for _, metric := range standardLoadMetrics {
		query, err := g.standardQuery(metric)
		if err != nil {
			return nil, err
		}
		standardQueries[string(metric)] = query
	}

	return standardQueries, nil
}

func (g *GeneratorQueryExecutor) standardQuery(standardMetric StandardLoadMetric) (GeneratorQueryFn, error) {
	switch standardMetric {
	case MedianLatency:
		medianFn := func(responses *wasp.SliceBuffer[wasp.Response]) (string, error) {
			var asMiliDuration []float64
			for _, response := range responses.Data {
				asMiliDuration = append(asMiliDuration, float64(response.Duration.Milliseconds()))
			}

			return fmt.Sprintf("%.4f", CalculatePercentile(asMiliDuration, 0.5)), nil
		}
		return medianFn, nil
	case Percentile95Latency:
		p95Fn := func(responses *wasp.SliceBuffer[wasp.Response]) (string, error) {
			var asMiliDuration []float64
			for _, response := range responses.Data {
				asMiliDuration = append(asMiliDuration, float64(response.Duration.Milliseconds()))
			}

			return fmt.Sprintf("%.4f", CalculatePercentile(asMiliDuration, 0.95)), nil
		}
		return p95Fn, nil
	case ErrorRate:
		errorRateFn := func(responses *wasp.SliceBuffer[wasp.Response]) (string, error) {
			failedCount := 0.0
			successfulCount := 0.0
			for _, response := range responses.Data {
				if response.Failed {
					failedCount = failedCount + 1
				} else {
					successfulCount = successfulCount + 1
				}
			}

			ratio := failedCount / (failedCount + successfulCount)

			return fmt.Sprintf("%.4f", ratio), nil
		}
		return errorRateFn, nil
	default:
		return nil, fmt.Errorf("unsupported standard metric %s", standardMetric)
	}
}

func (g *GeneratorQueryExecutor) MarshalJSON() ([]byte, error) {
	// we need custom marshalling to only include query names, since the functions are not serializable
	type QueryExecutor struct {
		Kind         string                 `json:"kind"`
		Generator    interface{}            `json:"generator_config"`
		Queries      []string               `json:"queries"`
		QueryResults map[string]interface{} `json:"query_results"`
	}

	return json.Marshal(&QueryExecutor{
		Kind: g.KindName,
		Generator: func() interface{} {
			if g.Generator != nil {
				return g.Generator.Cfg
			}
			return nil
		}(),
		Queries: func() []string {
			keys := make([]string, 0, len(g.Queries))
			for k := range g.Queries {
				keys = append(keys, k)
			}
			return keys
		}(),
		QueryResults: g.QueryResults,
	})
}

func (g *GeneratorQueryExecutor) UnmarshalJSON(data []byte) error {
	// helper struct with QueryExecutors as json.RawMessage and QueryResults as map[string]interface{}
	// and as actual types
	type Alias GeneratorQueryExecutor
	var raw struct {
		Alias
		GeneratorCfg wasp.Config            `json:"generator_config"`
		Queries      []json.RawMessage      `json:"queries"`
		QueryResults map[string]interface{} `json:"query_results"`
	}

	// unmarshal into the helper struct to populate other fields automatically
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	queries := make(map[string]GeneratorQueryFn)

	// unmarshall only query names
	for _, rawQuery := range raw.Queries {
		var queryName string
		if err := json.Unmarshal(rawQuery, &queryName); err != nil {
			return err
		}

		queries[queryName] = nil
	}

	// convert map[string]interface{} to map[string]actualType
	convertedTypes, conversionErr := convertQueryResults(raw.QueryResults)
	if conversionErr != nil {
		return conversionErr
	}

	*g = GeneratorQueryExecutor(raw.Alias)
	g.Queries = queries
	g.QueryResults = convertedTypes
	g.Generator = &wasp.Generator{
		Cfg: &raw.GeneratorCfg,
	}
	return nil
}
