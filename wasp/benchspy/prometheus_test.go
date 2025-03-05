package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBenchSpy_NewPrometheusQueryExecutor(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		queries := map[string]string{"test": "query"}

		executor, err := NewPrometheusQueryExecutor(queries, &PrometheusConfig{Url: "http://localhost:9090"})

		require.NoError(t, err)
		assert.NotNil(t, executor)
		assert.Equal(t, queries, executor.Queries)
	})
}

func TestBenchSpy_NewStandardPrometheusQueryExecutor(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()

		executor, err := NewStandardPrometheusQueryExecutor(startTime, endTime, &PrometheusConfig{Url: "http://localhost:9090", NameRegexPatterns: []string{"test.*"}})

		require.NoError(t, err)
		assert.NotNil(t, executor)
		assert.NotEmpty(t, executor.Queries)
		assert.Equal(t, 6, len(executor.Queries))

		queries := []string{}
		for name := range executor.Queries {
			queries = append(queries, name)
		}

		assert.Contains(t, queries, string(MedianCPUUsage))
		assert.Contains(t, queries, string(MedianMemUsage))
		assert.Contains(t, queries, string(P95MemUsage))
		assert.Contains(t, queries, string(P95CPUUsage))
	})
}

func TestBenchSpy_PrometheusQueryExecutor_Execute(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		// Create test data
		expectedQuery := "rate(test_metric[5m])"
		expectedTime := time.Now()
		expectedValue := &model.Vector{
			&model.Sample{
				Timestamp: model.Time(expectedTime.Unix()),
				Value:     42,
			},
		}

		// Create mock client with custom query behavior
		mockClient := &mockPrometheusClient{
			queryFn: func(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error) {
				assert.Equal(t, expectedQuery, query)
				return expectedValue, nil, nil
			},
		}

		executor := &PrometheusQueryExecutor{
			client:       mockClient,
			Queries:      map[string]string{"test_metric": expectedQuery},
			warnings:     make(map[string]v1.Warnings),
			QueryResults: make(map[string]interface{}),
			StartTime:    expectedTime.Add(-1 * time.Hour),
			EndTime:      expectedTime,
		}

		err := executor.Execute(context.Background())

		require.NoError(t, err)
		assert.NotEmpty(t, executor.QueryResults)
		assert.Contains(t, executor.QueryResults, "test_metric")

		// Verify the stored result
		result, ok := executor.QueryResults["test_metric"]
		require.True(t, ok)
		assert.Equal(t, expectedValue, result)
	})

	t.Run("handles warnings", func(t *testing.T) {
		// Create test data
		expectedQuery := "rate(test_metric[5m])"
		expectedTime := time.Now()
		expectedWarnings := v1.Warnings{"warning1", "warning2"}
		expectedValue := &model.Vector{
			&model.Sample{
				Timestamp: model.Time(expectedTime.Unix()),
				Value:     42,
			},
		}

		// Create mock client that returns warnings
		mockClient := &mockPrometheusClient{
			queryFn: func(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error) {
				return expectedValue, expectedWarnings, nil
			},
		}

		executor := &PrometheusQueryExecutor{
			client:       mockClient,
			Queries:      map[string]string{"test_metric": expectedQuery},
			warnings:     make(map[string]v1.Warnings),
			QueryResults: make(map[string]interface{}),
			StartTime:    expectedTime.Add(-1 * time.Hour),
			EndTime:      expectedTime,
		}

		err := executor.Execute(context.Background())

		// Verify execution
		require.NoError(t, err)
		assert.NotEmpty(t, executor.QueryResults)
		assert.Contains(t, executor.QueryResults, "test_metric")

		// Verify warnings were stored
		assert.Contains(t, executor.Warnings(), "test_metric")
		assert.Equal(t, expectedWarnings, executor.Warnings()["test_metric"])

		// Verify result was still stored despite warnings
		result, ok := executor.QueryResults["test_metric"]
		require.True(t, ok)
		assert.Equal(t, expectedValue, result)
	})

	t.Run("handles query error", func(t *testing.T) {
		mockClient := &mockPrometheusClient{
			queryFn: func(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error) {
				return nil, nil, fmt.Errorf("query failed")
			},
		}

		executor := &PrometheusQueryExecutor{
			client:       mockClient,
			Queries:      map[string]string{"test_metric": "invalid_query"},
			warnings:     make(map[string]v1.Warnings),
			QueryResults: make(map[string]interface{}),
			StartTime:    time.Now().Add(-1 * time.Hour),
			EndTime:      time.Now(),
		}

		err := executor.Execute(context.Background())
		require.Error(t, err)
		assert.Empty(t, executor.QueryResults)
	})
}

func TestBenchSpy_PrometheusQueryExecutor_Validate(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{
			client:    &mockPrometheusClient{},
			Queries:   map[string]string{"test": "query"},
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}

		err := executor.Validate()

		require.NoError(t, err)
	})

	t.Run("missing client", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{
			Queries:   map[string]string{"test": "query"},
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}

		err := executor.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "client is nil")
	})

	t.Run("empty queries", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{
			client:    &mockPrometheusClient{},
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}

		err := executor.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no queries provided")
	})
}

func TestBenchSpy_PrometheusQueryExecutor_IsComparable(t *testing.T) {
	t.Run("same queries", func(t *testing.T) {
		queries := map[string]string{"test": "query"}
		executor1 := &PrometheusQueryExecutor{Queries: queries}
		executor2 := &PrometheusQueryExecutor{Queries: queries}

		err := executor1.IsComparable(executor2)

		require.NoError(t, err)
	})

	t.Run("different query count", func(t *testing.T) {
		executor1 := &PrometheusQueryExecutor{
			Queries: map[string]string{
				"test1": "query1",
				"test2": "query2",
			},
		}
		executor2 := &PrometheusQueryExecutor{
			Queries: map[string]string{
				"test1": "query1",
			},
		}

		err := executor1.IsComparable(executor2)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "queries count is different")
	})

	t.Run("missing query", func(t *testing.T) {
		executor1 := &PrometheusQueryExecutor{
			Queries: map[string]string{
				"test1": "query1",
				"test2": "query2",
			},
		}
		executor2 := &PrometheusQueryExecutor{
			Queries: map[string]string{
				"test1": "query1",
				"test3": "query3",
			},
		}

		err := executor1.IsComparable(executor2)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "query test2 is missing")
	})

	t.Run("different query content", func(t *testing.T) {
		executor1 := &PrometheusQueryExecutor{
			Queries: map[string]string{
				"test1": "query1",
				"test2": "query2",
			},
		}
		executor2 := &PrometheusQueryExecutor{
			Queries: map[string]string{
				"test1": "query1",
				"test2": "different_query",
			},
		}

		err := executor1.IsComparable(executor2)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "query test2 is different")
	})

	t.Run("different types", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{}
		other := &mockQueryExecutor{
			kindName: "mock",
			queries:  make(map[string]string),
			results:  make(map[string]interface{}),
		}

		err := executor.IsComparable(other)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected type")
	})
}

func TestBenchSpy_PrometheusQueryExecutor_JSONMarshalling(t *testing.T) {
	t.Run("marshal and unmarshal", func(t *testing.T) {
		original := &PrometheusQueryExecutor{
			KindName: "prometheus",
			Queries: map[string]string{
				"test_metric":  "rate(test[5m])",
				"test_metric2": "histogram_quantile(0.95, test[5m])",
			},
			QueryResults: map[string]interface{}{
				"test_metric": &model.Vector{
					&model.Sample{
						Timestamp: model.Time(time.Now().Unix()),
						Value:     42,
					},
				},
			},
		}

		// Marshal to JSON
		data, err := json.Marshal(original)
		require.NoError(t, err)

		// Verify JSON contains expected fields
		jsonStr := string(data)
		assert.Contains(t, jsonStr, "prometheus")
		assert.Contains(t, jsonStr, "test_metric")
		assert.Contains(t, jsonStr, "rate(test[5m])")

		// Unmarshal back
		var decoded PrometheusQueryExecutor
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		// Verify fields
		assert.Equal(t, original.KindName, decoded.KindName)
		assert.Equal(t, original.Queries, decoded.Queries)
		assert.Equal(t, len(original.QueryResults), len(decoded.QueryResults))

		// Verify unexported fields are not marshalled
		assert.Zero(t, decoded.StartTime)
		assert.Zero(t, decoded.EndTime)
		assert.Nil(t, decoded.client)
	})
}

type mockQueryExecutor struct {
	kindName string
	queries  map[string]string
	results  map[string]interface{}
}

func (m *mockQueryExecutor) Execute(ctx context.Context) error {
	return nil
}

func (m *mockQueryExecutor) IsComparable(other QueryExecutor) error {
	return nil
}

func (m *mockQueryExecutor) Validate() error {
	return nil
}

func (m *mockQueryExecutor) Kind() string {
	return m.kindName
}

func (m *mockQueryExecutor) TimeRange(startTime, endTime time.Time) {
}

func (m *mockQueryExecutor) Results() map[string]interface{} {
	return m.results
}

// Mock Prometheus client for testing
type mockPrometheusClient struct {
	// Add test control fields
	queryFn func(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error)
}

// Default query implementation
func (m *mockPrometheusClient) Query(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error) {
	if m.queryFn != nil {
		return m.queryFn(ctx, query, ts, opts...)
	}
	return &model.Vector{}, v1.Warnings{}, nil
}

// Minimal implementations for remaining interface methods
func (m *mockPrometheusClient) Alerts(ctx context.Context) (v1.AlertsResult, error) {
	return v1.AlertsResult{}, nil
}

func (m *mockPrometheusClient) AlertManagers(ctx context.Context) (v1.AlertManagersResult, error) {
	return v1.AlertManagersResult{}, nil
}

func (m *mockPrometheusClient) CleanTombstones(ctx context.Context) error {
	return nil
}

func (m *mockPrometheusClient) Config(ctx context.Context) (v1.ConfigResult, error) {
	return v1.ConfigResult{}, nil
}

func (m *mockPrometheusClient) DeleteSeries(ctx context.Context, matches []string, startTime, endTime time.Time) error {
	return nil
}

func (m *mockPrometheusClient) Flags(ctx context.Context) (v1.FlagsResult, error) {
	return v1.FlagsResult{}, nil
}

func (m *mockPrometheusClient) LabelNames(ctx context.Context, matches []string, startTime, endTime time.Time, opts ...v1.Option) ([]string, v1.Warnings, error) {
	return []string{}, v1.Warnings{}, nil
}

func (m *mockPrometheusClient) LabelValues(ctx context.Context, label string, matches []string, startTime, endTime time.Time, opts ...v1.Option) (model.LabelValues, v1.Warnings, error) {
	return model.LabelValues{}, v1.Warnings{}, nil
}

func (m *mockPrometheusClient) QueryRange(ctx context.Context, query string, r v1.Range, opts ...v1.Option) (model.Value, v1.Warnings, error) {
	return &model.Vector{}, v1.Warnings{}, nil
}

func (m *mockPrometheusClient) QueryExemplars(ctx context.Context, query string, startTime, endTime time.Time) ([]v1.ExemplarQueryResult, error) {
	return []v1.ExemplarQueryResult{}, nil
}

func (m *mockPrometheusClient) Buildinfo(ctx context.Context) (v1.BuildinfoResult, error) {
	return v1.BuildinfoResult{}, nil
}

func (m *mockPrometheusClient) Runtimeinfo(ctx context.Context) (v1.RuntimeinfoResult, error) {
	return v1.RuntimeinfoResult{}, nil
}

func (m *mockPrometheusClient) Series(ctx context.Context, matches []string, startTime, endTime time.Time, opts ...v1.Option) ([]model.LabelSet, v1.Warnings, error) {
	return []model.LabelSet{}, v1.Warnings{}, nil
}

func (m *mockPrometheusClient) Snapshot(ctx context.Context, skipHead bool) (v1.SnapshotResult, error) {
	return v1.SnapshotResult{}, nil
}

func (m *mockPrometheusClient) Rules(ctx context.Context) (v1.RulesResult, error) {
	return v1.RulesResult{}, nil
}

func (m *mockPrometheusClient) Targets(ctx context.Context) (v1.TargetsResult, error) {
	return v1.TargetsResult{}, nil
}

func (m *mockPrometheusClient) TargetsMetadata(ctx context.Context, matchTarget, metric, limit string) ([]v1.MetricMetadata, error) {
	return []v1.MetricMetadata{}, nil
}

func (m *mockPrometheusClient) Metadata(ctx context.Context, metric, limit string) (map[string][]v1.Metadata, error) {
	return map[string][]v1.Metadata{}, nil
}

func (m *mockPrometheusClient) TSDB(ctx context.Context, opts ...v1.Option) (v1.TSDBResult, error) {
	return v1.TSDBResult{}, nil
}

func (m *mockPrometheusClient) WalReplay(ctx context.Context) (v1.WalReplayStatus, error) {
	return v1.WalReplayStatus{}, nil
}
