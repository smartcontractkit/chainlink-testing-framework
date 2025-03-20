package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestBenchSpy_NewDirectQueryExecutor(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				GenName: "my_gen",
				Schedule: []*wasp.Segment{
					{
						Type:     "plain",
						From:     1,
						Duration: time.Minute,
					},
				},
			},
		}
		executor, err := NewStandardDirectQueryExecutor(gen)
		assert.NoError(t, err)
		assert.Equal(t, "direct", executor.KindName)
		assert.Equal(t, gen, executor.Generator)
		assert.NotEmpty(t, executor.Queries)
		assert.NotNil(t, executor.QueryResults)

		err = executor.Validate()
		assert.NoError(t, err)
	})

	t.Run("nil generator", func(t *testing.T) {
		executor, err := NewStandardDirectQueryExecutor(nil)
		assert.NoError(t, err)

		err = executor.Validate()
		assert.Error(t, err)
	})
}

func TestBenchSpy_DirectQueryExecutor_Results(t *testing.T) {
	expected := map[string]interface{}{
		"test": "result",
	}
	executor := &DirectQueryExecutor{
		QueryResults: expected,
	}
	assert.Equal(t, expected, executor.Results())
}

func TestBenchSpy_DirectQueryExecutor_IsComparable(t *testing.T) {
	baseGen := &wasp.Generator{
		Cfg: &wasp.Config{
			GenName: "my_gen",
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: time.Minute,
				},
			},
		},
	}

	t.Run("same configs", func(t *testing.T) {
		exec1, _ := NewStandardDirectQueryExecutor(baseGen)
		exec2, _ := NewStandardDirectQueryExecutor(baseGen)
		assert.NoError(t, exec1.IsComparable(exec2))
	})

	t.Run("different query count", func(t *testing.T) {
		exec1, _ := NewStandardDirectQueryExecutor(baseGen)
		exec2, _ := NewStandardDirectQueryExecutor(baseGen)
		exec2.Queries = map[string]DirectQueryFn{"test": nil}

		assert.Error(t, exec1.IsComparable(exec2))
	})

	t.Run("different queries", func(t *testing.T) {
		exec1, _ := NewStandardDirectQueryExecutor(baseGen)
		exec2, _ := NewStandardDirectQueryExecutor(baseGen)
		exec2.Queries = map[string]DirectQueryFn{"test": nil, "test2": nil, "test3": nil}

		assert.Error(t, exec1.IsComparable(exec2))
	})

	t.Run("different configs", func(t *testing.T) {
		exec1, _ := NewStandardDirectQueryExecutor(baseGen)
		differentGen := &wasp.Generator{
			Cfg: &wasp.Config{
				GenName: "my_gen",
				Schedule: []*wasp.Segment{
					{
						Type:     "plain",
						From:     2,
						Duration: time.Minute,
					},
				},
			},
		}
		exec2, _ := NewStandardDirectQueryExecutor(differentGen)
		assert.Error(t, exec1.IsComparable(exec2))
	})

	t.Run("different types", func(t *testing.T) {
		exec1, _ := NewStandardDirectQueryExecutor(baseGen)
		exec2 := &LokiQueryExecutor{}
		assert.Error(t, exec1.IsComparable(exec2))
	})
}

func TestBenchSpy_DirectQueryExecutor_Validate(t *testing.T) {
	t.Run("valid case", func(t *testing.T) {
		executor, _ := NewStandardDirectQueryExecutor(&wasp.Generator{
			Cfg: &wasp.Config{},
		})
		assert.NoError(t, executor.Validate())
	})

	t.Run("missing generator", func(t *testing.T) {
		executor := &DirectQueryExecutor{
			Queries: map[string]DirectQueryFn{"test": nil},
		}
		assert.Error(t, executor.Validate())
	})

	t.Run("empty queries", func(t *testing.T) {
		executor := &DirectQueryExecutor{
			Generator: &wasp.Generator{},
			Queries:   map[string]DirectQueryFn{},
		}
		assert.Error(t, executor.Validate())
	})
}

type fakeGun struct {
	maxSuccesses   int
	successCounter int
	maxFailures    int
	failureCounter int
	successMutex   sync.Mutex
	failureMutex   sync.Mutex

	schedule *wasp.Segment
}

func (f *fakeGun) Call(l *wasp.Generator) *wasp.Response {
	now := time.Now()
	f.successMutex.Lock()
	defer f.successMutex.Unlock()
	if f.successCounter < f.maxSuccesses {
		f.successCounter++
		d := time.Duration(150) * time.Millisecond
		time.Sleep(d)
		finishedAt := time.Now()
		return &wasp.Response{
			StartedAt:  &now,
			Duration:   d,
			FinishedAt: &finishedAt,
			Failed:     false,
			Data:       fmt.Sprintf("success[%d]", f.successCounter),
		}
	}

	f.failureMutex.Lock()
	defer f.failureMutex.Unlock()
	if f.failureCounter < f.maxFailures {
		f.failureCounter++
		d := time.Duration(200) * time.Millisecond
		time.Sleep(d)
		finishedAt := time.Now()
		return &wasp.Response{
			StartedAt:  &now,
			Duration:   d,
			FinishedAt: &finishedAt,
			Failed:     true,
			Data:       fmt.Sprintf("failure[%d]", f.failureCounter),
		}
	}

	panic(fmt.Sprintf("fakeGun.Call called too many times (%d vs %d). Expected maxFailures: %d. Expected maxSuccesses: %d. do adjust your settings (sum of maxSuccesses and maxFilure should be greater than the duration of the segment in seconds)", f.maxFailures+f.maxSuccesses, f.maxFailures+f.maxSuccesses+1, f.maxFailures, f.maxSuccesses))
}

func TestBenchSpy_DirectQueryExecutor_Execute(t *testing.T) {
	t.Run("success case with mixed responses", func(t *testing.T) {
		cfg := &wasp.Config{
			GenName:  "my_gen",
			LoadType: wasp.RPS,
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: 5 * time.Second,
				},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 4,
			maxFailures:  3,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun

		gen, err := wasp.NewGenerator(cfg)
		require.NoError(t, err)

		gen.Run(true)

		require.Equal(t, 4, len(gen.GetData().OKResponses.Data), "expected 4 successful responses")
		require.GreaterOrEqual(t, len(gen.GetData().FailResponses.Data), 2, "expected >=2 failed responses")

		actualFailures := len(gen.GetData().FailResponses.Data)

		executor, err := NewStandardDirectQueryExecutor(gen)
		assert.NoError(t, err)

		err = executor.Execute(context.Background())
		assert.NoError(t, err)

		results := executor.Results()
		assert.NotEmpty(t, results)

		// 4 responses with ~150ms latency (150ms sleep + some execution overhead)
		// and 2-3 responses with ~200ms latency (200ms sleep + some execution overhead)
		// expected median latency: (150ms, 151ms>
		resultsAsFloats, err := ResultsAs(0.0, executor, string(MedianLatency), string(Percentile95Latency), string(ErrorRate))
		assert.NoError(t, err)
		require.Equal(t, 3, len(resultsAsFloats))
		require.InDelta(t, 151.0, resultsAsFloats[string(MedianLatency)], 1.0)

		// since we have 2-3 responses with 200-201ms latency, the 95th percentile should be (200ms, 201ms>
		require.InDelta(t, 201.0, resultsAsFloats[string(Percentile95Latency)], 1.0)

		errorRate, exists := resultsAsFloats[string(ErrorRate)]
		assert.True(t, exists)

		// error rate is the number of failures divided by the total number of responses
		expectedErrorRate := float64(actualFailures) / (float64(fakeGun.maxSuccesses) + float64(actualFailures))
		assert.Equal(t, expectedErrorRate, errorRate)
	})

	t.Run("all responses failed", func(t *testing.T) {
		cfg := &wasp.Config{
			GenName:  "my_gen",
			LoadType: wasp.RPS,
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: 5 * time.Second,
				},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 0,
			maxFailures:  10,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun

		gen, err := wasp.NewGenerator(cfg)
		assert.NoError(t, err)

		gen.Run(true)

		require.Equal(t, 0, len(gen.GetData().OKResponses.Data), "expected 0 successful responses")
		require.GreaterOrEqual(t, len(gen.GetData().FailResponses.Data), 5, "expected >=5 failed responses")

		executor, err := NewStandardDirectQueryExecutor(gen)
		assert.NoError(t, err)

		err = executor.Execute(context.Background())
		assert.NoError(t, err)

		results := executor.Results()
		assert.NotEmpty(t, results)

		errorRate, exists := results[string(ErrorRate)]
		assert.True(t, exists)
		assert.Equal(t, 1.0, errorRate)
	})

	t.Run("no responses", func(t *testing.T) {
		cfg := &wasp.Config{
			GenName:  "my_gen",
			LoadType: wasp.RPS,
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: 5 * time.Second,
				},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 0,
			maxFailures:  0,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun
		gen, err := wasp.NewGenerator(cfg)
		require.NoError(t, err)

		executor, _ := NewStandardDirectQueryExecutor(gen)
		err = executor.Execute(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no responses found for generator")
	})
}

func TestBenchSpy_DirectQueryExecutor_MarshalJSON(t *testing.T) {
	t.Run("marshal/unmarshal round trip", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				GenName: "my_gen",
				Schedule: []*wasp.Segment{
					{
						Type:     "plain",
						From:     1,
						Duration: time.Minute,
					},
				},
			},
		}
		original, _ := NewStandardDirectQueryExecutor(gen)
		original.QueryResults["test"] = 2.0
		original.QueryResults["test2"] = 12.1

		original.Queries = map[string]DirectQueryFn{
			"test": func(responses *wasp.SliceBuffer[*wasp.Response]) (float64, error) {
				return 2.0, nil
			},
			"test2": func(responses *wasp.SliceBuffer[*wasp.Response]) (float64, error) {
				return 12.1, nil
			},
		}

		data, err := json.Marshal(original)
		assert.NoError(t, err)

		var recovered DirectQueryExecutor
		err = json.Unmarshal(data, &recovered)
		assert.NoError(t, err)

		assert.Equal(t, original.KindName, recovered.KindName)
		assert.Equal(t, original.QueryResults, recovered.QueryResults)
		assert.Equal(t, len(original.Queries), len(recovered.Queries))
	})

	t.Run("marshal with nil generator", func(t *testing.T) {
		executor := &DirectQueryExecutor{}
		data, err := json.Marshal(executor)
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"generator_config":null`)
	})

	t.Run("unmarshal invalid JSON", func(t *testing.T) {
		var executor DirectQueryExecutor
		err := json.Unmarshal([]byte(`{invalid json}`), &executor)
		assert.Error(t, err)
	})
}

func TestBenchSpy_DirectQueryExecutor_TimeRange(t *testing.T) {
	executor := &DirectQueryExecutor{}
	start := time.Now()
	end := start.Add(time.Hour)

	// TimeRange should not modify the executor
	executor.TimeRange(start, end)
	assert.NotPanics(t, func() {
		executor.TimeRange(start, end)
	})
}
