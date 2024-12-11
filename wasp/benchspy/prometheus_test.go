package benchspy

import (
	"context"
	"fmt"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrometheusQueryExecutor(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		queries := map[string]string{"test": "query"}

		executor, err := NewPrometheusQueryExecutor("http://localhost:9090", startTime, endTime, queries)

		require.NoError(t, err)
		assert.NotNil(t, executor)
		assert.Equal(t, queries, executor.Queries)
		assert.Equal(t, startTime, executor.startTime)
		assert.Equal(t, endTime, executor.endTime)
	})
}

func TestNewStandardPrometheusQueryExecutor(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()

		executor, err := NewStandardPrometheusQueryExecutor("http://localhost:9090", startTime, endTime, "test.*")

		require.NoError(t, err)
		assert.NotNil(t, executor)
		assert.NotEmpty(t, executor.Queries)
		assert.Equal(t, 4, len(executor.Queries))

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

func TestPrometheusQueryExecutor_Execute(t *testing.T) {
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
			startTime:    expectedTime.Add(-1 * time.Hour),
			endTime:      expectedTime,
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
			startTime:    time.Now().Add(-1 * time.Hour),
			endTime:      time.Now(),
		}

		err := executor.Execute(context.Background())
		require.Error(t, err)
		assert.Empty(t, executor.QueryResults)
	})
}

func TestPrometheusQueryExecutor_Validate(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{
			client:    &mockPrometheusClient{},
			Queries:   map[string]string{"test": "query"},
			startTime: time.Now(),
			endTime:   time.Now(),
		}

		err := executor.Validate()

		require.NoError(t, err)
	})

	t.Run("missing client", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{
			Queries:   map[string]string{"test": "query"},
			startTime: time.Now(),
			endTime:   time.Now(),
		}

		err := executor.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "client is nil")
	})

	t.Run("empty queries", func(t *testing.T) {
		executor := &PrometheusQueryExecutor{
			client:    &mockPrometheusClient{},
			startTime: time.Now(),
			endTime:   time.Now(),
		}

		err := executor.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no queries provided")
	})
}

func TestPrometheusQueryExecutor_IsComparable(t *testing.T) {
	t.Run("same queries", func(t *testing.T) {
		queries := map[string]string{"test": "query"}
		executor1 := &PrometheusQueryExecutor{Queries: queries}
		executor2 := &PrometheusQueryExecutor{Queries: queries}

		err := executor1.IsComparable(executor2)

		require.NoError(t, err)
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

type mockQueryExecutor struct {
	kindName string
	queries  map[string]string
	results  map[string]interface{}
}

func (m *mockQueryExecutor) Execute(ctx context.Context) error {
	return nil
}

func (r *mockQueryExecutor) IsComparable(other QueryExecutor) error {
	return nil
}

func (m *mockQueryExecutor) Validate() error {
	return nil
}

func (m *mockQueryExecutor) Kind() string {
	return m.kindName
}

func (r *mockQueryExecutor) TimeRange(startTime, endTime time.Time) {
	return
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
