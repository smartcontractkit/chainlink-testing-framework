package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBenchSpy_NewStandardReport(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}

	t.Run("successful creation", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 1, len(report.QueryExecutors))
	})

	t.Run("missing branch label", func(t *testing.T) {
		invalidGen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen",
				Labels: map[string]string{
					"commit": "abc123",
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
				},
			},
		}
		_, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, invalidGen)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing branch or commit labels")
	})

	t.Run("missing commit label", func(t *testing.T) {
		invalidGen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen",
				Labels: map[string]string{
					"branch": "abc123",
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
				},
			},
		}
		_, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, invalidGen)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing branch or commit labels")
	})

	t.Run("nil generator", func(t *testing.T) {
		_, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker)
		require.Error(t, err)
	})
}

func TestBenchSpy_StandardReport_FetchData_WithMockExecutors(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}
	ctx := context.Background()

	t.Run("successful parallel execution", func(t *testing.T) {
		// Create mock executors that simulate successful execution
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.NoError(t, err)
		assert.True(t, exec1.ValidateCalled)
		assert.True(t, exec2.ValidateCalled)
		assert.True(t, exec1.ExecuteCalled)
		assert.True(t, exec2.ExecuteCalled)
	})

	t.Run("one executor fails validation", func(t *testing.T) {
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return fmt.Errorf("validation failed") },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.True(t, exec1.ValidateCalled)
		assert.True(t, exec2.ValidateCalled)
		assert.True(t, exec1.ExecuteCalled)
		assert.False(t, exec2.ExecuteCalled)
	})

	t.Run("one executor fails execution", func(t *testing.T) {
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return fmt.Errorf("execution failed") },
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "execution failed")
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)

		exec1 := &MockQueryExecutor{
			ExecuteFn: func(ctx context.Context) error {
				cancel()
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn: func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			},
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("executor returns error", func(t *testing.T) {
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return fmt.Errorf("execution failed") },
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "execution failed")
	})
}

// MockQueryExecutor implements QueryExecutor interface for testing
type MockQueryExecutor struct {
	ExecuteFn      func(context.Context) error
	ValidateFn     func() error
	TimeRangeFn    func(time.Time, time.Time)
	ValidateCalled bool
	ExecuteCalled  bool
}

func (m *MockQueryExecutor) Execute(ctx context.Context) error {
	m.ExecuteCalled = true
	return m.ExecuteFn(ctx)
}

func (m *MockQueryExecutor) Validate() error {
	m.ValidateCalled = true
	return m.ValidateFn()
}

func (m *MockQueryExecutor) TimeRange(start, end time.Time) {
	if m.TimeRangeFn != nil {
		m.TimeRangeFn(start, end)
	}
}

func (m *MockQueryExecutor) IsComparable(other QueryExecutor) error {
	return nil
}

func (m *MockQueryExecutor) Results() map[string][]string {
	return nil
}

func TestBenchSpy_StandardReport_UnmarshalJSON(t *testing.T) {
	t.Run("valid loki executor", func(t *testing.T) {
		jsonData := `{
			"test_name": "test1",
			"commit_or_tag": "abc123",
			"query_executors": [{
				"kind": "loki",
				"query": "test query"
			}]
		}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.NoError(t, err)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &LokiQueryExecutor{}, report.QueryExecutors[0])
	})

	t.Run("unknown executor type", func(t *testing.T) {
		jsonData := `{
			"test_name": "test1",
			"commit_or_tag": "abc123",
			"query_executors": [{
				"kind": "unknown",
				"query": "test query"
			}]
		}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown query executor type")
	})
}

func TestBenchSpy_StandardReport_FetchData(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}

	t.Run("valid fetch", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)

		mockExec := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		report.QueryExecutors = []QueryExecutor{mockExec}

		err = report.FetchData(context.Background())
		require.NoError(t, err)
		assert.True(t, mockExec.ExecuteCalled)
	})

	t.Run("missing loki config", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)
		report.QueryExecutors = []QueryExecutor{&LokiQueryExecutor{
			Queries: map[string]string{"test": "query"},
		}}

		err = report.FetchData(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loki config is missing")
	})
}

func TestBenchSpy_StandardReport_IsComparable(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}

	t.Run("matching reports", func(t *testing.T) {
		report1, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)
		report2, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)

		err = report1.IsComparable(report2)
		require.NoError(t, err)
	})

	t.Run("different report types", func(t *testing.T) {
		report1, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)

		// Create a mock reporter that implements Reporter interface
		mockReport := &MockReport{}

		err = report1.IsComparable(mockReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected type *StandardReport")
	})

	t.Run("different executors", func(t *testing.T) {
		report1, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)

		// Create second report with different executor
		diffGen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen-2", // different name
				Labels: map[string]string{
					"branch": "main",
					"commit": "def456", // different commit
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(2 * time.Hour)}, // different duration
				},
			},
		}
		report2, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, diffGen)
		require.NoError(t, err)

		err = report1.IsComparable(report2)
		require.Error(t, err)
	})
}

// MockReport implements Reporter interface for testing
type MockReport struct{}

func (m *MockReport) FetchData(ctx context.Context) error      { return nil }
func (m *MockReport) Store() (string, error)                   { return "", nil }
func (m *MockReport) Load(testName, commitOrTag string) error  { return nil }
func (m *MockReport) LoadLatest(testName string) error         { return nil }
func (m *MockReport) IsComparable(other Reporter) error        { return nil }
func (m *MockReport) FetchResources(ctx context.Context) error { return nil }

func TestBenchSpy_StandardReport_Store_Load(t *testing.T) {
	// Setup test directory with git repo
	tmpDir := t.TempDir()

	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}

	t.Run("store and load", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", ExecutionEnvironment_Docker, basicGen)
		require.NoError(t, err)

		storage := LocalStorage{Directory: tmpDir}
		report.LocalStorage = storage

		_, err = report.Store()
		require.NoError(t, err)

		loadedReport := &StandardReport{LocalStorage: storage}
		err = loadedReport.Load(report.TestName, report.CommitOrTag)
		require.NoError(t, err)

		assert.Equal(t, report.TestName, loadedReport.TestName)
		assert.Equal(t, report.CommitOrTag, loadedReport.CommitOrTag)
	})
}
