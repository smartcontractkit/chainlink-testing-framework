package benchspy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestBenchSpy_NewLokiQueryExecutor(t *testing.T) {
	queries := map[string]string{
		"query1": "test query 1",
		"query2": "test query 2",
	}
	config := &wasp.LokiConfig{
		URL:       "http://localhost:3100",
		TenantID:  "test",
		BasicAuth: "user:pass",
	}

	executor := NewLokiQueryExecutor("some_generator", queries, config)
	assert.Equal(t, "loki", executor.KindName)
	assert.Equal(t, queries, executor.Queries)
	assert.Equal(t, config, executor.Config)
	assert.NotNil(t, executor.QueryResults)
}

func TestBenchSpy_LokiQueryExecutor_Results(t *testing.T) {
	executor := &LokiQueryExecutor{
		QueryResults: map[string]interface{}{
			"test": []string{"result1", "result2"},
		},
	}
	results := executor.Results()
	assert.Equal(t, executor.QueryResults, results)
}

type anotherQueryExecutor struct{}

func (a *anotherQueryExecutor) Kind() string {
	return "another"
}

func (a *anotherQueryExecutor) Validate() error {
	return nil
}
func (a *anotherQueryExecutor) Execute(_ context.Context) error {
	return nil
}
func (a *anotherQueryExecutor) Results() map[string]interface{} {
	return nil
}
func (a *anotherQueryExecutor) IsComparable(_ QueryExecutor) error {
	return nil
}
func (a *anotherQueryExecutor) TimeRange(_, _ time.Time) {

}

func TestBenchSpy_LokiQueryExecutor_IsComparable(t *testing.T) {
	executor1 := &LokiQueryExecutor{
		GeneratorNameString: "generator",
		Queries:             map[string]string{"q1": "query1"},
	}
	executor2 := &LokiQueryExecutor{
		GeneratorNameString: "generator",
		Queries:             map[string]string{"q1": "query2"},
	}
	executor3 := &LokiQueryExecutor{
		GeneratorNameString: "generator",
		Queries:             map[string]string{"q2": "query1"},
	}
	executor4 := &LokiQueryExecutor{
		GeneratorNameString: "generator",
		Queries:             map[string]string{"q1": "query1", "q2": "query2"},
	}
	executor5 := &LokiQueryExecutor{
		GeneratorNameString: "generator",
		Queries:             map[string]string{"q1": "query1", "q3": "query3"},
	}
	executor6 := &LokiQueryExecutor{
		GeneratorNameString: "other",
		Queries:             map[string]string{"q1": "query1"},
	}

	t.Run("same queries", func(t *testing.T) {
		err := executor1.IsComparable(executor1)
		assert.NoError(t, err)
	})

	t.Run("different generator names", func(t *testing.T) {
		err := executor1.IsComparable(executor6)
		assert.Error(t, err)
	})

	t.Run("same queries, different names", func(t *testing.T) {
		err := executor1.IsComparable(executor3)
		assert.Error(t, err)
	})

	t.Run("same names, different queries", func(t *testing.T) {
		err := executor1.IsComparable(executor2)
		assert.Error(t, err)
	})

	t.Run("different types", func(t *testing.T) {
		invalidExecutor := &anotherQueryExecutor{}

		err := executor1.IsComparable(invalidExecutor)
		assert.Error(t, err)
	})

	t.Run("different query count", func(t *testing.T) {
		err := executor1.IsComparable(executor4)
		assert.Error(t, err)
	})

	t.Run("missing query", func(t *testing.T) {
		err := executor4.IsComparable(executor5)
		assert.Error(t, err)
	})
}

func TestBenchSpy_LokiQueryExecutor_Validate(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		executor := &LokiQueryExecutor{
			GeneratorNameString: "generator",
			Queries:             map[string]string{"q1": "query1"},
			Config:              &wasp.LokiConfig{},
		}
		err := executor.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing generator name", func(t *testing.T) {
		executor := &LokiQueryExecutor{
			Config: &wasp.LokiConfig{},
		}
		err := executor.Validate()
		assert.Error(t, err)
	})

	t.Run("missing queries", func(t *testing.T) {
		executor := &LokiQueryExecutor{
			Config: &wasp.LokiConfig{},
		}
		err := executor.Validate()
		assert.Error(t, err)
	})

	t.Run("missing config", func(t *testing.T) {
		executor := &LokiQueryExecutor{
			Queries: map[string]string{"q1": "query1"},
		}
		err := executor.Validate()
		assert.Error(t, err)
	})
}

func TestBenchSpy_LokiQueryExecutor_Execute(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/loki/api/v1/query_range", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"data": {
				"result": [
					{
						"stream": {"namespace": "test"},
						"values": [["1234567890", "Log message 1"]]
					}
				]
			}
		}`))
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	executor := &LokiQueryExecutor{
		Queries: map[string]string{"test_query": "test"},
		Config: &wasp.LokiConfig{
			URL:       mockServer.URL,
			TenantID:  "test",
			BasicAuth: "user:pass",
		},
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
	}

	err := executor.Execute(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, executor.QueryResults, "test_query")

	asStringSlice, ok := executor.QueryResults["test_query"].([]string)
	assert.True(t, ok)

	assert.Equal(t, "Log message 1", asStringSlice[0])
}

func TestBenchSpy_LokiQueryExecutor_TimeRange(t *testing.T) {
	executor := &LokiQueryExecutor{}
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()

	executor.TimeRange(start, end)
	assert.Equal(t, start, executor.StartTime)
	assert.Equal(t, end, executor.EndTime)
}

func TestBenchSpy_NewStandardMetricsLokiExecutor(t *testing.T) {
	config := &wasp.LokiConfig{
		URL:       "http://localhost:3100",
		TenantID:  "test",
		BasicAuth: "user:pass",
	}
	testName := "test"
	genName := "generator"
	branch := "main"
	commit := "abc123"
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()

	executor, err := NewStandardMetricsLokiExecutor(config, testName, genName, branch, commit, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, executor)
	assert.Equal(t, "loki", executor.KindName)
	assert.Len(t, executor.Queries, len(StandardLoadMetrics))
}

func TestBenchSpy_CalculateTimeRange(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"exact hours", 2 * time.Hour, "2h"},
		{"exact minutes", 30 * time.Minute, "30m"},
		{"seconds", 45 * time.Second, "45s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			end := start.Add(tt.duration)
			got := calculateTimeRange(start, end)
			assert.Equal(t, tt.want, got)
		})
	}
}
